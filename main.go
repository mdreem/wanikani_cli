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

	fmt.Print("Radicals\n")
	for idx, progression := range radicalProgression {
		fmt.Printf("%d: %v\n", idx, progression)
	}

	fmt.Print("Kanji\n")
	for idx, progression := range kanjiProgression {
		fmt.Printf("%d: %v\n", idx, progression)
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
