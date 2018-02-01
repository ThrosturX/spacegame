package spacegame

import (
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// This may be a prototype for a "subwindow" interface...
type Renderer interface {
    Bounds() pixel.Rect
	Center() pixel.Vec
	Clear()
	Render(renderable Entity, position pixel.Vec)
    ResourceManager() ResourceManager
	Update()
}

type PixelWindowRenderer struct {
	window          *pixelgl.Window
	resourceManager ResourceManager
}

func NewPixelWindowRenderer(window *pixelgl.Window, resourceManager ResourceManager) *PixelWindowRenderer {
	return &PixelWindowRenderer{
		window:          window,
		resourceManager: resourceManager,
	}
}

func (pwr *PixelWindowRenderer) Bounds() pixel.Rect {
     return pwr.window.Bounds()
}

func (pwr *PixelWindowRenderer) Center() pixel.Vec {
	return pwr.window.Bounds().Center()
}

func (pwr *PixelWindowRenderer) Clear() {
	pwr.window.Clear(colornames.Black)
}

func (pwr *PixelWindowRenderer) Render(renderable Entity, position pixel.Vec) {
	// return early if position is out of bounds
	if pwr.window.Bounds().Intersect(renderable.Bounds().Moved(position)).Area() == 0 {
		return
	}

	resource := pwr.resourceManager.Resource(renderable)
	sprite := resource.sprite

    bounds := resource.Bounds()
    unitScalerX, unitScalerY := 1 / bounds.W(), 1 / bounds.H()
    rect := resource.Entity().Bounds()

	matrix := pixel.IM
	matrix = matrix.ScaledXY(pixel.ZV, pixel.V(unitScalerX, unitScalerY))
	matrix = matrix.ScaledXY(pixel.ZV, pixel.V(rect.W(), rect.H()))
	matrix = matrix.Rotated(pixel.ZV, renderable.Angle())
	matrix = matrix.Moved(position)

	sprite.Draw(pwr.window, matrix)
}
 
func (pwr *PixelWindowRenderer) ResourceManager() ResourceManager {
    return pwr.resourceManager
}

func (pwr *PixelWindowRenderer) Update() {
	pwr.window.Update()
}
