package staking

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	ptypes "github.com/zeta-chain/node/precompiles/types"
)

func (c *Contract) MoveStake(
	ctx sdk.Context,
	evm *vm.EVM,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 4 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 4,
		})
	}

	stakerAddress, ok := args[0].(common.Address)
	if !ok {
		return nil, ptypes.ErrInvalidArgument{
			Got: args[0],
		}
	}

	if contract.CallerAddress != stakerAddress {
		return nil, fmt.Errorf("caller is not staker address")
	}

	validatorSrcAddress, ok := args[1].(string)
	if !ok {
		return nil, ptypes.ErrInvalidArgument{
			Got: args[1],
		}
	}

	validatorDstAddress, ok := args[2].(string)
	if !ok {
		return nil, ptypes.ErrInvalidArgument{
			Got: args[2],
		}
	}

	amount, ok := args[3].(*big.Int)
	if !ok {
		return nil, ptypes.ErrInvalidArgument{
			Got: args[3],
		}
	}

	msgServer := stakingkeeper.NewMsgServerImpl(&c.stakingKeeper)
	res, err := msgServer.BeginRedelegate(ctx, &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    sdk.AccAddress(stakerAddress.Bytes()).String(),
		ValidatorSrcAddress: validatorSrcAddress,
		ValidatorDstAddress: validatorDstAddress,
		Amount: sdk.Coin{
			Denom:  c.stakingKeeper.BondDenom(ctx),
			Amount: math.NewIntFromBigInt(amount),
		},
	})
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(ptypes.ExtStateDB)
	err = c.addMoveStakeLog(ctx, stateDB, stakerAddress, validatorSrcAddress, validatorDstAddress, amount)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(res.GetCompletionTime().UTC().Unix())
}
