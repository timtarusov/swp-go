package models

import (
	"math"
	"sort"
	"sync"
)

//BuildCostForNewHires calculates the train cost for the new hires
//It is a sum of costs for the 3 cheapest skills to train
func (r *Role) BuildCostForNewHires() float64 {

	n := 3 //number of skills to consider
	var costs []float64
	var totalCost float64

	for skillID, skill := range r.Profile {
		//get the cost to reach the target level
		skillCost := r.Team.CostOfSkills[skillID][skill.TargetLevel]
		costs = append(costs, skillCost)
	}
	//sort to get the cheapest skills upfront
	sort.Float64s(costs)
	for i := 0; i < n; i++ {
		totalCost += costs[i]
	}
	return totalCost
}

//CostOfStrategy has 2 fields: flat and delta.
//Flat - the cost for new hires, Delta - for existing employees
type CostOfStrategy struct {
	Flat  map[int]float64
	Delta map[int]float64
}

//Calculates buy cost. Flat - salary of new hires plus the cost of train
//Delta - difference between flat cost and the current salary of the role
func (r *Role) CalculateBuyCost(ms *MarketSalaries) *CostOfStrategy {
	cos := &CostOfStrategy{
		Flat:  make(map[int]float64),
		Delta: make(map[int]float64),
	}
	trainCost := r.BuildCostForNewHires()
	for y, sal := range ms.GetSalary(City(r.City), Grade(r.Grade)) {
		cos.Flat[y] = sal + trainCost
	}
	for y, fl := range cos.Flat {
		cos.Delta[y] = fl - r.Salary
	}
	return cos
}

//Calculates retain cost. Flat - salary of new hires plus the cost of train
//Delta - difference between flat cost and the current salary of the role
func (r *Role) CalculateRetainCost(ms *MarketSalaries, salRetainpct float64, salRetainpctAM float64) *CostOfStrategy {
	cos := &CostOfStrategy{
		Flat:  make(map[int]float64),
		Delta: make(map[int]float64),
	}
	buy_reference := ms.GetSalary(City(r.City), Grade(r.Grade))
	years := r.Team.years
	var sal_increased float64
	for idx, y := range years {
		// #если зарплата роли меньше зарплаты рынка
		if r.Salary < buy_reference[y] {
			// #каждый год увеличиваем на Х%
			sal_increased = math.Pow((r.Salary * (1 + salRetainpct)), float64(idx)+1)
			// #увеличиваем на процент, но не больше, чем медианная зарплата по рынку
			cos.Flat[y] = math.Min(sal_increased, buy_reference[y])
		} else {
			// #каждый год увеличиваем на Х% c коэффициентом Above Median
			sal_increased = math.Pow((r.Salary * (1 + salRetainpctAM)), float64(idx)+1)
			cos.Flat[y] = sal_increased
		}
		// в первый год сравниваем с текущей зарплатой
		if idx == 0 {
			cos.Delta[y] = cos.Flat[y] - r.Salary

		} else {
			// потом с зарплатой в прошлом году
			cos.Delta[y] = cos.Flat[y] - cos.Flat[y-1]
		}
	}

	return cos
}

//Calculates build cost. Flat - the cost to train employees in a given year
func (r *Role) CalculateBuildCost(suffThreshold float64) *CostOfStrategy {
	cos := &CostOfStrategy{
		Flat: make(map[int]float64),
	}
	cands := r.BuildCandidates(suffThreshold)
	var wg sync.WaitGroup
	var mx sync.Mutex

	for _, y := range r.Team.years {
		var totalCost float64
		if len(cands[y]) > 0 {

			for _, cnd := range cands[y] {
				wg.Add(1)
				go func(cnd *Employee, y int, wg *sync.WaitGroup, mx *sync.Mutex) {
					defer wg.Done()
					var candCost float64
					for skillId, skill := range r.PriorityProfile() {
						var skillCost float64
						emplSkill := cnd.EmployeeProfile[skillId].Forecast[y].Level
						for emplSkill < skill.TargetLevel {
							skillCost += r.Team.CostOfSkills[skillId][emplSkill]
							emplSkill += 1
						}
						candCost += skillCost
					}
					mx.Lock()
					totalCost += candCost
					mx.Unlock()
				}(cnd, y, &wg, &mx)
			}
		}
		wg.Wait()
		cos.Flat[y] = totalCost
	}
	return cos
}
