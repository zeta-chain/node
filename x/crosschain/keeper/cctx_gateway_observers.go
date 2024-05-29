package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

type CCTXGatewayObservers struct {
	crosschainKeeper Keeper
}

func NewCCTXGatewayObservers(crosschainKeeper Keeper) CCTXGatewayObservers {
	return CCTXGatewayObservers{
		crosschainKeeper: crosschainKeeper,
	}
}

func (c CCTXGatewayObservers) InitiateOutbound(ctx sdk.Context, cctx *types.CrossChainTx) error {
	tmpCtx, commit := ctx.CacheContext()
	outboundReceiverChainID := cctx.GetCurrentOutboundParam().ReceiverChainId
	err := func() error {
		err := c.crosschainKeeper.PayGasAndUpdateCctx(
			tmpCtx,
			outboundReceiverChainID,
			cctx,
			cctx.InboundParams.Amount,
			false,
		)
		if err != nil {
			return err
		}
		return c.crosschainKeeper.UpdateNonce(tmpCtx, outboundReceiverChainID, cctx)
	}()
	if err != nil {
		// do not commit anything here as the CCTX should be aborted
		cctx.SetAbort(err.Error())
		return nil
	}
	commit()
	cctx.SetPendingOutbound("")
	return nil
}
