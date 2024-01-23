package game

import (
	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
	"image"
	"image/color"
	"math"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Projectile struct {
	HP                 int
	position           config.Vector
	movement           config.Vector
	target             config.Vector
	rotation           float64
	owner              string
	wType              *config.WeaponType
	intercectAnimation *Animation
	instantAnimation   *Animation
}

type Beam struct {
	position config.Vector
	target   config.Vector
	rotation float64
	owner    string
	Damage   int
	Line     config.Line
	Steps    int
	Step     int
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
	EnemyShoot    func(e *Enemy)
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
			IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
			Velocity:                      400 + p.params.LightRocketVelocityMultiplier,
			Damage:                        3,
			TargetType:                    "straight",
			WeaponName:                    config.LightRocket,
		}
		lightR := Weapon{
			projectile: Projectile{
				wType: lightRType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (300 - p.params.LightRocketSpeedUpscale)),
			ammo:          100,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				w := p.sprite.Bounds().Dx()
				h := p.sprite.Bounds().Dy()
				px, py := p.position.X+float64(w)/2, p.position.Y+float64(h)/2
				spawnPos := config.Vector{
					X: px + ((p.position.X+halfW-px)*math.Cos(-p.rotation) - (py-p.position.Y)*math.Sin(-p.rotation)),
					Y: py - ((p.position.X+halfW-px)*math.Sin(-p.rotation) + (py-p.position.Y)*math.Cos(-p.rotation)),
				}
				animation := NewAnimation(config.Vector{}, lightRType.IntercectAnimationSpriteSheet, 1, 56, 60, false, "projectileBlow", 0)
				projectile := NewProjectile(config.Vector{}, spawnPos, p.rotation, lightRType, animation, 0)
				projectile.owner = "player"
				p.game.AddProjectile(projectile)
			},
		}
		return &lightR
	case config.DoubleLightRocket:
		boubleRType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.DoubleMissileSprite, 0.7),
			IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
			Velocity:                      400 + p.params.DoubleLightRocketVelocityMultiplier,
			Damage:                        3,
			TargetType:                    "straight",
			WeaponName:                    config.DoubleLightRocket,
		}
		doubleR := Weapon{
			projectile: Projectile{
				wType: boubleRType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (250 - p.params.DoubleLightRocketSpeedUpscale)),
			ammo:          50,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfWleft := float64(bounds.Dx()) / 4
				halfWright := float64(bounds.Dx()) - float64(bounds.Dx())/4
				w := p.sprite.Bounds().Dx()
				h := p.sprite.Bounds().Dy()
				px, py := p.position.X+float64(w)/2, p.position.Y+float64(h)/2
				spawnPosLeft := config.Vector{
					X: px + ((p.position.X+halfWleft-px)*math.Cos(-p.rotation) - (py-p.position.Y)*math.Sin(-p.rotation)),
					Y: py - ((p.position.X+halfWleft-px)*math.Sin(-p.rotation) + (py-p.position.Y)*math.Cos(-p.rotation)),
				}
				spawnPosRight := config.Vector{
					X: px + ((p.position.X+halfWright-px)*math.Cos(-p.rotation) - (py-p.position.Y)*math.Sin(-p.rotation)),
					Y: py - ((p.position.X+halfWright-px)*math.Sin(-p.rotation) + (py-p.position.Y)*math.Cos(-p.rotation)),
				}
				animationLeft := NewAnimation(config.Vector{}, boubleRType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileLeft := NewProjectile(config.Vector{}, spawnPosLeft, p.rotation, boubleRType, animationLeft, 0)
				animationRight := NewAnimation(config.Vector{}, boubleRType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileRight := NewProjectile(config.Vector{}, spawnPosRight, p.rotation, boubleRType, animationRight, 0)

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
				wType: laserCType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (500 - p.params.LaserCannonSpeedUpscale)),
			ammo:          30,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				w := p.sprite.Bounds().Dx()
				h := p.sprite.Bounds().Dy()
				px, py := p.position.X+float64(w)/2, p.position.Y+float64(h)/2
				spawnPos := config.Vector{
					X: px + ((p.position.X+halfW-px)*math.Cos(-p.rotation) - (py-p.position.Y)*math.Sin(-p.rotation)),
					Y: py - ((p.position.X+halfW-px)*math.Sin(-p.rotation) + (py-p.position.Y)*math.Cos(-p.rotation)),
				}
				beam := NewBeam(config.Vector{X: float64(x), Y: float64(y)}, p.rotation, spawnPos, laserCType)
				beam.owner = "player"
				p.game.AddBeam(beam)
				ba := beam.NewBeamAnimation()
				p.game.AddBeamAnimation(ba)
			},
		}
		return &laserC
	case config.DoubleLaserCannon:
		doubleLaserCType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.LaserCannon, 0.5),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Damage:                        2,
			TargetType:                    "straight",
			WeaponName:                    config.DoubleLaserCannon,
		}
		doubleLaserC := Weapon{
			projectile: Projectile{
				wType: doubleLaserCType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (460 - p.params.DoubleLaserCannonSpeedUpscale)),
			ammo:          30,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfWleft := float64(bounds.Dx()) / 4
				halfWright := float64(bounds.Dx()) - float64(bounds.Dx())/4
				heightShift := float64(bounds.Dy()) / 10
				w := p.sprite.Bounds().Dx()
				h := p.sprite.Bounds().Dy()
				px, py := p.position.X+float64(w)/2, p.position.Y+float64(h)/2
				spawnPosLeft := config.Vector{
					X: px + ((p.position.X+halfWleft-px)*math.Cos(-p.rotation) - (py-heightShift-p.position.Y)*math.Sin(-p.rotation)),
					Y: py - ((p.position.X+halfWleft-px)*math.Sin(-p.rotation) + (py-heightShift-p.position.Y)*math.Cos(-p.rotation)),
				}
				spawnPosRight := config.Vector{
					X: px + ((p.position.X+halfWright-px)*math.Cos(-p.rotation) - (py-heightShift-p.position.Y)*math.Sin(-p.rotation)),
					Y: py - ((p.position.X+halfWright-px)*math.Sin(-p.rotation) + (py-heightShift-p.position.Y)*math.Cos(-p.rotation)),
				}
				beamLeft := NewBeam(config.Vector{X: float64(x), Y: float64(y)}, p.rotation, spawnPosLeft, doubleLaserCType)
				beamLeft.owner = "player"
				p.game.AddBeam(beamLeft)
				baL := beamLeft.NewBeamAnimation()
				p.game.AddBeamAnimation(baL)
				beamRight := NewBeam(config.Vector{X: float64(x), Y: float64(y)}, p.rotation, spawnPosRight, doubleLaserCType)
				beamRight.owner = "player"
				p.game.AddBeam(beamRight)
				baR := beamRight.NewBeamAnimation()
				p.game.AddBeamAnimation(baR)
			},
		}
		return &doubleLaserC
	case config.ClusterMines:
		clusterMType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.ClusterMines, 0.5),
			IntercectAnimationSpriteSheet: assets.ClusterMinesBlowSpriteSheet,
			Velocity:                      240 + p.params.ClusterMinesVelocityMultiplier,
			Damage:                        3,
			TargetType:                    "straight",
			WeaponName:                    config.ClusterMines,
		}
		clusterM := Weapon{
			projectile: Projectile{
				wType: clusterMType,
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
					animation := NewAnimation(config.Vector{}, clusterMType.IntercectAnimationSpriteSheet, 1, 50, 50, false, "projectileBlow", 0)
					projectile := NewProjectile(config.Vector{}, spawnPos, angle, clusterMType, animation, 0)
					projectile.owner = "player"
					p.game.AddProjectile(projectile)
				}
			},
		}
		return &clusterM
	case config.BigBomb:
		bigBType := &config.WeaponType{
			Sprite:     objects.ScaleImg(assets.BigBomb, 0.8),
			Velocity:   200 + p.params.BigBombVelocityMultiplier,
			Damage:     10,
			TargetType: "straight",
			WeaponName: config.BigBomb,
		}
		bigB := Weapon{
			projectile: Projectile{
				wType: bigBType,
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
				projectile := NewProjectile(config.Vector{}, spawnPos, p.rotation, bigBType, animation, 0)
				projectile.owner = "player"
				p.game.AddProjectile(projectile)
			},
		}
		return &bigB
	case config.MachineGun:
		machineGType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.MachineGun, 0.3),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Velocity:                      600 + p.params.MachineGunVelocityMultiplier,
			Damage:                        1,
			TargetType:                    "straight",
			WeaponName:                    config.MachineGun,
		}
		machineG := Weapon{
			projectile: Projectile{
				wType: machineGType,
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
				projectile := NewProjectile(config.Vector{}, spawnPos, p.rotation, machineGType, animation, 0)
				projectile.owner = "player"
				p.game.AddProjectile(projectile)
			},
		}
		return &machineG
	case config.DoubleMachineGun:
		doubleMachineGType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.DoubleMachineGun, 0.28),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Velocity:                      600 + p.params.DoubleMachineGunVelocityMultiplier,
			Damage:                        1,
			TargetType:                    "straight",
			WeaponName:                    config.MachineGun,
		}
		doubleMachineG := Weapon{
			projectile: Projectile{
				wType: doubleMachineGType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (260 - p.params.DoubleMachineGunSpeedUpscale)),
			ammo:          99,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfWleft := float64(bounds.Dx()) / 4
				halfWright := float64(bounds.Dx()) - float64(bounds.Dx())/4
				w := p.sprite.Bounds().Dx()
				h := p.sprite.Bounds().Dy()
				px, py := p.position.X+float64(w)/2, p.position.Y+float64(h)/2
				spawnPosLeft := config.Vector{
					X: px + ((p.position.X+halfWleft-px)*math.Cos(-p.rotation) - (py-p.position.Y)*math.Sin(-p.rotation)),
					Y: py - ((p.position.X+halfWleft-px)*math.Sin(-p.rotation) + (py-p.position.Y)*math.Cos(-p.rotation)),
				}
				spawnPosRight := config.Vector{
					X: px + ((p.position.X+halfWright-px)*math.Cos(-p.rotation) - (py-p.position.Y)*math.Sin(-p.rotation)),
					Y: py - ((p.position.X+halfWright-px)*math.Sin(-p.rotation) + (py-p.position.Y)*math.Cos(-p.rotation)),
				}
				animationLeft := NewAnimation(config.Vector{}, doubleMachineGType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileLeft := NewProjectile(config.Vector{}, spawnPosLeft, p.rotation, doubleMachineGType, animationLeft, 0)
				animationRight := NewAnimation(config.Vector{}, doubleMachineGType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileRight := NewProjectile(config.Vector{}, spawnPosRight, p.rotation, doubleMachineGType, animationRight, 0)
				projectileLeft.owner = "player"
				projectileRight.owner = "player"
				p.game.AddProjectile(projectileLeft)
				p.game.AddProjectile(projectileRight)
			},
		}
		return &doubleMachineG
	case config.PlasmaGun:
		plasmaGType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.PlasmaGun, 0.8),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			InstantAnimationSpiteSheet:    assets.PlasmaGunProjectileSpriteSheet,
			Velocity:                      500 + p.params.PlasmaGunVelocityMultiplier,
			AnimationOnly:                 true,
			Damage:                        1,
			TargetType:                    "straight",
			WeaponName:                    config.PlasmaGun,
		}

		plasmaG := Weapon{
			projectile: Projectile{
				wType: plasmaGType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (760 - p.params.PlasmaGunSpeedUpscale)),
			ammo:          99,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				halfH := float64(bounds.Dy()) / 2

				spawnPos := config.Vector{
					X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
					Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
				}
				plasmaAnimation := NewAnimation(config.Vector{}, plasmaGType.InstantAnimationSpiteSheet, 1, 55, 50, true, "projectileInstant", 0)
				animation := NewAnimation(config.Vector{}, plasmaGType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectile := NewProjectile(config.Vector{}, spawnPos, p.rotation, plasmaGType, animation, 6)
				projectile.owner = "player"
				projectile.instantAnimation = plasmaAnimation
				p.game.AddProjectile(projectile)
				p.game.AddAnimation(plasmaAnimation)
			},
		}
		return &plasmaG
	case config.DoublePlasmaGun:
		doublePlasmaGType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.PlasmaGun, 1.2),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			InstantAnimationSpiteSheet:    assets.PlasmaGunProjectileSpriteSheet,
			Velocity:                      500 + p.params.DoublePlasmaGunVelocityMultiplier,
			AnimationOnly:                 true,
			Damage:                        1,
			TargetType:                    "straight",
			WeaponName:                    config.DoublePlasmaGun,
		}
		doublePlasmaG := Weapon{
			projectile: Projectile{
				position: config.Vector{},
				target:   config.Vector{},
				movement: config.Vector{},
				rotation: 0,
				wType:    doublePlasmaGType,
			},
			shootCooldown: config.NewTimer(time.Millisecond * (880 - p.params.DoublePlasmaGunSpeedUpscale)),
			ammo:          99,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfWleft := float64(bounds.Dx()) / 4
				halfWright := float64(bounds.Dx()) - float64(bounds.Dx())/4
				w := p.sprite.Bounds().Dx()
				h := p.sprite.Bounds().Dy()
				px, py := p.position.X+float64(w)/2, p.position.Y+float64(h)/2
				spawnPosLeft := config.Vector{
					X: px + ((p.position.X+halfWleft-px)*math.Cos(-p.rotation) - (py-p.position.Y)*math.Sin(-p.rotation)),
					Y: py - ((p.position.X+halfWleft-px)*math.Sin(-p.rotation) + (py-p.position.Y)*math.Cos(-p.rotation)),
				}
				spawnPosRight := config.Vector{
					X: px + ((p.position.X+halfWright-px)*math.Cos(-p.rotation) - (py-p.position.Y)*math.Sin(-p.rotation)),
					Y: py - ((p.position.X+halfWright-px)*math.Sin(-p.rotation) + (py-p.position.Y)*math.Cos(-p.rotation)),
				}
				plasmaAnimationLeft := NewAnimation(config.Vector{}, doublePlasmaGType.InstantAnimationSpiteSheet, 1, 55, 50, true, "projectileInstant", 0)
				animationLeft := NewAnimation(config.Vector{}, doublePlasmaGType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileLeft := NewProjectile(config.Vector{}, spawnPosLeft, p.rotation, doublePlasmaGType, animationLeft, 4)
				plasmaAnimationRight := NewAnimation(config.Vector{}, doublePlasmaGType.InstantAnimationSpiteSheet, 1, 55, 50, true, "projectileInstant", 0)
				animationRight := NewAnimation(config.Vector{}, doublePlasmaGType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileRight := NewProjectile(config.Vector{}, spawnPosRight, p.rotation, doublePlasmaGType, animationRight, 4)
				projectileLeft.owner = "player"
				projectileRight.owner = "player"
				projectileLeft.instantAnimation = plasmaAnimationLeft
				projectileRight.instantAnimation = plasmaAnimationRight
				p.game.AddProjectile(projectileLeft)
				p.game.AddProjectile(projectileRight)
				p.game.AddAnimation(plasmaAnimationLeft)
				p.game.AddAnimation(plasmaAnimationRight)
			},
		}
		return &doublePlasmaG
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
			IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
			Velocity:                      150,
			Damage:                        1,
			TargetType:                    "straight",
		},
	},
	ammo: 10,
	EnemyShoot: func(e *Enemy) {
		bounds := e.enemyType.Sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		spawnPos := config.Vector{
			X: e.position.X + halfW + math.Sin(e.rotation)*bulletSpawnOffset,
			Y: e.position.Y + halfH + math.Cos(e.rotation)*bulletSpawnOffset,
		}
		animation := NewAnimation(config.Vector{}, e.weapon.projectile.wType.IntercectAnimationSpriteSheet, 1, 56, 60, false, "projectileBlow", 0)
		projectile := NewProjectile(config.Vector{}, spawnPos, e.rotation, e.weapon.projectile.wType, animation, 0)
		projectile.owner = "enemy"
		e.game.AddProjectile(projectile)
	},
}

var enemyAutoLightRocket = Weapon{
	projectile: Projectile{
		wType: &config.WeaponType{
			Sprite:                        assets.EnemyAutoLightMissile,
			IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
			Velocity:                      3,
			Damage:                        5,
			TargetType:                    "auto",
		},
	},
	ammo: 3,
	EnemyShoot: func(e *Enemy) {
		bounds := e.enemyType.Sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		spawnPos := config.Vector{
			X: e.position.X + halfW + math.Sin(e.rotation)*bulletSpawnOffset,
			Y: e.position.Y + halfH + math.Cos(e.rotation)*bulletSpawnOffset,
		}
		animation := NewAnimation(config.Vector{}, e.weapon.projectile.wType.IntercectAnimationSpriteSheet, 1, 56, 60, false, "projectileBlow", 0)
		projectile := NewProjectile(config.Vector{}, spawnPos, e.rotation, e.weapon.projectile.wType, animation, 0)
		projectile.owner = "enemy"
		e.game.AddProjectile(projectile)
	},
}

func NewProjectile(target config.Vector, pos config.Vector, rotation float64, wType *config.WeaponType, animation *Animation, hp int) *Projectile {
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
		HP:                 hp,
	}

	return p
}

func (p *Projectile) AddAnimation(g *Game) {
	animation := NewAnimation(p.intercectAnimation.position, p.wType.IntercectAnimationSpriteSheet, p.intercectAnimation.speed, p.intercectAnimation.frameHeight, p.intercectAnimation.frameWidth, false, "projectileBlow", 0)
	g.animations = append(g.animations, animation)
}

func (p *Projectile) VelocityUpdate(player *Player) {
	switch p.wType.WeaponName {
	case config.LightRocket:
		p.wType.Velocity = 400 + player.params.LightRocketVelocityMultiplier
		player.curWeapon.shootCooldown.Restart(time.Millisecond * (300 - player.params.LightRocketSpeedUpscale))
	case config.DoubleLightRocket:
		p.wType.Velocity = 400 + player.params.DoubleLightRocketVelocityMultiplier
		player.curWeapon.shootCooldown.Restart(time.Millisecond * (250 - player.params.DoubleLightRocketSpeedUpscale))
	case config.LaserCannon:
		player.curWeapon.shootCooldown.Restart(time.Millisecond * (500 - player.params.LaserCannonSpeedUpscale))
	case config.DoubleLaserCannon:
		player.curWeapon.shootCooldown.Restart(time.Millisecond * (460 - player.params.DoubleLaserCannonSpeedUpscale))
	case config.MachineGun:
		p.wType.Velocity = 600 + player.params.MachineGunVelocityMultiplier
		player.curWeapon.shootCooldown.Restart(time.Millisecond * (160 - player.params.MachineGunSpeedUpscale))
	case config.DoubleMachineGun:
		p.wType.Velocity = 600 + player.params.DoubleMachineGunVelocityMultiplier
		player.curWeapon.shootCooldown.Restart(time.Millisecond * (260 - player.params.DoubleMachineGunSpeedUpscale))
	case config.PlasmaGun:
		p.wType.Velocity = 500 + player.params.PlasmaGunVelocityMultiplier
		player.curWeapon.shootCooldown.Restart(time.Millisecond * (760 - player.params.PlasmaGunSpeedUpscale))
	case config.DoublePlasmaGun:
		p.wType.Velocity = 500 + player.params.DoublePlasmaGunVelocityMultiplier
		player.curWeapon.shootCooldown.Restart(time.Millisecond * (880 - player.params.DoublePlasmaGunSpeedUpscale))
	case config.ClusterMines:
		p.wType.Velocity = 240 + player.params.ClusterMinesVelocityMultiplier
		player.curWeapon.shootCooldown.Restart(time.Millisecond * (400 - player.params.ClusterMinesSpeedUpscale))
	case config.BigBomb:
		p.wType.Velocity = 200 + player.params.BigBombVelocityMultiplier
		player.curWeapon.shootCooldown.Restart(time.Millisecond * (600 - player.params.BigBombSpeedUpscale))
	}
}
func (p *Projectile) Update() {
	speed := p.wType.Velocity / float64(ebiten.TPS())
	quant := speed
	if p.owner == "player" {
		quant = -speed
	}
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
		p.position.X += math.Sin(p.rotation) * speed
		p.position.Y += math.Cos(p.rotation) * quant
	}
	if p.instantAnimation != nil {
		p.instantAnimation.position.X = p.position.X
		p.instantAnimation.position.Y = p.position.Y
		p.instantAnimation.rotation = p.rotation
	}
}

func (p *Projectile) Draw(screen *ebiten.Image) {
	if !p.wType.AnimationOnly {
		objects.RotateAndTranslateObject(p.rotation, p.wType.Sprite, screen, p.position.X, p.position.Y)
	}
}

func (p *Projectile) Destroy(g *Game, i int) {
	if p.owner == "player" {
		g.projectiles = slices.Delete(g.projectiles, i, i+1)
	} else {
		g.enemyProjectiles = slices.Delete(g.enemyProjectiles, i, i+1)
	}
	if p.instantAnimation != nil {
		p.instantAnimation.looping = false
		p.instantAnimation.curTick = p.instantAnimation.speed - 1
		p.instantAnimation.currF = p.instantAnimation.numFrames - 1
	}
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
		Steps:    5,
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

func (b *Beam) Update() {
	if b.Step < b.Steps {
		b.Step++
	}
}

// func (b *Blow) Draw(screen *ebiten.Image) {
// 	vector.DrawFilledCircle(screen, float32(b.circle.X), float32(b.circle.Y), float32(b.circle.Radius), color.RGBA{255, 255, 255, 255}, false)
// }
