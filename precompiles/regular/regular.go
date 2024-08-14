package regular

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	ptypes "github.com/zeta-chain/zetacore/precompiles/types"
	"github.com/zeta-chain/zetacore/testutil/contracts"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

const (
	RegularCallMethodName     = "regularCall"
	Bech32ToHexAddrMethodName = "bech32ToHexAddr"
	Bech32ifyMethodName       = "bech32ify"
)

var (
	ABI                 abi.ABI
	ContractAddress     = common.BytesToAddress([]byte{101})
	GasRequiredByMethod = map[[4]byte]uint64{}
	ExampleABI          *abi.ABI
)

func init() {
	ABI, GasRequiredByMethod = initABI()
	ExampleABI, _ = contracts.ExampleMetaData.GetAbi()
}

var RegularModuleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"name\":\"bech32ToHexAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"prefix\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"bech32ify\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"method\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"regularCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"result\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

func initABI() (abi abi.ABI, gasRequiredByMethod map[[4]byte]uint64) {
	gasRequiredByMethod = map[[4]byte]uint64{}
	if err := abi.UnmarshalJSON([]byte(RegularModuleMetaData.ABI)); err != nil {
		panic(err)
	}
	for methodName := range abi.Methods {
		var methodID [4]byte
		copy(methodID[:], abi.Methods[methodName].ID[:4])
		switch methodName {
		case RegularCallMethodName:
			gasRequiredByMethod[methodID] = 10
		case Bech32ToHexAddrMethodName:
			gasRequiredByMethod[methodID] = 0
		case Bech32ifyMethodName:
			gasRequiredByMethod[methodID] = 0
		default:
			gasRequiredByMethod[methodID] = 0
		}
	}
	return abi, gasRequiredByMethod
}

type Contract struct {
	ptypes.BaseContract

	FungibleKeeper fungiblekeeper.Keeper
	cdc            codec.Codec
	kvGasConfig    storetypes.GasConfig
}

// NewRegularContract creates the precompiled contract to manage native tokens
func NewRegularContract(
	fungibleKeeper fungiblekeeper.Keeper,
	cdc codec.Codec,
	kvGasConfig storetypes.GasConfig,
) *Contract {
	return &Contract{
		BaseContract:   ptypes.NewBaseContract(ContractAddress),
		FungibleKeeper: fungibleKeeper,
		cdc:            cdc,
		kvGasConfig:    kvGasConfig,
	}
}

func (rc *Contract) Address() common.Address {
	return ContractAddress
}

func (rc *Contract) Abi() abi.ABI {
	return ABI
}

// RequiredGas calculates the contract gas use
func (rc *Contract) RequiredGas(input []byte) uint64 {
	// base cost to prevent large input size
	baseCost := uint64(len(input)) * rc.kvGasConfig.WriteCostPerByte
	var methodID [4]byte
	copy(methodID[:], input[:4])
	requiredGas, ok := GasRequiredByMethod[methodID]
	if ok {
		return requiredGas + baseCost
	}
	return baseCost
}

func (rc *Contract) Bech32ToHexAddr(method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(ptypes.ErrInvalidNumberOfArgs, 1, len(args))
	}

	address, ok := args[0].(string)
	if !ok || address == "" {
		return nil, fmt.Errorf("invalid bech32 address: %v", args[0])
	}

	bech32Prefix := strings.SplitN(address, "1", 2)[0]
	if bech32Prefix == address {
		return nil, fmt.Errorf("invalid bech32 address: %s", address)
	}

	addressBz, err := sdk.GetFromBech32(address, bech32Prefix)
	if err != nil {
		return nil, err
	}

	if err := sdk.VerifyAddressFormat(addressBz); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(common.BytesToAddress(addressBz))
}
func (rc *Contract) Bech32ify(method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(ptypes.ErrInvalidNumberOfArgs, 2, len(args))
	}

	cfg := sdk.GetConfig()
	prefix, _ := args[0].(string)
	if strings.TrimSpace(prefix) == "" {
		return nil, fmt.Errorf(
			"invalid bech32 human readable prefix (HRP). Please provide a either an account, validator or consensus address prefix (eg: %s, %s, %s)",
			cfg.GetBech32AccountAddrPrefix(),
			cfg.GetBech32ValidatorAddrPrefix(),
			cfg.GetBech32ConsensusAddrPrefix(),
		)
	}

	address, ok := args[1].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid hex address")
	}

	// NOTE: safety check, should not happen given that the address is is 20 bytes.
	if err := sdk.VerifyAddressFormat(address.Bytes()); err != nil {
		return nil, err
	}

	bech32Str, err := sdk.Bech32ifyAddressBytes(prefix, address.Bytes())
	if err != nil {
		return nil, err
	}

	addressBz, err := sdk.GetFromBech32(bech32Str, "zeta")
	if err != nil {
		return nil, err
	}

	if err := sdk.VerifyAddressFormat(addressBz); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(bech32Str)
}

func (rc *Contract) RegularCall(ctx sdk.Context, method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(ptypes.ErrInvalidNumberOfArgs, 2, len(args))
	}
	callingMethod, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf(ptypes.ErrInvalidMethod, args[0])
	}
	callingContract, ok := args[1].(common.Address)
	if !ok {
		return nil, fmt.Errorf(ptypes.ErrInvalidAddr, args[1])
	}

	res, err := rc.FungibleKeeper.CallEVM(
		ctx,
		*ExampleABI,
		fungibletypes.ModuleAddressEVM,
		callingContract,
		big.NewInt(0),
		nil,
		true,
		false,
		callingMethod,
	)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(
		ptypes.BytesToBigInt(res.Ret),
	)
}

func (rc *Contract) Run(evm *vm.EVM, contract *vm.Contract, _ bool) ([]byte, error) {
	// parse input
	methodID := contract.Input[:4]
	method, err := ABI.MethodById(methodID)
	if err != nil {
		return nil, err
	}
	args, err := method.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("fail to unpack input arguments")
	}

	stateDB := evm.StateDB.(ptypes.ExtStateDB)

	switch method.Name {
	case RegularCallMethodName:
		var res []byte
		if execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = rc.RegularCall(ctx, method, args)
			return err
		}); execErr != nil {
			return nil, err
		} else {
			return res, nil
		}

	case Bech32ToHexAddrMethodName:
		return rc.Bech32ToHexAddr(method, args)
	case Bech32ifyMethodName:
		return rc.Bech32ify(method, args)
	// case OtherMethods:
	// ..
	default:
		return nil, errors.New("unknown method")
	}
}
