package store

import "time"

type URlShortening struct {
	Id           *string    `db:"id"`
	OriginalUrl  *string    `db:"original_url"`
	ExpiryDate   *time.Time `db:"expiry_date"`
	ShortenedUrl *string    `db:"short_url"`
	CreatedDate  *time.Time `db:"expiry_date"`
	RequesterId  *string    `db:"requester_id"`
}

type URlShorteningParam struct {
	OriginalUrl  *string    `db:"original_url"`
	ExpiryDate   *time.Time `db:"expiry_date"`
	ShortenedUrl *string    `db:"short_url"`
	RequesterId  *string    `db:"requester_id"`
}
