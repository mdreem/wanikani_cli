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
