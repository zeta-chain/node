package orchestrator

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testrpc"
)

const (
	solanaGatewayAddress = "2kJndCL9NBR36ySiQ4bmArs4YgWQu67LmCDfLzk5Gb7s"
	tonGatewayAddress    = "0:997d889c815aeac21c47f86ae0e38383efc3c3463067582f6263ad48c5a1485b"
	tonMainnet           = "https://ton.org/global-config.json"
)

func TestCreateSignerMap(t *testing.T) {
	var (
		tss        = mocks.NewTSS(t)
		log        = zerolog.New(zerolog.NewTestWriter(t))
		baseLogger = base.Logger{Std: log, Compliance: log}
	)

	t.Run("CreateSignerMap", func(t *testing.T) {
		// ARRANGE
		// Given a zetaclient config with ETH, MATIC, and BTC chains
		cfg := config.New(false)

		// Given AppContext
		app := zctx.New(cfg, nil, log)
		ctx := zctx.WithAppContext(context.Background(), app)

		// Given chain & chainParams "fetched" from zetacore
		mustUpdateAppContextChainParams(t, app, []chains.Chain{
			chains.SolanaMainnet,
		})

		// ACT
		signers, err := CreateSignerMap(ctx, tss, baseLogger)

		// ASSERT
		assert.NoError(t, err)
		assert.NotEmpty(t, signers)

		// Okay, now we want to check that signer for EVM was created
		assert.Equal(t, 1, len(signers))
		hasSigner(t, signers, chains.Ethereum.ChainId)

		t.Run("Add polygon in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum, chains.Polygon,
			})

			// ACT
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, &signers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
		})

		t.Run("Disable ethereum in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Polygon, chains.BitcoinMainnet,
			})

			// ACT
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, &signers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 0, added)
			assert.Equal(t, 1, removed)

			missesSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
			missesSigner(t, signers, chains.BitcoinMainnet.ChainId)
		})

		t.Run("Re-enable ethereum in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum,
				chains.Polygon,
			})

			// ACT
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, &signers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 0, removed)

			hasSigner(t, signers, chains.Ethereum.ChainId)
			hasSigner(t, signers, chains.Polygon.ChainId)
		})

		t.Run("No changes", func(t *testing.T) {
			// ARRANGE
			before := len(signers)

			// ACT
			added, removed, err := syncSignerMap(ctx, tss, baseLogger, &signers)

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
		tss        = mocks.NewTSS(t)
		log        = zerolog.New(zerolog.NewTestWriter(t))
		baseLogger = base.Logger{Std: log, Compliance: log}
		client     = mocks.NewZetacoreClient(t)
		dbPath     = db.SqliteInMemory
	)

	mockZetacore(client)

	t.Run("CreateChainObserverMap", func(t *testing.T) {
		// ARRANGE
		// Given generic EVM RPC
		evmServer := testrpc.NewEVMServer(t)
		evmServer.SetBlockNumber(100)

		// Given SOL config
		_, solConfig := testrpc.NewSolanaServer(t)

		// Given TON config
		tonConfig := config.TONConfig{LiteClientConfigURL: tonMainnet}

		// Given a zetaclient config with ETH, MATIC, and BTC chains
		cfg := config.New(false)

		cfg.EVMChainConfigs[chains.Ethereum.ChainId] = config.EVMConfig{
			Endpoint: evmServer.Endpoint,
		}

		cfg.EVMChainConfigs[chains.Polygon.ChainId] = config.EVMConfig{
			Endpoint: evmServer.Endpoint,
		}

		cfg.SolanaConfig = solConfig
		cfg.TONConfig = tonConfig

		// Given AppContext
		app := zctx.New(cfg, nil, log)
		ctx := zctx.WithAppContext(context.Background(), app)

		// Given chain & chainParams "fetched" from zetacore
		// note that slice LACKS polygon & SOL chains on purpose
		// also note that BTC is handled by orchestrator v2
		mustUpdateAppContextChainParams(t, app, []chains.Chain{
			chains.Ethereum,
			chains.TONMainnet,
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
		hasObserver(t, observers, chains.TONMainnet.ChainId)
		missesObserver(t, observers, chains.BitcoinMainnet.ChainId)

		t.Run("Add polygon and remove TON in the runtime", func(t *testing.T) {
			// ARRANGE
			mustUpdateAppContextChainParams(t, app, []chains.Chain{
				chains.Ethereum, chains.BitcoinMainnet, chains.Polygon,
			})

			// ACT
			added, removed, err := syncObserverMap(ctx, client, tss, dbPath, baseLogger, ts, &observers)

			// ASSERT
			assert.NoError(t, err)
			assert.Equal(t, 1, added)
			assert.Equal(t, 1, removed)

			hasObserver(t, observers, chains.Ethereum.ChainId)
			hasObserver(t, observers, chains.Polygon.ChainId)
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

func TestBtcDatabaseFileName(t *testing.T) {
	tests := []struct {
		name     string
		chain    chains.Chain
		expected string
	}{
		{
			name:     "should use legacy file name for bitcoin mainnet",
			chain:    chains.BitcoinMainnet,
			expected: "btc_chain_client",
		},
		{
			name:     "should use legacy file name for bitcoin testnet3",
			chain:    chains.BitcoinTestnet,
			expected: "btc_chain_client",
		},
		{
			name:     "should use new file name for bitcoin regtest",
			chain:    chains.BitcoinRegtest,
			expected: "btc_chain_client_btc_regtest",
		},
		{
			name:     "should use new file name for bitcoin signet",
			chain:    chains.BitcoinSignetTestnet,
			expected: "btc_chain_client_btc_signet_testnet",
		},
		{
			name:     "should use new file name for bitcoin testnet4",
			chain:    chains.BitcoinTestnet4,
			expected: "btc_chain_client_btc_testnet4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, btcDatabaseFileName(tt.chain))
		})
	}
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

		if chains.IsEVMChain(chainID, nil) {
			params[chainID] = ptr.Ptr(mocks.MockChainParams(chainID, 100))
			continue
		}

		if chains.IsSolanaChain(chainID, nil) {
			p := mocks.MockChainParams(chainID, 100)
			p.GatewayAddress = solanaGatewayAddress
			params[chainID] = &p
			continue
		}

		if chains.IsTONChain(chainID, nil) {
			p := mocks.MockChainParams(chainID, 100)
			p.GatewayAddress = tonGatewayAddress
			params[chainID] = &p
			continue
		}

		panic("unknown chain: " + chain.String())
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
		chains,
		additionalChains,
		chainParams,
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
	t.Helper()

	signer, ok := observer[chainId]
	assert.True(t, ok, "missing observer for chain %d", chainId)
	assert.NotEmpty(t, signer)
}

func missesObserver(t *testing.T, observer map[int64]interfaces.ChainObserver, chainId int64) {
	_, ok := observer[chainId]
	assert.False(t, ok, "unexpected observer for chain %d", chainId)
}

// observer&signers have background tasks that rely on mocked calls.
// Ignorance results in FLAKY tests which fail silently with exit code 1.
func mockZetacore(client *mocks.ZetacoreClient) {
	// ctx context.Context, chain chains.Chain, gasPrice uint64, priorityFee uint64, blockNum uint64
	client.
		On("PostVoteGasPrice", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", nil).
		Maybe()
}
