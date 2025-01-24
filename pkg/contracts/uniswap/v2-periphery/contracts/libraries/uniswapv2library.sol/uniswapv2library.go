// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package uniswapv2library

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

// UniswapV2LibraryMetaData contains all meta data concerning the UniswapV2Library contract.
var UniswapV2LibraryMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212208421bd92ba607f7025ef11dd3ba3c30a46c62559adc890b6863b512efbd8984464736f6c63430006060033",
}

// UniswapV2LibraryABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV2LibraryMetaData.ABI instead.
var UniswapV2LibraryABI = UniswapV2LibraryMetaData.ABI

// UniswapV2LibraryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UniswapV2LibraryMetaData.Bin instead.
var UniswapV2LibraryBin = UniswapV2LibraryMetaData.Bin

// DeployUniswapV2Library deploys a new Ethereum contract, binding an instance of UniswapV2Library to it.
func DeployUniswapV2Library(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UniswapV2Library, error) {
	parsed, err := UniswapV2LibraryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UniswapV2LibraryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UniswapV2Library{UniswapV2LibraryCaller: UniswapV2LibraryCaller{contract: contract}, UniswapV2LibraryTransactor: UniswapV2LibraryTransactor{contract: contract}, UniswapV2LibraryFilterer: UniswapV2LibraryFilterer{contract: contract}}, nil
}

// UniswapV2Library is an auto generated Go binding around an Ethereum contract.
type UniswapV2Library struct {
	UniswapV2LibraryCaller     // Read-only binding to the contract
	UniswapV2LibraryTransactor // Write-only binding to the contract
	UniswapV2LibraryFilterer   // Log filterer for contract events
}

// UniswapV2LibraryCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapV2LibraryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2LibraryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapV2LibraryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2LibraryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapV2LibraryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2LibrarySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapV2LibrarySession struct {
	Contract     *UniswapV2Library // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniswapV2LibraryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapV2LibraryCallerSession struct {
	Contract *UniswapV2LibraryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// UniswapV2LibraryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapV2LibraryTransactorSession struct {
	Contract     *UniswapV2LibraryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// UniswapV2LibraryRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapV2LibraryRaw struct {
	Contract *UniswapV2Library // Generic contract binding to access the raw methods on
}

// UniswapV2LibraryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapV2LibraryCallerRaw struct {
	Contract *UniswapV2LibraryCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV2LibraryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapV2LibraryTransactorRaw struct {
	Contract *UniswapV2LibraryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV2Library creates a new instance of UniswapV2Library, bound to a specific deployed contract.
func NewUniswapV2Library(address common.Address, backend bind.ContractBackend) (*UniswapV2Library, error) {
	contract, err := bindUniswapV2Library(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV2Library{UniswapV2LibraryCaller: UniswapV2LibraryCaller{contract: contract}, UniswapV2LibraryTransactor: UniswapV2LibraryTransactor{contract: contract}, UniswapV2LibraryFilterer: UniswapV2LibraryFilterer{contract: contract}}, nil
}

// NewUniswapV2LibraryCaller creates a new read-only instance of UniswapV2Library, bound to a specific deployed contract.
func NewUniswapV2LibraryCaller(address common.Address, caller bind.ContractCaller) (*UniswapV2LibraryCaller, error) {
	contract, err := bindUniswapV2Library(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2LibraryCaller{contract: contract}, nil
}

// NewUniswapV2LibraryTransactor creates a new write-only instance of UniswapV2Library, bound to a specific deployed contract.
func NewUniswapV2LibraryTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV2LibraryTransactor, error) {
	contract, err := bindUniswapV2Library(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2LibraryTransactor{contract: contract}, nil
}

// NewUniswapV2LibraryFilterer creates a new log filterer instance of UniswapV2Library, bound to a specific deployed contract.
func NewUniswapV2LibraryFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV2LibraryFilterer, error) {
	contract, err := bindUniswapV2Library(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV2LibraryFilterer{contract: contract}, nil
}

// bindUniswapV2Library binds a generic wrapper to an already deployed contract.
func bindUniswapV2Library(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapV2LibraryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Library *UniswapV2LibraryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Library.Contract.UniswapV2LibraryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Library *UniswapV2LibraryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Library.Contract.UniswapV2LibraryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Library *UniswapV2LibraryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Library.Contract.UniswapV2LibraryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Library *UniswapV2LibraryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Library.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Library *UniswapV2LibraryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Library.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Library *UniswapV2LibraryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Library.Contract.contract.Transact(opts, method, params...)
}
