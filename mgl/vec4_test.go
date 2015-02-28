package mgl

import (
	"testing"
)

func Test_Vec4Add(t *testing.T) {
	v1 := Vec4{1.0, 2.5, 1.1, 2.0}
	v2 := Vec4{0.0, 1.0, 9.9, 100.0}

	v3 := v1.Add(v2)

	if !FloatEqual(v3[0], 1.0) || !FloatEqual(v3[1], 3.5) || !FloatEqual(v3[2], 11.0) || !FloatEqual(v3[3], 102.0) {
		t.Errorf("Add not adding properly")
	}

	v4 := v2.Add(v1)

	if !FloatEqual(v3[0], v4[0]) || !FloatEqual(v3[0], v4[0]) || !FloatEqual(v3[2], v4[2]) || !FloatEqual(v3[3], v4[3]) {
		t.Errorf("Addition is somehow not commutative")
	}
}

func Test_Vec4Sub(t *testing.T) {
	v1 := Vec4{1.0, 2.5, 1.1, 2.0}
	v2 := Vec4{0.0, 1.0, 9.9, 100.0}

	v3 := v1.Sub(v2)

	// 1.1-9.9 does stupid things to floats, so we need a more tolerant threshold
	if !FloatEqual(v3[0], 1.0) || !FloatEqual(v3[1], 1.5) || !FloatEqualThreshold(v3[2], -8.8, 1e-5) || !FloatEqual(v3[3], -98.0) {
		t.Errorf("Sub not subtracting properly [%f, %f, %f, %f]", v3[0], v3[1], v3[2], v3[3])
	}
}

func Test_Vec4Vec3(t *testing.T) {
	v1 := Vec4{1, 2, 3, 4}
	var v2 Vec3 = v1.Vec3()
	if !FloatEqual(v1.X(), v2.X()) || !FloatEqual(v1.Y(), v2.Y()) || !FloatEqual(v1.Z(), v2.Z()) {
		t.Error("Vec4.Vec3() not respecting coordinates")
	}
}

func Test_Vec4XYZW(t *testing.T) {
	v1 := Vec4{1, 2, 3, 4}
	if !FloatEqual(v1.X(), 1) || !FloatEqual(v1.Y(), 2) || !FloatEqual(v1.Z(), 3) || !FloatEqual(v1.W(), 4) {
		t.Error("Coord functions of Vec4 missmatches")
	}
}
