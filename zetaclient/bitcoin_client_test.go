package zetaclient

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/common"
	"testing"
)

type BitcoinClientTestSuite struct {
	suite.Suite
	BitcoinChainClient *BitcoinChainClient
}

func (suite *BitcoinClientTestSuite) SetupTest() {
	// test private key with EVM address
	//// EVM: 0x236C7f53a90493Bb423411fe4117Cb4c2De71DfB
	// BTC testnet3: muGe9prUBjQwEnX19zG26fVRHNi8z7kSPo
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	suite.Require().NoError(err)
	pkBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	suite.T().Logf("pubkey: %d", len(pkBytes))

	tss := TestSigner{
		PrivKey: privateKey,
	}
	client, err := NewBitcoinClient(common.BTCTestnetChain, nil, tss, "", nil)
	suite.Require().NoError(err)
	suite.BitcoinChainClient = client

	//suite.BitcoinChainClient.Start()
}

func (suite *BitcoinClientTestSuite) TearDownSuite() {

}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *BitcoinClientTestSuite) Test1() {
	feeResult, err := suite.BitcoinChainClient.rpcClient.EstimateSmartFee(1, nil)
	suite.Require().NoError(err)
	suite.T().Logf("fee result: %f", *feeResult.FeeRate)
	bn, err := suite.BitcoinChainClient.rpcClient.GetBlockCount()
	suite.Require().NoError(err)
	suite.T().Logf("block %d", bn)

	hashStr := "0000000000000032cb372f5d5d99c1ebf4430a3059b67c47a54dd626550fb50d"
	var hash chainhash.Hash
	err = chainhash.Decode(&hash, hashStr)
	suite.Require().NoError(err)

	//:= suite.BitcoinChainClient.rpcClient.GetBlock(&hash)
	block, err := suite.BitcoinChainClient.rpcClient.GetBlockVerboseTx(&hash)
	suite.Require().NoError(err)
	suite.T().Logf("block confirmation %d", block.Confirmations)
	suite.T().Logf("block txs len %d", len(block.Tx))

	inTxs := FilterAndParseIncomingTx(block.Tx, uint64(block.Height), "tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2", &log.Logger)

	suite.Require().Equal(1, len(inTxs))
	suite.Require().Equal(inTxs[0].Value, 0.0001)
	suite.Require().Equal(inTxs[0].ToAddress, "tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2")
	// the text memo is base64 std encoded string:DSRR1RmDCwWmxqY201/TMtsJdmA=
	// see https://blockstream.info/testnet/tx/889bfa69eaff80a826286d42ec3f725fd97c3338357ddc3a1f543c2d6266f797
	memo, err := hex.DecodeString("0d2451D519830B05a6C6a636d35fd332dB097660")
	suite.Require().NoError(err)
	suite.Require().Equal((inTxs[0].MemoBytes), memo)
	suite.Require().Equal(inTxs[0].FromAddress, "tb1qyslx2s8evalx67n88wf42yv7236303ezj3tm2l")
	suite.T().Logf("from: %s", inTxs[0].FromAddress)
	suite.Require().Equal(inTxs[0].BlockNumber, uint64(2406185))
	suite.Require().Equal(inTxs[0].TxHash, "889bfa69eaff80a826286d42ec3f725fd97c3338357ddc3a1f543c2d6266f797")
}

// a tx with memo around 81B (is this allowed1?)
func (suite *BitcoinClientTestSuite) Test2() {
	hashStr := "000000000000002fd8136dbf91708898da9d6ae61d7c354065a052568e2f2888"
	var hash chainhash.Hash
	err := chainhash.Decode(&hash, hashStr)
	suite.Require().NoError(err)

	//:= suite.BitcoinChainClient.rpcClient.GetBlock(&hash)
	block, err := suite.BitcoinChainClient.rpcClient.GetBlockVerboseTx(&hash)
	suite.Require().NoError(err)
	suite.T().Logf("block confirmation %d", block.Confirmations)
	suite.T().Logf("block height %d", block.Height)
	suite.T().Logf("block txs len %d", len(block.Tx))

	inTxs := FilterAndParseIncomingTx(block.Tx, uint64(block.Height), "tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2", &log.Logger)

	suite.Require().Equal(0, len(inTxs))
}

func TestBitcoinChainClient(t *testing.T) {
	suite.Run(t, new(BitcoinClientTestSuite))
}
