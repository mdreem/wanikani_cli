package main

import (
	"fmt"
	"time"
	"wanikani_cli/data"
)

type Progression struct {
	Characters  string
	UnlockedAt  time.Time
	PassedAt    time.Time
	SrsStage    int64
	SrsSystem   string
	AvailableAt time.Time
}

func fetchRadicalProgression(client data.Client, level string) []Progression {
	return fetchProgression(client, level, "radical")
}

func fetchKanjiProgression(client data.Client, level string) []Progression {
	return fetchProgression(client, level, "kanji")
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
		availanbleAt := parseTime(assignment.Data.AvailableAt)

		srsStage, err := assignment.Data.SrsStage.Int64()
		if err != nil {
			panic(fmt.Errorf("could not parse value %v: %v", assignment.Data.SrsStage, err))
		}

		progression := Progression{
			Characters:  relatedSubject.Characters,
			UnlockedAt:  unlockedAt,
			AvailableAt: availanbleAt,
			PassedAt:    passedAt,
			SrsStage:    srsStage,
			SrsSystem:   relatedSubject.SpacedRepetitionSystemId.String(),
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
	unlockedAt := ""
	if (o.UnlockedAt == time.Time{}) {
		unlockedAt = "Not unlocked yet"
	} else {
		unlockedAt = o.UnlockedAt.String()
	}

	passedAt := ""
	if (o.PassedAt == time.Time{}) {
		passedAt = "Not passed yet"
	} else {
		passedAt = o.PassedAt.String()
	}

	availableAt := 0
	if (o.AvailableAt == time.Time{}) {
		passedAt = "Not available yet"
	} else {
		passedAt = o.AvailableAt.String()
	}

	return fmt.Sprintf("%s: Unlocked: %s - Passed: %s - Available: %d [%d]", o.Characters, unlockedAt, passedAt, availableAt, o.SrsStage)
}
