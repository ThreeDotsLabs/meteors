package game

import (
	"astrogame/config"
	"astrogame/objects"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	position      config.Vector
	sprite        *ebiten.Image
	speed         int
	looping       bool
	run           bool
	numFrames     int
	startAt       int // which frame
	numberOfPlays int
	currF         int
	frameHeight   int
	frameWidth    int
}

func NewAnimation(position config.Vector, sprite *ebiten.Image, speed int, numFrames int, frameHeight int, frameWidth int) *Animation {
	return &Animation{
		position:      position,
		sprite:        sprite,
		speed:         speed,
		looping:       false,
		run:           true,
		numFrames:     numFrames,
		startAt:       0,
		numberOfPlays: 1,
		currF:         0,
		frameHeight:   frameHeight,
		frameWidth:    frameWidth,
	}
}

func (a *Animation) Update() {
	if !a.run {
		return
	}
	a.currF++
}

func (a *Animation) Draw(screen *ebiten.Image) {
	sprites := LoadSpritesheet(a.sprite, a.numFrames, a.frameWidth, a.frameHeight)
	objects.RotateAndTranslateObject(0, sprites[a.currF], screen, a.position.X, a.position.Y)
}

func LoadSpritesheet(sourceImg *ebiten.Image, n int, width int, height int) []*ebiten.Image {
	sprites := []*ebiten.Image{}
	numOfLines := sourceImg.Bounds().Dy() / height
	numFramesInLine := sourceImg.Bounds().Dx() / width
	for l := 0; l < numOfLines; l++ {
		for i := 0; i < numFramesInLine; i++ {
			dH := l * height
			if l == 0 {
				dH = height
			}
			dimensions := image.Rect(i*width, l*height, (i+1)*width, dH)
			sprite := sourceImg.SubImage(dimensions).(*ebiten.Image)
			sprites = append(sprites, sprite)
		}
	}

	return sprites
}
