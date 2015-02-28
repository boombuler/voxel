package rendering

import (
	"github.com/boombuler/voxel/mgl"
)

type VertexF struct {
	Color
	Norm mgl.Vec3
	Pos  mgl.Vec3
}

const vector3f_Size int = 3 * 4

const vertexF_Size int = color_Size + (2 * vector3f_Size)
