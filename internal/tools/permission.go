package tools

import (
	"net/http"
	"redGlow/internal/httpError"
	"redGlow/internal/model"

	"go.uber.org/zap"
)

func CheckIfAuthenticated(r *http.Request, logger *zap.Logger, userSessionContextKey string) (*model.UserSession, httpError.HTTPError) {
	userSession, ok := r.Context().Value(userSessionContextKey).(*model.UserSession)

	if !ok || userSession == nil {
		err := httpError.NewForbiddenError("You are't logged in")
		return nil, err
	}
	return userSession, nil
}