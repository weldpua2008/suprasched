package communicator

import (
	config "github.com/weldpua2008/suprasched/config"

	"errors"
	"testing"
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

// func loadCfg(t *testing.T, CfgFile string) {
// 	tmp := config.C
// 	defer func() {
// 		config.C = tmp
// 	}()
// 	config.C = config.Config{}
// 	config.CfgFile = CfgFile
// 	viper.SetConfigFile(config.CfgFile)
// 	// t.Logf("Loaded: %v", config.CfgFile)
// 	if err := viper.ReadInConfig(); err != nil {
// 		t.Errorf("Can't read config: %v", err)
// 	}
// 	err := viper.Unmarshal(&config.C)
// 	if err != nil {
// 		t.Errorf("unable to decode into struct, %v", err)
//
// 	}
//
// 	if config.C.ClientId == string("") {
// 		t.Errorf("Expected C.ClientId not empty got %v\n", config.C)
// 	}
// 	if config.C.ConfigVersion == string("") {
// 		t.Errorf("Expected C.ConfigVersion not empty got %v\n", config.C)
// 	}
//
// }

func TestGetSectionCommunicator(t *testing.T) {
	config.LoadCfgForTests(t, "fixtures/http.yml")

	cases := []struct {
		section string
		in      string
		want    error
	}{
		{
			section: "http",
			in:      "HTTP",
			want:    nil,
		},
		{
			section: "http",
			in:      "http",
			want:    nil,
		},
		{
			section: "http_capital",
			in:      "HTTP",
			want:    nil,
		},
		{
			section: "http_capital",
			in:      "http",
			want:    nil,
		},

		{
			section: "broken",
			in:      "broken",
			want:    ErrNoSuitableCommunicator,
		},
	}

	for _, tc := range cases {
		result, got := GetSectionCommunicator(tc.in)
		if (tc.want == nil) && (tc.want != got) {
			t.Errorf("want %v, got %v", tc.want, got)
        } else if if (tc.want == nil) && (!result.Configured()) {
				t.Errorf("want %v, got %v, res %v", true, result.Configured(), result)

		} else {
			if !errors.Is(got, tc.want) {
				t.Errorf("want %v, got %v, res %v", tc.want, got, result)
			}
		}
	}
}
