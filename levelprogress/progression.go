package levelprogress

import (
	"encoding/json"
	"fmt"
	"github.com/mdreem/wanikani_cli/wanikani"
	"time"
)

type Progression struct {
	SubjectID  string
	Characters string

	SrsStage  int64
	SrsSystem string

	UnlockedAt             time.Time
	PassedAt               time.Time
	AvailableAt            time.Time
	UnlockTimes            Unlocks
	PotentiallyAvailableAt time.Time

	UnlockByRadicalComputed bool

	AmalgamationSubjectIds []json.Number
}

type Progressions struct {
	RadicalProgression []Progression
	KanjiProgression   []Progression
}

func (progression Progression) isUnlocked() bool {
	return progression.UnlockedAt != time.Time{}
}

func FetchProgressions(client wanikani.RealClient, level string) Progressions {
	return Progressions{
		RadicalProgression: fetchProgression(client, level, "radical"),
		KanjiProgression:   fetchProgression(client, level, "kanji"),
	}
}

func fetchProgression(client wanikani.RealClient, level string, subjectType string) []Progression {
	assignments := client.FetchAssignments([]string{level}, []string{subjectType})
	subjects := getSubjects(client, level, subjectType)

	progressionList := make([]Progression, 0)
	for _, assignment := range assignments {
		relatedSubjectID := assignment.Data.SubjectID.String()
		relatedSubject := subjects[relatedSubjectID]

		passedAt := parseTime(assignment.Data.PassedAt)
		unlockedAt := parseTime(assignment.Data.UnlockedAt)
		availableAt := parseTime(assignment.Data.AvailableAt)

		srsStage, err := assignment.Data.SrsStage.Int64()
		if err != nil {
			panic(fmt.Errorf("could not parse value %v: %v", assignment.Data.SrsStage, err))
		}

		progression := Progression{
			SubjectID:               relatedSubjectID,
			Characters:              relatedSubject.Characters,
			SrsStage:                srsStage,
			SrsSystem:               relatedSubject.SrsSystem,
			UnlockedAt:              unlockedAt,
			PassedAt:                passedAt,
			AvailableAt:             availableAt,
			UnlockByRadicalComputed: subjectType == "radical",
			AmalgamationSubjectIds:  relatedSubject.AmalgamationSubjectIds,
		}
		progressionList = append(progressionList, progression)

		relatedSubject.HasAssigment = true
		subjects[relatedSubjectID] = relatedSubject
	}

	for subjectID, subject := range subjects {
		if !subject.HasAssigment {
			progression := Progression{
				SubjectID:               subjectID,
				Characters:              subject.Characters,
				SrsStage:                0,
				SrsSystem:               subject.SrsSystem,
				UnlockedAt:              time.Time{},
				PassedAt:                time.Time{},
				AvailableAt:             time.Time{},
				UnlockByRadicalComputed: false,
				AmalgamationSubjectIds:  subject.AmalgamationSubjectIds,
			}
			progressionList = append(progressionList, progression)
		}
	}

	return progressionList
}

type subjectForAssigment struct {
	Characters             string
	SrsSystem              string
	AmalgamationSubjectIds []json.Number
	HasAssigment           bool
}

func getSubjects(client wanikani.RealClient, level string, subjectType string) map[string]subjectForAssigment {
	subjectList := client.FetchSubjects(nil, []string{level}, []string{subjectType})

	subjects := make(map[string]subjectForAssigment)

	for _, subject := range subjectList {
		subjects[subject.ID.String()] = subjectForAssigment{
			Characters:             subject.Data.Characters,
			SrsSystem:              subject.Data.SpacedRepetitionSystemID.String(),
			AmalgamationSubjectIds: subject.Data.AmalgamationSubjectIds,
			HasAssigment:           false,
		}
	}
	return subjects
}

func parseTime(timeString string) time.Time {
	if timeString == "" {
		return time.Time{}
	}
	passedAt, err := time.Parse(time.RFC3339Nano, timeString)
	if err != nil {
		panic(fmt.Errorf("could not parse date %s: %v", timeString, err))
	}
	return passedAt
}

func (progression Progression) String() string {
	location := time.Now().Location()

	unlockedAt := ""
	if (progression.UnlockedAt == time.Time{}) {
		unlockedAt = "Not unlocked yet"
	} else {
		unlockedAt = progression.UnlockedAt.In(location).String()
	}

	passedAt := ""
	if (progression.PassedAt == time.Time{}) {
		passedAt = "Not passed yet"
	} else {
		passedAt = progression.PassedAt.In(location).String()
	}

	availableAt := ""
	if (progression.AvailableAt == time.Time{}) {
		availableAt = "Not available yet"
	} else {
		availableAt = progression.AvailableAt.In(location).String()
	}

	return fmt.Sprintf("%s: Unlocked: %s - Passed: %s - Available: %s [%d]", progression.Characters, unlockedAt, passedAt, availableAt, progression.SrsStage)
}
