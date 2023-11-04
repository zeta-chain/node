// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// DepositorMetaData contains all meta data concerning the Depositor contract.
var DepositorMetaData = &bind.MetaData{
	ABI: "null",
	Bin: "0xnull",
}

// DepositorABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositorMetaData.ABI instead.
var DepositorABI = DepositorMetaData.ABI

// DepositorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DepositorMetaData.Bin instead.
var DepositorBin = DepositorMetaData.Bin

// DeployDepositor deploys a new Ethereum contract, binding an instance of Depositor to it.
func DeployDepositor(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Depositor, error) {
	parsed, err := DepositorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DepositorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Depositor{DepositorCaller: DepositorCaller{contract: contract}, DepositorTransactor: DepositorTransactor{contract: contract}, DepositorFilterer: DepositorFilterer{contract: contract}}, nil
}

// Depositor is an auto generated Go binding around an Ethereum contract.
type Depositor struct {
	DepositorCaller     // Read-only binding to the contract
	DepositorTransactor // Write-only binding to the contract
	DepositorFilterer   // Log filterer for contract events
}

// DepositorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositorSession struct {
	Contract     *Depositor        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositorCallerSession struct {
	Contract *DepositorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DepositorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositorTransactorSession struct {
	Contract     *DepositorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DepositorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositorRaw struct {
	Contract *Depositor // Generic contract binding to access the raw methods on
}

// DepositorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositorCallerRaw struct {
	Contract *DepositorCaller // Generic read-only contract binding to access the raw methods on
}

// DepositorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositorTransactorRaw struct {
	Contract *DepositorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositor creates a new instance of Depositor, bound to a specific deployed contract.
func NewDepositor(address common.Address, backend bind.ContractBackend) (*Depositor, error) {
	contract, err := bindDepositor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Depositor{DepositorCaller: DepositorCaller{contract: contract}, DepositorTransactor: DepositorTransactor{contract: contract}, DepositorFilterer: DepositorFilterer{contract: contract}}, nil
}

// NewDepositorCaller creates a new read-only instance of Depositor, bound to a specific deployed contract.
func NewDepositorCaller(address common.Address, caller bind.ContractCaller) (*DepositorCaller, error) {
	contract, err := bindDepositor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositorCaller{contract: contract}, nil
}

// NewDepositorTransactor creates a new write-only instance of Depositor, bound to a specific deployed contract.
func NewDepositorTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositorTransactor, error) {
	contract, err := bindDepositor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositorTransactor{contract: contract}, nil
}

// NewDepositorFilterer creates a new log filterer instance of Depositor, bound to a specific deployed contract.
func NewDepositorFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositorFilterer, error) {
	contract, err := bindDepositor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositorFilterer{contract: contract}, nil
}

// bindDepositor binds a generic wrapper to an already deployed contract.
func bindDepositor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depositor *DepositorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Depositor.Contract.DepositorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depositor *DepositorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositor.Contract.DepositorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depositor *DepositorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depositor.Contract.DepositorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depositor *DepositorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Depositor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depositor *DepositorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depositor *DepositorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depositor.Contract.contract.Transact(opts, method, params...)
}
