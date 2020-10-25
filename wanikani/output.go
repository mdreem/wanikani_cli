package wanikani

import (
	"fmt"
	"strings"
	"time"
)

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
			res := element.In(location).Format("02.01.2006 15:04")
			times[idx] = res
		}
	}
	return strings.Join(times, "|")
}

func PrintTable(progressions Progressions, radicalProgression []Progression, kanjiProgression []Progression) {
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
		progressionTime := progression.PassedAt.In(location).Format("02.01.2006 15:04")

		col := formatColumn(progression.UnlockTimes)
		formattedColumn := fmt.Sprintf("%3d|%s|%5t|%s|%s|\n", idx, progression.Characters, progression.PassedAt != time.Time{}, progressionTime, col)
		fmt.Print(formattedColumn)
	}

	fmt.Print("---------- Kanji --------------\n")

	for idx, progression := range kanjiProgression {
		progressionTime := progression.PassedAt.In(location).Format("02.01.2006 15:04")

		col := formatColumn(progression.UnlockTimes)
		formattedColumn := fmt.Sprintf("%3d|%s|%5t|%s|%s|\n", idx, progression.Characters, progression.PassedAt != time.Time{}, progressionTime, col)
		fmt.Print(formattedColumn)
	}
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
