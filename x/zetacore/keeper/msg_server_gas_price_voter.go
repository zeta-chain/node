package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	clientconfig "github.com/zeta-chain/zetacore/zetaclient/config"
	"math/big"
	"sort"
)

func (k msgServer) GasPriceVoter(goCtx context.Context, msg *types.MsgGasPriceVoter) (*types.MsgGasPriceVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !IsBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	chain := msg.Chain
	gasPrice, isFound := k.GetGasPrice(ctx, chain)
	if !isFound {
		gasPrice = types.GasPrice{
			Creator:     msg.Creator,
			Index:       chain,
			Chain:       chain,
			Prices:      []uint64{msg.Price},
			BlockNums:   []uint64{msg.BlockNumber},
			Signers:     []string{msg.Creator},
			MedianIndex: 0,
		}
	} else {
		signers := gasPrice.Signers
		exist := false
		for i, s := range signers {
			if s == msg.Creator { // update existing entry
				gasPrice.BlockNums[i] = msg.BlockNumber
				gasPrice.Prices[i] = msg.Price
				exist = true
				break
			}
		}
		if !exist {
			gasPrice.Signers = append(gasPrice.Signers, msg.Creator)
			gasPrice.BlockNums = append(gasPrice.BlockNums, msg.BlockNumber)
			gasPrice.Prices = append(gasPrice.Prices, msg.Price)
		}
		// recompute the median gas price
		mi := medianOfArray(gasPrice.Prices)
		gasPrice.MedianIndex = uint64(mi)
	}
	k.SetGasPrice(ctx, gasPrice)

	chainid := clientconfig.Chains[chain].ChainID
	if chainid != nil {
		if err := k.fungibleKeeper.SetGasPrice(ctx, chainid, big.NewInt(int64(gasPrice.Prices[gasPrice.MedianIndex]))); err != nil {
			return nil, err
		}
	} else {
		k.Logger(ctx).Error("chainid not found", "chain", chain)
	}

	return &types.MsgGasPriceVoterResponse{}, nil
}

type indexValue struct {
	Index int
	Value uint64
}

func medianOfArray(values []uint64) int {
	array := make([]indexValue, len(values))
	for i, v := range values {
		array[i] = indexValue{Index: i, Value: v}
	}
	sort.SliceStable(array, func(i, j int) bool {
		return array[i].Value < array[j].Value
	})
	l := len(array)
	return array[l/2].Index
}
