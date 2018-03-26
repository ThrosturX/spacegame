package spacegame

import "github.com/faiface/pixel"

type Celestial interface {
	Entity
	ImagePath() string
	Radius() float64
    Land(ship *Ship) bool // Returns true if the ship was successfully docked with the celestial
}

// A Celestial is an Entity that can be selected and potentially allows ships to land/dock
// Celestials reside in Solar Systems
type BaseCelestial struct {
	name      string
	imagePath string
	bounds    pixel.Rect
	position  pixel.Vec
	radius    float64
}

func NewCelestial(name, imagePath string, position pixel.Vec) Celestial {
    return DockableCelestial{ // TODO: Return uninhabitable celestials by default, and wrap this
		BaseCelestial{
            name:      name,
            imagePath: imagePath,
            position:  position,
            radius:    64.0, // TODO: Get from imagePath? makes docking adaptive
        },
	}
}

func (c BaseCelestial) Name() string {
	return c.name
}

func (c BaseCelestial) ImagePath() string {
	return c.imagePath
}

func (c BaseCelestial) Angle() float64 {
	return 0.0
}

func (c BaseCelestial) Bounds() pixel.Rect {
	// center the celestial around its origin
	return pixel.R(-c.radius, -c.radius, c.radius, c.radius)
}

func (c BaseCelestial) Coordinates() pixel.Vec {
	return c.position
}

func (c BaseCelestial) Radius() float64 {
	return c.radius
}

func (c BaseCelestial) Velocity() pixel.Vec {
	return pixel.ZV
}

func (c BaseCelestial) Translate(by pixel.Vec) {
	c.position = c.position.Add(by)
}

type DockableCelestial struct {
    BaseCelestial
}

func (dc DockableCelestial) Land(ship *Ship) bool {
    // If the player is in range
    // TODO: Notify the game engine that the ship has landed so that a LandedScene can be rendered.

    // No docking available until the LandedScene has been implemented

    // Regardless, the player should receive "Docking granted" or "Docking requested" messages.

    // The player is not in range
    return false
}

type CelestialCollection []Celestial

type CelestialConfig struct {
	Name      string
	Coordinates  pixel.Vec
	ImagePath string
	Radius    float64
}

func (cs CelestialCollection) Config() []CelestialConfig {
	var collection []CelestialConfig

	for _, c := range cs {
		config := CelestialConfig{
			Name:      c.Name(),
			ImagePath: c.ImagePath(),
			Coordinates:  c.Coordinates(),
			Radius:    c.Radius(),
		}
		collection = append(collection, config)
	}

	return collection
}
