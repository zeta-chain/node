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

// SystemContractMetaData contains all meta data concerning the SystemContract contract.
var SystemContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wzeta_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"uniswapv2Factory_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"SetGasCoin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"SetGasPrice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"SetGasZetaPool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"SystemContractDeployed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc4\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"DepositAndCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FUNGIBLE_MODULE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasCoinERC4\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"gasZetaPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"zrc4\",\"type\":\"address\"}],\"name\":\"setGasCoinERC4\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"setGasPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"erc20\",\"type\":\"address\"}],\"name\":\"setGasZetaPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setUniswapv2FactoryAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setWZETAContractAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"uniswapv2FactoryAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"factory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"}],\"name\":\"uniswapv2PairFor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wzetaContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// SystemContractABI is the input ABI used to generate the binding from.
// Deprecated: Use SystemContractMetaData.ABI instead.
var SystemContractABI = SystemContractMetaData.ABI

// SystemContract is an auto generated Go binding around an Ethereum contract.
type SystemContract struct {
	SystemContractCaller     // Read-only binding to the contract
	SystemContractTransactor // Write-only binding to the contract
	SystemContractFilterer   // Log filterer for contract events
}

// SystemContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type SystemContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SystemContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SystemContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SystemContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SystemContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SystemContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SystemContractSession struct {
	Contract     *SystemContract   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SystemContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SystemContractCallerSession struct {
	Contract *SystemContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// SystemContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SystemContractTransactorSession struct {
	Contract     *SystemContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SystemContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type SystemContractRaw struct {
	Contract *SystemContract // Generic contract binding to access the raw methods on
}

// SystemContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SystemContractCallerRaw struct {
	Contract *SystemContractCaller // Generic read-only contract binding to access the raw methods on
}

// SystemContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SystemContractTransactorRaw struct {
	Contract *SystemContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSystemContract creates a new instance of SystemContract, bound to a specific deployed contract.
func NewSystemContract(address common.Address, backend bind.ContractBackend) (*SystemContract, error) {
	contract, err := bindSystemContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SystemContract{SystemContractCaller: SystemContractCaller{contract: contract}, SystemContractTransactor: SystemContractTransactor{contract: contract}, SystemContractFilterer: SystemContractFilterer{contract: contract}}, nil
}

// NewSystemContractCaller creates a new read-only instance of SystemContract, bound to a specific deployed contract.
func NewSystemContractCaller(address common.Address, caller bind.ContractCaller) (*SystemContractCaller, error) {
	contract, err := bindSystemContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SystemContractCaller{contract: contract}, nil
}

// NewSystemContractTransactor creates a new write-only instance of SystemContract, bound to a specific deployed contract.
func NewSystemContractTransactor(address common.Address, transactor bind.ContractTransactor) (*SystemContractTransactor, error) {
	contract, err := bindSystemContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SystemContractTransactor{contract: contract}, nil
}

// NewSystemContractFilterer creates a new log filterer instance of SystemContract, bound to a specific deployed contract.
func NewSystemContractFilterer(address common.Address, filterer bind.ContractFilterer) (*SystemContractFilterer, error) {
	contract, err := bindSystemContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SystemContractFilterer{contract: contract}, nil
}

// bindSystemContract binds a generic wrapper to an already deployed contract.
func bindSystemContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SystemContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SystemContract *SystemContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SystemContract.Contract.SystemContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SystemContract *SystemContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SystemContract.Contract.SystemContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SystemContract *SystemContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SystemContract.Contract.SystemContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SystemContract *SystemContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SystemContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SystemContract *SystemContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SystemContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SystemContract *SystemContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SystemContract.Contract.contract.Transact(opts, method, params...)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_SystemContract *SystemContractCaller) FUNGIBLEMODULEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SystemContract.contract.Call(opts, &out, "FUNGIBLE_MODULE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_SystemContract *SystemContractSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _SystemContract.Contract.FUNGIBLEMODULEADDRESS(&_SystemContract.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_SystemContract *SystemContractCallerSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _SystemContract.Contract.FUNGIBLEMODULEADDRESS(&_SystemContract.CallOpts)
}

// GasCoinERC4 is a free data retrieval call binding the contract method 0xadbe1596.
//
// Solidity: function gasCoinERC4(uint256 ) view returns(address)
func (_SystemContract *SystemContractCaller) GasCoinERC4(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SystemContract.contract.Call(opts, &out, "gasCoinERC4", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GasCoinERC4 is a free data retrieval call binding the contract method 0xadbe1596.
//
// Solidity: function gasCoinERC4(uint256 ) view returns(address)
func (_SystemContract *SystemContractSession) GasCoinERC4(arg0 *big.Int) (common.Address, error) {
	return _SystemContract.Contract.GasCoinERC4(&_SystemContract.CallOpts, arg0)
}

// GasCoinERC4 is a free data retrieval call binding the contract method 0xadbe1596.
//
// Solidity: function gasCoinERC4(uint256 ) view returns(address)
func (_SystemContract *SystemContractCallerSession) GasCoinERC4(arg0 *big.Int) (common.Address, error) {
	return _SystemContract.Contract.GasCoinERC4(&_SystemContract.CallOpts, arg0)
}

// GasPrice is a free data retrieval call binding the contract method 0x052c72a4.
//
// Solidity: function gasPrice(uint256 ) view returns(uint256)
func (_SystemContract *SystemContractCaller) GasPrice(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _SystemContract.contract.Call(opts, &out, "gasPrice", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GasPrice is a free data retrieval call binding the contract method 0x052c72a4.
//
// Solidity: function gasPrice(uint256 ) view returns(uint256)
func (_SystemContract *SystemContractSession) GasPrice(arg0 *big.Int) (*big.Int, error) {
	return _SystemContract.Contract.GasPrice(&_SystemContract.CallOpts, arg0)
}

// GasPrice is a free data retrieval call binding the contract method 0x052c72a4.
//
// Solidity: function gasPrice(uint256 ) view returns(uint256)
func (_SystemContract *SystemContractCallerSession) GasPrice(arg0 *big.Int) (*big.Int, error) {
	return _SystemContract.Contract.GasPrice(&_SystemContract.CallOpts, arg0)
}

// GasZetaPool is a free data retrieval call binding the contract method 0x05a6ee8f.
//
// Solidity: function gasZetaPool(uint256 ) view returns(address)
func (_SystemContract *SystemContractCaller) GasZetaPool(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SystemContract.contract.Call(opts, &out, "gasZetaPool", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GasZetaPool is a free data retrieval call binding the contract method 0x05a6ee8f.
//
// Solidity: function gasZetaPool(uint256 ) view returns(address)
func (_SystemContract *SystemContractSession) GasZetaPool(arg0 *big.Int) (common.Address, error) {
	return _SystemContract.Contract.GasZetaPool(&_SystemContract.CallOpts, arg0)
}

// GasZetaPool is a free data retrieval call binding the contract method 0x05a6ee8f.
//
// Solidity: function gasZetaPool(uint256 ) view returns(address)
func (_SystemContract *SystemContractCallerSession) GasZetaPool(arg0 *big.Int) (common.Address, error) {
	return _SystemContract.Contract.GasZetaPool(&_SystemContract.CallOpts, arg0)
}

// Uniswapv2FactoryAddress is a free data retrieval call binding the contract method 0xd936a012.
//
// Solidity: function uniswapv2FactoryAddress() view returns(address)
func (_SystemContract *SystemContractCaller) Uniswapv2FactoryAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SystemContract.contract.Call(opts, &out, "uniswapv2FactoryAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Uniswapv2FactoryAddress is a free data retrieval call binding the contract method 0xd936a012.
//
// Solidity: function uniswapv2FactoryAddress() view returns(address)
func (_SystemContract *SystemContractSession) Uniswapv2FactoryAddress() (common.Address, error) {
	return _SystemContract.Contract.Uniswapv2FactoryAddress(&_SystemContract.CallOpts)
}

// Uniswapv2FactoryAddress is a free data retrieval call binding the contract method 0xd936a012.
//
// Solidity: function uniswapv2FactoryAddress() view returns(address)
func (_SystemContract *SystemContractCallerSession) Uniswapv2FactoryAddress() (common.Address, error) {
	return _SystemContract.Contract.Uniswapv2FactoryAddress(&_SystemContract.CallOpts)
}

// Uniswapv2PairFor is a free data retrieval call binding the contract method 0xc63585cc.
//
// Solidity: function uniswapv2PairFor(address factory, address tokenA, address tokenB) pure returns(address pair)
func (_SystemContract *SystemContractCaller) Uniswapv2PairFor(opts *bind.CallOpts, factory common.Address, tokenA common.Address, tokenB common.Address) (common.Address, error) {
	var out []interface{}
	err := _SystemContract.contract.Call(opts, &out, "uniswapv2PairFor", factory, tokenA, tokenB)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Uniswapv2PairFor is a free data retrieval call binding the contract method 0xc63585cc.
//
// Solidity: function uniswapv2PairFor(address factory, address tokenA, address tokenB) pure returns(address pair)
func (_SystemContract *SystemContractSession) Uniswapv2PairFor(factory common.Address, tokenA common.Address, tokenB common.Address) (common.Address, error) {
	return _SystemContract.Contract.Uniswapv2PairFor(&_SystemContract.CallOpts, factory, tokenA, tokenB)
}

// Uniswapv2PairFor is a free data retrieval call binding the contract method 0xc63585cc.
//
// Solidity: function uniswapv2PairFor(address factory, address tokenA, address tokenB) pure returns(address pair)
func (_SystemContract *SystemContractCallerSession) Uniswapv2PairFor(factory common.Address, tokenA common.Address, tokenB common.Address) (common.Address, error) {
	return _SystemContract.Contract.Uniswapv2PairFor(&_SystemContract.CallOpts, factory, tokenA, tokenB)
}

// WzetaContractAddress is a free data retrieval call binding the contract method 0x6719e404.
//
// Solidity: function wzetaContractAddress() view returns(address)
func (_SystemContract *SystemContractCaller) WzetaContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SystemContract.contract.Call(opts, &out, "wzetaContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WzetaContractAddress is a free data retrieval call binding the contract method 0x6719e404.
//
// Solidity: function wzetaContractAddress() view returns(address)
func (_SystemContract *SystemContractSession) WzetaContractAddress() (common.Address, error) {
	return _SystemContract.Contract.WzetaContractAddress(&_SystemContract.CallOpts)
}

// WzetaContractAddress is a free data retrieval call binding the contract method 0x6719e404.
//
// Solidity: function wzetaContractAddress() view returns(address)
func (_SystemContract *SystemContractCallerSession) WzetaContractAddress() (common.Address, error) {
	return _SystemContract.Contract.WzetaContractAddress(&_SystemContract.CallOpts)
}

// DepositAndCall is a paid mutator transaction binding the contract method 0x12cc2a3b.
//
// Solidity: function DepositAndCall(address zrc4, uint256 amount, address target, bytes message) returns()
func (_SystemContract *SystemContractTransactor) DepositAndCall(opts *bind.TransactOpts, zrc4 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _SystemContract.contract.Transact(opts, "DepositAndCall", zrc4, amount, target, message)
}

// DepositAndCall is a paid mutator transaction binding the contract method 0x12cc2a3b.
//
// Solidity: function DepositAndCall(address zrc4, uint256 amount, address target, bytes message) returns()
func (_SystemContract *SystemContractSession) DepositAndCall(zrc4 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _SystemContract.Contract.DepositAndCall(&_SystemContract.TransactOpts, zrc4, amount, target, message)
}

// DepositAndCall is a paid mutator transaction binding the contract method 0x12cc2a3b.
//
// Solidity: function DepositAndCall(address zrc4, uint256 amount, address target, bytes message) returns()
func (_SystemContract *SystemContractTransactorSession) DepositAndCall(zrc4 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _SystemContract.Contract.DepositAndCall(&_SystemContract.TransactOpts, zrc4, amount, target, message)
}

// SetGasCoinERC4 is a paid mutator transaction binding the contract method 0xbbb49d4b.
//
// Solidity: function setGasCoinERC4(uint256 chainID, address zrc4) returns()
func (_SystemContract *SystemContractTransactor) SetGasCoinERC4(opts *bind.TransactOpts, chainID *big.Int, zrc4 common.Address) (*types.Transaction, error) {
	return _SystemContract.contract.Transact(opts, "setGasCoinERC4", chainID, zrc4)
}

// SetGasCoinERC4 is a paid mutator transaction binding the contract method 0xbbb49d4b.
//
// Solidity: function setGasCoinERC4(uint256 chainID, address zrc4) returns()
func (_SystemContract *SystemContractSession) SetGasCoinERC4(chainID *big.Int, zrc4 common.Address) (*types.Transaction, error) {
	return _SystemContract.Contract.SetGasCoinERC4(&_SystemContract.TransactOpts, chainID, zrc4)
}

// SetGasCoinERC4 is a paid mutator transaction binding the contract method 0xbbb49d4b.
//
// Solidity: function setGasCoinERC4(uint256 chainID, address zrc4) returns()
func (_SystemContract *SystemContractTransactorSession) SetGasCoinERC4(chainID *big.Int, zrc4 common.Address) (*types.Transaction, error) {
	return _SystemContract.Contract.SetGasCoinERC4(&_SystemContract.TransactOpts, chainID, zrc4)
}

// SetGasPrice is a paid mutator transaction binding the contract method 0xa7cb0507.
//
// Solidity: function setGasPrice(uint256 chainID, uint256 price) returns()
func (_SystemContract *SystemContractTransactor) SetGasPrice(opts *bind.TransactOpts, chainID *big.Int, price *big.Int) (*types.Transaction, error) {
	return _SystemContract.contract.Transact(opts, "setGasPrice", chainID, price)
}

// SetGasPrice is a paid mutator transaction binding the contract method 0xa7cb0507.
//
// Solidity: function setGasPrice(uint256 chainID, uint256 price) returns()
func (_SystemContract *SystemContractSession) SetGasPrice(chainID *big.Int, price *big.Int) (*types.Transaction, error) {
	return _SystemContract.Contract.SetGasPrice(&_SystemContract.TransactOpts, chainID, price)
}

// SetGasPrice is a paid mutator transaction binding the contract method 0xa7cb0507.
//
// Solidity: function setGasPrice(uint256 chainID, uint256 price) returns()
func (_SystemContract *SystemContractTransactorSession) SetGasPrice(chainID *big.Int, price *big.Int) (*types.Transaction, error) {
	return _SystemContract.Contract.SetGasPrice(&_SystemContract.TransactOpts, chainID, price)
}

// SetGasZetaPool is a paid mutator transaction binding the contract method 0x91dd645f.
//
// Solidity: function setGasZetaPool(uint256 chainID, address erc20) returns()
func (_SystemContract *SystemContractTransactor) SetGasZetaPool(opts *bind.TransactOpts, chainID *big.Int, erc20 common.Address) (*types.Transaction, error) {
	return _SystemContract.contract.Transact(opts, "setGasZetaPool", chainID, erc20)
}

// SetGasZetaPool is a paid mutator transaction binding the contract method 0x91dd645f.
//
// Solidity: function setGasZetaPool(uint256 chainID, address erc20) returns()
func (_SystemContract *SystemContractSession) SetGasZetaPool(chainID *big.Int, erc20 common.Address) (*types.Transaction, error) {
	return _SystemContract.Contract.SetGasZetaPool(&_SystemContract.TransactOpts, chainID, erc20)
}

// SetGasZetaPool is a paid mutator transaction binding the contract method 0x91dd645f.
//
// Solidity: function setGasZetaPool(uint256 chainID, address erc20) returns()
func (_SystemContract *SystemContractTransactorSession) SetGasZetaPool(chainID *big.Int, erc20 common.Address) (*types.Transaction, error) {
	return _SystemContract.Contract.SetGasZetaPool(&_SystemContract.TransactOpts, chainID, erc20)
}

// SetUniswapv2FactoryAddress is a paid mutator transaction binding the contract method 0x7c24b45c.
//
// Solidity: function setUniswapv2FactoryAddress(address addr) returns()
func (_SystemContract *SystemContractTransactor) SetUniswapv2FactoryAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _SystemContract.contract.Transact(opts, "setUniswapv2FactoryAddress", addr)
}

// SetUniswapv2FactoryAddress is a paid mutator transaction binding the contract method 0x7c24b45c.
//
// Solidity: function setUniswapv2FactoryAddress(address addr) returns()
func (_SystemContract *SystemContractSession) SetUniswapv2FactoryAddress(addr common.Address) (*types.Transaction, error) {
	return _SystemContract.Contract.SetUniswapv2FactoryAddress(&_SystemContract.TransactOpts, addr)
}

// SetUniswapv2FactoryAddress is a paid mutator transaction binding the contract method 0x7c24b45c.
//
// Solidity: function setUniswapv2FactoryAddress(address addr) returns()
func (_SystemContract *SystemContractTransactorSession) SetUniswapv2FactoryAddress(addr common.Address) (*types.Transaction, error) {
	return _SystemContract.Contract.SetUniswapv2FactoryAddress(&_SystemContract.TransactOpts, addr)
}

// SetWZETAContractAddress is a paid mutator transaction binding the contract method 0x97770dff.
//
// Solidity: function setWZETAContractAddress(address addr) returns()
func (_SystemContract *SystemContractTransactor) SetWZETAContractAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _SystemContract.contract.Transact(opts, "setWZETAContractAddress", addr)
}

// SetWZETAContractAddress is a paid mutator transaction binding the contract method 0x97770dff.
//
// Solidity: function setWZETAContractAddress(address addr) returns()
func (_SystemContract *SystemContractSession) SetWZETAContractAddress(addr common.Address) (*types.Transaction, error) {
	return _SystemContract.Contract.SetWZETAContractAddress(&_SystemContract.TransactOpts, addr)
}

// SetWZETAContractAddress is a paid mutator transaction binding the contract method 0x97770dff.
//
// Solidity: function setWZETAContractAddress(address addr) returns()
func (_SystemContract *SystemContractTransactorSession) SetWZETAContractAddress(addr common.Address) (*types.Transaction, error) {
	return _SystemContract.Contract.SetWZETAContractAddress(&_SystemContract.TransactOpts, addr)
}

// SystemContractSetGasCoinIterator is returned from FilterSetGasCoin and is used to iterate over the raw logs and unpacked data for SetGasCoin events raised by the SystemContract contract.
type SystemContractSetGasCoinIterator struct {
	Event *SystemContractSetGasCoin // Event containing the contract specifics and raw log

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
func (it *SystemContractSetGasCoinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SystemContractSetGasCoin)
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
		it.Event = new(SystemContractSetGasCoin)
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
func (it *SystemContractSetGasCoinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SystemContractSetGasCoinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SystemContractSetGasCoin represents a SetGasCoin event raised by the SystemContract contract.
type SystemContractSetGasCoin struct {
	Arg0 *big.Int
	Arg1 common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetGasCoin is a free log retrieval operation binding the contract event 0xd1b36d30f6248e97c473b4d1348ca164a4ef6759022f54a58ec200326c39c45d.
//
// Solidity: event SetGasCoin(uint256 arg0, address arg1)
func (_SystemContract *SystemContractFilterer) FilterSetGasCoin(opts *bind.FilterOpts) (*SystemContractSetGasCoinIterator, error) {

	logs, sub, err := _SystemContract.contract.FilterLogs(opts, "SetGasCoin")
	if err != nil {
		return nil, err
	}
	return &SystemContractSetGasCoinIterator{contract: _SystemContract.contract, event: "SetGasCoin", logs: logs, sub: sub}, nil
}

// WatchSetGasCoin is a free log subscription operation binding the contract event 0xd1b36d30f6248e97c473b4d1348ca164a4ef6759022f54a58ec200326c39c45d.
//
// Solidity: event SetGasCoin(uint256 arg0, address arg1)
func (_SystemContract *SystemContractFilterer) WatchSetGasCoin(opts *bind.WatchOpts, sink chan<- *SystemContractSetGasCoin) (event.Subscription, error) {

	logs, sub, err := _SystemContract.contract.WatchLogs(opts, "SetGasCoin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SystemContractSetGasCoin)
				if err := _SystemContract.contract.UnpackLog(event, "SetGasCoin", log); err != nil {
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

// ParseSetGasCoin is a log parse operation binding the contract event 0xd1b36d30f6248e97c473b4d1348ca164a4ef6759022f54a58ec200326c39c45d.
//
// Solidity: event SetGasCoin(uint256 arg0, address arg1)
func (_SystemContract *SystemContractFilterer) ParseSetGasCoin(log types.Log) (*SystemContractSetGasCoin, error) {
	event := new(SystemContractSetGasCoin)
	if err := _SystemContract.contract.UnpackLog(event, "SetGasCoin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SystemContractSetGasPriceIterator is returned from FilterSetGasPrice and is used to iterate over the raw logs and unpacked data for SetGasPrice events raised by the SystemContract contract.
type SystemContractSetGasPriceIterator struct {
	Event *SystemContractSetGasPrice // Event containing the contract specifics and raw log

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
func (it *SystemContractSetGasPriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SystemContractSetGasPrice)
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
		it.Event = new(SystemContractSetGasPrice)
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
func (it *SystemContractSetGasPriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SystemContractSetGasPriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SystemContractSetGasPrice represents a SetGasPrice event raised by the SystemContract contract.
type SystemContractSetGasPrice struct {
	Arg0 *big.Int
	Arg1 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetGasPrice is a free log retrieval operation binding the contract event 0x49f492222906ac486c3c1401fa545626df1f0c0e5a77a05597ea2ed66af9850d.
//
// Solidity: event SetGasPrice(uint256 arg0, uint256 arg1)
func (_SystemContract *SystemContractFilterer) FilterSetGasPrice(opts *bind.FilterOpts) (*SystemContractSetGasPriceIterator, error) {

	logs, sub, err := _SystemContract.contract.FilterLogs(opts, "SetGasPrice")
	if err != nil {
		return nil, err
	}
	return &SystemContractSetGasPriceIterator{contract: _SystemContract.contract, event: "SetGasPrice", logs: logs, sub: sub}, nil
}

// WatchSetGasPrice is a free log subscription operation binding the contract event 0x49f492222906ac486c3c1401fa545626df1f0c0e5a77a05597ea2ed66af9850d.
//
// Solidity: event SetGasPrice(uint256 arg0, uint256 arg1)
func (_SystemContract *SystemContractFilterer) WatchSetGasPrice(opts *bind.WatchOpts, sink chan<- *SystemContractSetGasPrice) (event.Subscription, error) {

	logs, sub, err := _SystemContract.contract.WatchLogs(opts, "SetGasPrice")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SystemContractSetGasPrice)
				if err := _SystemContract.contract.UnpackLog(event, "SetGasPrice", log); err != nil {
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

// ParseSetGasPrice is a log parse operation binding the contract event 0x49f492222906ac486c3c1401fa545626df1f0c0e5a77a05597ea2ed66af9850d.
//
// Solidity: event SetGasPrice(uint256 arg0, uint256 arg1)
func (_SystemContract *SystemContractFilterer) ParseSetGasPrice(log types.Log) (*SystemContractSetGasPrice, error) {
	event := new(SystemContractSetGasPrice)
	if err := _SystemContract.contract.UnpackLog(event, "SetGasPrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SystemContractSetGasZetaPoolIterator is returned from FilterSetGasZetaPool and is used to iterate over the raw logs and unpacked data for SetGasZetaPool events raised by the SystemContract contract.
type SystemContractSetGasZetaPoolIterator struct {
	Event *SystemContractSetGasZetaPool // Event containing the contract specifics and raw log

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
func (it *SystemContractSetGasZetaPoolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SystemContractSetGasZetaPool)
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
		it.Event = new(SystemContractSetGasZetaPool)
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
func (it *SystemContractSetGasZetaPoolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SystemContractSetGasZetaPoolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SystemContractSetGasZetaPool represents a SetGasZetaPool event raised by the SystemContract contract.
type SystemContractSetGasZetaPool struct {
	Arg0 *big.Int
	Arg1 common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetGasZetaPool is a free log retrieval operation binding the contract event 0x0ecec485166da6139b13bb7e033e9446e2d35348e80ebf1180d4afe2dba1704e.
//
// Solidity: event SetGasZetaPool(uint256 arg0, address arg1)
func (_SystemContract *SystemContractFilterer) FilterSetGasZetaPool(opts *bind.FilterOpts) (*SystemContractSetGasZetaPoolIterator, error) {

	logs, sub, err := _SystemContract.contract.FilterLogs(opts, "SetGasZetaPool")
	if err != nil {
		return nil, err
	}
	return &SystemContractSetGasZetaPoolIterator{contract: _SystemContract.contract, event: "SetGasZetaPool", logs: logs, sub: sub}, nil
}

// WatchSetGasZetaPool is a free log subscription operation binding the contract event 0x0ecec485166da6139b13bb7e033e9446e2d35348e80ebf1180d4afe2dba1704e.
//
// Solidity: event SetGasZetaPool(uint256 arg0, address arg1)
func (_SystemContract *SystemContractFilterer) WatchSetGasZetaPool(opts *bind.WatchOpts, sink chan<- *SystemContractSetGasZetaPool) (event.Subscription, error) {

	logs, sub, err := _SystemContract.contract.WatchLogs(opts, "SetGasZetaPool")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SystemContractSetGasZetaPool)
				if err := _SystemContract.contract.UnpackLog(event, "SetGasZetaPool", log); err != nil {
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

// ParseSetGasZetaPool is a log parse operation binding the contract event 0x0ecec485166da6139b13bb7e033e9446e2d35348e80ebf1180d4afe2dba1704e.
//
// Solidity: event SetGasZetaPool(uint256 arg0, address arg1)
func (_SystemContract *SystemContractFilterer) ParseSetGasZetaPool(log types.Log) (*SystemContractSetGasZetaPool, error) {
	event := new(SystemContractSetGasZetaPool)
	if err := _SystemContract.contract.UnpackLog(event, "SetGasZetaPool", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SystemContractSystemContractDeployedIterator is returned from FilterSystemContractDeployed and is used to iterate over the raw logs and unpacked data for SystemContractDeployed events raised by the SystemContract contract.
type SystemContractSystemContractDeployedIterator struct {
	Event *SystemContractSystemContractDeployed // Event containing the contract specifics and raw log

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
func (it *SystemContractSystemContractDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SystemContractSystemContractDeployed)
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
		it.Event = new(SystemContractSystemContractDeployed)
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
func (it *SystemContractSystemContractDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SystemContractSystemContractDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SystemContractSystemContractDeployed represents a SystemContractDeployed event raised by the SystemContract contract.
type SystemContractSystemContractDeployed struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterSystemContractDeployed is a free log retrieval operation binding the contract event 0x80699e81136d69cb8367ad52a994e25c722a86da654b561d0c14b61a777e7ac5.
//
// Solidity: event SystemContractDeployed()
func (_SystemContract *SystemContractFilterer) FilterSystemContractDeployed(opts *bind.FilterOpts) (*SystemContractSystemContractDeployedIterator, error) {

	logs, sub, err := _SystemContract.contract.FilterLogs(opts, "SystemContractDeployed")
	if err != nil {
		return nil, err
	}
	return &SystemContractSystemContractDeployedIterator{contract: _SystemContract.contract, event: "SystemContractDeployed", logs: logs, sub: sub}, nil
}

// WatchSystemContractDeployed is a free log subscription operation binding the contract event 0x80699e81136d69cb8367ad52a994e25c722a86da654b561d0c14b61a777e7ac5.
//
// Solidity: event SystemContractDeployed()
func (_SystemContract *SystemContractFilterer) WatchSystemContractDeployed(opts *bind.WatchOpts, sink chan<- *SystemContractSystemContractDeployed) (event.Subscription, error) {

	logs, sub, err := _SystemContract.contract.WatchLogs(opts, "SystemContractDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SystemContractSystemContractDeployed)
				if err := _SystemContract.contract.UnpackLog(event, "SystemContractDeployed", log); err != nil {
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

// ParseSystemContractDeployed is a log parse operation binding the contract event 0x80699e81136d69cb8367ad52a994e25c722a86da654b561d0c14b61a777e7ac5.
//
// Solidity: event SystemContractDeployed()
func (_SystemContract *SystemContractFilterer) ParseSystemContractDeployed(log types.Log) (*SystemContractSystemContractDeployed, error) {
	event := new(SystemContractSystemContractDeployed)
	if err := _SystemContract.contract.UnpackLog(event, "SystemContractDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
