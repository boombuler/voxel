package mgl

type Vec3I [3]int

func (v1 Vec3I) Add(v2 Vec3I) Vec3I {
	return Vec3I{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}
func (v1 Vec3I) Sub(v2 Vec3I) Vec3I {
	return Vec3I{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}
func (v1 Vec3I) Mul(val int) Vec3I {
	return Vec3I{v1[0] * val, v1[1] * val, v1[2] * val}
}
func (v1 Vec3I) Div(val int) Vec3I {
	return Vec3I{v1[0] / val, v1[1] / val, v1[2] / val}
}
func (v1 Vec3I) Mod(val int) Vec3I {
	return Vec3I{v1[0] % val, v1[1] % val, v1[2] % val}
}
func (v Vec3I) X() int {
	return v[0]
}
func (v Vec3I) Y() int {
	return v[1]
}
func (v Vec3I) Z() int {
	return v[2]
}
func (v Vec3I) Vec3() Vec3 {
	return Vec3{float32(v[0]), float32(v[1]), float32(v[2])}
}
func (v1 Vec3I) Equals(v2 Vec3I) bool {
	return v1[0] == v2[0] &&
		v1[1] == v2[1] &&
		v1[2] == v2[2]
}
