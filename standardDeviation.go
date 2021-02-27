package main

import (
	"math"
)

func standardDev(values []int64) float64 {
	var sum, mean, sd float64
	len := len(values)
	for i := 0; i < len; i++ {
		sum += float64(values[i])
	}
	mean = sum / float64(len)
	for j := 0; j < len; j++ {

		sd += math.Pow(float64(values[j])-mean, 2)
	}
	sd = math.Sqrt(sd / float64(len))

	return sd
}
