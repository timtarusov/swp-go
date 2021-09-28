package dbmodels

import (
	"database/sql"

	"github.com/SAP/go-hdb/driver"
)

//EmplrolesRow corresponds to the EMPROLE table in DB
type EmplrolesRow struct {
	Id         string `json:"id" db:"ID"`
	Employee   string `json:"employee" db:"EMPLOYEE"`
	Function   int    `json:"function" db:"FUNCTION"`
	Function_T string `json:"function_t" db:"FUNCTION_T"`
	Role       int    `json:"role" db:"ROLE"`
	Role_T     string `json:"role_t" db:"ROLE_T"`
	ProjectId  int    `json:"project_id" db:"PROJECT_ID"`
	CreatedAt  string `json:"created_at" db:"CREATED_AT"`
}

//RoleImpactRow corresponds to the ROLE_IMPACT table in DB
type RoleImpactRow struct {
	Id        int            `db:"ID"`
	Role      int            `db:"ROLE"`
	Target    int            `db:"TARGET"`
	ProjectId int            `db:"PROJECT_ID"`
	Impact    driver.Decimal `db:"IMPACT"`
	Handicap  driver.Decimal `db:"HANDICAP"`
	CreatedAt string         `db:"CREATED_AT"`
	Dept      sql.NullString `db:"DEPT,omitempty"`
}

//TargetsRow corresponds to the TARGETS table in DB
type TargetsRow struct {
	Id        int            `db:"ID"`
	Target    int            `db:"TARGET"`
	Target_T  string         `db:"TARGET_T"`
	Year      int            `db:"YEAR"`
	ProjectId int            `db:"PROJECT_ID"`
	Type      int            `db:"TYPE"`
	Type_T    string         `db:"TYPE_T"`
	Value     driver.Decimal `db:"VALUE"`
	CreatedAt string         `db:"CREATED_AT"`
	Units     sql.NullString `db:"UNITS"`
}

//ProjectConfigRow corresponds to the PROJECT_CONFIG table in DB
type ProjectConfigRow struct {
	Id              int            `db:"ID"`
	ProjectID       int            `db:"PROJECT_ID"`
	PersLeverage    driver.Decimal `db:"PERS_LEVERAGE"`
	SuffThreshold   driver.Decimal `db:"SUFFICIENCY_THRESHOLD"`
	SalRetainPerc   driver.Decimal `db:"SALARY_RETAIN_PERC"`
	SalRetainPercAM driver.Decimal `db:"SALARY_RETAIN_PERC_AM"`
	MrkSalInc       driver.Decimal `db:"MARKET_SALARIES_INC"`
	BuySR           driver.Decimal `db:"BUY_SUPPLY_RATIO"`
	BuildSR         driver.Decimal `db:"BUILD_SUPPLY_RATIO"`
	BorrowSR        driver.Decimal `db:"BORROW_SUPPLY_RATIO"`
	AttrGood        driver.Decimal `db:"ATTR_GOOD_RATIO"`
	CreatedAt       string         `db:"CREATED_AT"`
}

//RoleProfileRow corresponds to the ROLE_PROFILE table in DB
type RoleProfileRow struct {
	Id         int    `db:"ID"`
	ProjectId  int    `db:"PROJECT_ID"`
	Role       int    `db:"ROLE"`
	Skill      string `db:"SKILL"`
	Skill_T    string `db:"SKILL_T"`
	Type       int    `db:"TYPE"`
	Type_T     string `db:"TYPE_T"`
	TargetLvl  int    `db:"TARGET_LEVEL"`
	IsPriority string `db:"IS_PRIORITY"`
	CreatedAt  string `db:"CREATED_AT"`
}

//EmployeeScoresRow corresponds to the EMPLOYEE_SCORES table in DB
type EmployeeScoresRow struct {
	Role         int    `db:"ROLE"`
	Employee     string `db:"EMPLOYEE"`
	Skill        string `db:"SKILL"`
	CurrentLevel int    `db:"CURRENT_LEVEL"`
}

//AvailabilityRow corresponds to the AVAILABILITY table in DB
type AvailabilityRow struct {
	Id             int            `db:"ID"`
	ProjectId      int            `db:"PROJECT_ID"`
	Role           int            `db:"ROLE"`
	Year           int            `db:"YEAR"`
	PotentialHires driver.Decimal `db:"POTENTIAL_HIRES"`
	CreatedAt      string         `db:"CREATED_AT"`
}

//AvailabilityRankRow corresponds to the AVAILABILITY_RANKS table in DB
type AvailabilityRankRow struct {
	Id               int            `db:"ID"`
	ProjectId        int            `db:"PROJECT_ID"`
	Role             int            `db:"ROLE"`
	Rank             int            `db:"RANK"`
	PercentAvailable driver.Decimal `db:"PERCENT_AVAILABLE"`
	CreatedAt        string         `db:"CREATED_AT"`
}

//BusinessValueRow corresponds to the BUSINESS_VALUE table in DB
type BusinessValueRow struct {
	Id          int            `db:"ID"`
	ProjectId   int            `db:"PROJECT_ID"`
	Role        int            `db:"ROLE"`
	TotalGrowth driver.Decimal `db:"TOTAL_GROWTH"`
	Percentile  driver.Decimal `db:"PERCENTILE"`
	Rank        int            `db:"RANK"`
	CreatedAt   string         `db:"CREATED_AT"`
}

//MarketSalariesRow corresponds to the MARKET_SALARIES_CURRENT table in DB
type MarketSalariesRow struct {
	Id         int            `db:"ID"`
	FunctionId int            `db:"FUNCTION_ID"`
	City       string         `db:"CITY"`
	Grade      string         `db:"GRADE"`
	Salary     driver.Decimal `db:"SALARY"`
	CreatedAt  string         `db:"CREATED_AT"`
}
