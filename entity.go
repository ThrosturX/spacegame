package spacegame

import "github.com/faiface/pixel"

// TODO: Rotate(float64), replaces Turn for Ship etc.
type Entity interface {
	Name() string
	Angle() float64
	Bounds() pixel.Rect
	Coordinates() pixel.Vec
	Velocity() pixel.Vec
	Translate(pixel.Vec)
}

type BasicEntity struct {
	name        string
	angle       float64
	coordinates pixel.Vec
	bounds      pixel.Rect
}

func NewBasicEntity(name string, rect pixel.Rect) *BasicEntity {
	entity := BasicEntity{
		name:        name,
		angle:       0.0,
		coordinates: pixel.ZV,
		bounds:      rect,
	}

	return &entity
}

func (be BasicEntity) Name() string {
	return be.name
}

func (be BasicEntity) Angle() float64 {
	return be.angle
}

func (be BasicEntity) Bounds() pixel.Rect {
	return be.bounds
}

func (be BasicEntity) Coordinates() pixel.Vec {
	return be.coordinates
}

func (be BasicEntity) Velocity() pixel.Vec {
	return pixel.ZV

}
func (be BasicEntity) Translate(vec pixel.Vec) {
	be.coordinates = be.coordinates.Add(vec)
}
