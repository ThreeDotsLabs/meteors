package game

import (
	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Item struct {
	game       *Game
	position   config.Vector
	target     config.Vector
	rotation   float64
	TargetType string
	movement   config.Vector
	itemType   *config.Item
}

func NewItem(g *Game, target config.Vector, pos config.Vector, itemType *config.Item) *Item {
	direction := config.Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}
	normalizedDirection := direction.Normalize()

	movement := config.Vector{
		X: normalizedDirection.X * itemType.Velocity,
		Y: normalizedDirection.Y * itemType.Velocity,
	}

	modSprite := objects.ScaleImg(itemType.Sprite, 1.2)
	itemType.Sprite = modSprite

	i := &Item{
		game:     g,
		position: pos,
		movement: movement,
		itemType: itemType,
	}

	return i
}

func (i *Item) CollideWithPlayer(p *Player) {
	if i.itemType.AmmoType != nil {
		switch i.itemType.AmmoType.WeaponType {
		case config.LightRocket:
			for _, w := range p.weapons {
				if w.projectile.wType.WeaponName == config.LightRocket {
					w.ammo += i.itemType.AmmoType.Amount
				}
			}
		}
	} else if i.itemType.WeaponType != nil {
		switch i.itemType.WeaponType.WeaponName {
		case config.LightRocket:
			for _, w := range p.weapons {
				if w.projectile.wType.WeaponName == config.LightRocket {
					lightRocket := NewWeapon(config.LightRocket, p)
					w.ammo += lightRocket.ammo
				}
			}
		case config.DoubleLightRocket:
			persist := false
			for _, w := range p.weapons {
				if w.projectile.wType.WeaponName == config.DoubleLightRocket {
					doubleLightRocket := NewWeapon(config.DoubleLightRocket, p)
					w.ammo += doubleLightRocket.ammo
					persist = true
				}
			}
			if !persist {
				p.weapons = append(p.weapons, NewWeapon(config.DoubleLightRocket, p))
			}
		case config.LaserCannon:
			persist := false
			for _, w := range p.weapons {
				if w.projectile.wType.WeaponName == config.LaserCannon {
					laserCannon := NewWeapon(config.LaserCannon, p)
					w.ammo += laserCannon.ammo
					persist = true
				}
			}
			if !persist {
				p.weapons = append(p.weapons, NewWeapon(config.LaserCannon, p))
			}
		case config.MachineGun:
			persist := false
			for _, w := range p.weapons {
				if w.projectile.wType.WeaponName == config.MachineGun {
					machineGun := NewWeapon(config.MachineGun, p)
					w.ammo += machineGun.ammo
					persist = true
				}
			}
			if !persist {
				p.weapons = append(p.weapons, NewWeapon(config.MachineGun, p))
			}
		// Secondary weapons
		case config.ClusterMines:
			persist := false
			for _, w := range p.secondaryWeapons {
				if w.projectile.wType.WeaponName == config.ClusterMines {
					clusterMines := NewWeapon(config.ClusterMines, p)
					w.ammo += clusterMines.ammo
					persist = true
				}
			}
			if !persist {
				p.secondaryWeapons = append(p.secondaryWeapons, NewWeapon(config.ClusterMines, p))
				if p.curSecondaryWeapon == nil {
					p.curSecondaryWeapon = p.secondaryWeapons[0]
				}
			}
		case config.BigBomb:
			persist := false
			for _, w := range p.secondaryWeapons {
				if w.projectile.wType.WeaponName == config.ClusterMines {
					bigBomb := NewWeapon(config.BigBomb, p)
					w.ammo += bigBomb.ammo
					persist = true
				}
			}
			if !persist {
				p.secondaryWeapons = append(p.secondaryWeapons, NewWeapon(config.BigBomb, p))
				if p.curSecondaryWeapon == nil {
					p.curSecondaryWeapon = p.secondaryWeapons[0]
				}
			}
		}
	} else if i.itemType.HealType != nil {
		p.params.HP += i.itemType.HealType.HP
	} else if i.itemType.ShieldType != nil {
		if p.shield != nil {
			p.shield.HP += i.itemType.ShieldType.HP
		} else {
			w := p.sprite.Bounds().Dx()
			h := p.sprite.Bounds().Dy()
			px, py := p.position.X-float64(w)/2, p.position.Y-float64(h)/2
			p.shield = &Shield{
				position: p.position,
				HP:       i.itemType.ShieldType.HP,
				sprite:   i.itemType.ShieldType.Sprite,
			}
			animationPos := config.Vector{
				X: px,
				Y: py,
			}
			shieldAnimation := NewAnimation(animationPos, assets.ShieldSpriteSheet, 1, 192, 192, true, "shield", 0)
			p.animations = append(p.animations, shieldAnimation)
			p.game.AddAnimation(shieldAnimation)
		}
	}
}

func (i *Item) SetDirection(target config.Vector, pos config.Vector, itemType *config.Item) {
	direction := config.Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}

	normalizedDirection := direction.Normalize()

	i.movement = config.Vector{
		X: normalizedDirection.X * itemType.Velocity,
		Y: normalizedDirection.Y * itemType.Velocity,
	}
	i.position = pos
}

func (i *Item) Update() {
	i.position.X += i.movement.X
	i.position.Y += i.movement.Y
	i.rotation += i.itemType.RotationSpeed

	direction := config.Vector{
		X: i.target.X - i.position.X,
		Y: i.target.Y - i.position.Y,
	}
	normalizedDirection := direction.Normalize()

	movement := config.Vector{
		X: normalizedDirection.X * i.itemType.Velocity,
		Y: normalizedDirection.Y * i.itemType.Velocity,
	}
	i.movement = movement
}

func (i *Item) Draw(screen *ebiten.Image) {
	objects.RotateAndTranslateObject(i.rotation, i.itemType.Sprite, screen, i.position.X, i.position.Y)
}

func (i *Item) Collider() image.Rectangle {
	bounds := i.itemType.Sprite.Bounds()
	return image.Rectangle{
		Min: image.Point{
			X: int(i.position.X),
			Y: int(i.position.Y),
		},
		Max: image.Point{
			X: int(i.position.X + float64(bounds.Dx())),
			Y: int(i.position.Y + float64(bounds.Dy())),
		},
	}
}
