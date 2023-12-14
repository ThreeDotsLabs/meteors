package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"astrogame/config"
	"astrogame/objects"
)

type Enemy struct {
	position  config.Vector
	target    config.Vector
	rotation  float64
	movement  config.Vector
	enemyType *config.EnemyType
}

func NewEnemy(target config.Vector, pos config.Vector, enType config.EnemyType) *Enemy {
	direction := config.Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}
	normalizedDirection := direction.Normalize()

	movement := config.Vector{
		X: normalizedDirection.X * enType.Velocity,
		Y: normalizedDirection.Y * enType.Velocity,
	}

	modSprite := objects.ScaleImg(enType.Sprite, 0.5)
	enType.Sprite = modSprite

	e := &Enemy{
		position:  pos,
		movement:  movement,
		enemyType: &enType,
	}

	return e
}

func (e *Enemy) SetDirection(target config.Vector, pos config.Vector, enType config.EnemyType) {
	direction := config.Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}

	normalizedDirection := direction.Normalize()

	e.movement = config.Vector{
		X: normalizedDirection.X * enType.Velocity,
		Y: normalizedDirection.Y * enType.Velocity,
	}
	e.position = pos
}

func (e *Enemy) Update() {
	e.position.X += e.movement.X
	e.position.Y += e.movement.Y
	e.rotation += e.enemyType.RotationSpeed
	direction := config.Vector{
		X: e.target.X - e.position.X,
		Y: e.target.Y - e.position.Y,
	}
	normalizedDirection := direction.Normalize()

	movement := config.Vector{
		X: normalizedDirection.X * e.enemyType.Velocity,
		Y: normalizedDirection.Y * e.enemyType.Velocity,
	}
	e.movement = movement
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	objects.RotateAndTranslateObject(e.rotation, e.enemyType.Sprite, screen, e.position.X, e.position.Y)
}

func (e *Enemy) Collider() config.Rect {
	bounds := e.enemyType.Sprite.Bounds()

	return config.NewRect(
		e.position.X,
		e.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
