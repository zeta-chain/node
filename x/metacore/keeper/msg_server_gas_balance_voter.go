package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"math/big"
	"sort"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) GasBalanceVoter(goCtx context.Context, msg *types.MsgGasBalanceVoter) (*types.MsgGasBalanceVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !isBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	chain := msg.Chain
	gasBalance, isFound := k.GetGasBalance(ctx, chain)
	if !isFound {
		gasBalance = types.GasBalance{
			Creator:     msg.Creator,
			Index:       chain,
			Chain:       chain,
			Balances:      []string{msg.Balance},
			BlockNums:   []uint64{msg.BlockNumber},
			Signers:     []string{msg.Creator},
			MedianIndex: 0,
		}
	} else {
		signers := gasBalance.Signers
		exist := false
		for i, s := range signers {
			if s == msg.Creator { // update existing entry
				gasBalance.BlockNums[i] = msg.BlockNumber
				gasBalance.Balances[i] = msg.Balance
				exist = true
				break
			}
		}
		if !exist {
			gasBalance.Signers = append(gasBalance.Signers, msg.Creator)
			gasBalance.BlockNums = append(gasBalance.BlockNums, msg.BlockNumber)
			gasBalance.Balances = append(gasBalance.Balances, msg.Balance)
		}
		// recompute the median gas price
		mi := medianOfArrayBigInt(gasBalance.Balances)
		gasBalance.MedianIndex = uint64(mi)
	}
	k.SetGasBalance(ctx, gasBalance)

	return &types.MsgGasBalanceVoterResponse{}, nil
}

type IndexValueBigInt struct {
	Index int
	Value *big.Int
}

func medianOfArrayBigInt(values []string) int {
	var array []IndexValueBigInt
	for i, v := range values {
		vv, ok := big.NewInt(0).SetString(v, 10)
		if ok {
			array = append(array, IndexValueBigInt{Index: i, Value: vv})
		}
	}
	sort.SliceStable(array, func(i, j int) bool {
		return array[i].Value.Cmp(array[j].Value) < 0
	})
	l := len(array)
	return array[l/2].Index
}
