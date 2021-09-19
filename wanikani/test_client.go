//go:build mock
// +build mock

package wanikani

import (
	"github.com/mdreem/wanikani_cli/wanikani/data"
	"net/http"
)

type TestClient struct {
	Assignments []data.AssignmentEnvelope
}

func (testClient TestClient) FetchAssignments(levels []string, subjectTypes []string) []data.AssignmentEnvelope {
	return testClient.Assignments
}

func (testClient TestClient) FetchWanikaniDataFromEndpoint(endpoint string, data interface{}, parameters map[string]string) error {
	return nil
}

func (testClient TestClient) FetchWanikaniDataFromURL(url string, data interface{}) error {
	return nil
}

func (testClient TestClient) fetchWanikaniData(request *http.Request, data interface{}) error {
	return nil
}

func (testClient TestClient) createAuthorizedRequest(url string) (*http.Request, error) {
	return nil, nil
}

func (testClient TestClient) createRequest(endpoint string, parameters map[string]string) (*http.Request, error) {
	return nil, nil
}

func (testClient TestClient) convertResponse(response *http.Response, data interface{}) error {
	return nil
}

func (testClient TestClient) FetchSpacedRepetitionSystems() []data.SpacedRepetitionSystemEnvelope {
	return nil
}

func (testClient TestClient) FetchSubjects(ids []string, levels []string, types []string) []data.SubjectEnvelope {
	return nil
}

func (testClient TestClient) FetchUserInformation() data.User {
	return data.User{}
}

func CreateTestClient() {

}
