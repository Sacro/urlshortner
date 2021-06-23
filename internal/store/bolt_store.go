package store

import (
	"errors"
	"strings"

	bolt "go.etcd.io/bbolt"
)

type boltStore struct {
	db     bolt.DB
	bucket []byte
}

// ErrNotFound is returned when a lookup has no result
var ErrNotFound = errors.New("not found")

// NewBoltStore returns a new store backed by a BoltDB file for persistence
func NewBoltStore(db bolt.DB, bucket string) (store, error) {
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

func (s boltStore) InsertURL(shortcode, url string) error {
	s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucket)
		err := b.Put([]byte(shortcode), []byte(url))
		return err
	})

	return nil
}

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
