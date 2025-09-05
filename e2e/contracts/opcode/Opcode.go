// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package opcode

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

// OpcodeMetaData contains all meta data concerning the Opcode contract.
var OpcodeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"testPUSH0\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testTLOAD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b5060e78061001b5f395ff3fe6080604052348015600e575f80fd5b50600436106030575f3560e01c8063d1529e50146034578063eac28ca814604e575b5f80fd5b603a6068565b60405160459190609a565b60405180910390f35b60546072565b604051605f9190609a565b60405180910390f35b5f80805f5260205ff35b5f806112345f5d5f5c90508091505090565b5f819050919050565b6094816084565b82525050565b5f60208201905060ab5f830184608d565b9291505056fea2646970667358221220438d8394df6669497c00f8449340fa539c875be843a40a7d63ebc2d9377fb2f664736f6c634300081a0033",
}

// OpcodeABI is the input ABI used to generate the binding from.
// Deprecated: Use OpcodeMetaData.ABI instead.
var OpcodeABI = OpcodeMetaData.ABI

// OpcodeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OpcodeMetaData.Bin instead.
var OpcodeBin = OpcodeMetaData.Bin

// DeployOpcode deploys a new Ethereum contract, binding an instance of Opcode to it.
func DeployOpcode(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Opcode, error) {
	parsed, err := OpcodeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OpcodeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Opcode{OpcodeCaller: OpcodeCaller{contract: contract}, OpcodeTransactor: OpcodeTransactor{contract: contract}, OpcodeFilterer: OpcodeFilterer{contract: contract}}, nil
}

// Opcode is an auto generated Go binding around an Ethereum contract.
type Opcode struct {
	OpcodeCaller     // Read-only binding to the contract
	OpcodeTransactor // Write-only binding to the contract
	OpcodeFilterer   // Log filterer for contract events
}

// OpcodeCaller is an auto generated read-only Go binding around an Ethereum contract.
type OpcodeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpcodeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OpcodeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpcodeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OpcodeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpcodeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OpcodeSession struct {
	Contract     *Opcode           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OpcodeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OpcodeCallerSession struct {
	Contract *OpcodeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OpcodeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OpcodeTransactorSession struct {
	Contract     *OpcodeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OpcodeRaw is an auto generated low-level Go binding around an Ethereum contract.
type OpcodeRaw struct {
	Contract *Opcode // Generic contract binding to access the raw methods on
}

// OpcodeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OpcodeCallerRaw struct {
	Contract *OpcodeCaller // Generic read-only contract binding to access the raw methods on
}

// OpcodeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OpcodeTransactorRaw struct {
	Contract *OpcodeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOpcode creates a new instance of Opcode, bound to a specific deployed contract.
func NewOpcode(address common.Address, backend bind.ContractBackend) (*Opcode, error) {
	contract, err := bindOpcode(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Opcode{OpcodeCaller: OpcodeCaller{contract: contract}, OpcodeTransactor: OpcodeTransactor{contract: contract}, OpcodeFilterer: OpcodeFilterer{contract: contract}}, nil
}

// NewOpcodeCaller creates a new read-only instance of Opcode, bound to a specific deployed contract.
func NewOpcodeCaller(address common.Address, caller bind.ContractCaller) (*OpcodeCaller, error) {
	contract, err := bindOpcode(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OpcodeCaller{contract: contract}, nil
}

// NewOpcodeTransactor creates a new write-only instance of Opcode, bound to a specific deployed contract.
func NewOpcodeTransactor(address common.Address, transactor bind.ContractTransactor) (*OpcodeTransactor, error) {
	contract, err := bindOpcode(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OpcodeTransactor{contract: contract}, nil
}

// NewOpcodeFilterer creates a new log filterer instance of Opcode, bound to a specific deployed contract.
func NewOpcodeFilterer(address common.Address, filterer bind.ContractFilterer) (*OpcodeFilterer, error) {
	contract, err := bindOpcode(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OpcodeFilterer{contract: contract}, nil
}

// bindOpcode binds a generic wrapper to an already deployed contract.
func bindOpcode(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OpcodeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Opcode *OpcodeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Opcode.Contract.OpcodeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Opcode *OpcodeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Opcode.Contract.OpcodeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Opcode *OpcodeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Opcode.Contract.OpcodeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Opcode *OpcodeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Opcode.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Opcode *OpcodeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Opcode.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Opcode *OpcodeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Opcode.Contract.contract.Transact(opts, method, params...)
}

// TestPUSH0 is a paid mutator transaction binding the contract method 0xd1529e50.
//
// Solidity: function testPUSH0() returns(uint256)
func (_Opcode *OpcodeTransactor) TestPUSH0(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Opcode.contract.Transact(opts, "testPUSH0")
}

// TestPUSH0 is a paid mutator transaction binding the contract method 0xd1529e50.
//
// Solidity: function testPUSH0() returns(uint256)
func (_Opcode *OpcodeSession) TestPUSH0() (*types.Transaction, error) {
	return _Opcode.Contract.TestPUSH0(&_Opcode.TransactOpts)
}

// TestPUSH0 is a paid mutator transaction binding the contract method 0xd1529e50.
//
// Solidity: function testPUSH0() returns(uint256)
func (_Opcode *OpcodeTransactorSession) TestPUSH0() (*types.Transaction, error) {
	return _Opcode.Contract.TestPUSH0(&_Opcode.TransactOpts)
}

// TestTLOAD is a paid mutator transaction binding the contract method 0xeac28ca8.
//
// Solidity: function testTLOAD() returns(uint256)
func (_Opcode *OpcodeTransactor) TestTLOAD(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Opcode.contract.Transact(opts, "testTLOAD")
}

// TestTLOAD is a paid mutator transaction binding the contract method 0xeac28ca8.
//
// Solidity: function testTLOAD() returns(uint256)
func (_Opcode *OpcodeSession) TestTLOAD() (*types.Transaction, error) {
	return _Opcode.Contract.TestTLOAD(&_Opcode.TransactOpts)
}

// TestTLOAD is a paid mutator transaction binding the contract method 0xeac28ca8.
//
// Solidity: function testTLOAD() returns(uint256)
func (_Opcode *OpcodeTransactorSession) TestTLOAD() (*types.Transaction, error) {
	return _Opcode.Contract.TestTLOAD(&_Opcode.TransactOpts)
}
