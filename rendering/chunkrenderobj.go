package rendering

import (
	"github.com/go-gl-legacy/gl"
)

func NewRenderedChunk(c Chunk, opt Options) Renderer {
	mesh := CreateMeshFromChunk(c, opt)
	if len(mesh) == 0 {
		return RenderFunc(func() {})
	}
	if !opt.HasFlag(NO_VBO) {
		return NewCubeMesh(mesh)
	} else {
		return RenderFunc(func() {
			gl.Begin(gl.QUADS)
			defer gl.End()
			for _, v := range mesh {
				gl.Normal3f(v.Norm.X(), v.Norm.Y(), v.Norm.Z())
				gl.Color4f(v.Color.Red, v.Color.Green, v.Color.Blue, v.Color.Alpha)
				gl.Vertex3f(v.Pos.X(), v.Pos.Y(), v.Pos.Z())
			}
		})
	}
}
