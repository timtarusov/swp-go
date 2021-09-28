package tests

import (
	"testing"

	"github.com/ts.tarusov/swp/pkg/models"
)

func TestSliceToPercentiles(t *testing.T) {
	slice := []float64{0.3, 0.2, 0.1, 0.55, 0.6, 0.14, 0.9, 2.15}
	// nums := []int{}
	// for i := range slice {
	// 	nums = append(nums, i)
	// }
	// ts := models.TwoSlices{Main: nums, Other: slice}
	// t.Logf("Before: %v", ts.Main)
	// sort.Sort(models.SortByOther(ts))
	// t.Logf("After: %v", ts.Main)

	// res := models.SliceToPercentile(slice)
	// for i := range slice {
	// 	t.Logf("%v is in %v percentile\n", slice[i], res[i])
	// }

	models.SliceToPercentile(slice)
}

func TestBisectBusinessValue(t *testing.T) {
	slice := []float64{0.3, 0.2, 0.1, 0.55, 0.6, 0.14, 0.9, 2.15}
	percentiles := models.SliceToPercentile(slice)
	for _, pct := range percentiles {
		models.Bisect(pct)
		// t.Logf("%.2f is in %.2f percentile and ranks #%v", slice[i], pct, r)
	}
}
