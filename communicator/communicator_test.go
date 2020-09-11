package communicator

import (
	// "context"
	// "fmt"
	"errors"
	"testing"
	// config "github.com/weldpua2008/suprasched/config"
)

func TestGetCommunicator(t *testing.T) {
	cases := []struct {
		in   string
		want error
	}{
		{
			in:   "HTTP",
			want: nil,
		},
		{
			in:   "http",
			want: nil,
		},
		{
			in:   "broken",
			want: ErrNoSuitableCommunicator,
		},
	}

	for _, tc := range cases {
		result, got := GetCommunicator(tc.in)
		if (tc.want == nil) && (tc.want != got) {
			t.Errorf("want %v, got %v", tc.want, got)

		} else {
			if !errors.Is(got, tc.want) {
				t.Errorf("want %v, got %v, res %v", tc.want, got, result)
			}
		}
	}
}
