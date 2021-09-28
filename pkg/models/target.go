package models

import "fmt"

type Target struct {
	Name          int
	Type          string
	AbsValues     map[int]float64
	OverallGrowth float64
	PercentOfMax  float64
	years         []int
	YoYGrowth     map[int]float64
}

//NewTarget returns a new Target object
func NewTarget(n int, t string) *Target {
	return &Target{
		Name:          n,
		Type:          t,
		AbsValues:     make(map[int]float64),
		OverallGrowth: 0,
		PercentOfMax:  0,
		years:         []int{},
		YoYGrowth:     make(map[int]float64),
	}
}

//FirstYear returns the value of the first year for the target
func (t Target) FirstYear() float64 {
	return t.AbsValues[t.years[0]]
}

//SecondYear returns the value of the second year for the target
func (t Target) SecondYear() float64 {
	if len(t.years) > 1 {
		return t.AbsValues[t.years[1]]
	} else {
		return t.FirstYear()
	}
}

//LastYear returns the value of the last year for the target
func (t Target) LastYear() float64 {
	if len(t.years) > 1 {
		return t.AbsValues[t.years[len(t.years)-1]]
	} else {
		return t.FirstYear()
	}
}

//ThisYear returns the value of this year fro the target
func (t Target) ThisYear(y int) (float64, error) {
	for _, yr := range t.years {
		if y == yr {
			return t.AbsValues[y], nil
		}
	}
	return 0, fmt.Errorf("%v doesn't exits in target %d", y, t.Name)
}

func (t *Target) Years() []int {
	return t.years
}
