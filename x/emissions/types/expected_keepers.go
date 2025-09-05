package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

type ObserverKeeper interface {
	GetBallot(ctx sdk.Context, index string) (val observertypes.Ballot, found bool)
	GetMaturedBallots(ctx sdk.Context, maturityBlocks int64) (val observertypes.BallotListForHeight, found bool)
	ClearFinalizedMaturedBallots(ctx sdk.Context, maturityBlocks int64, deleteAllBallots bool)
	GetObserverSet(ctx sdk.Context) (val observertypes.ObserverSet, found bool)
	CheckObserverCanVote(ctx sdk.Context, address string) error
	GetSupportedChains(ctx sdk.Context) []chains.Chain
	GetNodeAccount(ctx sdk.Context, address string) (observertypes.NodeAccount, bool)
	GetAllNodeAccount(ctx sdk.Context) []observertypes.NodeAccount
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(
		ctx context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}

type StakingKeeper interface {
	BondedRatio(ctx context.Context) (sdkmath.LegacyDec, error)
}
