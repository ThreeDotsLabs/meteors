package game

import (
	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
	"fmt"
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	player            *Player
	meteorSpawnTimer  *config.Timer
	meteors           []*objects.Meteor
	projectiles       []*Projectile
	beams             []*Beam
	enemyBeams        []*Beam
	enemyProjectiles  []*Projectile
	enemies           []*Enemy
	items             []*Item
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
}

func NewGame() *Game {
	g := &Game{
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
	}

	g.player = NewPlayer(g)

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

func (g *Game) Update() error {
	g.BgMove()
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += config.MeteorSpeedUpAmount
	}

	g.player.Update()
	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()
		if g.CurStage.MeteorsCount > 0 {
			m := objects.NewMeteor(g.baseVelocity)
			g.meteors = append(g.meteors, m)
			g.CurStage.MeteorsCount--
		}
	}

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
			startPos = config.Vector{
				X: float64(config.ScreenWidth/2) - (float64(itemWidth) / 2),
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
			for i := 0; i < batch.Count; i++ {
				var target config.Vector
				var startPos config.Vector
				e := NewEnemy(g, target, startPos, *batch.Type)
				e.TargetType = batch.TargetType
				enemyWidth := e.enemyType.Sprite.Bounds().Dx()
				enemyHight := e.enemyType.Sprite.Bounds().Dy()
				switch batch.StartPositionType {
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

	for _, e := range g.enemies {
		if e.TargetType == "player" {
			e.target = config.Vector{
				X: g.player.position.X,
				Y: g.player.position.Y,
			}
		}
		e.Update()
	}

	for _, m := range g.meteors {
		m.Update()
	}

	for _, p := range g.projectiles {
		p.Update()
	}

	for _, p := range g.enemyProjectiles {
		if p.wType.TargetType == "auto" && p.owner == "enemy" {
			p.target = config.Vector{
				X: g.player.position.X,
				Y: g.player.position.Y,
			}
		}
		p.Update()
	}

	for _, i := range g.items {
		i.Update()
	}

	// Check for meteor/projectile collisions
	for i, m := range g.meteors {
		for j, b := range g.projectiles {
			if m.Collider().Intersects(b.Collider()) {
				if (i < len(g.meteors)-1) && (j < len(g.projectiles)-1) {
					g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
					g.projectiles = append(g.projectiles[:j], g.projectiles[j+1:]...)
					g.score++
				}
			}
		}
	}

	// Check for meteor/enemy projectile collisions
	for i, m := range g.meteors {
		for j, b := range g.enemyProjectiles {
			if m.Collider().Intersects(b.Collider()) {
				if (i < len(g.meteors)-1) && (j < len(g.projectiles)-1) {
					g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
					g.enemyProjectiles = append(g.enemyProjectiles[:j], g.enemyProjectiles[j+1:]...)
				}
			}
		}
	}

	// Check for enemy/projectile collisions
	for i, m := range g.enemies {
		for j, b := range g.projectiles {
			if m.Collider().Intersects(b.Collider()) && b.owner == "player" {
				m.HP -= b.wType.Damage
				if m.HP <= 0 {
					if (i < len(g.enemies)) && (j < len(g.projectiles)) {
						g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
						g.score++
					}
				}
				if j < len(g.projectiles) {
					g.projectiles = append(g.projectiles[:j], g.projectiles[j+1:]...)
				}
			}
		}
	}

	// Check for enemy projectile/player projectile collisions
	for i, m := range g.enemyProjectiles {
		for j, b := range g.projectiles {
			if m.Collider().Intersects(b.Collider()) {
				if (i < len(g.enemyProjectiles)) && (j < len(g.projectiles)-1) {
					g.enemyProjectiles = append(g.enemyProjectiles[:i], g.enemyProjectiles[i+1:]...)
					g.projectiles = append(g.projectiles[:j], g.projectiles[j+1:]...)
				}
			}
		}
	}

	// Check for projectiles/player collisions
	for i, p := range g.enemyProjectiles {
		if p.Collider().Intersects(g.player.Collider()) {
			g.player.hp -= p.wType.Damage
			if i < len(g.enemyProjectiles) {
				g.enemyProjectiles = append(g.enemyProjectiles[:i], g.enemyProjectiles[i+1:]...)
			}
			if g.player.hp <= 0 {
				g.Reset()
				break
			}
		}
	}

	// Check for meteor/player collisions
	for _, m := range g.meteors {
		if m.Collider().Intersects(g.player.Collider()) {
			g.Reset()
			break
		}
	}

	// Check for enemy/player collisions
	for _, e := range g.enemies {
		if e.Collider().Intersects(g.player.Collider()) {
			g.Reset()
			break
		}
	}

	// Check for item/player collisions
	for i, item := range g.items {
		if item.Collider().Intersects(g.player.Collider()) {
			item.CollideWithPlayer(g.player)
			if i < len(g.items) {
				g.items = append(g.items[:i], g.items[i+1:]...)
			}
		}
	}

	// Check for enemy/beam collisions
	// for i, m := range g.enemies {
	// 	for j, b := range g.beams {
	// 		if m.Collider().Intersects(b.Collider()) && b.owner == "player" {
	// 			m.HP -= b.Damage
	// 			if m.HP <= 0 {
	// 				if i < len(g.enemies) {
	// 					g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
	// 					g.score++
	// 				}
	// 			}
	// 		}
	// 		g.beams = append(g.beams[:j], g.beams[j+1:]...)
	// 	}
	// }

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	_, y16 := g.BgPosition()
	offsetY := float64(-y16) / 32

	// Draw bgImage on the screen repeatedly.
	const repeat = 3
	h := g.bgImage.Bounds().Dy()
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, float64(h*j))
			op.GeoM.Translate(0, offsetY)
			screen.DrawImage(g.bgImage, op)
		}
	}

	g.player.Draw(screen)

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

	for _, b := range g.beams {
		b.Draw(screen)
	}

	// Draw the hit points bar
	barX := config.ScreenWidth - 120
	vector.DrawFilledRect(screen, float32(barX-2), 38, 104, 24, color.RGBA{255, 255, 255, 255}, false)
	vector.DrawFilledRect(screen, float32(barX), 40, float32(g.player.hp)*10, 20, color.RGBA{179, 14, 14, 255}, false)

	// Draw weapons
	for i, w := range g.player.weapons {
		offset := 20
		object := w.projectile.wType.Sprite
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i*offset+offset), config.ScreenHeight-float64(60))
		screen.DrawImage(object, op)
	}

	// msg := fmt.Sprintf("StageEn: %v, StageId: %v", g.CurWave.EnemiesCount, g.CurStage.StageId)
	// ebitenutil.DebugPrint(screen, msg)
	text.Draw(screen, fmt.Sprintf("Level: %v Stage: %v Wave: %v Ammo: %v", g.curLevel.LevelId+1, g.CurStage.StageId+1, g.CurWave.WaveId+1, g.player.curWeapon.ammo), assets.InfoFont, 20, 50, color.White)
	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, config.ScreenWidth/2-100, 50, color.White)
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
		g.beams = append(g.beams, b)
	} else {
		g.enemyBeams = append(g.enemyBeams, b)
	}
}

func (g *Game) Reset() {
	g.player = NewPlayer(g)
	g.meteors = nil
	g.projectiles = nil
	g.enemyProjectiles = nil
	g.enemies = nil
	g.items = nil
	g.beams = nil
	g.score = 0
	var newLevels = config.NewLevels()
	g.curLevel = newLevels[0]
	g.CurStage = &g.curLevel.Stages[0]
	g.CurWave = &g.CurStage.Waves[0]
	g.levels = newLevels
	g.meteorSpawnTimer.Reset()
	g.batchesSpawnTimer.Reset()
	g.itemSpawnTimer.Reset()
	g.baseVelocity = config.BaseMeteorVelocity
	g.velocityTimer.Reset()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
