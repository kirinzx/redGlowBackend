package auth

import (
	"net/http"
	"redGlow/internal/handler"
	"redGlow/internal/service"
	"redGlow/internal/tools"

	"go.uber.org/zap"
)

var _ handler.Handler = (*logOutHandler)(nil)


type logOutHandler struct {
	logger *zap.Logger
	service service.AuthService
}

func NewLogOutHandler(logger *zap.Logger, service service.AuthService) *logOutHandler{
	return &logOutHandler{
		logger: logger,
		service: service,
	}
}

func (handler *logOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userSession, err := tools.CheckIfAuthenticated(r,handler.logger)
	if err != nil {
		tools.HandleErrors(w,err,handler.logger)
		return
	}

	err = handler.service.DeleteSession(r.Context(), userSession.HashedSessionID)
	
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