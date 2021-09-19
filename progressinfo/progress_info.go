package progressinfo

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/wanikani"
	"github.com/spf13/cobra"
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

func PrintProgressInfo(client wanikani.RealClient, _ *cobra.Command, _ []string) {
	radicalProgressInfo := ComputeProgressInfo(client, "radical")
	kanjiProgressInfo := ComputeProgressInfo(client, "kanji")
	vocabularyProgressInfo := ComputeProgressInfo(client, "vocabulary")

	fmt.Printf("Radical Progress:\n")
	printSrsLevelProgress(radicalProgressInfo)

	fmt.Printf("\nKanji Progress:\n")
	printSrsLevelProgress(kanjiProgressInfo)

	fmt.Printf("\nVocabulary Progress:\n")
	printSrsLevelProgress(vocabularyProgressInfo)
}

func printSrsLevelProgress(progressionData ProgressionData) {
	fmt.Printf("\tApprentice+:  %d\n", progressionData.NumApprentice)
	fmt.Printf("\tGuru+:        %d\n", progressionData.NumGuru)
	fmt.Printf("\tMaster+:      %d\n", progressionData.NumMaster)
	fmt.Printf("\tEnlightened+: %d\n", progressionData.NumEnlightened)
	fmt.Printf("\tBurned:       %d\n", progressionData.NumBurned)
}

func ComputeProgressInfo(client wanikani.Client, subjectType string) ProgressionData {
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

func fetchDistribution(client wanikani.Client, subjectType string) []int {
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
