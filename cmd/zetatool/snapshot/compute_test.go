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
//   - a module account (not-bonded pool) that must NOT be attributed
//   - a pinned WZETA contract that must NOT be attributed
//   - a zero-address delegation that completes the validator's shares
func TestCompute(t *testing.T) {
	// ARRANGE
	valBytes := addr20(0x02)
	delBytes := addr20(0x03)
	eoaBytes := addr20(0x01)
	unbBytes := addr20(0x04)
	wzetaBytes := addr20(0x05)
	zeroBytes := make([]byte, addressLen)
	poolBytes := authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).Bytes()

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
			Supply: coin(2175),
			Balances: []banktypes.Balance{
				{Address: accBech32(t, eoaBytes), Coins: coin(1000)},
				{Address: accBech32(t, valBytes), Coins: coin(500)}, // withdrawn commission
				{Address: accBech32(t, delBytes), Coins: coin(200)},
				{Address: accBech32(t, unbBytes), Coins: coin(10)},
				{Address: accBech32(t, poolBytes), Coins: coin(165)},  // module: not attributed
				{Address: accBech32(t, wzetaBytes), Coins: coin(300)}, // pinned: not attributed
			},
		},
		Staking: stakingtypes.GenesisState{
			Params:     stakingtypes.Params{BondDenom: BaseDenom},
			Validators: []stakingtypes.Validator{validator},
			Delegations: []stakingtypes.Delegation{
				// self-delegation: floor(40 * 90/100) = 36
				{DelegatorAddress: accBech32(t, valBytes), ValidatorAddress: valBech32(t, valBytes), Shares: sdkmath.LegacyNewDec(40)},
				// external: floor(33 * 90/100) = floor(29.7) = 29
				{DelegatorAddress: accBech32(t, delBytes), ValidatorAddress: valBech32(t, valBytes), Shares: sdkmath.LegacyNewDec(33)},
				// zero address: floor(27 * 90/100) = floor(24.3) = 24, not attributed
				{DelegatorAddress: accBech32(t, zeroBytes), ValidatorAddress: valBech32(t, valBytes), Shares: sdkmath.LegacyNewDec(27)},
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
	require.Equal(t, "65", res.TotalStaked.String())   // 36+29
	require.Equal(t, "75", res.TotalUnbonding.String())

	// invariants
	require.Equal(t, "2175", res.Supply.String())
	require.Equal(t, "1850", res.AttributedTotal().String())
	require.Equal(t, "325", res.Remainder.String())
	require.False(t, res.Remainder.IsNegative())
	require.Equal(t, res.Supply.String(), res.AttributedTotal().Add(res.Remainder).String())

	// remainder by category: 90 tokens - (36+29+24) staked leaves 1 of share dust
	require.Equal(t, "0", res.Breakdown.ModuleAccounts.String(), "the staking pool is in no category")
	require.Equal(t, "24", res.Breakdown.ZeroAddress.String())
	require.Equal(t, "0", res.Breakdown.NonStandard.String())
	require.Equal(t, "300", res.Breakdown.Pinned.String())
	require.Equal(t, "1", res.Breakdown.ShareDust.String())
	require.Equal(t, res.Remainder.String(), res.Breakdown.Total().String())
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
		Staking: stakingtypes.GenesisState{Params: stakingtypes.Params{BondDenom: BaseDenom}},
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
	require.Equal(t, "400", res.Breakdown.ZeroAddress.String())
	require.Equal(t, "0", res.Breakdown.ModuleAccounts.String())
	require.Equal(t, "0", res.Breakdown.NonStandard.String())
	require.Equal(t, "0", res.Breakdown.Pinned.String())
	require.Equal(t, "0", res.Breakdown.ShareDust.String())
}

// TestComputePinnedAccounting checks that a pin holding liquid, staked and
// unbonding azeta is kept out of attribution while its full contribution to the
// remainder is reported. The pin is passed twice (hex + bech32) to prove the two
// forms collapse to one entry.
func TestComputePinnedAccounting(t *testing.T) {
	// ARRANGE
	valBytes := addr20(0x02)
	eoaBytes := addr20(0x01)
	pinBytes := addr20(0x05)

	// slashed validator: 18 tokens backing 20 shares -> exchange rate 0.9
	validator := stakingtypes.Validator{
		OperatorAddress: valBech32(t, valBytes),
		Tokens:          sdkmath.NewInt(18),
		DelegatorShares: sdkmath.LegacyNewDec(20),
		Status:          stakingtypes.Bonded,
	}
	coin := func(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }

	gen := &Genesis{
		Bank: banktypes.GenesisState{
			Supply: coin(928),
			Balances: []banktypes.Balance{
				{Address: accBech32(t, eoaBytes), Coins: coin(600)},
				{Address: accBech32(t, pinBytes), Coins: coin(300)},
				{Address: moduleBech32(t, stakingtypes.BondedPoolName), Coins: coin(18)},
				{Address: moduleBech32(t, stakingtypes.NotBondedPoolName), Coins: coin(10)},
			},
		},
		Staking: stakingtypes.GenesisState{
			Params:     stakingtypes.Params{BondDenom: BaseDenom},
			Validators: []stakingtypes.Validator{validator},
			Delegations: []stakingtypes.Delegation{
				// pinned delegator: floor(20 * 18/20) = 18
				{
					DelegatorAddress: accBech32(t, pinBytes),
					ValidatorAddress: valBech32(t, valBytes),
					Shares:           sdkmath.LegacyNewDec(20),
				},
			},
			UnbondingDelegations: []stakingtypes.UnbondingDelegation{
				{
					DelegatorAddress: accBech32(t, pinBytes),
					ValidatorAddress: valBech32(t, valBytes),
					Entries: []stakingtypes.UnbondingDelegationEntry{
						{Balance: sdkmath.NewInt(6)},
						{Balance: sdkmath.NewInt(4)},
					},
				},
			},
		},
	}

	cfg := Config{Denom: BaseDenom, Pins: []string{canonHex(pinBytes), accBech32(t, pinBytes)}}

	// ACT
	res, err := Compute(gen, cfg)

	// ASSERT
	require.NoError(t, err)
	require.Len(t, res.Accounts, 1, "only the EOA is attributed")
	require.Equal(t, canonHex(eoaBytes), res.Accounts[0].Canonical)

	require.Len(t, res.Pinned, 1, "hex and bech32 forms of the same pin collapse")
	require.Equal(t, canonHex(pinBytes), res.Pinned[0].Canonical)
	require.Equal(t, "328", res.Pinned[0].Azeta.String(), "liquid 300 + staked 18 + unbonding 10")

	require.Equal(t, "600", res.AttributedTotal().String())
	require.Equal(t, "328", res.Remainder.String())
	require.Equal(t, res.Supply.String(), res.AttributedTotal().Add(res.Remainder).String())
	require.Equal(t, "328", res.Breakdown.Pinned.String())
	require.Equal(t, "0", res.Breakdown.ModuleAccounts.String())
	require.Equal(t, "0", res.Breakdown.ZeroAddress.String())
	require.Equal(t, "0", res.Breakdown.NonStandard.String())
	require.Equal(t, "0", res.Breakdown.ShareDust.String())
}

// TestComputePinNoHolder covers the typo case: a pin that matches nothing in the
// export is still reported, with zero azeta, so the caller can fail the run.
func TestComputePinNoHolder(t *testing.T) {
	// ARRANGE
	eoaBytes := addr20(0x01)
	missingBytes := addr20(0x07)
	coin := func(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }

	gen := &Genesis{
		Bank: banktypes.GenesisState{
			Supply: coin(1000),
			Balances: []banktypes.Balance{
				{Address: accBech32(t, eoaBytes), Coins: coin(600)},
				{Address: moduleBech32(t, authtypes.FeeCollectorName), Coins: coin(400)},
			},
		},
		Staking: stakingtypes.GenesisState{Params: stakingtypes.Params{BondDenom: BaseDenom}},
	}

	// ACT
	res, err := Compute(gen, Config{Denom: BaseDenom, Pins: []string{canonHex(missingBytes)}})

	// ASSERT
	require.NoError(t, err)
	require.Len(t, res.Pinned, 1)
	require.Equal(t, canonHex(missingBytes), res.Pinned[0].Canonical)
	require.True(t, res.Pinned[0].Azeta.IsZero(), "unmatched pin reports zero azeta")

	require.Len(t, res.Accounts, 1)
	require.Equal(t, "600", res.AttributedTotal().String())
	require.Equal(t, "400", res.Remainder.String())
}

// TestComputePinnedOrder pins three addresses and checks the reported order is
// amount descending with the canonical address breaking ties.
func TestComputePinnedOrder(t *testing.T) {
	// ARRANGE
	bigBytes := addr20(0x0c)
	tieLowBytes := addr20(0x0a)
	tieHighBytes := addr20(0x0b)
	coin := func(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }

	gen := &Genesis{
		Bank: banktypes.GenesisState{
			Supply: coin(110),
			Balances: []banktypes.Balance{
				{Address: accBech32(t, tieHighBytes), Coins: coin(5)},
				{Address: accBech32(t, bigBytes), Coins: coin(100)},
				{Address: accBech32(t, tieLowBytes), Coins: coin(5)},
			},
		},
		Staking: stakingtypes.GenesisState{Params: stakingtypes.Params{BondDenom: BaseDenom}},
	}
	cfg := Config{Denom: BaseDenom, Pins: []string{
		canonHex(tieHighBytes), canonHex(bigBytes), canonHex(tieLowBytes),
	}}

	// ACT
	res, err := Compute(gen, cfg)

	// ASSERT
	require.NoError(t, err)
	require.Empty(t, res.Accounts)
	require.Equal(t, []PinnedAddr{
		{Canonical: canonHex(bigBytes), Azeta: sdkmath.NewInt(100)},
		{Canonical: canonHex(tieLowBytes), Azeta: sdkmath.NewInt(5)},
		{Canonical: canonHex(tieHighBytes), Azeta: sdkmath.NewInt(5)},
	}, res.Pinned)
}

func TestComputeUnknownValidator(t *testing.T) {
	// ARRANGE: a delegation pointing at a validator not present in the set
	delegator, validator := accBech32(t, addr20(0x03)), valBech32(t, addr20(0x09))
	gen := &Genesis{
		Bank: banktypes.GenesisState{Supply: sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, 100))},
		Staking: stakingtypes.GenesisState{
			Params: stakingtypes.Params{BondDenom: BaseDenom},
			Delegations: []stakingtypes.Delegation{
				{DelegatorAddress: delegator, ValidatorAddress: validator, Shares: sdkmath.LegacyNewDec(10)},
			},
		},
	}

	// ACT
	_, err := Compute(gen, Config{Denom: BaseDenom})

	// ASSERT
	require.ErrorContains(t, err, "delegation from "+delegator+" references unknown validator "+validator)
}

// TestComputeRemainderCategories gives every non-claimable kind of holder liquid,
// stake and unbonding at once, and checks each lands in its own category and
// the categories add up to the remainder.
func TestComputeRemainderCategories(t *testing.T) {
	// ARRANGE
	valBytes := addr20(0x02)
	eoaBytes := addr20(0x01)
	pinBytes := addr20(0x05)
	zeroBytes := make([]byte, addressLen)
	nonStdBytes := make([]byte, 32)
	for i := range nonStdBytes {
		nonStdBytes[i] = 0x0e
	}
	coin := func(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }

	pin := accBech32(t, pinBytes)
	module := moduleBech32(t, authtypes.FeeCollectorName)
	zero := accBech32(t, zeroBytes)
	nonStd := accBech32(t, nonStdBytes)
	eoa := accBech32(t, eoaBytes)
	val := valBech32(t, valBytes)

	delegate := func(addr string, shares int64) stakingtypes.Delegation {
		return stakingtypes.Delegation{DelegatorAddress: addr, ValidatorAddress: val, Shares: sdkmath.LegacyNewDec(shares)}
	}
	unbond := func(addr string, amt int64) stakingtypes.UnbondingDelegation {
		return stakingtypes.UnbondingDelegation{
			DelegatorAddress: addr,
			ValidatorAddress: val,
			Entries:          []stakingtypes.UnbondingDelegationEntry{{Balance: sdkmath.NewInt(amt)}},
		}
	}

	// slashed validator: 90 tokens backing 100 shares -> exchange rate 0.9
	gen := &Genesis{
		Bank: banktypes.GenesisState{
			Supply: coin(1615),
			Balances: []banktypes.Balance{
				{Address: pin, Coins: coin(100)},
				{Address: module, Coins: coin(200)},
				{Address: zero, Coins: coin(300)},
				{Address: nonStd, Coins: coin(400)},
				{Address: eoa, Coins: coin(500)},
				{Address: moduleBech32(t, stakingtypes.BondedPoolName), Coins: coin(90)},
				{Address: moduleBech32(t, stakingtypes.NotBondedPoolName), Coins: coin(25)},
			},
		},
		Staking: stakingtypes.GenesisState{
			Params: stakingtypes.Params{BondDenom: BaseDenom},
			Validators: []stakingtypes.Validator{{
				OperatorAddress: val,
				Tokens:          sdkmath.NewInt(90),
				DelegatorShares: sdkmath.LegacyNewDec(100),
				Status:          stakingtypes.Bonded,
			}},
			// stakes 22, 22, 22, 13, 9 = 88 of 90 tokens -> share dust 2
			Delegations: []stakingtypes.Delegation{
				delegate(pin, 25), delegate(module, 25), delegate(zero, 25), delegate(nonStd, 15), delegate(eoa, 10),
			},
			UnbondingDelegations: []stakingtypes.UnbondingDelegation{
				unbond(pin, 3), unbond(module, 4), unbond(zero, 5), unbond(nonStd, 6), unbond(eoa, 7),
			},
		},
	}

	// ACT
	res, err := Compute(gen, Config{Denom: BaseDenom, Pins: []string{pin}})

	// ASSERT
	require.NoError(t, err)
	require.Len(t, res.Accounts, 1)
	require.Equal(t, "516", res.AttributedTotal().String())
	require.Equal(t, "1099", res.Remainder.String())

	require.Equal(t, "125", res.Breakdown.Pinned.String())
	require.Equal(t, "226", res.Breakdown.ModuleAccounts.String())
	require.Equal(t, "327", res.Breakdown.ZeroAddress.String())
	require.Equal(t, "419", res.Breakdown.NonStandard.String())
	require.Equal(t, "2", res.Breakdown.ShareDust.String())
	require.Equal(t, res.Remainder.String(), res.Breakdown.Total().String())

	require.Equal(t, 1, res.NonStandard, "non-standard count stays liquid-only")
	require.Equal(t, "400", res.NonStandardAzeta.String(), "non-standard azeta stays liquid-only")
}

func TestComputeRefusesStakingPoolPin(t *testing.T) {
	// ARRANGE
	pool := moduleBech32(t, stakingtypes.BondedPoolName)
	gen := &Genesis{Staking: stakingtypes.GenesisState{Params: stakingtypes.Params{BondDenom: BaseDenom}}}

	// ACT
	_, err := Compute(gen, Config{Denom: BaseDenom, Pins: []string{pool}})

	// ASSERT
	require.ErrorContains(t, err, "pin "+pool+" is the bonded_tokens_pool staking pool")
}

// TestComputeRefusesSupplyMismatch records a supply 1 above the bank balances.
func TestComputeRefusesSupplyMismatch(t *testing.T) {
	// ARRANGE
	coin := func(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }
	gen := &Genesis{
		Bank: banktypes.GenesisState{
			Supply: coin(1001),
			Balances: []banktypes.Balance{
				{Address: accBech32(t, addr20(0x01)), Coins: coin(600)},
				{Address: accBech32(t, make([]byte, addressLen)), Coins: coin(400)},
			},
		},
		Staking: stakingtypes.GenesisState{Params: stakingtypes.Params{BondDenom: BaseDenom}},
	}

	// ACT
	_, err := Compute(gen, Config{Denom: BaseDenom})

	// ASSERT
	require.ErrorContains(t, err, "bank balances sum 1000 != supply 1001")
}

// TestFinalizeRefusesUnexplainedRemainder hands finalize a breakdown that misses
// part of the remainder. Compute cannot produce one once the supply matches the
// bank and the staking records are validated, so the identity check is a
// backstop and is reached here directly.
func TestFinalizeRefusesUnexplainedRemainder(t *testing.T) {
	// ARRANGE
	coin := func(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }
	gen := &Genesis{
		Bank: banktypes.GenesisState{
			Supply:   coin(400),
			Balances: []banktypes.Balance{{Address: accBech32(t, make([]byte, addressLen)), Coins: coin(400)}},
		},
	}
	res := &Result{Breakdown: RemainderBreakdown{
		ModuleAccounts: sdkmath.ZeroInt(),
		ZeroAddress:    sdkmath.NewInt(399),
		NonStandard:    sdkmath.ZeroInt(),
		ShareDust:      sdkmath.ZeroInt(),
	}}

	// ACT
	_, err := finalize(res, map[string]*Account{}, map[string]sdkmath.Int{}, gen, BaseDenom)

	// ASSERT
	require.ErrorContains(t, err, "remainder 400 != categorised 399 (modules 0, zero address 399, non-20-byte 0")
}

// TestComputePinOverridesCategory pins the zero address and a module account,
// each holding liquid, stake and unbonding, and checks a pin takes precedence
// over both categories.
func TestComputePinOverridesCategory(t *testing.T) {
	// ARRANGE
	coin := func(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }
	zero := accBech32(t, make([]byte, addressLen))
	module := moduleBech32(t, authtypes.FeeCollectorName)
	val := valBech32(t, addr20(0x02))

	gen := &Genesis{
		Bank: banktypes.GenesisState{
			Supply: coin(368),
			Balances: []banktypes.Balance{
				{Address: zero, Coins: coin(100)},
				{Address: module, Coins: coin(200)},
				{Address: moduleBech32(t, stakingtypes.BondedPoolName), Coins: coin(50)},
				{Address: moduleBech32(t, stakingtypes.NotBondedPoolName), Coins: coin(18)},
			},
		},
		Staking: stakingtypes.GenesisState{
			Params: stakingtypes.Params{BondDenom: BaseDenom},
			Validators: []stakingtypes.Validator{{
				OperatorAddress: val,
				Tokens:          sdkmath.NewInt(50),
				DelegatorShares: sdkmath.LegacyNewDec(50),
				Status:          stakingtypes.Bonded,
			}},
			Delegations: []stakingtypes.Delegation{
				{DelegatorAddress: zero, ValidatorAddress: val, Shares: sdkmath.LegacyNewDec(20)},
				{DelegatorAddress: module, ValidatorAddress: val, Shares: sdkmath.LegacyNewDec(30)},
			},
			UnbondingDelegations: []stakingtypes.UnbondingDelegation{
				{
					DelegatorAddress: zero,
					ValidatorAddress: val,
					Entries:          []stakingtypes.UnbondingDelegationEntry{{Balance: sdkmath.NewInt(7)}},
				},
				{
					DelegatorAddress: module,
					ValidatorAddress: val,
					Entries:          []stakingtypes.UnbondingDelegationEntry{{Balance: sdkmath.NewInt(11)}},
				},
			},
		},
	}

	// ACT
	res, err := Compute(gen, Config{Denom: BaseDenom, Pins: []string{zero, module}})

	// ASSERT
	require.NoError(t, err)
	require.Empty(t, res.Accounts)
	require.Equal(t, []PinnedAddr{
		{Canonical: canonHex(authtypes.NewModuleAddress(authtypes.FeeCollectorName).Bytes()), Azeta: sdkmath.NewInt(241)},
		{Canonical: zeroAddress, Azeta: sdkmath.NewInt(127)},
	}, res.Pinned, "module 200+30+11, zero address 100+20+7")
	require.Equal(t, "368", res.Breakdown.Pinned.String())
	require.Equal(t, "0", res.Breakdown.ModuleAccounts.String())
	require.Equal(t, "0", res.Breakdown.ZeroAddress.String())
	require.Equal(t, "368", res.Remainder.String())
	require.Equal(t, res.Remainder.String(), res.Breakdown.Total().String())
}

func TestComputeConversionDust(t *testing.T) {
	// ARRANGE
	coin := func(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }
	gen := &Genesis{
		Bank: banktypes.GenesisState{
			Supply: coin(1_500_000_700),
			Balances: []banktypes.Balance{
				{Address: accBech32(t, addr20(0x01)), Coins: coin(1_500_000_000)},
				{Address: accBech32(t, make([]byte, addressLen)), Coins: coin(700)},
			},
		},
		Staking: stakingtypes.GenesisState{Params: stakingtypes.Params{BondDenom: BaseDenom}},
	}

	// ACT
	res, err := Compute(gen, Config{Denom: BaseDenom})

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, "1", res.SPLCap().String())
	require.Equal(t, "500000700", res.ConversionDust().String())
}

// TestComputeAuditExports runs the staking exports from the audit. Complete
// exports attribute the holder's stake or unbonding; dropping the delegation
// or the unbonding record is refused instead of moving it into the remainder.
func TestComputeAuditExports(t *testing.T) {
	const (
		holder = "0x1111111111111111111111111111111111111111"
		pin    = "0x3333333333333333333333333333333333333333"
	)
	tests := []struct {
		name    string
		wantErr string
	}{
		{name: "complete"},
		{name: "complete_unbonding"},
		{
			name: "omitted_delegations",
			wantErr: "validator 0x2222222222222222222222222222222222222222 delegation shares sum 0.000000000000000000 " +
				"!= delegator_shares 100000000000000000000",
		},
		{
			name:    "omitted_unbonding",
			wantErr: "not-bonded pool balance 100000000000000000000 != unbonding/unbonded validator tokens 0 + unbonding entries 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			gen := loadExport(t, tt.name)

			// ACT
			res, err := Compute(gen, Config{Denom: BaseDenom, Pins: []string{pin}})

			// ASSERT
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, res.Accounts, 1)
			require.Equal(t, holder, res.Accounts[0].Canonical)
			require.Equal(t, "110000000000000000000", res.Accounts[0].Total().String())
			require.Equal(t, "5000000000000000000", res.Remainder.String())
			require.Equal(t, "5000000000000000000", res.Breakdown.Pinned.String())
		})
	}
}
