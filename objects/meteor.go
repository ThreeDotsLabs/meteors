package objects

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"astrogame/assets"
	"astrogame/config"
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

func NewMeteor(baseVelocity float64) *Meteor {
	target := config.Vector{
		X: config.ScreenWidth / 2,
		Y: config.ScreenHeight / 2,
	}

	pos := config.Vector{
		X: config.ScreenWidth * rand.Float64(),
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

	modSprite := ScaleImg(assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))], float64(RandInt(5, 8))/10)

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
	RotateAndTranslateObject(m.rotation, m.sprite, screen, m.position.X, m.position.Y)
}

func (m *Meteor) Collider() config.Rect {
	bounds := m.sprite.Bounds()

	return config.NewRect(
		m.position.X,
		m.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
