package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
)

const (
	rotationPerSecond = math.Pi
	maxAngle          = 256
	bulletSpawnOffset = 50.0
)

type Player struct {
	game *Game

	position            config.Vector
	rotation            float64
	sprite              *ebiten.Image
	objectRotationSpeed float64
	weapons             []*Weapon
	curWeapon           *Weapon
	animations          []*Animation
	hp                  int
}

func NewPlayer(curgame *Game) *Player {
	sprite := objects.ScaleImg(assets.PlayerSprite, 0.5)
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := config.Vector{
		X: config.ScreenWidth/2 - halfW,
		Y: config.ScreenHeight/2 - halfH,
	}
	posFireburst := config.Vector{
		X: config.ScreenWidth/2 - halfW,
		Y: config.ScreenHeight/2 + float64(bounds.Dy()/2),
	}
	startWeapon := NewWeapon(config.LightRocket)
	engineFireburst := NewAnimation(posFireburst, assets.PlayerFireburstSpriteSheet, 3, 30, 192, 96, true, "engineFireburst", 0)
	curgame.AddAnimation(engineFireburst)
	p := &Player{
		game:                curgame,
		position:            pos,
		rotation:            0,
		sprite:              sprite,
		objectRotationSpeed: 1.2,
		hp:                  10,
		weapons: []*Weapon{
			startWeapon,
		},
		animations: []*Animation{
			engineFireburst,
		},
	}
	p.curWeapon = p.weapons[0]
	return p
}

func (p *Player) Update() {
	x, y := ebiten.CursorPosition()
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.position.X -= 10
		if p.position.X < 0 {
			p.position.X = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.position.X += 10
		if p.position.X > config.ScreenWidth {
			p.position.X = config.ScreenWidth
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.position.Y -= 10
		if p.position.Y < 0 {
			p.position.Y = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.position.Y += 10
		if p.position.Y > config.ScreenHeight {
			p.position.Y = config.ScreenHeight
		}
	}

	targetRotation := math.Atan2(float64(y-int(p.position.Y)), float64(x-int(p.position.X)))
	// Calculate the rotation difference between the current and target rotation
	rotationDiff := math.Mod(targetRotation-p.rotation, 2*math.Pi)
	// Normalize the rotation difference to the range [-Pi, Pi]
	if rotationDiff > math.Pi {
		rotationDiff -= 2 * math.Pi
	} else if rotationDiff < -math.Pi {
		rotationDiff += 2 * math.Pi
	}
	// Calculate the rotation increment based on the rotation speed
	rotationIncrement := p.objectRotationSpeed * math.Abs(rotationDiff)
	// Apply the rotation increment to smoothly rotate the object towards the target rotation
	if rotationDiff > 0 {
		p.rotation += rotationIncrement
	} else if rotationDiff < 0 {
		p.rotation -= rotationIncrement
	}

	for _, a := range p.animations {
		switch a.name {
		case "engineFireburst":
			w := p.sprite.Bounds().Dx()
			h := p.sprite.Bounds().Dy()
			px, py := p.position.X+float64(w)/2, p.position.Y+float64(h)/2
			a.position.X = px + ((p.position.X-px)*math.Cos(-p.rotation) - (py-p.position.Y)*math.Sin(-p.rotation))
			a.position.Y = py - ((p.position.X-px)*math.Sin(-p.rotation) + (py-p.position.Y)*math.Cos(-p.rotation))
			a.rotation = p.rotation
			a.rotationPoint.X = float64(w)
			a.rotationPoint.Y = float64(h)
		default:
			a.position = p.position
			a.rotation = p.rotation
		}
	}
	if ebiten.IsKeyPressed(ebiten.Key1) {
		p.curWeapon = p.weapons[0]
	} else if ebiten.IsKeyPressed(ebiten.Key2) && len(p.weapons) > 1 {
		p.curWeapon = p.weapons[1]
	} else if ebiten.IsKeyPressed(ebiten.Key3) && len(p.weapons) > 2 {
		p.curWeapon = p.weapons[2]
	}

	p.curWeapon.shootCooldown.Update()
	if p.curWeapon.shootCooldown.IsReady() && (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) {
		if p.curWeapon.ammo <= 0 {
			return
		}
		p.curWeapon.shootCooldown.Reset()
		switch p.curWeapon.projectile.wType.WeaponName {
		case config.DoubleLightRocket:
			bounds := p.sprite.Bounds()
			halfWleft := float64(bounds.Dx()) / 4
			halfWright := float64(bounds.Dx()) - float64(bounds.Dx())/4
			halfH := float64(bounds.Dy()) / 2

			spawnPosLeft := config.Vector{
				X: p.position.X + halfWleft + math.Sin(p.rotation)*bulletSpawnOffset,
				Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
			}
			projectileLeft := NewProjectile(config.Vector{}, spawnPosLeft, p.rotation, p.curWeapon.projectile.wType)

			spawnPosRight := config.Vector{
				X: p.position.X + halfWright + math.Sin(p.rotation)*bulletSpawnOffset,
				Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
			}
			projectileRight := NewProjectile(config.Vector{}, spawnPosRight, p.rotation, p.curWeapon.projectile.wType)

			projectileLeft.owner = "player"
			projectileRight.owner = "player"
			p.game.AddProjectile(projectileLeft)
			p.game.AddProjectile(projectileRight)
		case config.LaserCannon:
			bounds := p.sprite.Bounds()
			halfW := float64(bounds.Dx()) / 2
			halfH := float64(bounds.Dy()) / 2

			spawnPos := config.Vector{
				X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
				Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
			}
			beam := NewBeam(config.Vector{X: float64(x), Y: float64(y)}, p.rotation, spawnPos, p.curWeapon.projectile.wType)
			beam.owner = "player"
			p.game.AddBeam(beam)
		default:
			bounds := p.sprite.Bounds()
			halfW := float64(bounds.Dx()) / 2
			halfH := float64(bounds.Dy()) / 2

			spawnPos := config.Vector{
				X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
				Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
			}

			projectile := NewProjectile(config.Vector{}, spawnPos, p.rotation, p.curWeapon.projectile.wType)
			projectile.owner = "player"
			p.game.AddProjectile(projectile)
		}
		p.curWeapon.ammo--
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	objects.RotateAndTranslateObject(p.rotation, p.sprite, screen, p.position.X, p.position.Y)
}

func (p *Player) Collider() image.Rectangle {
	bounds := p.sprite.Bounds()
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
