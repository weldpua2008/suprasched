package core

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

// Trace keeps track of a set of "steps" and allows us to log a specific
// step if it took longer than its share of the total allowed time
type Trace struct {
	logger      *logrus.Entry
	name        string
	threshold   *time.Duration
	startTime   time.Time
	endTime     *time.Time
	traceSteps  []traceStep
	parentTrace *Trace
}

type traceStep struct {
	stepTime time.Time
	msg      string
	msgLog   string
}

// if the trace is incomplete, don't assume an end time
func (t *Trace) time() time.Time {
	if t.endTime != nil {
		return *t.endTime
	}
	return t.startTime
}

// New creates a Trace with the specified name. The name identifies the operation to be traced. The
// Fields add key value pairs to provide additional details about the trace, such as operation inputs.
func NewTracer(ctx context.Context, name string, l *logrus.Entry) *Trace {
	logger := LoggerFromContext(ctx, l)
	return &Trace{name: name, startTime: time.Now(), logger: logger}
}

// Step adds a new step with a specific message. Call this at the end of an execution step to record
// how long it took. The Fields add key value pairs to provide additional details about the trace
// step.
func (t *Trace) Step(msg string) {
	if t.traceSteps == nil {
		// do this to avoid more than a single allocation
		t.traceSteps = make([]traceStep, 0, 7)
	}
	msgLog, _ := t.logger.String()
	t.traceSteps = append(t.traceSteps, traceStep{stepTime: time.Now(), msg: msg, msgLog: msgLog})
	t.logger.Trace(msg)
}
