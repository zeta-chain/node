package keeper

import (
	"errors"

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

	// check if the cctx is aborted
	if cctx.CctxStatus.Status != types.CctxStatus_Aborted {
		return nil, errorsmod.Wrap(types.ErrInvalidStatus, "CCTX is not aborted")
	}
	// check if the cctx is not refunded
	if cctx.CctxStatus.IsAbortRefunded {
		return nil, errorsmod.Wrap(types.ErrUnableProcessRefund, "CCTX is already refunded")
	}

	// Check if aborted amount is available to maintain zeta accounting
	// NOTE: Need to verify if this check works / is required in athens 3
	if cctx.InboundTxParams.CoinType == common.CoinType_Zeta {
		err := k.RemoveZetaAbortedAmount(ctx, GetAbortedAmount(cctx))
		// if the zeta accounting is not found, it means the zeta accounting is not set yet and the refund should not be processed
		if errors.Is(err, types.ErrUnableToFindZetaAccounting) {
			return nil, errorsmod.Wrap(types.ErrUnableProcessRefund, err.Error())
		}
		// if the zeta accounting is found but the amount is insufficient, it means the refund can be processed but the zeta accounting is not maintained properly
		// aborted amounts for zeta accounting would need to be updated in the envionment via a migration script
		if errors.Is(err, types.ErrInsufficientZetaAmount) {
			ctx.Logger().Error("Zeta Accounting Error: ", err)
		}
	}

	refundAddress, err := GetRefundAddress(cctx, msg.RefundAddress)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidAddress, err.Error())
	}
	// refund the amount
	// use temporary context to avoid gas refunding issues and side effects
	tmpCtx, commit := ctx.CacheContext()
	err = k.RefundAbortedAmountOnZetaChain(tmpCtx, cctx, refundAddress)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUnableProcessRefund, err.Error())
	}
	commit()

	// set the cctx as refunded
	cctx.CctxStatus.AbortRefunded(ctx.BlockTime().Unix())

	k.SetCrossChainTx(ctx, cctx)

	return &types.MsgRefundAbortedCCTXResponse{}, nil
}

// Set the proper refund address.
// For BTC sender chain the refund address is the one provided in the message in the RefundAddress field.
// For EVM chain with coin type ERC20 the refund address is the sender , but can be overridden by the RefundAddress field in the message.
// For EVM chain with coin type Zeta the refund address is the tx origin, but can be overridden by the RefundAddress field in the message.
// For EVM chain with coin type Gas the refund address is the tx origin, but can be overridden by the RefundAddress field in the message.

func GetRefundAddress(cctx types.CrossChainTx, optionalRefundAddress string) (ethcommon.Address, error) {
	// make sure a separate refund address is provided for a bitcoin chain as we cannot refund to tx origin or sender in this case
	if common.IsBitcoinChain(cctx.InboundTxParams.SenderChainId) && optionalRefundAddress == "" {
		return ethcommon.Address{}, errorsmod.Wrap(types.ErrInvalidAddress, "refund address is required for bitcoin chain")
	}
	refundAddress := ethcommon.HexToAddress(cctx.InboundTxParams.TxOrigin)
	if cctx.InboundTxParams.CoinType == common.CoinType_ERC20 {
		refundAddress = ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
	}
	if optionalRefundAddress != "" {
		if !ethcommon.IsHexAddress(optionalRefundAddress) {
			return ethcommon.Address{}, errorsmod.Wrap(types.ErrInvalidAddress, "invalid refund address provided")
		}
		refundAddress = ethcommon.HexToAddress(optionalRefundAddress)
	}
	// Double check to make sure the refund address is valid
	if refundAddress == (ethcommon.Address{}) {
		return ethcommon.Address{}, errorsmod.Wrap(types.ErrInvalidAddress, "invalid refund address")
	}
	return refundAddress, nil
}
