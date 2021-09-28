package tests

import (
	"testing"

	"github.com/ts.tarusov/swp/pkg/models"
)

func TestAddMarketSalaries(t *testing.T) {
	ms := models.NewMarketSalaries(0.05)

	ms.AddPayline("100", 10, 10000, []int{2021, 2022, 2023, 2024, 2025, 2026})
	ms.AddPayline("100", 11, 15000, []int{2021, 2022, 2023, 2024, 2025, 2026})
	ms.AddPayline("100", 12, 20000, []int{2021, 2022, 2023, 2024, 2025, 2026})
	ms.AddPayline("100", 13, 40000, []int{2021, 2022, 2023, 2024, 2025, 2026})
	ms.AddPayline("100", 16, 90000, []int{2021, 2022, 2023, 2024, 2025, 2026})
	ms.AddPayline("100", 19, 200000, []int{2021, 2022, 2023, 2024, 2025, 2026})

	// t.Logf("%#v", ms.YearsGrid)
}
