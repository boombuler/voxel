package rle

import (
	"github.com/boombuler/voxel/mgl"
)

const ChunkSizeX = 64
const ChunkSizeY = 64
const ChunkSizeZ = 64

func vecToIdx(v mgl.Vec3I) int {
	return (((v.Z() * ChunkSizeY) + v.Y()) * ChunkSizeX) + v.X()
}

func idxToVec(i int) mgl.Vec3I {
	return mgl.Vec3I{
		i % ChunkSizeX,
		(i / ChunkSizeX) % ChunkSizeY,
		(i / (ChunkSizeX * ChunkSizeY)) % ChunkSizeZ,
	}
}

func codeInt(cnt uint) []byte {
	if cnt == 0 {
		return []byte{0}
	}
	b := make([]byte, 0)
	for cnt > 0 {
		v := byte(cnt & 0x7F)
		cnt >>= 7
		if cnt > 0 {
			v = v | 0x80
		}

		b = append(b, v)
	}
	return b
}

func decodeInt(data []byte) (uint, int) {
	res := uint(0)
	for i := 0; i < len(data); i++ {
		b := data[i]
		res = res | (uint(b&0x7F) << uint(7*i))
		if b&0x80 != 0x80 {
			return res, i + 1
		}
	}
	return res, len(data)
}
