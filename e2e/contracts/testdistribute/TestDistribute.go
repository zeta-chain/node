// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testdistribute

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

// TestDistributeMetaData contains all meta data concerning the TestDistribute contract.
var TestDistributeMetaData = &bind.MetaData{
	ABI: "[{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"claimRewardsThroughContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"distributeThroughContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"}],\"name\":\"getDelegatorValidatorsThroughContract\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"getRewardsThroughContract\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structDecCoin[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405260666000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561005157600080fd5b50610e2d806100616000396000f3fe6080604052600436106100435760003560e01c80630f4865ea1461004c57806350b54e8414610089578063834b902f146100c6578063cdc5ec4a146101035761004a565b3661004a57005b005b34801561005857600080fd5b50610073600480360381019061006e919061059d565b610140565b6040516100809190610614565b60405180910390f35b34801561009557600080fd5b506100b060048036038101906100ab9190610665565b6101e9565b6040516100bd9190610614565b60405180910390f35b3480156100d257600080fd5b506100ed60048036038101906100e8919061059d565b610292565b6040516100fa919061083b565b60405180910390f35b34801561010f57600080fd5b5061012a6004803603810190610125919061085d565b61033d565b604051610137919061094c565b60405180910390f35b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166354dbdc3884846040518363ffffffff1660e01b815260040161019e9291906109c7565b6020604051808303816000875af11580156101bd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101e19190610a23565b905092915050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fb93210884846040518363ffffffff1660e01b8152600401610247929190610a5f565b6020604051808303816000875af1158015610266573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061028a9190610a23565b905092915050565b606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16639342879284846040518363ffffffff1660e01b81526004016102ef9291906109c7565b600060405180830381865afa15801561030c573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052508101906103359190610c69565b905092915050565b606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b6a216ae836040518263ffffffff1660e01b81526004016103989190610cb2565b600060405180830381865afa1580156103b5573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052508101906103de9190610dae565b9050919050565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610424826103f9565b9050919050565b61043481610419565b811461043f57600080fd5b50565b6000813590506104518161042b565b92915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6104aa82610461565b810181811067ffffffffffffffff821117156104c9576104c8610472565b5b80604052505050565b60006104dc6103e5565b90506104e882826104a1565b919050565b600067ffffffffffffffff82111561050857610507610472565b5b61051182610461565b9050602081019050919050565b82818337600083830152505050565b600061054061053b846104ed565b6104d2565b90508281526020810184848401111561055c5761055b61045c565b5b61056784828561051e565b509392505050565b600082601f83011261058457610583610457565b5b813561059484826020860161052d565b91505092915050565b600080604083850312156105b4576105b36103ef565b5b60006105c285828601610442565b925050602083013567ffffffffffffffff8111156105e3576105e26103f4565b5b6105ef8582860161056f565b9150509250929050565b60008115159050919050565b61060e816105f9565b82525050565b60006020820190506106296000830184610605565b92915050565b6000819050919050565b6106428161062f565b811461064d57600080fd5b50565b60008135905061065f81610639565b92915050565b6000806040838503121561067c5761067b6103ef565b5b600061068a85828601610442565b925050602061069b85828601610650565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561070b5780820151818401526020810190506106f0565b8381111561071a576000848401525b50505050565b600061072b826106d1565b61073581856106dc565b93506107458185602086016106ed565b61074e81610461565b840191505092915050565b6107628161062f565b82525050565b600060408301600083015184820360008601526107858282610720565b915050602083015161079a6020860182610759565b508091505092915050565b60006107b18383610768565b905092915050565b6000602082019050919050565b60006107d1826106a5565b6107db81856106b0565b9350836020820285016107ed856106c1565b8060005b85811015610829578484038952815161080a85826107a5565b9450610815836107b9565b925060208a019950506001810190506107f1565b50829750879550505050505092915050565b6000602082019050818103600083015261085581846107c6565b905092915050565b600060208284031215610873576108726103ef565b5b600061088184828501610442565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b60006108c28383610720565b905092915050565b6000602082019050919050565b60006108e28261088a565b6108ec8185610895565b9350836020820285016108fe856108a6565b8060005b8581101561093a578484038952815161091b85826108b6565b9450610926836108ca565b925060208a01995050600181019050610902565b50829750879550505050505092915050565b6000602082019050818103600083015261096681846108d7565b905092915050565b61097781610419565b82525050565b600082825260208201905092915050565b6000610999826106d1565b6109a3818561097d565b93506109b38185602086016106ed565b6109bc81610461565b840191505092915050565b60006040820190506109dc600083018561096e565b81810360208301526109ee818461098e565b90509392505050565b610a00816105f9565b8114610a0b57600080fd5b50565b600081519050610a1d816109f7565b92915050565b600060208284031215610a3957610a386103ef565b5b6000610a4784828501610a0e565b91505092915050565b610a598161062f565b82525050565b6000604082019050610a74600083018561096e565b610a816020830184610a50565b9392505050565b600067ffffffffffffffff821115610aa357610aa2610472565b5b602082029050602081019050919050565b600080fd5b600080fd5b600080fd5b6000610ad6610ad1846104ed565b6104d2565b905082815260208101848484011115610af257610af161045c565b5b610afd8482856106ed565b509392505050565b600082601f830112610b1a57610b19610457565b5b8151610b2a848260208601610ac3565b91505092915050565b600081519050610b4281610639565b92915050565b600060408284031215610b5e57610b5d610ab9565b5b610b6860406104d2565b9050600082015167ffffffffffffffff811115610b8857610b87610abe565b5b610b9484828501610b05565b6000830152506020610ba884828501610b33565b60208301525092915050565b6000610bc7610bc284610a88565b6104d2565b90508083825260208201905060208402830185811115610bea57610be9610ab4565b5b835b81811015610c3157805167ffffffffffffffff811115610c0f57610c0e610457565b5b808601610c1c8982610b48565b85526020850194505050602081019050610bec565b5050509392505050565b600082601f830112610c5057610c4f610457565b5b8151610c60848260208601610bb4565b91505092915050565b600060208284031215610c7f57610c7e6103ef565b5b600082015167ffffffffffffffff811115610c9d57610c9c6103f4565b5b610ca984828501610c3b565b91505092915050565b6000602082019050610cc7600083018461096e565b92915050565b600067ffffffffffffffff821115610ce857610ce7610472565b5b602082029050602081019050919050565b6000610d0c610d0784610ccd565b6104d2565b90508083825260208201905060208402830185811115610d2f57610d2e610ab4565b5b835b81811015610d7657805167ffffffffffffffff811115610d5457610d53610457565b5b808601610d618982610b05565b85526020850194505050602081019050610d31565b5050509392505050565b600082601f830112610d9557610d94610457565b5b8151610da5848260208601610cf9565b91505092915050565b600060208284031215610dc457610dc36103ef565b5b600082015167ffffffffffffffff811115610de257610de16103f4565b5b610dee84828501610d80565b9150509291505056fea2646970667358221220d29e8c0ffd7f95c3ae2950ad56c9ec844a4f83f78ebf290ed1f2076d3fa1537864736f6c634300080a0033",
}

// TestDistributeABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDistributeMetaData.ABI instead.
var TestDistributeABI = TestDistributeMetaData.ABI

// TestDistributeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestDistributeMetaData.Bin instead.
var TestDistributeBin = TestDistributeMetaData.Bin

// DeployTestDistribute deploys a new Ethereum contract, binding an instance of TestDistribute to it.
func DeployTestDistribute(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestDistribute, error) {
	parsed, err := TestDistributeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDistributeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestDistribute{TestDistributeCaller: TestDistributeCaller{contract: contract}, TestDistributeTransactor: TestDistributeTransactor{contract: contract}, TestDistributeFilterer: TestDistributeFilterer{contract: contract}}, nil
}

// TestDistribute is an auto generated Go binding around an Ethereum contract.
type TestDistribute struct {
	TestDistributeCaller     // Read-only binding to the contract
	TestDistributeTransactor // Write-only binding to the contract
	TestDistributeFilterer   // Log filterer for contract events
}

// TestDistributeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestDistributeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDistributeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestDistributeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDistributeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestDistributeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDistributeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestDistributeSession struct {
	Contract     *TestDistribute   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestDistributeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestDistributeCallerSession struct {
	Contract *TestDistributeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// TestDistributeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestDistributeTransactorSession struct {
	Contract     *TestDistributeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// TestDistributeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestDistributeRaw struct {
	Contract *TestDistribute // Generic contract binding to access the raw methods on
}

// TestDistributeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestDistributeCallerRaw struct {
	Contract *TestDistributeCaller // Generic read-only contract binding to access the raw methods on
}

// TestDistributeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestDistributeTransactorRaw struct {
	Contract *TestDistributeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestDistribute creates a new instance of TestDistribute, bound to a specific deployed contract.
func NewTestDistribute(address common.Address, backend bind.ContractBackend) (*TestDistribute, error) {
	contract, err := bindTestDistribute(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestDistribute{TestDistributeCaller: TestDistributeCaller{contract: contract}, TestDistributeTransactor: TestDistributeTransactor{contract: contract}, TestDistributeFilterer: TestDistributeFilterer{contract: contract}}, nil
}

// NewTestDistributeCaller creates a new read-only instance of TestDistribute, bound to a specific deployed contract.
func NewTestDistributeCaller(address common.Address, caller bind.ContractCaller) (*TestDistributeCaller, error) {
	contract, err := bindTestDistribute(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestDistributeCaller{contract: contract}, nil
}

// NewTestDistributeTransactor creates a new write-only instance of TestDistribute, bound to a specific deployed contract.
func NewTestDistributeTransactor(address common.Address, transactor bind.ContractTransactor) (*TestDistributeTransactor, error) {
	contract, err := bindTestDistribute(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestDistributeTransactor{contract: contract}, nil
}

// NewTestDistributeFilterer creates a new log filterer instance of TestDistribute, bound to a specific deployed contract.
func NewTestDistributeFilterer(address common.Address, filterer bind.ContractFilterer) (*TestDistributeFilterer, error) {
	contract, err := bindTestDistribute(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestDistributeFilterer{contract: contract}, nil
}

// bindTestDistribute binds a generic wrapper to an already deployed contract.
func bindTestDistribute(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestDistributeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDistribute *TestDistributeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDistribute.Contract.TestDistributeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDistribute *TestDistributeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDistribute.Contract.TestDistributeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDistribute *TestDistributeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDistribute.Contract.TestDistributeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDistribute *TestDistributeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDistribute.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDistribute *TestDistributeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDistribute.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDistribute *TestDistributeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDistribute.Contract.contract.Transact(opts, method, params...)
}

// GetDelegatorValidatorsThroughContract is a free data retrieval call binding the contract method 0xcdc5ec4a.
//
// Solidity: function getDelegatorValidatorsThroughContract(address delegator) view returns(string[])
func (_TestDistribute *TestDistributeCaller) GetDelegatorValidatorsThroughContract(opts *bind.CallOpts, delegator common.Address) ([]string, error) {
	var out []interface{}
	err := _TestDistribute.contract.Call(opts, &out, "getDelegatorValidatorsThroughContract", delegator)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetDelegatorValidatorsThroughContract is a free data retrieval call binding the contract method 0xcdc5ec4a.
//
// Solidity: function getDelegatorValidatorsThroughContract(address delegator) view returns(string[])
func (_TestDistribute *TestDistributeSession) GetDelegatorValidatorsThroughContract(delegator common.Address) ([]string, error) {
	return _TestDistribute.Contract.GetDelegatorValidatorsThroughContract(&_TestDistribute.CallOpts, delegator)
}

// GetDelegatorValidatorsThroughContract is a free data retrieval call binding the contract method 0xcdc5ec4a.
//
// Solidity: function getDelegatorValidatorsThroughContract(address delegator) view returns(string[])
func (_TestDistribute *TestDistributeCallerSession) GetDelegatorValidatorsThroughContract(delegator common.Address) ([]string, error) {
	return _TestDistribute.Contract.GetDelegatorValidatorsThroughContract(&_TestDistribute.CallOpts, delegator)
}

// GetRewardsThroughContract is a free data retrieval call binding the contract method 0x834b902f.
//
// Solidity: function getRewardsThroughContract(address delegator, string validator) view returns((string,uint256)[])
func (_TestDistribute *TestDistributeCaller) GetRewardsThroughContract(opts *bind.CallOpts, delegator common.Address, validator string) ([]DecCoin, error) {
	var out []interface{}
	err := _TestDistribute.contract.Call(opts, &out, "getRewardsThroughContract", delegator, validator)

	if err != nil {
		return *new([]DecCoin), err
	}

	out0 := *abi.ConvertType(out[0], new([]DecCoin)).(*[]DecCoin)

	return out0, err

}

// GetRewardsThroughContract is a free data retrieval call binding the contract method 0x834b902f.
//
// Solidity: function getRewardsThroughContract(address delegator, string validator) view returns((string,uint256)[])
func (_TestDistribute *TestDistributeSession) GetRewardsThroughContract(delegator common.Address, validator string) ([]DecCoin, error) {
	return _TestDistribute.Contract.GetRewardsThroughContract(&_TestDistribute.CallOpts, delegator, validator)
}

// GetRewardsThroughContract is a free data retrieval call binding the contract method 0x834b902f.
//
// Solidity: function getRewardsThroughContract(address delegator, string validator) view returns((string,uint256)[])
func (_TestDistribute *TestDistributeCallerSession) GetRewardsThroughContract(delegator common.Address, validator string) ([]DecCoin, error) {
	return _TestDistribute.Contract.GetRewardsThroughContract(&_TestDistribute.CallOpts, delegator, validator)
}

// ClaimRewardsThroughContract is a paid mutator transaction binding the contract method 0x0f4865ea.
//
// Solidity: function claimRewardsThroughContract(address delegator, string validator) returns(bool)
func (_TestDistribute *TestDistributeTransactor) ClaimRewardsThroughContract(opts *bind.TransactOpts, delegator common.Address, validator string) (*types.Transaction, error) {
	return _TestDistribute.contract.Transact(opts, "claimRewardsThroughContract", delegator, validator)
}

// ClaimRewardsThroughContract is a paid mutator transaction binding the contract method 0x0f4865ea.
//
// Solidity: function claimRewardsThroughContract(address delegator, string validator) returns(bool)
func (_TestDistribute *TestDistributeSession) ClaimRewardsThroughContract(delegator common.Address, validator string) (*types.Transaction, error) {
	return _TestDistribute.Contract.ClaimRewardsThroughContract(&_TestDistribute.TransactOpts, delegator, validator)
}

// ClaimRewardsThroughContract is a paid mutator transaction binding the contract method 0x0f4865ea.
//
// Solidity: function claimRewardsThroughContract(address delegator, string validator) returns(bool)
func (_TestDistribute *TestDistributeTransactorSession) ClaimRewardsThroughContract(delegator common.Address, validator string) (*types.Transaction, error) {
	return _TestDistribute.Contract.ClaimRewardsThroughContract(&_TestDistribute.TransactOpts, delegator, validator)
}

// DistributeThroughContract is a paid mutator transaction binding the contract method 0x50b54e84.
//
// Solidity: function distributeThroughContract(address zrc20, uint256 amount) returns(bool)
func (_TestDistribute *TestDistributeTransactor) DistributeThroughContract(opts *bind.TransactOpts, zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestDistribute.contract.Transact(opts, "distributeThroughContract", zrc20, amount)
}

// DistributeThroughContract is a paid mutator transaction binding the contract method 0x50b54e84.
//
// Solidity: function distributeThroughContract(address zrc20, uint256 amount) returns(bool)
func (_TestDistribute *TestDistributeSession) DistributeThroughContract(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestDistribute.Contract.DistributeThroughContract(&_TestDistribute.TransactOpts, zrc20, amount)
}

// DistributeThroughContract is a paid mutator transaction binding the contract method 0x50b54e84.
//
// Solidity: function distributeThroughContract(address zrc20, uint256 amount) returns(bool)
func (_TestDistribute *TestDistributeTransactorSession) DistributeThroughContract(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestDistribute.Contract.DistributeThroughContract(&_TestDistribute.TransactOpts, zrc20, amount)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestDistribute *TestDistributeTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TestDistribute.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestDistribute *TestDistributeSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TestDistribute.Contract.Fallback(&_TestDistribute.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestDistribute *TestDistributeTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TestDistribute.Contract.Fallback(&_TestDistribute.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDistribute *TestDistributeTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDistribute.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDistribute *TestDistributeSession) Receive() (*types.Transaction, error) {
	return _TestDistribute.Contract.Receive(&_TestDistribute.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDistribute *TestDistributeTransactorSession) Receive() (*types.Transaction, error) {
	return _TestDistribute.Contract.Receive(&_TestDistribute.TransactOpts)
}
