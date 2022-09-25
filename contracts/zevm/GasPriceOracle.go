// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zevm

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

// GasPriceOracleMetaData contains all meta data concerning the GasPriceOracle contract.
var GasPriceOracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"FUNGIBLE_MODULE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasCoinERC4\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"erc4\",\"type\":\"address\"}],\"name\":\"setGasCoinERC4\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"setGasPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// GasPriceOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use GasPriceOracleMetaData.ABI instead.
var GasPriceOracleABI = GasPriceOracleMetaData.ABI

// GasPriceOracle is an auto generated Go binding around an Ethereum contract.
type GasPriceOracle struct {
	GasPriceOracleCaller     // Read-only binding to the contract
	GasPriceOracleTransactor // Write-only binding to the contract
	GasPriceOracleFilterer   // Log filterer for contract events
}

// GasPriceOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type GasPriceOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasPriceOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GasPriceOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasPriceOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GasPriceOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GasPriceOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GasPriceOracleSession struct {
	Contract     *GasPriceOracle   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GasPriceOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GasPriceOracleCallerSession struct {
	Contract *GasPriceOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// GasPriceOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GasPriceOracleTransactorSession struct {
	Contract     *GasPriceOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// GasPriceOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type GasPriceOracleRaw struct {
	Contract *GasPriceOracle // Generic contract binding to access the raw methods on
}

// GasPriceOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GasPriceOracleCallerRaw struct {
	Contract *GasPriceOracleCaller // Generic read-only contract binding to access the raw methods on
}

// GasPriceOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GasPriceOracleTransactorRaw struct {
	Contract *GasPriceOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGasPriceOracle creates a new instance of GasPriceOracle, bound to a specific deployed contract.
func NewGasPriceOracle(address common.Address, backend bind.ContractBackend) (*GasPriceOracle, error) {
	contract, err := bindGasPriceOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GasPriceOracle{GasPriceOracleCaller: GasPriceOracleCaller{contract: contract}, GasPriceOracleTransactor: GasPriceOracleTransactor{contract: contract}, GasPriceOracleFilterer: GasPriceOracleFilterer{contract: contract}}, nil
}

// NewGasPriceOracleCaller creates a new read-only instance of GasPriceOracle, bound to a specific deployed contract.
func NewGasPriceOracleCaller(address common.Address, caller bind.ContractCaller) (*GasPriceOracleCaller, error) {
	contract, err := bindGasPriceOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GasPriceOracleCaller{contract: contract}, nil
}

// NewGasPriceOracleTransactor creates a new write-only instance of GasPriceOracle, bound to a specific deployed contract.
func NewGasPriceOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*GasPriceOracleTransactor, error) {
	contract, err := bindGasPriceOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GasPriceOracleTransactor{contract: contract}, nil
}

// NewGasPriceOracleFilterer creates a new log filterer instance of GasPriceOracle, bound to a specific deployed contract.
func NewGasPriceOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*GasPriceOracleFilterer, error) {
	contract, err := bindGasPriceOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GasPriceOracleFilterer{contract: contract}, nil
}

// bindGasPriceOracle binds a generic wrapper to an already deployed contract.
func bindGasPriceOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GasPriceOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GasPriceOracle *GasPriceOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GasPriceOracle.Contract.GasPriceOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GasPriceOracle *GasPriceOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.GasPriceOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GasPriceOracle *GasPriceOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.GasPriceOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GasPriceOracle *GasPriceOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GasPriceOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GasPriceOracle *GasPriceOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GasPriceOracle *GasPriceOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.contract.Transact(opts, method, params...)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_GasPriceOracle *GasPriceOracleCaller) FUNGIBLEMODULEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GasPriceOracle.contract.Call(opts, &out, "FUNGIBLE_MODULE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_GasPriceOracle *GasPriceOracleSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _GasPriceOracle.Contract.FUNGIBLEMODULEADDRESS(&_GasPriceOracle.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_GasPriceOracle *GasPriceOracleCallerSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _GasPriceOracle.Contract.FUNGIBLEMODULEADDRESS(&_GasPriceOracle.CallOpts)
}

// GasCoinERC4 is a free data retrieval call binding the contract method 0xadbe1596.
//
// Solidity: function gasCoinERC4(uint256 ) view returns(address)
func (_GasPriceOracle *GasPriceOracleCaller) GasCoinERC4(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GasPriceOracle.contract.Call(opts, &out, "gasCoinERC4", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GasCoinERC4 is a free data retrieval call binding the contract method 0xadbe1596.
//
// Solidity: function gasCoinERC4(uint256 ) view returns(address)
func (_GasPriceOracle *GasPriceOracleSession) GasCoinERC4(arg0 *big.Int) (common.Address, error) {
	return _GasPriceOracle.Contract.GasCoinERC4(&_GasPriceOracle.CallOpts, arg0)
}

// GasCoinERC4 is a free data retrieval call binding the contract method 0xadbe1596.
//
// Solidity: function gasCoinERC4(uint256 ) view returns(address)
func (_GasPriceOracle *GasPriceOracleCallerSession) GasCoinERC4(arg0 *big.Int) (common.Address, error) {
	return _GasPriceOracle.Contract.GasCoinERC4(&_GasPriceOracle.CallOpts, arg0)
}

// GasPrice is a free data retrieval call binding the contract method 0x052c72a4.
//
// Solidity: function gasPrice(uint256 ) view returns(uint256)
func (_GasPriceOracle *GasPriceOracleCaller) GasPrice(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GasPriceOracle.contract.Call(opts, &out, "gasPrice", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GasPrice is a free data retrieval call binding the contract method 0x052c72a4.
//
// Solidity: function gasPrice(uint256 ) view returns(uint256)
func (_GasPriceOracle *GasPriceOracleSession) GasPrice(arg0 *big.Int) (*big.Int, error) {
	return _GasPriceOracle.Contract.GasPrice(&_GasPriceOracle.CallOpts, arg0)
}

// GasPrice is a free data retrieval call binding the contract method 0x052c72a4.
//
// Solidity: function gasPrice(uint256 ) view returns(uint256)
func (_GasPriceOracle *GasPriceOracleCallerSession) GasPrice(arg0 *big.Int) (*big.Int, error) {
	return _GasPriceOracle.Contract.GasPrice(&_GasPriceOracle.CallOpts, arg0)
}

// SetGasCoinERC4 is a paid mutator transaction binding the contract method 0xbbb49d4b.
//
// Solidity: function setGasCoinERC4(uint256 chainID, address erc4) returns()
func (_GasPriceOracle *GasPriceOracleTransactor) SetGasCoinERC4(opts *bind.TransactOpts, chainID *big.Int, erc4 common.Address) (*types.Transaction, error) {
	return _GasPriceOracle.contract.Transact(opts, "setGasCoinERC4", chainID, erc4)
}

// SetGasCoinERC4 is a paid mutator transaction binding the contract method 0xbbb49d4b.
//
// Solidity: function setGasCoinERC4(uint256 chainID, address erc4) returns()
func (_GasPriceOracle *GasPriceOracleSession) SetGasCoinERC4(chainID *big.Int, erc4 common.Address) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.SetGasCoinERC4(&_GasPriceOracle.TransactOpts, chainID, erc4)
}

// SetGasCoinERC4 is a paid mutator transaction binding the contract method 0xbbb49d4b.
//
// Solidity: function setGasCoinERC4(uint256 chainID, address erc4) returns()
func (_GasPriceOracle *GasPriceOracleTransactorSession) SetGasCoinERC4(chainID *big.Int, erc4 common.Address) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.SetGasCoinERC4(&_GasPriceOracle.TransactOpts, chainID, erc4)
}

// SetGasPrice is a paid mutator transaction binding the contract method 0xa7cb0507.
//
// Solidity: function setGasPrice(uint256 chainID, uint256 price) returns()
func (_GasPriceOracle *GasPriceOracleTransactor) SetGasPrice(opts *bind.TransactOpts, chainID *big.Int, price *big.Int) (*types.Transaction, error) {
	return _GasPriceOracle.contract.Transact(opts, "setGasPrice", chainID, price)
}

// SetGasPrice is a paid mutator transaction binding the contract method 0xa7cb0507.
//
// Solidity: function setGasPrice(uint256 chainID, uint256 price) returns()
func (_GasPriceOracle *GasPriceOracleSession) SetGasPrice(chainID *big.Int, price *big.Int) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.SetGasPrice(&_GasPriceOracle.TransactOpts, chainID, price)
}

// SetGasPrice is a paid mutator transaction binding the contract method 0xa7cb0507.
//
// Solidity: function setGasPrice(uint256 chainID, uint256 price) returns()
func (_GasPriceOracle *GasPriceOracleTransactorSession) SetGasPrice(chainID *big.Int, price *big.Int) (*types.Transaction, error) {
	return _GasPriceOracle.Contract.SetGasPrice(&_GasPriceOracle.TransactOpts, chainID, price)
}
