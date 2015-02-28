package rendering

import (
	"image/color"
	"math"
)

type Color struct {
	Red   float32
	Green float32
	Blue  float32
	Alpha float32
}

const color_Size int = 4 * 4

func fToUiColVal(f float32) uint32 {
	// ensure it is in range
	ff := math.Max(float64(0), math.Min(float64(1), float64(f)))
	return uint32(math.Trunc(ff * float64(math.MaxUint16)))
}

func uiToFColVal(f uint32) float32 {
	ff := math.Min(float64(math.MaxUint16), float64(f))
	return float32(ff) / float32(math.MaxUint16)
}

func (c *Color) RGBA() (r, g, b, a uint32) {
	if c == nil {
		return 0, 0, 0, 0
	}

	return fToUiColVal(c.Red),
		fToUiColVal(c.Green),
		fToUiColVal(c.Blue),
		fToUiColVal(c.Alpha)
}

var ColorModel = color.ModelFunc(glColorModel)

func glColorModel(c color.Color) color.Color {
	if c == nil {
		return nil
	}
	if _, ok := c.(*Color); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	if a == 0 {
		return nil
	}
	return &Color{
		uiToFColVal(r),
		uiToFColVal(g),
		uiToFColVal(b),
		uiToFColVal(a),
	}
}
