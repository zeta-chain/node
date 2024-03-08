package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
)

func TestKeeper_GetDistributions(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.EmissionsKeeper(t)

	val, obs, tss := k.GetDistributions(ctx)

	require.EqualValues(t, "4810474537037037037", val.String()) // 0.5 * block reward
	require.EqualValues(t, "2405237268518518518", obs.String()) // 0.25 * block reward
	require.EqualValues(t, "2405237268518518518", tss.String()) // 0.25 * block reward
}
