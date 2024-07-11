package rpc_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

type BitcoinObserverTestSuite struct {
	suite.Suite
	rpcClient *rpcclient.Client
}

func (suite *BitcoinObserverTestSuite) SetupTest() {
	// test private key with EVM address
	//// EVM: 0x236C7f53a90493Bb423411fe4117Cb4c2De71DfB
	// BTC testnet3: muGe9prUBjQwEnX19zG26fVRHNi8z7kSPo
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	suite.Require().NoError(err)
	pkBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	suite.T().Logf("pubkey: %d", len(pkBytes))

	tss := &mocks.TSS{
		PrivKey: privateKey,
	}

	// create mock arguments for constructor
	chain := chains.BitcoinMainnet
	params := mocks.MockChainParams(chain.ChainId, 10)
	btcClient := mocks.NewMockBTCRPCClient()

	// create observer
	ob, err := observer.NewObserver(chain, btcClient, params, nil, tss, testutils.SQLiteMemory,
		base.DefaultLogger(), nil)
	suite.Require().NoError(err)
	suite.Require().NotNil(ob)
	suite.rpcClient, err = createRPCClient(18332)
	suite.Require().NoError(err)
	skBytes, err := hex.DecodeString(skHex)
	suite.Require().NoError(err)
	suite.T().Logf("skBytes: %d", len(skBytes))

	_, err = btcClient.CreateWallet("e2e")
	suite.Require().NoError(err)
	addr, err := btcClient.GetNewAddress("test")
	suite.Require().NoError(err)
	suite.T().Logf("deployer address: %s", addr)
	//err = btc.ImportPrivKey(privkeyWIF)
	//suite.Require().NoError(err)

	btcClient.GenerateToAddress(101, addr, nil)
	suite.Require().NoError(err)

	bal, err := btcClient.GetBalance("*")
	suite.Require().NoError(err)
	suite.T().Logf("balance: %f", bal.ToBTC())

	utxo, err := btcClient.ListUnspent()
	suite.Require().NoError(err)
	suite.T().Logf("utxo: %d", len(utxo))
	for _, u := range utxo {
		suite.T().Logf("utxo: %s %f", u.Address, u.Amount)
	}
}

func (suite *BitcoinObserverTestSuite) TearDownSuite() {
}

// createRPCClient creates a new Bitcoin RPC client for given chainID
func createRPCClient(chainID int64) (*rpcclient.Client, error) {
	var connCfg *rpcclient.ConnConfig
	rpcMainnet := os.Getenv("BTC_RPC_MAINNET")
	rpcTestnet := os.Getenv("BTC_RPC_TESTNET")

	// mainnet
	if chainID == chains.BitcoinMainnet.ChainId {
		connCfg = &rpcclient.ConnConfig{
			Host:         rpcMainnet, // mainnet endpoint goes here
			User:         "user",
			Pass:         "pass",
			Params:       "mainnet",
			HTTPPostMode: true,
			DisableTLS:   true,
		}
	}
	// testnet3
	if chainID == chains.BitcoinTestnet.ChainId {
		connCfg = &rpcclient.ConnConfig{
			Host:         rpcTestnet, // testnet endpoint goes here
			User:         "user",
			Pass:         "pass",
			Params:       "testnet3",
			HTTPPostMode: true,
			DisableTLS:   true,
		}
	}
	return rpcclient.New(connCfg, nil)
}

func getFeeRate(
	client *rpcclient.Client,
	confTarget int64,
	estimateMode *btcjson.EstimateSmartFeeMode,
) (*big.Int, error) {
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
func (suite *BitcoinObserverTestSuite) Test1() {
	feeResult, err := suite.rpcClient.EstimateSmartFee(1, nil)
	suite.Require().NoError(err)
	suite.T().Logf("fee result: %f", *feeResult.FeeRate)
	bn, err := suite.rpcClient.GetBlockCount()
	suite.Require().NoError(err)
	suite.T().Logf("block %d", bn)

	hashStr := "0000000000000032cb372f5d5d99c1ebf4430a3059b67c47a54dd626550fb50d"
	var hash chainhash.Hash
	err = chainhash.Decode(&hash, hashStr)
	suite.Require().NoError(err)

	block, err := suite.rpcClient.GetBlockVerboseTx(&hash)
	suite.Require().NoError(err)
	suite.T().Logf("block confirmation %d", block.Confirmations)
	suite.T().Logf("block txs len %d", len(block.Tx))

	inbounds, err := observer.FilterAndParseIncomingTx(
		suite.rpcClient,
		block.Tx,
		uint64(block.Height),
		"tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2",
		log.Logger,
		&chaincfg.TestNet3Params,
		0.0,
	)
	suite.Require().NoError(err)
	suite.Require().Equal(1, len(inbounds))
	suite.Require().Equal(inbounds[0].Value, 0.0001)
	suite.Require().Equal(inbounds[0].ToAddress, "tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2")
	// the text memo is base64 std encoded string:DSRR1RmDCwWmxqY201/TMtsJdmA=
	// see https://blockstream.info/testnet/tx/889bfa69eaff80a826286d42ec3f725fd97c3338357ddc3a1f543c2d6266f797
	memo, err := hex.DecodeString("0d2451D519830B05a6C6a636d35fd332dB097660")
	suite.Require().NoError(err)
	suite.Require().Equal((inbounds[0].MemoBytes), memo)
	suite.Require().Equal(inbounds[0].FromAddress, "tb1qyslx2s8evalx67n88wf42yv7236303ezj3tm2l")
	suite.T().Logf("from: %s", inbounds[0].FromAddress)
	suite.Require().Equal(inbounds[0].BlockNumber, uint64(2406185))
	suite.Require().Equal(inbounds[0].TxHash, "889bfa69eaff80a826286d42ec3f725fd97c3338357ddc3a1f543c2d6266f797")
}

// a tx with memo around 81B (is this allowed1?)
func (suite *BitcoinObserverTestSuite) Test2() {
	hashStr := "000000000000002fd8136dbf91708898da9d6ae61d7c354065a052568e2f2888"
	var hash chainhash.Hash
	err := chainhash.Decode(&hash, hashStr)
	suite.Require().NoError(err)

	block, err := suite.rpcClient.GetBlockVerboseTx(&hash)
	suite.Require().NoError(err)
	suite.T().Logf("block confirmation %d", block.Confirmations)
	suite.T().Logf("block height %d", block.Height)
	suite.T().Logf("block txs len %d", len(block.Tx))

	inbounds, err := observer.FilterAndParseIncomingTx(
		suite.rpcClient,
		block.Tx,
		uint64(block.Height),
		"tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2",
		log.Logger,
		&chaincfg.TestNet3Params,
		0.0,
	)
	suite.Require().NoError(err)
	suite.Require().Equal(0, len(inbounds))
}

// TestBitcoinObserverLive is a phony test to run each live test individually
func TestBitcoinObserverLive(t *testing.T) {
	// suite.Run(t, new(BitcoinClientTestSuite))

	// LiveTestNewRPCClient(t)
	// LiveTestGetBlockHeightByHash(t)
	// LiveTestBitcoinFeeRate(t)
	// LiveTestAvgFeeRateMainnetMempoolSpace(t)
	// LiveTestAvgFeeRateTestnetMempoolSpace(t)
	// LiveTestGetRecentFeeRate(t)
	// LiveTestGetSenderByVin(t)
}

// LiveTestNewRPCClient creates a new Bitcoin RPC client
func LiveTestNewRPCClient(t *testing.T) {
	btcConfig := config.BTCConfig{
		RPCUsername: "user",
		RPCPassword: "pass",
		RPCHost:     "bitcoin.rpc.zetachain.com/6315704c-49bc-4649-8b9d-e9278a1dfeb3",
		RPCParams:   "mainnet",
	}

	// create Bitcoin RPC client
	client, err := rpc.NewRPCClient(btcConfig)
	require.NoError(t, err)

	// get block count
	bn, err := client.GetBlockCount()
	require.NoError(t, err)
	require.Greater(t, bn, int64(0))
}

// LiveTestGetBlockHeightByHash queries Bitcoin block height by hash
func LiveTestGetBlockHeightByHash(t *testing.T) {
	// setup Bitcoin client
	client, err := createRPCClient(chains.BitcoinMainnet.ChainId)
	require.NoError(t, err)

	// the block hashes to test
	expectedHeight := int64(835053)
	hash := "00000000000000000000994a5d12976ec5bda078a7b9c27981f0a4e7a6d46d23"
	invalidHash := "invalidhash"

	// get block by invalid hash
	_, err = rpc.GetBlockHeightByHash(client, invalidHash)
	require.ErrorContains(t, err, "error decoding block hash")

	// get block height by block hash
	height, err := rpc.GetBlockHeightByHash(client, hash)
	require.NoError(t, err)
	require.Equal(t, expectedHeight, height)
}

// LiveTestBitcoinFeeRate query Bitcoin mainnet fee rate every 5 minutes
// and compares Conservative and Economical fee rates for different block targets (1 and 2)
func LiveTestBitcoinFeeRate(t *testing.T) {
	// setup Bitcoin client
	client, err := createRPCClient(chains.BitcoinMainnet.ChainId)
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
	fmt.Printf(
		"Block: %d, Conservative-1 fee rate: %d, Economical-1 fee rate: %d\n",
		bn,
		feeRateConservative1.Uint64(),
		feeRateEconomical1.Uint64(),
	)
	fmt.Printf(
		"Block: %d, Conservative-2 fee rate: %d, Economical-2 fee rate: %d\n",
		bn,
		feeRateConservative2.Uint64(),
		feeRateEconomical2.Uint64(),
	)

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
		fmt.Printf(
			"Block: %d, Conservative-1 fee rate: %d, Economical-1 fee rate: %d\n",
			bn,
			feeRateConservative1.Uint64(),
			feeRateEconomical1.Uint64(),
		)
		fmt.Printf(
			"Block: %d, Conservative-2 fee rate: %d, Economical-2 fee rate: %d\n",
			bn,
			feeRateConservative2.Uint64(),
			feeRateEconomical2.Uint64(),
		)
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
			gasRate, err := bitcoin.CalcBlockAvgFeeRate(blockVb, netParams)
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

// LiveTestAvgFeeRateMainnetMempoolSpace compares calculated fee rate with mempool.space fee rate for mainnet
func LiveTestAvgFeeRateMainnetMempoolSpace(t *testing.T) {
	// setup Bitcoin client
	client, err := createRPCClient(chains.BitcoinMainnet.ChainId)
	require.NoError(t, err)

	// test against mempool.space API for 10000 blocks
	//startBlock := 210000 * 3 // 3rd halving
	startBlock := 829596
	endBlock := startBlock - 10000

	compareAvgFeeRate(t, client, startBlock, endBlock, false)
}

// LiveTestAvgFeeRateTestnetMempoolSpace compares calculated fee rate with mempool.space fee rate for testnet
func LiveTestAvgFeeRateTestnetMempoolSpace(t *testing.T) {
	// setup Bitcoin client
	client, err := createRPCClient(chains.BitcoinTestnet.ChainId)
	require.NoError(t, err)

	// test against mempool.space API for 10000 blocks
	//startBlock := 210000 * 12 // 12th halving
	startBlock := 2577600
	endBlock := startBlock - 10000

	compareAvgFeeRate(t, client, startBlock, endBlock, true)
}

// LiveTestGetRecentFeeRate gets the highest fee rate from recent blocks
func LiveTestGetRecentFeeRate(t *testing.T) {
	// setup Bitcoin testnet client
	client, err := createRPCClient(chains.BitcoinTestnet.ChainId)
	require.NoError(t, err)

	// get fee rate from recent blocks
	feeRate, err := rpc.GetRecentFeeRate(client, &chaincfg.TestNet3Params)
	require.NoError(t, err)
	require.Greater(t, feeRate, uint64(0))
}

// LiveTestGetSenderByVin gets sender address for each vin and compares with mempool.space sender address
func LiveTestGetSenderByVin(t *testing.T) {
	// setup Bitcoin client
	chainID := chains.BitcoinMainnet.ChainId
	client, err := createRPCClient(chainID)
	require.NoError(t, err)

	// net params
	net, err := chains.GetBTCChainParams(chainID)
	require.NoError(t, err)
	testnet := false
	if chainID == chains.BitcoinTestnet.ChainId {
		testnet = true
	}

	// calculates block range to test
	startBlock, err := client.GetBlockCount()
	require.NoError(t, err)
	endBlock := startBlock - 5000

	// loop through mempool.space blocks in descending order
BLOCKLOOP:
	for bn := startBlock; bn >= endBlock; {
		// get block hash
		blkHash, err := client.GetBlockHash(int64(bn))
		if err != nil {
			fmt.Printf("error GetBlockHash for block %d: %s\n", bn, err)
			time.Sleep(3 * time.Second)
			continue
		}

		// get mempool.space txs for the block
		mempoolTxs, err := testutils.GetBlockTxs(context.Background(), blkHash.String(), testnet)
		if err != nil {
			fmt.Printf("error GetBlockTxs %d: %s\n", bn, err)
			time.Sleep(10 * time.Second)
			continue
		}

		// loop through each tx in the block
		for i, mptx := range mempoolTxs {
			// sample 10 txs per block
			if i >= 10 {
				break
			}
			for _, mpvin := range mptx.Vin {
				// skip coinbase tx
				if mpvin.IsCoinbase {
					continue
				}
				// get sender address for each vin
				vin := btcjson.Vin{
					Txid: mpvin.TxID,
					Vout: mpvin.Vout,
				}
				senderAddr, err := observer.GetSenderAddressByVin(client, vin, net)
				if err != nil {
					fmt.Printf("error GetSenderAddressByVin for block %d, tx %s vout %d: %s\n", bn, vin.Txid, vin.Vout, err)
					time.Sleep(3 * time.Second)
					continue BLOCKLOOP // retry the block
				}
				if senderAddr != mpvin.Prevout.ScriptpubkeyAddress {
					t.Errorf("block %d, tx %s, vout %d: want %s, got %s\n", bn, vin.Txid, vin.Vout, mpvin.Prevout.ScriptpubkeyAddress, senderAddr)
				} else {
					fmt.Printf("block: %d sender address type: %s\n", bn, mpvin.Prevout.ScriptpubkeyType)
				}
			}
		}
		bn--
		time.Sleep(500 * time.Millisecond)
	}
}
