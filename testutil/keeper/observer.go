package keeper

import (
	"testing"

	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/stretchr/testify/assert"
	tmdb "github.com/tendermint/tm-db"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func initObserverKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	stakingKeeper stakingkeeper.Keeper,
	slashingKeeper slashingkeeper.Keeper,
	paramKeeper paramskeeper.Keeper,
) *keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, db)

	return keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		paramKeeper.Subspace(types.ModuleName),
		stakingKeeper,
		slashingKeeper,
	)
}

// ObserverKeeper instantiates an observer keeper for testing purposes
func ObserverKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create the observer keeper
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	assert.NoError(t, stateStore.LoadLatestVersion())

	ctx := NewContext(stateStore)

	// Initialize modules genesis
	sdkKeepers.InitGenesis(ctx)

	// Add a proposer to the context
	ctx = sdkKeepers.InitBlockProposer(t, ctx)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		sdkKeepers.ParamsKeeper.Subspace(types.ModuleName),
		sdkKeepers.StakingKeeper,
		sdkKeepers.SlashingKeeper,
	)

	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
