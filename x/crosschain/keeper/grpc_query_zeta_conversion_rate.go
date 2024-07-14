package keeper

import (
	"context"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) ConvertGasToZeta(
	context context.Context,
	request *types.QueryConvertGasToZetaRequest,
) (*types.QueryConvertGasToZetaResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)

	chain, found := chains.GetChainFromChainID(request.ChainId, k.GetAuthorityKeeper().GetAdditionalChainList(ctx))
	if !found {
		return nil, zetaObserverTypes.ErrSupportedChains
	}

	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chain.ChainId)
	if !isFound {
		return nil, status.Error(codes.InvalidArgument, "invalid request: param chain")
	}

	gasLimit := math.NewUintFromString(request.GasLimit)
	outTxGasFee := medianGasPrice.Mul(gasLimit)
	zrc20, err := k.fungibleKeeper.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chain.ChainId))
	if err != nil {
		return nil, status.Error(codes.NotFound, "zrc20 not found")
	}

	outTxGasFeeInZeta, err := k.fungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, outTxGasFee.BigInt(), zrc20)
	if err != nil {
		return nil, status.Error(codes.Internal, "zQueryUniswapv2RouterGetAmountsIn failed")
	}

	return &types.QueryConvertGasToZetaResponse{
		OutboundGasInZeta: outTxGasFeeInZeta.String(),
		ProtocolFeeInZeta: types.GetProtocolFee().String(),
		// #nosec G115 always positive
		ZetaBlockHeight: uint64(ctx.BlockHeight()),
	}, nil
}

func (k Keeper) ProtocolFee(
	_ context.Context,
	_ *types.QueryMessagePassingProtocolFeeRequest,
) (*types.QueryMessagePassingProtocolFeeResponse, error) {
	return &types.QueryMessagePassingProtocolFeeResponse{
		FeeInZeta: types.GetProtocolFee().String(),
	}, nil
}
