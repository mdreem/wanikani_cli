package wanikani

import (
	"encoding/json"
	"github.com/mdreem/wanikani_cli/data"
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
			var assignments []data.AssignmentEnvelope
			for _, srsLevel := range tt.srsLevels {
				assignments = append(assignments, createAssignment(srsLevel))
			}
			testClient := data.TestClient{
				Assignments: assignments,
			}

			if got := fetchDistribution(testClient, tt.subjectType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchDistribution() = %v\n, want = %v", got, tt.want)
			}
		})
	}
}

func createAssignment(srsLevel json.Number) data.AssignmentEnvelope {
	return data.AssignmentEnvelope{
		Data: data.Assignment{
			SrsStage: srsLevel,
		},
	}
}
