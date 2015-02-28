package mgl

import (
	"testing"
)

func Test_Vec3Equal(t *testing.T) {
	tests := []struct {
		v1 Vec3
		v2 Vec3
		r  bool
	}{
		{Vec3{13, 0, 0}, Vec3{13, 0, 0}, true},
		{Vec3{1, 2, 3}, Vec3{1, 2, 3}, true},
		{Vec3{0, 0, 0}, Vec3{1, 2, 3}, false},
		{Vec3{1, 2, 0}, Vec3{1, 2, 3}, false},
		{Vec3{1, 1, 3}, Vec3{1, 2, 3}, false},
		{Vec3{2, 2, 3}, Vec3{1, 2, 3}, false},
	}
	for _, test := range tests {
		if test.v1.Equals(test.v2) != test.r {
			t.Errorf("Equal Failed on: %v <-> %v", test.v1, test.v2)
		}
	}
}
func Test_Vec3Add(t *testing.T) {
	v1 := Vec3{1.0, 2.5, 1.1}
	v2 := Vec3{0.0, 1.0, 9.9}

	v3 := v1.Add(v2)

	if !FloatEqual(v3[0], 1.0) || !FloatEqual(v3[1], 3.5) || !FloatEqual(v3[2], 11.0) {
		t.Errorf("Add not adding properly")
	}

	v4 := v2.Add(v1)

	if !FloatEqual(v3[0], v4[0]) || !FloatEqual(v3[0], v4[0]) || !FloatEqual(v3[2], v4[2]) {
		t.Errorf("Addition is somehow not commutative")
	}
}
func Test_Vec3Sub(t *testing.T) {
	v1 := Vec3{1.0, 2.5, 1.1}
	v2 := Vec3{0.0, 1.0, 9.9}

	v3 := v1.Sub(v2)

	if !FloatEqual(v3[0], 1.0) || !FloatEqual(v3[1], 1.5) || !FloatEqualThreshold(v3[2], -8.8, 1e-5) {
		t.Errorf("Sub not subtracting properly [%f, %f, %f]", v3[0], v3[1], v3[2])
	}
}
func Test_Vec3Mul(t *testing.T) {
	v := Vec3{1.0, 0.0, 100.1}
	v = v.Mul(15.0)

	if !FloatEqual(v[0], 15.0) || !FloatEqual(v[1], 0.0) || !FloatEqual(v[2], 1501.5) {
		t.Errorf("Vec mul does something weird [%f, %f, %f]", v[0], v[1], v[2])
	}
}
func Test_Vec3Length(t *testing.T) {
	tests := []struct {
		v Vec3
		l float32
	}{
		{Vec3{13, 0, 0}, 13},
		{Vec3{2, 3, 4}, 5.38516480713450403125},
		{Vec3{1, 1, 1}, 1.73205080756887729352},
	}
	for _, test := range tests {
		if l := test.v.Len(); !FloatEqual(l, test.l) {
			t.Errorf("Invalid length for vector %v -> %f", test.v, l)
		}
	}
}
func Test_Vec3Normalize(t *testing.T) {
	v1 := Vec3{10, 0, 0}
	v2 := Vec3{23.31, 0, 0}

	v1n := v1.Normalize()
	v2n := v2.Normalize()

	vn := Vec3{1, 0, 0}

	if !v1n.Equals(vn) {
		t.Errorf("failed to normalize v1 got: %v", v1n)
	}

	if !v2n.Equals(vn) {
		t.Errorf("failed to normalize v2 got: %v", v1n)
	}

	if !FloatEqual((Vec3{12, 31, 10.023}).Normalize().Len(), 1) {
		t.Error("Normalized vector should have a length of 1")
	}
}
func Test_Vec3Vec3I(t *testing.T) {
	v := Vec3{0.4, 1.5, -5.9}
	ve := Vec3I{0, 2, -6}
	if !v.Vec3I().Equals(ve) {
		t.Error("Failed Vec3 to Vec3I")
	}
}
func Test_Vec3XYZ(t *testing.T) {
	v1 := Vec3{1, 2, 3}
	if !FloatEqual(v1.X(), 1) || !FloatEqual(v1.Y(), 2) || !FloatEqual(v1.Z(), 3) {
		t.Error("Coord functions of Vec3 missmatches")
	}
}
func Test_Vec3Vec4(t *testing.T) {
	tests := []struct {
		v Vec3
		w float32
		r Vec4
	}{
		{Vec3{13, 0, 0}, 1, Vec4{13, 0, 0, 1}},
		{Vec3{2, 3, 4}, 0.4, Vec4{2, 3, 4, 0.4}},
		{Vec3{1, 1, 1}, 0, Vec4{1, 1, 1, 0}},
	}
	for _, test := range tests {
		res := test.v.Vec4(test.w)
		for i := 0; i < 4; i++ {
			if res[i] != test.r[i] {
				t.Errorf("Failed To Convert Vec3 (%v) to Vec4 got: %v", test.v, res)
			}
		}
	}
}
func Test_Vec3Cross(t *testing.T) {
	v1 := Vec3{2, 6, 3}
	v2 := Vec3{2, 1, -2}
	vr := Vec3{-15, 10, -10}

	if !v1.Cross(v2).Equals(vr) {
		t.Error("v1 x v2 not correct")
	}
	if v2.Cross(v1).Equals(vr) {
		t.Error("Cross is somehow commutative")
	}
	if !v2.Cross(v1).Mul(-1).Equals(vr) {
		t.Error("-(v2 x v1) not correct")
	}
}
func Test_Vec3Dot(t *testing.T) {
	v1 := Vec3{1, 3, 5}
	v2 := Vec3{2, 0, 4}
	res := float32(22)

	if !FloatEqual(v1.Dot(v2), res) {
		t.Error("v1 * v2 not correct")
	}
	if !FloatEqual(v2.Dot(v1), res) {
		t.Error("Dot is somehow not commutative")
	}
}
