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

// WithdrawerMetaData contains all meta data concerning the Withdrawer contract.
var WithdrawerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"},{\"internalType\":\"contractIZRC20\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"runWithdraws\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610543806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063e3be6f6814610030575b600080fd5b61004a60048036038101906100459190610282565b61004c565b005b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd333084866100769190610339565b6040518463ffffffff1660e01b815260040161009493929190610399565b600060405180830381600087803b1580156100ae57600080fd5b505af11580156100c2573d6000803e3d6000fd5b5050505060005b81811015610165578373ffffffffffffffffffffffffffffffffffffffff1663c70126268787866040518463ffffffff1660e01b815260040161010e9392919061042e565b6020604051808303816000875af115801561012d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101519190610498565b50808061015d906104c5565b9150506100c9565b505050505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f84011261019c5761019b610177565b5b8235905067ffffffffffffffff8111156101b9576101b861017c565b5b6020830191508360018202830111156101d5576101d4610181565b5b9250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610207826101dc565b9050919050565b6000610219826101fc565b9050919050565b6102298161020e565b811461023457600080fd5b50565b60008135905061024681610220565b92915050565b6000819050919050565b61025f8161024c565b811461026a57600080fd5b50565b60008135905061027c81610256565b92915050565b60008060008060006080868803121561029e5761029d61016d565b5b600086013567ffffffffffffffff8111156102bc576102bb610172565b5b6102c888828901610186565b955095505060206102db88828901610237565b93505060406102ec8882890161026d565b92505060606102fd8882890161026d565b9150509295509295909350565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006103448261024c565b915061034f8361024c565b925082820261035d8161024c565b915082820484148315176103745761037361030a565b5b5092915050565b610384816101fc565b82525050565b6103938161024c565b82525050565b60006060820190506103ae600083018661037b565b6103bb602083018561037b565b6103c8604083018461038a565b949350505050565b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b600061040d83856103d0565b935061041a8385846103e1565b610423836103f0565b840190509392505050565b60006040820190508181036000830152610449818587610401565b9050610458602083018461038a565b949350505050565b60008115159050919050565b61047581610460565b811461048057600080fd5b50565b6000815190506104928161046c565b92915050565b6000602082840312156104ae576104ad61016d565b5b60006104bc84828501610483565b91505092915050565b60006104d08261024c565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036105025761050161030a565b5b60018201905091905056fea264697066735822122010a6b06e6d91c64b6322e1ab1f565375da35df355d317055a2ba7e00b60ad26764736f6c63430008150033",
}

// WithdrawerABI is the input ABI used to generate the binding from.
// Deprecated: Use WithdrawerMetaData.ABI instead.
var WithdrawerABI = WithdrawerMetaData.ABI

// WithdrawerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use WithdrawerMetaData.Bin instead.
var WithdrawerBin = WithdrawerMetaData.Bin

// DeployWithdrawer deploys a new Ethereum contract, binding an instance of Withdrawer to it.
func DeployWithdrawer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Withdrawer, error) {
	parsed, err := WithdrawerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WithdrawerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Withdrawer{WithdrawerCaller: WithdrawerCaller{contract: contract}, WithdrawerTransactor: WithdrawerTransactor{contract: contract}, WithdrawerFilterer: WithdrawerFilterer{contract: contract}}, nil
}

// Withdrawer is an auto generated Go binding around an Ethereum contract.
type Withdrawer struct {
	WithdrawerCaller     // Read-only binding to the contract
	WithdrawerTransactor // Write-only binding to the contract
	WithdrawerFilterer   // Log filterer for contract events
}

// WithdrawerCaller is an auto generated read-only Go binding around an Ethereum contract.
type WithdrawerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WithdrawerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WithdrawerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WithdrawerSession struct {
	Contract     *Withdrawer       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WithdrawerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WithdrawerCallerSession struct {
	Contract *WithdrawerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// WithdrawerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WithdrawerTransactorSession struct {
	Contract     *WithdrawerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// WithdrawerRaw is an auto generated low-level Go binding around an Ethereum contract.
type WithdrawerRaw struct {
	Contract *Withdrawer // Generic contract binding to access the raw methods on
}

// WithdrawerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WithdrawerCallerRaw struct {
	Contract *WithdrawerCaller // Generic read-only contract binding to access the raw methods on
}

// WithdrawerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WithdrawerTransactorRaw struct {
	Contract *WithdrawerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWithdrawer creates a new instance of Withdrawer, bound to a specific deployed contract.
func NewWithdrawer(address common.Address, backend bind.ContractBackend) (*Withdrawer, error) {
	contract, err := bindWithdrawer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Withdrawer{WithdrawerCaller: WithdrawerCaller{contract: contract}, WithdrawerTransactor: WithdrawerTransactor{contract: contract}, WithdrawerFilterer: WithdrawerFilterer{contract: contract}}, nil
}

// NewWithdrawerCaller creates a new read-only instance of Withdrawer, bound to a specific deployed contract.
func NewWithdrawerCaller(address common.Address, caller bind.ContractCaller) (*WithdrawerCaller, error) {
	contract, err := bindWithdrawer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WithdrawerCaller{contract: contract}, nil
}

// NewWithdrawerTransactor creates a new write-only instance of Withdrawer, bound to a specific deployed contract.
func NewWithdrawerTransactor(address common.Address, transactor bind.ContractTransactor) (*WithdrawerTransactor, error) {
	contract, err := bindWithdrawer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WithdrawerTransactor{contract: contract}, nil
}

// NewWithdrawerFilterer creates a new log filterer instance of Withdrawer, bound to a specific deployed contract.
func NewWithdrawerFilterer(address common.Address, filterer bind.ContractFilterer) (*WithdrawerFilterer, error) {
	contract, err := bindWithdrawer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WithdrawerFilterer{contract: contract}, nil
}

// bindWithdrawer binds a generic wrapper to an already deployed contract.
func bindWithdrawer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WithdrawerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Withdrawer *WithdrawerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Withdrawer.Contract.WithdrawerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Withdrawer *WithdrawerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Withdrawer.Contract.WithdrawerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Withdrawer *WithdrawerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Withdrawer.Contract.WithdrawerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Withdrawer *WithdrawerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Withdrawer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Withdrawer *WithdrawerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Withdrawer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Withdrawer *WithdrawerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Withdrawer.Contract.contract.Transact(opts, method, params...)
}

// RunWithdraws is a paid mutator transaction binding the contract method 0xe3be6f68.
//
// Solidity: function runWithdraws(bytes recipient, address asset, uint256 amount, uint256 count) returns()
func (_Withdrawer *WithdrawerTransactor) RunWithdraws(opts *bind.TransactOpts, recipient []byte, asset common.Address, amount *big.Int, count *big.Int) (*types.Transaction, error) {
	return _Withdrawer.contract.Transact(opts, "runWithdraws", recipient, asset, amount, count)
}

// RunWithdraws is a paid mutator transaction binding the contract method 0xe3be6f68.
//
// Solidity: function runWithdraws(bytes recipient, address asset, uint256 amount, uint256 count) returns()
func (_Withdrawer *WithdrawerSession) RunWithdraws(recipient []byte, asset common.Address, amount *big.Int, count *big.Int) (*types.Transaction, error) {
	return _Withdrawer.Contract.RunWithdraws(&_Withdrawer.TransactOpts, recipient, asset, amount, count)
}

// RunWithdraws is a paid mutator transaction binding the contract method 0xe3be6f68.
//
// Solidity: function runWithdraws(bytes recipient, address asset, uint256 amount, uint256 count) returns()
func (_Withdrawer *WithdrawerTransactorSession) RunWithdraws(recipient []byte, asset common.Address, amount *big.Int, count *big.Int) (*types.Transaction, error) {
	return _Withdrawer.Contract.RunWithdraws(&_Withdrawer.TransactOpts, recipient, asset, amount, count)
}
