package v1

import (
	"encoding/json"
	"errors"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
)

const ErrMsgBody = "unable to parse body"

func (a *UrlShortenerAPI) initAPIRoutes() {
	a.Router.
		Methods(http.MethodPost).
		Path("/url-shortener").
		Handler(negroni.New(
			a.contentTypeJSON,
			negroni.HandlerFunc(a.shortenUrl),
		))

	a.Router.
		Methods(http.MethodGet).
		Path("/url-shortener/{shortenedUrl}").
		Handler(negroni.New(
			a.contentTypeJSON,
			negroni.HandlerFunc(a.getShortenedUrl),
		))
}

func (a *UrlShortenerAPI) shortenUrl(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
	urlShortenRequest := UrlShortenRequest{}
	if err := json.NewDecoder(r.Body).Decode(&urlShortenRequest); err != nil {
		msg := ErrMsgBody
		log.Warn(msg, err.Error())
		ErrCode(w, &HttpErrorMessage{errors.New(msg).Error()}, http.StatusBadRequest)
		return
	}
	err := validateURLShorteningRequest(&urlShortenRequest)
	if err != nil {
		log.Error(err)
		ErrCode(w, &HttpErrorMessage{errors.New(err.Error()).Error()}, http.StatusBadRequest)
		return
	}
	urlShortenResponse, serviceError := a.UrlShortenerService.CreatedShortenedUrl(&urlShortenRequest)
	if serviceError != nil {
		log.Error(serviceError)
		ErrCode(w, &HttpErrorMessage{errors.New("could not create shortened url").Error()}, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	buf, marshalError := json.Marshal(urlShortenResponse)
	if marshalError != nil {
		ErrCode(w, &HttpErrorMessage{errors.New("bad json").Error()}, http.StatusInternalServerError)
		return
	}
	_, bufferWriteError := w.Write(buf)
	if bufferWriteError != nil {
		panic("could not write to response buffer")
	}
}

func (a *UrlShortenerAPI) getShortenedUrl(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
	shortUrlName, err := PathParamString("shortenedUrl", r)
	if err != nil {
		ErrCode(w, err, http.StatusBadRequest)
		return
	}
	urlShortenResponse, serviceError := a.UrlShortenerService.GetUrlByShortenedUrlName(&shortUrlName)
	if serviceError != nil {
		ErrCode(w, &HttpErrorMessage{errors.New(serviceError.Error()).Error()}, http.StatusInternalServerError)
		return
	}
	if urlShortenResponse.Id == nil {
		ErrCode(w, &HttpErrorMessage{errors.New("no matching original url found").Error()}, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	buf, marshalError := json.Marshal(urlShortenResponse)
	if marshalError != nil {
		ErrCode(w, &HttpErrorMessage{errors.New("bad json").Error()}, http.StatusInternalServerError)
		return
	}
	_, bufferWriteError := w.Write(buf)
	if bufferWriteError != nil {
		panic("could not write to response buffer")
	}
}

func NewAPI(router *mux.Router, urlShortenerService UrlShortenerService) *UrlShortenerAPI {
	api := UrlShortenerAPI{
		router,
		urlShortenerService,
		ContentTypeJSON{},
		jsonpb.Marshaler{EmitDefaults: true},
	}
	api.initAPIRoutes()
	return &api
}
