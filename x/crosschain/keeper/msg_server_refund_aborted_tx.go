package keeper

import (
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"golang.org/x/net/context"
)

func (k msgServer) RefundAbortedTx(goCtx context.Context, msg *types.MsgRefundAbortedCCTX) (*types.MsgRefundAbortedCCTXResponse, error) {

	return &types.MsgRefundAbortedCCTXResponse{}, nil
}
