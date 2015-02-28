package rle

import (
	"bytes"
	"github.com/boombuler/voxel/mgl"
	r "github.com/boombuler/voxel/rendering"
)

type ChunkData struct {
	Palette []r.Voxel
	buf     *bytes.Buffer
}

func (c *ChunkData) palIndex(vox r.Voxel) int {
	for i, v := range c.Palette {
		if v != nil {
			if v == vox {
				return i
			}
		} else if vox == nil {
			return i
		}
	}
	idx := len(c.Palette)
	c.Palette = append(c.Palette, vox)
	return idx
}

func (data *UncompressedChunkData) Compress() *ChunkData {
	result := &ChunkData{
		Palette: []r.Voxel{},
		buf:     new(bytes.Buffer),
	}

	curIdx := -1
	curCnt := uint(0)

	for _, v := range data {
		idx := result.palIndex(v)
		if idx == curIdx {
			curCnt++
		} else {
			if curCnt > 0 {
				result.buf.Write(codeInt(uint(curIdx)))
				result.buf.Write(codeInt(curCnt))
			}
			curCnt = 1
			curIdx = idx
		}
	}
	if curCnt > 0 {
		result.buf.Write(codeInt(uint(curIdx)))
		result.buf.Write(codeInt(curCnt))
	}
	return result
}

func (c *ChunkData) iterate(fn func(startIdx, cnt int, vox r.Voxel) bool) {
	d := c.buf.Bytes()
	i := 0

	for len(d) > 0 {
		idx, bytes := decodeInt(d)
		d = d[bytes:]
		cnt, bytes := decodeInt(d)
		d = d[bytes:]
		if !fn(i, int(cnt), c.Palette[idx]) {
			return
		}
		i += int(cnt)
	}
}

func (c *ChunkData) Uncompress() *UncompressedChunkData {
	cd := NewUncompressedChunkData()

	c.iterate(func(i, cnt int, vox r.Voxel) bool {
		for j := 0; j < cnt; j++ {
			cd[i+j] = vox
		}
		return true
	})

	return cd
}

func (c *ChunkData) ForeachVoxel(fn func(pos mgl.Vec3I, vox r.Voxel)) {
	c.iterate(func(i, cnt int, vox r.Voxel) bool {
		if vox != nil {
			for j := 0; j < cnt; j++ {
				fn(idxToVec(i+j), vox)
			}
		}
		return true
	})
}

func (c *ChunkData) Size() mgl.Vec3I {
	return mgl.Vec3I{ChunkSizeX, ChunkSizeY, ChunkSizeZ}
}

// At returns the voxel at the given position. (Warning this very slow!)
func (c *ChunkData) At(pos mgl.Vec3I) r.Voxel {
	tIdx := vecToIdx(pos)

	var rVox r.Voxel = nil
	c.iterate(func(i, cnt int, vox r.Voxel) bool {
		if i <= tIdx && tIdx < (i+cnt) {
			rVox = vox
			return false
		}
		return true
	})
	return rVox
}
