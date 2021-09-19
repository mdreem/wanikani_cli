package wanikani

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/wanikani/data"
)

func (o WanikaniClient) FetchAssignments(levels []string, subjectTypes []string) []data.AssignmentEnvelope {
	parameters := make(map[string]string)

	if levels != nil {
		parameters["levels"] = joinArrayToParameter(levels)
	}
	if subjectTypes != nil {
		parameters["subject_types"] = joinArrayToParameter(subjectTypes)
	}
	assignmentsEnvelope := data.AssignmentsEnvelope{}

	err := o.FetchWanikaniDataFromEndpoint("assignments", &assignmentsEnvelope, parameters)
	if err != nil {
		panic(fmt.Errorf("error fetching list of assignments: %v", err))
	}

	var assignmentEnvelopeDataList = assignmentsEnvelope.Data
	var nextURL = assignmentsEnvelope.Pages.NextURL
	for nextURL != "" {
		currentAssignmentsEnvelope := data.AssignmentsEnvelope{}

		err := o.FetchWanikaniDataFromURL(nextURL, &currentAssignmentsEnvelope)
		if err != nil {
			panic(fmt.Errorf("error fetching list of assignments: %v", err))
		}
		nextURL = currentAssignmentsEnvelope.Pages.NextURL
		assignmentEnvelopeDataList = append(assignmentEnvelopeDataList, currentAssignmentsEnvelope.Data...)
	}

	return assignmentEnvelopeDataList
}
