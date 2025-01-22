package orchestrator

import (
	"testing"
	"time"

	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			// disable btc for this test
			cfg.BTCChainConfigs = nil
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
}

func tasksHaveGroup(t *testing.T, tasks map[uuid.UUID]*scheduler.Task, group string) {
	var found bool
	for _, task := range tasks {
		// t.Logf("Task %s:%s", task.Group(), task.Name())
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

func mockBitcoinCalls(_ *testSuite, client *testrpc.BtcServer) {
	client.SetBlockCount(100)
}

func mockEthCalls(_ *testSuite, client *testrpc.EVMServer) {
	client.SetBlockNumber(100)
	client.SetChainID(1)

}

func mockZetacoreCalls(ts *testSuite) {
	blockChan := make(chan cometbfttypes.EventDataNewBlock)
	ts.zetacore.On("NewBlockSubscriber", mock.Anything).Return(blockChan, nil).Maybe()

	ts.zetacore.On("GetInboundTrackersForChain", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	ts.zetacore.On("GetPendingNoncesByChain", mock.Anything, mock.Anything).
		Return(observertypes.PendingNonces{}, nil).
		Maybe()
	ts.zetacore.On("GetAllOutboundTrackerByChain", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()
}
