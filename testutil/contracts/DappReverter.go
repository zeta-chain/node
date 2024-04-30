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
)

// DappReverterMetaData contains all meta data concerning the DappReverter contract.
var DappReverterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"onZetaMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"onZetaRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b50608180601d6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063705847b71460375780639a19074914603f575b600080fd5b603d6047565b005b60456049565b005b565b56fea26469706673582212202ae32d3809d629fd01d309562a51297d761f547f4633bf45c73bf33c9955651164736f6c63430008190033",
}

// DappReverterABI is the input ABI used to generate the binding from.
// Deprecated: Use DappReverterMetaData.ABI instead.
var DappReverterABI = DappReverterMetaData.ABI

// DappReverterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DappReverterMetaData.Bin instead.
var DappReverterBin = DappReverterMetaData.Bin

// DeployDappReverter deploys a new Ethereum contract, binding an instance of DappReverter to it.
func DeployDappReverter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DappReverter, error) {
	parsed, err := DappReverterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DappReverterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DappReverter{DappReverterCaller: DappReverterCaller{contract: contract}, DappReverterTransactor: DappReverterTransactor{contract: contract}, DappReverterFilterer: DappReverterFilterer{contract: contract}}, nil
}

// DappReverter is an auto generated Go binding around an Ethereum contract.
type DappReverter struct {
	DappReverterCaller     // Read-only binding to the contract
	DappReverterTransactor // Write-only binding to the contract
	DappReverterFilterer   // Log filterer for contract events
}

// DappReverterCaller is an auto generated read-only Go binding around an Ethereum contract.
type DappReverterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DappReverterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DappReverterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DappReverterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DappReverterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DappReverterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DappReverterSession struct {
	Contract     *DappReverter     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DappReverterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DappReverterCallerSession struct {
	Contract *DappReverterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// DappReverterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DappReverterTransactorSession struct {
	Contract     *DappReverterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// DappReverterRaw is an auto generated low-level Go binding around an Ethereum contract.
type DappReverterRaw struct {
	Contract *DappReverter // Generic contract binding to access the raw methods on
}

// DappReverterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DappReverterCallerRaw struct {
	Contract *DappReverterCaller // Generic read-only contract binding to access the raw methods on
}

// DappReverterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DappReverterTransactorRaw struct {
	Contract *DappReverterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDappReverter creates a new instance of DappReverter, bound to a specific deployed contract.
func NewDappReverter(address common.Address, backend bind.ContractBackend) (*DappReverter, error) {
	contract, err := bindDappReverter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DappReverter{DappReverterCaller: DappReverterCaller{contract: contract}, DappReverterTransactor: DappReverterTransactor{contract: contract}, DappReverterFilterer: DappReverterFilterer{contract: contract}}, nil
}

// NewDappReverterCaller creates a new read-only instance of DappReverter, bound to a specific deployed contract.
func NewDappReverterCaller(address common.Address, caller bind.ContractCaller) (*DappReverterCaller, error) {
	contract, err := bindDappReverter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DappReverterCaller{contract: contract}, nil
}

// NewDappReverterTransactor creates a new write-only instance of DappReverter, bound to a specific deployed contract.
func NewDappReverterTransactor(address common.Address, transactor bind.ContractTransactor) (*DappReverterTransactor, error) {
	contract, err := bindDappReverter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DappReverterTransactor{contract: contract}, nil
}

// NewDappReverterFilterer creates a new log filterer instance of DappReverter, bound to a specific deployed contract.
func NewDappReverterFilterer(address common.Address, filterer bind.ContractFilterer) (*DappReverterFilterer, error) {
	contract, err := bindDappReverter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DappReverterFilterer{contract: contract}, nil
}

// bindDappReverter binds a generic wrapper to an already deployed contract.
func bindDappReverter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DappReverterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DappReverter *DappReverterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DappReverter.Contract.DappReverterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DappReverter *DappReverterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DappReverter.Contract.DappReverterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DappReverter *DappReverterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DappReverter.Contract.DappReverterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DappReverter *DappReverterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DappReverter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DappReverter *DappReverterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DappReverter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DappReverter *DappReverterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DappReverter.Contract.contract.Transact(opts, method, params...)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x705847b7.
//
// Solidity: function onZetaMessage() returns()
func (_DappReverter *DappReverterTransactor) OnZetaMessage(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DappReverter.contract.Transact(opts, "onZetaMessage")
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x705847b7.
//
// Solidity: function onZetaMessage() returns()
func (_DappReverter *DappReverterSession) OnZetaMessage() (*types.Transaction, error) {
	return _DappReverter.Contract.OnZetaMessage(&_DappReverter.TransactOpts)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x705847b7.
//
// Solidity: function onZetaMessage() returns()
func (_DappReverter *DappReverterTransactorSession) OnZetaMessage() (*types.Transaction, error) {
	return _DappReverter.Contract.OnZetaMessage(&_DappReverter.TransactOpts)
}

// OnZetaRevert is a paid mutator transaction binding the contract method 0x9a190749.
//
// Solidity: function onZetaRevert() returns()
func (_DappReverter *DappReverterTransactor) OnZetaRevert(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DappReverter.contract.Transact(opts, "onZetaRevert")
}

// OnZetaRevert is a paid mutator transaction binding the contract method 0x9a190749.
//
// Solidity: function onZetaRevert() returns()
func (_DappReverter *DappReverterSession) OnZetaRevert() (*types.Transaction, error) {
	return _DappReverter.Contract.OnZetaRevert(&_DappReverter.TransactOpts)
}

// OnZetaRevert is a paid mutator transaction binding the contract method 0x9a190749.
//
// Solidity: function onZetaRevert() returns()
func (_DappReverter *DappReverterTransactorSession) OnZetaRevert() (*types.Transaction, error) {
	return _DappReverter.Contract.OnZetaRevert(&_DappReverter.TransactOpts)
}
