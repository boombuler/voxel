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
	b := make([]byte, 0)
	for cnt >= 0x80 {
		b = append(b, byte((cnt&0x7F)|0x80))
		cnt >>= 7
	}
	b = append(b, byte(cnt&0x7F))
	return b
}

func decodeInt(data []byte) (uint, int) {
	res := uint(0)
	shift := uint(0)

	for i := 0; i < len(data); i++ {
		b := data[i]

		res = res | (uint(b&0x7F) << shift)
		if b&0x80 != 0x80 {
			return res, i + 1
		}
		shift += 7
	}
	return res, len(data)
}
