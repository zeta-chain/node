// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package example

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

// ExamplezContext is an auto generated low-level Go binding around an user-defined struct.
type ExamplezContext struct {
	Sender    []byte
	SenderEVM common.Address
	ChainID   *big.Int
}

// ExampleMetaData contains all meta data concerning the Example contract.
var ExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Foo\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"bar\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevertWithMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevertWithRequire\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doSucceed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMessage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastSender\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"senderEVM\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structExample.zContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"senderEVM\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structExample.zContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5060008081905550610a64806100266000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063d720cb4511610066578063d720cb45146100fa578063dd8e556c14610104578063de43156e1461010e578063fd5ad9651461012a578063febb0f7e1461013457610093565b8063256fec881461009857806332970710146100b65780635bcfd616146100d4578063afc874d2146100f0575b600080fd5b6100a0610152565b6040516100ad9190610412565b60405180910390f35b6100be6101e0565b6040516100cb9190610412565b60405180910390f35b6100ee60048036038101906100e9919061055b565b61026e565b005b6100f86102ae565b005b6101026102e0565b005b61010c61031b565b005b6101286004803603810190610123919061055b565b61035e565b005b610132610372565b005b61013c61037c565b604051610149919061060e565b60405180910390f35b6002805461015f90610658565b80601f016020809104026020016040519081016040528092919081815260200182805461018b90610658565b80156101d85780601f106101ad576101008083540402835291602001916101d8565b820191906000526020600020905b8154815290600101906020018083116101bb57829003601f168201915b505050505081565b600180546101ed90610658565b80601f016020809104026020016040519081016040528092919081815260200182805461021990610658565b80156102665780601f1061023b57610100808354040283529160200191610266565b820191906000526020600020905b81548152906001019060200180831161024957829003601f168201915b505050505081565b8260008190555081816001918261028692919061086f565b50848060000190610297919061094e565b600291826102a692919061086f565b505050505050565b6040517fbfb4ebcf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161031290610a0e565b60405180910390fd5b600061035c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161035390610a0e565b60405180910390fd5b565b61036b858585858561026e565b5050505050565b6001600081905550565b60005481565b600081519050919050565b600082825260208201905092915050565b60005b838110156103bc5780820151818401526020810190506103a1565b60008484015250505050565b6000601f19601f8301169050919050565b60006103e482610382565b6103ee818561038d565b93506103fe81856020860161039e565b610407816103c8565b840191505092915050565b6000602082019050818103600083015261042c81846103d9565b905092915050565b600080fd5b600080fd5b600080fd5b6000606082840312156104595761045861043e565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061048d82610462565b9050919050565b61049d81610482565b81146104a857600080fd5b50565b6000813590506104ba81610494565b92915050565b6000819050919050565b6104d3816104c0565b81146104de57600080fd5b50565b6000813590506104f0816104ca565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f84011261051b5761051a6104f6565b5b8235905067ffffffffffffffff811115610538576105376104fb565b5b60208301915083600182028301111561055457610553610500565b5b9250929050565b60008060008060006080868803121561057757610576610434565b5b600086013567ffffffffffffffff81111561059557610594610439565b5b6105a188828901610443565b95505060206105b2888289016104ab565b94505060406105c3888289016104e1565b935050606086013567ffffffffffffffff8111156105e4576105e3610439565b5b6105f088828901610505565b92509250509295509295909350565b610608816104c0565b82525050565b600060208201905061062360008301846105ff565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061067057607f821691505b60208210810361068357610682610629565b5b50919050565b600082905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026107257fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826106e8565b61072f86836106e8565b95508019841693508086168417925050509392505050565b6000819050919050565b600061076c610767610762846104c0565b610747565b6104c0565b9050919050565b6000819050919050565b61078683610751565b61079a61079282610773565b8484546106f5565b825550505050565b600090565b6107af6107a2565b6107ba81848461077d565b505050565b5b818110156107de576107d36000826107a7565b6001810190506107c0565b5050565b601f821115610823576107f4816106c3565b6107fd846106d8565b8101602085101561080c578190505b610820610818856106d8565b8301826107bf565b50505b505050565b600082821c905092915050565b600061084660001984600802610828565b1980831691505092915050565b600061085f8383610835565b9150826002028217905092915050565b6108798383610689565b67ffffffffffffffff81111561089257610891610694565b5b61089c8254610658565b6108a78282856107e2565b6000601f8311600181146108d657600084156108c4578287013590505b6108ce8582610853565b865550610936565b601f1984166108e4866106c3565b60005b8281101561090c578489013582556001820191506020850194506020810190506108e7565b868310156109295784890135610925601f891682610835565b8355505b6001600288020188555050505b50505050505050565b600080fd5b600080fd5b600080fd5b6000808335600160200384360303811261096b5761096a61093f565b5b80840192508235915067ffffffffffffffff82111561098d5761098c610944565b5b6020830192506001820236038313156109a9576109a8610949565b5b509250929050565b600082825260208201905092915050565b7f666f6f0000000000000000000000000000000000000000000000000000000000600082015250565b60006109f86003836109b1565b9150610a03826109c2565b602082019050919050565b60006020820190508181036000830152610a27816109eb565b905091905056fea26469706673582212202895a4229fcbe4219ffcc5f2b1db026260e6cb3c4d25f7d462793e723d027a3864736f6c634300081a0033",
}

// ExampleABI is the input ABI used to generate the binding from.
// Deprecated: Use ExampleMetaData.ABI instead.
var ExampleABI = ExampleMetaData.ABI

// ExampleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ExampleMetaData.Bin instead.
var ExampleBin = ExampleMetaData.Bin

// DeployExample deploys a new Ethereum contract, binding an instance of Example to it.
func DeployExample(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Example, error) {
	parsed, err := ExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ExampleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Example{ExampleCaller: ExampleCaller{contract: contract}, ExampleTransactor: ExampleTransactor{contract: contract}, ExampleFilterer: ExampleFilterer{contract: contract}}, nil
}

// Example is an auto generated Go binding around an Ethereum contract.
type Example struct {
	ExampleCaller     // Read-only binding to the contract
	ExampleTransactor // Write-only binding to the contract
	ExampleFilterer   // Log filterer for contract events
}

// ExampleCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExampleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExampleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExampleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExampleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExampleSession struct {
	Contract     *Example          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExampleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExampleCallerSession struct {
	Contract *ExampleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ExampleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExampleTransactorSession struct {
	Contract     *ExampleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ExampleRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExampleRaw struct {
	Contract *Example // Generic contract binding to access the raw methods on
}

// ExampleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExampleCallerRaw struct {
	Contract *ExampleCaller // Generic read-only contract binding to access the raw methods on
}

// ExampleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExampleTransactorRaw struct {
	Contract *ExampleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExample creates a new instance of Example, bound to a specific deployed contract.
func NewExample(address common.Address, backend bind.ContractBackend) (*Example, error) {
	contract, err := bindExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Example{ExampleCaller: ExampleCaller{contract: contract}, ExampleTransactor: ExampleTransactor{contract: contract}, ExampleFilterer: ExampleFilterer{contract: contract}}, nil
}

// NewExampleCaller creates a new read-only instance of Example, bound to a specific deployed contract.
func NewExampleCaller(address common.Address, caller bind.ContractCaller) (*ExampleCaller, error) {
	contract, err := bindExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleCaller{contract: contract}, nil
}

// NewExampleTransactor creates a new write-only instance of Example, bound to a specific deployed contract.
func NewExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*ExampleTransactor, error) {
	contract, err := bindExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExampleTransactor{contract: contract}, nil
}

// NewExampleFilterer creates a new log filterer instance of Example, bound to a specific deployed contract.
func NewExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*ExampleFilterer, error) {
	contract, err := bindExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExampleFilterer{contract: contract}, nil
}

// bindExample binds a generic wrapper to an already deployed contract.
func bindExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Example *ExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Example.Contract.ExampleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Example *ExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.Contract.ExampleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Example *ExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Example.Contract.ExampleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Example *ExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Example.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Example *ExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Example *ExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Example.Contract.contract.Transact(opts, method, params...)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_Example *ExampleCaller) Bar(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Example.contract.Call(opts, &out, "bar")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_Example *ExampleSession) Bar() (*big.Int, error) {
	return _Example.Contract.Bar(&_Example.CallOpts)
}

// Bar is a free data retrieval call binding the contract method 0xfebb0f7e.
//
// Solidity: function bar() view returns(uint256)
func (_Example *ExampleCallerSession) Bar() (*big.Int, error) {
	return _Example.Contract.Bar(&_Example.CallOpts)
}

// LastMessage is a free data retrieval call binding the contract method 0x32970710.
//
// Solidity: function lastMessage() view returns(bytes)
func (_Example *ExampleCaller) LastMessage(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Example.contract.Call(opts, &out, "lastMessage")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// LastMessage is a free data retrieval call binding the contract method 0x32970710.
//
// Solidity: function lastMessage() view returns(bytes)
func (_Example *ExampleSession) LastMessage() ([]byte, error) {
	return _Example.Contract.LastMessage(&_Example.CallOpts)
}

// LastMessage is a free data retrieval call binding the contract method 0x32970710.
//
// Solidity: function lastMessage() view returns(bytes)
func (_Example *ExampleCallerSession) LastMessage() ([]byte, error) {
	return _Example.Contract.LastMessage(&_Example.CallOpts)
}

// LastSender is a free data retrieval call binding the contract method 0x256fec88.
//
// Solidity: function lastSender() view returns(bytes)
func (_Example *ExampleCaller) LastSender(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Example.contract.Call(opts, &out, "lastSender")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// LastSender is a free data retrieval call binding the contract method 0x256fec88.
//
// Solidity: function lastSender() view returns(bytes)
func (_Example *ExampleSession) LastSender() ([]byte, error) {
	return _Example.Contract.LastSender(&_Example.CallOpts)
}

// LastSender is a free data retrieval call binding the contract method 0x256fec88.
//
// Solidity: function lastSender() view returns(bytes)
func (_Example *ExampleCallerSession) LastSender() ([]byte, error) {
	return _Example.Contract.LastSender(&_Example.CallOpts)
}

// DoRevert is a paid mutator transaction binding the contract method 0xafc874d2.
//
// Solidity: function doRevert() returns()
func (_Example *ExampleTransactor) DoRevert(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "doRevert")
}

// DoRevert is a paid mutator transaction binding the contract method 0xafc874d2.
//
// Solidity: function doRevert() returns()
func (_Example *ExampleSession) DoRevert() (*types.Transaction, error) {
	return _Example.Contract.DoRevert(&_Example.TransactOpts)
}

// DoRevert is a paid mutator transaction binding the contract method 0xafc874d2.
//
// Solidity: function doRevert() returns()
func (_Example *ExampleTransactorSession) DoRevert() (*types.Transaction, error) {
	return _Example.Contract.DoRevert(&_Example.TransactOpts)
}

// DoRevertWithMessage is a paid mutator transaction binding the contract method 0xd720cb45.
//
// Solidity: function doRevertWithMessage() returns()
func (_Example *ExampleTransactor) DoRevertWithMessage(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "doRevertWithMessage")
}

// DoRevertWithMessage is a paid mutator transaction binding the contract method 0xd720cb45.
//
// Solidity: function doRevertWithMessage() returns()
func (_Example *ExampleSession) DoRevertWithMessage() (*types.Transaction, error) {
	return _Example.Contract.DoRevertWithMessage(&_Example.TransactOpts)
}

// DoRevertWithMessage is a paid mutator transaction binding the contract method 0xd720cb45.
//
// Solidity: function doRevertWithMessage() returns()
func (_Example *ExampleTransactorSession) DoRevertWithMessage() (*types.Transaction, error) {
	return _Example.Contract.DoRevertWithMessage(&_Example.TransactOpts)
}

// DoRevertWithRequire is a paid mutator transaction binding the contract method 0xdd8e556c.
//
// Solidity: function doRevertWithRequire() returns()
func (_Example *ExampleTransactor) DoRevertWithRequire(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "doRevertWithRequire")
}

// DoRevertWithRequire is a paid mutator transaction binding the contract method 0xdd8e556c.
//
// Solidity: function doRevertWithRequire() returns()
func (_Example *ExampleSession) DoRevertWithRequire() (*types.Transaction, error) {
	return _Example.Contract.DoRevertWithRequire(&_Example.TransactOpts)
}

// DoRevertWithRequire is a paid mutator transaction binding the contract method 0xdd8e556c.
//
// Solidity: function doRevertWithRequire() returns()
func (_Example *ExampleTransactorSession) DoRevertWithRequire() (*types.Transaction, error) {
	return _Example.Contract.DoRevertWithRequire(&_Example.TransactOpts)
}

// DoSucceed is a paid mutator transaction binding the contract method 0xfd5ad965.
//
// Solidity: function doSucceed() returns()
func (_Example *ExampleTransactor) DoSucceed(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "doSucceed")
}

// DoSucceed is a paid mutator transaction binding the contract method 0xfd5ad965.
//
// Solidity: function doSucceed() returns()
func (_Example *ExampleSession) DoSucceed() (*types.Transaction, error) {
	return _Example.Contract.DoSucceed(&_Example.TransactOpts)
}

// DoSucceed is a paid mutator transaction binding the contract method 0xfd5ad965.
//
// Solidity: function doSucceed() returns()
func (_Example *ExampleTransactorSession) DoSucceed() (*types.Transaction, error) {
	return _Example.Contract.DoSucceed(&_Example.TransactOpts)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_Example *ExampleTransactor) OnCall(opts *bind.TransactOpts, context ExamplezContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "onCall", context, zrc20, amount, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_Example *ExampleSession) OnCall(context ExamplezContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _Example.Contract.OnCall(&_Example.TransactOpts, context, zrc20, amount, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_Example *ExampleTransactorSession) OnCall(context ExamplezContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _Example.Contract.OnCall(&_Example.TransactOpts, context, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_Example *ExampleTransactor) OnCrossChainCall(opts *bind.TransactOpts, context ExamplezContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _Example.contract.Transact(opts, "onCrossChainCall", context, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_Example *ExampleSession) OnCrossChainCall(context ExamplezContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _Example.Contract.OnCrossChainCall(&_Example.TransactOpts, context, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_Example *ExampleTransactorSession) OnCrossChainCall(context ExamplezContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _Example.Contract.OnCrossChainCall(&_Example.TransactOpts, context, zrc20, amount, message)
}
