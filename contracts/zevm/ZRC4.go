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

// ZRC4MetaData contains all meta data concerning the ZRC4 contract.
var ZRC4MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"chainid_\",\"type\":\"uint256\"},{\"internalType\":\"enumCoinType\",\"name\":\"coinType_\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"systemContractAddress_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasfee\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CHAIN_ID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"COIN_TYPE\",\"outputs\":[{\"internalType\":\"enumCoinType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FUNGIBLE_MODULE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GAS_LIMIT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SYSTEM_CONTRACT_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"updateGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"updateSystemContractAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawGasFee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ZRC4ABI is the input ABI used to generate the binding from.
// Deprecated: Use ZRC4MetaData.ABI instead.
var ZRC4ABI = ZRC4MetaData.ABI

// ZRC4 is an auto generated Go binding around an Ethereum contract.
type ZRC4 struct {
	ZRC4Caller     // Read-only binding to the contract
	ZRC4Transactor // Write-only binding to the contract
	ZRC4Filterer   // Log filterer for contract events
}

// ZRC4Caller is an auto generated read-only Go binding around an Ethereum contract.
type ZRC4Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZRC4Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ZRC4Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZRC4Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZRC4Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZRC4Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZRC4Session struct {
	Contract     *ZRC4             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZRC4CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZRC4CallerSession struct {
	Contract *ZRC4Caller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ZRC4TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZRC4TransactorSession struct {
	Contract     *ZRC4Transactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZRC4Raw is an auto generated low-level Go binding around an Ethereum contract.
type ZRC4Raw struct {
	Contract *ZRC4 // Generic contract binding to access the raw methods on
}

// ZRC4CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZRC4CallerRaw struct {
	Contract *ZRC4Caller // Generic read-only contract binding to access the raw methods on
}

// ZRC4TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZRC4TransactorRaw struct {
	Contract *ZRC4Transactor // Generic write-only contract binding to access the raw methods on
}

// NewZRC4 creates a new instance of ZRC4, bound to a specific deployed contract.
func NewZRC4(address common.Address, backend bind.ContractBackend) (*ZRC4, error) {
	contract, err := bindZRC4(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZRC4{ZRC4Caller: ZRC4Caller{contract: contract}, ZRC4Transactor: ZRC4Transactor{contract: contract}, ZRC4Filterer: ZRC4Filterer{contract: contract}}, nil
}

// NewZRC4Caller creates a new read-only instance of ZRC4, bound to a specific deployed contract.
func NewZRC4Caller(address common.Address, caller bind.ContractCaller) (*ZRC4Caller, error) {
	contract, err := bindZRC4(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZRC4Caller{contract: contract}, nil
}

// NewZRC4Transactor creates a new write-only instance of ZRC4, bound to a specific deployed contract.
func NewZRC4Transactor(address common.Address, transactor bind.ContractTransactor) (*ZRC4Transactor, error) {
	contract, err := bindZRC4(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZRC4Transactor{contract: contract}, nil
}

// NewZRC4Filterer creates a new log filterer instance of ZRC4, bound to a specific deployed contract.
func NewZRC4Filterer(address common.Address, filterer bind.ContractFilterer) (*ZRC4Filterer, error) {
	contract, err := bindZRC4(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZRC4Filterer{contract: contract}, nil
}

// bindZRC4 binds a generic wrapper to an already deployed contract.
func bindZRC4(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ZRC4ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZRC4 *ZRC4Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZRC4.Contract.ZRC4Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZRC4 *ZRC4Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZRC4.Contract.ZRC4Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZRC4 *ZRC4Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZRC4.Contract.ZRC4Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZRC4 *ZRC4CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZRC4.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZRC4 *ZRC4TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZRC4.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZRC4 *ZRC4TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZRC4.Contract.contract.Transact(opts, method, params...)
}

// CHAINID is a free data retrieval call binding the contract method 0x85e1f4d0.
//
// Solidity: function CHAIN_ID() view returns(uint256)
func (_ZRC4 *ZRC4Caller) CHAINID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "CHAIN_ID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CHAINID is a free data retrieval call binding the contract method 0x85e1f4d0.
//
// Solidity: function CHAIN_ID() view returns(uint256)
func (_ZRC4 *ZRC4Session) CHAINID() (*big.Int, error) {
	return _ZRC4.Contract.CHAINID(&_ZRC4.CallOpts)
}

// CHAINID is a free data retrieval call binding the contract method 0x85e1f4d0.
//
// Solidity: function CHAIN_ID() view returns(uint256)
func (_ZRC4 *ZRC4CallerSession) CHAINID() (*big.Int, error) {
	return _ZRC4.Contract.CHAINID(&_ZRC4.CallOpts)
}

// COINTYPE is a free data retrieval call binding the contract method 0xa3413d03.
//
// Solidity: function COIN_TYPE() view returns(uint8)
func (_ZRC4 *ZRC4Caller) COINTYPE(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "COIN_TYPE")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// COINTYPE is a free data retrieval call binding the contract method 0xa3413d03.
//
// Solidity: function COIN_TYPE() view returns(uint8)
func (_ZRC4 *ZRC4Session) COINTYPE() (uint8, error) {
	return _ZRC4.Contract.COINTYPE(&_ZRC4.CallOpts)
}

// COINTYPE is a free data retrieval call binding the contract method 0xa3413d03.
//
// Solidity: function COIN_TYPE() view returns(uint8)
func (_ZRC4 *ZRC4CallerSession) COINTYPE() (uint8, error) {
	return _ZRC4.Contract.COINTYPE(&_ZRC4.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZRC4 *ZRC4Caller) FUNGIBLEMODULEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "FUNGIBLE_MODULE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZRC4 *ZRC4Session) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _ZRC4.Contract.FUNGIBLEMODULEADDRESS(&_ZRC4.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZRC4 *ZRC4CallerSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _ZRC4.Contract.FUNGIBLEMODULEADDRESS(&_ZRC4.CallOpts)
}

// GASLIMIT is a free data retrieval call binding the contract method 0x091d2788.
//
// Solidity: function GAS_LIMIT() view returns(uint256)
func (_ZRC4 *ZRC4Caller) GASLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "GAS_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GASLIMIT is a free data retrieval call binding the contract method 0x091d2788.
//
// Solidity: function GAS_LIMIT() view returns(uint256)
func (_ZRC4 *ZRC4Session) GASLIMIT() (*big.Int, error) {
	return _ZRC4.Contract.GASLIMIT(&_ZRC4.CallOpts)
}

// GASLIMIT is a free data retrieval call binding the contract method 0x091d2788.
//
// Solidity: function GAS_LIMIT() view returns(uint256)
func (_ZRC4 *ZRC4CallerSession) GASLIMIT() (*big.Int, error) {
	return _ZRC4.Contract.GASLIMIT(&_ZRC4.CallOpts)
}

// SYSTEMCONTRACTADDRESS is a free data retrieval call binding the contract method 0xf2441b32.
//
// Solidity: function SYSTEM_CONTRACT_ADDRESS() view returns(address)
func (_ZRC4 *ZRC4Caller) SYSTEMCONTRACTADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "SYSTEM_CONTRACT_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SYSTEMCONTRACTADDRESS is a free data retrieval call binding the contract method 0xf2441b32.
//
// Solidity: function SYSTEM_CONTRACT_ADDRESS() view returns(address)
func (_ZRC4 *ZRC4Session) SYSTEMCONTRACTADDRESS() (common.Address, error) {
	return _ZRC4.Contract.SYSTEMCONTRACTADDRESS(&_ZRC4.CallOpts)
}

// SYSTEMCONTRACTADDRESS is a free data retrieval call binding the contract method 0xf2441b32.
//
// Solidity: function SYSTEM_CONTRACT_ADDRESS() view returns(address)
func (_ZRC4 *ZRC4CallerSession) SYSTEMCONTRACTADDRESS() (common.Address, error) {
	return _ZRC4.Contract.SYSTEMCONTRACTADDRESS(&_ZRC4.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZRC4 *ZRC4Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZRC4 *ZRC4Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZRC4.Contract.Allowance(&_ZRC4.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZRC4 *ZRC4CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZRC4.Contract.Allowance(&_ZRC4.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZRC4 *ZRC4Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZRC4 *ZRC4Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZRC4.Contract.BalanceOf(&_ZRC4.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZRC4 *ZRC4CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZRC4.Contract.BalanceOf(&_ZRC4.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZRC4 *ZRC4Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZRC4 *ZRC4Session) Decimals() (uint8, error) {
	return _ZRC4.Contract.Decimals(&_ZRC4.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZRC4 *ZRC4CallerSession) Decimals() (uint8, error) {
	return _ZRC4.Contract.Decimals(&_ZRC4.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZRC4 *ZRC4Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZRC4 *ZRC4Session) Name() (string, error) {
	return _ZRC4.Contract.Name(&_ZRC4.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZRC4 *ZRC4CallerSession) Name() (string, error) {
	return _ZRC4.Contract.Name(&_ZRC4.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZRC4 *ZRC4Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZRC4 *ZRC4Session) Symbol() (string, error) {
	return _ZRC4.Contract.Symbol(&_ZRC4.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZRC4 *ZRC4CallerSession) Symbol() (string, error) {
	return _ZRC4.Contract.Symbol(&_ZRC4.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZRC4 *ZRC4Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZRC4 *ZRC4Session) TotalSupply() (*big.Int, error) {
	return _ZRC4.Contract.TotalSupply(&_ZRC4.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZRC4 *ZRC4CallerSession) TotalSupply() (*big.Int, error) {
	return _ZRC4.Contract.TotalSupply(&_ZRC4.CallOpts)
}

// WithdrawGasFee is a free data retrieval call binding the contract method 0xd9eeebed.
//
// Solidity: function withdrawGasFee() view returns(address, uint256)
func (_ZRC4 *ZRC4Caller) WithdrawGasFee(opts *bind.CallOpts) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _ZRC4.contract.Call(opts, &out, "withdrawGasFee")

	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// WithdrawGasFee is a free data retrieval call binding the contract method 0xd9eeebed.
//
// Solidity: function withdrawGasFee() view returns(address, uint256)
func (_ZRC4 *ZRC4Session) WithdrawGasFee() (common.Address, *big.Int, error) {
	return _ZRC4.Contract.WithdrawGasFee(&_ZRC4.CallOpts)
}

// WithdrawGasFee is a free data retrieval call binding the contract method 0xd9eeebed.
//
// Solidity: function withdrawGasFee() view returns(address, uint256)
func (_ZRC4 *ZRC4CallerSession) WithdrawGasFee() (common.Address, *big.Int, error) {
	return _ZRC4.Contract.WithdrawGasFee(&_ZRC4.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.Approve(&_ZRC4.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.Approve(&_ZRC4.TransactOpts, spender, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address to, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4Transactor) Deposit(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.contract.Transact(opts, "deposit", to, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address to, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4Session) Deposit(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.Deposit(&_ZRC4.TransactOpts, to, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address to, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4TransactorSession) Deposit(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.Deposit(&_ZRC4.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.Transfer(&_ZRC4.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.Transfer(&_ZRC4.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.TransferFrom(&_ZRC4.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.TransferFrom(&_ZRC4.TransactOpts, sender, recipient, amount)
}

// UpdateGasLimit is a paid mutator transaction binding the contract method 0xf687d12a.
//
// Solidity: function updateGasLimit(uint256 gasLimit) returns()
func (_ZRC4 *ZRC4Transactor) UpdateGasLimit(opts *bind.TransactOpts, gasLimit *big.Int) (*types.Transaction, error) {
	return _ZRC4.contract.Transact(opts, "updateGasLimit", gasLimit)
}

// UpdateGasLimit is a paid mutator transaction binding the contract method 0xf687d12a.
//
// Solidity: function updateGasLimit(uint256 gasLimit) returns()
func (_ZRC4 *ZRC4Session) UpdateGasLimit(gasLimit *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.UpdateGasLimit(&_ZRC4.TransactOpts, gasLimit)
}

// UpdateGasLimit is a paid mutator transaction binding the contract method 0xf687d12a.
//
// Solidity: function updateGasLimit(uint256 gasLimit) returns()
func (_ZRC4 *ZRC4TransactorSession) UpdateGasLimit(gasLimit *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.UpdateGasLimit(&_ZRC4.TransactOpts, gasLimit)
}

// UpdateSystemContractAddress is a paid mutator transaction binding the contract method 0xc835d7cc.
//
// Solidity: function updateSystemContractAddress(address addr) returns()
func (_ZRC4 *ZRC4Transactor) UpdateSystemContractAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ZRC4.contract.Transact(opts, "updateSystemContractAddress", addr)
}

// UpdateSystemContractAddress is a paid mutator transaction binding the contract method 0xc835d7cc.
//
// Solidity: function updateSystemContractAddress(address addr) returns()
func (_ZRC4 *ZRC4Session) UpdateSystemContractAddress(addr common.Address) (*types.Transaction, error) {
	return _ZRC4.Contract.UpdateSystemContractAddress(&_ZRC4.TransactOpts, addr)
}

// UpdateSystemContractAddress is a paid mutator transaction binding the contract method 0xc835d7cc.
//
// Solidity: function updateSystemContractAddress(address addr) returns()
func (_ZRC4 *ZRC4TransactorSession) UpdateSystemContractAddress(addr common.Address) (*types.Transaction, error) {
	return _ZRC4.Contract.UpdateSystemContractAddress(&_ZRC4.TransactOpts, addr)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc7012626.
//
// Solidity: function withdraw(bytes to, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4Transactor) Withdraw(opts *bind.TransactOpts, to []byte, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.contract.Transact(opts, "withdraw", to, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc7012626.
//
// Solidity: function withdraw(bytes to, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4Session) Withdraw(to []byte, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.Withdraw(&_ZRC4.TransactOpts, to, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc7012626.
//
// Solidity: function withdraw(bytes to, uint256 amount) returns(bool)
func (_ZRC4 *ZRC4TransactorSession) Withdraw(to []byte, amount *big.Int) (*types.Transaction, error) {
	return _ZRC4.Contract.Withdraw(&_ZRC4.TransactOpts, to, amount)
}

// ZRC4ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ZRC4 contract.
type ZRC4ApprovalIterator struct {
	Event *ZRC4Approval // Event containing the contract specifics and raw log

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
func (it *ZRC4ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZRC4Approval)
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
		it.Event = new(ZRC4Approval)
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
func (it *ZRC4ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZRC4ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZRC4Approval represents a Approval event raised by the ZRC4 contract.
type ZRC4Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZRC4 *ZRC4Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ZRC4ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZRC4.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ZRC4ApprovalIterator{contract: _ZRC4.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZRC4 *ZRC4Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ZRC4Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZRC4.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZRC4Approval)
				if err := _ZRC4.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZRC4 *ZRC4Filterer) ParseApproval(log types.Log) (*ZRC4Approval, error) {
	event := new(ZRC4Approval)
	if err := _ZRC4.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZRC4DepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the ZRC4 contract.
type ZRC4DepositIterator struct {
	Event *ZRC4Deposit // Event containing the contract specifics and raw log

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
func (it *ZRC4DepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZRC4Deposit)
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
		it.Event = new(ZRC4Deposit)
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
func (it *ZRC4DepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZRC4DepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZRC4Deposit represents a Deposit event raised by the ZRC4 contract.
type ZRC4Deposit struct {
	From  []byte
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab3.
//
// Solidity: event Deposit(bytes from, address indexed to, uint256 value)
func (_ZRC4 *ZRC4Filterer) FilterDeposit(opts *bind.FilterOpts, to []common.Address) (*ZRC4DepositIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZRC4.contract.FilterLogs(opts, "Deposit", toRule)
	if err != nil {
		return nil, err
	}
	return &ZRC4DepositIterator{contract: _ZRC4.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab3.
//
// Solidity: event Deposit(bytes from, address indexed to, uint256 value)
func (_ZRC4 *ZRC4Filterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *ZRC4Deposit, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZRC4.contract.WatchLogs(opts, "Deposit", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZRC4Deposit)
				if err := _ZRC4.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab3.
//
// Solidity: event Deposit(bytes from, address indexed to, uint256 value)
func (_ZRC4 *ZRC4Filterer) ParseDeposit(log types.Log) (*ZRC4Deposit, error) {
	event := new(ZRC4Deposit)
	if err := _ZRC4.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZRC4TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ZRC4 contract.
type ZRC4TransferIterator struct {
	Event *ZRC4Transfer // Event containing the contract specifics and raw log

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
func (it *ZRC4TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZRC4Transfer)
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
		it.Event = new(ZRC4Transfer)
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
func (it *ZRC4TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZRC4TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZRC4Transfer represents a Transfer event raised by the ZRC4 contract.
type ZRC4Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZRC4 *ZRC4Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ZRC4TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZRC4.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ZRC4TransferIterator{contract: _ZRC4.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZRC4 *ZRC4Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ZRC4Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZRC4.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZRC4Transfer)
				if err := _ZRC4.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZRC4 *ZRC4Filterer) ParseTransfer(log types.Log) (*ZRC4Transfer, error) {
	event := new(ZRC4Transfer)
	if err := _ZRC4.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZRC4WithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the ZRC4 contract.
type ZRC4WithdrawalIterator struct {
	Event *ZRC4Withdrawal // Event containing the contract specifics and raw log

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
func (it *ZRC4WithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZRC4Withdrawal)
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
		it.Event = new(ZRC4Withdrawal)
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
func (it *ZRC4WithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZRC4WithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZRC4Withdrawal represents a Withdrawal event raised by the ZRC4 contract.
type ZRC4Withdrawal struct {
	From   common.Address
	To     []byte
	Value  *big.Int
	Gasfee *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x6cf92e651ecb67795ff01e9bce83f166b16214b1dfb2669dd73c12b767ba4f92.
//
// Solidity: event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee)
func (_ZRC4 *ZRC4Filterer) FilterWithdrawal(opts *bind.FilterOpts, from []common.Address) (*ZRC4WithdrawalIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _ZRC4.contract.FilterLogs(opts, "Withdrawal", fromRule)
	if err != nil {
		return nil, err
	}
	return &ZRC4WithdrawalIterator{contract: _ZRC4.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x6cf92e651ecb67795ff01e9bce83f166b16214b1dfb2669dd73c12b767ba4f92.
//
// Solidity: event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee)
func (_ZRC4 *ZRC4Filterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *ZRC4Withdrawal, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _ZRC4.contract.WatchLogs(opts, "Withdrawal", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZRC4Withdrawal)
				if err := _ZRC4.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x6cf92e651ecb67795ff01e9bce83f166b16214b1dfb2669dd73c12b767ba4f92.
//
// Solidity: event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee)
func (_ZRC4 *ZRC4Filterer) ParseWithdrawal(log types.Log) (*ZRC4Withdrawal, error) {
	event := new(ZRC4Withdrawal)
	if err := _ZRC4.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
