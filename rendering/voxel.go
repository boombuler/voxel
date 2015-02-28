package rendering

import (
	"github.com/boombuler/voxel/mgl"
	"image/color"
)

type Voxel interface {
	Color() color.Color
	Equals(v Voxel) bool
}

type Chunk interface {
	Size() mgl.Vec3I
	At(pos mgl.Vec3I) Voxel
}

type IteratableChunk interface {
	Chunk
	ForeachVoxel(fn func(pos mgl.Vec3I, vox Voxel))
}
