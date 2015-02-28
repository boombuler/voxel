package rle

import (
	"github.com/boombuler/voxel/mgl"
	"testing"
)

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
