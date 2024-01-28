package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"astrogame/config"
	"astrogame/objects"
)

type Enemy struct {
	game       *Game
	position   config.Vector
	target     config.Vector
	rotation   float64
	TargetType string
	movement   config.Vector
	enemyType  *config.EnemyType
	weapon     Weapon
	HP         int
}

func NewEnemy(g *Game, target config.Vector, pos config.Vector, enType config.EnemyType) *Enemy {
	direction := config.Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}
	normalizedDirection := direction.Normalize()

	movement := config.Vector{
		X: normalizedDirection.X * enType.Velocity,
		Y: normalizedDirection.Y * enType.Velocity,
	}
	enType.Sprite = objects.ScaleImg(enType.Sprite, 0.5)

	e := &Enemy{
		game:      g,
		position:  pos,
		movement:  movement,
		enemyType: &enType,
		HP:        enType.StartHP,
	}
	if enType.WeaponType != nil {
		e.weapon = NewEnemyWeapon(enType.WeaponType)
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
	//e.rotation += e.enemyType.RotationSpeed
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
	if e.TargetType == config.TargetTypePlayer {
		e.rotation = math.Atan2(float64(e.target.Y-e.position.Y), float64(e.target.X-e.position.X))
		e.rotation -= (90 * math.Pi) / 180
	}
	if e.weapon.projectile.wType != nil {
		e.weapon.shootCooldown.Update()
		if e.weapon.shootCooldown.IsReady() {
			if e.weapon.ammo <= 0 {
				return
			}
			e.weapon.shootCooldown.Reset()
			e.weapon.EnemyShoot(e)
			e.weapon.ammo--
		}
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	objects.RotateAndTranslateObject(e.rotation, e.enemyType.Sprite, screen, e.position.X, e.position.Y)
}

func (e *Enemy) Collider() image.Rectangle {
	bounds := e.enemyType.Sprite.Bounds()
	return image.Rectangle{
		Min: image.Point{
			X: int(e.position.X),
			Y: int(e.position.Y),
		},
		Max: image.Point{
			X: int(e.position.X + float64(bounds.Dx())),
			Y: int(e.position.Y + float64(bounds.Dy())),
		},
	}
}
