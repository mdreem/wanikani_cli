package levelprogress

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/wanikani"
	"github.com/spf13/cobra"
	"time"
)

func PrintLevelProgress(_ *cobra.Command, _ []string) {
	client := wanikani.CreateClient()
	userInformation := client.FetchUserInformation()

	fmt.Printf("Fetching information for user '%s' at level %v\n", userInformation.Username, userInformation.Level)

	progressions := FetchProgressions(client, userInformation.Level.String())

	spacedRepetitionSystems := client.FetchSpacedRepetitionSystems()
	spacedRepetitionSystemMap := wanikani.CreateSpacedRepetitionSystemMap(spacedRepetitionSystems)

	UpdateOptimalUnlockTimes(spacedRepetitionSystemMap, &progressions)

	PrintTable(progressions, progressions.RadicalProgression, progressions.KanjiProgression)

	earliestProgression := FindTimeOfPassingRatio(progressions)
	location := time.Now().Location()
	fmt.Printf("\nEarliest progression time: %v", earliestProgression.In(location))
}
