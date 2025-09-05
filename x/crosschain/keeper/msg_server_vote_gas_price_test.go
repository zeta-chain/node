package keeper_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestMsgServer_VoteGasPrice(t *testing.T) {
	t.Run("should error if unsupported chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		keepertest.MockFailedGetSupportedChainFromChainID(observerMock, sample.Chain(5))
		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.VoteGasPrice(ctx, &types.MsgVoteGasPrice{
			ChainId: 5,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if not non tombstoned observer", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		keepertest.MockGetSupportedChainFromChainID(observerMock, sample.Chain(5))
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(errors.New("not an observer"))

		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.VoteGasPrice(ctx, &types.MsgVoteGasPrice{
			ChainId: 5,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if gas price not found and set gas price in fungible keeper fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		keepertest.MockGetSupportedChainFromChainID(observerMock, chains.Chain{
			ChainId: 5,
		})
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(nil)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("SetGasPrice", mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), errors.New("err"))
		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.VoteGasPrice(ctx, &types.MsgVoteGasPrice{
			ChainId: 5,
		})
		require.Error(t, err)
		require.Nil(t, res)
		_, found := k.GetGasPrice(ctx, 5)
		require.True(t, found)
	})

	t.Run("should not error if gas price not found and set gas price in fungible keeper succeeds", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		keepertest.MockGetSupportedChainFromChainID(observerMock, chains.Chain{
			ChainId: 5,
		})
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(nil)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("SetGasPrice", mock.Anything, mock.Anything, mock.Anything).Return(uint64(1), nil)
		msgServer := keeper.NewMsgServerImpl(*k)
		creator := sample.AccAddress()
		res, err := msgServer.VoteGasPrice(ctx, &types.MsgVoteGasPrice{
			Creator:     creator,
			ChainId:     5,
			Price:       1,
			BlockNumber: 1,
		})
		require.NoError(t, err)
		require.Empty(t, res)
		gp, found := k.GetGasPrice(ctx, 5)
		require.True(t, found)
		require.Equal(t, types.GasPrice{
			Creator:      creator,
			Index:        "5",
			ChainId:      5,
			Signers:      []string{creator},
			BlockNums:    []uint64{1},
			Prices:       []uint64{1},
			PriorityFees: []uint64{0},
			MedianIndex:  0,
		}, gp)
	})

	t.Run("should not error if gas price found and msg.creator in signers", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		keepertest.MockGetSupportedChainFromChainID(observerMock, chains.Chain{
			ChainId: 5,
		})
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(nil)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("SetGasPrice", mock.Anything, mock.Anything, mock.Anything).Return(uint64(1), nil)
		msgServer := keeper.NewMsgServerImpl(*k)

		creator := sample.AccAddress()
		k.SetGasPrice(ctx, types.GasPrice{
			Creator:      creator,
			ChainId:      5,
			Signers:      []string{creator},
			BlockNums:    []uint64{1},
			Prices:       []uint64{1},
			PriorityFees: []uint64{0},
		})

		res, err := msgServer.VoteGasPrice(ctx, &types.MsgVoteGasPrice{
			Creator:     creator,
			ChainId:     5,
			BlockNumber: 2,
			Price:       2,
		})
		require.NoError(t, err)
		require.Empty(t, res)
		gp, found := k.GetGasPrice(ctx, 5)
		require.True(t, found)
		require.Equal(t, types.GasPrice{
			Creator:      creator,
			Index:        "5",
			ChainId:      5,
			Signers:      []string{creator},
			BlockNums:    []uint64{2},
			Prices:       []uint64{2},
			PriorityFees: []uint64{0},
			MedianIndex:  0,
		}, gp)
	})

	t.Run("should not error if gas price found and msg.creator not in signers", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		keepertest.MockGetSupportedChainFromChainID(observerMock, chains.Chain{
			ChainId: 5,
		})
		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(nil)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("SetGasPrice", mock.Anything, mock.Anything, mock.Anything).Return(uint64(1), nil)
		msgServer := keeper.NewMsgServerImpl(*k)

		creator := sample.AccAddress()
		k.SetGasPrice(ctx, types.GasPrice{
			Creator:      creator,
			ChainId:      5,
			BlockNums:    []uint64{1},
			Prices:       []uint64{1},
			PriorityFees: []uint64{0},
		})

		res, err := msgServer.VoteGasPrice(ctx, &types.MsgVoteGasPrice{
			Creator:     creator,
			ChainId:     5,
			BlockNumber: 2,
			Price:       2,
		})
		require.NoError(t, err)
		require.Empty(t, res)
		gp, found := k.GetGasPrice(ctx, 5)
		require.True(t, found)
		require.Equal(t, types.GasPrice{
			Creator:      creator,
			Index:        "5",
			ChainId:      5,
			Signers:      []string{creator},
			BlockNums:    []uint64{1, 2},
			Prices:       []uint64{1, 2},
			PriorityFees: []uint64{0, 0},
			MedianIndex:  1,
		}, gp)
	})

	t.Run("works with a priority fee", func(t *testing.T) {
		// ARRANGE
		// Given a keeper with grpc server and some mocks
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		// Given a chain
		chain := chains.Chain{ChainId: 5}

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, chain.ChainId).
			Return(chain, true)

		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(nil)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("SetGasPrice", mock.Anything, mock.Anything, mock.Anything).Return(uint64(1), nil)

		msgServer := keeper.NewMsgServerImpl(*k)
		creator := sample.AccAddress()

		// Given an existing gas price
		_, err := msgServer.VoteGasPrice(ctx, &types.MsgVoteGasPrice{
			Creator:     creator,
			ChainId:     chain.ChainId,
			BlockNumber: 2,
			Price:       2,
		})
		require.NoError(t, err)

		// ACT
		// When a new gas price is voted with a priority fee
		_, err = msgServer.VoteGasPrice(ctx, &types.MsgVoteGasPrice{
			Creator:     creator,
			ChainId:     5,
			BlockNumber: 3,
			Price:       3,
			PriorityFee: 2,
		})

		// ASSERT
		require.NoError(t, err)

		// Then gas prices should be updated as well as priority fee
		gp, found := k.GetGasPrice(ctx, 5)

		assert.True(t, found)
		assert.Equal(t, []string{creator}, gp.Signers)
		assert.Equal(t, []uint64{3}, gp.BlockNums)
		assert.Equal(t, []uint64{3}, gp.Prices)
		assert.Equal(t, []uint64{2}, gp.PriorityFees)
	})

	t.Run("tolerates lack of priorityFee of the same signer", func(t *testing.T) {
		// ARRANGE
		// Given a keeper with grpc server and some mocks
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		// Given a chain
		chain := chains.Chain{ChainId: 5}

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.
			On("GetSupportedChainFromChainID", mock.Anything, chain.ChainId).
			Return(chain, true)

		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(nil)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("SetGasPrice", mock.Anything, mock.Anything, mock.Anything).Return(uint64(1), nil)

		msgServer := keeper.NewMsgServerImpl(*k)
		creator := sample.AccAddress()

		// Given an existing gas price
		// Note that it MISSES priorityFee
		k.SetGasPrice(ctx, types.GasPrice{
			Creator:   creator,
			ChainId:   5,
			Signers:   []string{creator},
			BlockNums: []uint64{100},
			Prices:    []uint64{3},
		})

		// ACT
		// When a new gas price is voted with a priority fee
		_, err := msgServer.VoteGasPrice(ctx, &types.MsgVoteGasPrice{
			Creator:     creator,
			ChainId:     5,
			BlockNumber: 101,
			Price:       4,
			PriorityFee: 2,
		})

		// ASSERT
		require.NoError(t, err)

		// Then gas prices should be updated as well as priority fee
		gp, found := k.GetGasPrice(ctx, 5)

		assert.True(t, found)
		assert.Equal(t, []string{creator}, gp.Signers)
		assert.Equal(t, []uint64{101}, gp.BlockNums)
		assert.Equal(t, []uint64{4}, gp.Prices)
		assert.Equal(t, []uint64{2}, gp.PriorityFees)
	})

	t.Run("tolerates lack of priorityFee of another signer", func(t *testing.T) {
		// ARRANGE
		// Given a keeper with grpc server and some mocks
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		// Given a chain
		chain := chains.Chain{ChainId: 5}

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.
			On("GetSupportedChainFromChainID", mock.Anything, chain.ChainId).
			Return(chain, true)

		observerMock.On("CheckObserverCanVote", mock.Anything, mock.Anything).Return(nil)

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("SetGasPrice", mock.Anything, mock.Anything, mock.Anything).Return(uint64(1), nil)

		msgServer := keeper.NewMsgServerImpl(*k)
		creator := sample.AccAddress()
		creator2 := sample.AccAddress()

		// Given an existing gas price
		// Note that it MISSES priorityFee
		k.SetGasPrice(ctx, types.GasPrice{
			Creator:   creator,
			ChainId:   5,
			Signers:   []string{creator},
			BlockNums: []uint64{100},
			Prices:    []uint64{3},
		})

		// ACT
		// When a new gas price is voted with a priority fee
		_, err := msgServer.VoteGasPrice(ctx, &types.MsgVoteGasPrice{
			Creator:     creator2,
			ChainId:     5,
			BlockNumber: 100,
			Price:       4,
			PriorityFee: 2,
		})

		// ASSERT
		require.NoError(t, err)

		// Then gas prices should be updated as well as priority fee
		gp, found := k.GetGasPrice(ctx, 5)

		assert.True(t, found)
		assert.Equal(t, []string{creator, creator2}, gp.Signers)
		assert.Equal(t, []uint64{100, 100}, gp.BlockNums)
		assert.Equal(t, []uint64{3, 4}, gp.Prices)
		assert.Equal(t, []uint64{0, 2}, gp.PriorityFees)
	})
}
