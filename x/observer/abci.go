package observer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func CleanupState(ctx sdk.Context, keeper keeper.Keeper) {
	pendingBallots := keeper.GetAllBallots(ctx)
	for _, ballot := range pendingBallots {
		if IsBallotExpired(ctx, ballot) {
			keeper.RemoveBallot(ctx, ballot.Index)
		}
	}
}

func IsBallotExpired(ctx sdk.Context, ballot *types.Ballot) bool {
	if ballot.CreationHeight+common.BlocksPerDay <= ctx.BlockHeight() {
		return true
	}
	return false
}
