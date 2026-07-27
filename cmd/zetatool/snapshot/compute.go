package snapshot

import (
	"fmt"
	"sort"

	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// SumBankDenom sums every bank balance for denom, including module accounts.
// Used for the "parser dropped nothing" check against the recorded supply.
func SumBankDenom(gen *Genesis, denom string) sdkmath.Int {
	total := sdkmath.ZeroInt()
	for i := range gen.Bank.Balances {
		total = total.Add(gen.Bank.Balances[i].Coins.AmountOf(denom))
	}
	return total
}

// nonClaimableSet builds the set of canonical addresses that must not be
// attributed: the module accounts plus the caller's pins (e.g. WZETA). The pins
// are also returned on their own, zero-seeded, so Compute can report how much
// azeta each one kept out of attribution (two pin inputs can canonicalize to the
// same account; the map dedupes them).
func nonClaimableSet(pins []string) (map[string]bool, map[string]sdkmath.Int, error) {
	set := make(map[string]bool, len(moduleAccountNames)+len(pins)+1)
	pinned := make(map[string]sdkmath.Int, len(pins))
	set[zeroAddress] = true
	for _, name := range moduleAccountNames {
		canon, err := canonicalFromBytes(authtypes.NewModuleAddress(name).Bytes())
		if err != nil {
			return nil, nil, fmt.Errorf("canonicalize module account %q: %w", name, err)
		}
		set[canon] = true
	}
	for _, pin := range pins {
		canon, err := Canonical(pin)
		if err != nil {
			return nil, nil, fmt.Errorf("canonicalize pin %q: %w", pin, err)
		}
		set[canon] = true
		if _, ok := pinned[canon]; !ok {
			pinned[canon] = sdkmath.ZeroInt()
		}
	}
	return set, pinned, nil
}

// Compute attributes native ZETA to each claimable address and folds everything
// else into the remainder. It is pure and deterministic.
//
//   - liquid    = bank balance in denom (already includes withdrawn rewards)
//   - staked    = sum of floor(shares * validator.tokens / validator.delegator_shares)
//     across delegations, using the slashing-aware exchange rate for every
//     validator status (bonded/unbonding/unbonded/jailed)
//   - unbonding = sum of unbonding-delegation entry balances
//
// remainder = supply - sum(attributed). It must be >= 0: staked/unbonding azeta
// physically sits in the (non-claimable) staking pool module accounts, so those
// pool balances cancel against the amounts credited to delegators here.
func Compute(gen *Genesis, cfg Config) (*Result, error) {
	nonClaimable, pinned, err := nonClaimableSet(cfg.Pins)
	if err != nil {
		return nil, err
	}

	// addPinned accumulates what a pin kept out of attribution. The other
	// non-claimable addresses (module accounts, the zero address) are not tracked.
	addPinned := func(canon string, amt sdkmath.Int) {
		if cur, isPin := pinned[canon]; isPin {
			pinned[canon] = cur.Add(amt)
		}
	}

	// index validators by canonical operator address for the exchange rate
	validators := make(map[string]stakingtypes.Validator, len(gen.Staking.Validators))
	for i := range gen.Staking.Validators {
		v := gen.Staking.Validators[i]
		canon, err := Canonical(v.OperatorAddress)
		if err != nil {
			return nil, fmt.Errorf("canonicalize validator %q: %w", v.OperatorAddress, err)
		}
		validators[canon] = v
	}

	accounts := make(map[string]*Account)
	get := func(canon string) *Account {
		acc, ok := accounts[canon]
		if !ok {
			class := classEOA
			if _, isVal := validators[canon]; isVal {
				class = classValidator
			}
			acc = &Account{
				Canonical: canon,
				Class:     class,
				Liquid:    sdkmath.ZeroInt(),
				Staked:    sdkmath.ZeroInt(),
				Unbonding: sdkmath.ZeroInt(),
			}
			accounts[canon] = acc
		}
		return acc
	}

	res := &Result{NonStandardAzeta: sdkmath.ZeroInt()}

	// liquid
	for i := range gen.Bank.Balances {
		bal := gen.Bank.Balances[i]
		canon, ok, err := Classify(bal.Address)
		if err != nil {
			return nil, fmt.Errorf("classify balance holder %q: %w", bal.Address, err)
		}
		amt := bal.Coins.AmountOf(cfg.Denom)
		// non-standard (non-20-byte) holders are not claimable and fall through
		// to the remainder along with their balance.
		if !ok {
			if !amt.IsZero() {
				res.NonStandard++
				res.NonStandardAzeta = res.NonStandardAzeta.Add(amt)
			}
			continue
		}
		if nonClaimable[canon] {
			addPinned(canon, amt)
			continue
		}
		if amt.IsZero() {
			continue
		}
		acc := get(canon)
		acc.Liquid = acc.Liquid.Add(amt)
	}

	// staked
	for i := range gen.Staking.Delegations {
		del := gen.Staking.Delegations[i]
		canon, ok, err := Classify(del.DelegatorAddress)
		if err != nil {
			return nil, fmt.Errorf("classify delegator %q: %w", del.DelegatorAddress, err)
		}
		// a pinned delegator skips attribution too, but only after the exchange
		// rate math below, so its stake is reported against the pin
		_, isPin := pinned[canon]
		if !ok || (nonClaimable[canon] && !isPin) {
			continue
		}
		valCanon, err := Canonical(del.ValidatorAddress)
		if err != nil {
			return nil, fmt.Errorf("canonicalize delegation validator %q: %w", del.ValidatorAddress, err)
		}
		val, ok := validators[valCanon]
		if !ok {
			return nil, fmt.Errorf("delegation references unknown validator %q", del.ValidatorAddress)
		}
		if val.DelegatorShares.IsZero() {
			return nil, fmt.Errorf("validator %q has zero delegator shares", del.ValidatorAddress)
		}
		// floor(shares * tokens / delegator_shares)
		staked := val.TokensFromShares(del.Shares).TruncateInt()
		if isPin {
			addPinned(canon, staked)
			continue
		}
		acc := get(canon)
		acc.Staked = acc.Staked.Add(staked)
	}

	// unbonding
	for i := range gen.Staking.UnbondingDelegations {
		ubd := gen.Staking.UnbondingDelegations[i]
		canon, ok, err := Classify(ubd.DelegatorAddress)
		if err != nil {
			return nil, fmt.Errorf("classify unbonding delegator %q: %w", ubd.DelegatorAddress, err)
		}
		if !ok {
			continue
		}
		sum := sdkmath.ZeroInt()
		for _, entry := range ubd.Entries {
			sum = sum.Add(entry.Balance)
		}
		if nonClaimable[canon] {
			addPinned(canon, sum)
			continue
		}
		acc := get(canon)
		acc.Unbonding = acc.Unbonding.Add(sum)
	}

	return finalize(res, accounts, pinned, gen, cfg.Denom)
}

// finalize sorts the accounts and pins, computes bucket totals and the
// remainder, and verifies the remainder is non-negative.
func finalize(
	res *Result,
	accounts map[string]*Account,
	pinned map[string]sdkmath.Int,
	gen *Genesis,
	denom string,
) (*Result, error) {
	res.Supply = gen.Bank.Supply.AmountOf(denom)
	res.TotalLiquid = sdkmath.ZeroInt()
	res.TotalStaked = sdkmath.ZeroInt()
	res.TotalUnbonding = sdkmath.ZeroInt()
	res.Accounts = make([]Account, 0, len(accounts))

	for _, acc := range accounts {
		res.TotalLiquid = res.TotalLiquid.Add(acc.Liquid)
		res.TotalStaked = res.TotalStaked.Add(acc.Staked)
		res.TotalUnbonding = res.TotalUnbonding.Add(acc.Unbonding)
		res.Accounts = append(res.Accounts, *acc)
	}
	sort.Slice(res.Accounts, func(i, j int) bool {
		return res.Accounts[i].Canonical < res.Accounts[j].Canonical
	})

	res.Pinned = make([]PinnedAddr, 0, len(pinned))
	for canon, amt := range pinned {
		res.Pinned = append(res.Pinned, PinnedAddr{Canonical: canon, Azeta: amt})
	}
	sort.Slice(res.Pinned, func(i, j int) bool {
		if !res.Pinned[i].Azeta.Equal(res.Pinned[j].Azeta) {
			return res.Pinned[i].Azeta.GT(res.Pinned[j].Azeta)
		}
		return res.Pinned[i].Canonical < res.Pinned[j].Canonical
	})

	res.Remainder = res.Supply.Sub(res.AttributedTotal())
	if res.Remainder.IsNegative() {
		return nil, fmt.Errorf(
			"remainder is negative: supply %s < attributed %s",
			res.Supply, res.AttributedTotal(),
		)
	}
	return res, nil
}
