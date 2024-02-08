package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func Test_GetRefundAddress(t *testing.T) {
	t.Run("should return refund address if provided coin-type gas", func(t *testing.T) {
		validEthAddress := sample.EthAddress()
		address, err := keeper.GetRefundAddress(crosschaintypes.CrossChainTx{
			InboundTxParams: &crosschaintypes.InboundTxParams{
				TxOrigin:      validEthAddress.String(),
				CoinType:      common.CoinType_Gas,
				SenderChainId: getValidEthChainID(t),
			}},
			"")
		require.NoError(t, err)
		require.Equal(t, validEthAddress, address)
	})
	t.Run("should return refund address if provided coin-type zeta", func(t *testing.T) {
		validEthAddress := sample.EthAddress()
		address, err := keeper.GetRefundAddress(crosschaintypes.CrossChainTx{
			InboundTxParams: &crosschaintypes.InboundTxParams{
				TxOrigin:      validEthAddress.String(),
				CoinType:      common.CoinType_Zeta,
				SenderChainId: getValidEthChainID(t),
			}},
			"")
		require.NoError(t, err)
		require.Equal(t, validEthAddress, address)
	})
	t.Run("should return refund address if provided coin-type erc20", func(t *testing.T) {
		validEthAddress := sample.EthAddress()
		address, err := keeper.GetRefundAddress(crosschaintypes.CrossChainTx{
			InboundTxParams: &crosschaintypes.InboundTxParams{
				Sender:        validEthAddress.String(),
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
			}},
			"")
		require.NoError(t, err)
		require.Equal(t, validEthAddress, address)
	})
	t.Run("should return refund address if provided coin-type gas for btc chain", func(t *testing.T) {
		validEthAddress := sample.EthAddress()
		address, err := keeper.GetRefundAddress(crosschaintypes.CrossChainTx{
			InboundTxParams: &crosschaintypes.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: getValidBtcChainID(),
			}},
			validEthAddress.String())
		require.NoError(t, err)
		require.Equal(t, validEthAddress, address)
	})
	t.Run("fail if refund address is not provided for btc chain", func(t *testing.T) {
		_, err := keeper.GetRefundAddress(crosschaintypes.CrossChainTx{
			InboundTxParams: &crosschaintypes.InboundTxParams{
				CoinType:      common.CoinType_Gas,
				SenderChainId: getValidBtcChainID(),
			}},
			"")
		require.ErrorContains(t, err, "refund address is required for bitcoin chain")
	})
	t.Run("address overridden if optional address is provided", func(t *testing.T) {
		validEthAddress := sample.EthAddress()
		address, err := keeper.GetRefundAddress(crosschaintypes.CrossChainTx{
			InboundTxParams: &crosschaintypes.InboundTxParams{
				Sender:        sample.EthAddress().String(),
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
			}},
			validEthAddress.String())
		require.NoError(t, err)
		require.Equal(t, validEthAddress, address)
		address, err = keeper.GetRefundAddress(crosschaintypes.CrossChainTx{
			InboundTxParams: &crosschaintypes.InboundTxParams{
				Sender:        sample.EthAddress().String(),
				CoinType:      common.CoinType_Zeta,
				SenderChainId: getValidEthChainID(t),
			}},
			validEthAddress.String())
		require.NoError(t, err)
		require.Equal(t, validEthAddress, address)
		address, err = keeper.GetRefundAddress(crosschaintypes.CrossChainTx{
			InboundTxParams: &crosschaintypes.InboundTxParams{
				Sender:        sample.EthAddress().String(),
				CoinType:      common.CoinType_Gas,
				SenderChainId: getValidEthChainID(t),
			}},
			validEthAddress.String())
		require.NoError(t, err)
		require.Equal(t, validEthAddress, address)
	})

}
func TestMsgServer_RefundAbortedCCTX(t *testing.T) {
	t.Run("Successfully refund tx for coin-type Gas", func(t *testing.T) {
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
			RefundAddress: "",
		})
		require.NoError(t, err)
		refundAddress := ethcommon.HexToAddress(cctx.InboundTxParams.TxOrigin)
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20, refundAddress)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundTxParams.Amount.Uint64(), balance.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})
	t.Run("Successfully refund tx for coin-type Zeta", func(t *testing.T) {
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
			RefundAddress: "",
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
	t.Run("Successfully refund to optional refund address if provided", func(t *testing.T) {
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
		require.Equal(t, cctx.InboundTxParams.Amount.Uint64(), balance.Amount.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})
	t.Run("Successfully refund tx for coin-type ERC20", func(t *testing.T) {
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
			RefundAddress: "",
		})
		require.NoError(t, err)
		refundAddress := ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20Addr, refundAddress)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundTxParams.Amount.Uint64(), balance.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})
	t.Run("Successfully refund tx for coin-type Gas with BTC sender", func(t *testing.T) {
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
		require.Equal(t, cctx.InboundTxParams.Amount.Uint64(), balance.Uint64())
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.True(t, c.CctxStatus.IsAbortRefunded)
	})
	t.Run("Fail refund if address provided is invalid", func(t *testing.T) {
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
	t.Run("Fail refund if address provided is invalid 2 ", func(t *testing.T) {
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
	t.Run("Fail refund if status is not aborted", func(t *testing.T) {
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
	t.Run("Fail refund if status cctx not found", func(t *testing.T) {
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
	t.Run("Fail refund if refund address not provided for BTC chain", func(t *testing.T) {
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
		require.ErrorContains(t, err, "refund address is required for bitcoin chain")
	})
	t.Run("Fail refund tx for coin-type Zeta if zeta accounting object is not present", func(t *testing.T) {
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
			RefundAddress: "",
		})
		require.ErrorContains(t, err, "unable to find zeta accounting")
	})
	t.Run("Fail refund tx for coin-type Zeta if zeta accounting does not have enough aborted amount", func(t *testing.T) {
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
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{AbortedZetaAmount: sdkmath.OneUint()})
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		_, err := msgServer.RefundAbortedCCTX(ctx, &crosschaintypes.MsgRefundAbortedCCTX{
			Creator:       admin,
			CctxIndex:     cctx.Index,
			RefundAddress: "",
		})
		require.ErrorContains(t, err, "insufficient zeta amount")
	})
}
