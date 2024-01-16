package game

import (
	"astrogame/config"
	"astrogame/objects"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	position      config.Vector
	rotationPoint config.Vector
	rotation      float64
	sprites       []*ebiten.Image
	speed         int
	looping       bool
	run           bool
	numFrames     int
	startAt       int // which frame
	numberOfPlays int
	currF         int
	frameHeight   int
	frameWidth    int
	curTick       int
	name          string
}

func NewAnimation(position config.Vector, sprite *ebiten.Image, speed int, numFrames int, frameHeight int, frameWidth int, looping bool, name string, rotation float64) *Animation {
	sprites := LoadSpritesheet(sprite, numFrames, frameWidth, frameHeight)
	return &Animation{
		position:      position,
		sprites:       sprites,
		speed:         speed,
		looping:       looping,
		run:           true,
		numFrames:     numFrames,
		startAt:       0,
		numberOfPlays: 1,
		currF:         0,
		curTick:       0,
		frameHeight:   frameHeight,
		frameWidth:    frameWidth,
		name:          name,
		rotation:      rotation,
	}
}

func (a *Animation) Update() {
	if !a.run {
		return
	}
	if a.curTick < a.speed {
		a.curTick++
		return
	}
	if a.looping && a.currF == a.numFrames-1 {
		a.currF = 0
		a.curTick = 0
	}
	a.curTick = 0
	a.currF++
}

func (a *Animation) Draw(screen *ebiten.Image) {
	switch a.name {
	case "engineFireburst":
		objects.RotateAndTranslateAnimation(a.rotation, a.sprites[a.currF], screen, a.position.X, a.position.Y)
	default:
		objects.RotateAndTranslateObject(a.rotation, a.sprites[a.currF], screen, a.position.X, a.position.Y)
	}
}

func LoadSpritesheet(sourceImg *ebiten.Image, n int, width int, height int) []*ebiten.Image {
	sprites := []*ebiten.Image{}
	numOfLines := sourceImg.Bounds().Dy() / height
	numFramesInLine := sourceImg.Bounds().Dx() / width
	for l := 0; l < numOfLines; l++ {
		for i := 0; i < numFramesInLine; i++ {
			dH := l*height + height
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
