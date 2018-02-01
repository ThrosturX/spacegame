package spacegame

import "github.com/faiface/pixel"

// A graph of solar systems
type Universe struct {
	systems     map[string]*SolarSystem
	coordinates map[string]pixel.Vec
}

func NewUniverse(rm ResourceManager) *Universe {
	systems := make(map[string]*SolarSystem)
	coordinates := make(map[string]pixel.Vec)

    // TODO: This section comes from somewhere else
    name := "Vera"
	veraPlanet := NewCelestial(name, "images/planets/planet27.png", pixel.V(30, 60))
//    rm.CreateResource(veraPlanet, "planets/planet27.png")

	veraSystem := NewSolarSystem(name, veraPlanet)

    err := veraSystem.SaveToFile("resources/universe/systems/Vera.json")
    if err != nil { panic(err) }
    systems[name] = veraSystem
	coordinates[name] = pixel.V(0, 0)

	universe := Universe{
		systems:     systems,
		coordinates: coordinates,
	}

	return &universe
}

func (uv *Universe) Systems() map[string]*SolarSystem {
	return uv.systems
}

func (uv *Universe) SystemCoordinates() map[string]pixel.Vec {
	return uv.coordinates
}

func (uv *Universe) addSystem(system *SolarSystem, coordinates pixel.Vec) {
	name := system.Name()
	uv.systems[name] = system
	uv.coordinates[name] = coordinates
}
