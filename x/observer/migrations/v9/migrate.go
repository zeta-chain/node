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

// MigrateStore migrates the x/observer module state from the consensus version 7 to 8
// It updates the indexing for chain nonces object to use chain ID instead of chain name
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	currentHeight := ctx.BlockHeight()
	// Maturity blocks is a parameter in the emissions module
	maturityBlocks := int64(100)
	if currentHeight < maturityBlocks {
		return nil
	}
	maturedHeight := currentHeight - maturityBlocks

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
