package dbmodels

type ProjectedHeadcount struct {
	Id        int     `db:"ID"`
	Function  int     `db:"FUNCTION"`
	Role      int     `db:"ROLE"`
	Type      string  `db:"TYPE"`
	Year      int     `db:"YEAR"`
	ProjectId int     `db:"PROJECT_ID"`
	Value     float64 `db:"VALUE"`
	CreatedAt string  `db:"CREATED_AT"`
}

type ProjectedGrowth struct {
	Id        int     `db:"ID"`
	Function  int     `db:"FUNCTION"`
	Role      int     `db:"ROLE"`
	Type      string  `db:"TYPE"`
	Year      int     `db:"YEAR"`
	ProjectId int     `db:"PROJECT_ID"`
	Value     float64 `db:"VALUE"`
	CreatedAt string  `db:"CREATED_AT"`
}
