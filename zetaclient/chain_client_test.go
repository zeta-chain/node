package zetaclient

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/types"
	. "gopkg.in/check.v1"
	"math/big"
	"os"
	"strings"
)

type ChainClientSuite struct {
	ethOB *ChainObserver
}

var _ = Suite(&ChainClientSuite{})

func (s *ChainClientSuite) SetUpTest(c *C) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

//
//func (s *ChainClientSuite) TestPolygonClient(c *C) {
//	client, err := ethclient.Dial(config.POLY_ENDPOINT)
//	c.Assert(err, IsNil)
//	bn, err := client.BlockNumber(context.TODO())
//	c.Assert(err, IsNil)
//	c.Logf("blocknum %d", bn)
//
//	gas, err := client.SuggestGasPrice(context.TODO())
//	c.Assert(err, IsNil)
//	c.Logf("gas price %d", gas)
//
//	receipt, err := client.TransactionReceipt(context.TODO(), ethcommon.HexToHash("0xa8ab7e7242ee1b00c7e4de581d9c87b2465bae76115bce086e7ff0e8d6a7e1ef"))
//	c.Assert(err, IsNil)
//	c.Log(receipt.Status, receipt.PostState, receipt.GasUsed, receipt.Logs[0], receipt.BlockNumber)
//
//	// non-existent txhash
//	_, _, err = client.TransactionByHash(context.TODO(), ethcommon.HexToHash("0x2c5d00aa638f04e49eb5d86499d1ea25d6ed1c62279008b93963d09b70fba270"))
//	c.Assert(err, NotNil)
//
//	_, err = client.TransactionReceipt(context.TODO(), ethcommon.HexToHash("0x2c5d00aa638f04e49eb5d86499d1ea25d6ed1c62279008b93963d09b70fba270"))
//	c.Assert(err, NotNil)
//
//	x := ethcommon.HexToHash("0x33")
//	c.Log(x)
//}
//
//func (s *ChainClientSuite) TestBSCClient(c *C) {
//	client, err := ethclient.Dial(config.BSC_ENDPOINT)
//	c.Assert(err, IsNil)
//	bn, err := client.BlockNumber(context.TODO())
//	c.Assert(err, IsNil)
//	c.Logf("blocknum %d", bn)
//
//	gas, err := client.SuggestGasPrice(context.TODO())
//	c.Assert(err, IsNil)
//	c.Logf("gas price %d", gas)
//
//	receipt, err := client.TransactionReceipt(context.TODO(), ethcommon.HexToHash("0x63326995eb00cc49df7d2aa249ec473dc351cea30f230001c4d310d6e6763490"))
//	c.Assert(err, IsNil)
//	c.Log(receipt.Status, receipt.PostState, receipt.GasUsed, receipt.Logs[0], receipt.BlockNumber)
//}
//
//func (s *ChainClientSuite) TestGoerliClient(c *C) {
//	client, err := ethclient.Dial(config.ETH_ENDPOINT)
//	c.Assert(err, IsNil)
//	query := ethereum.FilterQuery{
//		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.ETH_MPI_ADDRESS)},
//		FromBlock: big.NewInt(0).SetUint64(6013558), // LastBlock has been processed;
//		ToBlock:   big.NewInt(0).SetUint64(6013558),
//	}
//	logs, err := client.FilterLogs(context.Background(), query)
//	c.Assert(err, IsNil)
//	contractAbi, err := abi.JSON(strings.NewReader(config.ETH_ZETALOCK_ABI))
//	c.Assert(err, IsNil)
//	for _, vLog := range logs {
//		switch vLog.Topics[0].Hex() {
//		case logLockSendSignatureHash.Hex():
//			returnVal, err := contractAbi.Unpack("LockSend", vLog.Data)
//			c.Assert(err, IsNil)
//			fmt.Printf("Topic 0: %s\n", vLog.Topics[0])
//			fmt.Printf("Topic 1: %s\n", ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex())
//			fmt.Printf("# of data fields: %d\n", len(returnVal))
//			fmt.Printf("F0: receiver? %s\n", returnVal[0].(string))
//			fmt.Printf("F1: amount %d\n", returnVal[1].(*big.Int))
//			fmt.Printf("F2: wanted %d\n", returnVal[2].(*big.Int))
//			fmt.Printf("F3: chainid? %s\n", returnVal[3].(string))
//			fmt.Printf("F4: message %s\n", string(returnVal[4].([]byte)))
//		}
//	}
//}
//
//func (s *ChainClientSuite) TestGoerliFilterByHash(c *C) {
//	client, err := ethclient.Dial(config.ETH_ENDPOINT)
//	c.Assert(err, IsNil)
//	topics := make([][]ethcommon.Hash, 3)
//	topics[2] = make([]ethcommon.Hash, 1)
//	topics[2][0] = ethcommon.HexToHash("0xae2c0ed83822269932b4b0dc55fd677b03b75e7eed10016b155ac61b0dea21d0")
//	fmt.Printf("Goerli filter logs by topic: %s\n", topics)
//	query := ethereum.FilterQuery{
//		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.ETH_MPI_ADDRESS)},
//		FromBlock: big.NewInt(0), // LastBlock has been processed;
//		ToBlock:   nil,
//		Topics:    topics,
//	}
//	logs, err := client.FilterLogs(context.Background(), query)
//	c.Assert(err, IsNil)
//	contractAbi, err := abi.JSON(strings.NewReader(config.ETH_ZETALOCK_ABI))
//	c.Assert(err, IsNil)
//	c.Assert(len(logs), Equals, 1)
//	for _, vLog := range logs {
//		switch vLog.Topics[0].Hex() {
//		case logUnlockSignatureHash.Hex():
//			retval, err := contractAbi.Unpack("Unlock", vLog.Data)
//			if err != nil {
//				fmt.Println("error unpacking Unlock")
//				continue
//			}
//			c.Assert(err, IsNil)
//			fmt.Printf("Topic 0: %s\n", vLog.Topics[0])
//			fmt.Printf("Topic 1: %s\n", vLog.Topics[1])
//			fmt.Printf("Topic 2: %s\n", vLog.Topics[2])
//			fmt.Printf("data: %d\n", retval[0].(*big.Int))
//			fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)
//		}
//	}
//}
//
//func (s *ChainClientSuite) TestPolygonFilterByHash(c *C) {
//	client, err := ethclient.Dial(config.POLY_ENDPOINT)
//	c.Assert(err, IsNil)
//	topics := make([][]ethcommon.Hash, 3)
//	topics[2] = make([]ethcommon.Hash, 1)
//	topics[2][0] = ethcommon.HexToHash("0xdf79baae2602a70c4044575a7d62113f74453b75bfa72022e591dcd81f078956")
//	fmt.Printf("Polygon filter logs by topic: %s\n", topics)
//	query := ethereum.FilterQuery{
//		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.POLYGON_MPI_ADDRESS)},
//		FromBlock: big.NewInt(0), // LastBlock has been processed;
//		ToBlock:   nil,
//		Topics:    topics,
//	}
//	logs, err := client.FilterLogs(context.Background(), query)
//	c.Assert(err, IsNil)
//	c.Assert(len(logs), Equals, 1)
//	contractAbi, err := abi.JSON(strings.NewReader(config.NONETH_ZETA_ABI))
//	c.Assert(err, IsNil)
//
//	for _, vLog := range logs {
//		switch vLog.Topics[0].Hex() {
//		case logMMintedSignatureHash.Hex():
//			retval, err := contractAbi.Unpack("MMinted", vLog.Data)
//			if err != nil {
//				fmt.Println("error unpacking Unlock")
//				continue
//			}
//			c.Assert(err, IsNil)
//			fmt.Printf("Topic 0: %s\n", vLog.Topics[0])
//			fmt.Printf("Topic 1: %s\n", vLog.Topics[1])
//			fmt.Printf("Topic 2: %s\n", vLog.Topics[2])
//			fmt.Printf("data: %d\n", retval[0].(*big.Int))
//			fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)
//		}
//	}
//
//}
//
//func (s *ChainClientSuite) TestBSCFilterByHash(c *C) {
//	client, err := ethclient.Dial(config.BSC_ENDPOINT)
//	c.Assert(err, IsNil)
//	topics := make([][]ethcommon.Hash, 3)
//	topics[2] = make([]ethcommon.Hash, 1)
//	topics[2][0] = ethcommon.HexToHash("0x2ea2d1f24b36c236487d0b6d8b1ed894b78785db0af88b8e32616ed810b69a30")
//	fmt.Printf("Polygon filter logs by topic: %s\n", topics)
//	query := ethereum.FilterQuery{
//		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.BSC_MPI_ADDRESS)},
//		FromBlock: big.NewInt(0), // LastBlock has been processed;
//		ToBlock:   nil,
//		Topics:    topics,
//	}
//	logs, err := client.FilterLogs(context.Background(), query)
//	c.Assert(err, IsNil)
//	c.Assert(len(logs), Equals, 1)
//	contractAbi, err := abi.JSON(strings.NewReader(config.NONETH_ZETA_ABI))
//	c.Assert(err, IsNil)
//
//	for _, vLog := range logs {
//		switch vLog.Topics[0].Hex() {
//		case logMMintedSignatureHash.Hex():
//			retval, err := contractAbi.Unpack("MMinted", vLog.Data)
//			if err != nil {
//				fmt.Println("error unpacking Unlock")
//				continue
//			}
//			c.Assert(err, IsNil)
//			fmt.Printf("Topic 0: %s\n", vLog.Topics[0])
//			fmt.Printf("Topic 1: %s\n", vLog.Topics[1])
//			fmt.Printf("Topic 2: %s\n", vLog.Topics[2])
//			fmt.Printf("data: %d\n", retval[0].(*big.Int))
//			fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)
//		}
//	}
//
//}

func (s *ChainClientSuite) TestGoerliFilterByHash(c *C) {

	//c.Assert(logMPIReceiveSignatureHash.Hex(), Equals, "bf55560bd045bd8d3823023d2adb4614181aaefea493132a99e6a33ef11cd084"
	client, err := ethclient.Dial(config.ETH_ENDPOINT)
	c.Assert(err, IsNil)

	topics[0] = []ethcommon.Hash{logMPISendSignatureHash}
	fmt.Printf("Goerli filter logs by topic: %s\n", topics)
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.Chains["ETH"].MPIContractAddress)},
		FromBlock: big.NewInt(0), // LastBlock has been processed;
		ToBlock:   nil,
		Topics:    topics,
	}

	logs, err := client.FilterLogs(context.Background(), query)
	c.Assert(err, IsNil)
	contractAbi, err := abi.JSON(strings.NewReader(config.MPI_ABI_STRING))
	c.Assert(err, IsNil)
	//c.Assert(len(logs), Equals, 1)
	var log_message string
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logMPISendSignatureHash.Hex():
			vals, err := contractAbi.Unpack("ZetaMessageSendEvent", vLog.Data)
			c.Assert(err, IsNil)

			sender := vLog.Topics[1]
			log_message = fmt.Sprintf("sender %s", ethcommon.HexToAddress(sender.Hex()).Hex())
			log.Info().Msg(log_message)
			//
			destChainID := vals[0].(uint16)
			log_message = fmt.Sprintf("destChainID %d", destChainID)
			log.Info().Msg(log_message)

			destContract := vals[1].([]byte)
			log_message = fmt.Sprintf("destContract %x", destContract)
			log.Info().Msg(log_message)

			zetaAmount := vals[2].(*big.Int)
			log_message = fmt.Sprintf("zetaAmount %d", zetaAmount)
			log.Info().Msg(log_message)

			gasLimit := vals[3].(*big.Int)
			log_message = fmt.Sprintf("gasLimit %d", gasLimit)
			log.Info().Msg(log_message)

			message := vals[4].([]byte)
			log_message = fmt.Sprintf("message %s", hex.EncodeToString(message))
			log.Info().Msg(log_message)

			zetaParams := vals[5].([]byte)
			log_message = fmt.Sprintf("zetaParams %s", hex.EncodeToString(zetaParams[:]))
			log.Info().Msg(log_message)

		}
	}
}

func (s *ChainClientSuite) TestGoerliSendHash(c *C) {
	log.Debug().Msg("TestingGoerliSendHash")
	client, err := ethclient.Dial(config.ETH_ENDPOINT)
	c.Assert(err, IsNil)

	recvTopics := make([][]ethcommon.Hash, 4)
	recvTopics[3] = []ethcommon.Hash{ethcommon.HexToHash("0xcccb58610b0b65d5b1d8e5f16435254a787e324209d9b3877a8ece68859a0f55")}
	recvTopics[0] = []ethcommon.Hash{logMPIReceiveSignatureHash}
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.Chains["ETH"].MPIContractAddress)},
		FromBlock: big.NewInt(0), // LastBlock has been processed;
		ToBlock:   nil,
		Topics:    recvTopics,
	}
	logs, err := client.FilterLogs(context.Background(), query)
	c.Assert(err, IsNil)
	contractAbi, err := abi.JSON(strings.NewReader(config.MPI_ABI_STRING))
	c.Assert(err, IsNil)
	var log_message string
	for i, vLog := range logs {
		log.Info().Msgf("log %d", i)

		vals, err := contractAbi.Unpack("ZetaMessageReceiveEvent", vLog.Data)
		c.Assert(err, IsNil)

		srcChainID := vLog.Topics[1]
		log_message = fmt.Sprintf("srcChainID %d", srcChainID.Big())
		log.Info().Msg(log_message)
		//
		destContract := vLog.Topics[2]
		log_message = fmt.Sprintf("destContra %s", types.HashToAddress(destContract))
		log.Info().Msg(log_message)

		sendHash := vLog.Topics[3]
		log_message = fmt.Sprintf("sendHash %s", sendHash)
		log.Info().Msg(log_message)

		sender := vals[0].([]byte)
		log_message = fmt.Sprintf("sender %s", ethcommon.BytesToAddress(sender))
		log.Info().Msg(log_message)

		zetaAmount := vals[1].(*big.Int)
		log_message = fmt.Sprintf("zeta %d", zetaAmount)
		log.Info().Msg(log_message)

		message := vals[2].([]byte)
		log_message = fmt.Sprintf("message  %d", message)
		log.Info().Msg(log_message)
	}
}

func (s *ChainClientSuite) TestGoerliIsSendOutTxProcessed(c *C) {
	sendHash := "CCCB58610B0B65D5B1D8E5F16435254A787E324209D9B3877A8ECE68859A0F55"
	recvTopics := make([][]ethcommon.Hash, 4)
	recvTopics[3] = []ethcommon.Hash{ethcommon.HexToHash(sendHash)}
	recvTopics[0] = []ethcommon.Hash{logMPIReceiveSignatureHash}
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.Chains["ETH"].MPIContractAddress)},
		FromBlock: big.NewInt(0), // LastBlock has been processed;
		ToBlock:   nil,
		Topics:    recvTopics,
	}
	log.Debug().Msg("TestGoerliIsSendOutTxProcessed")
	client, err := ethclient.Dial(config.ETH_ENDPOINT)
	c.Assert(err, IsNil)
	logs, err := client.FilterLogs(context.Background(), query)
	c.Assert(err, IsNil)
	contractAbi, err := abi.JSON(strings.NewReader(config.MPI_ABI_STRING))
	c.Assert(err, IsNil)
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logMPIReceiveSignatureHash.Hex():
			fmt.Printf("Found sendHash %s on chain %s\n", sendHash, "ETH")
			retval, err := contractAbi.Unpack("ZetaMessageReceiveEvent", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking Unlock")
				continue
			}
			fmt.Printf("Topic 0: %s\n", vLog.Topics[0])
			fmt.Printf("Topic 1: %s\n", vLog.Topics[1])
			fmt.Printf("Topic 2: %s\n", vLog.Topics[2])
			//fmt.Printf("data: %d\n", retval[0].(*big.Int))
			fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)
			mMint := retval[1].(*big.Int).String()

			fmt.Printf("Confirmed! Sending PostConfirmation to zetacore...\n")
			sendhash := vLog.Topics[3].Hex()
			fmt.Printf("mMint: %s\n", mMint)
			fmt.Printf("sendHash %s\n", sendhash)
		}
	}
}
