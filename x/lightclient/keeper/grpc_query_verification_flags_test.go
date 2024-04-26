package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

func TestKeeper_VerificationFlags(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.VerificationFlags(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, _ := k.VerificationFlags(wctx, &types.QueryVerificationFlagsRequest{})
		require.Len(t, res.VerificationFlags, 0)
	})

	t.Run("should return if block header state is found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		vf := sample.VerificationFlags()
		for _, v := range vf {
			k.SetVerificationFlags(ctx, v)
		}

		res, err := k.VerificationFlags(wctx, &types.QueryVerificationFlagsRequest{})
		require.NoError(t, err)
		require.Equal(t, vf, res.VerificationFlags)
	})
}
