package main

import (
	"encoding/json"
	"fmt"
)

type Assignment struct {
	CreatedAt     string      `json:"created_at"`
	SubjectId     json.Number `json:"subject_id"`
	SubjectType   string      `json:"subject_type"`
	SrsStage      json.Number `json:"srs_stage"`
	UnlockedAt    string      `json:"unlocked_at"`
	StartedAt     string      `json:"started_at"`
	PassedAt      string      `json:"passed_at"`
	BurnedAt      string      `json:"burned_at"`
	AvailableAt   string      `json:"available_at"`
	ResurrectedAt string      `json:"resurrected_at"`
}

type AssignmentsEnvelope struct {
	Object         string               `json:"object"`
	Url            string               `json:"url"`
	Pages          Pages                `json:"pages"`
	TotalCount     json.Number          `json:"total_count"`
	DataUploadedAt string               `json:"data_updated_at"`
	Data           []AssignmentEnvelope `json:"data"`
}

type AssignmentEnvelope struct {
	Id             json.Number `json:"id"`
	Object         string      `json:"object"`
	Url            string      `json:"url"`
	DataUploadedAt string      `json:"data_updated_at"`
	Data           Assignment  `json:"data"`
}

func (o Client) FetchAssignments() []AssignmentEnvelope {
	assignmentsEnvelope := AssignmentsEnvelope{}

	parameters := make(map[string]string)
	parameters["levels"] = joinArrayToParameter([]string{"28"})

	err := o.FetchWanikaniData("assignments", &assignmentsEnvelope, parameters)
	if err != nil {
		panic(fmt.Errorf("error fetching list of assignments: %v", err))
	}

	return assignmentsEnvelope.Data
}
