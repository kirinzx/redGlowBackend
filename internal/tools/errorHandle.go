package tools

import (
	"encoding/json"
	"net/http"
	"redGlow/internal/httpError"
	"strings"

	"go.uber.org/zap"
)

type errorResponse struct {
	Errors []string `json:"errors"`
}

func HandleErrors(w http.ResponseWriter, err httpError.HTTPError, logger *zap.Logger) {
	w.WriteHeader(err.Status())
	var response []byte
	if 500 <= err.Status() && err.Status() < 600 {
		logger.Error(err.Error())
		response = []byte(`{"error":"Whoops. Something went wrong"}`)
	} else {
		errResponse := errorResponse{
			Errors: strings.Split(err.Error(), ";"),
		}
		response, _ = json.Marshal(errResponse)
	}
	w.Write(response)
}