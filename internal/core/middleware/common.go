package middleware

import (
	"Miji/internal/core/logger"
	"Miji/internal/core/transport"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const RequestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(RequestIDHeader, requestID)
			w.Header().Add(RequestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)
			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), logger.LogContextKey, l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					l := logger.FromContext(r.Context())
					transport.Error(w, r, l.Logger, panicErr(p), "panic during http request")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := logger.FromContext(r.Context())
			rw := transport.NewResponseWriter(w)

			before := time.Now()
			l.Debug(
				"handling request",
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			l.Debug(
				"request finished",
				zap.Int("status_code", rw.StatusCode()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}

func panicErr(p any) error {
	switch v := p.(type) {
	case error:
		return v
	default:
		return nil
	}
}
