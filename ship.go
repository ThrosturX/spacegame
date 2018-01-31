package spacegame

import "github.com/faiface/pixel"

type Ship struct {
	name       string
	imagePath  string
	angle      float64
	bounds pixel.Rect
}

func NewShip(name, imagePath string) *Ship {
	return &Ship{
		name:       name,
		imagePath:  imagePath,
		angle:      0.0,
		bounds: pixel.R(0, 0, 32, 32),
	}
}

func (ss *Ship) Name() string {
	return ss.name
}

func (ss *Ship) ImagePath() string {
	return ss.imagePath
}

func (ss *Ship) Angle() float64 {
	return ss.angle
}

func (ss *Ship) Bounds() pixel.Rect {
    return ss.bounds
}

func (ss *Ship) Turn(angle float64) {
	ss.angle += angle
}
