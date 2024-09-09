package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	"github.com/stretchr/testify/require"

	ibccrosschainmocks "github.com/zeta-chain/node/testutil/keeper/mocks/ibccrosschain"
	"github.com/zeta-chain/node/x/ibccrosschain/keeper"
	"github.com/zeta-chain/node/x/ibccrosschain/types"
)

type IBCCroscchainMockOptions struct {
	UseCrosschainMock  bool
	UseIBCTransferMock bool
}

var (
	IBCCrosschainMocksAll = IBCCroscchainMockOptions{
		UseCrosschainMock:  true,
		UseIBCTransferMock: true,
	}
	IBCCrosschainNoMocks = IBCCroscchainMockOptions{}
)

func initIBCCrosschainKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	crosschainKeeper types.CrosschainKeeper,
	ibcTransferKeeper types.IBCTransferKeeper,
	capabilityKeeper capabilitykeeper.Keeper,
) *keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, db)

	capabilityKeeper.ScopeToModule(types.ModuleName)

	return keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		crosschainKeeper,
		ibcTransferKeeper,
	)
}

func IBCCrosschainKeeperWithMocks(
	t testing.TB,
	mockOptions IBCCroscchainMockOptions,
) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create zeta keepers
	authorityKeeper := initAuthorityKeeper(cdc, stateStore)
	lightclientKeeper := initLightclientKeeper(cdc, stateStore, authorityKeeper)
	observerKeeper := initObserverKeeper(
		cdc,
		stateStore,
		sdkKeepers.StakingKeeper,
		sdkKeepers.SlashingKeeper,
		authorityKeeper,
		lightclientKeeper,
	)
	fungibleKeeper := initFungibleKeeper(
		cdc,
		stateStore,
		sdkKeepers.AuthKeeper,
		sdkKeepers.BankKeeper,
		sdkKeepers.EvmKeeper,
		observerKeeper,
		authorityKeeper,
	)
	crosschainKeeperTmp := initCrosschainKeeper(
		cdc,
		db,
		stateStore,
		sdkKeepers.StakingKeeper,
		sdkKeepers.AuthKeeper,
		sdkKeepers.BankKeeper,
		observerKeeper,
		fungibleKeeper,
		authorityKeeper,
		lightclientKeeper,
	)

	zetaKeepers := ZetaKeepers{
		ObserverKeeper:    observerKeeper,
		FungibleKeeper:    fungibleKeeper,
		AuthorityKeeper:   &authorityKeeper,
		LightclientKeeper: &lightclientKeeper,
		CrosschainKeeper:  crosschainKeeperTmp,
	}

	var crosschainKeeper types.CrosschainKeeper = crosschainKeeperTmp
	var ibcTransferKeeper types.IBCTransferKeeper = sdkKeepers.TransferKeeper

	// Create the ibccrosschain keeper
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := NewContext(stateStore)

	// Initialize modules genesis
	sdkKeepers.InitGenesis(ctx)
	zetaKeepers.InitGenesis(ctx)

	// Add a proposer to the context
	ctx = sdkKeepers.InitBlockProposer(t, ctx)

	// Initialize mocks for mocked keepers
	if mockOptions.UseCrosschainMock {
		crosschainKeeper = ibccrosschainmocks.NewLightclientCrosschainKeeper(t)
	}
	if mockOptions.UseIBCTransferMock {
		ibcTransferKeeper = ibccrosschainmocks.NewLightclientTransferKeeper(t)
	}

	sdkKeepers.CapabilityKeeper.ScopeToModule(types.ModuleName)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		crosschainKeeper,
		ibcTransferKeeper,
	)

	// seal the IBC router
	sdkKeepers.IBCKeeper.SetRouter(sdkKeepers.IBCRouter)

	return k, ctx, sdkKeepers, zetaKeepers
}

// IBCCrosschainKeeperAllMocks creates a new ibccrosschain keeper with all mocked keepers
func IBCCrosschainKeeperAllMocks(t testing.TB) (*keeper.Keeper, sdk.Context) {
	k, ctx, _, _ := IBCCrosschainKeeperWithMocks(t, IBCCrosschainMocksAll)
	return k, ctx
}

// IBCCrosschainKeeper creates a new ibccrosschain keeper with no mocked keepers
func IBCCrosschainKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	return IBCCrosschainKeeperWithMocks(t, IBCCrosschainNoMocks)
}
