package tests

import (
	"regexp"
	"testing"

	"github.com/ts.tarusov/swp/internal/dbmodels"
)

func TestNexVal(t *testing.T) {
	var pattern = regexp.MustCompile(`SWP_(.*)$`)
	tn := pattern.FindString(dbmodels.ProjectedHeadcountTable)

	t.Logf(`SELECT "YBS_HC.SWP.SEQ::SEQ_%s".NEXTVAL FROM DUMMY`, tn[4:])

}
