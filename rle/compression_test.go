package rle

import (
	"bytes"
	"github.com/boombuler/voxel/mgl"
	"image/color"
	"math/rand"
	"testing"
)

type testVoxel int

func (tv testVoxel) Color() color.Color {
	return color.Black
}

func Test_CompressEmpty(t *testing.T) {
	ucd := NewUncompressedChunkData()
	c := ucd.Compress()
	cnt := ChunkSizeX * ChunkSizeY * ChunkSizeZ

	expected := append([]byte{0}, codeInt(uint(cnt))...)

	if bytes.Compare(expected, c.buf.Bytes()) != 0 {
		t.Errorf("Compression error!\nGot: %v\nExpected:%v", c.buf.Bytes(), expected)
	}
}

func Test_AtCompressed(t *testing.T) {
	ucd := NewUncompressedChunkData()
	defer FreeUncompressedChunkData(&ucd)
	for x := 0; x < ChunkSizeX; x++ {
		for y := 0; y < ChunkSizeY; y++ {
			for z := 0; z < ChunkSizeZ; z++ {
				val := ((3 * x) + (17 * y) + z) % 3
				ucd.Set(mgl.Vec3I{x, y, z}, testVoxel(val))
			}
		}
	}
	cd := ucd.Compress()

	for x := 0; x < ChunkSizeX; x++ {
		for y := 0; y < ChunkSizeY; y++ {
			for z := 0; z < ChunkSizeZ; z++ {
				pos := mgl.Vec3I{x, y, z}
				uv := ucd.At(pos)
				cv := cd.At(pos)
				if uv != cv {
					t.Errorf("At failed at position %v. Got %v expected %v", pos, cv, uv)
					return
				}
			}
		}
	}
}

func BenchmarkAtCompressed(b *testing.B) {
	b.StopTimer()

	ucd := NewUncompressedChunkData()
	defer FreeUncompressedChunkData(&ucd)
	for x := 0; x < ChunkSizeX; x++ {
		for y := 0; y < ChunkSizeY; y++ {
			for z := 0; z < ChunkSizeZ; z++ {
				val := ((3 * x) + (17 * y) + z) % 3
				ucd.Set(mgl.Vec3I{x, y, z}, testVoxel(val))
			}
		}
	}
	cd := ucd.Compress()
	rand.Seed(1337)
	positions := make([]mgl.Vec3I, 0)
	for i := 0; i < 10000; i++ {
		positions = append(positions, idxToVec(int(rand.Int31n(ChunkSizeX*ChunkSizeY*ChunkSizeZ))))
	}

	for i := 0; i < b.N; i++ {
		pos := positions[i%len(positions)]
		b.StartTimer()
		_ = cd.At(pos)
		b.StopTimer()
	}
}

func BenchmarkCompress(b *testing.B) {
	b.StopTimer()
	ucd := NewUncompressedChunkData()

	for x := 0; x < ChunkSizeX; x++ {
		for y := 0; y < ChunkSizeY; y++ {
			for z := 0; z < ChunkSizeZ; z++ {
				val := ((3 * x) + (17 * y) + z) % 3
				ucd.Set(mgl.Vec3I{x, y, z}, testVoxel(val))
			}
		}
	}

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_ = ucd.Compress()
		b.StopTimer()
	}
	ucs := len(*ucd)
	c := ucd.Compress()
	cp := (float64(c.buf.Len()) / float64(ucs)) * 100
	b.Logf("Compressed To: %v bytes --> %v%%", c.buf.Len(), cp)
}

func BenchmarkDecompress(b *testing.B) {
	b.StopTimer()
	ucd := NewUncompressedChunkData()
	for x := 0; x < ChunkSizeX; x++ {
		for y := 0; y < ChunkSizeY; y++ {
			for z := 0; z < ChunkSizeZ; z++ {
				val := ((3 * x) + (17 * y) + z) % 3
				ucd.Set(mgl.Vec3I{x, y, z}, testVoxel(val))
			}
		}
	}
	cd := ucd.Compress()
	FreeUncompressedChunkData(&ucd)
	if ucd != nil {
		b.Fail()
	}
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		data := cd.Uncompress()
		b.StopTimer()
		FreeUncompressedChunkData(&data)
	}
}
