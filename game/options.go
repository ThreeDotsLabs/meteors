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
	scale := ebiten.DeviceScaleFactor()
	OptionsMenu := OptionsMenu{
		Game: g,
		Items: []*MenuItem{
			{
				Label:   "1024x768",
				Active:  true,
				Choosen: true,
				Pos:     0,
				Action: func(g *Game) error {
					g.Options.ScreenWidth = config.ScreenWidth1024X768 * scale
					g.Options.ScreenHeight = config.ScreenHeight1024X768 * scale
					g.Options.ScreenFontHeight = int(config.Screen1024X768FontHeight * scale)
					g.Options.ScreenFontWidth = int(config.Screen1024X768FontWidth * scale)
					g.Options.ScreenXProfileShift = int(config.Screen1024X768XProfileShift * scale)
					g.Options.ScreenYProfileShift = int(config.Screen1024X768YProfileShift * scale)
					g.Options.ScreenXMenuShift = int(config.Screen1024X768XMenuShift * scale)
					g.Options.ScreenYMenuShift = int(config.Screen1024X768YMenuShift * scale)
					g.Options.ScreenYMenuHeight = int(config.Screen1024X768YMenuHeight * scale)
					g.Options.ScreenYProfileMenuShift = int(config.Screen1024X768YProfileMenuShift * scale)
					g.Options.ScoreFont = assets.ScoreFont1024x768
					g.Options.InfoFont = assets.InfoFont1024x768
					g.Options.SmallFont = assets.SmallFont1024x768
					g.Options.ProfileFont = assets.ProfileFont1024x768
					g.Options.ProfileBigFont = assets.ProfileBigFont1024x768
					g.Options.ResolutionMultipler = 0.5
					g.Options.ProjectileResMulti = 1
					g.Options.ResolutionMultiplerX = 1
					g.Options.ResolutionMultiplerY = 1
					g.ResolutionChange = true
					ebiten.SetWindowSize(int(g.Options.ScreenWidth), int(g.Options.ScreenHeight))
					return nil
				},
			},
			{
				Label:   "1920x1080",
				Active:  true,
				Choosen: false,
				Pos:     1,
				Action: func(g *Game) error {
					g.Options.ScreenWidth = config.ScreenWidth1920x1080 * scale
					g.Options.ScreenHeight = config.ScreenHeight1920x1080 * scale
					g.Options.ScreenFontHeight = int(config.Screen1920x1080FontHeight * scale)               //config.Screen1920x1080FontHeight
					g.Options.ScreenFontWidth = int(config.Screen1920x1080FontWidth * scale)                 //config.Screen1920x1080FontWidth
					g.Options.ScreenXProfileShift = int(config.Screen1920x1080XProfileShift * scale)         //config.Screen1920x1080XProfileShift
					g.Options.ScreenYProfileShift = int(config.Screen1920x1080YProfileShift * scale)         //config.Screen1920x1080YProfileShift
					g.Options.ScreenXMenuShift = int(config.Screen1920x1080XMenuShift * scale)               //config.Screen1920x1080XMenuShift
					g.Options.ScreenYMenuShift = int(config.Screen1920x1080YMenuShift * scale)               //config.Screen1920x1080YMenuShift
					g.Options.ScreenYMenuHeight = int(config.Screen1920x1080YMenuHeight * scale)             //config.Screen1920x1080YMenuHeight
					g.Options.ScreenYProfileMenuShift = int(config.Screen1920x1080YProfileMenuShift * scale) //config.Screen1920x1080YProfileMenuShift
					g.Options.ScoreFont = assets.ScoreFont1920x1080
					g.Options.InfoFont = assets.InfoFont1920x1080
					g.Options.SmallFont = assets.SmallFont1920x1080
					g.Options.ProfileFont = assets.ProfileFont1920x1080
					g.Options.ProfileBigFont = assets.ProfileBigFont1920x1080
					g.Options.ResolutionMultipler = 1
					g.Options.ProjectileResMulti = 2
					g.Options.ResolutionMultiplerX = 1.875
					g.Options.ResolutionMultiplerY = 1.40625
					g.ResolutionChange = true
					ebiten.SetWindowSize(int(g.Options.ScreenWidth), int(g.Options.ScreenHeight))
					return nil
				},
			},
			{
				Label:   "fullscreen on/off",
				Active:  true,
				Choosen: false,
				Pos:     3,
				Action: func(g *Game) error {
					if g.Options.Fullscreen {
						g.Options.Fullscreen = false
						ebiten.SetFullscreen(false)
						return nil
					}
					g.Options.Fullscreen = true
					ebiten.SetFullscreen(true)
					fw, _ := ebiten.ScreenSizeInFullscreen()
					if g.Options.ScreenWidth < float64(fw) {
						xOffsetMod := (float64(fw) - float64(g.Options.ScreenWidth)) / 2
						ebiten.SetWindowPosition(int(xOffsetMod), 100)
					}
					return nil
				},
			},
			{
				Label:   "main menu",
				Active:  true,
				Choosen: false,
				Pos:     4,
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
