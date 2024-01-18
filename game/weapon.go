package game

import (
	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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
	rotation float64
	owner    string
	Damage   int
	Line     config.Line
}

type BeamAnimation struct {
	curRect  config.Rect
	rotation float64
	Steps    int
	Step     int
}

type Weapon struct {
	projectile    Projectile
	ammo          int
	shootCooldown *config.Timer
}

type Blow struct {
	circle config.Circle
	Damage int
	Steps  int
	Step   int
}

var screenDiag = math.Sqrt(config.ScreenWidth*config.ScreenWidth + config.ScreenHeight*config.ScreenHeight)

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
	case config.ClusterMines:
		return &Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType: &config.WeaponType{
					Velocity:   200,
					Sprite:     objects.ScaleImg(assets.ClusterMines, 0.5),
					Damage:     2,
					TargetType: "straight",
					WeaponName: config.ClusterMines,
				},
			},
			shootCooldown: config.NewTimer(time.Millisecond * 400),
			ammo:          5,
		}
	case config.BigBomb:
		return &Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType: &config.WeaponType{
					Velocity:   200,
					Sprite:     objects.ScaleImg(assets.BigBomb, 0.8),
					Damage:     5,
					TargetType: "straight",
					WeaponName: config.BigBomb,
				},
			},
			shootCooldown: config.NewTimer(time.Millisecond * 600),
			ammo:          20,
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

func (p *Projectile) Collider() image.Rectangle {
	bounds := p.wType.Sprite.Bounds()
	return image.Rectangle{
		Min: image.Point{
			X: int(p.position.X),
			Y: int(p.position.Y),
		},
		Max: image.Point{
			X: int(p.position.X + float64(bounds.Dx())),
			Y: int(p.position.Y + float64(bounds.Dy())),
		},
	}
}

func NewBeam(target config.Vector, rotation float64, pos config.Vector, wType *config.WeaponType) *Beam {
	bounds := wType.Sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	pos.X -= halfW
	pos.Y -= halfH

	line := config.NewLine(
		pos.X,
		pos.Y,
		math.Cos(rotation-math.Pi/2)*(screenDiag)+pos.X,
		math.Sin(rotation-math.Pi/2)*(screenDiag)+pos.Y,
	)
	b := &Beam{
		position: pos,
		target:   target,
		rotation: rotation,
		Damage:   wType.Damage,
		Line:     line,
	}

	return b
}

func (b *Beam) NewBeamAnimation() *BeamAnimation {
	rect := config.NewRectangle(
		b.position.X,
		b.position.Y,
		float64(1),
		float64(screenDiag),
	)
	return &BeamAnimation{
		curRect:  rect,
		Steps:    5,
		Step:     1,
		rotation: b.rotation + math.Pi,
	}
}

func (b *BeamAnimation) Update() {
	b.curRect.Width += float64(b.Step)
	b.curRect.X -= 1
	b.Step++
}

func (b *BeamAnimation) Draw(screen *ebiten.Image) {
	rectImage := ebiten.NewImage(int(b.curRect.Width), int(b.curRect.Height))
	rectImage.Fill(color.White)
	rotationOpts := &ebiten.DrawImageOptions{}
	rotationOpts.GeoM.Rotate(b.rotation)
	rotationOpts.GeoM.Translate(b.curRect.X, b.curRect.Y)
	screen.DrawImage(rectImage, rotationOpts)
}

func NewBlow(x, y, radius float64, damage int) *Blow {
	return &Blow{
		circle: config.Circle{
			X:      x,
			Y:      y,
			Radius: radius,
		},
		Damage: damage,
	}
}

func (b *Blow) Update() {
	if b.Step < b.Steps {
		b.Step++
	}
}

// func (b *Blow) Draw(screen *ebiten.Image) {
// 	vector.DrawFilledCircle(screen, float32(b.circle.X), float32(b.circle.Y), float32(b.circle.Radius), color.RGBA{255, 255, 255, 255}, false)
// }
