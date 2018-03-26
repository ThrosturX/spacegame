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
	controller  *Controller
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

	resourceManager := NewStandardResourceManager("resources")

	universe := NewUniverse(resourceManager)

	// TODO: Import(options GameOptions)
	resourceManager.ImportDefault()

	// TODO: ctor won't need resourceManager
	player := NewPlayer("Cap'n Hector", resourceManager)

	// select a suitable start location
	var startSystem *SolarSystem
	for _, system := range universe.Systems() {
		// select the first system for now // TODO
		startSystem = system
		break
	}

	renderer := NewPixelWindowRenderer(window, resourceManager)

	ge := GameEngine{
		player:     player,
		controller: NewPlayerController(window, player),
		universe:   universe,
		window:     window, // TODO: Event manager
		renderer:   renderer,
		scene:      NewSpaceScene(startSystem, player, renderer),
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

		go ge.controller.relay(dt)
		ge.tick(dt)

		// Render everything (refactor.. decouple)

		ge.scene.Render()

		// Draw extra UI elements
		// If extra UI elements...
		if false {
			ge.window.Update()
		}
	}

}

func (ge *GameEngine) tick(dt float64) {
	// Check key events and update game state
	go ge.player.tick() // TODO: Move to scene, send scene info as parameter or otherwise find way to give player a context
	ge.scene.tick(dt)
}
