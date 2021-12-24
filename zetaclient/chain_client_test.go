package zetaclient

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	. "gopkg.in/check.v1"
	"math/big"
	"strings"
)

type ChainClientSuite struct {
}

var _ = Suite(&ChainClientSuite{})

func (s *ChainClientSuite) SetUpTest(c *C) {

}

func (s *ChainClientSuite) TestPolygonClient(c *C) {
	client, err := ethclient.Dial(config.POLY_ENDPOINT)
	c.Assert(err, IsNil)
	bn, err := client.BlockNumber(context.TODO())
	c.Assert(err, IsNil)
	c.Logf("blocknum %d", bn)

	gas, err := client.SuggestGasPrice(context.TODO())
	c.Assert(err, IsNil)
	c.Logf("gas price %d", gas)

	receipt, err := client.TransactionReceipt(context.TODO(), ethcommon.HexToHash("0xa8ab7e7242ee1b00c7e4de581d9c87b2465bae76115bce086e7ff0e8d6a7e1ef"))
	c.Assert(err, IsNil)
	c.Log(receipt.Status, receipt.PostState, receipt.GasUsed, receipt.Logs[0], receipt.BlockNumber)

	// non-existent txhash
	_, _, err = client.TransactionByHash(context.TODO(), ethcommon.HexToHash("0x2c5d00aa638f04e49eb5d86499d1ea25d6ed1c62279008b93963d09b70fba270"))
	c.Assert(err, NotNil)

	_, err = client.TransactionReceipt(context.TODO(), ethcommon.HexToHash("0x2c5d00aa638f04e49eb5d86499d1ea25d6ed1c62279008b93963d09b70fba270"))
	c.Assert(err, NotNil)

	x := ethcommon.HexToHash("0x33")
	c.Log(x)
}

func (s *ChainClientSuite) TestBSCClient(c *C) {
	client, err := ethclient.Dial(config.BSC_ENDPOINT)
	c.Assert(err, IsNil)
	bn, err := client.BlockNumber(context.TODO())
	c.Assert(err, IsNil)
	c.Logf("blocknum %d", bn)

	gas, err := client.SuggestGasPrice(context.TODO())
	c.Assert(err, IsNil)
	c.Logf("gas price %d", gas)

	receipt, err := client.TransactionReceipt(context.TODO(), ethcommon.HexToHash("0x63326995eb00cc49df7d2aa249ec473dc351cea30f230001c4d310d6e6763490"))
	c.Assert(err, IsNil)
	c.Log(receipt.Status, receipt.PostState, receipt.GasUsed, receipt.Logs[0], receipt.BlockNumber)
}

func (s *ChainClientSuite) TestGoerliClient(c *C) {
	client, err := ethclient.Dial(config.ETH_ENDPOINT)
	c.Assert(err, IsNil)
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.ETH_ZETALOCK_ADDRESS)},
		FromBlock: big.NewInt(0).SetUint64(6013558), // LastBlock has been processed;
		ToBlock:   big.NewInt(0).SetUint64(6013558),
	}
	logs, err := client.FilterLogs(context.Background(), query)
	c.Assert(err, IsNil)
	contractAbi, err := abi.JSON(strings.NewReader(config.ETH_ZETALOCK_ABI))
	c.Assert(err, IsNil)
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logLockSendSignatureHash.Hex():
			returnVal, err := contractAbi.Unpack("LockSend", vLog.Data)
			c.Assert(err, IsNil)
			fmt.Printf("Topic 0: %s\n", vLog.Topics[0])
			fmt.Printf("Topic 1: %s\n", ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex())
			fmt.Printf("# of data fields: %d\n", len(returnVal))
			fmt.Printf("F0: receiver? %s\n", returnVal[0].(string))
			fmt.Printf("F1: amount %d\n", returnVal[1].(*big.Int))
			fmt.Printf("F2: wanted %d\n", returnVal[2].(*big.Int))
			fmt.Printf("F3: chainid? %s\n", returnVal[3].(string))
			fmt.Printf("F4: message %s\n", string(returnVal[4].([]byte)))
		}
	}
}

func (s *ChainClientSuite) TestGoerliFilterByHash(c *C) {
	client, err := ethclient.Dial(config.ETH_ENDPOINT)
	c.Assert(err, IsNil)
	topics := make([][]ethcommon.Hash, 3)
	topics[2] = make([]ethcommon.Hash, 1)
	topics[2][0] = ethcommon.HexToHash("0xae2c0ed83822269932b4b0dc55fd677b03b75e7eed10016b155ac61b0dea21d0")
	fmt.Printf("Goerli filter logs by topic: %s\n", topics)
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.ETH_ZETALOCK_ADDRESS)},
		FromBlock: big.NewInt(0), // LastBlock has been processed;
		ToBlock:   nil,
		Topics:    topics,
	}
	logs, err := client.FilterLogs(context.Background(), query)
	c.Assert(err, IsNil)
	contractAbi, err := abi.JSON(strings.NewReader(config.ETH_ZETALOCK_ABI))
	c.Assert(err, IsNil)
	c.Assert(len(logs), Equals, 1)
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logUnlockSignatureHash.Hex():
			retval, err := contractAbi.Unpack("Unlock", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking Unlock")
				continue
			}
			c.Assert(err, IsNil)
			fmt.Printf("Topic 0: %s\n", vLog.Topics[0])
			fmt.Printf("Topic 1: %s\n", vLog.Topics[1])
			fmt.Printf("Topic 2: %s\n", vLog.Topics[2])
			fmt.Printf("data: %d\n", retval[0].(*big.Int))
			fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)
		}
	}
}

func (s *ChainClientSuite) TestPolygonFilterByHash(c *C) {
	client, err := ethclient.Dial(config.POLY_ENDPOINT)
	c.Assert(err, IsNil)
	topics := make([][]ethcommon.Hash, 3)
	topics[2] = make([]ethcommon.Hash, 1)
	topics[2][0] = ethcommon.HexToHash("0xdf79baae2602a70c4044575a7d62113f74453b75bfa72022e591dcd81f078956")
	fmt.Printf("Polygon filter logs by topic: %s\n", topics)
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.POLYGON_TOKEN_ADDRESS)},
		FromBlock: big.NewInt(0), // LastBlock has been processed;
		ToBlock:   nil,
		Topics:    topics,
	}
	logs, err := client.FilterLogs(context.Background(), query)
	c.Assert(err, IsNil)
	c.Assert(len(logs), Equals, 1)
	contractAbi, err := abi.JSON(strings.NewReader(config.NONETH_ZETA_ABI))
	c.Assert(err, IsNil)

	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logMMintedSignatureHash.Hex():
			retval, err := contractAbi.Unpack("MMinted", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking Unlock")
				continue
			}
			c.Assert(err, IsNil)
			fmt.Printf("Topic 0: %s\n", vLog.Topics[0])
			fmt.Printf("Topic 1: %s\n", vLog.Topics[1])
			fmt.Printf("Topic 2: %s\n", vLog.Topics[2])
			fmt.Printf("data: %d\n", retval[0].(*big.Int))
			fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)
		}
	}

}

func (s *ChainClientSuite) TestBSCFilterByHash(c *C) {
	client, err := ethclient.Dial(config.BSC_ENDPOINT)
	c.Assert(err, IsNil)
	topics := make([][]ethcommon.Hash, 3)
	topics[2] = make([]ethcommon.Hash, 1)
	topics[2][0] = ethcommon.HexToHash("0x2ea2d1f24b36c236487d0b6d8b1ed894b78785db0af88b8e32616ed810b69a30")
	fmt.Printf("Polygon filter logs by topic: %s\n", topics)
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.BSC_TOKEN_ADDRESS)},
		FromBlock: big.NewInt(0), // LastBlock has been processed;
		ToBlock:   nil,
		Topics:    topics,
	}
	logs, err := client.FilterLogs(context.Background(), query)
	c.Assert(err, IsNil)
	c.Assert(len(logs), Equals, 1)
	contractAbi, err := abi.JSON(strings.NewReader(config.NONETH_ZETA_ABI))
	c.Assert(err, IsNil)

	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logMMintedSignatureHash.Hex():
			retval, err := contractAbi.Unpack("MMinted", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking Unlock")
				continue
			}
			c.Assert(err, IsNil)
			fmt.Printf("Topic 0: %s\n", vLog.Topics[0])
			fmt.Printf("Topic 1: %s\n", vLog.Topics[1])
			fmt.Printf("Topic 2: %s\n", vLog.Topics[2])
			fmt.Printf("data: %d\n", retval[0].(*big.Int))
			fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)
		}
	}

}
