package observer_test

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	"cosmossdk.io/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru"
	"github.com/onrik/ethrpc"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/keys"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

// the relative path to the testdata directory
var TestDataDir = "../../../"

// getZetacoreContext creates a zetacore context for unit tests
func getZetacoreContext(
	evmChain chains.Chain,
	endpoint string,
	evmChainParams *observertypes.ChainParams,
) (*zctx.AppContext, config.EVMConfig) {
	// use default endpoint if not provided
	if endpoint == "" {
		endpoint = "http://localhost:8545"
	}

	// create config
	cfg := config.New(false)
	cfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
		Chain:    evmChain,
		Endpoint: endpoint,
	}

	// create zetacore context
	appContext := zctx.New(cfg, zerolog.Nop())
	evmChainParamsMap := make(map[int64]*observertypes.ChainParams)
	evmChainParamsMap[evmChain.ChainId] = evmChainParams

	// feed chain params
	appContext.Update(
		&observertypes.Keygen{},
		[]chains.Chain{evmChain},
		evmChainParamsMap,
		nil,
		"",
		*sample.CrosschainFlags(),
		[]chains.Chain{},
		sample.HeaderSupportedChains(),
		true,
	)
	// create app context
	return appContext, cfg.EVMChainConfigs[evmChain.ChainId]
}

// MockEVMObserver creates a mock ChainObserver with custom chain, TSS, params etc
func MockEVMObserver(
	t *testing.T,
	chain chains.Chain,
	evmClient interfaces.EVMRPCClient,
	evmJSONRPC interfaces.EVMJSONRPCClient,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	dbpath string,
	lastBlock uint64,
	params observertypes.ChainParams,
) *observer.Observer {
	ctx := context.Background()

	// use default mock evm client if not provided
	if evmClient == nil {
		evmClient = mocks.NewMockEvmClient().WithBlockNumber(1000)
	}

	// use default mock zetacore client if not provided
	if zetacoreClient == nil {
		zetacoreClient = mocks.NewZetacoreClient(t).
			WithKeys(&keys.Keys{}).
			WithZetaChain().
			WithPostVoteInbound("", "").
			WithPostVoteOutbound("", "")
	}
	// use default mock tss if not provided
	if tss == nil {
		tss = mocks.NewTSSMainnet()
	}
	// create zetacore context
	_, evmCfg := getZetacoreContext(chain, "", &params)

	// create observer
	ob, err := observer.NewObserver(ctx, evmCfg, evmClient, params, zetacoreClient, tss, dbpath, base.Logger{}, nil)
	require.NoError(t, err)
	ob.WithEvmJSONRPC(evmJSONRPC)
	ob.WithLastBlock(lastBlock)

	return ob
}

func Test_NewObserver(t *testing.T) {
	ctx := context.Background()

	// use Ethereum chain for testing
	chain := chains.Ethereum
	params := mocks.MockChainParams(chain.ChainId, 10)

	// test cases
	tests := []struct {
		name        string
		evmCfg      config.EVMConfig
		chainParams observertypes.ChainParams
		evmClient   interfaces.EVMRPCClient
		tss         interfaces.TSSSigner
		dbpath      string
		logger      base.Logger
		ts          *metrics.TelemetryServer
		fail        bool
		message     string
	}{
		{
			name: "should be able to create observer",
			evmCfg: config.EVMConfig{
				Chain:    chain,
				Endpoint: "http://localhost:8545",
			},
			chainParams: params,
			evmClient:   mocks.NewMockEvmClient().WithBlockNumber(1000),
			tss:         mocks.NewTSSMainnet(),
			dbpath:      sample.CreateTempDir(t),
			logger:      base.Logger{},
			ts:          nil,
			fail:        false,
		},
		{
			name: "should fail on invalid dbpath",
			evmCfg: config.EVMConfig{
				Chain:    chain,
				Endpoint: "http://localhost:8545",
			},
			chainParams: params,
			evmClient:   mocks.NewMockEvmClient().WithBlockNumber(1000),
			tss:         mocks.NewTSSMainnet(),
			dbpath:      "/invalid/dbpath", // invalid dbpath
			logger:      base.Logger{},
			ts:          nil,
			fail:        true,
			message:     "error creating db path",
		},
		{
			name: "should fail if RPC call fails",
			evmCfg: config.EVMConfig{
				Chain:    chain,
				Endpoint: "http://localhost:8545",
			},
			chainParams: params,
			evmClient:   mocks.NewMockEvmClient().WithError(fmt.Errorf("error RPC")),
			tss:         mocks.NewTSSMainnet(),
			dbpath:      sample.CreateTempDir(t),
			logger:      base.Logger{},
			ts:          nil,
			fail:        true,
			message:     "error RPC",
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create zetacore context, client and tss
			//zetacoreCtx, _ := getZetacoreContext(tt.evmCfg.Chain, tt.evmCfg.Endpoint, &params)
			zetacoreClient := mocks.NewZetacoreClient(t)

			// create observer
			ob, err := observer.NewObserver(
				ctx,
				tt.evmCfg,
				tt.evmClient,
				tt.chainParams,
				zetacoreClient,
				tt.tss,
				tt.dbpath,
				tt.logger,
				tt.ts,
			)

			// check result
			if tt.fail {
				require.ErrorContains(t, err, tt.message)
				require.Nil(t, ob)
			} else {
				require.NoError(t, err)
				require.NotNil(t, ob)
			}
		})
	}
}

func Test_LoadDB(t *testing.T) {
	ctx := context.Background()

	// use Ethereum chain for testing
	chain := chains.Ethereum
	params := mocks.MockChainParams(chain.ChainId, 10)
	dbpath := sample.CreateTempDir(t)
	ob := MockEVMObserver(t, chain, nil, nil, nil, nil, dbpath, 1, params)

	t.Run("should load db successfully", func(t *testing.T) {
		err := ob.LoadDB(ctx, dbpath)
		require.NoError(t, err)
		require.EqualValues(t, 1000, ob.LastBlockScanned())
	})
	t.Run("should fail on invalid dbpath", func(t *testing.T) {
		// load db with empty dbpath
		err := ob.LoadDB(ctx, "")
		require.ErrorContains(t, err, "empty db path")

		// load db with invalid dbpath
		err = ob.LoadDB(ctx, "/invalid/dbpath")
		require.ErrorContains(t, err, "error OpenDB")
	})
	t.Run("should fail on invalid env var", func(t *testing.T) {
		// set invalid environment variable
		envvar := base.EnvVarLatestBlockByChain(chain)
		os.Setenv(envvar, "invalid")
		defer os.Unsetenv(envvar)

		// load db
		err := ob.LoadDB(ctx, dbpath)
		require.ErrorContains(t, err, "error LoadLastBlockScanned")
	})
	t.Run("should fail on RPC error", func(t *testing.T) {
		// create observer
		tempClient := mocks.NewMockEvmClient()
		ob := MockEVMObserver(t, chain, tempClient, nil, nil, nil, dbpath, 1, params)

		// set RPC error
		tempClient.WithError(fmt.Errorf("error RPC"))

		// load db
		err := ob.LoadDB(ctx, dbpath)
		require.ErrorContains(t, err, "error RPC")
	})
}

func Test_LoadLastBlockScanned(t *testing.T) {
	ctx := context.Background()

	// use Ethereum chain for testing
	chain := chains.Ethereum
	params := mocks.MockChainParams(chain.ChainId, 10)

	// create observer using mock evm client
	evmClient := mocks.NewMockEvmClient().WithBlockNumber(100)
	dbpath := sample.CreateTempDir(t)
	ob := MockEVMObserver(t, chain, evmClient, nil, nil, nil, dbpath, 1, params)

	t.Run("should load last block scanned", func(t *testing.T) {
		// create db and write 123 as last block scanned
		ob.WriteLastBlockScannedToDB(123)

		// load last block scanned
		err := ob.LoadLastBlockScanned(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 123, ob.LastBlockScanned())
	})
	t.Run("should fail on invalid env var", func(t *testing.T) {
		// set invalid environment variable
		envvar := base.EnvVarLatestBlockByChain(chain)
		os.Setenv(envvar, "invalid")
		defer os.Unsetenv(envvar)

		// load last block scanned
		err := ob.LoadLastBlockScanned(ctx)
		require.ErrorContains(t, err, "error LoadLastBlockScanned")
	})
	t.Run("should fail on RPC error", func(t *testing.T) {
		// create observer on separate path, as we need to reset last block scanned
		otherPath := sample.CreateTempDir(t)
		obOther := MockEVMObserver(t, chain, evmClient, nil, nil, nil, otherPath, 1, params)

		// reset last block scanned to 0 so that it will be loaded from RPC
		obOther.WithLastBlockScanned(0)

		// set RPC error
		evmClient.WithError(fmt.Errorf("error RPC"))

		// load last block scanned
		err := obOther.LoadLastBlockScanned(ctx)
		require.ErrorContains(t, err, "error RPC")
	})
}

func Test_BlockCache(t *testing.T) {
	t.Run("should get block from cache", func(t *testing.T) {
		// create observer
		ob := &observer.Observer{}
		blockCache, err := lru.New(100)
		require.NoError(t, err)
		ob.WithBlockCache(blockCache)

		// create mock evm client
		JSONRPC := mocks.NewMockJSONRPCClient()
		ob.WithEvmJSONRPC(JSONRPC)

		// feed block to JSON rpc client
		block := &ethrpc.Block{Number: 100}
		JSONRPC.WithBlock(block)

		// get block header from observer, fallback to JSON RPC
		result, err := ob.GetBlockByNumberCached(uint64(100))
		require.NoError(t, err)
		require.EqualValues(t, block, result)

		// get block header from cache
		result, err = ob.GetBlockByNumberCached(uint64(100))
		require.NoError(t, err)
		require.EqualValues(t, block, result)
	})
	t.Run("should fail if stored type is not block", func(t *testing.T) {
		// create observer
		ob := &observer.Observer{}
		blockCache, err := lru.New(100)
		require.NoError(t, err)
		ob.WithBlockCache(blockCache)

		// add a string to cache
		blockNumber := uint64(100)
		blockCache.Add(blockNumber, "a string value")

		// get result header from cache
		result, err := ob.GetBlockByNumberCached(blockNumber)
		require.ErrorContains(t, err, "cached value is not of type *ethrpc.Block")
		require.Nil(t, result)
	})
	t.Run("should be able to remove block from cache", func(t *testing.T) {
		// create observer
		ob := &observer.Observer{}
		blockCache, err := lru.New(100)
		require.NoError(t, err)
		ob.WithBlockCache(blockCache)

		// delete non-existing block should not panic
		blockNumber := uint64(123)
		ob.RemoveCachedBlock(blockNumber)

		// add a block
		block := &ethrpc.Block{Number: 123}
		blockCache.Add(blockNumber, block)
		ob.WithBlockCache(blockCache)

		// block should be in cache
		result, err := ob.GetBlockByNumberCached(blockNumber)
		require.NoError(t, err)
		require.EqualValues(t, block, result)

		// delete the block should not panic
		ob.RemoveCachedBlock(blockNumber)
	})
}

func Test_HeaderCache(t *testing.T) {
	ctx := context.Background()

	t.Run("should get block header from cache", func(t *testing.T) {
		// create observer
		ob := &observer.Observer{}
		headerCache, err := lru.New(100)
		require.NoError(t, err)
		ob.WithHeaderCache(headerCache)

		// create mock evm client
		evmClient := mocks.NewMockEvmClient()
		ob.WithEvmClient(evmClient)

		// feed block header to evm client
		header := &ethtypes.Header{Number: big.NewInt(100)}
		evmClient.WithHeader(header)

		// get block header from observer
		resHeader, err := ob.GetBlockHeaderCached(ctx, uint64(100))
		require.NoError(t, err)
		require.EqualValues(t, header, resHeader)
	})
	t.Run("should fail if stored type is not block header", func(t *testing.T) {
		// create observer
		ob := &observer.Observer{}
		headerCache, err := lru.New(100)
		require.NoError(t, err)
		ob.WithHeaderCache(headerCache)

		// add a string to cache
		blockNumber := uint64(100)
		headerCache.Add(blockNumber, "a string value")

		// get block header from cache
		header, err := ob.GetBlockHeaderCached(ctx, blockNumber)
		require.ErrorContains(t, err, "cached value is not of type *ethtypes.Header")
		require.Nil(t, header)
	})
}

func Test_CheckTxInclusion(t *testing.T) {
	// load archived evm outbound Gas
	// https://etherscan.io/tx/0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3
	chainID := int64(1)
	coinType := coin.CoinType_Gas
	outboundHash := "0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3"
	tx, receipt := testutils.LoadEVMOutboundNReceipt(t, TestDataDir, chainID, outboundHash, coinType)

	// load archived evm block
	// https://etherscan.io/block/19363323
	blockNumber := receipt.BlockNumber.Uint64()
	block := testutils.LoadEVMBlock(t, TestDataDir, chainID, blockNumber, true)

	// create client
	blockCache, err := lru.New(1000)
	require.NoError(t, err)
	ob := &observer.Observer{}

	// save block to cache
	blockCache.Add(blockNumber, block)
	ob.WithBlockCache(blockCache)

	t.Run("should pass for archived outbound", func(t *testing.T) {
		err := ob.CheckTxInclusion(tx, receipt)
		require.NoError(t, err)
	})
	t.Run("should fail on tx index out of range", func(t *testing.T) {
		// modify tx index to invalid number
		copyReceipt := *receipt
		// #nosec G115 non negative value
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
	// load archived evm outbound Gas
	// https://etherscan.io/tx/0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3
	chainID := int64(1)
	coinType := coin.CoinType_Gas
	outboundHash := "0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3"
	tx, receipt := testutils.LoadEVMOutboundNReceipt(t, TestDataDir, chainID, outboundHash, coinType)

	// load archived cctx
	cctx := testutils.LoadCctxByNonce(t, chainID, tx.Nonce())
	t.Run("outbound ballot should match cctx", func(t *testing.T) {
		msg := types.NewMsgVoteOutbound(
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
		ballotExpected := cctx.GetCurrentOutboundParam().BallotIndex
		require.Equal(t, ballotExpected, msg.Digest())
	})
}
