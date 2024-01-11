package ecsexample

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteRenderSystem struct {
	Registry *Registry
}

func (s *SpriteRenderSystem) Draw(screen *ebiten.Image) {

	for _, e := range s.Registry.Query(TransformType, SpriteType) {

		position := e.GetComponent(TransformType).(*TransformComponent)
		sprite := e.GetComponent(SpriteType).(*SpriteComponent)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(position.PosX, position.PosY)
		screen.DrawImage(sprite.Image, op)
	}
}

type AnimationSystem struct {
	Registry *Registry
}

func (s *AnimationSystem) Update() error {
	entities := s.Registry.Query(SpriteType, AnimationType)

	for _, e := range entities {
		a := e.GetComponent(AnimationType).(*AnimationComponent)
		s := e.GetComponent(SpriteType).(*SpriteComponent)

		// advance animation
		a.Count += a.AnimationSpeed
		a.CurrentFrameIndex = int(math.Floor(a.Count))

		if a.CurrentFrameIndex >= len(a.Frames) { // restart animation
			a.Count = 0
			a.CurrentFrameIndex = 0
		}

		// update image
		s.Image = a.Frames[a.CurrentFrameIndex]
	}
	return nil
}
