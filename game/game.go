package game

import (
	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
	"fmt"
	"image/color"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	menu              *MainMenu
	profile           *ProfileScreen
	state             config.GameState
	player            *Player
	meteorSpawnTimer  *config.Timer
	meteors           []*objects.Meteor
	projectiles       []*Projectile
	blows             []*Blow
	beam              *Beam
	beamAnimations    []*BeamAnimation
	enemyBeams        []*Beam
	enemyProjectiles  []*Projectile
	enemies           []*Enemy
	items             []*Item
	animations        []*Animation
	bgImage           *ebiten.Image
	score             int
	viewport          viewport
	levels            []*config.Level
	curLevel          *config.Level
	CurStage          *config.Stage
	baseVelocity      float64
	velocityTimer     *config.Timer
	enemySpawnTimer   *config.Timer
	batchesSpawnTimer *config.Timer
	itemSpawnTimer    *config.Timer
	CurWave           *config.Wave
	started           bool
}

func NewGame() *Game {
	g := &Game{
		state:             config.MainMenu,
		meteorSpawnTimer:  config.NewTimer(config.MeteorSpawnTime),
		baseVelocity:      config.BaseMeteorVelocity,
		velocityTimer:     config.NewTimer(config.MeteorSpeedUpTime),
		enemySpawnTimer:   config.NewTimer(config.Levels[0].Stages[0].Waves[0].Batches[0].Type.EnemySpawnTime),
		batchesSpawnTimer: config.NewTimer(config.Levels[0].Stages[0].Waves[0].Batches[0].BatchSpawnTime),
		itemSpawnTimer:    config.NewTimer(config.Levels[0].Stages[0].Items[0].ItemSpawnTime),
		bgImage:           config.Levels[0].BgImg,
		levels:            config.Levels,
		curLevel:          config.Levels[0],
		CurStage:          &config.Levels[0].Stages[0],
		CurWave:           &config.Levels[0].Stages[0].Waves[0],
		started:           false,
	}

	g.player = NewPlayer(g)
	g.menu = NewMainMenu(g)
	g.profile = NewPlayerProfile(g)

	return g
}

type viewport struct {
	X int
	Y int
}

func (g *Game) BgMove() {
	s := g.bgImage.Bounds().Size()
	maxX := s.X * 16
	maxY := s.Y * 16

	g.viewport.X += s.X / 32
	g.viewport.Y += s.Y / 32
	g.viewport.X %= maxX
	g.viewport.Y %= maxY
}

func (g *Game) BgPosition() (int, int) {
	return g.viewport.X, g.viewport.Y
}
func (g *Game) MoveBgPosition() {
	g.BgMove()
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += config.MeteorSpeedUpAmount
	}
}

func (g *Game) Update() error {
	g.MoveBgPosition()
	switch g.state {
	case config.Profile:
		g.profile.Update()
	case config.MainMenu:
		err := g.menu.Update()
		if err != nil {
			return err
		}
	case config.InGame:
		// Main menu logic
		if !g.started {
			g.started = true
			// Unlock the continue button
			for i := range g.menu.Items {
				if g.menu.Items[i].Label == "Continue game" {
					g.menu.Items[i].Active = true
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = config.MainMenu
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			for _, i := range g.profile.LeftBar.Items {
				i.UpdatePrevValue(g)
			}
			g.state = config.Profile
		}

		// Game logic
		g.player.Update()

		// Meteor spawning
		g.meteorSpawnTimer.Update()
		if g.meteorSpawnTimer.IsReady() {
			g.meteorSpawnTimer.Reset()
			if g.CurStage.MeteorsCount > 0 {
				m := objects.NewMeteor(g.baseVelocity)
				g.meteors = append(g.meteors, m)
				g.CurStage.MeteorsCount--
			}
		}

		// Item spawning
		g.itemSpawnTimer.Update()
		if g.itemSpawnTimer.IsReady() {
			if len(g.CurStage.Items) > 0 {
				var target config.Vector
				var startPos config.Vector
				itemParam := g.CurStage.Items[0]
				item := NewItem(g, target, startPos, &itemParam)
				if len(g.CurStage.Items) > 1 {
					g.itemSpawnTimer = config.NewTimer(g.CurStage.Items[1].ItemSpawnTime)
				}
				g.itemSpawnTimer.Reset()
				itemWidth := item.itemType.Sprite.Bounds().Dx()
				itemHight := item.itemType.Sprite.Bounds().Dy()
				r := rand.New(rand.NewSource(99))
				posX := r.Intn(config.ScreenWidth - (itemWidth + itemWidth/2))
				startPos = config.Vector{
					X: float64(posX + itemWidth),
					Y: -(float64(itemHight)),
				}
				target = config.Vector{
					X: startPos.X,
					Y: config.ScreenHeight + 10,
				}
				item.SetDirection(target, startPos, &itemParam)
				item.target = target
				g.items = append(g.items, item)
				g.CurStage.Items = slices.Delete(g.CurStage.Items, 0, 1)
			}
		}

		// Enemy spawning
		g.batchesSpawnTimer.Update()
		if g.batchesSpawnTimer.IsReady() {
			if len(g.CurWave.Batches) > 0 {
				batch := g.CurWave.Batches[0]
				if len(g.CurWave.Batches) > 1 {
					g.batchesSpawnTimer = config.NewTimer(g.CurWave.Batches[1].BatchSpawnTime)
				}
				g.batchesSpawnTimer.Reset()
				elemInLineCount := 0
				linesCount := 0.0
				var xOffsetMod float64
				if batch.StartPositionType == "centered" {
					eTst := NewEnemy(g, config.Vector{}, config.Vector{}, *batch.Type)
					enemyWidthTst := eTst.enemyType.Sprite.Bounds().Dx()
					xOffsetMod = (config.ScreenWidth - float64(batch.Count*(enemyWidthTst+int(batch.StartPosOffset))-int(batch.StartPosOffset))) / 2
				}
				for i := 0; i < batch.Count; i++ {
					var target config.Vector
					var startPos config.Vector
					e := NewEnemy(g, target, startPos, *batch.Type)
					e.TargetType = batch.TargetType
					enemyWidth := e.enemyType.Sprite.Bounds().Dx()
					enemyHight := e.enemyType.Sprite.Bounds().Dy()
					switch batch.StartPositionType {
					case "centered":
						xOffset := batch.StartPosOffset
						elemInLine := config.ScreenWidth / (enemyWidth + int(xOffset))
						if elemInLineCount == 0 {
							xOffset = 0.0
						}
						if elemInLineCount >= elemInLine {
							elemInLineCount = 0
							linesCount++
						}
						startPos = config.Vector{
							X: (float64(enemyWidth)+xOffset)*float64(elemInLineCount) + xOffsetMod,
							Y: -(float64(enemyHight)*linesCount + 10),
						}
					case "lines":
						xOffset := batch.StartPosOffset
						elemInLine := config.ScreenWidth / (enemyWidth + int(xOffset))
						if elemInLineCount == 0 {
							xOffset = 0.0
						}
						if elemInLineCount >= elemInLine {
							elemInLineCount = 0
							linesCount++
						}
						startPos = config.Vector{
							X: (float64(enemyWidth) + xOffset) * float64(elemInLineCount),
							Y: -(float64(enemyHight)*linesCount + 10),
						}
					case "checkmate":
						elemInLine := config.ScreenWidth / (enemyWidth * 2)
						if elemInLineCount >= elemInLine {
							elemInLineCount = 0
							linesCount++
						}
						startPos = config.Vector{
							X: float64(enemyWidth)*float64(elemInLineCount)*2 + float64(enemyWidth),
							Y: -(float64(enemyHight)*linesCount + 10),
						}
						if int(linesCount)%2 == 0 {
							startPos = config.Vector{
								X: float64(enemyWidth) * float64(elemInLineCount) * 2,
								Y: -(float64(enemyHight)*linesCount + 10),
							}
						}
					}
					switch batch.TargetType {
					case "player":
						target = config.Vector{
							X: g.player.position.X,
							Y: g.player.position.Y,
						}
					case "straight":
						target = config.Vector{
							X: startPos.X,
							Y: config.ScreenHeight + 10,
						}
					}
					elemInLineCount++
					e.SetDirection(target, startPos, *batch.Type)
					e.target = target
					g.enemies = append(g.enemies, e)
				}
				g.CurWave.Batches = slices.Delete(g.CurWave.Batches, 0, 1)
			}
		}

		if len(g.CurWave.Batches) == 0 {
			if g.CurWave.WaveId < len(g.CurStage.Waves)-1 {
				g.CurWave = &g.CurStage.Waves[g.CurWave.WaveId+1]
			} else {
				if g.CurStage.MeteorsCount == 0 && g.CurStage.StageId < len(g.curLevel.Stages)-1 && len(g.CurStage.Items) == 0 {
					g.CurStage = &g.curLevel.Stages[g.CurStage.StageId+1]
					g.CurWave = &g.CurStage.Waves[0]
				} else {
					if g.curLevel.LevelId < len(g.levels)-1 {
						g.curLevel = g.levels[g.curLevel.LevelId+1]
						g.CurStage = &g.curLevel.Stages[0]
						g.CurWave = &g.CurStage.Waves[0]
					} else {
						g.Reset()
					}
				}
			}
		}

		for i, m := range g.meteors {
			m.Update()
			if m.Collider().Min.Y >= config.ScreenHeight && i < len(g.meteors) {
				g.meteors = slices.Delete(g.meteors, i, i+1)
			}
		}

		for i, p := range g.projectiles {
			p.Update()
			if p.position.Y <= 0 && i < len(g.projectiles) {
				g.projectiles = slices.Delete(g.projectiles, i, i+1)
			}
		}

		for i, p := range g.enemyProjectiles {
			if p.wType.TargetType == "auto" && p.owner == "enemy" {
				p.target = config.Vector{
					X: g.player.position.X,
					Y: g.player.position.Y,
				}
			}
			p.Update()
			if p.position.Y >= config.ScreenHeight && i < len(g.enemyProjectiles) {
				g.enemyProjectiles = slices.Delete(g.enemyProjectiles, i, i+1)
			}
		}

		for i, item := range g.items {
			item.Update()
			if item.position.Y >= config.ScreenHeight && i < len(g.items) {
				g.items = slices.Delete(g.items, i, i+1)
			}
		}

		// Check for meteor/projectile collisions
		// Check for meteor/enemy projectile collisions
		for i, m := range g.meteors {
			for j, b := range g.projectiles {
				if config.IntersectRect(m.Collider(), b.Collider()) {
					if (i < len(g.meteors)-1) && (j < len(g.projectiles)-1) {
						g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
						g.projectiles = append(g.projectiles[:j], g.projectiles[j+1:]...)
						g.score++
					}
				}
			}

			for j, b := range g.enemyProjectiles {
				if config.IntersectRect(m.Collider(), b.Collider()) {
					if (i < len(g.meteors)-1) && (j < len(g.projectiles)-1) {
						g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
						g.enemyProjectiles = append(g.enemyProjectiles[:j], g.enemyProjectiles[j+1:]...)
					}
				}
			}
		}

		// Check for enemy/projectile collisions
		// Check for enemy/player collisions
		// Check for enemy/beam collisions
		// Check for enemy/blow collisions
		for i, m := range g.enemies {
			if m.TargetType == config.TargetTypePlayer {
				m.target = config.Vector{
					X: g.player.position.X,
					Y: g.player.position.Y,
				}
			}
			m.Update()
			if m.position.Y >= config.ScreenHeight && i < len(g.enemies) {
				g.enemies = slices.Delete(g.enemies, i, i+1)
			}

			for j, b := range g.projectiles {
				if config.IntersectRect(m.Collider(), b.Collider()) && b.owner == "player" {
					switch b.wType.WeaponName {
					case config.BigBomb:
						bounds := b.wType.Sprite.Bounds()
						blow := NewBlow(b.position.X+float64(bounds.Dx()/2), b.position.Y+float64(bounds.Dy()/2), float64(bounds.Dx())*4, b.wType.Damage)
						blow.Steps = 5
						g.AddBlow(blow, m.position)
					default:
						m.HP -= b.wType.Damage
						if m.HP <= 0 {
							if (i < len(g.enemies)) && (j < len(g.projectiles)) {
								g.KillEnemy(i)
							}
						}
					}

					if j < len(g.projectiles) {
						g.projectiles[j].intercectAnimation.position = b.position
						g.AddAnimation(g.projectiles[j].intercectAnimation)
						g.projectiles = append(g.projectiles[:j], g.projectiles[j+1:]...)
					}
				}
			}

			for _, blow := range g.blows {
				if config.IntersectCircle(m.Collider(), blow.circle) {
					m.HP -= blow.Damage
					if m.HP <= 0 {
						if i < len(g.enemies) {
							g.KillEnemy(i)
						}
					}
				}
			}

			if config.IntersectRect(m.Collider(), g.player.Collider()) {
				if g.player.shield != nil {
					g.player.shield = nil
				} else {
					g.Reset()
					break
				}
			}

			if g.beam != nil {
				if config.IntersectLine(g.beam.Line, m.Collider()) {
					m.HP -= g.beam.Damage
					if m.HP <= 0 {
						if i < len(g.enemies) {
							g.KillEnemy(i)
						}
					}
				}
			}
		}

		// Check for enemy projectile/player projectile collisions
		for i, m := range g.enemyProjectiles {
			for j, b := range g.projectiles {
				if config.IntersectRect(m.Collider(), b.Collider()) {
					if (i < len(g.enemyProjectiles)) && (j < len(g.projectiles)-1) {
						g.enemyProjectiles[i].intercectAnimation.position = m.position
						g.AddAnimation(g.enemyProjectiles[i].intercectAnimation)
						g.enemyProjectiles = append(g.enemyProjectiles[:i], g.enemyProjectiles[i+1:]...)
						g.projectiles[j].intercectAnimation.position = b.position
						g.AddAnimation(g.projectiles[j].intercectAnimation)
						g.projectiles = append(g.projectiles[:j], g.projectiles[j+1:]...)
					}
				}
			}
			if g.beam != nil {
				if config.IntersectLine(g.beam.Line, m.Collider()) {
					if i < len(g.enemyProjectiles) {
						g.enemyProjectiles[i].intercectAnimation.position = m.position
						g.AddAnimation(g.enemyProjectiles[i].intercectAnimation)
						g.enemyProjectiles = append(g.enemyProjectiles[:i], g.enemyProjectiles[i+1:]...)
					}
				}
			}
		}

		// Check for projectiles/player collisions
		for i, p := range g.enemyProjectiles {
			if config.IntersectRect(p.Collider(), g.player.Collider()) {
				if g.player.shield != nil {
					g.player.shield.HP -= p.wType.Damage
					if g.player.shield.HP <= 0 {
						g.player.shield = nil
					}
				} else {
					g.player.params.HP -= p.wType.Damage
				}
				if i < len(g.enemyProjectiles) {
					g.enemyProjectiles[i].intercectAnimation.position = p.position
					g.AddAnimation(g.enemyProjectiles[i].intercectAnimation)
					g.enemyProjectiles = append(g.enemyProjectiles[:i], g.enemyProjectiles[i+1:]...)
				}
				if g.player.params.HP <= 0 {
					g.Reset()
					break
				}
			}
		}

		// Check for meteor/player collisions
		for _, m := range g.meteors {
			if config.IntersectRect(m.Collider(), g.player.Collider()) {
				g.Reset()
				break
			}
		}

		// Check for item/player collisions
		for i, item := range g.items {
			if config.IntersectRect(item.Collider(), g.player.Collider()) {
				item.CollideWithPlayer(g.player)
				if i < len(g.items) {
					g.items = append(g.items[:i], g.items[i+1:]...)
				}
			}
		}

		if g.beam != nil {
			g.AddBeamAnimation(g.beam.NewBeamAnimation())
			g.beam = nil
		}

		for i, ba := range g.beamAnimations {
			ba.Update()
			if ba.Step >= ba.Steps && i < len(g.beamAnimations) {
				g.beamAnimations = slices.Delete(g.beamAnimations, i, i+1)
			}
		}

		for i, a := range g.animations {
			a.Update()
			if a.currF >= a.numFrames && i < len(g.animations) && !a.looping {
				g.animations = slices.Delete(g.animations, i, i+1)
			}
			if a.name == "shield" && g.player.shield == nil && i < len(g.animations) {
				g.animations = slices.Delete(g.animations, i, i+1)
			}
		}

		// Remove blows
		for k, b := range g.blows {
			b.Update()
			if b.Step >= b.Steps && k < len(g.blows) {
				g.blows = slices.Delete(g.blows, k, k+1)
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case config.Profile:
		g.profile.Draw(screen)
	case config.MainMenu:
		g.menu.Draw(screen)
	case config.InGame:
		_, y16 := g.BgPosition()
		offsetY := float64(-y16) / 64

		// Draw bgImage on the screen repeatedly.
		const repeat = 3
		h := g.bgImage.Bounds().Dy()
		for j := 0; j < repeat; j++ {
			for i := 0; i < repeat; i++ {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(0, -float64(h*j))
				op.GeoM.Translate(0, -offsetY)
				screen.DrawImage(g.bgImage, op)
			}
		}

		g.player.Draw(screen)

		for _, ba := range g.beamAnimations {
			ba.Draw(screen)
		}

		for _, e := range g.enemies {
			e.Draw(screen)
		}

		for _, m := range g.meteors {
			m.Draw(screen)
		}

		for _, p := range g.projectiles {
			p.Draw(screen)
		}

		for _, p := range g.enemyProjectiles {
			p.Draw(screen)
		}

		for _, i := range g.items {
			i.Draw(screen)
		}

		for i, a := range g.animations {
			if i < len(g.animations) {
				a.Draw(screen)
			}
		}

		// Draw the hit points bar
		barX := config.ScreenWidth - 120
		backWidth := float32(g.player.params.HP)*10 + 4
		if backWidth < 104 {
			backWidth = 104
		}

		vector.DrawFilledRect(screen, float32(barX-2), 38, backWidth, 24, color.RGBA{255, 255, 255, 255}, false)
		vector.DrawFilledRect(screen, float32(barX), 40, float32(g.player.params.HP)*10, 20, color.RGBA{179, 14, 14, 255}, false)

		// Draw shield bar
		if g.player.shield != nil {
			backWidth := float32(g.player.shield.HP)*10 + 4
			if backWidth < 104 {
				backWidth = 104
			}
			shiftX := int(backWidth) - 104
			vector.DrawFilledRect(screen, float32(barX-shiftX-2), 82, backWidth, 24, color.RGBA{255, 255, 255, 255}, false)
			vector.DrawFilledRect(screen, float32(barX-shiftX), 84, float32(g.player.shield.HP)*10, 20, color.RGBA{26, 14, 189, 255}, false)
		}

		// Draw weapons
		for i, w := range g.player.weapons {
			offset := 20
			object := w.projectile.wType.Sprite
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i*offset+offset), config.ScreenHeight-float64(60))
			if w.projectile.wType.WeaponName == g.player.curWeapon.projectile.wType.WeaponName {
				vector.DrawFilledRect(screen, float32(i*offset+offset-2), float32(config.ScreenHeight-float64(30)), float32(object.Bounds().Dx()+4), 3, color.RGBA{255, 255, 255, 255}, false)
			}
			text.Draw(screen, fmt.Sprintf("%v", w.ammo), assets.SmallFont, i*offset+offset, config.ScreenHeight-80, color.White)
			screen.DrawImage(object, op)
		}

		// Draw secondary weapons
		for i, w := range g.player.secondaryWeapons {
			startOffset := config.ScreenWidth - 40
			offset := 20
			object := w.projectile.wType.Sprite
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(startOffset-i*offset), config.ScreenHeight-float64(60))
			if g.player.curSecondaryWeapon != nil && w.projectile.wType.WeaponName == g.player.curSecondaryWeapon.projectile.wType.WeaponName {
				vector.DrawFilledRect(screen, float32(startOffset-i*offset-2), float32(config.ScreenHeight-float64(30)), float32(object.Bounds().Dx()+4), 3, color.RGBA{255, 255, 255, 255}, false)
			}
			text.Draw(screen, fmt.Sprintf("%v", w.ammo), assets.SmallFont, startOffset-i*offset, config.ScreenHeight-80, color.White)
			screen.DrawImage(object, op)
		}

		// if g.beam != nil {
		// 	gradRot := float64(180) / math.Pi * g.beam.rotation
		// 	gradRotPl := float64(180) / math.Pi * g.player.rotation
		// 	msg := fmt.Sprintf("Beams: %v, Rotation: %v, PlayerRotation: %v", g.beam.Line, gradRot, gradRotPl)
		// 	ebitenutil.DebugPrint(screen, msg)
		// }
		// for _, a := range g.animations {
		// 	if a.name == "engineFireburst" {
		// 		msg := fmt.Sprintf("X: %v, Y: %v, Angle: %v, Step: %v, Frame: %v", a.position.X, a.position.Y, a.rotation, a.currF, a.curTick)
		// 		ebitenutil.DebugPrint(screen, msg)
		// 	}
		// }

		text.Draw(screen, fmt.Sprintf("Level: %v Stage: %v Wave: %v", g.curLevel.LevelId+1, g.CurStage.StageId+1, g.CurWave.WaveId), assets.InfoFont, 20, 50, color.White)
		text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, config.ScreenWidth/2-100, 50, color.White)
	}
}

func (g *Game) AddProjectile(p *Projectile) {
	if p.owner == "player" {
		g.projectiles = append(g.projectiles, p)
	} else {
		g.enemyProjectiles = append(g.enemyProjectiles, p)
	}
}

func (g *Game) AddBeam(b *Beam) {
	if b.owner == "player" {
		g.beam = b
	} else {
		g.enemyBeams = append(g.enemyBeams, b)
	}
}

func (g *Game) AddBeamAnimation(b *BeamAnimation) {
	g.beamAnimations = append(g.beamAnimations, b)
}

func (g *Game) AddAnimation(a *Animation) {
	g.animations = append(g.animations, a)
}

func (g *Game) AddBlow(b *Blow, target config.Vector) {
	g.blows = append(g.blows, b)
	blowAnimation := NewAnimation(target, assets.BigBlowSpriteSheet, 1, 124, 128, false, "digBlow", 0)
	g.AddAnimation(blowAnimation)
}

func (g *Game) KillEnemy(i int) {
	enemyBlow := NewAnimation(g.enemies[i].position, assets.EnemyBlowSpriteSheet, 1, 73, 75, false, "enemyBlow", 0)
	g.AddAnimation(enemyBlow)
	g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
	g.score++
	g.profile.credits += 10
}

func (g *Game) Reset() {
	g.meteors = nil
	g.projectiles = nil
	g.enemyProjectiles = nil
	g.enemies = nil
	g.items = nil
	g.beamAnimations = nil
	g.animations = nil
	g.score = 0
	g.player = NewPlayer(g)
	var newLevels = config.NewLevels()
	g.curLevel = newLevels[0]
	g.CurStage = &g.curLevel.Stages[0]
	g.CurWave = &g.CurStage.Waves[0]
	g.levels = newLevels
	g.meteorSpawnTimer.Reset()
	g.batchesSpawnTimer = config.NewTimer(g.levels[0].Stages[0].Waves[0].Batches[0].BatchSpawnTime)
	g.itemSpawnTimer.Reset()
	g.baseVelocity = config.BaseMeteorVelocity
	g.velocityTimer.Reset()
	g.started = false
	g.menu = NewMainMenu(g)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
