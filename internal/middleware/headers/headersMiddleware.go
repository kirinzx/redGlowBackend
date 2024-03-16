package headers

import (
	"net/http"
)

type headerMiddleware struct {}

func NewHeaderMiddleware() *headerMiddleware{
	return &headerMiddleware{}
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func (hm *headerMiddleware) middlrwareFunc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        setHeaders(w)
        next.ServeHTTP(w, r)
    })
}

func (hm *headerMiddleware) GetMiddlewareFunc() func(http.Handler) http.Handler{
	return hm.middlrwareFunc
}

func (hm *headerMiddleware) Priority() int{
	return 2
}