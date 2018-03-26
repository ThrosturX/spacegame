package spacegame

import (
	"encoding/json"
	"os"
)

type SolarSystem struct {
	name       string
	celestials []Celestial
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

func NewSolarSystem(name string, celestials ...Celestial) *SolarSystem {
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
	var celestials []Celestial

	for _, c := range config.Celestials {
		celestial := NewCelestial(c.Name, c.ImagePath, c.Coordinates)
		celestials = append(celestials, celestial)
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

func (s SolarSystem) Celestials() []Celestial {
	return s.celestials
}
