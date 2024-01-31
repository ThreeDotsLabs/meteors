package game

import (
	"astrogame/config"
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type shipChoosingScreen struct {
	Game  *Game
	Items []*shipMenuItem
}
type shipMenuItem struct {
	Ship     *Ship
	MenuItem *MenuItem
}

func NewShipChoosingScreen(g *Game) *shipChoosingScreen {
	ships := []*Ship{AngryOcelot, MightyOrca, ShadyWeasel}
	var shipChoosingScreen shipChoosingScreen
	var menuItems []*shipMenuItem
	for _, ship := range ships {
		menuItems = append(menuItems, makeShipMenuItem(ship))
	}
	for i, item := range menuItems {
		cellWidth := item.Ship.Sprite.Bounds().Dx() / 2
		cellHeight := item.Ship.Sprite.Bounds().Dy() / 2
		stroke := 16
		item.MenuItem.vector = image.Rectangle{
			Min: image.Point{X: stroke + (cellWidth)*i, Y: int(g.Options.ScreenHeight/2) - cellHeight},
			Max: image.Point{X: stroke + (cellWidth)*i + cellWidth, Y: int(g.Options.ScreenHeight / 2)},
		}
	}
	shipChoosingScreen.Items = menuItems
	shipChoosingScreen.Game = g
	return &shipChoosingScreen
}

func (scs *shipChoosingScreen) Update() {
	shipChoosingMenuUpdate(scs.Game, scs.Items)
}

func (scs *shipChoosingScreen) Draw(screen *ebiten.Image) {
	scs.Game.DrawBg(screen)
	shipChoosingMenuDraw(scs.Game, scs.Items, screen)
}

func shipChoosingMenuUpdate(g *Game, Items []*shipMenuItem) {
	// Keyboard arrows input handling
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		for i := 0; i < len(Items); i++ {
			if Items[i].MenuItem.Choosen && i < len(Items)-1 {
				if Items[i+1].MenuItem.Active {
					Items[i].MenuItem.Choosen = false
					Items[i+1].MenuItem.Choosen = true
					break
				} else {
					Items[i].MenuItem.Choosen = false
					Items[i+1].MenuItem.Choosen = true
				}
			}
			continue
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		for i := len(Items) - 1; i >= 0; i-- {
			if Items[i].MenuItem.Choosen && i > 0 {
				if Items[i-1].MenuItem.Active {
					Items[i].MenuItem.Choosen = false
					Items[i-1].MenuItem.Choosen = true
					break
				} else {
					Items[i].MenuItem.Choosen = false
					Items[i-1].MenuItem.Choosen = true
				}
			}
			continue
		}
	}

	// Choose menu item
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		for i := 0; i < len(Items); i++ {
			if Items[i].MenuItem.Choosen {
				Items[i].MenuItem.Action(g)
				break
			}
		}
	}

	// Mouse hover on menu items
	mouseX, mouseY := ebiten.CursorPosition()
	for i, shipItem := range Items {
		if mouseX >= shipItem.MenuItem.vector.Min.X && mouseX <= shipItem.MenuItem.vector.Max.X && mouseY >= shipItem.MenuItem.vector.Min.Y && mouseY <= shipItem.MenuItem.vector.Max.Y && Items[i].MenuItem.Active {
			Items[i].MenuItem.Choosen = false
			if Items[i].MenuItem.Active {
				Items[i].MenuItem.Choosen = true
				for idx, menuItemInner := range Items {
					if idx != i {
						menuItemInner.MenuItem.Choosen = false
					}
				}
				break
			}
		}
	}

	// Mouse click on menu items
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for i, shipItem := range Items {
			if mouseX >= shipItem.MenuItem.vector.Min.X && mouseX <= shipItem.MenuItem.vector.Max.X && mouseY >= shipItem.MenuItem.vector.Min.Y && mouseY <= shipItem.MenuItem.vector.Max.Y && Items[i].MenuItem.Active {
				_ = Items[i].MenuItem.Action(g)
				break
			}
		}
	}

	// Return to game if started
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = config.InGame
	}
}

func shipChoosingMenuDraw(g *Game, Items []*shipMenuItem, screen *ebiten.Image) {
	g.DrawBg(screen)
	scale := ebiten.DeviceScaleFactor()
	for idx, i := range Items {
		vector.StrokeRect(screen, float32(i.MenuItem.vector.Min.X), float32(i.MenuItem.vector.Min.Y), float32(i.MenuItem.vector.Max.X), float32(i.MenuItem.vector.Max.Y), 1, color.RGBA{uint8(55 * idx), uint8(55 * idx), uint8(55 * idx), 255}, false)
	}

	for _, i := range Items {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i.MenuItem.vector.Min.X), float64(i.MenuItem.vector.Min.Y))
		colorr := color.RGBA{255, 255, 255, 255}
		if i.MenuItem.Choosen {
			colorr = color.RGBA{179, 14, 14, 255}
		} else if !i.MenuItem.Active {
			colorr = color.RGBA{100, 100, 100, 255}
		}
		screen.DrawImage(i.Ship.Sprite, op)
		text.Draw(screen, fmt.Sprintf("%v", i.MenuItem.Label), text.FaceWithLineHeight(g.Options.ScoreFont, 20*scale), i.MenuItem.vector.Max.X-2, i.MenuItem.vector.Min.Y-2, color.RGBA{0, 0, 0, 255})
		text.Draw(screen, fmt.Sprintf("%v", i.MenuItem.Label), text.FaceWithLineHeight(g.Options.ScoreFont, 20*scale), i.MenuItem.vector.Max.X, i.MenuItem.vector.Min.Y, colorr)
	}
}

func makeShipMenuItem(ship *Ship) *shipMenuItem {
	return &shipMenuItem{
		Ship: ship,
		MenuItem: &MenuItem{
			Label:   ship.Name,
			Active:  true,
			Choosen: false,
			vector:  image.Rect(0, 0, 0, 0),
			Action: func(g *Game) error {
				g.choosenStartShip = ship
				return nil
			},
		},
	}
}
