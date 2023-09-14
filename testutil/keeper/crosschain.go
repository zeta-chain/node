package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmdb "github.com/tendermint/tm-db"

	crosschainmocks "github.com/zeta-chain/zetacore/testutil/keeper/mocks/crosschain"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
)

type CrosschainMockOptions struct {
	UseBankMock     bool
	UseAccountMock  bool
	UseStakingMock  bool
	UseObserverMock bool
	UseFungibleMock bool
}

var (
	CrosschainMocksAll = CrosschainMockOptions{
		UseBankMock:     true,
		UseAccountMock:  true,
		UseStakingMock:  true,
		UseObserverMock: true,
		UseFungibleMock: true,
	}
	CrosschainNoMocks = CrosschainMockOptions{}
)

// CrosschainKeeper initializes a crosschain keeper for testing purposes with option to mock specific keepers
func CrosschainKeeperWithMocks(t testing.TB, mockOptions CrosschainMockOptions) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create observer keeper
	var observerKeeper types.ZetaObserverKeeper = initObserverKeeper(
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
	var stakingKeeper types.StakingKeeper = sdkKeepers.StakingKeeper
	var fungibleKeeper types.FungibleKeeper = &fungiblekeeper.Keeper{}
	if mockOptions.UseAccountMock {
		authKeeper = crosschainmocks.NewCrosschainAccountKeeper(t)
	}
	if mockOptions.UseBankMock {
		bankKeeper = crosschainmocks.NewCrosschainBankKeeper(t)
	}
	if mockOptions.UseStakingMock {
		stakingKeeper = crosschainmocks.NewCrosschainStakingKeeper(t)
	}
	if mockOptions.UseObserverMock {
		observerKeeper = crosschainmocks.NewCrosschainObserverKeeper(t)
	}
	if mockOptions.UseFungibleMock {
		fungibleKeeper = crosschainmocks.NewCrosschainFungibleKeeper(t)
	}

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		stakingKeeper,
		sdkKeepers.ParamsKeeper.Subspace(types.ModuleName),
		authKeeper,
		bankKeeper,
		observerKeeper,
		fungibleKeeper,
	)

	return k, ctx
}

// CrosschainKeeperAllMocks initializes a crosschain keeper for testing purposes with all mocks
func CrosschainKeeperAllMocks(t testing.TB) (*keeper.Keeper, sdk.Context) {
	return CrosschainKeeperWithMocks(t, CrosschainMocksAll)
}

// CrosschainKeeper initializes a crosschain keeper for testing purposes
func CrosschainKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	return CrosschainKeeperWithMocks(t, CrosschainNoMocks)
}

func GetCrosschainAccountMock(t testing.TB, keeper *keeper.Keeper) *crosschainmocks.CrosschainAccountKeeper {
	cak, ok := keeper.GetAuthKeeper().(*crosschainmocks.CrosschainAccountKeeper)
	require.True(t, ok)
	return cak
}

func GetCrosschainBankMock(t testing.TB, keeper *keeper.Keeper) *crosschainmocks.CrosschainBankKeeper {
	cbk, ok := keeper.GetBankKeeper().(*crosschainmocks.CrosschainBankKeeper)
	require.True(t, ok)
	return cbk
}

func GetCrosschainStakingMock(t testing.TB, keeper *keeper.Keeper) *crosschainmocks.CrosschainStakingKeeper {
	csk, ok := keeper.GetStakingKeeper().(*crosschainmocks.CrosschainStakingKeeper)
	require.True(t, ok)
	return csk
}

func GetCrosschainObserverMock(t testing.TB, keeper *keeper.Keeper) *crosschainmocks.CrosschainObserverKeeper {
	cok, ok := keeper.GetObserverKeeper().(*crosschainmocks.CrosschainObserverKeeper)
	require.True(t, ok)
	return cok
}

func GetCrosschainFungibleMock(t testing.TB, keeper *keeper.Keeper) *crosschainmocks.CrosschainFungibleKeeper {
	cfk, ok := keeper.GetFungibleKeeper().(*crosschainmocks.CrosschainFungibleKeeper)
	require.True(t, ok)
	return cfk
}
