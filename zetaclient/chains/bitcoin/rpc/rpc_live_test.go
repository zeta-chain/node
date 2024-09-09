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

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/node/zetaclient/common"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

// createRPCClient creates a new Bitcoin RPC client for given chainID
func createRPCClient(chainID int64) (*rpcclient.Client, error) {
	var connCfg *rpcclient.ConnConfig
	rpcMainnet := os.Getenv(common.EnvBtcRPCMainnet)
	rpcTestnet := os.Getenv(common.EnvBtcRPCTestnet)

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

// getFeeRate is a helper function to get fee rate for a given confirmation target
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

// getMempoolSpaceTxsByBlock gets mempool.space txs for a given block
func getMempoolSpaceTxsByBlock(
	t *testing.T,
	client *rpcclient.Client,
	blkNumber int64,
	testnet bool,
) (*chainhash.Hash, []testutils.MempoolTx, error) {
	blkHash, err := client.GetBlockHash(blkNumber)
	if err != nil {
		t.Logf("error GetBlockHash for block %d: %s\n", blkNumber, err)
		return nil, nil, err
	}

	// get mempool.space txs for the block
	mempoolTxs, err := testutils.GetBlockTxs(context.Background(), blkHash.String(), testnet)
	if err != nil {
		t.Logf("error GetBlockTxs %d: %s\n", blkNumber, err)
		return nil, nil, err
	}

	return blkHash, mempoolTxs, nil
}

// Test_BitcoinLive is a phony test to run each live test individually
func Test_BitcoinLive(t *testing.T) {
	// LiveTest_FilterAndParseIncomingTx(t)
	// LiveTest_FilterAndParseIncomingTx_Nop(t)
	// LiveTest_NewRPCClient(t)
	// LiveTest_GetBlockHeightByHash(t)
	// LiveTest_BitcoinFeeRate(t)
	// LiveTest_AvgFeeRateMainnetMempoolSpace(t)
	// LiveTest_AvgFeeRateTestnetMempoolSpace(t)
	// LiveTest_GetRecentFeeRate(t)
	// LiveTest_GetSenderByVin(t)
	// LiveTest_GetTransactionFeeAndRate(t)
	// LiveTest_CalcDepositorFeeV2(t)
}

func LiveTest_FilterAndParseIncomingTx(t *testing.T) {
	// setup Bitcoin client
	client, err := createRPCClient(chains.BitcoinTestnet.ChainId)
	require.NoError(t, err)

	// get the block that contains the incoming tx
	hashStr := "0000000000000032cb372f5d5d99c1ebf4430a3059b67c47a54dd626550fb50d"
	hash, err := chainhash.NewHashFromStr(hashStr)
	require.NoError(t, err)

	block, err := client.GetBlockVerboseTx(hash)
	require.NoError(t, err)

	// filter incoming tx
	inbounds, err := observer.FilterAndParseIncomingTx(
		client,
		block.Tx,
		uint64(block.Height),
		"tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2",
		log.Logger,
		&chaincfg.TestNet3Params,
		0.0,
	)
	require.NoError(t, err)
	require.Len(t, inbounds, 1)
	require.Equal(t, inbounds[0].Value, 0.0001)
	require.Equal(t, inbounds[0].ToAddress, "tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2")

	// the text memo is base64 std encoded string:DSRR1RmDCwWmxqY201/TMtsJdmA=
	// see https://blockstream.info/testnet/tx/889bfa69eaff80a826286d42ec3f725fd97c3338357ddc3a1f543c2d6266f797
	memo, err := hex.DecodeString("4453525231526d444377576d7871593230312f544d74734a646d413d")
	require.NoError(t, err)
	require.Equal(t, inbounds[0].MemoBytes, memo)
	require.Equal(t, inbounds[0].FromAddress, "tb1qyslx2s8evalx67n88wf42yv7236303ezj3tm2l")
	require.Equal(t, inbounds[0].BlockNumber, uint64(2406185))
	require.Equal(t, inbounds[0].TxHash, "889bfa69eaff80a826286d42ec3f725fd97c3338357ddc3a1f543c2d6266f797")
}

func LiveTest_FilterAndParseIncomingTx_Nop(t *testing.T) {
	// setup Bitcoin client
	client, err := createRPCClient(chains.BitcoinTestnet.ChainId)
	require.NoError(t, err)

	// get a block that contains no incoming tx
	hashStr := "000000000000002fd8136dbf91708898da9d6ae61d7c354065a052568e2f2888"
	hash, err := chainhash.NewHashFromStr(hashStr)
	require.NoError(t, err)

	block, err := client.GetBlockVerboseTx(hash)
	require.NoError(t, err)

	// filter incoming tx
	inbounds, err := observer.FilterAndParseIncomingTx(
		client,
		block.Tx,
		uint64(block.Height),
		"tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2",
		log.Logger,
		&chaincfg.TestNet3Params,
		0.0,
	)

	require.NoError(t, err)
	require.Empty(t, inbounds)
}

// TestBitcoinObserverLive is a phony test to run each live test individually
func TestBitcoinObserverLive(t *testing.T) {
	if !common.LiveTestEnabled() {
		return
	}

	LiveTest_NewRPCClient(t)
	LiveTest_CheckRPCStatus(t)
	LiveTest_GetBlockHeightByHash(t)
	LiveTest_BitcoinFeeRate(t)
	LiveTest_AvgFeeRateMainnetMempoolSpace(t)
	LiveTest_AvgFeeRateTestnetMempoolSpace(t)
	LiveTest_GetRecentFeeRate(t)
	LiveTest_GetSenderByVin(t)
}

// LiveTestNewRPCClient creates a new Bitcoin RPC client
func LiveTest_NewRPCClient(t *testing.T) {
	btcConfig := config.BTCConfig{
		RPCUsername: "user",
		RPCPassword: "pass",
		RPCHost:     os.Getenv(common.EnvBtcRPCTestnet),
		RPCParams:   "testnet3",
	}

	// create Bitcoin RPC client
	client, err := rpc.NewRPCClient(btcConfig)
	require.NoError(t, err)

	// get block count
	bn, err := client.GetBlockCount()
	require.NoError(t, err)
	require.Greater(t, bn, int64(0))
}

// Live_TestCheckRPCStatus checks the RPC status of the Bitcoin chain
func LiveTest_CheckRPCStatus(t *testing.T) {
	// setup Bitcoin client
	chainID := chains.BitcoinMainnet.ChainId
	client, err := createRPCClient(chainID)
	require.NoError(t, err)

	// decode tss address
	tssAddress, err := chains.DecodeBtcAddress(testutils.TSSAddressBTCMainnet, chainID)
	require.NoError(t, err)

	// check RPC status
	_, err = rpc.CheckRPCStatus(client, tssAddress)
	require.NoError(t, err)
}

// LiveTestGetBlockHeightByHash queries Bitcoin block height by hash
func LiveTest_GetBlockHeightByHash(t *testing.T) {
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
func LiveTest_BitcoinFeeRate(t *testing.T) {
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

	// monitor fee rate every 5 minutes, adjust the iteration count as needed
	for i := 0; i < 1; i++ {
		// please uncomment this interval for long running test
		//time.Sleep(time.Duration(5) * time.Minute)

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
func LiveTest_AvgFeeRateMainnetMempoolSpace(t *testing.T) {
	// setup Bitcoin client
	client, err := createRPCClient(chains.BitcoinMainnet.ChainId)
	require.NoError(t, err)

	// test against mempool.space API for 10000 blocks
	//startBlock := 210000 * 3 // 3rd halving
	startBlock := 829596
	endBlock := startBlock - 1 // go back to whatever block as needed

	compareAvgFeeRate(t, client, startBlock, endBlock, false)
}

// LiveTestAvgFeeRateTestnetMempoolSpace compares calculated fee rate with mempool.space fee rate for testnet
func LiveTest_AvgFeeRateTestnetMempoolSpace(t *testing.T) {
	// setup Bitcoin client
	client, err := createRPCClient(chains.BitcoinTestnet.ChainId)
	require.NoError(t, err)

	// test against mempool.space API for 10000 blocks
	//startBlock := 210000 * 12 // 12th halving
	startBlock := 2577600
	endBlock := startBlock - 1 // go back to whatever block as needed

	compareAvgFeeRate(t, client, startBlock, endBlock, true)
}

// LiveTestGetRecentFeeRate gets the highest fee rate from recent blocks
func LiveTest_GetRecentFeeRate(t *testing.T) {
	// setup Bitcoin testnet client
	client, err := createRPCClient(chains.BitcoinTestnet.ChainId)
	require.NoError(t, err)

	// get fee rate from recent blocks
	feeRate, err := bitcoin.GetRecentFeeRate(client, &chaincfg.TestNet3Params)
	require.NoError(t, err)
	require.Greater(t, feeRate, uint64(0))
}

// LiveTest_GetSenderByVin gets sender address for each vin and compares with mempool.space sender address
func LiveTest_GetSenderByVin(t *testing.T) {
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
	endBlock := startBlock - 1 // go back to whatever block as needed

	// loop through mempool.space blocks backwards
BLOCKLOOP:
	for bn := startBlock; bn >= endBlock; {
		// get mempool.space txs for the block
		_, mempoolTxs, err := getMempoolSpaceTxsByBlock(t, client, bn, testnet)
		if err != nil {
			time.Sleep(3 * time.Second)
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
		time.Sleep(100 * time.Millisecond)
	}
}

// LiveTestGetTransactionFeeAndRate gets the transaction fee and rate for each tx and compares with mempool.space fee rate
func LiveTest_GetTransactionFeeAndRate(t *testing.T) {
	// setup Bitcoin client
	chainID := chains.BitcoinTestnet.ChainId
	client, err := createRPCClient(chainID)
	require.NoError(t, err)

	// testnet or mainnet
	testnet := false
	if chainID == chains.BitcoinTestnet.ChainId {
		testnet = true
	}

	// calculates block range to test
	startBlock, err := client.GetBlockCount()
	require.NoError(t, err)
	endBlock := startBlock - 100 // go back whatever blocks as needed

	// loop through mempool.space blocks backwards
	for bn := startBlock; bn >= endBlock; {
		// get mempool.space txs for the block
		blkHash, mempoolTxs, err := getMempoolSpaceTxsByBlock(t, client, bn, testnet)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		// get the block from rpc client
		block, err := client.GetBlockVerboseTx(blkHash)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}

		// loop through each tx in the block (skip coinbase tx)
		for i := 1; i < len(block.Tx); {
			// sample 20 txs per block
			if i >= 20 {
				break
			}

			// the two txs from two different sources
			tx := block.Tx[i]
			mpTx := mempoolTxs[i]
			require.Equal(t, tx.Txid, mpTx.TxID)

			// get transaction fee rate for the raw result
			fee, feeRate, err := rpc.GetTransactionFeeAndRate(client, &tx)
			if err != nil {
				t.Logf("error GetTransactionFeeRate %s: %s\n", mpTx.TxID, err)
				continue
			}
			require.EqualValues(t, mpTx.Fee, fee)
			require.EqualValues(t, mpTx.Weight, tx.Weight)

			// calculate mempool.space fee rate
			vBytes := mpTx.Weight / blockchain.WitnessScaleFactor
			mpFeeRate := int64(mpTx.Fee / vBytes)

			// compare our fee rate with mempool.space fee rate
			var diff int64
			var diffPercent float64
			if feeRate == mpFeeRate {
				fmt.Printf("tx %s: [our rate] %5d == %5d [mempool.space]", mpTx.TxID, feeRate, mpFeeRate)
			} else if feeRate > mpFeeRate {
				diff = feeRate - mpFeeRate
				fmt.Printf("tx %s: [our rate] %5d >  %5d [mempool.space]", mpTx.TxID, feeRate, mpFeeRate)
			} else {
				diff = mpFeeRate - feeRate
				fmt.Printf("tx %s: [our rate] %5d <  %5d [mempool.space]", mpTx.TxID, feeRate, mpFeeRate)
			}

			// print the diff percentage
			diffPercent = float64(diff) / float64(mpFeeRate) * 100
			if diff > 0 {
				fmt.Printf(", diff: %f%%\n", diffPercent)
			} else {
				fmt.Printf("\n")
			}

			// the expected diff percentage should be within 5%
			if mpFeeRate >= 20 {
				require.LessOrEqual(t, diffPercent, 5.0)
			} else {
				// for small fee rate, the absolute diff should be within 1 satoshi/vByte
				require.LessOrEqual(t, diff, int64(1))
			}

			// next tx
			i++
		}

		bn--
		time.Sleep(100 * time.Millisecond)
	}
}

func LiveTest_CalcDepositorFeeV2(t *testing.T) {
	// setup Bitcoin client
	client, err := createRPCClient(chains.BitcoinMainnet.ChainId)
	require.NoError(t, err)

	// test tx hash
	// https://mempool.space/tx/8dc0d51f83810cec7fcb5b194caebfc5fc64b10f9fe21845dfecc621d2a28538
	hash, err := chainhash.NewHashFromStr("8dc0d51f83810cec7fcb5b194caebfc5fc64b10f9fe21845dfecc621d2a28538")
	require.NoError(t, err)

	// get the raw transaction result
	rawResult, err := client.GetRawTransactionVerbose(hash)
	require.NoError(t, err)

	t.Run("should return default depositor fee", func(t *testing.T) {
		depositorFee, err := bitcoin.CalcDepositorFeeV2(client, rawResult, &chaincfg.RegressionNetParams)
		require.NoError(t, err)
		require.Equal(t, bitcoin.DefaultDepositorFee, depositorFee)
	})

	t.Run("should return correct depositor fee for a given tx", func(t *testing.T) {
		depositorFee, err := bitcoin.CalcDepositorFeeV2(client, rawResult, &chaincfg.MainNetParams)
		require.NoError(t, err)

		// the actual fee rate is 860 sat/vByte
		// #nosec G115 always in range
		expectedRate := int64(float64(860) * common.BTCOutboundGasPriceMultiplier)
		expectedFee := bitcoin.DepositorFee(expectedRate)
		require.Equal(t, expectedFee, depositorFee)
	})
}
