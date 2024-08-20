package prototype

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	ptypes "github.com/zeta-chain/zetacore/precompiles/types"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
)

const (
	Bech32ToHexAddrMethodName      = "bech32ToHexAddr"
	Bech32ifyMethodName            = "bech32ify"
	GetGasStabilityPoolBalanceName = "getGasStabilityPoolBalance"
)

var (
	ABI                 abi.ABI
	ContractAddress     = common.HexToAddress("0x0000000000000000000000000000000000000065")
	GasRequiredByMethod = map[[4]byte]uint64{}

	//go:embed IPrototype.abi
	prototypeABI string
)

func init() {
	initABI()
}

func initABI() {
	if prototypeABI == "" {
		panic("missing prototype ABI")
	}

	var IPrototypeMetaData = &bind.MetaData{
		ABI: prototypeABI,
	}

	if err := ABI.UnmarshalJSON([]byte(IPrototypeMetaData.ABI)); err != nil {
		panic(err)
	}

	GasRequiredByMethod = map[[4]byte]uint64{}
	for methodName := range ABI.Methods {
		var methodID [4]byte
		copy(methodID[:], ABI.Methods[methodName].ID[:4])
		switch methodName {
		case Bech32ToHexAddrMethodName:
			GasRequiredByMethod[methodID] = 10000
		case Bech32ifyMethodName:
			GasRequiredByMethod[methodID] = 10000
		case GetGasStabilityPoolBalanceName:
			GasRequiredByMethod[methodID] = 10000
		}
	}
}

type Contract struct {
	ptypes.BaseContract

	fungibleKeeper fungiblekeeper.Keeper
	cdc            codec.Codec
	kvGasConfig    storetypes.GasConfig
}

func NewIPrototypeContract(
	fungibleKeeper *fungiblekeeper.Keeper,
	cdc codec.Codec,
	kvGasConfig storetypes.GasConfig,
) *Contract {
	return &Contract{
		BaseContract:   ptypes.NewBaseContract(ContractAddress),
		fungibleKeeper: *fungibleKeeper,
		cdc:            cdc,
		kvGasConfig:    kvGasConfig,
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

func (c *Contract) Bech32ToHexAddr(method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, &ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 1,
		}
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

func (c *Contract) Bech32ify(method *abi.Method, args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, &ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		}
	}

	cfg := sdk.GetConfig()
	prefix, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid bech32 human readable prefix (HRP): %v", args[0])
	}

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

	// NOTE: safety check, should not happen given that the address is 20 bytes.
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

func (c *Contract) GetGasStabilityPoolBalance(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 1 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 1,
		})
	}

	// Unwrap arguments. The chainID is the first and unique argument.
	chainID, ok := args[0].(int64)
	if !ok {
		return nil, ptypes.ErrInvalidArgument{
			Got: chainID,
		}
	}

	balance, err := c.fungibleKeeper.GetGasStabilityPoolBalance(ctx, chainID)
	if err != nil {
		return nil, fmt.Errorf("error calling fungible keeper: %s", err.Error())
	}

	return method.Outputs.Pack(balance)
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
	case GetGasStabilityPoolBalanceName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.GetGasStabilityPoolBalance(ctx, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil

	case Bech32ToHexAddrMethodName:
		return c.Bech32ToHexAddr(method, args)
	case Bech32ifyMethodName:
		return c.Bech32ify(method, args)
	default:
		return nil, ptypes.ErrInvalidMethod{
			Method: method.Name,
		}
	}
}
