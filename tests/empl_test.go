package tests

import (
	"testing"

	"github.com/ts.tarusov/swp/pkg/models"
)

func TestSetScore(t *testing.T) {
	role := models.NewRole(1, "")

	empl := models.NewEmployee("70004710")
	role.AddEmployee(empl)

	role.AddRoleSkill(&models.RoleSkill{
		SkillName:   "1",
		TargetLevel: 4,
		IsPriority:  true,
	})
	err := empl.SetScore("1", 3)
	if err != nil {
		t.Fatal(err)
	}

}
