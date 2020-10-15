package main

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

type Client struct {
	baseUrl string
	apiKey  string
	client  *http.Client
}

func main() {
	initializeConfiguration()

	client := CreateClient()
	userInformation := client.FetchUserInformation()

	fmt.Printf("userInformation: %v\n", userInformation)
	fmt.Printf("level: %v\n", userInformation.Level)

	radicalAssignments := client.FetchAssignments([]string{"28"}, []string{"radical"})
	kanjiAssignments := client.FetchAssignments([]string{"28"}, []string{"kanji"})

	radicalSubjectIds := make([]string, len(radicalAssignments))
	for idx, assignment := range radicalAssignments {
		radicalSubjectIds[idx] = assignment.Data.SubjectId.String()
	}

	kanjiSubjectIds := make([]string, len(kanjiAssignments))
	for idx, assignment := range kanjiAssignments {
		kanjiSubjectIds[idx] = assignment.Data.SubjectId.String()
	}

	radicalSubjects := client.FetchSubjects(radicalSubjectIds)
	kanjiSubjects := client.FetchSubjects(kanjiSubjectIds)

	for idx, subject := range radicalSubjects {
		fmt.Printf("%d: (%v) %s \n", idx, subject.Data.Level, subject.Data.Characters)
	}

	for idx, subject := range kanjiSubjects {
		fmt.Printf("%d: (%v) %s \n", idx, subject.Data.Level, subject.Data.Characters)
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

func CreateClient() Client {
	apiKey := GetApiKey()
	return Client{baseUrl: "https://api.wanikani.com/v2/", apiKey: apiKey, client: &http.Client{}}
}
