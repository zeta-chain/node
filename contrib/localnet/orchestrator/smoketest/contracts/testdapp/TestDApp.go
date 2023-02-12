// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testdapp

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

// TestDAppMetaData contains all meta data concerning the TestDApp contract.
var TestDAppMetaData = &bind.MetaData{
	ABI: "null",
}

// TestDAppABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDAppMetaData.ABI instead.
var TestDAppABI = TestDAppMetaData.ABI

// TestDApp is an auto generated Go binding around an Ethereum contract.
type TestDApp struct {
	TestDAppCaller     // Read-only binding to the contract
	TestDAppTransactor // Write-only binding to the contract
	TestDAppFilterer   // Log filterer for contract events
}

// TestDAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestDAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestDAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestDAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestDAppSession struct {
	Contract     *TestDApp         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestDAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestDAppCallerSession struct {
	Contract *TestDAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TestDAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestDAppTransactorSession struct {
	Contract     *TestDAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TestDAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestDAppRaw struct {
	Contract *TestDApp // Generic contract binding to access the raw methods on
}

// TestDAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestDAppCallerRaw struct {
	Contract *TestDAppCaller // Generic read-only contract binding to access the raw methods on
}

// TestDAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestDAppTransactorRaw struct {
	Contract *TestDAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestDApp creates a new instance of TestDApp, bound to a specific deployed contract.
func NewTestDApp(address common.Address, backend bind.ContractBackend) (*TestDApp, error) {
	contract, err := bindTestDApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestDApp{TestDAppCaller: TestDAppCaller{contract: contract}, TestDAppTransactor: TestDAppTransactor{contract: contract}, TestDAppFilterer: TestDAppFilterer{contract: contract}}, nil
}

// NewTestDAppCaller creates a new read-only instance of TestDApp, bound to a specific deployed contract.
func NewTestDAppCaller(address common.Address, caller bind.ContractCaller) (*TestDAppCaller, error) {
	contract, err := bindTestDApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppCaller{contract: contract}, nil
}

// NewTestDAppTransactor creates a new write-only instance of TestDApp, bound to a specific deployed contract.
func NewTestDAppTransactor(address common.Address, transactor bind.ContractTransactor) (*TestDAppTransactor, error) {
	contract, err := bindTestDApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppTransactor{contract: contract}, nil
}

// NewTestDAppFilterer creates a new log filterer instance of TestDApp, bound to a specific deployed contract.
func NewTestDAppFilterer(address common.Address, filterer bind.ContractFilterer) (*TestDAppFilterer, error) {
	contract, err := bindTestDApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestDAppFilterer{contract: contract}, nil
}

// bindTestDApp binds a generic wrapper to an already deployed contract.
func bindTestDApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestDAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDApp *TestDAppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDApp.Contract.TestDAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDApp *TestDAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDApp.Contract.TestDAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDApp *TestDAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDApp.Contract.TestDAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDApp *TestDAppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDApp *TestDAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDApp *TestDAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDApp.Contract.contract.Transact(opts, method, params...)
}
