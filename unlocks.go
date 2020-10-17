package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"wanikani_cli/data"
)

type Unlocks []time.Time

func (unlocks Unlocks) String() string {
	times := make([]string, len(unlocks))

	for idx, element := range unlocks {
		if (element == time.Time{}) {
			times[idx] = fmt.Sprintf("%s: P", getStageName(idx))
		} else {
			res := element.Format("02-01-2006 15:04")
			times[idx] = fmt.Sprintf("%s: %s", getStageName(idx), res)
		}
	}
	return strings.Join(times, " ")
}

func computeOptimalUnlocks(system data.SpacedRepetitionSystem, progression Progression) Unlocks {
	optimalUnlocks := make([]time.Time, len(system.Stages))
	for idx, stage := range system.Stages {
		if int64(idx) < progression.SrsStage+1 {
			optimalUnlocks[idx] = time.Time{}
		} else if int64(idx) == progression.SrsStage+1 {
			optimalUnlocks[idx] = progression.AvailableAt
		} else if int64(idx) > progression.SrsStage+1 {
			lastUnlock := optimalUnlocks[idx-1]
			intervalDuration := time.Duration(toIntOrPanic(stage.Interval))
			nextUnlock := lastUnlock.Add(intervalDuration * time.Second)
			optimalUnlocks[idx] = nextUnlock
		}
	}
	return optimalUnlocks
}

func toIntOrPanic(value json.Number) int64 {
	intValue, err := value.Int64()
	if err != nil {
		panic(fmt.Errorf("could not convert '%v' to int: %v", value, err))
	}
	return intValue
}
