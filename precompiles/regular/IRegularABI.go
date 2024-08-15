// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package regular

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

// RegularMetaData contains all meta data concerning the Regular contract.
var RegularMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"name\":\"bech32ToHexAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"prefix\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"bech32ify\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"method\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"regularCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"result\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RegularABI is the input ABI used to generate the binding from.
// Deprecated: Use RegularMetaData.ABI instead.
var RegularABI = RegularMetaData.ABI

// Regular is an auto generated Go binding around an Ethereum contract.
type Regular struct {
	RegularCaller     // Read-only binding to the contract
	RegularTransactor // Write-only binding to the contract
	RegularFilterer   // Log filterer for contract events
}

// RegularCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegularCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegularTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegularTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegularFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegularFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegularSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegularSession struct {
	Contract     *Regular          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegularCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegularCallerSession struct {
	Contract *RegularCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// RegularTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegularTransactorSession struct {
	Contract     *RegularTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// RegularRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegularRaw struct {
	Contract *Regular // Generic contract binding to access the raw methods on
}

// RegularCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegularCallerRaw struct {
	Contract *RegularCaller // Generic read-only contract binding to access the raw methods on
}

// RegularTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegularTransactorRaw struct {
	Contract *RegularTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegular creates a new instance of Regular, bound to a specific deployed contract.
func NewRegular(address common.Address, backend bind.ContractBackend) (*Regular, error) {
	contract, err := bindRegular(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Regular{RegularCaller: RegularCaller{contract: contract}, RegularTransactor: RegularTransactor{contract: contract}, RegularFilterer: RegularFilterer{contract: contract}}, nil
}

// NewRegularCaller creates a new read-only instance of Regular, bound to a specific deployed contract.
func NewRegularCaller(address common.Address, caller bind.ContractCaller) (*RegularCaller, error) {
	contract, err := bindRegular(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegularCaller{contract: contract}, nil
}

// NewRegularTransactor creates a new write-only instance of Regular, bound to a specific deployed contract.
func NewRegularTransactor(address common.Address, transactor bind.ContractTransactor) (*RegularTransactor, error) {
	contract, err := bindRegular(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegularTransactor{contract: contract}, nil
}

// NewRegularFilterer creates a new log filterer instance of Regular, bound to a specific deployed contract.
func NewRegularFilterer(address common.Address, filterer bind.ContractFilterer) (*RegularFilterer, error) {
	contract, err := bindRegular(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegularFilterer{contract: contract}, nil
}

// bindRegular binds a generic wrapper to an already deployed contract.
func bindRegular(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RegularMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Regular *RegularRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Regular.Contract.RegularCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Regular *RegularRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Regular.Contract.RegularTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Regular *RegularRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Regular.Contract.RegularTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Regular *RegularCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Regular.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Regular *RegularTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Regular.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Regular *RegularTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Regular.Contract.contract.Transact(opts, method, params...)
}

// Bech32ToHexAddr is a free data retrieval call binding the contract method 0xe4e2a4ec.
//
// Solidity: function bech32ToHexAddr(string bech32) view returns(address addr)
func (_Regular *RegularCaller) Bech32ToHexAddr(opts *bind.CallOpts, bech32 string) (common.Address, error) {
	var out []interface{}
	err := _Regular.contract.Call(opts, &out, "bech32ToHexAddr", bech32)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bech32ToHexAddr is a free data retrieval call binding the contract method 0xe4e2a4ec.
//
// Solidity: function bech32ToHexAddr(string bech32) view returns(address addr)
func (_Regular *RegularSession) Bech32ToHexAddr(bech32 string) (common.Address, error) {
	return _Regular.Contract.Bech32ToHexAddr(&_Regular.CallOpts, bech32)
}

// Bech32ToHexAddr is a free data retrieval call binding the contract method 0xe4e2a4ec.
//
// Solidity: function bech32ToHexAddr(string bech32) view returns(address addr)
func (_Regular *RegularCallerSession) Bech32ToHexAddr(bech32 string) (common.Address, error) {
	return _Regular.Contract.Bech32ToHexAddr(&_Regular.CallOpts, bech32)
}

// Bech32ify is a free data retrieval call binding the contract method 0x0615b74e.
//
// Solidity: function bech32ify(string prefix, address addr) view returns(string bech32)
func (_Regular *RegularCaller) Bech32ify(opts *bind.CallOpts, prefix string, addr common.Address) (string, error) {
	var out []interface{}
	err := _Regular.contract.Call(opts, &out, "bech32ify", prefix, addr)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Bech32ify is a free data retrieval call binding the contract method 0x0615b74e.
//
// Solidity: function bech32ify(string prefix, address addr) view returns(string bech32)
func (_Regular *RegularSession) Bech32ify(prefix string, addr common.Address) (string, error) {
	return _Regular.Contract.Bech32ify(&_Regular.CallOpts, prefix, addr)
}

// Bech32ify is a free data retrieval call binding the contract method 0x0615b74e.
//
// Solidity: function bech32ify(string prefix, address addr) view returns(string bech32)
func (_Regular *RegularCallerSession) Bech32ify(prefix string, addr common.Address) (string, error) {
	return _Regular.Contract.Bech32ify(&_Regular.CallOpts, prefix, addr)
}

// RegularCall is a paid mutator transaction binding the contract method 0x93e3663d.
//
// Solidity: function regularCall(string method, address addr) returns(uint256 result)
func (_Regular *RegularTransactor) RegularCall(opts *bind.TransactOpts, method string, addr common.Address) (*types.Transaction, error) {
	return _Regular.contract.Transact(opts, "regularCall", method, addr)
}

// RegularCall is a paid mutator transaction binding the contract method 0x93e3663d.
//
// Solidity: function regularCall(string method, address addr) returns(uint256 result)
func (_Regular *RegularSession) RegularCall(method string, addr common.Address) (*types.Transaction, error) {
	return _Regular.Contract.RegularCall(&_Regular.TransactOpts, method, addr)
}

// RegularCall is a paid mutator transaction binding the contract method 0x93e3663d.
//
// Solidity: function regularCall(string method, address addr) returns(uint256 result)
func (_Regular *RegularTransactorSession) RegularCall(method string, addr common.Address) (*types.Transaction, error) {
	return _Regular.Contract.RegularCall(&_Regular.TransactOpts, method, addr)
}
