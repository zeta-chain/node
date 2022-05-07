package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"sort"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k msgServer) ZetaConversionRateVoter(goCtx context.Context, msg *types.MsgZetaConversionRateVoter) (*types.MsgZetaConversionRateVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !isBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	chain, err := common.ParseChain(msg.Chain)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnsupportedChain, fmt.Sprintf("chain %s not supported", msg.Chain))
	}
	rate, isFound := k.GetZetaConversionRate(ctx, chain.String())
	nativeTokenSymbol := chain.GetNativeTokenSymbol()
	if !isFound {
		rate = types.ZetaConversionRate{
			Index:               chain.String(),
			Chain:               chain.String(),
			Signers:             []string{msg.Creator},
			BlockNums:           []uint64{msg.BlockNumber},
			ZetaConversionRates: []string{msg.ZetaConversionRate},
			NativeTokenSymbol:   nativeTokenSymbol,
			MedianIndex:         0,
		}
	} else {
		signers := rate.Signers
		exist := false
		for i, s := range signers {
			if s == msg.Creator { // update existing entry
				rate.BlockNums[i] = msg.BlockNumber
				rate.ZetaConversionRates[i] = msg.ZetaConversionRate
				exist = true
				break
			}
		}
		if !exist {
			rate.Signers = append(rate.Signers, msg.Creator)
			rate.BlockNums = append(rate.BlockNums, msg.BlockNumber)
			rate.ZetaConversionRates = append(rate.ZetaConversionRates, msg.ZetaConversionRate)
		}
		mi := medianOfArrayFloat(rate.ZetaConversionRates)
		rate.MedianIndex = uint64(mi)
	}
	k.SetZetaConversionRate(ctx, rate)

	return &types.MsgZetaConversionRateVoterResponse{}, nil
}

type indexFloatValue struct {
	Index int
	Value float64
}

func medianOfArrayFloat(values []string) int {
	var array []indexFloatValue
	for i, v := range values {
		f, _ := strconv.ParseFloat(v, 64)
		array = append(array, indexFloatValue{Index: i, Value: f})
	}
	sort.SliceStable(array, func(i, j int) bool {
		return array[i].Value < array[j].Value
	})
	l := len(array)
	return array[l/2].Index
}
