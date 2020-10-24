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

func formatColumn(unlocks Unlocks) string {
	times := make([]string, len(unlocks.UnlockTimes))

	location := timeNow().Location()
	for idx, element := range unlocks.UnlockTimes {
		if idx == 0 {
			continue
		}
		if (element == time.Time{}) {
			times[idx] = "N/A             "
		} else {
			res := element.In(location).Format("02-01-2006 15:04")
			times[idx] = res
		}
	}
	return strings.Join(times, "|")
}

func printTable(progressions Progressions, radicalProgression []Progression, kanjiProgression []Progression) {
	headings := make([]string, len(progressions.KanjiProgression[0].UnlockTimes.UnlockTimes))
	for idx := range progressions.KanjiProgression[0].UnlockTimes.UnlockTimes {
		if idx == 0 {
			continue
		}
		headings[idx] = fmt.Sprintf("%16s", getStageName(idx))
	}
	formattedHeader := fmt.Sprintf("   | |      |Passed time     |%s|\n", strings.Join(headings, "|"))
	fmt.Print(formattedHeader)

	location := timeNow().Location()

	fmt.Print("---------- Radicals --------------\n")

	for idx, progression := range radicalProgression {
		progressionTime := progression.PassedAt.In(location).Format("02-01-2006 15:04")

		col := formatColumn(progression.UnlockTimes)
		formattedColumn := fmt.Sprintf("%3d|%s|%5t|%s|%s|\n", idx, progression.Characters, progression.PassedAt != time.Time{}, progressionTime, col)
		fmt.Print(formattedColumn)
	}

	fmt.Print("---------- Kanji --------------\n")

	for idx, progression := range kanjiProgression {
		progressionTime := progression.PassedAt.In(location).Format("02-01-2006 15:04")

		col := formatColumn(progression.UnlockTimes)
		formattedColumn := fmt.Sprintf("%3d|%s|%5t|%s|%s|\n", idx, progression.Characters, progression.PassedAt != time.Time{}, progressionTime, col)
		fmt.Print(formattedColumn)
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
			if progression.AvailableAt.Before(now) {
				optimalUnlocks[idx] = now.Truncate(time.Hour).UTC()
			} else {
				optimalUnlocks[idx] = progression.AvailableAt
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
						kanji.AvailableAt = radicalProgression.UnlockTimes.UnlockTimes[5]
					} else {
						kanji.AvailableAt = radicalProgression.PassedAt
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

	printTable(progressions, progressions.RadicalProgression, kanjiProgression)

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
