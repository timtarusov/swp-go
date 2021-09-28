package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ts.tarusov/swp/internal/dbmodels"
	"github.com/ts.tarusov/swp/pkg/models"
)

func (srv *Service) PostAvailability(w http.ResponseWriter, r *http.Request) {
	prid, err := strconv.Atoi(r.URL.Query().Get("project_id"))
	if err != nil {
		srv.Error("wrong request: %s", err)
		http.Error(w, fmt.Sprintf("wrong request: %s", err), http.StatusBadRequest)
	}

	avls, err := srv.availability(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(avls) == 0 {
		onEmptyTable(srv, w, dbmodels.AvailabilityTable)
	}

	empls, err := srv.employees(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(empls) == 0 {
		onEmptyTable(srv, w, dbmodels.EmplrolesTable)
	}

	team := models.NewTeam("", empls[0].Function)

	for _, empl := range empls {
		if _, ok := team.Roles[empl.Role]; !ok {
			team.AddRole(models.NewRole(empl.Role, ""))
		}
		team.Roles[empl.Role].AddEmployee(models.NewEmployee(empl.Employee))
	}

	for _, avl := range avls {
		team.Roles[avl.Role].AddAvailability(avl.Year, decToFloat(avl.PotentialHires))
	}

	idx := srv.seqNextVal(dbmodels.AvailabilityRanksTable)
	stmt := srv.writeAvailabilityRanks()
	defer stmt.Close()

	for _, role := range team.Roles {
		role.AvailabilityRanks()
		_, err := stmt.Exec(
			idx,
			prid,
			role.Name,
			strconv.Itoa(role.Availability.Rank),
			floatToDec(role.Availability.PercentAvailable),
			time.Now(),
		)
		if err != nil {
			srv.Error("failed to prepare stmt: %v", err)
		}
		idx++
	}
	srv.deleteFromTableById(dbmodels.AvailabilityRanksTable, prid)
	_, err = stmt.Exec()
	if err != nil {
		srv.Error("failed to bulk insert: %v", err)
	}

	w.Write([]byte("OK"))
	srv.Info("Successfully handled")

}
