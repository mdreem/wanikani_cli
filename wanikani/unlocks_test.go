package wanikani

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
	"wanikani_cli/data"
)

func Test_computeOptimalUnlocks(t *testing.T) {
	srs := createSrs()

	type args struct {
		system        data.SpacedRepetitionSystem
		progression   Progression
		referenceTime time.Time
	}
	tests := []struct {
		name string
		args args
		want Unlocks
	}{
		{
			name: "level progression of unseen item",
			args: args{
				system: srs,
				progression: Progression{
					Characters:  "never seen",
					UnlockedAt:  time.Time{},
					PassedAt:    time.Time{},
					SrsStage:    0,
					SrsSystem:   "this",
					AvailableAt: time.Time{},
				},
				referenceTime: getTime("01-01-2020 00:00"),
			},
			want: Unlocks{
				UnlockTimes: make([]time.Time, 10),
				Unlocked:    false,
			},
		},
		{
			name: "item unlocked at stage 2",
			args: args{
				system: srs,
				progression: Progression{
					Characters:  "never seen",
					UnlockedAt:  getTime("01-01-2020 08:00"),
					PassedAt:    time.Time{},
					SrsStage:    2,
					SrsSystem:   "this",
					AvailableAt: getTime("01-01-2020 08:00"),
				},
				referenceTime: getTime("01-01-2020 00:00"),
			},
			want: Unlocks{
				UnlockTimes: createUnlockTimes(),
				Unlocked:    true,
			},
		},
		{
			name: "item unlocked at stage 2, but review starts later",
			args: args{
				system: srs,
				progression: Progression{
					Characters:  "never seen",
					UnlockedAt:  getTime("01-01-2020 08:00"),
					PassedAt:    time.Time{},
					SrsStage:    2,
					SrsSystem:   "this",
					AvailableAt: getTime("01-01-2020 08:00"),
				},
				referenceTime: getTime("02-02-2020 09:00"),
			},
			want: Unlocks{
				UnlockTimes: createUnlockTimesShifted(),
				Unlocked:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeNow = func() time.Time {
				return tt.args.referenceTime
			}
			if got := computeOptimalUnlocks(tt.args.system, tt.args.progression); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeOptimalUnlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateLockedKanji(t *testing.T) {
	type args struct {
		radicalProgressions []Progression
		kanjiProgressions   []Progression
	}
	tests := []struct {
		name string
		args args
		want []Unlocks
	}{
		{
			name: "Kanji is not yet unlocked, but radical is",
			args: args{
				radicalProgressions: []Progression{
					createRadicalProgression(getTime("09-01-2020 06:00"), true),
				},
				kanjiProgressions: []Progression{
					createKanjiProgression(false, make([]time.Time, 10)),
				},
			},
			want: []Unlocks{
				{
					UnlockTimes: createUnlockTimesWithUnlockSet(),
					Unlocked:    false,
				},
			},
		},
		{
			name: "Kanji is already unlocked, radical does not update kanji data",
			args: args{
				radicalProgressions: []Progression{
					createRadicalProgression(getTime("09-01-2020 06:00"), true),
				},
				kanjiProgressions: []Progression{
					createKanjiProgression(true, createUnlockTimes()),
				},
			},
			want: []Unlocks{
				{
					UnlockTimes: createUnlockTimes(),
					Unlocked:    true,
				},
			},
		},
		{
			name: "Kanji is not yet unlocked, radical is also not unlocked",
			args: args{
				radicalProgressions: []Progression{
					createRadicalProgression(time.Time{}, false),
				},
				kanjiProgressions: []Progression{
					createKanjiProgression(false, make([]time.Time, 10)),
				},
			},
			want: []Unlocks{
				{
					UnlockTimes: make([]time.Time, 10),
					Unlocked:    false,
				},
			},
		},
		{
			name: "Kanji is already unlocked, radical is not. Does not change kanji data",
			args: args{
				radicalProgressions: []Progression{
					createRadicalProgression(time.Time{}, false),
				},
				kanjiProgressions: []Progression{
					createKanjiProgression(true, createUnlockTimes()),
				},
			},
			want: []Unlocks{
				{
					UnlockTimes: createUnlockTimes(),
					Unlocked:    true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeNow = func() time.Time {
				return getTime("01-01-2020 00:00")
			}
			updateLockedKanji(&tt.args.radicalProgressions, &tt.args.kanjiProgressions)

			got := make([]Unlocks, len(tt.args.kanjiProgressions))
			for idx, progression := range tt.args.kanjiProgressions {
				got[idx] = progression.UnlockTimes
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateLockedKanji() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createSrs() data.SpacedRepetitionSystem {
	return data.SpacedRepetitionSystem{
		Stages: []data.Stage{
			{
				Interval:     "",
				Position:     "0",
				IntervalUnit: "",
			},
			{
				Interval:     "7200",
				Position:     "1",
				IntervalUnit: "seconds",
			},
			{
				Interval:     "14400",
				Position:     "2",
				IntervalUnit: "seconds",
			},
			{
				Interval:     "28800",
				Position:     "3",
				IntervalUnit: "seconds",
			},
			{
				Interval:     "82800",
				Position:     "4",
				IntervalUnit: "seconds",
			},
			{
				Interval:     "601200",
				Position:     "5",
				IntervalUnit: "seconds",
			},
			{
				Interval:     "1206000",
				Position:     "6",
				IntervalUnit: "seconds",
			},
			{
				Interval:     "2588400",
				Position:     "7",
				IntervalUnit: "seconds",
			},
			{
				Interval:     "10364400",
				Position:     "8",
				IntervalUnit: "seconds",
			},
			{
				Interval:     "",
				Position:     "9",
				IntervalUnit: "",
			},
		},
	}
}

func createUnlockTimes() []time.Time {
	unlockTimes := make([]time.Time, 10)

	unlockTimes[0] = time.Time{}
	unlockTimes[1] = time.Time{}
	unlockTimes[2] = time.Time{}
	unlockTimes[3] = getTime("01-01-2020 08:00")
	unlockTimes[4] = getTime("02-01-2020 07:00")
	unlockTimes[5] = getTime("09-01-2020 06:00")
	unlockTimes[6] = getTime("23-01-2020 05:00")
	unlockTimes[7] = getTime("22-02-2020 04:00")
	unlockTimes[8] = getTime("21-06-2020 03:00")
	unlockTimes[9] = time.Time{}

	return unlockTimes
}

func createUnlockTimesShifted() []time.Time {
	unlockTimes := make([]time.Time, 10)

	unlockTimes[0] = time.Time{}
	unlockTimes[1] = time.Time{}
	unlockTimes[2] = time.Time{}
	unlockTimes[3] = getTime("02-02-2020 09:00")
	unlockTimes[4] = getTime("03-02-2020 08:00")
	unlockTimes[5] = getTime("10-02-2020 07:00")
	unlockTimes[6] = getTime("24-02-2020 06:00")
	unlockTimes[7] = getTime("25-03-2020 05:00")
	unlockTimes[8] = getTime("23-07-2020 04:00")
	unlockTimes[9] = time.Time{}

	return unlockTimes
}

func createUnlockTimesWithUnlockSet() []time.Time {
	unlockTimes := make([]time.Time, 10)
	unlockTimes[0] = getTime("09-01-2020 06:00")
	return unlockTimes
}

func createRadicalProgression(passedAt time.Time, unlocked bool) Progression {
	return Progression{
		SubjectId:   "1",
		Characters:  "X",
		SrsStage:    0,
		SrsSystem:   "1",
		UnlockedAt:  time.Time{},
		PassedAt:    passedAt,
		AvailableAt: time.Time{},
		UnlockTimes: Unlocks{
			UnlockTimes: createUnlockTimes(),
			Unlocked:    unlocked,
		},
		AmalgamationSubjectIds: []json.Number{"2", "100"},
	}
}

func createKanjiProgression(unlocked bool, unlockTimes []time.Time) Progression {
	return Progression{
		SubjectId:   "2",
		Characters:  "Y",
		SrsStage:    0,
		SrsSystem:   "1",
		UnlockedAt:  time.Time{},
		PassedAt:    time.Time{},
		AvailableAt: time.Time{},
		UnlockTimes: Unlocks{
			UnlockTimes: unlockTimes,
			Unlocked:    unlocked,
		},
		AmalgamationSubjectIds: []json.Number{},
	}
}

func getTime(timeString string) time.Time {
	parsedTime, _ := time.Parse("02-01-2006 15:04", timeString)
	return parsedTime
}
