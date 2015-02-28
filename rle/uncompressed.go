package rle

import (
	"github.com/boombuler/voxel/mgl"
	r "github.com/boombuler/voxel/rendering"
	"sync"
)

var uncompressedChunkPool = sync.Pool{
	New: func() interface{} {
		return new(UncompressedChunkData)
	},
}

func NewUncompressedChunkData() *UncompressedChunkData {
	v, ok := uncompressedChunkPool.Get().(*UncompressedChunkData)
	if ok {
		return v
	} else {
		return new(UncompressedChunkData)
	}
}

func FreeUncompressedChunkData(data **UncompressedChunkData) {
	uncompressedChunkPool.Put(*data)
	*data = nil
}

type UncompressedChunkData [ChunkSizeX * ChunkSizeY * ChunkSizeZ]r.Voxel

func (u *UncompressedChunkData) ForeachVoxel(fn func(pos mgl.Vec3I, vox r.Voxel)) {
	for i, vox := range u {
		if vox != nil {
			fn(idxToVec(i), vox)
		}
	}
}

func (u *UncompressedChunkData) Set(pos mgl.Vec3I, vox r.Voxel) {
	u[vecToIdx(pos)] = vox
}

func (u *UncompressedChunkData) At(pos mgl.Vec3I) r.Voxel {
	return u[vecToIdx(pos)]
}

func (u *UncompressedChunkData) Size() mgl.Vec3I {
	return mgl.Vec3I{ChunkSizeX, ChunkSizeY, ChunkSizeZ}
}
