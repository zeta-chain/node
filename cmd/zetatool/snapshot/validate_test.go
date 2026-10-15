package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/app"
)

func azetaCoins(n int64) sdk.Coins { return sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, n)) }

func moduleBech32(t *testing.T, name string) string {
	return accBech32(t, authtypes.NewModuleAddress(name).Bytes())
}

// loadExport parses a ported audit export from testdata.
func loadExport(t *testing.T, name string) *Genesis {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", name+".json"))
	require.NoError(t, err)
	var doc struct {
		AppState map[string]json.RawMessage `json:"app_state"`
	}
	require.NoError(t, json.Unmarshal(raw, &doc))
	gen, err := ParseAppState(app.MakeEncodingConfig(7000).Codec, doc.AppState)
	require.NoError(t, err)
	return gen
}

// stakingFixture is an internally consistent export. Validator 0x02 is bonded
// and slashed (90 tokens backing 100 shares) with two fractional delegations
// whose floored stakes leave 1 azeta of share dust; validator 0x06 is unbonded
// with one delegation; 0x07 is unbonding 15 from 0x02. Both pools match.
func stakingFixture(t *testing.T) *Genesis {
	valA, valB := valBech32(t, addr20(0x02)), valBech32(t, addr20(0x06))
	return &Genesis{
		Bank: banktypes.GenesisState{
			Supply: azetaCoins(155),
			Balances: []banktypes.Balance{
				{Address: moduleBech32(t, stakingtypes.BondedPoolName), Coins: azetaCoins(90)},
				{Address: moduleBech32(t, stakingtypes.NotBondedPoolName), Coins: azetaCoins(65)},
			},
		},
		Staking: stakingtypes.GenesisState{
			Params: stakingtypes.Params{BondDenom: BaseDenom},
			Validators: []stakingtypes.Validator{
				{
					OperatorAddress: valA,
					Tokens:          sdkmath.NewInt(90),
					DelegatorShares: sdkmath.LegacyNewDec(100),
					Status:          stakingtypes.Bonded,
				},
				{
					OperatorAddress: valB,
					Tokens:          sdkmath.NewInt(50),
					DelegatorShares: sdkmath.LegacyNewDec(50),
					Status:          stakingtypes.Unbonded,
				},
			},
			Delegations: []stakingtypes.Delegation{
				{
					DelegatorAddress: accBech32(t, addr20(0x03)),
					ValidatorAddress: valA,
					Shares:           sdkmath.LegacyMustNewDecFromStr("33.5"),
				},
				{
					DelegatorAddress: accBech32(t, addr20(0x04)),
					ValidatorAddress: valA,
					Shares:           sdkmath.LegacyMustNewDecFromStr("66.5"),
				},
				{
					DelegatorAddress: accBech32(t, addr20(0x03)),
					ValidatorAddress: valB,
					Shares:           sdkmath.LegacyNewDec(50),
				},
			},
			UnbondingDelegations: []stakingtypes.UnbondingDelegation{
				{
					DelegatorAddress: accBech32(t, addr20(0x07)),
					ValidatorAddress: valA,
					Entries: []stakingtypes.UnbondingDelegationEntry{
						{Balance: sdkmath.NewInt(10)},
						{Balance: sdkmath.NewInt(5)},
					},
				},
			},
		},
	}
}

func TestValidateStakingAccepts(t *testing.T) {
	valGone := valBech32(t, addr20(0x09))
	tests := []struct {
		name   string
		mutate func(g *Genesis)
	}{
		{name: "complete export", mutate: func(*Genesis) {}},
		{name: "jailed validator", mutate: func(g *Genesis) { g.Staking.Validators[0].Jailed = true }},
		{name: "unbonding validator", mutate: func(g *Genesis) { g.Staking.Validators[1].Status = stakingtypes.Unbonding }},
		{
			name: "zero-token zero-share validator without delegations",
			mutate: func(g *Genesis) {
				g.Staking.Validators = append(g.Staking.Validators, stakingtypes.Validator{
					OperatorAddress: valBech32(t, addr20(0x08)),
					Tokens:          sdkmath.ZeroInt(),
					DelegatorShares: sdkmath.LegacyZeroDec(),
					Status:          stakingtypes.Unbonded,
				})
			},
		},
		{
			name: "unbonding from a validator absent from the export",
			mutate: func(g *Genesis) {
				g.Staking.UnbondingDelegations = append(g.Staking.UnbondingDelegations, stakingtypes.UnbondingDelegation{
					DelegatorAddress: accBech32(t, addr20(0x07)),
					ValidatorAddress: valGone,
					Entries:          []stakingtypes.UnbondingDelegationEntry{{Balance: sdkmath.NewInt(4)}},
				})
				g.Bank.Balances[1].Coins = azetaCoins(69)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			gen := stakingFixture(t)
			tt.mutate(gen)

			// ACT
			dust, err := validateStaking(gen, BaseDenom)

			// ASSERT
			require.NoError(t, err)
			require.Equal(t, "1", dust.String())
		})
	}
}

func TestValidateStakingRefuses(t *testing.T) {
	valA, valB := valBech32(t, addr20(0x02)), valBech32(t, addr20(0x06))
	del3, del4, unb7 := accBech32(t, addr20(0x03)), accBech32(t, addr20(0x04)), accBech32(t, addr20(0x07))
	valC := valBech32(t, addr20(0x08))
	valUnknown := valBech32(t, addr20(0x09))

	tests := []struct {
		name    string
		mutate  func(g *Genesis)
		wantErr string
	}{
		{
			name:    "omitted delegation",
			mutate:  func(g *Genesis) { g.Staking.Delegations = g.Staking.Delegations[1:] },
			wantErr: "validator " + valA + " delegation shares sum 66.500000000000000000 != delegator_shares 100",
		},
		{
			name: "omitted unbonding entry",
			mutate: func(g *Genesis) {
				g.Staking.UnbondingDelegations[0].Entries = g.Staking.UnbondingDelegations[0].Entries[:1]
			},
			wantErr: "not-bonded pool balance 65 != unbonding/unbonded validator tokens 50 + unbonding entries 10",
		},
		{
			name:    "bonded pool one over",
			mutate:  func(g *Genesis) { g.Bank.Balances[0].Coins = azetaCoins(91) },
			wantErr: "bonded pool balance 91 != bonded validator tokens 90",
		},
		{
			name:    "bonded pool one under",
			mutate:  func(g *Genesis) { g.Bank.Balances[0].Coins = azetaCoins(89) },
			wantErr: "bonded pool balance 89 != bonded validator tokens 90",
		},
		{
			name:    "not-bonded pool one over",
			mutate:  func(g *Genesis) { g.Bank.Balances[1].Coins = azetaCoins(66) },
			wantErr: "not-bonded pool balance 66 != unbonding/unbonded validator tokens 50 + unbonding entries 15",
		},
		{
			name:    "not-bonded pool one under",
			mutate:  func(g *Genesis) { g.Bank.Balances[1].Coins = azetaCoins(64) },
			wantErr: "not-bonded pool balance 64 != unbonding/unbonded validator tokens 50 + unbonding entries 15",
		},
		{
			name: "duplicate validator",
			mutate: func(g *Genesis) {
				g.Staking.Validators = append(g.Staking.Validators, g.Staking.Validators[1])
			},
			wantErr: "duplicate validator " + valB,
		},
		{
			name: "duplicate delegation",
			mutate: func(g *Genesis) {
				g.Staking.Delegations = append(g.Staking.Delegations, g.Staking.Delegations[0])
			},
			wantErr: "duplicate delegation " + del3 + " -> " + valA,
		},
		{
			name: "duplicate unbonding delegation",
			mutate: func(g *Genesis) {
				g.Staking.UnbondingDelegations = append(g.Staking.UnbondingDelegations, g.Staking.UnbondingDelegations[0])
			},
			wantErr: "duplicate unbonding delegation " + unb7 + " -> " + valA,
		},
		{
			name:    "unknown validator",
			mutate:  func(g *Genesis) { g.Staking.Delegations[1].ValidatorAddress = valUnknown },
			wantErr: "delegation from " + del4 + " references unknown validator " + valUnknown,
		},
		{
			name:    "invalid status",
			mutate:  func(g *Genesis) { g.Staking.Validators[1].Status = stakingtypes.Unspecified },
			wantErr: "validator " + valB + " has invalid status BOND_STATUS_UNSPECIFIED",
		},
		{
			name:    "negative delegation shares",
			mutate:  func(g *Genesis) { g.Staking.Delegations[1].Shares = sdkmath.LegacyNewDec(-1) },
			wantErr: "delegation " + del4 + " -> " + valA + " has negative shares",
		},
		{
			name: "negative unbonding entry",
			mutate: func(g *Genesis) {
				g.Staking.UnbondingDelegations[0].Entries[1].Balance = sdkmath.NewInt(-5)
			},
			wantErr: "unbonding delegation " + unb7 + " -> " + valA + " has negative entry balance -5",
		},
		{
			name:    "negative validator tokens",
			mutate:  func(g *Genesis) { g.Staking.Validators[1].Tokens = sdkmath.NewInt(-1) },
			wantErr: "validator " + valB + " has negative tokens -1",
		},
		{
			name: "negative validator delegator_shares",
			mutate: func(g *Genesis) {
				g.Staking.Validators[1].DelegatorShares = sdkmath.LegacyNewDec(-50)
			},
			wantErr: "validator " + valB + " has negative delegator_shares",
		},
		{
			name:    "wrong bond denom",
			mutate:  func(g *Genesis) { g.Staking.Params.BondDenom = "stake" },
			wantErr: `staking bond_denom "stake" != snapshot denom "azeta"`,
		},
		{
			name: "tokens without shares or delegations",
			mutate: func(g *Genesis) {
				g.Staking.Validators = append(g.Staking.Validators, stakingtypes.Validator{
					OperatorAddress: valC,
					Tokens:          sdkmath.NewInt(5),
					DelegatorShares: sdkmath.LegacyZeroDec(),
					Status:          stakingtypes.Unbonded,
				})
				g.Bank.Balances[1].Coins = azetaCoins(70)
			},
			wantErr: "validator " + valC + " share-rounding dust 5 exceeds its 0 delegations",
		},
		{
			name: "zero-share validator with a delegation",
			mutate: func(g *Genesis) {
				g.Staking.Validators = append(g.Staking.Validators, stakingtypes.Validator{
					OperatorAddress: valC,
					Tokens:          sdkmath.ZeroInt(),
					DelegatorShares: sdkmath.LegacyZeroDec(),
					Status:          stakingtypes.Unbonded,
				})
				g.Staking.Delegations = append(g.Staking.Delegations, stakingtypes.Delegation{
					DelegatorAddress: del3,
					ValidatorAddress: valC,
					Shares:           sdkmath.LegacyOneDec(),
				})
			},
			wantErr: "validator " + valC + " delegation shares sum 1.000000000000000000 != delegator_shares 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			gen := stakingFixture(t)
			tt.mutate(gen)

			// ACT
			_, err := validateStaking(gen, BaseDenom)

			// ASSERT
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

// distributionFixture is a valid --for-zero-height distribution state: every
// per-validator record is zero and the module holds the truncated community
// pool (10.5 -> 10).
func distributionFixture(t *testing.T) *Genesis {
	valA := valBech32(t, addr20(0x02))
	return &Genesis{
		Bank: banktypes.GenesisState{
			Balances: []banktypes.Balance{
				{Address: moduleBech32(t, distrtypes.ModuleName), Coins: azetaCoins(10)},
			},
		},
		Distribution: distrtypes.GenesisState{
			FeePool: distrtypes.FeePool{CommunityPool: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec(BaseDenom, sdkmath.LegacyMustNewDecFromStr("10.5")),
			)},
			OutstandingRewards: []distrtypes.ValidatorOutstandingRewardsRecord{{ValidatorAddress: valA}},
			ValidatorCurrentRewards: []distrtypes.ValidatorCurrentRewardsRecord{
				{ValidatorAddress: valA, Rewards: distrtypes.ValidatorCurrentRewards{Period: 2}},
			},
			ValidatorAccumulatedCommissions: []distrtypes.ValidatorAccumulatedCommissionRecord{
				{ValidatorAddress: valA},
			},
		},
	}
}

func TestValidateDistribution(t *testing.T) {
	valA := valBech32(t, addr20(0x02))
	reward := sdk.NewDecCoins(sdk.NewInt64DecCoin(BaseDenom, 3))

	tests := []struct {
		name    string
		mutate  func(g *Genesis)
		wantErr string
	}{
		{name: "zero-height export", mutate: func(*Genesis) {}},
		{name: "empty distribution and no module balance", mutate: func(g *Genesis) { *g = Genesis{} }},
		{
			name:    "outstanding rewards",
			mutate:  func(g *Genesis) { g.Distribution.OutstandingRewards[0].OutstandingRewards = reward },
			wantErr: "validator " + valA + " has non-zero outstanding_rewards",
		},
		{
			name:    "current rewards",
			mutate:  func(g *Genesis) { g.Distribution.ValidatorCurrentRewards[0].Rewards.Rewards = reward },
			wantErr: "validator " + valA + " has non-zero current_rewards",
		},
		{
			name: "accumulated commission",
			mutate: func(g *Genesis) {
				g.Distribution.ValidatorAccumulatedCommissions[0].Accumulated.Commission = reward
			},
			wantErr: "validator " + valA + " has non-zero accumulated_commission",
		},
		{
			name:    "module balance one over",
			mutate:  func(g *Genesis) { g.Bank.Balances[0].Coins = azetaCoins(11) },
			wantErr: "distribution module balance 11 != outstanding rewards + community pool 10",
		},
		{
			name:    "module balance one under",
			mutate:  func(g *Genesis) { g.Bank.Balances[0].Coins = azetaCoins(9) },
			wantErr: "distribution module balance 9 != outstanding rewards + community pool 10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			gen := distributionFixture(t)
			tt.mutate(gen)

			// ACT
			err := validateDistribution(gen, BaseDenom)

			// ASSERT
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

// TestValidateDistributionAuditExports runs the reward exports from the audit:
// an export taken without --for-zero-height, the same export with its reward
// records stripped while the module still holds the unpaid azeta, and a settled
// export.
func TestValidateDistributionAuditExports(t *testing.T) {
	t.Run("unwithdrawn rewards", func(t *testing.T) {
		// ARRANGE
		gen := loadExport(t, "unwithdrawn_rewards")

		// ACT
		err := validateDistribution(gen, BaseDenom)

		// ASSERT
		require.ErrorContains(t, err,
			"validator 0x2222222222222222222222222222222222222222 has non-zero outstanding_rewards")
	})

	t.Run("reward records stripped", func(t *testing.T) {
		// ARRANGE
		gen := loadExport(t, "unwithdrawn_rewards")
		gen.Distribution.OutstandingRewards = nil
		gen.Distribution.ValidatorCurrentRewards = nil
		gen.Distribution.ValidatorAccumulatedCommissions = nil

		// ACT
		err := validateDistribution(gen, BaseDenom)

		// ASSERT
		require.ErrorContains(t, err,
			"distribution module balance 10000000000000000000 != outstanding rewards + community pool 0")
	})

	t.Run("settled rewards", func(t *testing.T) {
		// ARRANGE
		gen := loadExport(t, "settled_rewards")

		// ACT
		err := validateDistribution(gen, BaseDenom)

		// ASSERT
		require.NoError(t, err)
	})
}
