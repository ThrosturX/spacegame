package spacegame

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

func NewSpaceScene(system *SolarSystem, player *Player, renderer Renderer) *SpaceScene {
	camera := NewChaseCamera(player.Ship())
	return &SpaceScene{
		camera:     camera,
		system:     system,
		playerShip: player.Ship(),
		entities:   []Entity{player.Ship()},
		renderer:   renderer,
		starscape:  NewStarscape(renderer, camera),
	}
}

func (ss *SpaceScene) Render() {
	//

	// Render the background
	ss.renderer.Clear()

	// Starscape
	ss.starscape.Render()

	// Render any planets in this scene
	for _, celestial := range ss.system.Celestials() {
		ss.camera.Render(ss.renderer, celestial)
	}

	// Render the player's ship last, in the middle
	ss.renderer.Render(ss.playerShip, ss.renderer.Center())

	ss.renderer.Update()
}

func (ss *SpaceScene) tick(dt float64) {
	// TODO: Beware dragons... hehe, learning opportunity
	go ss.starscape.Displace(dt)
	// Everything gets translated by its velocity
	for _, entity := range ss.entities {
		entity.Translate(entity.Velocity())
	}
}
