package zetaclient

import (
	"fmt"
	. "gopkg.in/check.v1"
	"regexp"
)

type BcastSuite struct {
	bridge *ZetaCoreBridge
}

var _ = Suite(&BcastSuite{})

func (s *BcastSuite) SetUpTest(c *C) {
	fmt.Println("hello")
}

func (s *BcastSuite) TestParsingSeqNumMismatch(c *C) {
	err_msg := "fail to broadcast to zetacore,code:32, log:account sequence mismatch, expected 386232, got 386230: incorrect account sequence"
	re := regexp.MustCompile(`account sequence mismatch, expected ([0-9]*), got ([0-9]*)`)
	fmt.Printf("%q\n", re.FindStringSubmatch(err_msg))
	err_msg2 := "hahah"
	fmt.Printf("%q\n", re.FindStringSubmatch(err_msg2))
}
