// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package iuniswapv3poolimmutables

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

// IUniswapV3PoolImmutablesMetaData contains all meta data concerning the IUniswapV3PoolImmutables contract.
var IUniswapV3PoolImmutablesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fee\",\"outputs\":[{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxLiquidityPerTick\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tickSpacing\",\"outputs\":[{\"internalType\":\"int24\",\"name\":\"\",\"type\":\"int24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token0\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token1\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IUniswapV3PoolImmutablesABI is the input ABI used to generate the binding from.
// Deprecated: Use IUniswapV3PoolImmutablesMetaData.ABI instead.
var IUniswapV3PoolImmutablesABI = IUniswapV3PoolImmutablesMetaData.ABI

// IUniswapV3PoolImmutables is an auto generated Go binding around an Ethereum contract.
type IUniswapV3PoolImmutables struct {
	IUniswapV3PoolImmutablesCaller     // Read-only binding to the contract
	IUniswapV3PoolImmutablesTransactor // Write-only binding to the contract
	IUniswapV3PoolImmutablesFilterer   // Log filterer for contract events
}

// IUniswapV3PoolImmutablesCaller is an auto generated read-only Go binding around an Ethereum contract.
type IUniswapV3PoolImmutablesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolImmutablesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IUniswapV3PoolImmutablesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolImmutablesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IUniswapV3PoolImmutablesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV3PoolImmutablesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IUniswapV3PoolImmutablesSession struct {
	Contract     *IUniswapV3PoolImmutables // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IUniswapV3PoolImmutablesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IUniswapV3PoolImmutablesCallerSession struct {
	Contract *IUniswapV3PoolImmutablesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// IUniswapV3PoolImmutablesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IUniswapV3PoolImmutablesTransactorSession struct {
	Contract     *IUniswapV3PoolImmutablesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// IUniswapV3PoolImmutablesRaw is an auto generated low-level Go binding around an Ethereum contract.
type IUniswapV3PoolImmutablesRaw struct {
	Contract *IUniswapV3PoolImmutables // Generic contract binding to access the raw methods on
}

// IUniswapV3PoolImmutablesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IUniswapV3PoolImmutablesCallerRaw struct {
	Contract *IUniswapV3PoolImmutablesCaller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV3PoolImmutablesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IUniswapV3PoolImmutablesTransactorRaw struct {
	Contract *IUniswapV3PoolImmutablesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV3PoolImmutables creates a new instance of IUniswapV3PoolImmutables, bound to a specific deployed contract.
func NewIUniswapV3PoolImmutables(address common.Address, backend bind.ContractBackend) (*IUniswapV3PoolImmutables, error) {
	contract, err := bindIUniswapV3PoolImmutables(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolImmutables{IUniswapV3PoolImmutablesCaller: IUniswapV3PoolImmutablesCaller{contract: contract}, IUniswapV3PoolImmutablesTransactor: IUniswapV3PoolImmutablesTransactor{contract: contract}, IUniswapV3PoolImmutablesFilterer: IUniswapV3PoolImmutablesFilterer{contract: contract}}, nil
}

// NewIUniswapV3PoolImmutablesCaller creates a new read-only instance of IUniswapV3PoolImmutables, bound to a specific deployed contract.
func NewIUniswapV3PoolImmutablesCaller(address common.Address, caller bind.ContractCaller) (*IUniswapV3PoolImmutablesCaller, error) {
	contract, err := bindIUniswapV3PoolImmutables(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolImmutablesCaller{contract: contract}, nil
}

// NewIUniswapV3PoolImmutablesTransactor creates a new write-only instance of IUniswapV3PoolImmutables, bound to a specific deployed contract.
func NewIUniswapV3PoolImmutablesTransactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV3PoolImmutablesTransactor, error) {
	contract, err := bindIUniswapV3PoolImmutables(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolImmutablesTransactor{contract: contract}, nil
}

// NewIUniswapV3PoolImmutablesFilterer creates a new log filterer instance of IUniswapV3PoolImmutables, bound to a specific deployed contract.
func NewIUniswapV3PoolImmutablesFilterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV3PoolImmutablesFilterer, error) {
	contract, err := bindIUniswapV3PoolImmutables(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV3PoolImmutablesFilterer{contract: contract}, nil
}

// bindIUniswapV3PoolImmutables binds a generic wrapper to an already deployed contract.
func bindIUniswapV3PoolImmutables(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUniswapV3PoolImmutablesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV3PoolImmutables.Contract.IUniswapV3PoolImmutablesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV3PoolImmutables.Contract.IUniswapV3PoolImmutablesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV3PoolImmutables.Contract.IUniswapV3PoolImmutablesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV3PoolImmutables.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV3PoolImmutables.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV3PoolImmutables.Contract.contract.Transact(opts, method, params...)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV3PoolImmutables.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesSession) Factory() (common.Address, error) {
	return _IUniswapV3PoolImmutables.Contract.Factory(&_IUniswapV3PoolImmutables.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCallerSession) Factory() (common.Address, error) {
	return _IUniswapV3PoolImmutables.Contract.Factory(&_IUniswapV3PoolImmutables.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint24)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV3PoolImmutables.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint24)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesSession) Fee() (*big.Int, error) {
	return _IUniswapV3PoolImmutables.Contract.Fee(&_IUniswapV3PoolImmutables.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint24)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCallerSession) Fee() (*big.Int, error) {
	return _IUniswapV3PoolImmutables.Contract.Fee(&_IUniswapV3PoolImmutables.CallOpts)
}

// MaxLiquidityPerTick is a free data retrieval call binding the contract method 0x70cf754a.
//
// Solidity: function maxLiquidityPerTick() view returns(uint128)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCaller) MaxLiquidityPerTick(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV3PoolImmutables.contract.Call(opts, &out, "maxLiquidityPerTick")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxLiquidityPerTick is a free data retrieval call binding the contract method 0x70cf754a.
//
// Solidity: function maxLiquidityPerTick() view returns(uint128)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesSession) MaxLiquidityPerTick() (*big.Int, error) {
	return _IUniswapV3PoolImmutables.Contract.MaxLiquidityPerTick(&_IUniswapV3PoolImmutables.CallOpts)
}

// MaxLiquidityPerTick is a free data retrieval call binding the contract method 0x70cf754a.
//
// Solidity: function maxLiquidityPerTick() view returns(uint128)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCallerSession) MaxLiquidityPerTick() (*big.Int, error) {
	return _IUniswapV3PoolImmutables.Contract.MaxLiquidityPerTick(&_IUniswapV3PoolImmutables.CallOpts)
}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCaller) TickSpacing(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV3PoolImmutables.contract.Call(opts, &out, "tickSpacing")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesSession) TickSpacing() (*big.Int, error) {
	return _IUniswapV3PoolImmutables.Contract.TickSpacing(&_IUniswapV3PoolImmutables.CallOpts)
}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCallerSession) TickSpacing() (*big.Int, error) {
	return _IUniswapV3PoolImmutables.Contract.TickSpacing(&_IUniswapV3PoolImmutables.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCaller) Token0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV3PoolImmutables.contract.Call(opts, &out, "token0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesSession) Token0() (common.Address, error) {
	return _IUniswapV3PoolImmutables.Contract.Token0(&_IUniswapV3PoolImmutables.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCallerSession) Token0() (common.Address, error) {
	return _IUniswapV3PoolImmutables.Contract.Token0(&_IUniswapV3PoolImmutables.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCaller) Token1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV3PoolImmutables.contract.Call(opts, &out, "token1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesSession) Token1() (common.Address, error) {
	return _IUniswapV3PoolImmutables.Contract.Token1(&_IUniswapV3PoolImmutables.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_IUniswapV3PoolImmutables *IUniswapV3PoolImmutablesCallerSession) Token1() (common.Address, error) {
	return _IUniswapV3PoolImmutables.Contract.Token1(&_IUniswapV3PoolImmutables.CallOpts)
}
