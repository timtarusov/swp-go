package models

import "fmt"

type PreferredStrategy struct {
	Quadrant  string
	Strategy1 string
	Strategy2 string
}

func NewPreferredStrategy() *PreferredStrategy {
	return &PreferredStrategy{}
}

func (r *Role) CalculatePreferredStrategy() error {
	if (r.Rank == 0) || (r.Availability.Rank == 0) {
		return fmt.Errorf("business value and availability are not set")
	}

	if r.Rank > 5 {
		if r.Availability.Rank <= 5 {
			r.PreferredStrategy.Quadrant = "1"
			r.PreferredStrategy.Strategy1 = "Retain"
			r.PreferredStrategy.Strategy2 = "Build"
		} else {
			r.PreferredStrategy.Quadrant = "4"
			r.PreferredStrategy.Strategy1 = "Retain"
			r.PreferredStrategy.Strategy2 = "Buy"
		}
	} else {
		if r.Availability.Rank <= 5 {
			r.PreferredStrategy.Quadrant = "2"
			r.PreferredStrategy.Strategy1 = "Borrow"
			r.PreferredStrategy.Strategy2 = "Build"
		} else {
			r.PreferredStrategy.Quadrant = "3"
			r.PreferredStrategy.Strategy1 = "Build"
			r.PreferredStrategy.Strategy2 = "Buy"
		}
	}
	return nil
}
