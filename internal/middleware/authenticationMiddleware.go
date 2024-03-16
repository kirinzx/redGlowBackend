package middleware

import (
	"context"
	"net/http"
	"redGlow/internal/service"

	"go.uber.org/zap"
)

type authMiddleware struct {
	service service.AuthService
	logger *zap.Logger
}

func NewAuthMiddleware(service service.AuthService, logger *zap.Logger) *authMiddleware{
	return &authMiddleware{
		service: service,
		logger: logger,
	}
}

func getContext(service service.AuthService, logger *zap.Logger, r *http.Request) context.Context{
	sessionID, err := r.Cookie("sessionID")
	
	if err != nil {
		logger.Error(err.Error())
		return r.Context()
	}

	userSession, err := service.GetSession(r.Context(),sessionID.Value)

	if err != nil && userSession != nil {
		logger.Error(err.Error())
		return r.Context()
	}
	
	return context.WithValue(r.Context(), "userSession", userSession)
}

func middlewareFunc(service service.AuthService, logger *zap.Logger) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			newCtx := getContext(service, logger, r)
			
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}


func (am *authMiddleware) GetMiddlewareFunc() func(http.Handler) http.Handler{
	return middlewareFunc(am.service, am.logger)
}