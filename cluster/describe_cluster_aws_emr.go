package cluster

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	aws_request "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/emr/emriface"
	"github.com/weldpua2008/suprasched/metrics"

	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"
	"os"
	"strings"
	"sync"
	"time"
)

func init() {
	DescriberConstructors[ConstructorsDescriberTypeAwsEMR] = DescriberTypeSpec{
		instance:    NewDescriberEMR,
		constructor: NewDescriberEMRFromSection,
		Summary: `
DescribeEMR is an implementation of ClustersDescriber for Amazon EMR clusters.`,
		Description: `
It supports the following params:
- ` + "`ClusterId`" + ` Cluster Identificator
- ` + "`ClusterPool`" + ` To differentiate clusters by Pools
- ` + "`ClusterProfile`" + ` To differentiate clusters by Accounts.`,
	}
}

type DescribeEMR struct {
	ClustersDescriber
	awsSessions map[string]*session.Session
	mu          sync.RWMutex
	t           string
	section     string
	getEmrSvc   func(*session.Session) emriface.EMRAPI
}

// DefaultGetEMR implements EMR Api wrapper for tests.
func DefaultGetEMR(sess *session.Session) emriface.EMRAPI {
	return emr.New(sess)
}

// NewDescriberEMR prepare struct communicator for EMR
func NewDescriberEMR() ClustersDescriber {
	return &DescribeEMR{
		awsSessions: make(map[string]*session.Session),
		t:           "DescribeEMR",
		getEmrSvc:   DefaultGetEMR,
	}
}

// NewFetchClustersDefault prepare struct FetchClustersDefault
func NewDescriberEMRFromSection(section string) (ClustersDescriber, error) {
	s := make(map[string]*session.Session)
	// log.Warningf("NewDescriberEMRFromSection %v", section)
	return &DescribeEMR{
		awsSessions: s,
		t:           "DescribeEMR",
		section:     section,
		getEmrSvc:   DefaultGetEMR,
	}, nil

}

func (c *DescribeEMR) getCachedAwsSession(key string) (*session.Session, error) {

	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.awsSessions != nil {
		if val, ok := c.awsSessions[key]; ok {
			return val, nil
		}
	}
	return nil, fmt.Errorf("Session %v is not in cache", key)
}

func (c *DescribeEMR) SupportedClusters() []*model.Cluster {
	def := []string{ConstructorsDescriberTypeAwsEMR}
	cfgSection := fmt.Sprintf("%v.%v", c.section, config.CFG_PREFIX_CLUSTER_SUPPORTED_TYPES)
	clusterTypes := config.GetGetStringSliceDefault(cfgSection, def)

	// log.Infof("GetGetStringSliceDefault %v cfg_section %v: %v", cluster_types, cfg_section, config.ClusterRegistry.Filter(cluster_types))

	return config.ClusterRegistry.Filter(clusterTypes)
}

// getAwsSession
func (c *DescribeEMR) getAwsSession(params map[string]interface{}) (*session.Session, error) {
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
	sessionKey := fmt.Sprintf("%v%v", Profile, Region)

	if val, err := c.getCachedAwsSession(sessionKey); err == nil {
		return val, nil
	}
	// Creating & adding the session to the cache
	c.mu.Lock()
	defer c.mu.Unlock()
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

	if c.awsSessions == nil {
		c.awsSessions = make(map[string]*session.Session)
	}
	if err == nil {
		sess.Handlers.Send.PushFront(func(r *aws_request.Request) {
			// Log every request made and its payload
			metrics.ApiCallsStatistics.WithLabelValues(
				"aws",
				fmt.Sprintf("%v.%v", Profile, Region),
				"emr",
				strings.ToLower(r.Operation.Name),
			).Inc()
		})
		c.awsSessions[sessionKey] = sess
	}
	return sess, err
}

// ClusterStatus return cluster status from AWS EMR Service.
// TODO:
// * Support multiple AWS Profiles.
func (c *DescribeEMR) ClusterStatus(params map[string]interface{}) (string, error) {
	var ClusterId string
	var ctx context.Context
	var clusterCtx context.Context
	var cancel context.CancelFunc
	ttr := 30

	for _, k := range []string{"ClusterId", "Clusterid", "clusterID", "ClusterID", "clusterId",
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
	sess, err := c.getAwsSession(params)
	if err != nil {
		return "", err
	}
	svc := c.getEmrSvc(sess)
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

// DescribeClusterRequest return cluster Request from AWS EMR Service.
func (c *DescribeEMR) DescribeClusterRequest(params map[string]interface{}) (out *emr.DescribeClusterOutput, err error) {
	var ClusterId string

	for _, k := range []string{"ClusterId", "Clusterid", "clusterID", "ClusterID", "clusterId",
		"clusterid", "JobFlowID", "JobFlowId", "JobflowID", "jobFlowId"} {
		if _, ok := params[k]; ok {
			ClusterId = params[k].(string)
			break
		}
	}
	if len(ClusterId) < 1 {
		return out, fmt.Errorf("ClusterID is empty")
	}
	sess, err := c.getAwsSession(params)
	if err != nil {
		return out, err
	}
	svc := emr.New(sess)

	clusterInput := &emr.DescribeClusterInput{
		ClusterId: aws.String(ClusterId),
	}
	req, resp := svc.DescribeClusterRequest(clusterInput)
	err = req.Send()
	out = resp
	return out, err
}
