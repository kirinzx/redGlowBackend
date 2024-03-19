package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type loggerMiddleware struct {
	logger *zap.Logger
}

func NewLoggerMiddleware(logger *zap.Logger) *loggerMiddleware{
	return &loggerMiddleware{
		logger:logger,
	}
}

func (lm *loggerMiddleware) GetMiddlewareFunc() func(http.Handler) http.Handler{
	logger := zap.New(lm.logger.Core(), zap.AddCallerSkip(1)).Named("")
	logger.Debug("zap.logger detected for chi")
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				logger.Info("served",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", ww.Status()),
					zap.String("reqId", middleware.GetReqID(r.Context())),
					zap.String("remoteAddr", r.RemoteAddr),
					zap.String("proto", r.Proto),
					zap.Duration("latency", time.Since(t1)),
					zap.Int("size", ww.BytesWritten()))
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (lm *loggerMiddleware) Priority() int {
	return 1
}