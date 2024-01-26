package main

import (
	"astrogame/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()
	// fw, _ := ebiten.ScreenSizeInFullscreen()
	if g.Options.Fullscreen {
		ebiten.SetFullscreen(true)
	}
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// if g.Options.Fullscreen && g.Options.ScreenWidth < float64(fw) {
	// 	xOffsetMod := (float64(fw) - float64(g.Options.ScreenWidth)) / 2
	// 	ebiten.SetWindowPosition(int(xOffsetMod), 0)
	// }
	ebiten.SetWindowSize(int(g.Options.ScreenWidth), int(g.Options.ScreenHeight))
	ebiten.SetWindowTitle("Astro Ship (Ebitengine Demo)")
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
