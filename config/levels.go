package config

import (
	"astrogame/assets"
	"time"
)

func NewLevels() []*Level {
	var lvl = Level{
		Name:    "Level 1",
		Number:  1,
		LevelId: 0,
		Stages: []Stage{
			{
				StageId:      0,
				MeteorsCount: 0,
				Items: []Item{
					{
						RotationSpeed: 0,
						Sprite:        assets.DoubleLaserCannon,
						Velocity:      1.2,
						WeaponType: &WeaponType{
							WeaponName: DoubleLaserCannon,
						},
						ItemSpawnTime: 1 * time.Second,
					},
					{
						RotationSpeed: 0,
						Sprite:        assets.ClusterMines,
						Velocity:      1.5,
						WeaponType: &WeaponType{
							WeaponName: ClusterMines,
						},
						ItemSpawnTime: 1 * time.Second,
					},
					{
						RotationSpeed: 0,
						Sprite:        assets.ShieldSprite,
						Velocity:      2,
						ShieldType: &ShieldType{
							HP:     5,
							Sprite: assets.ShieldSprite,
						},
						ItemSpawnTime: 1 * time.Second,
					},
					// {
					// 	RotationSpeed: 0,
					// 	Sprite:        assets.MachineGun,
					// 	Velocity:      1.7,
					// 	WeaponType: &WeaponType{
					// 		WeaponName: MachineGun,
					// 	},
					// 	ItemSpawnTime: 5 * time.Second,
					// },
					{
						RotationSpeed: 0,
						Sprite:        assets.BigBomb,
						Velocity:      2,
						WeaponType: &WeaponType{
							WeaponName: BigBomb,
						},
						ItemSpawnTime: 20 * time.Second,
					},
					// {
					// 	RotationSpeed: 0.2,
					// 	Sprite:        assets.MissileSprite,
					// 	Velocity:      1,
					// 	AmmoType: &AmmoType{
					// 		WeaponType: LightRocket,
					// 		Amount:     50,
					// 	},
					// 	ItemSpawnTime: 1 * time.Second,
					// },
					// {
					// 	RotationSpeed: 0,
					// 	Sprite:        assets.DoubleMissileSprite,
					// 	Velocity:      1.3,
					// 	WeaponType: &WeaponType{
					// 		WeaponName: DoubleLightRocket,
					// 	},
					// 	ItemSpawnTime: 5 * time.Second,
					// },
				},
				Waves: []Wave{
					{
						WaveId: 0,
						Batches: []EnemyBatch{
							{
								Type: &EnemyType{
									RotationSpeed: 0,
									Sprite:        assets.HighSpeedFollowPlayerEnemySprite,
									Velocity:      2,
									StartHP:       6,
								},
								Count:             6,
								TargetType:        TargetTypePlayer,
								BatchSpawnTime:    1 * time.Second,
								StartPositionType: "centered",
								StartPosOffset:    40.0,
							},
							{
								Type: &EnemyType{
									RotationSpeed: 0,
									Sprite:        assets.LowSpeedEnemyLightMissile,
									Velocity:      1,
									WeaponTypeStr: LightRocket,
									StartHP:       1,
								},
								Count:             20,
								TargetType:        "straight",
								BatchSpawnTime:    5 * time.Second,
								StartPositionType: "checkmate",
							},
							{
								Type: &EnemyType{
									RotationSpeed: 0,
									Sprite:        assets.LowSpeedEnemyAutoLightMissile,
									Velocity:      1,
									WeaponTypeStr: AutoLightRocket,
									StartHP:       2,
								},
								Count:             4,
								TargetType:        "straight",
								BatchSpawnTime:    10 * time.Second,
								StartPositionType: "checkmate",
							},
						},
					},
					{
						WaveId: 1,
						Batches: []EnemyBatch{
							{
								Type: &EnemyType{
									RotationSpeed: 0,
									Sprite:        assets.HighSpeedFollowPlayerEnemySprite,
									Velocity:      2,
									StartHP:       3,
								},
								Count:             8,
								TargetType:        TargetTypePlayer,
								BatchSpawnTime:    1 * time.Second,
								StartPositionType: "lines",
								StartPosOffset:    20.0,
							},
							{
								Type: &EnemyType{
									RotationSpeed: 0,
									Sprite:        assets.LowSpeedEnemyLightMissile,
									Velocity:      1,
									WeaponTypeStr: LightRocket,
									StartHP:       3,
								},
								Count:             10,
								TargetType:        "straight",
								BatchSpawnTime:    5 * time.Second,
								StartPositionType: "checkmate",
							},
							// {
							// 	Type: &EnemyType{
							// 		RotationSpeed: 0,
							// 		Sprite:        assets.LowSpeedEnemyAutoLightMissile,
							// 		Velocity:      1,
							// 		WeaponTypeStr: AutoLightRocket,
							// 		StartHP:       3,
							// 	},
							// 	Count:             5,
							// 	TargetType:        "straight",
							// 	BatchSpawnTime:    10 * time.Second,
							// 	StartPositionType: "checkmate",
							// },
							// {
							// 	Type: &EnemyType{
							// 		RotationSpeed: 0,
							// 		Sprite:        assets.LowSpeedEnemyAutoLightMissile,
							// 		Velocity:      1,
							// 		WeaponTypeStr: AutoLightRocket,
							// 		StartHP:       3,
							// 	},
							// 	Count:             5,
							// 	TargetType:        "straight",
							// 	BatchSpawnTime:    10 * time.Second,
							// 	StartPositionType: "checkmate",
							// },
						},
					},
				},
			},
		},
		BgImg: assets.FirstLevelBg,
	}
	lvls := []*Level{
		&lvl,
	}
	return lvls
}

var Levels = NewLevels()
