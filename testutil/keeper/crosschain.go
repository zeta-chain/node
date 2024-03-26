package keeper

import (
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	tmdb "github.com/tendermint/tm-db"
	"github.com/zeta-chain/zetacore/common"
	crosschainmocks "github.com/zeta-chain/zetacore/testutil/keeper/mocks/crosschain"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type CrosschainMockOptions struct {
	UseBankMock      bool
	UseAccountMock   bool
	UseStakingMock   bool
	UseObserverMock  bool
	UseFungibleMock  bool
	UseAuthorityMock bool
}

var (
	CrosschainMocksAll = CrosschainMockOptions{
		UseBankMock:      true,
		UseAccountMock:   true,
		UseStakingMock:   true,
		UseObserverMock:  true,
		UseFungibleMock:  true,
		UseAuthorityMock: true,
	}
	CrosschainNoMocks = CrosschainMockOptions{}
)

// CrosschainKeeperWithMocks initializes a crosschain keeper for testing purposes with option to mock specific keepers
func CrosschainKeeperWithMocks(
	t testing.TB,
	mockOptions CrosschainMockOptions,
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
	authorityKeeperTmp := initAuthorityKeeper(cdc, db, stateStore)
	observerKeeperTmp := initObserverKeeper(
		cdc,
		db,
		stateStore,
		sdkKeepers.StakingKeeper,
		sdkKeepers.SlashingKeeper,
		sdkKeepers.ParamsKeeper,
		authorityKeeperTmp,
	)
	fungibleKeeperTmp := initFungibleKeeper(
		cdc,
		db,
		stateStore,
		sdkKeepers.ParamsKeeper,
		sdkKeepers.AuthKeeper,
		sdkKeepers.BankKeeper,
		sdkKeepers.EvmKeeper,
		observerKeeperTmp,
		authorityKeeperTmp,
	)
	zetaKeepers := ZetaKeepers{
		ObserverKeeper:  observerKeeperTmp,
		FungibleKeeper:  fungibleKeeperTmp,
		AuthorityKeeper: &authorityKeeperTmp,
	}
	var authorityKeeper types.AuthorityKeeper = authorityKeeperTmp
	var observerKeeper types.ObserverKeeper = observerKeeperTmp
	var fungibleKeeper types.FungibleKeeper = fungibleKeeperTmp

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
		authorityKeeper,
	)

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

func MockGetSupportedChainFromChainID(m *crosschainmocks.CrosschainObserverKeeper, senderChain *common.Chain) {
	m.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
		Return(senderChain).Once()

}
func MockGetRevertGasLimitForERC20(m *crosschainmocks.CrosschainFungibleKeeper, asset string, senderChain common.Chain) {
	m.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
		Return(fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().String(),
		}, true).Once()
	m.On("QueryGasLimit", mock.Anything, mock.Anything).
		Return(big.NewInt(100), nil).Once()

}
func MockPayGasAndUpdateCCTX(m *crosschainmocks.CrosschainFungibleKeeper, m2 *crosschainmocks.CrosschainObserverKeeper, ctx sdk.Context, k keeper.Keeper, senderChain common.Chain, asset string) {
	m2.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
		Return(&senderChain).Twice()
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
		Return([]*big.Int{big.NewInt(0), big.NewInt(1), big.NewInt(1000)}, nil).Once()
	m.On("CallZRC20Burn", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

}

func MockUpdateNonce(m *crosschainmocks.CrosschainObserverKeeper, senderChain common.Chain) (nonce uint64) {
	nonce = uint64(1)
	tss := sample.Tss()
	m.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
		Return(senderChain)
	m.On("GetChainNonces", mock.Anything, senderChain.ChainName.String()).
		Return(observertypes.ChainNonces{Nonce: nonce}, true)
	m.On("GetTSS", mock.Anything).
		Return(tss, true)
	m.On("GetPendingNonces", mock.Anything, tss.TssPubkey, mock.Anything).
		Return(observertypes.PendingNonces{NonceHigh: int64(nonce)}, true)
	m.On("SetChainNonces", mock.Anything, mock.Anything)
	m.On("SetPendingNonces", mock.Anything, mock.Anything)
	return
}

func MockRevertForHandleEVMDeposit(m *crosschainmocks.CrosschainFungibleKeeper, receiver ethcommon.Address, amount *big.Int, senderChainID int64, errDeposit error) {
	m.On(
		"ZRC20DepositAndCallContract",
		mock.Anything,
		mock.Anything,
		receiver,
		amount,
		senderChainID,
		mock.Anything,
		common.CoinType_ERC20,
		mock.Anything,
	).Return(&evmtypes.MsgEthereumTxResponse{VmError: "reverted"}, false, errDeposit)
}

func MockVoteOnOutboundSuccessBallot(m *crosschainmocks.CrosschainObserverKeeper, ctx sdk.Context, cctx *types.CrossChainTx, senderChain common.Chain, observer string) {
	m.On("VoteOnOutboundBallot", ctx, mock.Anything, cctx.GetCurrentOutTxParam().ReceiverChainId, common.ReceiveStatus_Success, observer).
		Return(true, true, observertypes.Ballot{BallotStatus: observertypes.BallotStatus_BallotFinalized_SuccessObservation}, senderChain.ChainName.String(), nil).Once()
}

func MockVoteOnOutboundFailedBallot(m *crosschainmocks.CrosschainObserverKeeper, ctx sdk.Context, cctx *types.CrossChainTx, senderChain common.Chain, observer string) {
	m.On("VoteOnOutboundBallot", ctx, mock.Anything, cctx.GetCurrentOutTxParam().ReceiverChainId, common.ReceiveStatus_Failed, observer).
		Return(true, true, observertypes.Ballot{BallotStatus: observertypes.BallotStatus_BallotFinalized_FailureObservation}, senderChain.ChainName.String(), nil).Once()
}

func MockGetOutBound(m *crosschainmocks.CrosschainObserverKeeper, ctx sdk.Context) {
	m.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()
}

func MockSaveOutBound(m *crosschainmocks.CrosschainObserverKeeper, ctx sdk.Context, cctx *types.CrossChainTx, tss observertypes.TSS) {
	m.On("RemoveFromPendingNonces",
		ctx, tss.TssPubkey, cctx.GetCurrentOutTxParam().ReceiverChainId, mock.Anything).
		Return().Once()
	m.On("GetTSS", ctx).Return(observertypes.TSS{}, true)
}

func MockSaveOutBoundNewRevertCreated(m *crosschainmocks.CrosschainObserverKeeper, ctx sdk.Context, cctx *types.CrossChainTx, tss observertypes.TSS) {
	m.On("RemoveFromPendingNonces",
		ctx, tss.TssPubkey, cctx.GetCurrentOutTxParam().ReceiverChainId, mock.Anything).
		Return().Once()
	m.On("GetTSS", ctx).Return(observertypes.TSS{}, true)
	m.On("SetNonceToCctx", mock.Anything, mock.Anything).Return().Once()
}
