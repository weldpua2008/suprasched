package utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	aws_request "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/emr/emriface"
	//
	"github.com/weldpua2008/suprasched/metrics"
	//
	"os"
	"strings"
	"time"
)

func GetCachedAwsSession(key string) (*session.Session, error) {

	mu.RLock()
	defer mu.RUnlock()
	if aws_sessions != nil {
		if val, ok := aws_sessions[key]; ok {
			return val, nil
		}
	}
	return nil, fmt.Errorf("Session %v is not in cache", key)
}

func GetAwsSession(params map[string]interface{}) (*session.Session, error) {
	var Profile string
	Region := "us-east-1"

	if len(os.Getenv("AWS_DEFAULT_REGION")) > 0 {
		Region = os.Getenv("AWS_DEFAULT_REGION")
	}

	for _, k := range []string{"AWS_PROFILE", "PROFILE", "aws_profile", "profile", "Profile", "ClusterProfile"} {
		if _, ok := params[k]; ok {
			Profile = params[k].(string)
			break
		}
	}
	for _, k := range []string{"Region", "AWS_Region", "AWS_REGION", "aws_region", "region", "ClusterRegion"} {
		if _, ok := params[k]; ok {
			Region = params[k].(string)
			break
		}
	}
	session_key := fmt.Sprintf("%v%v", Profile, Region)

	if val, err := GetCachedAwsSession(session_key); err == nil {
		return val, nil
	}
	// Creating & adding the session to the cache
	mu.Lock()
	defer mu.Unlock()
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
	sess.Handlers.Send.PushFront(func(r *aws_request.Request) {
		// Log every request made and its payload
		metrics.ApiCallsStatistics.WithLabelValues(
			"aws",
			fmt.Sprintf("%v.%v", Profile, Region),
			"emr",
			strings.ToLower(r.Operation.Name),
		).Inc()
	})

	if aws_sessions == nil {
		aws_sessions = make(map[string]*session.Session)
	}
	if err == nil {
		aws_sessions[session_key] = sess
	}
	return sess, err
}

// ClusterStatus return cluster status from AWS EMR Service.
// TODO:
// * Support multiple AWS Profiles.
func EmrClusterStatus(params map[string]interface{}, getemr func(*session.Session) emriface.EMRAPI) (string, error) {
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
	if len(ClusterId) < 1 {
		return "", ErrEmptyClusterId
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
	sess, err := GetAwsSession(params)
	if err != nil {
		return "", err
	}
	svc := getemr(sess)
	clusterInput := &emr.DescribeClusterInput{
		ClusterId: aws.String(ClusterId),
	}
	cl, err := svc.DescribeClusterWithContext(clusterCtx, clusterInput)

	if err != nil {
		if strings.Contains(err.Error(), "InvalidRequestException: Cluster id") && strings.Contains(err.Error(), "is not valid") {
			return "", fmt.Errorf("%w '%v'", ErrClusterIdIsNotValid, ClusterId)
		}
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

// EmrClusterTerminate return true if AWS EMR was terminated.
// TODO:
// * Support multiple AWS Profiles.
func EmrClusterTerminate(params map[string]interface{}, getemr func(*session.Session) emriface.EMRAPI) (bool, error) {
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
	if len(ClusterId) < 1 {
		return false, ErrEmptyClusterId
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
	sess, err := GetAwsSession(params)
	if err != nil {
		return false, err
	}
	svc := getemr(sess)
	clusterInput := &emr.TerminateJobFlowsInput{
		JobFlowIds: []*string{&ClusterId},
	}
	_, err = svc.TerminateJobFlowsWithContext(clusterCtx, clusterInput)

	if err != nil {
		if strings.Contains(err.Error(), "InvalidRequestException: Cluster id") && strings.Contains(err.Error(), "is not valid") {
			return false, fmt.Errorf("%w '%v'", ErrClusterIdIsNotValid, ClusterId)
		}
		return false, err
	}

	return true, nil
}
