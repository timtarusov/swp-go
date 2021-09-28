package models

type MarketSalaries struct {
	IndexRate float64
	YearsGrid map[int]SalaryGrid
}

type City string
type Grade int
type Salary float64
type PayLines map[Grade]Salary
type SalaryGrid map[City]PayLines

func NewMarketSalaries(indexRate float64) *MarketSalaries {
	return &MarketSalaries{
		IndexRate: indexRate,
		YearsGrid: make(map[int]SalaryGrid),
	}
}

//AddPayline adds a new payline for marketsalaries and calculates
//its indexation over the years
func (ms *MarketSalaries) AddPayline(c City, g Grade, s Salary, yrs []int) {
	for idx, y := range yrs {
		if _, ok := ms.YearsGrid[y]; !ok {
			ms.YearsGrid[y] = make(SalaryGrid)
		}
		if _, ok := ms.YearsGrid[y][c]; !ok {
			ms.YearsGrid[y][c] = make(PayLines)
		}
		if idx == 0 {
			ms.YearsGrid[y][c][g] = s
		} else {
			for lg, ls := range ms.YearsGrid[y-1][c] {
				ms.YearsGrid[y][c][lg] = ls * Salary(1+ms.IndexRate)
			}
		}
	}
}

//GetSalary returns the market salaries for all years given a city, and a grade
func (ms *MarketSalaries) GetSalary(c City, g Grade) map[int]float64 {
	yearSal := make(map[int]float64)
	for y, sg := range ms.YearsGrid {
		yearSal[y] = float64(sg[c][g])
	}
	return yearSal
}
