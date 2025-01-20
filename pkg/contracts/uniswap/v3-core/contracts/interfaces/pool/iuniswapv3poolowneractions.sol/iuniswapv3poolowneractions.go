// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package iuniswapv3poolowneractions

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

// IUniswapV3PoolOwnerActionsMetaData contains all meta data concerning the IUniswapV3PoolOwnerActions contract.
var IUniswapV3PoolOwnerActionsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"amount0Requested\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amount1Requested\",\"type\":\"uint128\"}],\"name\":\"collectProtocol\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"amount0\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amount1\",\"type\":\"uint128\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"feeProtocol0\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"feeProtocol1\",\"type\":\"uint8\"}],\"name\":\"setFeeProtocol\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IUniswapV3PoolOwnerActionsABI is the input ABI used to generate the binding from.
// Deprecated: Use IUniswapV3PoolOwnerActionsMetaData.ABI instead.
var IUniswapV3PoolOwnerActionsABI = IUniswapV3PoolOwnerActionsMetaData.ABI

// IUniswapV3PoolOwnerActions is an auto generated Go binding around an Ethereum contract.
type IUniswapV3PoolOwnerActions struct {
	IUniswapV3PoolOwnerActionsCaller     // Read-only binding to the contract
	IUniswapV3PoolOwnerActionsTransactor // Write-only binding to the contract
	IUniswapV3PoolOwnerActionsFilterer   // Log filterer for contract events
}

// IUniswapV3PoolOwnerActionsCaller is an auto generated read-only Go binding around an Ethereum contract.
type IUniswapV3PoolOwnerActionsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolOwnerActionsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IUniswapV3PoolOwnerActionsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolOwnerActionsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IUniswapV3PoolOwnerActionsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolOwnerActionsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IUniswapV3PoolOwnerActionsSession struct {
	Contract     *IUniswapV3PoolOwnerActions // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IUniswapV3PoolOwnerActionsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IUniswapV3PoolOwnerActionsCallerSession struct {
	Contract *IUniswapV3PoolOwnerActionsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// IUniswapV3PoolOwnerActionsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IUniswapV3PoolOwnerActionsTransactorSession struct {
	Contract     *IUniswapV3PoolOwnerActionsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// IUniswapV3PoolOwnerActionsRaw is an auto generated low-level Go binding around an Ethereum contract.
type IUniswapV3PoolOwnerActionsRaw struct {
	Contract *IUniswapV3PoolOwnerActions // Generic contract binding to access the raw methods on
}

// IUniswapV3PoolOwnerActionsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IUniswapV3PoolOwnerActionsCallerRaw struct {
	Contract *IUniswapV3PoolOwnerActionsCaller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV3PoolOwnerActionsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IUniswapV3PoolOwnerActionsTransactorRaw struct {
	Contract *IUniswapV3PoolOwnerActionsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV3PoolOwnerActions creates a new instance of IUniswapV3PoolOwnerActions, bound to a specific deployed contract.
func NewIUniswapV3PoolOwnerActions(address common.Address, backend bind.ContractBackend) (*IUniswapV3PoolOwnerActions, error) {
	contract, err := bindIUniswapV3PoolOwnerActions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolOwnerActions{IUniswapV3PoolOwnerActionsCaller: IUniswapV3PoolOwnerActionsCaller{contract: contract}, IUniswapV3PoolOwnerActionsTransactor: IUniswapV3PoolOwnerActionsTransactor{contract: contract}, IUniswapV3PoolOwnerActionsFilterer: IUniswapV3PoolOwnerActionsFilterer{contract: contract}}, nil
}

// NewIUniswapV3PoolOwnerActionsCaller creates a new read-only instance of IUniswapV3PoolOwnerActions, bound to a specific deployed contract.
func NewIUniswapV3PoolOwnerActionsCaller(address common.Address, caller bind.ContractCaller) (*IUniswapV3PoolOwnerActionsCaller, error) {
	contract, err := bindIUniswapV3PoolOwnerActions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolOwnerActionsCaller{contract: contract}, nil
}

// NewIUniswapV3PoolOwnerActionsTransactor creates a new write-only instance of IUniswapV3PoolOwnerActions, bound to a specific deployed contract.
func NewIUniswapV3PoolOwnerActionsTransactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV3PoolOwnerActionsTransactor, error) {
	contract, err := bindIUniswapV3PoolOwnerActions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolOwnerActionsTransactor{contract: contract}, nil
}

// NewIUniswapV3PoolOwnerActionsFilterer creates a new log filterer instance of IUniswapV3PoolOwnerActions, bound to a specific deployed contract.
func NewIUniswapV3PoolOwnerActionsFilterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV3PoolOwnerActionsFilterer, error) {
	contract, err := bindIUniswapV3PoolOwnerActions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolOwnerActionsFilterer{contract: contract}, nil
}

// bindIUniswapV3PoolOwnerActions binds a generic wrapper to an already deployed contract.
func bindIUniswapV3PoolOwnerActions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUniswapV3PoolOwnerActionsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV3PoolOwnerActions.Contract.IUniswapV3PoolOwnerActionsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV3PoolOwnerActions.Contract.IUniswapV3PoolOwnerActionsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV3PoolOwnerActions.Contract.IUniswapV3PoolOwnerActionsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV3PoolOwnerActions.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV3PoolOwnerActions.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV3PoolOwnerActions.Contract.contract.Transact(opts, method, params...)
}

// CollectProtocol is a paid mutator transaction binding the contract method 0x85b66729.
//
// Solidity: function collectProtocol(address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsTransactor) CollectProtocol(opts *bind.TransactOpts, recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolOwnerActions.contract.Transact(opts, "collectProtocol", recipient, amount0Requested, amount1Requested)
}

// CollectProtocol is a paid mutator transaction binding the contract method 0x85b66729.
//
// Solidity: function collectProtocol(address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsSession) CollectProtocol(recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolOwnerActions.Contract.CollectProtocol(&_IUniswapV3PoolOwnerActions.TransactOpts, recipient, amount0Requested, amount1Requested)
}

// CollectProtocol is a paid mutator transaction binding the contract method 0x85b66729.
//
// Solidity: function collectProtocol(address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsTransactorSession) CollectProtocol(recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _IUniswapV3PoolOwnerActions.Contract.CollectProtocol(&_IUniswapV3PoolOwnerActions.TransactOpts, recipient, amount0Requested, amount1Requested)
}

// SetFeeProtocol is a paid mutator transaction binding the contract method 0x8206a4d1.
//
// Solidity: function setFeeProtocol(uint8 feeProtocol0, uint8 feeProtocol1) returns()
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsTransactor) SetFeeProtocol(opts *bind.TransactOpts, feeProtocol0 uint8, feeProtocol1 uint8) (*types.Transaction, error) {
	return _IUniswapV3PoolOwnerActions.contract.Transact(opts, "setFeeProtocol", feeProtocol0, feeProtocol1)
}

// SetFeeProtocol is a paid mutator transaction binding the contract method 0x8206a4d1.
//
// Solidity: function setFeeProtocol(uint8 feeProtocol0, uint8 feeProtocol1) returns()
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsSession) SetFeeProtocol(feeProtocol0 uint8, feeProtocol1 uint8) (*types.Transaction, error) {
	return _IUniswapV3PoolOwnerActions.Contract.SetFeeProtocol(&_IUniswapV3PoolOwnerActions.TransactOpts, feeProtocol0, feeProtocol1)
}

// SetFeeProtocol is a paid mutator transaction binding the contract method 0x8206a4d1.
//
// Solidity: function setFeeProtocol(uint8 feeProtocol0, uint8 feeProtocol1) returns()
func (_IUniswapV3PoolOwnerActions *IUniswapV3PoolOwnerActionsTransactorSession) SetFeeProtocol(feeProtocol0 uint8, feeProtocol1 uint8) (*types.Transaction, error) {
	return _IUniswapV3PoolOwnerActions.Contract.SetFeeProtocol(&_IUniswapV3PoolOwnerActions.TransactOpts, feeProtocol0, feeProtocol1)
}
