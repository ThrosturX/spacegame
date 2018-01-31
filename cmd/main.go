package main

import (
	"spacegame"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	gameEngine := spacegame.NewGame()
	gameEngine.Run()
}

func main() {
	pixelgl.Run(run)
}
