package data

import (
	"encoding/json"
)

type User struct {
	ID                       string       `json:"id"`
	Username                 string       `json:"username"`
	Level                    json.Number  `json:"level"`
	ProfileURL               string       `json:"profile_url"`
	StartedAt                string       `json:"started_at"`
	CurrentVacationStartedAt string       `json:"current_vacation_started_at"`
	Subscription             Subscription `json:"subscription"`
	Preferences              Preferences  `json:"preferences"`
}

type Preferences struct {
	DefaultVoiceActorID        json.Number `json:"default_voice_actor_id"`
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
	URL            string `json:"url"`
	DataUploadedAt string `json:"data_updated_at"`
	Data           User   `json:"data"`
}
