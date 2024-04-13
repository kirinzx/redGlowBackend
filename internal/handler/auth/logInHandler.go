package auth

import (
	"encoding/json"
	"net/http"
	"redGlow/internal/config"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
	"redGlow/internal/service"
	"redGlow/internal/tools"
	"redGlow/internal/validation"

	"go.uber.org/zap"
)

type logInHandler struct{
	service service.AuthService
	logger *zap.Logger
	cfg *config.Config
	validator *validation.CustomValidator
}

func NewLogInHandler(service service.AuthService, logger *zap.Logger, cfg *config.Config, validator *validation.CustomValidator) *logInHandler{
	return &logInHandler{
		service: service,
		logger: logger,
		cfg: cfg,
		validator: validator,
	}
}

func (handler *logInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userSession, _ := tools.CheckIfAuthenticated(r,handler.logger,handler.cfg.AuthSettings.UserSessionContextKey)
	if userSession != nil {
		tools.HandleErrors(w,httpError.NewForbiddenError("You are already logged in"),handler.logger)
		return
	}
	var userCreds model.Credentials
	
	if err := json.NewDecoder(r.Body).Decode(&userCreds); err != nil {
		newError := httpError.NewInternalServerError(err.Error())
		tools.HandleErrors(w, newError, handler.logger)
		return
	}

	if err := handler.validator.Validate(&userCreds); err != nil {
		tools.HandleErrors(w, httpError.NewBadRequestError(handler.validator.MakePrettyErrors(err)), handler.logger)
		return
	}

	userData, err := handler.service.Login(r.Context(), &userCreds, w)
	if err != nil{
		tools.HandleErrors(w, err, handler.logger)
		return
	}
	dataToResponse, marshalErr := json.Marshal(userData)
	if marshalErr != nil{
		newError := httpError.NewInternalServerError(marshalErr.Error())
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