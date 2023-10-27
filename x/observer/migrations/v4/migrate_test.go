package v4_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	v4 "github.com/zeta-chain/zetacore/x/observer/migrations/v4"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateStore(t *testing.T) {

	k, ctx := keepertest.ObserverKeeper(t)
	store := prefix.NewStore(ctx.KVStore(k.StoreKey()), types.KeyPrefix(types.CrosschainFlagsKey))
	legacyFlags := types.LegacyCrosschainFlags{
		IsInboundEnabled:      false,
		IsOutboundEnabled:     false,
		GasPriceIncreaseFlags: &types.DefaultGasPriceIncreaseFlags,
	}
	val := k.Codec().MustMarshal(&legacyFlags)
	store.Set([]byte{0}, val)
	err := v4.MigrateStore(ctx, k.StoreKey(), k.Codec())
	assert.NoError(t, err)
	flags, found := k.GetCrosschainFlags(ctx)
	assert.True(t, found)
	assert.True(t, flags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled)
	assert.True(t, flags.BlockHeaderVerificationFlags.IsEthTypeChainEnabled)
}
