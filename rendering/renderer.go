package rendering

type Renderer interface {
	Render()
}

type RenderCloser interface {
	Renderer
	Close()
}

type RenderFunc func()

func (rf RenderFunc) Render() {
	rf()
}

type Options int

const (
	NONE       Options = 0
	NO_CULLING Options = 1 << iota
	NO_MESHING
	NO_VBO
)

func (o Options) HasFlag(opt Options) bool {
	return o&opt == opt
}
