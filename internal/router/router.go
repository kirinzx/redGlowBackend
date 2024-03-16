package router

import (
	"fmt"
	"redGlow/internal/handler"
	"redGlow/internal/middleware"
	"slices"

	"github.com/go-chi/chi/v5"
)

func NewChiRouter(handlers []handler.Handler, middlewares []middleware.Middleware) chi.Router {
	router := chi.NewRouter()
	slices.SortFunc(middlewares,func(mw1 middleware.Middleware, mw2 middleware.Middleware) int {
		if mw1.Priority() < mw2.Priority() {
			return -214
		}
		if mw1.Priority() > mw2.Priority() {
			return 214
		}
		return 0
	})
	for _, middleware := range middlewares{
		router.Use(middleware.GetMiddlewareFunc())
	}
	
	for _, handler := range handlers{
		router.Method(handler.HTTPMethod(),fmt.Sprintf("/%s",handler.Pattern()),handler)
	}
	
	return router
}