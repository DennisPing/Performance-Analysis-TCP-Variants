package common

import "math"

// Get the sum from a slice of float64
func Sum(arr []float64) float64 {
	var sum float64
	for _, v := range arr {
		sum += v
	}
	return sum
}

// Get the max from a slice of float64
func Max(arr []float64) float64 {
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max
}

// Get the min from a slice of float64
func Min(arr []float64) float64 {
	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}
	return min
}

// Get the avg from a slice of float64
func Mean(arr []float64) float64 {
	var sum float64
	for _, v := range arr {
		sum += v
	}
	return sum / float64(len(arr))
}

// Get the standard deviation from a slice of float64
func StdDev(arr []float64) float64 {
	avg := Mean(arr)
	var sum float64
	for _, v := range arr {
		sum += math.Pow(v-avg, 2)
	}
	return math.Sqrt(sum / float64(len(arr)))
}
