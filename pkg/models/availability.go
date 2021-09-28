package models

type Availability struct {
	PotentialHires   map[int]float64
	PercentAvailable float64
	Rank             int
}

func NewAvailability() *Availability {
	return &Availability{
		PotentialHires:   make(map[int]float64),
		PercentAvailable: 0.0,
		Rank:             0,
	}
}

func (r *Role) SetAvailabilityRank(rank int) {
	r.Availability.Rank = rank
}

//AddAvailability adds the number of potential hires for the role for a given yearz.
func (r *Role) AddAvailability(year int, potHires float64) {
	r.Availability.PotentialHires[year] = potHires
}

// AvailabilityRanks calculates role's rank on the availability scale
func (r *Role) AvailabilityRanks() {
	var firstYear int = 9999
	for y := range r.Availability.PotentialHires {
		if y < firstYear {
			firstYear = y
		}
	}

	r.Availability.PercentAvailable = r.Availability.PotentialHires[firstYear] / float64(len(r.employees))
	r.Availability.Rank = Bisect(r.Availability.PercentAvailable)
}
