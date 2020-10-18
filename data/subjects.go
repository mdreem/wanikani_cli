package data

import (
	"encoding/json"
	"fmt"
)

type Meaning struct {
	Meaning        string `json:"meaning"`
	Primary        bool   `json:"primary"`
	AcceptedAnswer bool   `json:"accepted_answer"`
}

type Reading struct {
	Type           string `json:"type"`
	Primary        bool   `json:"primary"`
	AcceptedAnswer bool   `json:"accepted_answer"`
	Reading        string `json:"reading"`
}

type Subject struct {
	CreatedAt   string      `json:"created_at"`
	Level       json.Number `json:"level"`
	Slug        string      `json:"slug"`
	HiddenAt    string      `json:"hidden_at"`
	DocumentUrl string      `json:"document_url"`
	Characters  string      `json:"characters"`

	Meanings []Meaning `json:"meanings"`
	Readings []Reading `json:"readings"`

	ComponentSubjectIds       []json.Number `json:"component_subject_ids"`
	AmalgamationSubjectIds    []json.Number `json:"amalgamation_subject_ids"`
	VisuallySimilarSubjectIds []json.Number `json:"visually_similar_subject_ids"`

	MeaningMnemonic string `json:"meaning_mnemonic"`
	MeaningHint     string `json:"meaning_hint"`
	ReadingMnemonic string `json:"reading_mnemonic"`
	ReadingHint     string `json:"reading_hint"`

	LessonPosition           json.Number `json:"lesson_position"`
	SpacedRepetitionSystemId json.Number `json:"spaced_repetition_system_id"`
}

type SubjectsEnvelope struct {
	Object         string            `json:"object"`
	Url            string            `json:"url"`
	Pages          Pages             `json:"pages"`
	TotalCount     json.Number       `json:"total_count"`
	DataUploadedAt string            `json:"data_updated_at"`
	Data           []SubjectEnvelope `json:"data"`
}

type SubjectEnvelope struct {
	Id             json.Number `json:"id"`
	Object         string      `json:"object"`
	Url            string      `json:"url"`
	DataUploadedAt string      `json:"data_updated_at"`
	Data           Subject     `json:"data"`
}

func (o Client) FetchSubjects(ids []string) []SubjectEnvelope {
	subjectEnvelope := SubjectsEnvelope{}

	parameters := make(map[string]string)
	parameters["ids"] = joinArrayToParameter(ids)

	err := o.FetchWanikaniData("subjects", &subjectEnvelope, parameters)
	if err != nil {
		panic(fmt.Errorf("error fetching list of subjects: %v", err))
	}

	return subjectEnvelope.Data
}