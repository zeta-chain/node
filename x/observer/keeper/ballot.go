package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

func (k Keeper) SetBallot(ctx sdk.Context, ballot *types.Ballot) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	b := k.cdc.MustMarshal(ballot)
	store.Set([]byte(ballot.BallotIdentifier), b)
}

func (k Keeper) SetBallotList(ctx sdk.Context, ballotlist *types.BallotListForHeight) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BallotListKey))
	b := k.cdc.MustMarshal(ballotlist)
	store.Set(types.BallotListKeyPrefix(ballotlist.Height), b)
}

func (k Keeper) DeleteBallot(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	store.Delete([]byte(index))
}

func (k Keeper) DeleteBallotListForHeight(ctx sdk.Context, height int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BallotListKey))
	store.Delete(types.BallotListKeyPrefix(height))
}

func (k Keeper) GetBallot(ctx sdk.Context, index string) (val types.Ballot, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetBallotListForHeight(ctx sdk.Context, height int64) (val types.BallotListForHeight, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BallotListKey))
	b := store.Get(types.BallotListKeyPrefix(height))
	if b == nil {
		return types.BallotListForHeight{
			Height:           height,
			BallotsIndexList: nil,
		}, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetMaturedBallots(ctx sdk.Context, maturityBlocks int64) (val types.BallotListForHeight, found bool) {
	return k.GetBallotListForHeight(ctx, getMaturedBallotHeight(ctx, maturityBlocks))
}

func (k Keeper) GetAllBallots(ctx sdk.Context) (voters []*types.Ballot) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Ballot
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		voters = append(voters, &val)
	}
	return
}

// AddBallotToList adds a ballot to the list of ballots for a given height.
func (k Keeper) AddBallotToList(ctx sdk.Context, ballot types.Ballot) {
	list, found := k.GetBallotListForHeight(ctx, ballot.BallotCreationHeight)
	if !found {
		list = types.BallotListForHeight{Height: ballot.BallotCreationHeight, BallotsIndexList: []string{}}
	}
	list.BallotsIndexList = append(list.BallotsIndexList, ballot.BallotIdentifier)
	k.SetBallotList(ctx, &list)
}

// ClearFinalizedMaturedBallots deletes all matured and finalized ballots for a given height.

// ClearFinalizedMaturedBallots deletes matured ballots and optionally the ballot list for a given height.
// If deleteAllBallots is true, it deletes all ballots and always deletes the ballot list.
// If deleteAllBallots is false, it only deletes finalized ballots and only deletes the ballot list if all ballots are deleted.
// It emits an event for each ballot deleted.
func (k Keeper) ClearFinalizedMaturedBallots(ctx sdk.Context, maturityBlocks int64, forceDeleteAllBallots bool) {
	maturedBallotsHeight := getMaturedBallotHeight(ctx, maturityBlocks)
	maturedBallots, found := k.GetBallotListForHeight(ctx, maturedBallotsHeight)
	if !found {
		return
	}

	ballotsDeleted := 0
	for _, ballotIndex := range maturedBallots.BallotsIndexList {
		ballot, foundBallot := k.GetBallot(ctx, ballotIndex)
		if !foundBallot {
			continue
		}

		// Skip non-finalized ballots if we're only deleting finalized ballots
		if !forceDeleteAllBallots && !ballot.IsFinalized() {
			continue
		}

		k.DeleteBallot(ctx, ballotIndex)
		logBallotDeletion(ctx, ballot)
		ballotsDeleted++
	}

	// Delete ballot list if:
	// 1. We're deleting all ballots (regardless of how many were actually deleted), OR
	// 2. We deleted only finalized ballots but all ballots in the list were deleted
	if forceDeleteAllBallots || ballotsDeleted == len(maturedBallots.BallotsIndexList) {
		k.DeleteBallotListForHeight(ctx, maturedBallotsHeight)
	}
}

// getMaturedBallotHeight returns the height at which a ballot is considered matured.
func getMaturedBallotHeight(ctx sdk.Context, maturityBlocks int64) int64 {
	return ctx.BlockHeight() - maturityBlocks
}

func logBallotDeletion(ctx sdk.Context, ballot types.Ballot) {
	if len(ballot.VoterList) != len(ballot.Votes) {
		ctx.Logger().
			Error(fmt.Sprintf("voter list and votes list length mismatch for deleted ballot %s", ballot.BallotIdentifier))
		return
	}

	votersList := ""
	for i := range ballot.VoterList {
		votersList += fmt.Sprintf("Voter : %s | Vote : %s\n", ballot.VoterList[i], ballot.Votes[i])
	}

	ctx.Logger().
		Debug(fmt.Sprintf("ballotIdentifier: %s,ballotType: %s,voterList: %s", ballot.BallotIdentifier, ballot.ObservationType.String(), votersList))
}
