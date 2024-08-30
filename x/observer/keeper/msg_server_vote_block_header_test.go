package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/proofs"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	mocks "github.com/zeta-chain/node/testutil/keeper/mocks/observer"
	"github.com/zeta-chain/node/testutil/sample"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func mockCheckNewBlockHeader(m *mocks.ObserverLightclientKeeper, err error) {
	m.On(
		"CheckNewBlockHeader",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(sample.Hash().Bytes(), err)
}

func mockAddBlockHeader(m *mocks.ObserverLightclientKeeper) {
	m.On(
		"AddBlockHeader",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	)
}

func TestMsgServer_VoteBlockHeader(t *testing.T) {
	one, err := sdk.NewDecFromStr("1.0")
	require.NoError(t, err)

	t.Run("fails if the chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		_, err := srv.VoteBlockHeader(ctx, &types.MsgVoteBlockHeader{
			Creator:   sample.AccAddress(),
			ChainId:   9999,
			BlockHash: sample.Hash().Bytes(),
			Height:    42,
			Header:    proofs.HeaderData{},
		})

		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("fails if the observer is not in the observer set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:         chains.GoerliLocalnet.ChainId,
					IsSupported:     true,
					BallotThreshold: one,
				},
			},
		})

		_, err := srv.VoteBlockHeader(ctx, &types.MsgVoteBlockHeader{
			Creator:   sample.AccAddress(),
			ChainId:   chains.GoerliLocalnet.ChainId,
			BlockHash: sample.Hash().Bytes(),
			Height:    42,
			Header:    proofs.HeaderData{},
		})

		require.ErrorIs(t, err, types.ErrNotObserver)
	})

	t.Run("fails if the new block header is invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)
		srv := keeper.NewMsgServerImpl(*k)
		observer := sample.AccAddress()

		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		lightclientMock := keepertest.GetObserverLightclientMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:         chains.GoerliLocalnet.ChainId,
					IsSupported:     true,
					BallotThreshold: one,
				},
			},
		})

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})

		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)
		mockCheckNewBlockHeader(lightclientMock, errors.New("foo"))

		_, err := srv.VoteBlockHeader(ctx, &types.MsgVoteBlockHeader{
			Creator:   observer,
			ChainId:   chains.GoerliLocalnet.ChainId,
			BlockHash: sample.Hash().Bytes(),
			Height:    42,
			Header:    proofs.HeaderData{},
		})

		require.ErrorIs(t, err, lightclienttypes.ErrInvalidBlockHeader)
	})

	t.Run("can create a new ballot, vote and finalize", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)
		srv := keeper.NewMsgServerImpl(*k)
		observer := sample.AccAddress()

		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		lightclientMock := keepertest.GetObserverLightclientMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:         chains.GoerliLocalnet.ChainId,
					IsSupported:     true,
					BallotThreshold: one,
				},
			},
		})

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})

		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)
		mockCheckNewBlockHeader(lightclientMock, nil)
		mockAddBlockHeader(lightclientMock)

		// there is a single node account, so the ballot will be created and finalized in a single vote
		res, err := srv.VoteBlockHeader(ctx, &types.MsgVoteBlockHeader{
			Creator:   observer,
			ChainId:   chains.GoerliLocalnet.ChainId,
			BlockHash: sample.Hash().Bytes(),
			Height:    42,
			Header:    proofs.HeaderData{},
		})

		require.NoError(t, err)
		require.True(t, res.VoteFinalized)
		require.True(t, res.BallotCreated)
	})

	t.Run("can create a new ballot, vote without finalizing, then add vote and finalizing", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)
		srv := keeper.NewMsgServerImpl(*k)
		observer1 := sample.AccAddress()
		observer2 := sample.AccAddress()
		observer3 := sample.AccAddress()
		blockHash := sample.Hash().Bytes()

		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		lightclientMock := keepertest.GetObserverLightclientMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:         chains.GoerliLocalnet.ChainId,
					IsSupported:     true,
					BallotThreshold: one,
				},
			},
		})

		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer1, observer2, observer3},
		})

		// first observer, created, not finalized
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)
		mockCheckNewBlockHeader(lightclientMock, nil)
		res, err := srv.VoteBlockHeader(ctx, &types.MsgVoteBlockHeader{
			Creator:   observer1,
			ChainId:   chains.GoerliLocalnet.ChainId,
			BlockHash: blockHash,
			Height:    42,
			Header:    proofs.HeaderData{},
		})

		require.NoError(t, err)
		require.False(t, res.VoteFinalized)
		require.True(t, res.BallotCreated)

		// second observer, found, not finalized
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)
		mockCheckNewBlockHeader(lightclientMock, nil)
		res, err = srv.VoteBlockHeader(ctx, &types.MsgVoteBlockHeader{
			Creator:   observer2,
			ChainId:   chains.GoerliLocalnet.ChainId,
			BlockHash: blockHash,
			Height:    42,
			Header:    proofs.HeaderData{},
		})

		require.NoError(t, err)
		require.False(t, res.VoteFinalized)
		require.False(t, res.BallotCreated)

		// third observer, found, finalized, add block header called
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)
		mockCheckNewBlockHeader(lightclientMock, nil)
		mockAddBlockHeader(lightclientMock)
		res, err = srv.VoteBlockHeader(ctx, &types.MsgVoteBlockHeader{
			Creator:   observer3,
			ChainId:   chains.GoerliLocalnet.ChainId,
			BlockHash: blockHash,
			Height:    42,
			Header:    proofs.HeaderData{},
		})

		require.NoError(t, err)
		require.True(t, res.VoteFinalized)
		require.False(t, res.BallotCreated)
	})

	t.Run("fail if voting fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)
		srv := keeper.NewMsgServerImpl(*k)
		observer := sample.AccAddress()
		blockHash := sample.Hash().Bytes()

		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		lightclientMock := keepertest.GetObserverLightclientMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:         chains.GoerliLocalnet.ChainId,
					IsSupported:     true,
					BallotThreshold: one,
				},
			},
		})

		// add multiple observers to not finalize the vote
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer, sample.AccAddress()},
		})

		// vote once
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)
		mockCheckNewBlockHeader(lightclientMock, nil)
		_, err := srv.VoteBlockHeader(ctx, &types.MsgVoteBlockHeader{
			Creator:   observer,
			ChainId:   chains.GoerliLocalnet.ChainId,
			BlockHash: blockHash,
			Height:    42,
			Header:    proofs.HeaderData{},
		})
		require.NoError(t, err)

		// vote a second time should make voting fail
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)
		mockCheckNewBlockHeader(lightclientMock, nil)
		_, err = srv.VoteBlockHeader(ctx, &types.MsgVoteBlockHeader{
			Creator:   observer,
			ChainId:   chains.GoerliLocalnet.ChainId,
			BlockHash: blockHash,
			Height:    42,
			Header:    proofs.HeaderData{},
		})
		require.ErrorIs(t, err, types.ErrUnableToAddVote)
	})
}
