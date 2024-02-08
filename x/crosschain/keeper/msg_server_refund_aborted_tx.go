package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"golang.org/x/net/context"
)

func (k msgServer) RefundAbortedCCTX(goCtx context.Context, msg *types.MsgRefundAbortedCCTX) (*types.MsgRefundAbortedCCTXResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if authorized
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_group2) {
		return nil, observertypes.ErrNotAuthorized
	}

	// check if the cctx exists
	cctx, found := k.GetCrossChainTx(ctx, msg.CctxIndex)
	if !found {
		return nil, types.ErrCannotFindCctx
	}
	// make sure separate refund address is provided for bitcoin chain as we cannot refund to tx origin or sender in this case
	if common.IsBitcoinChain(cctx.InboundTxParams.SenderChainId) && msg.RefundAddress == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidAddress, "invalid refund address")
	}

	// check if the cctx is aborted
	if cctx.CctxStatus.Status != types.CctxStatus_Aborted {
		return nil, errorsmod.Wrap(types.ErrInvalidStatus, "CCTX is not aborted")
	}
	// check if the cctx is not refunded
	if cctx.CctxStatus.IsAbortRefunded {
		return nil, errorsmod.Wrap(types.ErrUnableProcessRefund, "CCTX is already refunded")
	}

	// Set the proper refund address.
	// For BTC sender chain the refund address is the one provided in the message in the RefundAddress field.
	// For EVM chain with coin type ERC20 the refund address is the sender , but can be overridden by the RefundAddress field in the message.
	// For EVM chain with coin type Zeta the refund address is the tx origin, but can be overridden by the RefundAddress field in the message.
	// For EVM chain with coin type Gas the refund address is the tx origin, but can be overridden by the RefundAddress field in the message.

	refundAddress := ethcommon.HexToAddress(cctx.InboundTxParams.TxOrigin)
	if cctx.InboundTxParams.CoinType == common.CoinType_ERC20 {
		refundAddress = ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
	}
	if msg.RefundAddress != "" {
		refundAddress = ethcommon.HexToAddress(msg.RefundAddress)
	}
	// Make sure the refund address is valid
	if refundAddress == (ethcommon.Address{}) {
		return nil, errorsmod.Wrap(types.ErrInvalidAddress, "invalid refund address")
	}

	// refund the amount
	err := k.RefundAbortedAmountOnZetaChain(ctx, cctx, refundAddress)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUnableProcessRefund, err.Error())
	}

	// set the cctx as refunded
	cctx.CctxStatus.IsAbortRefunded = true
	k.SetCrossChainTx(ctx, cctx)

	// Include the refunded amount in ZetaAccount, so we can now remove it from the ZetaAbortedAmount counter.
	if cctx.GetCurrentOutTxParam().CoinType == common.CoinType_Zeta {
		k.RemoveZetaAbortedAmount(ctx, cctx.GetCurrentOutTxParam().Amount)
	}
	return &types.MsgRefundAbortedCCTXResponse{}, nil
}
