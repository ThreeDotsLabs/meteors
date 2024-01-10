package game

import (
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
					lightRocket := NewWeapon(config.LightRocket)
					w.ammo += lightRocket.ammo
				}
			}
		case config.DoubleLightRocket:
			persist := false
			for _, w := range p.weapons {
				if w.projectile.wType.WeaponName == config.DoubleLightRocket {
					doubleLightRocket := NewWeapon(config.DoubleLightRocket)
					w.ammo += doubleLightRocket.ammo
					persist = true
				}
			}
			if !persist {
				p.weapons = append(p.weapons, NewWeapon(config.DoubleLightRocket))
			}
		case config.LaserCannon:
			persist := false
			for _, w := range p.weapons {
				if w.projectile.wType.WeaponName == config.LaserCannon {
					laserCannon := NewWeapon(config.LaserCannon)
					w.ammo += laserCannon.ammo
					persist = true
				}
			}
			if !persist {
				p.weapons = append(p.weapons, NewWeapon(config.LaserCannon))
			}
		}
	} else if i.itemType.HealType != nil {
		p.hp += i.itemType.HealType.HP
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
	// return config.NewRect(
	// 	i.position.X,
	// 	i.position.Y,
	// 	float64(bounds.Dx()),
	// 	float64(bounds.Dy()),
	// )
}
