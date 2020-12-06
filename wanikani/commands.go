package wanikani

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/data"
	"time"
)

type Command int

const (
	LevelProgress Command = iota
	UserInfo
	ProgressInfo
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
	location := time.Now().Location()
	fmt.Printf("\nEarliest progression time: %v", earliestProgression.In(location))
}

func PrintUserInfo() {
	client := CreateClient()
	userInformation := client.FetchUserInformation()

	fmt.Printf("Username: %s\n", userInformation.Username)
	fmt.Printf("Level: %s\n", userInformation.Level)
	fmt.Printf("Profile URL: %s\n", userInformation.ProfileURL)

	fmt.Printf("CurrentVacationStartedAt: %s\n", userInformation.CurrentVacationStartedAt)
	fmt.Println("Subscription:")
	fmt.Printf("\tActive: %t\n", userInformation.Subscription.Active)
	fmt.Printf("\tMax level granted: %s\n", userInformation.Subscription.MaxLevelGranted.String())
	fmt.Printf("\tPeriod ends at: %s\n", userInformation.Subscription.PeriodEndsAt)
	fmt.Printf("\tSubscription type: %s\n", userInformation.Subscription.Type)
}

func PrintProgressInfo() {
	client := CreateClient()

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
