package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

func Test_GetRefundAddress(t *testing.T) {
	t.Run("should return refund address if provided coin-type gas", func(t *testing.T) {
		validEthAddress := sample.EthAddress()
		address, err := keeper.GetRefundAddress(validEthAddress.String())
		require.NoError(t, err)
		require.Equal(t, validEthAddress, address)
	})
	t.Run("should fail if refund address is empty", func(t *testing.T) {
		address, err := keeper.GetRefundAddress("")
		require.ErrorIs(t, crosschaintypes.ErrInvalidAddress, err)
		require.Equal(t, ethcommon.Address{}, address)
	})
	t.Run("should fail if refund address is invalid", func(t *testing.T) {
		address, err := keeper.GetRefundAddress("invalid-address")
		require.ErrorIs(t, crosschaintypes.ErrInvalidAddress, err)
		require.Equal(t, ethcommon.Address{}, address)
	})

}
func TestMsgServer_RefundAbortedCCTX(t *testing.T) {
	t.Run("successfully refund tx for coin-type gas", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			cctx.InboundParams.SenderChainId,
			"foobar",
			"foobar",
		)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundParams.TxOrigin)
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20, refundAddress)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount.Uint64(), balance.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})

	t.Run("successfully refund tx for coin-type zeta", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(
			ctx,
			crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.GetCurrentOutboundParam().Amount},
		)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundParams.TxOrigin)
		refundAddressCosmos := sdk.AccAddress(refundAddress.Bytes())
		balance := sdkk.BankKeeper.GetBalance(ctx, refundAddressCosmos, config.BaseDenom)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount.Uint64(), balance.Amount.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})

	t.Run("successfully refund tx to inbound amount if outbound is not found for coin-type zeta", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.OutboundParams = nil
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(
			ctx,
			crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.GetCurrentOutboundParam().Amount},
		)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundParams.TxOrigin)
		refundAddressCosmos := sdk.AccAddress(refundAddress.Bytes())
		balance := sdkk.BankKeeper.GetBalance(ctx, refundAddressCosmos, config.BaseDenom)
		require.Equal(t, cctx.InboundParams.Amount.Uint64(), balance.Amount.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})

	t.Run("should error if already refunded", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = true
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.OutboundParams = nil
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(
			ctx,
			crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.GetCurrentOutboundParam().Amount},
		)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("should error if refund fails", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Cmd
		cctx.OutboundParams = nil
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(
			ctx,
			crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.GetCurrentOutboundParam().Amount},
		)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.Error(t, err)
	})

	t.Run("successfully refund to optional refund address if provided", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")
		refundAddress := sample.EthAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.InboundParams.Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: refundAddress.String(),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.NoError(t, err)

		refundAddressCosmos := sdk.AccAddress(refundAddress.Bytes())
		balance := sdkk.BankKeeper.GetBalance(ctx, refundAddressCosmos, config.BaseDenom)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount.Uint64(), balance.Amount.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})

	t.Run("successfully refund tx for coin-type ERC20", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		asset := sample.EthAddress().String()
		cctx := sample.CrossChainTx(t, "sample-index")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_ERC20
		cctx.InboundParams.Asset = asset
		k.SetCrossChainTx(ctx, *cctx)
		// deploy zrc20
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20Addr := deployZRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			asset,
			"bar",
		)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundParams.Sender)
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20Addr, refundAddress)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount.Uint64(), balance.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})

	t.Run("successfully refund tx for coin-type Gas with BTC sender", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidBtcChainID()
		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			cctx.InboundParams.SenderChainId,
			"foobar",
			"foobar",
		)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.TxOrigin,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundParams.TxOrigin)
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20, refundAddress)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutboundParam().Amount.Uint64(), balance.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})

	t.Run("fail refund if address provided is invalid", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.InboundParams.Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "invalid-address",
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.ErrorContains(t, err, "invalid refund address")
	})

	t.Run("fail refund if address provided is null ", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.InboundParams.Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: constant.EVMZeroAddress,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.ErrorContains(t, err, "invalid refund address")
	})

	t.Run("fail refund if status is not aborted", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		cctx := sample.CrossChainTx(t, "sample-index")

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "",
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.ErrorContains(t, err, "CCTX is not aborted")
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.False(t, c.CctxStatus.IsAbortRefunded)
	})

	t.Run("fail refund if status cctx not found", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "",
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.ErrorContains(t, err, "cannot find cctx")
	})

	t.Run("fail refund if refund address not provided", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidBtcChainID()
		cctx := sample.CrossChainTx(t, "sample-index")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		_ = setupGasCoin(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			cctx.InboundParams.SenderChainId,
			"foobar",
			"foobar",
		)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "",
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.ErrorContains(t, err, "refund address is required")
	})

	t.Run("fail refund tx for coin-type Zeta if zeta accounting object is not present", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.ErrorContains(t, err, "unable to find zeta accounting")
	})

	t.Run("fail refund if non admin account is the creator", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		cctx := sample.CrossChainTx(t, "sample-index")
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		_ = setupGasCoin(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			cctx.InboundParams.SenderChainId,
			"foobar",
			"foobar",
		)

		msg := crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.RefundAbortedCCTX(ctx, &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
