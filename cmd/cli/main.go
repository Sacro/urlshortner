package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Sacro/urlshortner/internal/store"
	"github.com/lithammer/shortuuid/v3"
	bolt "go.etcd.io/bbolt"
)

const domain = `http://example.com`

func main() {
	if len(os.Args) != 3 {
		log.Println("You must pass a two arguments, create/retrieve and the URL")
		return
	}

	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		log.Fatal("Unable to create database")
	}

	s, err := store.NewBoltStore(db, "bucket")
	if err != nil {
		log.Fatal("Unable to build store")
	}

	action, key := os.Args[1], os.Args[2]

	switch action {
	case "create":
		u, err := url.Parse(key)
		if err != nil {
			log.Fatalf("Unable to parse URL: %s", u)
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			log.Fatalf("Not a valid URL: %s", u)
		}

		code := shortuuid.New()
		key := fmt.Sprintf("%s/%s", domain, code)

		if err := s.InsertURL(key, u.String()); err != nil {
			log.Fatal("Unable to insert URL")
		}

		log.Printf("code: %s", key)
		return

	case "retrieve":
		code, err := s.RetrieveURL(key)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				log.Printf("code: %s not found", key)
				return
			}

			log.Fatal(err)
		}

		log.Printf("url: %s", code)
		return

	default:
		log.Fatalf("Unknown action: %s", action)
	}
}
