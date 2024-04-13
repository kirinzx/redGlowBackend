package auth

import (
	"context"
	"encoding/json"
	"redGlow/internal/database"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
	"redGlow/internal/repository"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var _ repository.AuthRepository = (*authRepository)(nil)

type authRepository struct{
	PostgresDB *database.PostgresDB
	RedisDB *database.RedisDB
	logger *zap.Logger
	txCfg pgx.TxOptions
}

func NewAuthRepository(pdb *database.PostgresDB, rdb *database.RedisDB, logger *zap.Logger) *authRepository{
	return &authRepository{
		PostgresDB: pdb,
		RedisDB: rdb,
		logger: logger,
		txCfg: pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted,
		},
	}
}

func (repo *authRepository) CreateUser(ctx context.Context, userSignUp *model.UserSignUp) (int, httpError.HTTPError){
	conn, err := repo.PostgresDB.Connect(ctx)
	if err != nil {
		return 0, httpError.NewInternalServerError(err.Error())
	}
	defer conn.Conn().Close(ctx)
	repo.txCfg.AccessMode = pgx.ReadWrite
	tx, err := conn.BeginTx(ctx, repo.txCfg)
	if err != nil {
		return 0, httpError.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)
	query := "SELECT * FROM create_user($1,$2,$3);"
	var inserted_id int
	err = tx.QueryRow(ctx, query, userSignUp.Username,userSignUp.Password,userSignUp.Email).Scan(&inserted_id)
	pdbErr := repo.PostgresDB.HandleErrors(err)
	return inserted_id, pdbErr
}

func (repo *authRepository) CommitUser(ctx context.Context, conf *model.Confirmation) (*model.UserSession, httpError.HTTPError){
	conn, err := repo.PostgresDB.Connect(ctx)
	
	if err != nil{
		return nil, httpError.NewInternalServerError(err.Error())
	}
	defer conn.Conn().Close(ctx)
	repo.txCfg.AccessMode = pgx.ReadWrite
	tx, err := conn.BeginTx(ctx, repo.txCfg)

	if err != nil {
		return nil, httpError.NewInternalServerError(err.Error())
	}

	defer tx.Rollback(ctx)

	var userSession model.UserSession
	
	sqlQuery := "SELECT * FROM commit_user($1);"
	err = tx.QueryRow(ctx, sqlQuery, conf.Email).Scan(
		&userSession.UserData.Username,
		&userSession.UserData.PhotoPath,
		&userSession.UserData.BackgroundPath,
		&userSession.UserData.Email,
		&userSession.UserData.SteamID,
	)
	
	pdbErr := repo.PostgresDB.HandleErrors(err)
	return &userSession, pdbErr
}

func (repo *authRepository) CheckUserByCredentials(ctx context.Context, creds *model.Credentials) (*model.UserSession, httpError.HTTPError){
	conn, err := repo.PostgresDB.Connect(ctx)
	
	if err != nil{
		return nil, httpError.NewInternalServerError(err.Error())
	}
	defer conn.Conn().Close(ctx)
	repo.txCfg.AccessMode = pgx.ReadOnly
	tx, err := conn.BeginTx(ctx, repo.txCfg)

	if err != nil {
		return nil, httpError.NewInternalServerError(err.Error())
	}

	defer tx.Rollback(ctx)

	var userSession model.UserSession
	
	sqlQuery := "SELECT * FROM check_user_by_credits($1,$2);"
	err = tx.QueryRow(ctx, sqlQuery, creds.Email, creds.Password).Scan(
		&userSession.UserData.Username,
		&userSession.UserData.PhotoPath,
		&userSession.UserData.BackgroundPath,
		&userSession.UserData.Email,
		&userSession.UserData.SteamID,
	)
	
	pdbErr := repo.PostgresDB.HandleErrors(err)
	return &userSession, pdbErr
}

func (repo *authRepository) ChangeUserPassword(ctx context.Context, passChange *model.PasswordChange, session *model.UserSession) httpError.HTTPError{
	conn, err := repo.PostgresDB.Connect(ctx)
	
	if err != nil{
		return httpError.NewInternalServerError(err.Error())
	}
	defer conn.Conn().Close(ctx)
	repo.txCfg.AccessMode = pgx.ReadWrite
	tx, err := conn.BeginTx(ctx, repo.txCfg)

	if err != nil {
		return httpError.NewInternalServerError(err.Error())
	}

	defer tx.Rollback(ctx)
	
	sqlQuery := "CALL change_password($1,$2,$3)"
	_, err = tx.Exec(ctx, sqlQuery, session.UserData.Email, &passChange.OldPassword, passChange.NewPassword)

	pdbErr := repo.PostgresDB.HandleErrors(err)
	return pdbErr
}

func (repo *authRepository) GetUsernameByEmail(ctx context.Context, confData *model.ConfirmationData) (string, httpError.HTTPError) {
	conn, err := repo.PostgresDB.Connect(ctx)
	
	if err != nil{
		return "", httpError.NewInternalServerError(err.Error())
	}
	defer conn.Conn().Close(ctx)
	repo.txCfg.AccessMode = pgx.ReadWrite
	tx, err := conn.BeginTx(ctx, repo.txCfg)

	if err != nil {
		return "", httpError.NewInternalServerError(err.Error())
	}

	defer tx.Rollback(ctx)
	
	sqlQuery := "SELECT * FROM check_email($1)"
	var username string
	err = tx.QueryRow(ctx, sqlQuery, confData.Email).Scan(&username)
	
	pdbErr := repo.PostgresDB.HandleErrors(err)
	if pdbErr != nil {
		return "", pdbErr	
	}
	return username, pdbErr
}

func (repo *authRepository) ChangeUserForgottenPassword(ctx context.Context, passChange *model.PasswordChange, userConfirm *model.Confirmation) httpError.HTTPError{
	conn, err := repo.PostgresDB.Connect(ctx)
	
	if err != nil{
		return httpError.NewInternalServerError(err.Error())
	}
	defer conn.Conn().Close(ctx)
	repo.txCfg.AccessMode = pgx.ReadWrite
	tx, err := conn.BeginTx(ctx, repo.txCfg)

	if err != nil {
		return httpError.NewInternalServerError(err.Error())
	}

	defer tx.Rollback(ctx)
	
	sqlQuery := "CALL change_forgotten_password($1,$2)"
	_, err = tx.Exec(ctx, sqlQuery, userConfirm.Email, passChange.NewPassword)

	pdbErr := repo.PostgresDB.HandleErrors(err)
	return pdbErr
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

func (repo *authRepository) GetSession(ctx context.Context, hashedSessionID string) (*model.UserSession, error) {
	
	var userSession model.UserSession
	err := repo.RedisDB.Get(ctx, hashedSessionID, &userSession)
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

func (repo *authRepository) SaveConfirmationCode(ctx context.Context, hashedConf *model.HashedConfirmation, exp time.Duration) httpError.HTTPError{
	toSave, err := json.Marshal(hashedConf)
	if err != nil {
		return httpError.NewInternalServerError(err.Error())
	}
	err = repo.RedisDB.Set(ctx, hashedConf.HashedCode, toSave, exp)
	if err != nil {
		return httpError.NewInternalServerError(err.Error())
	}
	return nil
}

func (repo *authRepository) GetConfirmationCode(ctx context.Context, hashedCode string) (*model.HashedConfirmation, httpError.HTTPError){
	var confValue model.HashedConfirmation
	err := repo.RedisDB.Get(ctx, hashedCode, &confValue)
	if err != nil {
		return nil, httpError.NewInternalServerError(err.Error())
	}
	return &confValue, nil
}

func (repo *authRepository) DeleteConfirmationCode(ctx context.Context, hashedCode string) httpError.HTTPError{
	err := repo.RedisDB.Del(ctx, hashedCode)

	if err != nil {
		return httpError.NewInternalServerError(err.Error())
	}

	return nil
}