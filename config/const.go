package config

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth               = 1024
	ScreenHeight              = 768
	Screen1024X768XMenuShift  = 120
	Screen1024X768YMenuShift  = 120
	Screen1024X768YMenuHeight = 60

	MeteorSpawnTime     = 2 * time.Second
	EnemySpawnTime      = 1 * time.Second
	BaseMeteorVelocity  = 0.25
	MeteorSpeedUpAmount = 0.1
	MeteorSpeedUpTime   = 5 * time.Second
	TargetTypePlayer    = "player"
	LightRocket         = "lightRocket"
	AutoLightRocket     = "autoLightRocket"
	DoubleLightRocket   = "doubleLightRocket"
	LaserCannon         = "lightCannon"
	ClusterMines        = "clusterMines"
	BigBomb             = "bigBomb"
	MachineGun          = "machineGun"
)

type GameState string

var (
	MainMenu GameState = "mainMenu"
	InGame   GameState = "inGame"
)

type EnemyType struct {
	RotationSpeed   float64
	Sprite          *ebiten.Image
	Velocity        float64
	EnemiesStartPos Vector
	EnemySpawnTime  time.Duration
	WeaponTypeStr   string
	StartHP         int
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
	WaveId       int
	EnemiesCount int
	EnemyType    *EnemyType
	Batches      []EnemyBatch
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
	Sprite     *ebiten.Image
	Velocity   float64
	Damage     int
	Target     Vector
	TargetType string
	WeaponName string
}

type AmmoType struct {
	WeaponType string
	Amount     int
}

type HealType struct {
	HP int
}

type Item struct {
	AmmoType      *AmmoType
	WeaponType    *WeaponType
	HealType      *HealType
	RotationSpeed float64
	Sprite        *ebiten.Image
	Velocity      float64
	Target        Vector
	ItemSpawnTime time.Duration
}
