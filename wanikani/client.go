package wanikani

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/data"
	"github.com/spf13/viper"
	"net/http"
)

func GetAPIKey() string {
	apiConfig := viper.GetStringMapString("api")
	if apiConfig == nil {
		panic(fmt.Errorf("cannot find section 'api' in config file"))
	}
	apiKey := apiConfig["api_key"]
	return apiKey
}

func CreateClient() data.WanikaniClient {
	apiKey := GetAPIKey()
	return data.WanikaniClient{BaseURL: "https://api.wanikani.com/v2/", APIKey: apiKey, Client: &http.Client{}}
}
