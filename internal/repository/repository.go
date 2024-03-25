package repository

import (
	"context"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
	"time"
)

type AuthRepository interface{
	CreateUser(ctx context.Context, userSignUp *model.UserSignUp) (int, httpError.HTTPError)
	CommitUser(ctx context.Context, userConfirm *model.EmailConfirmation) (*model.UserSession, httpError.HTTPError)
	CheckUserByCredentials(ctx context.Context, creds *model.Credentials) (*model.UserSession, httpError.HTTPError)
	ChangeUserPassword(ctx context.Context, passChange *model.PasswordChange, session *model.UserSession) httpError.HTTPError
	ChangeUserForgottenPassword(ctx context.Context, passChange *model.PasswordChange, session *model.UserSession) httpError.HTTPError
	SaveSessionID(ctx context.Context, session *model.UserSession,expiration time.Duration) httpError.HTTPError
	GetSession(ctx context.Context, sessionID string) (*model.UserSession, error)
	DeleteSession(ctx context.Context, hashedSessionID string) httpError.HTTPError
}