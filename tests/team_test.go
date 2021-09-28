package tests

import (
	"reflect"
	"testing"

	"github.com/ts.tarusov/swp/pkg/models"
)

func TestAddRole(t *testing.T) {
	team := models.NewTeam("Go team", 1)
	rl1 := models.NewRole(1, "")
	rl2 := models.NewRole(2, "")
	rl3 := models.NewRole(2, "")
	if err := team.AddRole(rl1); err != nil {
		t.Error(err)
	}
	if err := team.AddRole(rl2); err != nil {
		t.Error(err)
	}
	if err := team.AddRole(rl3); err == nil {
		t.Error("duplicate values")
	}
	if len(team.Roles) != 2 {
		t.Errorf("expected length of 2, got %d", len(team.Roles))
	}
}

// func TestExtractYears(t *testing.T) {
// 	team := models.NewTeam("Go team", "dev")
// team.AddTarget("t1", "short", map[int]float64{
// 	2020: 100,
// 	2021: 200,
// 	2022: 300,
// 	2025: 400,
// 	2024: 300,
// 	2023: 200,
// })
// team.AddTarget("t2", "short", map[int]float64{
// 	2020: 100,
// 	2021: 200,
// 	2022: 300,
// 	2025: 400,
// 	2024: 300,
// 	2023: 200,
// })
// team.AddTarget("t3", "short", map[int]float64{
// 	2020: 100,
// 	2021: 200,
// 	2022: 300,
// 	2025: 400,
// 	2024: 300,
// 	2023: 200,
// })
// 	team.ExtractYears()
// 	t.Log(team.Years())
// 	if len(team.Years()) != 6 {
// 		t.Errorf("expected length of 6, got %d", len(team.Years()))
// 	}
// }

func TestAddTarget(t *testing.T) {
	team := models.NewTeam("Go team", 1)
	years := []int{2020, 2021, 2022, 2023, 2024, 2025}
	vals := []float64{100, 200, 300, 200, 300, 400}
	for i := range years {
		team.AddTarget(1, "short", years[i], vals[i])
	}
	for i := range years {
		team.AddTarget(2, "short", years[i], vals[i])
	}
	for i := range years {
		team.AddTarget(3, "short", years[i], vals[i])
	}

	if err := team.GenerateYoYTargets(); err != nil {
		t.Errorf("couldn't generate yoy targets: %s", err)
	}
	tg1 := team.Targets[1]

	etal := map[int]float64{
		2020: 0,
		2021: 1,
		2022: 0.5,
		2023: -0.33333333333333337,
		2024: 0.5,
		2025: 0.33333333333333326,
	}
	eq := reflect.DeepEqual(tg1.YoYGrowth, etal)
	if !eq {
		t.Error("wrong calculations")
	}

}

func TestGenerateRolePercents(t *testing.T) {
	team := models.NewTeam("Go team", 1)
	years := []int{2020, 2021, 2022, 2023, 2024, 2025}
	vals := []float64{100, 200, 300, 200, 300, 400}
	for i := range years {
		team.AddTarget(1, "short", years[i], vals[i])
	}
	for i := range years {
		team.AddTarget(2, "short", years[i], vals[i])
	}
	for i := range years {
		team.AddTarget(3, "short", years[i], vals[i])
	}

	if err := team.GenerateYoYTargets(); err != nil {
		t.Errorf("couldn't generate yoy targets: %s", err)
	}

	rl1 := models.NewRole(1, "")
	rl2 := models.NewRole(2, "")
	rl3 := models.NewRole(3, "")
	if err := team.AddRole(rl1); err != nil {
		t.Error(err)
	}
	if err := team.AddRole(rl2); err != nil {
		t.Error(err)
	}
	if err := team.AddRole(rl3); err != nil {
		t.Error(err)
	}

	if err := rl1.SetImpactForTarget(1, 0.5); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := rl1.SetImpactForTarget(2, 0.1); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := rl1.SetImpactForTarget(3, 0.6); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}

	if err := rl2.SetImpactForTarget(1, 0.1); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := rl2.SetImpactForTarget(2, 0.9); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := rl2.SetImpactForTarget(3, 0.4); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := rl3.SetImpactForTarget(1, 0.2); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := rl3.SetImpactForTarget(2, 0.3); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := rl3.SetImpactForTarget(3, 0.89); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := team.GenerateRolesPercent(); err != nil {
		t.Errorf("couldn't generate percent: %v", err)
	}
	for _, r := range team.Roles {
		t.Logf("Role: %d, TG: %.2f, PCT: %.2f, RANK: %d\n", r.Name, r.TotalGrowth, r.Percentile, r.Rank)
	}
}
