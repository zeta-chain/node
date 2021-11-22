package keeper

import (
	"context"
	"math/big"
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
			Creator: msg.Creator,
			Index:   chain,
			Chain:   chain,
			Prices: map[string]uint64{
				msg.Creator: msg.Price,
			},
			Median: msg.Price,
			BlockNum: map[string]uint64{
				msg.Creator: msg.BlockNumber,
			},
			MedianBlock: msg.BlockNumber,
			// TODO: fix supply data type; cannot use pointer
			// Otherwise conesensus failure
			//Supply: map[string]*types.GasPrice_ValueBlockPair{
			//	msg.Creator: {
			//		Value:    msg.Supply,
			//		BlockNum: msg.BlockNumber,
			//	},
			//},
			//MedianSupply: &types.GasPrice_ValueBlockPair{
			//	Value: msg.Supply,
			//	BlockNum: msg.BlockNumber,
			//},
		}
	} else {
		signer := msg.Creator
		gasPrice.Prices[signer] = msg.Price
		gasPrice.BlockNum[signer] = msg.BlockNumber
		//gasPrice.Supply[signer] = &types.GasPrice_ValueBlockPair{
		//	Value: msg.Supply,
		//	BlockNum: msg.BlockNumber,
		//}
		gasPrice.Median, gasPrice.MedianBlock = calMedian(gasPrice.Prices, gasPrice.BlockNum)
		//gasPrice.MedianSupply = calMedianSupply(gasPrice.Supply)
	}
	k.SetGasPrice(ctx, gasPrice)

	return &types.MsgGasPriceVoterResponse{}, nil
}

//func calMedianSupply(supplyMap map[string]*types.GasPrice_ValueBlockPair) *types.GasPrice_ValueBlockPair {
//	if len(supplyMap) == 0 {
//		return nil
//	}
//	p := []SignerValue{}
//	for signer, valueblock := range supplyMap {
//		supply, ok := big.NewInt(0).SetString(valueblock.Value, 10)
//		if !ok {
//			log.Error().Msgf("calMedianSupply SetString error %s ", valueblock.Value)
//		}
//		p = append(p, SignerValue{signer, supply, valueblock.BlockNum})
//	}
//	sort.SliceStable(p, func(i, j int) bool {
//		return p[i].Value.Cmp(p[j].Value) < 0
//	})
//	return &types.GasPrice_ValueBlockPair{p[len(p)/2].Value.String(), p[len(p)/2].BlockNum}
//}

//TODO: Remove calMedian; replace with calMedianSupply.
//This is allocation intensive.
func calMedian(prices map[string]uint64, blocks map[string]uint64) (uint64, uint64) {

	p := []SignerValue{}
	for signer, price := range prices {
		p = append(p, SignerValue{signer, big.NewInt(0).SetUint64(price), blocks[signer]})
	}
	sort.SliceStable(p, func(i, j int) bool {
		return p[i].Value.Cmp(p[j].Value) < 0
	})
	median := p[len(p)/2]

	return median.Value.Uint64(), median.BlockNum
}

type SignerValue struct {
	Signer   string
	Value    *big.Int
	BlockNum uint64
}
