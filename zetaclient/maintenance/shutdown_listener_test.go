package maintenance

import (
	"context"
	"testing"
	"time"

	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func assertChannelNotClosed[T any](t *testing.T, ch <-chan T) {
	select {
	case <-ch:
		t.Errorf("failed: channel was closed")
	default:
	}
}

func newBlockEventHeightOnly(height int64) cometbfttypes.EventDataNewBlock {
	return cometbfttypes.EventDataNewBlock{
		Block: &cometbfttypes.Block{
			Header: cometbfttypes.Header{
				Height: height,
			},
		},
	}
}

func TestShutdownListener(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := zerolog.New(zerolog.NewTestWriter(t))

	t.Run("scheduled shutdown", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)

		listener := NewShutdownListener(client, logger)

		client.Mock.On("GetOperationalFlags", ctx).Return(observertypes.OperationalFlags{
			RestartHeight: 10,
		}, nil)
		client.Mock.On("GetBlockHeight", ctx).Return(int64(8), nil)
		blockChan := make(chan cometbfttypes.EventDataNewBlock)
		client.Mock.On("NewBlockSubscriber", ctx).Return(blockChan, nil)
		client.Mock.On("GetSyncStatus", ctx).Return(false, nil)

		complete := make(chan interface{})
		listener.Listen(ctx, func() {
			close(complete)
		})

		assertChannelNotClosed(t, complete)

		blockChan <- newBlockEventHeightOnly(9)
		assertChannelNotClosed(t, complete)

		blockChan <- newBlockEventHeightOnly(10)
		<-complete
	})

	t.Run("no shutdown", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)

		listener := NewShutdownListener(client, logger)

		client.Mock.On("GetOperationalFlags", ctx).Return(observertypes.OperationalFlags{}, nil)
		// GetBlockHeight is not mocked because we want the test to panic if it's called
		// NewBlockSubscriber is not mocked because we want the test to panic if it's called
		complete := make(chan interface{})
		client.Mock.On("GetSyncStatus", ctx).Return(false, nil)
		listener.Listen(ctx, func() {
			close(complete)
		})

		require.Eventually(t, func() bool {
			return len(client.Calls) == 2
		}, time.Second, time.Millisecond)

		assertChannelNotClosed(t, complete)
	})

	t.Run("shutdown if zetacore is syncing", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)

		listener := NewShutdownListener(client, logger)

		client.Mock.On("GetOperationalFlags", ctx).Return(observertypes.OperationalFlags{}, nil)
		client.Mock.On("GetSyncStatus", ctx).Return(true, nil)
		complete := make(chan interface{})
		listener.Listen(ctx, func() {
			close(complete)
		})

		require.Eventually(t, func() bool {
			return len(client.Calls) == 2
		}, time.Second, time.Millisecond)

		<-complete
	})

	t.Run("shutdown height missed", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)

		listener := NewShutdownListener(client, logger)

		client.Mock.On("GetOperationalFlags", ctx).Return(observertypes.OperationalFlags{
			RestartHeight: 10,
		}, nil)
		client.Mock.On("GetBlockHeight", ctx).Return(int64(11), nil)
		client.Mock.On("GetSyncStatus", ctx).Return(false, nil)

		complete := make(chan interface{})
		listener.Listen(ctx, func() {
			close(complete)
		})

		require.Eventually(t, func() bool {
			return len(client.Calls) == 3
		}, time.Second, time.Millisecond)
		assertChannelNotClosed(t, complete)
	})

	t.Run("minimum version ok", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)

		listener := NewShutdownListener(client, logger)
		listener.getVersion = func() string {
			return "v1.1.2"
		}

		client.Mock.On("GetOperationalFlags", ctx).Return(observertypes.OperationalFlags{
			MinimumVersion: "v1.1.1",
		}, nil)
		client.Mock.On("GetSyncStatus", ctx).Return(false, nil)

		// pre start checks passed
		err := listener.RunPreStartCheck(ctx)
		require.NoError(t, err)

		// listener also does not shutdown
		complete := make(chan interface{})
		listener.Listen(ctx, func() {
			close(complete)
		})

		require.Eventually(t, func() bool {
			return len(client.Calls) == 3
		}, time.Second, time.Millisecond)
		assertChannelNotClosed(t, complete)
	})

	t.Run("minimum version failed", func(t *testing.T) {
		client := mocks.NewZetacoreClient(t)

		listener := NewShutdownListener(client, logger)
		listener.getVersion = func() string {
			return "v1.1.1"
		}

		client.Mock.On("GetOperationalFlags", ctx).Return(observertypes.OperationalFlags{
			MinimumVersion: "v1.1.2",
		}, nil)
		client.Mock.On("GetSyncStatus", ctx).Return(false, nil)

		// pre start checks would return error
		err := listener.RunPreStartCheck(ctx)
		require.Error(t, err)

		// listener would also shutdown
		complete := make(chan interface{})
		listener.Listen(ctx, func() {
			close(complete)
		})

		require.Eventually(t, func() bool {
			return len(client.Calls) == 3
		}, time.Second, time.Millisecond)
		<-complete
	})

	// avoid Log in goroutine after TestShutdownListener has completed
	cancel()
	time.Sleep(time.Millisecond * 100)
}
