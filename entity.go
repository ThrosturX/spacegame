package spacegame

import "github.com/faiface/pixel"

// TODO: Rotate(float64), replaces Turn for Ship etc.
// An Entity describes a
type Entity interface {
	Name() string
	Angle() float64
	Bounds() pixel.Rect
	Coordinates() pixel.Vec
	Velocity() pixel.Vec
	Translate(pixel.Vec)
}

type BaseEntity struct {
	name        string
	angle       float64
	coordinates pixel.Vec
	bounds      pixel.Rect
}

func NewBaseEntity(name string, rect pixel.Rect) *BaseEntity {
	entity := BaseEntity{
		name:        name,
		angle:       0.0,
		coordinates: pixel.ZV,
		bounds:      rect,
	}

	return &entity
}

func (be BaseEntity) Name() string {
	return be.name
}

func (be BaseEntity) Angle() float64 {
	return be.angle
}

func (be BaseEntity) Bounds() pixel.Rect {
	return be.bounds
}

func (be BaseEntity) Coordinates() pixel.Vec {
	return be.coordinates
}

func (be BaseEntity) Velocity() pixel.Vec {
	return pixel.ZV

}
func (be BaseEntity) Translate(vec pixel.Vec) {
	be.coordinates = be.coordinates.Add(vec)
}
