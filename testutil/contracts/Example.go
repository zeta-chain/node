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

// ExampleMetaData contains all meta data concerning the Example contract.
var ExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Foo\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"bar\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevertWithMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevertWithRequire\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doSucceed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_bar\",\"type\":\"uint256\"}],\"name\":\"setBar\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50600080819055506102d8806100276000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c8063352d3fba14610067578063afc874d214610083578063d720cb451461008d578063dd8e556c14610097578063fd5ad965146100a1578063febb0f7e146100ab575b600080fd5b610081600480360381019061007c91906101ce565b6100c9565b005b61008b6100d3565b005b610095610105565b005b61009f610140565b005b6100a9610183565b005b6100b361018d565b6040516100c0919061020a565b60405180910390f35b8060008190555050565b6040517fbfb4ebcf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161013790610282565b60405180910390fd5b6000610181576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161017890610282565b60405180910390fd5b565b6001600081905550565b60005481565b600080fd5b6000819050919050565b6101ab81610198565b81146101b657600080fd5b50565b6000813590506101c8816101a2565b92915050565b6000602082840312156101e4576101e3610193565b5b60006101f2848285016101b9565b91505092915050565b61020481610198565b82525050565b600060208201905061021f60008301846101fb565b92915050565b600082825260208201905092915050565b7f666f6f0000000000000000000000000000000000000000000000000000000000600082015250565b600061026c600383610225565b915061027782610236565b602082019050919050565b6000602082019050818103600083015261029b8161025f565b905091905056fea2646970667358221220b829d5de8f1b16b82dbf677072a361c740900d3b01c6eefba386729591f9958e64736f6c63430008150033",
}

// ExampleABI is the input ABI used to generate the binding from.
// Deprecated: Use ExampleMetaData.ABI instead.
var ExampleABI = ExampleMetaData.ABI

// ExampleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ExampleMetaData.Bin instead.
var ExampleBin = ExampleMetaData.Bin

// DeployExample deploys a new Ethereum contract, binding an instance of Example to it.
func DeployExample(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Example, error) {
	parsed, err := ExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ExampleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Example{ExampleCaller: ExampleCaller{contract: contract}, ExampleTransactor: ExampleTransactor{contract: contract}, ExampleFilterer: ExampleFilterer{contract: contract}}, nil
}

// Example is an auto generated Go binding around an Ethereum contract.
type Example struct {
	ExampleCaller     // Read-only binding to the contract
	ExampleTransactor // Write-only binding to the contract
	ExampleFilterer   // Log filterer for contract events
}

// ExampleCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExampleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExampleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExampleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExampleSession struct {
	Contract     *Example          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExampleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExampleCallerSession struct {
	Contract *ExampleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ExampleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExampleTransactorSession struct {
	Contract     *ExampleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ExampleRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExampleRaw struct {
	Contract *Example // Generic contract binding to access the raw methods on
}

// ExampleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExampleCallerRaw struct {
	Contract *ExampleCaller // Generic read-only contract binding to access the raw methods on
}

// ExampleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExampleTransactorRaw struct {
	Contract *ExampleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExample creates a new instance of Example, bound to a specific deployed contract.
func NewExample(address common.Address, backend bind.ContractBackend) (*Example, error) {
	contract, err := bindExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Example{ExampleCaller: ExampleCaller{contract: contract}, ExampleTransactor: ExampleTransactor{contract: contract}, ExampleFilterer: ExampleFilterer{contract: contract}}, nil
}

// NewExampleCaller creates a new read-only instance of Example, bound to a specific deployed contract.
func NewExampleCaller(address common.Address, caller bind.ContractCaller) (*ExampleCaller, error) {
	contract, err := bindExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleCaller{contract: contract}, nil
}

// NewExampleTransactor creates a new write-only instance of Example, bound to a specific deployed contract.
func NewExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*ExampleTransactor, error) {
	contract, err := bindExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleTransactor{contract: contract}, nil
}

// NewExampleFilterer creates a new log filterer instance of Example, bound to a specific deployed contract.
func NewExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*ExampleFilterer, error) {
	contract, err := bindExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExampleFilterer{contract: contract}, nil
}

// bindExample binds a generic wrapper to an already deployed contract.
func bindExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Example *ExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Example.Contract.ExampleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Example *ExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.Contract.ExampleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Example *ExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Example.Contract.ExampleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Example *ExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Example.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Example *ExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Example *ExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Example.Contract.contract.Transact(opts, method, params...)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_Example *ExampleCaller) Bar(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Example.contract.Call(opts, &out, "bar")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_Example *ExampleSession) Bar() (*big.Int, error) {
	return _Example.Contract.Bar(&_Example.CallOpts)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_Example *ExampleCallerSession) Bar() (*big.Int, error) {
	return _Example.Contract.Bar(&_Example.CallOpts)
}

// DoRevert is a paid mutator transaction binding the contract method 0xafc874d2.
//
// Solidity: function doRevert() returns()
func (_Example *ExampleTransactor) DoRevert(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "doRevert")
}

// DoRevert is a paid mutator transaction binding the contract method 0xafc874d2.
//
// Solidity: function doRevert() returns()
func (_Example *ExampleSession) DoRevert() (*types.Transaction, error) {
	return _Example.Contract.DoRevert(&_Example.TransactOpts)
}

// DoRevert is a paid mutator transaction binding the contract method 0xafc874d2.
//
// Solidity: function doRevert() returns()
func (_Example *ExampleTransactorSession) DoRevert() (*types.Transaction, error) {
	return _Example.Contract.DoRevert(&_Example.TransactOpts)
}

// DoRevertWithMessage is a paid mutator transaction binding the contract method 0xd720cb45.
//
// Solidity: function doRevertWithMessage() returns()
func (_Example *ExampleTransactor) DoRevertWithMessage(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "doRevertWithMessage")
}

// DoRevertWithMessage is a paid mutator transaction binding the contract method 0xd720cb45.
//
// Solidity: function doRevertWithMessage() returns()
func (_Example *ExampleSession) DoRevertWithMessage() (*types.Transaction, error) {
	return _Example.Contract.DoRevertWithMessage(&_Example.TransactOpts)
}

// DoRevertWithMessage is a paid mutator transaction binding the contract method 0xd720cb45.
//
// Solidity: function doRevertWithMessage() returns()
func (_Example *ExampleTransactorSession) DoRevertWithMessage() (*types.Transaction, error) {
	return _Example.Contract.DoRevertWithMessage(&_Example.TransactOpts)
}

// DoRevertWithRequire is a paid mutator transaction binding the contract method 0xdd8e556c.
//
// Solidity: function doRevertWithRequire() returns()
func (_Example *ExampleTransactor) DoRevertWithRequire(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "doRevertWithRequire")
}

// DoRevertWithRequire is a paid mutator transaction binding the contract method 0xdd8e556c.
//
// Solidity: function doRevertWithRequire() returns()
func (_Example *ExampleSession) DoRevertWithRequire() (*types.Transaction, error) {
	return _Example.Contract.DoRevertWithRequire(&_Example.TransactOpts)
}

// DoRevertWithRequire is a paid mutator transaction binding the contract method 0xdd8e556c.
//
// Solidity: function doRevertWithRequire() returns()
func (_Example *ExampleTransactorSession) DoRevertWithRequire() (*types.Transaction, error) {
	return _Example.Contract.DoRevertWithRequire(&_Example.TransactOpts)
}

// DoSucceed is a paid mutator transaction binding the contract method 0xfd5ad965.
//
// Solidity: function doSucceed() returns()
func (_Example *ExampleTransactor) DoSucceed(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "doSucceed")
}

// DoSucceed is a paid mutator transaction binding the contract method 0xfd5ad965.
//
// Solidity: function doSucceed() returns()
func (_Example *ExampleSession) DoSucceed() (*types.Transaction, error) {
	return _Example.Contract.DoSucceed(&_Example.TransactOpts)
}

// DoSucceed is a paid mutator transaction binding the contract method 0xfd5ad965.
//
// Solidity: function doSucceed() returns()
func (_Example *ExampleTransactorSession) DoSucceed() (*types.Transaction, error) {
	return _Example.Contract.DoSucceed(&_Example.TransactOpts)
}

// SetBar is a paid mutator transaction binding the contract method 0x352d3fba.
//
// Solidity: function setBar(uint256 _bar) returns()
func (_Example *ExampleTransactor) SetBar(opts *bind.TransactOpts, _bar *big.Int) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "setBar", _bar)
}

// SetBar is a paid mutator transaction binding the contract method 0x352d3fba.
//
// Solidity: function setBar(uint256 _bar) returns()
func (_Example *ExampleSession) SetBar(_bar *big.Int) (*types.Transaction, error) {
	return _Example.Contract.SetBar(&_Example.TransactOpts, _bar)
}

// SetBar is a paid mutator transaction binding the contract method 0x352d3fba.
//
// Solidity: function setBar(uint256 _bar) returns()
func (_Example *ExampleTransactorSession) SetBar(_bar *big.Int) (*types.Transaction, error) {
	return _Example.Contract.SetBar(&_Example.TransactOpts, _bar)
}
