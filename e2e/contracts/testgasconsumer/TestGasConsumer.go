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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_targetGas\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structTestGasConsumer.zContext\",\"name\":\"_context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a0604052348015600f57600080fd5b506040516104383803806104388339818101604052810190602f91906072565b806080818152505050609a565b600080fd5b6000819050919050565b6052816041565b8114605c57600080fd5b50565b600081519050606c81604b565b92915050565b6000602082840312156085576084603c565b5b6000609184828501605f565b91505092915050565b6080516103846100b46000396000606701526103846000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80635bcfd61614610030575b600080fd5b61004a6004803603810190610045919061024a565b61004c565b005b61005461005b565b5050505050565b6000614e2090506000817f0000000000000000000000000000000000000000000000000000000000000000610090919061031d565b905060005b818110156100d35760008190806001815401808255809150506001900390600052602060002001600090919091909150558080600101915050610095565b506000806100e191906100e5565b5050565b50805460008255906000526020600020908101906101039190610106565b50565b5b8082111561011f576000816000905550600101610107565b5090565b600080fd5b600080fd5b600080fd5b6000606082840312156101485761014761012d565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061017c82610151565b9050919050565b61018c81610171565b811461019757600080fd5b50565b6000813590506101a981610183565b92915050565b6000819050919050565b6101c2816101af565b81146101cd57600080fd5b50565b6000813590506101df816101b9565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f84011261020a576102096101e5565b5b8235905067ffffffffffffffff811115610227576102266101ea565b5b602083019150836001820283011115610243576102426101ef565b5b9250929050565b60008060008060006080868803121561026657610265610123565b5b600086013567ffffffffffffffff81111561028457610283610128565b5b61029088828901610132565b95505060206102a18882890161019a565b94505060406102b2888289016101d0565b935050606086013567ffffffffffffffff8111156102d3576102d2610128565b5b6102df888289016101f4565b92509250509295509295909350565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000610328826101af565b9150610333836101af565b925082610343576103426102ee565b5b82820490509291505056fea264697066735822122074511fb2a3edc8019fca5d2e2d94fe97609fa1383042562826b2fee4e67c513e64736f6c634300081a0033",
}

// TestGasConsumerABI is the input ABI used to generate the binding from.
// Deprecated: Use TestGasConsumerMetaData.ABI instead.
var TestGasConsumerABI = TestGasConsumerMetaData.ABI

// TestGasConsumerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestGasConsumerMetaData.Bin instead.
var TestGasConsumerBin = TestGasConsumerMetaData.Bin

// DeployTestGasConsumer deploys a new Ethereum contract, binding an instance of TestGasConsumer to it.
func DeployTestGasConsumer(auth *bind.TransactOpts, backend bind.ContractBackend, _targetGas *big.Int) (common.Address, *types.Transaction, *TestGasConsumer, error) {
	parsed, err := TestGasConsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestGasConsumerBin), backend, _targetGas)
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
