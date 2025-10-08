package observer

import (
	"context"
	"fmt"
	"os"
	"testing"

	"cosmossdk.io/math"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/ptr"
	"github.com/zeta-chain/node/zetaclient/chains/evm/client"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/mode"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testrpc"
)

// the relative path to the testdata directory
var TestDataDir = "../../../"

// getAppContext creates an AppContext for unit tests
func getAppContext(
	t *testing.T,
	evmChain chains.Chain,
	endpoint string,
	evmChainParams *observertypes.ChainParams,
) (*zctx.AppContext, config.EVMConfig) {
	// use default endpoint if not provided
	if endpoint == "" {
		endpoint = "http://localhost:8545"
	}

	require.Equal(t, evmChain.ChainId, evmChainParams.ChainId, "chain id mismatch between chain and params")

	// create config
	cfg := config.New(false)
	cfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
		Endpoint: endpoint,
	}

	logger := zerolog.New(zerolog.NewTestWriter(t))

	// create AppContext
	appContext := zctx.New(cfg, nil, logger)
	chainParams := map[int64]*observertypes.ChainParams{
		evmChain.ChainId: evmChainParams,
		chains.ZetaChainMainnet.ChainId: ptr.Ptr(
			mocks.MockChainParams(chains.ZetaChainMainnet.ChainId, 10),
		),
	}

	// feed chain params
	err := appContext.Update(
		[]chains.Chain{evmChain, chains.ZetaChainMainnet},
		nil,
		chainParams,
		*sample.CrosschainFlags(),
		sample.OperationalFlags(),
	)
	require.NoError(t, err)

	// create AppContext
	return appContext, cfg.EVMChainConfigs[evmChain.ChainId]
}

func Test_NewObserver(t *testing.T) {
	ctx := context.Background()

	// use Ethereum chain for testing
	chain := chains.Ethereum
	params := mocks.MockChainParams(chain.ChainId, 10)

	// create evm client with mocked block number 1000
	evmServer := testrpc.NewEVMServer(t)
	evmServer.SetChainID(int(chain.ChainId))
	evmServer.SetBlockNumber(1000)

	evmClient, err := client.NewFromEndpoint(ctx, evmServer.Endpoint)
	require.NoError(t, err)

	// test cases
	tests := []struct {
		name        string
		evmCfg      config.EVMConfig
		chainParams observertypes.ChainParams
		evmClient   *client.Client
		tssSigner   interfaces.TSSSigner
		logger      base.Logger
		before      func()
		after       func()
		ts          *metrics.TelemetryServer
		fail        bool
		message     string
	}{
		{
			name: "should be able to create observer",
			evmCfg: config.EVMConfig{
				Endpoint: "http://localhost:8545",
			},
			chainParams: params,
			evmClient:   evmClient,
			tssSigner:   mocks.NewTSS(t),
			logger:      base.Logger{},
			ts:          nil,
			fail:        false,
		},
		{
			name: "should fail if RPC call fails",
			evmCfg: config.EVMConfig{
				Endpoint: "http://localhost:8545",
			},
			chainParams: params,
			evmClient: func() *client.Client {
				// create mock evm client with RPC error
				evmServer := testrpc.NewEVMServer(t)
				evmServer.SetChainID(int(chain.ChainId))
				evmServer.SetBlockNumberFailure(fmt.Errorf("error RPC"))

				c, err := client.NewFromEndpoint(ctx, evmServer.Endpoint)
				require.NoError(t, err)

				return c
			}(),
			tssSigner: mocks.NewTSS(t),
			logger:    base.Logger{},
			ts:        nil,
			fail:      true,
			message:   "json-rpc error",
		},
		{
			name: "should fail on invalid ENV var",
			evmCfg: config.EVMConfig{
				Endpoint: "http://localhost:8545",
			},
			chainParams: params,
			evmClient:   evmClient,
			tssSigner:   mocks.NewTSS(t),
			before: func() {
				envVar := base.EnvVarLatestBlockByChain(chain)
				os.Setenv(envVar, "invalid")
			},
			after: func() {
				envVar := base.EnvVarLatestBlockByChain(chain)
				os.Unsetenv(envVar)
			},
			logger:  base.Logger{},
			ts:      nil,
			fail:    true,
			message: "unable to load last block scanned",
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create mock zetacore client
			zetacoreClient := mocks.NewZetacoreClient(t)

			database, err := db.NewFromSqliteInMemory(true)
			require.NoError(t, err)

			if tt.before != nil {
				tt.before()
			}
			if tt.after != nil {
				defer tt.after()
			}

			// create observer
			baseObserver, err := base.NewObserver(
				chain,
				tt.chainParams,
				zrepo.New(zetacoreClient, chain, mode.StandardMode),
				tt.tssSigner,
				1000,
				tt.ts,
				database,
				tt.logger,
			)
			require.NoError(t, err)
			ob, err := New(baseObserver, tt.evmClient)

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

func Test_LoadLastBlockScanned(t *testing.T) {
	ctx := context.Background()

	// create observer using mock evm client
	ob := newTestSuite(t)

	t.Run("should load last block scanned", func(t *testing.T) {
		// create db and write 123 as last block scanned
		ob.WriteLastBlockScannedToDB(123)

		// load last block scanned
		err := ob.loadLastBlockScanned(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 123, ob.LastBlockScanned())
	})
	t.Run("should fail on invalid env var", func(t *testing.T) {
		// set invalid environment variable
		envvar := base.EnvVarLatestBlockByChain(ob.Chain())
		os.Setenv(envvar, "invalid")
		defer os.Unsetenv(envvar)

		// load last block scanned
		err := ob.loadLastBlockScanned(ctx)
		require.ErrorContains(t, err, "error LoadLastBlockScanned")
	})
	t.Run("should fail on RPC error", func(t *testing.T) {
		// create observer on separate path, as we need to reset last block scanned
		obOther := newTestSuite(t)

		// reset last block scanned to 0 so that it will be loaded from RPC
		obOther.WithLastBlockScanned(0)

		// attach mock evm client to observer
		obOther.evmMock.On("BlockNumber", mock.Anything).Unset()
		obOther.evmMock.On("BlockNumber", mock.Anything).Return(uint64(0), fmt.Errorf("error RPC"))

		// load last block scanned
		err := obOther.loadLastBlockScanned(ctx)
		require.ErrorContains(t, err, "error RPC")
	})
}

func Test_BlockCache(t *testing.T) {
	t.Run("should get block from cache", func(t *testing.T) {
		// create observer
		ts := newTestSuite(t)

		// feed block to JSON rpc client
		block := &client.Block{Number: 100}
		ts.evmMock.On("BlockByNumberCustom", mock.Anything, mock.Anything).Return(block, nil)

		// get block header from observer, fallback to JSON RPC
		result, err := ts.Observer.GetBlockByNumberCached(ts.ctx, uint64(100))
		require.NoError(t, err)
		require.EqualValues(t, block, result)

		// get block header from cache
		result, err = ts.Observer.GetBlockByNumberCached(ts.ctx, uint64(100))
		require.NoError(t, err)
		require.EqualValues(t, block, result)
	})
	t.Run("should fail if stored type is not block", func(t *testing.T) {
		// create observer
		ts := newTestSuite(t)

		// add a string to cache
		blockNumber := uint64(100)
		ts.BlockCache().Add(blockNumber, "a string value")

		// get result header from cache
		result, err := ts.Observer.GetBlockByNumberCached(ts.ctx, blockNumber)
		require.ErrorContains(t, err, "cached value is not of type *client.Block")
		require.Nil(t, result)
	})
	t.Run("should be able to remove block from cache", func(t *testing.T) {
		// create observer
		ts := newTestSuite(t)

		// delete non-existing block should not panic
		blockNumber := uint64(123)
		ts.removeCachedBlock(blockNumber)

		// add a block
		block := &client.Block{Number: 123}
		ts.BlockCache().Add(blockNumber, block)

		// block should be in cache
		result, err := ts.GetBlockByNumberCached(ts.ctx, blockNumber)
		require.NoError(t, err)
		require.EqualValues(t, block, result)

		// delete the block should not panic
		ts.removeCachedBlock(blockNumber)
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

	// create observer
	ts := newTestSuite(t)

	// save block to cache
	ts.BlockCache().Add(blockNumber, block)

	t.Run("should pass for archived outbound", func(t *testing.T) {
		err := ts.checkTxInclusion(ts.ctx, tx, receipt)
		require.NoError(t, err)
	})
	t.Run("should fail on tx index out of range", func(t *testing.T) {
		// modify tx index to invalid number
		copyReceipt := *receipt
		// #nosec G115 non negative value
		copyReceipt.TransactionIndex = uint(len(block.Transactions))
		err := ts.checkTxInclusion(ts.ctx, tx, &copyReceipt)
		require.ErrorContains(t, err, "out of range")
	})
	t.Run("should fail on tx hash mismatch", func(t *testing.T) {
		// change the tx at position 'receipt.TransactionIndex' to a different tx
		priorTx := block.Transactions[receipt.TransactionIndex-1]
		block.Transactions[receipt.TransactionIndex] = priorTx
		ts.BlockCache().Add(blockNumber, block)

		// check inclusion should fail
		err := ts.checkTxInclusion(ts.ctx, tx, receipt)
		require.ErrorContains(t, err, "has different hash")

		// wrong block should be removed from cache
		_, ok := ts.BlockCache().Get(blockNumber)
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
			crosschaintypes.ConfirmationMode_SAFE,
		)
		ballotExpected := cctx.GetCurrentOutboundParam().BallotIndex
		require.Equal(t, ballotExpected, msg.Digest())
	})
}

type testSuite struct {
	*Observer
	ctx         context.Context
	appContext  *zctx.AppContext
	chainParams *observertypes.ChainParams
	tss         *mocks.TSS
	zetacore    *mocks.ZetacoreClient
	evmMock     *mocks.EVMClient
}

type testSuiteConfig struct {
	chain *chains.Chain
}

func newTestSuite(t *testing.T, opts ...func(*testSuiteConfig)) *testSuite {
	var cfg testSuiteConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	chain := chains.Ethereum
	if cfg.chain != nil {
		chain = *cfg.chain
	}

	chainParams := mocks.MockChainParams(chain.ChainId, 10)

	appContext, _ := getAppContext(t, chain, "", &chainParams)
	ctx := zctx.WithAppContext(context.Background(), appContext)

	evmMock := mocks.NewEVMClient(t)
	evmMock.On("BlockNumber", mock.Anything).Return(uint64(1000), nil).Maybe()

	zetacore := mocks.NewZetacoreClient(t).
		WithKeys(&keys.Keys{}).
		WithZetaChain().
		WithPostVoteInbound("", "").
		WithPostVoteOutbound("", "")

	tss := mocks.NewTSS(t).FakePubKey(testutils.TSSPubKeyMainnet)

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	log := zerolog.New(zerolog.NewTestWriter(t)).With().Caller().Logger()
	logger := base.Logger{Std: log, Compliance: log}

	baseObserver, err := base.NewObserver(chain, chainParams,
		zrepo.New(zetacore, chain, mode.StandardMode), tss, 1000, nil, database, logger)
	require.NoError(t, err)

	ob, err := New(baseObserver, evmMock)
	require.NoError(t, err)
	ob.WithLastBlock(1)

	return &testSuite{
		Observer:    ob,
		ctx:         ctx,
		appContext:  appContext,
		chainParams: &chainParams,
		tss:         tss,
		zetacore:    zetacore,
		evmMock:     evmMock,
	}
}
