package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	observermocks "github.com/zeta-chain/node/testutil/keeper/mocks/observer"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// ObserverMockOptions represents options for instantiating an observer keeper with mocks
type ObserverMockOptions struct {
	UseStakingMock     bool
	UseSlashingMock    bool
	UseAuthorityMock   bool
	UseLightclientMock bool
}

var (
	ObserverMocksAll = ObserverMockOptions{
		UseStakingMock:     true,
		UseSlashingMock:    true,
		UseAuthorityMock:   true,
		UseLightclientMock: true,
	}
	ObserverNoMocks = ObserverMockOptions{}
)

func initObserverKeeper(
	cdc codec.Codec,
	ss store.CommitMultiStore,
	stakingKeeper stakingkeeper.Keeper,
	slashingKeeper slashingkeeper.Keeper,
	authorityKeeper types.AuthorityKeeper,
	lightclientKeeper types.LightclientKeeper,
) *keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, nil)

	return keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		stakingKeeper,
		slashingKeeper,
		authorityKeeper,
		lightclientKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

// ObserverKeeperWithMocks instantiates an observer keeper for testing purposes with the option to mock specific keepers
func ObserverKeeperWithMocks(
	t testing.TB,
	mockOptions ObserverMockOptions,
) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := rootmulti.NewStore(db, log.NewNopLogger())
	cdc := NewCodec()

	authorityKeeperTmp := initAuthorityKeeper(cdc, stateStore)
	lightclientKeeperTmp := initLightclientKeeper(cdc, stateStore, authorityKeeperTmp)

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create the observer keeper
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
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
	var authorityKeeper types.AuthorityKeeper = authorityKeeperTmp
	var lightclientKeeper types.LightclientKeeper = lightclientKeeperTmp
	if mockOptions.UseStakingMock {
		stakingKeeper = observermocks.NewObserverStakingKeeper(t)
	}
	if mockOptions.UseSlashingMock {
		slashingKeeper = observermocks.NewObserverSlashingKeeper(t)
	}
	if mockOptions.UseAuthorityMock {
		authorityKeeper = observermocks.NewObserverAuthorityKeeper(t)
	}
	if mockOptions.UseLightclientMock {
		lightclientKeeper = observermocks.NewObserverLightclientKeeper(t)
	}

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		stakingKeeper,
		slashingKeeper,
		authorityKeeper,
		lightclientKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	return k, ctx, sdkKeepers, ZetaKeepers{
		AuthorityKeeper: &authorityKeeperTmp,
	}
}

// ObserverKeeper instantiates an observer keeper for testing purposes
func ObserverKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	return ObserverKeeperWithMocks(t, ObserverNoMocks)
}

// GetObserverLightclientMock returns a new observer lightclient keeper mock
func GetObserverLightclientMock(t testing.TB, keeper *keeper.Keeper) *observermocks.ObserverLightclientKeeper {
	cok, ok := keeper.GetLightclientKeeper().(*observermocks.ObserverLightclientKeeper)
	require.True(t, ok)
	return cok
}

// GetObserverAuthorityMock returns a new observer authority keeper mock
func GetObserverAuthorityMock(t testing.TB, keeper *keeper.Keeper) *observermocks.ObserverAuthorityKeeper {
	cok, ok := keeper.GetAuthorityKeeper().(*observermocks.ObserverAuthorityKeeper)
	require.True(t, ok)
	return cok
}

// GetObserverStakingMock returns a new observer staking keeper mock
func GetObserverStakingMock(t testing.TB, keeper *keeper.Keeper) *ObserverMockStakingKeeper {
	k, ok := keeper.GetStakingKeeper().(*observermocks.ObserverStakingKeeper)
	require.True(t, ok)
	return &ObserverMockStakingKeeper{
		ObserverStakingKeeper: k,
	}
}

// ObserverMockStakingKeeper is a wrapper of the observer staking keeper mock that add methods to mock the GetValidator method
type ObserverMockStakingKeeper struct {
	*observermocks.ObserverStakingKeeper
}

func (m *ObserverMockStakingKeeper) MockGetValidator(validator stakingtypes.Validator) {
	m.On("GetValidator", mock.Anything, mock.Anything).Return(validator, true)
}

// GetObserverSlashingMock returns a new observer slashing keeper mock
func GetObserverSlashingMock(t testing.TB, keeper *keeper.Keeper) *ObserverMockSlashingKeeper {
	k, ok := keeper.GetSlashingKeeper().(*observermocks.ObserverSlashingKeeper)
	require.True(t, ok)
	return &ObserverMockSlashingKeeper{
		ObserverSlashingKeeper: k,
	}
}

// ObserverMockSlashingKeeper is a wrapper of the observer slashing keeper mock that add methods to mock the IsTombstoned method
type ObserverMockSlashingKeeper struct {
	*observermocks.ObserverSlashingKeeper
}

func (m *ObserverMockSlashingKeeper) MockIsTombstoned(isTombstoned bool) {
	m.On("IsTombstoned", mock.Anything, mock.Anything).Return(isTombstoned)
}
