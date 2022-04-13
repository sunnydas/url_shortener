package v1

import (
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type UrlShortenRequest struct {
	OriginalUrl *string    `json:"originalUrl"`
	ExpiryDate  *time.Time `json:"expiryDate"`
	RequesterId *string    `json:"requesterId"`
}

type UrlShortenResponse struct {
	Id           *string    `json:"id"`
	OriginalUrl  *string    `json:"originalUrl"`
	ExpiryDate   *time.Time `json:"expiryDate"`
	ShortenedUrl *string    `json:"shortenedUrl"`
	CreatedDate  *time.Time `json:"createdDate"`
	RequesterId  *string    `json:"requesterId"`
	Content      *string    `json:"content"`
}

type UrlShortenerAPI struct {
	Router              *mux.Router
	UrlShortenerService UrlShortenerService
	contentTypeJSON     ContentTypeJSON
	json                jsonpb.Marshaler
}

type ContentTypeJSON struct{}

func (l ContentTypeJSON) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	next(rw, r)
}
