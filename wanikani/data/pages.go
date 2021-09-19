package data

import "encoding/json"

type Pages struct {
	PerPage     json.Number `json:"per_page"`
	NextURL     string      `json:"next_url"`
	PreviousURL string      `json:"previous_url"`
}
