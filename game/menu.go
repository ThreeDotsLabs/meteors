package game

import (
	"astrogame/assets"
	"astrogame/config"
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

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
	g.state = config.InGame
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
				Label:   "Options",
				Action:  StartGame,
				Active:  false,
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
			Min: image.Point{X: config.ScreenWidth/2 - config.Screen1024X768XMenuShift, Y: config.ScreenHeight/2 - config.Screen1024X768YMenuShift + config.Screen1024X768YMenuHeight*idx - config.Screen1024X768FontHeight},
			Max: image.Point{X: config.ScreenWidth/2 - config.Screen1024X768XMenuShift + chars*config.Screen1024X768FontWidth, Y: config.ScreenHeight/2 - config.Screen1024X768YMenuShift + config.Screen1024X768YMenuHeight*idx + config.Screen1024X768FontHeight},
		}
	}
	return &MainMenu
}

func (m *MainMenu) Update() error {
	// Keyboard arrows input handling
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		for i := 0; i < len(m.Items); i++ {
			if m.Items[i].Choosen && i < len(m.Items)-1 {
				if m.Items[i+1].Active {
					m.Items[i].Choosen = false
					m.Items[i+1].Choosen = true
					break
				} else {
					m.Items[i].Choosen = false
					m.Items[i+1].Choosen = true
				}
			}
			continue
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		for i := len(m.Items) - 1; i >= 0; i-- {
			if m.Items[i].Choosen && i > 0 {
				if m.Items[i-1].Active {
					m.Items[i].Choosen = false
					m.Items[i-1].Choosen = true
					break
				} else {
					m.Items[i].Choosen = false
					m.Items[i-1].Choosen = true
				}
			}
			continue
		}
	}

	// Choose menu item
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		for i := 0; i < len(m.Items); i++ {
			if m.Items[i].Choosen {
				err := m.Items[i].Action(m.Game)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	// Mouse hover on menu items
	mouseX, mouseY := ebiten.CursorPosition()
	for i, menuItem := range m.Items {
		if mouseX >= menuItem.vector.Min.X && mouseX <= menuItem.vector.Max.X && mouseY >= menuItem.vector.Min.Y && mouseY <= menuItem.vector.Max.Y && m.Items[i].Active {
			m.Items[i].Choosen = false
			if m.Items[i].Active {
				m.Items[i].Choosen = true
				for idx, menuItemInner := range m.Items {
					if idx != i {
						menuItemInner.Choosen = false
					}
				}
				break
			}
		}
	}

	// Mouse click on menu items
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		for i, menuItem := range m.Items {
			if mouseX >= menuItem.vector.Min.X && mouseX <= menuItem.vector.Max.X && mouseY >= menuItem.vector.Min.Y && mouseY <= menuItem.vector.Max.Y && m.Items[i].Active {
				err := m.Items[i].Action(m.Game)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	// Return to game if started
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && m.Game.started {
		err := ContinueGame(m.Game)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MainMenu) Draw(screen *ebiten.Image) {
	m.Game.DrawBg(screen)

	// for _, i := range m.Items {
	// 	chars := len([]rune(i.Label))
	// 	vector.StrokeRect(screen, float32(i.vector.Min.X), float32(i.vector.Min.Y), float32(chars*config.Screen1024X768FontWidth), config.Screen1024X768FontHeight, 1, color.RGBA{255, 255, 255, 255}, false)
	// }

	for index, i := range m.Items {
		colorr := color.RGBA{255, 255, 255, 255}
		if i.Choosen {
			colorr = color.RGBA{179, 14, 14, 255}
		} else if !i.Active {
			colorr = color.RGBA{100, 100, 100, 255}
		}
		text.Draw(screen, fmt.Sprintf("%v", i.Label), assets.ScoreFont, config.ScreenWidth/2-config.Screen1024X768XMenuShift-2, config.ScreenHeight/2-config.Screen1024X768YMenuShift+config.Screen1024X768YMenuHeight*index-2, color.RGBA{0, 0, 0, 255})
		text.Draw(screen, fmt.Sprintf("%v", i.Label), assets.ScoreFont, config.ScreenWidth/2-config.Screen1024X768XMenuShift, config.ScreenHeight/2-config.Screen1024X768YMenuShift+config.Screen1024X768YMenuHeight*index, colorr)
	}
}
