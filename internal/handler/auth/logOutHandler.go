package auth

import (
	"net/http"
	"redGlow/internal/config"
	"redGlow/internal/handler"
	"redGlow/internal/service"
	"redGlow/internal/tools"

	"go.uber.org/zap"
)

var _ handler.Handler = (*logOutHandler)(nil)


type logOutHandler struct {
	logger *zap.Logger
	service service.AuthService
	cfg *config.Config
}

func NewLogOutHandler(logger *zap.Logger, service service.AuthService, cfg *config.Config) *logOutHandler{
	return &logOutHandler{
		logger: logger,
		service: service,
		cfg: cfg,
	}
}

func (handler *logOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userSession, err := tools.CheckIfAuthenticated(r,handler.logger, handler.cfg.AuthSettings.UserSessionContextKey)
	if err != nil {
		tools.HandleErrors(w,err,handler.logger)
		return
	}

	err = handler.service.LogOut(r.Context(), userSession.HashedSessionID)
	
	if err != nil{
		tools.HandleErrors(w, err, handler.logger)
		return
	}
	w.WriteHeader(200)
}


func (handler *logOutHandler) Pattern() string {
	return "user/session/"
}

func (handler *logOutHandler) HTTPMethod() string {
	return "PUT"
}