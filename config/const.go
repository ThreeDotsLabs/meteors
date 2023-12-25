package config

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 1024
	ScreenHeight = 768

	MeteorSpawnTime     = 2 * time.Second
	EnemySpawnTime      = 1 * time.Second
	BaseMeteorVelocity  = 0.25
	MeteorSpeedUpAmount = 0.1
	MeteorSpeedUpTime   = 5 * time.Second
	LightRocket         = "lightRocket"
	AutoLightRocket     = "autoLightRocket"
)

type EnemyType struct {
	RotationSpeed   float64
	Sprite          *ebiten.Image
	Velocity        float64
	EnemiesStartPos Vector
	EnemySpawnTime  time.Duration
	WeaponTypeStr   string
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
	ItemSpawnTime time.Duration
}
