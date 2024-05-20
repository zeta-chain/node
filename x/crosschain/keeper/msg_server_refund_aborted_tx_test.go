package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
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
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, cctx.InboundParams.SenderChainId, "foobar", "foobar")

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		})
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
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.GetCurrentOutboundParam().Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		})
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
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.OutboundParams = nil
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.GetCurrentOutboundParam().Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		})
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundParams.TxOrigin)
		refundAddressCosmos := sdk.AccAddress(refundAddress.Bytes())
		balance := sdkk.BankKeeper.GetBalance(ctx, refundAddressCosmos, config.BaseDenom)
		require.Equal(t, cctx.InboundParams.Amount.Uint64(), balance.Amount.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})

	t.Run("should error if aleady refunded", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = true
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		cctx.OutboundParams = nil
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.GetCurrentOutboundParam().Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		})
		require.Error(t, err)
	})

	t.Run("should error if refund fails", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Cmd
		cctx.OutboundParams = nil
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.GetCurrentOutboundParam().Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		})
		require.Error(t, err)
	})

	t.Run("successfully refund to optional refund address if provided", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.InboundParams.Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		refundAddress := sample.EthAddress()
		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: refundAddress.String(),
		})
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
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
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

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		})
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
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, cctx.InboundParams.SenderChainId, "foobar", "foobar")

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.TxOrigin,
		})
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
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.InboundParams.Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "invalid-address",
		})
		require.ErrorContains(t, err, "invalid refund address")
	})

	t.Run("fail refund if address provided is null ", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.InboundParams.Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "0x0000000000000000000000000000000000000000",
		})
		require.ErrorContains(t, err, "invalid refund address")
	})

	t.Run("fail refund if status is not aborted", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "",
		})
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
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "",
		})
		require.ErrorContains(t, err, "cannot find cctx")
	})

	t.Run("fail refund if refund address not provided", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidBtcChainID()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, cctx.InboundParams.SenderChainId, "foobar", "foobar")

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "",
		})
		require.ErrorContains(t, err, "refund address is required")
	})

	t.Run("fail refund tx for coin-type Zeta if zeta accounting object is not present", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		})
		require.ErrorContains(t, err, "unable to find zeta accounting")
	})

	t.Run("fail refund if non admin account is the creator", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		chainID := getValidEthChainID()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, false)

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundParams.TxOrigin = cctx.InboundParams.Sender
		cctx.InboundParams.SenderChainId = chainID
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, cctx.InboundParams.SenderChainId, "foobar", "foobar")

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundParams.Sender,
		})
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
