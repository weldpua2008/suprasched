package core

import (
	"sort"
	"testing"
)

func TestSortJobSteps(t *testing.T) {
	cases := []struct {
		jobSteps []JobStep
		rightCmd []string
	}{
		{
			jobSteps: []JobStep{
				{CMD: "1", StepPriority: 1},
				{CMD: "2", StepPriority: 2},
			},
			rightCmd: []string{"1", "2"},
		},
		{
			jobSteps: []JobStep{
				{CMD: "2", StepPriority: 2},
				{CMD: "1", StepPriority: 1},
			},
			rightCmd: []string{"1", "2"},
		},
	}
	for _, tc := range cases {
		if len(tc.jobSteps) != len(tc.rightCmd) {
			t.Fatalf("len mismatch %d != %d", len(tc.jobSteps), len(tc.rightCmd))
		}
		sort.Sort(JobStepsSorter(tc.jobSteps))

		for i, val := range tc.rightCmd {
			if val != tc.jobSteps[i].CMD {
				t.Errorf("[%d] %s != %s", i, val, tc.jobSteps[i].CMD)
			}
		}

	}
}
