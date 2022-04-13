package store

import (
	"context"
	"database/sql"
	"github.com/VividCortex/mysqlerr"
	mysqlErrUtil "github.com/go-mysql/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const errBeginTransactionFailed = "Failed to create a transactional lock in the database"
const ErrFailedInsertingUrlShortening = "Failed to generate shortened url"
const ErrFailedGettingUrlShortening = "Failed to get shortened url"

type LockDBError struct {
	Err error
}

func (e LockDBError) Error() string {
	return errBeginTransactionFailed + ": " + e.Err.Error()
}

type URLShortenerStore struct {
	Database *sql.DB
}

var shortenedUrlInsertCommand = `INSERT INTO shortened_urls 
			   (original_url, shortened_url, requester_id, expiry_date)
			   VALUES (?, ?, ?, ?)`

var getShortenedUrlById = `SELECT id, original_url, shortened_url, created_date, requester_id, expiry_date from shortened_urls where id = ?`

var getShortenedUrlByShortUrlStr = `SELECT id, original_url, shortened_url, created_date, requester_id, expiry_date from shortened_urls where shortened_url = ?`

func NewStore(db *sql.DB) *URLShortenerStore {
	return &URLShortenerStore{
		Database: db,
	}
}

func isImmediatelyRetryableMySQLErr(err error) bool {
	mySQLErr := findMySQLErr(err)
	if mysqlErrUtil.MySQLErrorCode(mySQLErr) == mysqlerr.ER_LOCK_DEADLOCK {
		return true
	}
	return false
}

func findMySQLErr(err error) error {
	target := &mysql.MySQLError{}
	isTarget := errors.As(err, &target)
	if isTarget {
		return target
	}
	return nil
}

func (store *URLShortenerStore) WithTransaction(ctx context.Context, doer func(*sql.Tx) error) (err error) {

	const MaxRetries = 3
	err = store.runTransaction(ctx, doer)

	for retries := 0; isImmediatelyRetryableMySQLErr(err) && retries < MaxRetries; retries++ {
		log.Warningf("Deadlock found, retrying (%v/%v): %+v", retries, MaxRetries, err)
		err = store.runTransaction(ctx, doer)
	}

	return err
}

func (store *URLShortenerStore) runTransaction(ctx context.Context, doer func(*sql.Tx) error) error {
	tx, err := store.Database.BeginTx(ctx, nil)
	if err != nil {
		return errors.WithStack(LockDBError{Err: err})
	}

	doErr := doer(tx)

	if doErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Errorf("Failed to rollback transaction %+v", rollbackErr)
		}
		return doErr
	}

	return tx.Commit()
}

func doInsertShortenedUrl(tx *sql.Tx, param *URlShorteningParam) (*URlShortening, error) {
	result, createError := tx.Exec(shortenedUrlInsertCommand, param.OriginalUrl, param.ShortenedUrl, param.RequesterId, param.ExpiryDate)
	if createError != nil {
		return nil, createError
	}
	id, createError := result.LastInsertId()
	if createError != nil {
		return nil, createError
	}
	return getShortenedURLData(tx, id)
}

func getShortenedURLData(tx *sql.Tx, id int64) (*URlShortening, error) {
	var urlShortening = &URlShortening{}
	rows, fetchError := tx.Query(getShortenedUrlById, id)
	if fetchError != nil {
		return nil, fetchError
	}
	for rows.Next() {
		scanError := rows.Scan(&urlShortening.Id, &urlShortening.OriginalUrl, &urlShortening.ShortenedUrl, &urlShortening.CreatedDate, &urlShortening.RequesterId, &urlShortening.ExpiryDate)
		if scanError != nil {
			return nil, scanError
		}
	}
	return urlShortening, nil
}

func (store *URLShortenerStore) CreateShortenedURL(param *URlShorteningParam) (*URlShortening, error) {
	var urlShortening *URlShortening
	err := store.WithTransaction(context.Background(), func(tx *sql.Tx) (err error) {
		urlShortening, err = doInsertShortenedUrl(tx, param)
		return
	})
	return urlShortening, errors.WithMessagef(err, "%v with requester ID %v", ErrFailedInsertingUrlShortening, param.RequesterId)
}

func (store *URLShortenerStore) GetShortenedURL(shortenedURLName *string) (*URlShortening, error) {
	var urlShortening = &URlShortening{}
	rows, err := store.Database.Query(getShortenedUrlByShortUrlStr, shortenedURLName)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		scanError := rows.Scan(&urlShortening.Id, &urlShortening.OriginalUrl, &urlShortening.ShortenedUrl, &urlShortening.CreatedDate, &urlShortening.RequesterId, &urlShortening.ExpiryDate)
		if scanError != nil {
			return nil, scanError
		}
	}
	return urlShortening, errors.WithMessagef(err, "%v with shortened url %v", ErrFailedGettingUrlShortening, shortenedURLName)
}
