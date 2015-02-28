package mgl

import "math"

type Vec3 [3]float32

func (v1 Vec3) Normalize() Vec3 {
	l := 1.0 / v1.Len()
	return Vec3{v1[0] * l, v1[1] * l, v1[2] * l}
}
func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}
func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}
func (v1 Vec3) Mul(c float32) Vec3 {
	return Vec3{v1[0] * c, v1[1] * c, v1[2] * c}
}
func (v1 Vec3) Dot(v2 Vec3) float32 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}
func (v1 Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{v1[1]*v2[2] - v1[2]*v2[1], v1[2]*v2[0] - v1[0]*v2[2], v1[0]*v2[1] - v1[1]*v2[0]}
}
func (v1 Vec3) Len() float32 {
	return float32(math.Sqrt(float64(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2])))
}
func (v1 Vec3) Vec4(w float32) Vec4 {
	return Vec4{v1[0], v1[1], v1[2], w}
}
func (v Vec3) Vec3I() Vec3I {
	return Vec3I{
		int(Round(v[0], 0)),
		int(Round(v[1], 0)),
		int(Round(v[2], 0)),
	}
}
func (v Vec3) X() float32 {
	return v[0]
}
func (v Vec3) Y() float32 {
	return v[1]
}
func (v Vec3) Z() float32 {
	return v[2]
}
func (v1 Vec3) Equals(v2 Vec3) bool {
	for i := range v1 {
		if !FloatEqual(v1[i], v2[i]) {
			return false
		}
	}
	return true
}
