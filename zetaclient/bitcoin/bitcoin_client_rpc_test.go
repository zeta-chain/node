package bitcoin

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/zeta-chain/zetacore/common"
	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

type BitcoinClientTestSuite struct {
	suite.Suite
	BitcoinChainClient *BTCChainClient
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
	appContext := appcontext.NewAppContext(&corecontext.ZeraCoreContext{}, &config.Config{}, log.Logger)
	client, err := NewBitcoinClient(appContext, common.BtcRegtestChain(), nil, tss, "/tmp", nil, nil)
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

func getRPCClient(chainID int64) (*rpcclient.Client, error) {
	var connCfg *rpcclient.ConnConfig
	// mainnet
	if chainID == 8332 {
		connCfg = &rpcclient.ConnConfig{
			Host:         "127.0.0.1:8332", // mainnet endpoint goes here
			User:         "user",
			Pass:         "pass",
			Params:       "mainnet",
			HTTPPostMode: true,
			DisableTLS:   true,
		}
	}
	// testnet3
	if chainID == 18332 {
		connCfg = &rpcclient.ConnConfig{
			Host:         "127.0.0.1:8332", // testnet endpoint goes here
			User:         "user",
			Pass:         "pass",
			Params:       "testnet3",
			HTTPPostMode: true,
			DisableTLS:   true,
		}
	}
	return rpcclient.New(connCfg, nil)
}

func getFeeRate(client *rpcclient.Client, confTarget int64, estimateMode *btcjson.EstimateSmartFeeMode) (*big.Int, error) {
	feeResult, err := client.EstimateSmartFee(confTarget, estimateMode)
	if err != nil {
		return nil, err
	}
	if feeResult.Errors != nil {
		return nil, errors.New(strings.Join(feeResult.Errors, ", "))
	}
	if feeResult.FeeRate == nil {
		return nil, errors.New("fee rate is nil")
	}
	return new(big.Int).SetInt64(int64(*feeResult.FeeRate * 1e8)), nil
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
		&chaincfg.TestNet3Params,
		0.0,
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
		&chaincfg.TestNet3Params,
		0.0,
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

// func TestBitcoinChainClient(t *testing.T) {
// 	suite.Run(t, new(BitcoinClientTestSuite))
// }

// Remove prefix "Live" to run this live test
func LiveTestBitcoinFeeRate(t *testing.T) {
	// setup Bitcoin client
	client, err := getRPCClient(8332)
	require.NoError(t, err)
	bn, err := client.GetBlockCount()
	if err != nil {
		t.Error(err)
	}

	// get fee rate for 1 block target
	feeRateConservative1, errCon1 := getFeeRate(client, 1, &btcjson.EstimateModeConservative)
	if errCon1 != nil {
		t.Error(errCon1)
	}
	feeRateEconomical1, errEco1 := getFeeRate(client, 1, &btcjson.EstimateModeEconomical)
	if errEco1 != nil {
		t.Error(errEco1)
	}
	// get fee rate for 2 block target
	feeRateConservative2, errCon2 := getFeeRate(client, 2, &btcjson.EstimateModeConservative)
	if errCon2 != nil {
		t.Error(errCon2)
	}
	feeRateEconomical2, errEco2 := getFeeRate(client, 2, &btcjson.EstimateModeEconomical)
	if errEco2 != nil {
		t.Error(errEco2)
	}
	fmt.Printf("Block: %d, Conservative-1 fee rate: %d, Economical-1 fee rate: %d\n", bn, feeRateConservative1.Uint64(), feeRateEconomical1.Uint64())
	fmt.Printf("Block: %d, Conservative-2 fee rate: %d, Economical-2 fee rate: %d\n", bn, feeRateConservative2.Uint64(), feeRateEconomical2.Uint64())

	// monitor fee rate every 5 minutes
	for {
		time.Sleep(time.Duration(5) * time.Minute)
		bn, err = client.GetBlockCount()
		feeRateConservative1, errCon1 = getFeeRate(client, 1, &btcjson.EstimateModeConservative)
		feeRateEconomical1, errEco1 = getFeeRate(client, 1, &btcjson.EstimateModeEconomical)
		feeRateConservative2, errCon2 = getFeeRate(client, 2, &btcjson.EstimateModeConservative)
		feeRateEconomical2, errEco2 = getFeeRate(client, 2, &btcjson.EstimateModeEconomical)
		if err != nil || errCon1 != nil || errEco1 != nil || errCon2 != nil || errEco2 != nil {
			continue
		}
		require.True(t, feeRateConservative1.Uint64() >= feeRateEconomical1.Uint64())
		require.True(t, feeRateConservative2.Uint64() >= feeRateEconomical2.Uint64())
		require.True(t, feeRateConservative1.Uint64() >= feeRateConservative2.Uint64())
		require.True(t, feeRateEconomical1.Uint64() >= feeRateEconomical2.Uint64())
		fmt.Printf("Block: %d, Conservative-1 fee rate: %d, Economical-1 fee rate: %d\n", bn, feeRateConservative1.Uint64(), feeRateEconomical1.Uint64())
		fmt.Printf("Block: %d, Conservative-2 fee rate: %d, Economical-2 fee rate: %d\n", bn, feeRateConservative2.Uint64(), feeRateEconomical2.Uint64())
	}
}

// compareAvgFeeRate compares fee rate with mempool.space for blocks [startBlock, endBlock]
func compareAvgFeeRate(t *testing.T, client *rpcclient.Client, startBlock int, endBlock int, testnet bool) {
	// mempool.space return 15 blocks [bn-14, bn] per request
	for bn := startBlock; bn >= endBlock; {
		// get mempool.space return blocks in descending order [bn, bn-14]
		mempoolBlocks, err := testutils.GetBlocks(context.Background(), bn, testnet)
		if err != nil {
			fmt.Printf("error GetBlocks %d: %s\n", bn, err)
			time.Sleep(10 * time.Second)
			continue
		}

		// calculate gas rate for each block
		for _, mb := range mempoolBlocks {
			// stop on end block
			if mb.Height < endBlock {
				break
			}
			bn = int(mb.Height) - 1

			// get block hash
			blkHash, err := client.GetBlockHash(int64(mb.Height))
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}
			// get block
			blockVb, err := client.GetBlockVerboseTx(blkHash)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}
			// calculate gas rate
			netParams := &chaincfg.MainNetParams
			if testnet {
				netParams = &chaincfg.TestNet3Params
			}
			gasRate, err := CalcBlockAvgFeeRate(blockVb, netParams)
			require.NoError(t, err)

			// compare with mempool.space
			if int(gasRate) == mb.Extras.AvgFeeRate {
				fmt.Printf("block %d: gas rate %d == mempool.space gas rate\n", mb.Height, gasRate)
			} else if int(gasRate) > mb.Extras.AvgFeeRate {
				fmt.Printf("block %d: gas rate %d >  mempool.space gas rate %d, diff: %f percent\n",
					mb.Height, gasRate, mb.Extras.AvgFeeRate, float64(int(gasRate)-mb.Extras.AvgFeeRate)/float64(mb.Extras.AvgFeeRate)*100)
			} else {
				fmt.Printf("block %d: gas rate %d <  mempool.space gas rate %d, diff: %f percent\n",
					mb.Height, gasRate, mb.Extras.AvgFeeRate, float64(mb.Extras.AvgFeeRate-int(gasRate))/float64(mb.Extras.AvgFeeRate)*100)
			}
		}
	}
}

// Remove prefix "Live" to run this live test
func LiveTestAvgFeeRateMainnetMempoolSpace(t *testing.T) {
	// setup Bitcoin client
	client, err := getRPCClient(8332)
	require.NoError(t, err)

	// test against mempool.space API for 10000 blocks
	//startBlock := 210000 * 3 // 3rd halving
	startBlock := 829596
	endBlock := startBlock - 10000

	compareAvgFeeRate(t, client, startBlock, endBlock, false)
}

// Remove prefix "Live" to run this live test
func LiveTestAvgFeeRateTestnetMempoolSpace(t *testing.T) {
	// setup Bitcoin client
	client, err := getRPCClient(18332)
	require.NoError(t, err)

	// test against mempool.space API for 10000 blocks
	//startBlock := 210000 * 12 // 12th halving
	startBlock := 2577600
	endBlock := startBlock - 10000

	compareAvgFeeRate(t, client, startBlock, endBlock, true)
}
