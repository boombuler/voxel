package rendering

import "github.com/go-gl/gl"

type CubeMesh struct {
	buf    gl.Buffer
	length int
}

func NewCubeMesh(verts []VertexF) *CubeMesh {
	r := new(CubeMesh)
	r.length = len(verts)

	r.buf = gl.GenBuffer()
	r.buf.Bind(gl.ARRAY_BUFFER)
	defer r.buf.Unbind(gl.ARRAY_BUFFER)

	cnt := r.length * vertexF_Size
	gl.BufferData(gl.ARRAY_BUFFER, cnt, verts, gl.STATIC_DRAW)

	return r
}

func (v *CubeMesh) Close() {
	v.buf.Delete()
}

func (v *CubeMesh) Render() {
	v.buf.Bind(gl.ARRAY_BUFFER)
	gl.EnableClientState(gl.VERTEX_ARRAY)
	defer gl.DisableClientState(gl.VERTEX_ARRAY)

	gl.InterleavedArrays(gl.C4F_N3F_V3F, vertexF_Size, nil)
	gl.DrawArrays(gl.QUADS, 0, v.length)
}
