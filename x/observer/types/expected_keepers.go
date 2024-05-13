package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
)

type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool)
	SetValidator(ctx sdk.Context, validator stakingtypes.Validator)
}

type SlashingKeeper interface {
	IsTombstoned(ctx sdk.Context, addr sdk.ConsAddress) bool
	SetValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo)
}

type StakingHooks interface {
	AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error
	AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error
	BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error
}

type AuthorityKeeper interface {
	CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error

	// SetPolicies is solely used for the migration of policies from observer to authority
	SetPolicies(ctx sdk.Context, policies authoritytypes.Policies)
}

type LightclientKeeper interface {
	CheckNewBlockHeader(
		ctx sdk.Context,
		chainID int64,
		blockHash []byte,
		height int64,
		header proofs.HeaderData,
	) ([]byte, error)
	AddBlockHeader(
		ctx sdk.Context,
		chainID int64,
		height int64,
		blockHash []byte,
		header proofs.HeaderData,
		parentHash []byte,
	)
}
