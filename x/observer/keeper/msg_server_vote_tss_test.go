package keeper_test

import (
	"fmt"
	"math"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_VoteTSS(t *testing.T) {
	t.Run("fail if node account not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		// ACT
		_, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          sample.AccAddress(),
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})

		// ASSERT
		require.ErrorIs(t, err, sdkerrors.ErrorInvalidSigner)
	})

	t.Run("fail if keygen is not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state
		nodeAcc := sample.NodeAccount()
		k.SetNodeAccount(ctx, *nodeAcc)

		// ACT
		_, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})

		// ASSERT
		require.ErrorIs(t, err, types.ErrKeygenNotFound)
	})

	t.Run("fail if keygen already completed ", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state
		nodeAcc := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.Status = types.KeygenStatus_KeyGenSuccess
		keygen.BlockNumber = 42
		k.SetNodeAccount(ctx, *nodeAcc)
		k.SetKeygen(ctx, *keygen)
		tss := sample.Tss()
		// ACT
		_, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})

		// ASSERT
		// keygen is already completed, but the vote can still be added if the operator has not voted yet
		require.NoError(t, err)
		ballot, found := k.GetBallot(ctx, fmt.Sprintf("%d-%s-%s", 42, tss.TssPubkey, "tss-keygen"))
		require.True(t, found)
		require.EqualValues(t, types.BallotStatus_BallotFinalized_SuccessObservation, ballot.BallotStatus)
		require.True(t, ballot.HasVoted(nodeAcc.Operator))
	})

	t.Run("can create a new ballot, vote success and finalize", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		finalizingHeight := int64(55)
		ctx = ctx.WithBlockHeight(finalizingHeight)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state
		nodeAcc := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.Status = types.KeygenStatus_PendingKeygen
		keygen.BlockNumber = 42
		k.SetNodeAccount(ctx, *nodeAcc)
		k.SetKeygen(ctx, *keygen)

		// ACT
		// there is a single node account, so the ballot will be created and finalized in a single vote
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})

		// ASSERT
		require.NoError(t, err)
		// check response
		require.True(t, res.BallotCreated)
		require.True(t, res.VoteFinalized)
		require.True(t, res.KeygenSuccess)

		// check keygen updated
		newKeygen, found := k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_KeyGenSuccess, newKeygen.Status)
		require.EqualValues(t, finalizingHeight, newKeygen.BlockNumber)

		// check tss updated
		tss, found := k.GetTSS(ctx)
		require.True(t, found)
		require.Equal(t, tss.KeyGenZetaHeight, int64(42))
		require.Equal(t, tss.FinalizedZetaHeight, finalizingHeight)
	})

	t.Run("can create a new ballot, vote failure and finalize", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ctx = ctx.WithBlockHeight(42)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state
		nodeAcc := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.BlockNumber = 42
		keygen.Status = types.KeygenStatus_PendingKeygen
		k.SetNodeAccount(ctx, *nodeAcc)
		k.SetKeygen(ctx, *keygen)

		// ACT
		// there is a single node account, so the ballot will be created and finalized in a single vote
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_failed,
		})

		// ASSERT
		require.NoError(t, err)
		// check response
		require.True(t, res.BallotCreated)
		require.True(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen updated
		newKeygen, found := k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_KeyGenFailed, newKeygen.Status)
		require.EqualValues(t, math.MaxInt64, newKeygen.BlockNumber)
	})

	t.Run(
		"can create a new ballot, vote without finalizing, then add final vote to update keygen and set tss",
		func(t *testing.T) {
			// ARRANGE
			k, ctx, _, _ := keepertest.ObserverKeeper(t)
			finalizingHeight := int64(55)
			ctx = ctx.WithBlockHeight(finalizingHeight)
			srv := keeper.NewMsgServerImpl(*k)

			// setup state with 3 node accounts
			nodeAcc1 := sample.NodeAccount()
			nodeAcc2 := sample.NodeAccount()
			nodeAcc3 := sample.NodeAccount()
			keygen := sample.Keygen(t)
			keygen.BlockNumber = 42
			keygen.Status = types.KeygenStatus_PendingKeygen
			tss := sample.Tss()
			k.SetNodeAccount(ctx, *nodeAcc1)
			k.SetNodeAccount(ctx, *nodeAcc2)
			k.SetNodeAccount(ctx, *nodeAcc3)
			k.SetKeygen(ctx, *keygen)

			// ACT
			// 1st vote: created ballot, but not finalized
			res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
				Creator:          nodeAcc1.Operator,
				TssPubkey:        tss.TssPubkey,
				KeygenZetaHeight: 42,
				Status:           chains.ReceiveStatus_success,
			})
			require.NoError(t, err)

			// check response
			require.True(t, res.BallotCreated)
			require.False(t, res.VoteFinalized)
			require.False(t, res.KeygenSuccess)

			// check keygen not updated
			newKeygen, found := k.GetKeygen(ctx)
			require.True(t, found)
			require.EqualValues(t, types.KeygenStatus_PendingKeygen, newKeygen.Status)

			// 2nd vote: already created ballot, and not finalized
			res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
				Creator:          nodeAcc2.Operator,
				TssPubkey:        tss.TssPubkey,
				KeygenZetaHeight: 42,
				Status:           chains.ReceiveStatus_success,
			})
			require.NoError(t, err)

			// check response
			require.False(t, res.BallotCreated)
			require.False(t, res.VoteFinalized)
			require.False(t, res.KeygenSuccess)

			// check keygen not updated
			newKeygen, found = k.GetKeygen(ctx)
			require.True(t, found)
			require.EqualValues(t, types.KeygenStatus_PendingKeygen, newKeygen.Status)

			// 3rd vote: finalize the ballot
			res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
				Creator:          nodeAcc3.Operator,
				TssPubkey:        tss.TssPubkey,
				KeygenZetaHeight: 42,
				Status:           chains.ReceiveStatus_success,
			})
			require.NoError(t, err)

			// ASSERT
			// check response
			require.False(t, res.BallotCreated)
			require.True(t, res.VoteFinalized)
			require.True(t, res.KeygenSuccess)

			// check keygen updated
			newKeygen, found = k.GetKeygen(ctx)
			require.True(t, found)
			require.EqualValues(t, types.KeygenStatus_KeyGenSuccess, newKeygen.Status)
			require.EqualValues(t, ctx.BlockHeight(), newKeygen.BlockNumber)

			// check tss updated
			tss, found = k.GetTSS(ctx)
			require.True(t, found)
			require.Equal(t, tss.KeyGenZetaHeight, int64(42))
			require.Equal(t, tss.FinalizedZetaHeight, finalizingHeight)
		},
	)

	t.Run("fail if voting fails", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ctx = ctx.WithBlockHeight(42)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state with two node accounts to not finalize the ballot
		nodeAcc := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.Status = types.KeygenStatus_PendingKeygen
		k.SetNodeAccount(ctx, *nodeAcc)
		k.SetNodeAccount(ctx, *sample.NodeAccount())
		k.SetKeygen(ctx, *keygen)
		tss := sample.Tss()

		// add a first vote
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)
		require.False(t, res.VoteFinalized)

		// ACT
		// vote again: voting should fail
		_, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})

		// ASSERT
		require.ErrorIs(t, err, types.ErrUnableToAddVote)
	})

	t.Run("can create a new ballot, without finalizing the older and then finalize older ballot", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ctx = ctx.WithBlockHeight(42)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state with 3 node accounts
		nodeAcc1 := sample.NodeAccount()
		nodeAcc2 := sample.NodeAccount()
		nodeAcc3 := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.Status = types.KeygenStatus_PendingKeygen
		tss := sample.Tss()
		k.SetNodeAccount(ctx, *nodeAcc1)
		k.SetNodeAccount(ctx, *nodeAcc2)
		k.SetNodeAccount(ctx, *nodeAcc3)
		k.SetKeygen(ctx, *keygen)

		// ACT
		// 1st vote: created ballot, but not finalized
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc1.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.True(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found := k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_PendingKeygen, newKeygen.Status)

		// 2nd vote: already created ballot, and not finalized
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc2.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.False(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_PendingKeygen, newKeygen.Status)

		keygen.Status = types.KeygenStatus_PendingKeygen
		keygen.BlockNumber = 52
		k.SetKeygen(ctx, *keygen)

		// Start voting on a new ballot
		tss2 := sample.Tss()
		// 1st Vote on new ballot (acc1)
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc1.Operator,
			TssPubkey:        tss2.TssPubkey,
			KeygenZetaHeight: 52,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.True(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// 2nd vote on new ballot: already created ballot, and not finalized (acc3)
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc3.Operator,
			TssPubkey:        tss2.TssPubkey,
			KeygenZetaHeight: 52,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.False(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen status
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_PendingKeygen, newKeygen.Status)

		// At this point, we have two ballots
		// 1. Ballot for keygen 42 Voted : (acc1, acc2)
		// 2. Ballot for keygen 52 Voted : (acc1, acc3)

		// 3rd vote on ballot 1: finalize the older ballot

		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc3.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// ASSERT
		// Check response
		require.False(t, res.BallotCreated)
		require.True(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)
		// Older ballot should be finalized which still keep keygen in pending state.
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_PendingKeygen, newKeygen.Status)

		_, found = k.GetTSS(ctx)
		require.False(t, found)

		oldBallot, found := k.GetBallot(ctx, fmt.Sprintf("%d-%s-%s", 42, tss.TssPubkey, "tss-keygen"))
		require.True(t, found)
		require.EqualValues(t, types.BallotStatus_BallotFinalized_SuccessObservation, oldBallot.BallotStatus)

		newBallot, found := k.GetBallot(ctx, fmt.Sprintf("%d-%s-%s", 52, tss2.TssPubkey, "tss-keygen"))
		require.True(t, found)
		require.EqualValues(t, types.BallotStatus_BallotInProgress, newBallot.BallotStatus)
	})

	t.Run("can create a new ballot, vote without finalizing,then finalize newer ballot", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ctx = ctx.WithBlockHeight(42)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state with 3 node accounts
		nodeAcc1 := sample.NodeAccount()
		nodeAcc2 := sample.NodeAccount()
		nodeAcc3 := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.Status = types.KeygenStatus_PendingKeygen
		keygen.BlockNumber = 42
		tss := sample.Tss()
		k.SetNodeAccount(ctx, *nodeAcc1)
		k.SetNodeAccount(ctx, *nodeAcc2)
		k.SetNodeAccount(ctx, *nodeAcc3)
		k.SetKeygen(ctx, *keygen)

		// ACT
		// 1st vote: created ballot, but not finalized
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc1.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.True(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found := k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_PendingKeygen, newKeygen.Status)

		// 2nd vote: already created ballot, and not finalized
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc2.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.False(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_PendingKeygen, newKeygen.Status)

		// Update keygen to 52 and start voting on new ballot
		keygen.Status = types.KeygenStatus_PendingKeygen
		keygen.BlockNumber = 52
		k.SetKeygen(ctx, *keygen)

		// Start voting on a new ballot
		tss2 := sample.Tss()
		// 1st Vote on new ballot (acc1)
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc1.Operator,
			TssPubkey:        tss2.TssPubkey,
			KeygenZetaHeight: 52,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.True(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// 2nd vote on new ballot: already created ballot, and not finalized (acc3)
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc3.Operator,
			TssPubkey:        tss2.TssPubkey,
			KeygenZetaHeight: 52,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.False(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen status
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_PendingKeygen, newKeygen.Status)

		// At this point, we have two ballots
		// 1. Ballot for keygen 42 Voted : (acc1, acc2)
		// 2. Ballot for keygen 52 Voted : (acc1, acc3)

		// 3rd vote on ballot 2: finalize the newer ballot

		finalizingHeight := int64(55)
		ctx = ctx.WithBlockHeight(finalizingHeight)
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc2.Operator,
			TssPubkey:        tss2.TssPubkey,
			KeygenZetaHeight: 52,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// ASSERT
		require.False(t, res.BallotCreated)
		require.True(t, res.VoteFinalized)
		require.True(t, res.KeygenSuccess)
		// Newer ballot should be finalized which make keygen success
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, finalizingHeight, newKeygen.BlockNumber)
		require.EqualValues(t, types.KeygenStatus_KeyGenSuccess, newKeygen.Status)

		tssQueried, found := k.GetTSS(ctx)
		require.True(t, found)
		require.Equal(t, tssQueried.KeyGenZetaHeight, int64(52))
		require.Equal(t, tssQueried.FinalizedZetaHeight, finalizingHeight)

		oldBallot, found := k.GetBallot(ctx, fmt.Sprintf("%d-%s-%s", 42, tss.TssPubkey, "tss-keygen"))
		require.True(t, found)
		require.EqualValues(t, types.BallotStatus_BallotInProgress, oldBallot.BallotStatus)

		newBallot, found := k.GetBallot(ctx, fmt.Sprintf("%d-%s-%s", 52, tss2.TssPubkey, "tss-keygen"))
		require.True(t, found)
		require.EqualValues(t, types.BallotStatus_BallotFinalized_SuccessObservation, newBallot.BallotStatus)
	})

	t.Run("add vote to a successful keygen", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ctx = ctx.WithBlockHeight(42)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state with 3 node accounts
		nodeAcc1 := sample.NodeAccount()
		nodeAcc2 := sample.NodeAccount()
		nodeAcc3 := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.Status = types.KeygenStatus_KeyGenSuccess
		tss := sample.Tss()
		k.SetNodeAccount(ctx, *nodeAcc1)
		k.SetNodeAccount(ctx, *nodeAcc2)
		k.SetNodeAccount(ctx, *nodeAcc3)
		k.SetKeygen(ctx, *keygen)

		// ACT
		// 1st vote: created ballot, but not finalized
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc1.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.True(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found := k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_KeyGenSuccess, newKeygen.Status)

		// 2nd vote: already created ballot, and not finalized
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc2.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.False(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_KeyGenSuccess, newKeygen.Status)

		// 3nd vote: already created ballot, and not finalized (acc3)
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc3.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.False(t, res.BallotCreated)
		require.True(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_KeyGenSuccess, newKeygen.Status)
	})

	t.Run("add vote to a failed keygen ", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ctx = ctx.WithBlockHeight(42)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state with 3 node accounts
		nodeAcc1 := sample.NodeAccount()
		nodeAcc2 := sample.NodeAccount()
		nodeAcc3 := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.Status = types.KeygenStatus_KeyGenFailed
		tss := sample.Tss()
		k.SetNodeAccount(ctx, *nodeAcc1)
		k.SetNodeAccount(ctx, *nodeAcc2)
		k.SetNodeAccount(ctx, *nodeAcc3)
		k.SetKeygen(ctx, *keygen)

		// ACT
		// 1st vote: created ballot, but not finalized
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc1.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_failed,
		})
		require.NoError(t, err)

		// check response
		require.True(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found := k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_KeyGenFailed, newKeygen.Status)

		// 2nd vote: already created ballot, and not finalized
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc2.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_failed,
		})
		require.NoError(t, err)

		// check response
		require.False(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_KeyGenFailed, newKeygen.Status)

		// 3nd vote: already created ballot, and not finalized (acc3)
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc3.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_failed,
		})
		require.NoError(t, err)

		// check response
		require.False(t, res.BallotCreated)
		require.True(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_KeyGenFailed, newKeygen.Status)
	})

	t.Run("unable to finalize tss if pubkey is different", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ctx = ctx.WithBlockHeight(42)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state with 3 node accounts
		nodeAcc1 := sample.NodeAccount()
		nodeAcc2 := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.BlockNumber = 42
		keygen.Status = types.KeygenStatus_PendingKeygen
		tss := sample.Tss()
		k.SetNodeAccount(ctx, *nodeAcc1)
		k.SetNodeAccount(ctx, *nodeAcc2)
		k.SetKeygen(ctx, *keygen)

		// ACT
		// Add first vote
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc1.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.True(t, res.BallotCreated)
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// Add second vote with different pubkey should not finalize the tss
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc2.Operator,
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// ASSERT
		require.True(t, res.BallotCreated) // New ballot created as pubkey is different
		require.False(t, res.VoteFinalized)
		require.False(t, res.KeygenSuccess)

		// Add the second vote with correct pubkey should finalize the tss
		finalizingHeight := int64(55)
		ctx = ctx.WithBlockHeight(finalizingHeight)
		res, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc2.Operator,
			TssPubkey:        tss.TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// ASSERT
		require.False(t, res.BallotCreated)
		require.True(t, res.VoteFinalized)
		require.True(t, res.KeygenSuccess)

		tssQueried, found := k.GetTSS(ctx)
		require.True(t, found)

		require.Equal(t, finalizingHeight, tssQueried.FinalizedZetaHeight)
		require.Equal(t, tss.TssPubkey, tssQueried.TssPubkey)
	})
}
