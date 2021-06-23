package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Sacro/urlshortner/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	logger  *logrus.Logger
	store   store.Store
	codeGen func() string
}

func NewRepository(logger *logrus.Logger, store store.Store, codeGen func() string) Repository {
	return Repository{
		logger:  logger,
		store:   store,
		codeGen: codeGen,
	}
}

func (r Repository) CreateHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		r.logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	key := r.codeGen()
	if err := r.store.InsertURL(key, string(body)); err != nil {
		r.logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(key)); err != nil {
		r.logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (r Repository) RetrieveHandler(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "key")

	url, err := r.store.RetrieveURL(key)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		r.logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	http.Redirect(w, req, url, http.StatusSeeOther)
}
