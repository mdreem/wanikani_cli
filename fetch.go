package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (o Client) FetchWanikaniData(endpoint string, data interface{}, parameters map[string]string) error {
	request, err := o.createRequest(endpoint, parameters)
	if err != nil {
		fmt.Printf("an error occured when creating the request: %v\n", err)
		return err
	}

	response, err := o.client.Do(request)
	if err != nil {
		fmt.Printf("an error occured when executing the request: %v\n", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Did receive status code %d\n", response.StatusCode)
		return errors.New("wrong response code received")
	}

	err = o.convertResponse(response, data)
	if err != nil {
		fmt.Printf("an error occured while converting the response: %v\n", err)
		return err
	}

	return nil
}

func (o Client) createRequest(endpoint string, parameters map[string]string) (*http.Request, error) {
	request, err := http.NewRequest("GET", o.baseUrl+endpoint, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+o.apiKey)

	if parameters != nil {
		q := request.URL.Query()
		for key, value := range parameters {
			q.Add(key, value)
		}
		request.URL.RawQuery = q.Encode()
	}

	return request, nil
}

func (o Client) convertResponse(response *http.Response, data interface{}) error {
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
