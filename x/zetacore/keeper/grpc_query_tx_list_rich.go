package keeper

import (
	"context"

	"github.com/zeta-chain/zetacore/x/zetacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TxListRich(goCtx context.Context, req *types.QueryTxListRichRequest) (*types.QueryTxListRichResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetTxList(ctx)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}
	txlist := val.Tx
	len := len(txlist)
	from := len - int(req.Last)
	if from < 0 {
		from = 0
	}

	var sendlist []*types.Send
	for i := from; i < len; i++ {
		tx := txlist[i]
		sendHash := tx.SendHash
		send, found := k.GetSend(ctx, sendHash)
		if found {
			sendlist = append(sendlist, &send)
		}
	}

	return &types.QueryTxListRichResponse{Tx: sendlist, Length: uint64(len)}, nil
}
