package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.etcd.io/bbolt"
)

const testfile = "test.db"

type BoltStoreTestSuite struct {
	suite.Suite
	db    bbolt.DB
	store store
}

// SetupTest creates the DB and sets up the store
func (suite *BoltStoreTestSuite) SetupTest() {
	db, err := bbolt.Open(testfile, 0600, nil)
	if err != nil {
		suite.FailNow("Unable to open database")
	}

	suite.db = *db

	store, err := NewBoltStore(*db, "test")
	if err != nil {
		suite.FailNow("Unable to create store")
	}

	suite.store = store
}

// TearDownTest closes the DB and removes the file
func (suite *BoltStoreTestSuite) TearDownTest() {
	if err := suite.db.Close(); err != nil {
		suite.FailNow("Unable to close database")
	}

	if err := os.Remove(testfile); err != nil {
		suite.FailNowf("Unable to remove %s", testfile)
	}
}

func (suite *BoltStoreTestSuite) TestRetrevialFound() {
	code := "SHORTCODE"
	url := "http://www.example.com"

	err := suite.store.InsertURL(code, url)
	suite.Nil(err)

	u, err := suite.store.RetrieveURL(code)
	suite.Nil(err)
	suite.Equal(url, u)
}

func (suite *BoltStoreTestSuite) TestRetrievelNotFound() {
	code := "SHORTCODE"

	url, err := suite.store.RetrieveURL(code)
	suite.NotNil(err)
	suite.ErrorIs(err, ErrNotFound)
	suite.Equal("", url)
}

func TestBoltStoreTestSuite(t *testing.T) {
	suite.Run(t, new(BoltStoreTestSuite))
}
