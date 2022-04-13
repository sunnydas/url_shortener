package v1

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

func PathParamString(name string, r *http.Request) (string, error) {

	vars := mux.Vars(r)
	val, ok := vars[name]
	if !ok {
		const msgBase = "Missing url parameter: "
		msg := fmt.Sprintf(msgBase+"'%s'", name)
		return "", &HttpErrorMessage{msg}
	}

	return val, nil
}

func validateURLShorteningRequest(urlShorteningRequest *UrlShortenRequest) error {
	originalUrl := urlShorteningRequest.OriginalUrl
	if originalUrl == nil || !isUrl(originalUrl) || !IsWebsiteUp(originalUrl) {
		return errors.New("invalid original url ")
	}
	if urlShorteningRequest.ExpiryDate == nil || urlShorteningRequest.ExpiryDate.Before(time.Now()) {
		return errors.New("expiry date cannot be in the past ")
	}
	return nil
}

func isUrl(str *string) bool {
	u, err := url.ParseRequestURI(*str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func IsWebsiteUp(url *string) bool {
	client := &http.Client{}
	_, err := client.Head(*url)
	if err != nil {
		log.Error("Invalid url: %url ", url)
		return false
	}
	return true
}
