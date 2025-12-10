package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/node/pkg/coin"
	cctxerror "github.com/zeta-chain/node/pkg/errors"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// CCTXGatewayZEVM is implementation of CCTXGateway interface for ZEVM
type CCTXGatewayZEVM struct {
	crosschainKeeper Keeper
}

// NewCCTXGatewayZEVM returns new instance of CCTXGatewayZEVM
func NewCCTXGatewayZEVM(crosschainKeeper Keeper) CCTXGatewayZEVM {
	return CCTXGatewayZEVM{
		crosschainKeeper: crosschainKeeper,
	}
}

// InitiateOutbound handles evm deposit and immediately validates pending outbound
func (c CCTXGatewayZEVM) InitiateOutbound(
	ctx sdk.Context,
	config InitiateOutboundConfig,
) (newCCTXStatus types.CctxStatus, err error) {
	switch config.CCTX.InboundParams.Status {
	case types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE:
		// abort if CCTX has insufficient depositor fee for Bitcoin, the CCTX can't be reverted in this case
		// because there is no fund to pay for the revert tx
		depositErr := errors.New(config.CCTX.InboundParams.ErrorMessage)
		c.crosschainKeeper.ProcessAbort(ctx, config.CCTX, types.StatusMessages{
			ErrorMessageOutbound: depositErr.Error(),
			StatusMessage:        "inbound observation failed",
		})
		return types.CctxStatus_Aborted, nil
	case types.InboundStatus_INVALID_MEMO, types.InboundStatus_EXCESSIVE_NOASSETCALL_FUNDS:
		// when invalid memo or excessive funds is reported, the CCTX is reverted to the sender
		depositErr := errors.New(config.CCTX.InboundParams.ErrorMessage)
		newCCTXStatus = c.crosschainKeeper.ValidateOutboundZEVM(ctx, config.CCTX, depositErr, true)
		return newCCTXStatus, nil
	case types.InboundStatus_SUCCESS:
		// process the deposit normally
		if config.CCTX.InboundParams.CoinType == coin.CoinType_Zeta && config.CCTX.ProtocolContractVersion == types.ProtocolContractVersion_V2 {
			config.CCTX.SetAbort(types.StatusMessages{
				StatusMessage: types.ErrZetaThroughGateway.Error(),
			})
			return types.CctxStatus_Aborted, nil
		}
		tmpCtx, commit := ctx.CacheContext()
		isContractReverted, err := c.crosschainKeeper.HandleEVMDeposit(tmpCtx, config.CCTX)

		if err != nil && !isContractReverted {
			// exceptional case; internal error; should abort CCTX
			// use ctx as tmpCtx is dismissed to not save any side effects performed during the evm deposit
			c.crosschainKeeper.ProcessAbort(ctx, config.CCTX, types.StatusMessages{
				StatusMessage:        "outbound failed but the universal contract did not revert",
				ErrorMessageOutbound: cctxerror.NewCCTXErrorJSONMessage("failed to deposit tokens in ZEVM", err),
			})
			return types.CctxStatus_Aborted, err
		}

		newCCTXStatus = c.crosschainKeeper.ValidateOutboundZEVM(ctx, config.CCTX, err, isContractReverted)
		if newCCTXStatus == types.CctxStatus_OutboundMined || newCCTXStatus == types.CctxStatus_PendingRevert {
			commit()
		}

		return newCCTXStatus, nil
	default:
		// unknown observation status, abort the CCTX
		c.crosschainKeeper.ProcessAbort(ctx, config.CCTX, types.StatusMessages{
			ErrorMessageOutbound: fmt.Sprintf("invalid observation status %d", config.CCTX.InboundParams.Status),
			StatusMessage:        "inbound observation failed",
		})
		return types.CctxStatus_Aborted, nil
	}
}
