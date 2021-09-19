package wanikani

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/wanikani/data"
)

func (o RealClient) FetchSubjects(ids []string, levels []string, types []string) []data.SubjectEnvelope {
	subjectEnvelope := data.SubjectsEnvelope{}

	parameters := make(map[string]string)

	if ids != nil {
		parameters["ids"] = joinArrayToParameter(ids)
	}
	if levels != nil {
		parameters["levels"] = joinArrayToParameter(levels)
	}
	if types != nil {
		parameters["types"] = joinArrayToParameter(types)
	}

	err := o.FetchWanikaniDataFromEndpoint("subjects", &subjectEnvelope, parameters)
	if err != nil {
		panic(fmt.Errorf("error fetching list of subjects: %v", err))
	}

	return subjectEnvelope.Data
}
