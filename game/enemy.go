package game

import (
	"math"
	"time"

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

	modSprite := objects.ScaleImg(enType.Sprite, 0.5)
	enType.Sprite = modSprite

	e := &Enemy{
		game:      g,
		position:  pos,
		movement:  movement,
		enemyType: &enType,
		HP:        enType.StartHP,
	}

	switch enType.WeaponTypeStr {
	case config.LightRocket:
		startWeapon := enemyLightRocket
		e.weapon = startWeapon
		e.weapon.shootCooldown = config.NewTimer(time.Millisecond * 2000)
	case config.AutoLightRocket:
		startWeapon := enemyAutoLightRocket
		e.weapon = startWeapon
		e.weapon.shootCooldown = config.NewTimer(time.Millisecond * 2300)
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
	if e.weapon.projectile.wType != nil {
		e.weapon.shootCooldown.Update()
		if e.weapon.shootCooldown.IsReady() {
			if e.weapon.ammo <= 0 {
				return
			}
			e.weapon.shootCooldown.Reset()
			bounds := e.enemyType.Sprite.Bounds()
			halfW := float64(bounds.Dx()) / 2
			halfH := float64(bounds.Dy()) / 2

			spawnPos := config.Vector{
				X: e.position.X + halfW + math.Sin(e.rotation)*bulletSpawnOffset,
				Y: e.position.Y + halfH + math.Cos(e.rotation)*bulletSpawnOffset,
			}

			projectile := NewProjectile(config.Vector{}, spawnPos, e.rotation, e.weapon.projectile.wType)
			projectile.owner = "enemy"
			e.game.AddProjectile(projectile)
			e.weapon.ammo--
		}
	}
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
