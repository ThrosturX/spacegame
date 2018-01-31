package spacegame

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GameEngine struct {
	universe    *Universe
	localSystem *SolarSystem
	player      *Player
	scene       Scene
	window      *pixelgl.Window // REPLACE WITH EVENT MANAGER!!
	renderer    Renderer
}

// TODO: Options parameter
func NewGame() GameEngine {
	// TODO: Get bounds from monitor
	cfg := pixelgl.WindowConfig{
		Title:  "Space Game!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	window.SetSmooth(true)

	universe := NewUniverse()

	player := NewPlayer("Cap'n Hector")

	// select a suitable start location
	var startSystem *SolarSystem
	for _, system := range universe.Systems() {
		// select the first system for now // TODO
		startSystem = system
		break
	}

	ge := GameEngine{
		player:   player,
		universe: universe,
		window:   window, // TODO: Event manager
		renderer: NewPixelWindowRenderer(window, "resources/images"),
		scene:    NewSpaceScene(startSystem, player, pixel.V(0, 0)),
	}

	return ge
}

func (ge *GameEngine) Run() {
	last := time.Now()
	for !ge.window.Closed() { // TODO: event manager
		dt := time.Since(last).Seconds()
		last = time.Now()

		if ge.window.Pressed(pixelgl.KeyEscape) {
			// TODO: Are you sure you want to quit or other menu
			break
		}

		ge.tick(dt)

		// Render everything (refactor.. decouple)

		ge.scene.Render(ge.renderer)

		// Draw extra UI elements
		// If extra UI elements...
		if false {
			ge.window.Update()
		}
	}

}

func (ge *GameEngine) tick(dt float64) {
	// Check key events and update game state

	var da float64
	if ge.window.Pressed(pixelgl.KeyLeft) {
		da = +3 * dt
	} else if ge.window.Pressed(pixelgl.KeyRight) {
		da = -3 * dt
	}
	if da != 0 {
		ge.player.ship.Turn(da)
	}
}
