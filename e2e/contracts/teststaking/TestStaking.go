// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package teststaking

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

// TestStakingMetaData contains all meta data concerning the TestStaking contract.
var TestStakingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"name\":\"bech32CallFn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"name\":\"bech32Fn\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"bech32\",\"type\":\"string\"}],\"name\":\"bech32StaticFn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"getShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405260666000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506065600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506065600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503480156100d557600080fd5b50610cde806100e56000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80630d1b3daf1461005c5780632e7f454b1461008c57806340d3051f146100bc57806369163d55146100ec5780639494c8ae1461011d575b600080fd5b610076600480360381019061007191906106fc565b610139565b6040516100839190610a77565b60405180910390f35b6100a660048036038101906100a19190610785565b6101ef565b6040516100b39190610961565b60405180910390f35b6100d660048036038101906100d19190610785565b6102a3565b6040516100e391906109ea565b60405180910390f35b61010660048036038101906101019190610785565b6103c7565b604051610114929190610a05565b60405180910390f35b610137600480360381019061013291906107ce565b6104f6565b005b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630d1b3daf84846040518363ffffffff1660e01b815260040161019792919061097c565b60206040518083038186803b1580156101af57600080fd5b505afa1580156101c3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101e7919061082a565b905092915050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e4e2a4ec836040518263ffffffff1660e01b815260040161024c9190610a35565b60206040518083038186803b15801561026457600080fd5b505afa158015610278573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061029c91906106cf565b9050919050565b600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16836040516024016102f09190610a35565b6040516020818303038152906040527fe4e2a4ec000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161037a919061094a565b600060405180830381855afa9150503d80600081146103b5576040519150601f19603f3d011682016040523d82523d6000602084013e6103ba565b606091505b5050905080915050919050565b60006060600080600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16856040516024016104189190610a35565b6040516020818303038152906040527fe4e2a4ec000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516104a2919061094a565b6000604051808303816000865af19150503d80600081146104df576040519150601f19603f3d011682016040523d82523d6000602084013e6104e4565b606091505b50915091508181935093505050915091565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166390b8436f3385856040518463ffffffff1660e01b8152600401610556939291906109ac565b602060405180830381600087803b15801561057057600080fd5b505af1158015610584573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105a89190610758565b905060011515811515146105f1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105e890610a57565b60405180910390fd5b505050565b600061060961060484610ab7565b610a92565b90508281526020810184848401111561062557610624610c1a565b5b610630848285610b73565b509392505050565b60008135905061064781610c63565b92915050565b60008151905061065c81610c63565b92915050565b60008151905061067181610c7a565b92915050565b600082601f83011261068c5761068b610c15565b5b813561069c8482602086016105f6565b91505092915050565b6000813590506106b481610c91565b92915050565b6000815190506106c981610c91565b92915050565b6000602082840312156106e5576106e4610c24565b5b60006106f38482850161064d565b91505092915050565b6000806040838503121561071357610712610c24565b5b600061072185828601610638565b925050602083013567ffffffffffffffff81111561074257610741610c1f565b5b61074e85828601610677565b9150509250929050565b60006020828403121561076e5761076d610c24565b5b600061077c84828501610662565b91505092915050565b60006020828403121561079b5761079a610c24565b5b600082013567ffffffffffffffff8111156107b9576107b8610c1f565b5b6107c584828501610677565b91505092915050565b600080604083850312156107e5576107e4610c24565b5b600083013567ffffffffffffffff81111561080357610802610c1f565b5b61080f85828601610677565b9250506020610820858286016106a5565b9150509250929050565b6000602082840312156108405761083f610c24565b5b600061084e848285016106ba565b91505092915050565b61086081610b2b565b82525050565b61086f81610b3d565b82525050565b600061088082610ae8565b61088a8185610afe565b935061089a818560208601610b82565b6108a381610c29565b840191505092915050565b60006108b982610ae8565b6108c38185610b0f565b93506108d3818560208601610b82565b80840191505092915050565b60006108ea82610af3565b6108f48185610b1a565b9350610904818560208601610b82565b61090d81610c29565b840191505092915050565b6000610925600e83610b1a565b915061093082610c3a565b602082019050919050565b61094481610b69565b82525050565b600061095682846108ae565b915081905092915050565b60006020820190506109766000830184610857565b92915050565b60006040820190506109916000830185610857565b81810360208301526109a381846108df565b90509392505050565b60006060820190506109c16000830186610857565b81810360208301526109d381856108df565b90506109e2604083018461093b565b949350505050565b60006020820190506109ff6000830184610866565b92915050565b6000604082019050610a1a6000830185610866565b8181036020830152610a2c8184610875565b90509392505050565b60006020820190508181036000830152610a4f81846108df565b905092915050565b60006020820190508181036000830152610a7081610918565b9050919050565b6000602082019050610a8c600083018461093b565b92915050565b6000610a9c610aad565b9050610aa88282610bb5565b919050565b6000604051905090565b600067ffffffffffffffff821115610ad257610ad1610be6565b5b610adb82610c29565b9050602081019050919050565b600081519050919050565b600081519050919050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b6000610b3682610b49565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b83811015610ba0578082015181840152602081019050610b85565b83811115610baf576000848401525b50505050565b610bbe82610c29565b810181811067ffffffffffffffff82111715610bdd57610bdc610be6565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f7374616b696e67206661696c6564000000000000000000000000000000000000600082015250565b610c6c81610b2b565b8114610c7757600080fd5b50565b610c8381610b3d565b8114610c8e57600080fd5b50565b610c9a81610b69565b8114610ca557600080fd5b5056fea26469706673582212204446df03d7396eb13f306dd2746e225c6d97f151c9a79b5b9051640e5c3ce28164736f6c63430008070033",
}

// TestStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use TestStakingMetaData.ABI instead.
var TestStakingABI = TestStakingMetaData.ABI

// TestStakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestStakingMetaData.Bin instead.
var TestStakingBin = TestStakingMetaData.Bin

// DeployTestStaking deploys a new Ethereum contract, binding an instance of TestStaking to it.
func DeployTestStaking(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestStaking, error) {
	parsed, err := TestStakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestStakingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestStaking{TestStakingCaller: TestStakingCaller{contract: contract}, TestStakingTransactor: TestStakingTransactor{contract: contract}, TestStakingFilterer: TestStakingFilterer{contract: contract}}, nil
}

// TestStaking is an auto generated Go binding around an Ethereum contract.
type TestStaking struct {
	TestStakingCaller     // Read-only binding to the contract
	TestStakingTransactor // Write-only binding to the contract
	TestStakingFilterer   // Log filterer for contract events
}

// TestStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestStakingSession struct {
	Contract     *TestStaking      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestStakingCallerSession struct {
	Contract *TestStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// TestStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestStakingTransactorSession struct {
	Contract     *TestStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// TestStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestStakingRaw struct {
	Contract *TestStaking // Generic contract binding to access the raw methods on
}

// TestStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestStakingCallerRaw struct {
	Contract *TestStakingCaller // Generic read-only contract binding to access the raw methods on
}

// TestStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestStakingTransactorRaw struct {
	Contract *TestStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestStaking creates a new instance of TestStaking, bound to a specific deployed contract.
func NewTestStaking(address common.Address, backend bind.ContractBackend) (*TestStaking, error) {
	contract, err := bindTestStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestStaking{TestStakingCaller: TestStakingCaller{contract: contract}, TestStakingTransactor: TestStakingTransactor{contract: contract}, TestStakingFilterer: TestStakingFilterer{contract: contract}}, nil
}

// NewTestStakingCaller creates a new read-only instance of TestStaking, bound to a specific deployed contract.
func NewTestStakingCaller(address common.Address, caller bind.ContractCaller) (*TestStakingCaller, error) {
	contract, err := bindTestStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestStakingCaller{contract: contract}, nil
}

// NewTestStakingTransactor creates a new write-only instance of TestStaking, bound to a specific deployed contract.
func NewTestStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*TestStakingTransactor, error) {
	contract, err := bindTestStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestStakingTransactor{contract: contract}, nil
}

// NewTestStakingFilterer creates a new log filterer instance of TestStaking, bound to a specific deployed contract.
func NewTestStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*TestStakingFilterer, error) {
	contract, err := bindTestStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestStakingFilterer{contract: contract}, nil
}

// bindTestStaking binds a generic wrapper to an already deployed contract.
func bindTestStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestStakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestStaking *TestStakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestStaking.Contract.TestStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestStaking *TestStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestStaking.Contract.TestStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestStaking *TestStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestStaking.Contract.TestStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestStaking *TestStakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestStaking *TestStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestStaking *TestStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestStaking.Contract.contract.Transact(opts, method, params...)
}

// Bech32Fn is a free data retrieval call binding the contract method 0x2e7f454b.
//
// Solidity: function bech32Fn(string bech32) view returns(address addr)
func (_TestStaking *TestStakingCaller) Bech32Fn(opts *bind.CallOpts, bech32 string) (common.Address, error) {
	var out []interface{}
	err := _TestStaking.contract.Call(opts, &out, "bech32Fn", bech32)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bech32Fn is a free data retrieval call binding the contract method 0x2e7f454b.
//
// Solidity: function bech32Fn(string bech32) view returns(address addr)
func (_TestStaking *TestStakingSession) Bech32Fn(bech32 string) (common.Address, error) {
	return _TestStaking.Contract.Bech32Fn(&_TestStaking.CallOpts, bech32)
}

// Bech32Fn is a free data retrieval call binding the contract method 0x2e7f454b.
//
// Solidity: function bech32Fn(string bech32) view returns(address addr)
func (_TestStaking *TestStakingCallerSession) Bech32Fn(bech32 string) (common.Address, error) {
	return _TestStaking.Contract.Bech32Fn(&_TestStaking.CallOpts, bech32)
}

// Bech32StaticFn is a free data retrieval call binding the contract method 0x40d3051f.
//
// Solidity: function bech32StaticFn(string bech32) view returns(bool)
func (_TestStaking *TestStakingCaller) Bech32StaticFn(opts *bind.CallOpts, bech32 string) (bool, error) {
	var out []interface{}
	err := _TestStaking.contract.Call(opts, &out, "bech32StaticFn", bech32)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Bech32StaticFn is a free data retrieval call binding the contract method 0x40d3051f.
//
// Solidity: function bech32StaticFn(string bech32) view returns(bool)
func (_TestStaking *TestStakingSession) Bech32StaticFn(bech32 string) (bool, error) {
	return _TestStaking.Contract.Bech32StaticFn(&_TestStaking.CallOpts, bech32)
}

// Bech32StaticFn is a free data retrieval call binding the contract method 0x40d3051f.
//
// Solidity: function bech32StaticFn(string bech32) view returns(bool)
func (_TestStaking *TestStakingCallerSession) Bech32StaticFn(bech32 string) (bool, error) {
	return _TestStaking.Contract.Bech32StaticFn(&_TestStaking.CallOpts, bech32)
}

// GetShares is a free data retrieval call binding the contract method 0x0d1b3daf.
//
// Solidity: function getShares(address staker, string validator) view returns(uint256 shares)
func (_TestStaking *TestStakingCaller) GetShares(opts *bind.CallOpts, staker common.Address, validator string) (*big.Int, error) {
	var out []interface{}
	err := _TestStaking.contract.Call(opts, &out, "getShares", staker, validator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetShares is a free data retrieval call binding the contract method 0x0d1b3daf.
//
// Solidity: function getShares(address staker, string validator) view returns(uint256 shares)
func (_TestStaking *TestStakingSession) GetShares(staker common.Address, validator string) (*big.Int, error) {
	return _TestStaking.Contract.GetShares(&_TestStaking.CallOpts, staker, validator)
}

// GetShares is a free data retrieval call binding the contract method 0x0d1b3daf.
//
// Solidity: function getShares(address staker, string validator) view returns(uint256 shares)
func (_TestStaking *TestStakingCallerSession) GetShares(staker common.Address, validator string) (*big.Int, error) {
	return _TestStaking.Contract.GetShares(&_TestStaking.CallOpts, staker, validator)
}

// Bech32CallFn is a paid mutator transaction binding the contract method 0x69163d55.
//
// Solidity: function bech32CallFn(string bech32) returns(bool, bytes)
func (_TestStaking *TestStakingTransactor) Bech32CallFn(opts *bind.TransactOpts, bech32 string) (*types.Transaction, error) {
	return _TestStaking.contract.Transact(opts, "bech32CallFn", bech32)
}

// Bech32CallFn is a paid mutator transaction binding the contract method 0x69163d55.
//
// Solidity: function bech32CallFn(string bech32) returns(bool, bytes)
func (_TestStaking *TestStakingSession) Bech32CallFn(bech32 string) (*types.Transaction, error) {
	return _TestStaking.Contract.Bech32CallFn(&_TestStaking.TransactOpts, bech32)
}

// Bech32CallFn is a paid mutator transaction binding the contract method 0x69163d55.
//
// Solidity: function bech32CallFn(string bech32) returns(bool, bytes)
func (_TestStaking *TestStakingTransactorSession) Bech32CallFn(bech32 string) (*types.Transaction, error) {
	return _TestStaking.Contract.Bech32CallFn(&_TestStaking.TransactOpts, bech32)
}

// Stake is a paid mutator transaction binding the contract method 0x9494c8ae.
//
// Solidity: function stake(string validator, uint256 amount) returns()
func (_TestStaking *TestStakingTransactor) Stake(opts *bind.TransactOpts, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.contract.Transact(opts, "stake", validator, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x9494c8ae.
//
// Solidity: function stake(string validator, uint256 amount) returns()
func (_TestStaking *TestStakingSession) Stake(validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.Stake(&_TestStaking.TransactOpts, validator, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x9494c8ae.
//
// Solidity: function stake(string validator, uint256 amount) returns()
func (_TestStaking *TestStakingTransactorSession) Stake(validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.Stake(&_TestStaking.TransactOpts, validator, amount)
}
