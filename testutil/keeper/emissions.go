package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	tmdb "github.com/tendermint/tm-db"
	emissionsmocks "github.com/zeta-chain/zetacore/testutil/keeper/mocks/emissions"
	"github.com/zeta-chain/zetacore/x/emissions/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

type EmissionMockOptions struct {
	UseBankMock       bool
	UseStakingMock    bool
	UseObserverMock   bool
	UseAccountMock    bool
	UseParamStoreMock bool
}

func EmissionsKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	return EmissionKeeperWithMockOptions(t, EmissionMockOptions{})
}
func EmissionKeeperWithMockOptions(
	t testing.TB,
	mockOptions EmissionMockOptions,
) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	SetConfig(false)
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create zeta keepers
	observerKeeperTmp := initObserverKeeper(
		cdc,
		db,
		stateStore,
		sdkKeepers.StakingKeeper,
		sdkKeepers.SlashingKeeper,
		sdkKeepers.ParamsKeeper,
		initAuthorityKeeper(cdc, db, stateStore),
	)

	zetaKeepers := ZetaKeepers{
		ObserverKeeper: observerKeeperTmp,
	}
	var observerKeeper types.ObserverKeeper = observerKeeperTmp

	// Create the fungible keeper
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
	var authKeeper types.AccountKeeper = sdkKeepers.AuthKeeper
	var bankKeeper types.BankKeeper = sdkKeepers.BankKeeper
	var stakingKeeper types.StakingKeeper = sdkKeepers.StakingKeeper
	if mockOptions.UseAccountMock {
		authKeeper = emissionsmocks.NewEmissionAccountKeeper(t)
	}
	if mockOptions.UseBankMock {
		bankKeeper = emissionsmocks.NewEmissionBankKeeper(t)
	}
	if mockOptions.UseStakingMock {
		stakingKeeper = emissionsmocks.NewEmissionStakingKeeper(t)
	}
	if mockOptions.UseObserverMock {
		observerKeeper = emissionsmocks.NewEmissionObserverKeeper(t)
	}

	var paramStore types.ParamStore
	if mockOptions.UseParamStoreMock {
		mock := emissionsmocks.NewEmissionParamStore(t)
		// mock this method for the keeper constructor
		mock.On("HasKeyTable").Maybe().Return(true)
		paramStore = mock
	} else {
		paramStore = sdkKeepers.ParamsKeeper.Subspace(types.ModuleName)
	}

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramStore,
		authtypes.FeeCollectorName,
		bankKeeper,
		stakingKeeper,
		observerKeeper,
		authKeeper,
	)

	if !mockOptions.UseParamStoreMock {
		k.SetParams(ctx, types.DefaultParams())
	}

	return k, ctx, sdkKeepers, zetaKeepers
}

func GetEmissionsBankMock(t testing.TB, keeper *keeper.Keeper) *emissionsmocks.EmissionBankKeeper {
	cbk, ok := keeper.GetBankKeeper().(*emissionsmocks.EmissionBankKeeper)
	require.True(t, ok)
	return cbk
}

func GetEmissionsParamStoreMock(t testing.TB, keeper *keeper.Keeper) *emissionsmocks.EmissionParamStore {
	m, ok := keeper.GetParamStore().(*emissionsmocks.EmissionParamStore)
	require.True(t, ok)
	return m
}
