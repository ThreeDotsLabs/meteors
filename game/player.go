package game

import (
	"image"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
)

const (
	rotationPerSecond = math.Pi
	maxAngle          = 256
	bulletSpawnOffset = 50.0
)

type PlayerParams struct {
	Ship  *Ship
	Level int
	HP    int
	speed float64

	LightRocketSpeedUpscale       time.Duration
	AutoLightRocketSpeedUpscale   time.Duration
	DoubleLightRocketSpeedUpscale time.Duration
	LaserCanonSpeedUpscale        time.Duration
	DoubleLaserCanonSpeedUpscale  time.Duration
	ClusterMinesSpeedUpscale      time.Duration
	BigBombSpeedUpscale           time.Duration
	MachineGunSpeedUpscale        time.Duration
	DoubleMachineGunSpeedUpscale  time.Duration
	PlasmaGunSpeedUpscale         time.Duration
	DoublePlasmaGunSpeedUpscale   time.Duration
	PentaLaserSpeedUpscale        time.Duration

	LightRocketVelocityMultiplier       float64
	AutoLightRocketVelocityMultiplier   float64
	DoubleLightRocketVelocityMultiplier float64
	ClusterMinesVelocityMultiplier      float64
	BigBombVelocityMultiplier           float64
	MachineGunVelocityMultiplier        float64
	DoubleMachineGunVelocityMultiplier  float64
	PlasmaGunVelocityMultiplier         float64
	DoublePlasmaGunVelocityMultiplier   float64
}

func (p *PlayerParams) GetHealthPoints() int {
	return p.HP
}
func (p *PlayerParams) IncreaseHP(i int) {
	p.HP += i
}

func (p *PlayerParams) GetSpeed() int {
	return int(p.speed)
}
func (p *PlayerParams) IncreaseSpeed(i int) {
	p.speed += float64(i)
}

func (p *PlayerParams) GetLightRocketSpeedUpscale() int {
	return int(p.LightRocketSpeedUpscale)
}

func (p *PlayerParams) IncreaseLightRocketSpeedUpscale(t int) {
	if t > 0 {
		p.LightRocketSpeedUpscale++
	} else {
		p.LightRocketSpeedUpscale--
	}
}

func (p *PlayerParams) GetLightRocketVelocityMultiplier() int {
	return int(p.LightRocketVelocityMultiplier)
}

func (p *PlayerParams) IncreaseLightRocketVelocityMultiplier(t int) {
	p.LightRocketVelocityMultiplier += float64(t)
}

func (p *PlayerParams) GetAutoLightRocketSpeedUpscale() int {
	return int(p.AutoLightRocketSpeedUpscale)
}

func (p *PlayerParams) IncreaseAutoLightRocketSpeedUpscale(t int) {
	if t > 0 {
		p.AutoLightRocketSpeedUpscale++
	} else {
		p.AutoLightRocketSpeedUpscale--
	}
}

func (p *PlayerParams) GetAutoLightRocketVelocityMultiplier() int {
	return int(p.AutoLightRocketVelocityMultiplier)
}

func (p *PlayerParams) IncreaseAutoLightRocketVelocityMultiplier(t int) {
	p.AutoLightRocketVelocityMultiplier += float64(t)
}

func (p *PlayerParams) GetDoubleLightRocketSpeedUpscale() int {
	return int(p.DoubleLightRocketSpeedUpscale)
}

func (p *PlayerParams) IncreaseDoubleLightRocketSpeedUpscale(t int) {
	if t > 0 {
		p.DoubleLightRocketSpeedUpscale++
	} else {
		p.DoubleLightRocketSpeedUpscale--
	}
}

func (p *PlayerParams) GetDoubleLightRocketVelocityMultiplier() int {
	return int(p.DoubleLightRocketVelocityMultiplier)
}

func (p *PlayerParams) IncreaseDoubleLightRocketVelocityMultiplier(t int) {
	p.DoubleLightRocketVelocityMultiplier += float64(t)
}

func (p *PlayerParams) GetClusterMinesSpeedUpscale() int {
	return int(p.ClusterMinesSpeedUpscale)
}

func (p *PlayerParams) IncreaseClusterMinesSpeedUpscale(t int) {
	if t > 0 {
		p.ClusterMinesSpeedUpscale++
	} else {
		p.ClusterMinesSpeedUpscale--
	}
}

func (p *PlayerParams) GetClusterMinesVelocityMultiplier() int {
	return int(p.ClusterMinesVelocityMultiplier)
}

func (p *PlayerParams) IncreaseClusterMinesVelocityMultiplier(t int) {
	p.ClusterMinesVelocityMultiplier += float64(t)
}

func (p *PlayerParams) GetPentaLaserSpeedUpscale() int {
	return int(p.PentaLaserSpeedUpscale)
}

func (p *PlayerParams) IncreasePentaLaserSpeedUpscale(t int) {
	if t > 0 {
		p.PentaLaserSpeedUpscale++
	} else {
		p.PentaLaserSpeedUpscale--
	}
}

func (p *PlayerParams) GetLaserCanonSpeedUpscale() int {
	return int(p.LaserCanonSpeedUpscale)
}

func (p *PlayerParams) IncreaseLaserCanonSpeedUpscale(t int) {
	if t > 0 {
		p.LaserCanonSpeedUpscale++
	} else {
		p.LaserCanonSpeedUpscale--
	}
}

func (p *PlayerParams) GetDoubleLaserCanonSpeedUpscale() int {
	return int(p.DoubleLaserCanonSpeedUpscale)
}

func (p *PlayerParams) IncreaseDoubleLaserCanonSpeedUpscale(t int) {
	if t > 0 {
		p.DoubleLaserCanonSpeedUpscale++
	} else {
		p.DoubleLaserCanonSpeedUpscale--
	}
}

func (p *PlayerParams) GetBigBombSpeedUpscale() int {
	return int(p.BigBombSpeedUpscale)
}

func (p *PlayerParams) IncreaseBigBombSpeedUpscale(t int) {
	if t > 0 {
		p.BigBombSpeedUpscale++
	} else {
		p.BigBombSpeedUpscale--
	}
}

func (p *PlayerParams) GetBigBombVelocityMultiplier() int {
	return int(p.BigBombVelocityMultiplier)
}

func (p *PlayerParams) IncreaseBigBombVelocityMultiplier(t int) {
	p.BigBombVelocityMultiplier += float64(t)
}

func (p *PlayerParams) GetMachineGunSpeedUpscale() int {
	return int(p.MachineGunSpeedUpscale)
}

func (p *PlayerParams) IncreaseMachineGunSpeedUpscale(t int) {
	if t > 0 {
		p.MachineGunSpeedUpscale++
	} else {
		p.MachineGunSpeedUpscale--
	}
}

func (p *PlayerParams) GetMachineGunVelocityMultiplier() int {
	return int(p.MachineGunVelocityMultiplier)
}

func (p *PlayerParams) IncreaseMachineGunVelocityMultiplier(t int) {
	p.MachineGunVelocityMultiplier += float64(t)
}

func (p *PlayerParams) GetDoubleMachineGunSpeedUpscale() int {
	return int(p.DoubleMachineGunSpeedUpscale)
}

func (p *PlayerParams) IncreaseDoubleMachineGunSpeedUpscale(t int) {
	if t > 0 {
		p.DoubleMachineGunSpeedUpscale++
	} else {
		p.DoubleMachineGunSpeedUpscale--
	}
}

func (p *PlayerParams) GetDoubleMachineGunVelocityMultiplier() int {
	return int(p.DoubleMachineGunVelocityMultiplier)
}

func (p *PlayerParams) IncreaseDoubleMachineGunVelocityMultiplier(t int) {
	p.DoubleMachineGunVelocityMultiplier += float64(t)
}

func (p *PlayerParams) GetPlasmaGunVelocityMultiplier() int {
	return int(p.PlasmaGunVelocityMultiplier)
}

func (p *PlayerParams) IncreasePlasmaGunVelocityMultiplier(t int) {
	p.PlasmaGunVelocityMultiplier += float64(t)
}

func (p *PlayerParams) GetPlasmaGunSpeedUpscale() int {
	return int(p.PlasmaGunSpeedUpscale)
}

func (p *PlayerParams) IncreasePlasmaGunSpeedUpscale(t int) {
	if t > 0 {
		p.PlasmaGunSpeedUpscale++
	} else {
		p.PlasmaGunSpeedUpscale--
	}
}

func (p *PlayerParams) GetDoublePlasmaGunSpeedUpscale() int {
	return int(p.DoublePlasmaGunSpeedUpscale)
}

func (p *PlayerParams) IncreaseDoublePlasmaGunSpeedUpscale(t int) {
	if t > 0 {
		p.DoublePlasmaGunSpeedUpscale++
	} else {
		p.DoublePlasmaGunSpeedUpscale--
	}
}

func (p *PlayerParams) GetDoublePlasmaGunVelocityMultiplier() int {
	return int(p.DoublePlasmaGunVelocityMultiplier)
}

func (p *PlayerParams) IncreaseDoublePlasmaGunVelocityMultiplier(t int) {
	p.DoublePlasmaGunVelocityMultiplier += float64(t)
}

type Player struct {
	game                *Game
	params              *PlayerParams
	position            config.Vector
	rotation            float64
	sprite              *ebiten.Image
	objectRotationSpeed float64
	weapons             []*Weapon
	curWeapon           *Weapon
	secondaryWeapons    []*Weapon
	curSecondaryWeapon  *Weapon
	animations          []*Animation
	shield              *Shield
}

func (p *Player) SetShip(s *Ship) {
	p.params.HP += s.HP
	p.params.speed += s.Velocity
	p.sprite = objects.ScaleImg(s.Sprite, p.game.Options.ResolutionMultipler)
	p.weapons = append(p.weapons, s.UniqueWeapon)
	p.curWeapon = s.UniqueWeapon
	p.curWeapon.ammo = s.UniqueWeapon.projectile.wType.StartAmmo
	p.params.Ship = s
}

type Shield struct {
	position config.Vector
	sprite   *ebiten.Image
	HP       int
}

func NewPlayer(curgame *Game) *Player {
	sprite := objects.ScaleImg(assets.PlayerSprite, curgame.Options.ResolutionMultipler)
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := config.Vector{
		X: curgame.Options.ScreenWidth/2 - halfW,
		Y: curgame.Options.ScreenHeight/2 - halfH,
	}
	posFireburst := config.Vector{
		X: curgame.Options.ScreenWidth/2 - halfW,
		Y: curgame.Options.ScreenHeight/2 + float64(bounds.Dy()/2),
	}

	engineFireburst := NewAnimation(posFireburst, assets.PlayerFireburstSpriteSheet, 1, 90, 85, true, "engineFireburst", 0)
	curgame.AddAnimation(engineFireburst)
	p := &Player{
		params: &PlayerParams{
			Level: 1,
			HP:    10,
			speed: 10,
			Ship: &Ship{
				HP:                0,
				Velocity:          0,
				HPMod:             0,
				VelocityMod:       1,
				WeaponFireRateMod: 1,
				WeaponDamageMod:   1,
			},
			LightRocketVelocityMultiplier:       1,
			AutoLightRocketVelocityMultiplier:   1,
			DoubleLightRocketVelocityMultiplier: 1,
			ClusterMinesVelocityMultiplier:      1,
			BigBombVelocityMultiplier:           1,
			MachineGunVelocityMultiplier:        1,
			PlasmaGunVelocityMultiplier:         1,
			DoublePlasmaGunVelocityMultiplier:   1,
		},
		game:                curgame,
		position:            pos,
		rotation:            0,
		sprite:              sprite,
		objectRotationSpeed: 1.2,
		animations: []*Animation{
			engineFireburst,
		},
	}
	return p
}

func (p *Player) Update() {
	if p.game.ResolutionChange {
		p.sprite = objects.ScaleImg(assets.PlayerSprite, p.game.Options.ResolutionMultipler)
		//p.game.ResolutionChange = false
	}

	x, y := ebiten.CursorPosition()
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.position.X -= p.params.speed
		if p.position.X < 0 {
			p.position.X = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.position.X += p.params.speed
		if p.position.X > p.game.Options.ScreenWidth {
			p.position.X = p.game.Options.ScreenWidth
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.position.Y -= p.params.speed
		if p.position.Y < 0 {
			p.position.Y = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.position.Y += p.params.speed
		if p.position.Y > p.game.Options.ScreenHeight {
			p.position.Y = p.game.Options.ScreenHeight
		}
	}

	if p.shield != nil {
		p.shield.position.X = p.position.X
		p.shield.position.Y = p.position.Y
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
		case "shield":
			w := p.sprite.Bounds().Dx()
			h := p.sprite.Bounds().Dy()
			px, py := p.position.X-float64(w)/2, p.position.Y-float64(h)/2
			a.position.X = px
			a.position.Y = py
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
	} else if ebiten.IsKeyPressed(ebiten.Key4) && len(p.weapons) > 3 {
		p.curWeapon = p.weapons[3]
	} else if ebiten.IsKeyPressed(ebiten.Key5) && len(p.weapons) > 4 {
		p.curWeapon = p.weapons[4]
	} else if ebiten.IsKeyPressed(ebiten.Key6) && len(p.weapons) > 5 {
		p.curWeapon = p.weapons[5]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) && len(p.weapons) > 1 {
		for i := range p.weapons {
			if i < len(p.weapons)-1 && p.curWeapon == p.weapons[i] {
				p.curWeapon = p.weapons[(i + 1)]
				break
			} else if i == len(p.weapons)-1 && p.curWeapon == p.weapons[i] {
				p.curWeapon = p.weapons[0]
				break
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.Key7) && len(p.secondaryWeapons) > 0 {
		p.curSecondaryWeapon = p.secondaryWeapons[0]
	} else if ebiten.IsKeyPressed(ebiten.Key8) && len(p.secondaryWeapons) > 1 {
		p.curSecondaryWeapon = p.secondaryWeapons[1]
	} else if ebiten.IsKeyPressed(ebiten.Key9) && len(p.secondaryWeapons) > 2 {
		p.curSecondaryWeapon = p.secondaryWeapons[2]
	} else if ebiten.IsKeyPressed(ebiten.Key0) && len(p.secondaryWeapons) > 3 {
		p.curSecondaryWeapon = p.secondaryWeapons[3]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) && len(p.secondaryWeapons) > 1 {
		for i := range p.secondaryWeapons {
			if i < len(p.secondaryWeapons)-1 && p.curWeapon == p.secondaryWeapons[i] {
				p.curSecondaryWeapon = p.secondaryWeapons[(i + 1)]
				break
			} else if i == len(p.secondaryWeapons)-1 && p.curWeapon == p.secondaryWeapons[i] {
				p.curSecondaryWeapon = p.secondaryWeapons[0]
				break
			}
		}
	}

	p.curWeapon.shootCooldown.Update()
	if p.curWeapon.shootCooldown.IsReady() && (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) {
		if p.curWeapon.ammo <= 0 {
			return
		}
		p.curWeapon.shootCooldown.Reset()
		p.curWeapon.Shoot(p)
	}
	if p.curSecondaryWeapon != nil {
		p.curSecondaryWeapon.shootCooldown.Update()
		if p.curSecondaryWeapon.shootCooldown.IsReady() && ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			if p.curSecondaryWeapon.ammo <= 0 {
				return
			}
			p.curSecondaryWeapon.shootCooldown.Reset()
			p.curSecondaryWeapon.Shoot(p)
		}
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
