package wanikani

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (o RealClient) FetchWanikaniDataFromEndpoint(endpoint string, data interface{}, parameters map[string]string) error {
	request, err := o.createRequest(endpoint, parameters)
	if err != nil {
		fmt.Printf("an error occurred when creating the request: %v\n", err)
		return err
	}
	err = o.fetchWanikaniData(request, data)
	if err != nil {
		fmt.Printf("an error occurred when fetching data: %v\n", err)
		return err
	}

	return nil
}

func (o RealClient) FetchWanikaniDataFromURL(url string, data interface{}) error {
	request, err := o.createAuthorizedRequest(url)
	if err != nil {
		fmt.Printf("an error occurred when creating the request: %v\n", err)
		return err
	}
	err = o.fetchWanikaniData(request, data)
	if err != nil {
		fmt.Printf("an error occurred when fetching data: %v\n", err)
		return err
	}

	return nil
}

func (o RealClient) fetchWanikaniData(request *http.Request, data interface{}) error {
	response, err := o.Client.Do(request)
	if err != nil {
		fmt.Printf("an error occurred when executing the request: %v\n", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Did receive status code %d\n", response.StatusCode)
		return errors.New("wrong response code received")
	}

	err = o.convertResponse(response, data)
	if err != nil {
		fmt.Printf("an error occurred while converting the response: %v\n", err)
		return err
	}
	return nil
}

func (o RealClient) createAuthorizedRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+o.APIKey)
	return request, nil
}

func (o RealClient) createRequest(endpoint string, parameters map[string]string) (*http.Request, error) {
	request, err := o.createAuthorizedRequest(o.BaseURL + endpoint)
	if err != nil {
		return nil, err
	}

	if parameters != nil {
		q := request.URL.Query()
		for key, value := range parameters {
			q.Add(key, value)
		}
		request.URL.RawQuery = q.Encode()
	}

	return request, nil
}

func (o RealClient) convertResponse(response *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}
	return nil
}

func joinArrayToParameter(values []string) string {
	return strings.Join(values, ",")
}
