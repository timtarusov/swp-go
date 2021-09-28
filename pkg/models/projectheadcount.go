package models

type ProjectedHeadcount struct {
	Growth    map[int]float64
	Headcount map[int]float64
}

func NewProjectedHeadcount() *ProjectedHeadcount {
	return &ProjectedHeadcount{
		Growth:    make(map[int]float64),
		Headcount: make(map[int]float64),
	}
}

// ProjectAdditionalDemand return the additional demand in workforce for the year
func (r *Role) ProjectAdditionalDemand(year int) float64 {
	//в первый год доп потребность ноль
	if year == r.Team.years[0] {
		return 0.
	} else {
		//в последующие смотрим разницу между прогнозом на этот год и предыдущим
		return r.ProjectedHeadcount.Headcount[year] - r.ProjectedHeadcount.Headcount[year-1]
	}

}
