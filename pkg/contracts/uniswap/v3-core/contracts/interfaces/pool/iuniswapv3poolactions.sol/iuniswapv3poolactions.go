// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package iuniswapv3poolactions

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

// IUniswapV3PoolActionsMetaData contains all meta data concerning the IUniswapV3PoolActions contract.
var IUniswapV3PoolActionsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"uint128\",\"name\":\"amount\",\"type\":\"uint128\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"uint128\",\"name\":\"amount0Requested\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amount1Requested\",\"type\":\"uint128\"}],\"name\":\"collect\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"amount0\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amount1\",\"type\":\"uint128\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"flash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"observationCardinalityNext\",\"type\":\"uint16\"}],\"name\":\"increaseObservationCardinalityNext\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"uint128\",\"name\":\"amount\",\"type\":\"uint128\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"zeroForOne\",\"type\":\"bool\"},{\"internalType\":\"int256\",\"name\":\"amountSpecified\",\"type\":\"int256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"amount0\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IUniswapV3PoolActionsABI is the input ABI used to generate the binding from.
// Deprecated: Use IUniswapV3PoolActionsMetaData.ABI instead.
var IUniswapV3PoolActionsABI = IUniswapV3PoolActionsMetaData.ABI

// IUniswapV3PoolActions is an auto generated Go binding around an Ethereum contract.
type IUniswapV3PoolActions struct {
	IUniswapV3PoolActionsCaller     // Read-only binding to the contract
	IUniswapV3PoolActionsTransactor // Write-only binding to the contract
	IUniswapV3PoolActionsFilterer   // Log filterer for contract events
}

// IUniswapV3PoolActionsCaller is an auto generated read-only Go binding around an Ethereum contract.
type IUniswapV3PoolActionsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolActionsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IUniswapV3PoolActionsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolActionsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IUniswapV3PoolActionsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolActionsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IUniswapV3PoolActionsSession struct {
	Contract     *IUniswapV3PoolActions // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// IUniswapV3PoolActionsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IUniswapV3PoolActionsCallerSession struct {
	Contract *IUniswapV3PoolActionsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// IUniswapV3PoolActionsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IUniswapV3PoolActionsTransactorSession struct {
	Contract     *IUniswapV3PoolActionsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// IUniswapV3PoolActionsRaw is an auto generated low-level Go binding around an Ethereum contract.
type IUniswapV3PoolActionsRaw struct {
	Contract *IUniswapV3PoolActions // Generic contract binding to access the raw methods on
}

// IUniswapV3PoolActionsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IUniswapV3PoolActionsCallerRaw struct {
	Contract *IUniswapV3PoolActionsCaller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV3PoolActionsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IUniswapV3PoolActionsTransactorRaw struct {
	Contract *IUniswapV3PoolActionsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV3PoolActions creates a new instance of IUniswapV3PoolActions, bound to a specific deployed contract.
func NewIUniswapV3PoolActions(address common.Address, backend bind.ContractBackend) (*IUniswapV3PoolActions, error) {
	contract, err := bindIUniswapV3PoolActions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolActions{IUniswapV3PoolActionsCaller: IUniswapV3PoolActionsCaller{contract: contract}, IUniswapV3PoolActionsTransactor: IUniswapV3PoolActionsTransactor{contract: contract}, IUniswapV3PoolActionsFilterer: IUniswapV3PoolActionsFilterer{contract: contract}}, nil
}

// NewIUniswapV3PoolActionsCaller creates a new read-only instance of IUniswapV3PoolActions, bound to a specific deployed contract.
func NewIUniswapV3PoolActionsCaller(address common.Address, caller bind.ContractCaller) (*IUniswapV3PoolActionsCaller, error) {
	contract, err := bindIUniswapV3PoolActions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolActionsCaller{contract: contract}, nil
}

// NewIUniswapV3PoolActionsTransactor creates a new write-only instance of IUniswapV3PoolActions, bound to a specific deployed contract.
func NewIUniswapV3PoolActionsTransactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV3PoolActionsTransactor, error) {
	contract, err := bindIUniswapV3PoolActions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolActionsTransactor{contract: contract}, nil
}

// NewIUniswapV3PoolActionsFilterer creates a new log filterer instance of IUniswapV3PoolActions, bound to a specific deployed contract.
func NewIUniswapV3PoolActionsFilterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV3PoolActionsFilterer, error) {
	contract, err := bindIUniswapV3PoolActions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolActionsFilterer{contract: contract}, nil
}

// bindIUniswapV3PoolActions binds a generic wrapper to an already deployed contract.
func bindIUniswapV3PoolActions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUniswapV3PoolActionsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV3PoolActions.Contract.IUniswapV3PoolActionsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.IUniswapV3PoolActionsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.IUniswapV3PoolActionsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV3PoolActions.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.contract.Transact(opts, method, params...)
}

// Burn is a paid mutator transaction binding the contract method 0xa34123a7.
//
// Solidity: function burn(int24 tickLower, int24 tickUpper, uint128 amount) returns(uint256 amount0, uint256 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactor) Burn(opts *bind.TransactOpts, tickLower *big.Int, tickUpper *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.contract.Transact(opts, "burn", tickLower, tickUpper, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xa34123a7.
//
// Solidity: function burn(int24 tickLower, int24 tickUpper, uint128 amount) returns(uint256 amount0, uint256 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsSession) Burn(tickLower *big.Int, tickUpper *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Burn(&_IUniswapV3PoolActions.TransactOpts, tickLower, tickUpper, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xa34123a7.
//
// Solidity: function burn(int24 tickLower, int24 tickUpper, uint128 amount) returns(uint256 amount0, uint256 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactorSession) Burn(tickLower *big.Int, tickUpper *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Burn(&_IUniswapV3PoolActions.TransactOpts, tickLower, tickUpper, amount)
}

// Collect is a paid mutator transaction binding the contract method 0x4f1eb3d8.
//
// Solidity: function collect(address recipient, int24 tickLower, int24 tickUpper, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactor) Collect(opts *bind.TransactOpts, recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.contract.Transact(opts, "collect", recipient, tickLower, tickUpper, amount0Requested, amount1Requested)
}

// Collect is a paid mutator transaction binding the contract method 0x4f1eb3d8.
//
// Solidity: function collect(address recipient, int24 tickLower, int24 tickUpper, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsSession) Collect(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Collect(&_IUniswapV3PoolActions.TransactOpts, recipient, tickLower, tickUpper, amount0Requested, amount1Requested)
}

// Collect is a paid mutator transaction binding the contract method 0x4f1eb3d8.
//
// Solidity: function collect(address recipient, int24 tickLower, int24 tickUpper, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactorSession) Collect(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Collect(&_IUniswapV3PoolActions.TransactOpts, recipient, tickLower, tickUpper, amount0Requested, amount1Requested)
}

// Flash is a paid mutator transaction binding the contract method 0x490e6cbc.
//
// Solidity: function flash(address recipient, uint256 amount0, uint256 amount1, bytes data) returns()
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactor) Flash(opts *bind.TransactOpts, recipient common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.contract.Transact(opts, "flash", recipient, amount0, amount1, data)
}

// Flash is a paid mutator transaction binding the contract method 0x490e6cbc.
//
// Solidity: function flash(address recipient, uint256 amount0, uint256 amount1, bytes data) returns()
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsSession) Flash(recipient common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Flash(&_IUniswapV3PoolActions.TransactOpts, recipient, amount0, amount1, data)
}

// Flash is a paid mutator transaction binding the contract method 0x490e6cbc.
//
// Solidity: function flash(address recipient, uint256 amount0, uint256 amount1, bytes data) returns()
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactorSession) Flash(recipient common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Flash(&_IUniswapV3PoolActions.TransactOpts, recipient, amount0, amount1, data)
}

// IncreaseObservationCardinalityNext is a paid mutator transaction binding the contract method 0x32148f67.
//
// Solidity: function increaseObservationCardinalityNext(uint16 observationCardinalityNext) returns()
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactor) IncreaseObservationCardinalityNext(opts *bind.TransactOpts, observationCardinalityNext uint16) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.contract.Transact(opts, "increaseObservationCardinalityNext", observationCardinalityNext)
}

// IncreaseObservationCardinalityNext is a paid mutator transaction binding the contract method 0x32148f67.
//
// Solidity: function increaseObservationCardinalityNext(uint16 observationCardinalityNext) returns()
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsSession) IncreaseObservationCardinalityNext(observationCardinalityNext uint16) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.IncreaseObservationCardinalityNext(&_IUniswapV3PoolActions.TransactOpts, observationCardinalityNext)
}

// IncreaseObservationCardinalityNext is a paid mutator transaction binding the contract method 0x32148f67.
//
// Solidity: function increaseObservationCardinalityNext(uint16 observationCardinalityNext) returns()
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactorSession) IncreaseObservationCardinalityNext(observationCardinalityNext uint16) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.IncreaseObservationCardinalityNext(&_IUniswapV3PoolActions.TransactOpts, observationCardinalityNext)
}

// Initialize is a paid mutator transaction binding the contract method 0xf637731d.
//
// Solidity: function initialize(uint160 sqrtPriceX96) returns()
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactor) Initialize(opts *bind.TransactOpts, sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.contract.Transact(opts, "initialize", sqrtPriceX96)
}

// Initialize is a paid mutator transaction binding the contract method 0xf637731d.
//
// Solidity: function initialize(uint160 sqrtPriceX96) returns()
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsSession) Initialize(sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Initialize(&_IUniswapV3PoolActions.TransactOpts, sqrtPriceX96)
}

// Initialize is a paid mutator transaction binding the contract method 0xf637731d.
//
// Solidity: function initialize(uint160 sqrtPriceX96) returns()
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactorSession) Initialize(sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Initialize(&_IUniswapV3PoolActions.TransactOpts, sqrtPriceX96)
}

// Mint is a paid mutator transaction binding the contract method 0x3c8a7d8d.
//
// Solidity: function mint(address recipient, int24 tickLower, int24 tickUpper, uint128 amount, bytes data) returns(uint256 amount0, uint256 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactor) Mint(opts *bind.TransactOpts, recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.contract.Transact(opts, "mint", recipient, tickLower, tickUpper, amount, data)
}

// Mint is a paid mutator transaction binding the contract method 0x3c8a7d8d.
//
// Solidity: function mint(address recipient, int24 tickLower, int24 tickUpper, uint128 amount, bytes data) returns(uint256 amount0, uint256 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsSession) Mint(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Mint(&_IUniswapV3PoolActions.TransactOpts, recipient, tickLower, tickUpper, amount, data)
}

// Mint is a paid mutator transaction binding the contract method 0x3c8a7d8d.
//
// Solidity: function mint(address recipient, int24 tickLower, int24 tickUpper, uint128 amount, bytes data) returns(uint256 amount0, uint256 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactorSession) Mint(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Mint(&_IUniswapV3PoolActions.TransactOpts, recipient, tickLower, tickUpper, amount, data)
}

// Swap is a paid mutator transaction binding the contract method 0x128acb08.
//
// Solidity: function swap(address recipient, bool zeroForOne, int256 amountSpecified, uint160 sqrtPriceLimitX96, bytes data) returns(int256 amount0, int256 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactor) Swap(opts *bind.TransactOpts, recipient common.Address, zeroForOne bool, amountSpecified *big.Int, sqrtPriceLimitX96 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.contract.Transact(opts, "swap", recipient, zeroForOne, amountSpecified, sqrtPriceLimitX96, data)
}

// Swap is a paid mutator transaction binding the contract method 0x128acb08.
//
// Solidity: function swap(address recipient, bool zeroForOne, int256 amountSpecified, uint160 sqrtPriceLimitX96, bytes data) returns(int256 amount0, int256 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsSession) Swap(recipient common.Address, zeroForOne bool, amountSpecified *big.Int, sqrtPriceLimitX96 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Swap(&_IUniswapV3PoolActions.TransactOpts, recipient, zeroForOne, amountSpecified, sqrtPriceLimitX96, data)
}

// Swap is a paid mutator transaction binding the contract method 0x128acb08.
//
// Solidity: function swap(address recipient, bool zeroForOne, int256 amountSpecified, uint160 sqrtPriceLimitX96, bytes data) returns(int256 amount0, int256 amount1)
func (_IUniswapV3PoolActions *IUniswapV3PoolActionsTransactorSession) Swap(recipient common.Address, zeroForOne bool, amountSpecified *big.Int, sqrtPriceLimitX96 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV3PoolActions.Contract.Swap(&_IUniswapV3PoolActions.TransactOpts, recipient, zeroForOne, amountSpecified, sqrtPriceLimitX96, data)
}
