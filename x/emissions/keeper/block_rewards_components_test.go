package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestKeeper_GetReservesFactor(t *testing.T) {
	t.Run("successfully get reserves factor", func(t *testing.T) {
		//Arrange
		k, ctx, sdkK, _ := keepertest.EmissionsKeeper(t)
		amount := sdkmath.NewInt(100000000000000000)
		err := sdkK.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)))
		require.NoError(t, err)
		//Act
		reserveAmount := k.GetReservesFactor(ctx)
		//Assert
		require.Equal(t, amount.ToLegacyDec(), reserveAmount)
	})
}
