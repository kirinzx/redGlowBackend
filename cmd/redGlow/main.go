package main

import (
	"context"
	"net/http"
	"redGlow/internal/config"
	"redGlow/internal/database"
	"redGlow/internal/handler"
	authHandler "redGlow/internal/handler/auth"
	"redGlow/internal/httpserver"
	"redGlow/internal/middleware"
	"redGlow/internal/repository"
	authRepository "redGlow/internal/repository/auth"
	"redGlow/internal/router"
	"redGlow/internal/service"
	authService "redGlow/internal/service/auth"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)


func main(){
    fx.New(
        
        fx.WithLogger(func (logger *zap.Logger) fxevent.Logger  {
            return &fxevent.ZapLogger{Logger: logger}
        }),
        
        fx.Provide(
            config.NewConfig,
            zap.NewProduction,
            context.Background,
            database.NewRedisDB,
            database.NewPostgresDB,
            fx.Annotate(
                authService.NewAuthService,
                fx.As(new(service.AuthService)),
            ),
            fx.Annotate(
                authRepository.NewAuthRepository,
                fx.As(new(repository.AuthRepository)),
            ),
            AsMiddleware(middleware.NewLoggerMiddleware),
            AsMiddleware(middleware.NewHeaderMiddleware),
            AsMiddleware(middleware.NewAuthMiddleware),
            AsHandler(authHandler.NewLogInHandler),  
            AsHandler(authHandler.NewLogOutHandler),
            fx.Annotate(
                router.NewChiRouter,
                fx.ParamTags(`group:"handlers"`, `group:"middlewares"`),
            ),

            httpserver.NewHTTPServer,
        ),
        fx.Invoke(func(*http.Server) {}),
    ).Run()
}

func AsHandler(f any) any {
  return fx.Annotate(
    f,
    fx.As(new(handler.Handler)),
    fx.ResultTags(`group:"handlers"`),
  )
}

func AsMiddleware(f any) any{
    return fx.Annotate(
        f,
        fx.As(new(middleware.Middleware)),
        fx.ResultTags(`group:"middlewares"`),
    )
}