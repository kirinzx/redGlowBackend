package auth

import (
	"net/http"
	"redGlow/internal/handler"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
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
	userSession, ok := r.Context().Value("userSession").(*model.UserSession)

	if !ok || userSession == nil {
		tools.HandleErrors(w, httpError.NewForbiddenError("You are't logged in"), handler.logger)
		return
	}

	err := handler.service.DeleteSession(r.Context(), userSession.HashedSessionID)
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