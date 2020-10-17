package main

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"wanikani_cli/data"
)

func main() {
	initializeConfiguration()

	client := CreateClient()
	userInformation := client.FetchUserInformation()

	fmt.Printf("userInformation: %v\n", userInformation)
	fmt.Printf("level: %v\n", userInformation.Level)

	radicalProgression := fetchRadicalProgression(client, "28")
	kanjiProgression := fetchKanjiProgression(client, "28")

	spacedRepetitionSystems := client.FetchSpacedRepetitionSystems()
	spacedRepetitionSystemMap := data.CreateSpacedRepetitionSystemMap(spacedRepetitionSystems)

	fmt.Print("Radicals\n")
	for idx, progression := range radicalProgression {
		fmt.Printf("%d: %v\n", idx, progression)
	}

	fmt.Print("Kanji\n")
	for idx, progression := range kanjiProgression {
		system := spacedRepetitionSystemMap[progression.SrsSystem]

		optimalUnlocks := computeOptimalUnlocks(system, progression)

		fmt.Printf("%d: %v\n", idx, progression)
		fmt.Printf("\t%v\n", optimalUnlocks)
	}
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

func getStageName(stage int) string {
	switch stage {
	case 0:
		return "Not started"
	case 1:
		return "Apprentice 1"
	case 2:
		return "Apprentice 2"
	case 3:
		return "Apprentice 3"
	case 4:
		return "Apprentice 4"
	case 5:
		return "Guru 1"
	case 6:
		return "Guru 2"
	case 7:
		return "Master"
	case 8:
		return "Enlightened"
	case 9:
		return "Burned"
	default:
		return "Unknown"
	}
}
