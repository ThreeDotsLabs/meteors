package config

import (
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
	TargetTypePlayer    = "player"
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
	MainMenu GameState = "mainMenu"
	InGame   GameState = "inGame"
	Profile  GameState = "profile"
	Options  GameState = "options"
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

type LevelTemplate struct {
	Stages []*StageTemplate
	BgImg  *ebiten.Image
	Name   string
}

type StageTemplate struct {
	Waves []*WaveTemplate
}
type WaveTemplate struct {
	Batches []*BatchTemplate
}

type BatchTemplate struct {
	Enemies []*EnemyTemplate
}

type EnemyTemplate struct {
	Sprite         *ebiten.Image
	Velocity       float64
	EnemySpawnTime time.Duration
	WeaponTypeStr  string
	StartHP        int
	StartPosOffset float64
}

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
