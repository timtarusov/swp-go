package tests

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ts.tarusov/swp/pkg/models"
)

func benchBuildCandidates(qnt int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		team := models.NewTeam("go team", 1)
		team.SetYears([]int{2021, 2022, 2023, 2024, 2025})
		for id := 0; id < qnt; id++ {
			team.AddRole(models.NewRole(id, "some role"))
			if id == 0 {
				team.Roles[id].Grade = 17

			} else {
				team.Roles[id].Grade = 15
			}
			for g := 0; g < qnt; g++ {

				team.Roles[id].AddEmployee(models.NewEmployee(fmt.Sprintf("ID%d", g)))
			}
			for j := 0; j < qnt; j++ {
				if err := team.Roles[id].AddRoleSkill(&models.RoleSkill{
					SkillName:   fmt.Sprintf("Skill_%d", j),
					TargetLevel: 4,
					IsPriority:  true,
				}); err != nil {
					b.Fatal(err)
				}
			}

			for h := 0; h < qnt; h++ {
				empl := team.Roles[id].Employees()[fmt.Sprintf("ID%d", h)]
				// t.Log(empl)
				for k := 0; k < qnt; k++ {
					empl.SetScore(fmt.Sprintf("Skill_%d", k), rand.Intn(5))
				}
				empl.ProjectSkillsDevelopment(team.Years(), 1)
			}
		}
		// cnd := team.Roles[0].deprecatedBuildCandidates(0.6)
		team.Roles[0].BuildCandidates(0.6)
		// b.Log(len(cnd.Cl))
	}
}

func BenchmarkBuildCandidates10(b *testing.B) {
	benchBuildCandidates(10, b)
}

func BenchmarkBuildCandidates20(b *testing.B) {
	benchBuildCandidates(20, b)
}
func BenchmarkBuildCandidates30(b *testing.B) {
	benchBuildCandidates(30, b)
}
func BenchmarkBuildCandidates40(b *testing.B) {
	benchBuildCandidates(40, b)
}
