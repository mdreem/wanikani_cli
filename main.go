package main

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"wanikani_cli/data"
	"wanikani_cli/wanikani"
)

func main() {
	initializeConfiguration()

	client := CreateClient()
	userInformation := client.FetchUserInformation()

	fmt.Printf("userInformation: %v\n", userInformation)
	fmt.Printf("level: %v\n", userInformation.Level)

	progressions := wanikani.FetchProgressions(client, "29")

	spacedRepetitionSystems := client.FetchSpacedRepetitionSystems()
	spacedRepetitionSystemMap := data.CreateSpacedRepetitionSystemMap(spacedRepetitionSystems)

	wanikani.UpdateOptimalUnlockTimes(spacedRepetitionSystemMap, &progressions)

	fmt.Print("Radicals\n")
	for idx, progression := range progressions.RadicalProgression {
		fmt.Printf("%d: %v\n", idx, progression)
	}

	fmt.Print("Kanji\n")
	for idx, progression := range progressions.KanjiProgression {
		optimalUnlocks := progression.UnlockTimes

		fmt.Printf("%d: %v\n", idx, progression)
		fmt.Printf("\t%v\n", optimalUnlocks.UnlockTimes)
	}

	earliestProgression := wanikani.FindTimeOfPassingRatio(progressions)
	fmt.Printf("Earliest progression time: %v", earliestProgression)
}

func initializeConfiguration() {
	viper.SetConfigName("wanikani")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error: %s \n", err))
	}
}

func GetApiKey() string {
	apiConfig := viper.GetStringMapString("api")
	if apiConfig == nil {
		panic(fmt.Errorf("cannot find section 'api' in config file"))
	}
	apiKey := apiConfig["api_key"]
	return apiKey
}

func CreateClient() data.Client {
	apiKey := GetApiKey()
	return data.Client{BaseUrl: "https://api.wanikani.com/v2/", ApiKey: apiKey, Client: &http.Client{}}
}
