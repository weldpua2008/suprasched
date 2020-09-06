package cluster

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
	"sync"
	"time"
)

type DescribeEMR struct {
	ClusterDescriber
	aws_sessions map[string]*session.Session
	mu           sync.RWMutex
}

// NewDescribeEMR prepare struct communicator for EMR
func NewDescribeEMR() *DescribeEMR {
	s := make(map[string]*session.Session)
	return &DescribeEMR{
		aws_sessions: s,
	}
}

// getAwsSession
func (c *DescribeEMR) getAwsSession(params map[string]interface{}) (*session.Session, error) {
	var Profile string
	Region := "us-east-1"
	for _, k := range []string{"AWS_PROFILE", "PROFILE", "aws_profile", "profile", "Profile"} {
		if _, ok := params[k]; ok {
			Profile = params[k].(string)
			break
		}
	}
	for _, k := range []string{"Region", "AWS_Region", "AWS_REGION", "aws_region", "region"} {
		if _, ok := params[k]; ok {
			Region = params[k].(string)
			break
		}
	}
	session_key := fmt.Sprintf("%v%v", Profile, Region)
	c.mu.Lock()
	defer c.mu.Unlock()

	if val, ok := c.aws_sessions[session_key]; ok {
		return val, nil
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		Profile: Profile,
		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region: aws.String(Region),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})
	if err == nil {
		c.aws_sessions[session_key] = sess
	}
	return sess, err
}

// DescribeCluster
func (c *DescribeEMR) DescribeCluster(params map[string]interface{}) (string, error) {
	var ClusterId string
	var ctx context.Context
	var clusterCtx context.Context
	var cancel context.CancelFunc
	ttr := 30

	for _, k := range []string{"ClusterId", "clusterID", "ClusterID", "clusterId",
		"clusterid", "JobFlowID", "JobFlowId", "JobflowID", "jobFlowId"} {
		if _, ok := params[k]; ok {
			ClusterId = params[k].(string)
			break
		}
	}
	for _, k := range []string{"context", "ctx"} {
		if _, ok := params[k]; ok {
			if v, ok := params[k].(context.Context); ok {
				ctx = v
				break
			}
		}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	clusterCtx, cancel = context.WithTimeout(ctx, time.Duration(ttr)*time.Second)
	defer cancel() // cancel when we are getting the kill signal or exit
	sess, err := c.getAwsSession(params)
	if err != nil {
		return "", err
	}
	svc := emr.New(sess)

	clusterInput := &emr.DescribeClusterInput{
		ClusterId: aws.String(ClusterId),
	}
	cl, err := svc.DescribeClusterWithContext(clusterCtx, clusterInput)
	if err != nil {
		return "", err
	}
	status := cl.Cluster.Status.State
	result := *status
	switch *status {
	case emr.ClusterStateStarting, emr.ClusterStateBootstrapping:
		result = "STARTING"
	case emr.ClusterStateRunning, emr.ClusterStateWaiting:
		result = "RUNNING"
	case emr.ClusterStateTerminated, emr.ClusterStateTerminatedWithErrors:
		result = "TERMINATED"
	}

	return result, nil
}
