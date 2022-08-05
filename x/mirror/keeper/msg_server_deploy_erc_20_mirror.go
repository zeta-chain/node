package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/mirror/types"
)

func (k msgServer) DeployERC20Mirror(goCtx context.Context, msg *types.MsgDeployERC20Mirror) (*types.MsgDeployERC20MirrorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tokenPairs, found := k.GetERC20TokenPairs(ctx)
	if !found {
		tokenPairs = types.ERC20TokenPairs{TokenPairs: []*types.ERC20TokenPair{}}
	}
	for _, tokenPair := range tokenPairs.TokenPairs {
		if tokenPair.HomeERC20ContractAddress == msg.HomeERC20ContractAddress {
			return nil, sdkerrors.Wrap(types.ErrTOkenPairAlreadyExists, "toke pair already exists")
		}
	}
	tp := types.ERC20TokenPair{
		HomeERC20ContractAddress:   msg.HomeERC20ContractAddress,
		MirrorERC20ContractAddress: "",
		Name:                       msg.Name,
		Symbol:                     msg.Symbol,
		Decimals:                   msg.Decimals,
	}

	addr, err := k.DeployERC20Contract(ctx, msg.Name, msg.Symbol, uint8(msg.Decimals))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "deploy erc20 mirror error")
	}
	tp.MirrorERC20ContractAddress = addr.Hex()

	tokenPairs.TokenPairs = append(tokenPairs.TokenPairs, &tp)
	k.SetERC20TokenPairs(ctx, tokenPairs)

	return &types.MsgDeployERC20MirrorResponse{}, nil
}
