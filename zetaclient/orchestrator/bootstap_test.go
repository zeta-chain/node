package orchestrator

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/ptr"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testrpc"
)

const solanaGatewayAddress = "2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s"

func TestCreateSignerMap(t *testing.T) {
	var (
		ts         = metrics.NewTelemetryServer()
		tss        = mocks.NewTSSMainnet()
		log        = zerolog.New(zerolog.NewTestWriter(t))
		baseLogger = base.Logger{Std: log, Compliance: log}
	)

	t.Run("CreateSignerMap", func(t *testing.T) {
		// ARRANGE
		// Given a BTC server
		_, btcConfig := testrpc.NewBtcServer(t)

		// Given a zetaclient config with ETH, MATIC, and BTC chains
		cfg := config.New(false)

		cfg.EVMChainConfigs[chains.Ethereum.ChainId] = config.EVMConfig{
			Chain:    chains.Ethereum,
			Endpoint: testutils.MockEVMRPCEndpoint,
		}

		cfg.EVMChainConfigs[chains.Polygon.ChainId] = config.EVMConfig{
			Chain:    chains.Polygon,
			Endpoint: testutils.MockEVMRPCEndpoint,
		}

		cfg.BitcoinConfig = btcConfig

		// Given AppContext
		app := zctx.New(cfg, nil, log)
		ctx := zctx.WithAppContext(context.Background(), app)

		// Given chain & chainParams "fetched" from zetacore
		// (note that slice LACKS polygon chain on purpose)
		mustUpdateAppContextChainParams(t, app, []chains.Chain{
			chains.Ethereum, chains.BitcoinMainnet,
		})

		// ACT
		signers, err := CreateSignerMap(ctx, tss, baseLogger, ts)

		// ASSERT
		assert.NoError(t, err)
		assert.NotEmpty(t, signers)

		// Okay, now we want to check that signers for EVM and BTC were created
		assert.Equal(t, 2, len(signers))
		hasSigner(t, signers, chains.Ethereum.ChainId)
		hasSigner(t, signers, chains.BitcoinMainnet.ChainId)

		t.Run("Add polygon in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum, chains.BitcoinMainnet, chains.Polygon,
			})

			// ACT
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &signers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			hasSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Disable ethereum in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Polygon, chains.BitcoinMainnet,
			})

			// ACT
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &signers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 0, added)
			assert.Equal(t, 1, removed)

			missesSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			hasSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Re-enable ethereum in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum,
				chains.Polygon,
				chains.BitcoinMainnet,
			})

			// ACT
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &signers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			hasSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Disable btc in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum,
				chains.Polygon,
			})

			// ACT
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &signers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 0, added)
			assert.Equal(t, 1, removed)

			hasSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			missesSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Re-enable btc in the runtime", func(t *testing.T) {
			// ARRANGE
			// Given updated data from zetacore containing polygon chain
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum,
				chains.Polygon,
				chains.BitcoinMainnet,
			})

			// ACT
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &signers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			hasSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("No changes", func(t *testing.T) {
			// ARRANGE
			before := len(signers)

			// ACT
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, ts, &signers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 0, added)
			assert.Equal(t, 0, removed)
			assert.Equal(t, before, len(signers))
		})
	})
}

func TestCreateChainObserverMap(t *testing.T) {
	var (
		ts         = metrics.NewTelemetryServer()
		tss        = mocks.NewTSSMainnet()
		log        = zerolog.New(zerolog.NewTestWriter(t))
		baseLogger = base.Logger{Std: log, Compliance: log}
		client     = mocks.NewZetacoreClient(t)
		dbPath     = db.SqliteInMemory
	)

	t.Run("CreateChainObserverMap", func(t *testing.T) {
		// ARRANGE
		// Given a BTC server
		btcServer, btcConfig := testrpc.NewBtcServer(t)

		btcServer.SetBlockCount(123)

		// Given generic EVM RPC
		evmServer := testrpc.NewEVMServer(t)
		evmServer.SetBlockNumber(100)

		// Given generic SOL RPC
		_, solConfig := testrpc.NewSolanaServer(t)

		// Given a zetaclient config with ETH, MATIC, and BTC chains
		cfg := config.New(false)

		cfg.EVMChainConfigs[chains.Ethereum.ChainId] = config.EVMConfig{
			Chain:    chains.Ethereum,
			Endpoint: evmServer.Endpoint,
		}

		cfg.EVMChainConfigs[chains.Polygon.ChainId] = config.EVMConfig{
			Chain:    chains.Polygon,
			Endpoint: evmServer.Endpoint,
		}

		cfg.BitcoinConfig = btcConfig
		cfg.SolanaConfig = solConfig

		// Given AppContext
		app := zctx.New(cfg, nil, log)
		ctx := zctx.WithAppContext(context.Background(), app)

		// Given chain & chainParams "fetched" from zetacore
		// (note that slice LACKS polygon & SOL chains on purpose)
		mustUpdateAppContextChainParams(t, app, []chains.Chain{
			chains.Ethereum,
			chains.BitcoinMainnet,
		})

		// ACT
		observers, err := CreateChainObserverMap(ctx, client, tss, dbPath, baseLogger, ts)

		// ASSERT
		assert.NoError(t, err)
		assert.NotEmpty(t, observers)

		// Okay, now we want to check that signers for EVM and BTC were created
		assert.Equal(t, 2, len(observers))
		hasObserver(t, observers, chains.Ethereum.ChainId)
		hasObserver(t, observers, chains.BitcoinMainnet.ChainId)

		t.Run("Add polygon in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum, chains.BitcoinMainnet, chains.Polygon,
			})

			// ACT
			added, removed, err := syncObserverMap(ctx, client, tss, dbPath, baseLogger, ts, &observers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasObserver(t, observers, chains.Ethereum.ChainId)
			hasObserver(t, observers, chains.Polygon.ChainId)
			hasObserver(t, observers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Add solana in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum,
				chains.BitcoinMainnet,
				chains.Polygon,
				chains.SolanaMainnet,
			})

			// ACT
			added, removed, err := syncObserverMap(ctx, client, tss, dbPath, baseLogger, ts, &observers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasObserver(t, observers, chains.Ethereum.ChainId)
			hasObserver(t, observers, chains.Polygon.ChainId)
			hasObserver(t, observers, chains.BitcoinMainnet.ChainId)
			hasObserver(t, observers, chains.SolanaMainnet.ChainId)
		})

		t.Run("Disable ethereum and solana in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.BitcoinMainnet,
				chains.Polygon,
			})

			// ACT
			added, removed, err := syncObserverMap(ctx, client, tss, dbPath, baseLogger, ts, &observers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 0, added)
			assert.Equal(t, 2, removed)

			missesObserver(t, observers, chains.Ethereum.ChainId)
			hasObserver(t, observers, chains.Polygon.ChainId)
			hasObserver(t, observers, chains.BitcoinMainnet.ChainId)
			missesObserver(t, observers, chains.SolanaMainnet.ChainId)
		})

		t.Run("Re-enable ethereum in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum, chains.BitcoinMainnet, chains.Polygon,
			})

			// ACT
			added, removed, err := syncObserverMap(ctx, client, tss, dbPath, baseLogger, ts, &observers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasObserver(t, observers, chains.Ethereum.ChainId)
			hasObserver(t, observers, chains.Polygon.ChainId)
			hasObserver(t, observers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Disable btc in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum, chains.Polygon,
			})

			// ACT
			added, removed, err := syncObserverMap(ctx, client, tss, dbPath, baseLogger, ts, &observers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 0, added)
			assert.Equal(t, 1, removed)

			hasObserver(t, observers, chains.Ethereum.ChainId)
			hasObserver(t, observers, chains.Polygon.ChainId)
			missesObserver(t, observers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Re-enable btc in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.BitcoinMainnet, chains.Ethereum, chains.Polygon,
			})

			// ACT
			added, removed, err := syncObserverMap(ctx, client, tss, dbPath, baseLogger, ts, &observers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasObserver(t, observers, chains.Ethereum.ChainId)
			hasObserver(t, observers, chains.Polygon.ChainId)
			hasObserver(t, observers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("No changes", func(t *testing.T) {
			// ARRANGE
			before := len(observers)

			// ACT
			added, removed, err := syncObserverMap(ctx, client, tss, dbPath, baseLogger, ts, &observers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 0, added)
			assert.Equal(t, 0, removed)
			assert.Equal(t, before, len(observers))
		})
	})
}

func chainParams(supportedChains []chains.Chain) ([]chains.Chain, map[int64]*observertypes.ChainParams) {
	params := make(map[int64]*observertypes.ChainParams)

	for _, chain := range supportedChains {
		chainID := chain.ChainId
		if chains.IsBitcoinChain(chainID, nil) {
			p := mocks.MockChainParams(chainID, 100)
			params[chainID] = &p
			continue
		}

		if chains.IsSolanaChain(chainID, nil) {
			p := mocks.MockChainParams(chainID, 100)
			p.GatewayAddress = solanaGatewayAddress
			params[chainID] = &p
			continue
		}

		if chains.IsEVMChain(chainID, nil) {
			params[chainID] = ptr.Ptr(mocks.MockChainParams(chainID, 100))
			continue
		}
	}

	return supportedChains, params
}

func mustUpdateAppContextChainParams(t *testing.T, app *zctx.AppContext, chains []chains.Chain) {
	supportedChain, params := chainParams(chains)
	mustUpdateAppContext(t, app, supportedChain, nil, params)
}

func mustUpdateAppContext(
	t *testing.T,
	app *zctx.AppContext,
	chains, additionalChains []chains.Chain,
	chainParams map[int64]*observertypes.ChainParams,
) {
	err := app.Update(
		app.GetKeygen(),
		chains,
		additionalChains,
		chainParams,
		"tssPubKey",
		app.GetCrossChainFlags(),
	)

	require.NoError(t, err)
}

func hasSigner(t *testing.T, signers map[int64]interfaces.ChainSigner, chainId int64) {
	signer, ok := signers[chainId]
	assert.True(t, ok, "missing signer for chain %d", chainId)
	assert.NotEmpty(t, signer)
}

func missesSigner(t *testing.T, signers map[int64]interfaces.ChainSigner, chainId int64) {
	_, ok := signers[chainId]
	assert.False(t, ok, "unexpected signer for chain %d", chainId)
}

func hasObserver(t *testing.T, observer map[int64]interfaces.ChainObserver, chainId int64) {
	signer, ok := observer[chainId]
	assert.True(t, ok, "missing observer for chain %d", chainId)
	assert.NotEmpty(t, signer)
}

func missesObserver(t *testing.T, observer map[int64]interfaces.ChainObserver, chainId int64) {
	_, ok := observer[chainId]
	assert.False(t, ok, "unexpected observer for chain %d", chainId)
}
