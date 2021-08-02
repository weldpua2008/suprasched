# suprasched
[![Build Status](https://travis-ci.org/weldpua2008/suprasched.svg?branch=master)](https://travis-ci.org/weldpua2008/suprasched) ![GitHub All Releases](https://img.shields.io/github/downloads/weldpua2008/suprasched/total) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Go Report Card](https://goreportcard.com/badge/github.com/weldpua2008/suprasched)](https://goreportcard.com/report/github.com/weldpua2008/suprasched) [![Docker Pulls](https://img.shields.io/docker/pulls/weldpua2008/suprasched)](https://hub.docker.com/r/weldpua2008/suprasched) ![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/weldpua2008/suprasched?label=docker%20image)

Suprasched is rescheduler for Supra Platform.

### How it works?
Rescheduler is the service to rescue stuck jobs and help to reach a terminal state.

When the service is over-engineering?:
* 1 job per cluster. Cluster bootstrap time is 0. Job starts and ends always
* it's totally fine if jobs scheduled on the cluster won't finish due the cluster termination

Motivation:
* help utilize the clusters
* resolve any stuck states as for jobs and for clusters
* got the best option of the following?:
1). Spin a new job on the existing cluster  
2).Spin a new job on a new cluster

### High-level overview

a). Enriching Job info
1. Fetch a new Pending Job. Assess the number of resources.
2. Fetch clusters utilization (background)  
3. Fetch other Job statuses, resource usage and approximate finish time (background)
4. Calculate probabilities (when cluster will be shut down)
5. Prediction service based on the history data:
- Prediction Job execution time
- Cluster issues based on any current issue
- Anticipate the future load and pre-warm cluster resources (e.g. if cluster creation > 3-5 minutes)

b). Plan the execution plan in regards to the full state snapshot.
Determine a Job execution time vector (e.g. on top of a new cluster, reprioritize some task, create a new cluster, wait for some task to finish)

c). Choose the best option based on Job requirements.

d). Schedule the Job, Prepare Cluster

e). In case the Job failed - back to the (a)

f). In case the Job succeeded/canceled - free resources and back to (b)


### Events

Job creation        -> Assign to Cluster
                    -> Forecast load
Job Status          -> Change Cluster Capacity
                    -> Modify Capacity
                   (->) Reschedule existing job


Job Time            -> Cancel/Reschedule

Cluster Death       -> Cluster Creation
                    -> Job rescheduling/cancelation

Cluster Capacity    -> Job rescheduling
                    -> Cluster Creation

Forecast            -> Cluster Creation/Modify Capacity
                    -> Job rescheduling
                    -> Put cluster in 'Terminating' mode ()


### Relationship
New Job --> unassigned cluster
unassigned cluster --> New cluster
New Job --> Cluster
Cluster (termination): Job --> unassigned cluster


### Installing from source

1. install [Go](http://golang.org) `v1.14+`
1. clone this down into your `$GOPATH`
  * `mkdir -p $GOPATH/src/github.com/weldpua2008`
  * `git clone https://github.com/weldpua2008/suprasched $GOPATH/src/github.com/weldpua2008/suprasched`
  * `cd $GOPATH/src/github.com/weldpua2008/suprasched`
1. install [golangci-lint](https://github.com/golangci/golangci-lint#install) for linting + static analysis
  * Lint: `docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.24.0 golangci-lint run -v`
### Running tests

*  expires all test results

```bash
$ go clean -testcache
```
* run all tests

```bash
$ go test -bench= -test.v  ./...

$ go test  -bench=. -benchmem -v  ./...

```
