package main

type logger struct {
	fields map[string]interface{}
}

func New() *logger {
	return &logger{
		fields: map[string]interface{}{},
	}
}

func (l *logger) With(fieldName string, arg interface{}) *logger {
	newLogger := New()
	newLogger.fields = make(map[string]interface{})
	for key, value := range l.fields {
		newLogger.fields[key] = value
	}
	newLogger.fields[fieldName] = arg
	return newLogger
}

func (l *logger) Trace(format string, args ...interface{}) {
	log(TraceLevel, l.fields, format, args...)
}

func (l *logger) Debug(format string, args ...interface{}) {
	log(DebugLevel, l.fields, format, args...)
}

func (l *logger) Info(format string, args ...interface{}) {
	log(InfoLevel, l.fields, format, args...)
}

func (l *logger) Warn(format string, args ...interface{}) {
	log(WarnLevel, l.fields, format, args...)
}

func (l *logger) Error(format string, args ...interface{}) {
	log(ErrorLevel, l.fields, format, args...)
}

func (l *logger) Fatal(format string, args ...interface{}) {
	log(FatalLevel, l.fields, format, args...)
}

func (l *logger) Panic(format string, args ...interface{}) {
	log(PanicLevel, l.fields, format, args...)
}
