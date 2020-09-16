package config

import (
	"testing"
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

func TestGetStringMapStringTemplatedDefault(t *testing.T) {
    tmp := C
	defer func() {
		C = tmp
	}()
    C =LoadCfgForTests(t, CommonCfgFile)
    def :=make(map[string]string)
    def["example"] = "example1"
    cases := []struct {
        key   string
        want string
    }{
        {
            key:   "example",
            want: def["example"],
        },
        {
            key:   "a",
            want: "a",
        },
        {
            key:   "c",
            want: "1",
        },
        {
            // WARNING: viper downgrade to lower case
            key:   "clientid",
            want: "clientId",
        },

    }
    res:=GetStringMapStringTemplatedDefault("TestGetStringMapStringTemplatedDefault","param", def)
    for _, tc := range cases {
        if got, ok:= res[tc.key];ok{
            if res[tc.key] != tc.want {
                t.Errorf("want %v, got %v", tc.want, got)
            }
        }else {
            t.Errorf("want key %v, value %v in res %v", tc.key, tc.want, res)
        }
	}
}
