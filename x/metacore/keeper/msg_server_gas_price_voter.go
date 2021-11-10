package keeper

import (
	"context"
	"sort"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) GasPriceVoter(goCtx context.Context, msg *types.MsgGasPriceVoter) (*types.MsgGasPriceVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	chain := msg.Chain
	gasPrice, isFound := k.GetGasPrice(ctx, chain)
	if !isFound {
		gasPrice = types.GasPrice{
			Creator:     msg.Creator,
			Index:       chain,
			Chain:       chain,
			Prices:      map[string]uint64{
				msg.Creator: msg.Price,
			},
			Median:      msg.Price,
			BlockNum:    map[string]uint64{
				msg.Creator: msg.BlockNumber,
			},
			MedianBlock: msg.BlockNumber,
		}
	} else {
		signer := msg.Creator
		gasPrice.Prices[signer] = msg.Price
		gasPrice.BlockNum[signer] = msg.BlockNumber
		gasPrice.Median, gasPrice.MedianBlock = calMedian(gasPrice.Prices, gasPrice.BlockNum)

	}
	k.SetGasPrice(ctx, gasPrice)

	return &types.MsgGasPriceVoterResponse{}, nil
}

func calMedian(prices map[string]uint64, blocks map[string]uint64) (uint64, uint64) {

	p := []SignerPrice{}
	for signer, price := range prices {
		p = append(p, SignerPrice{signer, price})
	}
	sort.Slice(p, func(i,j int) bool {
		return p[i].Price < p[j].Price
	})
	median := p[len(p)/2]

	return median.Price, blocks[median.Signer]
}

type SignerPrice struct {
	Signer string
	Price uint64
}
