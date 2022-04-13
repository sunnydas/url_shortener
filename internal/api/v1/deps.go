package v1

type UrlShortenerService interface {
	GetUrlByShortenedUrlName(shortenedUrlName *string) (*UrlShortenResponse, error)
	CreatedShortenedUrl(urlShortenedRequest *UrlShortenRequest) (*UrlShortenResponse, error)
}
