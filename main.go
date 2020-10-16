package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
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
		fmt.Printf("\t%v", optimalUnlocks)
	}
}

type Unlocks []time.Time

func (unlocks Unlocks) String() string {
	times := make([]string, len(unlocks))

	for idx, element := range unlocks {
		if (element == time.Time{}) {
			times[idx] = fmt.Sprintf("%d: P", idx)
		} else {
			res := element.Format("02-01-2006 15:04")
			times[idx] = fmt.Sprintf("%d: %s", idx, res)
		}
	}
	return strings.Join(times, " ")
}

func computeOptimalUnlocks(system data.SpacedRepetitionSystem, progression Progression) Unlocks {
	optimalUnlocks := make([]time.Time, len(system.Stages))
	for idx, stage := range system.Stages {
		if int64(idx) < progression.SrsStage {
			optimalUnlocks[idx] = time.Time{}
		}
		if int64(idx) == progression.SrsStage {
			optimalUnlocks[idx] = progression.AvailableAt
		}
		if int64(idx) > progression.SrsStage {
			lastUnlock := optimalUnlocks[idx-1]
			intervalDuration := time.Duration(toIntOrPanic(stage.Interval))
			nextUnlock := lastUnlock.Add(intervalDuration * time.Second)
			optimalUnlocks[idx] = nextUnlock
		}
	}
	return optimalUnlocks
}

func toIntOrPanic(value json.Number) int64 {
	intValue, err := value.Int64()
	if err != nil {
		panic(fmt.Errorf("could not convert '%v' to int: %v", value, err))
	}
	return intValue
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
