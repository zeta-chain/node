package observer_test

import (
	"sync"
	"testing"

	"cosmossdk.io/math"

	lru "github.com/hashicorp/golang-lru"
	"github.com/onrik/ethrpc"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

// the relative path to the testdata directory
var TestDataDir = "../../../"

// getAppContext creates an app context for unit tests
func getAppContext(evmChain chains.Chain, evmChainParams *observertypes.ChainParams) (*context.AppContext, config.EVMConfig) {
	// create config
	cfg := config.NewConfig()
	cfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
		Chain:    evmChain,
		Endpoint: "http://localhost:8545",
	}
	// create core context
	coreCtx := context.NewZetaCoreContext(cfg)
	evmChainParamsMap := make(map[int64]*observertypes.ChainParams)
	evmChainParamsMap[evmChain.ChainId] = evmChainParams

	// feed chain params
	coreCtx.Update(
		&observertypes.Keygen{},
		[]chains.Chain{evmChain},
		evmChainParamsMap,
		nil,
		"",
		*sample.CrosschainFlags(),
		sample.HeaderSupportedChains(),
		true,
		zerolog.Logger{},
	)
	// create app context
	appCtx := context.NewAppContext(coreCtx, cfg)
	return appCtx, cfg.EVMChainConfigs[evmChain.ChainId]
}

// MockEVMObserver creates a mock ChainObserver with custom chain, TSS, params etc
func MockEVMObserver(
	t *testing.T,
	chain chains.Chain,
	evmClient interfaces.EVMRPCClient,
	evmJSONRPC interfaces.EVMJSONRPCClient,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	lastBlock uint64,
	params observertypes.ChainParams) *observer.Observer {
	// use default mock zetacore client if not provided
	if zetacoreClient == nil {
		zetacoreClient = mocks.NewMockZetaCoreClient()
	}
	// use default mock tss if not provided
	if tss == nil {
		tss = mocks.NewTSSMainnet()
	}
	// create app context
	appCtx, evmCfg := getAppContext(chain, &params)

	// create chain observer
	client, err := observer.NewObserver(appCtx, zetacoreClient, tss, "", common.ClientLogger{}, evmCfg, nil)
	require.NoError(t, err)
	client.WithEvmClient(evmClient)
	client.WithEvmJSONRPC(evmJSONRPC)
	client.SetLastBlockHeight(lastBlock)

	return client
}

func Test_BlockCache(t *testing.T) {
	// create client
	blockCache, err := lru.New(1000)
	require.NoError(t, err)
	ob := &observer.Observer{Mu: &sync.Mutex{}}
	ob.WithBlockCache(blockCache)

	// delete non-existing block should not panic
	blockNumber := uint64(10388180)
	ob.RemoveCachedBlock(blockNumber)

	// add a block
	block := &ethrpc.Block{
		// #nosec G701 always in range
		Number: int(blockNumber),
	}
	blockCache.Add(blockNumber, block)
	ob.WithBlockCache(blockCache)

	// block should be in cache
	_, err = ob.GetBlockByNumberCached(blockNumber)
	require.NoError(t, err)

	// delete the block should not panic
	ob.RemoveCachedBlock(blockNumber)
}

func Test_CheckTxInclusion(t *testing.T) {
	// load archived evm outtx Gas
	// https://etherscan.io/tx/0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3
	chainID := int64(1)
	coinType := coin.CoinType_Gas
	outtxHash := "0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3"
	tx, receipt := testutils.LoadEVMOuttxNReceipt(t, TestDataDir, chainID, outtxHash, coinType)

	// load archived evm block
	// https://etherscan.io/block/19363323
	blockNumber := receipt.BlockNumber.Uint64()
	block := testutils.LoadEVMBlock(t, TestDataDir, chainID, blockNumber, true)

	// create client
	blockCache, err := lru.New(1000)
	require.NoError(t, err)
	ob := &observer.Observer{Mu: &sync.Mutex{}}

	// save block to cache
	blockCache.Add(blockNumber, block)
	ob.WithBlockCache(blockCache)

	t.Run("should pass for archived outtx", func(t *testing.T) {
		err := ob.CheckTxInclusion(tx, receipt)
		require.NoError(t, err)
	})
	t.Run("should fail on tx index out of range", func(t *testing.T) {
		// modify tx index to invalid number
		copyReceipt := *receipt
		// #nosec G701 non negative value
		copyReceipt.TransactionIndex = uint(len(block.Transactions))
		err := ob.CheckTxInclusion(tx, &copyReceipt)
		require.ErrorContains(t, err, "out of range")
	})
	t.Run("should fail on tx hash mismatch", func(t *testing.T) {
		// change the tx at position 'receipt.TransactionIndex' to a different tx
		priorTx := block.Transactions[receipt.TransactionIndex-1]
		block.Transactions[receipt.TransactionIndex] = priorTx
		blockCache.Add(blockNumber, block)
		ob.WithBlockCache(blockCache)

		// check inclusion should fail
		err := ob.CheckTxInclusion(tx, receipt)
		require.ErrorContains(t, err, "has different hash")

		// wrong block should be removed from cache
		_, ok := blockCache.Get(blockNumber)
		require.False(t, ok)
	})
}

func Test_VoteOutboundBallot(t *testing.T) {
	// load archived evm outtx Gas
	// https://etherscan.io/tx/0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3
	chainID := int64(1)
	coinType := coin.CoinType_Gas
	outtxHash := "0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3"
	tx, receipt := testutils.LoadEVMOuttxNReceipt(t, TestDataDir, chainID, outtxHash, coinType)

	// load archived cctx
	cctx := testutils.LoadCctxByNonce(t, chainID, tx.Nonce())

	t.Run("outtx ballot should match cctx", func(t *testing.T) {
		msg := types.NewMsgVoteOnObservedOutboundTx(
			"anyCreator",
			cctx.Index,
			receipt.TxHash.Hex(),
			receipt.BlockNumber.Uint64(),
			receipt.GasUsed,
			math.NewIntFromBigInt(tx.GasPrice()),
			tx.Gas(),
			math.NewUintFromBigInt(tx.Value()),
			chains.ReceiveStatus_success,
			chainID,
			tx.Nonce(),
			coinType,
		)
		ballotExpected := cctx.GetCurrentOutTxParam().OutboundTxBallotIndex
		require.Equal(t, ballotExpected, msg.Digest())
	})
}
