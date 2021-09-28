package models

import (
	"sort"
)

type TwoSlices struct {
	Main  []int
	Other []float64
}

type SortByOther TwoSlices

func (sbo SortByOther) Len() int {
	return len(sbo.Main)
}

func (sbo SortByOther) Less(i, j int) bool {
	return sbo.Other[i] < sbo.Other[j]
}

func (sbo SortByOther) Swap(i, j int) {
	sbo.Main[i], sbo.Main[j] = sbo.Main[j], sbo.Main[i]
	sbo.Other[i], sbo.Other[j] = sbo.Other[j], sbo.Other[i]
}

//SliceToPercentile converts the list of float values to their
//respective percentile values within this list
func SliceToPercentile(sl []float64) []float64 {
	var nums []int
	for id := range sl {
		nums = append(nums, id)
	}
	var original []int
	copy(original, nums)
	ts := TwoSlices{Main: nums, Other: sl}
	sort.Sort(SortByOther(ts))

	result := make([]float64, len(sl))
	for idx, rank := range nums {
		original_idx := nums[idx]
		result[original_idx] = float64(rank) / float64(len(nums)-1)

	}
	return result

}

//BisectBusinessValue finds the rank of the role by its percentile
//compared to the predefined interval
func Bisect(t float64) int {
	// percentiles := SliceToPercentile(tg)
	var intervals = [...]float64{0, 0.05, 0.1, 0.2, 0.4, 0.6, 0.7, 0.8, 0.9, 1}
	var ranks = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	bv := sort.Search(len(ranks), func(i int) bool { return intervals[i] >= t })
	return ranks[bv]

}

//SliceIndex returns the index of an element in the slice
//https://stackoverflow.com/questions/8307478/how-to-find-out-element-position-in-slice
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
