package spacegame

import (
	"fmt"
	"sync"

	"github.com/faiface/pixel"
)

type Scene interface {
	Render()
	tick(float64)
}

// TODO: Extract Ship into Player...
type SpaceScene struct {
	renderer   Renderer
	camera     Camera
	system     *SolarSystem
	playerShip *Ship
	entities   []Entity
	starscape  Background
}

type SceneInformation struct {
	Entities   []Entity
	Celestials []Celestial
}

func NewSpaceScene(system *SolarSystem, player *Player, renderer Renderer) *SpaceScene {
	camera := NewChaseCamera(player.Ship())
	return &SpaceScene{
		camera:     camera,
		system:     system,
		playerShip: player.Ship(),
		entities:   []Entity{player.Ship()},
		renderer:   renderer,
		starscape:  NewStarscape(renderer, camera, 0.2),
	}
}

// TODO: Render with camera
func (ss *SpaceScene) Render() {
	//

	// Render the background
	ss.renderer.Clear()

	// Starscape
	ss.starscape.Render()

	// Directional arrow TODO

	// Render any planets in this scene
	for _, celestial := range ss.system.Celestials() {
		ss.camera.Render(ss.renderer, celestial)
	}

	// Asteroids (if any) TODO

	// Friendly players/NPCs TODO

	// Projectiles TODO

	// Neutral or enemy players/NPCs TODO

	// Render the player's ship last, in the middle
	ss.camera.Render(ss.renderer, ss.playerShip)

    // Special effects, explosions

	// Very basic HUD 
    // Make it better TODO: The ship "renders" the HUD! And can therefore respond appropriately to stimuli -- needs some engineering
    pos := ss.playerShip.Coordinates()
	hudTxt := fmt.Sprintf("Position: %.0f, %.0f\nVelocity: %4.2f", pos.X, pos.Y, ss.playerShip.Velocity().Len())
	ss.renderer.Text(hudTxt, pixel.V(30, 30))

    // Get HUD elements from player's ship

    // Apply HUD elements' renderer methods

	ss.renderer.Update()
}

func (ss *SpaceScene) tick(dt float64) {
	// TODO: Beware dragons... hehe, learning opportunity
	go ss.starscape.Displace(dt)

	si := SceneInformation{
		Celestials: ss.system.Celestials(),
		Entities:   ss.entities,
	}

	var wg sync.WaitGroup
	for _, entity := range ss.entities {
		// all pilotable ships get updated
		if ship, ok := entity.(PilotableShip); ok {
			wg.Add(1)
			go func(s PilotableShip) {
				defer wg.Done()
				s.Update(si)
			}(ship)
		}

		// Everything gets translated by its velocity
		entity.Translate(entity.Velocity())
	}
	wg.Wait()
}
