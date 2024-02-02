package types

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalBlockNumberOrHash(t *testing.T) {
	bnh := new(BlockNumberOrHash)

	testCases := []struct {
		msg      string
		input    []byte
		malleate func()
		expPass  bool
	}{
		{
			"JSON input with block hash",
			[]byte("{\"blockHash\": \"0x579917054e325746fda5c3ee431d73d26255bc4e10b51163862368629ae19739\"}"),
			func() {
				assert.Equal(t, *bnh.BlockHash, common.HexToHash("0x579917054e325746fda5c3ee431d73d26255bc4e10b51163862368629ae19739"))
				assert.Nil(t, bnh.BlockNumber)
			},
			true,
		},
		{
			"JSON input with block number",
			[]byte("{\"blockNumber\": \"0x35\"}"),
			func() {
				assert.Equal(t, *bnh.BlockNumber, BlockNumber(0x35))
				assert.Nil(t, bnh.BlockHash)
			},
			true,
		},
		{
			"JSON input with block number latest",
			[]byte("{\"blockNumber\": \"latest\"}"),
			func() {
				assert.Equal(t, *bnh.BlockNumber, EthLatestBlockNumber)
				assert.Nil(t, bnh.BlockHash)
			},
			true,
		},
		{
			"JSON input with both block hash and block number",
			[]byte("{\"blockHash\": \"0x579917054e325746fda5c3ee431d73d26255bc4e10b51163862368629ae19739\", \"blockNumber\": \"0x35\"}"),
			func() {
			},
			false,
		},
		{
			"String input with block hash",
			[]byte("\"0x579917054e325746fda5c3ee431d73d26255bc4e10b51163862368629ae19739\""),
			func() {
				assert.Equal(t, *bnh.BlockHash, common.HexToHash("0x579917054e325746fda5c3ee431d73d26255bc4e10b51163862368629ae19739"))
				assert.Nil(t, bnh.BlockNumber)
			},
			true,
		},
		{
			"String input with block number",
			[]byte("\"0x35\""),
			func() {
				assert.Equal(t, *bnh.BlockNumber, BlockNumber(0x35))
				assert.Nil(t, bnh.BlockHash)
			},
			true,
		},
		{
			"String input with block number latest",
			[]byte("\"latest\""),
			func() {
				assert.Equal(t, *bnh.BlockNumber, EthLatestBlockNumber)
				assert.Nil(t, bnh.BlockHash)
			},
			true,
		},
		{
			"String input with block number overflow",
			[]byte("\"0xffffffffffffffffffffffffffffffffffffff\""),
			func() {
			},
			false,
		},
	}

	for _, tc := range testCases {
		fmt.Printf("Case %s", tc.msg)
		// reset input
		bnh = new(BlockNumberOrHash)
		err := bnh.UnmarshalJSON(tc.input)
		tc.malleate()
		if tc.expPass {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
