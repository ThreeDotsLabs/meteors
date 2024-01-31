package config

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth1920x1080             = 1920
	ScreenHeight1920x1080            = 1080
	Screen1920x1080XMenuShift        = 225
	Screen1920x1080YMenuShift        = 169
	Screen1920x1080XProfileShift     = 225
	Screen1920x1080YProfileShift     = 169
	Screen1920x1080YMenuHeight       = 80
	Screen1920x1080FontHeight        = 36
	Screen1920x1080FontWidth         = 32
	Screen1920x1080YProfileMenuShift = 42
	ScreenWidth1024X768              = 1024
	ScreenHeight1024X768             = 768
	Screen1024X768XMenuShift         = 120
	Screen1024X768YMenuShift         = 120
	Screen1024X768XProfileShift      = 120
	Screen1024X768YProfileShift      = 120
	Screen1024X768YMenuHeight        = 40
	Screen1024X768FontHeight         = 20
	Screen1024X768FontWidth          = 16
	Screen1024X768YProfileMenuShift  = 22

	MeteorSpawnTime     = 2 * time.Second
	EnemySpawnTime      = 1 * time.Second
	BaseMeteorVelocity  = 0.25
	MeteorSpeedUpAmount = 0.1
	MeteorSpeedUpTime   = 5 * time.Second
	OwnerEnemy          = "enemy"
	OwnerPlayer         = "player"
	TargetTypePlayer    = "player"
	TargetTypeStraight  = "straight"
	LightRocket         = "lightRocket"
	AutoLightRocket     = "autoLightRocket"
	DoubleLightRocket   = "doubleLightRocket"
	LaserCanon          = "lightCanon"
	DoubleLaserCanon    = "doubleLaserCanon"
	ClusterMines        = "clusterMines"
	BigBomb             = "bigBomb"
	MachineGun          = "machineGun"
	DoubleMachineGun    = "doubleMachineGun"
	PlasmaGun           = "plasmaGun"
	DoublePlasmaGun     = "doublePlasmaGun"
	PentaLaser          = "pentaLaser"
)

type GameState string

var (
	MainMenu           GameState = "mainMenu"
	InGame             GameState = "inGame"
	Profile            GameState = "profile"
	Options            GameState = "options"
	ShipChoosingWindow GameState = "shipChoosing"
)

type ItemTemplate struct {
	Sprite           *ebiten.Image
	Velocity         float64
	ItemSpawnTime    time.Duration
	AmmoType         *AmmoType
	WeaponType       *WeaponType
	SecondWeaponType *WeaponType
	HealType         *HealType
	ShieldType       *ShieldType
}

func (it *ItemTemplate) toItem() Item {
	return Item{
		AmmoType:         it.AmmoType,
		WeaponType:       it.WeaponType,
		SecondWeaponType: it.SecondWeaponType,
		HealType:         it.HealType,
		ShieldType:       it.ShieldType,
		RotationSpeed:    0,
		Sprite:           it.Sprite,
		Velocity:         it.Velocity,
		ItemSpawnTime:    it.ItemSpawnTime,
	}
}

func ConvertItems(items []*ItemTemplate) []Item {
	var res []Item
	for _, item := range items {
		res = append(res, item.toItem())
	}
	return res
}

type Level struct {
	Stages   []Stage
	CurStage *Stage
	BgImg    *ebiten.Image
	Name     string
	Number   int
	LevelId  int
}

type Stage struct {
	MeteorsCount int
	StageId      int
	Items        []Item
	Waves        []Wave
}

type Wave struct {
	WaveId  int
	Batches []EnemyBatch
}

type EnemyBatch struct {
	Type              *EnemyType
	TargetType        string
	StartPositionType string
	BatchSpawnTime    time.Duration
	Count             int
	StartPosOffset    float64
}

type WeaponType struct {
	Sprite                        *ebiten.Image
	ItemSprite                    *ebiten.Image
	IntercectAnimationSpriteSheet *ebiten.Image
	InstantAnimationSpiteSheet    *ebiten.Image
	Scale                         float64
	AnimationOnly                 bool
	Velocity                      float64
	StartVelocity                 float64
	Damage                        int
	Target                        Vector
	TargetType                    string
	WeaponName                    string
	cost                          int
	StartTime                     time.Duration
	StartAmmo                     int
}

type AmmoType struct {
	WeaponName string
	Amount     int
}

type HealType struct {
	HP int
}

type ShieldType struct {
	HP     int
	Sprite *ebiten.Image
}

type Item struct {
	AmmoType         *AmmoType
	WeaponType       *WeaponType
	SecondWeaponType *WeaponType
	HealType         *HealType
	ShieldType       *ShieldType
	RotationSpeed    float64
	Sprite           *ebiten.Image
	Velocity         float64
	Target           Vector
	ItemSpawnTime    time.Duration
}

type EnemyType struct {
	RotationSpeed   float64
	Sprite          *ebiten.Image
	Velocity        float64
	EnemiesStartPos Vector
	EnemySpawnTime  time.Duration
	WeaponTypeStr   string
	WeaponType      *WeaponType
	StartHP         int
}

type EnemyBody struct {
	cost           int
	sprite         *ebiten.Image
	velocity       float64
	startHP        int
	targetType     string
	enemySpawnTime time.Duration
}

type LevelTemplate struct {
	Stages []*StageTemplate
	BgImg  *ebiten.Image
	Name   string
}

var startPosTypes = []string{"centered", "lines", "checkmate"}

func (l *LevelTemplate) ToLevel() *Level {
	var stages []Stage
	var level *Level = &Level{
		Stages: stages,
		BgImg:  l.BgImg,
		Name:   l.Name,
	}
	for s := range l.Stages {
		level.Stages = append(level.Stages, Stage{StageId: s})
	}
	for s := range level.Stages {
		var items []Item
		level.Stages[s].Items = items
		conVItems := ConvertItems(l.Stages[s].Items)
		// for i := range conVItems {
		// 	level.Stages[s].Items = append(level.Stages[s].Items, conVItems[i])
		// }
		level.Stages[s].Items = append(level.Stages[s].Items, conVItems...)
		level.Stages[s].Items = append(level.Stages[s].Items, ConvertItems(l.Stages[s].Items)...)
		for w := range l.Stages[s].Waves {
			level.Stages[s].Waves = append(level.Stages[s].Waves, Wave{WaveId: w})
		}
	}
	for s, stage := range level.Stages {
		for w := range stage.Waves {
			for _, batch := range l.Stages[s].Waves[w].Batches {
				randPosType := startPosTypes[rand.Intn(len(startPosTypes))]
				enemyCount := len(batch.Enemies) * 10
				if batch.Enemies[0].TargetType == TargetTypePlayer {
					enemyCount = len(batch.Enemies) * 6
				}
				startPosOffset := batch.Enemies[0].Sprite.Bounds().Dx() / 2
				if enemyCount*(startPosOffset+startPosOffset)-int(startPosOffset) <= ScreenWidth1024X768 {
					randPosType = "centered"
				} else if randPosType == "centered" && enemyCount*(startPosOffset+startPosOffset)-int(startPosOffset) > ScreenWidth1024X768 {
					randPosType = "checkmate"
				}
				stage.Waves[w].Batches = append(level.Stages[s].Waves[w].Batches, EnemyBatch{
					Type:              batch.Enemies[0].ToEnemy(),
					TargetType:        batch.Enemies[0].TargetType,
					StartPositionType: randPosType,
					BatchSpawnTime:    batch.Enemies[0].EnemySpawnTime,
					Count:             enemyCount,
					StartPosOffset:    float64(startPosOffset),
				})
			}
		}
	}
	return level
}

type StageTemplate struct {
	Waves []*WaveTemplate
	Items []*ItemTemplate
}
type WaveTemplate struct {
	Batches []*BatchTemplate
}

type BatchTemplate struct {
	Enemies []*EnemyTemplate
}

type EnemyTemplate struct {
	CurCost        int
	Sprite         *ebiten.Image
	Velocity       float64
	EnemySpawnTime time.Duration
	WeaponType     *WeaponType
	StartHP        int
	StartPosOffset float64
	CritChance     float64
	TargetType     string
}

func (e *EnemyTemplate) ToEnemy() *EnemyType {
	return &EnemyType{
		Sprite:         e.Sprite,
		Velocity:       e.Velocity,
		EnemySpawnTime: e.EnemySpawnTime,
		WeaponType:     e.WeaponType,
		StartHP:        e.StartHP,
	}
}

func (e *EnemyTemplate) SetBody(body *EnemyBody) {
	if e.CurCost >= body.cost {
		e.CurCost -= body.cost
		e.Sprite = body.sprite
		e.Velocity = body.velocity
		e.StartHP = body.startHP
		e.TargetType = body.targetType
		e.EnemySpawnTime = body.enemySpawnTime
	}
}

func (e *EnemyTemplate) SetWeapon(w *WeaponType) {
	if e.CurCost >= w.cost {
		e.CurCost -= w.cost
		e.WeaponType = w
	}
}

func (e *EnemyTemplate) AddHP() {
	if e.CurCost >= 200 {
		e.CurCost -= 200
		e.StartHP++
	}
}

func (e *EnemyTemplate) AddVelocity() {
	if e.CurCost >= 12 {
		e.CurCost -= 12
		e.Velocity += 0.1
	}
}

func (e *EnemyTemplate) AddWeaponProjectileVelocity(velocity float64) {
	if e.CurCost >= 100 && e.WeaponType != nil {
		e.CurCost -= 100
		e.WeaponType.Velocity += velocity
		e.WeaponType.StartVelocity += velocity

	}
}

func (e *EnemyTemplate) AddWeaponProjectileFireRate() {
	if e.CurCost >= 200 && e.WeaponType != nil {
		e.CurCost -= 200
		e.WeaponType.StartTime++
	}
}

func (e *EnemyTemplate) AddWeaponProjectileDamage() {
	if e.CurCost >= 140 && e.WeaponType != nil {
		e.CurCost -= 140
		e.WeaponType.Damage++
	}
}

func (e *EnemyTemplate) AddWeaponAmmo() {
	if e.CurCost >= 1 && e.WeaponType != nil {
		e.CurCost--
		e.WeaponType.StartAmmo++
	}
}

func (e *EnemyTemplate) DecreaseCost() {
	if e.CurCost > 0 {
		e.CurCost--
	}
}
