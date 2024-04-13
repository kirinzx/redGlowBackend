package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"redGlow/internal/config"
	"redGlow/internal/customsmtp"
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
	cs *customsmtp.CustomSMTP
}

func NewAuthService(repo repository.AuthRepository, cfg *config.Config, cs *customsmtp.CustomSMTP) *authService{
	return &authService{
		repo: repo,
		cfg: cfg,
		cs: cs,
	}
}


func (service *authService) SignUp(ctx context.Context, userSignUp *model.UserSignUp) httpError.HTTPError{
	_, err := service.repo.CreateUser(ctx, userSignUp)
	if err != nil {
		return err
	}
	
	code := service.generateRandomString()
	conf := &model.Confirmation{
		Email: userSignUp.Email,
		Code: code,
	}
	err = service.saveCode(ctx, conf)
	if err != nil {
		return err
	}
	link := service.generateLink(conf, service.cfg.AuthSettings.FrontSignUp)
	go service.sendMail(&model.EmailConfirmation{
		Username: userSignUp.Username,
		ForWhat: "подтверждения регистрации",
		Email: userSignUp.Email, 
		Link: link,
	})
	return nil
}

func (service *authService) ConfirmSignUp(ctx context.Context, conf *model.Confirmation, w http.ResponseWriter) (*model.UserSession, httpError.HTTPError){
	hashedCode := service.hashCode(conf)
	hashConf, _ := service.repo.GetConfirmationCode(
		ctx,
		hashedCode,
	)

	if hashConf == nil {
		return nil, httpError.NewNotFoundError("Not found")
	}
	
	userSession, err := service.repo.CommitUser(ctx, conf)
	if err != nil {
		return nil, err
	}

	err = service.perfomeLogin(ctx, userSession, w)
	if err != nil {
		return nil, err
	}

	err = service.repo.DeleteConfirmationCode(ctx, hashedCode)
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (service *authService) Login(ctx context.Context, creds *model.Credentials, w http.ResponseWriter) (*model.UserGeneralInfo, httpError.HTTPError){
	userSession, err := service.repo.CheckUserByCredentials(ctx, creds)
	if err != nil {
		return nil, err
	}
	
	err = service.perfomeLogin(ctx, userSession, w)

	return &userSession.UserData, err
}

func (service *authService) ChangePassword(ctx context.Context, passwordChange *model.PasswordChange, session *model.UserSession) httpError.HTTPError{
	return service.repo.ChangeUserPassword(ctx, passwordChange, session)
}

func (service *authService) SendCodeForPasswordRecovery(ctx context.Context, confData *model.ConfirmationData) httpError.HTTPError {
	username, err := service.repo.GetUsernameByEmail(ctx, confData)
	if err != nil {
		return err
	}
	code := service.generateRandomString()
	conf := &model.Confirmation{
		Email: confData.Email,
		Code: code,
	}
	err = service.saveCode(ctx, conf)
	if err != nil{
		return err
	}
	link := service.generateLink(conf, service.cfg.AuthSettings.FrontPassRecovery)
	
	go service.sendMail(&model.EmailConfirmation{
		Username: username,
		ForWhat: "смены пароля",
		Email: confData.Email, 
		Link: link,
	})
	return nil
}

func (service *authService) ChangeForgottenPassword(ctx context.Context, passwordChange *model.PasswordChange, conf *model.Confirmation) httpError.HTTPError{
	hashedCode := service.hashCode(conf)
	hashedConf, _ := service.repo.GetConfirmationCode(ctx, hashedCode)
	if hashedConf == nil {
		return httpError.NewNotFoundError("Not found")
	}

	err := service.repo.ChangeUserForgottenPassword(ctx, passwordChange, conf)
	if err != nil {
		return err
	}
	err = service.repo.DeleteConfirmationCode(ctx, hashedCode)
	if err != nil {
		return err
	}
	return nil
}

func (service *authService) generateSessionID(session *model.UserSession) string{
	randomBytes := make([]byte, 32)

    _,_ = rand.Read(randomBytes)


	bytesAsStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("%d__%s", session.UserData.ID, bytesAsStr)
}

func (service *authService) LogOut(ctx context.Context, sessionID string) httpError.HTTPError {
	return service.repo.DeleteSession(ctx, sessionID)
}

func (service *authService) GetUserSession(ctx context.Context, sessionID string) (*model.UserSession, error) {
	return service.repo.GetSession(ctx, tools.HashString(sessionID))
}

func (service *authService) generateRandomString() string{
	randomBytes := make([]byte, 32)

    _,_ = rand.Read(randomBytes)


	bytesAsStr := hex.EncodeToString(randomBytes)
	return bytesAsStr
}

func (service *authService) generateLink(conf *model.Confirmation, uri string) string {
	params := url.Values{}
	params.Set("code",conf.Code)
	params.Set("email",conf.Email)
	u := &url.URL{
		Scheme: service.cfg.FrontServer.Scheme,
		Host: service.cfg.FrontServer.Host,
		Path: uri,
		RawQuery: params.Encode(),
	}
	return u.String()
}

func (service *authService) sendMail(emailConf *model.EmailConfirmation) error{
	messageText := fmt.Sprintf(
		"Перейдите по ссылке для %s: %s\nЕсли вы этого не делали, то просто проигнорируйте письмо",
		emailConf.Link,
		emailConf.ForWhat,
	)
	message, err := service.cs.PrepareMessage([]string{emailConf.Username, messageText})
	if err != nil {
		return err
	}
	
	err = service.cs.SendMail([]string{emailConf.Email},message)
	if err != nil {
		return err
	}
	return nil
}

func (service *authService) setCookie(unhashedSessionID, unhashedCsrfToken string, w http.ResponseWriter) {
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
}

func (service *authService) perfomeLogin(ctx context.Context, userSession *model.UserSession, w http.ResponseWriter) httpError.HTTPError{
	unhashedSessionID := service.generateSessionID(userSession)

	unhashedCsrfToken := service.generateRandomString()

	userSession.HashedSessionID = tools.HashString(unhashedSessionID)
	userSession.HashedCsrfToken = tools.HashString(unhashedCsrfToken)

	if err := service.repo.SaveSessionID(ctx, userSession, service.cfg.AuthSettings.SessionExpiration); err != nil {
		return err
	}
	service.setCookie(unhashedSessionID, unhashedCsrfToken, w)
	return nil
}

func (service *authService) saveCode(ctx context.Context, conf *model.Confirmation) httpError.HTTPError{
	hashConfirm := &model.HashedConfirmation{
		RawCode: conf.Code,
		HashedCode: service.hashCode(conf),
	}
	err := service.repo.SaveConfirmationCode(ctx, hashConfirm, service.cfg.AuthSettings.CodeExpiration)

	return err
}

func (service *authService) hashCode(conf *model.Confirmation) string {
	return tools.HashString(fmt.Sprintf(
		"%s__%s",
		conf.Code,
		conf.Email,
	))
}