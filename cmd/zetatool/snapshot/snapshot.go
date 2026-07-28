// Package snapshot holds the pure, deterministic compute for the ZETA->Solana
// migration snapshot (Stage 2). It turns parsed bank/staking genesis state into
// a per-address native-ZETA balance list for the Solana migration mint.
//
// Everything here is free of file, RPC and cobra dependencies so it can be
// unit-tested with hand-built fixtures. Reward math is intentionally absent:
// the Stage 1 export runs with --for-zero-height, so withdrawn staking rewards
// and validator commission are already folded into liquid bank balances.
package snapshot

import (
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

const (
	// BaseDenom is the native ZETA denom (18 decimals). Only this denom is summed.
	BaseDenom = "azeta"

	// classEOA / classValidator / classRemainder label each output row.
	classEOA       = "eoa"
	classValidator = "validator"
	classRemainder = "remainder"

	// claimStatusClaimable is the initial status for an attributed address; the
	// remainder is routed to the migration multisig instead of being claimable.
	claimStatusClaimable = "claimable"
	claimStatusMultisig  = "multisig"

	// zeroAddress is the null 20-byte account. It has no keypair to claim with,
	// so it is non-claimable and its balance folds into the remainder.
	zeroAddress = "0x0000000000000000000000000000000000000000"
)

// moduleAccountNames are the module accounts that hold azeta but must never be
// attributed to a claimant. Their balances (including the staking pools that
// physically hold delegated azeta) fall through to the remainder, which cancels
// against the staked/unbonding amounts credited to delegators.
var moduleAccountNames = []string{
	authtypes.FeeCollectorName,
	distrtypes.ModuleName,
	stakingtypes.BondedPoolName,
	stakingtypes.NotBondedPoolName,
	govtypes.ModuleName,
	crosschaintypes.ModuleName,
	evmtypes.ModuleName,
	fungibletypes.ModuleName,
	emissionstypes.ModuleName,
	emissionstypes.UndistributedObserverRewardsPool,
	emissionstypes.UndistributedTSSRewardsPool,
	feemarkettypes.ModuleName,
}

// oneGZeta is 1e9 azeta, the divisor to fold 18-decimal azeta down to the
// 9-decimal SPL representation used by the Solana mint.
var oneGZeta = sdkmath.NewInt(1_000_000_000)

// Config parameterizes a snapshot compute run.
type Config struct {
	// Denom is the native denom to snapshot. Callers pass BaseDenom.
	Denom string
	// Pins are extra non-claimable addresses (e.g. WZETA) that route to the
	// remainder instead of being attributed. Accepts zeta1.../zetavaloper1.../0x forms.
	Pins []string
}

// Account is a single attributed (claimable) address.
type Account struct {
	Canonical string
	Class     string
	Liquid    sdkmath.Int
	Staked    sdkmath.Int
	Unbonding sdkmath.Int
}

// Total returns liquid + staked + unbonding in azeta.
func (a Account) Total() sdkmath.Int {
	return a.Liquid.Add(a.Staked).Add(a.Unbonding)
}

// Total9Dec folds the account total from 18-decimal azeta down to 9 decimals.
func (a Account) Total9Dec() sdkmath.Int {
	return a.Total().Quo(oneGZeta)
}

// PinnedAddr is a pinned non-claimable address and the azeta it contributed to
// the remainder. Reported so an operator can confirm each pin actually matched
// a holder in the export; a zero Azeta means the pin matched nothing.
type PinnedAddr struct {
	Canonical string
	Azeta     sdkmath.Int
}

// Result is the full snapshot output: attributed accounts plus the remainder.
type Result struct {
	Accounts  []Account
	Remainder sdkmath.Int
	Supply    sdkmath.Int

	// Pinned is one entry per pinned address, amount descending, holding the
	// azeta that pin kept out of attribution.
	Pinned []PinnedAddr

	TotalLiquid    sdkmath.Int
	TotalStaked    sdkmath.Int
	TotalUnbonding sdkmath.Int

	// NonStandard counts well-formed but non-20-byte holders (e.g. module-derived
	// / group / interchain accounts) that have no eth-style keypair to claim with.
	// Their azeta is not attributed and is swept into the remainder.
	NonStandard      int
	NonStandardAzeta sdkmath.Int
}

// AttributedTotal is the sum over attributed accounts (liquid+staked+unbonding).
func (r *Result) AttributedTotal() sdkmath.Int {
	return r.TotalLiquid.Add(r.TotalStaked).Add(r.TotalUnbonding)
}

// SPLCap is the total 9-decimal supply cap: the sum of floored per-row totals
// across every claimable account plus the remainder row.
func (r *Result) SPLCap() sdkmath.Int {
	total := r.Remainder.Quo(oneGZeta)
	for _, a := range r.Accounts {
		total = total.Add(a.Total9Dec())
	}
	return total
}
