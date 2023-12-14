package objects

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"astrogame/assets"
	"astrogame/config"
)

const (
	bulletSpeedPerSecond = 350.0
)

type Bullet struct {
	position config.Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewBullet(pos config.Vector, rotation float64) *Bullet {
	sprite := ScaleImg(assets.MissileSprite, 0.7)
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH

	b := &Bullet{
		position: pos,
		rotation: rotation,
		sprite:   sprite,
	}

	return b
}

func (b *Bullet) Update() {
	speed := bulletSpeedPerSecond / float64(ebiten.TPS())

	b.position.X += math.Sin(b.rotation) * speed
	b.position.Y += math.Cos(b.rotation) * -speed
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	RotateAndTranslateObject(b.rotation, b.sprite, screen, b.position.X, b.position.Y)
}

func (b *Bullet) Collider() config.Rect {
	bounds := b.sprite.Bounds()

	return config.NewRect(
		b.position.X,
		b.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
