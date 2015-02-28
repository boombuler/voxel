package rendering

import (
	"github.com/boombuler/voxel/mgl"
	"github.com/go-gl/gl"
)

type plane struct {
	mgl.Vec3
	D float32
}

func GLProjection() mgl.Mat4 {
	var projM [16]float32
	gl.GetFloatv(gl.PROJECTION_MATRIX, projM[:])
	return mgl.Mat4(projM)
}

func GLModelView() mgl.Mat4 {
	var modM [16]float32
	gl.GetFloatv(gl.MODELVIEW_MATRIX, modM[:])
	return mgl.Mat4(modM)
}

func (p *plane) Assign(v mgl.Vec4) {
	p.Vec3 = v.Vec3()
	p.D = v.W()
	p.Normalize()
}

func (p *plane) Normalize() {
	mag := float32(1) / p.Vec3.Len()
	p.Vec3 = p.Vec3.Mul(mag)
	p.D = p.D * mag
}

type Frustum struct {
	planes [6]*plane
}

const (
	pRight  = 0
	pLeft   = 1
	pBottom = 2
	pTop    = 3
	pBack   = 4
	pFront  = 5

	ptA = 0
	ptB = 1
	ptC = 2
	ptD = 3
)

func NewFrustum() *Frustum {
	res := new(Frustum)
	for i := 0; i < 6; i++ {
		res.planes[i] = new(plane)
	}
	return res
}

func (f *Frustum) Update() {
	clip := GLProjection().MulMat4(GLModelView())
	// Die Seiten des Frustums aus der berechneten Clippingmatrix extrahieren
	r3 := clip.Row(3)
	f.planes[pLeft].Assign(r3.Add(clip.Row(0)))
	f.planes[pRight].Assign(r3.Sub(clip.Row(0)))

	f.planes[pBottom].Assign(r3.Add(clip.Row(1)))
	f.planes[pTop].Assign(r3.Sub(clip.Row(1)))

	f.planes[pFront].Assign(r3.Add(clip.Row(2)))
	f.planes[pBack].Assign(r3.Sub(clip.Row(2)))
}

func (f *Frustum) IsPointWithin(pt mgl.Vec3) bool {
	for _, pl := range f.planes {
		if (pl.X()*pt.X() + pl.Y()*pt.Y() + pl.Z()*pt.Z() + pl.D) <= 0 {
			return false
		}
	}
	return true
}

func (f *Frustum) IsSphereWithin(center mgl.Vec3, pRadius float32) bool {
	for _, pl := range f.planes {
		if (pl.X()*center.X() + pl.Y()*center.Y() + pl.Z()*center.Z() + pl.D) <= -pRadius {
			return false
		}
	}
	return true
}

func (f *Frustum) IsCubeWithin(pt mgl.Vec3, size mgl.Vec3) bool {
	dCenter := size.Mul(0.5)
	center := pt.Add(dCenter)
	pRadius := dCenter.Len()

	for _, pl := range f.planes {
		dPl := (pl.X()*center.X() + pl.Y()*center.Y() + pl.Z()*center.Z() + pl.D)

		if dPl <= -pRadius {
			return false
		}
	}
	return true
}
