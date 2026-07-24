package snapshot

import (
	"encoding/json"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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

	bankRaw, err := cdc.MarshalJSON(&bank)
	require.NoError(t, err)
	stakingRaw, err := cdc.MarshalJSON(&staking)
	require.NoError(t, err)
	appState := map[string]json.RawMessage{
		banktypes.ModuleName:    bankRaw,
		stakingtypes.ModuleName: stakingRaw,
	}

	// ACT
	gen, err := ParseAppState(cdc, appState)

	// ASSERT
	require.NoError(t, err)
	require.Len(t, gen.Bank.Balances, 2)
	require.Len(t, gen.Staking.Validators, 1)
	require.Equal(t, "1500", gen.Bank.Supply.AmountOf(BaseDenom).String())
	require.Equal(t, "1500", SumBankDenom(gen, BaseDenom).String())
}

func TestParseAppStateMissingModule(t *testing.T) {
	// ARRANGE / ACT / ASSERT
	cdc := app.MakeEncodingConfig(7001).Codec
	_, err := ParseAppState(cdc, map[string]json.RawMessage{})
	require.ErrorContains(t, err, "bank")
}
