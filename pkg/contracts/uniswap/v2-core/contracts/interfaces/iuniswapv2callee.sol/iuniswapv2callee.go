// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package iuniswapv2callee

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

// IUniswapV2CalleeMetaData contains all meta data concerning the IUniswapV2Callee contract.
var IUniswapV2CalleeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"uniswapV2Call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IUniswapV2CalleeABI is the input ABI used to generate the binding from.
// Deprecated: Use IUniswapV2CalleeMetaData.ABI instead.
var IUniswapV2CalleeABI = IUniswapV2CalleeMetaData.ABI

// IUniswapV2Callee is an auto generated Go binding around an Ethereum contract.
type IUniswapV2Callee struct {
	IUniswapV2CalleeCaller     // Read-only binding to the contract
	IUniswapV2CalleeTransactor // Write-only binding to the contract
	IUniswapV2CalleeFilterer   // Log filterer for contract events
}

// IUniswapV2CalleeCaller is an auto generated read-only Go binding around an Ethereum contract.
type IUniswapV2CalleeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2CalleeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IUniswapV2CalleeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2CalleeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IUniswapV2CalleeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2CalleeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IUniswapV2CalleeSession struct {
	Contract     *IUniswapV2Callee // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IUniswapV2CalleeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IUniswapV2CalleeCallerSession struct {
	Contract *IUniswapV2CalleeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// IUniswapV2CalleeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IUniswapV2CalleeTransactorSession struct {
	Contract     *IUniswapV2CalleeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IUniswapV2CalleeRaw is an auto generated low-level Go binding around an Ethereum contract.
type IUniswapV2CalleeRaw struct {
	Contract *IUniswapV2Callee // Generic contract binding to access the raw methods on
}

// IUniswapV2CalleeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IUniswapV2CalleeCallerRaw struct {
	Contract *IUniswapV2CalleeCaller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV2CalleeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IUniswapV2CalleeTransactorRaw struct {
	Contract *IUniswapV2CalleeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV2Callee creates a new instance of IUniswapV2Callee, bound to a specific deployed contract.
func NewIUniswapV2Callee(address common.Address, backend bind.ContractBackend) (*IUniswapV2Callee, error) {
	contract, err := bindIUniswapV2Callee(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2Callee{IUniswapV2CalleeCaller: IUniswapV2CalleeCaller{contract: contract}, IUniswapV2CalleeTransactor: IUniswapV2CalleeTransactor{contract: contract}, IUniswapV2CalleeFilterer: IUniswapV2CalleeFilterer{contract: contract}}, nil
}

// NewIUniswapV2CalleeCaller creates a new read-only instance of IUniswapV2Callee, bound to a specific deployed contract.
func NewIUniswapV2CalleeCaller(address common.Address, caller bind.ContractCaller) (*IUniswapV2CalleeCaller, error) {
	contract, err := bindIUniswapV2Callee(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2CalleeCaller{contract: contract}, nil
}

// NewIUniswapV2CalleeTransactor creates a new write-only instance of IUniswapV2Callee, bound to a specific deployed contract.
func NewIUniswapV2CalleeTransactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV2CalleeTransactor, error) {
	contract, err := bindIUniswapV2Callee(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2CalleeTransactor{contract: contract}, nil
}

// NewIUniswapV2CalleeFilterer creates a new log filterer instance of IUniswapV2Callee, bound to a specific deployed contract.
func NewIUniswapV2CalleeFilterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV2CalleeFilterer, error) {
	contract, err := bindIUniswapV2Callee(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2CalleeFilterer{contract: contract}, nil
}

// bindIUniswapV2Callee binds a generic wrapper to an already deployed contract.
func bindIUniswapV2Callee(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUniswapV2CalleeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Callee *IUniswapV2CalleeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Callee.Contract.IUniswapV2CalleeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Callee *IUniswapV2CalleeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.IUniswapV2CalleeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Callee *IUniswapV2CalleeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.IUniswapV2CalleeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Callee *IUniswapV2CalleeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Callee.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Callee *IUniswapV2CalleeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Callee *IUniswapV2CalleeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.contract.Transact(opts, method, params...)
}

// UniswapV2Call is a paid mutator transaction binding the contract method 0x10d1e85c.
//
// Solidity: function uniswapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_IUniswapV2Callee *IUniswapV2CalleeTransactor) UniswapV2Call(opts *bind.TransactOpts, sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV2Callee.contract.Transact(opts, "uniswapV2Call", sender, amount0, amount1, data)
}

// UniswapV2Call is a paid mutator transaction binding the contract method 0x10d1e85c.
//
// Solidity: function uniswapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_IUniswapV2Callee *IUniswapV2CalleeSession) UniswapV2Call(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.UniswapV2Call(&_IUniswapV2Callee.TransactOpts, sender, amount0, amount1, data)
}

// UniswapV2Call is a paid mutator transaction binding the contract method 0x10d1e85c.
//
// Solidity: function uniswapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_IUniswapV2Callee *IUniswapV2CalleeTransactorSession) UniswapV2Call(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.UniswapV2Call(&_IUniswapV2Callee.TransactOpts, sender, amount0, amount1, data)
}
