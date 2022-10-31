package gnalog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestFormatter struct {
}

func (f *TestFormatter) Format(l *Log) ([]byte, error) {
	return []byte(fmt.Sprintf(l.Format, l.Args...)), nil
}

func mockOutputs(out io.Writer) func() {
	TraceOut = out
	DebugOut = out
	InfoOut = out
	WarnOut = out
	ErrorOut = out
	FatalOut = out
	PanicOut = out
	SetLevel(TraceLevel)
	return func() {
		TraceOut = os.Stdout
		DebugOut = os.Stdout
		InfoOut = os.Stdout
		WarnOut = os.Stdout
		ErrorOut = os.Stdout
		FatalOut = os.Stdout
		PanicOut = os.Stdout
		SetLevel(DebugLevel)
	}
}

func TestLoggerImplementsTheInterface(t *testing.T) {
	var _ Logger = New()
}

func TestLogger(t *testing.T) {
	testCases := []struct {
		name   string
		logger func(format string, args ...interface{})
		level  Level
		format string
		args   []interface{}
	}{
		{name: "Trace", logger: Trace, level: TraceLevel, format: "trace msg"},
		{
			name:   "TraceWithArgs",
			logger: Trace,
			level:  TraceLevel,
			format: "trace msg %s",
			args:   []interface{}{"with args"},
		},
		{name: "Debug", logger: Debug, level: DebugLevel, format: "debug msg"},
		{
			name:   "DebugWithArgs",
			logger: Debug,
			level:  DebugLevel,
			format: "debug msg %s",
			args:   []interface{}{"with args"},
		},
		{name: "Info", logger: Info, level: InfoLevel, format: "info msg"},
		{
			name:   "InfoWithArgs",
			logger: Info,
			level:  InfoLevel,
			format: "info msg %s",
			args:   []interface{}{"with args"},
		},
		{name: "Warn", logger: Warn, level: WarnLevel, format: "warn msg"},
		{
			name:   "WarnWithArgs",
			logger: Warn,
			level:  WarnLevel,
			format: "warn msg %s",
			args:   []interface{}{"with args"},
		},
		{name: "Error", logger: Error, level: ErrorLevel, format: "error msg"},
		{
			name:   "ErrorWithArgs",
			logger: Error,
			level:  ErrorLevel,
			format: "error msg %s",
			args:   []interface{}{"with args"},
		},
		{name: "Fatal", logger: Fatal, level: FatalLevel, format: "fatal msg"},
		{
			name:   "FatalWithArgs",
			logger: Fatal,
			level:  FatalLevel,
			format: "fatal msg %s",
			args:   []interface{}{"with args"},
		},
		{name: "Panic", logger: Panic, level: PanicLevel, format: "panic msg"},
		{
			name:   "PanicWithArgs",
			logger: Panic,
			level:  PanicLevel,
			format: "panic msg %s",
			args:   []interface{}{"with args"},
		},
	}

	f := &TestFormatter{}
	SetFormatter(f)
	SetLevel(TraceLevel)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			cleanup := mockOutputs(out)
			defer cleanup()

			tc.logger(tc.format, tc.args...)

			expected := fmt.Sprintf(tc.format, tc.args...)
			actual := out.String()
			require.Equal(t, expected, actual)
		})
	}
}

func TestDifferentOutputsByLevel(t *testing.T) {
	cleanup := mockOutputs(nil)
	defer cleanup()

	f := &TestFormatter{}
	SetFormatter(f)
	SetLevel(TraceLevel)

	traceOut := bytes.NewBuffer(nil)
	debugOut := bytes.NewBuffer(nil)
	infoOut := bytes.NewBuffer(nil)
	warnOut := bytes.NewBuffer(nil)
	errorOut := bytes.NewBuffer(nil)
	fatalOut := bytes.NewBuffer(nil)
	panicOut := bytes.NewBuffer(nil)

	TraceOut = traceOut
	DebugOut = debugOut
	InfoOut = infoOut
	WarnOut = warnOut
	ErrorOut = errorOut
	FatalOut = fatalOut
	PanicOut = panicOut

	logs := []struct {
		level  Level
		logger func(format string, args ...interface{})
		output *bytes.Buffer
		msg    string
	}{
		{level: TraceLevel, logger: Trace, output: traceOut, msg: "trace"},
		{level: DebugLevel, logger: Debug, output: debugOut, msg: "debug"},
		{level: InfoLevel, logger: Info, output: infoOut, msg: "info"},
		{level: WarnLevel, logger: Warn, output: warnOut, msg: "warn"},
		{level: ErrorLevel, logger: Error, output: errorOut, msg: "error"},
		{level: FatalLevel, logger: Fatal, output: fatalOut, msg: "fatal"},
		{level: PanicLevel, logger: Panic, output: panicOut, msg: "panic"},
	}

	for _, l := range logs {
		l.logger(l.msg)
	}

	for _, l := range logs {
		require.Equal(t, l.msg, l.output.String())
	}
}

func TestLevels(t *testing.T) {
	// Levels in order.
	levels := []Level{
		TraceLevel,
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
		PanicLevel,
	}

	for i := 0; i < len(levels); i++ {
		currentLevel := levels[i]
		t.Run(currentLevel.String(), func(t *testing.T) {
			out := bytes.NewBuffer(nil)
			cleanup := mockOutputs(out)
			defer cleanup()

			f := &TestFormatter{}
			SetFormatter(f)
			SetLevel(currentLevel)

			msg := "custom msg"
			Trace(msg)
			Debug(msg)
			Info(msg)
			Warn(msg)
			Error(msg)
			Fatal(msg)
			Panic(msg)

			logMsg := out.String()
			num := strings.Count(logMsg, msg)
			require.Equal(t, len(levels)-i, num)
		})
	}
}

func TestWithFields(t *testing.T) {
	out := bytes.NewBuffer(nil)
	cleanup := mockOutputs(out)
	defer cleanup()

	SetFormatter(&LineFormatter{})

	logger := With("data", "value")
	logger.Debug("debug msg")

	logMsg := out.String()
	require.True(t, strings.Contains(logMsg, "data=\"value\""))
	require.True(t, strings.Contains(logMsg, "\"debug msg\""))
}
