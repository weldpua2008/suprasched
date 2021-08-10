# suprasched
[![Build Status](https://travis-ci.org/weldpua2008/suprasched.svg?branch=master)](https://travis-ci.org/weldpua2008/suprasched) ![GitHub All Releases](https://img.shields.io/github/downloads/weldpua2008/suprasched/total) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Go Report Card](https://goreportcard.com/badge/github.com/weldpua2008/suprasched)](https://goreportcard.com/report/github.com/weldpua2008/suprasched) [![Docker Pulls](https://img.shields.io/docker/pulls/weldpua2008/suprasched)](https://hub.docker.com/r/weldpua2008/suprasched) ![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/weldpua2008/suprasched?label=docker%20image)

Suprasched is schduler  for Supra Platform

### How it works?
Suprasched is using [etcd](https://etcd.io/).
You can bring from one to three instances. Only one will be active. 

### Use cases
a). Assigning Jobs

Once we have a new job it should be assigned ASAP to a cluster

1. Get notifications about new Unassigned Job from Storage Layer - enrich the metadata: 
- Assess the number of required resources.
- Define supported clusters types
2. Send an event 
3. _Cluster Event Handler_ checks whether we have a ready cluster. Sends a new event in case we do not have a ready cluster
4. _Assigning Event Handler_ assignes the job to a ready cluster once 

b). Timeout Jobs on failed Clusters

Once the cluster failed - we should cancel the jos that are not finished otherwise they will stuck

1. Refresh the cluster status and fire an event if cluster entered failed state
2. Check all jobs assigned to the cluster that are not finished
3. Cancel the jobs

c). Scale up/down clusters

From time to time we will be in the state when we have a queue of the jobs on the same cluster.
In some use-cases we will decide to create and move the PENDING jobs to the new cluster

c). Scale up/down cluster's size

Each job has its requrements. And the cluster has its capacity.
We should scale in/out the cluster in regards to the job that are running and queued one.


### Control Plane Components 
The control plane's components make global decisions about the clusters (for example, scheduling), as well as detecting 
and responding to cluster events (for example, starting up a new cluster when all current are overprovisioned).

Control plane components can be run on any machine. However, for simplicity, set up scripts typically start all control plane components on the same machine. 

#### etcd
Consistent and highly-available key value store used as Suprascheduler' backing store for all cluster data.
Make sure you have a [back up plan for those data](https://etcd.io/docs/v2.3/admin_guide/#backing-up-the-datastore).

#### Joblet

#### apiserver
The API server is a front end component of the control plane that exposes the API. 

#### scheduler
Control plane component that watches for newly created Jobs with no assigned cluster, and selects a cluster for them to run on.

Factors taken into account for scheduling decisions include: individual and collective resource requirements, hardware/software/policy constraints, 
affinity and anti-affinity specifications, regions, inter-workload interference, and deadlines.

#### controller-manager
Control plane component that runs controller processes.

Some types of these controllers are:
* Cluster controller: Responsible for noticing and responding when cluster go down, has no sufficient resources.
* Job controller: Creates and watches to completion for Job objects that represent one-off tasks.


