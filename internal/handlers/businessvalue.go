package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ts.tarusov/swp/internal/dbmodels"
	"github.com/ts.tarusov/swp/pkg/models"
)

func (srv *Service) PostBusinessValue(w http.ResponseWriter, r *http.Request) {
	prid, err := strconv.Atoi(r.URL.Query().Get("project_id"))
	if err != nil {
		srv.Error("wrong request: %s", err)
		http.Error(w, fmt.Sprintf("wrong request: %s", err), http.StatusBadRequest)
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

	team := models.NewTeam("", 0)
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
		team.AddRole(models.NewRole(impact.Role, ""))
		team.Roles[impact.Role].SetImpactForTarget(impact.Target, decToFloat(impact.Impact))
	}

	if err := team.GenerateRolesPercent(); err != nil {
		srv.Error("failed to generate roles percent: %s", err)
		http.Error(w, fmt.Sprintf("failed to generate roles percent: %s", err), http.StatusInternalServerError)
	}

	idx := srv.seqNextVal(dbmodels.BusinessValueTable)
	stmt := srv.writeBusinessValue()
	defer stmt.Close()
	for _, role := range team.Roles {
		_, err := stmt.Exec(
			idx,
			prid,
			role.Name,
			floatToDec(role.TotalGrowth),
			floatToDec(role.Percentile),
			strconv.Itoa(role.Rank),
			time.Now(),
		)
		if err != nil {
			srv.Error("failed to prepare stmt: %v", err)
		}
		idx++
	}
	srv.deleteFromTableById(dbmodels.BusinessValueTable, prid)
	_, err = stmt.Exec()
	if err != nil {
		srv.Error("failed to bulk insert: %v", err)
	}
	w.Write([]byte("OK"))
	srv.Info("Successfully handled")
}
