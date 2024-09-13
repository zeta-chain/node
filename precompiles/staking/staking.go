package staking

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	ptypes "github.com/zeta-chain/node/precompiles/types"
)

// method names
const (
	// write
	StakeMethodName     = "stake"
	UnstakeMethodName   = "unstake"
	MoveStakeMethodName = "moveStake"

	// read
	GetAllValidatorsMethodName = "getAllValidators"
	GetSharesMethodName        = "getShares"
)

var (
	ABI                 abi.ABI
	ContractAddress     = common.HexToAddress("0x0000000000000000000000000000000000000066")
	GasRequiredByMethod = map[[4]byte]uint64{}
	ViewMethod          = map[[4]byte]bool{}
)

func init() {
	initABI()
}

func initABI() {
	if err := ABI.UnmarshalJSON([]byte(IStakingMetaData.ABI)); err != nil {
		panic(err)
	}

	GasRequiredByMethod = map[[4]byte]uint64{}
	for methodName := range ABI.Methods {
		var methodID [4]byte
		copy(methodID[:], ABI.Methods[methodName].ID[:4])
		switch methodName {
		// TODO: https://github.com/zeta-chain/node/issues/2812
		// just temporary flat values, double check these flat values
		// can we just use WriteCostFlat/ReadCostFlat from gas config for flat values?
		case StakeMethodName:
			GasRequiredByMethod[methodID] = 10000
		case UnstakeMethodName:
			GasRequiredByMethod[methodID] = 10000
		case MoveStakeMethodName:
			GasRequiredByMethod[methodID] = 10000
		case GetAllValidatorsMethodName:
			GasRequiredByMethod[methodID] = 0
			ViewMethod[methodID] = true
		case GetSharesMethodName:
			GasRequiredByMethod[methodID] = 0
			ViewMethod[methodID] = true
		default:
			GasRequiredByMethod[methodID] = 0
		}
	}
}

type Contract struct {
	ptypes.BaseContract

	stakingKeeper stakingkeeper.Keeper
	cdc           codec.Codec
	kvGasConfig   storetypes.GasConfig
}

func NewIStakingContract(
	stakingKeeper *stakingkeeper.Keeper,
	cdc codec.Codec,
	kvGasConfig storetypes.GasConfig,
) *Contract {
	return &Contract{
		BaseContract:  ptypes.NewBaseContract(ContractAddress),
		stakingKeeper: *stakingKeeper,
		cdc:           cdc,
		kvGasConfig:   kvGasConfig,
	}
}

// Address() is required to implement the PrecompiledContract interface.
func (c *Contract) Address() common.Address {
	return ContractAddress
}

// Abi() is required to implement the PrecompiledContract interface.
func (c *Contract) Abi() abi.ABI {
	return ABI
}

// RequiredGas is required to implement the PrecompiledContract interface.
// The gas has to be calculated deterministically based on the input.
func (c *Contract) RequiredGas(input []byte) uint64 {
	// get methodID (first 4 bytes)
	var methodID [4]byte
	copy(methodID[:], input[:4])
	// base cost to prevent large input size
	baseCost := uint64(len(input)) * c.kvGasConfig.WriteCostPerByte
	if ViewMethod[methodID] {
		baseCost = uint64(len(input)) * c.kvGasConfig.ReadCostPerByte
	}

	if requiredGas, ok := GasRequiredByMethod[methodID]; ok {
		return requiredGas + baseCost
	}

	// Can not happen, but return 0 if the method is not found.
	return 0
}

func (c *Contract) GetAllValidators(
	ctx sdk.Context,
	method *abi.Method,
) ([]byte, error) {
	validators := c.stakingKeeper.GetAllValidators(ctx)

	validatorsRes := make([]Validator, len(validators))
	for i, v := range validators {
		validatorsRes[i] = Validator{
			OperatorAddress: v.OperatorAddress,
			ConsensusPubKey: v.ConsensusPubkey.String(),
			BondStatus:      uint8(v.Status),
			Jailed:          v.Jailed,
		}
	}

	return method.Outputs.Pack(validatorsRes)
}

func (c *Contract) GetShares(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 2 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}
	stakerAddress, ok := args[0].(common.Address)
	if !ok {
		return nil, ptypes.ErrInvalidArgument{
			Got: args[0],
		}
	}

	validatorAddress, ok := args[1].(string)
	if !ok {
		return nil, ptypes.ErrInvalidArgument{
			Got: args[1],
		}
	}

	validator, err := sdk.ValAddressFromBech32(validatorAddress)
	if err != nil {
		return nil, err
	}

	delegation := c.stakingKeeper.Delegation(ctx, sdk.AccAddress(stakerAddress.Bytes()), validator)
	shares := big.NewInt(0)
	if delegation != nil {
		shares = delegation.GetShares().BigInt()
	}

	return method.Outputs.Pack(shares)
}

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

	err = c.AddStakeLog(ctx, stateDB, stakerAddress, validatorAddress, amount)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

func (c *Contract) Unstake(
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
	res, err := msgServer.Undelegate(ctx, &stakingtypes.MsgUndelegate{
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

	stateDB := evm.StateDB.(ptypes.ExtStateDB)
	err = c.AddUnstakeLog(ctx, stateDB, stakerAddress, validatorAddress, amount)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(res.GetCompletionTime().UTC().Unix())
}

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
	err = c.AddMoveStakeLog(ctx, stateDB, stakerAddress, validatorSrcAddress, validatorDstAddress, amount)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(res.GetCompletionTime().UTC().Unix())
}

// Run is the entrypoint of the precompiled contract, it switches over the input method,
// and execute them accordingly.
func (c *Contract) Run(evm *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	method, err := ABI.MethodById(contract.Input[:4])
	if err != nil {
		return nil, err
	}

	args, err := method.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(ptypes.ExtStateDB)

	switch method.Name {
	case GetAllValidatorsMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.GetAllValidators(ctx, method)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case GetSharesMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.GetShares(ctx, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case StakeMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.Stake(ctx, evm, contract, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case UnstakeMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.Unstake(ctx, evm, contract, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case MoveStakeMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.MoveStake(ctx, evm, contract, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil

	default:
		return nil, ptypes.ErrInvalidMethod{
			Method: method.Name,
		}
	}
}
