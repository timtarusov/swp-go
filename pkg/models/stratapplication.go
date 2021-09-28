package models

type stratApplication struct {
	AdditionalDemand map[string]float64
	Attrition        map[string]float64
}

func NewStratApplication() *stratApplication {
	return &stratApplication{
		AdditionalDemand: make(map[string]float64),
		Attrition:        make(map[string]float64),
	}
}
