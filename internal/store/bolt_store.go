package store

import (
	"errors"
	"strings"

	bolt "go.etcd.io/bbolt"
)

// boltStore is an implementation of the store interface backed by bbolt
type boltStore struct {
	db     *bolt.DB
	bucket []byte
}

// ErrNotFound is returned when a lookup has no result
var ErrNotFound = errors.New("not found")

// NewBoltStore returns a new store backed by a BoltDB file for persistence
func NewBoltStore(db *bolt.DB, bucket string) (store, error) {
	// Ensure we don't have a blank bucket name
	if strings.TrimSpace(bucket) == "" {
		bucket = "bucket"
	}

	// Create the bucket if it's not already there
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	}); err != nil {
		return nil, err
	}

	return boltStore{
		db:     db,
		bucket: []byte(bucket),
	}, nil
}

// InsertURL puts the shortcode and URL into the store
func (s boltStore) InsertURL(shortcode, url string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		err := b.Put([]byte(shortcode), []byte(url))
		return err
	})
}

// RetrieveURL returns the URL for a shortcode, or an error if not found
func (s boltStore) RetrieveURL(shortcode string) (string, error) {
	var value string

	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		v := b.Get([]byte(shortcode))

		if v == nil {
			return ErrNotFound
		}

		value = string(v)
		return nil
	}); err != nil {
		return "", err
	}

	return value, nil
}
