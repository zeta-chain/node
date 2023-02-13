package keeper

import (
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

type msgServer struct {
	Keeper
	bankkeeper types.BankKeeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, bank types.BankKeeper) types.MsgServer {
	return &msgServer{
		Keeper:     keeper,
		bankkeeper: bank,
	}
}

var _ types.MsgServer = msgServer{}
