package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Sacro/urlshortner/internal/repository"
	"github.com/Sacro/urlshortner/internal/store"
	"github.com/lithammer/shortuuid/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"
)

const testfile = "test.db"

type HandlersTestSuite struct {
	suite.Suite
	db   *bbolt.DB
	repo HandlerRepository
}

// SetupTest creates the DB and sets up the store
func (suite *HandlersTestSuite) SetupTest() {
	db, err := bbolt.Open(testfile, 0600, nil)
	if err != nil {
		suite.FailNow("Unable to open database")
	}

	suite.db = db

	store, err := store.NewBoltStore(db, "test")
	if err != nil {
		suite.FailNow("Unable to create store")
	}

	suite.repo = NewHandlerRepository(repository.Repository{
		Logger:  logrus.New(),
		Store:   store,
		Codegen: shortuuid.New,
	})
}

// TearDownTest closes the DB and removes the file
func (suite *HandlersTestSuite) TearDownTest() {
	if err := suite.db.Close(); err != nil {
		suite.FailNow("Unable to close database")
	}

	if err := os.Remove(testfile); err != nil {
		suite.FailNowf("Unable to remove %s", testfile)
	}
}

func (suite *HandlersTestSuite) TestRetrevialFound() {
	// Example URL for testing
	url := "http://www.example.com"

	// Create POST request to add shortcode
	postReq, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(url))
	suite.Nil(err)

	// Set up handler
	createRecorder := httptest.NewRecorder()
	createHandler := http.HandlerFunc(suite.repo.CreateHandler)

	// Send request to handler
	createHandler.ServeHTTP(createRecorder, postReq)

	// Check record was created
	suite.Equal(http.StatusCreated, createRecorder.Code)

	// Check body for short code
	createBody, err := io.ReadAll(createRecorder.Result().Body)
	suite.Nil(err)
	suite.NotEmpty(createBody)
	code := string(createBody)

	// Create GET request to check shortcode
	getReq, err := http.NewRequest(http.MethodGet, code, nil)
	suite.Nil(err)

	// Set up handler
	retrieveRecorder := httptest.NewRecorder()
	retrieveHandler := http.HandlerFunc(suite.repo.RetrieveHandler)

	// Send request to handler
	retrieveHandler.ServeHTTP(retrieveRecorder, getReq)

	// Check response for correct location and redirect
	suite.Equal(http.StatusSeeOther, retrieveRecorder.Code)
	location, err := retrieveRecorder.Result().Location()
	suite.Nil(err)
	suite.Equal(url, location.String())
}

func (suite *HandlersTestSuite) TestRetrievelNotFound() {
	// Example shortcode for testing
	code := "SHORTCODE"

	// Create GET request
	req, err := http.NewRequest(http.MethodGet, code, nil)
	suite.Nil(err)

	// Set up handler
	rr := httptest.NewRecorder()
	retrieveHandler := http.HandlerFunc(suite.repo.RetrieveHandler)

	// Send request to handler
	retrieveHandler.ServeHTTP(rr, req)

	// Check response for not found
	suite.Equal(http.StatusNotFound, rr.Code)
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
