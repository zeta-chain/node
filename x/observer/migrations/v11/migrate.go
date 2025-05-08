package v11

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

type observerKeeper interface {
	GetBallotListForHeight(ctx sdk.Context, height int64) (val types.BallotListForHeight, found bool)
	DeleteBallot(ctx sdk.Context, index string)
	DeleteBallotListForHeight(ctx sdk.Context, height int64)
}

const MaturityBlocks = int64(100)

// MigrateStore migrates the x/observer module state from consensus version 10 to version 11.
// The migration deletes all ballots and ballot lists for height 0.
// If ballots are not present for this height, it does nothing.
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	currentHeight := ctx.BlockHeight()
	// Maturity blocks is a parameter in the emissions module
	if currentHeight < MaturityBlocks {
		return nil
	}

	deleteHeight := int64(0)
	ballotList, found := observerKeeper.GetBallotListForHeight(ctx, deleteHeight)
	if !found {
		return nil
	}
	for _, ballotIndex := range ballotList.BallotsIndexList {
		observerKeeper.DeleteBallot(ctx, ballotIndex)
	}
	observerKeeper.DeleteBallotListForHeight(ctx, deleteHeight)

	return nil
}
