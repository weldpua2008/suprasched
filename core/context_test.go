package core

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLoggerFromContext(t *testing.T) {
	cases := []struct {
		gotCtx    context.Context
		gotLogger *logrus.Entry
	}{
		{
			gotCtx: context.WithValue(context.Background(), "k", "v"),
		},
	}
	for _, tc := range cases {

		val := LoggerFromContext(tc.gotCtx, tc.gotLogger)
		if val == nil {
			t.Errorf("val %v == nil ", val)
		}
	}
}
