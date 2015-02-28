package rendering

import (
	"fmt"
	"github.com/boombuler/voxel/mgl"
	"math"
	"sync"
	"time"
)

type faceDirection byte

const (
	left  faceDirection = 0
	right faceDirection = 1

	bottom faceDirection = 2
	top    faceDirection = 3

	back  faceDirection = 4
	front faceDirection = 5
)

func isVoxelSolid(v Voxel) bool {
	if v == nil {
		return false
	}
	c := v.Color()
	if c == nil {
		return false
	}
	_, _, _, a := c.RGBA()
	return a >= uint32(math.MaxUint16)
}

func isVoxelInvisible(v Voxel) bool {
	if v == nil {
		return true
	}
	c := v.Color()
	if c == nil {
		return true
	}
	_, _, _, a := c.RGBA()
	return a == 0
}

var (
	vLeftOf   = mgl.Vec3I{-1, 0, 0}
	vRightOf  = mgl.Vec3I{1, 0, 0}
	vBottomOf = mgl.Vec3I{0, -1, 0}
	vTopOf    = mgl.Vec3I{0, 1, 0}
	vBackOf   = mgl.Vec3I{0, 0, -1}
	vFrontOf  = mgl.Vec3I{0, 0, 1}
)

func defaultVoxelIterator(c Chunk) func(fn func(p mgl.Vec3I, v Voxel)) {
	s := c.Size()
	return func(fn func(pos mgl.Vec3I, vox Voxel)) {
		for x := 0; x < s.X(); x++ {
			for y := 0; y < s.Y(); y++ {
				for z := 0; z < s.Z(); z++ {
					p := mgl.Vec3I{x, y, z}
					vox := c.At(p)
					fn(p, vox)
				}
			}
		}
	}
}

func performCulling(c Chunk, noCulling bool) map[faceDirection]map[mgl.Vec3I]Voxel {
	result := make(map[faceDirection]map[mgl.Vec3I]Voxel)
	for f := faceDirection(0); f < faceDirection(6); f++ {
		result[f] = make(map[mgl.Vec3I]Voxel)
	}
	bounds := c.Size()
	var it func(fn func(p mgl.Vec3I, v Voxel))
	if itChunk, ok := c.(IteratableChunk); ok {
		it = itChunk.ForeachVoxel
	} else {
		it = defaultVoxelIterator(c)
	}

	it(func(p mgl.Vec3I, vox Voxel) {
		if !isVoxelInvisible(vox) {
			if noCulling || p.X() == 0 || !isVoxelSolid(c.At(p.Add(vLeftOf))) {
				result[left][p] = vox
			}
			if noCulling || p.X() == bounds.X()-1 || !isVoxelSolid(c.At(p.Add(vRightOf))) {
				result[right][p] = vox
			}
			if noCulling || p.Y() == 0 || !isVoxelSolid(c.At(p.Add(vBottomOf))) {
				result[bottom][p] = vox
			}
			if noCulling || p.Y() == bounds.Y()-1 || !isVoxelSolid(c.At(p.Add(vTopOf))) {
				result[top][p] = vox
			}
			if noCulling || p.Z() == 0 || !isVoxelSolid(c.At(p.Add(vBackOf))) {
				result[back][p] = vox
			}
			if noCulling || p.Z() == bounds.Z()-1 || !isVoxelSolid(c.At(p.Add(vFrontOf))) {
				result[front][p] = vox
			}
		}
	})
	return result
}

type meshingDirectionInfo struct {
	d1     mgl.Vec3I
	d2     mgl.Vec3I
	n      mgl.Vec3
	offset mgl.Vec3I
}

var meshingDirections map[faceDirection]meshingDirectionInfo = map[faceDirection]meshingDirectionInfo{
	front: meshingDirectionInfo{
		d1:     mgl.Vec3I{1, 0, 0},
		d2:     mgl.Vec3I{0, 1, 0},
		offset: mgl.Vec3I{0, 0, 1},
		n:      mgl.Vec3{0, 0, 1},
	},
	back: meshingDirectionInfo{
		d1:     mgl.Vec3I{1, 0, 0},
		d2:     mgl.Vec3I{0, 1, 0},
		offset: mgl.Vec3I{0, 0, 0},
		n:      mgl.Vec3{0, 0, -1},
	},
	top: meshingDirectionInfo{
		d1:     mgl.Vec3I{0, 0, 1},
		d2:     mgl.Vec3I{1, 0, 0},
		offset: mgl.Vec3I{0, 1, 0},
		n:      mgl.Vec3{0, 1, 0},
	},
	bottom: meshingDirectionInfo{
		d1:     mgl.Vec3I{0, 0, 1},
		d2:     mgl.Vec3I{1, 0, 0},
		offset: mgl.Vec3I{0, 0, 0},
		n:      mgl.Vec3{0, -1, 0},
	},
	right: meshingDirectionInfo{
		d1:     mgl.Vec3I{0, 0, 1},
		d2:     mgl.Vec3I{0, 1, 0},
		offset: mgl.Vec3I{1, 0, 0},
		n:      mgl.Vec3{1, 0, 0},
	},
	left: meshingDirectionInfo{
		d1:     mgl.Vec3I{0, 0, 1},
		d2:     mgl.Vec3I{0, 1, 0},
		offset: mgl.Vec3I{0, 0, 0},
		n:      mgl.Vec3{-1, 0, 0},
	},
}

func perfomMeshing(sides map[mgl.Vec3I]Voxel, dir faceDirection) (result []VertexF) {
	result = make([]VertexF, 0, len(sides))
	dinf := meshingDirections[dir]
	d1, d2, n, offset := dinf.d1, dinf.d2, dinf.n, dinf.offset

	d1s := d1.Mul(-1)
	d2s := d2.Mul(-1)

	for len(sides) > 0 {
		var startPos mgl.Vec3I
		var checkVal Voxel
		for k, v := range sides {
			startPos = k
			checkVal = v
			break
		}

		for {
			cpn := startPos.Add(d1s)
			v, ok := sides[cpn]
			if ok && v.Equals(checkVal) {
				startPos = cpn
			} else {
				break
			}
		}
		for {
			cpn := startPos.Add(d2s)
			v, ok := sides[cpn]
			if ok && v.Equals(checkVal) {
				startPos = cpn
			} else {
				break
			}
		}

		delete(sides, startPos)
		width := 1
		for {
			t := startPos.Add(d1.Mul(width))
			v, ok := sides[t]
			if ok && v.Equals(checkVal) {
				delete(sides, t)
				width++
			} else {
				break
			}
		}

		height := 1
		for {
			startRow := startPos.Add(d2.Mul(height))
			allOk := true
			for i := 0; i < width; i++ {
				t := startRow.Add(d1.Mul(i))
				v, ok := sides[t]
				if !ok || !v.Equals(checkVal) {
					allOk = false
					break
				}
			}
			if allOk {
				height++
				for i := 0; i < width; i++ {
					t := startRow.Add(d1.Mul(i))
					delete(sides, t)
				}
			} else {
				break
			}
		}

		pColor, ok := ColorModel.Convert(checkVal.Color()).(*Color)
		if !ok {
			continue
		}
		color := *pColor

		dm1 := mgl.Identity().TranslateVec3(d1.Mul(width).Vec3())
		dm2 := mgl.Identity().TranslateVec3(d2.Mul(height).Vec3())
		dm3 := dm1.MulMat4(dm2)
		startPosV3 := startPos.Add(offset).Vec3()
		startPosV4 := startPosV3.Vec4(1)
		result = append(result,
			VertexF{color, n, startPosV3},
			VertexF{color, n, dm1.MulVec4(startPosV4).Vec3()},
			VertexF{color, n, dm3.MulVec4(startPosV4).Vec3()},
			VertexF{color, n, dm2.MulVec4(startPosV4).Vec3()})
	}
	return result
}

func dontPerfomMeshing(sides map[mgl.Vec3I]Voxel, dir faceDirection) (result []VertexF) {
	result = make([]VertexF, 0, len(sides)*4)
	dinf := meshingDirections[dir]
	d1, d2, n, offset := dinf.d1, dinf.d2, dinf.n, dinf.offset

	for pos, vox := range sides {
		pColor, ok := ColorModel.Convert(vox.Color()).(*Color)
		if !ok {
			continue
		}
		color := *pColor
		startPos := pos.Add(offset)

		result = append(result,
			VertexF{color, n, startPos.Vec3()},
			VertexF{color, n, startPos.Add(d1).Vec3()},
			VertexF{color, n, startPos.Add(d1).Add(d2).Vec3()},
			VertexF{color, n, startPos.Add(d2).Vec3()})
	}
	return
}

func CreateMeshFromChunk(c Chunk, o Options) []VertexF {
	t0 := time.Now()
	culled := performCulling(c, o.HasFlag(NO_CULLING))
	t1 := time.Now()

	wg := new(sync.WaitGroup)
	wg.Add(6)
	results := make([][]VertexF, 6, 6)
	for face, items := range culled {
		f := face
		i := items
		go func() {
			if o.HasFlag(NO_MESHING) {
				results[f] = dontPerfomMeshing(i, f)
			} else {
				results[f] = perfomMeshing(i, f)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	cnt := 0
	for i := faceDirection(0); i < faceDirection(6); i++ {
		cnt += len(results[i])
	}
	result := make([]VertexF, 0, cnt)
	for i := faceDirection(0); i < faceDirection(6); i++ {
		result = append(result, results[i]...)
	}
	t2 := time.Now()
	fmt.Println("Quad Count: ", cnt/4)
	fmt.Println("Culling took:", t1.Sub(t0))
	fmt.Println("Greedy took:", t2.Sub(t1))

	return result
}
