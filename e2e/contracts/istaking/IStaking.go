// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package istaking

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

// Coin is an auto generated low-level Go binding around an user-defined struct.
type Coin struct {
	Denom  string
	Amount *big.Int
}

// CommissionRates is an auto generated low-level Go binding around an user-defined struct.
type CommissionRates struct {
	Rate          *big.Int
	MaxRate       *big.Int
	MaxChangeRate *big.Int
}

// Description is an auto generated low-level Go binding around an user-defined struct.
type Description struct {
	Moniker         string
	Identity        string
	Website         string
	SecurityContact string
	Details         string
}

// PageRequest is an auto generated low-level Go binding around an user-defined struct.
type PageRequest struct {
	Key        []byte
	Offset     uint64
	Limit      uint64
	CountTotal bool
	Reverse    bool
}

// PageResponse is an auto generated low-level Go binding around an user-defined struct.
type PageResponse struct {
	NextKey []byte
	Total   uint64
}

// Redelegation is an auto generated low-level Go binding around an user-defined struct.
type Redelegation struct {
	DelegatorAddress    string
	ValidatorSrcAddress string
	ValidatorDstAddress string
	Entries             []RedelegationEntry
}

// RedelegationEntry is an auto generated low-level Go binding around an user-defined struct.
type RedelegationEntry struct {
	CreationHeight int64
	CompletionTime int64
	InitialBalance *big.Int
	SharesDst      *big.Int
}

// RedelegationEntryResponse is an auto generated low-level Go binding around an user-defined struct.
type RedelegationEntryResponse struct {
	RedelegationEntry RedelegationEntry
	Balance           *big.Int
}

// RedelegationOutput is an auto generated low-level Go binding around an user-defined struct.
type RedelegationOutput struct {
	DelegatorAddress    string
	ValidatorSrcAddress string
	ValidatorDstAddress string
	Entries             []RedelegationEntry
}

// RedelegationResponse is an auto generated low-level Go binding around an user-defined struct.
type RedelegationResponse struct {
	Redelegation Redelegation
	Entries      []RedelegationEntryResponse
}

// UnbondingDelegationEntry is an auto generated low-level Go binding around an user-defined struct.
type UnbondingDelegationEntry struct {
	CreationHeight          int64
	CompletionTime          int64
	InitialBalance          *big.Int
	Balance                 *big.Int
	UnbondingId             uint64
	UnbondingOnHoldRefCount int64
}

// UnbondingDelegationOutput is an auto generated low-level Go binding around an user-defined struct.
type UnbondingDelegationOutput struct {
	DelegatorAddress string
	ValidatorAddress string
	Entries          []UnbondingDelegationEntry
}

// Validator is an auto generated low-level Go binding around an user-defined struct.
type Validator struct {
	OperatorAddress   string
	ConsensusPubkey   string
	Jailed            bool
	Status            uint8
	Tokens            *big.Int
	DelegatorShares   *big.Int
	Description       string
	UnbondingHeight   int64
	UnbondingTime     int64
	Commission        *big.Int
	MinSelfDelegation *big.Int
}

// IStakingMetaData contains all meta data concerning the IStaking contract.
var IStakingMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"creationHeight\",\"type\":\"uint256\"}],\"name\":\"CancelUnbondingDelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"CreateValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newShares\",\"type\":\"uint256\"}],\"name\":\"Delegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"commissionRate\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"minSelfDelegation\",\"type\":\"int256\"}],\"name\":\"EditValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorSrcAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorDstAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"Redelegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"Unbond\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validatorAddress\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"creationHeight\",\"type\":\"uint256\"}],\"name\":\"cancelUnbondingDelegation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"moniker\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"website\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"securityContact\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"details\",\"type\":\"string\"}],\"internalType\":\"structDescription\",\"name\":\"description\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxChangeRate\",\"type\":\"uint256\"}],\"internalType\":\"structCommissionRates\",\"name\":\"commissionRates\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"minSelfDelegation\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"pubkey\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"createValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validatorAddress\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"delegate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validatorAddress\",\"type\":\"string\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCoin\",\"name\":\"balance\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"moniker\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"website\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"securityContact\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"details\",\"type\":\"string\"}],\"internalType\":\"structDescription\",\"name\":\"description\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"int256\",\"name\":\"commissionRate\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"minSelfDelegation\",\"type\":\"int256\"}],\"name\":\"editValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validatorSrcAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorDstAddress\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"redelegate\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"srcValidatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dstValidatorAddress\",\"type\":\"string\"}],\"name\":\"redelegation\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"delegatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorSrcAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorDstAddress\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"creationHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sharesDst\",\"type\":\"uint256\"}],\"internalType\":\"structRedelegationEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"internalType\":\"structRedelegationOutput\",\"name\":\"redelegation\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"srcValidatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dstValidatorAddress\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offset\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"limit\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"countTotal\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reverse\",\"type\":\"bool\"}],\"internalType\":\"structPageRequest\",\"name\":\"pageRequest\",\"type\":\"tuple\"}],\"name\":\"redelegations\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"delegatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorSrcAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorDstAddress\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"creationHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sharesDst\",\"type\":\"uint256\"}],\"internalType\":\"structRedelegationEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"internalType\":\"structRedelegation\",\"name\":\"redelegation\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"int64\",\"name\":\"creationHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"sharesDst\",\"type\":\"uint256\"}],\"internalType\":\"structRedelegationEntry\",\"name\":\"redelegationEntry\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structRedelegationEntryResponse[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"internalType\":\"structRedelegationResponse[]\",\"name\":\"response\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"nextKey\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"total\",\"type\":\"uint64\"}],\"internalType\":\"structPageResponse\",\"name\":\"pageResponse\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validatorAddress\",\"type\":\"string\"}],\"name\":\"unbondingDelegation\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"delegatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorAddress\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"int64\",\"name\":\"creationHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"initialBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"unbondingId\",\"type\":\"uint64\"},{\"internalType\":\"int64\",\"name\":\"unbondingOnHoldRefCount\",\"type\":\"int64\"}],\"internalType\":\"structUnbondingDelegationEntry[]\",\"name\":\"entries\",\"type\":\"tuple[]\"}],\"internalType\":\"structUnbondingDelegationOutput\",\"name\":\"unbondingDelegation\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegatorAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validatorAddress\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"name\":\"validator\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"operatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"jailed\",\"type\":\"bool\"},{\"internalType\":\"enumBondStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegatorShares\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"int64\",\"name\":\"unbondingHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"unbondingTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"commission\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSelfDelegation\",\"type\":\"uint256\"}],\"internalType\":\"structValidator\",\"name\":\"validator\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"status\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"offset\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"limit\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"countTotal\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"reverse\",\"type\":\"bool\"}],\"internalType\":\"structPageRequest\",\"name\":\"pageRequest\",\"type\":\"tuple\"}],\"name\":\"validators\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"operatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"consensusPubkey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"jailed\",\"type\":\"bool\"},{\"internalType\":\"enumBondStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegatorShares\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"int64\",\"name\":\"unbondingHeight\",\"type\":\"int64\"},{\"internalType\":\"int64\",\"name\":\"unbondingTime\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"commission\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSelfDelegation\",\"type\":\"uint256\"}],\"internalType\":\"structValidator[]\",\"name\":\"validators\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"nextKey\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"total\",\"type\":\"uint64\"}],\"internalType\":\"structPageResponse\",\"name\":\"pageResponse\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// Delegation is a free data retrieval call binding the contract method 0x241774e6.
//
// Solidity: function delegation(address delegatorAddress, string validatorAddress) view returns(uint256 shares, (string,uint256) balance)
func (_IStaking *IStakingCaller) Delegation(opts *bind.CallOpts, delegatorAddress common.Address, validatorAddress string) (struct {
	Shares  *big.Int
	Balance Coin
}, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "delegation", delegatorAddress, validatorAddress)

	outstruct := new(struct {
		Shares  *big.Int
		Balance Coin
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Shares = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Balance = *abi.ConvertType(out[1], new(Coin)).(*Coin)

	return *outstruct, err

}

// Delegation is a free data retrieval call binding the contract method 0x241774e6.
//
// Solidity: function delegation(address delegatorAddress, string validatorAddress) view returns(uint256 shares, (string,uint256) balance)
func (_IStaking *IStakingSession) Delegation(delegatorAddress common.Address, validatorAddress string) (struct {
	Shares  *big.Int
	Balance Coin
}, error) {
	return _IStaking.Contract.Delegation(&_IStaking.CallOpts, delegatorAddress, validatorAddress)
}

// Delegation is a free data retrieval call binding the contract method 0x241774e6.
//
// Solidity: function delegation(address delegatorAddress, string validatorAddress) view returns(uint256 shares, (string,uint256) balance)
func (_IStaking *IStakingCallerSession) Delegation(delegatorAddress common.Address, validatorAddress string) (struct {
	Shares  *big.Int
	Balance Coin
}, error) {
	return _IStaking.Contract.Delegation(&_IStaking.CallOpts, delegatorAddress, validatorAddress)
}

// Redelegation is a free data retrieval call binding the contract method 0x7d9f939c.
//
// Solidity: function redelegation(address delegatorAddress, string srcValidatorAddress, string dstValidatorAddress) view returns((string,string,string,(int64,int64,uint256,uint256)[]) redelegation)
func (_IStaking *IStakingCaller) Redelegation(opts *bind.CallOpts, delegatorAddress common.Address, srcValidatorAddress string, dstValidatorAddress string) (RedelegationOutput, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "redelegation", delegatorAddress, srcValidatorAddress, dstValidatorAddress)

	if err != nil {
		return *new(RedelegationOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(RedelegationOutput)).(*RedelegationOutput)

	return out0, err

}

// Redelegation is a free data retrieval call binding the contract method 0x7d9f939c.
//
// Solidity: function redelegation(address delegatorAddress, string srcValidatorAddress, string dstValidatorAddress) view returns((string,string,string,(int64,int64,uint256,uint256)[]) redelegation)
func (_IStaking *IStakingSession) Redelegation(delegatorAddress common.Address, srcValidatorAddress string, dstValidatorAddress string) (RedelegationOutput, error) {
	return _IStaking.Contract.Redelegation(&_IStaking.CallOpts, delegatorAddress, srcValidatorAddress, dstValidatorAddress)
}

// Redelegation is a free data retrieval call binding the contract method 0x7d9f939c.
//
// Solidity: function redelegation(address delegatorAddress, string srcValidatorAddress, string dstValidatorAddress) view returns((string,string,string,(int64,int64,uint256,uint256)[]) redelegation)
func (_IStaking *IStakingCallerSession) Redelegation(delegatorAddress common.Address, srcValidatorAddress string, dstValidatorAddress string) (RedelegationOutput, error) {
	return _IStaking.Contract.Redelegation(&_IStaking.CallOpts, delegatorAddress, srcValidatorAddress, dstValidatorAddress)
}

// Redelegations is a free data retrieval call binding the contract method 0x10a2851c.
//
// Solidity: function redelegations(address delegatorAddress, string srcValidatorAddress, string dstValidatorAddress, (bytes,uint64,uint64,bool,bool) pageRequest) view returns(((string,string,string,(int64,int64,uint256,uint256)[]),((int64,int64,uint256,uint256),uint256)[])[] response, (bytes,uint64) pageResponse)
func (_IStaking *IStakingCaller) Redelegations(opts *bind.CallOpts, delegatorAddress common.Address, srcValidatorAddress string, dstValidatorAddress string, pageRequest PageRequest) (struct {
	Response     []RedelegationResponse
	PageResponse PageResponse
}, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "redelegations", delegatorAddress, srcValidatorAddress, dstValidatorAddress, pageRequest)

	outstruct := new(struct {
		Response     []RedelegationResponse
		PageResponse PageResponse
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Response = *abi.ConvertType(out[0], new([]RedelegationResponse)).(*[]RedelegationResponse)
	outstruct.PageResponse = *abi.ConvertType(out[1], new(PageResponse)).(*PageResponse)

	return *outstruct, err

}

// Redelegations is a free data retrieval call binding the contract method 0x10a2851c.
//
// Solidity: function redelegations(address delegatorAddress, string srcValidatorAddress, string dstValidatorAddress, (bytes,uint64,uint64,bool,bool) pageRequest) view returns(((string,string,string,(int64,int64,uint256,uint256)[]),((int64,int64,uint256,uint256),uint256)[])[] response, (bytes,uint64) pageResponse)
func (_IStaking *IStakingSession) Redelegations(delegatorAddress common.Address, srcValidatorAddress string, dstValidatorAddress string, pageRequest PageRequest) (struct {
	Response     []RedelegationResponse
	PageResponse PageResponse
}, error) {
	return _IStaking.Contract.Redelegations(&_IStaking.CallOpts, delegatorAddress, srcValidatorAddress, dstValidatorAddress, pageRequest)
}

// Redelegations is a free data retrieval call binding the contract method 0x10a2851c.
//
// Solidity: function redelegations(address delegatorAddress, string srcValidatorAddress, string dstValidatorAddress, (bytes,uint64,uint64,bool,bool) pageRequest) view returns(((string,string,string,(int64,int64,uint256,uint256)[]),((int64,int64,uint256,uint256),uint256)[])[] response, (bytes,uint64) pageResponse)
func (_IStaking *IStakingCallerSession) Redelegations(delegatorAddress common.Address, srcValidatorAddress string, dstValidatorAddress string, pageRequest PageRequest) (struct {
	Response     []RedelegationResponse
	PageResponse PageResponse
}, error) {
	return _IStaking.Contract.Redelegations(&_IStaking.CallOpts, delegatorAddress, srcValidatorAddress, dstValidatorAddress, pageRequest)
}

// UnbondingDelegation is a free data retrieval call binding the contract method 0xa03ffee1.
//
// Solidity: function unbondingDelegation(address delegatorAddress, string validatorAddress) view returns((string,string,(int64,int64,uint256,uint256,uint64,int64)[]) unbondingDelegation)
func (_IStaking *IStakingCaller) UnbondingDelegation(opts *bind.CallOpts, delegatorAddress common.Address, validatorAddress string) (UnbondingDelegationOutput, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "unbondingDelegation", delegatorAddress, validatorAddress)

	if err != nil {
		return *new(UnbondingDelegationOutput), err
	}

	out0 := *abi.ConvertType(out[0], new(UnbondingDelegationOutput)).(*UnbondingDelegationOutput)

	return out0, err

}

// UnbondingDelegation is a free data retrieval call binding the contract method 0xa03ffee1.
//
// Solidity: function unbondingDelegation(address delegatorAddress, string validatorAddress) view returns((string,string,(int64,int64,uint256,uint256,uint64,int64)[]) unbondingDelegation)
func (_IStaking *IStakingSession) UnbondingDelegation(delegatorAddress common.Address, validatorAddress string) (UnbondingDelegationOutput, error) {
	return _IStaking.Contract.UnbondingDelegation(&_IStaking.CallOpts, delegatorAddress, validatorAddress)
}

// UnbondingDelegation is a free data retrieval call binding the contract method 0xa03ffee1.
//
// Solidity: function unbondingDelegation(address delegatorAddress, string validatorAddress) view returns((string,string,(int64,int64,uint256,uint256,uint64,int64)[]) unbondingDelegation)
func (_IStaking *IStakingCallerSession) UnbondingDelegation(delegatorAddress common.Address, validatorAddress string) (UnbondingDelegationOutput, error) {
	return _IStaking.Contract.UnbondingDelegation(&_IStaking.CallOpts, delegatorAddress, validatorAddress)
}

// Validator is a free data retrieval call binding the contract method 0x223b3b7a.
//
// Solidity: function validator(address validatorAddress) view returns((string,string,bool,uint8,uint256,uint256,string,int64,int64,uint256,uint256) validator)
func (_IStaking *IStakingCaller) Validator(opts *bind.CallOpts, validatorAddress common.Address) (Validator, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "validator", validatorAddress)

	if err != nil {
		return *new(Validator), err
	}

	out0 := *abi.ConvertType(out[0], new(Validator)).(*Validator)

	return out0, err

}

// Validator is a free data retrieval call binding the contract method 0x223b3b7a.
//
// Solidity: function validator(address validatorAddress) view returns((string,string,bool,uint8,uint256,uint256,string,int64,int64,uint256,uint256) validator)
func (_IStaking *IStakingSession) Validator(validatorAddress common.Address) (Validator, error) {
	return _IStaking.Contract.Validator(&_IStaking.CallOpts, validatorAddress)
}

// Validator is a free data retrieval call binding the contract method 0x223b3b7a.
//
// Solidity: function validator(address validatorAddress) view returns((string,string,bool,uint8,uint256,uint256,string,int64,int64,uint256,uint256) validator)
func (_IStaking *IStakingCallerSession) Validator(validatorAddress common.Address) (Validator, error) {
	return _IStaking.Contract.Validator(&_IStaking.CallOpts, validatorAddress)
}

// Validators is a free data retrieval call binding the contract method 0x186b2167.
//
// Solidity: function validators(string status, (bytes,uint64,uint64,bool,bool) pageRequest) view returns((string,string,bool,uint8,uint256,uint256,string,int64,int64,uint256,uint256)[] validators, (bytes,uint64) pageResponse)
func (_IStaking *IStakingCaller) Validators(opts *bind.CallOpts, status string, pageRequest PageRequest) (struct {
	Validators   []Validator
	PageResponse PageResponse
}, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "validators", status, pageRequest)

	outstruct := new(struct {
		Validators   []Validator
		PageResponse PageResponse
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Validators = *abi.ConvertType(out[0], new([]Validator)).(*[]Validator)
	outstruct.PageResponse = *abi.ConvertType(out[1], new(PageResponse)).(*PageResponse)

	return *outstruct, err

}

// Validators is a free data retrieval call binding the contract method 0x186b2167.
//
// Solidity: function validators(string status, (bytes,uint64,uint64,bool,bool) pageRequest) view returns((string,string,bool,uint8,uint256,uint256,string,int64,int64,uint256,uint256)[] validators, (bytes,uint64) pageResponse)
func (_IStaking *IStakingSession) Validators(status string, pageRequest PageRequest) (struct {
	Validators   []Validator
	PageResponse PageResponse
}, error) {
	return _IStaking.Contract.Validators(&_IStaking.CallOpts, status, pageRequest)
}

// Validators is a free data retrieval call binding the contract method 0x186b2167.
//
// Solidity: function validators(string status, (bytes,uint64,uint64,bool,bool) pageRequest) view returns((string,string,bool,uint8,uint256,uint256,string,int64,int64,uint256,uint256)[] validators, (bytes,uint64) pageResponse)
func (_IStaking *IStakingCallerSession) Validators(status string, pageRequest PageRequest) (struct {
	Validators   []Validator
	PageResponse PageResponse
}, error) {
	return _IStaking.Contract.Validators(&_IStaking.CallOpts, status, pageRequest)
}

// CancelUnbondingDelegation is a paid mutator transaction binding the contract method 0x12d58dfe.
//
// Solidity: function cancelUnbondingDelegation(address delegatorAddress, string validatorAddress, uint256 amount, uint256 creationHeight) returns(bool success)
func (_IStaking *IStakingTransactor) CancelUnbondingDelegation(opts *bind.TransactOpts, delegatorAddress common.Address, validatorAddress string, amount *big.Int, creationHeight *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "cancelUnbondingDelegation", delegatorAddress, validatorAddress, amount, creationHeight)
}

// CancelUnbondingDelegation is a paid mutator transaction binding the contract method 0x12d58dfe.
//
// Solidity: function cancelUnbondingDelegation(address delegatorAddress, string validatorAddress, uint256 amount, uint256 creationHeight) returns(bool success)
func (_IStaking *IStakingSession) CancelUnbondingDelegation(delegatorAddress common.Address, validatorAddress string, amount *big.Int, creationHeight *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.CancelUnbondingDelegation(&_IStaking.TransactOpts, delegatorAddress, validatorAddress, amount, creationHeight)
}

// CancelUnbondingDelegation is a paid mutator transaction binding the contract method 0x12d58dfe.
//
// Solidity: function cancelUnbondingDelegation(address delegatorAddress, string validatorAddress, uint256 amount, uint256 creationHeight) returns(bool success)
func (_IStaking *IStakingTransactorSession) CancelUnbondingDelegation(delegatorAddress common.Address, validatorAddress string, amount *big.Int, creationHeight *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.CancelUnbondingDelegation(&_IStaking.TransactOpts, delegatorAddress, validatorAddress, amount, creationHeight)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xf7cd5516.
//
// Solidity: function createValidator((string,string,string,string,string) description, (uint256,uint256,uint256) commissionRates, uint256 minSelfDelegation, address validatorAddress, string pubkey, uint256 value) returns(bool success)
func (_IStaking *IStakingTransactor) CreateValidator(opts *bind.TransactOpts, description Description, commissionRates CommissionRates, minSelfDelegation *big.Int, validatorAddress common.Address, pubkey string, value *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "createValidator", description, commissionRates, minSelfDelegation, validatorAddress, pubkey, value)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xf7cd5516.
//
// Solidity: function createValidator((string,string,string,string,string) description, (uint256,uint256,uint256) commissionRates, uint256 minSelfDelegation, address validatorAddress, string pubkey, uint256 value) returns(bool success)
func (_IStaking *IStakingSession) CreateValidator(description Description, commissionRates CommissionRates, minSelfDelegation *big.Int, validatorAddress common.Address, pubkey string, value *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.CreateValidator(&_IStaking.TransactOpts, description, commissionRates, minSelfDelegation, validatorAddress, pubkey, value)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xf7cd5516.
//
// Solidity: function createValidator((string,string,string,string,string) description, (uint256,uint256,uint256) commissionRates, uint256 minSelfDelegation, address validatorAddress, string pubkey, uint256 value) returns(bool success)
func (_IStaking *IStakingTransactorSession) CreateValidator(description Description, commissionRates CommissionRates, minSelfDelegation *big.Int, validatorAddress common.Address, pubkey string, value *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.CreateValidator(&_IStaking.TransactOpts, description, commissionRates, minSelfDelegation, validatorAddress, pubkey, value)
}

// Delegate is a paid mutator transaction binding the contract method 0x53266bbb.
//
// Solidity: function delegate(address delegatorAddress, string validatorAddress, uint256 amount) returns(bool success)
func (_IStaking *IStakingTransactor) Delegate(opts *bind.TransactOpts, delegatorAddress common.Address, validatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "delegate", delegatorAddress, validatorAddress, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x53266bbb.
//
// Solidity: function delegate(address delegatorAddress, string validatorAddress, uint256 amount) returns(bool success)
func (_IStaking *IStakingSession) Delegate(delegatorAddress common.Address, validatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Delegate(&_IStaking.TransactOpts, delegatorAddress, validatorAddress, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x53266bbb.
//
// Solidity: function delegate(address delegatorAddress, string validatorAddress, uint256 amount) returns(bool success)
func (_IStaking *IStakingTransactorSession) Delegate(delegatorAddress common.Address, validatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Delegate(&_IStaking.TransactOpts, delegatorAddress, validatorAddress, amount)
}

// EditValidator is a paid mutator transaction binding the contract method 0xa50f05ac.
//
// Solidity: function editValidator((string,string,string,string,string) description, address validatorAddress, int256 commissionRate, int256 minSelfDelegation) returns(bool success)
func (_IStaking *IStakingTransactor) EditValidator(opts *bind.TransactOpts, description Description, validatorAddress common.Address, commissionRate *big.Int, minSelfDelegation *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "editValidator", description, validatorAddress, commissionRate, minSelfDelegation)
}

// EditValidator is a paid mutator transaction binding the contract method 0xa50f05ac.
//
// Solidity: function editValidator((string,string,string,string,string) description, address validatorAddress, int256 commissionRate, int256 minSelfDelegation) returns(bool success)
func (_IStaking *IStakingSession) EditValidator(description Description, validatorAddress common.Address, commissionRate *big.Int, minSelfDelegation *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.EditValidator(&_IStaking.TransactOpts, description, validatorAddress, commissionRate, minSelfDelegation)
}

// EditValidator is a paid mutator transaction binding the contract method 0xa50f05ac.
//
// Solidity: function editValidator((string,string,string,string,string) description, address validatorAddress, int256 commissionRate, int256 minSelfDelegation) returns(bool success)
func (_IStaking *IStakingTransactorSession) EditValidator(description Description, validatorAddress common.Address, commissionRate *big.Int, minSelfDelegation *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.EditValidator(&_IStaking.TransactOpts, description, validatorAddress, commissionRate, minSelfDelegation)
}

// Redelegate is a paid mutator transaction binding the contract method 0x54b826f5.
//
// Solidity: function redelegate(address delegatorAddress, string validatorSrcAddress, string validatorDstAddress, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingTransactor) Redelegate(opts *bind.TransactOpts, delegatorAddress common.Address, validatorSrcAddress string, validatorDstAddress string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "redelegate", delegatorAddress, validatorSrcAddress, validatorDstAddress, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x54b826f5.
//
// Solidity: function redelegate(address delegatorAddress, string validatorSrcAddress, string validatorDstAddress, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingSession) Redelegate(delegatorAddress common.Address, validatorSrcAddress string, validatorDstAddress string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Redelegate(&_IStaking.TransactOpts, delegatorAddress, validatorSrcAddress, validatorDstAddress, amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x54b826f5.
//
// Solidity: function redelegate(address delegatorAddress, string validatorSrcAddress, string validatorDstAddress, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingTransactorSession) Redelegate(delegatorAddress common.Address, validatorSrcAddress string, validatorDstAddress string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Redelegate(&_IStaking.TransactOpts, delegatorAddress, validatorSrcAddress, validatorDstAddress, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x3edab33c.
//
// Solidity: function undelegate(address delegatorAddress, string validatorAddress, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingTransactor) Undelegate(opts *bind.TransactOpts, delegatorAddress common.Address, validatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "undelegate", delegatorAddress, validatorAddress, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x3edab33c.
//
// Solidity: function undelegate(address delegatorAddress, string validatorAddress, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingSession) Undelegate(delegatorAddress common.Address, validatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Undelegate(&_IStaking.TransactOpts, delegatorAddress, validatorAddress, amount)
}

// Undelegate is a paid mutator transaction binding the contract method 0x3edab33c.
//
// Solidity: function undelegate(address delegatorAddress, string validatorAddress, uint256 amount) returns(int64 completionTime)
func (_IStaking *IStakingTransactorSession) Undelegate(delegatorAddress common.Address, validatorAddress string, amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.Undelegate(&_IStaking.TransactOpts, delegatorAddress, validatorAddress, amount)
}

// IStakingCancelUnbondingDelegationIterator is returned from FilterCancelUnbondingDelegation and is used to iterate over the raw logs and unpacked data for CancelUnbondingDelegation events raised by the IStaking contract.
type IStakingCancelUnbondingDelegationIterator struct {
	Event *IStakingCancelUnbondingDelegation // Event containing the contract specifics and raw log

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
func (it *IStakingCancelUnbondingDelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingCancelUnbondingDelegation)
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
		it.Event = new(IStakingCancelUnbondingDelegation)
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
func (it *IStakingCancelUnbondingDelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingCancelUnbondingDelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingCancelUnbondingDelegation represents a CancelUnbondingDelegation event raised by the IStaking contract.
type IStakingCancelUnbondingDelegation struct {
	DelegatorAddress common.Address
	ValidatorAddress common.Address
	Amount           *big.Int
	CreationHeight   *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterCancelUnbondingDelegation is a free log retrieval operation binding the contract event 0x6dbe2fb6b2613bdd8e3d284a6111592e06c3ab0af846ff89b6688d48f408dbb5.
//
// Solidity: event CancelUnbondingDelegation(address indexed delegatorAddress, address indexed validatorAddress, uint256 amount, uint256 creationHeight)
func (_IStaking *IStakingFilterer) FilterCancelUnbondingDelegation(opts *bind.FilterOpts, delegatorAddress []common.Address, validatorAddress []common.Address) (*IStakingCancelUnbondingDelegationIterator, error) {

	var delegatorAddressRule []interface{}
	for _, delegatorAddressItem := range delegatorAddress {
		delegatorAddressRule = append(delegatorAddressRule, delegatorAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "CancelUnbondingDelegation", delegatorAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &IStakingCancelUnbondingDelegationIterator{contract: _IStaking.contract, event: "CancelUnbondingDelegation", logs: logs, sub: sub}, nil
}

// WatchCancelUnbondingDelegation is a free log subscription operation binding the contract event 0x6dbe2fb6b2613bdd8e3d284a6111592e06c3ab0af846ff89b6688d48f408dbb5.
//
// Solidity: event CancelUnbondingDelegation(address indexed delegatorAddress, address indexed validatorAddress, uint256 amount, uint256 creationHeight)
func (_IStaking *IStakingFilterer) WatchCancelUnbondingDelegation(opts *bind.WatchOpts, sink chan<- *IStakingCancelUnbondingDelegation, delegatorAddress []common.Address, validatorAddress []common.Address) (event.Subscription, error) {

	var delegatorAddressRule []interface{}
	for _, delegatorAddressItem := range delegatorAddress {
		delegatorAddressRule = append(delegatorAddressRule, delegatorAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "CancelUnbondingDelegation", delegatorAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingCancelUnbondingDelegation)
				if err := _IStaking.contract.UnpackLog(event, "CancelUnbondingDelegation", log); err != nil {
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

// ParseCancelUnbondingDelegation is a log parse operation binding the contract event 0x6dbe2fb6b2613bdd8e3d284a6111592e06c3ab0af846ff89b6688d48f408dbb5.
//
// Solidity: event CancelUnbondingDelegation(address indexed delegatorAddress, address indexed validatorAddress, uint256 amount, uint256 creationHeight)
func (_IStaking *IStakingFilterer) ParseCancelUnbondingDelegation(log types.Log) (*IStakingCancelUnbondingDelegation, error) {
	event := new(IStakingCancelUnbondingDelegation)
	if err := _IStaking.contract.UnpackLog(event, "CancelUnbondingDelegation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingCreateValidatorIterator is returned from FilterCreateValidator and is used to iterate over the raw logs and unpacked data for CreateValidator events raised by the IStaking contract.
type IStakingCreateValidatorIterator struct {
	Event *IStakingCreateValidator // Event containing the contract specifics and raw log

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
func (it *IStakingCreateValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingCreateValidator)
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
		it.Event = new(IStakingCreateValidator)
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
func (it *IStakingCreateValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingCreateValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingCreateValidator represents a CreateValidator event raised by the IStaking contract.
type IStakingCreateValidator struct {
	ValidatorAddress common.Address
	Value            *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterCreateValidator is a free log retrieval operation binding the contract event 0x9bdb560f8135cb46033a55410c14e14b1a7bc2d3f3e9973f4b49533e176468b0.
//
// Solidity: event CreateValidator(address indexed validatorAddress, uint256 value)
func (_IStaking *IStakingFilterer) FilterCreateValidator(opts *bind.FilterOpts, validatorAddress []common.Address) (*IStakingCreateValidatorIterator, error) {

	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "CreateValidator", validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &IStakingCreateValidatorIterator{contract: _IStaking.contract, event: "CreateValidator", logs: logs, sub: sub}, nil
}

// WatchCreateValidator is a free log subscription operation binding the contract event 0x9bdb560f8135cb46033a55410c14e14b1a7bc2d3f3e9973f4b49533e176468b0.
//
// Solidity: event CreateValidator(address indexed validatorAddress, uint256 value)
func (_IStaking *IStakingFilterer) WatchCreateValidator(opts *bind.WatchOpts, sink chan<- *IStakingCreateValidator, validatorAddress []common.Address) (event.Subscription, error) {

	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "CreateValidator", validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingCreateValidator)
				if err := _IStaking.contract.UnpackLog(event, "CreateValidator", log); err != nil {
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

// ParseCreateValidator is a log parse operation binding the contract event 0x9bdb560f8135cb46033a55410c14e14b1a7bc2d3f3e9973f4b49533e176468b0.
//
// Solidity: event CreateValidator(address indexed validatorAddress, uint256 value)
func (_IStaking *IStakingFilterer) ParseCreateValidator(log types.Log) (*IStakingCreateValidator, error) {
	event := new(IStakingCreateValidator)
	if err := _IStaking.contract.UnpackLog(event, "CreateValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingDelegateIterator is returned from FilterDelegate and is used to iterate over the raw logs and unpacked data for Delegate events raised by the IStaking contract.
type IStakingDelegateIterator struct {
	Event *IStakingDelegate // Event containing the contract specifics and raw log

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
func (it *IStakingDelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingDelegate)
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
		it.Event = new(IStakingDelegate)
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
func (it *IStakingDelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingDelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingDelegate represents a Delegate event raised by the IStaking contract.
type IStakingDelegate struct {
	DelegatorAddress common.Address
	ValidatorAddress common.Address
	Amount           *big.Int
	NewShares        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDelegate is a free log retrieval operation binding the contract event 0x500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b.
//
// Solidity: event Delegate(address indexed delegatorAddress, address indexed validatorAddress, uint256 amount, uint256 newShares)
func (_IStaking *IStakingFilterer) FilterDelegate(opts *bind.FilterOpts, delegatorAddress []common.Address, validatorAddress []common.Address) (*IStakingDelegateIterator, error) {

	var delegatorAddressRule []interface{}
	for _, delegatorAddressItem := range delegatorAddress {
		delegatorAddressRule = append(delegatorAddressRule, delegatorAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "Delegate", delegatorAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &IStakingDelegateIterator{contract: _IStaking.contract, event: "Delegate", logs: logs, sub: sub}, nil
}

// WatchDelegate is a free log subscription operation binding the contract event 0x500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b.
//
// Solidity: event Delegate(address indexed delegatorAddress, address indexed validatorAddress, uint256 amount, uint256 newShares)
func (_IStaking *IStakingFilterer) WatchDelegate(opts *bind.WatchOpts, sink chan<- *IStakingDelegate, delegatorAddress []common.Address, validatorAddress []common.Address) (event.Subscription, error) {

	var delegatorAddressRule []interface{}
	for _, delegatorAddressItem := range delegatorAddress {
		delegatorAddressRule = append(delegatorAddressRule, delegatorAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "Delegate", delegatorAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingDelegate)
				if err := _IStaking.contract.UnpackLog(event, "Delegate", log); err != nil {
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

// ParseDelegate is a log parse operation binding the contract event 0x500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b.
//
// Solidity: event Delegate(address indexed delegatorAddress, address indexed validatorAddress, uint256 amount, uint256 newShares)
func (_IStaking *IStakingFilterer) ParseDelegate(log types.Log) (*IStakingDelegate, error) {
	event := new(IStakingDelegate)
	if err := _IStaking.contract.UnpackLog(event, "Delegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingEditValidatorIterator is returned from FilterEditValidator and is used to iterate over the raw logs and unpacked data for EditValidator events raised by the IStaking contract.
type IStakingEditValidatorIterator struct {
	Event *IStakingEditValidator // Event containing the contract specifics and raw log

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
func (it *IStakingEditValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingEditValidator)
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
		it.Event = new(IStakingEditValidator)
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
func (it *IStakingEditValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingEditValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingEditValidator represents a EditValidator event raised by the IStaking contract.
type IStakingEditValidator struct {
	ValidatorAddress  common.Address
	CommissionRate    *big.Int
	MinSelfDelegation *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterEditValidator is a free log retrieval operation binding the contract event 0xdce27cf2792bd8d8f28df5d2cdf379cd593414f21332370ca808c1e703eb4e1f.
//
// Solidity: event EditValidator(address indexed validatorAddress, int256 commissionRate, int256 minSelfDelegation)
func (_IStaking *IStakingFilterer) FilterEditValidator(opts *bind.FilterOpts, validatorAddress []common.Address) (*IStakingEditValidatorIterator, error) {

	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "EditValidator", validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &IStakingEditValidatorIterator{contract: _IStaking.contract, event: "EditValidator", logs: logs, sub: sub}, nil
}

// WatchEditValidator is a free log subscription operation binding the contract event 0xdce27cf2792bd8d8f28df5d2cdf379cd593414f21332370ca808c1e703eb4e1f.
//
// Solidity: event EditValidator(address indexed validatorAddress, int256 commissionRate, int256 minSelfDelegation)
func (_IStaking *IStakingFilterer) WatchEditValidator(opts *bind.WatchOpts, sink chan<- *IStakingEditValidator, validatorAddress []common.Address) (event.Subscription, error) {

	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "EditValidator", validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingEditValidator)
				if err := _IStaking.contract.UnpackLog(event, "EditValidator", log); err != nil {
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

// ParseEditValidator is a log parse operation binding the contract event 0xdce27cf2792bd8d8f28df5d2cdf379cd593414f21332370ca808c1e703eb4e1f.
//
// Solidity: event EditValidator(address indexed validatorAddress, int256 commissionRate, int256 minSelfDelegation)
func (_IStaking *IStakingFilterer) ParseEditValidator(log types.Log) (*IStakingEditValidator, error) {
	event := new(IStakingEditValidator)
	if err := _IStaking.contract.UnpackLog(event, "EditValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingRedelegateIterator is returned from FilterRedelegate and is used to iterate over the raw logs and unpacked data for Redelegate events raised by the IStaking contract.
type IStakingRedelegateIterator struct {
	Event *IStakingRedelegate // Event containing the contract specifics and raw log

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
func (it *IStakingRedelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingRedelegate)
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
		it.Event = new(IStakingRedelegate)
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
func (it *IStakingRedelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingRedelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingRedelegate represents a Redelegate event raised by the IStaking contract.
type IStakingRedelegate struct {
	DelegatorAddress    common.Address
	ValidatorSrcAddress common.Address
	ValidatorDstAddress common.Address
	Amount              *big.Int
	CompletionTime      *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterRedelegate is a free log retrieval operation binding the contract event 0x82b07f2421474f1e3f1e0b34738cb5ffb925273f408e7591d9c803dcae8da657.
//
// Solidity: event Redelegate(address indexed delegatorAddress, address indexed validatorSrcAddress, address indexed validatorDstAddress, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) FilterRedelegate(opts *bind.FilterOpts, delegatorAddress []common.Address, validatorSrcAddress []common.Address, validatorDstAddress []common.Address) (*IStakingRedelegateIterator, error) {

	var delegatorAddressRule []interface{}
	for _, delegatorAddressItem := range delegatorAddress {
		delegatorAddressRule = append(delegatorAddressRule, delegatorAddressItem)
	}
	var validatorSrcAddressRule []interface{}
	for _, validatorSrcAddressItem := range validatorSrcAddress {
		validatorSrcAddressRule = append(validatorSrcAddressRule, validatorSrcAddressItem)
	}
	var validatorDstAddressRule []interface{}
	for _, validatorDstAddressItem := range validatorDstAddress {
		validatorDstAddressRule = append(validatorDstAddressRule, validatorDstAddressItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "Redelegate", delegatorAddressRule, validatorSrcAddressRule, validatorDstAddressRule)
	if err != nil {
		return nil, err
	}
	return &IStakingRedelegateIterator{contract: _IStaking.contract, event: "Redelegate", logs: logs, sub: sub}, nil
}

// WatchRedelegate is a free log subscription operation binding the contract event 0x82b07f2421474f1e3f1e0b34738cb5ffb925273f408e7591d9c803dcae8da657.
//
// Solidity: event Redelegate(address indexed delegatorAddress, address indexed validatorSrcAddress, address indexed validatorDstAddress, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) WatchRedelegate(opts *bind.WatchOpts, sink chan<- *IStakingRedelegate, delegatorAddress []common.Address, validatorSrcAddress []common.Address, validatorDstAddress []common.Address) (event.Subscription, error) {

	var delegatorAddressRule []interface{}
	for _, delegatorAddressItem := range delegatorAddress {
		delegatorAddressRule = append(delegatorAddressRule, delegatorAddressItem)
	}
	var validatorSrcAddressRule []interface{}
	for _, validatorSrcAddressItem := range validatorSrcAddress {
		validatorSrcAddressRule = append(validatorSrcAddressRule, validatorSrcAddressItem)
	}
	var validatorDstAddressRule []interface{}
	for _, validatorDstAddressItem := range validatorDstAddress {
		validatorDstAddressRule = append(validatorDstAddressRule, validatorDstAddressItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "Redelegate", delegatorAddressRule, validatorSrcAddressRule, validatorDstAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingRedelegate)
				if err := _IStaking.contract.UnpackLog(event, "Redelegate", log); err != nil {
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

// ParseRedelegate is a log parse operation binding the contract event 0x82b07f2421474f1e3f1e0b34738cb5ffb925273f408e7591d9c803dcae8da657.
//
// Solidity: event Redelegate(address indexed delegatorAddress, address indexed validatorSrcAddress, address indexed validatorDstAddress, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) ParseRedelegate(log types.Log) (*IStakingRedelegate, error) {
	event := new(IStakingRedelegate)
	if err := _IStaking.contract.UnpackLog(event, "Redelegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingUnbondIterator is returned from FilterUnbond and is used to iterate over the raw logs and unpacked data for Unbond events raised by the IStaking contract.
type IStakingUnbondIterator struct {
	Event *IStakingUnbond // Event containing the contract specifics and raw log

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
func (it *IStakingUnbondIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingUnbond)
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
		it.Event = new(IStakingUnbond)
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
func (it *IStakingUnbondIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingUnbondIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingUnbond represents a Unbond event raised by the IStaking contract.
type IStakingUnbond struct {
	DelegatorAddress common.Address
	ValidatorAddress common.Address
	Amount           *big.Int
	CompletionTime   *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUnbond is a free log retrieval operation binding the contract event 0x4bf8087be3b8a59c2662514df2ed4a3dcaf9ca22f442340cfc05a4e52343d18e.
//
// Solidity: event Unbond(address indexed delegatorAddress, address indexed validatorAddress, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) FilterUnbond(opts *bind.FilterOpts, delegatorAddress []common.Address, validatorAddress []common.Address) (*IStakingUnbondIterator, error) {

	var delegatorAddressRule []interface{}
	for _, delegatorAddressItem := range delegatorAddress {
		delegatorAddressRule = append(delegatorAddressRule, delegatorAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "Unbond", delegatorAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &IStakingUnbondIterator{contract: _IStaking.contract, event: "Unbond", logs: logs, sub: sub}, nil
}

// WatchUnbond is a free log subscription operation binding the contract event 0x4bf8087be3b8a59c2662514df2ed4a3dcaf9ca22f442340cfc05a4e52343d18e.
//
// Solidity: event Unbond(address indexed delegatorAddress, address indexed validatorAddress, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) WatchUnbond(opts *bind.WatchOpts, sink chan<- *IStakingUnbond, delegatorAddress []common.Address, validatorAddress []common.Address) (event.Subscription, error) {

	var delegatorAddressRule []interface{}
	for _, delegatorAddressItem := range delegatorAddress {
		delegatorAddressRule = append(delegatorAddressRule, delegatorAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "Unbond", delegatorAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingUnbond)
				if err := _IStaking.contract.UnpackLog(event, "Unbond", log); err != nil {
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

// ParseUnbond is a log parse operation binding the contract event 0x4bf8087be3b8a59c2662514df2ed4a3dcaf9ca22f442340cfc05a4e52343d18e.
//
// Solidity: event Unbond(address indexed delegatorAddress, address indexed validatorAddress, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) ParseUnbond(log types.Log) (*IStakingUnbond, error) {
	event := new(IStakingUnbond)
	if err := _IStaking.contract.UnpackLog(event, "Unbond", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
