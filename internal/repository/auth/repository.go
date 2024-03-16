package auth

import (
	"context"
	"encoding/json"
	"errors"
	"redGlow/internal/database"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
	"redGlow/internal/repository"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

var _ repository.AuthRepository = (*authRepository)(nil)

type authRepository struct{
	PostgresDB *database.PostgresDB
	RedisDB *database.RedisDB
	logger *zap.Logger
}

func NewAuthRepository(pdb *database.PostgresDB, rdb *database.RedisDB, logger *zap.Logger) *authRepository{
	return &authRepository{
		PostgresDB: pdb,
		RedisDB: rdb,
		logger: logger,
	}
}

func (repo *authRepository) CreateUser(ctx context.Context, email, username, password string) httpError.HTTPError{
	return nil
}

func (repo *authRepository) CommitUser(ctx context.Context, email string) (*model.UserSession, httpError.HTTPError){
	return nil, nil
}

func (repo *authRepository) CheckUserByCredentials(ctx context.Context, creds *model.Credentials) (*model.UserSession, httpError.HTTPError){
	conn, err := repo.PostgresDB.Connect(ctx)
	
	if err != nil{
		return nil, httpError.NewInternalServerError(err.Error())
	}
	defer conn.Conn().Close(ctx)

	var userSession model.UserSession

	sqlQuery := "SELECT * FROM check_user_by_credits($1,$2);"
	err = conn.QueryRow(ctx, sqlQuery, creds.Email, creds.Password).Scan(&userSession.UserData.Username,
		&userSession.UserData.PhotoPath,
		&userSession.UserData.BackgroundPath,
		&userSession.UserData.SteamID,
	)
	
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil,httpError.NewBadRequestError(pgErr.Message)
		}
		return nil,httpError.NewInternalServerError(pgErr.Message)
	}
	return &userSession, nil
}

func (repo *authRepository) ChangeUserPassword(ctx context.Context, passChange *model.PasswordChange, session *model.UserSession) httpError.HTTPError{
	return nil
}

func (repo *authRepository) ChangeUserForgottenPassword(ctx context.Context, passChange *model.PasswordChange, email string) httpError.HTTPError{
	return nil
}

func (repo *authRepository) SaveSessionID(ctx context.Context, session *model.UserSession, expiration time.Duration) httpError.HTTPError{
	toSave, err := json.Marshal(session)
	if err != nil {
		return httpError.NewInternalServerError(err.Error())
	}
	err = repo.RedisDB.Set(ctx, session.HashedSessionID, toSave, expiration)
	if err != nil {
		return httpError.NewInternalServerError(err.Error())
	}
	return nil
}

func (repo *authRepository) GetSession(ctx context.Context, sessionID string) (*model.UserSession, error) {
	
	var userSession model.UserSession
	err := repo.RedisDB.Get(ctx, sessionID, &userSession)
	if err != nil {
		return nil, err
	}
	return &userSession, nil
}

func (repo *authRepository) DeleteSession(ctx context.Context, hashedSessionID string) httpError.HTTPError{ 
	err := repo.RedisDB.Del(ctx, hashedSessionID)

	if err != nil {
		return httpError.NewInternalServerError(err.Error())
	}

	return nil
}