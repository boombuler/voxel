package rendering

import (
	"github.com/boombuler/voxel/mgl"
)

type Object interface {
	Position() mgl.Vec3
	Size() mgl.Vec3
	Renderer() Renderer
}
