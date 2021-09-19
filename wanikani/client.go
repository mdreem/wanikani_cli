package wanikani

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/wanikani/data"
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

func CreateClient() WanikaniClient {
	apiKey := GetAPIKey()
	return WanikaniClient{BaseURL: "https://api.wanikani.com/v2/", APIKey: apiKey, Client: &http.Client{}}
}

type WanikaniClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

type Client interface {
	FetchAssignments(levels []string, subjectTypes []string) []data.AssignmentEnvelope
	FetchWanikaniDataFromEndpoint(endpoint string, data interface{}, parameters map[string]string) error
	FetchWanikaniDataFromURL(url string, data interface{}) error
	fetchWanikaniData(request *http.Request, data interface{}) error
	createAuthorizedRequest(url string) (*http.Request, error)
	createRequest(endpoint string, parameters map[string]string) (*http.Request, error)
	convertResponse(response *http.Response, data interface{}) error
	FetchSpacedRepetitionSystems() []data.SpacedRepetitionSystemEnvelope
	FetchSubjects(ids []string, levels []string, types []string) []data.SubjectEnvelope
	FetchUserInformation() data.User
}
