package game

import (
	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Ship struct {
	Name                        string
	Sprite                      *ebiten.Image
	HP                          int
	Velocity                    float64
	HPMod                       float64
	VelocityMod                 float64
	WeaponFireRateMod           float64
	WeaponDamageMod             float64
	WeaponProjectileVelocityMod float64
	UniqueWeapon                *Weapon
}

var AngryOcelot = &Ship{
	Name:                        "Angry ocelot",
	Sprite:                      assets.ShipAngryOcelotSprite,
	HP:                          5,
	Velocity:                    1.0,
	WeaponFireRateMod:           1.0,
	WeaponDamageMod:             1.0,
	WeaponProjectileVelocityMod: 1.0,
	HPMod:                       1.0,
	VelocityMod:                 1.0,
}

var MightyOrca = &Ship{
	Name:                        "Mighty orca",
	Sprite:                      assets.ShipMightyOrcaSprite,
	HP:                          20,
	Velocity:                    -4.0,
	WeaponFireRateMod:           1.2,
	WeaponDamageMod:             1.0,
	WeaponProjectileVelocityMod: 1.0,
	HPMod:                       1.6,
	VelocityMod:                 0.7,
}

var ShadyWeasel = &Ship{
	Name:                        "Shady weasel",
	Sprite:                      assets.ShipShadyWeaselSprite,
	HP:                          0,
	Velocity:                    4.0,
	WeaponFireRateMod:           1.2,
	WeaponDamageMod:             1,
	WeaponProjectileVelocityMod: 1.5,
	HPMod:                       0.7,
	VelocityMod:                 1.5,
}

func NewAngryOcelotWeapon(p *Player) *Weapon {
	weaponType := &config.WeaponType{
		Sprite:                        objects.ScaleImg(assets.MissileSprite, 0.8),
		ItemSprite:                    objects.ScaleImg(assets.ItemMissileSprite, 0.5),
		IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
		Velocity:                      600,
		Damage:                        3,
		TargetType:                    config.TargetTypeStraight,
		WeaponName:                    "angryOcelotWeapon",
		Scale:                         p.game.Options.ResolutionMultipler,
		StartTime:                     250,
		StartAmmo:                     120,
	}
	weapon := Weapon{
		projectile: Projectile{
			wType: weaponType,
		},
		UpdateParams: func(player *Player, w *Weapon) {
		},
		shootCooldown: config.NewTimer(time.Millisecond * weaponType.StartTime),
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
			animation := NewAnimation(config.Vector{}, weaponType.IntercectAnimationSpriteSheet, 1, 56, 60, false, "projectileBlow", 0)
			projectile := NewProjectile(p.game, spawnPos, p.rotation, weaponType, animation, 0)
			projectile.owner = config.OwnerPlayer
			p.game.AddProjectile(projectile)
			p.curWeapon.ammo--
		},
	}
	return &weapon
}

func NewMightyOrcaWeapon(p *Player) *Weapon {
	weaponType := &config.WeaponType{
		Sprite:                        objects.ScaleImg(assets.MachineGun, p.game.Options.ResolutionMultipler),
		ItemSprite:                    objects.ScaleImg(assets.ItemMachineGunSprite, p.game.Options.ResolutionMultipler),
		IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
		Velocity:                      900,
		Damage:                        1,
		TargetType:                    config.TargetTypeStraight,
		WeaponName:                    "mightyOrcaWeapon",
		StartTime:                     180,
		StartAmmo:                     180,
	}
	weapon := Weapon{
		projectile: Projectile{
			wType: weaponType,
		},
		UpdateParams: func(player *Player, w *Weapon) {
		},
		shootCooldown: config.NewTimer(time.Millisecond * weaponType.StartTime),
		Shoot: func(p *Player) {
			bounds := p.sprite.Bounds()
			halfW := float64(bounds.Dx()) / 2
			halfH := float64(bounds.Dy()) / 2

			spawnPos := config.Vector{
				X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
				Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
			}
			animation := NewAnimation(config.Vector{}, weaponType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
			projectile := NewProjectile(p.game, spawnPos, p.rotation, weaponType, animation, 0)
			projectile.owner = config.OwnerPlayer
			p.game.AddProjectile(projectile)
			p.curWeapon.ammo--
		},
	}
	return &weapon
}

func NewShadyWeaselWeapon(p *Player) *Weapon {
	weaponType := &config.WeaponType{
		Sprite:                        objects.ScaleImg(assets.PlasmaGun, 0.8*p.game.Options.ResolutionMultipler),
		ItemSprite:                    objects.ScaleImg(assets.ItemPlasmaGunSprite, p.game.Options.ResolutionMultipler),
		IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
		InstantAnimationSpiteSheet:    assets.PlasmaGunProjectileSpriteSheet,
		Velocity:                      600,
		AnimationOnly:                 true,
		Damage:                        3,
		TargetType:                    config.TargetTypeStraight,
		WeaponName:                    "shadyWeaselWeapon",
		StartTime:                     480,
		StartAmmo:                     80,
	}

	weapon := Weapon{
		projectile: Projectile{
			wType: weaponType,
		},
		UpdateParams: func(player *Player, w *Weapon) {
			// w.projectile.wType.Velocity = (600 + player.params.PlasmaGunVelocityMultiplier)
			// w.shootCooldown.Restart(time.Millisecond * (w.projectile.wType.StartTime - player.params.PlasmaGunSpeedUpscale))
		},
		shootCooldown: config.NewTimer(time.Millisecond * (weaponType.StartTime - p.params.PlasmaGunSpeedUpscale)),
		Shoot: func(p *Player) {
			bounds := p.sprite.Bounds()
			halfW := float64(bounds.Dx()) / 2
			halfH := float64(bounds.Dy()) / 2

			spawnPos := config.Vector{
				X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
				Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
			}
			plasmaAnimation := NewAnimation(config.Vector{}, weaponType.InstantAnimationSpiteSheet, 1, 55, 50, true, "projectileInstant", 0)
			animation := NewAnimation(config.Vector{}, weaponType.IntercectAnimationSpriteSheet, 1, 40, 40, false, "projectileBlow", 0)
			projectile := NewProjectile(p.game, spawnPos, p.rotation, weaponType, animation, 3)
			projectile.owner = config.OwnerPlayer
			projectile.instantAnimation = plasmaAnimation
			p.game.AddProjectile(projectile)
			p.game.AddAnimation(plasmaAnimation)
			p.curWeapon.ammo--
		},
	}
	return &weapon
}
