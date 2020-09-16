package config

import (
	"github.com/spf13/viper"

	"testing"
)

func LoadCfgForTests(t *testing.T, CfgFile string) Config {
	tmp := C
	defer func() {
		C = tmp
	}()
	C = Config{}
	CfgFile = CfgFile
	viper.SetConfigFile(CfgFile)
	// t.Logf("Loaded: %v", CfgFile)
	if err := viper.ReadInConfig(); err != nil {
		t.Errorf("Can't read config: %v", err)
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		t.Errorf("unable to decode into struct, %v", err)

	}

	if C.ClientId == string("") {
		t.Errorf("Expected C.ClientId not empty got %v\n", C)
	}
	if C.ConfigVersion == string("") {
		t.Errorf("Expected C.ConfigVersion not empty got %v\n", C)
	}
	return C
}
