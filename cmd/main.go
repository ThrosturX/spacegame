package main

import (
	"spacegame"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
    // TODO: Start menu that creates the engine
	gameEngine := spacegame.NewGame()
	gameEngine.Run()
}

func main() {
	pixelgl.Run(run)
}
