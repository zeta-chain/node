//go:build btc_regtest
// +build btc_regtest

package bitcoin

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/zeta-chain/zetacore/zetaclient/interfaces"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/common"
)

type BitcoinClientTestSuite struct {
	suite.Suite
	BitcoinChainClient *ChainClient
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

	tss := interfaces.TestSigner{
		PrivKey: privateKey,
	}
	//client, err := NewBitcoinClient(common.BtcTestNetChain(), nil, tss, "", nil)
	client, err := NewBitcoinClient(common.BtcRegtestChain(), nil, tss, "/tmp", nil)
	suite.Require().NoError(err)
	suite.BitcoinChainClient = client
	skBytes, err := hex.DecodeString(skHex)
	suite.Require().NoError(err)
	suite.T().Logf("skBytes: %d", len(skBytes))

	btc := client.rpcClient

	_, err = btc.CreateWallet("smoketest")
	suite.Require().NoError(err)
	addr, err := btc.GetNewAddress("test")
	suite.Require().NoError(err)
	suite.T().Logf("deployer address: %s", addr)
	//err = btc.ImportPrivKey(privkeyWIF)
	//suite.Require().NoError(err)

	btc.GenerateToAddress(101, addr, nil)
	suite.Require().NoError(err)

	bal, err := btc.GetBalance("*")
	suite.Require().NoError(err)
	suite.T().Logf("balance: %f", bal.ToBTC())

	utxo, err := btc.ListUnspent()
	suite.Require().NoError(err)
	suite.T().Logf("utxo: %d", len(utxo))
	for _, u := range utxo {
		suite.T().Logf("utxo: %s %f", u.Address, u.Amount)
	}
}

func (suite *BitcoinClientTestSuite) TearDownSuite() {

}

func (suite *BitcoinClientTestSuite) Test0() {

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

	inTxs := FilterAndParseIncomingTx(
		block.Tx,
		uint64(block.Height),
		"tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2",
		&log.Logger,
		common.BtcRegtestChain().ChainId,
	)

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

	inTxs := FilterAndParseIncomingTx(
		block.Tx,
		uint64(block.Height),
		"tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2",
		&log.Logger,
		common.BtcRegtestChain().ChainId,
	)

	suite.Require().Equal(0, len(inTxs))
}

func (suite *BitcoinClientTestSuite) Test3() {
	client := suite.BitcoinChainClient.rpcClient
	res, err := client.EstimateSmartFee(1, &btcjson.EstimateModeConservative)
	suite.Require().NoError(err)
	suite.T().Logf("fee: %f", *res.FeeRate)
	suite.T().Logf("blocks: %d", res.Blocks)
	suite.T().Logf("errors: %s", res.Errors)
	gasPrice := big.NewFloat(0)
	gasPriceU64, _ := gasPrice.Mul(big.NewFloat(*res.FeeRate), big.NewFloat(1e8)).Uint64()
	suite.T().Logf("gas price: %d", gasPriceU64)

	bn, err := client.GetBlockCount()
	suite.Require().NoError(err)
	suite.T().Logf("block number %d", bn)
}
func TestBitcoinChainClient(t *testing.T) {
	suite.Run(t, new(BitcoinClientTestSuite))
}
