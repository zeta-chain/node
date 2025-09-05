package keeper_test

import (
	"encoding/base64"
	"errors"
	"math/big"
	"testing"

	"github.com/zeta-chain/node/e2e/contracts/dapp"

	cosmoserror "cosmossdk.io/errors"
	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_ValidateSuccessfulOutbound(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	cctx := sample.CrossChainTx(t, "test")
	// transition to reverted if pending revert
	cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
	err := k.ValidateOutboundObservers(
		ctx,
		cctx,
		observertypes.BallotStatus_BallotFinalized_SuccessObservation,
		sample.String(),
	)
	require.NoError(t, err)
	require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Reverted)

	// transition to outbound mined if pending outbound
	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	err = k.ValidateOutboundObservers(
		ctx,
		cctx,
		observertypes.BallotStatus_BallotFinalized_SuccessObservation,
		sample.String(),
	)
	require.NoError(t, err)
	require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_OutboundMined)

	// do nothing if it's in any other state
	err = k.ValidateOutboundObservers(
		ctx,
		cctx,
		observertypes.BallotStatus_BallotFinalized_SuccessObservation,
		sample.String(),
	)
	require.NoError(t, err)
	require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_OutboundMined)
}

func TestKeeper_ValidateFailedOutboundV2(t *testing.T) {
	t.Run("cctx aborted if can't fetch ZetaChain chain ID", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound

		ctx = ctx.WithChainID("invalid")

		// act
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)

		// assert
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("cctx aborted if CCTX source chain is not ZetaChain", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.InboundParams.SenderChainId = chains.Ethereum.ChainId

		// act
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)

		// assert
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("cctx aborted if fail to add revert outbound", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTxV2(t, "test")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId
		// fail to add revert outbound if it already exists
		cctx.OutboundParams = []*types.OutboundParams{
			sample.OutboundParams(sample.Rand()),
			sample.OutboundParams(sample.Rand()),
		}

		// act
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)

		// assert
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("cctx aborted if the revert call get reverted", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		cctx := sample.CrossChainTx(t, "test")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)

		cctx.OutboundParams = []*types.OutboundParams{sample.OutboundParams(sample.Rand())}

		revertErr := vm.ErrExecutionReverted
		keepertest.MockProcessV2RevertDeposit(fungibleMock, nil, revertErr)

		// act
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)

		// assert
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("cctx aborted if the revert call can't be processed", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		cctx := sample.CrossChainTx(t, "test")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)

		cctx.OutboundParams = []*types.OutboundParams{sample.OutboundParams(sample.Rand())}

		keepertest.MockProcessV2RevertDeposit(fungibleMock, nil, errors.New("error that is not onRevert reverts"))

		// act
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)

		// assert
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("process failed outbound with a revert", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		cctx := sample.CrossChainTx(t, "test")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)

		cctx.OutboundParams = []*types.OutboundParams{sample.OutboundParams(sample.Rand())}

		keepertest.MockProcessV2RevertDeposit(fungibleMock, &evmtypes.MsgEthereumTxResponse{}, nil)

		// act
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)

		// assert
		require.NoError(t, err)
		require.Len(t, cctx.OutboundParams, 2)
		require.EqualValues(t, types.CctxStatus_Reverted, cctx.CctxStatus.Status)
	})

	t.Run("process aborting the cctx if is pending revert", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		cctx := sample.CrossChainTx(t, "test")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId

		// act
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)

		// assert
		require.NoError(t, err)
		require.EqualValues(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("cctx aborted if cctx status is invalid", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		cctx := sample.CrossChainTx(t, "test")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.CctxStatus.Status = types.CctxStatus_OutboundMined

		// act
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)

		// assert
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})
}

func TestKeeper_ValidateFailedOutbound(t *testing.T) {
	t.Run("successfully validates failed outbound set to aborted for type cmd", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundParams.CoinType = coin.CoinType_Cmd
		cctx.InboundParams.SenderChainId = chains.Ethereum.ChainId
		cctx.OutboundParams[0].ReceiverChainId = chains.Ethereum.ChainId
		cctx.OutboundParams[1].ReceiverChainId = chains.Ethereum.ChainId
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutboundParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("set failed zevm outbound of cointype ERC20 to aborted", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundParams.CoinType = coin.CoinType_ERC20
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId
		cctx.OutboundParams[0].ReceiverChainId = chains.Ethereum.ChainId
		cctx.OutboundParams[1].ReceiverChainId = chains.Ethereum.ChainId
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutboundParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("set failed zevm outbound of cointype Gas to aborted", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId
		cctx.OutboundParams[0].ReceiverChainId = chains.Ethereum.ChainId
		cctx.OutboundParams[1].ReceiverChainId = chains.Ethereum.ChainId
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutboundParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("successfully validate failed outbound if original sender is a address", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.CrosschainKeeper(t)
		receiver := sample.EthAddress()
		cctx := GetERC20Cctx(t, receiver, chains.Goerli, "", big.NewInt(42))
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		err := sdkk.EvmKeeper.SetAccount(
			ctx,
			ethcommon.HexToAddress(cctx.InboundParams.Sender),
			*statedb.NewEmptyAccount(),
		)
		require.NoError(t, err)
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId
		err = k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Reverted, cctx.CctxStatus.Status)
		require.Equal(t, cctx.GetCurrentOutboundParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})
	t.Run("cctx aborted if GetCCTXIndexBytes fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		receiver := sample.EthAddress()
		cctx := GetERC20Cctx(t, receiver, chains.Goerli, "", big.NewInt(42))
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.Index = ""
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("cctx aborted if Adding Revert fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("cctx aborted if ZETARevertAndCallContract fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		errorFailedZETARevertAndCallContract := cosmoserror.New("test", 999, "failed ZETARevertAndCallContract")
		cctx := GetERC20Cctx(t, receiver, chains.Goerli, "", big.NewInt(42))
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId
		fungibleMock.On("ZETARevertAndCallContract", mock.Anything,
			ethcommon.HexToAddress(cctx.InboundParams.Sender),
			ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver),
			cctx.InboundParams.SenderChainId,
			cctx.GetCurrentOutboundParam().ReceiverChainId,
			cctx.GetCurrentOutboundParam().Amount.BigInt(),
			mock.Anything,
			mock.Anything).Return(nil, errorFailedZETARevertAndCallContract).Once()
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("successfully revert failed outbound if original sender is a contract", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		_ = zk.FungibleKeeper.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := GetERC20Cctx(t, sample.EthAddress(), chains.Goerli, "", big.NewInt(42))
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.RelayedMessage = base64.StdEncoding.EncodeToString([]byte("sample message"))

		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		dAppContract, err := zk.FungibleKeeper.DeployContract(ctx, dapp.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)
		cctx.InboundParams.Sender = dAppContract.String()
		cctx.InboundParams.SenderChainId = chains.ZetaChainMainnet.ChainId

		err = k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Reverted, cctx.CctxStatus.Status)
		require.Equal(t, cctx.GetCurrentOutboundParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)

		dappAbi, err := dapp.DappMetaData.GetAbi()
		require.NoError(t, err)
		res, err := zk.FungibleKeeper.CallEVM(
			ctx,
			*dappAbi,
			fungibletypes.ModuleAddressEVM,
			dAppContract,
			big.NewInt(0),
			nil,
			false,
			false,
			"zetaTxSenderAddress",
		)
		require.NoError(t, err)
		unpacked, err := dappAbi.Unpack("zetaTxSenderAddress", res.Ret)
		require.NoError(t, err)
		require.NotZero(t, len(unpacked))
		valSenderAddress, ok := unpacked[0].([]byte)
		require.True(t, ok)
		require.Equal(t, dAppContract.Bytes(), valSenderAddress)
	})

	t.Run("successfully validate failed outbound set to pending revert", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, senderChain, asset)

		// mock successful UpdateNonce
		_ = keepertest.MockUpdateNonce(observerMock, senderChain)

		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingRevert)
		require.Equal(t, types.TxFinalizationStatus_NotFinalized, cctx.GetCurrentOutboundParam().TxFinalizationStatus)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundParams[0].TxFinalizationStatus)
	})

	t.Run("successfully validate failed outbound set to pending revert if gas limit is 0", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 0)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, senderChain, asset)

		// mock successful UpdateNonce
		_ = keepertest.MockUpdateNonce(observerMock, senderChain)

		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingRevert)
		require.Equal(t, types.TxFinalizationStatus_NotFinalized, cctx.GetCurrentOutboundParam().TxFinalizationStatus)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.OutboundParams[0].TxFinalizationStatus)
	})

	t.Run("aborted if update nonce fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, senderChain, asset)

		// mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock failed UpdateNonce
		observerMock.On("GetChainNonces", mock.Anything, senderChain.ChainId).
			Return(observertypes.ChainNonces{}, false)

		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("aborted if PayGasAndUpdateCctx fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)

		// mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock failed failed GetSupportedChainFromChainID to fail PayGasAndUpdateCctx
		keepertest.MockFailedGetSupportedChainFromChainID(observerMock, senderChain)

		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("aborted if GetRevertGasLimit fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""

		// mock failed GetRevertGasLimit for ERC20
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{
				Zrc20ContractAddress: sample.EthAddress().String(),
			}, false).Once()

		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})
}

func TestKeeper_ValidateOutboundObservers(t *testing.T) {
	t.Run("successfully validate outbound with ballot finalized to success", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := GetERC20Cctx(t, sample.EthAddress(), chains.Goerli, "", big.NewInt(42))
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_SuccessObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_OutboundMined)
	})

	t.Run("successfully validate outbound with ballot finalized to failed and old status is Pending Revert",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.CrosschainKeeper(t)
			cctx := GetERC20Cctx(t, sample.EthAddress(), chains.Goerli, "", big.NewInt(42))
			cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
			err := k.ValidateOutboundObservers(
				ctx,
				cctx,
				observertypes.BallotStatus_BallotFinalized_FailureObservation,
				sample.String(),
			)
			require.NoError(t, err)
			require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
			require.Equal(t, cctx.GetCurrentOutboundParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
		},
	)

	t.Run("successfully validate outbound with ballot finalized to failed and coin-type is CMD", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := GetERC20Cctx(t, sample.EthAddress(), chains.Goerli, "", big.NewInt(42))
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.InboundParams.CoinType = coin.CoinType_Cmd
		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		require.Equal(t, cctx.GetCurrentOutboundParam().TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("do not validate if cctx invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := GetERC20Cctx(t, sample.EthAddress(), chains.Goerli, "", big.NewInt(42))
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		cctx.InboundParams = nil
		err := k.ValidateOutboundObservers(ctx, cctx, observertypes.BallotStatus_BallotInProgress, sample.String())
		require.Error(t, err)
	})

	t.Run("no new outbound created when aborted ", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""

		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		oldOutboundParamsLen := len(cctx.OutboundParams)
		// mock failed GetRevertGasLimit for ERC20
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{
				Zrc20ContractAddress: sample.EthAddress().String(),
			}, false).Once()

		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)

		// New outbound not added and the old outbound is not finalized
		require.Len(t, cctx.OutboundParams, oldOutboundParamsLen)
		require.Equal(t, cctx.GetCurrentOutboundParam().TxFinalizationStatus, types.TxFinalizationStatus_NotFinalized)
	})

	t.Run("cctx aborted if the cctx has already been reverted once", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""

		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.OutboundParams = append(cctx.OutboundParams, sample.OutboundParams(sample.Rand()))
		cctx.OutboundParams[1].ReceiverChainId = 5
		cctx.OutboundParams[1].BallotIndex = ""
		cctx.OutboundParams[1].Hash = ""

		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)

		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
	})

	t.Run("successfully revert a outbound and create a new revert tx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""

		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		oldOutboundParamsLen := len(cctx.OutboundParams)
		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, senderChain, asset)

		// mock successful UpdateNonce
		_ = keepertest.MockUpdateNonce(observerMock, senderChain)

		err := k.ValidateOutboundObservers(
			ctx,
			cctx,
			observertypes.BallotStatus_BallotFinalized_FailureObservation,
			sample.String(),
		)
		require.NoError(t, err)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_PendingRevert)
		// New outbound added for revert and the old outbound is finalized
		require.Len(t, cctx.OutboundParams, oldOutboundParamsLen+1)
		require.Equal(t, cctx.GetCurrentOutboundParam().TxFinalizationStatus, types.TxFinalizationStatus_NotFinalized)
		require.Equal(
			t,
			cctx.OutboundParams[oldOutboundParamsLen-1].TxFinalizationStatus,
			types.TxFinalizationStatus_Executed,
		)
	})
}
