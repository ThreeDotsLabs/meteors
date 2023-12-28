package game

import (
	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Projectile struct {
	position config.Vector
	movement config.Vector
	target   config.Vector
	rotation float64
	owner    string
	wType    *config.WeaponType
}

type Beam struct {
	position config.Vector
	target   config.Vector
	owner    string
	Damage   int
}

type BeamAnimation struct {
	curRect config.Rect
	Steps   int
	Step    int
}

type Weapon struct {
	projectile    Projectile
	ammo          int
	shootCooldown *config.Timer
}

func NewWeapon(wType string) *Weapon {
	switch wType {
	case config.LightRocket:
		return &Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType: &config.WeaponType{
					Sprite:     objects.ScaleImg(assets.MissileSprite, 0.7),
					Velocity:   400,
					Damage:     1,
					TargetType: "straight",
					WeaponName: config.LightRocket,
				},
			},
			shootCooldown: config.NewTimer(time.Millisecond * 250),
			ammo:          100,
		}
	case config.DoubleLightRocket:
		return &Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType: &config.WeaponType{
					Sprite:     objects.ScaleImg(assets.DoubleMissileSprite, 0.7),
					Velocity:   400,
					Damage:     1,
					TargetType: "straight",
					WeaponName: config.DoubleLightRocket,
				},
			},
			shootCooldown: config.NewTimer(time.Millisecond * 250),
			ammo:          50,
		}
	case config.LaserCannon:
		return &Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType: &config.WeaponType{
					Sprite:     objects.ScaleImg(assets.LaserCannon, 0.5),
					Damage:     2,
					TargetType: "straight",
					WeaponName: config.LaserCannon,
				},
			},
			shootCooldown: config.NewTimer(time.Millisecond * 300),
			ammo:          30,
		}
	}
	return nil
}

var enemyLightRocket = Weapon{
	projectile: Projectile{
		position: config.Vector{},
		target:   config.Vector{},
		movement: config.Vector{},
		rotation: 0,
		wType: &config.WeaponType{
			Sprite:     assets.EnemyLightMissile,
			Velocity:   150,
			Damage:     1,
			TargetType: "straight",
		},
	},
	ammo: 10,
}

var enemyAutoLightRocket = Weapon{
	projectile: Projectile{
		position: config.Vector{},
		target:   config.Vector{},
		movement: config.Vector{},
		rotation: 0,
		wType: &config.WeaponType{
			Sprite:     assets.EnemyAutoLightMissile,
			Velocity:   3,
			Damage:     5,
			TargetType: "auto",
		},
	},
	ammo: 3,
}

func NewProjectile(target config.Vector, pos config.Vector, rotation float64, wType *config.WeaponType) *Projectile {
	bounds := wType.Sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH

	p := &Projectile{
		position: pos,
		rotation: rotation,
		target:   target,
		wType:    wType,
	}

	return p
}

func (p *Projectile) Update() {
	if p.wType.TargetType == "auto" {
		p.position.X += p.movement.X
		p.position.Y += p.movement.Y

		direction := config.Vector{
			X: p.target.X - p.position.X,
			Y: p.target.Y - p.position.Y,
		}
		normalizedDirection := direction.Normalize()

		movement := config.Vector{
			X: normalizedDirection.X * p.wType.Velocity,
			Y: normalizedDirection.Y * p.wType.Velocity,
		}
		p.movement = movement
		p.rotation = math.Atan2(float64(p.target.Y-p.position.Y), float64(p.target.X-p.position.X))
		p.rotation -= (90 * math.Pi) / 180
	} else {
		speed := p.wType.Velocity / float64(ebiten.TPS())
		quant := speed
		if p.owner == "player" {
			quant = -speed
		}
		p.position.X += math.Sin(p.rotation) * speed
		p.position.Y += math.Cos(p.rotation) * quant
	}
}

func (p *Projectile) Draw(screen *ebiten.Image) {
	objects.RotateAndTranslateObject(p.rotation, p.wType.Sprite, screen, p.position.X, p.position.Y)
}

func (p *Projectile) Collider() config.Rect {
	bounds := p.wType.Sprite.Bounds()

	return config.NewRect(
		p.position.X,
		p.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}

func (b *Beam) NewBeamAnimation() *BeamAnimation {
	return &BeamAnimation{
		curRect: b.Collider(),
		Steps:   400,
		Step:    1,
	}
}

func NewBeam(target config.Vector, pos config.Vector, wType *config.WeaponType) *Beam {
	bounds := wType.Sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH

	b := &Beam{
		position: pos,
		target:   target,
		Damage:   wType.Damage,
	}

	return b
}

func (b *Beam) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(b.position.X), float32(b.position.Y), 4, -float32(config.ScreenHeight-(config.ScreenHeight-b.position.Y)), color.RGBA{255, 255, 255, 255}, false)
}

func (b *Beam) Collider() config.Rect {
	return config.NewRect(
		b.position.X,
		0,
		float64(4),
		float64(config.ScreenHeight-(config.ScreenHeight-b.position.Y)),
	)
}

func (b *BeamAnimation) Update() {
	b.curRect.Width += float64(b.Step)
	b.Step++
}

func (b *BeamAnimation) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(b.curRect.X), float32(b.curRect.Y), float32(b.curRect.Width), -float32(config.ScreenHeight-(config.ScreenHeight-b.curRect.Y)), color.RGBA{255, 255, 255, (255 - uint8(b.Step*20))}, false)
}
