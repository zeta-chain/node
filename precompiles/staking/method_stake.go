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

func (c *Contract) Stake(
	ctx sdk.Context,
	evm *vm.EVM,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 3 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 3,
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

	validatorAddress, ok := args[1].(string)
	if !ok {
		return nil, ptypes.ErrInvalidArgument{
			Got: args[1],
		}
	}

	amount, ok := args[2].(*big.Int)
	if !ok {
		return nil, ptypes.ErrInvalidArgument{
			Got: args[2],
		}
	}

	msgServer := stakingkeeper.NewMsgServerImpl(&c.stakingKeeper)
	_, err := msgServer.Delegate(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(stakerAddress.Bytes()).String(),
		ValidatorAddress: validatorAddress,
		Amount: sdk.Coin{
			Denom:  c.stakingKeeper.BondDenom(ctx),
			Amount: math.NewIntFromBigInt(amount),
		},
	})
	if err != nil {
		return nil, err
	}

	// if caller is not the same as origin it means call is coming through smart contract,
	// and because state of smart contract calling precompile might be updated as well
	// manually reduce amount in stateDB, so it is properly reflected in bank module
	stateDB := evm.StateDB.(ptypes.ExtStateDB)
	if contract.CallerAddress != evm.Origin {
		stateDB.SubBalance(stakerAddress, amount)
	}

	err = c.addStakeLog(ctx, stateDB, stakerAddress, validatorAddress, amount)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}
