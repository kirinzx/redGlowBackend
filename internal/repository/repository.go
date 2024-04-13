package repository

import (
	"context"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
	"time"
)

type AuthRepository interface{
	CreateUser(ctx context.Context, userSignUp *model.UserSignUp) (int, httpError.HTTPError)
	CommitUser(ctx context.Context, conf *model.Confirmation) (*model.UserSession, httpError.HTTPError)
	CheckUserByCredentials(ctx context.Context, creds *model.Credentials) (*model.UserSession, httpError.HTTPError)
	ChangeUserPassword(ctx context.Context, passChange *model.PasswordChange, session *model.UserSession) httpError.HTTPError
	GetUsernameByEmail(ctx context.Context, confData *model.ConfirmationData) (string, httpError.HTTPError)
	ChangeUserForgottenPassword(ctx context.Context, passChange *model.PasswordChange, conf *model.Confirmation) httpError.HTTPError
	SaveSessionID(ctx context.Context, session *model.UserSession, expiration time.Duration) httpError.HTTPError
	GetSession(ctx context.Context, hashedSessionID string) (*model.UserSession, error)
	DeleteSession(ctx context.Context, hashedSessionID string) httpError.HTTPError
	SaveConfirmationCode(ctx context.Context, hashedConf *model.HashedConfirmation, exp time.Duration) httpError.HTTPError
	GetConfirmationCode(ctx context.Context, hashedCode string) (*model.HashedConfirmation, httpError.HTTPError)
	DeleteConfirmationCode(ctx context.Context, hashedCode string) httpError.HTTPError
}