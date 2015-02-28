package mgl

import (
	"math"
	"testing"
)

func Test_DegToRad(t *testing.T) {
	if !FloatEqual(float32(Degree(180).ToRadian()), float32(math.Pi)) {
		t.Fail()
	}
	if !FloatEqual(float32(Degree(90).ToRadian()), float32(math.Pi/2)) {
		t.Fail()
	}
}

func Test_RadToDeg(t *testing.T) {
	if !FloatEqual(float32(Radian(math.Pi).ToDegree()), 180) {
		t.Fail()
	}
	if !FloatEqual(float32(Radian(math.Pi/2).ToDegree()), 90) {
		t.Fail()
	}
}
