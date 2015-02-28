package magica

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/boombuler/voxel/mgl"
	"github.com/boombuler/voxel/rendering"
	"image/color"
	"io"
)

type colorVoxel color.RGBA

func (cv colorVoxel) Color() color.Color {
	return color.RGBA(cv)
}

func (cv colorVoxel) Equals(v rendering.Voxel) bool {
	cvo, ok := v.(colorVoxel)
	if ok {
		return cv.R == cvo.R && cv.G == cvo.G && cv.B == cvo.B && cv.A == cvo.A
	}
	return false
}

type VoxFileModel struct {
	palette []colorVoxel
	size    mgl.Vec3I
	content map[mgl.Vec3I]byte
}

var defaultPalette []colorVoxel = []colorVoxel{defColor(32767), defColor(25599), defColor(19455), defColor(13311), defColor(7167), defColor(1023), defColor(32543), defColor(25375), defColor(19231), defColor(13087), defColor(6943), defColor(799), defColor(32351), defColor(25183),
	defColor(19039), defColor(12895), defColor(6751), defColor(607), defColor(32159), defColor(24991), defColor(18847), defColor(12703), defColor(6559), defColor(415), defColor(31967), defColor(24799), defColor(18655), defColor(12511), defColor(6367), defColor(223), defColor(31775), defColor(24607), defColor(18463), defColor(12319), defColor(6175), defColor(31),
	defColor(32760), defColor(25592), defColor(19448), defColor(13304), defColor(7160), defColor(1016), defColor(32536), defColor(25368), defColor(19224), defColor(13080), defColor(6936), defColor(792), defColor(32344), defColor(25176), defColor(19032), defColor(12888), defColor(6744), defColor(600), defColor(32152), defColor(24984), defColor(18840),
	defColor(12696), defColor(6552), defColor(408), defColor(31960), defColor(24792), defColor(18648), defColor(12504), defColor(6360), defColor(216), defColor(31768), defColor(24600), defColor(18456), defColor(12312), defColor(6168), defColor(24), defColor(32754), defColor(25586), defColor(19442), defColor(13298), defColor(7154), defColor(1010), defColor(32530),
	defColor(25362), defColor(19218), defColor(13074), defColor(6930), defColor(786), defColor(32338), defColor(25170), defColor(19026), defColor(12882), defColor(6738), defColor(594), defColor(32146), defColor(24978), defColor(18834), defColor(12690), defColor(6546), defColor(402), defColor(31954), defColor(24786), defColor(18642), defColor(12498), defColor(6354),
	defColor(210), defColor(31762), defColor(24594), defColor(18450), defColor(12306), defColor(6162), defColor(18), defColor(32748), defColor(25580), defColor(19436), defColor(13292), defColor(7148), defColor(1004), defColor(32524), defColor(25356), defColor(19212), defColor(13068), defColor(6924), defColor(780), defColor(32332), defColor(25164), defColor(19020),
	defColor(12876), defColor(6732), defColor(588), defColor(32140), defColor(24972), defColor(18828), defColor(12684), defColor(6540), defColor(396), defColor(31948), defColor(24780), defColor(18636), defColor(12492), defColor(6348), defColor(204), defColor(31756), defColor(24588), defColor(18444), defColor(12300), defColor(6156), defColor(12), defColor(32742),
	defColor(25574), defColor(19430), defColor(13286), defColor(7142), defColor(998), defColor(32518), defColor(25350), defColor(19206), defColor(13062), defColor(6918), defColor(774), defColor(32326), defColor(25158), defColor(19014), defColor(12870), defColor(6726), defColor(582), defColor(32134), defColor(24966), defColor(18822), defColor(12678), defColor(6534),
	defColor(390), defColor(31942), defColor(24774), defColor(18630), defColor(12486), defColor(6342), defColor(198), defColor(31750), defColor(24582), defColor(18438), defColor(12294), defColor(6150), defColor(6), defColor(32736), defColor(25568), defColor(19424), defColor(13280), defColor(7136), defColor(992), defColor(32512), defColor(25344), defColor(19200),
	defColor(13056), defColor(6912), defColor(768), defColor(32320), defColor(25152), defColor(19008), defColor(12864), defColor(6720), defColor(576), defColor(32128), defColor(24960), defColor(18816), defColor(12672), defColor(6528), defColor(384), defColor(31936), defColor(24768), defColor(18624), defColor(12480), defColor(6336), defColor(192), defColor(31744),
	defColor(24576), defColor(18432), defColor(12288), defColor(6144), defColor(28), defColor(26), defColor(22), defColor(20), defColor(16), defColor(14), defColor(10), defColor(8), defColor(4), defColor(2), defColor(896), defColor(832), defColor(704), defColor(640), defColor(512), defColor(448), defColor(320), defColor(256), defColor(128), defColor(64), defColor(28672), defColor(26624), defColor(22528), defColor(20480),
	defColor(16384), defColor(14336), defColor(10240), defColor(8192), defColor(4096), defColor(2048), defColor(29596), defColor(27482), defColor(23254), defColor(21140), defColor(16912), defColor(14798), defColor(10570), defColor(8456), defColor(4228), defColor(2114)}

func defColor(val uint16) colorVoxel {
	d := uint32(val)
	pmax := uint32(0x1f)
	b := uint8((d&pmax)<<3 | 0x07)
	g := uint8(((d>>5)&pmax)<<3 | 0x07)
	r := uint8(((d>>10)&pmax)<<3 | 0x07)
	return colorVoxel{r, g, b, 0xff}
}

func (vfm *VoxFileModel) String() string {
	return fmt.Sprintf("VOX File Model (%vx%vx%v)", vfm.size.X, vfm.size.Y, vfm.size.Z)
}

func (vfm *VoxFileModel) ColorModel() color.Model {
	return color.RGBAModel
}

func (vfm *VoxFileModel) Size() mgl.Vec3I {
	return vfm.size
}

func (vfm *VoxFileModel) ForeachVoxel(fn func(pos mgl.Vec3I, vox rendering.Voxel)) {
	for k, v := range vfm.content {
		if v > 0 {
			fn(k, vfm.palette[v-1])
		}
	}
}

func (vfm *VoxFileModel) At(pos mgl.Vec3I) rendering.Voxel {
	idx, ok := vfm.content[pos]
	if ok || idx > 0 {
		return vfm.palette[idx-1]
	}
	return nil

}

func (vfm *VoxFileModel) setRawVoxels(rv []rawVoxel) error {
	if rv == nil {
		return errors.New("no voxel data")
	}
	orgSize := vfm.Size()
	vfm.size = mgl.Vec3I{
		orgSize.X(),
		orgSize.Z(),
		orgSize.Y(),
	}
	rotateCoords := func(pos mgl.Vec3I) mgl.Vec3I {
		return mgl.Vec3I{
			pos.X(),
			orgSize.Z() - pos.Z() - 1,
			pos.Y(),
		}
	}

	vfm.content = make(map[mgl.Vec3I]byte)
	for _, v := range rv {
		p := rotateCoords(v.Vec3I)
		if p.X() >= vfm.size.X() || p.Y() >= vfm.size.Y() || p.Z() >= vfm.size.Z() {
			return fmt.Errorf("voxel position out of range: %v size: %v orgSize:", p, vfm.size, orgSize)
		}

		vfm.content[p] = v.idx
	}
	return nil
}

const (
	head_file       = "VOX "
	chunk_main      = "MAIN"
	chunk_palette   = "RGBA"
	chunk_voxels    = "XYZI"
	chunk_size      = "SIZE"
	current_version = 150
)

type voxReader struct {
	*bufio.Reader
}

func errorOrTxt(err error, txt string) error {
	if err == nil {
		return errors.New(txt)
	}
	return err
}

func (vr *voxReader) readBytes(cnt int) ([]byte, error) {
	buf := make([]byte, cnt)
	n, err := vr.Read(buf)
	if n != cnt {
		return nil, errors.New("unexpected end of file")
	}
	return buf, err
}

func (vr *voxReader) readStr4() (string, error) {
	txt, err := vr.readBytes(4)
	if err != nil {
		return "", err
	}
	return string(txt), err
}

func (vr *voxReader) skip(size int) error {
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

func (vr *voxReader) readInt() (int, error) {
	b, err := vr.readBytes(4)
	if err != nil {
		return 0, err
	}
	return int(int32(b[0]) | (int32(b[1]) << 8) | (int32(b[2]) << 16) | (int32(b[3]) << 24)), nil
}

func readPalette(vr *voxReader, contentSize, childChunkSize int) ([]colorVoxel, error) {
	if contentSize != (256 * 4) {
		return nil, fmt.Errorf("invalid palette size: %v", contentSize)
	}
	if childChunkSize != 0 {
		return nil, errors.New("unexpected child chunks in palette data")
	}
	p := make([]colorVoxel, 0, 255)
	for i := 0; i < 256; i++ {
		v, err := vr.readBytes(4)
		if err != nil {
			return nil, err
		}
		p = append(p, colorVoxel{v[0], v[1], v[2], v[3]})
	}

	return p, nil
}

var vec3IZero = mgl.Vec3I{0, 0, 0}

func readSize(vr *voxReader, contentSize, childChunkSize int) (mgl.Vec3I, error) {
	if contentSize != 12 || childChunkSize != 0 {
		return vec3IZero, errors.New("invalid child chunk")
	}
	xx, err := vr.readInt()
	if err != nil {
		return vec3IZero, err
	}
	yy, err := vr.readInt()
	if err != nil {
		return vec3IZero, err
	}
	zz, err := vr.readInt()
	if err != nil {
		return vec3IZero, err
	}
	if xx > 256 || yy > 256 || zz > 256 {
		return vec3IZero, fmt.Errorf("invalid size: %vx%vx%v", xx, yy, zz)
	}
	return mgl.Vec3I{xx, yy, zz}, nil
}

type rawVoxel struct {
	mgl.Vec3I
	idx byte
}

func readVoxels(vr *voxReader, contentSize, childChunkSize int) ([]rawVoxel, error) {
	if childChunkSize > 0 {
		return nil, errors.New("unexpected child chunk in voxels")
	}
	cnt, err := vr.readInt()
	if err != nil {
		return nil, err
	}
	if (cnt*4) != (contentSize-4) || cnt > 256*256*256 {
		return nil, errors.New("voxel count missmatches data size")
	}
	result := make([]rawVoxel, 0, cnt)
	for i := 0; i < cnt; i++ {
		d, err := vr.readBytes(4)
		if err != nil {
			return nil, err
		}
		result = append(result, rawVoxel{mgl.Vec3I{int(d[0]), int(d[1]), int(d[2])}, d[3]})
	}
	return result, nil
}

func Read(rd io.Reader) (*VoxFileModel, error) {
	vr := &voxReader{bufio.NewReader(rd)}
	if head, err := vr.readStr4(); err != nil || head != head_file {
		return nil, errorOrTxt(err, "invalid file format")
	}
	if ver, err := vr.readInt(); err != nil || ver > current_version { // read version number
		return nil, errorOrTxt(err, "unsupported vox file version")
	}
	if cn, err := vr.readStr4(); err != nil || cn != chunk_main {
		return nil, errorOrTxt(err, "invalid file: expected main chunk")
	}
	contentSize, err := vr.readInt()
	if err != nil {
		return nil, err
	}
	totalChunkSize, err := vr.readInt()
	if err != nil {
		return nil, err
	}
	if contentSize > 0 {
		if err := vr.skip(contentSize); err != nil {
			return nil, err
		}
	}
	result := new(VoxFileModel)
	var voxels []rawVoxel
	for totalChunkSize > 0 {
		chunkName, err := vr.readStr4()

		contentSize, err := vr.readInt()
		if err != nil {
			return nil, err
		}
		childChunkSize, err := vr.readInt()
		if err != nil {
			return nil, err
		}

		totalChunkSize -= 12 + contentSize + childChunkSize

		if err != nil {
			return nil, err
		}
		switch chunkName {
		case chunk_size:
			result.size, err = readSize(vr, contentSize, childChunkSize)
		case chunk_palette:
			result.palette, err = readPalette(vr, contentSize, childChunkSize)
		case chunk_voxels:
			voxels, err = readVoxels(vr, contentSize, childChunkSize)
		default:
			// skip unknown chunk
			err = vr.skip(contentSize + childChunkSize)
		}

		if err != nil {
			return nil, err
		}
	}
	if result.palette == nil {
		result.palette = defaultPalette
	}
	return result, result.setRawVoxels(voxels)
}
