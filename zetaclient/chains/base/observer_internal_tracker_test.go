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
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

var (
	// validCCTXErr is a valid error from zetacore client query 'GetCctxByHash'
	validCCTXErr = grpcstatus.Error(grpccodes.InvalidArgument, "a valid GRPC error (CCTX)")
)

func Test_GetInboundInternalTrackers(t *testing.T) {
	ctx := context.Background()
	chain := chains.Ethereum

	t.Run("should return empty internal tracker list", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ACT
		trackers := ob.GetInboundInternalTrackers(ctx, time.Now())

		// ASSERT
		require.Empty(t, trackers)
	})

	t.Run("should return non-empty internal tracker list", func(t *testing.T) {
		// ARRANGE
		ob := newTestSuite(t, chain)

		// mock cctx and ballot vote queries
		ob.zetacore.On("GetCctxByHash", ctx, mock.Anything).Return(nil, validCCTXErr)
		ob.zetacore.On("HasVoted", ctx, mock.Anything, mock.Anything).Return(false, nil)

		// add a failed inbound vote
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
		ob.AddInternalInboundTracker(ctx, &msg)

		// ACT
		retryTime := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime)

		// ASSERT
		require.Len(t, trackers, 1)
		require.Equal(t, msg.InboundTracker(), trackers[0])
	})

	t.Run("should remove internal tracker if the ballot is finalized", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ARRANGE
		// create a sample failed inbound vote
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)

		// mock cctx and ballot vote queries
		voterAddress := ob.ZetaRepo().GetOperatorAddress()
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, validCCTXErr).Twice()
		ob.zetacore.On("HasVoted", ctx, msg.Digest(), voterAddress).Return(false, nil)

		// add a failed inbound vote
		ob.AddInternalInboundTracker(ctx, &msg)

		// ACT 1
		retryTime1 := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime1)

		// ASSERT 1
		require.Len(t, trackers, 1)
		require.EqualValues(t, 1, len(trackers))

		// mock ballot as finalized
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(sample.CrossChainTx(t, msg.Digest()), nil).Once()

		// ACT 2
		retryTime2 := retryTime1.Add(internalTrackerRetryInterval)
		trackers = ob.GetInboundInternalTrackers(ctx, retryTime2)

		// ASSERT 2
		// should have removed internal tracker
		require.Empty(t, trackers)
	})

	t.Run("should remove internal tracker if the voter has already voted", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ARRANGE
		// create a sample failed inbound vote
		msg := sample.InboundVote(coin.CoinType_ERC20, 1, 7000)

		// mock cctx and ballot vote queries
		voterAddress := ob.ZetaRepo().GetOperatorAddress()
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, validCCTXErr)
		ob.zetacore.On("HasVoted", ctx, msg.Digest(), voterAddress).Return(false, nil).Twice()

		// add a failed inbound vote
		ob.AddInternalInboundTracker(ctx, &msg)

		// ACT 1
		retryTime1 := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime1)

		// ASSERT 1
		require.Len(t, trackers, 1)
		require.EqualValues(t, 1, len(trackers))

		// mock ballot as voted
		ob.zetacore.On("HasVoted", ctx, msg.Digest(), voterAddress).Return(true, nil).Once()

		// ACT 2
		retryTime2 := retryTime1.Add(internalTrackerRetryInterval)
		trackers = ob.GetInboundInternalTrackers(ctx, retryTime2)

		// ASSERT 2
		// should have removed internal tracker
		require.Empty(t, trackers)
	})

	t.Run("should skip recently retried internal tracker", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ARRANGE
		// create a sample failed inbound vote
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)

		// mock cctx and ballot vote queries
		voterAddress := ob.ZetaRepo().GetOperatorAddress()
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, validCCTXErr)
		ob.zetacore.On("HasVoted", ctx, msg.Digest(), voterAddress).Return(false, nil)

		// add a failed inbound vote
		ob.AddInternalInboundTracker(ctx, &msg)

		// ACT 1
		retryTime1 := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime1)

		// ASSERT 1
		require.Len(t, trackers, 1)
		require.EqualValues(t, 1, len(trackers))

		// retry with shorter interval
		retryTime2 := retryTime1.Add(internalTrackerRetryInterval - 1*time.Second)
		trackers = ob.GetInboundInternalTrackers(ctx, retryTime2)

		// ASSERT 2
		require.Empty(t, trackers)
	})

	t.Run("should skip internal tracker if unable to check ballot status", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ARRANGE
		// create a sample failed inbound vote
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)

		// mock cctx and ballot vote queries
		voterAddress := ob.ZetaRepo().GetOperatorAddress()
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, validCCTXErr).Twice()
		ob.zetacore.On("HasVoted", ctx, msg.Digest(), voterAddress).Return(false, nil).Twice()

		// add a failed inbound vote
		ob.AddInternalInboundTracker(ctx, &msg)

		// ACT 1
		retryTime1 := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime1)

		// ASSERT 1
		require.Len(t, trackers, 1)
		require.EqualValues(t, 1, len(trackers))

		// mock ballot status as unknown by returning an error
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, errors.New("invalid GRPC error")).Once()

		// retry with shorter interval
		retryTime2 := retryTime1.Add(internalTrackerRetryInterval)
		trackers = ob.GetInboundInternalTrackers(ctx, retryTime2)

		// ASSERT 2
		require.Empty(t, trackers)
	})

	t.Run("should only update timestamp for the first MaxInboundTrackersPerScan trackers", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ARRANGE
		// create MaxInboundTrackersPerScan + 1 sample failed inbound votes and mock up queries
		voterAddress := ob.ZetaRepo().GetOperatorAddress()
		msgs := make([]crosschaintypes.MsgVoteInbound, config.MaxInboundTrackersPerScan+1)
		for i := range config.MaxInboundTrackersPerScan + 1 {
			msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)
			msgs[i] = msg

			ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, validCCTXErr)
			ob.zetacore.On("HasVoted", ctx, msg.Digest(), voterAddress).Return(false, nil)
			ob.AddInternalInboundTracker(ctx, &msgs[i])
		}

		// ACT 1
		// should update timestamp for MaxInboundTrackersPerScan trackers
		retryTime1 := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime1)

		// ASSERT
		// should return all MaxInboundTrackersPerScan + 1 trackers
		require.Len(t, trackers, config.MaxInboundTrackersPerScan+1)

		// ACT 2
		// should return only one last remaining tracker
		trackers = ob.GetInboundInternalTrackers(ctx, retryTime1)

		// ASSERT 2
		require.Len(t, trackers, 1)

		// ACT 3
		// should return all MaxInboundTrackersPerScan + 1 trackers after another retry interval
		retryTime2 := retryTime1.Add(internalTrackerRetryInterval)
		trackers = ob.GetInboundInternalTrackers(ctx, retryTime2)

		// ASSERT 3
		require.Len(t, trackers, config.MaxInboundTrackersPerScan+1)
	})
}

func Test_AddInternalInboundTracker(t *testing.T) {
	ctx := context.Background()
	chain := chains.Ethereum

	t.Run("should not add internal inbound tracker if ballot is finalized", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ARRANGE
		// create a sample failed inbound vote
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)

		// mock cctx as finalized
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(sample.CrossChainTx(t, msg.Digest()), nil)

		// add a failed inbound vote
		ob.AddInternalInboundTracker(ctx, &msg)

		// ACT
		retryTime := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime)

		// ASSERT
		require.Empty(t, trackers)
	})

	t.Run("should not add internal inbound tracker if ballot is voted", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ARRANGE
		// create a sample failed inbound vote
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)

		// mock ballot as voted
		voterAddress := ob.ZetaRepo().GetOperatorAddress()
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, validCCTXErr)
		ob.zetacore.On("HasVoted", ctx, msg.Digest(), voterAddress).Return(true, nil)

		// add a failed inbound vote
		ob.AddInternalInboundTracker(ctx, &msg)

		// ACT
		retryTime := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime)

		// ASSERT
		require.Empty(t, trackers)
	})

	t.Run("should still add internal inbound tracker if failed to check cctx existence", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ARRANGE
		// create a sample failed inbound vote
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)

		// mock unknown ballot status
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, errors.New("invalid GRPC error")).Once()

		// add a failed inbound vote
		ob.AddInternalInboundTracker(ctx, &msg)

		// mock non-existent cctx and ballot as not voted
		voterAddress := ob.ZetaRepo().GetOperatorAddress()
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, validCCTXErr).Once()
		ob.zetacore.On("HasVoted", ctx, msg.Digest(), voterAddress).Return(false, nil).Once()

		// ACT
		retryTime := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime)

		// ASSERT
		require.Len(t, trackers, 1)
		require.Equal(t, msg.InboundTracker(), trackers[0])
	})

	t.Run("should still add internal inbound tracker if failed to check ballot vote status", func(t *testing.T) {
		ob := newTestSuite(t, chain)

		// ARRANGE
		// create a sample failed inbound vote
		msg := sample.InboundVote(coin.CoinType_Gas, 1, 7000)

		// mock unknown ballot vote status
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, validCCTXErr).Once()
		ob.zetacore.On("HasVoted", ctx, msg.Digest(), mock.Anything).Return(false, errors.New("some error")).Once()

		// add a failed inbound vote
		ob.AddInternalInboundTracker(ctx, &msg)

		// mock non-existent cctx and ballot as not voted
		voterAddress := ob.ZetaRepo().GetOperatorAddress()
		ob.zetacore.On("GetCctxByHash", ctx, msg.Digest()).Return(nil, validCCTXErr).Once()
		ob.zetacore.On("HasVoted", ctx, msg.Digest(), voterAddress).Return(false, nil).Once()

		// ACT
		retryTime := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime)

		// ASSERT
		require.Len(t, trackers, 1)
		require.Equal(t, msg.InboundTracker(), trackers[0])
	})
}

func Test_WatchMonitoringError(t *testing.T) {
	ctx := context.Background()
	chain := chains.Ethereum

	t.Run("should handle monitoring error", func(t *testing.T) {
		// ARRANGE
		ob := newTestSuite(t, chain)
		monitorErrCh := make(chan zetaerrors.ErrTxMonitor, 1)

		// mock cctx and ballot vote queries
		ob.zetacore.On("GetCctxByHash", ctx, mock.Anything).Return(nil, validCCTXErr)
		ob.zetacore.On("HasVoted", ctx, mock.Anything, mock.Anything).Return(false, nil)

		// ACT
		// start the monitoring error handler
		go func() {
			ob.WatchMonitoringError(ctx, monitorErrCh, "zetaHash")
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
		retryTime := time.Now().Add(internalTrackerRetryInterval)
		trackers := ob.GetInboundInternalTrackers(ctx, retryTime)
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

		// ACT
		ob.WatchMonitoringError(ctxTimeout, monitorErrCh, "zetaHash")

		// ASSERT
		trackers := ob.GetInboundInternalTrackers(ctx, time.Now())
		require.Empty(t, trackers)
	})
}
