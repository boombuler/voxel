package main

import (
	"fmt"
	"github.com/boombuler/voxel/rendering"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"runtime"
)

type Options struct {
	GlfwErrorCallback func(err glfw.ErrorCode, desc string)
	WindowTitle       string
	WindowWidth       int
	WindowHeight      int

	LoadFunc   func(e *Engine)
	UpdateFunc func(dt float64, e *Engine)
}

type Engine struct {
	RenderObjects []rendering.Object
}

func (e *Engine) renderObjects(fr *rendering.Frustum) {
	visibleObjects := make(chan rendering.Object)
	scaleF := float32(0.01)
	go func() {
		for _, obj := range e.RenderObjects {
			renderer := obj.Renderer()
			if renderer != nil && fr.IsCubeWithin(obj.Position(), obj.Size().Mul(scaleF)) {
				visibleObjects <- obj
			}
		}
		close(visibleObjects)
	}()
	for obj := range visibleObjects {
		gl.PushMatrix()
		p := obj.Position()
		gl.Translatef(p.X(), p.Y(), p.Z())
		gl.Scalef(scaleF, scaleF, scaleF)
		r := obj.Renderer()
		r.Render()
		gl.PopMatrix()
	}
}

func fallBackErrorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func StartEngine(options Options) {
	n := runtime.NumCPU()
	if n > 1 {
		runtime.GOMAXPROCS(n)
		runtime.LockOSThread()
	}
	if options.GlfwErrorCallback != nil {
		glfw.SetErrorCallback(options.GlfwErrorCallback)
	} else {
		glfw.SetErrorCallback(fallBackErrorCallback)
	}
	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	wnd, err := glfw.CreateWindow(options.WindowWidth, options.WindowHeight, options.WindowTitle, nil, nil)
	if err != nil {
		panic(err)
	}

	wnd.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	cam := NewCamera()

	gl.Viewport(0, 0, options.WindowWidth, options.WindowHeight)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	wnd.SetInputMode(glfw.Cursor, glfw.CursorDisabled)
	frustum := rendering.NewFrustum()
	engine := &Engine{
		RenderObjects: make([]rendering.Object, 0),
	}
	if options.LoadFunc != nil {
		options.LoadFunc(engine)
	}

	curTime := glfw.GetTime()
	for !wnd.ShouldClose() {
		nTime := glfw.GetTime()
		dt := nTime - curTime
		curTime = nTime
		if cam.update(wnd, dt) {
			frustum.Update()
		}
		options.UpdateFunc(dt, engine)
		// Draw Scene
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		engine.renderObjects(frustum)

		// Finish
		wnd.SwapBuffers()
		glfw.PollEvents()
	}
}
