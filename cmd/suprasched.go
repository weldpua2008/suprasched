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
	"math/rand"
	// communicator "github.com/weldpua2008/suprasched/communicator"
	config "github.com/weldpua2008/suprasched/config"

	handlers "github.com/weldpua2008/suprasched/handlers"
	job "github.com/weldpua2008/suprasched/job"
	metrics "github.com/weldpua2008/suprasched/metrics"
	// model "github.com/weldpua2008/suprasched/model"

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
	verbose              bool
	traceFlag            bool
	enableHealthCheckUri bool
	log                  = logrus.WithFields(logrus.Fields{"package": "cmd"})
	//numWorkers           int = 5
)

func init() {

	// Define Persistent Flags and configuration settings, which, if defined here,
	// will be global for application.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
	rootCmd.PersistentFlags().BoolVarP(&traceFlag, "trace", "t", false, "trace")
	rootCmd.PersistentFlags().BoolVarP(&enableHealthCheckUri, "healthcheck", "p", true, "healthcheck")

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
		rand.Seed(time.Now().UnixNano())
		sigs := make(chan os.Signal, 1)
		shutdownChan := make(chan bool, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // cancel when we are getting the kill signal or exit
		jobs := make(chan bool)
		jobsTimeOut := make(chan bool)
		clusters := make(chan bool)
		describers := make(chan bool)
		terminators := make(chan bool)
		// var wg sync.WaitGroup
		// jobs := make(chan *model.Job, 1)
		log.Infof("Starting suprasched\n")
		go func() {
			sig := <-sigs
			log.Infof("Shutting down - got %v signal", sig)
			cancel()
			shutdownChan <- true
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
		if enableHealthCheckUri {
			addr := config.GetStringTemplatedDefault("healthcheck.listen", ":8080")
			metricsUri := config.GetStringTemplatedDefault("healthcheck.uri", "/health/is_alive")
			metrics.StartHealthCheck(addr, metricsUri)

		}
		prometheusAddr := config.GetStringTemplatedDefault("prometheus.listen", ":8080")
		prometheusUri := config.GetStringTemplatedDefault("prometheus.uri", "/metrics")
		metrics.AddPrometheusMetricsHandler(prometheusAddr, prometheusUri)
		metrics.StartAll()
		defer metrics.StopAll(ctx)
		// Init Bus & Handlers
		handlers.Init()
		defer handlers.Deregister()
		defer config.EvenBusTearDown()

		go func() {
			if err := cluster.StartGenerateClusters(ctx, clusters, config.GetTimeDuration(
				fmt.Sprintf(
					"%s.fetch",
					config.CFG_PREFIX_CLUSTER,
				))); err != nil {
				log.Warningf("StartGenerateClusters returned error %v", err)
			}
		}()

		go func() {
			if err := job.StartFetchJobs(ctx, jobs, config.GetTimeDuration(
				fmt.Sprintf(
					"%s.%s",
					config.CFG_PREFIX_JOBS,
					config.CFG_PREFIX_JOBS_FETCHER,
				))); err != nil {
				log.Warningf("StartFetchJobs returned error %v", err)
			}
		}()

		go func() {
			if err := cluster.StartUpdateClustersMetadata(ctx, describers, config.GetTimeDuration(
				fmt.Sprintf(
					"%s.%s",
					config.CFG_PREFIX_CLUSTER,
					config.CFG_PREFIX_DESCRIBERS,
				))); err != nil {
				log.Warningf("StartUpdateClustersMetadata returned error %v", err)
			}
		}()

		go func() {
			terminationDelay := 2*config.GetTimeDuration(
				fmt.Sprintf(
					"%s.%s",
					config.CFG_PREFIX_JOBS,
					config.CFG_PREFIX_JOBS_FETCHER,
				)) + 2*config.GetTimeDuration(
				fmt.Sprintf(
					"%s.fetch",
					config.CFG_PREFIX_CLUSTER,
				))

			if err := cluster.StartTerminateClusters(ctx, terminators, config.GetTimeDuration(
				fmt.Sprintf(
					"%s.%s",
					config.CFG_PREFIX_CLUSTER,
					config.CFG_PREFIX_TERMINATORS,
				)),
				terminationDelay); err != nil {
				log.Warningf("StartTerminateClusters returned error %v", err)
			}
		}()
		go func() {
			if err := job.CancelTimeoutJobs(ctx, jobsTimeOut,
				config.GetTimeDuration(
					fmt.Sprintf(
						"%s.%s",
						config.CFG_PREFIX_JOBS,
						config.CFG_PREFIX_JOBS_TIMEOUT,
					)),
				config.GetTimeDurationDefault(
					fmt.Sprintf(
						"%s.%s",
						config.CFG_PREFIX_JOBS,
						config.CFG_PREFIX_JOBS_TIMEOUT,
					),
					config.CFG_PREFIX_JOB_TIMEOUT_DURATION,
					time.Duration(30)*time.Minute)); err != nil {
				log.Warningf("CancelTimeoutJobs returned error %v", err)
			}
		}()

		select {
		case <-jobsTimeOut:
			log.Infof("Jobs Timeout stopped")

		case <-jobs:
			log.Infof("Jobs stopped")
		case <-clusters:
			log.Infof("Clusters stopped")
		case <-describers:
			log.Infof("Describers stopped")
		case <-terminators:
			log.Infof("Terminators stopped")

		case <-ctx.Done():
			if ctx.Err() != nil {
				log.Tracef("Context cancelled")
			}
		}
		log.Infof("Gracefully Shutdown")
		time.Sleep(150 * time.Millisecond)

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
