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
	model "github.com/weldpua2008/suprasched/model"
	"html/template"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	// ProjectName defines project name
	ProjectName = "suprasched"
	// JobsSection for the config
	JobsSection = "jobs"
	//CFG_PREFIX_COMMUNICATOR defines parameter in the config for Communicators
	CFG_PREFIX_COMMUNICATOR     = "communicator"
	CFG_PREFIX_COMMUNICATORS    = "communicators"
	CFG_PREFIX_CLUSTER          = "cluster"
	CFG_PREFIX_FETCHER          = "fetch"
	CFG_COMMUNICATOR_PARAMS_KEY = "params"
)

var (
	JobsRegistry    = model.NewRegistry()
	ClusterRegistry = model.NewClusterRegistry()
)

// Config is top level Configuration structure
type Config struct {
	// Indentification for the process
	ClientId string `mapstructure:"clientId"`
	// delay between API calls to prevent Denial-of-service
	CallAPIDelaySec int `mapstructure:"api_delay_sec"`
	// Config version
	ConfigVersion string `mapstructure:"version"`
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
	InitEvenBus()

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

// GetSliceStringMapStringTemplatedDefault returns slice of [sting]sting maps templated & enriched by default.
func GetSliceStringMapStringTemplatedDefault(section string, param string, def map[string]string) []map[string]string {
	ret := make([]map[string]string, 0)
	sections_values := viper.GetStringMap(fmt.Sprintf("%s.%s", section, param))
	for _, section_value := range sections_values {
		if section_value == nil {
			continue
		}
		if params, ok := section_value.(map[string]interface{}); ok {
			c := make(map[string]string)
			for k, v := range def {
				c[k] = v
			}
			for k, v := range params {
				var tplBytes bytes.Buffer
				tpl := template.Must(template.New("params").Parse(fmt.Sprintf("%v", v)))
				if err := tpl.Execute(&tplBytes, C); err != nil {
					log.Tracef("params executing template for %v got %s", v, err)
					continue
				}
				c[k] = tplBytes.String()
			}
			ret = append(ret, c)
		}

	}
	return ret
}

func GetStringMapStringTemplatedDefault(section string, param string, def map[string]string) map[string]string {
	c := make(map[string]string)
	for k, v := range def {
		c[k] = v
	}
	params := viper.GetStringMapString(fmt.Sprintf("%s.%s", section, param))
	for k, v := range params {

		var tplBytes bytes.Buffer
        // WARNING: will panic:
        // tpl := template.Must(template.New("params").Parse(v))
        // we can preserve failed templated string
        c[k] = v
		tpl, err1 := template.New("params").Parse(v)
        if err1!=nil {
            continue
        }
		err := tpl.Execute(&tplBytes, C)
		if err != nil {
			log.Tracef("params executing template: %s", err)
			continue
		}
		c[k] = tplBytes.String()
	}
	return c
}

func ConvertMapStringToInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range in {
		out[k] = v
	}
	return out
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
