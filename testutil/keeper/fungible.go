package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmdb "github.com/tendermint/tm-db"

	fungiblemocks "github.com/zeta-chain/zetacore/testutil/keeper/mocks/fungible"
	fungiblemodule "github.com/zeta-chain/zetacore/x/fungible"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

type FungibleMockOptions struct {
	UseBankMock     bool
	UseAccountMock  bool
	UseObserverMock bool
	UseEVMMock      bool
}

var (
	FungibleMocksAll = FungibleMockOptions{
		UseBankMock:     true,
		UseAccountMock:  true,
		UseObserverMock: true,
		UseEVMMock:      true,
	}
	FungibleNoMocks = FungibleMockOptions{}
)

// FungibleKeeperWithMocks initializes a fungible keeper for testing purposes with option to mock specific keepers
func FungibleKeeperWithMocks(t testing.TB, mockOptions FungibleMockOptions) (*keeper.Keeper, sdk.Context, SDKKeepers) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create observer keeper
	var observerKeeper types.ObserverKeeper = initObserverKeeper(
		cdc,
		db,
		stateStore,
		sdkKeepers.StakingKeeper,
		sdkKeepers.ParamsKeeper,
	)

	// Create the fungible keeper
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := NewContext(stateStore)

	// Initialize modules genesis
	sdkKeepers.InitGenesis(ctx)

	// Add a proposer to the context
	ctx = sdkKeepers.InitBlockProposer(t, ctx)

	// Initialize mocks for mocked keepers
	var authKeeper types.AccountKeeper = sdkKeepers.AuthKeeper
	var bankKeeper types.BankKeeper = sdkKeepers.BankKeeper
	var evmKeeper types.EVMKeeper = sdkKeepers.EvmKeeper
	if mockOptions.UseAccountMock {
		authKeeper = fungiblemocks.NewFungibleAccountKeeper(t)
	}
	if mockOptions.UseBankMock {
		bankKeeper = fungiblemocks.NewFungibleBankKeeper(t)
	}
	if mockOptions.UseObserverMock {
		observerKeeper = fungiblemocks.NewFungibleObserverKeeper(t)
	}
	if mockOptions.UseEVMMock {
		evmKeeper = fungiblemocks.NewFungibleEVMKeeper(t)
	}

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		sdkKeepers.ParamsKeeper.Subspace(types.ModuleName),
		authKeeper,
		evmKeeper,
		bankKeeper,
		observerKeeper,
	)

	fungiblemodule.InitGenesis(ctx, *k, *types.DefaultGenesis())

	return k, ctx, sdkKeepers
}

// FungibleKeeperAllMocks initializes a fungible keeper for testing purposes with all keeper mocked
func FungibleKeeperAllMocks(t testing.TB) (*keeper.Keeper, sdk.Context) {
	k, ctx, _ := FungibleKeeperWithMocks(t, FungibleMocksAll)
	return k, ctx
}

// FungibleKeeper initializes a fungible keeper for testing purposes
func FungibleKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, SDKKeepers) {
	k, ctx, sdkk := FungibleKeeperWithMocks(t, FungibleNoMocks)
	return k, ctx, sdkk
}

func GetFungibleAccountMock(t testing.TB, keeper *keeper.Keeper) *fungiblemocks.FungibleAccountKeeper {
	fak, ok := keeper.GetAuthKeeper().(*fungiblemocks.FungibleAccountKeeper)
	assert.True(t, ok)
	return fak
}

func GetFungibleBankMock(t testing.TB, keeper *keeper.Keeper) *fungiblemocks.FungibleBankKeeper {
	fbk, ok := keeper.GetBankKeeper().(*fungiblemocks.FungibleBankKeeper)
	assert.True(t, ok)
	return fbk
}

func GetFungibleObserverMock(t testing.TB, keeper *keeper.Keeper) *fungiblemocks.FungibleObserverKeeper {
	fok, ok := keeper.GetObserverKeeper().(*fungiblemocks.FungibleObserverKeeper)
	assert.True(t, ok)
	return fok
}

func GetFungibleEVMMock(t testing.TB, keeper *keeper.Keeper) *fungiblemocks.FungibleEVMKeeper {
	fek, ok := keeper.GetEVMKeeper().(*fungiblemocks.FungibleEVMKeeper)
	assert.True(t, ok)
	return fek
}
