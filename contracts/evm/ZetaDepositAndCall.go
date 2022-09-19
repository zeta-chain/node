// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package evm

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

// ZetaDepositAndCallMetaData contains all meta data concerning the ZetaDepositAndCall contract.
var ZetaDepositAndCallMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fungibleModule\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc4\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"DepositAndCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FUNGIBLE_MODULE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ZetaDepositAndCallABI is the input ABI used to generate the binding from.
// Deprecated: Use ZetaDepositAndCallMetaData.ABI instead.
var ZetaDepositAndCallABI = ZetaDepositAndCallMetaData.ABI

// ZetaDepositAndCall is an auto generated Go binding around an Ethereum contract.
type ZetaDepositAndCall struct {
	ZetaDepositAndCallCaller     // Read-only binding to the contract
	ZetaDepositAndCallTransactor // Write-only binding to the contract
	ZetaDepositAndCallFilterer   // Log filterer for contract events
}

// ZetaDepositAndCallCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZetaDepositAndCallCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZetaDepositAndCallTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZetaDepositAndCallTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZetaDepositAndCallFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZetaDepositAndCallFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZetaDepositAndCallSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZetaDepositAndCallSession struct {
	Contract     *ZetaDepositAndCall // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ZetaDepositAndCallCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZetaDepositAndCallCallerSession struct {
	Contract *ZetaDepositAndCallCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ZetaDepositAndCallTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZetaDepositAndCallTransactorSession struct {
	Contract     *ZetaDepositAndCallTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ZetaDepositAndCallRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZetaDepositAndCallRaw struct {
	Contract *ZetaDepositAndCall // Generic contract binding to access the raw methods on
}

// ZetaDepositAndCallCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZetaDepositAndCallCallerRaw struct {
	Contract *ZetaDepositAndCallCaller // Generic read-only contract binding to access the raw methods on
}

// ZetaDepositAndCallTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZetaDepositAndCallTransactorRaw struct {
	Contract *ZetaDepositAndCallTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZetaDepositAndCall creates a new instance of ZetaDepositAndCall, bound to a specific deployed contract.
func NewZetaDepositAndCall(address common.Address, backend bind.ContractBackend) (*ZetaDepositAndCall, error) {
	contract, err := bindZetaDepositAndCall(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZetaDepositAndCall{ZetaDepositAndCallCaller: ZetaDepositAndCallCaller{contract: contract}, ZetaDepositAndCallTransactor: ZetaDepositAndCallTransactor{contract: contract}, ZetaDepositAndCallFilterer: ZetaDepositAndCallFilterer{contract: contract}}, nil
}

// NewZetaDepositAndCallCaller creates a new read-only instance of ZetaDepositAndCall, bound to a specific deployed contract.
func NewZetaDepositAndCallCaller(address common.Address, caller bind.ContractCaller) (*ZetaDepositAndCallCaller, error) {
	contract, err := bindZetaDepositAndCall(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZetaDepositAndCallCaller{contract: contract}, nil
}

// NewZetaDepositAndCallTransactor creates a new write-only instance of ZetaDepositAndCall, bound to a specific deployed contract.
func NewZetaDepositAndCallTransactor(address common.Address, transactor bind.ContractTransactor) (*ZetaDepositAndCallTransactor, error) {
	contract, err := bindZetaDepositAndCall(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZetaDepositAndCallTransactor{contract: contract}, nil
}

// NewZetaDepositAndCallFilterer creates a new log filterer instance of ZetaDepositAndCall, bound to a specific deployed contract.
func NewZetaDepositAndCallFilterer(address common.Address, filterer bind.ContractFilterer) (*ZetaDepositAndCallFilterer, error) {
	contract, err := bindZetaDepositAndCall(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZetaDepositAndCallFilterer{contract: contract}, nil
}

// bindZetaDepositAndCall binds a generic wrapper to an already deployed contract.
func bindZetaDepositAndCall(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ZetaDepositAndCallABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZetaDepositAndCall *ZetaDepositAndCallRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZetaDepositAndCall.Contract.ZetaDepositAndCallCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZetaDepositAndCall *ZetaDepositAndCallRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZetaDepositAndCall.Contract.ZetaDepositAndCallTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZetaDepositAndCall *ZetaDepositAndCallRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZetaDepositAndCall.Contract.ZetaDepositAndCallTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZetaDepositAndCall *ZetaDepositAndCallCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZetaDepositAndCall.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZetaDepositAndCall *ZetaDepositAndCallTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZetaDepositAndCall.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZetaDepositAndCall *ZetaDepositAndCallTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZetaDepositAndCall.Contract.contract.Transact(opts, method, params...)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZetaDepositAndCall *ZetaDepositAndCallCaller) FUNGIBLEMODULEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZetaDepositAndCall.contract.Call(opts, &out, "FUNGIBLE_MODULE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZetaDepositAndCall *ZetaDepositAndCallSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _ZetaDepositAndCall.Contract.FUNGIBLEMODULEADDRESS(&_ZetaDepositAndCall.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZetaDepositAndCall *ZetaDepositAndCallCallerSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _ZetaDepositAndCall.Contract.FUNGIBLEMODULEADDRESS(&_ZetaDepositAndCall.CallOpts)
}

// DepositAndCall is a paid mutator transaction binding the contract method 0x12cc2a3b.
//
// Solidity: function DepositAndCall(address zrc4, uint256 amount, address target, bytes message) returns()
func (_ZetaDepositAndCall *ZetaDepositAndCallTransactor) DepositAndCall(opts *bind.TransactOpts, zrc4 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _ZetaDepositAndCall.contract.Transact(opts, "DepositAndCall", zrc4, amount, target, message)
}

// DepositAndCall is a paid mutator transaction binding the contract method 0x12cc2a3b.
//
// Solidity: function DepositAndCall(address zrc4, uint256 amount, address target, bytes message) returns()
func (_ZetaDepositAndCall *ZetaDepositAndCallSession) DepositAndCall(zrc4 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _ZetaDepositAndCall.Contract.DepositAndCall(&_ZetaDepositAndCall.TransactOpts, zrc4, amount, target, message)
}

// DepositAndCall is a paid mutator transaction binding the contract method 0x12cc2a3b.
//
// Solidity: function DepositAndCall(address zrc4, uint256 amount, address target, bytes message) returns()
func (_ZetaDepositAndCall *ZetaDepositAndCallTransactorSession) DepositAndCall(zrc4 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _ZetaDepositAndCall.Contract.DepositAndCall(&_ZetaDepositAndCall.TransactOpts, zrc4, amount, target, message)
}
