package keeper_test

import (
	"testing"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/keeper"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestMsgServer_UpdatePolicies(t *testing.T) {
	t.Run("can't update policies with invalid signer", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		policies := sample.Policies()

		_, err := msgServer.UpdatePolicies(sdk.WrapSDKContext(ctx), &types.MsgUpdatePolicies{
			Signer:   sample.AccAddress(),
			Policies: policies,
		})
		require.ErrorIs(t, err, govtypes.ErrInvalidSigner)
	})

	t.Run("can update policies", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		policies := sample.Policies()

		res, err := msgServer.UpdatePolicies(sdk.WrapSDKContext(ctx), &types.MsgUpdatePolicies{
			Signer:   keepertest.AuthorityGovAddress.String(),
			Policies: policies,
		})
		require.NotNil(t, res)
		require.NoError(t, err)

		// Check policy is set
		got, found := k.GetPolicies(ctx)
		require.True(t, found)
		require.Equal(t, policies, got)
	})
}
