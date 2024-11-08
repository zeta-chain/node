// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testgasconsumer

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

// TestGasConsumerzContext is an auto generated low-level Go binding around an user-defined struct.
type TestGasConsumerzContext struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}

// TestGasConsumerMetaData contains all meta data concerning the TestGasConsumer contract.
var TestGasConsumerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structTestGasConsumer.zContext\",\"name\":\"_context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5061036d8061001f6000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80635bcfd61614610030575b600080fd5b61004a60048036038101906100459190610233565b61004c565b005b61005461005b565b5050505050565b6000624c4b4090506000614e209050600081836100789190610306565b905060005b818110156100bb576000819080600181540180825580915050600190039060005260206000200160009091909190915055808060010191505061007d565b506000806100c991906100ce565b505050565b50805460008255906000526020600020908101906100ec91906100ef565b50565b5b808211156101085760008160009055506001016100f0565b5090565b600080fd5b600080fd5b600080fd5b60006060828403121561013157610130610116565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006101658261013a565b9050919050565b6101758161015a565b811461018057600080fd5b50565b6000813590506101928161016c565b92915050565b6000819050919050565b6101ab81610198565b81146101b657600080fd5b50565b6000813590506101c8816101a2565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126101f3576101f26101ce565b5b8235905067ffffffffffffffff8111156102105761020f6101d3565b5b60208301915083600182028301111561022c5761022b6101d8565b5b9250929050565b60008060008060006080868803121561024f5761024e61010c565b5b600086013567ffffffffffffffff81111561026d5761026c610111565b5b6102798882890161011b565b955050602061028a88828901610183565b945050604061029b888289016101b9565b935050606086013567ffffffffffffffff8111156102bc576102bb610111565b5b6102c8888289016101dd565b92509250509295509295909350565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061031182610198565b915061031c83610198565b92508261032c5761032b6102d7565b5b82820490509291505056fea2646970667358221220e1d03a34090a8a647a128849d9f9434831ba3b1e4d28a514d9c9dc922068351e64736f6c634300081a0033",
}

// TestGasConsumerABI is the input ABI used to generate the binding from.
// Deprecated: Use TestGasConsumerMetaData.ABI instead.
var TestGasConsumerABI = TestGasConsumerMetaData.ABI

// TestGasConsumerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestGasConsumerMetaData.Bin instead.
var TestGasConsumerBin = TestGasConsumerMetaData.Bin

// DeployTestGasConsumer deploys a new Ethereum contract, binding an instance of TestGasConsumer to it.
func DeployTestGasConsumer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestGasConsumer, error) {
	parsed, err := TestGasConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestGasConsumerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestGasConsumer{TestGasConsumerCaller: TestGasConsumerCaller{contract: contract}, TestGasConsumerTransactor: TestGasConsumerTransactor{contract: contract}, TestGasConsumerFilterer: TestGasConsumerFilterer{contract: contract}}, nil
}

// TestGasConsumer is an auto generated Go binding around an Ethereum contract.
type TestGasConsumer struct {
	TestGasConsumerCaller     // Read-only binding to the contract
	TestGasConsumerTransactor // Write-only binding to the contract
	TestGasConsumerFilterer   // Log filterer for contract events
}

// TestGasConsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestGasConsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestGasConsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestGasConsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestGasConsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestGasConsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestGasConsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestGasConsumerSession struct {
	Contract     *TestGasConsumer  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestGasConsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestGasConsumerCallerSession struct {
	Contract *TestGasConsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// TestGasConsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestGasConsumerTransactorSession struct {
	Contract     *TestGasConsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// TestGasConsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestGasConsumerRaw struct {
	Contract *TestGasConsumer // Generic contract binding to access the raw methods on
}

// TestGasConsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestGasConsumerCallerRaw struct {
	Contract *TestGasConsumerCaller // Generic read-only contract binding to access the raw methods on
}

// TestGasConsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestGasConsumerTransactorRaw struct {
	Contract *TestGasConsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestGasConsumer creates a new instance of TestGasConsumer, bound to a specific deployed contract.
func NewTestGasConsumer(address common.Address, backend bind.ContractBackend) (*TestGasConsumer, error) {
	contract, err := bindTestGasConsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestGasConsumer{TestGasConsumerCaller: TestGasConsumerCaller{contract: contract}, TestGasConsumerTransactor: TestGasConsumerTransactor{contract: contract}, TestGasConsumerFilterer: TestGasConsumerFilterer{contract: contract}}, nil
}

// NewTestGasConsumerCaller creates a new read-only instance of TestGasConsumer, bound to a specific deployed contract.
func NewTestGasConsumerCaller(address common.Address, caller bind.ContractCaller) (*TestGasConsumerCaller, error) {
	contract, err := bindTestGasConsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestGasConsumerCaller{contract: contract}, nil
}

// NewTestGasConsumerTransactor creates a new write-only instance of TestGasConsumer, bound to a specific deployed contract.
func NewTestGasConsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*TestGasConsumerTransactor, error) {
	contract, err := bindTestGasConsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestGasConsumerTransactor{contract: contract}, nil
}

// NewTestGasConsumerFilterer creates a new log filterer instance of TestGasConsumer, bound to a specific deployed contract.
func NewTestGasConsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*TestGasConsumerFilterer, error) {
	contract, err := bindTestGasConsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestGasConsumerFilterer{contract: contract}, nil
}

// bindTestGasConsumer binds a generic wrapper to an already deployed contract.
func bindTestGasConsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestGasConsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestGasConsumer *TestGasConsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestGasConsumer.Contract.TestGasConsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestGasConsumer *TestGasConsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestGasConsumer.Contract.TestGasConsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestGasConsumer *TestGasConsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestGasConsumer.Contract.TestGasConsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestGasConsumer *TestGasConsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestGasConsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestGasConsumer *TestGasConsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestGasConsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestGasConsumer *TestGasConsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestGasConsumer.Contract.contract.Transact(opts, method, params...)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) _context, address _zrc20, uint256 _amount, bytes _message) returns()
func (_TestGasConsumer *TestGasConsumerTransactor) OnCall(opts *bind.TransactOpts, _context TestGasConsumerzContext, _zrc20 common.Address, _amount *big.Int, _message []byte) (*types.Transaction, error) {
	return _TestGasConsumer.contract.Transact(opts, "onCall", _context, _zrc20, _amount, _message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) _context, address _zrc20, uint256 _amount, bytes _message) returns()
func (_TestGasConsumer *TestGasConsumerSession) OnCall(_context TestGasConsumerzContext, _zrc20 common.Address, _amount *big.Int, _message []byte) (*types.Transaction, error) {
	return _TestGasConsumer.Contract.OnCall(&_TestGasConsumer.TransactOpts, _context, _zrc20, _amount, _message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) _context, address _zrc20, uint256 _amount, bytes _message) returns()
func (_TestGasConsumer *TestGasConsumerTransactorSession) OnCall(_context TestGasConsumerzContext, _zrc20 common.Address, _amount *big.Int, _message []byte) (*types.Transaction, error) {
	return _TestGasConsumer.Contract.OnCall(&_TestGasConsumer.TransactOpts, _context, _zrc20, _amount, _message)
}
