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
	"github.com/zeta-chain/node/pkg/scheduler"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/testutils/testrpc"
)

func TestBootstrap(t *testing.T) {
	t.Run("Bitcoin", func(t *testing.T) {
		// ARRANGE
		// Given orchestrator
		ts := newTestSuite(t)

		// Given BTC client
		btcServer, btcConfig := testrpc.NewBtcServer(t)

		ts.UpdateConfig(func(cfg *config.Config) {
			cfg.BTCChainConfigs[chains.BitcoinMainnet.ChainId] = btcConfig
		})

		mockBitcoinCalls(ts, btcServer)

		// ACT
		// Start the orchestrator and wait for BTC observerSigner to bootstrap
		require.NoError(t, ts.Start(ts.ctx))

		// ASSERT
		// Check that btc observerSigner is bootstrapped.
		check := func() bool {
			ts.V2.mu.RLock()
			defer ts.V2.mu.RUnlock()

			_, ok := ts.V2.chains[chains.BitcoinMainnet.ChainId]
			return ok
		}

		assert.Eventually(t, check, 5*time.Second, 100*time.Millisecond)

		// Check that the scheduler has some tasks for this
		tasksHaveGroup(t, ts.scheduler.Tasks(), "btc:8332")

		assert.Contains(t, ts.Log.String(), `"chain":8332,"chain_network":"btc","message":"Added observer-signer"`)
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

func mockBitcoinCalls(ts *testSuite, client *testrpc.BtcServer) {
	client.SetBlockCount(100)

	blockChan := make(chan cometbfttypes.EventDataNewBlock)
	ts.zetacore.On("NewBlockSubscriber", mock.Anything).Return(blockChan, nil)

	ts.zetacore.On("GetInboundTrackersForChain", mock.Anything, mock.Anything).Return(nil, nil)
	ts.zetacore.On("GetPendingNoncesByChain", mock.Anything, mock.Anything).Return(observertypes.PendingNonces{}, nil)
	ts.zetacore.On("GetAllOutboundTrackerByChain", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
}
