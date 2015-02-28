package mgl

import "math"

const (
	defaultEpsilon = float32(1e-10)
	minNormal      = float32(1.1754943508222875e-38) // 1 / 2**(127 - 1)
)

func Round(val float32, digits int) float32 {
	pow := math.Pow10(digits)
	valx := pow * float64(val)
	_, div := math.Modf(valx)
	if div >= 0.5 {
		valx = math.Ceil(valx)
	} else {
		valx = math.Floor(valx)
	}
	return float32(valx / pow)
}

func Abs(a float32) float32 {
	if a < 0 {
		return -a
	} else if a == 0 {
		return 0
	}

	return a
}

func FloatEqual(a, b float32) bool {
	return FloatEqualThreshold(a, b, defaultEpsilon)
}

func FloatEqualThreshold(a, b, epsilon float32) bool {
	if a == b { // Handles the case of inf or shortcuts the loop when no significant error has accumulated
		return true
	}

	diff := Abs(a - b)
	if a*b == 0 || diff < minNormal { // If a or b are 0 or both are extremely close to it
		return diff < epsilon*epsilon
	}

	// Else compare difference
	return diff/(Abs(a)+Abs(b)) < epsilon
}
