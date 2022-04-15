package main

import (
	"fmt"
	"strings"
	"time"
)

type LineFormatter struct {
}

func (f *LineFormatter) Format(l *Log) ([]byte, error) {
	fields := []string{}
	for key, value := range l.Fields {
		fields = append(fields, fmt.Sprintf("%s=\"%s\"", key, value))
	}
	logMsg := fmt.Sprintf(
		"Time=\"%s\" Level=\"%s\" Msg=\"%s\" %s",
		l.Time.Format(time.RFC3339),
		l.Level,
		l.Msg(),
		strings.Join(fields, " "),
	)
	return []byte(logMsg), nil
}
