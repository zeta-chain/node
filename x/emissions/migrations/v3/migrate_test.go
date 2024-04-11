package v3_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/exported"
	v3 "github.com/zeta-chain/zetacore/x/emissions/migrations/v3"
	"github.com/zeta-chain/zetacore/x/emissions/types"
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
	k, ctx, _, _ := keepertest.EmissionsKeeper(t)

	legacyParams := types.Params{
		MaxBondFactor:               "1",
		MinBondFactor:               "0.50",
		AvgBlockTime:                "5.00",
		TargetBondRatio:             "00.50",
		ValidatorEmissionPercentage: "00.50",
		ObserverEmissionPercentage:  "00.35",
		TssSignerEmissionPercentage: "00.15",
		DurationFactorConstant:      "0.001877876953694702",
		ObserverSlashAmount:         sdk.NewInt(10),
	}
	legacySubspace := newMockSubspace(legacyParams)

	require.NoError(t, v3.MigrateStore(ctx, k, legacySubspace))

	params := k.GetParams(ctx)
	require.Equal(t, legacyParams, params)
}
