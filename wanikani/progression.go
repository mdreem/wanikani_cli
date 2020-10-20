package wanikani

import (
	"encoding/json"
	"fmt"
	"time"
	"wanikani_cli/data"
)

type Progression struct {
	SubjectId  string
	Characters string

	SrsStage  int64
	SrsSystem string

	UnlockedAt  time.Time
	PassedAt    time.Time
	AvailableAt time.Time
	UnlockTimes Unlocks

	AmalgamationSubjectIds []json.Number
}

type Progressions struct {
	RadicalProgression []Progression
	KanjiProgression   []Progression
}

func FetchProgressions(client data.Client, level string) Progressions {
	return Progressions{
		RadicalProgression: fetchProgression(client, level, "radical"),
		KanjiProgression:   fetchProgression(client, level, "kanji"),
	}
}

func fetchProgression(client data.Client, level string, subjectType string) []Progression {
	assignments := client.FetchAssignments([]string{level}, []string{subjectType})
	subjects := getSubjects(client, assignments)

	progressionList := make([]Progression, 0)
	for _, assignment := range assignments {
		relatedSubjectId := assignment.Data.SubjectId.String()
		relatedSubject := subjects[relatedSubjectId]

		passedAt := parseTime(assignment.Data.PassedAt)
		unlockedAt := parseTime(assignment.Data.UnlockedAt)
		availableAt := parseTime(assignment.Data.AvailableAt)

		srsStage, err := assignment.Data.SrsStage.Int64()
		if err != nil {
			panic(fmt.Errorf("could not parse value %v: %v", assignment.Data.SrsStage, err))
		}

		progression := Progression{
			SubjectId:              relatedSubjectId,
			Characters:             relatedSubject.Characters,
			SrsStage:               srsStage,
			SrsSystem:              relatedSubject.SpacedRepetitionSystemId.String(),
			UnlockedAt:             unlockedAt,
			PassedAt:               passedAt,
			AvailableAt:            availableAt,
			AmalgamationSubjectIds: relatedSubject.AmalgamationSubjectIds,
		}
		progressionList = append(progressionList, progression)
	}

	return progressionList
}

func getSubjects(client data.Client, assignments []data.AssignmentEnvelope) map[string]data.Subject {
	subjectIds := make([]string, len(assignments))

	for idx, assignment := range assignments {
		subjectIds[idx] = assignment.Data.SubjectId.String()
	}

	subjectList := client.FetchSubjects(subjectIds)

	subjects := make(map[string]data.Subject)
	for _, subject := range subjectList {
		subjects[subject.Id.String()] = subject.Data
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

func (o Progression) String() string {
	location := time.Now().Location()

	unlockedAt := ""
	if (o.UnlockedAt == time.Time{}) {
		unlockedAt = "Not unlocked yet"
	} else {
		unlockedAt = o.UnlockedAt.In(location).String()
	}

	passedAt := ""
	if (o.PassedAt == time.Time{}) {
		passedAt = "Not passed yet"
	} else {
		passedAt = o.PassedAt.In(location).String()
	}

	availableAt := ""
	if (o.AvailableAt == time.Time{}) {
		availableAt = "Not available yet"
	} else {
		availableAt = o.AvailableAt.In(location).String()
	}

	return fmt.Sprintf("%s: Unlocked: %s - Passed: %s - Available: %s [%d]", o.Characters, unlockedAt, passedAt, availableAt, o.SrsStage)
}
