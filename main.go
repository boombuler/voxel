package main

import (
	_ "github.com/boombuler/voxel/kv6"
	kv6 "github.com/boombuler/voxel/magica"
	"github.com/boombuler/voxel/mgl"
	"github.com/boombuler/voxel/rendering"
	"os"
)

var (
	obj rendering.Renderer
	cam *camera
)

func main() {
	StartEngine(Options{
		WindowWidth:  640,
		WindowHeight: 480,
		WindowTitle:  "Voxel",
		LoadFunc:     LoadObjects,
		UpdateFunc:   Update,
	})
}

func LoadObjects(e *Engine) {

	obj, err := loadModelFile()
	if err != nil {
		panic(err.Error())
	}
	e.RenderObjects = append(e.RenderObjects, obj)
}

func Update(dt float64, e *Engine) {
	// Nothing to do here :)
}

type ChunkObj struct {
	size     mgl.Vec3
	pos      mgl.Vec3
	renderer rendering.Renderer
}

func (co *ChunkObj) Position() mgl.Vec3 {
	return co.pos
}
func (co *ChunkObj) Size() mgl.Vec3 {
	return co.size
}
func (co *ChunkObj) Renderer() rendering.Renderer {
	return co.renderer
}

func loadModelFile() (rendering.Object, error) {
	fn := "chr_knight.vox"
	if len(os.Args) > 1 {
		fn = os.Args[1]
	}
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	vf, err := kv6.Read(f)
	if err != nil {
		return nil, err
	}
	ro := rendering.NewRenderedChunk(vf, rendering.NONE)

	return &ChunkObj{
		size:     vf.Size().Vec3(),
		pos:      mgl.Vec3{0, 0, 0},
		renderer: ro,
	}, nil
}
