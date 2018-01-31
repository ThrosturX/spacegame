package spacegame

import (
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Renderer interface {
    Center() pixel.Vec
	Clear()
	Render(renderable Entity, position pixel.Vec)
	Update()
}

// TODO: Use a handle on the window from a controller
type PixelWindowRenderer struct {
	window          *pixelgl.Window
	resourceManager ResourceManager
}

func NewPixelWindowRenderer(window *pixelgl.Window, baseResourcePath string) *PixelWindowRenderer {
	return &PixelWindowRenderer{
		window:          window,
		resourceManager: NewStandardResourceManager(baseResourcePath),
	}
}

func (pwr *PixelWindowRenderer) Center() pixel.Vec {
    return pwr.window.Bounds().Center()
}

func (pwr *PixelWindowRenderer) Clear() {
	pwr.window.Clear(colornames.Black)
}

func (pwr *PixelWindowRenderer) Render(renderable Entity, position pixel.Vec) {
    resource := pwr.resourceManager.Resource(renderable)
	sprite := resource.sprite

    scale := resource.Scale()

	// TODO: remove magic scale constant Vec
	matrix := pixel.IM
	matrix = matrix.ScaledXY(pixel.ZV, pixel.V(scale, scale))
	matrix = matrix.Rotated(pixel.ZV, renderable.Angle())
	//	fmt.Printf("Rendering %s at angle %v", renderable.Name(), renderable.Angle())
	matrix = matrix.Moved(position)

	sprite.Draw(pwr.window, matrix)
}

func (pwr *PixelWindowRenderer) Update() {
	pwr.window.Update()
}
