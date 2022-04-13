package service

import "github.com/url-shortener/internal/store"

type URLShortenerStore interface {
	CreateShortenedURL(param *store.URlShorteningParam) (*store.URlShortening, error)
	GetShortenedURL(shortenedURLName *string) (*store.URlShortening, error)
}
