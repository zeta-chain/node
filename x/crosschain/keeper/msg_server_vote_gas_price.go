package keeper

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// VoteGasPrice submits information about the connected chain's gas price at a specific block
// height. Gas price submitted by each validator is recorded separately and a
// median index is updated.
//
// Only observer validators are authorized to broadcast this message.
func (k msgServer) VoteGasPrice(
	goCtx context.Context,
	msg *types.MsgVoteGasPrice,
) (*types.MsgVoteGasPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	chain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if !found {
		return nil, cosmoserrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("ChainID : %d ", msg.ChainId))
	}
	if ok := k.zetaObserverKeeper.IsNonTombstonedObserver(ctx, msg.Creator); !ok {
		return nil, observertypes.ErrNotObserver
	}

	gasPrice, isFound := k.GetGasPrice(ctx, chain.ChainId)
	if !isFound {
		gasPrice = types.GasPrice{
			Creator:     msg.Creator,
			Index:       strconv.FormatInt(chain.ChainId, 10), // TODO : Not needed index set at keeper
			ChainId:     chain.ChainId,
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
		// #nosec G115 always positive
		gasPrice.MedianIndex = uint64(mi)
	}
	k.SetGasPrice(ctx, gasPrice)
	chainIDBigINT := big.NewInt(chain.ChainId)

	gasUsed, err := k.fungibleKeeper.SetGasPrice(
		ctx,
		chainIDBigINT,
		math.NewUint(gasPrice.Prices[gasPrice.MedianIndex]).BigInt(),
	)
	if err != nil {
		return nil, err
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
