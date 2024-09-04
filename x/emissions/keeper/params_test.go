package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
)

func TestKeeper_GetParams(t *testing.T) {
	tests := []struct {
		name         string
		params       emissionstypes.Params
		constainsErr string
	}{
		{
			name: "Successfully set params",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "",
		},
		{
			name: "negative observer slashed amount",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(-100000000000000000),
			},
			constainsErr: "slash amount cannot be less than 0",
		},
		{
			name: "MaxBondFactor too high",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "max bond factor cannot be higher that 1.25",
		},
		{
			name: "MinBondFactor too low",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "min bond factor cannot be lower that 0.75",
		},
		{
			name: "invalid block time",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "invalid block time",
		},
		{
			name: "invalid block time less than 0",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "block time cannot be less than or equal to 0",
		},
		{
			name: "bond ratio too high",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "target bond ratio cannot be more than 100 percent",
		},
		{
			name: "bond ratio too low",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "target bond ratio cannot be less than 0 percent",
		},
		{
			name: "validator emission percentage too high",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "1.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "validator emission percentage cannot be more than 100 percent",
		},
		{
			name: "validator emission percentage too low",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "-1.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "validator emission percentage cannot be less than 0 percent",
		},
		{
			name: "observer percentage too low",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "-00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "observer emission percentage cannot be less than 0 percent",
		},
		{
			name: "observer percentage too high",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "150.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "observer emission percentage cannot be more than 100 percent",
		},
		{
			name: "tss signer percentage too high",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "102.22",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "tss emission percentage cannot be more than 100 percent",
		},
		{
			name: "tss signer percentage too low",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "-102.22",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			constainsErr: "tss emission percentage cannot be less than 0 percent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx, _, _ := keepertest.EmissionsKeeper(t)
			err := k.SetParams(ctx, tt.params)
			if tt.constainsErr != "" {
				require.ErrorContains(t, err, tt.constainsErr)
			} else {
				require.NoError(t, err)
				params, found := k.GetParams(ctx)
				require.True(t, found)
				require.Equal(t, tt.params, params)
			}
		})
	}
}

func TestKeeper_GetParamsIfParamsNotSet(t *testing.T) {
	k, ctx, _, _ := keepertest.EmissionKeeperWithMockOptions(t, keepertest.EmissionMockOptions{SkipSettingParams: true})
	params, found := k.GetParams(ctx)
	require.False(t, found)
	require.Empty(t, params)
}
