package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
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
		assert.Equal(t, ethcommon.Address{}, address)
	})
	t.Run("should fail if refund address is invalid", func(t *testing.T) {
		address, err := keeper.GetRefundAddress("invalid-address")
		require.ErrorIs(t, crosschaintypes.ErrInvalidAddress, err)
		assert.Equal(t, ethcommon.Address{}, address)
	})

}
func TestMsgServer_RefundAbortedCCTX(t *testing.T) {
	t.Run("successfully refund tx for coin-type gas", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, cctx.InboundTxParams.SenderChainId, "foobar", "foobar")

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundTxParams.Sender,
		})
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundTxParams.TxOrigin)
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20, refundAddress)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutTxParam().Amount.Uint64(), balance.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})
	t.Run("successfully refund tx for coin-type zeta", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.GetCurrentOutTxParam().Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundTxParams.Sender,
		})
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundTxParams.TxOrigin)
		refundAddressCosmos := sdk.AccAddress(refundAddress.Bytes())
		balance := sdkk.BankKeeper.GetBalance(ctx, refundAddressCosmos, config.BaseDenom)
		require.Equal(t, cctx.GetCurrentOutTxParam().Amount.Uint64(), balance.Amount.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})
	t.Run("successfully refund tx to inbound amount if outbound is not found for coin-type zeta", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Zeta
		cctx.OutboundTxParams = nil
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.GetCurrentOutTxParam().Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundTxParams.Sender,
		})
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundTxParams.TxOrigin)
		refundAddressCosmos := sdk.AccAddress(refundAddress.Bytes())
		balance := sdkk.BankKeeper.GetBalance(ctx, refundAddressCosmos, config.BaseDenom)
		require.Equal(t, cctx.InboundTxParams.Amount.Uint64(), balance.Amount.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})
	t.Run("successfully refund to optional refund address if provided", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.InboundTxParams.Amount})
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
		require.Equal(t, cctx.GetCurrentOutTxParam().Amount.Uint64(), balance.Amount.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})
	t.Run("successfully refund tx for coin-type ERC20", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		asset := sample.EthAddress().String()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_ERC20
		cctx.InboundTxParams.Asset = asset
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
			RefundAddress: cctx.InboundTxParams.Sender,
		})
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20Addr, refundAddress)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutTxParam().Amount.Uint64(), balance.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})
	t.Run("successfully refund tx for coin-type Gas with BTC sender", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidBtcChainID()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, cctx.InboundTxParams.SenderChainId, "foobar", "foobar")

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundTxParams.TxOrigin,
		})
		require.NoError(t, err)

		refundAddress := ethcommon.HexToAddress(cctx.InboundTxParams.TxOrigin)
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20, refundAddress)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutTxParam().Amount.Uint64(), balance.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})
	t.Run("fail refund if address provided is invalid", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.InboundTxParams.Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "invalid-address",
		})
		require.ErrorContains(t, err, "invalid refund address")
	})
	t.Run("fail refund if address provided is null ", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: cctx.InboundTxParams.Amount})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "0x0000000000000000000000000000000000000000",
		})
		require.ErrorContains(t, err, "invalid refund address")
	})
	t.Run("fail refund if status is not aborted", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Gas
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
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_PendingOutbound
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Gas
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "",
		})
		require.ErrorContains(t, err, "cannot find cctx")
	})
	t.Run("fail refund if refund address not provided", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidBtcChainID()
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, cctx.InboundTxParams.SenderChainId, "foobar", "foobar")

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "",
		})
		require.ErrorContains(t, err, "refund address is required")
	})
	t.Run("fail refund tx for coin-type Zeta if zeta accounting object is not present", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Zeta
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundTxParams.Sender,
		})
		require.ErrorContains(t, err, "unable to find zeta accounting")
	})
	t.Run("fail refund if non admin account is the creator", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		chainID := getValidEthChainID(t)
		setAdminPolicies(ctx, zk, admin)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		cctx := sample.CrossChainTx(t, "sample-index")
		cctx.CctxStatus.Status = crosschaintypes.CctxStatus_Aborted
		cctx.CctxStatus.IsAbortRefunded = false
		cctx.InboundTxParams.TxOrigin = cctx.InboundTxParams.Sender
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = common.CoinType_Gas
		k.SetCrossChainTx(ctx, *cctx)
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, cctx.InboundTxParams.SenderChainId, "foobar", "foobar")

		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       sample.AccAddress(),
			CctxIndex:     cctx.Index,
			RefundAddress: cctx.InboundTxParams.Sender,
		})
		require.ErrorIs(t, err, observertypes.ErrNotAuthorized)
	})
}
