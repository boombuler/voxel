package mgl

import (
	"testing"
)

func Test_Vec3IEquals(t *testing.T) {
	v1 := Vec3I{1, 2, 3}
	v2 := Vec3I{4, 5, 6}
	if !v1.Equals(v1) || v1.Equals(v2) || v2.Equals(v1) || !v2.Equals(v2) {
		t.Error("Vec3I Equals failed badly.")
	}
}

func Test_Vec3IAdd(t *testing.T) {
	v1 := Vec3I{1, 2, 3}
	v2 := Vec3I{4, 5, 6}

	v3 := v1.Add(v2)

	if v3[0] != 5 || v3[1] != 7 || v3[2] != 9 {
		t.Errorf("Add not adding properly")
	}

	v4 := v2.Add(v1)

	if !v3.Equals(v4) {
		t.Errorf("Addition is somehow not commutative")
	}
}

func Test_Vec3IXYZ(t *testing.T) {
	v1 := Vec3I{1, 2, 3}
	if v1.X() != 1 || v1.Y() != 2 || v1.Z() != 3 {
		t.Error("Coord functions of Vec3I missmatches")
	}
}

func Test_Vec3IMul(t *testing.T) {
	v1 := Vec3I{0, -2, 3}
	v2 := v1.Mul(-2)

	if !(Vec3I{0, 4, -6}).Equals(v2) {
		t.Errorf("Multiplication of Vec3I failed")
	}
}

func Test_Vec3ISub(t *testing.T) {
	v1 := Vec3I{10, 25, 11}
	v2 := Vec3I{0, 10, 55}

	v3 := v1.Sub(v2)

	if !(Vec3I{10, 15, -44}).Equals(v3) {
		t.Error("Sub not subtracting properly")
	}
}

func Test_Vec3IVec3(t *testing.T) {
	v1 := Vec3I{10, 25, 11}
	v2 := Vec3{10, 25, 11}

	if !v1.Vec3().Equals(v2) {
		t.Error("Sub not subtracting properly")
	}
}
