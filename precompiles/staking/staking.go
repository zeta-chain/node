package staking

import (
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"cosmossdk.io/math"
	ptypes "github.com/zeta-chain/zetacore/precompiles/types"
)

const (
	DelegateMethodName   = "delegate"
	UndelegateMethodName = "undelegate"
	RedelegateMethodName = "redelegate"
)

var (
	ABI                 abi.ABI
	ContractAddress     = common.HexToAddress("0x0000000000000000000000000000000000000066")
	GasRequiredByMethod = map[[4]byte]uint64{}
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
		// TODO: just temporary values, double check these flat values
		// can we just use WriteCostFlat from gas config?
		case DelegateMethodName:
			GasRequiredByMethod[methodID] = 10000
		case UndelegateMethodName:
			GasRequiredByMethod[methodID] = 10000
		case RedelegateMethodName:
			GasRequiredByMethod[methodID] = 10000
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
	// base cost to prevent large input size
	baseCost := uint64(len(input)) * c.kvGasConfig.WriteCostPerByte

	// get methodID (first 4 bytes)
	var methodID [4]byte
	copy(methodID[:], input[:4])

	if requiredGas, ok := GasRequiredByMethod[methodID]; ok {
		return requiredGas + baseCost
	}

	// Can not happen, but return 0 if the method is not found.
	return 0
}

func (c *Contract) Delegate(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 3 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 3,
		})
	}

	msgServer := stakingkeeper.NewMsgServerImpl(&c.stakingKeeper)

	delegatorAddress, ok := args[0].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid argument, wanted a string, got: %T", args[0])
	}

	validatorAddress, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid argument, wanted a string, got: %T", args[1])
	}

	amount, ok := args[2].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid argument, wanted an int64, got %T", args[2])
	}

	res, err := msgServer.Delegate(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: sdk.AccAddress(delegatorAddress.Bytes()).String(),
		ValidatorAddress: validatorAddress,
		Amount: sdk.Coin{
			Denom:  c.stakingKeeper.BondDenom(ctx),
			Amount: math.NewIntFromBigInt(big.NewInt(amount)),
		},
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("delegate res", res)

	return method.Outputs.Pack(true)
}

func (c *Contract) Undelegate(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 3 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 3,
		})
	}

	msgServer := stakingkeeper.NewMsgServerImpl(&c.stakingKeeper)

	delegatorAddress, ok := args[0].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid argument, wanted a string, got: %T", args[0])
	}

	validatorAddress, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid argument, wanted a string, got: %T", args[1])
	}

	amount, ok := args[2].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid argument, wanted an int64, got %T", args[2])
	}

	res, err := msgServer.Undelegate(ctx, &stakingtypes.MsgUndelegate{
		DelegatorAddress: sdk.AccAddress(delegatorAddress.Bytes()).String(),
		ValidatorAddress: validatorAddress,
		Amount: sdk.Coin{
			Denom:  c.stakingKeeper.BondDenom(ctx),
			Amount: math.NewIntFromBigInt(big.NewInt(amount)),
		},
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("undelegate res", res)

	return method.Outputs.Pack(res.GetCompletionTime().UTC().Unix())
}

func (c *Contract) Redelegate(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 4 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 4,
		})
	}

	msgServer := stakingkeeper.NewMsgServerImpl(&c.stakingKeeper)

	delegatorAddress, ok := args[0].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid argument, wanted a string, got: %T", args[0])
	}

	validatorSrcAddress, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid argument, wanted a string, got: %T", args[1])
	}

	validatorDstAddress, ok := args[2].(string)
	if !ok {
		return nil, fmt.Errorf("invalid argument, wanted a string, got: %T", args[1])
	}

	amount, ok := args[3].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid argument, wanted an int64, got %T", args[2])
	}

	res, err := msgServer.BeginRedelegate(ctx, &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    sdk.AccAddress(delegatorAddress.Bytes()).String(),
		ValidatorSrcAddress: validatorSrcAddress,
		ValidatorDstAddress: validatorDstAddress,
		Amount: sdk.Coin{
			Denom:  c.stakingKeeper.BondDenom(ctx),
			Amount: math.NewIntFromBigInt(big.NewInt(amount)),
		},
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("redelegate res", res)

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
	case DelegateMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.Delegate(ctx, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case UndelegateMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.Undelegate(ctx, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil
	case RedelegateMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.Redelegate(ctx, method, args)
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
