package keeper_test

import (
	"math"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_VoteTSS(t *testing.T) {
	t.Run("fail if node account not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		_, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          sample.AccAddress(),
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.ErrorIs(t, err, sdkerrors.ErrorInvalidSigner)
	})

	t.Run("fail if keygen is not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state
		nodeAcc := sample.NodeAccount()
		k.SetNodeAccount(ctx, *nodeAcc)

		_, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.ErrorIs(t, err, types.ErrKeygenNotFound)
	})

	t.Run("fail if keygen already completed ", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state
		nodeAcc := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.Status = types.KeygenStatus_KeyGenSuccess
		k.SetNodeAccount(ctx, *nodeAcc)
		k.SetKeygen(ctx, *keygen)

		_, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.ErrorIs(t, err, types.ErrKeygenCompleted)
	})

	t.Run("can create a new ballot, vote success and finalize", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ctx = ctx.WithBlockHeight(42)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state
		nodeAcc := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.Status = types.KeygenStatus_PendingKeygen
		k.SetNodeAccount(ctx, *nodeAcc)
		k.SetKeygen(ctx, *keygen)

		// there is a single node account, so the ballot will be created and finalized in a single vote
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)

		// check response
		require.True(t, res.BallotCreated)
		require.True(t, res.VoteFinalized)
		require.True(t, res.KeygenSuccess)

		// check keygen updated
		newKeygen, found := k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_KeyGenSuccess, newKeygen.Status)
		require.EqualValues(t, ctx.BlockHeight(), newKeygen.BlockNumber)
	})

	t.Run("can create a new ballot, vote failure and finalize", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		ctx = ctx.WithBlockHeight(42)
		srv := keeper.NewMsgServerImpl(*k)

		// setup state
		nodeAcc := sample.NodeAccount()
		keygen := sample.Keygen(t)
		keygen.Status = types.KeygenStatus_PendingKeygen
		k.SetNodeAccount(ctx, *nodeAcc)
		k.SetKeygen(ctx, *keygen)

		// there is a single node account, so the ballot will be created and finalized in a single vote
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_failed,
		})
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

	t.Run("can create a new ballot, vote without finalizing, then add vote and finalizing", func(t *testing.T) {
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

		// check response
		require.False(t, res.BallotCreated)
		require.True(t, res.VoteFinalized)
		require.True(t, res.KeygenSuccess)

		// check keygen not updated
		newKeygen, found = k.GetKeygen(ctx)
		require.True(t, found)
		require.EqualValues(t, types.KeygenStatus_KeyGenSuccess, newKeygen.Status)
		require.EqualValues(t, ctx.BlockHeight(), newKeygen.BlockNumber)
	})

	t.Run("fail if voting fails", func(t *testing.T) {
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

		// add a first vote
		res, err := srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.NoError(t, err)
		require.False(t, res.VoteFinalized)

		// vote again: voting should fail
		_, err = srv.VoteTSS(ctx, &types.MsgVoteTSS{
			Creator:          nodeAcc.Operator,
			TssPubkey:        sample.Tss().TssPubkey,
			KeygenZetaHeight: 42,
			Status:           chains.ReceiveStatus_success,
		})
		require.ErrorIs(t, err, types.ErrUnableToAddVote)
	})
}
