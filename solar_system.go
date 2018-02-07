package spacegame

import (
	"encoding/json"
	"os"

	"github.com/faiface/pixel"
)

// A Celestial is an Entity that can be selected
type Celestial struct {
	name      string
	imagePath string
	bounds    pixel.Rect
	position  pixel.Vec
	radius    float64
}

func NewCelestial(name, imagePath string, position pixel.Vec) Celestial {
	return Celestial{
		name:      name,
		imagePath: imagePath,
		position:  position,
		radius:    64.0,
	}
}

func (c Celestial) Name() string {
	return c.name
}

func (c Celestial) ImagePath() string {
	return c.imagePath
}

func (c Celestial) Angle() float64 {
	return 0.0
}

func (c Celestial) Bounds() pixel.Rect {
	return pixel.R(0.0, 0.0, 2*c.radius, 2*c.radius)
}

func (c Celestial) Coordinates() pixel.Vec {
	return c.position
}

func (c Celestial) Velocity() pixel.Vec {
	return pixel.ZV
}

func (c Celestial) Translate(by pixel.Vec) {
	c.position = c.position.Add(by)
}

type CelestialCollection []*Celestial

type CelestialConfig struct {
	Name      string
	Position  pixel.Vec
	ImagePath string
	Radius    float64
}

func (cs CelestialCollection) Config() []CelestialConfig {
	var collection []CelestialConfig

	for _, c := range cs {
		config := CelestialConfig{
			Name:      c.name,
			ImagePath: c.imagePath,
			Position:  c.position,
			Radius:    c.radius,
		}
		collection = append(collection, config)
	}

	return collection
}

type SolarSystem struct {
	name       string
	celestials []*Celestial
}

type SolarSystemConfig struct {
	Name       string
	Celestials []CelestialConfig
}

func (s SolarSystem) Config() SolarSystemConfig {
	return SolarSystemConfig{
		Name:       s.name,
		Celestials: CelestialCollection(s.celestials).Config(),
	}
}

func NewSolarSystem(name string, celestials ...*Celestial) *SolarSystem {
	s := SolarSystem{
		name:       name,
		celestials: celestials,
	}
	return &s
}

func LoadSystem(path string) (*SolarSystem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var solarSystemConfig SolarSystemConfig
	decoder := json.NewDecoder(file)
	decoder.Decode(&solarSystemConfig)
	return solarSystemConfig.load(), nil
}

func (config SolarSystemConfig) load() *SolarSystem {
	var celestials []*Celestial

	for _, c := range config.Celestials {
		celestial := NewCelestial(c.Name, c.ImagePath, c.Position)
		celestials = append(celestials, &celestial)
	}

	system := &SolarSystem{
		name:       config.Name,
		celestials: celestials,
	}
	return system
}

func (s SolarSystem) SaveToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.Encode(s.Config())

	return nil
}

func (s SolarSystem) Name() string {
	return s.name
}

func (s SolarSystem) Celestials() []*Celestial {
	return s.celestials
}
