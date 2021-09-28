package models

import (
	"fmt"
)

type Employee struct {
	ID              string
	role            *Role
	EmployeeProfile map[string]*EmployeeScore
	Sufficiency     map[int]float64
}

//NewEmployee returns a new employee
func NewEmployee(id string) *Employee {
	return &Employee{
		ID:              id,
		EmployeeProfile: make(map[string]*EmployeeScore),
		Sufficiency:     make(map[int]float64),
	}
}

//SetScore sets the current level of skill for an employee
//and calculates the gap and sufficiency
func (e *Employee) SetScore(skill string, cl int) error {
	if _, ok := e.EmployeeProfile[skill]; !ok {
		return fmt.Errorf("skill %s is not in the role profile", skill)
	}
	e.EmployeeProfile[skill].CurrentLevel.Level = cl
	if (e.EmployeeProfile[skill].RoleSkill.TargetLevel - cl) > 0 {
		e.EmployeeProfile[skill].CurrentLevel.Gap = (e.EmployeeProfile[skill].RoleSkill.TargetLevel - cl)
		e.EmployeeProfile[skill].CurrentLevel.Sufficient = false
	} else {
		e.EmployeeProfile[skill].CurrentLevel.Gap = 0
		e.EmployeeProfile[skill].CurrentLevel.Sufficient = true
	}
	return nil
}

//ProjectSkillsDevelopment calculates the forecast for learning a skill
//by an employee given the years and a speed
func (e *Employee) ProjectSkillsDevelopment(yrs []int, speed int) {
	for _, skill := range e.EmployeeProfile {
		skill.Forecast = make(map[int]Level)
		for i, y := range yrs {
			forecast := skill.CurrentLevel.Level + i*speed
			if forecast < skill.RoleSkill.TargetLevel {
				skill.Forecast[y] = Level{
					Level:      forecast,
					Gap:        skill.RoleSkill.TargetLevel - forecast,
					Sufficient: false,
				}
			} else {
				skill.Forecast[y] = Level{
					Level:      skill.RoleSkill.TargetLevel,
					Gap:        0,
					Sufficient: true,
				}
			}
		}
	}
}

//ProjectEmployeeSufficiency calculates the overall sufficiency of an employee over the years
func (e *Employee) ProjectEmployeeSufficiency(yrs []int) {
	for _, y := range yrs {
		var totalTarget float64
		var totalCurrent float64
		for _, skill := range e.EmployeeProfile {
			totalTarget += float64(skill.RoleSkill.TargetLevel)
			totalCurrent += float64(skill.Forecast[y].Level)
		}
		e.Sufficiency[y] = totalCurrent / totalTarget
	}
}
