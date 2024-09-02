// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking

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

// Validator is an auto generated low-level Go binding around an user-defined struct.
type Validator struct {
	OperatorAddress string
	ConsensusPubKey string
	Jailed          bool
	BondStatus      uint8
}

// IStakingMetaData contains all meta data concerning the IStaking contract.
var IStakingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getAllValidators\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"operatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"consensusPubKey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"jailed\",\"type\":\"bool\"},{\"internalType\":\"enumBondStatus\",\"name\":\"bondStatus\",\"type\":\"uint8\"}],\"internalType\":\"structValidator[]\",\"name\":\"validators\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"getShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validatorSrc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorDst\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"moveStake\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use IStakingMetaData.ABI instead.
var IStakingABI = IStakingMetaData.ABI

// IStaking is an auto generated Go binding around an Ethereum contract.
type IStaking struct {
	IStakingCaller     // Read-only binding to the contract
	IStakingTransactor // Write-only binding to the contract
	IStakingFilterer   // Log filterer for contract events
}

// IStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type IStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IStakingSession struct {
	Contract     *IStaking         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IStakingCallerSession struct {
	Contract *IStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// IStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IStakingTransactorSession struct {
	Contract     *IStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type IStakingRaw struct {
	Contract *IStaking // Generic contract binding to access the raw methods on
}

// IStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IStakingCallerRaw struct {
	Contract *IStakingCaller // Generic read-only contract binding to access the raw methods on
}

// IStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IStakingTransactorRaw struct {
	Contract *IStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIStaking creates a new instance of IStaking, bound to a specific deployed contract.
func NewIStaking(address common.Address, backend bind.ContractBackend) (*IStaking, error) {
	contract, err := bindIStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IStaking{IStakingCaller: IStakingCaller{contract: contract}, IStakingTransactor: IStakingTransactor{contract: contract}, IStakingFilterer: IStakingFilterer{contract: contract}}, nil
}

// NewIStakingCaller creates a new read-only instance of IStaking, bound to a specific deployed contract.
func NewIStakingCaller(address common.Address, caller bind.ContractCaller) (*IStakingCaller, error) {
	contract, err := bindIStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IStakingCaller{contract: contract}, nil
}

// NewIStakingTransactor creates a new write-only instance of IStaking, bound to a specific deployed contract.
func NewIStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*IStakingTransactor, error) {
	contract, err := bindIStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IStakingTransactor{contract: contract}, nil
}

// NewIStakingFilterer creates a new log filterer instance of IStaking, bound to a specific deployed contract.
func NewIStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*IStakingFilterer, error) {
	contract, err := bindIStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IStakingFilterer{contract: contract}, nil
}

// bindIStaking binds a generic wrapper to an already deployed contract.
func bindIStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IStakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IStaking *IStakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IStaking.Contract.IStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IStaking *IStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IStaking.Contract.IStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IStaking *IStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IStaking.Contract.IStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IStaking *IStakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IStaking *IStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IStaking *IStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IStaking.Contract.contract.Transact(opts, method, params...)
}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((string,string,bool,uint8)[] validators)
func (_IStaking *IStakingCaller) GetAllValidators(opts *bind.CallOpts) ([]Validator, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "getAllValidators")

	if err != nil {
		return *new([]Validator), err
	}

	out0 := *abi.ConvertType(out[0], new([]Validator)).(*[]Validator)

	return out0, err

}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((string,string,bool,uint8)[] validators)
func (_IStaking *IStakingSession) GetAllValidators() ([]Validator, error) {
	return _IStaking.Contract.GetAllValidators(&_IStaking.CallOpts)
}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((string,string,bool,uint8)[] validators)
func (_IStaking *IStakingCallerSession) GetAllValidators() ([]Validator, error) {
	return _IStaking.Contract.GetAllValidators(&_IStaking.CallOpts)
}

// GetShares is a free data retrieval call binding the contract method 0x0d1b3daf.
//
// Solidity: function getShares(address staker, string validator) view returns(uint256 shares)
func (_IStaking *IStakingCaller) GetShares(opts *bind.CallOpts, staker common.Address, validator string) (*big.Int, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "getShares", staker, validator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetShares is a free data retrieval call binding the contract method 0x0d1b3daf.
//
// Solidity: function getShares(address staker, string validator) view returns(uint256 shares)
func (_IStaking *IStakingSession) GetShares(staker common.Address, validator string) (*big.Int, error) {
	return _IStaking.Contract.GetShares(&_IStaking.CallOpts, staker, validator)
}

// GetShares is a free data retrieval call binding the contract method 0x0d1b3daf.
//
// Solidity: function getShares(address staker, string validator) view returns(uint256 shares)
func (_IStaking *IStakingCallerSession) GetShares(staker common.Address, validator string) (*big.Int, error) {
	return _IStaking.Contract.GetShares(&_IStaking.CallOpts, staker, validator)
}

// MoveStake is a paid mutator transaction binding the contract method 0xd11a93d0.
//
// Solidity: function moveStake(address staker, string validatorSrc, string validatorDst, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingTransactor) MoveStake(opts *bind.TransactOpts, staker common.Address, validatorSrc string, validatorDst string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "moveStake", staker, validatorSrc, validatorDst, amount)
}

// MoveStake is a paid mutator transaction binding the contract method 0xd11a93d0.
//
// Solidity: function moveStake(address staker, string validatorSrc, string validatorDst, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingSession) MoveStake(staker common.Address, validatorSrc string, validatorDst string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.MoveStake(&_IStaking.TransactOpts, staker, validatorSrc, validatorDst, amount)
}

// MoveStake is a paid mutator transaction binding the contract method 0xd11a93d0.
//
// Solidity: function moveStake(address staker, string validatorSrc, string validatorDst, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingTransactorSession) MoveStake(staker common.Address, validatorSrc string, validatorDst string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.MoveStake(&_IStaking.TransactOpts, staker, validatorSrc, validatorDst, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x90b8436f.
//
// Solidity: function stake(address staker, string validator, uint256 amount) returns(bool success)
func (_IStaking *IStakingTransactor) Stake(opts *bind.TransactOpts, staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "stake", staker, validator, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x90b8436f.
//
// Solidity: function stake(address staker, string validator, uint256 amount) returns(bool success)
func (_IStaking *IStakingSession) Stake(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Stake(&_IStaking.TransactOpts, staker, validator, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x90b8436f.
//
// Solidity: function stake(address staker, string validator, uint256 amount) returns(bool success)
func (_IStaking *IStakingTransactorSession) Stake(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Stake(&_IStaking.TransactOpts, staker, validator, amount)
}

// Unstake is a paid mutator transaction binding the contract method 0x57c6ea3e.
//
// Solidity: function unstake(address staker, string validator, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingTransactor) Unstake(opts *bind.TransactOpts, staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "unstake", staker, validator, amount)
}

// Unstake is a paid mutator transaction binding the contract method 0x57c6ea3e.
//
// Solidity: function unstake(address staker, string validator, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingSession) Unstake(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Unstake(&_IStaking.TransactOpts, staker, validator, amount)
}

// Unstake is a paid mutator transaction binding the contract method 0x57c6ea3e.
//
// Solidity: function unstake(address staker, string validator, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingTransactorSession) Unstake(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Unstake(&_IStaking.TransactOpts, staker, validator, amount)
}
