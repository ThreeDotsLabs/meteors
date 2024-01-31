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
	sprite             *ebiten.Image
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
	game     *Game
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
	UpdateParams  func(player *Player, w *Weapon)
	Shoot         func(p *Player)
	EnemyShoot    func(e *Enemy)
}

type Blow struct {
	circle config.Circle
	Damage int
	Steps  int
	Step   int
}

func NewEnemyWeapon(wType *config.WeaponType) Weapon {
	weapon := Weapon{
		projectile: Projectile{
			wType: wType,
			owner: config.OwnerEnemy,
		},
		shootCooldown: config.NewTimer(wType.StartTime),
		ammo:          wType.StartAmmo,
		EnemyShoot: func(e *Enemy) {
			bounds := e.enemyType.Sprite.Bounds()
			halfW := float64(bounds.Dx()) / 2
			halfH := float64(bounds.Dy()) / 2
			spawnPos := config.Vector{
				X: e.position.X + halfW + math.Sin(e.rotation)*bulletSpawnOffset,
				Y: e.position.Y + halfH + math.Cos(e.rotation)*bulletSpawnOffset,
			}
			animation := NewAnimation(config.Vector{}, e.weapon.projectile.wType.IntercectAnimationSpriteSheet, 1, 56, 60, false, "projectileBlow", 0)
			projectile := NewProjectile(e.game, spawnPos, e.rotation, e.weapon.projectile.wType, animation, 0)
			projectile.owner = config.OwnerEnemy
			e.game.AddProjectile(projectile)
		},
	}
	return weapon
}

func NewWeapon(wType string, p *Player) *Weapon {
	x, y := ebiten.CursorPosition()
	switch wType {
	case config.LightRocket:
		lightRType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.MissileSprite, 0.8),
			ItemSprite:                    objects.ScaleImg(assets.ItemMissileSprite, 0.5),
			IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
			Velocity:                      (400 + p.params.LightRocketVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod,
			Damage:                        int(3 * p.params.Ship.WeaponDamageMod),
			TargetType:                    config.TargetTypeStraight,
			WeaponName:                    config.LightRocket,
			Scale:                         p.game.Options.ResolutionMultipler,
			StartTime:                     time.Duration(p.params.Ship.WeaponFireRateMod * 300),
		}
		lightR := Weapon{
			projectile: Projectile{
				wType: lightRType,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.projectile.wType.Velocity = (400 + player.params.LightRocketVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod
				w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.LightRocketSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.LightRocketSpeedUpscale))
			},
			shootCooldown: config.NewTimer(time.Millisecond * (lightRType.StartTime - p.params.LightRocketSpeedUpscale)),
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
				projectile := NewProjectile(p.game, spawnPos, p.rotation, lightRType, animation, 0)
				projectile.owner = config.OwnerPlayer
				p.curWeapon.ammo--
				p.game.AddProjectile(projectile)
			},
		}
		return &lightR
	case config.DoubleLightRocket:
		boubleRType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.DoubleMissileSprite, 0.8),
			ItemSprite:                    objects.ScaleImg(assets.ItemDoubleMissileSprite, 0.5),
			IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
			Velocity:                      (400 + p.params.DoubleLightRocketVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod,
			Damage:                        int(3 * p.params.Ship.WeaponDamageMod),
			TargetType:                    config.TargetTypeStraight,
			WeaponName:                    config.DoubleLightRocket,
			StartTime:                     time.Duration(p.params.Ship.WeaponFireRateMod * 300),
		}
		doubleR := Weapon{
			projectile: Projectile{
				wType: boubleRType,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.projectile.wType.Velocity = (400 + player.params.DoubleLightRocketVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod
				w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.DoubleLightRocketSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.DoubleLightRocketSpeedUpscale))
			},
			shootCooldown: config.NewTimer(time.Millisecond * (boubleRType.StartTime - p.params.DoubleLightRocketSpeedUpscale)),
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
				projectileLeft := NewProjectile(p.game, spawnPosLeft, p.rotation, boubleRType, animationLeft, 0)
				animationRight := NewAnimation(config.Vector{}, boubleRType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileRight := NewProjectile(p.game, spawnPosRight, p.rotation, boubleRType, animationRight, 0)

				projectileLeft.owner = config.OwnerPlayer
				projectileRight.owner = config.OwnerPlayer
				p.game.AddProjectile(projectileLeft)
				p.game.AddProjectile(projectileRight)
				p.curWeapon.ammo--
			},
		}
		return &doubleR
	case config.LaserCanon:
		laserCType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.LaserCanon, 0.5*p.game.Options.ResolutionMultipler),
			ItemSprite:                    objects.ScaleImg(assets.ItemLaserCanonSprite, 0.5),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Damage:                        3,
			TargetType:                    config.TargetTypeStraight,
			WeaponName:                    config.LaserCanon,
			StartTime:                     time.Duration(p.params.Ship.WeaponFireRateMod * 500),
		}
		laserC := Weapon{
			projectile: Projectile{
				wType: laserCType,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.LaserCanonSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.LaserCanonSpeedUpscale))
			},
			shootCooldown: config.NewTimer(time.Millisecond * (laserCType.StartTime - p.params.LaserCanonSpeedUpscale)),
			ammo:          40,
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
				beam := NewBeam(config.Vector{X: float64(x), Y: float64(y)}, p.rotation, spawnPos, laserCType, p.game)
				beam.owner = config.OwnerPlayer
				p.game.AddBeam(beam)
				ba := beam.NewBeamAnimation()
				p.game.AddBeamAnimation(ba)
				p.curWeapon.ammo--
			},
		}
		return &laserC
	case config.DoubleLaserCanon:
		doubleLaserCType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.LaserCanon, 0.5*p.game.Options.ResolutionMultipler),
			ItemSprite:                    objects.ScaleImg(assets.ItemDoubleLaserCanonSprite, 0.5),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Damage:                        int(2 * p.params.Ship.WeaponDamageMod),
			TargetType:                    config.TargetTypeStraight,
			WeaponName:                    config.DoubleLaserCanon,
			StartTime:                     time.Duration(p.params.Ship.WeaponFireRateMod * 460),
		}
		doubleLaserC := Weapon{
			projectile: Projectile{
				wType: doubleLaserCType,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.DoubleLaserCanonSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.DoubleLaserCanonSpeedUpscale))
			},
			shootCooldown: config.NewTimer(time.Millisecond * (doubleLaserCType.StartTime - p.params.DoubleLaserCanonSpeedUpscale)),
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
				beamLeft := NewBeam(config.Vector{X: float64(x), Y: float64(y)}, p.rotation, spawnPosLeft, doubleLaserCType, p.game)
				beamLeft.owner = config.OwnerPlayer
				p.game.AddBeam(beamLeft)
				baL := beamLeft.NewBeamAnimation()
				p.game.AddBeamAnimation(baL)
				beamRight := NewBeam(config.Vector{X: float64(x), Y: float64(y)}, p.rotation, spawnPosRight, doubleLaserCType, p.game)
				beamRight.owner = config.OwnerPlayer
				p.game.AddBeam(beamRight)
				baR := beamRight.NewBeamAnimation()
				p.game.AddBeamAnimation(baR)
				p.curWeapon.ammo--
			},
		}
		return &doubleLaserC
	case config.ClusterMines:
		clusterMType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.ClusterMines, 0.5*p.game.Options.ResolutionMultipler),
			IntercectAnimationSpriteSheet: assets.ClusterMinesBlowSpriteSheet,
			Velocity:                      (360 + p.params.ClusterMinesVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod,
			Damage:                        int(3 * p.params.Ship.WeaponDamageMod),
			TargetType:                    config.TargetTypeStraight,
			WeaponName:                    config.ClusterMines,
			StartTime:                     time.Duration(p.params.Ship.WeaponFireRateMod * 400),
		}
		clusterM := Weapon{
			projectile: Projectile{
				wType: clusterMType,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.projectile.wType.Velocity = (360 + player.params.ClusterMinesVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod
				w.shootCooldown.Restart(time.Millisecond * (clusterMType.StartTime - player.params.ClusterMinesSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (400 - player.params.ClusterMinesSpeedUpscale))

			},
			shootCooldown: config.NewTimer(time.Millisecond * (clusterMType.StartTime - p.params.ClusterMinesSpeedUpscale)),
			ammo:          25,
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
					projectile := NewProjectile(p.game, spawnPos, angle, clusterMType, animation, 0)
					projectile.owner = config.OwnerPlayer
					p.game.AddProjectile(projectile)
					p.curSecondaryWeapon.ammo--
				}
			},
		}
		return &clusterM
	case config.BigBomb:
		bigBType := &config.WeaponType{
			Sprite:     objects.ScaleImg(assets.BigBomb, 0.8*p.game.Options.ResolutionMultipler),
			Velocity:   (200 + p.params.BigBombVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod,
			Damage:     int(10 * p.params.Ship.WeaponDamageMod),
			TargetType: config.TargetTypeStraight,
			WeaponName: config.BigBomb,
			StartTime:  time.Duration(p.params.Ship.WeaponFireRateMod * 600),
		}
		bigB := Weapon{
			projectile: Projectile{
				wType: bigBType,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.projectile.wType.Velocity = (200 + player.params.BigBombVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod
				w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.BigBombSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (600 - player.params.BigBombSpeedUpscale))
			},
			shootCooldown: config.NewTimer(time.Millisecond * (bigBType.StartTime - p.params.BigBombSpeedUpscale)),
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
				projectile := NewProjectile(p.game, spawnPos, p.rotation, bigBType, animation, 0)
				projectile.owner = config.OwnerPlayer
				p.game.AddProjectile(projectile)
				p.curSecondaryWeapon.ammo--
			},
		}
		return &bigB
	case config.MachineGun:
		machineGType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.MachineGun, p.game.Options.ResolutionMultipler),
			ItemSprite:                    objects.ScaleImg(assets.ItemMachineGunSprite, p.game.Options.ResolutionMultipler),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Velocity:                      (850 + p.params.MachineGunVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod,
			Damage:                        1,
			TargetType:                    config.TargetTypeStraight,
			WeaponName:                    config.MachineGun,
			StartTime:                     time.Duration(p.params.Ship.WeaponFireRateMod * 160),
		}
		machineG := Weapon{
			projectile: Projectile{
				wType: machineGType,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.projectile.wType.Velocity = (850 + player.params.MachineGunVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod
				w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.MachineGunSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (160 - player.params.MachineGunSpeedUpscale))
			},
			shootCooldown: config.NewTimer(time.Millisecond * (machineGType.StartTime - p.params.MachineGunSpeedUpscale)),
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
				projectile := NewProjectile(p.game, spawnPos, p.rotation, machineGType, animation, 0)
				projectile.owner = config.OwnerPlayer
				p.game.AddProjectile(projectile)
				p.curWeapon.ammo--
			},
		}
		return &machineG
	case config.DoubleMachineGun:
		doubleMachineGType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.DoubleMachineGun, p.game.Options.ResolutionMultipler),
			ItemSprite:                    objects.ScaleImg(assets.ItemDoubleMachineGunSprite, p.game.Options.ResolutionMultipler),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			Velocity:                      (850 + p.params.DoubleMachineGunVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod,
			Damage:                        1,
			TargetType:                    config.TargetTypeStraight,
			WeaponName:                    config.MachineGun,
			StartTime:                     time.Duration(p.params.Ship.WeaponFireRateMod * 260),
		}
		doubleMachineG := Weapon{
			projectile: Projectile{
				wType: doubleMachineGType,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.projectile.wType.Velocity = (850 + player.params.DoubleMachineGunVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod
				w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.DoubleMachineGunSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (260 - player.params.DoubleMachineGunSpeedUpscale))
			},
			shootCooldown: config.NewTimer(time.Millisecond * (doubleMachineGType.StartTime - p.params.DoubleMachineGunSpeedUpscale)),
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
				projectileLeft := NewProjectile(p.game, spawnPosLeft, p.rotation, doubleMachineGType, animationLeft, 0)
				animationRight := NewAnimation(config.Vector{}, doubleMachineGType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileRight := NewProjectile(p.game, spawnPosRight, p.rotation, doubleMachineGType, animationRight, 0)
				projectileLeft.owner = config.OwnerPlayer
				projectileRight.owner = config.OwnerPlayer
				p.game.AddProjectile(projectileLeft)
				p.game.AddProjectile(projectileRight)
				p.curWeapon.ammo--
			},
		}
		return &doubleMachineG
	case config.PlasmaGun:
		plasmaGType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.PlasmaGun, 0.8*p.game.Options.ResolutionMultipler),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			InstantAnimationSpiteSheet:    assets.PlasmaGunProjectileSpriteSheet,
			Velocity:                      (500 + p.params.PlasmaGunVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod,
			AnimationOnly:                 true,
			Damage:                        4,
			TargetType:                    config.TargetTypeStraight,
			WeaponName:                    config.PlasmaGun,
			StartTime:                     time.Duration(p.params.Ship.WeaponFireRateMod * 560),
		}

		plasmaG := Weapon{
			projectile: Projectile{
				wType: plasmaGType,
				HP:    4,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.projectile.wType.Velocity = (500 + player.params.PlasmaGunVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod
				w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.PlasmaGunSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (760 - player.params.PlasmaGunSpeedUpscale))
			},
			shootCooldown: config.NewTimer(time.Millisecond * (plasmaGType.StartTime - p.params.PlasmaGunSpeedUpscale)),
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
				projectile := NewProjectile(p.game, spawnPos, p.rotation, plasmaGType, animation, 4)
				projectile.owner = config.OwnerPlayer
				projectile.instantAnimation = plasmaAnimation
				p.game.AddProjectile(projectile)
				p.game.AddAnimation(plasmaAnimation)
				p.curWeapon.ammo--
			},
		}
		return &plasmaG
	case config.DoublePlasmaGun:
		doublePlasmaGType := &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.PlasmaGun, 1.2*p.game.Options.ResolutionMultipler),
			IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
			InstantAnimationSpiteSheet:    assets.PlasmaGunProjectileSpriteSheet,
			Velocity:                      (500 + p.params.DoublePlasmaGunVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod,
			AnimationOnly:                 true,
			Damage:                        4,
			TargetType:                    config.TargetTypeStraight,
			WeaponName:                    config.DoublePlasmaGun,
			StartTime:                     time.Duration(p.params.Ship.WeaponFireRateMod * 620),
		}
		doublePlasmaG := Weapon{
			projectile: Projectile{
				wType: doublePlasmaGType,
				HP:    4,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.projectile.wType.Velocity = (500 + player.params.DoublePlasmaGunVelocityMultiplier) * p.params.Ship.WeaponProjectileVelocityMod
				w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.DoublePlasmaGunSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (760 - player.params.DoublePlasmaGunSpeedUpscale))
			},
			shootCooldown: config.NewTimer(time.Millisecond * (doublePlasmaGType.StartTime - p.params.DoublePlasmaGunSpeedUpscale)),
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
				projectileLeft := NewProjectile(p.game, spawnPosLeft, p.rotation, doublePlasmaGType, animationLeft, 4)
				plasmaAnimationRight := NewAnimation(config.Vector{}, doublePlasmaGType.InstantAnimationSpiteSheet, 1, 55, 50, true, "projectileInstant", 0)
				animationRight := NewAnimation(config.Vector{}, doublePlasmaGType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
				projectileRight := NewProjectile(p.game, spawnPosRight, p.rotation, doublePlasmaGType, animationRight, 4)
				projectileLeft.owner = config.OwnerPlayer
				projectileRight.owner = config.OwnerPlayer
				projectileLeft.instantAnimation = plasmaAnimationLeft
				projectileRight.instantAnimation = plasmaAnimationRight
				p.game.AddProjectile(projectileLeft)
				p.game.AddProjectile(projectileRight)
				p.game.AddAnimation(plasmaAnimationLeft)
				p.game.AddAnimation(plasmaAnimationRight)
				p.curWeapon.ammo--
			},
		}
		return &doublePlasmaG
	case config.PentaLaser:
		pentaLaserType := &config.WeaponType{
			Sprite:     objects.ScaleImg(assets.PentaLaser, 0.8*p.game.Options.ResolutionMultipler),
			Damage:     int(8 * p.params.Ship.WeaponDamageMod),
			TargetType: config.TargetTypeStraight,
			WeaponName: config.PentaLaser,
			StartTime:  time.Duration(p.params.Ship.WeaponFireRateMod * 1000),
		}
		pentaLaser := Weapon{
			projectile: Projectile{
				wType: pentaLaserType,
			},
			UpdateParams: func(player *Player, w *Weapon) {
				w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.PentaLaserSpeedUpscale))
				//player.curWeapon.shootCooldown.Restart(time.Millisecond * (1000 - player.params.PentaLaserSpeedUpscale))
			},
			shootCooldown: config.NewTimer(time.Millisecond * (pentaLaserType.StartTime - p.params.PentaLaserSpeedUpscale)),
			ammo:          50,
			Shoot: func(p *Player) {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				halfH := float64(bounds.Dy()) / 2
				for i := 0; i < 5; i++ {
					angle := p.rotation - (math.Pi/180)*float64(i*24) + (math.Pi/180)*float64(45)
					spawnPos := config.Vector{
						X: p.position.X + halfW,
						Y: p.position.Y + halfH - halfW/2,
					}
					beam := NewBeam(config.Vector{X: float64(x), Y: float64(y)}, angle, spawnPos, pentaLaserType, p.game)
					beam.owner = config.OwnerPlayer
					p.game.AddBeam(beam)
					ba := beam.NewBeamAnimation()
					p.game.AddBeamAnimation(ba)
					p.curSecondaryWeapon.ammo--
				}
			},
		}
		return &pentaLaser
	}
	return nil
}

func NewProjectile(g *Game, pos config.Vector, rotation float64, wType *config.WeaponType, animation *Animation, hp int) *Projectile {
	bounds := wType.Sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH
	spriteImg := ebiten.NewImageFromImage(wType.Sprite)
	sprite := objects.ScaleImg(spriteImg, g.Options.ProjectileResMulti-0.2)
	p := &Projectile{
		sprite:             sprite,
		position:           pos,
		rotation:           rotation,
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

func (p *Projectile) Update() {
	speed := p.wType.Velocity / float64(ebiten.TPS())
	if p.wType.TargetType == config.TargetTypePlayer {
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
	} else if p.wType.TargetType == config.TargetTypeStraight {
		if p.owner == config.OwnerPlayer {
			p.position.X += math.Sin(p.rotation) * speed
			p.position.Y -= math.Cos(p.rotation) * speed
		} else {
			p.position.X += math.Sin(-p.rotation) * speed
			p.position.Y += math.Cos(p.rotation) * speed
		}
	}
	if p.instantAnimation != nil {
		p.instantAnimation.position.X = p.position.X
		p.instantAnimation.position.Y = p.position.Y
		p.instantAnimation.rotation = p.rotation
	}
}

func (p *Projectile) Draw(screen *ebiten.Image) {
	if !p.wType.AnimationOnly {
		objects.RotateAndTranslateObject(p.rotation, p.sprite, screen, p.position.X, p.position.Y)
	}
}

func (p *Projectile) Destroy(g *Game, i int) {
	if p.owner == config.OwnerPlayer {
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

func NewBeam(target config.Vector, rotation float64, pos config.Vector, wType *config.WeaponType, g *Game) *Beam {
	screenDiag := math.Sqrt(g.Options.ScreenWidth*g.Options.ScreenWidth + g.Options.ScreenHeight*g.Options.ScreenHeight)
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
		game:     g,
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
	screenDiag := math.Sqrt(b.game.Options.ScreenWidth*b.game.Options.ScreenWidth + b.game.Options.ScreenHeight*b.game.Options.ScreenHeight)
	rect := config.NewRectangle(
		b.position.X,
		b.position.Y,
		float64(1),
		float64(screenDiag),
	)
	return &BeamAnimation{
		curRect:  rect,
		Steps:    5,
		Step:     0,
		rotation: b.rotation + math.Pi,
	}
}

func (b *BeamAnimation) Update() {
	b.curRect.Width += float64(b.Step)
	b.curRect.X += 1
	b.Step++
}

func (b *BeamAnimation) Draw(screen *ebiten.Image) {
	rectImage := ebiten.NewImage(int(b.curRect.Width), int(b.curRect.Height))
	rectImage.Fill(color.White)
	rotationOpts := &ebiten.DrawImageOptions{}
	rotationOpts.GeoM.Rotate(b.rotation)
	rotationOpts.GeoM.Translate(b.curRect.X, b.curRect.Y)
	color := color.RGBA{255 - uint8(b.Step)*10, 255, 0, 0 - uint8(b.Step)*25}
	rotationOpts.ColorScale.ScaleWithColor(color)
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
