package main

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/data"
	"github.com/mdreem/wanikani_cli/wanikani"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	initializeConfiguration()

	client := CreateClient()
	userInformation := client.FetchUserInformation()

	fmt.Printf("Fetching information for user '%s' at level %v\n", userInformation.Username, userInformation.Level)

	progressions := wanikani.FetchProgressions(client, "29")

	spacedRepetitionSystems := client.FetchSpacedRepetitionSystems()
	spacedRepetitionSystemMap := data.CreateSpacedRepetitionSystemMap(spacedRepetitionSystems)

	wanikani.UpdateOptimalUnlockTimes(spacedRepetitionSystemMap, &progressions)

	wanikani.PrintTable(progressions, progressions.RadicalProgression, progressions.KanjiProgression)

	earliestProgression := wanikani.FindTimeOfPassingRatio(progressions)
	fmt.Printf("Earliest progression time: %v", earliestProgression)
}

func initializeConfiguration() {
	viper.SetConfigName("wanikani")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error: %s", err))
	}
}

func GetAPIKey() string {
	apiConfig := viper.GetStringMapString("api")
	if apiConfig == nil {
		panic(fmt.Errorf("cannot find section 'api' in config file"))
	}
	apiKey := apiConfig["api_key"]
	return apiKey
}

func CreateClient() data.Client {
	apiKey := GetAPIKey()
	return data.Client{BaseURL: "https://api.wanikani.com/v2/", APIKey: apiKey, Client: &http.Client{}}
}
