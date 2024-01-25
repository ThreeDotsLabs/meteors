package game

import (
	"astrogame/config"
	"image"
	"sort"

	"astrogame/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type OptionsMenu struct {
	Game  *Game
	Items []*MenuItem
}

func NewOptionsMenu(g *Game) *OptionsMenu {
	OptionsMenu := OptionsMenu{
		Game: g,
		Items: []*MenuItem{
			{
				Label:   "1024x768",
				Active:  true,
				Choosen: true,
				Pos:     0,
				Action: func(g *Game) error {
					g.Options.ScreenWidth = config.ScreenWidth1024X768
					g.Options.ScreenHeight = config.ScreenHeight1024X768
					g.Options.ScreenFontHeight = config.Screen1024X768FontHeight
					g.Options.ScreenFontWidth = config.Screen1024X768FontWidth
					g.Options.ScreenXProfileShift = config.Screen1024X768XProfileShift
					g.Options.ScreenYProfileShift = config.Screen1024X768YProfileShift
					g.Options.ScreenXMenuShift = config.Screen1024X768XMenuShift
					g.Options.ScreenYMenuShift = config.Screen1024X768YMenuShift
					g.Options.ScreenYMenuHeight = config.Screen1024X768YMenuHeight
					g.Options.ScoreFont = assets.ScoreFont1024x768
					g.Options.InfoFont = assets.InfoFont1024x768
					g.Options.SmallFont = assets.SmallFont1024x768
					g.Options.ProfileFont = assets.ProfileFont1024x768
					g.Options.ProfileBigFont = assets.ProfileBigFont1024x768
					g.Options.ResolutionMultipler = 1
					return nil
				},
			},
			{
				Label:   "1920x1080",
				Active:  true,
				Choosen: false,
				Pos:     1,
				Action: func(g *Game) error {
					g.Options.ScreenWidth = config.ScreenWidth1920x1080
					g.Options.ScreenHeight = config.ScreenHeight1920x1080
					g.Options.ScreenFontHeight = config.Screen1920x1080FontHeight
					g.Options.ScreenFontWidth = config.Screen1920x1080FontWidth
					g.Options.ScreenXProfileShift = config.Screen1920x1080XProfileShift
					g.Options.ScreenYProfileShift = config.Screen1920x1080YProfileShift
					g.Options.ScreenXMenuShift = config.Screen1920x1080XMenuShift
					g.Options.ScreenYMenuShift = config.Screen1920x1080YMenuShift
					g.Options.ScreenYMenuHeight = config.Screen1920x1080YMenuHeight
					g.Options.ScoreFont = assets.ScoreFont1920x1080
					g.Options.InfoFont = assets.InfoFont1920x1080
					g.Options.SmallFont = assets.SmallFont1920x1080
					g.Options.ProfileFont = assets.ProfileFont1920x1080
					g.Options.ProfileBigFont = assets.ProfileBigFont1920x1080
					g.Options.ResolutionMultipler = 2
					return nil
				},
			},
			{
				Label:   "Return to main menu",
				Active:  true,
				Choosen: false,
				Pos:     2,
				Action: func(g *Game) error {
					g.state = config.MainMenu
					return nil
				},
			},
		},
	}
	sort.Slice(OptionsMenu.Items, func(i, j int) bool { return OptionsMenu.Items[i].Pos < OptionsMenu.Items[j].Pos })
	for idx, i := range OptionsMenu.Items {
		chars := len([]rune(i.Label))
		i.vector = image.Rectangle{
			Min: image.Point{X: int(g.Options.ScreenWidth/2) - g.Options.ScreenXMenuShift, Y: int(g.Options.ScreenHeight/2) - g.Options.ScreenYMenuShift + g.Options.ScreenYMenuHeight*idx - g.Options.ScreenFontHeight},
			Max: image.Point{X: int(g.Options.ScreenWidth/2) - g.Options.ScreenXMenuShift + chars*g.Options.ScreenFontWidth, Y: int(g.Options.ScreenHeight/2) - g.Options.ScreenYMenuShift + g.Options.ScreenYMenuHeight*idx + g.Options.ScreenFontHeight},
		}
	}
	return &OptionsMenu
}

func (m *OptionsMenu) Update() error {
	return MenuUpdate(m.Game, m.Items)
}

func (m *OptionsMenu) Draw(screen *ebiten.Image) {
	MenuDraw(m.Game, m.Items, screen)
}
