package logger

import (
	"net/http"

	"github.com/leosunmo/zapchi"
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
	return zapchi.Logger(lm.logger,"")
}

func (lm *loggerMiddleware) Priority() int {
	return 1
}