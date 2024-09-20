// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testgatewayzevmcaller

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

// CallOptions is an auto generated low-level Go binding around an user-defined struct.
type CallOptions struct {
	GasLimit        *big.Int
	IsArbitraryCall bool
}

// RevertOptions is an auto generated low-level Go binding around an user-defined struct.
type RevertOptions struct {
	RevertAddress    common.Address
	CallOnRevert     bool
	AbortAddress     common.Address
	RevertMessage    []byte
	OnRevertGasLimit *big.Int
}

// TestGatewayZEVMCallerMetaData contains all meta data concerning the TestGatewayZEVMCaller contract.
var TestGatewayZEVMCallerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"gatewayZEVMAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isArbitraryCall\",\"type\":\"bool\"}],\"internalType\":\"structCallOptions\",\"name\":\"callOptions\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"revertAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"callOnRevert\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"abortAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"onRevertGasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structRevertOptions\",\"name\":\"revertOptions\",\"type\":\"tuple\"}],\"name\":\"callGatewayZEVM\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610a57380380610a57833981810160405281019061003291906100db565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050610108565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100a88261007d565b9050919050565b6100b88161009d565b81146100c357600080fd5b50565b6000815190506100d5816100af565b92915050565b6000602082840312156100f1576100f0610078565b5b60006100ff848285016100c6565b91505092915050565b610940806101176000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806325859e6214610030575b600080fd5b61004a600480360381019061004591906103eb565b61004c565b005b8473ffffffffffffffffffffffffffffffffffffffff1663095ea7b360008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1667016345785d8a00006040518363ffffffff1660e01b81526004016100af92919061051b565b6020604051808303816000875af11580156100ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100f2919061057c565b5060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166306cb89838787878787876040518763ffffffff1660e01b8152600401610156969594939291906108a0565b600060405180830381600087803b15801561017057600080fd5b505af1158015610184573d6000803e3d6000fd5b50505050505050505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6101f7826101ae565b810181811067ffffffffffffffff82111715610216576102156101bf565b5b80604052505050565b6000610229610190565b905061023582826101ee565b919050565b600067ffffffffffffffff821115610255576102546101bf565b5b61025e826101ae565b9050602081019050919050565b82818337600083830152505050565b600061028d6102888461023a565b61021f565b9050828152602081018484840111156102a9576102a86101a9565b5b6102b484828561026b565b509392505050565b600082601f8301126102d1576102d06101a4565b5b81356102e184826020860161027a565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610315826102ea565b9050919050565b6103258161030a565b811461033057600080fd5b50565b6000813590506103428161031c565b92915050565b600080fd5b600080fd5b60008083601f840112610368576103676101a4565b5b8235905067ffffffffffffffff81111561038557610384610348565b5b6020830191508360018202830111156103a1576103a061034d565b5b9250929050565b600080fd5b6000604082840312156103c3576103c26103a8565b5b81905092915050565b600060a082840312156103e2576103e16103a8565b5b81905092915050565b60008060008060008060c087890312156104085761040761019a565b5b600087013567ffffffffffffffff8111156104265761042561019f565b5b61043289828a016102bc565b965050602061044389828a01610333565b955050604087013567ffffffffffffffff8111156104645761046361019f565b5b61047089828a01610352565b9450945050606061048389828a016103ad565b92505060a087013567ffffffffffffffff8111156104a4576104a361019f565b5b6104b089828a016103cc565b9150509295509295509295565b6104c68161030a565b82525050565b6000819050919050565b6000819050919050565b6000819050919050565b60006105056105006104fb846104cc565b6104e0565b6104d6565b9050919050565b610515816104ea565b82525050565b600060408201905061053060008301856104bd565b61053d602083018461050c565b9392505050565b60008115159050919050565b61055981610544565b811461056457600080fd5b50565b60008151905061057681610550565b92915050565b6000602082840312156105925761059161019a565b5b60006105a084828501610567565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b838110156105e35780820151818401526020810190506105c8565b838111156105f2576000848401525b50505050565b6000610603826105a9565b61060d81856105b4565b935061061d8185602086016105c5565b610626816101ae565b840191505092915050565b600061063d83856105b4565b935061064a83858461026b565b610653836101ae565b840190509392505050565b610667816104d6565b811461067257600080fd5b50565b6000813590506106848161065e565b92915050565b60006106996020840184610675565b905092915050565b6106aa816104d6565b82525050565b6000813590506106bf81610550565b92915050565b60006106d460208401846106b0565b905092915050565b6106e581610544565b82525050565b604082016106fc600083018361068a565b61070960008501826106a1565b5061071760208301836106c5565b61072460208501826106dc565b50505050565b60006107396020840184610333565b905092915050565b61074a8161030a565b82525050565b600080fd5b600080fd5b600080fd5b6000808335600160200384360303811261077c5761077b61075a565b5b83810192508235915060208301925067ffffffffffffffff8211156107a4576107a3610750565b5b6001820236038413156107ba576107b9610755565b5b509250929050565b600082825260208201905092915050565b60006107df83856107c2565b93506107ec83858461026b565b6107f5836101ae565b840190509392505050565b600060a08301610813600084018461072a565b6108206000860182610741565b5061082e60208401846106c5565b61083b60208601826106dc565b50610849604084018461072a565b6108566040860182610741565b50610864606084018461075f565b85830360608701526108778382846107d3565b92505050610888608084018461068a565b61089560808601826106a1565b508091505092915050565b600060c08201905081810360008301526108ba81896105f8565b90506108c960208301886104bd565b81810360408301526108dc818688610631565b90506108eb60608301856106eb565b81810360a08301526108fd8184610800565b905097965050505050505056fea264697066735822122066a4b53d3b94f8a07a5f39ad45cbdfe04b0de8079b873728f16505cb20c6639064736f6c634300080a0033",
}

// TestGatewayZEVMCallerABI is the input ABI used to generate the binding from.
// Deprecated: Use TestGatewayZEVMCallerMetaData.ABI instead.
var TestGatewayZEVMCallerABI = TestGatewayZEVMCallerMetaData.ABI

// TestGatewayZEVMCallerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestGatewayZEVMCallerMetaData.Bin instead.
var TestGatewayZEVMCallerBin = TestGatewayZEVMCallerMetaData.Bin

// DeployTestGatewayZEVMCaller deploys a new Ethereum contract, binding an instance of TestGatewayZEVMCaller to it.
func DeployTestGatewayZEVMCaller(auth *bind.TransactOpts, backend bind.ContractBackend, gatewayZEVMAddress common.Address) (common.Address, *types.Transaction, *TestGatewayZEVMCaller, error) {
	parsed, err := TestGatewayZEVMCallerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestGatewayZEVMCallerBin), backend, gatewayZEVMAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestGatewayZEVMCaller{TestGatewayZEVMCallerCaller: TestGatewayZEVMCallerCaller{contract: contract}, TestGatewayZEVMCallerTransactor: TestGatewayZEVMCallerTransactor{contract: contract}, TestGatewayZEVMCallerFilterer: TestGatewayZEVMCallerFilterer{contract: contract}}, nil
}

// TestGatewayZEVMCaller is an auto generated Go binding around an Ethereum contract.
type TestGatewayZEVMCaller struct {
	TestGatewayZEVMCallerCaller     // Read-only binding to the contract
	TestGatewayZEVMCallerTransactor // Write-only binding to the contract
	TestGatewayZEVMCallerFilterer   // Log filterer for contract events
}

// TestGatewayZEVMCallerCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestGatewayZEVMCallerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestGatewayZEVMCallerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestGatewayZEVMCallerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestGatewayZEVMCallerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestGatewayZEVMCallerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestGatewayZEVMCallerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestGatewayZEVMCallerSession struct {
	Contract     *TestGatewayZEVMCaller // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// TestGatewayZEVMCallerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestGatewayZEVMCallerCallerSession struct {
	Contract *TestGatewayZEVMCallerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// TestGatewayZEVMCallerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestGatewayZEVMCallerTransactorSession struct {
	Contract     *TestGatewayZEVMCallerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// TestGatewayZEVMCallerRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestGatewayZEVMCallerRaw struct {
	Contract *TestGatewayZEVMCaller // Generic contract binding to access the raw methods on
}

// TestGatewayZEVMCallerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestGatewayZEVMCallerCallerRaw struct {
	Contract *TestGatewayZEVMCallerCaller // Generic read-only contract binding to access the raw methods on
}

// TestGatewayZEVMCallerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestGatewayZEVMCallerTransactorRaw struct {
	Contract *TestGatewayZEVMCallerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestGatewayZEVMCaller creates a new instance of TestGatewayZEVMCaller, bound to a specific deployed contract.
func NewTestGatewayZEVMCaller(address common.Address, backend bind.ContractBackend) (*TestGatewayZEVMCaller, error) {
	contract, err := bindTestGatewayZEVMCaller(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestGatewayZEVMCaller{TestGatewayZEVMCallerCaller: TestGatewayZEVMCallerCaller{contract: contract}, TestGatewayZEVMCallerTransactor: TestGatewayZEVMCallerTransactor{contract: contract}, TestGatewayZEVMCallerFilterer: TestGatewayZEVMCallerFilterer{contract: contract}}, nil
}

// NewTestGatewayZEVMCallerCaller creates a new read-only instance of TestGatewayZEVMCaller, bound to a specific deployed contract.
func NewTestGatewayZEVMCallerCaller(address common.Address, caller bind.ContractCaller) (*TestGatewayZEVMCallerCaller, error) {
	contract, err := bindTestGatewayZEVMCaller(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestGatewayZEVMCallerCaller{contract: contract}, nil
}

// NewTestGatewayZEVMCallerTransactor creates a new write-only instance of TestGatewayZEVMCaller, bound to a specific deployed contract.
func NewTestGatewayZEVMCallerTransactor(address common.Address, transactor bind.ContractTransactor) (*TestGatewayZEVMCallerTransactor, error) {
	contract, err := bindTestGatewayZEVMCaller(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestGatewayZEVMCallerTransactor{contract: contract}, nil
}

// NewTestGatewayZEVMCallerFilterer creates a new log filterer instance of TestGatewayZEVMCaller, bound to a specific deployed contract.
func NewTestGatewayZEVMCallerFilterer(address common.Address, filterer bind.ContractFilterer) (*TestGatewayZEVMCallerFilterer, error) {
	contract, err := bindTestGatewayZEVMCaller(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestGatewayZEVMCallerFilterer{contract: contract}, nil
}

// bindTestGatewayZEVMCaller binds a generic wrapper to an already deployed contract.
func bindTestGatewayZEVMCaller(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestGatewayZEVMCallerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestGatewayZEVMCaller *TestGatewayZEVMCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestGatewayZEVMCaller.Contract.TestGatewayZEVMCallerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestGatewayZEVMCaller *TestGatewayZEVMCallerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestGatewayZEVMCaller.Contract.TestGatewayZEVMCallerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestGatewayZEVMCaller *TestGatewayZEVMCallerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestGatewayZEVMCaller.Contract.TestGatewayZEVMCallerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestGatewayZEVMCaller *TestGatewayZEVMCallerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestGatewayZEVMCaller.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestGatewayZEVMCaller *TestGatewayZEVMCallerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestGatewayZEVMCaller.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestGatewayZEVMCaller *TestGatewayZEVMCallerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestGatewayZEVMCaller.Contract.contract.Transact(opts, method, params...)
}

// CallGatewayZEVM is a paid mutator transaction binding the contract method 0x25859e62.
//
// Solidity: function callGatewayZEVM(bytes receiver, address zrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_TestGatewayZEVMCaller *TestGatewayZEVMCallerTransactor) CallGatewayZEVM(opts *bind.TransactOpts, receiver []byte, zrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _TestGatewayZEVMCaller.contract.Transact(opts, "callGatewayZEVM", receiver, zrc20, message, callOptions, revertOptions)
}

// CallGatewayZEVM is a paid mutator transaction binding the contract method 0x25859e62.
//
// Solidity: function callGatewayZEVM(bytes receiver, address zrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_TestGatewayZEVMCaller *TestGatewayZEVMCallerSession) CallGatewayZEVM(receiver []byte, zrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _TestGatewayZEVMCaller.Contract.CallGatewayZEVM(&_TestGatewayZEVMCaller.TransactOpts, receiver, zrc20, message, callOptions, revertOptions)
}

// CallGatewayZEVM is a paid mutator transaction binding the contract method 0x25859e62.
//
// Solidity: function callGatewayZEVM(bytes receiver, address zrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_TestGatewayZEVMCaller *TestGatewayZEVMCallerTransactorSession) CallGatewayZEVM(receiver []byte, zrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _TestGatewayZEVMCaller.Contract.CallGatewayZEVM(&_TestGatewayZEVMCaller.TransactOpts, receiver, zrc20, message, callOptions, revertOptions)
}
