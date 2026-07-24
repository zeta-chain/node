package snapshot

import (
	"encoding/hex"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func addr20(seed byte) []byte {
	b := make([]byte, addressLen)
	for i := range b {
		b[i] = seed
	}
	return b
}

func accBech32(t *testing.T, b []byte) string {
	t.Helper()
	s, err := bech32.ConvertAndEncode("zeta", b)
	require.NoError(t, err)
	return s
}

func valBech32(t *testing.T, b []byte) string {
	t.Helper()
	s, err := bech32.ConvertAndEncode("zetavaloper", b)
	require.NoError(t, err)
	return s
}

func canonHex(b []byte) string { return "0x" + hex.EncodeToString(b) }

// TestCompute exercises the attribution rules with a hand-built fixture:
//   - an EOA with only liquid
//   - a slashed validator (tokens != delegator_shares) with a self-delegation
//     and one external delegation, whose commission liquid + self-stake fold
//     into the operator account
//   - an external unbonding delegation
//   - a module account (bonded pool) that must NOT be attributed
//   - a pinned WZETA contract that must NOT be attributed
func TestCompute(t *testing.T) {
	// ARRANGE
	valBytes := addr20(0x02)
	delBytes := addr20(0x03)
	eoaBytes := addr20(0x01)
	unbBytes := addr20(0x04)
	wzetaBytes := addr20(0x05)
	poolBytes := authtypes.NewModuleAddress(stakingtypes.BondedPoolName).Bytes()

	// slashed validator: 90 tokens backing 100 shares -> exchange rate 0.9
	validator := stakingtypes.Validator{
		OperatorAddress: valBech32(t, valBytes),
		Tokens:          sdkmath.NewInt(90),
		DelegatorShares: sdkmath.LegacyNewDec(100),
		Status:          stakingtypes.Unbonded,
		Jailed:          true, // still counted regardless of status
	}

	coin := func(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }

	gen := &Genesis{
		Bank: banktypes.GenesisState{
			Supply: coin(2100),
			Balances: []banktypes.Balance{
				{Address: accBech32(t, eoaBytes), Coins: coin(1000)},
				{Address: accBech32(t, valBytes), Coins: coin(500)}, // withdrawn commission
				{Address: accBech32(t, delBytes), Coins: coin(200)},
				{Address: accBech32(t, unbBytes), Coins: coin(10)},
				{Address: accBech32(t, poolBytes), Coins: coin(90)}, // module: not attributed
				{Address: accBech32(t, wzetaBytes), Coins: coin(300)}, // pinned: not attributed
			},
		},
		Staking: stakingtypes.GenesisState{
			Validators: []stakingtypes.Validator{validator},
			Delegations: []stakingtypes.Delegation{
				// self-delegation: floor(40 * 90/100) = 36
				{DelegatorAddress: accBech32(t, valBytes), ValidatorAddress: valBech32(t, valBytes), Shares: sdkmath.LegacyNewDec(40)},
				// external: floor(33 * 90/100) = floor(29.7) = 29
				{DelegatorAddress: accBech32(t, delBytes), ValidatorAddress: valBech32(t, valBytes), Shares: sdkmath.LegacyNewDec(33)},
			},
			UnbondingDelegations: []stakingtypes.UnbondingDelegation{
				{
					DelegatorAddress: accBech32(t, unbBytes),
					ValidatorAddress: valBech32(t, valBytes),
					Entries: []stakingtypes.UnbondingDelegationEntry{
						{Balance: sdkmath.NewInt(50)},
						{Balance: sdkmath.NewInt(25)},
					},
				},
			},
		},
	}

	cfg := Config{Denom: BaseDenom, Pins: []string{canonHex(wzetaBytes)}}

	// ACT
	res, err := Compute(gen, cfg)

	// ASSERT
	require.NoError(t, err)

	byAddr := map[string]Account{}
	for _, a := range res.Accounts {
		byAddr[a.Canonical] = a
	}

	require.Len(t, res.Accounts, 4, "only the 4 claimable addresses are attributed")
	require.NotContains(t, byAddr, canonHex(poolBytes), "module account must not be attributed")
	require.NotContains(t, byAddr, canonHex(wzetaBytes), "pinned WZETA must not be attributed")

	eoa := byAddr[canonHex(eoaBytes)]
	require.Equal(t, classEOA, eoa.Class)
	require.Equal(t, "1000", eoa.Total().String())

	val := byAddr[canonHex(valBytes)]
	require.Equal(t, classValidator, val.Class)
	require.Equal(t, "500", val.Liquid.String())
	require.Equal(t, "36", val.Staked.String(), "self-delegation folds into operator account")
	require.Equal(t, "536", val.Total().String())

	del := byAddr[canonHex(delBytes)]
	require.Equal(t, "200", del.Liquid.String())
	require.Equal(t, "29", del.Staked.String(), "slashing-aware exchange rate is floored")

	unb := byAddr[canonHex(unbBytes)]
	require.Equal(t, "75", unb.Unbonding.String())
	require.Equal(t, "85", unb.Total().String())

	// bucket totals
	require.Equal(t, "1710", res.TotalLiquid.String()) // 1000+500+200+10
	require.Equal(t, "65", res.TotalStaked.String())    // 36+29
	require.Equal(t, "75", res.TotalUnbonding.String())

	// invariants
	require.Equal(t, "2100", res.Supply.String())
	require.Equal(t, "1850", res.AttributedTotal().String())
	require.Equal(t, "250", res.Remainder.String())
	require.False(t, res.Remainder.IsNegative())
	require.Equal(t, res.Supply.String(), res.AttributedTotal().Add(res.Remainder).String())
}

func TestComputeZeroAddressNotAttributed(t *testing.T) {
	// ARRANGE: an EOA plus a balance held by the null 20-byte address
	eoaBytes := addr20(0x01)
	zeroBytes := make([]byte, addressLen)
	coin := func(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }

	gen := &Genesis{
		Bank: banktypes.GenesisState{
			Supply: coin(1000),
			Balances: []banktypes.Balance{
				{Address: accBech32(t, eoaBytes), Coins: coin(600)},
				{Address: accBech32(t, zeroBytes), Coins: coin(400)}, // zero address: not claimable
			},
		},
	}

	// ACT
	res, err := Compute(gen, Config{Denom: BaseDenom})

	// ASSERT
	require.NoError(t, err)
	require.Len(t, res.Accounts, 1, "only the EOA is attributed")
	require.Equal(t, canonHex(eoaBytes), res.Accounts[0].Canonical)
	for _, a := range res.Accounts {
		require.NotEqual(t, zeroAddress, a.Canonical, "zero address must not be attributed")
	}
	require.Equal(t, "600", res.AttributedTotal().String())
	require.Equal(t, "400", res.Remainder.String(), "zero-address balance folds into remainder")
}

func TestComputeUnknownValidator(t *testing.T) {
	// ARRANGE: a delegation pointing at a validator not present in the set
	gen := &Genesis{
		Bank: banktypes.GenesisState{Supply: sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, 100))},
		Staking: stakingtypes.GenesisState{
			Delegations: []stakingtypes.Delegation{
				{
					DelegatorAddress: accBech32(t, addr20(0x03)),
					ValidatorAddress: valBech32(t, addr20(0x09)),
					Shares:           sdkmath.LegacyNewDec(10),
				},
			},
		},
	}

	// ACT
	_, err := Compute(gen, Config{Denom: BaseDenom})

	// ASSERT
	require.ErrorContains(t, err, "unknown validator")
}
