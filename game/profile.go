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
	lineHeight := 25
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
	var profileItems = profileItemsType{
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
	}
	for i, profItem := range profileItems {
		newItem := ProfileItem{
			Icon:     profItem.icon,
			IconPos:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*6+20, barStroke+config.Screen1024X768YProfileShift+lineHeight*i, config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*i+lineHeight),
			Label:    profItem.label,
			LabelPos: image.Rect(barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight*i, config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*i+lineHeight),
			ValuePos: image.Rect(config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight*i, config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*i+lineHeight),
		}
		minusFunc := newItem.MakeMinusAction(i, profItem.barType, profItem.creditsCost, profItem.getter, profItem.increase)
		plusFunc := newItem.MakePlusAction(profItem.creditsCost, profItem.getter, profItem.increase)
		newItem.Buttons = append(newItem.Buttons, newItem.MakeButton(image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight*i, config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*i+lineHeight), "-", minusFunc))
		newItem.Buttons = append(newItem.Buttons, newItem.MakeButton(image.Rect(config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight*i, config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*i+lineHeight), "+", plusFunc))
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
	return &profileScreen
	// return &ProfileScreen{
	// 	Game:    g,
	// 	credits: 50,
	// 	returnButton: &MenuItem{
	// 		Active:  true,
	// 		Choosen: false,
	// 		Pos:     0,
	// 		Label:   "return in game",
	// 		vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, config.ScreenHeight-config.Screen1024X768YProfileShift, config.Screen1024X768XProfileShift+barStroke+section*7+220, config.ScreenHeight-config.Screen1024X768YProfileShift+28),
	// 		Action: func(g *Game) error {
	// 			g.state = config.InGame
	// 			return nil
	// 		},
	// 	},
	// 	LeftBar: &Bar{
	// 		Items: []*ProfileItem{
	// 			{
	// 				Label:     "Health points",
	// 				ValueInt:  g.player.params.HP,
	// 				PrevValue: g.player.params.HP,
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[0].ValueInt = g.player.params.HP
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[0].PrevValue = g.player.params.HP
	// 				},
	// 				LabelPos: image.Rect(barStroke, barStroke+config.Screen1024X768YProfileShift, config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight),
	// 				ValuePos: image.Rect(config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift, config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: true,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift, config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.player.params.HP > g.profile.LeftBar.Items[0].PrevValue {
	// 								g.player.params.HP--
	// 								g.profile.credits += 10
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift, config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 10 {
	// 								g.player.params.HP++
	// 								g.profile.credits -= 10
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Label:    "Ship's speed",
	// 				ValueInt: int(g.player.params.speed),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[1].ValueInt = int(g.player.params.speed)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[1].PrevValue = int(g.player.params.speed)
	// 				},
	// 				LabelPos: image.Rect(barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight, config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*2),
	// 				ValuePos: image.Rect(config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight, config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight, config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.speed) > g.profile.LeftBar.Items[1].PrevValue {
	// 								g.player.params.speed--
	// 								g.profile.credits += 50
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight, config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 50 {
	// 								g.player.params.speed++
	// 								g.profile.credits -= 50
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Label:    "Light missile fire rate X",
	// 				Icon:     objects.ScaleImg(assets.MissileSprite, 0.6),
	// 				IconPos:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*6+20, barStroke+config.Screen1024X768YProfileShift+lineHeight*2, config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2+lineHeight),
	// 				ValueInt: int(g.player.params.LightRocketSpeedUpscale),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[2].ValueInt = int(g.player.params.LightRocketSpeedUpscale)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[2].PrevValue = int(g.player.params.LightRocketSpeedUpscale)
	// 				},
	// 				LabelPos: image.Rect(barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight*2, config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*2+lineHeight),
	// 				ValuePos: image.Rect(config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight*2, config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight*2, config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.LightRocketSpeedUpscale) > g.profile.LeftBar.Items[2].PrevValue {
	// 								g.player.params.LightRocketSpeedUpscale--
	// 								g.profile.credits += 2
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight*2, config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 2 {
	// 								g.player.params.LightRocketSpeedUpscale++
	// 								g.profile.credits -= 2
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Label:    "Light missile velocity X",
	// 				ValueInt: int(g.player.params.LightRocketVelocityMultiplier),
	// 				Icon:     objects.ScaleImg(assets.MissileSprite, 0.6),
	// 				IconPos:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*6+20, barStroke+config.Screen1024X768YProfileShift+lineHeight*3, config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*3+lineHeight),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[3].ValueInt = int(g.player.params.LightRocketVelocityMultiplier)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[3].PrevValue = int(g.player.params.LightRocketVelocityMultiplier)
	// 				},
	// 				LabelPos: image.Rect(barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight*3, config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*3+lineHeight),
	// 				ValuePos: image.Rect(config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight*3, config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*3+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight*3, config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*3+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.LightRocketVelocityMultiplier) > g.profile.LeftBar.Items[3].PrevValue {
	// 								g.player.params.LightRocketVelocityMultiplier--
	// 								g.profile.credits++
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight*3, config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*3+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 1 {
	// 								g.player.params.LightRocketVelocityMultiplier++
	// 								g.profile.credits--
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Label:    "Double missile fire rate X",
	// 				ValueInt: int(g.player.params.DoubleLightRocketSpeedUpscale),
	// 				Icon:     objects.ScaleImg(assets.DoubleMissileSprite, 0.45),
	// 				IconPos:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*6+22, barStroke+config.Screen1024X768YProfileShift+lineHeight*4, config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*4+lineHeight),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[4].ValueInt = int(g.player.params.DoubleLightRocketSpeedUpscale)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[4].PrevValue = int(g.player.params.DoubleLightRocketSpeedUpscale)
	// 				},
	// 				LabelPos: image.Rect(barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight*4, config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*4+lineHeight),
	// 				ValuePos: image.Rect(config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight*4, config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*4+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight*4, config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*4+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.DoubleLightRocketSpeedUpscale) > g.profile.LeftBar.Items[4].PrevValue {
	// 								g.player.params.DoubleLightRocketSpeedUpscale--
	// 								g.profile.credits += 4
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight*4, config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*4+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 4 {
	// 								g.player.params.DoubleLightRocketSpeedUpscale++
	// 								g.profile.credits -= 4
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Label:    "Double missile velocity X",
	// 				ValueInt: int(g.player.params.DoubleLightRocketVelocityMultiplier),
	// 				Icon:     objects.ScaleImg(assets.DoubleMissileSprite, 0.45),
	// 				IconPos:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*6+22, barStroke+config.Screen1024X768YProfileShift+lineHeight*5, config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*5+lineHeight),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[5].ValueInt = int(g.player.params.DoubleLightRocketVelocityMultiplier)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[5].PrevValue = int(g.player.params.DoubleLightRocketVelocityMultiplier)
	// 				},
	// 				LabelPos: image.Rect(barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight*5, config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*5+lineHeight),
	// 				ValuePos: image.Rect(config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight*5, config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*5+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight*5, config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*5+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.DoubleLightRocketVelocityMultiplier) > g.profile.LeftBar.Items[5].PrevValue {
	// 								g.player.params.DoubleLightRocketVelocityMultiplier--
	// 								g.profile.credits += 4
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight*5, config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*5+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 4 {
	// 								g.player.params.DoubleLightRocketVelocityMultiplier++
	// 								g.profile.credits -= 4
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Label:    "Machine gun fire rate X",
	// 				ValueInt: int(g.player.params.MachineGunSpeedUpscale),
	// 				Icon:     objects.ScaleImg(assets.MachineGun, 0.45),
	// 				IconPos:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*6+22, barStroke+config.Screen1024X768YProfileShift+lineHeight*6, config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*6+lineHeight),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[6].ValueInt = int(g.player.params.MachineGunSpeedUpscale)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[6].PrevValue = int(g.player.params.MachineGunSpeedUpscale)
	// 				},
	// 				LabelPos: image.Rect(barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight*6, config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*6+lineHeight),
	// 				ValuePos: image.Rect(config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight*6, config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*6+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight*6, config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*6+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.MachineGunSpeedUpscale) > g.profile.LeftBar.Items[6].PrevValue {
	// 								g.player.params.MachineGunSpeedUpscale--
	// 								g.profile.credits += 4
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight*6, config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*6+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 4 {
	// 								g.player.params.MachineGunSpeedUpscale++
	// 								g.profile.credits -= 4
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Label:    "Machine gun velocity X",
	// 				ValueInt: int(g.player.params.MachineGunVelocityMultiplier),
	// 				Icon:     objects.ScaleImg(assets.MachineGun, 0.45),
	// 				IconPos:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*6+22, barStroke+config.Screen1024X768YProfileShift+lineHeight*7, config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*7+lineHeight),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[7].ValueInt = int(g.player.params.MachineGunVelocityMultiplier)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.LeftBar.Items[7].PrevValue = int(g.player.params.MachineGunVelocityMultiplier)
	// 				},
	// 				LabelPos: image.Rect(barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight*7, config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*7+lineHeight),
	// 				ValuePos: image.Rect(config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight*7, config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*7+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight*7, config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*7+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.MachineGunVelocityMultiplier) > g.profile.LeftBar.Items[7].PrevValue {
	// 								g.player.params.MachineGunVelocityMultiplier--
	// 								g.profile.credits += 10
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight*7, config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*7+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 10 {
	// 								g.player.params.MachineGunVelocityMultiplier++
	// 								g.profile.credits -= 10
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	RightBar: &Bar{
	// 		Items: []*ProfileItem{
	// 			{
	// 				Label:    "Cluster mines fire rate X",
	// 				ValueInt: int(g.player.params.ClusterMinesSpeedUpscale),
	// 				Icon:     objects.ScaleImg(assets.ClusterMines, 0.7),
	// 				IconPos:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*6+16, barStroke+config.Screen1024X768YProfileShift, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.RightBar.Items[0].ValueInt = int(g.player.params.ClusterMinesSpeedUpscale)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.RightBar.Items[0].PrevValue = int(g.player.params.ClusterMinesSpeedUpscale)
	// 				},
	// 				LabelPos: image.Rect(barWidth+barStroke*3+barStroke, barStroke+config.Screen1024X768YProfileShift, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight),
	// 				ValuePos: image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.ClusterMinesSpeedUpscale) > g.profile.RightBar.Items[0].PrevValue {
	// 								g.player.params.ClusterMinesSpeedUpscale--
	// 								g.profile.credits += 10
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 10 {
	// 								g.player.params.ClusterMinesSpeedUpscale++
	// 								g.profile.credits -= 10
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Label:    "Cluster mines velocity X",
	// 				ValueInt: int(g.player.params.ClusterMinesVelocityMultiplier),
	// 				Icon:     objects.ScaleImg(assets.ClusterMines, 0.7),
	// 				IconPos:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*6+16, barStroke+config.Screen1024X768YProfileShift+lineHeight, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight+lineHeight),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.RightBar.Items[1].ValueInt = int(g.player.params.ClusterMinesVelocityMultiplier)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.RightBar.Items[1].PrevValue = int(g.player.params.ClusterMinesVelocityMultiplier)
	// 				},
	// 				LabelPos: image.Rect(barWidth+barStroke*3+barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight+lineHeight),
	// 				ValuePos: image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.ClusterMinesVelocityMultiplier) > g.profile.RightBar.Items[1].PrevValue {
	// 								g.player.params.ClusterMinesVelocityMultiplier--
	// 								g.profile.credits += 12
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 12 {
	// 								g.player.params.ClusterMinesVelocityMultiplier++
	// 								g.profile.credits -= 12
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Label:    "Big bomb fire rate X",
	// 				ValueInt: int(g.player.params.BigBombSpeedUpscale),
	// 				Icon:     objects.ScaleImg(assets.BigBomb, 0.7),
	// 				IconPos:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*6+16, barStroke+config.Screen1024X768YProfileShift+lineHeight*2, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2+lineHeight),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.RightBar.Items[2].ValueInt = int(g.player.params.BigBombSpeedUpscale)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.RightBar.Items[2].PrevValue = int(g.player.params.BigBombSpeedUpscale)
	// 				},
	// 				LabelPos: image.Rect(barWidth+barStroke*3+barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight*2, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*2+lineHeight),
	// 				ValuePos: image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight*2, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight*2, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.BigBombSpeedUpscale) > g.profile.RightBar.Items[2].PrevValue {
	// 								g.player.params.BigBombSpeedUpscale--
	// 								g.profile.credits += 16
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight*2, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*2+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 16 {
	// 								g.player.params.BigBombSpeedUpscale++
	// 								g.profile.credits -= 16
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Label:    "Big bomb velocity X",
	// 				ValueInt: int(g.player.params.BigBombVelocityMultiplier),
	// 				Icon:     objects.ScaleImg(assets.BigBomb, 0.7),
	// 				IconPos:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*6+16, barStroke+config.Screen1024X768YProfileShift+lineHeight*3, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*6+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*3+lineHeight),
	// 				UpdateValue: func(g *Game) {
	// 					g.profile.RightBar.Items[3].ValueInt = int(g.player.params.BigBombVelocityMultiplier)
	// 				},
	// 				UpdatePrevValue: func(g *Game) {
	// 					g.profile.RightBar.Items[3].PrevValue = int(g.player.params.BigBombVelocityMultiplier)
	// 				},
	// 				LabelPos: image.Rect(barWidth+barStroke*3+barStroke, barStroke+config.Screen1024X768YProfileShift+lineHeight*3, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+barWidth, barStroke+config.Screen1024X768YProfileShift+lineHeight*3+lineHeight),
	// 				ValuePos: image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*8, barStroke+config.Screen1024X768YProfileShift+lineHeight*3, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*8+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*3+lineHeight),
	// 				Buttons: []*MenuItem{
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     0,
	// 						Label:   "-",
	// 						vector:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*7, barStroke+config.Screen1024X768YProfileShift+lineHeight*3, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*7+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*3+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if int(g.player.params.BigBombVelocityMultiplier) > g.profile.RightBar.Items[3].PrevValue {
	// 								g.player.params.BigBombVelocityMultiplier--
	// 								g.profile.credits += 20
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 					{
	// 						Active:  true,
	// 						Choosen: false,
	// 						Pos:     1,
	// 						Label:   "+",
	// 						vector:  image.Rect(barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*9, barStroke+config.Screen1024X768YProfileShift+lineHeight*3, barWidth+barStroke*3+config.Screen1024X768XProfileShift+barStroke+section*9+section, barStroke+config.Screen1024X768YProfileShift+lineHeight*3+lineHeight),
	// 						Action: func(g *Game) error {
	// 							if g.profile.credits >= 20 {
	// 								g.player.params.BigBombVelocityMultiplier++
	// 								g.profile.credits -= 20
	// 							}
	// 							return nil
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
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
