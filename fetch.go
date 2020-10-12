package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (o Client) FetchWanikaniData(endpoint string, data interface{}) error {
	request, err := o.createRequest(endpoint)
	if err != nil {
		fmt.Printf("an error occured when creating the request: %v", err)
		return err
	}

	response, err := o.client.Do(request)
	if err != nil {
		fmt.Printf("an error occured when executing the request: %v", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Did receive status code %d", response.StatusCode)
		return errors.New("wrong response code received")
	}

	err = o.convertResponse(response, data)
	if err != nil {
		fmt.Printf("an error occured while converting the response: %v", err)
		return err
	}

	fmt.Printf("data: %+v", data)

	return nil
}

func (o Client) createRequest(endpoint string) (*http.Request, error) {
	request, err := http.NewRequest("GET", o.baseUrl+endpoint, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+o.apiKey)
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
