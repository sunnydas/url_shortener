package app

import (
	"github.com/gorilla/mux"
	v1 "github.com/url-shortener/internal/api/v1"
)

type App struct {
	Router          *mux.Router
	URlShortenerAPI *v1.UrlShortenerAPI
}

func NewApp(
	r *mux.Router,
	api *v1.UrlShortenerAPI,

) *App {
	return &App{
		Router:          r,
		URlShortenerAPI: api,
	}
}
