package spacegame

import "github.com/faiface/pixel"

type Entity interface {
	Name() string
	ImagePath() string
	Angle() float64
	Bounds() pixel.Rect
}
