package game

import (
	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type ProfileScreen struct {
	Game         *Game
	credits      int
	returnButton *MenuItem
	LeftBar      *Bar
	RightBar     *Bar
}

type Bar struct {
	Side  string
	Items []*ProfileItem
}

type ProfileItem struct {
	Label           string
	ValueInt        int
	PrevValue       int
	UpdateValue     func(g *Game)
	UpdatePrevValue func(g *Game)
	Icon            *ebiten.Image
	IconPos         image.Rectangle
	LabelPos        image.Rectangle
	ValuePos        image.Rectangle
	Buttons         []*MenuItem
}
type profileItemTemplate struct {
	label       string
	barType     *Bar
	creditsCost int
	icon        *ebiten.Image
	getter      func() int
	increase    func(int)
}

type profileItemsType []profileItemTemplate

func NewPlayerProfile(g *Game) *ProfileScreen {
	barStroke := 16
	barWidth := (config.ScreenWidth - (config.Screen1024X768XProfileShift*2 + barStroke*2)) / 2
	section := barWidth / 10
	profileScreen := ProfileScreen{
		Game:    g,
		credits: 50,
		returnButton: &MenuItem{
			Active:  true,
			Choosen: false,
			Pos:     0,
			Label:   "return in game",
			vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, config.ScreenHeight-config.Screen1024X768YProfileShift, config.Screen1024X768XProfileShift+barStroke+section*7+220, config.ScreenHeight-config.Screen1024X768YProfileShift+28),
			Action: func(g *Game) error {
				g.state = config.InGame
				return nil
			},
		},
		LeftBar: &Bar{
			Side: "left",
		},
		RightBar: &Bar{
			Side: "right",
		},
	}
	var profileItemsLeft = profileItemsType{
		{
			label:       "Helth points",
			barType:     profileScreen.LeftBar,
			creditsCost: 10,
			icon:        nil,
			getter:      g.player.params.GetHealthPoints,
			increase:    g.player.params.IncreaseHP,
		},
		{
			label:       "Speed",
			barType:     profileScreen.LeftBar,
			creditsCost: 10,
			icon:        nil,
			getter:      g.player.params.GetSpeed,
			increase:    g.player.params.IncreaseSpeed,
		},
		{
			label:       "Light missile fire rate X",
			barType:     profileScreen.LeftBar,
			creditsCost: 20,
			icon:        objects.ScaleImg(assets.MissileSprite, 0.6),
			getter:      g.player.params.GetLightRocketSpeedUpscale,
			increase:    g.player.params.IncreaseLightRocketSpeedUpscale,
		},
		{
			label:       "Light missile velocity X",
			barType:     profileScreen.LeftBar,
			creditsCost: 10,
			icon:        objects.ScaleImg(assets.MissileSprite, 0.6),
			getter:      g.player.params.GetLightRocketVelocityMultiplier,
			increase:    g.player.params.IncreaseLightRocketVelocityMultiplier,
		},
		{
			label:       "Double missile fire rate X",
			barType:     profileScreen.LeftBar,
			creditsCost: 30,
			icon:        objects.ScaleImg(assets.DoubleMissileSprite, 0.6),
			getter:      g.player.params.GetDoubleLightRocketSpeedUpscale,
			increase:    g.player.params.IncreaseDoubleLightRocketSpeedUpscale,
		},
		{
			label:       "Double missile velocity X",
			barType:     profileScreen.LeftBar,
			creditsCost: 15,
			icon:        objects.ScaleImg(assets.DoubleMissileSprite, 0.6),
			getter:      g.player.params.GetDoubleLightRocketVelocityMultiplier,
			increase:    g.player.params.IncreaseDoubleLightRocketVelocityMultiplier,
		},
		{
			label:       "Machine gun fire rate X",
			barType:     profileScreen.LeftBar,
			creditsCost: 35,
			icon:        objects.ScaleImg(assets.MachineGun, 0.6),
			getter:      g.player.params.GetMachineGunSpeedUpscale,
			increase:    g.player.params.IncreaseMachineGunSpeedUpscale,
		},
		{
			label:       "Machine gun projectile velocity X",
			barType:     profileScreen.LeftBar,
			creditsCost: 20,
			icon:        objects.ScaleImg(assets.MachineGun, 0.6),
			getter:      g.player.params.GetMachineGunVelocityMultiplier,
			increase:    g.player.params.IncreaseMachineGunVelocityMultiplier,
		},
		{
			label:       "Laser canon fire rate X",
			barType:     profileScreen.LeftBar,
			creditsCost: 35,
			icon:        objects.ScaleImg(assets.LaserCanon, 0.6),
			getter:      g.player.params.GetLaserCanonSpeedUpscale,
			increase:    g.player.params.IncreaseLaserCanonSpeedUpscale,
		},
		{
			label:       "Double laser canon fire rate X",
			barType:     profileScreen.LeftBar,
			creditsCost: 45,
			icon:        objects.ScaleImg(assets.DoubleLaserCanon, 0.6),
			getter:      g.player.params.GetDoubleLaserCanonSpeedUpscale,
			increase:    g.player.params.IncreaseDoubleLaserCanonSpeedUpscale,
		},
		{
			label:       "Plasma gun fire rate X",
			barType:     profileScreen.LeftBar,
			creditsCost: 52,
			icon:        objects.ScaleImg(assets.PlasmaGun, 0.6),
			getter:      g.player.params.GetPlasmaGunSpeedUpscale,
			increase:    g.player.params.IncreasePlasmaGunSpeedUpscale,
		},
		{
			label:       "Plasma gun projectile velocity X",
			barType:     profileScreen.LeftBar,
			creditsCost: 40,
			icon:        objects.ScaleImg(assets.DoubleLaserCanon, 0.6),
			getter:      g.player.params.GetPlasmaGunVelocityMultiplier,
			increase:    g.player.params.IncreasePlasmaGunVelocityMultiplier,
		},
		{
			label:       "Double plasma gun fire rate X",
			barType:     profileScreen.LeftBar,
			creditsCost: 64,
			icon:        objects.ScaleImg(assets.PlasmaGun, 0.6),
			getter:      g.player.params.GetDoublePlasmaGunSpeedUpscale,
			increase:    g.player.params.IncreaseDoublePlasmaGunSpeedUpscale,
		},
		{
			label:       "Double plasma gun velocity X",
			barType:     profileScreen.LeftBar,
			creditsCost: 48,
			icon:        objects.ScaleImg(assets.PlasmaGun, 0.6),
			getter:      g.player.params.GetDoublePlasmaGunVelocityMultiplier,
			increase:    g.player.params.IncreaseDoublePlasmaGunVelocityMultiplier,
		},
	}
	var profileItemsRight = profileItemsType{
		{
			label:       "Big Bomb fire rate X",
			barType:     profileScreen.RightBar,
			creditsCost: 50,
			icon:        objects.ScaleImg(assets.BigBomb, 0.6),
			getter:      g.player.params.GetBigBombSpeedUpscale,
			increase:    g.player.params.IncreaseBigBombSpeedUpscale,
		},
		{
			label:       "Big Bomb velocity X",
			barType:     profileScreen.RightBar,
			creditsCost: 40,
			icon:        objects.ScaleImg(assets.BigBomb, 0.6),
			getter:      g.player.params.GetBigBombVelocityMultiplier,
			increase:    g.player.params.IncreaseBigBombVelocityMultiplier,
		},
		{
			label:       "Cluster mines fire rate X",
			barType:     profileScreen.RightBar,
			creditsCost: 50,
			icon:        objects.ScaleImg(assets.ClusterMines, 0.6),
			getter:      g.player.params.GetClusterMinesSpeedUpscale,
			increase:    g.player.params.IncreaseClusterMinesSpeedUpscale,
		},
		{
			label:       "Cluster mines velocity X",
			barType:     profileScreen.RightBar,
			creditsCost: 50,
			icon:        objects.ScaleImg(assets.ClusterMines, 0.6),
			getter:      g.player.params.GetClusterMinesVelocityMultiplier,
			increase:    g.player.params.IncreaseClusterMinesVelocityMultiplier,
		},
		{
			label:       "Penta laser fire rate X",
			barType:     profileScreen.RightBar,
			creditsCost: 80,
			icon:        objects.ScaleImg(assets.PentaLaser, 0.6),
			getter:      g.player.params.GetPentaLaserSpeedUpscale,
			increase:    g.player.params.IncreasePentaLaserSpeedUpscale,
		},
	}
	for i, profItem := range profileItemsLeft {
		prepareMenuItem(i, &profItem, &profileScreen)
	}
	for i, profItem := range profileItemsRight {
		prepareMenuItem(i, &profItem, &profileScreen)
	}
	return &profileScreen
}
func (p *ProfileScreen) Update() {
	for _, w := range p.Game.player.weapons {
		w.projectile.VelocityUpdate(p.Game.player)
	}
	// Mouse hover on menu items
	mouseX, mouseY := ebiten.CursorPosition()
	for _, profileItem := range p.LeftBar.Items {
		profileItem.UpdateValue(p.Game)
		for i, menuItem := range profileItem.Buttons {
			profileItem.Buttons[i].Choosen = false
			if mouseX >= menuItem.vector.Min.X && mouseX <= menuItem.vector.Max.X && mouseY >= menuItem.vector.Min.Y && mouseY <= menuItem.vector.Max.Y && profileItem.Buttons[i].Active {
				if profileItem.Buttons[i].Active {
					profileItem.Buttons[i].Choosen = true
					break
				}
			}
		}
	}
	for _, profileItem := range p.RightBar.Items {
		profileItem.UpdateValue(p.Game)
		for i, menuItem := range profileItem.Buttons {
			profileItem.Buttons[i].Choosen = false
			if mouseX >= menuItem.vector.Min.X && mouseX <= menuItem.vector.Max.X && mouseY >= menuItem.vector.Min.Y && mouseY <= menuItem.vector.Max.Y && profileItem.Buttons[i].Active {
				if profileItem.Buttons[i].Active {
					profileItem.Buttons[i].Choosen = true
					break
				}
			}
		}
	}
	if p.returnButton.Active {
		p.returnButton.Choosen = false
	}
	if mouseX >= p.returnButton.vector.Min.X && mouseX <= p.returnButton.vector.Max.X && mouseY >= p.returnButton.vector.Min.Y && mouseY <= p.returnButton.vector.Max.Y && p.returnButton.Active {
		if p.returnButton.Active {
			p.returnButton.Choosen = true
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			err := p.returnButton.Action(p.Game)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Mouse click on menu items
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, profileItem := range p.LeftBar.Items {
			for i, menuItem := range profileItem.Buttons {
				if mouseX >= menuItem.vector.Min.X && mouseX <= menuItem.vector.Max.X && mouseY >= menuItem.vector.Min.Y && mouseY <= menuItem.vector.Max.Y && profileItem.Buttons[i].Active {
					_ = profileItem.Buttons[i].Action(p.Game)
					break
				}
			}
		}
		for _, profileItem := range p.RightBar.Items {
			for i, menuItem := range profileItem.Buttons {
				if mouseX >= menuItem.vector.Min.X && mouseX <= menuItem.vector.Max.X && mouseY >= menuItem.vector.Min.Y && mouseY <= menuItem.vector.Max.Y && profileItem.Buttons[i].Active {
					_ = profileItem.Buttons[i].Action(p.Game)
					break
				}
			}
		}
	}

	// Return to game if started
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && p.Game.started {
		err := ContinueGame(p.Game)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (p *ProfileScreen) Draw(screen *ebiten.Image) {
	p.Game.DrawBg(screen)
	barStroke := 16
	barWidth := (config.ScreenWidth - (config.Screen1024X768XProfileShift*2 + barStroke*2)) / 2
	section := barWidth / 10
	text.Draw(screen, fmt.Sprintf("Available credits: %v", p.credits), assets.ProfileBigFont, config.ScreenWidth/2-200, barStroke+50-2, color.RGBA{0, 0, 0, 255})
	text.Draw(screen, fmt.Sprintf("Available credits: %v", p.credits), assets.ProfileBigFont, config.ScreenWidth/2-200, barStroke+50, color.RGBA{255, 255, 255, 255})
	for _, i := range p.LeftBar.Items {
		text.Draw(screen, fmt.Sprintf("%v", i.Label), assets.ProfileFont, i.LabelPos.Min.X+(section*2)-2, i.LabelPos.Min.Y+barStroke-2, color.RGBA{0, 0, 0, 255})
		text.Draw(screen, fmt.Sprintf("%v", i.Label), assets.ProfileFont, i.LabelPos.Min.X+(section*2), i.LabelPos.Min.Y+barStroke, color.RGBA{255, 255, 255, 255})
		if i.Icon != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i.IconPos.Min.X), float64(i.IconPos.Min.Y))
			screen.DrawImage(i.Icon, op)
		}
		text.Draw(screen, fmt.Sprintf("%v", i.ValueInt), assets.ProfileFont, i.ValuePos.Min.X+section/2-8, i.ValuePos.Min.Y+barStroke, color.RGBA{255, 255, 255, 255})
		for _, menuItem := range i.Buttons {
			colorr := color.RGBA{255, 255, 255, 255}
			if menuItem.Choosen {
				colorr = color.RGBA{179, 14, 14, 255}
			} else if !menuItem.Active {
				colorr = color.RGBA{100, 100, 100, 255}
			}
			text.Draw(screen, fmt.Sprintf("%v", menuItem.Label), assets.SmallFont, menuItem.vector.Min.X+(section/2-4)-2, menuItem.vector.Min.Y+16-2, color.RGBA{0, 0, 0, 255})
			text.Draw(screen, fmt.Sprintf("%v", menuItem.Label), assets.SmallFont, menuItem.vector.Min.X+(section/2-4), menuItem.vector.Min.Y+16, colorr)
		}
	}
	for _, i := range p.RightBar.Items {
		text.Draw(screen, fmt.Sprintf("%v", i.Label), assets.ProfileFont, i.LabelPos.Min.X+(section*2)-2, i.LabelPos.Min.Y+barStroke-2, color.RGBA{0, 0, 0, 255})
		text.Draw(screen, fmt.Sprintf("%v", i.Label), assets.ProfileFont, i.LabelPos.Min.X+(section*2), i.LabelPos.Min.Y+barStroke, color.RGBA{255, 255, 255, 255})
		if i.Icon != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i.IconPos.Min.X), float64(i.IconPos.Min.Y))
			screen.DrawImage(i.Icon, op)
		}
		text.Draw(screen, fmt.Sprintf("%v", i.ValueInt), assets.ProfileFont, i.ValuePos.Min.X+section/2-8, i.ValuePos.Min.Y+barStroke, color.RGBA{255, 255, 255, 255})
		for _, menuItem := range i.Buttons {
			colorr := color.RGBA{255, 255, 255, 255}
			if menuItem.Choosen {
				colorr = color.RGBA{179, 14, 14, 255}
			} else if !menuItem.Active {
				colorr = color.RGBA{100, 100, 100, 255}
			}
			text.Draw(screen, fmt.Sprintf("%v", menuItem.Label), assets.SmallFont, menuItem.vector.Min.X+(section/2-4)-2, menuItem.vector.Min.Y+16-2, color.RGBA{0, 0, 0, 255})
			text.Draw(screen, fmt.Sprintf("%v", menuItem.Label), assets.SmallFont, menuItem.vector.Min.X+(section/2-4), menuItem.vector.Min.Y+16, colorr)
		}
	}
	colorr := color.RGBA{255, 255, 255, 255}
	if p.returnButton.Choosen {
		colorr = color.RGBA{179, 14, 14, 255}
	} else if !p.returnButton.Active {
		colorr = color.RGBA{100, 100, 100, 255}
	}
	//vector.StrokeRect(screen, float32(p.returnButton.vector.Min.X), float32(p.returnButton.vector.Min.Y), float32(p.returnButton.vector.Max.X-p.returnButton.vector.Min.X), float32(p.returnButton.vector.Max.Y-p.returnButton.vector.Min.Y), 1, color.RGBA{255, 255, 255, 255}, false)
	text.Draw(screen, fmt.Sprintf("%v", p.returnButton.Label), assets.ScoreFont, p.returnButton.vector.Min.X+1-2, p.returnButton.vector.Min.Y+22-2, color.RGBA{0, 0, 0, 255})
	text.Draw(screen, fmt.Sprintf("%v", p.returnButton.Label), assets.ScoreFont, p.returnButton.vector.Min.X+1, p.returnButton.vector.Min.Y+22, colorr)
}

func (i *ProfileItem) MakeButton(rect image.Rectangle, label string, action func(g *Game) error) *MenuItem {
	return &MenuItem{
		Active:  true,
		Choosen: false,
		Label:   label,
		vector:  rect,
		Action:  action,
	}
}

func (i *ProfileItem) MakeMinusAction(idx int, barType *Bar, creditsCost int, getter func() int, increase func(int)) func(g *Game) error {
	return func(g *Game) error {
		if getter() > barType.Items[idx].PrevValue {
			increase(-1)
			g.profile.credits += creditsCost
		}
		return nil
	}
}

func (i *ProfileItem) MakePlusAction(creditsCost int, getter func() int, increase func(int)) func(g *Game) error {
	return func(g *Game) error {
		if g.profile.credits >= creditsCost {
			increase(1)
			g.profile.credits -= creditsCost
		}
		return nil
	}
}

func (i *ProfileItem) MakeUpdateValueFunc(idx int, getter func() int, barType *Bar) func(g *Game) {
	return func(g *Game) {
		barType.Items[idx].ValueInt = getter()
	}
}

func (i *ProfileItem) MakeUpdatePrevValueFunc(idx int, getter func() int, barType *Bar) func(g *Game) {
	return func(g *Game) {
		barType.Items[idx].PrevValue = getter()
	}
}

func prepareMenuItem(i int, profItem *profileItemTemplate, profileScreen *ProfileScreen) {
	barStroke := 16
	barWidth := (config.ScreenWidth - (config.Screen1024X768XProfileShift*2 + barStroke*2)) / 2
	section := barWidth / 10
	lineHeight := 25
	rightBarXmod := 0
	if profItem.barType.Side == profileScreen.RightBar.Side {
		rightBarXmod = barWidth + barStroke*3
	}
	newItem := ProfileItem{
		Icon:     profItem.icon,
		IconPos:  image.Rect(rightBarXmod+config.Screen1024X768XProfileShift+barStroke+section*6+20, barStroke+config.Screen1024X768YProfileShift+lineHeight*i, rightBarXmod+config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*i+lineHeight),
		Label:    profItem.label,
		LabelPos: image.Rect(rightBarXmod+barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight*i, rightBarXmod+config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*i+lineHeight),
		ValuePos: image.Rect(rightBarXmod+config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight*i, rightBarXmod+config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*i+lineHeight),
	}
	minusFunc := newItem.MakeMinusAction(i, profItem.barType, profItem.creditsCost, profItem.getter, profItem.increase)
	plusFunc := newItem.MakePlusAction(profItem.creditsCost, profItem.getter, profItem.increase)
	newItem.Buttons = append(newItem.Buttons, newItem.MakeButton(image.Rect(rightBarXmod+config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight*i, rightBarXmod+config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*i+lineHeight), "-", minusFunc))
	newItem.Buttons = append(newItem.Buttons, newItem.MakeButton(image.Rect(rightBarXmod+config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight*i, rightBarXmod+config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*i+lineHeight), "+", plusFunc))
	updValue := newItem.MakeUpdateValueFunc(i, profItem.getter, profItem.barType)
	updPrevValue := newItem.MakeUpdatePrevValueFunc(i, profItem.getter, profItem.barType)
	newItem.UpdateValue = updValue
	newItem.UpdatePrevValue = updPrevValue
	if profItem.barType.Side == profileScreen.LeftBar.Side {
		profileScreen.LeftBar.Items = append(profileScreen.LeftBar.Items, &newItem)
	} else {
		profileScreen.RightBar.Items = append(profileScreen.RightBar.Items, &newItem)
	}
}
