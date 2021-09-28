package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/SAP/go-hdb/driver"
	"github.com/ts.tarusov/swp/internal/dbmodels"
	"github.com/ts.tarusov/swp/pkg/models"
)

func (srv *Service) PostProjectedHeadcount(w http.ResponseWriter, r *http.Request) {
	prid, err := strconv.Atoi(r.URL.Query().Get("project_id"))
	if err != nil {
		srv.Error("wrong request: %s", err)
		http.Error(w, fmt.Sprintf("wrong request: %s", err), http.StatusBadRequest)
	}
	empls, err := srv.employees(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(empls) == 0 {
		onEmptyTable(srv, w, dbmodels.EmplrolesTable)
	}
	targets, err := srv.targets(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(targets) == 0 {
		onEmptyTable(srv, w, dbmodels.TargetsTable)
	}
	impacts, err := srv.roleImpacts(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(impacts) == 0 {
		onEmptyTable(srv, w, dbmodels.RoleImpactTable)
	}

	config, err := srv.projectConfig(prid)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't read rows: %s", err), http.StatusInternalServerError)
	}
	// if config == (dbmodels.ProjectConfigRow{}) {
	// 	http.Error(w, fmt.Sprintf("table %s is empty", dbmodels.ProjectConfigTable), http.StatusBadRequest)
	// }

	team := models.NewTeam("", empls[0].Function)

	for _, empl := range empls {
		if _, ok := team.Roles[empl.Role]; !ok {
			team.AddRole(models.NewRole(empl.Role, ""))
		}
		team.Roles[empl.Role].AddEmployee(models.NewEmployee(empl.Employee))
	}

	for _, tgt := range targets {
		team.AddTarget(tgt.Target, tgt.Type_T, tgt.Year, decToFloat(tgt.Value))
	}
	if err := team.RecalculateStrategy(); err != nil {
		srv.Error("failed to recalculate strategy: %s", err)
		http.Error(w, fmt.Sprintf("failed to recalculate strategy: %s", err), http.StatusInternalServerError)
	}
	if err := team.GenerateYoYTargets(); err != nil {
		srv.Error("failed to generate yoy targets: %s", err)
		http.Error(w, fmt.Sprintf("failed to generate yoy targets: %s", err), http.StatusInternalServerError)
	}

	for _, impact := range impacts {
		team.Roles[impact.Role].SetImpactForTarget(impact.Target, decToFloat(impact.Impact))
	}

	idx := srv.seqNextVal(dbmodels.ProjectedHeadcountTable)
	stmt := srv.writeProjectHeadcount()
	defer stmt.Close()
	for _, role := range team.Roles {
		role.ProjectHeadcount(decToFloat(config.PersLeverage))
		for _, year := range team.Years() {
			_, err := stmt.Exec(
				idx,
				team.Function,
				role.Name,
				"1",
				strconv.Itoa(year),
				prid,
				floatToDec(role.ProjectedHeadcount.Headcount[year]),
				time.Now(),
			)
			if err != nil {
				srv.Error("failed to prepare stmt: %v", err)
			}
			idx++
			_, err = stmt.Exec(
				idx,
				team.Function,
				role.Name,
				"2",
				strconv.Itoa(year),
				prid,
				floatToDec(role.ProjectedHeadcount.Growth[year]),
				time.Now(),
			)
			if err != nil {
				srv.Error("failed to prepare stmt: %v", err)
			}
			idx++
		}
	}
	srv.deleteFromTableById(dbmodels.ProjectedHeadcountTable, prid)
	_, err = stmt.Exec()
	if err != nil {
		srv.Error("failed to bulk insert: %v", err)
	}

	// resp, err := json.Marshal(out)
	// if err != nil {
	// 	http.Error(w, "couldn't marshall rows", http.StatusInternalServerError)
	// 	log.Println("couldn't marshall rows")
	// }
	w.Write([]byte("OK"))
	srv.Info("Successfully handled")
}
