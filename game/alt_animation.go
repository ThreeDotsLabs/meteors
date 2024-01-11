package game

import (
	"astrogame/config"
	"astrogame/ecsexample"

	"github.com/hajimehoshi/ebiten/v2"
)

type AnimationExample struct {
	spriteRenderSystem *ecsexample.SpriteRenderSystem
	animationSystem    *ecsexample.AnimationSystem
}

func (a *AnimationExample) Draw(screen *ebiten.Image) {
	a.spriteRenderSystem.Draw(screen)
}

func (a *AnimationExample) Update() error {
	return a.animationSystem.Update()
}

func (a *AnimationExample) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
