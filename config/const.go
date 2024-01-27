package config

import (
	"astrogame/assets"
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
	WeaponType      *WeaponType
	StartHP         int
}

type LevelTemplate struct {
	Stages []*StageTemplate
	BgImg  *ebiten.Image
	Name   string
}

var startPosTypes = []string{"centered", "lines", "checkmate"}

func (l *LevelTemplate) ToLevel() *Level {
	var level *Level
	for s, _ := range l.Stages {
		level.Stages = append(level.Stages, Stage{StageId: s})
	}
	for s, _ := range level.Stages {
		for w, _ := range l.Stages[s].Waves {
			level.Stages[s].Waves = append(level.Stages[s].Waves, Wave{WaveId: w})
		}
	}
	for s, stage := range level.Stages {
		for w := range stage.Waves {
			for _, batch := range l.Stages[s].Waves[w].Batches {
				randPosType := startPosTypes[rand.Intn(len(startPosTypes))]
				stage.Waves[w].Batches = append(level.Stages[s].Waves[w].Batches, EnemyBatch{
					Type:              batch.Enemies[0].ToEnemy(),
					TargetType:        batch.Enemies[0].TargetType,
					StartPositionType: randPosType,
					BatchSpawnTime:    batch.Enemies[0].EnemySpawnTime,
					Count:             len(batch.Enemies) * 10,
					StartPosOffset:    batch.Enemies[0].StartPosOffset,
				})
			}
		}
	}
	return level
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

type EnemyBody struct {
	cost           int
	sprite         *ebiten.Image
	velocity       float64
	startHP        int
	targetType     string
	enemySpawnTime time.Duration
}

func NewEnemyBodies() []*EnemyBody {
	var enemyBodies []*EnemyBody
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           20,
		sprite:         assets.Enemy1,
		velocity:       1,
		startHP:        1,
		targetType:     TargetTypeStraight,
		enemySpawnTime: 2 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           30,
		sprite:         assets.Enemy2,
		velocity:       1,
		startHP:        3,
		targetType:     TargetTypePlayer,
		enemySpawnTime: 4 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           42,
		sprite:         assets.Enemy3,
		velocity:       1.5,
		startHP:        3,
		targetType:     TargetTypeStraight,
		enemySpawnTime: 4 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           65,
		sprite:         assets.Enemy4,
		velocity:       1.4,
		startHP:        6,
		targetType:     TargetTypePlayer,
		enemySpawnTime: 4 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           66,
		sprite:         assets.Enemy5,
		velocity:       1.8,
		startHP:        4,
		targetType:     TargetTypeStraight,
		enemySpawnTime: 5 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           78,
		sprite:         assets.Enemy6,
		velocity:       1.2,
		startHP:        7,
		targetType:     TargetTypeStraight,
		enemySpawnTime: 5 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           90,
		sprite:         assets.Enemy7,
		velocity:       2,
		startHP:        5,
		targetType:     TargetTypeStraight,
		enemySpawnTime: 5 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           120,
		sprite:         assets.Enemy8,
		velocity:       1.3,
		startHP:        8,
		targetType:     TargetTypePlayer,
		enemySpawnTime: 6 * time.Second,
	})
	return enemyBodies
}

func NewWeaponTypes() []*WeaponType {
	var weaponTypes []*WeaponType
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:          6,
		Sprite:        assets.EnemyLightMissile,
		Velocity:      150,
		StartVelocity: 150,
		Damage:        2,
		TargetType:    TargetTypeStraight,
		AnimationOnly: false,
		StartTime:     1400,
		StartAmmo:     30,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:          20,
		Sprite:        assets.EnemyAutoLightMissile,
		Velocity:      3,
		StartVelocity: 3,
		Damage:        2,
		TargetType:    TargetTypePlayer,
		AnimationOnly: false,
		StartTime:     2000,
		StartAmmo:     5,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:          28,
		Sprite:        assets.EnemyHeavyMissile,
		Velocity:      130,
		StartVelocity: 130,
		Damage:        5,
		TargetType:    TargetTypeStraight,
		AnimationOnly: false,
		StartTime:     1700,
		StartAmmo:     10,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:          42,
		Sprite:        assets.EnemyHeavyMissile,
		Velocity:      4,
		StartVelocity: 4,
		Damage:        5,
		TargetType:    TargetTypePlayer,
		AnimationOnly: false,
		StartTime:     2600,
		StartAmmo:     5,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:          36,
		Sprite:        assets.EnemyGunProjevtile,
		Velocity:      280,
		StartVelocity: 280,
		Damage:        3,
		TargetType:    TargetTypeStraight,
		AnimationOnly: false,
		StartTime:     1600,
		StartAmmo:     16,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:          46,
		Sprite:        assets.EnemyBullet,
		Velocity:      350,
		StartVelocity: 350,
		Damage:        1,
		TargetType:    TargetTypeStraight,
		AnimationOnly: false,
		StartTime:     400,
		StartAmmo:     100,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:          52,
		Sprite:        assets.EnemyMidMissile,
		Velocity:      4,
		StartVelocity: 4,
		Damage:        4,
		TargetType:    TargetTypePlayer,
		AnimationOnly: false,
		StartTime:     1200,
		StartAmmo:     6,
	})
	return weaponTypes
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
	if e.CurCost >= 20 {
		e.CurCost -= 20
		e.StartHP++
	}
}

func (e *EnemyTemplate) AddVelocity() {
	if e.CurCost >= 12 {
		e.CurCost -= 12
		e.Velocity++
	}
}

func (e *EnemyTemplate) AddWeaponProjectileVelocity() {
	if e.CurCost >= 1 && e.WeaponType != nil {
		e.CurCost--
		e.WeaponType.Velocity++
		e.WeaponType.StartVelocity++

	}
}

func (e *EnemyTemplate) AddWeaponProjectileFireRate() {
	if e.CurCost >= 1 && e.WeaponType != nil {
		e.CurCost--
		e.WeaponType.StartTime++
	}
}

func (e *EnemyTemplate) AddWeaponProjectileDamage() {
	if e.CurCost >= 40 && e.WeaponType != nil {
		e.CurCost -= 40
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
