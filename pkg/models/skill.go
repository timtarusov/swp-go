package models

import (
	"math"
	"sync"
)

//RoleSkill struct that makes up the profile
//of the role
type RoleSkill struct {
	SkillName   string
	TargetLevel int
	IsPriority  bool
}

//EmployeeScore represents the individual command of the skill by an employee
type EmployeeScore struct {
	RoleSkill    *RoleSkill
	CurrentLevel Level
	Forecast     map[int]Level
}

type Level struct {
	Level      int
	Sufficient bool
	Gap        int
}

//BuildCandidates returns the list of Employees who are suitable to be promoted for the next role
// func (r *Role) deprecatedBuildCandidates(suffThreshold float64) map[int][]*Employee {
// 	candList := make(map[int][]*Employee)
// 	for _, role := range r.Team.Roles {
// 		// #кандидаты есть в других ролях, которые ниже грейдом
// 		if role.Grade < r.Grade {
// 			// #для каждого сотрудника считаем % соответствия
// 			// на каждый год
// 			for _, y := range r.Team.years {
// 				for _, empl := range role.employees {
// 					var nom float64
// 					var denom float64
// 					for skillName, skill := range r.PriorityProfile() {
// 						level := float64(empl.EmployeeProfile[skillName].Forecast[y].Level)
// 						nom += math.Min(level, float64(skill.TargetLevel))
// 						denom += float64(skill.TargetLevel)
// 					}
// 					if suff := nom / denom; suff >= suffThreshold {
// 						candList[y] = append(candList[y], empl)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return candList
// }

type candList struct {
	sync.Mutex
	Cl map[int][]*Employee
}

func newCandList() *candList {
	return &candList{
		Cl: make(map[int][]*Employee),
	}
}

//BuildCandidates returns the list of Employees who are suitable to be promoted for the next role
func (r *Role) BuildCandidates(suffThreshold float64) map[int][]*Employee {
	candList := newCandList()
	var wg sync.WaitGroup
	wg.Add(len(r.Team.Roles))
	for roleId := range r.Team.Roles {
		//first set of goroutines to loop through roles
		go func(role *Role, wg *sync.WaitGroup) {
			defer wg.Done()
			// #кандидаты есть в других ролях, которые ниже грейдом
			if role.Grade < r.Grade {
				// #для каждого сотрудника считаем % соответствия
				// на каждый год
				var wg_y sync.WaitGroup
				wg_y.Add(len(r.Team.years))
				for _, y := range r.Team.years {
					// second set of goroutines to loop through years
					go func(y int, wg *sync.WaitGroup) {
						defer wg.Done()
						for _, empl := range role.employees {
							var nom float64
							var denom float64
							for skillName, skill := range r.PriorityProfile() {
								level := float64(empl.EmployeeProfile[skillName].Forecast[y].Level)
								nom += math.Min(level, float64(skill.TargetLevel))
								denom += float64(skill.TargetLevel)
							}
							//Lock to prevent from concurrent map write
							if suff := nom / denom; suff >= suffThreshold {
								candList.Lock()
								candList.Cl[y] = append(candList.Cl[y], empl)
								candList.Unlock()
							}
						}
					}(y, &wg_y)
				}
				wg_y.Wait()
			}
		}(r.Team.Roles[roleId], &wg)
	}
	wg.Wait()
	return candList.Cl
}

func (r *Role) CurrentGoodAndBadEmployees(year int, suffThreshold float64) (cge float64, cbe float64) {
	for _, empl := range r.employees {
		if empl.Sufficiency[year] >= suffThreshold {
			cge++
		} else {
			cbe++
		}
	}
	return
}
