package snapshot

import (
	"encoding/json"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/app"
)

func TestParseAppState(t *testing.T) {
	// ARRANGE
	cdc := app.MakeEncodingConfig(7001).Codec

	bank := banktypes.GenesisState{
		Params: banktypes.DefaultParams(),
		Supply: sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, 1500)),
		Balances: []banktypes.Balance{
			{Address: accBech32(t, addr20(0x01)), Coins: sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, 1000))},
			{Address: accBech32(t, addr20(0x02)), Coins: sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, 500))},
		},
	}
	staking := stakingtypes.GenesisState{
		Params: stakingtypes.DefaultParams(),
		Validators: []stakingtypes.Validator{{
			OperatorAddress: valBech32(t, addr20(0x02)),
			Tokens:          sdkmath.NewInt(500),
			DelegatorShares: sdkmath.LegacyNewDec(500),
			Status:          stakingtypes.Bonded,
		}},
	}

	distr := distrtypes.GenesisState{
		Params:  distrtypes.DefaultParams(),
		FeePool: distrtypes.FeePool{CommunityPool: sdk.NewDecCoins(sdk.NewInt64DecCoin(BaseDenom, 7))},
		OutstandingRewards: []distrtypes.ValidatorOutstandingRewardsRecord{
			{ValidatorAddress: valBech32(t, addr20(0x02))},
		},
	}

	appState := marshalAppState(t, cdc, &bank, &staking, &distr)

	// ACT
	gen, err := ParseAppState(cdc, appState)

	// ASSERT
	require.NoError(t, err)
	require.Len(t, gen.Bank.Balances, 2)
	require.Len(t, gen.Staking.Validators, 1)
	require.Len(t, gen.Distribution.OutstandingRewards, 1)
	require.Equal(t, "7", gen.Distribution.FeePool.CommunityPool.AmountOf(BaseDenom).TruncateInt().String())
	require.Equal(t, "1500", gen.Bank.Supply.AmountOf(BaseDenom).String())
	require.Equal(t, "1500", SumBankDenom(gen, BaseDenom).String())
}

func TestParseAppStateMissingModule(t *testing.T) {
	// ARRANGE / ACT / ASSERT
	cdc := app.MakeEncodingConfig(7001).Codec
	_, err := ParseAppState(cdc, map[string]json.RawMessage{})
	require.ErrorContains(t, err, "bank")
}

func TestParseAppStateMissingDistribution(t *testing.T) {
	// ARRANGE
	cdc := app.MakeEncodingConfig(7001).Codec
	appState := marshalAppState(t, cdc, &banktypes.GenesisState{}, &stakingtypes.GenesisState{}, nil)

	// ACT
	_, err := ParseAppState(cdc, appState)

	// ASSERT
	require.ErrorContains(t, err, `missing the "distribution" module`)
}

// marshalAppState encodes the given module genesis states into an app_state
// map; a nil distr leaves the distribution module out.
func marshalAppState(
	t *testing.T,
	cdc codec.JSONCodec,
	bank *banktypes.GenesisState,
	staking *stakingtypes.GenesisState,
	distr *distrtypes.GenesisState,
) map[string]json.RawMessage {
	t.Helper()
	bankRaw, err := cdc.MarshalJSON(bank)
	require.NoError(t, err)
	stakingRaw, err := cdc.MarshalJSON(staking)
	require.NoError(t, err)
	appState := map[string]json.RawMessage{
		banktypes.ModuleName:    bankRaw,
		stakingtypes.ModuleName: stakingRaw,
	}
	if distr != nil {
		distrRaw, err := cdc.MarshalJSON(distr)
		require.NoError(t, err)
		appState[distrtypes.ModuleName] = distrRaw
	}
	return appState
}
