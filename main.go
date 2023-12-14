package main

import (
	"astrogame/config"
	"astrogame/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Astro Ship (Ebitengine Demo)")
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
