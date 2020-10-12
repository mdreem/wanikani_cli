package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	initializeConfiguration()

	apiKey := GetApiKey()
	fmt.Printf("key: %s", apiKey)
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
