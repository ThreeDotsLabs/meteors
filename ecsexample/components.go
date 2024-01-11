package ecsexample

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ComponentType defines the supported component types in a user readable format
type ComponentType string

const (
	AnimationType ComponentType = "ANIMATION"
	SpriteType    ComponentType = "SPRITE"
	TransformType ComponentType = "TRANSFORM"
)

// ComponentTyper returns the type of a component
type ComponentTyper interface{ Type() ComponentType }

type TransformComponent struct {
	PosX, PosY float64
}

func (t *TransformComponent) Type() ComponentType { return TransformType }

type SpriteComponent struct {
	Image *ebiten.Image
}

func (t *SpriteComponent) Type() ComponentType { return SpriteType }

type AnimationComponent struct {
	Frames            []*ebiten.Image
	CurrentFrameIndex int
	Count             float64
	AnimationSpeed    float64
}

func (a AnimationComponent) Type() ComponentType { return AnimationType }
