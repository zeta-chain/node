// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testprototype

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

// TestPrototypeMetaData contains all meta data concerning the TestPrototype contract.
var TestPrototypeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"name\":\"bech32ToHexAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"prefix\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"bech32ify\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int64\",\"name\":\"chainID\",\"type\":\"int64\"}],\"name\":\"getGasStabilityPoolBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"result\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405260656000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561005157600080fd5b50610878806100616000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80630615b74e146100465780633ee8d1a414610076578063e4e2a4ec146100a6575b600080fd5b610060600480360381019061005b9190610481565b6100d6565b60405161006d9190610565565b60405180910390f35b610090600480360381019061008b91906105c0565b610181565b60405161009d9190610606565b60405180910390f35b6100c060048036038101906100bb9190610621565b610225565b6040516100cd9190610679565b60405180910390f35b606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630615b74e84846040518363ffffffff1660e01b8152600401610133929190610694565b600060405180830381865afa158015610150573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052508101906101799190610734565b905092915050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633ee8d1a4836040518263ffffffff1660e01b81526004016101dd919061078c565b602060405180830381865afa1580156101fa573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061021e91906107d3565b9050919050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e4e2a4ec836040518263ffffffff1660e01b81526004016102819190610565565b602060405180830381865afa15801561029e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102c29190610815565b9050919050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610330826102e7565b810181811067ffffffffffffffff8211171561034f5761034e6102f8565b5b80604052505050565b60006103626102c9565b905061036e8282610327565b919050565b600067ffffffffffffffff82111561038e5761038d6102f8565b5b610397826102e7565b9050602081019050919050565b82818337600083830152505050565b60006103c66103c184610373565b610358565b9050828152602081018484840111156103e2576103e16102e2565b5b6103ed8482856103a4565b509392505050565b600082601f83011261040a576104096102dd565b5b813561041a8482602086016103b3565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061044e82610423565b9050919050565b61045e81610443565b811461046957600080fd5b50565b60008135905061047b81610455565b92915050565b60008060408385031215610498576104976102d3565b5b600083013567ffffffffffffffff8111156104b6576104b56102d8565b5b6104c2858286016103f5565b92505060206104d38582860161046c565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b60005b838110156105175780820151818401526020810190506104fc565b83811115610526576000848401525b50505050565b6000610537826104dd565b61054181856104e8565b93506105518185602086016104f9565b61055a816102e7565b840191505092915050565b6000602082019050818103600083015261057f818461052c565b905092915050565b60008160070b9050919050565b61059d81610587565b81146105a857600080fd5b50565b6000813590506105ba81610594565b92915050565b6000602082840312156105d6576105d56102d3565b5b60006105e4848285016105ab565b91505092915050565b6000819050919050565b610600816105ed565b82525050565b600060208201905061061b60008301846105f7565b92915050565b600060208284031215610637576106366102d3565b5b600082013567ffffffffffffffff811115610655576106546102d8565b5b610661848285016103f5565b91505092915050565b61067381610443565b82525050565b600060208201905061068e600083018461066a565b92915050565b600060408201905081810360008301526106ae818561052c565b90506106bd602083018461066a565b9392505050565b60006106d76106d284610373565b610358565b9050828152602081018484840111156106f3576106f26102e2565b5b6106fe8482856104f9565b509392505050565b600082601f83011261071b5761071a6102dd565b5b815161072b8482602086016106c4565b91505092915050565b60006020828403121561074a576107496102d3565b5b600082015167ffffffffffffffff811115610768576107676102d8565b5b61077484828501610706565b91505092915050565b61078681610587565b82525050565b60006020820190506107a1600083018461077d565b92915050565b6107b0816105ed565b81146107bb57600080fd5b50565b6000815190506107cd816107a7565b92915050565b6000602082840312156107e9576107e86102d3565b5b60006107f7848285016107be565b91505092915050565b60008151905061080f81610455565b92915050565b60006020828403121561082b5761082a6102d3565b5b600061083984828501610800565b9150509291505056fea26469706673582212203599aa5bcc18cda9492a410a6619dfd803fa572b0624c47efa4d4159a8f2430f64736f6c634300080a0033",
}

// TestPrototypeABI is the input ABI used to generate the binding from.
// Deprecated: Use TestPrototypeMetaData.ABI instead.
var TestPrototypeABI = TestPrototypeMetaData.ABI

// TestPrototypeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestPrototypeMetaData.Bin instead.
var TestPrototypeBin = TestPrototypeMetaData.Bin

// DeployTestPrototype deploys a new Ethereum contract, binding an instance of TestPrototype to it.
func DeployTestPrototype(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestPrototype, error) {
	parsed, err := TestPrototypeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestPrototypeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestPrototype{TestPrototypeCaller: TestPrototypeCaller{contract: contract}, TestPrototypeTransactor: TestPrototypeTransactor{contract: contract}, TestPrototypeFilterer: TestPrototypeFilterer{contract: contract}}, nil
}

// TestPrototype is an auto generated Go binding around an Ethereum contract.
type TestPrototype struct {
	TestPrototypeCaller     // Read-only binding to the contract
	TestPrototypeTransactor // Write-only binding to the contract
	TestPrototypeFilterer   // Log filterer for contract events
}

// TestPrototypeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestPrototypeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestPrototypeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestPrototypeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestPrototypeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestPrototypeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestPrototypeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestPrototypeSession struct {
	Contract     *TestPrototype    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestPrototypeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestPrototypeCallerSession struct {
	Contract *TestPrototypeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TestPrototypeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestPrototypeTransactorSession struct {
	Contract     *TestPrototypeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TestPrototypeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestPrototypeRaw struct {
	Contract *TestPrototype // Generic contract binding to access the raw methods on
}

// TestPrototypeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestPrototypeCallerRaw struct {
	Contract *TestPrototypeCaller // Generic read-only contract binding to access the raw methods on
}

// TestPrototypeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestPrototypeTransactorRaw struct {
	Contract *TestPrototypeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestPrototype creates a new instance of TestPrototype, bound to a specific deployed contract.
func NewTestPrototype(address common.Address, backend bind.ContractBackend) (*TestPrototype, error) {
	contract, err := bindTestPrototype(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestPrototype{TestPrototypeCaller: TestPrototypeCaller{contract: contract}, TestPrototypeTransactor: TestPrototypeTransactor{contract: contract}, TestPrototypeFilterer: TestPrototypeFilterer{contract: contract}}, nil
}

// NewTestPrototypeCaller creates a new read-only instance of TestPrototype, bound to a specific deployed contract.
func NewTestPrototypeCaller(address common.Address, caller bind.ContractCaller) (*TestPrototypeCaller, error) {
	contract, err := bindTestPrototype(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestPrototypeCaller{contract: contract}, nil
}

// NewTestPrototypeTransactor creates a new write-only instance of TestPrototype, bound to a specific deployed contract.
func NewTestPrototypeTransactor(address common.Address, transactor bind.ContractTransactor) (*TestPrototypeTransactor, error) {
	contract, err := bindTestPrototype(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestPrototypeTransactor{contract: contract}, nil
}

// NewTestPrototypeFilterer creates a new log filterer instance of TestPrototype, bound to a specific deployed contract.
func NewTestPrototypeFilterer(address common.Address, filterer bind.ContractFilterer) (*TestPrototypeFilterer, error) {
	contract, err := bindTestPrototype(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestPrototypeFilterer{contract: contract}, nil
}

// bindTestPrototype binds a generic wrapper to an already deployed contract.
func bindTestPrototype(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestPrototypeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestPrototype *TestPrototypeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestPrototype.Contract.TestPrototypeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestPrototype *TestPrototypeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestPrototype.Contract.TestPrototypeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestPrototype *TestPrototypeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestPrototype.Contract.TestPrototypeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestPrototype *TestPrototypeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestPrototype.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestPrototype *TestPrototypeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestPrototype.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestPrototype *TestPrototypeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestPrototype.Contract.contract.Transact(opts, method, params...)
}

// Bech32ToHexAddr is a free data retrieval call binding the contract method 0xe4e2a4ec.
//
// Solidity: function bech32ToHexAddr(string bech32) view returns(address addr)
func (_TestPrototype *TestPrototypeCaller) Bech32ToHexAddr(opts *bind.CallOpts, bech32 string) (common.Address, error) {
	var out []interface{}
	err := _TestPrototype.contract.Call(opts, &out, "bech32ToHexAddr", bech32)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bech32ToHexAddr is a free data retrieval call binding the contract method 0xe4e2a4ec.
//
// Solidity: function bech32ToHexAddr(string bech32) view returns(address addr)
func (_TestPrototype *TestPrototypeSession) Bech32ToHexAddr(bech32 string) (common.Address, error) {
	return _TestPrototype.Contract.Bech32ToHexAddr(&_TestPrototype.CallOpts, bech32)
}

// Bech32ToHexAddr is a free data retrieval call binding the contract method 0xe4e2a4ec.
//
// Solidity: function bech32ToHexAddr(string bech32) view returns(address addr)
func (_TestPrototype *TestPrototypeCallerSession) Bech32ToHexAddr(bech32 string) (common.Address, error) {
	return _TestPrototype.Contract.Bech32ToHexAddr(&_TestPrototype.CallOpts, bech32)
}

// Bech32ify is a free data retrieval call binding the contract method 0x0615b74e.
//
// Solidity: function bech32ify(string prefix, address addr) view returns(string bech32)
func (_TestPrototype *TestPrototypeCaller) Bech32ify(opts *bind.CallOpts, prefix string, addr common.Address) (string, error) {
	var out []interface{}
	err := _TestPrototype.contract.Call(opts, &out, "bech32ify", prefix, addr)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Bech32ify is a free data retrieval call binding the contract method 0x0615b74e.
//
// Solidity: function bech32ify(string prefix, address addr) view returns(string bech32)
func (_TestPrototype *TestPrototypeSession) Bech32ify(prefix string, addr common.Address) (string, error) {
	return _TestPrototype.Contract.Bech32ify(&_TestPrototype.CallOpts, prefix, addr)
}

// Bech32ify is a free data retrieval call binding the contract method 0x0615b74e.
//
// Solidity: function bech32ify(string prefix, address addr) view returns(string bech32)
func (_TestPrototype *TestPrototypeCallerSession) Bech32ify(prefix string, addr common.Address) (string, error) {
	return _TestPrototype.Contract.Bech32ify(&_TestPrototype.CallOpts, prefix, addr)
}

// GetGasStabilityPoolBalance is a free data retrieval call binding the contract method 0x3ee8d1a4.
//
// Solidity: function getGasStabilityPoolBalance(int64 chainID) view returns(uint256 result)
func (_TestPrototype *TestPrototypeCaller) GetGasStabilityPoolBalance(opts *bind.CallOpts, chainID int64) (*big.Int, error) {
	var out []interface{}
	err := _TestPrototype.contract.Call(opts, &out, "getGasStabilityPoolBalance", chainID)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGasStabilityPoolBalance is a free data retrieval call binding the contract method 0x3ee8d1a4.
//
// Solidity: function getGasStabilityPoolBalance(int64 chainID) view returns(uint256 result)
func (_TestPrototype *TestPrototypeSession) GetGasStabilityPoolBalance(chainID int64) (*big.Int, error) {
	return _TestPrototype.Contract.GetGasStabilityPoolBalance(&_TestPrototype.CallOpts, chainID)
}

// GetGasStabilityPoolBalance is a free data retrieval call binding the contract method 0x3ee8d1a4.
//
// Solidity: function getGasStabilityPoolBalance(int64 chainID) view returns(uint256 result)
func (_TestPrototype *TestPrototypeCallerSession) GetGasStabilityPoolBalance(chainID int64) (*big.Int, error) {
	return _TestPrototype.Contract.GetGasStabilityPoolBalance(&_TestPrototype.CallOpts, chainID)
}
