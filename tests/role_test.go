package tests

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/ts.tarusov/swp/pkg/models"
)

func TestAddRoleSkill(t *testing.T) {
	myFirstRole := models.NewRole(1, "")
	if err := myFirstRole.AddRoleSkill(&models.RoleSkill{
		SkillName:   "go programming language",
		TargetLevel: 4,
		IsPriority:  true,
	}); err != nil {
		t.Log(err)
	}
	if err := myFirstRole.AddRoleSkill(&models.RoleSkill{
		SkillName:   "design patterns",
		TargetLevel: 4,
		IsPriority:  true,
	}); err != nil {
		t.Log(err)
	}
	if len(myFirstRole.Profile) != 2 {
		t.Errorf("Expected length 2, got %d", len(myFirstRole.Profile))
	}
}

func TestGetPriorityProfile(t *testing.T) {
	role := models.NewRole(1, "")
	if err := role.AddRoleSkill(&models.RoleSkill{
		SkillName:   "go programming language",
		TargetLevel: 4,
		IsPriority:  true,
	}); err != nil {
		t.Log(err)
	}
	if err := role.AddRoleSkill(&models.RoleSkill{
		SkillName:   "design patterns",
		TargetLevel: 4,
		IsPriority:  false,
	}); err != nil {
		t.Log(err)
	}
	pprofile := role.PriorityProfile()
	if len(pprofile) != 1 {
		t.Errorf("Expected length 1, got %d", len(pprofile))
	}
}

func TestAddEmployee(t *testing.T) {
	role := models.NewRole(1, "")

	empl1 := &models.Employee{
		ID: "7000",
	}
	empl2 := &models.Employee{
		ID: "7001",
	}
	empl3 := &models.Employee{
		ID: "7000",
	}

	if err := role.AddEmployee(empl1); err != nil {
		t.Error(err)
	}
	if err := role.AddEmployee(empl2); err != nil {
		t.Error(err)
	}
	if err := role.AddEmployee(empl3); err == nil {
		t.Error("Expected to fail with duplicate IDs")
	}

	empls := role.Employees()
	if len(empls) != 2 {
		t.Errorf("Expected length 2, got %d", len(empls))
	}
}

func TestSetImpact(t *testing.T) {
	role := models.NewRole(1, "")
	team := models.NewTeam("go team", 1)
	role.Team = team
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

	if err := role.SetImpactForTarget(1, 0.5); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if role.Impact[1] != 0.5 {
		t.Error("failed to set impact")
	}
	if err := role.SetImpactForTarget(2, -0.5); err == nil {
		t.Error("expected to fail - negative value")
	}
	if err := role.SetImpactForTarget(3, 2); err == nil {
		t.Error("expected to fail - value too high")
	}
	if err := role.SetImpactForTarget(1, 1); err == nil {
		t.Error("expected to fail - duplicate value")
	}

}

func TestProjectHeadcount(t *testing.T) {
	role := models.NewRole(1, "")
	team := models.NewTeam("go team", 1)
	role.Team = team

	for i := 1; i < 100; i++ {
		if err := role.AddEmployee(models.NewEmployee(strconv.Itoa(i))); err != nil {
			t.Errorf("failed to add an employee: %v", err)
		}
	}
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
	if err := role.SetImpactForTarget(1, 0.5); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := role.SetImpactForTarget(2, 0.7); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := role.SetImpactForTarget(3, 0.8); err != nil {
		t.Errorf("couldn't set impact: %v", err)
	}
	if err := team.GenerateYoYTargets(); err != nil {
		t.Errorf("couldn't generate yoy targets: %s", err)
	}
	if err := role.ProjectHeadcount(0.15); err != nil {
		t.Errorf("couldn't project headcount: %s", err)
	}
	t.Log(role.ProjectedHeadcount.Growth)
	t.Log(role.ProjectedHeadcount.Headcount)
}

func TestBuildCandidates(t *testing.T) {
	team := models.NewTeam("go team", 1)
	team.SetYears([]int{2021, 2022, 2023, 2024, 2025})
	for id := 0; id < 10; id++ {
		team.AddRole(models.NewRole(id, "some role"))
		if id == 0 {
			team.Roles[id].Grade = 17

		} else {
			team.Roles[id].Grade = 15
		}
		for g := 0; g < 10; g++ {

			team.Roles[id].AddEmployee(models.NewEmployee(fmt.Sprintf("ID%d", g)))
		}
		for j := 0; j < 10; j++ {
			if err := team.Roles[id].AddRoleSkill(&models.RoleSkill{
				SkillName:   fmt.Sprintf("Skill_%d", j),
				TargetLevel: 4,
				IsPriority:  true,
			}); err != nil {
				t.Fatal(err)
			}
		}

		for h := 0; h < 10; h++ {
			empl := team.Roles[id].Employees()[fmt.Sprintf("ID%d", h)]
			// t.Log(empl)
			for k := 0; k < 10; k++ {
				empl.SetScore(fmt.Sprintf("Skill_%d", k), rand.Intn(5))
			}
			empl.ProjectSkillsDevelopment(team.Years(), 1)
		}
	}
	// cnd := team.Roles[0].deprecatedBuildCandidates(0.6)
	cnd := team.Roles[0].BuildCandidates(0.6)
	// t.Log(cnd)
	for year, c := range cnd {
		t.Logf("In %d role %d has %d candidates", year, 0, len(c))
	}

}

func TestCostOfBuild(t *testing.T) {
	team := models.NewTeam("go team", 1)
	team.SetYears([]int{2021, 2022, 2023, 2024, 2025})
	for id := 0; id < 10; id++ {
		team.AddRole(models.NewRole(id, "some role"))
		if id == 0 {
			team.Roles[id].Grade = 17

		} else {
			team.Roles[id].Grade = 15
		}
		for g := 0; g < 10; g++ {

			team.Roles[id].AddEmployee(models.NewEmployee(fmt.Sprintf("ID%d", g)))
		}
		for j := 0; j < 10; j++ {
			if err := team.Roles[id].AddRoleSkill(&models.RoleSkill{
				SkillName:   fmt.Sprintf("Skill_%d", j),
				TargetLevel: 4,
				IsPriority:  true,
			}); err != nil {
				t.Fatal(err)
			}
			team.AddCostOfSkill(
				fmt.Sprintf("Skill_%d", j),
				math.Round(rand.Float64()*0),
				math.Round(rand.Float64()*30000),
				math.Round(rand.Float64()*40000),
				math.Round(rand.Float64()*50000),
			)
		}

		for h := 0; h < 10; h++ {
			empl := team.Roles[id].Employees()[fmt.Sprintf("ID%d", h)]
			// t.Log(empl)
			for k := 0; k < 10; k++ {
				empl.SetScore(fmt.Sprintf("Skill_%d", k), rand.Intn(5))
			}
			empl.ProjectSkillsDevelopment(team.Years(), 1)
		}
	}
	cos := team.Roles[0].CalculateBuildCost(0.7)
	t.Logf("%#v", cos)
}
