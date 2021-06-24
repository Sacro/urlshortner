package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sacro/urlshortner/internal/repository"
	"github.com/Sacro/urlshortner/internal/store"
)

type HandlerRepository repository.Repository

func NewHandlerRepository(repo repository.Repository) HandlerRepository {
	return HandlerRepository{
		Logger:  repo.Logger,
		Store:   repo.Store,
		Codegen: repo.Codegen,
	}
}

func (r HandlerRepository) CreateHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		r.Logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	key := r.Codegen()
	if err := r.Store.InsertURL(key, string(body)); err != nil {
		r.Logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(key)); err != nil {
		r.Logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (r HandlerRepository) RetrieveHandler(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimLeft(req.URL.Path, "/")

	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	url, err := r.Store.RetrieveURL(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		r.Logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	http.Redirect(w, req, url, http.StatusSeeOther)
}
