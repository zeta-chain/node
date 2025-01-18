// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package iuniswapv3poolderivedstate

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IUniswapV3PoolDerivedStateMetaData contains all meta data concerning the IUniswapV3PoolDerivedState contract.
var IUniswapV3PoolDerivedStateMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"secondsAgos\",\"type\":\"uint32[]\"}],\"name\":\"observe\",\"outputs\":[{\"internalType\":\"int56[]\",\"name\":\"tickCumulatives\",\"type\":\"int56[]\"},{\"internalType\":\"uint160[]\",\"name\":\"secondsPerLiquidityCumulativeX128s\",\"type\":\"uint160[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"}],\"name\":\"snapshotCumulativesInside\",\"outputs\":[{\"internalType\":\"int56\",\"name\":\"tickCumulativeInside\",\"type\":\"int56\"},{\"internalType\":\"uint160\",\"name\":\"secondsPerLiquidityInsideX128\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"secondsInside\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IUniswapV3PoolDerivedStateABI is the input ABI used to generate the binding from.
// Deprecated: Use IUniswapV3PoolDerivedStateMetaData.ABI instead.
var IUniswapV3PoolDerivedStateABI = IUniswapV3PoolDerivedStateMetaData.ABI

// IUniswapV3PoolDerivedState is an auto generated Go binding around an Ethereum contract.
type IUniswapV3PoolDerivedState struct {
	IUniswapV3PoolDerivedStateCaller     // Read-only binding to the contract
	IUniswapV3PoolDerivedStateTransactor // Write-only binding to the contract
	IUniswapV3PoolDerivedStateFilterer   // Log filterer for contract events
}

// IUniswapV3PoolDerivedStateCaller is an auto generated read-only Go binding around an Ethereum contract.
type IUniswapV3PoolDerivedStateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolDerivedStateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IUniswapV3PoolDerivedStateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolDerivedStateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IUniswapV3PoolDerivedStateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolDerivedStateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IUniswapV3PoolDerivedStateSession struct {
	Contract     *IUniswapV3PoolDerivedState // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IUniswapV3PoolDerivedStateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IUniswapV3PoolDerivedStateCallerSession struct {
	Contract *IUniswapV3PoolDerivedStateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// IUniswapV3PoolDerivedStateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IUniswapV3PoolDerivedStateTransactorSession struct {
	Contract     *IUniswapV3PoolDerivedStateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// IUniswapV3PoolDerivedStateRaw is an auto generated low-level Go binding around an Ethereum contract.
type IUniswapV3PoolDerivedStateRaw struct {
	Contract *IUniswapV3PoolDerivedState // Generic contract binding to access the raw methods on
}

// IUniswapV3PoolDerivedStateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IUniswapV3PoolDerivedStateCallerRaw struct {
	Contract *IUniswapV3PoolDerivedStateCaller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV3PoolDerivedStateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IUniswapV3PoolDerivedStateTransactorRaw struct {
	Contract *IUniswapV3PoolDerivedStateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV3PoolDerivedState creates a new instance of IUniswapV3PoolDerivedState, bound to a specific deployed contract.
func NewIUniswapV3PoolDerivedState(address common.Address, backend bind.ContractBackend) (*IUniswapV3PoolDerivedState, error) {
	contract, err := bindIUniswapV3PoolDerivedState(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolDerivedState{IUniswapV3PoolDerivedStateCaller: IUniswapV3PoolDerivedStateCaller{contract: contract}, IUniswapV3PoolDerivedStateTransactor: IUniswapV3PoolDerivedStateTransactor{contract: contract}, IUniswapV3PoolDerivedStateFilterer: IUniswapV3PoolDerivedStateFilterer{contract: contract}}, nil
}

// NewIUniswapV3PoolDerivedStateCaller creates a new read-only instance of IUniswapV3PoolDerivedState, bound to a specific deployed contract.
func NewIUniswapV3PoolDerivedStateCaller(address common.Address, caller bind.ContractCaller) (*IUniswapV3PoolDerivedStateCaller, error) {
	contract, err := bindIUniswapV3PoolDerivedState(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolDerivedStateCaller{contract: contract}, nil
}

// NewIUniswapV3PoolDerivedStateTransactor creates a new write-only instance of IUniswapV3PoolDerivedState, bound to a specific deployed contract.
func NewIUniswapV3PoolDerivedStateTransactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV3PoolDerivedStateTransactor, error) {
	contract, err := bindIUniswapV3PoolDerivedState(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolDerivedStateTransactor{contract: contract}, nil
}

// NewIUniswapV3PoolDerivedStateFilterer creates a new log filterer instance of IUniswapV3PoolDerivedState, bound to a specific deployed contract.
func NewIUniswapV3PoolDerivedStateFilterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV3PoolDerivedStateFilterer, error) {
	contract, err := bindIUniswapV3PoolDerivedState(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolDerivedStateFilterer{contract: contract}, nil
}

// bindIUniswapV3PoolDerivedState binds a generic wrapper to an already deployed contract.
func bindIUniswapV3PoolDerivedState(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUniswapV3PoolDerivedStateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV3PoolDerivedState.Contract.IUniswapV3PoolDerivedStateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV3PoolDerivedState.Contract.IUniswapV3PoolDerivedStateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV3PoolDerivedState.Contract.IUniswapV3PoolDerivedStateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV3PoolDerivedState.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV3PoolDerivedState.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV3PoolDerivedState.Contract.contract.Transact(opts, method, params...)
}

// Observe is a free data retrieval call binding the contract method 0x883bdbfd.
//
// Solidity: function observe(uint32[] secondsAgos) view returns(int56[] tickCumulatives, uint160[] secondsPerLiquidityCumulativeX128s)
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateCaller) Observe(opts *bind.CallOpts, secondsAgos []uint32) (struct {
	TickCumulatives                    []*big.Int
	SecondsPerLiquidityCumulativeX128s []*big.Int
}, error) {
	var out []interface{}
	err := _IUniswapV3PoolDerivedState.contract.Call(opts, &out, "observe", secondsAgos)

	outstruct := new(struct {
		TickCumulatives                    []*big.Int
		SecondsPerLiquidityCumulativeX128s []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TickCumulatives = *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	outstruct.SecondsPerLiquidityCumulativeX128s = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Observe is a free data retrieval call binding the contract method 0x883bdbfd.
//
// Solidity: function observe(uint32[] secondsAgos) view returns(int56[] tickCumulatives, uint160[] secondsPerLiquidityCumulativeX128s)
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateSession) Observe(secondsAgos []uint32) (struct {
	TickCumulatives                    []*big.Int
	SecondsPerLiquidityCumulativeX128s []*big.Int
}, error) {
	return _IUniswapV3PoolDerivedState.Contract.Observe(&_IUniswapV3PoolDerivedState.CallOpts, secondsAgos)
}

// Observe is a free data retrieval call binding the contract method 0x883bdbfd.
//
// Solidity: function observe(uint32[] secondsAgos) view returns(int56[] tickCumulatives, uint160[] secondsPerLiquidityCumulativeX128s)
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateCallerSession) Observe(secondsAgos []uint32) (struct {
	TickCumulatives                    []*big.Int
	SecondsPerLiquidityCumulativeX128s []*big.Int
}, error) {
	return _IUniswapV3PoolDerivedState.Contract.Observe(&_IUniswapV3PoolDerivedState.CallOpts, secondsAgos)
}

// SnapshotCumulativesInside is a free data retrieval call binding the contract method 0xa38807f2.
//
// Solidity: function snapshotCumulativesInside(int24 tickLower, int24 tickUpper) view returns(int56 tickCumulativeInside, uint160 secondsPerLiquidityInsideX128, uint32 secondsInside)
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateCaller) SnapshotCumulativesInside(opts *bind.CallOpts, tickLower *big.Int, tickUpper *big.Int) (struct {
	TickCumulativeInside          *big.Int
	SecondsPerLiquidityInsideX128 *big.Int
	SecondsInside                 uint32
}, error) {
	var out []interface{}
	err := _IUniswapV3PoolDerivedState.contract.Call(opts, &out, "snapshotCumulativesInside", tickLower, tickUpper)

	outstruct := new(struct {
		TickCumulativeInside          *big.Int
		SecondsPerLiquidityInsideX128 *big.Int
		SecondsInside                 uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TickCumulativeInside = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SecondsPerLiquidityInsideX128 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.SecondsInside = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// SnapshotCumulativesInside is a free data retrieval call binding the contract method 0xa38807f2.
//
// Solidity: function snapshotCumulativesInside(int24 tickLower, int24 tickUpper) view returns(int56 tickCumulativeInside, uint160 secondsPerLiquidityInsideX128, uint32 secondsInside)
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateSession) SnapshotCumulativesInside(tickLower *big.Int, tickUpper *big.Int) (struct {
	TickCumulativeInside          *big.Int
	SecondsPerLiquidityInsideX128 *big.Int
	SecondsInside                 uint32
}, error) {
	return _IUniswapV3PoolDerivedState.Contract.SnapshotCumulativesInside(&_IUniswapV3PoolDerivedState.CallOpts, tickLower, tickUpper)
}

// SnapshotCumulativesInside is a free data retrieval call binding the contract method 0xa38807f2.
//
// Solidity: function snapshotCumulativesInside(int24 tickLower, int24 tickUpper) view returns(int56 tickCumulativeInside, uint160 secondsPerLiquidityInsideX128, uint32 secondsInside)
func (_IUniswapV3PoolDerivedState *IUniswapV3PoolDerivedStateCallerSession) SnapshotCumulativesInside(tickLower *big.Int, tickUpper *big.Int) (struct {
	TickCumulativeInside          *big.Int
	SecondsPerLiquidityInsideX128 *big.Int
	SecondsInside                 uint32
}, error) {
	return _IUniswapV3PoolDerivedState.Contract.SnapshotCumulativesInside(&_IUniswapV3PoolDerivedState.CallOpts, tickLower, tickUpper)
}
