package tools

import (
	"fmt"
	"net/http"
	"redGlow/internal/httpError"

	"go.uber.org/zap"
)

func HandleErrors(w http.ResponseWriter, err httpError.HTTPError, logger *zap.Logger) {
	w.WriteHeader(err.Status())
	var response []byte
	if 500 <= err.Status() && err.Status() < 600 {
		logger.Error(err.Error())
		response = []byte(`{"error":"Whoops. Something went wrong"}`)
	} else {
		response = []byte(fmt.Sprintf(`{"error":"%s"}`,err.Error()))
	}
	w.Write(response)
}