package game

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
)

const (
	rotationSpeedMin = -0.02
	rotationSpeedMax = 0.02
)

type Meteor struct {
	position      config.Vector
	rotation      float64
	movement      config.Vector
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewMeteor(baseVelocity float64, g *Game) *Meteor {
	target := config.Vector{
		X: g.Options.ScreenWidth / 2,
		Y: g.Options.ScreenHeight / 2,
	}

	pos := config.Vector{
		X: g.Options.ScreenWidth * rand.Float64(),
		Y: -10,
	}

	velocity := baseVelocity + rand.Float64()*1.5

	direction := config.Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}
	normalizedDirection := direction.Normalize()

	movement := config.Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	modSprite := objects.ScaleImg(assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))], float64(objects.RandInt(5, 8))/10)

	m := &Meteor{
		position:      pos,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        modSprite,
	}
	return m
}

func (m *Meteor) Update() {
	m.position.X += m.movement.X
	m.position.Y += m.movement.Y
	m.rotation += m.rotationSpeed
}

func (m *Meteor) Draw(screen *ebiten.Image) {
	objects.RotateAndTranslateObject(m.rotation, m.sprite, screen, m.position.X, m.position.Y)
}

func (m *Meteor) Collider() image.Rectangle {
	bounds := m.sprite.Bounds()
	return image.Rectangle{
		Min: image.Point{
			X: int(m.position.X),
			Y: int(m.position.Y),
		},
		Max: image.Point{
			X: int(m.position.X + float64(bounds.Dx())),
			Y: int(m.position.Y + float64(bounds.Dy())),
		},
	}
	// return config.NewRect(
	// 	m.position.X,
	// 	m.position.Y,
	// 	float64(bounds.Dx()),
	// 	float64(bounds.Dy()),
	// )
}
