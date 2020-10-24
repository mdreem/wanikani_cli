package wanikani

import (
	"encoding/json"
	"fmt"
	"github.com/mdreem/wanikani_cli/data"
	"math"
	"sort"
	"strings"
	"time"
)

var (
	timeNow = time.Now
)

type Unlocks struct {
	UnlockTimes []time.Time
	Unlocked    bool
}

func (unlocks Unlocks) String() string {
	times := make([]string, len(unlocks.UnlockTimes))

	location := timeNow().Location()
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
		return fmt.Sprintf("unlocked [%s]", joinedTimes)
	} else {
		return fmt.Sprintf("not unlocked [%s]", joinedTimes)
	}
}

func computeOptimalUnlocks(system data.SpacedRepetitionSystem, progression Progression) Unlocks {
	optimalUnlocks := make([]time.Time, len(system.Stages))

	if !progression.isUnlocked() && !progression.UnlockByRadicalComputed {
		return Unlocks{
			UnlockTimes: optimalUnlocks,
			Unlocked:    false,
		}
	}

	for idx, stage := range system.Stages {
		if int64(idx) < progression.SrsStage+1 {
			optimalUnlocks[idx] = time.Time{}
		} else if int64(idx) == progression.SrsStage+1 {
			now := timeNow()

			referenceTime := getReferenceTimeForAvailability(progression)

			if referenceTime.Before(now) {
				optimalUnlocks[idx] = now.Truncate(time.Hour).UTC()
			} else {
				optimalUnlocks[idx] = referenceTime
			}
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

	if (progression.AvailableAt == time.Time{}) {
		return Unlocks{
			UnlockTimes: optimalUnlocks,
			Unlocked:    false,
		}
	} else {
		return Unlocks{
			UnlockTimes: optimalUnlocks,
			Unlocked:    true,
		}
	}
}

func getReferenceTimeForAvailability(progression Progression) time.Time {
	if (progression.AvailableAt != time.Time{}) {
		return progression.AvailableAt
	} else if (progression.PotentiallyAvailableAt != time.Time{}) {
		return progression.PotentiallyAvailableAt
	} else {
		return time.Time{}
	}
}

func UpdateOptimalUnlockTimes(spacedRepetitionSystems map[string]data.SpacedRepetitionSystem, progressions *Progressions) {
	updateUnlockTimes(spacedRepetitionSystems, &(progressions.RadicalProgression))
	updateLockedKanji(&(progressions.RadicalProgression), &(progressions.KanjiProgression))
	updateUnlockTimes(spacedRepetitionSystems, &(progressions.KanjiProgression))
}

func updateLockedKanji(radicalProgressions *[]Progression, kanjiProgressions *[]Progression) {
	kanjiProgressionMap := make(map[string]int)
	for idx, kanjiProgression := range *kanjiProgressions {
		kanjiProgressionMap[kanjiProgression.SubjectId] = idx
	}

	for _, radicalProgression := range *radicalProgressions {
		for _, containingKanji := range radicalProgression.AmalgamationSubjectIds {
			kanjiIdx, ok := kanjiProgressionMap[containingKanji.String()]
			if ok {
				kanji := (*kanjiProgressions)[kanjiIdx]
				if !kanji.UnlockTimes.Unlocked {
					if kanji.UnlockTimes.UnlockTimes == nil {
						kanji.UnlockTimes.UnlockTimes = make([]time.Time, len(radicalProgression.UnlockTimes.UnlockTimes))
					}
					if (radicalProgression.PassedAt == time.Time{}) {
						kanji.PotentiallyAvailableAt = radicalProgression.UnlockTimes.UnlockTimes[5]
					} else {
						kanji.PotentiallyAvailableAt = radicalProgression.PassedAt
					}
					kanji.UnlockByRadicalComputed = true
					(*kanjiProgressions)[kanjiIdx] = kanji
				}
			}
		}
	}
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
		var firstUnlockTime, secondUnlockTime time.Time
		firstProgressionTime := kanjiProgression[i].UnlockTimes.UnlockTimes[5]
		secondProgressionTime := kanjiProgression[j].UnlockTimes.UnlockTimes[5]

		if (kanjiProgression[i].PassedAt != time.Time{}) {
			firstUnlockTime = kanjiProgression[i].PassedAt
		} else {
			firstUnlockTime = firstProgressionTime
		}

		if (kanjiProgression[j].PassedAt != time.Time{}) {
			secondUnlockTime = kanjiProgression[j].PassedAt
		} else {
			secondUnlockTime = secondProgressionTime
		}

		return firstUnlockTime.Before(secondUnlockTime)
	})

	ninetyPercentPoint := int(math.Ceil(0.9*float64(len(kanjiProgression)))) - 1
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
