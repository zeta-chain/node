package orchestrator

import (
	"bytes"
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/testutil/sample"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// [x] todo: fix UNEXPECTED test failure (exit code 1)
// [ ] log helper for debugging (zerolog)
// [ ] todo: add v2 to start.go
// [ ] todo: run e2e tests.
// [ ] todo: add btc observer&signer

func TestOrchestratorV2(t *testing.T) {
	t.Run("updates app context", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// ACT
		err := ts.Start(ts.ctx)

		// Mimic zetacore update
		ts.MockChainParams(chains.Ethereum, mocks.MockChainParams(chains.Ethereum.ChainId, 100))

		// ASSERT
		require.NoError(t, err)

		// Check that eventually appContext would contain only desired chains
		check := func() bool {
			listChains := ts.appContext.ListChains()
			if len(listChains) != 1 {
				return false
			}

			return listChains[0].ID() == chains.Ethereum.ChainId
		}

		assert.Eventually(t, check, 5*time.Second, 100*time.Millisecond)

		assert.Contains(t, ts.logBuffer.String(), "Chain list changed at the runtime!")
		assert.Contains(t, ts.logBuffer.String(), `"chains.new":[1]`)
	})
}

type testSuite struct {
	*V2

	t *testing.T

	ctx        context.Context
	appContext *zctx.AppContext

	chains      []chains.Chain
	chainParams []*observertypes.ChainParams

	zetacore  *mocks.ZetacoreClient
	scheduler *scheduler.Scheduler

	logBuffer *bytes.Buffer
	logger    zerolog.Logger

	mu sync.Mutex
}

var defaultChainsWithParams = []any{
	chains.Ethereum,
	chains.BitcoinMainnet,
	chains.SolanaMainnet,
	chains.TONMainnet,

	mocks.MockChainParams(chains.Ethereum.ChainId, 100),
	mocks.MockChainParams(chains.BitcoinMainnet.ChainId, 3),
	mocks.MockChainParams(chains.SolanaMainnet.ChainId, 10),
	mocks.MockChainParams(chains.TONMainnet.ChainId, 1),
}

func newTestSuite(t *testing.T) *testSuite {
	var (
		logBuffer = &bytes.Buffer{}
		logger    = zerolog.New(io.MultiWriter(zerolog.NewTestWriter(t), logBuffer))

		chainList, chainParams = parseChainsWithParams(t, defaultChainsWithParams...)
		ctx, appCtx            = newAppContext(t, logger, chainList, chainParams)

		schedulerService = scheduler.New(logger)
		zetacore         = mocks.NewZetacoreClient(t)
	)

	ts := &testSuite{
		V2: NewV2(zetacore, schedulerService, logger),

		t: t,

		ctx:        ctx,
		appContext: appCtx,

		chains:      chainList,
		chainParams: chainParams,

		scheduler: schedulerService,
		zetacore:  zetacore,

		logBuffer: logBuffer,
		logger:    logger,
	}

	// Mock basic zetacore methods
	zetacore.On("GetBlockHeight", mock.Anything).Return(int64(123), nil).Maybe()
	zetacore.On("GetUpgradePlan", mock.Anything).Return(nil, nil).Maybe()
	zetacore.On("GetAdditionalChains", mock.Anything).Return(nil, nil).Maybe()
	zetacore.On("GetCrosschainFlags", mock.Anything).Return(appCtx.GetCrossChainFlags(), nil).Maybe()

	//Mock chain-related methods as dynamic getters
	zetacore.On("GetSupportedChains", mock.Anything).Return(ts.getSupportedChains).Maybe()
	zetacore.On("GetChainParams", mock.Anything).Return(ts.getChainParams).Maybe()

	t.Cleanup(ts.Stop)

	return ts
}

func (ts *testSuite) MockChainParams(newValues ...any) {
	chainList, chainParams := parseChainsWithParams(ts.t, newValues...)

	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.chains = chainList
	ts.chainParams = chainParams
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
			cfg.TONConfig = config.TONConfig{LiteClientConfigURL: "localhost"}
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

	err := appContext.Update(chainList, nil, params, *ccFlags)
	require.NoError(t, err, "failed to update app context")

	ctx := zctx.WithAppContext(context.Background(), appContext)

	return ctx, appContext
}
