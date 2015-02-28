package kv6

import (
	"bufio"
	"errors"
	"github.com/boombuler/voxel/mgl"
	r "github.com/boombuler/voxel/rendering"
	"image/color"
	"io"
)

type kv6Vox struct {
	R, G, B byte
}

type kv6Block struct {
	kv6Vox
	ZPos int
}

func (b kv6Vox) Color() color.Color {
	return color.RGBA{b.R, b.G, b.B, 255}
}

type KV6File struct {
	size    mgl.Vec3I
	content []*kv6Block
	xpos    []int
	xypos   []int
}

// OpenGL Coords -> File Coords
func (f *KV6File) unrotateCoords(pos mgl.Vec3I) mgl.Vec3I {
	return mgl.Vec3I{
		pos.X(),
		pos.Z(),
		pos.Y() - f.size.Z() - 1,
	}
}

// File Coords -> OpenGL Coords
func (f *KV6File) rotateCoords(pos mgl.Vec3I) mgl.Vec3I {
	return mgl.Vec3I{
		pos.X(),
		f.size.Z() - pos.Z() - 1,
		pos.Y(),
	}
}

func (f *KV6File) ForeachVoxel(fn func(pos mgl.Vec3I, vox r.Voxel)) {
	x, y := 0, 0

	xMax := f.xpos[x]
	yMax := f.xypos[(x*f.size.Y())+y]
	for _, block := range f.content {
		for xMax == 0 {
			x++
			if x >= f.size.X() {
				break
			}
			y = 0
			yMax = f.xypos[(x*f.size.Y())+y]
			xMax = f.xpos[x]
		}
		for yMax == 0 {
			y++
			if y >= f.size.Y() {
				y = 0
				break
			}
			yMax = f.xypos[(x*f.size.Y())+y]
		}

		fn(f.rotateCoords(mgl.Vec3I{x, y, block.ZPos}), block.kv6Vox)
		yMax--

		xMax--
	}
}

func (f *KV6File) Size() mgl.Vec3I {
	return mgl.Vec3I{
		f.size.X(),
		f.size.Z(),
		f.size.Y(),
	}
}

func (f *KV6File) At(pos mgl.Vec3I) r.Voxel {
	pos = f.unrotateCoords(pos)
	idx := 0
	for x := 0; x < pos.X(); x++ {
		idx += f.xpos[x]
	}
	for y := 0; y < pos.Y(); y++ {
		idx += f.xypos[(pos.X()*f.size.Y())+y]
	}
	cnt := f.xypos[(pos.X()*f.size.Y())+pos.Y()]

	for i := 0; i < cnt; i++ {
		blk := f.content[idx+i]
		if blk.ZPos == pos.Z() {
			return blk.kv6Vox
		}
		if blk.ZPos > pos.Z() {
			return nil
		}
	}
	return nil
}

type reader struct {
	*bufio.Reader
}

func (vr *reader) readStr4() (string, error) {
	txt, err := vr.readBytes(4)
	if err != nil {
		return "", err
	}
	return string(txt), err
}
func (vr *reader) skip(size int) error {
	for size > 512 {
		_, err := vr.readBytes(512)
		if err != nil {
			return err
		}
		size -= 512
	}
	if size > 0 {
		_, err := vr.readBytes(size)
		if err != nil {
			return err
		}
	}
	return nil
}
func errorOrTxt(err error, txt string) error {
	if err == nil {
		return errors.New(txt)
	}
	return err
}
func (vr *reader) readI32() (uint32, error) {
	buf, err := vr.readBytes(4)
	if err != nil {
		return 0, err
	}
	return uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24, nil
}
func (vr *reader) readI16() (uint16, error) {
	buf, err := vr.readBytes(2)
	if err != nil {
		return 0, err
	}
	return uint16(buf[0]) | uint16(buf[1])<<8, nil
}
func (vr *reader) readBytes(cnt int) ([]byte, error) {
	buf := make([]byte, cnt)
	n, err := vr.Read(buf)
	if n != cnt {
		return nil, errors.New("unexpected end of file")
	}
	return buf, err
}

const fHeader = "Kvxl"

func Read(rd io.Reader) (*KV6File, error) {
	r := &reader{bufio.NewReader(rd)}
	if head, err := r.readStr4(); err != nil || head != fHeader {
		return nil, errorOrTxt(err, "Invalid file format")
	}
	xSize, err := r.readI32()
	if err != nil {
		return nil, err
	}
	ySize, err := r.readI32()
	if err != nil {
		return nil, err
	}
	zSize, err := r.readI32()
	if err != nil {
		return nil, err
	}
	err = r.skip(3 * 4) // Skip Pivot XYZ
	if err != nil {
		return nil, err
	}
	blkLen, err := r.readI32()
	if err != nil {
		return nil, err
	}
	blocks := make([]*kv6Block, 0, blkLen)
	for i := 0; i < int(blkLen); i++ {
		b := new(kv6Block)
		color, err := r.readBytes(4)
		if err != nil {
			return nil, err
		}
		b.R = color[0]
		b.G = color[1]
		b.B = color[2]
		z, err := r.readI16()
		b.ZPos = int(z)
		if err != nil {
			return nil, err
		}

		err = r.skip(2)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, b)
	}
	xoffsets := make([]int, 0, xSize)
	for i := 0; i < int(xSize); i++ {
		v, err := r.readI32()
		if err != nil {
			return nil, err
		}
		xoffsets = append(xoffsets, int(v))
	}
	offsets := make([]int, 0, xSize*ySize)
	for i := 0; i < int(xSize)*int(ySize); i++ {
		v, err := r.readI16()
		if err != nil {
			return nil, err
		}
		offsets = append(offsets, int(v))
	}

	return &KV6File{
		mgl.Vec3I{int(xSize), int(ySize), int(zSize)},
		blocks,
		xoffsets,
		offsets,
	}, nil
}
