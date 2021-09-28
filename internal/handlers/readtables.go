package handlers

import (
	"fmt"

	"github.com/ts.tarusov/swp/internal/dbmodels"
)

func (srv *Service) employees(prid int) ([]dbmodels.EmplrolesRow, error) {
	var empls []dbmodels.EmplrolesRow
	if err := srv.Db.Select(&empls, dbmodels.PrepQuery(dbmodels.EmplrolesTable, prid)); err != nil {
		return nil, err
	}
	return empls, nil
}
func (srv *Service) roleImpacts(prid int) ([]dbmodels.RoleImpactRow, error) {
	var impacts []dbmodels.RoleImpactRow
	if err := srv.Db.Select(&impacts, dbmodels.PrepQuery(dbmodels.RoleImpactTable, prid)); err != nil {
		return nil, err
	}
	return impacts, nil
}
func (srv *Service) targets(prid int) ([]dbmodels.TargetsRow, error) {
	var tgts []dbmodels.TargetsRow
	if err := srv.Db.Select(&tgts,
		fmt.Sprintf(`SELECT * FROM "%s" WHERE "PROJECT_ID" = %d ORDER BY "TARGET", "YEAR"`,
			dbmodels.TargetsTable, prid)); err != nil {
		return nil, err
	}
	return tgts, nil
}
func (srv *Service) projectConfig(prid int) (dbmodels.ProjectConfigRow, error) {
	var cfg dbmodels.ProjectConfigRow
	if err := srv.Db.Get(&cfg, dbmodels.PrepQuery(dbmodels.ProjectConfigTable, prid)); err != nil {
		return dbmodels.ProjectConfigRow{}, err
	}
	return cfg, nil
}

func (srv *Service) roleProfile(prid int) ([]dbmodels.RoleProfileRow, error) {
	var rps []dbmodels.RoleProfileRow
	if err := srv.Db.Select(&rps, dbmodels.PrepQuery(dbmodels.RoleProfileTable, prid)); err != nil {
		return nil, err
	}
	return rps, nil
}

func (srv *Service) years(prid int) ([]int, error) {
	var yrs []int
	if err := srv.Db.Select(&yrs,
		fmt.Sprintf(`SELECT DISTINCT "YEAR" FROM "%s" WHERE "PROJECT_ID" = %d ORDER BY "YEAR" ASC`,
			dbmodels.TargetsTable, prid)); err != nil {
		return nil, err
	}
	return yrs, nil
}

func (srv *Service) emplScores(prid int) ([]dbmodels.EmployeeScoresRow, error) {
	var emplscores []dbmodels.EmployeeScoresRow
	if err := srv.Db.Select(&emplscores,
		fmt.Sprintf(`SELECT 
				r."ROLE" "ROLE",
				e."EMPLOYEE" "EMPLOYEE",
				e."CURRENT_LEVEL" "CURRENT_LEVEL",
				e."SKILL" "SKILL"
				FROM "%s" e
				LEFT JOIN "%s" r
				ON e."EMPLOYEE" = r."EMPLOYEE"
				WHERE e."PROJECT_ID" = %d	`,
			dbmodels.EmployeeScoresTable, dbmodels.EmplrolesTable, prid)); err != nil {
		return nil, err
	}
	return emplscores, nil
}

func (srv *Service) availability(prid int) ([]dbmodels.AvailabilityRow, error) {
	var avls []dbmodels.AvailabilityRow
	if err := srv.Db.Select(&avls, dbmodels.PrepQuery(dbmodels.AvailabilityTable, prid)); err != nil {
		return nil, err
	}
	return avls, nil
}

func (srv *Service) availabilityRank(prid int) ([]dbmodels.AvailabilityRankRow, error) {
	var avls []dbmodels.AvailabilityRankRow
	if err := srv.Db.Select(&avls, dbmodels.PrepQuery(dbmodels.AvailabilityRanksTable, prid)); err != nil {
		return nil, err
	}
	return avls, nil
}

func (srv *Service) businessValue(prid int) ([]dbmodels.BusinessValueRow, error) {
	var bvs []dbmodels.BusinessValueRow
	if err := srv.Db.Select(&bvs, dbmodels.PrepQuery(dbmodels.BusinessValueTable, prid)); err != nil {
		return nil, err
	}
	return bvs, nil
}

func (srv *Service) roles(prid int) ([]int, error) {
	var roles []int
	if err := srv.Db.Select(&roles, fmt.Sprintf(`
	SELECT DISTINCT ROLE 
	FROM "%s"
	WHERE "PROJECT_ID" = %d
	`, dbmodels.EmplrolesTable, prid)); err != nil {
		return nil, err
	}
	return roles, nil
}

func (srv *Service) marketsalaries(fid int) ([]dbmodels.MarketSalariesRow, error) {
	var mss []dbmodels.MarketSalariesRow
	if err := srv.Db.Select(&mss, fmt.Sprintf(`SELECT * FROM "%s" WHERE "FUNCTION_ID"=%d`,
		dbmodels.MarketSalariesTable, fid)); err != nil {
		return nil, err
	}
	return mss, nil
}
