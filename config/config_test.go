package config

import (
	"fmt"
	"testing"
	"time"
)

const CommonCfgFile = "fixtures/test_load.yaml"

func TestLoadConfig(t *testing.T) {
	tmp := C
	defer func() {
		C = tmp
	}()
	C = Config{}
	CfgFile = CommonCfgFile
	initConfig()
	t.Logf("Loaded: %v", CfgFile)
	if C.ClientId == string("") {
		t.Errorf("Expected C.ClientId not empty got %v\n", C)
	}
	if C.ConfigVersion == string("") {
		t.Errorf("Expected C.ConfigVersion not empty got %v\n", C)
	}

}

func TestGetStringMapStringTemplatedFromMap(t *testing.T) {
	tmp := C
	defer func() {
		C = tmp
	}()
	C = LoadCfgForTests(t, CommonCfgFile)
	from := make(map[string]string)
    from["Example"] = "example1"
    from["ClientId"] = "clientId"

	cases := []struct {
		key  string
		want string
	}{
		{
			key:  "example",
			want: from["Example"],
		},
		{
			key:  "a",
			want: "a",
		},
		{
			key:  "c",
			want: "1",
		},
		{
			// WARNING: viper downgrade to lower case
			key:  "clientid",
			want: from["ClientId"],
		},
	}
	res := GetStringMapStringTemplatedFromMap("TestGetStringMapStringTemplatedFromMap", "param", from)
	for _, tc := range cases {
		if got, ok := res[tc.key]; ok {
			if res[tc.key] != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		} else {
			t.Errorf("want key %v, value %v in res %v", tc.key, tc.want, res)
		}
	}
}



func TestGetStringMapStringTemplatedDefault(t *testing.T) {
	tmp := C
	defer func() {
		C = tmp
	}()
	C = LoadCfgForTests(t, CommonCfgFile)
	def := make(map[string]string)
	def["example"] = "example1"
	cases := []struct {
		key  string
		want string
	}{
		{
			key:  "example",
			want: def["example"],
		},
		{
			key:  "a",
			want: "a",
		},
		{
			key:  "c",
			want: "1",
		},
		{
			// WARNING: viper downgrade to lower case
			key:  "clientid",
			want: "clientId",
		},
	}
	res := GetStringMapStringTemplatedDefault("TestGetStringMapStringTemplatedDefault", "param", def)
	for _, tc := range cases {
		if got, ok := res[tc.key]; ok {
			if res[tc.key] != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		} else {
			t.Errorf("want key %v, value %v in res %v", tc.key, tc.want, res)
		}
	}
}

func TestGetStringMapStringTemplated(t *testing.T) {
	tmp := C
	defer func() {
		C = tmp
	}()
	C = LoadCfgForTests(t, CommonCfgFile)
	cases := []struct {
		key  string
		want string
	}{
		{
			key:  "a",
			want: "a",
		},
		{
			key:  "c",
			want: "1",
		},
		{
			// WARNING: viper downgrade to lower case
			key:  "clientid",
			want: "clientId",
		},
	}
	res := GetStringMapStringTemplated("TestGetStringMapStringTemplatedDefault", "param")
	for _, tc := range cases {
		if got, ok := res[tc.key]; ok {
			if res[tc.key] != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		} else {
			t.Errorf("want key %v, value %v in res %v", tc.key, tc.want, res)
		}
	}
}

func TestGetTimeDuration(t *testing.T) {
	tmp := C
	defer func() {
		C = tmp
	}()
	C = LoadCfgForTests(t, CommonCfgFile)
	cases := []struct {
		key  string
		want time.Duration
	}{
		{
			key:  "one",
			want: 15 * time.Second,
		},
		{
			key:  "two",
			want: 15 * time.Hour,
		},
		{
			key:  "three",
			want: 100 * time.Second,
		},
	}
	for _, tc := range cases {
		got := GetTimeDuration(fmt.Sprintf("%v.%v.%v", "TestGetTimeDuration", tc.key, CFG_INTERVAL_PARAMETER))
		if got != tc.want {
			t.Errorf("want %v, got %v", tc.want, got)
		}
	}
}

func TestGetTimeDurationEmpty(t *testing.T) {
	tmp := C
	defer func() {
		C = tmp
	}()
	C = LoadCfgForTests(t, "fixtures/empty.yaml")

	want := 1 * time.Second
	got := GetTimeDuration("")
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}
