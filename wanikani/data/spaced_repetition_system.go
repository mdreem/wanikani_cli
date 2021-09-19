package data

import (
	"encoding/json"
)

type Stage struct {
	Interval     json.Number `json:"interval"`
	Position     json.Number `json:"position"`
	IntervalUnit string      `json:"interval_unit"`
}

type SpacedRepetitionSystem struct {
	CreatedAt   string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	UnlockingStagePosition json.Number `json:"unlocking_stage_position"`
	StartingStagePosition  json.Number `json:"starting_stage_position"`
	PassingStagePosition   json.Number `json:"passing_stage_position"`
	BurningStagePosition   json.Number `json:"burning_stage_position"`

	Stages []Stage `json:"stages"`
}

type SpacedRepetitionSystemsEnvelope struct {
	Object         string                           `json:"object"`
	URL            string                           `json:"url"`
	Pages          Pages                            `json:"pages"`
	TotalCount     json.Number                      `json:"total_count"`
	DataUploadedAt string                           `json:"data_updated_at"`
	Data           []SpacedRepetitionSystemEnvelope `json:"data"`
}

type SpacedRepetitionSystemEnvelope struct {
	ID             json.Number            `json:"id"`
	Object         string                 `json:"object"`
	URL            string                 `json:"url"`
	DataUploadedAt string                 `json:"data_updated_at"`
	Data           SpacedRepetitionSystem `json:"data"`
}
