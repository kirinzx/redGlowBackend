package service

import (
	"context"
	"net/http"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
)

type AuthService interface{
	SignUp(ctx context.Context, userSignUp *model.UserSignUp) httpError.HTTPError
	ConfirmSignUp(ctx context.Context, confirmation *model.Confirmation) (*model.UserSession, httpError.HTTPError)
	Login(ctx context.Context, creds *model.Credentials, w http.ResponseWriter) (*model.UserGeneralInfo, httpError.HTTPError)
	ChangePassword(ctx context.Context, passwordChange *model.PasswordChange, session *model.UserSession) httpError.HTTPError
	ChangeForgottenPassword(ctx context.Context, passwordChange *model.PasswordChange, email string) httpError.HTTPError
	GenerateSessionID(session *model.UserSession) (string, httpError.HTTPError)
	LogOut(ctx context.Context, sessionID string) httpError.HTTPError
	GetUserSession(ctx context.Context, sessionID string) (*model.UserSession, error)
	GenerateCsrfToken() (string, httpError.HTTPError)
}