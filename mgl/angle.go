package mgl

import (
	"math"
)

type Radian float64

const radToDeg = 180 / math.Pi

func (r Radian) ToDegree() Degree {
	return Degree(float64(r) * radToDeg)
}

type Degree float64

func (d Degree) ToRadian() Radian {
	return Radian(float64(d) / radToDeg)
}
