package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetTssAddress(goCtx context.Context, req *types.QueryGetTssAddressRequest) (*types.QueryGetTssAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var tssPubKey string
	if req.TssPubKey == "" {
		tss, found := k.GetTSS(ctx)
		if !found {
			return nil, status.Error(codes.NotFound, "current tss not set")
		}
		tssPubKey = tss.TssPubkey
	} else {
		tssList := k.GetAllTSS(ctx)
		for _, t := range tssList {
			if t.TssPubkey == req.TssPubKey {
				tssPubKey = t.TssPubkey
				break
			}
		}
	}
	ethAddress, err := common.GetTssAddrEVM(tssPubKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	bitcoinParams := common.BitcoinRegnetParams
	if req.BitcoinChainId != 0 {
		bitcoinParams, err = common.BitcoinNetParamsFromChainID(req.BitcoinChainId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	btcAddress, err := common.GetTssAddrBTC(tssPubKey, bitcoinParams)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetTssAddressResponse{
		Eth: ethAddress.String(),
		Btc: btcAddress,
	}, nil
}
