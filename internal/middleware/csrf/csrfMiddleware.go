package csrf

import (
	"net/http"
	"redGlow/internal/config"
	"redGlow/internal/httpError"
	"redGlow/internal/model"
	"redGlow/internal/service"
	"redGlow/internal/tools"
	"slices"

	"go.uber.org/zap"
)

type csrfMiddleware struct {
	service service.AuthService
	logger *zap.Logger
	cfg *config.Config
}

func NewCsrfMiddleware(service service.AuthService, logger *zap.Logger, cfg *config.Config) *csrfMiddleware{
	return &csrfMiddleware{
		service: service,
		logger: logger,
		cfg: cfg,
	}
}

func checkCsrfToken(csrftokenProvided string, userSession *model.UserSession) bool{
	return tools.HashString(csrftokenProvided) == userSession.HashedCsrfToken 
}

func middlewareFunc(cs *csrfMiddleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
			requestMethod := r.Method
			HTTPmethodArray := []string{"POST","PUT","PATCH","DELETE"}
			URLPathArray := []string{"/user/session/","users"}

			if slices.Contains(URLPathArray, r.URL.Path) && r.Method == "POST"{
				next.ServeHTTP(w,r)
				return
			}

			if !slices.Contains(HTTPmethodArray,requestMethod) {
				next.ServeHTTP(w,r)
				return
			}
			csrftokenProvided := r.Header.Get(cs.cfg.AuthSettings.CSRFTokenHeaderName)

			if csrftokenProvided == "" {
				tools.HandleErrors(w, httpError.NewForbiddenError("CSRF Token wasn't provided"),cs.logger)
				return
			}
			userSession, ok := r.Context().Value(cs.cfg.AuthSettings.UserSessionContextKey).(*model.UserSession)
			if !ok {
				tools.HandleErrors(w, httpError.NewForbiddenError("CSRF Token wasn't provided"),cs.logger)
				return
			}

			if !checkCsrfToken(csrftokenProvided, userSession) {
				tools.HandleErrors(w, httpError.NewForbiddenError("CSRF Token is incorrect"),cs.logger)
				return
			}
			next.ServeHTTP(w,r)
		})
	}
}

func (cs *csrfMiddleware) GetMiddlewareFunc() func(http.Handler) http.Handler{
	return middlewareFunc(cs)
}

func (cs *csrfMiddleware) Priority() int{
	return 4
}