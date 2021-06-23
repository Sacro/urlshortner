package main

import (
	"net/http"

	"github.com/Sacro/urlshortner/internal/handlers"
	"github.com/Sacro/urlshortner/internal/store"
	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/lithammer/shortuuid/v3"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

func main() {
	log := logrus.New()

	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		log.Fatal("Unable to create database")
	}

	s, err := store.NewBoltStore(db, "bucket")
	if err != nil {
		log.Fatal("Unable to build store")
	}

	repo := handlers.NewRepository(log, s, shortuuid.New)

	r := chi.NewRouter()
	r.Use(logger.Logger("router", log))

	r.Post("/", repo.CreateHandler)
	r.Get("/{key}", repo.RetrieveHandler)

	logrus.Panic(http.ListenAndServe(":3000", r))
}
