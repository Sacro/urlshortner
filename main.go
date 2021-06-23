package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/Sacro/urlshortner/internal/store"
	"github.com/lithammer/shortuuid/v3"
	bolt "go.etcd.io/bbolt"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("You must pass a single argument, either a URL (starting with http) or a short code")
		return
	}

	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		log.Fatal("Unable to create database")
	}

	s, err := store.NewBoltStore(*db, "bucket")
	if err != nil {
		log.Fatal("Unable to build store")
	}

	arg := os.Args[1]

	if strings.HasPrefix(arg, "http") {
		// Is a URL
		code := shortuuid.New()
		s.InsertURL(code, arg)

		log.Printf("code: %s", code)
	} else {
		// Is not a URL
		code, err := s.RetrieveURL(arg)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				log.Printf("code: %s not found", arg)
				return
			}

			log.Fatal(err)
		}

		log.Printf("url: %s", code)
	}

}
