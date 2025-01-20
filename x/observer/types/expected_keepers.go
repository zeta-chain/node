package types

import (
	context "context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/proofs"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
)

type StakingKeeper interface {
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	GetDelegation(
		ctx context.Context,
		delAddr sdk.AccAddress,
		valAddr sdk.ValAddress,
	) (delegation stakingtypes.Delegation, err error)
	SetValidator(ctx context.Context, validator stakingtypes.Validator) error
	GetAllValidators(ctx context.Context) (validators []stakingtypes.Validator, err error)
}

type SlashingKeeper interface {
	IsTombstoned(ctx context.Context, addr sdk.ConsAddress) bool
	SetValidatorSigningInfo(ctx context.Context, address sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) error
}

type StakingHooks interface {
	AfterValidatorRemoved(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error
	AfterValidatorBeginUnbonding(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error
	AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error
	BeforeValidatorSlashed(ctx context.Context, valAddr sdk.ValAddress, fraction sdkmath.LegacyDec) error
}

type AuthorityKeeper interface {
	CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error
	GetAdditionalChainList(ctx sdk.Context) (list []chains.Chain)

	// SetPolicies is solely used for the migration of policies from observer to authority
	SetPolicies(ctx sdk.Context, policies authoritytypes.Policies)
	GetPolicies(ctx sdk.Context) (val authoritytypes.Policies, found bool)
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

type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}
