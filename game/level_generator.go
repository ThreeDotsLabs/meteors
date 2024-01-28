package game

import (
	"astrogame/assets"
	"astrogame/config"
	"astrogame/objects"
	"fmt"
	"math"
	"time"
)

type levelTemplatesGen struct {
	lvls []*config.LevelTemplate
}

func GenerateLevels() []*config.Level {
	var levels []*config.Level
	lGen := GenerateLevelStructure()
	lGen.DecorateLevels()
	for _, l := range lGen.lvls {
		levels = append(levels, l.ToLevel())
	}
	return levels
}
func (lvlTpl levelTemplatesGen) DecorateLevels() {
	for i, l := range lvlTpl.lvls {
		l.Name = "Level " + fmt.Sprint(i+1)
		for k, s := range l.Stages {
			for j, w := range s.Waves {
				for _, b := range w.Batches {
					costLvlIdx := i + 1
					costStageIdx := k + 1
					costWaveIdx := j + 1
					cost := costLvlIdx*20 + 10*costStageIdx + costWaveIdx*2
					for _, e := range b.Enemies {
						DecorateEnemyTemplate(e, i, k, j, cost, &lvlTpl)
					}
				}
			}
		}
	}
}

func DecorateEnemyTemplate(e *config.EnemyTemplate, l int, s int, w int, cost int, lTpl *levelTemplatesGen) {
	// thirdPart := s / 3
	// wavesThirdPart := w / 3
	// costLvlIdx := l
	// costStageIdx := s
	// costWaveIdx := w
	// if l == 0 {
	// 	costLvlIdx = 1
	// }
	// if s == 0 {
	// 	costStageIdx = 1
	// }
	// if w == 0 {
	// 	costWaveIdx = 1
	// }
	// if l <= 3 {
	// 	if s <= thirdPart {
	// 		if w <= wavesThirdPart {
	// 			e.CurCost = cost //+ 600/costStageIdx
	// 		} else if w > wavesThirdPart && w <= wavesThirdPart*2 {
	// 			e.CurCost = cost //+ 600/costStageIdx
	// 		} else if w > wavesThirdPart*2 {
	// 			e.CurCost = cost //+ 600/costStageIdx
	// 		}
	// 	} else if s > thirdPart && s <= thirdPart*2 {
	// 		if w <= wavesThirdPart {
	// 			e.CurCost = cost //+ 900/costStageIdx
	// 		} else if w > wavesThirdPart && w <= wavesThirdPart*2 {
	// 			e.CurCost = cost //+ 900/costStageIdx
	// 		} else if w > wavesThirdPart*2 {
	// 			e.CurCost = cost //+ 900/costStageIdx
	// 		}
	// 	} else if s > thirdPart*2 {
	// 		if w <= wavesThirdPart {
	// 			e.CurCost = cost //+ 1600/costStageIdx
	// 		} else if w > wavesThirdPart && w <= wavesThirdPart*2 {
	// 			e.CurCost = cost //+ 1600/costStageIdx
	// 		} else if w > wavesThirdPart*2 {
	// 			e.CurCost = cost //+ 1600/costStageIdx
	// 		}
	// 	}
	// } else if l > 3 && l <= 7 {
	// 	if s <= thirdPart {
	// 		if w <= wavesThirdPart {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart && w <= wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		}
	// 	} else if s > thirdPart && s <= thirdPart*2 {
	// 		if w <= wavesThirdPart {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart && w <= wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		}
	// 	} else if s > thirdPart*2 {
	// 		if w <= wavesThirdPart {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart && w <= wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		}
	// 	}
	// } else if l > 7 {
	// 	if s <= thirdPart {
	// 		if w <= wavesThirdPart {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart && w <= wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		}
	// 	} else if s > thirdPart && s <= thirdPart*2 {
	// 		if w <= wavesThirdPart {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart && w <= wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		}
	// 	} else if s > thirdPart*2 {
	// 		if w <= wavesThirdPart {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart && w <= wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		} else if w > wavesThirdPart*2 {
	// 			e.CurCost = cost
	// 		}
	// 	}
	// }
	e.CurCost = cost
	bodyAdded := false
	weaponAdded := false
	weaponMinCost := 6
	eBodies := config.NewEnemyBodies()
	eWeapons := config.NewWeaponTypes()
	for {
		if e.CurCost <= 0 {
			break
		}
		if !bodyAdded {
			randBodyCount := objects.RandInt(0, len(eBodies)-1)
			e.SetBody(eBodies[randBodyCount])
			if e.StartHP != 0 {
				bodyAdded = true
			}
		}
		if bodyAdded && e.CurCost >= weaponMinCost && !weaponAdded {
			randWeaponCount := objects.RandInt(0, len(eWeapons)-1)
			e.SetWeapon(eWeapons[randWeaponCount])
			weaponAdded = true
		}
		if bodyAdded {
			e.AddHP()
			e.AddVelocity()
		}
		if bodyAdded && weaponAdded {
			e.AddWeaponProjectileDamage()
			e.AddWeaponProjectileVelocity()
			e.AddWeaponProjectileFireRate()
			e.AddWeaponAmmo()
		}
		if bodyAdded {
			e.DecreaseCost()
		}
	}
}
func GenerateLevelStructure() levelTemplatesGen {
	var structure levelTemplatesGen
	for l := 0; l < 10; l++ {
		var stageCountLLimit int
		var stageCountRLimit int
		if l <= 3 {
			stageCountLLimit = 3
			stageCountRLimit = 5
		} else if l > 3 && l <= 7 {
			stageCountLLimit = 4
			stageCountRLimit = 7
		} else if l > 7 {
			stageCountLLimit = 5
			stageCountRLimit = 8
		}
		randStageCount := objects.RandInt(stageCountLLimit, stageCountRLimit)
		stages := generateStages(l, randStageCount)
		var levels config.LevelTemplate
		levels.Stages = append(levels.Stages, stages...)
		structure.lvls = append(structure.lvls, &levels)
	}
	for i, l := range structure.lvls {
		for _, s := range l.Stages {
			for _, w := range s.Waves {
				var batchCountLLimit int
				var batchCountRLimit int
				if i <= 3 {
					batchCountLLimit = 1
					batchCountRLimit = 3
				} else if i > 3 && i <= 7 {
					batchCountLLimit = 2
					batchCountRLimit = 5
				} else if i > 7 {
					batchCountLLimit = 4
					batchCountRLimit = 7
				}
				randWaveCount := objects.RandInt(batchCountLLimit, batchCountRLimit)
				w.Batches = generateBatches(i, randWaveCount)
			}
		}
	}
	return structure
}

func generateStages(l int, count int) []*config.StageTemplate {
	var stages []*config.StageTemplate
	thirdPart := count / 3
	for w := 0; w < count; w++ {
		var waveCountLLimit int
		var waveCountRLimit int
		if l <= 3 {
			if count <= thirdPart {
				waveCountLLimit = 3
				waveCountRLimit = 4
			} else if count > thirdPart && count <= thirdPart*2 {
				waveCountLLimit = 4
				waveCountRLimit = 6
			} else if count > thirdPart*2 {
				waveCountLLimit = 5
				waveCountRLimit = 6
			}
		} else if l > 3 && l <= 7 {
			if count <= thirdPart {
				waveCountLLimit = 4
				waveCountRLimit = 6
			} else if count > thirdPart && count <= thirdPart*2 {
				waveCountLLimit = 5
				waveCountRLimit = 8
			} else if count > thirdPart*2 {
				waveCountLLimit = 8
				waveCountRLimit = 10
			}
		} else if l > 7 {
			if count <= thirdPart {
				waveCountLLimit = 4
				waveCountRLimit = 8
			} else if count > thirdPart && count <= thirdPart*2 {
				waveCountLLimit = 6
				waveCountRLimit = 12
			} else if count > thirdPart*2 {
				waveCountLLimit = 10
				waveCountRLimit = 18
			}
		}
		randWaveCount := objects.RandInt(waveCountLLimit, waveCountRLimit)
		randItemCount := objects.RandInt(waveCountLLimit/2, waveCountRLimit/2)
		var stage config.StageTemplate
		waves := generateWaves(randWaveCount)
		stage.Waves = waves
		items := generateItems(l, randItemCount)
		stage.Items = items
		stages = append(stages, &stage)
	}
	return stages
}

func generateWaves(count int) []*config.WaveTemplate {
	var waves []*config.WaveTemplate
	for w := 0; w < count; w++ {
		var wave config.WaveTemplate
		waves = append(waves, &wave)
	}
	return waves
}

func generateBatches(l int, count int) []*config.BatchTemplate {
	var batches []*config.BatchTemplate
	for w := 0; w < count; w++ {
		var batch config.BatchTemplate
		batches = append(batches, &batch)
	}

	for _, b := range batches {
		minEnemyCount := 1
		enemyBatchIdx := 3
		if l > 3 && l <= 6 {
			enemyBatchIdx = 5
			minEnemyCount = 3
		} else if l > 6 && l <= 10 {
			enemyBatchIdx = 6
			minEnemyCount = 3
		} else if l > 10 {
			enemyBatchIdx = 8
			minEnemyCount = 4
		}
		randEnemyCount := objects.RandInt(minEnemyCount, enemyBatchIdx)
		b.Enemies = generateEnemies(randEnemyCount)
	}
	return batches
}

func generateEnemies(count int) []*config.EnemyTemplate {
	var enemies []*config.EnemyTemplate
	for w := 0; w < count; w++ {
		var enemy config.EnemyTemplate
		enemies = append(enemies, &enemy)
	}
	return enemies
}

func generateItems(l int, count int) []*config.ItemTemplate {
	var items []*config.ItemTemplate
	for w := 0; w < count; w++ {
		itemsForLvl := config.NewItemTypes(l)
		itemRandNumber := objects.RandInt(0, len(itemsForLvl)-1)
		items = append(items, itemsForLvl[itemRandNumber])
	}
	return items
}

var enemyLightRocket = Weapon{
	projectile: Projectile{
		wType: &config.WeaponType{
			Sprite:                        objects.ScaleImg(assets.EnemyLightMissile, 0.7),
			IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
			Velocity:                      150,
			Damage:                        1,
			TargetType:                    "straight",
			StartTime:                     time.Duration(1000) * time.Millisecond,
			StartAmmo:                     10,
		},
	},
	ammo: 10,
	EnemyShoot: func(e *Enemy) {
		bounds := e.enemyType.Sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		spawnPos := config.Vector{
			X: e.position.X + halfW + math.Sin(e.rotation)*bulletSpawnOffset,
			Y: e.position.Y + halfH + math.Cos(e.rotation)*bulletSpawnOffset,
		}
		animation := NewAnimation(config.Vector{}, e.weapon.projectile.wType.IntercectAnimationSpriteSheet, 1, 56, 60, false, "projectileBlow", 0)
		projectile := NewProjectile(e.game, spawnPos, e.rotation, e.weapon.projectile.wType, animation, 0)
		projectile.owner = config.OwnerEnemy
		e.game.AddProjectile(projectile)
	},
}

var enemyAutoLightRocket = Weapon{
	projectile: Projectile{
		wType: &config.WeaponType{
			Sprite:                        assets.EnemyAutoLightMissile,
			IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
			Velocity:                      3,
			Damage:                        5,
			TargetType:                    config.TargetTypePlayer,
		},
	},
	ammo: 3,
	EnemyShoot: func(e *Enemy) {
		bounds := e.enemyType.Sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		spawnPos := config.Vector{
			X: e.position.X + halfW + math.Sin(e.rotation)*bulletSpawnOffset,
			Y: e.position.Y + halfH + math.Cos(e.rotation)*bulletSpawnOffset,
		}
		animation := NewAnimation(config.Vector{}, e.weapon.projectile.wType.IntercectAnimationSpriteSheet, 1, 56, 60, false, "projectileBlow", 0)
		projectile := NewProjectile(e.game, spawnPos, e.rotation, e.weapon.projectile.wType, animation, 0)
		projectile.owner = config.OwnerEnemy
		e.game.AddProjectile(projectile)
	},
}
