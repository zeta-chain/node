package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_Prove(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.Prove(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if invalid hash", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.Prove(wctx, &types.QueryProveRequest{
			ChainId:   987,
			BlockHash: sample.Hash().String(),
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if header not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.Prove(wctx, &types.QueryProveRequest{
			ChainId:   5,
			BlockHash: sample.Hash().String(),
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if proof not valid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		hash := sample.Hash()
		bh := proofs.BlockHeader{
			Height:     1,
			Hash:       hash.Bytes(),
			ParentHash: sample.Hash().Bytes(),
			ChainId:    1,
			Header:     proofs.HeaderData{},
		}
		k.SetBlockHeader(ctx, bh)

		res, err := k.Prove(wctx, &types.QueryProveRequest{
			ChainId:   5,
			BlockHash: hash.String(),
			Proof:     &proofs.Proof{},
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	// TODO: // https://github.com/zeta-chain/node/issues/1875 add more tests
}
