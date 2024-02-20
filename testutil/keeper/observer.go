package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	tmdb "github.com/tendermint/tm-db"
	observermocks "github.com/zeta-chain/zetacore/testutil/keeper/mocks/observer"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// ObserverMockOptions represents options for instantiating an observer keeper with mocks
type ObserverMockOptions struct {
	UseStakingMock  bool
	UseSlashingMock bool
}

var (
	ObserverMocksAll = ObserverMockOptions{
		UseStakingMock:  true,
		UseSlashingMock: true,
	}
	ObserverNoMocks = ObserverMockOptions{}
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

// ObserverKeeperWithMocks instantiates an observer keeper for testing purposes with the option to mock specific keepers
func ObserverKeeperWithMocks(t testing.TB, mockOptions ObserverMockOptions) (*keeper.Keeper, sdk.Context, SDKKeepers) {
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
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := NewContext(stateStore)

	// Initialize modules genesis
	sdkKeepers.InitGenesis(ctx)

	// Add a proposer to the context
	ctx = sdkKeepers.InitBlockProposer(t, ctx)

	// Initialize mocks for mocked keepers
	var stakingKeeper types.StakingKeeper = sdkKeepers.StakingKeeper
	var slashingKeeper types.SlashingKeeper = sdkKeepers.SlashingKeeper
	if mockOptions.UseStakingMock {
		stakingKeeper = observermocks.NewObserverStakingKeeper(t)
	}
	if mockOptions.UseSlashingMock {
		slashingKeeper = observermocks.NewObserverSlashingKeeper(t)
	}

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		sdkKeepers.ParamsKeeper.Subspace(types.ModuleName),
		stakingKeeper,
		slashingKeeper,
	)

	k.SetParams(ctx, types.DefaultParams())

	return k, ctx, sdkKeepers
}

// ObserverKeeper instantiates an observer keeper for testing purposes
func ObserverKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, SDKKeepers) {
	return ObserverKeeperWithMocks(t, ObserverNoMocks)
}

// GetObserverStakingMock returns a new observer staking keeper mock
func GetObserverStakingMock(t testing.TB, keeper *keeper.Keeper) *ObserverMockStakingKeeper {
	k, ok := keeper.GetStakingKeeper().(*observermocks.ObserverStakingKeeper)
	require.True(t, ok)
	return &ObserverMockStakingKeeper{
		ObserverStakingKeeper: k,
	}
}

// GetObserverSlashingMock returns a new observer slashing keeper mock
func GetObserverSlashingMock(t testing.TB, keeper *keeper.Keeper) *ObserverMockSlashingKeeper {
	k, ok := keeper.GetSlashingKeeper().(*observermocks.ObserverSlashingKeeper)
	require.True(t, ok)
	return &ObserverMockSlashingKeeper{
		ObserverSlashingKeeper: k,
	}
}

// ObserverMockStakingKeeper is a wrapper of the observer staking keeper mock that add methods to mock the GetValidator method
type ObserverMockStakingKeeper struct {
	*observermocks.ObserverStakingKeeper
}

func (m *ObserverMockStakingKeeper) MockGetValidator(validator stakingtypes.Validator) {
	m.On("GetValidator", mock.Anything, mock.Anything).Return(validator, true)
}

// ObserverMockSlashingKeeper is a wrapper of the observer slashing keeper mock that add methods to mock the IsTombstoned method
type ObserverMockSlashingKeeper struct {
	*observermocks.ObserverSlashingKeeper
}

func (m *ObserverMockSlashingKeeper) MockIsTombstoned(isTombstoned bool) {
	m.On("IsTombstoned", mock.Anything, mock.Anything).Return(isTombstoned)
}
