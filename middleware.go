package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/google/uuid"
)

var loggerCtxKey = "GNALOG_CTX_LOGGER"
var requestIDKey = "RequestID"

func SetCtxKey(key string) {
	loggerCtxKey = key
}

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type Middleware struct {
	handlerToWrap http.Handler
	withRequestID bool
}

func NewMiddleware(handlerToWrap http.Handler) *Middleware {
	return &Middleware{
		handlerToWrap: handlerToWrap,
		withRequestID: false,
	}
}

func NewMiddlewareWithRequestID(handlerToWrap http.Handler) *Middleware {
	return &Middleware{
		handlerToWrap: handlerToWrap,
		withRequestID: true,
	}
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := New()

	if m.withRequestID {
		id := uuid.New()
		l = l.With(requestIDKey, id.String())
	}

	recorder := &statusRecorder{
		ResponseWriter: w,
		Status:         http.StatusOK,
	}

	defer func() {
		l.With("status", recorder.Status).
			With("method", r.Method).
			With("client_ip", r.RemoteAddr).
			With("path", r.URL.Path).Debug("request")
	}()

	ctx := context.WithValue(r.Context(), loggerCtxKey, l)
	r = r.WithContext(ctx)
	m.handlerToWrap.ServeHTTP(recorder, r)
}

func GetLogger(r *http.Request) Logger {
	ctx := r.Context()
	ctxValue := ctx.Value(loggerCtxKey)
	if ctxValue == nil {
		panic("the context does not contain any logger")
	}
	logger, ok := ctxValue.(Logger)
	if !ok {
		panic(fmt.Sprintf("invalid logger data type: %s", reflect.TypeOf(logger).Name()))
	}
	return logger
}
