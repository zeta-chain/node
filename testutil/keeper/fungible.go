package keeper

import (
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	tmdb "github.com/tendermint/tm-db"
	fungiblemocks "github.com/zeta-chain/zetacore/testutil/keeper/mocks/fungible"
	"github.com/zeta-chain/zetacore/testutil/sample"
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

func initFungibleKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	paramKeeper paramskeeper.Keeper,
	authKeeper types.AccountKeeper,
	bankKeepr types.BankKeeper,
	evmKeeper types.EVMKeeper,
	observerKeeper types.ObserverKeeper,
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
		authKeeper,
		evmKeeper,
		bankKeepr,
		observerKeeper,
	)
}

// FungibleKeeperWithMocks initializes a fungible keeper for testing purposes with option to mock specific keepers
func FungibleKeeperWithMocks(t testing.TB, mockOptions FungibleMockOptions) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create observer keeper
	observerKeeperTmp := initObserverKeeper(
		cdc,
		db,
		stateStore,
		sdkKeepers.StakingKeeper,
		sdkKeepers.SlashingKeeper,
		sdkKeepers.ParamsKeeper,
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
	m.MockEVMSuccessCallTimes(5)
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
	).Return(ret, nil).Times(times)
}

func (m *FungibleMockEVMKeeper) MockEVMFailCallOnce() {
	m.On(
		"ApplyMessage",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&evmtypes.MsgEthereumTxResponse{}, sample.ErrSample).Once()
}
