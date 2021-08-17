/*
Copyright 2021 The Suprasched Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package app

import (
	"github.com/sirupsen/logrus"
	"time"
)

var (
	// Number of active jobs
	NumActiveJobs int64
	// Number of processed jobs
	NumProcessedJobs int64
	logFields        = logrus.Fields{"package": "joblet"}
	log              = logrus.WithFields(logFields)
)

const TimeoutJobsAfter5MinInTerminalState = 5 * time.Minute
const StopReadJobsOutputAfter5Min = 5 * time.Minute
const TimeoutAppendLogStreams = 10 * time.Minute
