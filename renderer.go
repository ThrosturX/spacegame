package spacegame

import (
	"fmt"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

// This may be a prototype for a "subwindow" interface...
type Renderer interface {
	Bounds() pixel.Rect
	Center() pixel.Vec
	Clear()
	Render(renderable Entity, position pixel.Vec)
	ResourceManager() ResourceManager
	Text(txt string, position pixel.Vec) error
	Update()
}

type PixelWindowRenderer struct {
	window          *pixelgl.Window
	resourceManager ResourceManager
	atlas           *text.Atlas
}

func NewPixelWindowRenderer(window *pixelgl.Window, resourceManager ResourceManager) *PixelWindowRenderer {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	return &PixelWindowRenderer{
		window:          window,
		resourceManager: resourceManager,
		atlas:           atlas,
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
	unitScalerX, unitScalerY := 1/bounds.W(), 1/bounds.H()
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

func (pwr *PixelWindowRenderer) Text(s string, position pixel.Vec) error {
	txt := text.New(position, pwr.atlas)
	_, err := fmt.Fprintln(txt, s)
	if err != nil {
		return err
	}

	txt.Draw(pwr.window, pixel.IM)
    return nil
}

func (pwr *PixelWindowRenderer) Update() {
	pwr.window.Update()
}
