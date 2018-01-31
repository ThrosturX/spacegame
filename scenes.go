package spacegame

import "github.com/faiface/pixel"

type Scene interface {
	Render(Renderer)
}

// TODO: Extract Ship into Player...
type SpaceScene struct {
	system         *SolarSystem
	playerShip     *Ship
	playerLocation pixel.Vec
}

func NewSpaceScene(system *SolarSystem, player *Player, playerLocation pixel.Vec) *SpaceScene {
	return &SpaceScene{
		system:         system,
		playerShip:     player.Ship(),
		playerLocation: playerLocation,
	}
}

func (ss *SpaceScene) Render(renderer Renderer) {
	// Render the background
	renderer.Clear()
	// TODO: Starscape

	// Render any planets in this scene
	for _, celestial := range ss.system.celestials {
		renderer.Render(celestial, celestial.position)
	}

	// Render the player's ship last, in the middle
	renderer.Render(ss.playerShip, renderer.Center())

	renderer.Update()
}
