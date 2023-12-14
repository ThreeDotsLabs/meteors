package game

import (
	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
	"fmt"
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	player            *Player
	meteorSpawnTimer  *config.Timer
	meteors           []*objects.Meteor
	bullets           []*objects.Bullet
	enemies           []*Enemy
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
	CurWave           *config.Wave
}

func NewGame() *Game {
	g := &Game{
		meteorSpawnTimer:  config.NewTimer(config.MeteorSpawnTime),
		baseVelocity:      config.BaseMeteorVelocity,
		velocityTimer:     config.NewTimer(config.MeteorSpeedUpTime),
		enemySpawnTimer:   config.NewTimer(config.Levels[0].Stages[0].Waves[0].Batches[0].Type.EnemySpawnTime),
		batchesSpawnTimer: config.NewTimer(config.Levels[0].Stages[0].Waves[0].Batches[0].BatchSpawnTime),
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

	g.batchesSpawnTimer.Update()
	if g.batchesSpawnTimer.IsReady() {
		g.batchesSpawnTimer.Reset()
		if len(g.CurWave.Batches) > 0 {
			batch := g.CurWave.Batches[0]
			elemInLineCount := 0
			linesCount := 0.0
			for i := 0; i < batch.Count; i++ {
				var target config.Vector
				var startPos config.Vector
				e := NewEnemy(target, startPos, *batch.Type)
				enemyWidth := e.enemyType.Sprite.Bounds().Dx()
				enemyHight := e.enemyType.Sprite.Bounds().Dy()
				switch batch.TargetType {
				case "player":
					target = config.Vector{
						X: g.player.position.X,
						Y: g.player.position.Y,
					}
				case "straight":
					target = config.Vector{
						X: batch.Type.EnemiesStartPos.X,
						Y: batch.Type.EnemiesStartPos.Y + config.ScreenHeight,
					}
				}
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

				elemInLineCount++
				e.SetDirection(target, startPos, *batch.Type)
				g.enemies = append(g.enemies, e)
			}
			g.CurWave.Batches = slices.Delete(g.CurWave.Batches, 0, 1)
		}
	}

	if len(g.CurWave.Batches) == 0 && len(g.enemies) == 0 {
		if g.CurWave.WaveId < len(g.CurStage.Waves)-1 {
			g.CurWave = &g.CurStage.Waves[g.CurWave.WaveId+1]
			g.enemySpawnTimer = config.NewTimer(g.CurWave.EnemyType.EnemySpawnTime)
		} else {
			if g.CurStage.MeteorsCount == 0 && g.CurStage.StageId < len(g.curLevel.Stages)-1 {
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

	// Check for meteor/bullet collisions
	for i, m := range g.meteors {
		for j, b := range g.bullets {
			if m.Collider().Intersects(b.Collider()) {
				if (i < len(g.meteors)-1) && (j < len(g.bullets)-1) {
					g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
					g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
					g.score++
				}
			}
		}
	}

	// Check for enemy/bullet collisions
	for i, m := range g.enemies {
		for j, b := range g.bullets {
			if m.Collider().Intersects(b.Collider()) {
				if (i < len(g.enemies)) && (j < len(g.bullets)-1) {
					g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
					g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
					g.score++
				}
			}
		}
	}
	target := config.Vector{
		X: g.player.position.X,
		Y: g.player.position.Y,
	}
	for _, e := range g.enemies {
		e.Update(target)
	}

	for _, m := range g.meteors {
		m.Update()
	}

	for _, b := range g.bullets {
		b.Update()
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

	for _, b := range g.bullets {
		b.Draw(screen)
	}
	msg := fmt.Sprintf("StageEn: %v, StageId: %v", g.CurWave.EnemiesCount, g.CurStage.StageId)
	ebitenutil.DebugPrint(screen, msg)
	text.Draw(screen, fmt.Sprintf("Level: %v Stage: %v Wave: %v", g.curLevel.LevelId+1, g.CurStage.StageId+1, g.CurWave.WaveId+1), assets.InfoFont, 20, 50, color.White)
	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, config.ScreenWidth/2-100, 50, color.White)
}

func (g *Game) AddBullet(b *objects.Bullet) {
	g.bullets = append(g.bullets, b)
}

func (g *Game) Reset() {
	g.player = NewPlayer(g)
	g.meteors = nil
	g.bullets = nil
	g.enemies = nil
	g.score = 0
	var newLevels = config.NewLevels()
	g.curLevel = newLevels[0]
	g.CurStage = &g.curLevel.Stages[0]
	g.CurWave = &g.CurStage.Waves[0]
	g.levels = newLevels
	g.enemySpawnTimer.Reset()
	g.meteorSpawnTimer.Reset()
	g.batchesSpawnTimer.Reset()
	g.baseVelocity = config.BaseMeteorVelocity
	g.velocityTimer.Reset()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
