package wanikani

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/data"
)

type Command int

const (
	LevelProgress Command = iota
	UserInfo
)

type CommandInfo struct {
	Command Command
}

func PrintLevelProgress() {
	client := CreateClient()
	userInformation := client.FetchUserInformation()

	fmt.Printf("Fetching information for user '%s' at level %v\n", userInformation.Username, userInformation.Level)

	progressions := FetchProgressions(client, userInformation.Level.String())

	spacedRepetitionSystems := client.FetchSpacedRepetitionSystems()
	spacedRepetitionSystemMap := data.CreateSpacedRepetitionSystemMap(spacedRepetitionSystems)

	UpdateOptimalUnlockTimes(spacedRepetitionSystemMap, &progressions)

	PrintTable(progressions, progressions.RadicalProgression, progressions.KanjiProgression)

	earliestProgression := FindTimeOfPassingRatio(progressions)
	fmt.Printf("Earliest progression time: %v", earliestProgression)
}
