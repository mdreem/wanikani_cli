package data

import (
	"encoding/json"
	"fmt"
	"sort"
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

func (o WanikaniClient) FetchSpacedRepetitionSystems() []SpacedRepetitionSystemEnvelope {
	spacedRepetitionSystemsEnvelope := SpacedRepetitionSystemsEnvelope{}

	err := o.FetchWanikaniDataFromEndpoint("spaced_repetition_systems", &spacedRepetitionSystemsEnvelope, nil)
	if err != nil {
		panic(fmt.Errorf("error fetching list of spaced repetition systems: %v", err))
	}

	return spacedRepetitionSystemsEnvelope.Data
}

func toIntOrPanic(value json.Number) int64 {
	if value == "" {
		return -1
	}
	intValue, err := value.Int64()
	if err != nil {
		panic(fmt.Errorf("could not convert '%v' to int: %v", value, err))
	}
	return intValue
}

func CreateSpacedRepetitionSystemMap(spacedRepetitionSystemList []SpacedRepetitionSystemEnvelope) map[string]SpacedRepetitionSystem {
	spacedRepetitionSystems := make(map[string]SpacedRepetitionSystem, len(spacedRepetitionSystemList))

	for _, spacedRepetitionSystem := range spacedRepetitionSystemList {
		spacedRepetitionSystems[spacedRepetitionSystem.ID.String()] = spacedRepetitionSystem.Data
		sort.Slice(spacedRepetitionSystem.Data.Stages, func(i, j int) bool {
			return toIntOrPanic(spacedRepetitionSystem.Data.Stages[i].Interval) < toIntOrPanic(spacedRepetitionSystem.Data.Stages[j].Interval)
		})
	}

	return spacedRepetitionSystems
}
