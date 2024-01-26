package game

import (
	"astrogame/config"
	"astrogame/objects"
)

func GenerateLevelsStructure() []*config.LevelTemplate {
	var structure []*config.LevelTemplate
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
		structure = append(structure, &levels)
	}

	return structure
}

func generateWaves(count int) []*config.WaveTemplate {
	var waves []*config.WaveTemplate
	for w := 0; w < count; w++ {
		var wave config.WaveTemplate
		waves = append(waves, &wave)
	}
	return waves
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
		var stage config.StageTemplate
		waves := generateWaves(randWaveCount)
		stage.Waves = waves
		stages = append(stages, &stage)
	}
	return stages
}
