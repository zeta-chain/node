package query

import (
	. "gopkg.in/check.v1"
	"testing"
)

const (
	zetaNodeIP = "3.20.194.40"
)

type TxSuite struct {
	zq *ZetaQuerier
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&TxSuite{})

//
//func (ts *TxSuite) SetUpSuite(c *C) {
//	zq, err := NewZetaQuerier(zetaNodeIP)
//	c.Assert(err, IsNil)
//	ts.zq = zq
//}
//
//func (ts *TxSuite) Test1(c *C) {
//	cnt := 0
//	total, err := ts.zq.VisitAllTxEvents("SendFinalized", 0, func(res *sdk.TxResponse) error {
//		cnt += 1
//		fmt.Printf("txhash: %s\n", res.TxHash)
//		return nil
//	})
//
//	c.Assert(err, IsNil)
//	fmt.Printf("total: %d, cnt: %d\n", total, cnt)
//}
