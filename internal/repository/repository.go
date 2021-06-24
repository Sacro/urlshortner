package repository

import (
	"github.com/Sacro/urlshortner/internal/store"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	Logger  *logrus.Logger
	Store   store.Store
	Codegen func() string
}

func NewRepository(logger *logrus.Logger, store store.Store, codeGen func() string) Repository {
	return Repository{
		Logger:  logger,
		Store:   store,
		Codegen: codeGen,
	}
}
