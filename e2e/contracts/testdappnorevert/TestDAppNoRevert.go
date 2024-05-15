// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testdappnorevert

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

// ZetaInterfacesZetaMessage is an auto generated low-level Go binding around an user-defined struct.
type ZetaInterfacesZetaMessage struct {
	ZetaTxSenderAddress []byte
	SourceChainId       *big.Int
	DestinationAddress  common.Address
	ZetaValue           *big.Int
	Message             []byte
}

// TestDAppNoRevertMetaData contains all meta data concerning the TestDAppNoRevert contract.
var TestDAppNoRevertMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_connector\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_zetaToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ErrorTransferringZeta\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMessageType\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"HelloWorldEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"RevertedHelloWorldEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"HELLO_WORLD_MESSAGE_TYPE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"connector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"zetaTxSenderAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"destinationAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"zetaValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"internalType\":\"structZetaInterfaces.ZetaMessage\",\"name\":\"zetaMessage\",\"type\":\"tuple\"}],\"name\":\"onZetaMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"destinationAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"doRevert\",\"type\":\"bool\"}],\"name\":\"sendHelloWorld\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zeta\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b50604051610c8b380380610c8b83398181016040528101906100319190610115565b815f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050610153565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100e4826100bb565b9050919050565b6100f4816100da565b81146100fe575f80fd5b50565b5f8151905061010f816100eb565b92915050565b5f806040838503121561012b5761012a6100b7565b5b5f61013885828601610101565b925050602061014985828601610101565b9150509250929050565b610b2b806101605f395ff3fe608060405260043610610049575f3560e01c80633749c51a1461004d5780637caca3041461007557806383f3084f146100915780638ac44a3f146100bb578063e8f9cb3a146100e5575b5f80fd5b348015610058575f80fd5b50610073600480360381019061006e9190610514565b61010f565b005b61008f600480360381019061008a919061061d565b6101a6565b005b34801561009c575f80fd5b506100a561047e565b6040516100b29190610690565b60405180910390f35b3480156100c6575f80fd5b506100cf6104a1565b6040516100dc91906106c1565b60405180910390f35b3480156100f0575f80fd5b506100f96104c5565b6040516101069190610690565b60405180910390f35b5f81806080019061012091906106e6565b81019061012d9190610772565b9150505f151581151514610176576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161016d9061080a565b60405180910390fd5b7f3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d160405160405180910390a15050565b5f60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663095ea7b35f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff16856040518363ffffffff1660e01b8152600401610222929190610837565b6020604051808303815f875af115801561023e573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906102629190610872565b90505f60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330876040518463ffffffff1660e01b81526004016102c39392919061089d565b6020604051808303815f875af11580156102df573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103039190610872565b905081801561030f5750805b610345576040517f2bd0ba5000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ec0269016040518060c00160405280888152602001896040516020016103a29190610917565b60405160208183030381529060405281526020016203d09081526020017f6e0182194bb1deba01849afd3e035a0b70ce7cb069e482ee663519c76cf569b4876040516020016103f2929190610940565b604051602081830303815290604052815260200187815260200160405160200161041b9061098a565b6040516020818303038152906040528152506040518263ffffffff1660e01b81526004016104499190610ad5565b5f604051808303815f87803b158015610460575f80fd5b505af1158015610472573d5f803e3d5ffd5b50505050505050505050565b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b7f6e0182194bb1deba01849afd3e035a0b70ce7cb069e482ee663519c76cf569b481565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f80fd5b5f80fd5b5f80fd5b5f60a0828403121561050b5761050a6104f2565b5b81905092915050565b5f60208284031215610529576105286104ea565b5b5f82013567ffffffffffffffff811115610546576105456104ee565b5b610552848285016104f6565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6105848261055b565b9050919050565b6105948161057a565b811461059e575f80fd5b50565b5f813590506105af8161058b565b92915050565b5f819050919050565b6105c7816105b5565b81146105d1575f80fd5b50565b5f813590506105e2816105be565b92915050565b5f8115159050919050565b6105fc816105e8565b8114610606575f80fd5b50565b5f81359050610617816105f3565b92915050565b5f805f8060808587031215610635576106346104ea565b5b5f610642878288016105a1565b9450506020610653878288016105d4565b9350506040610664878288016105d4565b925050606061067587828801610609565b91505092959194509250565b61068a8161057a565b82525050565b5f6020820190506106a35f830184610681565b92915050565b5f819050919050565b6106bb816106a9565b82525050565b5f6020820190506106d45f8301846106b2565b92915050565b5f80fd5b5f80fd5b5f80fd5b5f8083356001602003843603038112610702576107016106da565b5b80840192508235915067ffffffffffffffff821115610724576107236106de565b5b6020830192506001820236038313156107405761073f6106e2565b5b509250929050565b610751816106a9565b811461075b575f80fd5b50565b5f8135905061076c81610748565b92915050565b5f8060408385031215610788576107876104ea565b5b5f6107958582860161075e565b92505060206107a685828601610609565b9150509250929050565b5f82825260208201905092915050565b7f6d657373616765207361797320726576657274000000000000000000000000005f82015250565b5f6107f46013836107b0565b91506107ff826107c0565b602082019050919050565b5f6020820190508181035f830152610821816107e8565b9050919050565b610831816105b5565b82525050565b5f60408201905061084a5f830185610681565b6108576020830184610828565b9392505050565b5f8151905061086c816105f3565b92915050565b5f60208284031215610887576108866104ea565b5b5f6108948482850161085e565b91505092915050565b5f6060820190506108b05f830186610681565b6108bd6020830185610681565b6108ca6040830184610828565b949350505050565b5f8160601b9050919050565b5f6108e8826108d2565b9050919050565b5f6108f9826108de565b9050919050565b61091161090c8261057a565b6108ef565b82525050565b5f6109228284610900565b60148201915081905092915050565b61093a816105e8565b82525050565b5f6040820190506109535f8301856106b2565b6109606020830184610931565b9392505050565b50565b5f6109755f836107b0565b915061098082610967565b5f82019050919050565b5f6020820190508181035f8301526109a18161096a565b9050919050565b6109b1816105b5565b82525050565b5f81519050919050565b5f82825260208201905092915050565b5f5b838110156109ee5780820151818401526020810190506109d3565b5f8484015250505050565b5f601f19601f8301169050919050565b5f610a13826109b7565b610a1d81856109c1565b9350610a2d8185602086016109d1565b610a36816109f9565b840191505092915050565b5f60c083015f830151610a565f8601826109a8565b5060208301518482036020860152610a6e8282610a09565b9150506040830151610a8360408601826109a8565b5060608301518482036060860152610a9b8282610a09565b9150506080830151610ab060808601826109a8565b5060a083015184820360a0860152610ac88282610a09565b9150508091505092915050565b5f6020820190508181035f830152610aed8184610a41565b90509291505056fea264697066735822122069ddd13287d121449d54a6d0dcea4df3a319a78e1dbcd04aff8d2553fb893bd164736f6c63430008170033",
}

// TestDAppNoRevertABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDAppNoRevertMetaData.ABI instead.
var TestDAppNoRevertABI = TestDAppNoRevertMetaData.ABI

// TestDAppNoRevertBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestDAppNoRevertMetaData.Bin instead.
var TestDAppNoRevertBin = TestDAppNoRevertMetaData.Bin

// DeployTestDAppNoRevert deploys a new Ethereum contract, binding an instance of TestDAppNoRevert to it.
func DeployTestDAppNoRevert(auth *bind.TransactOpts, backend bind.ContractBackend, _connector common.Address, _zetaToken common.Address) (common.Address, *types.Transaction, *TestDAppNoRevert, error) {
	parsed, err := TestDAppNoRevertMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDAppNoRevertBin), backend, _connector, _zetaToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestDAppNoRevert{TestDAppNoRevertCaller: TestDAppNoRevertCaller{contract: contract}, TestDAppNoRevertTransactor: TestDAppNoRevertTransactor{contract: contract}, TestDAppNoRevertFilterer: TestDAppNoRevertFilterer{contract: contract}}, nil
}

// TestDAppNoRevert is an auto generated Go binding around an Ethereum contract.
type TestDAppNoRevert struct {
	TestDAppNoRevertCaller     // Read-only binding to the contract
	TestDAppNoRevertTransactor // Write-only binding to the contract
	TestDAppNoRevertFilterer   // Log filterer for contract events
}

// TestDAppNoRevertCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestDAppNoRevertCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppNoRevertTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestDAppNoRevertTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppNoRevertFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestDAppNoRevertFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppNoRevertSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestDAppNoRevertSession struct {
	Contract     *TestDAppNoRevert // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestDAppNoRevertCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestDAppNoRevertCallerSession struct {
	Contract *TestDAppNoRevertCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// TestDAppNoRevertTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestDAppNoRevertTransactorSession struct {
	Contract     *TestDAppNoRevertTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// TestDAppNoRevertRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestDAppNoRevertRaw struct {
	Contract *TestDAppNoRevert // Generic contract binding to access the raw methods on
}

// TestDAppNoRevertCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestDAppNoRevertCallerRaw struct {
	Contract *TestDAppNoRevertCaller // Generic read-only contract binding to access the raw methods on
}

// TestDAppNoRevertTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestDAppNoRevertTransactorRaw struct {
	Contract *TestDAppNoRevertTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestDAppNoRevert creates a new instance of TestDAppNoRevert, bound to a specific deployed contract.
func NewTestDAppNoRevert(address common.Address, backend bind.ContractBackend) (*TestDAppNoRevert, error) {
	contract, err := bindTestDAppNoRevert(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevert{TestDAppNoRevertCaller: TestDAppNoRevertCaller{contract: contract}, TestDAppNoRevertTransactor: TestDAppNoRevertTransactor{contract: contract}, TestDAppNoRevertFilterer: TestDAppNoRevertFilterer{contract: contract}}, nil
}

// NewTestDAppNoRevertCaller creates a new read-only instance of TestDAppNoRevert, bound to a specific deployed contract.
func NewTestDAppNoRevertCaller(address common.Address, caller bind.ContractCaller) (*TestDAppNoRevertCaller, error) {
	contract, err := bindTestDAppNoRevert(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevertCaller{contract: contract}, nil
}

// NewTestDAppNoRevertTransactor creates a new write-only instance of TestDAppNoRevert, bound to a specific deployed contract.
func NewTestDAppNoRevertTransactor(address common.Address, transactor bind.ContractTransactor) (*TestDAppNoRevertTransactor, error) {
	contract, err := bindTestDAppNoRevert(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevertTransactor{contract: contract}, nil
}

// NewTestDAppNoRevertFilterer creates a new log filterer instance of TestDAppNoRevert, bound to a specific deployed contract.
func NewTestDAppNoRevertFilterer(address common.Address, filterer bind.ContractFilterer) (*TestDAppNoRevertFilterer, error) {
	contract, err := bindTestDAppNoRevert(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevertFilterer{contract: contract}, nil
}

// bindTestDAppNoRevert binds a generic wrapper to an already deployed contract.
func bindTestDAppNoRevert(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestDAppNoRevertMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDAppNoRevert *TestDAppNoRevertRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDAppNoRevert.Contract.TestDAppNoRevertCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDAppNoRevert *TestDAppNoRevertRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.TestDAppNoRevertTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDAppNoRevert *TestDAppNoRevertRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.TestDAppNoRevertTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDAppNoRevert *TestDAppNoRevertCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDAppNoRevert.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDAppNoRevert *TestDAppNoRevertTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDAppNoRevert *TestDAppNoRevertTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.contract.Transact(opts, method, params...)
}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDAppNoRevert *TestDAppNoRevertCaller) HELLOWORLDMESSAGETYPE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TestDAppNoRevert.contract.Call(opts, &out, "HELLO_WORLD_MESSAGE_TYPE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDAppNoRevert *TestDAppNoRevertSession) HELLOWORLDMESSAGETYPE() ([32]byte, error) {
	return _TestDAppNoRevert.Contract.HELLOWORLDMESSAGETYPE(&_TestDAppNoRevert.CallOpts)
}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDAppNoRevert *TestDAppNoRevertCallerSession) HELLOWORLDMESSAGETYPE() ([32]byte, error) {
	return _TestDAppNoRevert.Contract.HELLOWORLDMESSAGETYPE(&_TestDAppNoRevert.CallOpts)
}

// Connector is a free data retrieval call binding the contract method 0x83f3084f.
//
// Solidity: function connector() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertCaller) Connector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestDAppNoRevert.contract.Call(opts, &out, "connector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Connector is a free data retrieval call binding the contract method 0x83f3084f.
//
// Solidity: function connector() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertSession) Connector() (common.Address, error) {
	return _TestDAppNoRevert.Contract.Connector(&_TestDAppNoRevert.CallOpts)
}

// Connector is a free data retrieval call binding the contract method 0x83f3084f.
//
// Solidity: function connector() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertCallerSession) Connector() (common.Address, error) {
	return _TestDAppNoRevert.Contract.Connector(&_TestDAppNoRevert.CallOpts)
}

// Zeta is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function zeta() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertCaller) Zeta(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestDAppNoRevert.contract.Call(opts, &out, "zeta")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Zeta is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function zeta() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertSession) Zeta() (common.Address, error) {
	return _TestDAppNoRevert.Contract.Zeta(&_TestDAppNoRevert.CallOpts)
}

// Zeta is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function zeta() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertCallerSession) Zeta() (common.Address, error) {
	return _TestDAppNoRevert.Contract.Zeta(&_TestDAppNoRevert.CallOpts)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onZetaMessage((bytes,uint256,address,uint256,bytes) zetaMessage) returns()
func (_TestDAppNoRevert *TestDAppNoRevertTransactor) OnZetaMessage(opts *bind.TransactOpts, zetaMessage ZetaInterfacesZetaMessage) (*types.Transaction, error) {
	return _TestDAppNoRevert.contract.Transact(opts, "onZetaMessage", zetaMessage)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onZetaMessage((bytes,uint256,address,uint256,bytes) zetaMessage) returns()
func (_TestDAppNoRevert *TestDAppNoRevertSession) OnZetaMessage(zetaMessage ZetaInterfacesZetaMessage) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.OnZetaMessage(&_TestDAppNoRevert.TransactOpts, zetaMessage)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onZetaMessage((bytes,uint256,address,uint256,bytes) zetaMessage) returns()
func (_TestDAppNoRevert *TestDAppNoRevertTransactorSession) OnZetaMessage(zetaMessage ZetaInterfacesZetaMessage) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.OnZetaMessage(&_TestDAppNoRevert.TransactOpts, zetaMessage)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDAppNoRevert *TestDAppNoRevertTransactor) SendHelloWorld(opts *bind.TransactOpts, destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDAppNoRevert.contract.Transact(opts, "sendHelloWorld", destinationAddress, destinationChainId, value, doRevert)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDAppNoRevert *TestDAppNoRevertSession) SendHelloWorld(destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.SendHelloWorld(&_TestDAppNoRevert.TransactOpts, destinationAddress, destinationChainId, value, doRevert)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDAppNoRevert *TestDAppNoRevertTransactorSession) SendHelloWorld(destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.SendHelloWorld(&_TestDAppNoRevert.TransactOpts, destinationAddress, destinationChainId, value, doRevert)
}

// TestDAppNoRevertHelloWorldEventIterator is returned from FilterHelloWorldEvent and is used to iterate over the raw logs and unpacked data for HelloWorldEvent events raised by the TestDAppNoRevert contract.
type TestDAppNoRevertHelloWorldEventIterator struct {
	Event *TestDAppNoRevertHelloWorldEvent // Event containing the contract specifics and raw log

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
func (it *TestDAppNoRevertHelloWorldEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestDAppNoRevertHelloWorldEvent)
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
		it.Event = new(TestDAppNoRevertHelloWorldEvent)
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
func (it *TestDAppNoRevertHelloWorldEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestDAppNoRevertHelloWorldEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestDAppNoRevertHelloWorldEvent represents a HelloWorldEvent event raised by the TestDAppNoRevert contract.
type TestDAppNoRevertHelloWorldEvent struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterHelloWorldEvent is a free log retrieval operation binding the contract event 0x3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d1.
//
// Solidity: event HelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) FilterHelloWorldEvent(opts *bind.FilterOpts) (*TestDAppNoRevertHelloWorldEventIterator, error) {

	logs, sub, err := _TestDAppNoRevert.contract.FilterLogs(opts, "HelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevertHelloWorldEventIterator{contract: _TestDAppNoRevert.contract, event: "HelloWorldEvent", logs: logs, sub: sub}, nil
}

// WatchHelloWorldEvent is a free log subscription operation binding the contract event 0x3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d1.
//
// Solidity: event HelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) WatchHelloWorldEvent(opts *bind.WatchOpts, sink chan<- *TestDAppNoRevertHelloWorldEvent) (event.Subscription, error) {

	logs, sub, err := _TestDAppNoRevert.contract.WatchLogs(opts, "HelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestDAppNoRevertHelloWorldEvent)
				if err := _TestDAppNoRevert.contract.UnpackLog(event, "HelloWorldEvent", log); err != nil {
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

// ParseHelloWorldEvent is a log parse operation binding the contract event 0x3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d1.
//
// Solidity: event HelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) ParseHelloWorldEvent(log types.Log) (*TestDAppNoRevertHelloWorldEvent, error) {
	event := new(TestDAppNoRevertHelloWorldEvent)
	if err := _TestDAppNoRevert.contract.UnpackLog(event, "HelloWorldEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestDAppNoRevertRevertedHelloWorldEventIterator is returned from FilterRevertedHelloWorldEvent and is used to iterate over the raw logs and unpacked data for RevertedHelloWorldEvent events raised by the TestDAppNoRevert contract.
type TestDAppNoRevertRevertedHelloWorldEventIterator struct {
	Event *TestDAppNoRevertRevertedHelloWorldEvent // Event containing the contract specifics and raw log

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
func (it *TestDAppNoRevertRevertedHelloWorldEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestDAppNoRevertRevertedHelloWorldEvent)
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
		it.Event = new(TestDAppNoRevertRevertedHelloWorldEvent)
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
func (it *TestDAppNoRevertRevertedHelloWorldEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestDAppNoRevertRevertedHelloWorldEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestDAppNoRevertRevertedHelloWorldEvent represents a RevertedHelloWorldEvent event raised by the TestDAppNoRevert contract.
type TestDAppNoRevertRevertedHelloWorldEvent struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRevertedHelloWorldEvent is a free log retrieval operation binding the contract event 0x4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a521.
//
// Solidity: event RevertedHelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) FilterRevertedHelloWorldEvent(opts *bind.FilterOpts) (*TestDAppNoRevertRevertedHelloWorldEventIterator, error) {

	logs, sub, err := _TestDAppNoRevert.contract.FilterLogs(opts, "RevertedHelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevertRevertedHelloWorldEventIterator{contract: _TestDAppNoRevert.contract, event: "RevertedHelloWorldEvent", logs: logs, sub: sub}, nil
}

// WatchRevertedHelloWorldEvent is a free log subscription operation binding the contract event 0x4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a521.
//
// Solidity: event RevertedHelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) WatchRevertedHelloWorldEvent(opts *bind.WatchOpts, sink chan<- *TestDAppNoRevertRevertedHelloWorldEvent) (event.Subscription, error) {

	logs, sub, err := _TestDAppNoRevert.contract.WatchLogs(opts, "RevertedHelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestDAppNoRevertRevertedHelloWorldEvent)
				if err := _TestDAppNoRevert.contract.UnpackLog(event, "RevertedHelloWorldEvent", log); err != nil {
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

// ParseRevertedHelloWorldEvent is a log parse operation binding the contract event 0x4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a521.
//
// Solidity: event RevertedHelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) ParseRevertedHelloWorldEvent(log types.Log) (*TestDAppNoRevertRevertedHelloWorldEvent, error) {
	event := new(TestDAppNoRevertRevertedHelloWorldEvent)
	if err := _TestDAppNoRevert.contract.UnpackLog(event, "RevertedHelloWorldEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
