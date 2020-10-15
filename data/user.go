package data

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Id                       string       `json:"id"`
	Username                 string       `json:"username"`
	Level                    json.Number  `json:"level"`
	ProfileUrl               string       `json:"profile_url"`
	StartedAt                string       `json:"started_at"`
	CurrentVacationStartedAt string       `json:"current_vacation_started_at"`
	Subscription             Subscription `json:"subscription"`
	Preferences              Preferences  `json:"preferences"`
}

type Preferences struct {
	DefaultVoiceActorId        json.Number `json:"default_voice_actor_id"`
	LessonsAutoplayAudio       bool        `json:"lessons_autoplay_audio"`
	LessonsBatchSize           json.Number `json:"lessons_batch_size"`
	LessonsPresentationOrder   string      `json:"lessons_presentation_order"`
	ReviewsAutoplayAudio       bool        `json:"reviews_autoplay_audio"`
	ReviewsDisplaySrsIndicator bool        `json:"reviews_display_srs_indicator"`
}

type Subscription struct {
	Active          bool        `json:"active"`
	Type            string      `json:"type"`
	MaxLevelGranted json.Number `json:"max_level_granted"`
	PeriodEndsAt    string      `json:"period_ends_at"`
}

type UserEnvelope struct {
	Object         string `json:"object"`
	Url            string `json:"url"`
	DataUploadedAt string `json:"data_updated_at"`
	Data           User   `json:"data"`
}

func (o Client) FetchUserInformation() User {
	userEnvelope := UserEnvelope{}
	err := o.FetchWanikaniData("user", &userEnvelope, nil)
	if err != nil {
		panic(fmt.Errorf("error fetching user data: %v", err))
	}
	return userEnvelope.Data
}
