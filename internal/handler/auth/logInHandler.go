package auth

import (
	"encoding/json"
	"net/http"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
	"redGlow/internal/service"
	"redGlow/internal/tools"

	"go.uber.org/zap"
)

type logInHandler struct{
	service service.AuthService
	logger *zap.Logger
}

func NewLogInHandler(service service.AuthService, logger *zap.Logger) *logInHandler{
	return &logInHandler{
		service: service,
		logger: logger,
	}
}

func (handler *logInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var userCreds model.Credentials
	err := json.NewDecoder(r.Body).Decode(&userCreds)
	if err != nil {
		newError := httpError.NewInternalServerError(err.Error())
		tools.HandleErrors(w, newError, handler.logger)
		return
	}

	userData, err := handler.service.Login(r.Context(), &userCreds, w)
	if err != nil{
		err := err.(httpError.HTTPError)
		tools.HandleErrors(w, err, handler.logger)
		return
	}
	dataToResponse, err := json.Marshal(userData)
	if err != nil{
		newError := httpError.NewInternalServerError(err.Error())
		tools.HandleErrors(w, newError, handler.logger)
		return
	}

	w.Write(dataToResponse)
}

func (*logInHandler) Pattern() string {
 	return "user/session/"
}

func (*logInHandler) HTTPMethod() string {
	return "POST"
}