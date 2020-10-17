package wanikani

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
	"wanikani_cli/data"
)

type Unlocks struct {
	UnlockTimes []time.Time
	Unlocked    bool
}

func (unlocks Unlocks) String() string {
	times := make([]string, len(unlocks.UnlockTimes))

	location := time.Now().Location()
	for idx, element := range unlocks.UnlockTimes {
		if (element == time.Time{}) {
			times[idx] = fmt.Sprintf("%s: P", getStageName(idx))
		} else {
			res := element.In(location).Format("02-01-2006 15:04")
			times[idx] = fmt.Sprintf("%s: %s", getStageName(idx), res)
		}
	}
	joinedTimes := strings.Join(times, ", ")

	if unlocks.Unlocked {
		return joinedTimes
	} else {
		return "not unlocked"
	}
}

func computeOptimalUnlocks(system data.SpacedRepetitionSystem, progression Progression) Unlocks {
	optimalUnlocks := make([]time.Time, len(system.Stages))

	if (progression.AvailableAt == time.Time{}) {
		return Unlocks{
			UnlockTimes: optimalUnlocks,
			Unlocked:    false,
		}
	}

	for idx, stage := range system.Stages {
		if int64(idx) < progression.SrsStage+1 {
			optimalUnlocks[idx] = time.Time{}
		} else if int64(idx) == progression.SrsStage+1 {
			optimalUnlocks[idx] = progression.AvailableAt
		} else if int64(idx) > progression.SrsStage+1 {
			lastUnlock := optimalUnlocks[idx-1]
			intervalDuration := time.Duration(toIntOrPanic(stage.Interval))
			var nextUnlock time.Time
			if intervalDuration < 0 {
				nextUnlock = time.Time{}
			} else {
				nextUnlock = lastUnlock.Add(intervalDuration * intervalUnitFactor(stage.IntervalUnit))
			}
			optimalUnlocks[idx] = nextUnlock
		}
	}
	return Unlocks{
		UnlockTimes: optimalUnlocks,
		Unlocked:    true,
	}
}

func UpdateOptimalUnlockTimes(spacedRepetitionSystems map[string]data.SpacedRepetitionSystem, progressions *Progressions) {
	updateUnlockTimes(spacedRepetitionSystems, &(progressions.RadicalProgression))
	updateUnlockTimes(spacedRepetitionSystems, &(progressions.KanjiProgression))
}

func updateUnlockTimes(spacedRepetitionSystems map[string]data.SpacedRepetitionSystem, progressions *[]Progression) {
	for idx, progression := range *progressions {
		system := spacedRepetitionSystems[progression.SrsSystem]

		optimalUnlocks := computeOptimalUnlocks(system, progression)
		(*progressions)[idx].UnlockTimes = optimalUnlocks
	}
}

func FindTimeOfPassingRatio(progressions Progressions) time.Time {
	kanjiProgression := progressions.KanjiProgression

	sort.Slice(kanjiProgression, func(i, j int) bool {
		firstUnlockTime := kanjiProgression[i].UnlockTimes.UnlockTimes[5]
		secondUnlockTime := kanjiProgression[j].UnlockTimes.UnlockTimes[5]
		return firstUnlockTime.Before(secondUnlockTime)
	})

	ninetyPercentPoint := int(math.Ceil(0.9 * float64(len(kanjiProgression))))
	return kanjiProgression[ninetyPercentPoint].UnlockTimes.UnlockTimes[5]
}

func intervalUnitFactor(intervalUnit string) time.Duration {
	switch intervalUnit {
	case "milliseconds":
		return time.Millisecond
	case "seconds":
		return time.Second
	case "minutes":
		return time.Minute
	case "hours":
		return time.Hour
	case "days":
		return 24 * time.Hour
	case "weeks":
		return 7 * 24 * time.Hour
	case "":
		return -1
	default:
		panic(fmt.Errorf("unknown interval unit %s", intervalUnit))
	}
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

func getStageName(stage int) string {
	switch stage {
	case 0:
		return "Not started"
	case 1:
		return "Apprentice 1"
	case 2:
		return "Apprentice 2"
	case 3:
		return "Apprentice 3"
	case 4:
		return "Apprentice 4"
	case 5:
		return "Guru 1"
	case 6:
		return "Guru 2"
	case 7:
		return "Master"
	case 8:
		return "Enlightened"
	case 9:
		return "Burned"
	default:
		return "Unknown"
	}
}
