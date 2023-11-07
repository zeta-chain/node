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

// DepositorMetaData contains all meta data concerning the Depositor contract.
var DepositorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"custody_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"},{\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"runDeposits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b506040516106d83803806106d8833981810160405281019061003291906100cf565b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff1681525050506100fc565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061009c82610071565b9050919050565b6100ac81610091565b81146100b757600080fd5b50565b6000815190506100c9816100a3565b92915050565b6000602082840312156100e5576100e461006c565b5b60006100f3848285016100ba565b91505092915050565b6080516105bc61011c60003960008181606a015260f101526105bc6000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80633d496c9314610030575b600080fd5b61004a600480360381019061004591906102b6565b61004c565b005b8473ffffffffffffffffffffffffffffffffffffffff1663095ea7b37f0000000000000000000000000000000000000000000000000000000000000000838761009591906103a1565b6040518363ffffffff1660e01b81526004016100b2929190610401565b600060405180830381600087803b1580156100cc57600080fd5b505af11580156100e0573d6000803e3d6000fd5b5050505060005b81811015610197577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663e609055e8989898989896040518763ffffffff1660e01b8152600401610152969594939291906104e7565b600060405180830381600087803b15801561016c57600080fd5b505af1158015610180573d6000803e3d6000fd5b50505050808061018f9061053e565b9150506100e7565b5050505050505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f8401126101d0576101cf6101ab565b5b8235905067ffffffffffffffff8111156101ed576101ec6101b0565b5b602083019150836001820283011115610209576102086101b5565b5b9250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061023b82610210565b9050919050565b600061024d82610230565b9050919050565b61025d81610242565b811461026857600080fd5b50565b60008135905061027a81610254565b92915050565b6000819050919050565b61029381610280565b811461029e57600080fd5b50565b6000813590506102b08161028a565b92915050565b600080600080600080600060a0888a0312156102d5576102d46101a1565b5b600088013567ffffffffffffffff8111156102f3576102f26101a6565b5b6102ff8a828b016101ba565b975097505060206103128a828b0161026b565b95505060406103238a828b016102a1565b945050606088013567ffffffffffffffff811115610344576103436101a6565b5b6103508a828b016101ba565b935093505060806103638a828b016102a1565b91505092959891949750929550565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006103ac82610280565b91506103b783610280565b92508282026103c581610280565b915082820484148315176103dc576103db610372565b5b5092915050565b6103ec81610230565b82525050565b6103fb81610280565b82525050565b600060408201905061041660008301856103e3565b61042360208301846103f2565b9392505050565b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b6000610467838561042a565b935061047483858461043b565b61047d8361044a565b840190509392505050565b6000819050919050565b60006104ad6104a86104a384610210565b610488565b610210565b9050919050565b60006104bf82610492565b9050919050565b60006104d1826104b4565b9050919050565b6104e1816104c6565b82525050565b6000608082019050818103600083015261050281888a61045b565b905061051160208301876104d8565b61051e60408301866103f2565b818103606083015261053181848661045b565b9050979650505050505050565b600061054982610280565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361057b5761057a610372565b5b60018201905091905056fea26469706673582212205c3c5f90d68ba0aa39770c8be50163c94d7c5e76e252bf5ef39e9ace318f1fef64736f6c63430008150033",
}

// DepositorABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositorMetaData.ABI instead.
var DepositorABI = DepositorMetaData.ABI

// DepositorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DepositorMetaData.Bin instead.
var DepositorBin = DepositorMetaData.Bin

// DeployDepositor deploys a new Ethereum contract, binding an instance of Depositor to it.
func DeployDepositor(auth *bind.TransactOpts, backend bind.ContractBackend, custody_ common.Address) (common.Address, *types.Transaction, *Depositor, error) {
	parsed, err := DepositorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DepositorBin), backend, custody_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Depositor{DepositorCaller: DepositorCaller{contract: contract}, DepositorTransactor: DepositorTransactor{contract: contract}, DepositorFilterer: DepositorFilterer{contract: contract}}, nil
}

// Depositor is an auto generated Go binding around an Ethereum contract.
type Depositor struct {
	DepositorCaller     // Read-only binding to the contract
	DepositorTransactor // Write-only binding to the contract
	DepositorFilterer   // Log filterer for contract events
}

// DepositorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositorSession struct {
	Contract     *Depositor        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositorCallerSession struct {
	Contract *DepositorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DepositorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositorTransactorSession struct {
	Contract     *DepositorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DepositorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositorRaw struct {
	Contract *Depositor // Generic contract binding to access the raw methods on
}

// DepositorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositorCallerRaw struct {
	Contract *DepositorCaller // Generic read-only contract binding to access the raw methods on
}

// DepositorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositorTransactorRaw struct {
	Contract *DepositorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositor creates a new instance of Depositor, bound to a specific deployed contract.
func NewDepositor(address common.Address, backend bind.ContractBackend) (*Depositor, error) {
	contract, err := bindDepositor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Depositor{DepositorCaller: DepositorCaller{contract: contract}, DepositorTransactor: DepositorTransactor{contract: contract}, DepositorFilterer: DepositorFilterer{contract: contract}}, nil
}

// NewDepositorCaller creates a new read-only instance of Depositor, bound to a specific deployed contract.
func NewDepositorCaller(address common.Address, caller bind.ContractCaller) (*DepositorCaller, error) {
	contract, err := bindDepositor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositorCaller{contract: contract}, nil
}

// NewDepositorTransactor creates a new write-only instance of Depositor, bound to a specific deployed contract.
func NewDepositorTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositorTransactor, error) {
	contract, err := bindDepositor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositorTransactor{contract: contract}, nil
}

// NewDepositorFilterer creates a new log filterer instance of Depositor, bound to a specific deployed contract.
func NewDepositorFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositorFilterer, error) {
	contract, err := bindDepositor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositorFilterer{contract: contract}, nil
}

// bindDepositor binds a generic wrapper to an already deployed contract.
func bindDepositor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depositor *DepositorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Depositor.Contract.DepositorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depositor *DepositorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositor.Contract.DepositorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depositor *DepositorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depositor.Contract.DepositorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depositor *DepositorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Depositor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depositor *DepositorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depositor *DepositorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depositor.Contract.contract.Transact(opts, method, params...)
}

// RunDeposits is a paid mutator transaction binding the contract method 0x3d496c93.
//
// Solidity: function runDeposits(bytes recipient, address asset, uint256 amount, bytes message, uint256 count) returns()
func (_Depositor *DepositorTransactor) RunDeposits(opts *bind.TransactOpts, recipient []byte, asset common.Address, amount *big.Int, message []byte, count *big.Int) (*types.Transaction, error) {
	return _Depositor.contract.Transact(opts, "runDeposits", recipient, asset, amount, message, count)
}

// RunDeposits is a paid mutator transaction binding the contract method 0x3d496c93.
//
// Solidity: function runDeposits(bytes recipient, address asset, uint256 amount, bytes message, uint256 count) returns()
func (_Depositor *DepositorSession) RunDeposits(recipient []byte, asset common.Address, amount *big.Int, message []byte, count *big.Int) (*types.Transaction, error) {
	return _Depositor.Contract.RunDeposits(&_Depositor.TransactOpts, recipient, asset, amount, message, count)
}

// RunDeposits is a paid mutator transaction binding the contract method 0x3d496c93.
//
// Solidity: function runDeposits(bytes recipient, address asset, uint256 amount, bytes message, uint256 count) returns()
func (_Depositor *DepositorTransactorSession) RunDeposits(recipient []byte, asset common.Address, amount *big.Int, message []byte, count *big.Int) (*types.Transaction, error) {
	return _Depositor.Contract.RunDeposits(&_Depositor.TransactOpts, recipient, asset, amount, message, count)
}
