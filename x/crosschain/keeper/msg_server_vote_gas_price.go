package keeper

import (
	"context"
	"math/big"
	"sort"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// VoteGasPrice submits information about the connected chain's gas price at a specific block
// height. Gas price submitted by each validator is recorded separately and a
// median index is updated.
//
// Only observer validators are authorized to broadcast this message.
func (k msgServer) VoteGasPrice(
	cc context.Context,
	msg *types.MsgVoteGasPrice,
) (*types.MsgVoteGasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(cc)

	chain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if !found {
		return nil, cosmoserrors.Wrapf(types.ErrUnsupportedChain, "chain id %d", msg.ChainId)
	}

	err := k.zetaObserverKeeper.CheckObserverCanVote(ctx, msg.Creator)
	if err != nil {
		return nil, err
	}

	gasPrice, isFound := k.GetGasPrice(ctx, chain.ChainId)
	if !isFound {
		return k.voteGasPrice(ctx, chain, types.GasPrice{
			Creator:      msg.Creator,
			ChainId:      chain.ChainId,
			Prices:       []uint64{msg.Price},
			PriorityFees: []uint64{msg.PriorityFee},
			BlockNums:    []uint64{msg.BlockNumber},
			Signers:      []string{msg.Creator},
		})
	}

	// Now we either want to update the gas price or add a new entry
	var exists bool
	for i, s := range gasPrice.Signers {
		if s != msg.Creator {
			continue
		}

		// update existing entry
		gasPrice.BlockNums[i] = msg.BlockNumber
		gasPrice.Prices[i] = msg.Price
		gasPrice.PriorityFees[i] = msg.PriorityFee
		exists = true
		break
	}

	if !exists {
		gasPrice.Signers = append(gasPrice.Signers, msg.Creator)
		gasPrice.BlockNums = append(gasPrice.BlockNums, msg.BlockNumber)
		gasPrice.Prices = append(gasPrice.Prices, msg.Price)
		gasPrice.PriorityFees = append(gasPrice.PriorityFees, msg.PriorityFee)
	}

	// recompute the median gas price
	mi := medianOfArray(gasPrice.Prices)

	// #nosec G115 always positive
	gasPrice.MedianIndex = uint64(mi)

	return k.voteGasPrice(ctx, chain, gasPrice)
}

func (k msgServer) voteGasPrice(
	ctx sdk.Context,
	chain chains.Chain,
	entity types.GasPrice,
) (*types.MsgVoteGasPriceResponse, error) {
	var (
		chainID  = big.NewInt(chain.ChainId)
		gasPrice = math.NewUint(entity.Prices[entity.MedianIndex]).BigInt()
	)

	// set gas price in this module
	k.SetGasPrice(ctx, entity)

	// set gas price in fungible keeper (also calls EVM)
	gasUsed, err := k.fungibleKeeper.SetGasPrice(ctx, chainID, gasPrice)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set gas price in fungible keeper")
	}

	// reset the gas count
	k.ResetGasMeterAndConsumeGas(ctx, gasUsed)

	return &types.MsgVoteGasPriceResponse{}, nil
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

// ResetGasMeterAndConsumeGas reset first the gas meter consumed value to zero and set it back to the new value
// 'gasUsed'
func (k Keeper) ResetGasMeterAndConsumeGas(ctx sdk.Context, gasUsed uint64) {
	// reset the gas count
	ctx.GasMeter().RefundGas(ctx.GasMeter().GasConsumed(), "reset the gas count")
	ctx.GasMeter().ConsumeGas(gasUsed, "apply evm transaction")
}
