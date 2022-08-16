package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger interface {
	Trace(format string, args ...string)
	Debug(format string, args ...string)
	Info(format string, args ...string)
	Warn(format string, args ...string)
	Error(format string, args ...string)
	Fatal(format string, args ...string)
	Panic(format string, args ...string)
}

type Formatter interface {
	Format(l *Log) ([]byte, error)
}

type Level int

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	}
	panic(fmt.Sprintf("invalid level '%d'", int(l)))
}

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

var (
	TraceOut io.Writer = os.Stdout
	DebugOut io.Writer = os.Stdout
	InfoOut  io.Writer = os.Stdout
	WarnOut  io.Writer = os.Stdout
	ErrorOut io.Writer = os.Stdout
	FatalOut io.Writer = os.Stdout
	PanicOut io.Writer = os.Stdout

	currentLevel     Level     = DebugLevel
	currentFormatter Formatter = &JSONFormatter{}
)

type Log struct {
	Time   time.Time
	Level  Level
	Format string
	Args   []interface{}
	Fields map[string]interface{}
}

func (l *Log) Msg() string {
	return fmt.Sprintf(l.Format, l.Args...)
}

func SetLevel(l Level) {
	_ = l.String() // Panic if the level does not exist
	currentLevel = l
}

func SetFormatter(f Formatter) {
	if f == nil {
		panic("invalid nil formatter")
	}
	currentFormatter = f
}

func getOutput(l Level) io.Writer {
	switch l {
	case TraceLevel:
		return TraceOut
	case DebugLevel:
		return DebugOut
	case InfoLevel:
		return InfoOut
	case WarnLevel:
		return WarnOut
	case ErrorLevel:
		return ErrorOut
	case FatalLevel:
		return FatalOut
	case PanicLevel:
		return PanicOut
	}
	panic(fmt.Sprintf("invalid level '%d'", int(l)))
}

func isLevel(expected, actual Level) bool {
	return actual >= expected
}

func log(level Level, fields map[string]interface{}, format string, args ...interface{}) {
	if isLevel(currentLevel, level) {
		// Select output by level.
		out := getOutput(level)

		// Create message.
		msg, err := currentFormatter.Format(
			&Log{Level: level, Format: format, Args: args, Fields: fields, Time: time.Now().UTC()},
		)
		if err != nil {
			panic(fmt.Errorf("the logger cannot create the message :: %v", err))
		}

		// Write the message.
		_, err = out.Write([]byte(msg))
		if err != nil {
			panic(fmt.Errorf("cannot write in %s logger :: %v", level, err))
		}
	}
}

func With(fieldName string, arg interface{}) *logger {
	logger := New()
	logger.fields[fieldName] = arg
	return logger
}

func Trace(format string, args ...interface{}) {
	log(TraceLevel, nil, format, args...)
}

func Debug(format string, args ...interface{}) {
	log(DebugLevel, nil, format, args...)
}

func Info(format string, args ...interface{}) {
	log(InfoLevel, nil, format, args...)
}

func Warn(format string, args ...interface{}) {
	log(WarnLevel, nil, format, args...)
}

func Error(format string, args ...interface{}) {
	log(ErrorLevel, nil, format, args...)
}

func Fatal(format string, args ...interface{}) {
	log(FatalLevel, nil, format, args...)
}

func Panic(format string, args ...interface{}) {
	log(PanicLevel, nil, format, args...)
}
