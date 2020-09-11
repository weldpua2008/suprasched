// Copyright 2020 Valeriy Soloviov. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// license that can be found in the LICENSE file.

// Package cmd provides CLI interfaces for the `suprasched` application.
package cmd

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cluster "github.com/weldpua2008/suprasched/cluster"
	communicator "github.com/weldpua2008/suprasched/communicator"
	config "github.com/weldpua2008/suprasched/config"

	handlers "github.com/weldpua2008/suprasched/handlers"
	job "github.com/weldpua2008/suprasched/job"
	model "github.com/weldpua2008/suprasched/model"

	// worker "github.com/weldpua2008/suprasched/worker"
	"time"
	// "html/template"
	"os"
	"os/signal"
	// "sync"
	// "github.com/mustafaturan/bus"
	"syscall"
)

var (
	verbose    bool
	traceFlag  bool
	log            = logrus.WithFields(logrus.Fields{"package": "cmd"})
	numWorkers int = 5
)

func init() {

	// Define Persistent Flags and configuration settings, which, if defined here,
	// will be global for application.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
	rootCmd.PersistentFlags().BoolVarP(&traceFlag, "trace", "t", false, "trace")
	// rootCmd.PersistentFlags().StringVar(&config.ClientId, "clientId", "", "ClientId (default is suprasched)")

	// rootCmd.PersistentFlags().IntVarP(&numWorkers, "workers", "w", 5, "Number of workers")
	// local flags, which will only run
	// when this action is called directly.

	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	// Only log the warning severity or above.
	logrus.SetLevel(logrus.InfoLevel)
}

// This represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "suprasched",
	Short: "suprasched is mastermind for jobs",
	Long: `A Fast and Flexible Abstraction around jobs rescheduler built with
                love by weldpua2008 and friends in Go.
                Complete documentation is available at github.com/weldpua2008/suprasched/cmd`,
	Version: FormattedVersion(),
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 1)
		shutchan := make(chan bool, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // cancel when we are getting the kill signal or exit
		jobs := make(chan *model.Job, 1)
		clusters := make(chan *model.Cluster, 1)
		// var wg sync.WaitGroup
		// jobs := make(chan *model.Job, 1)
		log.Infof("Starting suprasched\n")
		log.Infof("%v", config.GetStringMapStringTemplated("cluster", "param"))
		go func() {
			sig := <-sigs
			log.Infof("Shutting down - got %v signal", sig)
			cancel()
			shutchan <- true
		}()

		if traceFlag {
			logrus.SetLevel(logrus.TraceLevel)
		} else if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}

		// load config
		// if errCnf := model.ReinitializeConfig(); errCnf != nil {
		// 	log.Tracef("Failed ReinitializeConfig %v\n", errCnf)
		// }
		config.ReinitializeConfig()
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Trace("Config file changed:", e.Name)
			// if errCnf := model.ReinitializeConfig(); errCnf != nil {
			// log.Tracef("Failed model.ReinitializeConfig %v\n", errCnf)
			// }
			config.ReinitializeConfig()
		})

		log.Trace("Config file:", viper.ConfigFileUsed())

		handlers.Init()

		go func() {
			// StartGenerateClusters(ctx context.Context, clusters chan *model.Cluster, interval time.Duration) error
			if err := cluster.StartGenerateClusters(ctx, clusters, config.GetApiDelayForSection(
				fmt.Sprintf(
					"%s.fetch.delay",
					config.CFG_PREFIX_CLUSTER,
				))); err != nil {
				log.Tracef("StartGenerateClusters returned error %v", err)
			}
		}()

		communicator_type := config.GetStringDefault(fmt.Sprintf("%s.fetch.communicator", config.JobsSection), "http")
		comm, err_com := communicator.GetCommunicator(communicator_type)
		if err_com == nil {
			go func() {
				log.Trace("StartFetchJobs ")
				if err := job.StartFetchJobs(
					ctx, comm, jobs, config.GetApiDelayForSection(
						fmt.Sprintf(
							"%s.fetch.delay",
							config.JobsSection,
						)),
				); err != nil {
					log.Tracef("StartFetchJobs returned error %v", err)
				}
			}()
		} else {
			close(jobs)
		}

		//
		// for w := 1; w <= numWorkers; w++ {
		// 	wg.Add(1)
		// 	go worker.StartWorker(w, jobs, &wg)
		// }
		//
		// wg.Wait()
		time.Sleep(150 * time.Millisecond)
		time.Sleep(6500 * time.Millisecond)

	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// return error
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
