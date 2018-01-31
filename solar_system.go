package spacegame

import "github.com/faiface/pixel"

// A Celestial is an Entity that can be selected
type Celestial struct {
	name       string
	imagePath  string
	bounds pixel.Rect
	position   pixel.Vec
}

func NewCelestial(name, imagePath string, position pixel.Vec) *Celestial {
	return &Celestial{
		name:      name,
		imagePath: imagePath,
		position:  position,
	}
}

func (cc Celestial) Name() string {
	return cc.name
}

func (cc Celestial) ImagePath() string {
	return cc.imagePath
}

func (cc Celestial) Angle() float64 {
	return 0.0
}

func (cc Celestial) Bounds() pixel.Rect {
	return pixel.R(0.0, 0.0, 128.0, 128.0)
}

type SolarSystem struct {
	name       string
	celestials []*Celestial
}

func NewSolarSystem(name string, celestials ...*Celestial) *SolarSystem {
	ss := SolarSystem{
		name:       name,
		celestials: celestials,
	}
	return &ss
}

func (ss SolarSystem) Name() string {
	return ss.name
}

func (ss SolarSystem) Celestials() []*Celestial {
	return ss.celestials
}
