package gnalog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggerWithFields(t *testing.T) {
	SetFormatter(&LineFormatter{})

	l1 := With("f1", "v1")
	l2 := l1.With("f2", "v2")
	l3 := l2.With("f3", "v3")

	testCases := []struct {
		name        string
		logger      Logger
		contains    []string
		notContains []string
	}{
		{
			name:        "L1",
			logger:      l1,
			contains:    []string{"f1=\"v1\""},
			notContains: []string{"f2=\"v2\"", "f3=\"v3\""},
		},
		{
			name:        "L2",
			logger:      l2,
			contains:    []string{"f1=\"v1\"", "f2=\"v2\""},
			notContains: []string{"f3=\"v3\""},
		},
		{
			name:        "L3",
			logger:      l3,
			contains:    []string{"f1=\"v1\"", "f2=\"v2\"", "f3=\"v3\""},
			notContains: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			cleanup := mockOutputs(out)
			defer cleanup()

			tc.logger.Debug("msg")

			outMsg := out.String()
			for _, str := range tc.contains {
				require.True(t, strings.Contains(outMsg, str))
			}
			for _, str := range tc.notContains {
				require.False(t, strings.Contains(outMsg, str))
			}
		})
	}
}

func TestLoggerImplementTheInterface(t *testing.T) {
	var logger interface{} = New()
	_, ok := logger.(Logger)
	require.True(t, ok)
}
