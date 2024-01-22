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
	position           config.Vector
	movement           config.Vector
	target             config.Vector
	rotation           float64
	owner              string
	wType              *config.WeaponType
	intercectAnimation *Animation
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
	Shoot         func(p *Player)
}

type Blow struct {
	circle config.Circle
	Damage int
	Steps  int
	Step   int
}

var screenDiag = math.Sqrt(config.ScreenWidth*config.ScreenWidth + config.ScreenHeight*config.ScreenHeight)

func NewWeapon(wType string, p *Player) *Weapon {
	x, y := ebiten.CursorPosition()
	switch wType {
	case config.LightRocket:
		lightRType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.MissileSprite, 0.7),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Velocity:                      400 * p.params.LightRocketVelocityMultiplier,
			Damage:                        3,
			TargetType:                    "straight",
			WeaponName:                    config.LightRocket,
		}
		lightR := Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType:    lightRType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (250 - p.params.LightRocketSpeedUpscale)),
			ammo:          100,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				halfH := float64(bounds.Dy()) / 2

				spawnPos := config.Vector{
					X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
					Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
				}
				animation := NewAnimation(config.Vector{}, lightRType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectile := NewProjectile(config.Vector{}, spawnPos, p.rotation, lightRType, animation)
				projectile.owner = "player"
				p.game.AddProjectile(projectile)
			},
		}
		return &lightR
	case config.DoubleLightRocket:
		boubleRType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.DoubleMissileSprite, 0.7),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Velocity:                      400 * p.params.DoubleLightRocketVelocityMultiplier,
			Damage:                        3,
			TargetType:                    "straight",
			WeaponName:                    config.DoubleLightRocket,
		}
		doubleR := Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType:    boubleRType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (250 - p.params.DoubleLightRocketSpeedUpscale)),
			ammo:          50,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfWleft := float64(bounds.Dx()) / 4
				halfWright := float64(bounds.Dx()) - float64(bounds.Dx())/4
				halfH := float64(bounds.Dy()) / 2

				spawnPosLeft := config.Vector{
					X: p.position.X + halfWleft + math.Sin(p.rotation)*bulletSpawnOffset,
					Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
				}
				animationLeft := NewAnimation(config.Vector{}, boubleRType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileLeft := NewProjectile(config.Vector{}, spawnPosLeft, p.rotation, boubleRType, animationLeft)

				spawnPosRight := config.Vector{
					X: p.position.X + halfWright + math.Sin(p.rotation)*bulletSpawnOffset,
					Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
				}
				animationRight := NewAnimation(config.Vector{}, boubleRType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileRight := NewProjectile(config.Vector{}, spawnPosRight, p.rotation, boubleRType, animationRight)

				projectileLeft.owner = "player"
				projectileRight.owner = "player"
				p.game.AddProjectile(projectileLeft)
				p.game.AddProjectile(projectileRight)
			},
		}
		return &doubleR
	case config.LaserCannon:
		laserCType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.LaserCannon, 0.5),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Damage:                        3,
			TargetType:                    "straight",
			WeaponName:                    config.LaserCannon,
		}
		laserC := Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType:    laserCType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (300 - p.params.LaserCannonSpeedUpscale)),
			ammo:          30,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				halfH := float64(bounds.Dy()) / 2

				spawnPos := config.Vector{
					X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
					Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
				}
				beam := NewBeam(config.Vector{X: float64(x), Y: float64(y)}, p.rotation, spawnPos, laserCType)
				beam.owner = "player"
				p.game.AddBeam(beam)
			},
		}
		return &laserC
	case config.ClusterMines:
		clusterMType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.ClusterMines, 0.5),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Velocity:                      240 * p.params.ClusterMinesVelocityMultiplier,
			Damage:                        3,
			TargetType:                    "straight",
			WeaponName:                    config.ClusterMines,
		}
		clusterM := Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType:    clusterMType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (400 - p.params.ClusterMinesSpeedUpscale)),
			ammo:          5,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				halfH := float64(bounds.Dy()) / 2

				for i := 0; i < 20; i++ {
					angle := (math.Pi / 180) * float64(i*18)
					spawnPos := config.Vector{
						X: p.position.X + halfW + math.Sin(angle),
						Y: p.position.Y + halfH + math.Cos(angle),
					}
					animation := NewAnimation(config.Vector{}, clusterMType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
					projectile := NewProjectile(config.Vector{}, spawnPos, angle, clusterMType, animation)
					projectile.owner = "player"
					p.game.AddProjectile(projectile)
				}
			},
		}
		return &clusterM
	case config.BigBomb:
		bigBType := &config.WeaponType{
			Sprite:     objects.ScaleImg(assets.BigBomb, 0.8),
			Velocity:   200 * p.params.BigBombVelocityMultiplier,
			Damage:     10,
			TargetType: "straight",
			WeaponName: config.BigBomb,
		}
		bigB := Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType:    bigBType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (600 - p.params.BigBombSpeedUpscale)),
			ammo:          20,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				halfH := float64(bounds.Dy()) / 2

				spawnPos := config.Vector{
					X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
					Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
				}
				animation := NewAnimation(config.Vector{}, bigBType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectile := NewProjectile(config.Vector{}, spawnPos, p.rotation, bigBType, animation)
				projectile.owner = "player"
				p.game.AddProjectile(projectile)
			},
		}
		return &bigB
	case config.MachineGun:
		machineGType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.MachineGun, 0.3),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Velocity:                      600 * p.params.MachineGunVelocityMultiplier,
			Damage:                        1,
			TargetType:                    "straight",
			WeaponName:                    config.MachineGun,
		}
		machineG := Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType:    machineGType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (160 - p.params.MachineGunSpeedUpscale)),
			ammo:          99,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				halfH := float64(bounds.Dy()) / 2

				spawnPos := config.Vector{
					X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
					Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
				}
				animation := NewAnimation(config.Vector{}, machineGType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectile := NewProjectile(config.Vector{}, spawnPos, p.rotation, machineGType, animation)
				projectile.owner = "player"
				p.game.AddProjectile(projectile)
			},
		}
		return &machineG
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
			Sprite:                        assets.EnemyLightMissile,
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Velocity:                      150,
			Damage:                        1,
			TargetType:                    "straight",
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
			Sprite:                        assets.EnemyAutoLightMissile,
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Velocity:                      3,
			Damage:                        5,
			TargetType:                    "auto",
		},
	},
	ammo: 3,
}

func NewProjectile(target config.Vector, pos config.Vector, rotation float64, wType *config.WeaponType, animation *Animation) *Projectile {
	bounds := wType.Sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH

	p := &Projectile{
		position:           pos,
		rotation:           rotation,
		target:             target,
		wType:              wType,
		intercectAnimation: animation,
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
