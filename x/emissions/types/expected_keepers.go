package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
	// Methods imported from account should be defined here
}

type ObserverKeeper interface {
	GetMaturedBallotList(ctx sdk.Context) []string
	GetObserverSet(ctx sdk.Context) (val observertypes.ObserverSet, found bool)
	SetBallot(ctx sdk.Context, ballot *observertypes.Ballot)
	GetBallot(ctx sdk.Context, index string) (val observertypes.Ballot, found bool)
	GetAllBallots(ctx sdk.Context) (voters []*observertypes.Ballot)
	GetParams(ctx sdk.Context) (params observertypes.Params)
	GetCoreParamsByChainID(ctx sdk.Context, chainID int64) (params *observertypes.CoreParams, found bool)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	// Methods imported from bank should be defined here
}

type StakingKeeper interface {
	BondedRatio(ctx sdk.Context) sdk.Dec
}
