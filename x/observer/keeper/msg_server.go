package keeper

import (
	"context"

	"github.com/zeta-chain/zetacore/x/observer/types"
)

type msgServer struct {
	Keeper
}

func (k msgServer) EnableCCTXFlags(ctx context.Context, flags *types.MsgEnableCCTXFlags) (*types.MsgEnableCCTXFlagsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k msgServer) DisableCCTXFlags(ctx context.Context, flags *types.MsgDisableCCTXFlags) (*types.MsgDisableCCTXFlagsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k msgServer) UpdateGasPriceIncreaseFlags(ctx context.Context, flags *types.MsgUpdateGasPriceIncreaseFlags) (*types.MsgUpdateGasPriceIncreaseFlagsResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
