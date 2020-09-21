package handlers

import (
	"testing"
)

func TestTracing(t *testing.T) {
	startTestingHandler(Trace)
	defer stopTestingHandler()
	if err := emitTestingData(make(map[string]string)); err != nil {
		t.Errorf("%v", err)
	}
}
