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
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	crosschainmocks "github.com/zeta-chain/node/testutil/keeper/mocks/crosschain"
	"github.com/zeta-chain/node/testutil/sample"
	authoritykeeper "github.com/zeta-chain/node/x/authority/keeper"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	lightclientkeeper "github.com/zeta-chain/node/x/lightclient/keeper"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observerkeeper "github.com/zeta-chain/node/x/observer/keeper"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

type CrosschainMockOptions struct {
	UseBankMock          bool
	UseAccountMock       bool
	UseStakingMock       bool
	UseObserverMock      bool
	UseFungibleMock      bool
	UseAuthorityMock     bool
	UseLightclientMock   bool
	UseIBCCrosschainMock bool
}

var (
	CrosschainMocksAll = CrosschainMockOptions{
		UseBankMock:          true,
		UseAccountMock:       true,
		UseStakingMock:       true,
		UseObserverMock:      true,
		UseFungibleMock:      true,
		UseAuthorityMock:     true,
		UseLightclientMock:   true,
		UseIBCCrosschainMock: true,
	}
	CrosschainNoMocks = CrosschainMockOptions{}
)

func initCrosschainKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	stakingKeeper types.StakingKeeper,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	observerKeeper types.ObserverKeeper,
	fungibleKeeper types.FungibleKeeper,
	authorityKeeper types.AuthorityKeeper,
	lightclientKeeper types.LightclientKeeper,
) *keeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, db)

	return keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		stakingKeeper,
		authKeeper,
		bankKeeper,
		observerKeeper,
		fungibleKeeper,
		authorityKeeper,
		lightclientKeeper,
	)
}

// CrosschainKeeperWithMocks initializes a crosschain keeper for testing purposes with option to mock specific keepers
func CrosschainKeeperWithMocks(
	t testing.TB,
	mockOptions CrosschainMockOptions,
) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	keys, memKeys, tkeys, _ := StoreKeys()
	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepersWithKeys(cdc, keys, memKeys, tkeys)

	// Create zeta keepers
	authorityKeeperTmp := authoritykeeper.NewKeeper(
		cdc,
		keys[authoritytypes.StoreKey],
		memKeys[authoritytypes.MemStoreKey],
		AuthorityGovAddress,
	)

	lightclientKeeperTmp := lightclientkeeper.NewKeeper(
		cdc,
		keys[lightclienttypes.StoreKey],
		memKeys[lightclienttypes.MemStoreKey],
		authorityKeeperTmp,
	)

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

	fungibleKeeperTmp := fungiblekeeper.NewKeeper(
		cdc,
		keys[types.StoreKey],
		memKeys[types.MemStoreKey],
		sdkKeepers.AuthKeeper,
		sdkKeepers.EvmKeeper,
		sdkKeepers.BankKeeper,
		observerKeeperTmp,
		authorityKeeperTmp,
	)

	zetaKeepers := ZetaKeepers{
		ObserverKeeper:  observerKeeperTmp,
		FungibleKeeper:  fungibleKeeperTmp,
		AuthorityKeeper: &authorityKeeperTmp,
	}
	var lightclientKeeper types.LightclientKeeper = lightclientKeeperTmp
	var authorityKeeper types.AuthorityKeeper = authorityKeeperTmp
	var observerKeeper types.ObserverKeeper = observerKeeperTmp
	var fungibleKeeper types.FungibleKeeper = fungibleKeeperTmp

	// Initialize local store
	db := tmdb.NewMemDB()
	logger := log.NewNopLogger()
	stateStore := rootmulti.NewStore(db, logger, metrics.NewNoOpMetrics())
	for _, key := range keys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	}
	for _, key := range tkeys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeTransient, nil)
	}
	for _, key := range memKeys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeMemory, nil)
	}

	// Initialize mocks for mocked keepers
	var authKeeper types.AccountKeeper = sdkKeepers.AuthKeeper
	var bankKeeper types.BankKeeper = sdkKeepers.BankKeeper
	var stakingKeeper types.StakingKeeper = sdkKeepers.StakingKeeper
	if mockOptions.UseAccountMock {
		authKeeper = crosschainmocks.NewCrosschainAccountKeeper(t)
	}
	if mockOptions.UseBankMock {
		bankKeeper = crosschainmocks.NewCrosschainBankKeeper(t)
	}
	if mockOptions.UseStakingMock {
		stakingKeeper = crosschainmocks.NewCrosschainStakingKeeper(t)
	}

	if mockOptions.UseAuthorityMock {
		authorityKeeper = crosschainmocks.NewCrosschainAuthorityKeeper(t)
	}
	if mockOptions.UseObserverMock {
		observerKeeper = crosschainmocks.NewCrosschainObserverKeeper(t)
	}
	if mockOptions.UseFungibleMock {
		fungibleKeeper = crosschainmocks.NewCrosschainFungibleKeeper(t)
	}
	if mockOptions.UseLightclientMock {
		lightclientKeeper = crosschainmocks.NewCrosschainLightclientKeeper(t)
	}

	// create crosschain keeper
	k := keeper.NewKeeper(
		cdc,
		keys[types.StoreKey],
		memKeys[types.MemStoreKey],
		stakingKeeper,
		authKeeper,
		bankKeeper,
		observerKeeper,
		fungibleKeeper,
		authorityKeeper,
		lightclientKeeper,
	)

	// initialize ibccrosschain keeper and set it to the crosschain keeper
	// there is a circular dependency between the two keepers, crosschain keeper must be initialized first
	var ibcCrosschainKeeperTmp types.IBCCrosschainKeeper = initIBCCrosschainKeeper(
		cdc,
		db,
		stateStore,
		k,
		sdkKeepers.TransferKeeper,
		*sdkKeepers.CapabilityKeeper,
	)
	if mockOptions.UseIBCCrosschainMock {
		ibcCrosschainKeeperTmp = crosschainmocks.NewCrosschainIBCCrosschainKeeper(t)
	}
	k.SetIBCCrosschainKeeper(ibcCrosschainKeeperTmp)

	// seal the IBC router
	// sdkKeepers.IBCKeeper.SetRouter(sdkKeepers.IBCRouter)

	// load the latest version of the state store
	require.NoError(t, stateStore.LoadLatestVersion())
	ctx := NewContext(stateStore)

	// initialize modules genesis
	sdkKeepers.InitGenesis(ctx)
	zetaKeepers.InitGenesis(ctx)

	// add a proposer to the context
	ctx = sdkKeepers.InitBlockProposer(t, ctx)

	return k, ctx, sdkKeepers, zetaKeepers
}

// CrosschainKeeperAllMocks initializes a crosschain keeper for testing purposes with all mocks
func CrosschainKeeperAllMocks(t testing.TB) (*keeper.Keeper, sdk.Context) {
	k, ctx, _, _ := CrosschainKeeperWithMocks(t, CrosschainMocksAll)
	return k, ctx
}

// CrosschainKeeper initializes a crosschain keeper for testing purposes
func CrosschainKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	return CrosschainKeeperWithMocks(t, CrosschainNoMocks)
}

// GetCrosschainLightclientMock returns a new crosschain lightclient keeper mock
func GetCrosschainLightclientMock(t testing.TB, keeper *keeper.Keeper) *crosschainmocks.CrosschainLightclientKeeper {
	lk, ok := keeper.GetLightclientKeeper().(*crosschainmocks.CrosschainLightclientKeeper)
	require.True(t, ok)
	return lk
}

// GetCrosschainAuthorityMock returns a new crosschain authority keeper mock
func GetCrosschainAuthorityMock(t testing.TB, keeper *keeper.Keeper) *crosschainmocks.CrosschainAuthorityKeeper {
	cok, ok := keeper.GetAuthorityKeeper().(*crosschainmocks.CrosschainAuthorityKeeper)
	require.True(t, ok)
	return cok
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

func MockGetSupportedChainFromChainID(m *crosschainmocks.CrosschainObserverKeeper, senderChain chains.Chain) {
	m.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
		Return(senderChain, true).Once()
}

func MockFailedGetSupportedChainFromChainID(m *crosschainmocks.CrosschainObserverKeeper, senderChain chains.Chain) {
	m.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
		Return(chains.Chain{}, false).Once()
}

func MockProcessV2RevertDeposit(
	m *crosschainmocks.CrosschainFungibleKeeper,
	retEVMTxResponse *evmtypes.MsgEthereumTxResponse,
	retErr error,
) {
	// ctx types.Context
	// inboundSender string
	// amount *big.Int
	// chainID int64
	// coinType coin.FungibleTokenCoinType
	// asset string
	// revertAddress common.Address
	// callOnRevert bool
	// revertMessage []byte
	m.On(
		"ProcessRevert",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).
		Return(retEVMTxResponse, retErr).Once()
}

func MockGetRevertGasLimitForERC20(
	m *crosschainmocks.CrosschainFungibleKeeper,
	asset string,
	senderChain chains.Chain,
	returnVal int64,
) {
	m.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
		Return(fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().String(),
		}, true).Once()
	m.On("QueryGasLimit", mock.Anything, mock.Anything).
		Return(big.NewInt(returnVal), nil).Once()
}

func MockPayGasAndUpdateCCTX(
	m *crosschainmocks.CrosschainFungibleKeeper,
	m2 *crosschainmocks.CrosschainObserverKeeper,
	ctx sdk.Context,
	k keeper.Keeper,
	senderChain chains.Chain,
	asset string,
) {
	m2.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
		Return(senderChain, true).Twice()
	m.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
		Return(fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().String(),
		}, true).Once()
	m.On("QuerySystemContractGasCoinZRC20", mock.Anything, mock.Anything).
		Return(ethcommon.Address{}, nil).Once()
	m.On("QueryGasLimit", mock.Anything, mock.Anything).
		Return(big.NewInt(100), nil).Once()
	m.On("QueryProtocolFlatFee", mock.Anything, mock.Anything).
		Return(big.NewInt(1), nil).Once()
	k.SetGasPrice(ctx, types.GasPrice{
		ChainId:     senderChain.ChainId,
		MedianIndex: 0,
		Prices:      []uint64{1},
	})

	m.On("QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(big.NewInt(0), nil).Once()
	m.On("DepositZRC20", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&evmtypes.MsgEthereumTxResponse{}, nil)
	m.On("GetUniswapV2Router02Address", mock.Anything).
		Return(ethcommon.Address{}, nil).Once()
	m.On("CallZRC20Approve", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	m.On("CallUniswapV2RouterSwapExactTokensForTokens", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(1000)}, nil).
		Once()
	m.On("CallZRC20Burn", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
}

func MockUpdateNonce(m *crosschainmocks.CrosschainObserverKeeper, senderChain chains.Chain) (nonce uint64) {
	nonce = uint64(1)
	tss := sample.Tss()
	m.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
		Return(senderChain, true).Once()
	m.On("GetChainNonces", mock.Anything, senderChain.ChainId).
		Return(observertypes.ChainNonces{Nonce: nonce}, true)
	m.On("GetTSS", mock.Anything).
		Return(tss, true)
	m.On("GetPendingNonces", mock.Anything, tss.TssPubkey, mock.Anything).
		Return(observertypes.PendingNonces{NonceHigh: int64(nonce)}, true)
	m.On("SetChainNonces", mock.Anything, mock.Anything)
	m.On("SetPendingNonces", mock.Anything, mock.Anything)
	return
}

func MockRevertForHandleEVMDeposit(
	m *crosschainmocks.CrosschainFungibleKeeper,
	receiver ethcommon.Address,
	amount *big.Int,
	senderChainID int64,
	errDeposit error,
) {
	m.On(
		"ZRC20DepositAndCallContract",
		mock.Anything,
		mock.Anything,
		receiver,
		amount,
		senderChainID,
		mock.Anything,
		coin.CoinType_ERC20,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&evmtypes.MsgEthereumTxResponse{VmError: "reverted"}, false, errDeposit)
}

func MockVoteOnOutboundSuccessBallot(
	m *crosschainmocks.CrosschainObserverKeeper,
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	senderChain chains.Chain,
	observer string,
) {
	m.On("VoteOnOutboundBallot", ctx, mock.Anything, cctx.GetCurrentOutboundParam().ReceiverChainId, chains.ReceiveStatus_success, observer).
		Return(true, true, observertypes.Ballot{BallotStatus: observertypes.BallotStatus_BallotFinalized_SuccessObservation}, senderChain.Name, nil).
		Once()
}

func MockVoteOnOutboundFailedBallot(
	m *crosschainmocks.CrosschainObserverKeeper,
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	senderChain chains.Chain,
	observer string,
) {
	m.On("VoteOnOutboundBallot", ctx, mock.Anything, cctx.GetCurrentOutboundParam().ReceiverChainId, chains.ReceiveStatus_failed, observer).
		Return(true, true, observertypes.Ballot{BallotStatus: observertypes.BallotStatus_BallotFinalized_FailureObservation}, senderChain.Name, nil).
		Once()
}

func MockSaveOutbound(
	m *crosschainmocks.CrosschainObserverKeeper,
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	tss observertypes.TSS,
	expectedNumberOfOutboundParams int,
) {
	m.On("RemoveFromPendingNonces",
		ctx, tss.TssPubkey, cctx.GetCurrentOutboundParam().ReceiverChainId, mock.Anything).
		Return().Times(expectedNumberOfOutboundParams)
}

func MockSaveOutboundNewRevertCreated(
	m *crosschainmocks.CrosschainObserverKeeper,
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	tss observertypes.TSS,
	expectedNumberOfOutboundParams int,
) {
	m.On("RemoveFromPendingNonces",
		ctx, tss.TssPubkey, cctx.GetCurrentOutboundParam().ReceiverChainId, mock.Anything).
		Return().Times(expectedNumberOfOutboundParams)

	m.On("GetTSS", ctx).Return(observertypes.TSS{}, true)
	m.On("SetNonceToCctx", mock.Anything, mock.Anything).Return().Once()
}

// MockCctxByNonce is a utility function using observer mock to returns a cctx of the given status from crosschain keeper
// mocks the methods called by CctxByNonce to directly return the given cctx or error
func MockCctxByNonce(
	t *testing.T,
	ctx sdk.Context,
	k keeper.Keeper,
	observerKeeper *crosschainmocks.CrosschainObserverKeeper,
	cctxStatus types.CctxStatus,
	isErr bool,
) string {
	if isErr {
		// return error on GetTSS to make CctxByNonce return error
		observerKeeper.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, false).Once()
		return ""
	}

	cctx := sample.CrossChainTx(t, sample.StringRandom(sample.Rand(), 10))
	cctx.CctxStatus = &types.Status{
		Status: cctxStatus,
	}
	k.SetCrossChainTx(ctx, *cctx)

	observerKeeper.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, true).Once()
	observerKeeper.On("GetNonceToCctx", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(observertypes.NonceToCctx{
			CctxIndex: cctx.Index,
		}, true).
		Once()

	return cctx.Index
}
