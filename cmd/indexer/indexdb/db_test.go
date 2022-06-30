package indexdb

import (
	"context"
	"encoding/json"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/mattn/go-sqlite3" // this registers a sql driver
	"github.com/rs/zerolog/log"
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

type DBSuite struct {
	idb *IndexDB
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&DBSuite{})

//func (ts *DBSuite) SetUpSuite(c *C) {
//	querier, err := query.NewZetaQuerier("3.20.194.40")
//	c.Assert(err, IsNil)
//	//os.Remove("test.db")
//	db, err := sql.Open("postgres", "test.db")
//	c.Assert(err, IsNil)
//	ts.idb, err = NewIndexDB(db, querier)
//	c.Assert(err, IsNil)
//	ts.Rebuid(c)
//}
//
//func (ts *DBSuite) Rebuid(c *C) {
//	idb := ts.idb
//	err := idb.db.Ping()
//	c.Assert(err, IsNil)
//	err = idb.Rebuild()
//	c.Assert(err, IsNil)
//}
//
//func (ts *DBSuite) TestWatchEvent(c *C) {
//	err := ts.idb.processBlock(123)
//	c.Assert(err, IsNil)
//}

func (ts *DBSuite) TestExternalTx(c *C) {
	client, err := ethclient.Dial("https://speedy-nodes-nyc.moralis.io/a8e1aba2f554b36ba23e3bd0/polygon/mumbai")
	c.Assert(err, IsNil)
	TxHash := "0x4aafe96aaaae53000a78c112c2173ed92d7812fc0e8e71ca32f1d0a212275753"
	transaction, _, err := client.TransactionByHash(context.TODO(), ethcommon.HexToHash(TxHash))
	if err != nil {
		log.Error().Err(err)
	}
	receipt, err := client.TransactionReceipt(context.TODO(), ethcommon.HexToHash(TxHash))
	if err != nil {
		log.Error().Err(err)
	}
	sender, err := client.TransactionSender(context.TODO(), transaction, receipt.BlockHash, receipt.TransactionIndex)
	if err != nil {
		log.Error().Err(err)
	}
	block, err := client.BlockByHash(context.TODO(), receipt.BlockHash)
	log.Info().Msgf("TX %s %s", "GOERLI", TxHash)
	log.Info().Msgf("sender: %s", sender)
	log.Info().Msgf("to: %s", transaction.To().Hex())
	log.Info().Msgf("status: %d", receipt.Status)
	log.Info().Msgf("gas %d, gas price %d", receipt.GasUsed, transaction.GasPrice())
	log.Info().Msgf("timestamp (unix) %s", time.Unix(int64(block.Time()), 0))
	b, err := json.Marshal(receipt.Logs)
	log.Info().Msgf("logs %s", string(b))
}
