package snapshot

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// validateDistribution refuses an export that still holds unpaid staking
// rewards. A --for-zero-height export pays out every reward and commission and
// resets the per-validator records to zero, so any non-zero record means the
// rewards were never credited to a holder. The distribution module's balance
// must then be exactly the community pool (the SDK's own import rule), which
// catches an export whose reward records were dropped while the azeta stayed.
func validateDistribution(gen *Genesis, denom string) error {
	distr := gen.Distribution
	for _, r := range distr.OutstandingRewards {
		if !r.OutstandingRewards.IsZero() {
			return fmt.Errorf(
				"validator %s has non-zero outstanding_rewards %s; the export is not --for-zero-height",
				r.ValidatorAddress, r.OutstandingRewards,
			)
		}
	}
	for _, r := range distr.ValidatorCurrentRewards {
		if !r.Rewards.Rewards.IsZero() {
			return fmt.Errorf(
				"validator %s has non-zero current_rewards %s; the export is not --for-zero-height",
				r.ValidatorAddress, r.Rewards.Rewards,
			)
		}
	}
	for _, r := range distr.ValidatorAccumulatedCommissions {
		if !r.Accumulated.Commission.IsZero() {
			return fmt.Errorf(
				"validator %s has non-zero accumulated_commission %s; the export is not --for-zero-height",
				r.ValidatorAddress, r.Accumulated.Commission,
			)
		}
	}

	// outstanding rewards are all zero by now, so the SDK's
	// TruncateDecimal(outstanding + community pool) is just the community pool
	holdings, _ := distr.FeePool.CommunityPool.TruncateDecimal()
	balance, err := moduleBalance(gen, distrtypes.ModuleName, denom)
	if err != nil {
		return err
	}
	if want := holdings.AmountOf(denom); !balance.Equal(want) {
		return fmt.Errorf("distribution module balance %s != outstanding rewards + community pool %s", balance, want)
	}
	return nil
}

// stakingTally accumulates one validator's delegations during validateStaking.
type stakingTally struct {
	val         stakingtypes.Validator
	shares      sdkmath.LegacyDec
	delegations int
	staked      sdkmath.Int
}

// delegationKey identifies a (delegator, validator) staking record. The
// delegator is its raw address bytes so non-20-byte holders are keyed too.
type delegationKey struct {
	delegator string
	validator string
}

// validateStaking proves the staking records are complete: every delegation
// and unbonding entry the chain holds is present, so none of their azeta can
// silently fall into the remainder. It checks the full population, including
// module accounts, pins and non-20-byte holders. It returns the total
// share-rounding dust: validator tokens no delegator is credited with because
// each delegation's stake is floored.
func validateStaking(gen *Genesis, denom string) (sdkmath.Int, error) {
	st := gen.Staking
	if st.Params.BondDenom != denom {
		return sdkmath.Int{}, fmt.Errorf("staking bond_denom %q != snapshot denom %q", st.Params.BondDenom, denom)
	}

	tallies := make([]*stakingTally, 0, len(st.Validators))
	byOperator := make(map[string]*stakingTally, len(st.Validators))
	bondedTokens, notBondedTokens := sdkmath.ZeroInt(), sdkmath.ZeroInt()
	for i := range st.Validators {
		v := st.Validators[i]
		canon, err := Canonical(v.OperatorAddress)
		if err != nil {
			return sdkmath.Int{}, fmt.Errorf("canonicalize validator %q: %w", v.OperatorAddress, err)
		}
		if _, dup := byOperator[canon]; dup {
			return sdkmath.Int{}, fmt.Errorf("duplicate validator %s", v.OperatorAddress)
		}
		if v.Tokens.IsNegative() {
			return sdkmath.Int{}, fmt.Errorf("validator %s has negative tokens %s", v.OperatorAddress, v.Tokens)
		}
		if v.DelegatorShares.IsNegative() {
			return sdkmath.Int{}, fmt.Errorf(
				"validator %s has negative delegator_shares %s",
				v.OperatorAddress, v.DelegatorShares,
			)
		}
		switch v.Status {
		case stakingtypes.Bonded:
			bondedTokens = bondedTokens.Add(v.Tokens)
		case stakingtypes.Unbonding, stakingtypes.Unbonded:
			notBondedTokens = notBondedTokens.Add(v.Tokens)
		default:
			return sdkmath.Int{}, fmt.Errorf("validator %s has invalid status %s", v.OperatorAddress, v.Status)
		}
		tally := &stakingTally{val: v, shares: sdkmath.LegacyZeroDec(), staked: sdkmath.ZeroInt()}
		tallies = append(tallies, tally)
		byOperator[canon] = tally
	}

	owners := make([]*stakingTally, len(st.Delegations))
	seen := make(map[delegationKey]bool, len(st.Delegations))
	for i := range st.Delegations {
		del := st.Delegations[i]
		key, err := stakingRecordKey(del.DelegatorAddress, del.ValidatorAddress)
		if err != nil {
			return sdkmath.Int{}, err
		}
		if del.Shares.IsNegative() {
			return sdkmath.Int{}, fmt.Errorf(
				"delegation %s -> %s has negative shares %s",
				del.DelegatorAddress, del.ValidatorAddress, del.Shares,
			)
		}
		tally, ok := byOperator[key.validator]
		if !ok {
			return sdkmath.Int{}, fmt.Errorf(
				"delegation from %s references unknown validator %s",
				del.DelegatorAddress, del.ValidatorAddress,
			)
		}
		if seen[key] {
			return sdkmath.Int{}, fmt.Errorf(
				"duplicate delegation %s -> %s",
				del.DelegatorAddress, del.ValidatorAddress,
			)
		}
		seen[key] = true
		tally.shares = tally.shares.Add(del.Shares)
		tally.delegations++
		owners[i] = tally
	}

	// every share sum must match before any stake is computed: TokensFromShares
	// divides by delegator_shares, which a mismatched validator may have at zero
	for _, tally := range tallies {
		if !tally.shares.Equal(tally.val.DelegatorShares) {
			return sdkmath.Int{}, fmt.Errorf(
				"validator %s delegation shares sum %s != delegator_shares %s",
				tally.val.OperatorAddress, tally.shares, tally.val.DelegatorShares,
			)
		}
	}
	for i := range st.Delegations {
		tally := owners[i]
		if tally.val.DelegatorShares.IsZero() {
			continue
		}
		staked := tally.val.TokensFromShares(st.Delegations[i].Shares).TruncateInt()
		tally.staked = tally.staked.Add(staked)
	}
	shareDust := sdkmath.ZeroInt()
	for _, tally := range tallies {
		dust := tally.val.Tokens.Sub(tally.staked)
		if dust.GT(sdkmath.NewInt(int64(tally.delegations))) {
			return sdkmath.Int{}, fmt.Errorf(
				"validator %s share-rounding dust %s exceeds its %d delegations",
				tally.val.OperatorAddress, dust, tally.delegations,
			)
		}
		shareDust = shareDust.Add(dust)
	}

	// unbonding records may point at a validator that has since been removed
	unbonding := sdkmath.ZeroInt()
	seen = make(map[delegationKey]bool, len(st.UnbondingDelegations))
	for i := range st.UnbondingDelegations {
		ubd := st.UnbondingDelegations[i]
		key, err := stakingRecordKey(ubd.DelegatorAddress, ubd.ValidatorAddress)
		if err != nil {
			return sdkmath.Int{}, err
		}
		if seen[key] {
			return sdkmath.Int{}, fmt.Errorf(
				"duplicate unbonding delegation %s -> %s",
				ubd.DelegatorAddress, ubd.ValidatorAddress,
			)
		}
		seen[key] = true
		for _, entry := range ubd.Entries {
			if entry.Balance.IsNegative() {
				return sdkmath.Int{}, fmt.Errorf(
					"unbonding delegation %s -> %s has negative entry balance %s",
					ubd.DelegatorAddress, ubd.ValidatorAddress, entry.Balance,
				)
			}
			unbonding = unbonding.Add(entry.Balance)
		}
	}

	bondedPool, err := moduleBalance(gen, stakingtypes.BondedPoolName, denom)
	if err != nil {
		return sdkmath.Int{}, err
	}
	if !bondedPool.Equal(bondedTokens) {
		return sdkmath.Int{}, fmt.Errorf(
			"bonded pool balance %s != bonded validator tokens %s",
			bondedPool, bondedTokens,
		)
	}
	notBondedPool, err := moduleBalance(gen, stakingtypes.NotBondedPoolName, denom)
	if err != nil {
		return sdkmath.Int{}, err
	}
	if !notBondedPool.Equal(notBondedTokens.Add(unbonding)) {
		return sdkmath.Int{}, fmt.Errorf(
			"not-bonded pool balance %s != unbonding/unbonded validator tokens %s + unbonding entries %s",
			notBondedPool, notBondedTokens, unbonding,
		)
	}
	return shareDust, nil
}

// stakingRecordKey decodes a delegation or unbonding record's addresses into
// its (delegator, validator) key.
func stakingRecordKey(delegator, validator string) (delegationKey, error) {
	raw, err := decodeAddr(delegator)
	if err != nil {
		return delegationKey{}, fmt.Errorf("decode delegator %q: %w", delegator, err)
	}
	valCanon, err := Canonical(validator)
	if err != nil {
		return delegationKey{}, fmt.Errorf("canonicalize validator %q: %w", validator, err)
	}
	return delegationKey{delegator: string(raw), validator: valCanon}, nil
}

// moduleBalance sums the denom balance of a module account across every bank
// row that resolves to it; an absent row is a zero balance.
func moduleBalance(gen *Genesis, moduleName, denom string) (sdkmath.Int, error) {
	want, err := moduleCanonical(moduleName)
	if err != nil {
		return sdkmath.Int{}, err
	}
	total := sdkmath.ZeroInt()
	for i := range gen.Bank.Balances {
		bal := gen.Bank.Balances[i]
		canon, ok, err := Classify(bal.Address)
		if err != nil {
			return sdkmath.Int{}, fmt.Errorf("classify balance holder %q: %w", bal.Address, err)
		}
		if ok && canon == want {
			total = total.Add(bal.Coins.AmountOf(denom))
		}
	}
	return total, nil
}
