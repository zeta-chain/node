// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package iuniswapv2factory

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

// IUniswapV2FactoryMetaData contains all meta data concerning the IUniswapV2Factory contract.
var IUniswapV2FactoryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"PairCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allPairs\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"allPairsLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"}],\"name\":\"createPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeToSetter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"}],\"name\":\"getPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"setFeeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"setFeeToSetter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IUniswapV2FactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use IUniswapV2FactoryMetaData.ABI instead.
var IUniswapV2FactoryABI = IUniswapV2FactoryMetaData.ABI

// IUniswapV2Factory is an auto generated Go binding around an Ethereum contract.
type IUniswapV2Factory struct {
	IUniswapV2FactoryCaller     // Read-only binding to the contract
	IUniswapV2FactoryTransactor // Write-only binding to the contract
	IUniswapV2FactoryFilterer   // Log filterer for contract events
}

// IUniswapV2FactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type IUniswapV2FactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2FactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IUniswapV2FactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2FactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IUniswapV2FactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2FactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IUniswapV2FactorySession struct {
	Contract     *IUniswapV2Factory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IUniswapV2FactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IUniswapV2FactoryCallerSession struct {
	Contract *IUniswapV2FactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// IUniswapV2FactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IUniswapV2FactoryTransactorSession struct {
	Contract     *IUniswapV2FactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// IUniswapV2FactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type IUniswapV2FactoryRaw struct {
	Contract *IUniswapV2Factory // Generic contract binding to access the raw methods on
}

// IUniswapV2FactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IUniswapV2FactoryCallerRaw struct {
	Contract *IUniswapV2FactoryCaller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV2FactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IUniswapV2FactoryTransactorRaw struct {
	Contract *IUniswapV2FactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV2Factory creates a new instance of IUniswapV2Factory, bound to a specific deployed contract.
func NewIUniswapV2Factory(address common.Address, backend bind.ContractBackend) (*IUniswapV2Factory, error) {
	contract, err := bindIUniswapV2Factory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2Factory{IUniswapV2FactoryCaller: IUniswapV2FactoryCaller{contract: contract}, IUniswapV2FactoryTransactor: IUniswapV2FactoryTransactor{contract: contract}, IUniswapV2FactoryFilterer: IUniswapV2FactoryFilterer{contract: contract}}, nil
}

// NewIUniswapV2FactoryCaller creates a new read-only instance of IUniswapV2Factory, bound to a specific deployed contract.
func NewIUniswapV2FactoryCaller(address common.Address, caller bind.ContractCaller) (*IUniswapV2FactoryCaller, error) {
	contract, err := bindIUniswapV2Factory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2FactoryCaller{contract: contract}, nil
}

// NewIUniswapV2FactoryTransactor creates a new write-only instance of IUniswapV2Factory, bound to a specific deployed contract.
func NewIUniswapV2FactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV2FactoryTransactor, error) {
	contract, err := bindIUniswapV2Factory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2FactoryTransactor{contract: contract}, nil
}

// NewIUniswapV2FactoryFilterer creates a new log filterer instance of IUniswapV2Factory, bound to a specific deployed contract.
func NewIUniswapV2FactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV2FactoryFilterer, error) {
	contract, err := bindIUniswapV2Factory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2FactoryFilterer{contract: contract}, nil
}

// bindIUniswapV2Factory binds a generic wrapper to an already deployed contract.
func bindIUniswapV2Factory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUniswapV2FactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Factory *IUniswapV2FactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Factory.Contract.IUniswapV2FactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Factory *IUniswapV2FactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.IUniswapV2FactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Factory *IUniswapV2FactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.IUniswapV2FactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Factory *IUniswapV2FactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Factory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Factory *IUniswapV2FactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Factory *IUniswapV2FactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.contract.Transact(opts, method, params...)
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryCaller) AllPairs(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Factory.contract.Call(opts, &out, "allPairs", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactorySession) AllPairs(arg0 *big.Int) (common.Address, error) {
	return _IUniswapV2Factory.Contract.AllPairs(&_IUniswapV2Factory.CallOpts, arg0)
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryCallerSession) AllPairs(arg0 *big.Int) (common.Address, error) {
	return _IUniswapV2Factory.Contract.AllPairs(&_IUniswapV2Factory.CallOpts, arg0)
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_IUniswapV2Factory *IUniswapV2FactoryCaller) AllPairsLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2Factory.contract.Call(opts, &out, "allPairsLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_IUniswapV2Factory *IUniswapV2FactorySession) AllPairsLength() (*big.Int, error) {
	return _IUniswapV2Factory.Contract.AllPairsLength(&_IUniswapV2Factory.CallOpts)
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_IUniswapV2Factory *IUniswapV2FactoryCallerSession) AllPairsLength() (*big.Int, error) {
	return _IUniswapV2Factory.Contract.AllPairsLength(&_IUniswapV2Factory.CallOpts)
}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactoryCaller) FeeTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Factory.contract.Call(opts, &out, "feeTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactorySession) FeeTo() (common.Address, error) {
	return _IUniswapV2Factory.Contract.FeeTo(&_IUniswapV2Factory.CallOpts)
}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactoryCallerSession) FeeTo() (common.Address, error) {
	return _IUniswapV2Factory.Contract.FeeTo(&_IUniswapV2Factory.CallOpts)
}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactoryCaller) FeeToSetter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Factory.contract.Call(opts, &out, "feeToSetter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactorySession) FeeToSetter() (common.Address, error) {
	return _IUniswapV2Factory.Contract.FeeToSetter(&_IUniswapV2Factory.CallOpts)
}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactoryCallerSession) FeeToSetter() (common.Address, error) {
	return _IUniswapV2Factory.Contract.FeeToSetter(&_IUniswapV2Factory.CallOpts)
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address tokenA, address tokenB) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryCaller) GetPair(opts *bind.CallOpts, tokenA common.Address, tokenB common.Address) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Factory.contract.Call(opts, &out, "getPair", tokenA, tokenB)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address tokenA, address tokenB) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactorySession) GetPair(tokenA common.Address, tokenB common.Address) (common.Address, error) {
	return _IUniswapV2Factory.Contract.GetPair(&_IUniswapV2Factory.CallOpts, tokenA, tokenB)
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address tokenA, address tokenB) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryCallerSession) GetPair(tokenA common.Address, tokenB common.Address) (common.Address, error) {
	return _IUniswapV2Factory.Contract.GetPair(&_IUniswapV2Factory.CallOpts, tokenA, tokenB)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryTransactor) CreatePair(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.contract.Transact(opts, "createPair", tokenA, tokenB)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactorySession) CreatePair(tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.CreatePair(&_IUniswapV2Factory.TransactOpts, tokenA, tokenB)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryTransactorSession) CreatePair(tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.CreatePair(&_IUniswapV2Factory.TransactOpts, tokenA, tokenB)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactoryTransactor) SetFeeTo(opts *bind.TransactOpts, arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.contract.Transact(opts, "setFeeTo", arg0)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactorySession) SetFeeTo(arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.SetFeeTo(&_IUniswapV2Factory.TransactOpts, arg0)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactoryTransactorSession) SetFeeTo(arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.SetFeeTo(&_IUniswapV2Factory.TransactOpts, arg0)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactoryTransactor) SetFeeToSetter(opts *bind.TransactOpts, arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.contract.Transact(opts, "setFeeToSetter", arg0)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactorySession) SetFeeToSetter(arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.SetFeeToSetter(&_IUniswapV2Factory.TransactOpts, arg0)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactoryTransactorSession) SetFeeToSetter(arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.SetFeeToSetter(&_IUniswapV2Factory.TransactOpts, arg0)
}

// IUniswapV2FactoryPairCreatedIterator is returned from FilterPairCreated and is used to iterate over the raw logs and unpacked data for PairCreated events raised by the IUniswapV2Factory contract.
type IUniswapV2FactoryPairCreatedIterator struct {
	Event *IUniswapV2FactoryPairCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IUniswapV2FactoryPairCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUniswapV2FactoryPairCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IUniswapV2FactoryPairCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IUniswapV2FactoryPairCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUniswapV2FactoryPairCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUniswapV2FactoryPairCreated represents a PairCreated event raised by the IUniswapV2Factory contract.
type IUniswapV2FactoryPairCreated struct {
	Token0 common.Address
	Token1 common.Address
	Pair   common.Address
	Arg3   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPairCreated is a free log retrieval operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_IUniswapV2Factory *IUniswapV2FactoryFilterer) FilterPairCreated(opts *bind.FilterOpts, token0 []common.Address, token1 []common.Address) (*IUniswapV2FactoryPairCreatedIterator, error) {

	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}

	logs, sub, err := _IUniswapV2Factory.contract.FilterLogs(opts, "PairCreated", token0Rule, token1Rule)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2FactoryPairCreatedIterator{contract: _IUniswapV2Factory.contract, event: "PairCreated", logs: logs, sub: sub}, nil
}

// WatchPairCreated is a free log subscription operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_IUniswapV2Factory *IUniswapV2FactoryFilterer) WatchPairCreated(opts *bind.WatchOpts, sink chan<- *IUniswapV2FactoryPairCreated, token0 []common.Address, token1 []common.Address) (event.Subscription, error) {

	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}

	logs, sub, err := _IUniswapV2Factory.contract.WatchLogs(opts, "PairCreated", token0Rule, token1Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUniswapV2FactoryPairCreated)
				if err := _IUniswapV2Factory.contract.UnpackLog(event, "PairCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePairCreated is a log parse operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_IUniswapV2Factory *IUniswapV2FactoryFilterer) ParsePairCreated(log types.Log) (*IUniswapV2FactoryPairCreated, error) {
	event := new(IUniswapV2FactoryPairCreated)
	if err := _IUniswapV2Factory.contract.UnpackLog(event, "PairCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
