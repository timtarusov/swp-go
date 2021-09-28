package models

import (
	"fmt"
	"sort"
)

//Team struct
type Team struct {
	Name         string
	Function     int
	Targets      map[int]*Target
	Roles        map[int]*Role
	years        []int
	CostOfSkills map[string]costOfSkill
}

type costOfSkill map[int]float64

//NewTeam returns a new Team instance
func NewTeam(tname string, tfunc int) *Team {
	return &Team{
		Name:         tname,
		Function:     tfunc,
		Targets:      make(map[int]*Target),
		Roles:        make(map[int]*Role),
		years:        []int{},
		CostOfSkills: make(map[string]costOfSkill),
	}
}

//AddRole adds a new role to the team
func (t *Team) AddRole(r *Role) error {
	if _, ok := t.Roles[r.Name]; ok {
		return fmt.Errorf("%d already exists", r.Name)
	}
	t.Roles[r.Name] = r
	r.Team = t
	return nil
}

//AddTarget adds a new target to the team strategy
func (t *Team) AddTarget(n int, tp string, y int, val float64) error {
	if _, ok := t.Targets[n]; !ok {
		t.Targets[n] = NewTarget(n, tp)
	}
	t.Targets[n].AbsValues[y] = val
	t.Targets[n].years = append(t.Targets[n].years, y)
	return nil

}

func (t *Team) extractYears() {
	years := make(map[int]bool) //set
	for _, tgt := range t.Targets {
		for y := range tgt.AbsValues {
			years[y] = true
		}
	}
	var yearsList = make([]int, 0, len(years))
	for y := range years {
		yearsList = append(yearsList, y)
	}
	sort.Ints(yearsList)
	t.years = yearsList
}

//Years returns the list of years for the team
func (t *Team) Years() []int {
	if len(t.years) == 0 {
		t.extractYears()
		return t.years
	}
	return t.years
}

//SetYears sets the list of years for the team
func (t *Team) SetYears(years []int) {
	t.years = years
}

func (t *Team) RecalculateStrategy() error {
	t.extractYears()

	var maxTgtVal float64
	var maxTgt *Target
	for _, tgt := range t.Targets {
		if tgt.Type != "t2m" {
			if tgt.LastYear() > maxTgtVal {
				maxTgtVal = tgt.LastYear()
				maxTgt = tgt
			}
		} else {
			tgt.PercentOfMax = 1
		}
	}
	maxTgt.OverallGrowth = maxTgt.LastYear() / maxTgt.FirstYear()

	for _, tgt := range t.Targets {
		if tgt.FirstYear() != 0 {
			tgt.OverallGrowth = tgt.LastYear() / tgt.FirstYear()
		} else {
			tgt.OverallGrowth = tgt.LastYear() / tgt.SecondYear()
		}
		if tgt.Type != "t2m" {
			tgt.PercentOfMax = tgt.LastYear() / maxTgt.LastYear()
		}
	}

	return nil
}

//GenerateYoYTargets generates target growth for each year.
//Growth depends on the type of target
func (t *Team) GenerateYoYTargets() error {
	for _, tgt := range t.Targets {
		yoy := make(map[int]float64)
		for i, y := range tgt.years {
			var yoyValue float64
			if i == 0 {
				yoyValue = 0
			} else {
				thisYear, err := tgt.ThisYear(y)
				if err != nil {
					return fmt.Errorf("couldn't get the value of this year: %v", err)
				}
				lastYear, err := tgt.ThisYear(y - 1)
				if err != nil {
					return fmt.Errorf("couldn't get the value of last year: %v", err)
				}
				nextYear, err := tgt.ThisYear(y + 1)
				if err != nil {
					nextYear = thisYear
				}
				nextYearPlusOne, err := tgt.ThisYear(y + 2)
				if err != nil {
					nextYearPlusOne = nextYear
				}

				switch tgt.Type {
				case "short":
					//Year over year growth =  значение этого года разделить на значение предыдущего, минус 1
					yoyValue = (thisYear / lastYear) - 1
				case "middle":
					//За два года
					yoyValue = (thisYear+nextYear)/(lastYear+thisYear) - 1
				case "long":
					//За три года
					yoyValue = (thisYear+nextYear+nextYearPlusOne)/(lastYear+thisYear+nextYear) - 1
				case "t2m":
					//Если цель time 2 market, то делим значение предыдущего года, то значение текущего, минус 1
					yoyValue = (lastYear / thisYear) - 1
				}
			}
			yoy[y] = yoyValue
		}
		tgt.YoYGrowth = yoy
	}
	return nil
}

//GenerateRolesPercent generates the values of how much
//the role should grow
func (t *Team) GenerateRolesPercent() error {
	var sliceOfGrowths []float64
	var ranks []int
	for _, role := range t.Roles {
		var totalGrowth float64
		for tid, impc := range role.Impact {
			totalGrowth += impc * t.Targets[tid].OverallGrowth * t.Targets[tid].PercentOfMax
		}
		role.TotalGrowth = totalGrowth
		sliceOfGrowths = append(sliceOfGrowths, totalGrowth)
	}

	percentiles := SliceToPercentile(sliceOfGrowths)
	for _, pct := range percentiles {
		r := Bisect(pct)
		ranks = append(ranks, r)
	}

	for _, role := range t.Roles {
		var idx int
		for i, val := range sliceOfGrowths {
			if role.TotalGrowth == val {
				idx = i
				break
			}
		}
		role.Percentile = percentiles[idx]
		role.Rank = ranks[idx]
	}

	return nil
}

//AddCostOfSkill adds the cost of skill for 4 levels of expertise
func (t *Team) AddCostOfSkill(skillName string, lvl1 float64, lvl2 float64, lvl3 float64, lvl4 float64) {
	if _, ok := t.CostOfSkills[skillName]; !ok {
		t.CostOfSkills[skillName] = make(costOfSkill)
	}
	t.CostOfSkills[skillName][1] = lvl1
	t.CostOfSkills[skillName][2] = lvl2
	t.CostOfSkills[skillName][3] = lvl3
	t.CostOfSkills[skillName][4] = lvl4
}
