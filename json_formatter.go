package main

import (
	"encoding/json"
)

type JSONFormatter struct {
}

func (f *JSONFormatter) Format(l *Log) ([]byte, error) {
	fields := make(map[string]interface{})
	for key, value := range l.Fields {
		fields[key] = value
	}
	fields["level"] = l.Level.String()
	fields["msg"] = l.Msg()
	return json.Marshal(fields)
}
