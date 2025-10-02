package base

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func Test_GetInboundTrackersWithBacklog(t *testing.T) {
	ctx := context.Background()
	chain := chains.Ethereum

	t.Run("should return empty trackers if no failed inbounds", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// mock inbound trackers
		ob.zetacore.On("GetInboundTrackersForChain", ctx, chain.ChainId).Return([]crosschaintypes.InboundTracker{}, nil)

		trackers, err := ob.GetInboundTrackersWithBacklog(ctx)
		require.NoError(t, err)
		require.Empty(t, trackers)
	})

	t.Run("should return inbound trackers from zetacore and backlog", func(t *testing.T) {
		// ARRANGE
		ob := newTestSuite(t, chain)

		// mock inbound tracker and cctx queries
		tracker1 := sample.InboundTracker(t, "tracker1")
		ob.zetacore.On("GetCctxByHash", ctx, mock.Anything).Return(nil, errors.New("not found"))
		ob.zetacore.On("GetInboundTrackersForChain", ctx, chain.ChainId).Return([]crosschaintypes.InboundTracker{tracker1}, nil)

		// add an additional failed inbound
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
		ob.AddFailedInbound(&msg)

		// ACT
		trackers, err := ob.GetInboundTrackersWithBacklog(ctx)

		// ASSERT
		require.NoError(t, err)
		require.Len(t, trackers, 2)
		require.Equal(t, tracker1, trackers[0])
		require.Equal(t, msg.InboundTracker(), trackers[1])
	})

	t.Run("should remove failed inbound from backlog if it is finalized", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// mock inbound trackers and cctx queries
		tracker := sample.InboundTracker(t, "tracker")
		ob.zetacore.On("GetCctxByHash", ctx, mock.Anything).Return(nil, errors.New("not found")).Once()
		ob.zetacore.On("GetInboundTrackersForChain", ctx, chain.ChainId).Return([]crosschaintypes.InboundTracker{tracker}, nil)

		// add a failed inbound
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
		ob.AddFailedInbound(&msg)

		// ACT 1
		trackers, err := ob.GetInboundTrackersWithBacklog(ctx)

		// ASSERT 1
		require.NoError(t, err)
		require.Len(t, trackers, 2)

		// mock finalized cctx query
		ob.zetacore.On("GetCctxByHash", ctx, mock.Anything).Return(nil, nil)

		// ACT 2
		trackers, err = ob.GetInboundTrackersWithBacklog(ctx)

		// ASSERT 2
		// should remove failed inbound from backlog
		require.NoError(t, err)
		require.Len(t, trackers, 1)
		require.Equal(t, tracker, trackers[0])
	})
}

func Test_handleMonitoringError(t *testing.T) {
	ctx := context.Background()
	chain := chains.Ethereum

	t.Run("should handle monitoring error", func(t *testing.T) {
		// ARRANGE
		ob := newTestSuite(t, chain)
		monitorErrCh := make(chan zetaerrors.ErrTxMonitor, 1)

		// mock inbound trackers query
		ob.zetacore.On("GetCctxByHash", ctx, mock.Anything).Return(nil, errors.New("not found"))
		ob.zetacore.On("GetInboundTrackersForChain", ctx, chain.ChainId).Return([]crosschaintypes.InboundTracker{}, nil)

		// ACT
		// start the monitoring error handler
		go func() {
			ob.handleMonitoringError(ctx, monitorErrCh, "zetaHash")
		}()

		// feed an error to the channel
		vote := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
		monitorErrCh <- zetaerrors.ErrTxMonitor{
			Err:        errors.New("monitoring error"),
			Msg:        vote,
			ZetaTxHash: vote.InboundHash,
		}

		// wait for the failed inbound to be added to the backlog
		time.Sleep(1 * time.Second)

		// ASSERT
		trackers, err := ob.GetInboundTrackersWithBacklog(ctx)
		require.NoError(t, err)
		require.Len(t, trackers, 1)
		require.Equal(t, vote.InboundTracker(), trackers[0])
	})

	t.Run("should time out if no error is received", func(t *testing.T) {
		// ARRANGE
		ob := newTestSuite(t, chain)

		// create a context with timeout and a monitor error channel
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
		_ = cancel
		monitorErrCh := make(chan zetaerrors.ErrTxMonitor, 1)

		// mock inbound trackers query
		ob.zetacore.On("GetInboundTrackersForChain", ctx, chain.ChainId).Return([]crosschaintypes.InboundTracker{}, nil)

		// ACT
		ob.handleMonitoringError(ctxTimeout, monitorErrCh, "zetaHash")

		// ASSERT
		trackers, err := ob.GetInboundTrackersWithBacklog(ctx)
		require.NoError(t, err)
		require.Empty(t, trackers)
	})
}
