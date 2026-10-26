package main

import "math"

type Stats struct {
	Count         int
	Sum, Min, Max int64
}

func Calc(nums []int64) Stats {
	var sum, min1, min2, max1, max2, diff int64 = 0, 0, math.MaxInt64, 0, math.MinInt64, 0
	count := 0
	var result Stats
	if len(nums) == 0 || len(nums) == 1 {
		result.Sum = 0
		result.Min = 0
		result.Max = 0

		return result
	}

	for i := 1; i < len(nums); i++ {
		diff = nums[i] - nums[i-1]
		min1, max1 = diff, diff
		if min1 < min2 {
			min2 = min1
		}
		if max1 > max2 {
			max2 = max1
		}
		sum += diff
		count++
	}

	result.Count = count
	result.Sum = sum
	result.Min = min2
	result.Max = max2

	return result
}
