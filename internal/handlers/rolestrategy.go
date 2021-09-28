package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ts.tarusov/swp/internal/dbmodels"
	"github.com/ts.tarusov/swp/pkg/models"
)

func (srv *Service) PostRoleStrategy(w http.ResponseWriter, r *http.Request) {
	prid, err := strconv.Atoi(r.URL.Query().Get("project_id"))
	if err != nil {
		srv.Error("wrong request: %s", err)
		http.Error(w, fmt.Sprintf("wrong request: %s", err), http.StatusBadRequest)
	}
	avls, err := srv.availabilityRank(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(avls) == 0 {
		onEmptyTable(srv, w, dbmodels.AvailabilityRanksTable)
	}
	bvs, err := srv.businessValue(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(bvs) == 0 {
		onEmptyTable(srv, w, dbmodels.BusinessValueTable)
	}

	roles, err := srv.roles(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(roles) == 0 {
		onEmptyTable(srv, w, dbmodels.EmplrolesTable)
	}

	team := models.NewTeam("", 0)
	for _, role := range roles {
		team.AddRole(models.NewRole(role, ""))
	}
	for _, avl := range avls {
		team.Roles[avl.Role].SetAvailabilityRank(avl.Rank)
	}
	for _, bv := range bvs {
		team.Roles[bv.Role].SetBusinessValueRank(bv.Rank)
	}

	idx := srv.seqNextVal(dbmodels.PreferredStrategyTable)
	stmt := srv.writePreferredStrategy()
	defer stmt.Close()
	for _, role := range team.Roles {
		if err := role.CalculatePreferredStrategy(); err != nil {
			srv.Error("failed to calculate strategy: %v", err)
		}
		_, err := stmt.Exec(
			idx,
			prid,
			role.Name,
			role.PreferredStrategy.Quadrant,
			role.PreferredStrategy.Strategy1,
			role.PreferredStrategy.Strategy2,
			time.Now(),
		)
		if err != nil {
			srv.Error("failed to prepare stmt: %v", err)
		}
		idx++
	}

	srv.deleteFromTableById(dbmodels.PreferredStrategyTable, prid)
	_, err = stmt.Exec()
	if err != nil {
		srv.Error("failed to bulk insert: %v", err)
	}
	w.Write([]byte("OK"))
	srv.Info("Successfully handled")
}
