package utils

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/session"
	"sync"
)

var (
	awsSessions            = make(map[string]*session.Session)
	mu                     sync.RWMutex
	ErrEmptyClusterId      = errors.New("Cluster Id is empty")
	ErrClusterIdIsNotValid = errors.New("InvalidRequestException: Cluster id is not valid")
)
