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
	"github.com/zeta-chain/node/zetaclient/config"
)

func Test_GetInboundInternalTrackers(t *testing.T) {
	ctx := context.Background()
	chain := chains.Ethereum

	t.Run("should return empty internal tracker list", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ACT
		trackers, totalCount := ob.GetInboundInternalTrackers(ctx)

		// ASSERT
		require.Empty(t, trackers)
		require.Zero(t, totalCount)
	})

	t.Run("should return non-empty internal tracker list", func(t *testing.T) {
		// ARRANGE
		ob := newTestSuite(t, chain)

		// mock cctx query
		ob.zetacore.On("GetCctxByHash", ctx, mock.Anything).Return(nil, errors.New("not found"))

		// add a failed inbound vote
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
		ob.AddInternalInboundTracker(&msg)

		// ACT
		trackers, totalCount := ob.GetInboundInternalTrackers(ctx)

		// ASSERT
		require.Len(t, trackers, 1)
		require.Equal(t, msg.InboundTracker(), trackers[0])
		require.EqualValues(t, 1, totalCount)
	})

	t.Run("should remove internal tracker if the ballot is finalized", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// add a failed inbound vote
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
		ob.AddInternalInboundTracker(&msg)

		// mock cctx query
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, errors.New("not found")).Once()

		// ACT 1
		trackers, totalCount := ob.GetInboundInternalTrackers(ctx)

		// ASSERT 1
		require.Len(t, trackers, 1)
		require.EqualValues(t, 1, totalCount)

		// mock ballot as finalized
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, nil).Once()

		// ACT 2
		trackers, totalCount = ob.GetInboundInternalTrackers(ctx)

		// ASSERT 2
		// should remove internal tracker
		require.Len(t, trackers, 0)
		require.EqualValues(t, 0, totalCount)
	})

	t.Run("should return at most 'MaxInternalTrackersPerScan' internal trackers", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// add more than 'MaxInternalTrackersPerScan' internal trackers
		msgs := addNInternalTrackers(ob, config.MaxInternalTrackersPerScan+1)

		// mock cctx queries
		for _, msg := range msgs {
			ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Maybe().Return(nil, errors.New("not found"))
		}

		// ACT
		trackers, totalCount := ob.GetInboundInternalTrackers(ctx)
		require.Len(t, trackers, config.MaxInternalTrackersPerScan)
		require.EqualValues(t, config.MaxInternalTrackersPerScan+1, totalCount)
	})
}

func Test_handleMonitoringError(t *testing.T) {
	ctx := context.Background()
	chain := chains.Ethereum

	t.Run("should handle monitoring error", func(t *testing.T) {
		// ARRANGE
		ob := newTestSuite(t, chain)
		monitorErrCh := make(chan zetaerrors.ErrTxMonitor, 1)

		// mock cctx query
		ob.zetacore.On("GetCctxByHash", ctx, mock.Anything).Return(nil, errors.New("not found"))

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

		// wait for the internal tracker to be added
		time.Sleep(1 * time.Second)

		// ASSERT
		trackers, totalCount := ob.GetInboundInternalTrackers(ctx)
		require.Len(t, trackers, 1)
		require.Equal(t, vote.InboundTracker(), trackers[0])
		require.EqualValues(t, 1, totalCount)
	})

	t.Run("should time out if no error is received", func(t *testing.T) {
		// ARRANGE
		ob := newTestSuite(t, chain)

		// create a context with timeout and a monitor error channel
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
		_ = cancel
		monitorErrCh := make(chan zetaerrors.ErrTxMonitor, 1)

		// ACT
		ob.handleMonitoringError(ctxTimeout, monitorErrCh, "zetaHash")

		// ASSERT
		trackers, totalCount := ob.GetInboundInternalTrackers(ctx)
		require.Empty(t, trackers)
		require.EqualValues(t, 0, totalCount)
	})
}

// addNInternalTrackers adds n internal trackers to the observer
func addNInternalTrackers(ob *testSuite, n int) []crosschaintypes.MsgVoteInbound {
	msgs := make([]crosschaintypes.MsgVoteInbound, 0, n)
	for range n {
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
		ob.AddInternalInboundTracker(&msg)
		msgs = append(msgs, msg)
	}
	return msgs
}
