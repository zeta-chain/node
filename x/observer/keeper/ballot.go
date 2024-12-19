package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

var GetMaturedBallotHeightFunc = getMaturedBallotHeight

func (k Keeper) SetBallot(ctx sdk.Context, ballot *types.Ballot) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VoterKey))
	ballot.Index = ballot.BallotIdentifier
	b := k.cdc.MustMarshal(ballot)
	store.Set([]byte(ballot.Index), b)
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
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
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

// ClearMaturedBallotsAndBallotList deletes all matured ballots and the list of ballots for a given height.
// It also emits an event for each ballot deleted.
func (k Keeper) ClearMaturedBallotsAndBallotList(ctx sdk.Context, maturityBlocksParam int64) {
	maturedBallotsHeight := getMaturedBallotHeight(ctx, maturityBlocksParam)

	// Fetch all the matured ballots, return if no matured ballots are found
	// For the current implementation, this should never happen as ClearMaturedBallotsAndBallotList is only called after the Distribution of the rewards,
	//	which means that there are matured ballots to be deleted
	maturedBallots, found := k.GetBallotListForHeight(ctx, maturedBallotsHeight)
	if !found {
		return
	}
	// Delete all the matured ballots and emit an event for each ballot deleted
	for _, ballotIndex := range maturedBallots.BallotsIndexList {
		ballot, found := k.GetBallot(ctx, ballotIndex)
		if !found {
			continue
		}
		k.DeleteBallot(ctx, ballotIndex)
		EmitEventBallotDeleted(ctx, ballot)
	}
	// Delete the list of matured ballots
	k.DeleteBallotListForHeight(ctx, maturedBallotsHeight)
	return
}

// getMaturedBallotHeight returns the height at which a ballot is considered matured.
func getMaturedBallotHeight(ctx sdk.Context, maturityBlocks int64) int64 {
	return ctx.BlockHeight() - maturityBlocks
}
