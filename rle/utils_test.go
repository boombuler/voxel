package rle

import (
	"bytes"
	"github.com/boombuler/voxel/mgl"
	"testing"
)

func Test_CodeInt(t *testing.T) {
	tests := []struct {
		n uint
		r []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{128, []byte{0x80, 0x01}},
		{256, []byte{0x80, 0x02}},
	}
	for _, tst := range tests {
		res := codeInt(tst.n)
		if bytes.Compare(res, tst.r) != 0 {
			t.Errorf("Failed to code int %v\nGot: %v\nExpected:%v", tst.n, res, tst.r)
		}
	}
}

func Test_IntToVec_VecToInt(t *testing.T) {
	idx_Exp := 0
	for z := 0; z < ChunkSizeZ; z++ {
		for y := 0; y < ChunkSizeY; y++ {
			for x := 0; x < ChunkSizeX; x++ {
				var vecTest = mgl.Vec3I{x, y, z}
				idx := vecToIdx(vecTest)
				if idx != idx_Exp {
					t.Errorf("vecToIdx failed got %v expected %v", idx, idx_Exp)
					return
				}
				idx_Exp++

				vec := idxToVec(idx)
				if !vec.Equals(vecTest) {
					t.Errorf("idxToVec failed got %v expected %v", vec, vecTest)
				}
			}
		}
	}
}
