package service

import (
	"context"
	"net/http"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
)

type AuthService interface{
	SignUp(ctx context.Context, userSignUp *model.UserSignUp) httpError.HTTPError
	ConfirmSignUp(ctx context.Context, confirmation *model.Confirmation, w http.ResponseWriter) (*model.UserSession, httpError.HTTPError)
	Login(ctx context.Context, creds *model.Credentials, w http.ResponseWriter) (*model.UserGeneralInfo, httpError.HTTPError)
	ChangePassword(ctx context.Context, passwordChange *model.PasswordChange, session *model.UserSession) httpError.HTTPError
	SendCodeForPasswordRecovery(ctx context.Context, confData *model.ConfirmationData) httpError.HTTPError
	ChangeForgottenPassword(ctx context.Context, passwordChange *model.PasswordChange, conf *model.Confirmation) httpError.HTTPError
	LogOut(ctx context.Context, sessionID string) httpError.HTTPError
	GetUserSession(ctx context.Context, sessionID string) (*model.UserSession, error)
}