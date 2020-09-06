// Copyright 2020 Valeriy Soloviov. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// license that can be found in the LICENSE file.

// Package config provides configuration for `suprasched` application.
package config

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"html/template"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	// ProjectName defines project name
	ProjectName = "suprasched"
	// JobsFetchSection
	JobsSection = "jobs"
)

// Config is top level Configuration structure
type Config struct {
	// Indentification for the process
	ClientId string `mapstructure:"clientId"`
	// delay between API calls to prevent Denial-of-service
	CallAPIDelaySec int `mapstructure:"api_delay_sec"`

	HeartBeat ApiOperations `mapstructure:"heartbeat"`
	// Config version
	ConfigVersion string `mapstructure:"version"`
}

// ApiOperations is defines operations structure
type ApiOperations struct {
	Run         UrlConf `mapstructure:"run"`         // defines how to run item
	Cancelation UrlConf `mapstructure:"cancelation"` // defines how to cancel item

	LogStreams UrlConf `mapstructure:"logstream"` // defines how to get item

	Get    UrlConf `mapstructure:"get"`    // defines how to get item
	Lock   UrlConf `mapstructure:"lock"`   // defines how to lock item
	Update UrlConf `mapstructure:"update"` // defines how to update item
	Unlock UrlConf `mapstructure:"unlock"` // defines how to unlock item
	Finish UrlConf `mapstructure:"finish"` // defines how to finish item
	Failed UrlConf `mapstructure:"failed"` // defines how to update on failed
	Cancel UrlConf `mapstructure:"cancel"` // defines how to update on cancel

}

// UrlConf defines all params for request.
type UrlConf struct {
	Url             string            `mapstructure:"url"`
	Method          string            `mapstructure:"method"`
	Headers         map[string]string `mapstructure:"headers"`
	PreservedFields map[string]string `mapstructure:"preservedfields"`
	Params          map[string]string `mapstructure:"params"`
}

var (
	// CfgFile defines Path to the config
	CfgFile string
	// ClientId defines Indentification for the instance.
	ClientId string
	// C defines main configuration structure.
	C Config = Config{
		CallAPIDelaySec: int(2),
	}
	log = logrus.WithFields(logrus.Fields{"package": "config"})
)

// Init configuration
func init() {
	cobra.OnInitialize(initConfig)

}

// ReinitializeConfig on load or file change
func ReinitializeConfig() {
	if len(ClientId) > 0 {
		C.ClientId = ClientId
	}
	if len(C.ClientId) < 1 {
		C.ClientId = "suprasched"
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Don't forget to read config either from CfgFile or from home directory!
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		lProjectName := strings.ToLower(ProjectName)
		log.Debug("Searching for config with project", ProjectName)
		viper.AddConfigPath(".")
		viper.AddConfigPath("..")
		switch runtime.GOOS {
		case "windows":
			if userprofile := os.Getenv("USERPROFILE"); userprofile != "" {
				viper.AddConfigPath(userprofile)
			}
		default:
			// freebsd, openbsd, darwin, linux
			// plan9, windows...
			viper.AddConfigPath("$HOME/")
			viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", lProjectName))
			viper.AddConfigPath("/etc/")
			viper.AddConfigPath(fmt.Sprintf("/etc/%s/", lProjectName))

		}

		if conf := os.Getenv(fmt.Sprintf("%s_CFG", strings.ToUpper(ProjectName))); conf != "" {
			viper.SetConfigName(conf)
		} else {
			viper.SetConfigType("yaml")
			viper.SetConfigName(lProjectName)
		}
	}
	viper.AutomaticEnv()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("Can't read config:", err)
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("unable to decode into struct, %v", err))

	}
	log.Debug(viper.ConfigFileUsed())

}

func GetStringMapStringTemplatedDefault(section string, param string, def map[string]string) map[string]string {
	// log.Tracef("Calling GetParamsFromSection(%s,%s)",section, param)
	c := make(map[string]string)
	for k, v := range def {
		c[k] = v
	}
	params := viper.GetStringMapString(fmt.Sprintf("%s.%s", section, param))
	for k, v := range params {
		var tplBytes bytes.Buffer
		tpl := template.Must(template.New("params").Parse(v))
		err := tpl.Execute(&tplBytes, C)
		if err != nil {
			log.Tracef("params executing template: %s", err)
			continue
		}
		c[k] = tplBytes.String()
	}
	return c
}

func GetStringMapStringTemplated(section string, param string) map[string]string {
	c := make(map[string]string)
	return GetStringMapStringTemplatedDefault(section, param, c)
}

// GetStringDefault return section string or default
func GetStringDefault(section string, def string) string {
	if val := viper.GetString(section); len(val) > 0 {
		return val
	}
	return def
}

// GetApiDelayForSection return api call delay in seconds for the section
func GetApiDelayForSection(section string) (interval time.Duration) {
	def := "api_delay_sec"
	delay := int64(viper.GetInt(section))
	if delay < 1 {
		delay = int64(viper.GetInt(def))
		if delay < 1 {
			delay = 1
		}
	}
	interval = time.Duration(delay) * time.Second

	return interval
}
