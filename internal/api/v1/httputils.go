package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type HttpErrorMessage struct {
	Message string `json:"message"`
}

func (e HttpErrorMessage) Error() string {
	return e.Message
}

type HttpServiceError struct {
	StatusCode int         `json:"statusCode"`
	Service    string      `json:"service,omitempty"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func (h HttpServiceError) Error() string {
	return h.Message
}

func ErrCode(w http.ResponseWriter, err error, code int) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var hse *HttpServiceError
	if errors.As(err, &hse) {
		// This is an HttpServiceError, encode it verbatim
		w.WriteHeader(hse.StatusCode)
		_ = json.NewEncoder(w).Encode(hse)
		return true
	}

	var hem *HttpErrorMessage
	if errors.As(err, &hem) {
		// This is an unsuppressed error, encode the message and set error code
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(hem)
		return true
	}

	// This is another error, possibly originating from some low-level stack and/or include sensitive information.
	// Suppress and log it. Then return a generic response.
	unixTimestamp := time.Now().Unix()
	log.Infof("Unexpected error (%d): %s", unixTimestamp, err.Error())
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(HttpErrorMessage{Message: fmt.Sprintf("Server Error at %d", unixTimestamp)})

	return true
}
