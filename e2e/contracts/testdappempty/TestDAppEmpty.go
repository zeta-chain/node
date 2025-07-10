// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testdappempty

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

// TestDAppEmptyMessageContext is an auto generated low-level Go binding around an user-defined struct.
type TestDAppEmptyMessageContext struct {
	Sender common.Address
}

// TestDAppEmptyRevertContext is an auto generated low-level Go binding around an user-defined struct.
type TestDAppEmptyRevertContext struct {
	Sender        common.Address
	Asset         common.Address
	Amount        *big.Int
	RevertMessage []byte
}

// TestDAppEmptyzContext is an auto generated low-level Go binding around an user-defined struct.
type TestDAppEmptyzContext struct {
	Sender    []byte
	SenderEVM common.Address
	ChainID   *big.Int
}

// TestDAppEmptyMetaData contains all meta data concerning the TestDAppEmpty contract.
var TestDAppEmptyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"senderEVM\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structTestDAppEmpty.zContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structTestDAppEmpty.MessageContext\",\"name\":\"messageContext\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"}],\"internalType\":\"structTestDAppEmpty.RevertContext\",\"name\":\"revertContext\",\"type\":\"tuple\"}],\"name\":\"onRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506104858061001f6000396000f3fe6080604052600436106100385760003560e01c80635bcfd61614610044578063676cc0541461006d578063c9028a361461009d5761003f565b3661003f57005b600080fd5b34801561005057600080fd5b5061006b60048036038101906100669190610212565b6100c6565b005b610087600480360381019061008291906102d5565b6100cd565b60405161009491906103c5565b60405180910390f35b3480156100a957600080fd5b506100c460048036038101906100bf9190610406565b6100e8565b005b5050505050565b60606040518060200160405280600081525090509392505050565b50565b600080fd5b600080fd5b600080fd5b6000606082840312156101105761010f6100f5565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061014482610119565b9050919050565b61015481610139565b811461015f57600080fd5b50565b6000813590506101718161014b565b92915050565b6000819050919050565b61018a81610177565b811461019557600080fd5b50565b6000813590506101a781610181565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126101d2576101d16101ad565b5b8235905067ffffffffffffffff8111156101ef576101ee6101b2565b5b60208301915083600182028301111561020b5761020a6101b7565b5b9250929050565b60008060008060006080868803121561022e5761022d6100eb565b5b600086013567ffffffffffffffff81111561024c5761024b6100f0565b5b610258888289016100fa565b955050602061026988828901610162565b945050604061027a88828901610198565b935050606086013567ffffffffffffffff81111561029b5761029a6100f0565b5b6102a7888289016101bc565b92509250509295509295909350565b6000602082840312156102cc576102cb6100f5565b5b81905092915050565b6000806000604084860312156102ee576102ed6100eb565b5b60006102fc868287016102b6565b935050602084013567ffffffffffffffff81111561031d5761031c6100f0565b5b610329868287016101bc565b92509250509250925092565b600081519050919050565b600082825260208201905092915050565b60005b8381101561036f578082015181840152602081019050610354565b60008484015250505050565b6000601f19601f8301169050919050565b600061039782610335565b6103a18185610340565b93506103b1818560208601610351565b6103ba8161037b565b840191505092915050565b600060208201905081810360008301526103df818461038c565b905092915050565b6000608082840312156103fd576103fc6100f5565b5b81905092915050565b60006020828403121561041c5761041b6100eb565b5b600082013567ffffffffffffffff81111561043a576104396100f0565b5b610446848285016103e7565b9150509291505056fea26469706673582212200137f7f41c83c091d3d75f57459459833158a855c08d2ba7413098558ee6322964736f6c634300081a0033",
}

// TestDAppEmptyABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDAppEmptyMetaData.ABI instead.
var TestDAppEmptyABI = TestDAppEmptyMetaData.ABI

// TestDAppEmptyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestDAppEmptyMetaData.Bin instead.
var TestDAppEmptyBin = TestDAppEmptyMetaData.Bin

// DeployTestDAppEmpty deploys a new Ethereum contract, binding an instance of TestDAppEmpty to it.
func DeployTestDAppEmpty(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestDAppEmpty, error) {
	parsed, err := TestDAppEmptyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDAppEmptyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestDAppEmpty{TestDAppEmptyCaller: TestDAppEmptyCaller{contract: contract}, TestDAppEmptyTransactor: TestDAppEmptyTransactor{contract: contract}, TestDAppEmptyFilterer: TestDAppEmptyFilterer{contract: contract}}, nil
}

// TestDAppEmpty is an auto generated Go binding around an Ethereum contract.
type TestDAppEmpty struct {
	TestDAppEmptyCaller     // Read-only binding to the contract
	TestDAppEmptyTransactor // Write-only binding to the contract
	TestDAppEmptyFilterer   // Log filterer for contract events
}

// TestDAppEmptyCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestDAppEmptyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppEmptyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestDAppEmptyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppEmptyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestDAppEmptyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppEmptySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestDAppEmptySession struct {
	Contract     *TestDAppEmpty    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestDAppEmptyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestDAppEmptyCallerSession struct {
	Contract *TestDAppEmptyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TestDAppEmptyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestDAppEmptyTransactorSession struct {
	Contract     *TestDAppEmptyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TestDAppEmptyRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestDAppEmptyRaw struct {
	Contract *TestDAppEmpty // Generic contract binding to access the raw methods on
}

// TestDAppEmptyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestDAppEmptyCallerRaw struct {
	Contract *TestDAppEmptyCaller // Generic read-only contract binding to access the raw methods on
}

// TestDAppEmptyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestDAppEmptyTransactorRaw struct {
	Contract *TestDAppEmptyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestDAppEmpty creates a new instance of TestDAppEmpty, bound to a specific deployed contract.
func NewTestDAppEmpty(address common.Address, backend bind.ContractBackend) (*TestDAppEmpty, error) {
	contract, err := bindTestDAppEmpty(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestDAppEmpty{TestDAppEmptyCaller: TestDAppEmptyCaller{contract: contract}, TestDAppEmptyTransactor: TestDAppEmptyTransactor{contract: contract}, TestDAppEmptyFilterer: TestDAppEmptyFilterer{contract: contract}}, nil
}

// NewTestDAppEmptyCaller creates a new read-only instance of TestDAppEmpty, bound to a specific deployed contract.
func NewTestDAppEmptyCaller(address common.Address, caller bind.ContractCaller) (*TestDAppEmptyCaller, error) {
	contract, err := bindTestDAppEmpty(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppEmptyCaller{contract: contract}, nil
}

// NewTestDAppEmptyTransactor creates a new write-only instance of TestDAppEmpty, bound to a specific deployed contract.
func NewTestDAppEmptyTransactor(address common.Address, transactor bind.ContractTransactor) (*TestDAppEmptyTransactor, error) {
	contract, err := bindTestDAppEmpty(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppEmptyTransactor{contract: contract}, nil
}

// NewTestDAppEmptyFilterer creates a new log filterer instance of TestDAppEmpty, bound to a specific deployed contract.
func NewTestDAppEmptyFilterer(address common.Address, filterer bind.ContractFilterer) (*TestDAppEmptyFilterer, error) {
	contract, err := bindTestDAppEmpty(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestDAppEmptyFilterer{contract: contract}, nil
}

// bindTestDAppEmpty binds a generic wrapper to an already deployed contract.
func bindTestDAppEmpty(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestDAppEmptyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDAppEmpty *TestDAppEmptyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDAppEmpty.Contract.TestDAppEmptyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDAppEmpty *TestDAppEmptyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.TestDAppEmptyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDAppEmpty *TestDAppEmptyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.TestDAppEmptyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDAppEmpty *TestDAppEmptyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDAppEmpty.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDAppEmpty *TestDAppEmptyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDAppEmpty *TestDAppEmptyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.contract.Transact(opts, method, params...)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) context, address _zrc20, uint256 amount, bytes message) returns()
func (_TestDAppEmpty *TestDAppEmptyTransactor) OnCall(opts *bind.TransactOpts, context TestDAppEmptyzContext, _zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppEmpty.contract.Transact(opts, "onCall", context, _zrc20, amount, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) context, address _zrc20, uint256 amount, bytes message) returns()
func (_TestDAppEmpty *TestDAppEmptySession) OnCall(context TestDAppEmptyzContext, _zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.OnCall(&_TestDAppEmpty.TransactOpts, context, _zrc20, amount, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) context, address _zrc20, uint256 amount, bytes message) returns()
func (_TestDAppEmpty *TestDAppEmptyTransactorSession) OnCall(context TestDAppEmptyzContext, _zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.OnCall(&_TestDAppEmpty.TransactOpts, context, _zrc20, amount, message)
}

// OnCall0 is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppEmpty *TestDAppEmptyTransactor) OnCall0(opts *bind.TransactOpts, messageContext TestDAppEmptyMessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppEmpty.contract.Transact(opts, "onCall0", messageContext, message)
}

// OnCall0 is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppEmpty *TestDAppEmptySession) OnCall0(messageContext TestDAppEmptyMessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.OnCall0(&_TestDAppEmpty.TransactOpts, messageContext, message)
}

// OnCall0 is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppEmpty *TestDAppEmptyTransactorSession) OnCall0(messageContext TestDAppEmptyMessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.OnCall0(&_TestDAppEmpty.TransactOpts, messageContext, message)
}

// OnRevert is a paid mutator transaction binding the contract method 0xc9028a36.
//
// Solidity: function onRevert((address,address,uint256,bytes) revertContext) returns()
func (_TestDAppEmpty *TestDAppEmptyTransactor) OnRevert(opts *bind.TransactOpts, revertContext TestDAppEmptyRevertContext) (*types.Transaction, error) {
	return _TestDAppEmpty.contract.Transact(opts, "onRevert", revertContext)
}

// OnRevert is a paid mutator transaction binding the contract method 0xc9028a36.
//
// Solidity: function onRevert((address,address,uint256,bytes) revertContext) returns()
func (_TestDAppEmpty *TestDAppEmptySession) OnRevert(revertContext TestDAppEmptyRevertContext) (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.OnRevert(&_TestDAppEmpty.TransactOpts, revertContext)
}

// OnRevert is a paid mutator transaction binding the contract method 0xc9028a36.
//
// Solidity: function onRevert((address,address,uint256,bytes) revertContext) returns()
func (_TestDAppEmpty *TestDAppEmptyTransactorSession) OnRevert(revertContext TestDAppEmptyRevertContext) (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.OnRevert(&_TestDAppEmpty.TransactOpts, revertContext)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDAppEmpty *TestDAppEmptyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDAppEmpty.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDAppEmpty *TestDAppEmptySession) Receive() (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.Receive(&_TestDAppEmpty.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDAppEmpty *TestDAppEmptyTransactorSession) Receive() (*types.Transaction, error) {
	return _TestDAppEmpty.Contract.Receive(&_TestDAppEmpty.TransactOpts)
}
