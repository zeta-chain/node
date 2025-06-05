package orchestrator

import (
	"testing"
	"time"

	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/scheduler"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testrpc"
)

func TestBootstrap(t *testing.T) {
	t.Run("Bitcoin", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		// Given orchestrator
		ts := newTestSuite(t)

		// Given BTC client
		btcServer, btcConfig := testrpc.NewBtcServer(t)

		ts.UpdateConfig(func(cfg *config.Config) {
			clearChainConfigs(cfg)
			cfg.BTCChainConfigs[chains.BitcoinMainnet.ChainId] = btcConfig
		})

		mockBitcoinCalls(ts, btcServer)
		mockZetacoreCalls(ts)

		// ACT
		// Start the orchestrator and wait for BTC observerSigner to bootstrap
		require.NoError(t, ts.Start(ts.ctx))

		// ASSERT
		// Check that btc observerSigner is bootstrapped.
		check := func() bool {
			return ts.HasObserverSigner(chains.BitcoinMainnet.ChainId)
		}

		assert.Eventually(t, check, 5*time.Second, 100*time.Millisecond)

		// Check that the scheduler has some tasks for this
		tasksHaveGroup(t, ts.scheduler.Tasks(), "btc:8332")

		assert.Contains(t, ts.Log.String(), `"chain":8332,"chain_network":"btc","message":"Added observer-signer"`)
	})

	t.Run("EVM", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		// Given orchestrator
		ts := newTestSuite(t)

		// Given ETH RPC
		ethServer := testrpc.NewEVMServer(t)
		mockEthCalls(ts, ethServer)

		maticServer := testrpc.NewEVMServer(t)
		mockEthCalls(ts, maticServer)

		ts.UpdateConfig(func(cfg *config.Config) {
			clearChainConfigs(cfg)

			cfg.EVMChainConfigs[chains.Ethereum.ChainId] = config.EVMConfig{
				Endpoint: ethServer.Endpoint,
			}
			cfg.EVMChainConfigs[chains.Polygon.ChainId] = config.EVMConfig{
				Endpoint: maticServer.Endpoint,
			}
		})

		// Mock zetacore calls
		mockZetacoreCalls(ts)

		// ACT #1
		// Start the orchestrator and wait for Ethereum observerSigner to bootstrap
		require.NoError(t, ts.Start(ts.ctx))

		// ASSERT #1
		// Ethereum observerSigner is added. Polygon is not
		check := func() bool {
			return ts.HasObserverSigner(chains.Ethereum.ChainId) &&
				!ts.HasObserverSigner(chains.Polygon.ChainId)
		}

		assert.Eventually(t, check, 5*time.Second, 100*time.Millisecond)

		tasksHaveGroup(t, ts.scheduler.Tasks(), "evm:1")
		assert.Contains(t, ts.Log.String(), `"chain":1,"chain_network":"eth","message":"Added observer-signer"`)

		// ACT #2
		// Enable polygon, remove ETH
		ts.MockChainParams(
			chains.Polygon, mocks.MockChainParams(chains.Polygon.ChainId, 100),
		)

		// ASSERT #2
		// Has only 1 chain
		check = func() bool {
			return !ts.HasObserverSigner(chains.Ethereum.ChainId) && ts.HasObserverSigner(chains.Polygon.ChainId)
		}

		assert.Eventually(t, check, 3*constant.ZetaBlockTime, 100*time.Millisecond)

		tasksHaveGroup(t, ts.scheduler.Tasks(), "evm:137")
		assert.Contains(t, ts.Log.String(), `"chain":137,"chain_network":"polygon","message":"Added observer-signer"`)

		tasksMissGroup(t, ts.scheduler.Tasks(), "evm:1")
		assert.Contains(t, ts.Log.String(), `"chain":1,"message":"Removed observer-signer"`)
	})

	t.Run("Solana", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		// Given orchestrator
		ts := newTestSuite(t)

		// Given Solana RPC
		solServer, solConfig := testrpc.NewSolanaServer(t)
		mockSolanaCalls(ts, solServer)

		ts.UpdateConfig(func(cfg *config.Config) {
			clearChainConfigs(cfg)
			cfg.SolanaConfig = solConfig
		})

		// Mock zetacore calls
		mockZetacoreCalls(ts)

		// ACT
		// Start the orchestrator and wait for SOL observerSigner to bootstrap
		require.NoError(t, ts.Start(ts.ctx))

		// ASSERT
		// Solana observerSigner is added
		check := func() bool {
			return ts.HasObserverSigner(chains.SolanaMainnet.ChainId)
		}

		assert.Eventually(t, check, 5*time.Second, 100*time.Millisecond)

		tasksHaveGroup(t, ts.scheduler.Tasks(), "sol:900")
		assert.Contains(t, ts.Log.String(), `"chain":900,"chain_network":"solana","message":"Added observer-signer"`)
	})

	t.Run("Sui", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		// Given orchestrator
		ts := newTestSuite(t)

		suiServer, suiConfig := testrpc.NewSuiServer(t)
		mockSuiCalls(ts, suiServer)

		// Disable other chains
		ts.UpdateConfig(func(cfg *config.Config) {
			clearChainConfigs(cfg)
			cfg.SuiConfig = suiConfig
		})

		// Mock zetacore calls
		mockZetacoreCalls(ts)

		// ACT
		// Start the orchestrator and wait for Sui observerSigner to bootstrap
		require.NoError(t, ts.Start(ts.ctx))

		// ASSERT
		// Solana observerSigner is added
		check := func() bool {
			return ts.HasObserverSigner(chains.SuiMainnet.ChainId)
		}

		assert.Eventually(t, check, 5*time.Second, 100*time.Millisecond)

		tasksHaveGroup(t, ts.scheduler.Tasks(), "sui:105")
		assert.Contains(t, ts.Log.String(), `"chain":105,"chain_network":"sui","message":"Added observer-signer"`)
	})

	t.Run("TON", func(t *testing.T) {
		t.Parallel()

		// ARRANGE
		// Given orchestrator
		ts := newTestSuite(t)

		// Given TON lite-server mock ...
		tonRPC, tonConfig := testrpc.NewTONServer(t)
		mockTONCalls(ts, tonRPC)

		// Given TON rpc URL
		ts.UpdateConfig(func(cfg *config.Config) {
			clearChainConfigs(cfg)
			cfg.TONConfig = tonConfig
		})

		// Mock zetacore calls
		mockZetacoreCalls(ts)

		// ACT
		// Start the orchestrator and wait for TON observerSigner to bootstrap
		require.NoError(t, ts.Start(ts.ctx))

		// ASSERT
		check := func() bool {
			return ts.HasObserverSigner(chains.TONMainnet.ChainId)
		}

		assert.Eventually(t, check, 3*constant.ZetaBlockTime, 100*time.Millisecond)

		tasksHaveGroup(t, ts.scheduler.Tasks(), "ton:2015140")
		assert.Contains(t, ts.Log.String(), `"chain":2015140,"chain_network":"ton","message":"Added observer-signer"`)
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

func tasksHaveGroup(t *testing.T, tasks map[uuid.UUID]*scheduler.Task, group string) {
	var found bool
	for _, task := range tasks {
		if !found && task.Group() == scheduler.Group(group) {
			found = true
		}
	}

	assert.True(t, found, "Group %s not found in tasks", group)
}

func tasksMissGroup(t *testing.T, tasks map[uuid.UUID]*scheduler.Task, group string) {
	var found bool
	for _, task := range tasks {
		// t.Logf("Task %s:%s", task.Group(), task.Name())
		if !found && task.Group() == scheduler.Group(group) {
			found = true
		}
	}

	assert.False(t, found, "Group %s found in tasks", group)
}

// isolates config from other chains
func clearChainConfigs(cfg *config.Config) {
	cfg.BTCChainConfigs = make(map[int64]config.BTCConfig)
	cfg.EVMChainConfigs = make(map[int64]config.EVMConfig)
	cfg.SolanaConfig.Endpoint = ""
	cfg.SuiConfig.Endpoint = ""
	cfg.TONConfig.Endpoint = ""
}

func mockZetacoreCalls(ts *testSuite) {
	blockChan := make(chan cometbfttypes.EventDataNewBlock)

	on(ts.zetacore, "NewBlockSubscriber", 1).Return(blockChan, nil)
	on(ts.zetacore, "GetInboundTrackersForChain", 2).Return(nil, nil).Maybe()
	on(ts.zetacore, "GetPendingNoncesByChain", 2).Return(observertypes.PendingNonces{}, nil).Maybe()
	on(ts.zetacore, "GetAllOutboundTrackerByChain", 3).Return(nil, nil).Maybe()
	on(ts.zetacore, "PostVoteGasPrice", 5).Return("", nil).Maybe()
}

func mockBitcoinCalls(_ *testSuite, client *testrpc.BtcServer) {
	client.SetBlockCount(100)
}

func mockEthCalls(_ *testSuite, client *testrpc.EVMServer) {
	client.SetBlockNumber(100)
	client.SetChainID(1)
}

// might be useful
func mockSolanaCalls(_ *testSuite, _ *testrpc.SolanaServer) {}
func mockSuiCalls(_ *testSuite, _ *testrpc.SuiServer)       {}
func mockTONCalls(_ *testSuite, _ *testrpc.TONServer)       {}
