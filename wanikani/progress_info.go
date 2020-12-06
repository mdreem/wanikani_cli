package wanikani

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/data"
)

const NumSrsLevels = 10

type ProgressionData struct {
	NumNotStarted  int
	NumApprentice  int
	NumGuru        int
	NumMaster      int
	NumEnlightened int
	NumBurned      int
}

func ComputeProgressInfo(client data.Client, subjectType string) ProgressionData {
	srsDistribution := fetchDistribution(client, subjectType)
	accumulatedDistribution := make([]int, NumSrsLevels)
	progressionData := computeProgressionData(srsDistribution, accumulatedDistribution)
	return progressionData
}

func computeProgressionData(kanjiSrsDistribution []int, accumulatedDistribution []int) ProgressionData {
	for idx := range kanjiSrsDistribution {
		accumulatedDistribution[NumSrsLevels-idx-1] = kanjiSrsDistribution[NumSrsLevels-idx-1]
		if idx > 0 {
			accumulatedDistribution[NumSrsLevels-idx-1] += accumulatedDistribution[NumSrsLevels-idx]
		}
	}
	progressionData := ProgressionData{
		NumNotStarted:  accumulatedDistribution[0],
		NumApprentice:  accumulatedDistribution[1],
		NumGuru:        accumulatedDistribution[5],
		NumMaster:      accumulatedDistribution[7],
		NumEnlightened: accumulatedDistribution[8],
		NumBurned:      accumulatedDistribution[9],
	}
	return progressionData
}

func fetchDistribution(client data.Client, subjectType string) []int {
	assignments := client.FetchAssignments(nil, []string{subjectType})
	srsDistribution := make([]int, NumSrsLevels)
	for _, assignment := range assignments {
		srsStage, err := assignment.Data.SrsStage.Int64()
		if err != nil {
			panic(fmt.Errorf("could not convert '%v' to int: %v", assignment.Data.SrsStage, err))
		}
		srsDistribution[srsStage]++
	}
	return srsDistribution
}
