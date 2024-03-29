package data

import (
	"encoding/json"
)

type Assignment struct {
	CreatedAt     string      `json:"created_at"`
	SubjectID     json.Number `json:"subject_id"`
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
	URL            string               `json:"url"`
	Pages          Pages                `json:"pages"`
	TotalCount     json.Number          `json:"total_count"`
	DataUploadedAt string               `json:"data_updated_at"`
	Data           []AssignmentEnvelope `json:"data"`
}

type AssignmentEnvelope struct {
	ID             json.Number `json:"id"`
	Object         string      `json:"object"`
	URL            string      `json:"url"`
	DataUploadedAt string      `json:"data_updated_at"`
	Data           Assignment  `json:"data"`
}
