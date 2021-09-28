package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ts.tarusov/swp/internal/dbmodels"
	"github.com/ts.tarusov/swp/pkg/models"
)

func (srv *Service) PostSkillsDevelopment(w http.ResponseWriter, r *http.Request) {
	prid, err := strconv.Atoi(r.URL.Query().Get("project_id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("wrong request: %s", err), http.StatusBadRequest)
		srv.Error("wrong request: %s", err)
	}

	empls, err := srv.employees(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(empls) == 0 {
		onEmptyTable(srv, w, dbmodels.EmplrolesTable)

	}
	years, err := srv.years(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(years) == 0 {
		onEmptyTable(srv, w, dbmodels.TargetsTable)
	}

	rps, err := srv.roleProfile(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(rps) == 0 {
		onEmptyTable(srv, w, dbmodels.RoleProfileTable)
	}

	emplscores, err := srv.emplScores(prid)
	if err != nil {
		onReadFail(srv, err, w)
	}
	if len(emplscores) == 0 {
		onEmptyTable(srv, w, dbmodels.EmployeeScoresTable)
	}
	config, err := srv.projectConfig(prid)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't read rows: %s", err), http.StatusInternalServerError)
	}

	team := models.NewTeam("", empls[0].Function)

	for _, empl := range empls {
		if _, ok := team.Roles[empl.Role]; !ok {
			team.AddRole(models.NewRole(empl.Role, ""))
		}
		team.Roles[empl.Role].AddEmployee(models.NewEmployee(empl.Employee))
	}

	for _, rp := range rps {
		var pr bool
		if rp.IsPriority == "X" {
			pr = true
		}
		team.Roles[rp.Role].AddRoleSkill(&models.RoleSkill{
			SkillName:   rp.Skill,
			TargetLevel: rp.TargetLvl,
			IsPriority:  pr,
		})
	}

	for _, es := range emplscores {
		team.Roles[es.Role].Employees()[es.Employee].SetScore(es.Skill, es.CurrentLevel)
	}

	stmtSL := srv.writeProjectSkillLevels()
	idxSL := srv.seqNextVal(dbmodels.ProjectSkillLevelsTable)
	defer stmtSL.Close()

	stmtSf := srv.writeProjectSufficiency()
	idxSf := srv.seqNextVal(dbmodels.ProjectSufficiencyTable)
	defer stmtSf.Close()

	srv.Debug("years: %v", years)
	srv.Debug("roles: %d", len(team.Roles))

	var empllen int
	var sklen int
	for _, r := range team.Roles {
		empllen += len(r.Employees())
		sklen += len(r.Profile)
	}
	srv.Debug("employees: %d", empllen)
	srv.Debug("skills: %d", sklen/empllen)
	for _, role := range team.Roles {
		for _, empl := range role.Employees() {
			empl.ProjectSkillsDevelopment(years, 1)
			empl.ProjectEmployeeSufficiency(years)
			for _, year := range years {
				for skillID, skill := range empl.EmployeeProfile {
					var suffStr string
					if skill.Forecast[year].Sufficient {
						suffStr = "X"
					}
					_, err := stmtSL.Exec(
						idxSL,
						prid,
						empl.ID,
						skillID,
						strconv.Itoa(skill.Forecast[year].Level),
						strconv.Itoa(year),
						suffStr,
						time.Now(),
					)
					if err != nil {
						srv.Error("failed to prepare stmt: %v", err)
					}
					idxSL++
				}

				var emplSuf string
				if empl.Sufficiency[year] >= decToFloat(config.SuffThreshold) {
					emplSuf = "X"
				}
				_, err := stmtSf.Exec(
					idxSf,
					prid,
					empl.ID,
					strconv.Itoa(year),
					floatToDec(empl.Sufficiency[year]),
					emplSuf,
					time.Now(),
				)
				if err != nil {
					srv.Error("failed to prepare stmt: %v", err)
				}
				idxSf++
			}

		}
	}
	srv.deleteFromTableById(dbmodels.ProjectSkillLevelsTable, prid)
	srv.deleteFromTableById(dbmodels.ProjectSufficiencyTable, prid)
	_, err = stmtSL.Exec()
	if err != nil {
		srv.Error("failed to bulk insert stmtSL: %v", err)
	}
	_, err = stmtSf.Exec()
	if err != nil {
		srv.Error("failed to bulk insert stmtSf: %v", err)
	}
	w.Write([]byte("OK"))
	srv.Info("Successfully handled")

}
