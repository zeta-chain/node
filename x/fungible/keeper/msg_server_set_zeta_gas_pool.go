package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	clientconfig "github.com/zeta-chain/zetacore/zetaclient/config"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func (k msgServer) SetZetaGasPool(goCtx context.Context, msg *types.MsgSetZetaGasPool) (*types.MsgSetZetaGasPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	chain, found := clientconfig.Chains[msg.Chain]
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrChainNotFound, "chain %s not found", msg.Chain)
	}
	poolAddress := common.HexToAddress(msg.Address)
	if poolAddress == (common.Address{}) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAddress, "address %s is invalid", msg.Address)
	}
	if err := k.SetGasZetaPool(ctx, chain.ChainID, poolAddress); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrContractCall, "failed to set zeta gas pool: %s", err.Error())
	}

	return &types.MsgSetZetaGasPoolResponse{}, nil
}
