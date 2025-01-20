package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	// Methods imported from account should be defined here
}

type ObserverKeeper interface {
	GetBallot(ctx sdk.Context, index string) (val observertypes.Ballot, found bool)
	GetMaturedBallots(ctx sdk.Context, maturityBlocks int64) (val observertypes.BallotListForHeight, found bool)
	ClearMaturedBallotsAndBallotList(ctx sdk.Context, maturityBlocks int64)
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
	// Methods imported from bank should be defined here
}

type StakingKeeper interface {
	BondedRatio(ctx context.Context) (sdkmath.LegacyDec, error)
}
