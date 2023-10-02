// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// ReverterzContext is an auto generated low-level Go binding around an user-defined struct.
type ReverterzContext struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}

// ReverterMetaData contains all meta data concerning the Reverter contract.
var ReverterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"Foo\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structReverter.zContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061027f806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063de43156e14610030575b600080fd5b61004a600480360381019061004591906101a5565b61004c565b005b6040517fbfb4ebcf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080fd5b600080fd5b600080fd5b6000606082840312156100a3576100a2610088565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100d7826100ac565b9050919050565b6100e7816100cc565b81146100f257600080fd5b50565b600081359050610104816100de565b92915050565b6000819050919050565b61011d8161010a565b811461012857600080fd5b50565b60008135905061013a81610114565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f84011261016557610164610140565b5b8235905067ffffffffffffffff81111561018257610181610145565b5b60208301915083600182028301111561019e5761019d61014a565b5b9250929050565b6000806000806000608086880312156101c1576101c061007e565b5b600086013567ffffffffffffffff8111156101df576101de610083565b5b6101eb8882890161008d565b95505060206101fc888289016100f5565b945050604061020d8882890161012b565b935050606086013567ffffffffffffffff81111561022e5761022d610083565b5b61023a8882890161014f565b9250925050929550929590935056fea2646970667358221220793398f92335fd817a4bf49381e3e4cbdfa00223720869567e31300bf71f0c9e64736f6c63430008150033",
}

// ReverterABI is the input ABI used to generate the binding from.
// Deprecated: Use ReverterMetaData.ABI instead.
var ReverterABI = ReverterMetaData.ABI

// ReverterBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ReverterMetaData.Bin instead.
var ReverterBin = ReverterMetaData.Bin

// DeployReverter deploys a new Ethereum contract, binding an instance of Reverter to it.
func DeployReverter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Reverter, error) {
	parsed, err := ReverterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ReverterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Reverter{ReverterCaller: ReverterCaller{contract: contract}, ReverterTransactor: ReverterTransactor{contract: contract}, ReverterFilterer: ReverterFilterer{contract: contract}}, nil
}

// Reverter is an auto generated Go binding around an Ethereum contract.
type Reverter struct {
	ReverterCaller     // Read-only binding to the contract
	ReverterTransactor // Write-only binding to the contract
	ReverterFilterer   // Log filterer for contract events
}

// ReverterCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReverterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReverterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReverterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReverterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReverterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReverterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReverterSession struct {
	Contract     *Reverter         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReverterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReverterCallerSession struct {
	Contract *ReverterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ReverterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReverterTransactorSession struct {
	Contract     *ReverterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ReverterRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReverterRaw struct {
	Contract *Reverter // Generic contract binding to access the raw methods on
}

// ReverterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReverterCallerRaw struct {
	Contract *ReverterCaller // Generic read-only contract binding to access the raw methods on
}

// ReverterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReverterTransactorRaw struct {
	Contract *ReverterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReverter creates a new instance of Reverter, bound to a specific deployed contract.
func NewReverter(address common.Address, backend bind.ContractBackend) (*Reverter, error) {
	contract, err := bindReverter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Reverter{ReverterCaller: ReverterCaller{contract: contract}, ReverterTransactor: ReverterTransactor{contract: contract}, ReverterFilterer: ReverterFilterer{contract: contract}}, nil
}

// NewReverterCaller creates a new read-only instance of Reverter, bound to a specific deployed contract.
func NewReverterCaller(address common.Address, caller bind.ContractCaller) (*ReverterCaller, error) {
	contract, err := bindReverter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReverterCaller{contract: contract}, nil
}

// NewReverterTransactor creates a new write-only instance of Reverter, bound to a specific deployed contract.
func NewReverterTransactor(address common.Address, transactor bind.ContractTransactor) (*ReverterTransactor, error) {
	contract, err := bindReverter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReverterTransactor{contract: contract}, nil
}

// NewReverterFilterer creates a new log filterer instance of Reverter, bound to a specific deployed contract.
func NewReverterFilterer(address common.Address, filterer bind.ContractFilterer) (*ReverterFilterer, error) {
	contract, err := bindReverter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReverterFilterer{contract: contract}, nil
}

// bindReverter binds a generic wrapper to an already deployed contract.
func bindReverter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReverterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reverter *ReverterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Reverter.Contract.ReverterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reverter *ReverterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reverter.Contract.ReverterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reverter *ReverterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reverter.Contract.ReverterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reverter *ReverterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Reverter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reverter *ReverterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reverter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reverter *ReverterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reverter.Contract.contract.Transact(opts, method, params...)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_Reverter *ReverterTransactor) OnCrossChainCall(opts *bind.TransactOpts, context ReverterzContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _Reverter.contract.Transact(opts, "onCrossChainCall", context, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_Reverter *ReverterSession) OnCrossChainCall(context ReverterzContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _Reverter.Contract.OnCrossChainCall(&_Reverter.TransactOpts, context, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_Reverter *ReverterTransactorSession) OnCrossChainCall(context ReverterzContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _Reverter.Contract.OnCrossChainCall(&_Reverter.TransactOpts, context, zrc20, amount, message)
}
