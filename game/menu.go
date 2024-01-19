package game

import (
	"astrogame/assets"
	"astrogame/config"
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MainMenu struct {
	Game  *Game
	Items []*MenuItem
}

type MenuItem struct {
	vector  *image.Rectangle
	Label   string
	Choosen bool
	Pos     int
	Action  func(g *Game)
}

func StartGame(g *Game) {
	if g.started {
		g.Reset()
	}
	g.state = config.InGame
}

func ExitGame(g *Game) {
	g.state = config.MainMenu
}

func ContinueGame(g *Game) {
	g.state = config.InGame
}

func NewMainMenu(g *Game) *MainMenu {
	return &MainMenu{
		Game: g,
		Items: []*MenuItem{
			{
				vector:  &image.Rectangle{Min: image.Point{X: config.ScreenWidth/2 - 100, Y: config.ScreenHeight/2 - 20}, Max: image.Point{X: config.ScreenWidth/2 + 230, Y: config.ScreenHeight / 2}},
				Label:   "Start new game",
				Action:  StartGame,
				Choosen: true,
				Pos:     0,
			},
			{
				vector:  &image.Rectangle{Min: image.Point{X: config.ScreenWidth/2 - 100, Y: config.ScreenHeight/2 + 60}, Max: image.Point{X: config.ScreenWidth/2 + 138, Y: config.ScreenHeight/2 + 80}},
				Label:   "Exit game",
				Action:  ExitGame,
				Choosen: false,
				Pos:     10,
			},
		},
	}
}

func (m *MainMenu) Update() {
	sort.Slice(m.Items, func(i, j int) bool { return m.Items[i].Pos < m.Items[j].Pos })
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		for i := 0; i < len(m.Items); i++ {
			if m.Items[i].Choosen && i < len(m.Items)-1 {
				m.Items[i].Choosen = false
				m.Items[(i + 1)].Choosen = true
				break
			}
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		for i := 0; i < len(m.Items); i++ {
			if m.Items[i].Choosen && i > 0 {
				m.Items[i].Choosen = false
				m.Items[(i - 1)].Choosen = true
				break
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		for i := 0; i < len(m.Items); i++ {
			if m.Items[i].Choosen {
				m.Items[i].Action(m.Game)
				break
			}
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()

		// Check if the mouse click is within the bounds of any menu item
		for i, menuItem := range m.Items {
			if mouseX >= menuItem.vector.Min.X && mouseX <= menuItem.vector.Max.X && mouseY >= menuItem.vector.Min.Y && mouseY <= menuItem.vector.Max.Y {
				m.Items[i].Action(m.Game)
				break
			}
		}
	}

	// if ebiten.IsKeyPressed(ebiten.KeyEscape) && m.Game.started {
	// 	ContinueGame(m.Game)
	// }
}

func (m *MainMenu) Draw(screen *ebiten.Image) {
	_, y16 := m.Game.BgPosition()
	offsetY := float64(-y16) / 64
	// Draw bgImage on the screen repeatedly.
	const repeat = 3
	h := m.Game.bgImage.Bounds().Dy()
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, -float64(h*j))
			op.GeoM.Translate(0, -offsetY)
			screen.DrawImage(m.Game.bgImage, op)
		}
	}
	vector.StrokeRect(screen, config.ScreenWidth/2-100, config.ScreenHeight/2-20, 230, 20, 1, color.RGBA{255, 255, 255, 255}, false)
	vector.StrokeRect(screen, config.ScreenWidth/2-100, config.ScreenHeight/2+60, 138, 20, 1, color.RGBA{255, 255, 255, 255}, false)
	for index, i := range m.Items {
		colorr := color.RGBA{255, 255, 255, 255}
		if i.Choosen {
			colorr = color.RGBA{179, 14, 14, 255}
		}
		text.Draw(screen, fmt.Sprintf("%v", i.Label), assets.ScoreFont, config.ScreenWidth/2-98, config.ScreenHeight/2+80*index-2, color.RGBA{0, 0, 0, 255})
		text.Draw(screen, fmt.Sprintf("%v", i.Label), assets.ScoreFont, config.ScreenWidth/2-100, config.ScreenHeight/2+80*index, colorr)
	}
}
