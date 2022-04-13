package service

import (
	"errors"
	"github.com/catinello/base62"
	log "github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
	v1 "github.com/url-shortener/internal/api/v1"
	"github.com/url-shortener/internal/store"
	"io/ioutil"
	"net/http"
)

var defaultResponse = v1.UrlShortenResponse{
	Id:           nil,
	OriginalUrl:  nil,
	ExpiryDate:   nil,
	ShortenedUrl: nil,
	CreatedDate:  nil,
	RequesterId:  nil,
	Content:      nil,
}

type UrlShortenerService struct {
	URLShortenerStore URLShortenerStore
}

func NewService(store URLShortenerStore) *UrlShortenerService {
	return &UrlShortenerService{
		URLShortenerStore: store,
	}
}

func (urls *UrlShortenerService) GetUrlByShortenedUrlName(shortenedUrlName *string) (*v1.UrlShortenResponse, error) {
	response, err := urls.URLShortenerStore.GetShortenedURL(shortenedUrlName)
	if err != nil {
		return nil, err
	}
	if response == nil || response.OriginalUrl == nil {
		return &defaultResponse, nil
	}
	urlShortResponse, err := getUrlShortenResponse(response)
	if err != nil {
		return nil, err
	}
	return urlShortResponse, nil
}

func GetContentFromOriginalUrl(originalUrl *string) (*string, error) {
	if originalUrl != nil {
		resp, err := http.Get(*originalUrl)
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		sb := string(body)
		return &sb, nil
	} else {
		return nil, errors.New("invalid original url")
	}
}

func (urls *UrlShortenerService) CreatedShortenedUrl(urlShortenedRequest *v1.UrlShortenRequest) (*v1.UrlShortenResponse, error) {
	param, err := getCreateURLShortParam(urlShortenedRequest)
	if err != nil {
		return nil, err
	}
	response, err := urls.URLShortenerStore.CreateShortenedURL(param)
	if err != nil {
		return nil, err
	}
	urlShortenResponse, err := getUrlShortenResponse(response)
	if err != nil {
		return nil, err
	}
	return urlShortenResponse, nil
}

func getCreateURLShortParam(urlShortenedRequest *v1.UrlShortenRequest) (*store.URlShorteningParam, error) {
	urlShorteningParam := store.URlShorteningParam{
		OriginalUrl: urlShortenedRequest.OriginalUrl,
		ExpiryDate:  urlShortenedRequest.ExpiryDate,
		RequesterId: urlShortenedRequest.RequesterId,
	}
	shortenedUrl, err := generateTinyUrlString()
	if err != nil {
		return nil, err
	}
	urlShorteningParam.ShortenedUrl = shortenedUrl
	return &urlShorteningParam, nil
}

func getUrlShortenResponse(urlShortening *store.URlShortening) (*v1.UrlShortenResponse, error) {
	urlShortenedResponse := &v1.UrlShortenResponse{
		Id:           urlShortening.Id,
		OriginalUrl:  urlShortening.OriginalUrl,
		ExpiryDate:   urlShortening.ExpiryDate,
		ShortenedUrl: urlShortening.ShortenedUrl,
		CreatedDate:  urlShortening.CreatedDate,
		RequesterId:  urlShortening.RequesterId,
	}
	content, err := GetContentFromOriginalUrl(urlShortening.OriginalUrl)
	if err != nil {
		log.Errorf("Could not connect to url ")
		return nil, errors.New("could not get data from url")
	}
	urlShortenedResponse.Content = content
	return urlShortenedResponse, nil
}

func generateTinyUrlString() (*string, error) {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("Unable to generate tiny url failed with %s\n", err)
		return nil, err
	}
	shortUrl := base62.Encode(int(id))
	return &shortUrl, nil
}
