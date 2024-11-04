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

// DecCoin is an auto generated low-level Go binding around an user-defined struct.
type DecCoin struct {
	Denom  string
	Amount *big.Int
}

// Validator is an auto generated low-level Go binding around an user-defined struct.
type Validator struct {
	OperatorAddress string
	ConsensusPubKey string
	Jailed          bool
	BondStatus      uint8
}

// IStakingMetaData contains all meta data concerning the IStaking contract.
var IStakingMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"claim_address\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zrc20_token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ClaimedRewards\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zrc20_distributor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zrc20_token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Distributed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorSrc\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorDst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MoveStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Unstake\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"claimRewards\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"distribute\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllValidators\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"operatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"consensusPubKey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"jailed\",\"type\":\"bool\"},{\"internalType\":\"enumBondStatus\",\"name\":\"bondStatus\",\"type\":\"uint8\"}],\"internalType\":\"structValidator[]\",\"name\":\"validators\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"}],\"name\":\"getDelegatorValidators\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"validators\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"getRewards\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structDecCoin[]\",\"name\":\"rewards\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"getShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validatorSrc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorDst\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"moveStake\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// GetDelegatorValidators is a free data retrieval call binding the contract method 0xb6a216ae.
//
// Solidity: function getDelegatorValidators(address delegator) view returns(string[] validators)
func (_IStaking *IStakingCaller) GetDelegatorValidators(opts *bind.CallOpts, delegator common.Address) ([]string, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "getDelegatorValidators", delegator)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetDelegatorValidators is a free data retrieval call binding the contract method 0xb6a216ae.
//
// Solidity: function getDelegatorValidators(address delegator) view returns(string[] validators)
func (_IStaking *IStakingSession) GetDelegatorValidators(delegator common.Address) ([]string, error) {
	return _IStaking.Contract.GetDelegatorValidators(&_IStaking.CallOpts, delegator)
}

// GetDelegatorValidators is a free data retrieval call binding the contract method 0xb6a216ae.
//
// Solidity: function getDelegatorValidators(address delegator) view returns(string[] validators)
func (_IStaking *IStakingCallerSession) GetDelegatorValidators(delegator common.Address) ([]string, error) {
	return _IStaking.Contract.GetDelegatorValidators(&_IStaking.CallOpts, delegator)
}

// GetRewards is a free data retrieval call binding the contract method 0x93428792.
//
// Solidity: function getRewards(address delegator, string validator) view returns((string,uint256)[] rewards)
func (_IStaking *IStakingCaller) GetRewards(opts *bind.CallOpts, delegator common.Address, validator string) ([]DecCoin, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "getRewards", delegator, validator)

	if err != nil {
		return *new([]DecCoin), err
	}

	out0 := *abi.ConvertType(out[0], new([]DecCoin)).(*[]DecCoin)

	return out0, err

}

// GetRewards is a free data retrieval call binding the contract method 0x93428792.
//
// Solidity: function getRewards(address delegator, string validator) view returns((string,uint256)[] rewards)
func (_IStaking *IStakingSession) GetRewards(delegator common.Address, validator string) ([]DecCoin, error) {
	return _IStaking.Contract.GetRewards(&_IStaking.CallOpts, delegator, validator)
}

// GetRewards is a free data retrieval call binding the contract method 0x93428792.
//
// Solidity: function getRewards(address delegator, string validator) view returns((string,uint256)[] rewards)
func (_IStaking *IStakingCallerSession) GetRewards(delegator common.Address, validator string) ([]DecCoin, error) {
	return _IStaking.Contract.GetRewards(&_IStaking.CallOpts, delegator, validator)
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

// ClaimRewards is a paid mutator transaction binding the contract method 0x54dbdc38.
//
// Solidity: function claimRewards(address delegator, string validator) returns(bool success)
func (_IStaking *IStakingTransactor) ClaimRewards(opts *bind.TransactOpts, delegator common.Address, validator string) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "claimRewards", delegator, validator)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x54dbdc38.
//
// Solidity: function claimRewards(address delegator, string validator) returns(bool success)
func (_IStaking *IStakingSession) ClaimRewards(delegator common.Address, validator string) (*types.Transaction, error) {
	return _IStaking.Contract.ClaimRewards(&_IStaking.TransactOpts, delegator, validator)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x54dbdc38.
//
// Solidity: function claimRewards(address delegator, string validator) returns(bool success)
func (_IStaking *IStakingTransactorSession) ClaimRewards(delegator common.Address, validator string) (*types.Transaction, error) {
	return _IStaking.Contract.ClaimRewards(&_IStaking.TransactOpts, delegator, validator)
}

// Distribute is a paid mutator transaction binding the contract method 0xfb932108.
//
// Solidity: function distribute(address zrc20, uint256 amount) returns(bool success)
func (_IStaking *IStakingTransactor) Distribute(opts *bind.TransactOpts, zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "distribute", zrc20, amount)
}

// Distribute is a paid mutator transaction binding the contract method 0xfb932108.
//
// Solidity: function distribute(address zrc20, uint256 amount) returns(bool success)
func (_IStaking *IStakingSession) Distribute(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Distribute(&_IStaking.TransactOpts, zrc20, amount)
}

// Distribute is a paid mutator transaction binding the contract method 0xfb932108.
//
// Solidity: function distribute(address zrc20, uint256 amount) returns(bool success)
func (_IStaking *IStakingTransactorSession) Distribute(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Distribute(&_IStaking.TransactOpts, zrc20, amount)
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

// IStakingClaimedRewardsIterator is returned from FilterClaimedRewards and is used to iterate over the raw logs and unpacked data for ClaimedRewards events raised by the IStaking contract.
type IStakingClaimedRewardsIterator struct {
	Event *IStakingClaimedRewards // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IStakingClaimedRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingClaimedRewards)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IStakingClaimedRewards)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IStakingClaimedRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingClaimedRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingClaimedRewards represents a ClaimedRewards event raised by the IStaking contract.
type IStakingClaimedRewards struct {
	ClaimAddress common.Address
	Zrc20Token   common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterClaimedRewards is a free log retrieval operation binding the contract event 0x2ef606d064225d24c1514dc94907c134faee1237445c2f63f410cce0852b2054.
//
// Solidity: event ClaimedRewards(address indexed claim_address, address indexed zrc20_token, uint256 amount)
func (_IStaking *IStakingFilterer) FilterClaimedRewards(opts *bind.FilterOpts, claim_address []common.Address, zrc20_token []common.Address) (*IStakingClaimedRewardsIterator, error) {

	var claim_addressRule []interface{}
	for _, claim_addressItem := range claim_address {
		claim_addressRule = append(claim_addressRule, claim_addressItem)
	}
	var zrc20_tokenRule []interface{}
	for _, zrc20_tokenItem := range zrc20_token {
		zrc20_tokenRule = append(zrc20_tokenRule, zrc20_tokenItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "ClaimedRewards", claim_addressRule, zrc20_tokenRule)
	if err != nil {
		return nil, err
	}
	return &IStakingClaimedRewardsIterator{contract: _IStaking.contract, event: "ClaimedRewards", logs: logs, sub: sub}, nil
}

// WatchClaimedRewards is a free log subscription operation binding the contract event 0x2ef606d064225d24c1514dc94907c134faee1237445c2f63f410cce0852b2054.
//
// Solidity: event ClaimedRewards(address indexed claim_address, address indexed zrc20_token, uint256 amount)
func (_IStaking *IStakingFilterer) WatchClaimedRewards(opts *bind.WatchOpts, sink chan<- *IStakingClaimedRewards, claim_address []common.Address, zrc20_token []common.Address) (event.Subscription, error) {

	var claim_addressRule []interface{}
	for _, claim_addressItem := range claim_address {
		claim_addressRule = append(claim_addressRule, claim_addressItem)
	}
	var zrc20_tokenRule []interface{}
	for _, zrc20_tokenItem := range zrc20_token {
		zrc20_tokenRule = append(zrc20_tokenRule, zrc20_tokenItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "ClaimedRewards", claim_addressRule, zrc20_tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingClaimedRewards)
				if err := _IStaking.contract.UnpackLog(event, "ClaimedRewards", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClaimedRewards is a log parse operation binding the contract event 0x2ef606d064225d24c1514dc94907c134faee1237445c2f63f410cce0852b2054.
//
// Solidity: event ClaimedRewards(address indexed claim_address, address indexed zrc20_token, uint256 amount)
func (_IStaking *IStakingFilterer) ParseClaimedRewards(log types.Log) (*IStakingClaimedRewards, error) {
	event := new(IStakingClaimedRewards)
	if err := _IStaking.contract.UnpackLog(event, "ClaimedRewards", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingDistributedIterator is returned from FilterDistributed and is used to iterate over the raw logs and unpacked data for Distributed events raised by the IStaking contract.
type IStakingDistributedIterator struct {
	Event *IStakingDistributed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IStakingDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingDistributed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IStakingDistributed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IStakingDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingDistributed represents a Distributed event raised by the IStaking contract.
type IStakingDistributed struct {
	Zrc20Distributor common.Address
	Zrc20Token       common.Address
	Amount           *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDistributed is a free log retrieval operation binding the contract event 0xad4a9acf26d8bba7a8cf1a41160d59be042ee554578e256c98d2ab74cdd43542.
//
// Solidity: event Distributed(address indexed zrc20_distributor, address indexed zrc20_token, uint256 amount)
func (_IStaking *IStakingFilterer) FilterDistributed(opts *bind.FilterOpts, zrc20_distributor []common.Address, zrc20_token []common.Address) (*IStakingDistributedIterator, error) {

	var zrc20_distributorRule []interface{}
	for _, zrc20_distributorItem := range zrc20_distributor {
		zrc20_distributorRule = append(zrc20_distributorRule, zrc20_distributorItem)
	}
	var zrc20_tokenRule []interface{}
	for _, zrc20_tokenItem := range zrc20_token {
		zrc20_tokenRule = append(zrc20_tokenRule, zrc20_tokenItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "Distributed", zrc20_distributorRule, zrc20_tokenRule)
	if err != nil {
		return nil, err
	}
	return &IStakingDistributedIterator{contract: _IStaking.contract, event: "Distributed", logs: logs, sub: sub}, nil
}

// WatchDistributed is a free log subscription operation binding the contract event 0xad4a9acf26d8bba7a8cf1a41160d59be042ee554578e256c98d2ab74cdd43542.
//
// Solidity: event Distributed(address indexed zrc20_distributor, address indexed zrc20_token, uint256 amount)
func (_IStaking *IStakingFilterer) WatchDistributed(opts *bind.WatchOpts, sink chan<- *IStakingDistributed, zrc20_distributor []common.Address, zrc20_token []common.Address) (event.Subscription, error) {

	var zrc20_distributorRule []interface{}
	for _, zrc20_distributorItem := range zrc20_distributor {
		zrc20_distributorRule = append(zrc20_distributorRule, zrc20_distributorItem)
	}
	var zrc20_tokenRule []interface{}
	for _, zrc20_tokenItem := range zrc20_token {
		zrc20_tokenRule = append(zrc20_tokenRule, zrc20_tokenItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "Distributed", zrc20_distributorRule, zrc20_tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingDistributed)
				if err := _IStaking.contract.UnpackLog(event, "Distributed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDistributed is a log parse operation binding the contract event 0xad4a9acf26d8bba7a8cf1a41160d59be042ee554578e256c98d2ab74cdd43542.
//
// Solidity: event Distributed(address indexed zrc20_distributor, address indexed zrc20_token, uint256 amount)
func (_IStaking *IStakingFilterer) ParseDistributed(log types.Log) (*IStakingDistributed, error) {
	event := new(IStakingDistributed)
	if err := _IStaking.contract.UnpackLog(event, "Distributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingMoveStakeIterator is returned from FilterMoveStake and is used to iterate over the raw logs and unpacked data for MoveStake events raised by the IStaking contract.
type IStakingMoveStakeIterator struct {
	Event *IStakingMoveStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IStakingMoveStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingMoveStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IStakingMoveStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IStakingMoveStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingMoveStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingMoveStake represents a MoveStake event raised by the IStaking contract.
type IStakingMoveStake struct {
	Staker       common.Address
	ValidatorSrc common.Address
	ValidatorDst common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMoveStake is a free log retrieval operation binding the contract event 0x4dda2c731d442025256e6e47fbb109592bcd8baf3cf25996ebd09f1da7ec902b.
//
// Solidity: event MoveStake(address indexed staker, address indexed validatorSrc, address indexed validatorDst, uint256 amount)
func (_IStaking *IStakingFilterer) FilterMoveStake(opts *bind.FilterOpts, staker []common.Address, validatorSrc []common.Address, validatorDst []common.Address) (*IStakingMoveStakeIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorSrcRule []interface{}
	for _, validatorSrcItem := range validatorSrc {
		validatorSrcRule = append(validatorSrcRule, validatorSrcItem)
	}
	var validatorDstRule []interface{}
	for _, validatorDstItem := range validatorDst {
		validatorDstRule = append(validatorDstRule, validatorDstItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "MoveStake", stakerRule, validatorSrcRule, validatorDstRule)
	if err != nil {
		return nil, err
	}
	return &IStakingMoveStakeIterator{contract: _IStaking.contract, event: "MoveStake", logs: logs, sub: sub}, nil
}

// WatchMoveStake is a free log subscription operation binding the contract event 0x4dda2c731d442025256e6e47fbb109592bcd8baf3cf25996ebd09f1da7ec902b.
//
// Solidity: event MoveStake(address indexed staker, address indexed validatorSrc, address indexed validatorDst, uint256 amount)
func (_IStaking *IStakingFilterer) WatchMoveStake(opts *bind.WatchOpts, sink chan<- *IStakingMoveStake, staker []common.Address, validatorSrc []common.Address, validatorDst []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorSrcRule []interface{}
	for _, validatorSrcItem := range validatorSrc {
		validatorSrcRule = append(validatorSrcRule, validatorSrcItem)
	}
	var validatorDstRule []interface{}
	for _, validatorDstItem := range validatorDst {
		validatorDstRule = append(validatorDstRule, validatorDstItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "MoveStake", stakerRule, validatorSrcRule, validatorDstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingMoveStake)
				if err := _IStaking.contract.UnpackLog(event, "MoveStake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMoveStake is a log parse operation binding the contract event 0x4dda2c731d442025256e6e47fbb109592bcd8baf3cf25996ebd09f1da7ec902b.
//
// Solidity: event MoveStake(address indexed staker, address indexed validatorSrc, address indexed validatorDst, uint256 amount)
func (_IStaking *IStakingFilterer) ParseMoveStake(log types.Log) (*IStakingMoveStake, error) {
	event := new(IStakingMoveStake)
	if err := _IStaking.contract.UnpackLog(event, "MoveStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the IStaking contract.
type IStakingStakeIterator struct {
	Event *IStakingStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IStakingStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IStakingStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IStakingStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingStake represents a Stake event raised by the IStaking contract.
type IStakingStake struct {
	Staker    common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStake is a free log retrieval operation binding the contract event 0x99039fcf0a98f484616c5196ee8b2ecfa971babf0b519848289ea4db381f85f7.
//
// Solidity: event Stake(address indexed staker, address indexed validator, uint256 amount)
func (_IStaking *IStakingFilterer) FilterStake(opts *bind.FilterOpts, staker []common.Address, validator []common.Address) (*IStakingStakeIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "Stake", stakerRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &IStakingStakeIterator{contract: _IStaking.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0x99039fcf0a98f484616c5196ee8b2ecfa971babf0b519848289ea4db381f85f7.
//
// Solidity: event Stake(address indexed staker, address indexed validator, uint256 amount)
func (_IStaking *IStakingFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *IStakingStake, staker []common.Address, validator []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "Stake", stakerRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingStake)
				if err := _IStaking.contract.UnpackLog(event, "Stake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStake is a log parse operation binding the contract event 0x99039fcf0a98f484616c5196ee8b2ecfa971babf0b519848289ea4db381f85f7.
//
// Solidity: event Stake(address indexed staker, address indexed validator, uint256 amount)
func (_IStaking *IStakingFilterer) ParseStake(log types.Log) (*IStakingStake, error) {
	event := new(IStakingStake)
	if err := _IStaking.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingUnstakeIterator is returned from FilterUnstake and is used to iterate over the raw logs and unpacked data for Unstake events raised by the IStaking contract.
type IStakingUnstakeIterator struct {
	Event *IStakingUnstake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IStakingUnstakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingUnstake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IStakingUnstake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IStakingUnstakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingUnstakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingUnstake represents a Unstake event raised by the IStaking contract.
type IStakingUnstake struct {
	Staker    common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnstake is a free log retrieval operation binding the contract event 0x390b1276974b9463e5d66ab10df69b6f3d7b930eb066a0e66df327edd2cc811c.
//
// Solidity: event Unstake(address indexed staker, address indexed validator, uint256 amount)
func (_IStaking *IStakingFilterer) FilterUnstake(opts *bind.FilterOpts, staker []common.Address, validator []common.Address) (*IStakingUnstakeIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "Unstake", stakerRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &IStakingUnstakeIterator{contract: _IStaking.contract, event: "Unstake", logs: logs, sub: sub}, nil
}

// WatchUnstake is a free log subscription operation binding the contract event 0x390b1276974b9463e5d66ab10df69b6f3d7b930eb066a0e66df327edd2cc811c.
//
// Solidity: event Unstake(address indexed staker, address indexed validator, uint256 amount)
func (_IStaking *IStakingFilterer) WatchUnstake(opts *bind.WatchOpts, sink chan<- *IStakingUnstake, staker []common.Address, validator []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "Unstake", stakerRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingUnstake)
				if err := _IStaking.contract.UnpackLog(event, "Unstake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnstake is a log parse operation binding the contract event 0x390b1276974b9463e5d66ab10df69b6f3d7b930eb066a0e66df327edd2cc811c.
//
// Solidity: event Unstake(address indexed staker, address indexed validator, uint256 amount)
func (_IStaking *IStakingFilterer) ParseUnstake(log types.Log) (*IStakingUnstake, error) {
	event := new(IStakingUnstake)
	if err := _IStaking.contract.UnpackLog(event, "Unstake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
