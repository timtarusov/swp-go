package dbmodels

import (
	"fmt"
	"os"

	_ "github.com/SAP/go-hdb/driver"
)

const (
	DriverName                 = "hdb"
	Host                       = "hanabwprod-01.severstal.severstalgroup.com:30015"
	EmplrolesTable             = "YBS_HC.SWP.TABLE::SWP_EMPROLE"
	RoleImpactTable            = "YBS_HC.SWP.TABLE::SWP_ROLE_IMPACT"
	TargetsTable               = "YBS_HC.SWP.TABLE::SWP_TARGETS"
	ProjectConfigTable         = "YBS_HC.SWP.TABLE::SWP_PROJECT_CONFIG"
	ProjectedHeadcountTable    = "YBS_HC.SWP.TABLE::SWP_PROJECTED_HEADCOUNT"
	RoleProfileTable           = "YBS_HC.SWP.TABLE::SWP_ROLE_PROFILE"
	EmployeeScoresTable        = "YBS_HC.SWP.TABLE::SWP_EMPLOYEE_SCORES"
	BusinessValueTable         = "YBS_HC.SWP.TABLE::SWP_ROLE_BUSINESS_VALUE"
	ProjectSkillLevelsTable    = "YBS_HC.SWP.TABLE::SWP_PROJECTED_SKILL_LEVELS"
	ProjectSufficiencyTable    = "YBS_HC.SWP.TABLE::SWP_PROJECTED_SUFFICIENCY"
	AvailabilityTable          = "YBS_HC.SWP.TABLE::SWP_ROLE_AVAILABILITY"
	AvailabilityRanksTable     = "YBS_HC.SWP.TABLE::SWP_ROLE_AVAILABILITY_RANKS"
	PreferredStrategyTable     = "YBS_HC.SWP.TABLE::SWP_ROLE_STRATEGY"
	MarketSalariesTable        = "YBS_HC.SWP.TABLE::SWP_MARKET_SALARIES_CURRENT"
	MarketSalariesIndexedTable = "YBS_HC.SWP.TABLE::SWP_MARKET_SALARIES_INDEXED"
)

//HdbDsn returns a connection string for SAP Hana. It uses the environment variable HANA_PASSWORD to authorize
func HdbDsn() string {
	return fmt.Sprintf("hdb://Z_HRSWPWRITE:%s@%s", os.Getenv("HANA_PASSWORD"), Host)
}

func PrepQuery(table string, prid int) string {
	return fmt.Sprintf(`SELECT * FROM "%s" WHERE "PROJECT_ID" = %d`, table, prid)
}
