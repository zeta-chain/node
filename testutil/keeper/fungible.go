package keeper

import (
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	fungiblemocks "github.com/zeta-chain/node/testutil/keeper/mocks/fungible"
	"github.com/zeta-chain/node/testutil/sample"
	authoritykeeper "github.com/zeta-chain/node/x/authority/keeper"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	fungiblemodule "github.com/zeta-chain/node/x/fungible"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
	lightclientkeeper "github.com/zeta-chain/node/x/lightclient/keeper"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observerkeeper "github.com/zeta-chain/node/x/observer/keeper"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

type FungibleMockOptions struct {
	UseBankMock      bool
	UseAccountMock   bool
	UseObserverMock  bool
	UseEVMMock       bool
	UseAuthorityMock bool
}

var (
	FungibleMocksAll = FungibleMockOptions{
		UseBankMock:      true,
		UseAccountMock:   true,
		UseObserverMock:  true,
		UseEVMMock:       true,
		UseAuthorityMock: true,
	}
	FungibleNoMocks = FungibleMockOptions{}
)

func initFungibleKeeper(
	cdc codec.Codec,
	ss store.CommitMultiStore,
	authKeeper types.AccountKeeper,
	bankKeepr types.BankKeeper,
	evmKeeper types.EVMKeeper,
	observerKeeper types.ObserverKeeper,
	authorityKeeper types.AuthorityKeeper,
) *keeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, nil)

	return keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		authKeeper,
		evmKeeper,
		bankKeepr,
		observerKeeper,
		authorityKeeper,
	)
}

// FungibleKeeperWithMocks initializes a fungible keeper for testing purposes with option to mock specific keepers
func FungibleKeeperWithMocks(
	t testing.TB,
	mockOptions FungibleMockOptions,
) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	keys, memKeys, tkeys, _ := StoreKeys()

	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepersWithKeys(cdc, keys, memKeys, tkeys)

	// Create authority keeper
	authorityKeeperTmp := authoritykeeper.NewKeeper(
		cdc,
		keys[authoritytypes.StoreKey],
		memKeys[authoritytypes.MemStoreKey],
		AuthorityGovAddress,
	)

	// Create lightclient keeper
	lightclientKeeperTmp := lightclientkeeper.NewKeeper(
		cdc,
		keys[lightclienttypes.StoreKey],
		memKeys[lightclienttypes.MemStoreKey],
		authorityKeeperTmp,
	)

	// Create observer keeper
	observerKeeperTmp := observerkeeper.NewKeeper(
		cdc,
		keys[observertypes.StoreKey],
		memKeys[observertypes.MemStoreKey],
		sdkKeepers.StakingKeeper,
		sdkKeepers.SlashingKeeper,
		authorityKeeperTmp,
		lightclientKeeperTmp,
		sdkKeepers.BankKeeper,
		sdkKeepers.AuthKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	zetaKeepers := ZetaKeepers{
		ObserverKeeper:    observerKeeperTmp,
		AuthorityKeeper:   &authorityKeeperTmp,
		LightclientKeeper: &lightclientKeeperTmp,
	}
	var observerKeeper types.ObserverKeeper = observerKeeperTmp
	var authorityKeeper types.AuthorityKeeper = authorityKeeperTmp

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := rootmulti.NewStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	// Create the fungible keeper
	for _, key := range keys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	}
	for _, key := range tkeys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeTransient, nil)
	}
	for _, key := range memKeys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeMemory, nil)
	}

	require.NoError(t, stateStore.LoadLatestVersion())

	// Initialize modules genesis
	ctx := NewContext(stateStore)
	sdkKeepers.InitGenesis(ctx)
	zetaKeepers.InitGenesis(ctx)

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
	if mockOptions.UseAuthorityMock {
		authorityKeeper = fungiblemocks.NewFungibleAuthorityKeeper(t)
	}

	k := keeper.NewKeeper(
		cdc,
		keys[types.StoreKey],
		memKeys[types.MemStoreKey],
		authKeeper,
		evmKeeper,
		bankKeeper,
		observerKeeper,
		authorityKeeper,
	)

	fungiblemodule.InitGenesis(ctx, *k, *types.DefaultGenesis())

	return k, ctx, sdkKeepers, zetaKeepers
}

// FungibleKeeperAllMocks initializes a fungible keeper for testing purposes with all keeper mocked
func FungibleKeeperAllMocks(t testing.TB) (*keeper.Keeper, sdk.Context) {
	k, ctx, _, _ := FungibleKeeperWithMocks(t, FungibleMocksAll)
	return k, ctx
}

// FungibleKeeper initializes a fungible keeper for testing purposes
func FungibleKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	k, ctx, sdkk, zk := FungibleKeeperWithMocks(t, FungibleNoMocks)
	return k, ctx, sdkk, zk
}

// GetFungibleAuthorityMock returns a new fungible authority keeper mock
func GetFungibleAuthorityMock(t testing.TB, keeper *keeper.Keeper) *fungiblemocks.FungibleAuthorityKeeper {
	cok, ok := keeper.GetAuthorityKeeper().(*fungiblemocks.FungibleAuthorityKeeper)
	require.True(t, ok)
	return cok
}

func GetFungibleAccountMock(t testing.TB, keeper *keeper.Keeper) *fungiblemocks.FungibleAccountKeeper {
	fak, ok := keeper.GetAuthKeeper().(*fungiblemocks.FungibleAccountKeeper)
	require.True(t, ok)
	return fak
}

func GetFungibleBankMock(t testing.TB, keeper *keeper.Keeper) *fungiblemocks.FungibleBankKeeper {
	fbk, ok := keeper.GetBankKeeper().(*fungiblemocks.FungibleBankKeeper)
	require.True(t, ok)
	return fbk
}

func GetFungibleObserverMock(t testing.TB, keeper *keeper.Keeper) *fungiblemocks.FungibleObserverKeeper {
	fok, ok := keeper.GetObserverKeeper().(*fungiblemocks.FungibleObserverKeeper)
	require.True(t, ok)
	return fok
}

func GetFungibleEVMMock(t testing.TB, keeper *keeper.Keeper) *FungibleMockEVMKeeper {
	fek, ok := keeper.GetEVMKeeper().(*fungiblemocks.FungibleEVMKeeper)
	require.True(t, ok)
	return &FungibleMockEVMKeeper{
		FungibleEVMKeeper: fek,
	}
}

type FungibleMockEVMKeeper struct {
	*fungiblemocks.FungibleEVMKeeper
}

func (m *FungibleMockEVMKeeper) SetupMockEVMKeeperForSystemContractDeployment() {
	gasRes := &evmtypes.EstimateGasResponse{Gas: 1000}
	m.On("WithChainID", mock.Anything).Maybe().Return(mock.Anything)
	m.On("ChainID").Maybe().Return(big.NewInt(1))
	m.On(
		"EstimateGas",
		mock.Anything,
		mock.Anything,
	).Return(gasRes, nil)
	m.MockEVMSuccessCallTimes(7)
	m.On(
		"GetAccount",
		mock.Anything,
		mock.Anything,
	).Return(&statedb.Account{
		Nonce: 1,
	})
	m.On(
		"GetCode",
		mock.Anything,
		mock.Anything,
	).Return([]byte{1, 2, 3})
}

func (m *FungibleMockEVMKeeper) MockEVMSuccessCallOnce() {
	m.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{})
}

func (m *FungibleMockEVMKeeper) MockEVMSuccessCallTimes(times int) {
	m.MockEVMSuccessCallTimesWithReturn(&evmtypes.MsgEthereumTxResponse{}, times)
}

func (m *FungibleMockEVMKeeper) MockEVMSuccessCallOnceWithReturn(ret *evmtypes.MsgEthereumTxResponse) {
	m.MockEVMSuccessCallTimesWithReturn(ret, 1)
}

func (m *FungibleMockEVMKeeper) MockEVMSuccessCallTimesWithReturn(ret *evmtypes.MsgEthereumTxResponse, times int) {
	if ret == nil {
		ret = &evmtypes.MsgEthereumTxResponse{}
	}
	m.On(
		"ApplyMessage",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(ret, nil).Times(times)
}

func (m *FungibleMockEVMKeeper) MockEVMFailCallOnce() {
	m.On(
		"ApplyMessage",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&evmtypes.MsgEthereumTxResponse{}, sample.ErrSample).Once()
}
