package indexdb

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // this registers a sql driver
	"github.com/zeta-chain/zetacore/cmd/indexer/query"
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

type DBSuite struct {
	idb *IndexDB
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&DBSuite{})

func (ts *DBSuite) SetUpSuite(c *C) {
	querier, err := query.NewZetaQuerier("3.20.194.40")
	c.Assert(err, IsNil)
	os.Remove("test.db")
	db, err := sql.Open("sqlite3", "test.db")
	c.Assert(err, IsNil)
	ts.idb, err = NewIndexDB(db, querier)
	c.Assert(err, IsNil)

}

func (ts *DBSuite) TestRebuid(c *C) {
	idb := ts.idb
	err := idb.Rebuild()
	c.Assert(err, IsNil)
	err = idb.db.Ping()
	c.Assert(err, IsNil)

}
