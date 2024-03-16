package auth

import (
	"context"
	"net/http"
	"redGlow/internal/config"
	"redGlow/internal/service"

	"go.uber.org/zap"
)

type authMiddleware struct {
	service service.AuthService
	logger *zap.Logger
	cfg *config.Config
}

func NewAuthMiddleware(service service.AuthService, logger *zap.Logger, cfg *config.Config) *authMiddleware{
	return &authMiddleware{
		service: service,
		logger: logger,
		cfg: cfg,
	}
}

func getContext(am *authMiddleware, r *http.Request) context.Context{
	sessionID, err := r.Cookie(am.cfg.AuthSettings.SessionCookieName)
	
	if err != nil {
		am.logger.Error(err.Error())
		return r.Context()
	}

	userSession, _ := am.service.GetUserSession(r.Context(),sessionID.Value)

	if userSession == nil {
		return r.Context()
	}
	
	return context.WithValue(r.Context(), am.cfg.AuthSettings.UserSessionContextKey, userSession)
}

func middlewareFunc(am *authMiddleware) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			newCtx := getContext(am, r)
			
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}


func (am *authMiddleware) GetMiddlewareFunc() func(http.Handler) http.Handler{
	return middlewareFunc(am)
}

func (am *authMiddleware) Priority() int{
	return 3
}