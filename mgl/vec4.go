package mgl

type Vec4 [4]float32

func (v Vec4) X() float32 {
	return v[0]
}
func (v Vec4) Y() float32 {
	return v[1]
}
func (v Vec4) Z() float32 {
	return v[2]
}
func (v Vec4) W() float32 {
	return v[3]
}
func (v1 Vec4) Vec3() Vec3 {
	return Vec3{v1[0], v1[1], v1[2]}
}

func (v1 Vec4) Add(v2 Vec4) Vec4 {
	return Vec4{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2], v1[3] + v2[3]}
}
func (v1 Vec4) Sub(v2 Vec4) Vec4 {
	return Vec4{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2], v1[3] - v2[3]}
}
