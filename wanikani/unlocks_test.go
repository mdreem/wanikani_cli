package wanikani

import (
	"reflect"
	"testing"
	"time"
	"wanikani_cli/data"
)

func TestComputeOptimalUnlocks(t *testing.T) {
	srs := data.SpacedRepetitionSystem{
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

	type args struct {
		system      data.SpacedRepetitionSystem
		progression Progression
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
			},
			want: Unlocks{
				UnlockTimes: createUnlockTimes(),
				Unlocked:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeNow = func() time.Time {
				return getTime("01-01-2020 00:00")
			}
			if got := computeOptimalUnlocks(tt.args.system, tt.args.progression); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeOptimalUnlocks() = %v, want %v", got, tt.want)
			}
		})
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

func getTime(timeString string) time.Time {
	parsedTime, _ := time.Parse("02-01-2006 15:04", timeString)
	return parsedTime
}
