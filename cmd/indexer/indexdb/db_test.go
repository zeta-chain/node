package indexdb

import (
	_ "github.com/mattn/go-sqlite3" // this registers a sql driver
	. "gopkg.in/check.v1"
	"testing"
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
//
//func (ts *DBSuite) TestExternalTx(c *C) {
//	client, err := ethclient.Dial(config.Chains[common.MumbaiChain.String()].Endpoint)
//	c.Assert(err, IsNil)
//	TxHash := "0xae19f542968253a9c9856be7226c45fbb5ad256cf54e58c4a4bfcc000f953c38"
//	transaction, _, err := client.TransactionByHash(context.TODO(), ethcommon.HexToHash(TxHash))
//	if err != nil {
//		log.Error().Err(err)
//	}
//	c.Assert(transaction, NotNil)
//	receipt, err := client.TransactionReceipt(context.TODO(), ethcommon.HexToHash(TxHash))
//	if err != nil {
//		log.Error().Err(err)
//	}
//	c.Assert(receipt, NotNil)
//	sender, err := client.TransactionSender(context.TODO(), transaction, receipt.BlockHash, receipt.TransactionIndex)
//	if err != nil {
//		log.Error().Err(err)
//	}
//	c.Assert(sender, NotNil)
//	block, err := client.BlockByHash(context.TODO(), receipt.BlockHash)
//	log.Info().Msgf("TX %s %s", "GOERLI", TxHash)
//	log.Info().Msgf("sender: %s", sender)
//	log.Info().Msgf("to: %s", transaction.To().Hex())
//	log.Info().Msgf("status: %d", receipt.Status)
//	log.Info().Msgf("gas %d, gas price %d", receipt.GasUsed, transaction.GasPrice())
//	log.Info().Msgf("timestamp (unix) %s", time.Unix(int64(block.Time()), 0))
//	b, err := json.Marshal(receipt.Logs)
//	log.Info().Msgf("logs %s", string(b))
//}
