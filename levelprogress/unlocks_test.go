package levelprogress

import (
	"encoding/json"
	"github.com/mdreem/wanikani_cli/wanikani/data"
	"reflect"
	"testing"
	"time"
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
				t.Errorf("computeOptimalUnlocks() = %v\n, want = %v", got, tt.want)
			}
		})
	}
}

func Test_updateLockedKanji(t *testing.T) {
	type unlockData struct {
		Unlocks          []Unlocks
		AvailableAtTimes []time.Time
	}
	type args struct {
		radicalProgressions []Progression
		kanjiProgressions   []Progression
	}
	tests := []struct {
		name string
		args args
		want unlockData
	}{
		{
			name: "Kanji is not yet unlocked, but radical is. Kanjis availability is updated.",
			args: args{
				radicalProgressions: []Progression{
					createRadicalProgression(getTime("09-01-2020 06:00"), time.Time{}, true),
				},
				kanjiProgressions: []Progression{
					createKanjiProgression(false, make([]time.Time, 10)),
				},
			},
			want: unlockData{
				Unlocks: []Unlocks{
					{
						UnlockTimes: make([]time.Time, 10),
						Unlocked:    false,
					},
				},
				AvailableAtTimes: []time.Time{getTime("09-01-2020 06:00")},
			},
		},
		{
			name: "Kanji is already unlocked, radical does not update kanji data",
			args: args{
				radicalProgressions: []Progression{
					createRadicalProgression(getTime("09-01-2020 06:00"), time.Time{}, true),
				},
				kanjiProgressions: []Progression{
					createKanjiProgression(true, createUnlockTimes()),
				},
			},
			want: unlockData{
				Unlocks: []Unlocks{
					{
						UnlockTimes: createUnlockTimes(),
						Unlocked:    true,
					},
				},
				AvailableAtTimes: make([]time.Time, 1),
			},
		},
		{
			name: "Kanji is not yet unlocked, radical is also not unlocked. Kanji availability is estimated.",
			args: args{
				radicalProgressions: []Progression{
					createRadicalProgression(time.Time{}, time.Time{}, false),
				},
				kanjiProgressions: []Progression{
					createKanjiProgression(false, make([]time.Time, 10)),
				},
			},
			want: unlockData{
				Unlocks: []Unlocks{
					{
						UnlockTimes: make([]time.Time, 10),
						Unlocked:    false,
					},
				},
				AvailableAtTimes: []time.Time{getTime("09-01-2020 06:00")},
			},
		},
		{
			name: "Kanji is already unlocked, radical is not. Does not change kanji data",
			args: args{
				radicalProgressions: []Progression{
					createRadicalProgression(time.Time{}, time.Time{}, false),
				},
				kanjiProgressions: []Progression{
					createKanjiProgression(true, createUnlockTimes()),
				},
			},
			want: unlockData{
				Unlocks: []Unlocks{
					{
						UnlockTimes: createUnlockTimes(),
						Unlocked:    true,
					},
				},
				AvailableAtTimes: make([]time.Time, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeNow = func() time.Time {
				return getTime("01-01-2020 00:00")
			}
			updateLockedKanji(&tt.args.radicalProgressions, &tt.args.kanjiProgressions)

			gotUnlocks := make([]Unlocks, len(tt.args.kanjiProgressions))
			for idx, progression := range tt.args.kanjiProgressions {
				gotUnlocks[idx] = progression.UnlockTimes
			}

			gotAvailableAtTimes := make([]time.Time, len(tt.args.kanjiProgressions))
			for idx, progression := range tt.args.kanjiProgressions {
				gotAvailableAtTimes[idx] = progression.PotentiallyAvailableAt
			}

			if !reflect.DeepEqual(gotUnlocks, tt.want.Unlocks) {
				t.Errorf("updateLockedKanji() = %v, want %v", gotUnlocks, tt.want.Unlocks)
			}

			if !reflect.DeepEqual(gotAvailableAtTimes, tt.want.AvailableAtTimes) {
				t.Errorf("updateLockedKanji() = %v, want %v", gotAvailableAtTimes, tt.want.AvailableAtTimes)
			}
		})
	}
}

func createSrsMap() map[string]data.SpacedRepetitionSystem {
	srsMap := make(map[string]data.SpacedRepetitionSystem)
	srsMap["1"] = createSrs()
	return srsMap
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

func createRadicalProgression(passedAt time.Time, potentiallyAvailableAt time.Time, unlocked bool) Progression {
	var srsStage int64
	if unlocked {
		srsStage = 2
	} else {
		srsStage = 0
	}
	return Progression{
		SubjectID:  "1",
		Characters: "X",
		SrsStage:   srsStage,
		SrsSystem:  "1",

		UnlockedAt:  time.Time{},
		PassedAt:    passedAt,
		AvailableAt: time.Time{},
		UnlockTimes: Unlocks{
			UnlockTimes: createUnlockTimes(),
			Unlocked:    unlocked,
		},

		AmalgamationSubjectIds: []json.Number{"2", "100"},

		PotentiallyAvailableAt:  potentiallyAvailableAt,
		UnlockByRadicalComputed: potentiallyAvailableAt != time.Time{},
	}
}

func createKanjiProgression(unlocked bool, unlockTimes []time.Time) Progression {
	var srsStage int64
	if unlocked {
		srsStage = 2
	} else {
		srsStage = 0
	}
	return Progression{
		SubjectID:              "2",
		Characters:             "Y",
		SrsStage:               srsStage,
		SrsSystem:              "1",
		UnlockedAt:             time.Time{},
		PassedAt:               time.Time{},
		AvailableAt:            time.Time{},
		PotentiallyAvailableAt: time.Time{},
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

func TestUpdateOptimalUnlockTimes(t *testing.T) {
	type args struct {
		spacedRepetitionSystems map[string]data.SpacedRepetitionSystem
		progressions            Progressions
	}
	tests := []struct {
		name string
		args args
		want []Unlocks
	}{
		{
			name: "Optimal unlock times with unlocked radical and locked kanji",
			args: args{
				spacedRepetitionSystems: createSrsMap(),
				progressions: Progressions{
					RadicalProgression: []Progression{
						createRadicalProgression(getTime("02-02-2020 00:00"), time.Time{}, true),
					},
					KanjiProgression: []Progression{
						createKanjiProgression(false, make([]time.Time, 10)),
					},
				},
			},
			want: []Unlocks{
				{
					UnlockTimes: createUnlockTimesForIndirectlyUpdatedKanji(),
					Unlocked:    false,
				},
			},
		},
		{
			name: "Optimal unlock times with locked radical and locked kanji",
			args: args{
				spacedRepetitionSystems: createSrsMap(),
				progressions: Progressions{
					RadicalProgression: []Progression{
						createRadicalProgression(time.Time{}, getTime("03-03-2020 00:00"), false),
					},
					KanjiProgression: []Progression{
						createKanjiProgression(false, make([]time.Time, 10)),
					},
				},
			},
			want: []Unlocks{
				{
					UnlockTimes: createUnlockTimesForIndirectlyUpdatedKanjiViaLockedRadical(),
					Unlocked:    false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeNow = func() time.Time {
				return getTime("01-01-2020 00:00")
			}
			UpdateOptimalUnlockTimes(tt.args.spacedRepetitionSystems, &tt.args.progressions)

			got := make([]Unlocks, len(tt.args.progressions.KanjiProgression))
			for idx, progression := range tt.args.progressions.KanjiProgression {
				got[idx] = progression.UnlockTimes
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateOptimalUnlockTimes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createUnlockTimesForIndirectlyUpdatedKanji() []time.Time {
	unlockTimes := make([]time.Time, 10)

	unlockTimes[0] = time.Time{}
	unlockTimes[1] = getTime("02-02-2020 00:00")
	unlockTimes[2] = getTime("02-02-2020 04:00")
	unlockTimes[3] = getTime("02-02-2020 12:00")
	unlockTimes[4] = getTime("03-02-2020 11:00")
	unlockTimes[5] = getTime("10-02-2020 10:00")
	unlockTimes[6] = getTime("24-02-2020 09:00")
	unlockTimes[7] = getTime("25-03-2020 08:00")
	unlockTimes[8] = getTime("23-07-2020 07:00")
	unlockTimes[9] = time.Time{}

	return unlockTimes
}

func createUnlockTimesForIndirectlyUpdatedKanjiViaLockedRadical() []time.Time {
	unlockTimes := make([]time.Time, 10)

	unlockTimes[0] = time.Time{}
	unlockTimes[1] = getTime("11-03-2020 10:00")
	unlockTimes[2] = getTime("11-03-2020 14:00")
	unlockTimes[3] = getTime("11-03-2020 22:00")
	unlockTimes[4] = getTime("12-03-2020 21:00")
	unlockTimes[5] = getTime("19-03-2020 20:00")
	unlockTimes[6] = getTime("02-04-2020 19:00")
	unlockTimes[7] = getTime("02-05-2020 18:00")
	unlockTimes[8] = getTime("30-08-2020 17:00")
	unlockTimes[9] = time.Time{}

	return unlockTimes
}
