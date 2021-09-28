package models

import (
	"fmt"
	"math"
)

//Role struct
type Role struct {
	Name               int
	Description        string
	Profile            map[string]*RoleSkill
	employees          map[string]*Employee
	Team               *Team
	handicap           float64
	Impact             map[int]float64
	ProjectedHeadcount *ProjectedHeadcount
	TotalGrowth        float64
	Percentile         float64
	Rank               int
	Availability       *Availability
	PreferredStrategy  *PreferredStrategy
	Grade              int
	City               string
	Salary             float64
}

//NewRole returns a new role struct
func NewRole(r int, d string) *Role {
	return &Role{
		Name:               r,
		Description:        d,
		Impact:             make(map[int]float64),
		employees:          make(map[string]*Employee),
		Profile:            make(map[string]*RoleSkill),
		ProjectedHeadcount: NewProjectedHeadcount(),
		handicap:           1,
		TotalGrowth:        0,
		Percentile:         0,
		Rank:               0,
		Availability:       NewAvailability(),
		PreferredStrategy:  NewPreferredStrategy(),
	}
}

//AddRoleSkill adds a new skill to the role profile
func (r *Role) AddRoleSkill(rs *RoleSkill) error {
	if _, ok := r.Profile[rs.SkillName]; ok {
		return fmt.Errorf("%s already exists", rs.SkillName)
	}
	r.Profile[rs.SkillName] = rs
	for _, empl := range r.employees {
		empl.EmployeeProfile[rs.SkillName] = &EmployeeScore{
			RoleSkill: rs,
		}
		// log.Printf("%#v", empl.EmployeeProfile)
	}
	return nil
}

//PriorityProfile returns those skills in the profile
//which are considered to be in priority
func (r *Role) PriorityProfile() map[string]*RoleSkill {
	pprofile := make(map[string]*RoleSkill)
	for name, rs := range r.Profile {
		if rs.IsPriority {
			pprofile[name] = rs
		}
	}
	return pprofile
}

//Employees returns the list of employees
func (r *Role) Employees() map[string]*Employee {
	return r.employees
}

//AddEmployee adds a new employee to the role.
//If the role already has the employee with this ID,
//then raises an error
func (r *Role) AddEmployee(empl *Employee) error {
	if _, ok := r.employees[empl.ID]; ok {
		return fmt.Errorf("%s already exists", empl.ID)
	}
	r.employees[empl.ID] = empl
	empl.role = r
	return nil

}

//SetImpactForTarget sets the impact value for a specific target
func (r *Role) SetImpactForTarget(target int, impact float64) error {
	if (impact < 0) || (impact > 1) {
		return fmt.Errorf("expected value between 0 and 1, got %v", impact)
	}

	if _, ok := r.Impact[target]; ok {
		return fmt.Errorf("%v already exists", target)
	}
	r.Impact[target] = impact
	return nil
}

//SetHandicap sets the restriction in growth for the role
func (r *Role) SetHandicap(hndcp float64) error {
	if (hndcp < 0) && (hndcp > 1) {
		return fmt.Errorf("expected value between 0 and 1, got %v", hndcp)
	}
	r.handicap = hndcp
	return nil
}

//ProjectHeadcoutn calculates two maps.
//Growth - the percent of YoY growth of people in the role.
//And Headcount - the actual fte for the role
func (r *Role) ProjectHeadcount(pl float64) error {
	for _, year := range r.Team.years {
		var yearGrowth float64
		var runningGrowth float64
		//Годовой прирост роли = прирост цели * влияние роли на цель. Складываем
		//влияние на все цели
		for _, tgt := range r.Team.Targets {
			yearGrowth += tgt.YoYGrowth[year] * r.Impact[tgt.Name]
		}
		//Умножаем прирост на ограничитель роли и на рычаг персонала
		yearGrowth = yearGrowth * r.handicap * pl
		r.ProjectedHeadcount.Growth[year] = yearGrowth
		for _, g := range r.ProjectedHeadcount.Growth {
			runningGrowth += g
		}
		//Рост численности в роли = текущая численность * накопленный прирост.
		//Округляем до целого человека
		r.ProjectedHeadcount.Headcount[year] = math.Round(float64(len(r.employees)) * (1 + runningGrowth))
	}
	return nil
}

func (r *Role) SetBusinessValueRank(rank int) {
	r.Rank = rank
}
