package gnalog

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/google/uuid"
)

type key int

const (
	loggerCtxKey key = iota
)

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type Middleware struct {
	handlerToWrap   http.Handler
	requestIDHeader string
}

func NewMiddleware(handlerToWrap http.Handler) *Middleware {
	return &Middleware{
		handlerToWrap: handlerToWrap,
	}
}

func NewMiddlewareWithRequestID(handlerToWrap http.Handler, requestIDHeader string) *Middleware {
	return &Middleware{
		handlerToWrap:   handlerToWrap,
		requestIDHeader: requestIDHeader,
	}
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := New()

	if m.requestIDHeader != "" {
		reqID := m.getRequestID(r)
		l = l.With("request_id", reqID)
		w.Header().Add(m.requestIDHeader, reqID)
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

func (m *Middleware) getRequestID(r *http.Request) string {
	reqID := r.Header.Get(m.requestIDHeader)
	if reqID == "" {
		reqID = uuid.New().String()
	}
	return reqID
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
