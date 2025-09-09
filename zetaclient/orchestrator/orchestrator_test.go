package orchestrator

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"cosmossdk.io/math"
	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/testutil/sample"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
)

func TestOrchestratorV2(t *testing.T) {
	t.Run("updates app context", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// ACT #1
		// Start orchestrator
		err := ts.Start(ts.ctx)

		// Mimic zetacore update
		ts.MockChainParams(chains.Ethereum, mocks.MockChainParams(chains.Ethereum.ChainId, 100))

		// ASSERT #1
		require.NoError(t, err)

		// Check that eventually appContext would contain only desired chains
		check := func() bool {
			list := ts.appContext.ListChains()
			return len(list) == 1 && chainsContain(list, chains.Ethereum.ChainId)
		}

		assert.Eventually(t, check, 5*time.Second, 100*time.Millisecond)

		assert.Contains(t, ts.Log.String(), "chain list changed at the runtime")
		assert.Contains(t, ts.Log.String(), `"chains_new":[1]`)

		// ACT #2
		// Mimic zetacore update that adds bitcoin chain with chain params
		ts.MockChainParams(
			chains.Ethereum,
			mocks.MockChainParams(chains.Ethereum.ChainId, 100),
			chains.BitcoinMainnet,
			mocks.MockChainParams(chains.BitcoinMainnet.ChainId, 100),
		)

		check = func() bool {
			list := ts.appContext.ListChains()
			return len(list) == 2 && chainsContain(list, chains.Ethereum.ChainId, chains.BitcoinMainnet.ChainId)
		}

		assert.Eventually(t, check, 5*time.Second, 100*time.Millisecond)

		assert.Contains(t, ts.Log.String(), `"chains_new":[1,8332],"message":"chain list changed at the runtime"`)
	})
}

type testSuite struct {
	*Orchestrator
	*testlog.Log

	t *testing.T

	ctx        context.Context
	appContext *zctx.AppContext

	chains      []chains.Chain
	chainParams []*observertypes.ChainParams

	zetacore  *mocks.ZetacoreClient
	scheduler *scheduler.Scheduler
	tss       *mocks.TSS

	mu sync.Mutex
}

var defaultChainsWithParams = []any{
	chains.Ethereum,
	chains.BitcoinMainnet,
	chains.SolanaMainnet,
	chains.SuiMainnet,
	chains.TONMainnet,

	mocks.MockChainParams(chains.Ethereum.ChainId, 100),
	mocks.MockChainParams(chains.BitcoinMainnet.ChainId, 3),
	mocks.MockChainParams(chains.SolanaMainnet.ChainId, 10),
	mocks.MockChainParams(chains.SuiMainnet.ChainId, 1),
	mocks.MockChainParams(chains.TONMainnet.ChainId, 1),
}

func newTestSuite(t *testing.T) *testSuite {
	logger := testlog.New(t)
	baseLogger := base.Logger{
		Std:        logger.Logger,
		Compliance: logger.Logger,
	}

	chainList, chainParams := parseChainsWithParams(t, defaultChainsWithParams...)

	ctx, appCtx := newAppContext(t, logger.Logger, chainList, chainParams)

	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	// Services
	var (
		schedulerService = scheduler.New(logger.Logger, time.Second)
		zetacore         = mocks.NewZetacoreClient(t)
		tss              = mocks.NewTSS(t)
	)

	deps := &Dependencies{
		Zetacore:  zetacore,
		TSS:       tss,
		DBPath:    db.SqliteInMemory,
		Telemetry: metrics.NewTelemetryServer(),
	}

	v2, err := New(schedulerService, deps, baseLogger)
	require.NoError(t, err)

	ts := &testSuite{
		Orchestrator: v2,
		Log:          logger,

		t: t,

		ctx:        ctx,
		appContext: appCtx,

		chains:      chainList,
		chainParams: chainParams,

		scheduler: schedulerService,
		zetacore:  zetacore,
		tss:       tss,
	}

	// Mock basic zetacore methods
	on(zetacore, "GetBlockHeight", 1).Return(int64(123), nil).Maybe()
	on(zetacore, "GetUpgradePlan", 1).Return(nil, nil).Maybe()
	on(zetacore, "GetAdditionalChains", 1).Return(nil, nil).Maybe()
	on(zetacore, "GetCrosschainFlags", 1).Return(appCtx.GetCrossChainFlags(), nil).Maybe()
	on(zetacore, "GetOperationalFlags", 1).Return(appCtx.GetOperationalFlags(), nil).Maybe()
	on(zetacore, "GetZetaHotKeyBalance", 1).Return(math.NewInt(123), nil).Maybe()

	// Mock chain-related methods as dynamic getters
	on(zetacore, "GetSupportedChains", 1).Return(ts.getSupportedChains).Maybe()
	on(zetacore, "GetChainParams", 1).Return(ts.getChainParams).Maybe()

	// Mock zetacore blocks
	on(zetacore, "NewBlockSubscriber", 1).Return(ts.blockProducer).Maybe()

	// Mock CCTX-related calls (stubs for now)
	on(zetacore, "ListPendingCCTX", 2).Return(nil, uint64(0), nil).Maybe()
	on(zetacore, "GetInboundTrackersForChain", 2).Return(nil, nil).Maybe()
	on(zetacore, "GetAllOutboundTrackerByChain", 3).Return(nil, nil).Maybe()

	t.Cleanup(ts.Stop)

	return ts
}

func (ts *testSuite) HasObserverSigner(chainID int64) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	_, ok := ts.Orchestrator.chains[chainID]
	return ok
}

func (ts *testSuite) MockChainParams(newValues ...any) {
	chainList, chainParams := parseChainsWithParams(ts.t, newValues...)

	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.chains = chainList
	ts.chainParams = chainParams
}

// UpdateConfig updates "global" config.Config for test suite.
func (ts *testSuite) UpdateConfig(fn func(cfg *config.Config)) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	cfg := ts.appContext.Config()
	fn(&cfg)

	// The config is sealed i.e. we can't alter it after starting zetaclientd.
	// But for test purposes we use `reflect` to mimic
	// that it was set by the validator *before* starting the app.
	field := reflect.ValueOf(ts.appContext).Elem().FieldByName("config")
	ptr := unsafe.Pointer(field.UnsafeAddr())
	configPtr := (*config.Config)(ptr)

	*configPtr = cfg
}

func (ts *testSuite) getSupportedChains(_ context.Context) ([]chains.Chain, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.chains, nil
}

func (ts *testSuite) getChainParams(_ context.Context) ([]*observertypes.ChainParams, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.chainParams, nil
}

func (ts *testSuite) blockProducer(_ context.Context) (chan cometbfttypes.EventDataNewBlock, error) {
	closeCh := make(chan struct{})
	ts.t.Cleanup(func() { close(closeCh) })

	blocksChan := make(chan cometbfttypes.EventDataNewBlock)
	blockNumber := int64(100)

	go func() {
		for {
			block := cometbfttypes.EventDataNewBlock{
				Block: &cometbfttypes.Block{
					Header: cometbfttypes.Header{
						Height: atomic.AddInt64(&blockNumber, 1),
						Time:   time.Now().UTC(),
					},
				},
			}

			select {
			case blocksChan <- block:
				time.Sleep(200 * time.Millisecond)
			case <-closeCh:
				close(blocksChan)
				return
			}
		}
	}()

	return blocksChan, nil
}

func newAppContext(
	t *testing.T,
	logger zerolog.Logger,
	chainList []chains.Chain,
	chainParams []*observertypes.ChainParams,
) (context.Context, *zctx.AppContext) {
	// Mock config
	cfg := config.New(false)

	cfg.ConfigUpdateTicker = 1

	for _, c := range chainList {
		switch {
		case chains.IsEVMChain(c.ChainId, nil):
			cfg.EVMChainConfigs[c.ChainId] = config.EVMConfig{Endpoint: "localhost"}
		case chains.IsBitcoinChain(c.ChainId, nil):
			cfg.BTCChainConfigs[c.ChainId] = config.BTCConfig{RPCHost: "localhost"}
		case chains.IsSolanaChain(c.ChainId, nil):
			cfg.SolanaConfig = config.SolanaConfig{Endpoint: "localhost"}
		case chains.IsTONChain(c.ChainId, nil):
			cfg.TONConfig = config.TONConfig{Endpoint: "localhost"}
		case chains.IsSuiChain(c.ChainId, nil):
			cfg.SuiConfig = config.SuiConfig{Endpoint: "localhost"}
		default:
			t.Fatalf("create app context: unsupported chain %d", c.ChainId)
		}
	}

	// chain params
	params := map[int64]*observertypes.ChainParams{}
	for i := range chainParams {
		cp := chainParams[i]
		params[cp.ChainId] = cp
	}

	// new AppContext
	appContext := zctx.New(cfg, nil, logger)

	ccFlags := sample.CrosschainFlags()
	opFlags := sample.OperationalFlags()

	err := appContext.Update(chainList, nil, params, *ccFlags, opFlags)
	require.NoError(t, err, "failed to update app context")

	ctx := zctx.WithAppContext(context.Background(), appContext)

	return ctx, appContext
}

func parseChainsWithParams(t *testing.T, chainsOrParams ...any) ([]chains.Chain, []*observertypes.ChainParams) {
	var (
		supportedChains = make([]chains.Chain, 0, len(chainsOrParams))
		obsParams       = make([]*observertypes.ChainParams, 0, len(chainsOrParams))
	)

	for _, something := range chainsOrParams {
		switch tt := something.(type) {
		case *chains.Chain:
			supportedChains = append(supportedChains, *tt)
		case chains.Chain:
			supportedChains = append(supportedChains, tt)
		case *observertypes.ChainParams:
			obsParams = append(obsParams, tt)
		case observertypes.ChainParams:
			obsParams = append(obsParams, &tt)
		default:
			t.Fatalf("parse chains and params: unsupported type %T (%+v)", tt, tt)
		}
	}

	return supportedChains, obsParams
}

func chainsContain(list []zctx.Chain, ids ...int64) bool {
	set := make(map[int64]struct{}, len(list))
	for _, chain := range list {
		set[chain.ID()] = struct{}{}
	}

	for _, chainID := range ids {
		if _, found := set[chainID]; !found {
			return false
		}
	}

	return true
}

type mockOn interface {
	On(methodName string, arguments ...any) *mock.Call
}

// handy wrapper for concise calls
func on(m mockOn, method string, nArgs int) *mock.Call {
	args := make([]any, nArgs)
	for i := range nArgs {
		args[i] = mock.Anything
	}

	return m.On(method, args...)
}
