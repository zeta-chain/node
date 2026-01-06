package client_test

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
	types "github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	btc "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/common"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
)

// TestClientLive runs tests on a real note.
// Note that t.Parallel() is avoided due to potential rate limiting.
//
// You can get a free btc mainnet & testnet @ nownodes.io
// - mainnet: "https://btc.nownodes.io/<token>
// - testnet: "https://btc-testnet.nownodes.io/<token>
func TestClientLive(t *testing.T) {
	if !common.LiveTestEnabled() {
		t.Skip("skipping live test")
	}

	mainnetConfig := config.BTCConfig{
		RPCHost:   os.Getenv(common.EnvBtcRPCMainnet),
		RPCParams: "mainnet",
	}

	signetConfig := config.BTCConfig{
		RPCHost:   os.Getenv(common.EnvBtcRPCSignet),
		RPCParams: "signet",
	}

	testnet4Config := config.BTCConfig{
		RPCHost:   os.Getenv(common.EnvBtcRPCTestnet4),
		RPCParams: "testnet4",
	}

	t.Run("Healthcheck", func(t *testing.T) {
		t.Skip("most rpc won't allow private methods e.g. listUnspentMinMaxAddresses")

		// ARRANGE
		ts := newTestSuite(t, mainnetConfig)

		// ACT
		_, err := ts.Healthcheck(ts.ctx)

		// ASSERT
		require.NoError(t, err)
	})

	t.Run("GetBlockCount", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, mainnetConfig)

		// ACT
		bn, err := ts.GetBlockCount(ts.ctx)

		// ASSERT
		require.NoError(t, err)
		require.True(t, bn > 879_088)
	})

	t.Run("GetBlockHash", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, mainnetConfig)

		// ACT
		hash, err := ts.GetBlockHash(ts.ctx, 879088)

		// ASSERT
		require.NoError(t, err)
		require.NotEmpty(t, hash)

		// ACT #2
		block, err := ts.GetBlockHeader(ts.ctx, hash)

		// ASSERT #2
		require.NoError(t, err)
		require.NotEmpty(t, block)
	})

	t.Run("GetBlockHeightByStr", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, mainnetConfig)

		// the block hashes to test
		expectedHeight := int64(835053)
		hash := "00000000000000000000994a5d12976ec5bda078a7b9c27981f0a4e7a6d46d23"
		invalidHash := "invalidhash"

		// ACT #1
		// get block by invalid has
		_, err := ts.GetBlockHeightByStr(ts.ctx, invalidHash)
		require.ErrorContains(t, err, "unable to create btc hash from string")

		// ACT #2
		// get block height by block hash
		height, err := ts.GetBlockHeightByStr(ts.ctx, hash)
		require.NoError(t, err)
		require.Equal(t, expectedHeight, height)
	})

	t.Run("FilterAndParseIncomingTx", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, mainnetConfig)

		// net params
		net, err := chains.GetBTCChainParams(ts.chain.ChainId)
		require.NoError(t, err)

		// get the block that contains the incoming tx
		hashStr := "000000000000000000004a8ff16cd6be6c9410768c40c253868137f72f020f2d"

		block, err := ts.GetBlockVerboseByStr(ts.ctx, hashStr)
		require.NoError(t, err)

		// get fee rate multiplier
		feeRateMultiplier := observertypes.DefaultBTCOutboundGasPriceMultiplier.MustFloat64()

		inbounds, err := observer.FilterAndParseIncomingTx(
			ts.ctx,
			ts.Client,
			block.Tx,
			uint64(block.Height),
			testutils.TSSAddressBTCMainnet,
			feeRateMultiplier,
			ts.Logger,
			net,
		)

		require.NoError(t, err)
		require.Len(t, inbounds, 1)
		require.Equal(t, inbounds[0].Value+inbounds[0].DepositorFee, 0.001)
		require.Equal(t, inbounds[0].ToAddress, testutils.TSSAddressBTCMainnet)

		// the text memo is base64 std encoded string:DSRR1RmDCwWmxqY201/TMtsJdmA=
		// see https://mempool.space/tx/fb626645382defa1a68f4c07300ede722147c3b27a2fded698f0bf54b1110232
		memo, err := hex.DecodeString("6a956f0aed3b8625f20d696a5e934a5de8c27a2c")
		require.NoError(t, err)
		require.Equal(t, inbounds[0].MemoBytes, memo)
		require.Equal(t, inbounds[0].FromAddress, "bc1q4hffmv7m7yzglh3ryl3fcenwas46jyql7jpzdk")
		require.Equal(t, inbounds[0].BlockNumber, uint64(899645))
		require.Equal(t, inbounds[0].TxHash, "fb626645382defa1a68f4c07300ede722147c3b27a2fded698f0bf54b1110232")
	})

	t.Run("FilterAndParseIncomingTxNoop", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, mainnetConfig)

		// get a block that contains no incoming tx
		hashStr := "000000000000000000004a8ff16cd6be6c9410768c40c253868137f72f020f2d"

		block, err := ts.GetBlockVerboseByStr(ts.ctx, hashStr)
		require.NoError(t, err)

		// get fee rate multiplier
		feeRateMultiplier := observertypes.DefaultBTCOutboundGasPriceMultiplier.MustFloat64()

		// filter incoming tx
		inbounds, err := observer.FilterAndParseIncomingTx(
			ts.ctx,
			ts.Client,
			block.Tx,
			uint64(block.Height),
			"tb1qsa222mn2rhdq9cruxkz8p2teutvxuextx3ees2",
			feeRateMultiplier,
			ts.Logger,
			&chaincfg.TestNet3Params,
		)

		require.NoError(t, err)
		require.Empty(t, inbounds)
	})

	t.Run("GetRecentFeeRate", func(t *testing.T) {
		// ARRANGE
		// setup Bitcoin testnet client
		ts := newTestSuite(t, testnet4Config)

		// net params
		net, err := chains.GetBTCChainParams(ts.chain.ChainId)
		require.NoError(t, err)

		// ACT
		// get fee rate from recent blocks
		feeRate, err := btc.GetRecentFeeRate(ts.ctx, ts.Client, net)

		// ASSERT
		require.NoError(t, err)
		require.Greater(t, feeRate, uint64(0))
	})

	// LiveTestBitcoinFeeRate query Bitcoin mainnet fee rate every 5 minutes
	// and compares Conservative and Economical fee rates for different block targets (1 and 2)
	t.Run("BitcoinFeeRate", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, mainnetConfig)
		bn, err := ts.GetBlockCount(ts.ctx)
		require.NoError(t, err)

		// get fee rate for 1 block target
		feeRateConservative1, errCon1 := ts.getFeeRate(1, &types.EstimateModeConservative)
		if errCon1 != nil {
			t.Error(errCon1)
		}

		feeRateEconomical1, errEco1 := ts.getFeeRate(1, &types.EstimateModeEconomical)
		if errEco1 != nil {
			t.Error(errEco1)
		}

		// get fee rate for 2 block target
		feeRateConservative2, errCon2 := ts.getFeeRate(2, &types.EstimateModeConservative)
		if errCon2 != nil {
			t.Error(errCon2)
		}

		feeRateEconomical2, errEco2 := ts.getFeeRate(2, &types.EstimateModeEconomical)
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

			bn, err = ts.GetBlockCount(ts.ctx)
			feeRateConservative1, errCon1 = ts.getFeeRate(1, &types.EstimateModeConservative)
			feeRateEconomical1, errEco1 = ts.getFeeRate(1, &types.EstimateModeEconomical)
			feeRateConservative2, errCon2 = ts.getFeeRate(2, &types.EstimateModeConservative)
			feeRateEconomical2, errEco2 = ts.getFeeRate(2, &types.EstimateModeEconomical)
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
	})

	t.Run("GetTransactionInputSpender", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, mainnetConfig)

		// calculates block range to test
		startBlock, err := ts.GetBlockCount(ts.ctx)
		require.NoError(t, err)

		// go back to whatever block as needed
		endBlock := startBlock - 1

		// loop through mempool.space blocks backwards
	BLOCKLOOP:
		for bn := startBlock; bn >= endBlock; {
			// get mempool.space txs for the block
			_, mempoolTxs, err := ts.getMemPoolSpaceTxsByBlock(bn, false)
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
					vin := types.Vin{
						Txid: mpvin.TxID,
						Vout: mpvin.Vout,
					}
					senderAddr, err := ts.GetTransactionInputSpender(ts.ctx, mpvin.TxID, mpvin.Vout)
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
	})

	t.Run("GetTransactionInitiator", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, signetConfig)

		// Given a inscription commit tx
		commitTxSender := "tb1quy58ql4ss9q2f00p8xlhtl05ncp3rp6dvtxdtp"
		commitTxid := "e8072dcd232fe6f89037ee1d6a4ab75ed4c119798743b6edfbfe9af18811ad3e"

		// ACT
		initiator, err := ts.GetTransactionInitiator(ts.ctx, commitTxid)
		require.NoError(t, err)

		// ASSERT
		require.Equal(t, commitTxSender, initiator)
	})

	t.Run("GetTransactionFeeAndRate", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, testnet4Config)

		// calculates block range to test
		startBlock, err := ts.GetBlockCount(ts.ctx)
		require.NoError(t, err)

		// go back whatever blocks as needed
		endBlock := startBlock - 1

		// loop through mempool.space blocks backwards
		for bn := startBlock; bn > endBlock; {
			// get mempool.space txs for the block
			blkHash, mempoolTxs, err := ts.getMemPoolSpaceTxsByBlock(bn, true)
			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			}

			// get the block from rpc client
			block, err := ts.GetBlockVerbose(ts.ctx, blkHash)
			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			}

			// loop through each tx in the block (skip coinbase tx)
			for i := 1; i < len(block.Tx); {
				// sample 10 txs per block
				if i >= 10 {
					break
				}

				// the two txs from two different sources
				tx := block.Tx[i]
				mpTx := mempoolTxs[i]
				require.Equal(t, tx.Txid, mpTx.TxID)

				// get transaction fee rate for the raw result
				fee, feeRate, err := ts.GetTransactionFeeAndRate(ts.ctx, &tx)
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
	})

	t.Run("AvgFeeRateMainnetMempoolSpace", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, mainnetConfig)

		// test against mempool.space API for 10000 blocks
		// startBlock := 210000 * 3 // 3rd halving
		startBlock := 829596
		endBlock := startBlock - 1 // go back to whatever block as needed

		// ACT
		ts.compareAvgFeeRate(startBlock, endBlock, false)
	})

	t.Run("AvgFeeRateTestnetMempoolSpace", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, testnet4Config)

		// test against mempool.space API
		startBlock := 85240
		endBlock := startBlock - 1 // go back to whatever block as needed

		// ACT
		ts.compareAvgFeeRate(startBlock, endBlock, true)
	})

	t.Run("CalcDepositorFee", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t, mainnetConfig)

		// test tx hash
		// https://mempool.space/tx/8dc0d51f83810cec7fcb5b194caebfc5fc64b10f9fe21845dfecc621d2a28538
		hash, err := chainhash.NewHashFromStr("8dc0d51f83810cec7fcb5b194caebfc5fc64b10f9fe21845dfecc621d2a28538")
		require.NoError(t, err)

		// get the raw transaction result
		rawResult, err := ts.GetRawTransactionVerbose(ts.ctx, hash)
		require.NoError(t, err)

		// get fee rate multiplier
		feeRateMultiplier := observertypes.DefaultBTCOutboundGasPriceMultiplier.MustFloat64()

		t.Run("should return default depositor fee", func(t *testing.T) {
			depositorFee, err := btc.CalcDepositorFee(ts.ctx, ts.Client, rawResult, feeRateMultiplier, &chaincfg.RegressionNetParams)
			require.NoError(t, err)
			require.Equal(t, btc.DefaultDepositorFee, depositorFee)
		})

		t.Run("should return correct depositor fee for a given tx", func(t *testing.T) {
			depositorFee, err := btc.CalcDepositorFee(ts.ctx, ts.Client, rawResult, feeRateMultiplier, &chaincfg.MainNetParams)
			require.NoError(t, err)

			// the actual fee rate is 860 sat/vByte
			// #nosec G115 always in range
			expectedRate := int64(float64(860) * observertypes.DefaultBTCOutboundGasPriceMultiplier.MustFloat64())
			expectedFee := btc.DepositorFee(expectedRate)
			require.Equal(t, expectedFee, depositorFee)
		})
	})
}

type testSuite struct {
	t *testing.T
	*testlog.Log
	*client.Client
	ctx   context.Context
	chain chains.Chain
}

func newTestSuite(t *testing.T, cfg config.BTCConfig) *testSuite {
	logger := testlog.New(t)

	var chain chains.Chain
	switch cfg.RPCParams {
	case "mainnet":
		chain = chains.BitcoinMainnet
	case "signet":
		// signet also uses 'testnet3' RPC params, override it
		cfg.RPCParams = "testnet3"
		chain = chains.BitcoinSignetTestnet
	case "testnet4":
		// testnet4 also uses 'testnet3' RPC params, override it
		cfg.RPCParams = "testnet3"
		chain = chains.BitcoinTestnet4
	default:
		t.Fatalf("invalid RPC params: %s", cfg.RPCParams)
	}

	c, err := client.New(cfg, chain.ChainId, logger.Logger)
	require.NoError(t, err)

	return &testSuite{
		t:      t,
		Log:    logger,
		Client: c,
		ctx:    context.Background(),
		chain:  chain,
	}
}

// getMemPoolSpaceTxsByBlock gets mempool.space txs for a given block
func (ts *testSuite) getMemPoolSpaceTxsByBlock(
	blkNumber int64,
	testnet bool,
) (*chainhash.Hash, []testutils.MempoolTx, error) {
	blkHash, err := ts.GetBlockHash(ts.ctx, blkNumber)
	if err != nil {
		return nil, nil, err
	}

	// get mempool.space txs for the block
	mempoolTxs, err := testutils.GetBlockTxs(ts.ctx, blkHash.String(), testnet)
	if err != nil {
		return nil, nil, err
	}

	return blkHash, mempoolTxs, nil
}

func (ts *testSuite) getFeeRate(confTarget int64, estimateMode *types.EstimateSmartFeeMode) (*big.Int, error) {
	feeResult, err := ts.EstimateSmartFee(ts.ctx, confTarget, estimateMode)
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

func (ts *testSuite) compareAvgFeeRate(startBlock int, endBlock int, testnet bool) {
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
			blkHash, err := ts.GetBlockHash(ts.ctx, int64(mb.Height))
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}
			// get block
			blockVb, err := ts.GetBlockVerbose(ts.ctx, blkHash)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}
			// calculate gas rate
			netParams := &chaincfg.MainNetParams
			if testnet {
				netParams = &chaincfg.TestNet3Params
			}
			gasRate, err := btc.CalcBlockAvgFeeRate(blockVb, netParams)
			require.NoError(ts.t, err)

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
