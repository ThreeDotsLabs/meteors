package main

import (
	"astrogame/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(int(g.Options.ScreenWidth), int(g.Options.ScreenHeight))
	ebiten.SetWindowTitle("Astro Ship (Ebitengine Demo)")
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
