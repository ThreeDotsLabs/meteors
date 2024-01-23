package config

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth                 = 1024
	ScreenHeight                = 768
	Screen1024X768XMenuShift    = 120
	Screen1024X768YMenuShift    = 120
	Screen1024X768XProfileShift = 120
	Screen1024X768YProfileShift = 120
	Screen1024X768YMenuHeight   = 40
	Screen1024X768FontHeight    = 20
	Screen1024X768FontWidth     = 16

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
	DoubleLaserCannon   = "doubleLaserCannon"
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
	MainMenu GameState = "mainMenu"
	InGame   GameState = "inGame"
	Profile  GameState = "profile"
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
	Sprite                        *ebiten.Image
	IntercectAnimationSpriteSheet *ebiten.Image
	InstantAnimationSpiteSheet    *ebiten.Image
	AnimationOnly                 bool
	Velocity                      float64
	StartVelocity                 float64
	Damage                        int
	Target                        Vector
	TargetType                    string
	WeaponName                    string
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
