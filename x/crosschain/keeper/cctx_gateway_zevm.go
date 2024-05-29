package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

type CCTXGatewayZEVM struct {
	crosschainKeeper Keeper
}

func NewCCTXGatewayZEVM(crosschainKeeper Keeper) CCTXGatewayZEVM {
	return CCTXGatewayZEVM{
		crosschainKeeper: crosschainKeeper,
	}
}

func (c CCTXGatewayZEVM) InitiateOutbound(ctx sdk.Context, cctx *types.CrossChainTx) error {
	tmpCtx, commit := ctx.CacheContext()
	isContractReverted, err := c.crosschainKeeper.HandleEVMDeposit(tmpCtx, cctx)

	// further processing will be in validateOutbound(...), for now keeping it here
	if err != nil && !isContractReverted { // exceptional case; internal error; should abort CCTX
		cctx.SetAbort(err.Error())
		return err
	} else if err != nil && isContractReverted { // contract call reverted; should refund via a revert tx
		revertMessage := err.Error()
		senderChain := c.crosschainKeeper.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, cctx.InboundParams.SenderChainId)
		if senderChain == nil {
			cctx.SetAbort(fmt.Sprintf("invalid sender chain id %d", cctx.InboundParams.SenderChainId))
			return nil
		}
		gasLimit, err := c.crosschainKeeper.GetRevertGasLimit(ctx, *cctx)
		if err != nil {
			cctx.SetAbort(fmt.Sprintf("revert gas limit error: %s", err.Error()))
			return nil
		}
		if gasLimit == 0 {
			// use same gas limit of outbound as a fallback -- should not be required
			gasLimit = cctx.GetCurrentOutboundParam().GasLimit
		}

		err = cctx.AddRevertOutbound(gasLimit)
		if err != nil {
			cctx.SetAbort(fmt.Sprintf("revert outbound error: %s", err.Error()))
			return nil
		}
		// we create a new cached context, and we don't commit the previous one with EVM deposit
		tmpCtxRevert, commitRevert := ctx.CacheContext()
		err = func() error {
			err := c.crosschainKeeper.PayGasAndUpdateCctx(
				tmpCtxRevert,
				senderChain.ChainId,
				cctx,
				cctx.InboundParams.Amount,
				false,
			)
			if err != nil {
				return err
			}
			// Update nonce using senderchain id as this is a revert tx and would go back to the original sender
			return c.crosschainKeeper.UpdateNonce(tmpCtxRevert, senderChain.ChainId, cctx)
		}()
		if err != nil {
			cctx.SetAbort(fmt.Sprintf("deposit revert message: %s err : %s", revertMessage, err.Error()))
			return nil
		}
		commitRevert()
		cctx.SetPendingRevert(revertMessage)
		return nil
	}
	// successful HandleEVMDeposit;
	commit()
	cctx.SetOutBoundMined("Remote omnichain contract call completed")
	return nil
}
