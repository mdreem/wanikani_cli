package progressinfo

import (
	"encoding/json"
	"github.com/mdreem/wanikani_cli/wanikani"
	data2 "github.com/mdreem/wanikani_cli/wanikani/data"
	"reflect"
	"testing"
)

func Test_fetchDistribution(t *testing.T) {
	tests := []struct {
		name        string
		subjectType string
		srsLevels   []json.Number
		want        []int
	}{
		{
			name:        "compute distribution with srsLevel 0",
			subjectType: "kanji",
			srsLevels:   []json.Number{"0", "0"},
			want:        []int{2, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:        "compute distribution with srsLevel 9",
			subjectType: "kanji",
			srsLevels:   []json.Number{"9", "9"},
			want:        []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
		},
		{
			name:        "compute distribution with every srsLevel set",
			subjectType: "kanji",
			srsLevels:   []json.Number{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
			want:        []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var assignments []data2.AssignmentEnvelope
			for _, srsLevel := range tt.srsLevels {
				assignments = append(assignments, createAssignment(srsLevel))
			}
			testClient := wanikani.TestClient{
				Assignments: assignments,
			}

			if got := fetchDistribution(testClient, tt.subjectType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchDistribution() = %v\n, want = %v", got, tt.want)
			}
		})
	}
}

func createAssignment(srsLevel json.Number) data2.AssignmentEnvelope {
	return data2.AssignmentEnvelope{
		Data: data2.Assignment{
			SrsStage: srsLevel,
		},
	}
}

func Test_ComputeProgressInfo(t *testing.T) {
	tests := []struct {
		name        string
		subjectType string
		srsLevels   []json.Number
		want        ProgressionData
	}{
		{
			name:        "compute progress info with srsLevel 0",
			subjectType: "kanji",
			srsLevels:   []json.Number{"0", "0"},
			want: ProgressionData{
				NumNotStarted:  2,
				NumApprentice:  0,
				NumGuru:        0,
				NumMaster:      0,
				NumEnlightened: 0,
				NumBurned:      0,
			},
		},
		{
			name:        "compute progress info with srsLevel 9",
			subjectType: "kanji",
			srsLevels:   []json.Number{"9", "9"},
			want: ProgressionData{
				NumNotStarted:  2,
				NumApprentice:  2,
				NumGuru:        2,
				NumMaster:      2,
				NumEnlightened: 2,
				NumBurned:      2,
			},
		},
		{
			name:        "compute progress info with every srsLevel set",
			subjectType: "kanji",
			srsLevels:   []json.Number{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
			want: ProgressionData{
				NumNotStarted:  10,
				NumApprentice:  9,
				NumGuru:        5,
				NumMaster:      3,
				NumEnlightened: 2,
				NumBurned:      1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var assignments []data2.AssignmentEnvelope
			for _, srsLevel := range tt.srsLevels {
				assignments = append(assignments, createAssignment(srsLevel))
			}
			testClient := wanikani.TestClient{
				Assignments: assignments,
			}

			if got := ComputeProgressInfo(testClient, tt.subjectType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeProgressInfo() = %v\n, want = %v", got, tt.want)
			}
		})
	}
}
