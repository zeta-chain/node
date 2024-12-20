package v9

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

// MigrateStore migrates the x/observer module state from the consensus version 8 to version 9.
// The migration deletes all the ballots and ballot lists that are older than MaturityBlocks.
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	currentHeight := ctx.BlockHeight()
	// Maturity blocks is a parameter in the emissions module
	if currentHeight < MaturityBlocks {
		return nil
	}
	maturedHeight := currentHeight - MaturityBlocks
	for i := maturedHeight; i > 0; i-- {
		ballotList, found := observerKeeper.GetBallotListForHeight(ctx, i)
		if !found {
			continue
		}
		for _, ballotIndex := range ballotList.BallotsIndexList {
			observerKeeper.DeleteBallot(ctx, ballotIndex)
		}
		observerKeeper.DeleteBallotListForHeight(ctx, i)
	}

	return nil
}
