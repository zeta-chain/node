package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	emissionsmocks "github.com/zeta-chain/node/testutil/keeper/mocks/emissions"
	"github.com/zeta-chain/node/x/emissions/keeper"
	"github.com/zeta-chain/node/x/emissions/types"
)

type EmissionMockOptions struct {
	UseBankMock       bool
	UseStakingMock    bool
	UseObserverMock   bool
	UseAccountMock    bool
	SkipSettingParams bool
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
	stateStore := rootmulti.NewStore(db, log.NewNopLogger())
	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	authorityKeeper := initAuthorityKeeper(cdc, stateStore)

	// Create zeta keepers
	observerKeeperTmp := initObserverKeeper(
		cdc,
		stateStore,
		sdkKeepers.StakingKeeper,
		sdkKeepers.SlashingKeeper,
		authorityKeeper,
		initLightclientKeeper(cdc, stateStore, authorityKeeper),
	)

	zetaKeepers := ZetaKeepers{
		ObserverKeeper: observerKeeperTmp,
	}
	var observerKeeper types.ObserverKeeper = observerKeeperTmp

	// Create the fungible keeper
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
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

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		authtypes.FeeCollectorName,
		bankKeeper,
		stakingKeeper,
		observerKeeper,
		authKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	if !mockOptions.SkipSettingParams {
		err := k.SetParams(ctx, types.DefaultParams())
		require.NoError(t, err)
	}

	return k, ctx, sdkKeepers, zetaKeepers
}

func GetEmissionsBankMock(t testing.TB, keeper *keeper.Keeper) *emissionsmocks.EmissionBankKeeper {
	cbk, ok := keeper.GetBankKeeper().(*emissionsmocks.EmissionBankKeeper)
	require.True(t, ok)
	return cbk
}

func GetEmissionsStakingMock(t testing.TB, keeper *keeper.Keeper) *emissionsmocks.EmissionStakingKeeper {
	cbk, ok := keeper.GetStakingKeeper().(*emissionsmocks.EmissionStakingKeeper)
	require.True(t, ok)
	return cbk
}
