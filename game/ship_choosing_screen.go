package game

import (
	"astrogame/config"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/math/fixed"
)

type shipChoosingScreen struct {
	Game         *Game
	Items        []*shipMenuItem
	returnButton *MenuItem
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
	menuItems[0].MenuItem.Choosen = true
	strokeX := ships[0].Sprite.Bounds().Dx() / 2
	shiftX := (int(g.Options.ScreenWidth) - (len(ships)*ships[0].Sprite.Bounds().Dx() + strokeX*(len(ships)-1))) / 2
	for i, item := range menuItems {
		cellWidth := item.Ship.Sprite.Bounds().Dx()
		cellHeight := item.Ship.Sprite.Bounds().Dy()
		stroke := cellWidth / 2
		item.MenuItem.vector = image.Rectangle{
			Min: image.Point{X: shiftX + (cellWidth+stroke)*i, Y: int(g.Options.ScreenHeight / 3)},
			Max: image.Point{X: shiftX + (cellWidth+stroke)*i + cellWidth, Y: int(g.Options.ScreenHeight/3) + cellHeight},
		}
	}
	var dot fixed.Point26_6
	glifImg, _, _, _, _ := g.Options.ScoreFont.Glyph(dot, 'a')
	charWidth := glifImg.Bounds().Dx()
	charHeight := glifImg.Bounds().Dy()
	labelReturn := "return to main menu"
	chars := len([]rune(labelReturn))
	fontShiftX := (int(g.Options.ScreenWidth) - charWidth*chars) / 2
	shipChoosingScreen.Items = menuItems
	shipChoosingScreen.Game = g
	shipChoosingScreen.returnButton = &MenuItem{
		Label:   labelReturn,
		Active:  true,
		Choosen: false,
		vector: image.Rectangle{
			Min: image.Point{X: fontShiftX, Y: int(g.Options.ScreenHeight) - g.Options.ScreenYProfileShift},
			Max: image.Point{X: fontShiftX + charWidth*chars, Y: int(g.Options.ScreenHeight) - g.Options.ScreenYProfileShift + charHeight},
		},
		Action: func(g *Game) error {
			g.state = config.MainMenu
			return nil
		},
	}
	return &shipChoosingScreen
}

func (scs *shipChoosingScreen) Update() {
	scs.shipChoosingMenuUpdate()
}

func (scs *shipChoosingScreen) Draw(screen *ebiten.Image) {
	scs.Game.DrawBg(screen)
	scs.shipChoosingMenuDraw(screen)
}

func (scs *shipChoosingScreen) shipChoosingMenuUpdate() {
	Items := scs.Items
	g := scs.Game
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
				_ = Items[i].MenuItem.Action(g)
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
	scs.returnButton.Choosen = false

	if mouseX >= scs.returnButton.vector.Min.X && mouseX <= scs.returnButton.vector.Max.X && mouseY >= scs.returnButton.vector.Min.Y && mouseY <= scs.returnButton.vector.Max.Y && scs.returnButton.Active {
		if scs.returnButton.Active {
			scs.returnButton.Choosen = true
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			err := scs.returnButton.Action(scs.Game)
			if err != nil {
				log.Fatal(err)
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

func (scs *shipChoosingScreen) shipChoosingMenuDraw(screen *ebiten.Image) {
	Items := scs.Items
	g := scs.Game
	g.DrawBg(screen)
	scale := ebiten.DeviceScaleFactor()
	// for idx, i := range Items {
	// 	vector.StrokeRect(screen, float32(i.MenuItem.vector.Min.X), float32(i.MenuItem.vector.Min.Y), float32(i.MenuItem.vector.Max.X-i.MenuItem.vector.Min.X), float32(i.MenuItem.vector.Max.Y-i.MenuItem.vector.Min.Y), 1, color.RGBA{uint8(55 * idx), uint8(55 * idx), uint8(55 * idx), 255}, false)
	// }
	var dot fixed.Point26_6
	glifImg, _, _, _, _ := g.Options.ProfileFont.Glyph(dot, 'A')
	charWidth := glifImg.Bounds().Dx()
	charHeight := glifImg.Bounds().Dy()
	glifRetButImg, _, _, _, _ := g.Options.ScoreFont.Glyph(dot, 'a')
	charRetButHeight := glifRetButImg.Bounds().Dy()
	for _, i := range Items {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i.MenuItem.vector.Min.X), float64(i.MenuItem.vector.Min.Y))
		colorr := color.RGBA{255, 255, 255, 255}
		if i.MenuItem.Choosen {
			colorr = color.RGBA{179, 14, 14, 255}
		} else if !i.MenuItem.Active {
			colorr = color.RGBA{100, 100, 100, 255}
		}
		chars := len([]rune(i.MenuItem.Label))
		fontShiftX := (i.MenuItem.vector.Max.X - i.MenuItem.vector.Min.X - chars*charWidth) / 2
		screen.DrawImage(i.Ship.Sprite, op)
		text.Draw(screen, fmt.Sprintf("%v", i.MenuItem.Label), text.FaceWithLineHeight(g.Options.ProfileFont, 20*scale), i.MenuItem.vector.Min.X+fontShiftX-2, i.MenuItem.vector.Max.Y+charHeight*3-2, color.RGBA{0, 0, 0, 255})
		text.Draw(screen, fmt.Sprintf("%v", i.MenuItem.Label), text.FaceWithLineHeight(g.Options.ProfileFont, 20*scale), i.MenuItem.vector.Min.X+fontShiftX, i.MenuItem.vector.Max.Y+charHeight*3, colorr)
	}

	colorr := color.RGBA{255, 255, 255, 255}
	if scs.returnButton.Choosen {
		colorr = color.RGBA{179, 14, 14, 255}
	} else if !scs.returnButton.Active {
		colorr = color.RGBA{100, 100, 100, 255}
	}
	//vector.StrokeRect(screen, float32(scs.returnButton.vector.Min.X), float32(scs.returnButton.vector.Min.Y), float32(scs.returnButton.vector.Max.X-scs.returnButton.vector.Min.X), float32(scs.returnButton.vector.Max.Y-scs.returnButton.vector.Min.Y), 1, color.RGBA{255, 255, 255, 255}, false)
	text.Draw(screen, fmt.Sprintf("%v", scs.returnButton.Label), scs.Game.Options.ScoreFont, scs.returnButton.vector.Min.X+1-2, scs.returnButton.vector.Min.Y+charRetButHeight-2, color.RGBA{0, 0, 0, 255})
	text.Draw(screen, fmt.Sprintf("%v", scs.returnButton.Label), scs.Game.Options.ScoreFont, scs.returnButton.vector.Min.X+1, scs.returnButton.vector.Min.Y+charRetButHeight, colorr)
}

func makeShipMenuItem(ship *Ship) *shipMenuItem {
	return &shipMenuItem{
		Ship: ship,
		MenuItem: &MenuItem{
			Label:   ship.Name,
			Active:  true,
			Choosen: false,
			Action: func(g *Game) error {
				switch ship.Name {
				case "Angry ocelot":
					ship.UniqueWeapon = NewAngryOcelotWeapon(g.player)
				case "Mighty orca":
					ship.UniqueWeapon = NewMightyOrcaWeapon(g.player)
				case "Shady weasel":
					ship.UniqueWeapon = NewShadyWeaselWeapon(g.player)
				}
				g.choosenStartShip = ship
				g.player.SetShip(ship)
				g.state = config.InGame
				return nil
			},
		},
	}
}
