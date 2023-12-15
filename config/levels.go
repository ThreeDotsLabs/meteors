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
				Waves: []Wave{
					{
						WaveId: 0,
						Batches: []EnemyBatch{
							{
								Type: &EnemyType{
									RotationSpeed: 0,
									Sprite:        assets.HighSpeedFollowPlayerEnemySprite,
									Velocity:      2,
								},
								Count:             5,
								TargetType:        "player",
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
								},
								Count:             10,
								TargetType:        "straight",
								BatchSpawnTime:    40 * time.Second,
								StartPositionType: "checkmate",
							},
							{
								Type: &EnemyType{
									RotationSpeed: 0,
									Sprite:        assets.LowSpeedEnemyAutoLightMissile,
									Velocity:      1,
									WeaponTypeStr: AutoLightRocket,
								},
								Count:             4,
								TargetType:        "straight",
								BatchSpawnTime:    100 * time.Second,
								StartPositionType: "checkmate",
							},
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
