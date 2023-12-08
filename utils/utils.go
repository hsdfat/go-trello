package utils

import "math"

func GetYValue(a float64, x int, b int32) float64 {
	return float64(x)*a + float64(b)
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
