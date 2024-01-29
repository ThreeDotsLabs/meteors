package config

import (
	"astrogame/assets"
	"astrogame/objects"
	"time"
)

func NewEnemyBodies() []*EnemyBody {
	var enemyBodies []*EnemyBody
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           20,
		sprite:         assets.Enemy1,
		velocity:       1,
		startHP:        1,
		targetType:     TargetTypeStraight,
		enemySpawnTime: 4 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           30,
		sprite:         assets.Enemy2,
		velocity:       1.2,
		startHP:        3,
		targetType:     TargetTypePlayer,
		enemySpawnTime: 5 * time.Second,
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
		cost:           76,
		sprite:         assets.Enemy4,
		velocity:       1.4,
		startHP:        6,
		targetType:     TargetTypePlayer,
		enemySpawnTime: 4 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           120,
		sprite:         assets.Enemy5,
		velocity:       1.8,
		startHP:        4,
		targetType:     TargetTypeStraight,
		enemySpawnTime: 5 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           178,
		sprite:         assets.Enemy6,
		velocity:       1.2,
		startHP:        7,
		targetType:     TargetTypeStraight,
		enemySpawnTime: 5 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           290,
		sprite:         assets.Enemy7,
		velocity:       2,
		startHP:        5,
		targetType:     TargetTypeStraight,
		enemySpawnTime: 5 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           520,
		sprite:         assets.Enemy8,
		velocity:       1.3,
		startHP:        8,
		targetType:     TargetTypePlayer,
		enemySpawnTime: 6 * time.Second,
	})
	enemyBodies = append(enemyBodies, &EnemyBody{
		cost:           640,
		sprite:         assets.Enemy9,
		velocity:       1.5,
		startHP:        12,
		targetType:     TargetTypePlayer,
		enemySpawnTime: 8 * time.Second,
	})
	return enemyBodies
}

func NewWeaponTypes() []*WeaponType {
	var weaponTypes []*WeaponType
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:                          6,
		Sprite:                        assets.EnemyLightMissile,
		IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
		Velocity:                      150,
		StartVelocity:                 150,
		Damage:                        2,
		TargetType:                    TargetTypeStraight,
		AnimationOnly:                 false,
		StartTime:                     time.Duration(1400) * time.Millisecond,
		StartAmmo:                     30,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:                          20,
		Sprite:                        assets.EnemyAutoLightMissile,
		IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
		Velocity:                      1.5,
		StartVelocity:                 1.5,
		Damage:                        2,
		TargetType:                    TargetTypePlayer,
		AnimationOnly:                 false,
		StartTime:                     time.Duration(2000) * time.Millisecond,
		StartAmmo:                     5,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:                          40,
		Sprite:                        assets.EnemyHeavyMissile,
		IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
		Velocity:                      130,
		StartVelocity:                 130,
		Damage:                        5,
		TargetType:                    TargetTypeStraight,
		AnimationOnly:                 false,
		StartTime:                     time.Duration(1700)*time.Millisecond + 1700,
		StartAmmo:                     10,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:                          70,
		Sprite:                        assets.EnemyHeavyMissile,
		IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
		Velocity:                      2.2,
		StartVelocity:                 2.2,
		Damage:                        5,
		TargetType:                    TargetTypePlayer,
		AnimationOnly:                 false,
		StartTime:                     2600,
		StartAmmo:                     5,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:                          56,
		Sprite:                        assets.EnemyGunProjevtile,
		IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
		Velocity:                      280,
		StartVelocity:                 280,
		Damage:                        3,
		TargetType:                    TargetTypeStraight,
		AnimationOnly:                 false,
		StartTime:                     1600,
		StartAmmo:                     16,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:                          86,
		Sprite:                        assets.EnemyBullet,
		IntercectAnimationSpriteSheet: assets.ProjectileBlowSpriteSheet,
		Velocity:                      350,
		StartVelocity:                 350,
		Damage:                        1,
		TargetType:                    TargetTypeStraight,
		AnimationOnly:                 false,
		StartTime:                     400,
		StartAmmo:                     100,
	})
	weaponTypes = append(weaponTypes, &WeaponType{
		cost:                          100,
		Sprite:                        assets.EnemyMidMissile,
		IntercectAnimationSpriteSheet: assets.LightMissileBlowSpriteSheet,
		Velocity:                      2.4,
		StartVelocity:                 2.4,
		Damage:                        4,
		TargetType:                    TargetTypePlayer,
		AnimationOnly:                 false,
		StartTime:                     1200,
		StartAmmo:                     6,
	})
	return weaponTypes
}

func NewItemTypes(l int) []*ItemTemplate {
	var itemTypes []*ItemTemplate
	if l < 3 {
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        assets.ItemMissileSprite,
			Velocity:      1.2,
			ItemSpawnTime: 5 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: LightRocket,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        assets.ItemDoubleMissileSprite,
			Velocity:      1.3,
			ItemSpawnTime: 6 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: DoubleLightRocket,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        assets.ItemLaserCanonSprite,
			Velocity:      1.5,
			ItemSpawnTime: 7 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: LaserCanon,
			},
		})

		if l > 0 {
			for idx := 1; idx < l*2; idx++ {
				itemTypes = append(itemTypes, &ItemTemplate{
					Sprite:        objects.ScaleImg(assets.Heal, 0.5),
					Velocity:      1.4,
					ItemSpawnTime: 6 * time.Second,
					HealType: &HealType{
						HP: l * 5,
					},
				})
			}
			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.MachineGun, 1),
				Velocity:      1.4,
				ItemSpawnTime: 6 * time.Second,
				WeaponType: &WeaponType{
					WeaponName: MachineGun,
				},
			})
		}

		if l > 1 {
			for idx := 1; idx < l; idx++ {
				itemTypes = append(itemTypes, &ItemTemplate{
					Sprite:        objects.ScaleImg(assets.ShieldSprite, 0.8),
					Velocity:      1.8,
					ItemSpawnTime: 5 * time.Second,
					ShieldType: &ShieldType{
						HP:     l * 4,
						Sprite: assets.ShieldSprite,
					},
				})
			}

			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.ClusterMines, 1),
				Velocity:      1.5,
				ItemSpawnTime: 5 * time.Second,
				SecondWeaponType: &WeaponType{
					WeaponName: ClusterMines,
				},
			})
		}
	} else if l >= 3 && l <= 6 {
		for idx := 1; idx < l*2; idx++ {
			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.Heal, 0.5),
				Velocity:      1.4,
				ItemSpawnTime: 6 * time.Second,
				HealType: &HealType{
					HP: 8 * (l / 2),
				},
			})
		}
		for idx := 1; idx < l; idx++ {
			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.ShieldSprite, 0.8),
				Velocity:      1.8,
				ItemSpawnTime: 5 * time.Second,
				ShieldType: &ShieldType{
					HP:     5 * (l / 2),
					Sprite: assets.ShieldSprite,
				},
			})
		}
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        assets.ItemMissileSprite,
			Velocity:      1.2,
			ItemSpawnTime: 5 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: LightRocket,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.MissileSprite, 0.75),
			Velocity:      1.4,
			ItemSpawnTime: 5 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: LightRocket,
				Amount:     75 * (l / 2),
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        assets.ItemDoubleMissileSprite,
			Velocity:      1.3,
			ItemSpawnTime: 6 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: DoubleLightRocket,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.DoubleMissileSprite, 0.75),
			Velocity:      1.4,
			ItemSpawnTime: 7 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: DoubleLightRocket,
				Amount:     50 * l,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        assets.ItemLaserCanonSprite,
			Velocity:      1.5,
			ItemSpawnTime: 7 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: LaserCanon,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.LaserCanon, 0.75),
			Velocity:      1.6,
			ItemSpawnTime: 7 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: LaserCanon,
				Amount:     48 * (l / 2),
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        assets.ItemDoubleLaserCanonSprite,
			Velocity:      1.5,
			ItemSpawnTime: 7 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: DoubleLaserCanon,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.MachineGun, 1),
			Velocity:      1.4,
			ItemSpawnTime: 6 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: MachineGun,
			},
		})

		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.ClusterMines, 1),
			Velocity:      1.6,
			ItemSpawnTime: 7 * time.Second,
			SecondWeaponType: &WeaponType{
				WeaponName: ClusterMines,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.ClusterMines, 0.75),
			Velocity:      1.7,
			ItemSpawnTime: 8 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: ClusterMines,
				Amount:     20 * (l / 2),
			},
		})

		if l > 4 {
			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.PlasmaGun, 1),
				Velocity:      1.4,
				ItemSpawnTime: 6 * time.Second,
				WeaponType: &WeaponType{
					WeaponName: PlasmaGun,
				},
			})

			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.BigBomb, 1),
				Velocity:      1.6,
				ItemSpawnTime: 7 * time.Second,
				SecondWeaponType: &WeaponType{
					WeaponName: BigBomb,
				},
			})
		}
	} else if l > 6 && l <= 10 {
		for idx := 1; idx < l; idx++ {
			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.Heal, 0.5),
				Velocity:      1.4,
				ItemSpawnTime: 6 * time.Second,
				HealType: &HealType{
					HP: 12 * (l / 3),
				},
			})
		}
		for idx := 1; idx < l/2; idx++ {
			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.ShieldSprite, 0.8),
				Velocity:      1.8,
				ItemSpawnTime: 5 * time.Second,
				ShieldType: &ShieldType{
					HP:     14 * (l / 3),
					Sprite: assets.ShieldSprite,
				},
			})
		}
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.MissileSprite, 1),
			Velocity:      1.5,
			ItemSpawnTime: 5 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: LightRocket,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.MissileSprite, 0.75),
			Velocity:      1.6,
			ItemSpawnTime: 6 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: LightRocket,
				Amount:     100 * (l / 3),
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        assets.ItemDoubleMissileSprite,
			Velocity:      1.4,
			ItemSpawnTime: 6 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: DoubleLightRocket,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.DoubleMissileSprite, 0.75),
			Velocity:      1.6,
			ItemSpawnTime: 7 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: DoubleLightRocket,
				Amount:     80 * (l / 3),
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        assets.ItemLaserCanonSprite,
			Velocity:      1.5,
			ItemSpawnTime: 6 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: LaserCanon,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.LaserCanon, 0.75),
			Velocity:      1.6,
			ItemSpawnTime: 7 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: LaserCanon,
				Amount:     60 * (l / 3),
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        assets.ItemDoubleLaserCanonSprite,
			Velocity:      1.6,
			ItemSpawnTime: 7 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: DoubleLaserCanon,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.DoubleLaserCanon, 0.75),
			Velocity:      1.6,
			ItemSpawnTime: 7 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: DoubleLaserCanon,
				Amount:     60 * (l / 3),
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.MachineGun, 1),
			Velocity:      1.5,
			ItemSpawnTime: 6 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: MachineGun,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.MachineGun, 0.75),
			Velocity:      1.6,
			ItemSpawnTime: 7 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: MachineGun,
				Amount:     60 * (l / 3),
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.PlasmaGun, 1),
			Velocity:      1.6,
			ItemSpawnTime: 6 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: PlasmaGun,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.PlasmaGun, 0.75),
			Velocity:      1.7,
			ItemSpawnTime: 7 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: PlasmaGun,
				Amount:     30 * (l / 3),
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.DoubleMachineGun, 1),
			Velocity:      1.7,
			ItemSpawnTime: 8 * time.Second,
			WeaponType: &WeaponType{
				WeaponName: DoubleMachineGun,
			},
		})

		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.ClusterMines, 1),
			Velocity:      1.6,
			ItemSpawnTime: 8 * time.Second,
			SecondWeaponType: &WeaponType{
				WeaponName: ClusterMines,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.ClusterMines, 0.75),
			Velocity:      1.7,
			ItemSpawnTime: 9 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: ClusterMines,
				Amount:     36 * (l / 3),
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.BigBomb, 1),
			Velocity:      1.8,
			ItemSpawnTime: 8 * time.Second,
			SecondWeaponType: &WeaponType{
				WeaponName: BigBomb,
			},
		})
		itemTypes = append(itemTypes, &ItemTemplate{
			Sprite:        objects.ScaleImg(assets.BigBomb, 0.75),
			Velocity:      1.9,
			ItemSpawnTime: 10 * time.Second,
			AmmoType: &AmmoType{
				WeaponName: BigBomb,
				Amount:     20 * (l / 3),
			},
		})

		if l > 7 {
			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.ShieldSprite, 1),
				Velocity:      1.9,
				ItemSpawnTime: 7 * time.Second,
				ShieldType: &ShieldType{
					HP:     26 * (l / 2),
					Sprite: assets.ShieldSprite,
				},
			})

			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.PlasmaGun, 1),
				Velocity:      1.6,
				ItemSpawnTime: 8 * time.Second,
				WeaponType: &WeaponType{
					WeaponName: DoublePlasmaGun,
				},
			})
			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.PlasmaGun, 0.75),
				Velocity:      1.8,
				ItemSpawnTime: 9 * time.Second,
				AmmoType: &AmmoType{
					WeaponName: DoublePlasmaGun,
					Amount:     20 * (l / 3),
				},
			})

			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.PentaLaser, 1),
				Velocity:      1.6,
				ItemSpawnTime: 8 * time.Second,
				SecondWeaponType: &WeaponType{
					WeaponName: PentaLaser,
				},
			})
			itemTypes = append(itemTypes, &ItemTemplate{
				Sprite:        objects.ScaleImg(assets.PentaLaser, 0.75),
				Velocity:      1.7,
				ItemSpawnTime: 9 * time.Second,
				AmmoType: &AmmoType{
					WeaponName: PentaLaser,
					Amount:     36 * (l / 3),
				},
			})
		}
	}
	return itemTypes
}
