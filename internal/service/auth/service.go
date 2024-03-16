package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"redGlow/internal/config"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
	"redGlow/internal/repository"
	"redGlow/internal/service"
	"redGlow/internal/tools"
)

var _ service.AuthService = (*authService)(nil)

type authService struct {
	repo repository.AuthRepository
	cfg *config.Config
}

func NewAuthService(repo repository.AuthRepository, cfg *config.Config) *authService{
	return &authService{
		repo: repo,
		cfg: cfg,
	}
}


func (service *authService) SignUp(ctx context.Context, userSignUp *model.UserSignUp) httpError.HTTPError{
	return nil
}

func (service *authService) ConfirmSignUp(ctx context.Context, confirmation *model.Confirmation) (*model.UserSession, httpError.HTTPError){
	return nil,nil
}

func (service *authService) Login(ctx context.Context, creds *model.Credentials, w http.ResponseWriter) (*model.UserGeneralInfo, httpError.HTTPError){
	userSession, err := service.repo.CheckUserByCredentials(ctx, creds)
	if err != nil {
		return nil, err
	}
	unhashedSessionID, err := service.GenerateSessionID(userSession)
	if err != nil {
		return nil, err
	}
	unhashedCsrfToken, err := service.GenerateCsrfToken()
	if err != nil {
		return nil, err
	}

	userSession.HashedSessionID = tools.HashString(unhashedSessionID)
	userSession.HashedCsrfToken = tools.HashString(unhashedCsrfToken)

	if err := service.repo.SaveSessionID(ctx, userSession, service.cfg.AuthSettings.SessionExpiration); err != nil {
		return nil, err
	}

	sessionIDCookie := &http.Cookie{
		Name: service.cfg.AuthSettings.SessionCookieName,
		Value: unhashedSessionID,
		MaxAge: int(service.cfg.AuthSettings.SessionExpiration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w,sessionIDCookie)

	csrfTokenCookie := &http.Cookie{
		Name: service.cfg.AuthSettings.CSRFTokenCookiename,
		Value: unhashedCsrfToken,
		MaxAge: int(service.cfg.AuthSettings.SessionExpiration.Seconds()),
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w,csrfTokenCookie)

	return &userSession.UserData, nil
}

func (service *authService) ChangePassword(ctx context.Context, passwordChange *model.PasswordChange, session *model.UserSession) httpError.HTTPError{
	return nil
}

func (service *authService) ChangeForgottenPassword(ctx context.Context, passwordChange *model.PasswordChange, email string) httpError.HTTPError{
	return nil
}

func (service *authService) GenerateSessionID(session *model.UserSession) (string, httpError.HTTPError){
	randomBytes := make([]byte, 32)

    _,err := rand.Read(randomBytes)
	if err != nil{
		return "", httpError.NewInternalServerError(err.Error())
	}

	bytesAsStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("%d%s%d", session.UserData.ID, bytesAsStr, session.UserMetaData.ID), nil
}

func (service *authService) LogOut(ctx context.Context, sessionID string) httpError.HTTPError {
	return service.repo.DeleteSession(ctx, sessionID)
}

func (service *authService) GetUserSession(ctx context.Context, sessionID string) (*model.UserSession, error) {
	return service.repo.GetSession(ctx, tools.HashString(sessionID))
}

func (service *authService) GenerateCsrfToken() (string, httpError.HTTPError){
	randomBytes := make([]byte, 32)

    _,err := rand.Read(randomBytes)
	if err != nil{
		return "", httpError.NewInternalServerError(err.Error())
	}

	bytesAsStr := hex.EncodeToString(randomBytes)
	return bytesAsStr, nil
}