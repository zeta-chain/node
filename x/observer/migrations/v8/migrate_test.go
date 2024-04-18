package v8_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/observer/exported"
	v8 "github.com/zeta-chain/zetacore/x/observer/migrations/v8"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

type mockSubspace struct {
	ps types.Params
}

func newMockSubspace(ps types.Params) mockSubspace {
	return mockSubspace{ps: ps}
}

func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps exported.ParamSet) {
	*ps.(*types.Params) = ms.ps
}

func TestMigrate(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)

	legacyParams := types.Params{
		BallotMaturityBlocks: 42,
	}
	legacySubspace := newMockSubspace(legacyParams)

	require.NoError(t, v8.MigrateStore(ctx, k, legacySubspace))

	params, found := k.GetParams(ctx)
	require.True(t, found)
	require.Equal(t, legacyParams, params)
}
