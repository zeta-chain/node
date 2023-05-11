package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/zeta-chain/zetacore/common"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type StakingKeeper interface {
	GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) types.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, types.ModuleAccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

type ZetaObserverKeeper interface {
	SetObserverMapper(ctx sdk.Context, om *zetaObserverTypes.ObserverMapper)
	GetObserverMapper(ctx sdk.Context, chain *common.Chain) (val zetaObserverTypes.ObserverMapper, found bool)
	GetAllObserverMappers(ctx sdk.Context) (mappers []*zetaObserverTypes.ObserverMapper)
	SetBallot(ctx sdk.Context, ballot *zetaObserverTypes.Ballot)
	GetBallot(ctx sdk.Context, index string) (val zetaObserverTypes.Ballot, found bool)
	GetAllBallots(ctx sdk.Context) (voters []*zetaObserverTypes.Ballot)
	GetParams(ctx sdk.Context) (params zetaObserverTypes.Params)
	GetCoreParamsByChainID(ctx sdk.Context, chainID int64) (params *zetaObserverTypes.CoreParams, found bool)
}
