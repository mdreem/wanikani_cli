package wanikani

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/wanikani/data"
)

func (o RealClient) FetchUserInformation() data.User {
	userEnvelope := data.UserEnvelope{}
	err := o.FetchWanikaniDataFromEndpoint("user", &userEnvelope, nil)
	if err != nil {
		panic(fmt.Errorf("error fetching user data: %v", err))
	}
	return userEnvelope.Data
}
