package game

import (
	"astrogame/config"
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Menu interface {
	Update(g *Game, Items []*MenuItem) error
	Draw(g *Game, Items []*MenuItem, screen *ebiten.Image)
}
type MainMenu struct {
	Game  *Game
	Items []*MenuItem
}

type MenuItem struct {
	vector  image.Rectangle
	Label   string
	Active  bool
	Choosen bool
	Pos     int
	Action  func(g *Game) error
}

func StartGame(g *Game) error {
	if g.started {
		g.Reset()
	}
	g.state = config.ShipChoosingWindow
	return nil
}

func ExitGame(g *Game) error {
	return ebiten.Termination
}

func ContinueGame(g *Game) error {
	g.state = config.InGame
	return nil
}

func NewMainMenu(g *Game) *MainMenu {
	MainMenu := MainMenu{
		Game: g,
		Items: []*MenuItem{
			{
				Label:   "Start new game",
				Action:  StartGame,
				Active:  true,
				Choosen: true,
				Pos:     0,
			},
			{
				Label:   "Continue game",
				Action:  ContinueGame,
				Active:  false,
				Choosen: false,
				Pos:     1,
			},
			{
				Label: "Options",
				Action: func(g *Game) error {
					g.state = config.Options
					return nil
				},
				Active:  true,
				Choosen: false,
				Pos:     2,
			},
			{
				Label:   "Exit game",
				Action:  ExitGame,
				Active:  true,
				Choosen: false,
				Pos:     3,
			},
		},
	}
	sort.Slice(MainMenu.Items, func(i, j int) bool { return MainMenu.Items[i].Pos < MainMenu.Items[j].Pos })
	for idx, i := range MainMenu.Items {
		chars := len([]rune(i.Label))
		i.vector = image.Rectangle{
			Min: image.Point{X: int(g.Options.ScreenWidth/2) - g.Options.ScreenXMenuShift, Y: int(g.Options.ScreenHeight/2) - g.Options.ScreenYMenuShift + g.Options.ScreenYMenuHeight*idx - g.Options.ScreenFontHeight},
			Max: image.Point{X: int(g.Options.ScreenWidth/2) - g.Options.ScreenXMenuShift + chars*g.Options.ScreenFontWidth, Y: int(g.Options.ScreenHeight/2) - g.Options.ScreenYMenuShift + g.Options.ScreenYMenuHeight*idx + g.Options.ScreenFontHeight},
		}
	}
	return &MainMenu
}

func (m *MainMenu) Update() error {
	return MenuUpdate(m.Game, m.Items)
}

func (m *MainMenu) Draw(screen *ebiten.Image) {
	m.Game.DrawBg(screen)
	MenuDraw(m.Game, m.Items, screen)
}

func MenuDraw(g *Game, Items []*MenuItem, screen *ebiten.Image) {
	//g.DrawBg(screen)
	scale := ebiten.DeviceScaleFactor()
	// for _, i := range m.Items {
	// 	chars := len([]rune(i.Label))
	// 	vector.StrokeRect(screen, float32(i.vector.Min.X), float32(i.vector.Min.Y), float32(chars*config.Screen1024X768FontWidth), config.Screen1024X768FontHeight, 1, color.RGBA{255, 255, 255, 255}, false)
	// }

	for index, i := range Items {
		colorr := color.RGBA{255, 255, 255, 255}
		if i.Choosen {
			colorr = color.RGBA{179, 14, 14, 255}
		} else if !i.Active {
			colorr = color.RGBA{100, 100, 100, 255}
		}
		text.Draw(screen, fmt.Sprintf("%v", i.Label), text.FaceWithLineHeight(g.Options.ScoreFont, 20*scale), int(g.Options.ScreenWidth/2)-g.Options.ScreenXMenuShift-2, int(g.Options.ScreenHeight/2)-g.Options.ScreenYMenuShift+g.Options.ScreenYMenuHeight*index-2, color.RGBA{0, 0, 0, 255})
		text.Draw(screen, fmt.Sprintf("%v", i.Label), text.FaceWithLineHeight(g.Options.ScoreFont, 20*scale), int(g.Options.ScreenWidth/2)-g.Options.ScreenXMenuShift, int(g.Options.ScreenHeight/2)-g.Options.ScreenYMenuShift+g.Options.ScreenYMenuHeight*index, colorr)
	}
}

func MenuUpdate(g *Game, Items []*MenuItem) error {
	// Keyboard arrows input handling
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		for i := 0; i < len(Items); i++ {
			if Items[i].Choosen && i < len(Items)-1 {
				if Items[i+1].Active {
					Items[i].Choosen = false
					Items[i+1].Choosen = true
					break
				} else {
					Items[i].Choosen = false
					Items[i+1].Choosen = true
				}
			}
			continue
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		for i := len(Items) - 1; i >= 0; i-- {
			if Items[i].Choosen && i > 0 {
				if Items[i-1].Active {
					Items[i].Choosen = false
					Items[i-1].Choosen = true
					break
				} else {
					Items[i].Choosen = false
					Items[i-1].Choosen = true
				}
			}
			continue
		}
	}

	// Choose menu item
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		for i := 0; i < len(Items); i++ {
			if Items[i].Choosen {
				err := Items[i].Action(g)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	// Mouse hover on menu items
	mouseX, mouseY := ebiten.CursorPosition()
	for i, menuItem := range Items {
		if mouseX >= menuItem.vector.Min.X && mouseX <= menuItem.vector.Max.X && mouseY >= menuItem.vector.Min.Y && mouseY <= menuItem.vector.Max.Y && Items[i].Active {
			Items[i].Choosen = false
			if Items[i].Active {
				Items[i].Choosen = true
				for idx, menuItemInner := range Items {
					if idx != i {
						menuItemInner.Choosen = false
					}
				}
				break
			}
		}
	}

	// Mouse click on menu items
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for i, menuItem := range Items {
			if mouseX >= menuItem.vector.Min.X && mouseX <= menuItem.vector.Max.X && mouseY >= menuItem.vector.Min.Y && mouseY <= menuItem.vector.Max.Y && Items[i].Active {
				err := Items[i].Action(g)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	// Return to game if started
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && g.started {
		err := ContinueGame(g)
		if err != nil {
			return err
		}
	}
	return nil
}
