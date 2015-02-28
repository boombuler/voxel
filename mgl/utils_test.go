package mgl

import (
	"math"
	"testing"
)

func Test_Round(t *testing.T) {
	tests := []struct {
		Value     float32
		Precision int
		Expected  float32
	}{
		{0.5, 0, 1},
		{0.123, 2, 0.12},
		{9.99999999, 6, 10},
		{-9.99999999, 6, -10},
		{-0.000099, 4, -0.0001},
	}

	for _, c := range tests {
		if r := Round(c.Value, c.Precision); r != c.Expected {
			t.Errorf("Round(%v, %v) != %v (got %v)", c.Value, c.Precision, c.Expected, r)
		}
	}
}

func Test_Abs(t *testing.T) {
	if !FloatEqual(0.5, Abs(0.5)) {
		t.Fail()
	}
	if !FloatEqual(0.5, Abs(-0.5)) {
		t.Fail()
	}
	if !FloatEqual(0, Abs(3-3)) {
		t.Fail()
	}
}

func Test_Equal(t *testing.T) {
	t.Parallel()

	var a float32 = 1.5
	var b float32 = 1.0 + .5

	if !FloatEqual(a, a) {
		t.Errorf("Float Equal fails on comparing a number with itself")
	}

	if !FloatEqual(a, b) {
		t.Errorf("Float Equal fails to compare two equivalent numbers with minimal drift")
	} else if !FloatEqual(b, a) {
		t.Errorf("Float Equal is not symmetric for some reason")
	}

	if !FloatEqual(0.0, 0.0) {
		t.Errorf("Float Equal fails to compare zero values correctly")
	}

	if FloatEqual(1.5, 1.51) {
		t.Errorf("Float Equal gives false positive on large difference")
	}

	if FloatEqual(1.5, 1.5000001) {
		t.Errorf("Float Equal gives false positive on small difference")
	}

	if FloatEqual(1.5, 0.0) {
		t.Errorf("Float Equal gives false positive comparing with zero")
	}
}

func Test_EqualThreshold(t *testing.T) {
	t.Parallel()

	// |1.0 - 1.01| < .1
	if !FloatEqualThreshold(1.0, 1.01, 1e-1) {
		t.Errorf("Thresholded equal returns negative on threshold")
	}

	// Comes out to |1.0 - 1.01| < .0001
	if FloatEqualThreshold(1.0, 1.01, 1e-3) {
		t.Errorf("Thresholded equal returns false positive on tolerant threshold")
	}
}

func Test_EqualThresholdTable(t *testing.T) {
	// http://floating-point-gui.de/errors/NearlyEqualsTest.java
	InfPos := float32(math.Inf(1))
	InfNeg := float32(math.Inf(-1))
	NaN := float32(math.NaN())
	MinValue := float32(math.SmallestNonzeroFloat32)
	MaxValue := float32(math.MaxFloat32)
	tests := []struct {
		A, B, Ep float32
		Expected bool
	}{
		{1.0, 1.01, 1e-1, true},
		{1.0, 1.01, 1e-3, false},

		// Regular large numbers
		{1000000.0, 1000001.0, 0.00001, true},
		{1000001.0, 1000000.0, 0.00001, true},
		{10000.0, 10001.0, 0.00001, false},
		{10001.0, 10000.0, 0.00001, false},

		// Negative large numbers
		{-1000000.0, -1000001.0, 0.00001, true},
		{-1000001.0, -1000000.0, 0.00001, true},
		{-10000.0, -10001.0, 0.00001, false},
		{-10001.0, -10000.0, 0.00001, false},

		// Numbers around 1
		{1.0000001, 1.0000002, 0.00001, true},
		{1.0000002, 1.0000001, 0.00001, true},
		{1.0002, 1.0001, 0.00001, false},
		{1.0001, 1.0002, 0.00001, false},

		// Numbers around -1
		{-1.000001, -1.000002, 0.00001, true},
		{-1.000002, -1.000001, 0.00001, true},
		{-1.0001, -1.0002, 0.00001, false},
		{-1.0002, -1.0001, 0.00001, false},

		// Numbers between 1 and 0
		{0.000000001000001, 0.000000001000002, 0.00001, true},
		{0.000000001000002, 0.000000001000001, 0.00001, true},
		{0.000000000001002, 0.000000000001001, 0.00001, false},
		{0.000000000001001, 0.000000000001002, 0.00001, false},

		// Numbers between -1 and 0
		{-0.000000001000001, -0.000000001000002, 0.00001, true},
		{-0.000000001000002, -0.000000001000001, 0.00001, true},
		{-0.000000000001002, -0.000000000001001, 0.00001, false},
		{-0.000000000001001, -0.000000000001002, 0.00001, false},

		// Comparisons involving zero
		{0.0, 0.0, 0.00001, true},
		{0.0, -0.0, 0.00001, true},
		{-0.0, -0.0, 0.00001, true},
		{0.00000001, 0.0, 0.00001, false},
		{0.0, 0.00000001, 0.00001, false},
		{-0.00000001, 0.0, 0.00001, false},
		{0.0, -0.00000001, 0.00001, false},

		// Comparisons involving infinities
		{InfPos, InfPos, 0.00001, true},
		{InfNeg, InfNeg, 0.00001, true},
		{InfNeg, InfPos, 0.00001, false},
		{InfPos, MaxValue, 0.00001, false},
		{InfNeg, -MaxValue, 0.00001, false},

		// Comparisons involving NaN values
		{NaN, NaN, 0.00001, false},
		{0.0, NaN, 0.00001, false},
		{NaN, 0.0, 0.00001, false},
		{-0.0, NaN, 0.00001, false},
		{NaN, -0.0, 0.00001, false},
		{NaN, InfPos, 0.00001, false},
		{InfPos, NaN, 0.00001, false},
		{NaN, InfNeg, 0.00001, false},
		{InfNeg, NaN, 0.00001, false},
		{NaN, MaxValue, 0.00001, false},
		{MaxValue, NaN, 0.00001, false},
		{NaN, -MaxValue, 0.00001, false},
		{-MaxValue, NaN, 0.00001, false},
		{NaN, MinValue, 0.00001, false},
		{MinValue, NaN, 0.00001, false},
		{NaN, -MinValue, 0.00001, false},
		{-MinValue, NaN, 0.00001, false},

		// Comparisons of numbers on opposite sides of 0
		{1.000000001, -1.0, 0.00001, false},
		{-1.0, 1.000000001, 0.00001, false},
		{-1.000000001, 1.0, 0.00001, false},
		{1.0, -1.000000001, 0.00001, false},
		{10 * MinValue, 10 * -MinValue, 0.00001, true},
		{10000 * MinValue, 10000 * -MinValue, 0.00001, true},

		// Comparisons of numbers very close to zero
		{MinValue, -MinValue, 0.00001, true},
		{-MinValue, MinValue, 0.00001, true},
		{MinValue, 0, 0.00001, true},
		{0, MinValue, 0.00001, true},
		{-MinValue, 0, 0.00001, true},
		{0, -MinValue, 0.00001, true},
		{0.000000001, -MinValue, 0.00001, false},
		{0.000000001, MinValue, 0.00001, false},
		{MinValue, 0.000000001, 0.00001, false},
		{-MinValue, 0.000000001, 0.00001, false},
	}

	for _, c := range tests {
		if r := FloatEqualThreshold(c.A, c.B, c.Ep); r != c.Expected {
			t.Errorf("FloatEqualThreshold(%v, %v, %v) != %v (got %v)", c.A, c.B, c.Ep, c.Expected, r)
		}
	}
}
