package wanikani

import (
	"encoding/json"
	"fmt"
	"github.com/mdreem/wanikani_cli/wanikani/data"
	"sort"
)

func (o WanikaniClient) FetchSpacedRepetitionSystems() []data.SpacedRepetitionSystemEnvelope {
	spacedRepetitionSystemsEnvelope := data.SpacedRepetitionSystemsEnvelope{}

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

func CreateSpacedRepetitionSystemMap(spacedRepetitionSystemList []data.SpacedRepetitionSystemEnvelope) map[string]data.SpacedRepetitionSystem {
	spacedRepetitionSystems := make(map[string]data.SpacedRepetitionSystem, len(spacedRepetitionSystemList))

	for _, spacedRepetitionSystem := range spacedRepetitionSystemList {
		spacedRepetitionSystems[spacedRepetitionSystem.ID.String()] = spacedRepetitionSystem.Data
		sort.Slice(spacedRepetitionSystem.Data.Stages, func(i, j int) bool {
			return toIntOrPanic(spacedRepetitionSystem.Data.Stages[i].Interval) < toIntOrPanic(spacedRepetitionSystem.Data.Stages[j].Interval)
		})
	}

	return spacedRepetitionSystems
}
