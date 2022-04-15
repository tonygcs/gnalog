package main

import (
	"encoding/json"
	"time"
)

type JSONFormatter struct {
}

func (f *JSONFormatter) Format(l *Log) ([]byte, error) {
	fields := make(map[string]interface{})
	for key, value := range l.Fields {
		fields[key] = value
	}
	fields["time"] = l.Time.Format(time.RFC3339)
	fields["level"] = l.Level.String()
	fields["msg"] = l.Msg()
	return json.Marshal(fields)
}
