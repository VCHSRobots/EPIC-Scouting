/*
Package calc provides functions for calculating various game statistics.
*/
package calc

import (
	"math"
	"sort"
)

//DemocraticOutliers - Finds outliers which are not in the majority opinion i.e. are not the same value as most of the others
func DemocraticOutliers(data []float64) []float64 {
	//This will mark repeated values as outliers if there is another value more common than them
	outliers := make([]float64, 0, len(data))
	maxentries := 0
	counter := make(map[float64]int)
	for _, num := range data {
		_, ok := counter[num]
		if ok != true {
			counter[num] = 1
		} else {
			counter[num]++
		}
	}
	//Gets max number of occurances for any value
	for _, v := range counter {
		if v > maxentries {
			maxentries = v
		}
	}
	//Checks for outliers and appends them to list
	for k, v := range counter {
		if v < maxentries {
			outliers = append(outliers, k)
		}
	}
	return outliers
}

func DataSize(bytes)

//StatisticalOutliers - Checks for outliers outside +/-1.5 of the interquartile range
func StatisticalOutliers(data []float64) []float64 {
	outliers := make([]float64, 0, len(data))
	mean := mean(data)
	iqr := findIQR(data)
	for _, num := range data {
		if num > mean+iqr*1.5 || num < mean-iqr*1.5 {
			outliers = append(outliers, num)
		}
	}
	return outliers
}

func findIQR(data []float64) float64 {
	var half1 []float64
	var half2 []float64
	sort.Float64s(data)
	if len(data)%2 == 1 {
		half1 = data[:len(data)/2]
		half2 = data[len(data)/2+1:]
	} else {
		half1 = data[:len(data)/2]
		half2 = data[len(data)/2:]
	}
	return findMedian(half2) - findMedian(half1)
}

func findMedian(data []float64) float64 {
	var median float64
	sort.Float64s(data)
	if len(data)%2 == 1 {
		median = data[len(data)-1]
	} else {
		median = (data[len(data)/2] + data[len(data)/2-1]) / 2
	}
	return median
}

func standardDeviation(data []float64) float64 {
	mean := mean(data)
	var total float64 = 0
	for _, num := range data {
		total += math.Pow(num-mean, 2)
	}
	return math.Sqrt(total / float64(len(data)))
}

func mean(data []float64) float64 {
	var total float64 = 0
	for _, num := range data {
		total += num
	}
	return total / float64(len(data))
}
