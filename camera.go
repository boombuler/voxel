package main

import (
	"math"

	"github.com/boombuler/voxel/mgl"
	"github.com/go-gl-legacy/gl"
	"github.com/go-gl-legacy/glu"
	"github.com/go-gl/glfw/v3.0/glfw"
)

type camera struct {
	pos        mgl.Vec3
	speed      float32
	mouseSpeed float64
	hAngle     float64
	vAngle     float64
	mouseX     float64
	mouseY     float64
}

func NewCamera() *camera {
	return &camera{
		pos:        mgl.Vec3{0, 0, -8},
		speed:      1.5,
		mouseSpeed: 1.1,
	}
}

func (c *camera) update(w *glfw.Window, dt float64) bool {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	glu.Perspective(45, 4.0/3.0, 0.1, 100.0)

	if c.mouseX == 0 && c.mouseY == 0 {
		c.mouseX, c.mouseY = w.GetCursorPosition()
	}
	x, y := w.GetCursorPosition()
	dx := c.mouseX - x
	dy := c.mouseY - y
	result := dx != 0 || dy != 0
	c.mouseX, c.mouseY = x, y
	c.hAngle += c.mouseSpeed * dt * dx
	c.vAngle += c.mouseSpeed * dt * dy

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	direction := mgl.Vec3{
		float32(math.Cos(c.vAngle) * math.Sin(c.hAngle)),
		float32(math.Sin(c.vAngle)),
		float32(math.Cos(c.vAngle) * math.Cos(c.hAngle)),
	}
	right := mgl.Vec3{
		float32(math.Sin(c.hAngle - math.Pi/2.0)),
		0,
		float32(math.Cos(c.hAngle - math.Pi/2.0)),
	}
	up := direction.Cross(right)
	center := direction.Add(c.pos)
	glu.LookAt(
		float64(c.pos.X()), float64(c.pos.Y()), float64(c.pos.Z()),
		float64(center.X()), float64(center.Y()), float64(center.Z()),
		float64(up.X()), float64(up.Y()), float64(up.Z()))

	dTime := float32(dt)

	if w.GetKey(glfw.KeyW) == glfw.Press {
		c.pos = c.pos.Add(direction.Mul(dTime * c.speed))
		result = true
	}
	if w.GetKey(glfw.KeyS) == glfw.Press {
		c.pos = c.pos.Add(direction.Mul(dTime * -c.speed))
		result = true
	}
	if w.GetKey(glfw.KeyA) == glfw.Press {
		c.pos = c.pos.Add(right.Mul(dTime * c.speed))
		result = true
	}
	if w.GetKey(glfw.KeyD) == glfw.Press {
		c.pos = c.pos.Add(right.Mul(dTime * -c.speed))
		result = true
	}
	return result
}
