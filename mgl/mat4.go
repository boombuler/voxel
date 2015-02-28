package mgl


import (
	"bytes"
	"fmt"
	"math"
	"text/tabwriter"
)

type Mat4 [16]float32

func Identity() Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Pretty prints the matrix
func (m Mat4) String() string {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 4, 4, 1, ' ', tabwriter.AlignRight)
	for i := 0; i < 4; i++ {
		for _, col := range m.Row(i) {
			fmt.Fprintf(w, "%f\t", col)
		}

		fmt.Fprintln(w, "")
	}
	w.Flush()

	return buf.String()
}

func (m Mat4) Row(row int) Vec4 {
	return Vec4{m[row+0], m[row+4], m[row+8], m[row+12]}
}
func (m Mat4) Rows() (row0, row1, row2, row3 Vec4) {
	return m.Row(0), m.Row(1), m.Row(2), m.Row(3)
}
func (m Mat4) Col(col int) Vec4 {
	return Vec4{m[col*4+0], m[col*4+1], m[col*4+2], m[col*4+3]}
}
func (m Mat4) Cols() (col0, col1, col2, col3 Vec4) {
	return m.Col(0), m.Col(1), m.Col(2), m.Col(3)
}

func (m1 Mat4) MulMat4(m2 Mat4) Mat4 {
	return Mat4{
		m1[0]*m2[0] + m1[4]*m2[1] + m1[8]*m2[2] + m1[12]*m2[3],
		m1[1]*m2[0] + m1[5]*m2[1] + m1[9]*m2[2] + m1[13]*m2[3],
		m1[2]*m2[0] + m1[6]*m2[1] + m1[10]*m2[2] + m1[14]*m2[3],
		m1[3]*m2[0] + m1[7]*m2[1] + m1[11]*m2[2] + m1[15]*m2[3],
		m1[0]*m2[4] + m1[4]*m2[5] + m1[8]*m2[6] + m1[12]*m2[7],
		m1[1]*m2[4] + m1[5]*m2[5] + m1[9]*m2[6] + m1[13]*m2[7],
		m1[2]*m2[4] + m1[6]*m2[5] + m1[10]*m2[6] + m1[14]*m2[7],
		m1[3]*m2[4] + m1[7]*m2[5] + m1[11]*m2[6] + m1[15]*m2[7],
		m1[0]*m2[8] + m1[4]*m2[9] + m1[8]*m2[10] + m1[12]*m2[11],
		m1[1]*m2[8] + m1[5]*m2[9] + m1[9]*m2[10] + m1[13]*m2[11],
		m1[2]*m2[8] + m1[6]*m2[9] + m1[10]*m2[10] + m1[14]*m2[11],
		m1[3]*m2[8] + m1[7]*m2[9] + m1[11]*m2[10] + m1[15]*m2[11],
		m1[0]*m2[12] + m1[4]*m2[13] + m1[8]*m2[14] + m1[12]*m2[15],
		m1[1]*m2[12] + m1[5]*m2[13] + m1[9]*m2[14] + m1[13]*m2[15],
		m1[2]*m2[12] + m1[6]*m2[13] + m1[10]*m2[14] + m1[14]*m2[15],
		m1[3]*m2[12] + m1[7]*m2[13] + m1[11]*m2[14] + m1[15]*m2[15]}
}
func (mat Mat4) MulVec4(vec Vec4) Vec4 {
	return Vec4{
		mat[0]*vec[0] + mat[4]*vec[1] + mat[8]*vec[2] + mat[12]*vec[3],
		mat[1]*vec[0] + mat[5]*vec[1] + mat[9]*vec[2] + mat[13]*vec[3],
		mat[2]*vec[0] + mat[6]*vec[1] + mat[10]*vec[2] + mat[14]*vec[3],
		mat[3]*vec[0] + mat[7]*vec[1] + mat[11]*vec[2] + mat[15]*vec[3]}
}
func (m Mat4) Translate(x, y, z float32) Mat4 {
	return m.MulMat4(Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, z, 1,
	})
}
func (m Mat4) TranslateVec3(v Vec3) Mat4 {
	return m.MulMat4(Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		v[0], v[1], v[2], 1,
	})
}
func (m Mat4) Scale(x, y, z float32) Mat4 {
	return m.MulMat4(Mat4{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	})
}

func (m Mat4) Rotate(angle Radian, axis Vec3) Mat4 {
	a := axis.Normalize()

	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))
	d := 1 - c

	return m.MulMat4(Mat4{
		c + d*a[0]*a[0],
		0 + d*a[1]*a[0] - s*a[2],
		0 + d*a[2]*a[0] + s*a[1],
		0,

		0 + d*a[0]*a[1] + s*a[2],
		c + d*a[1]*a[1],
		0 + d*a[2]*a[1] - s*a[0],
		0,

		0 + d*a[0]*a[1] - s*a[1],
		0 + d*a[1]*a[2] + s*a[0],
		c + d*a[2]*a[2],
		0,

		0, 0, 0, 1,
	})
}
