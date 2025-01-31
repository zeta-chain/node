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

// ExamplezContext is an auto generated low-level Go binding around an user-defined struct.
type ExamplezContext struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}

// ExampleMetaData contains all meta data concerning the Example contract.
var ExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Foo\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"bar\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevertWithMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevertWithRequire\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doSucceed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMessage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structExample.zContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structExample.zContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506000808190555061091b806100266000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063dd8e556c1161005b578063dd8e556c146100db578063de43156e146100e5578063fd5ad96514610101578063febb0f7e1461010b57610088565b8063329707101461008d5780635bcfd616146100ab578063afc874d2146100c7578063d720cb45146100d1575b600080fd5b610095610129565b6040516100a2919061033b565b60405180910390f35b6100c560048036038101906100c09190610484565b6101b7565b005b6100cf6101d7565b005b6100d9610209565b005b6100e3610244565b005b6100ff60048036038101906100fa9190610484565b610287565b005b61010961029b565b005b6101136102a5565b6040516101209190610537565b60405180910390f35b6001805461013690610581565b80601f016020809104026020016040519081016040528092919081815260200182805461016290610581565b80156101af5780601f10610184576101008083540402835291602001916101af565b820191906000526020600020905b81548152906001019060200180831161019257829003601f168201915b505050505081565b826000819055508181600191826101cf929190610798565b505050505050565b6040517fbfb4ebcf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161023b906108c5565b60405180910390fd5b6000610285576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161027c906108c5565b60405180910390fd5b565b61029485858585856101b7565b5050505050565b6001600081905550565b60005481565b600081519050919050565b600082825260208201905092915050565b60005b838110156102e55780820151818401526020810190506102ca565b60008484015250505050565b6000601f19601f8301169050919050565b600061030d826102ab565b61031781856102b6565b93506103278185602086016102c7565b610330816102f1565b840191505092915050565b600060208201905081810360008301526103558184610302565b905092915050565b600080fd5b600080fd5b600080fd5b60006060828403121561038257610381610367565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006103b68261038b565b9050919050565b6103c6816103ab565b81146103d157600080fd5b50565b6000813590506103e3816103bd565b92915050565b6000819050919050565b6103fc816103e9565b811461040757600080fd5b50565b600081359050610419816103f3565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126104445761044361041f565b5b8235905067ffffffffffffffff81111561046157610460610424565b5b60208301915083600182028301111561047d5761047c610429565b5b9250929050565b6000806000806000608086880312156104a05761049f61035d565b5b600086013567ffffffffffffffff8111156104be576104bd610362565b5b6104ca8882890161036c565b95505060206104db888289016103d4565b94505060406104ec8882890161040a565b935050606086013567ffffffffffffffff81111561050d5761050c610362565b5b6105198882890161042e565b92509250509295509295909350565b610531816103e9565b82525050565b600060208201905061054c6000830184610528565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061059957607f821691505b6020821081036105ac576105ab610552565b5b50919050565b600082905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b60006008830261064e7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610611565b6106588683610611565b95508019841693508086168417925050509392505050565b6000819050919050565b600061069561069061068b846103e9565b610670565b6103e9565b9050919050565b6000819050919050565b6106af8361067a565b6106c36106bb8261069c565b84845461061e565b825550505050565b600090565b6106d86106cb565b6106e38184846106a6565b505050565b5b81811015610707576106fc6000826106d0565b6001810190506106e9565b5050565b601f82111561074c5761071d816105ec565b61072684610601565b81016020851015610735578190505b61074961074185610601565b8301826106e8565b50505b505050565b600082821c905092915050565b600061076f60001984600802610751565b1980831691505092915050565b6000610788838361075e565b9150826002028217905092915050565b6107a283836105b2565b67ffffffffffffffff8111156107bb576107ba6105bd565b5b6107c58254610581565b6107d082828561070b565b6000601f8311600181146107ff57600084156107ed578287013590505b6107f7858261077c565b86555061085f565b601f19841661080d866105ec565b60005b8281101561083557848901358255600182019150602085019450602081019050610810565b86831015610852578489013561084e601f89168261075e565b8355505b6001600288020188555050505b50505050505050565b600082825260208201905092915050565b7f666f6f0000000000000000000000000000000000000000000000000000000000600082015250565b60006108af600383610868565b91506108ba82610879565b602082019050919050565b600060208201905081810360008301526108de816108a2565b905091905056fea26469706673582212209cfc763a2051af7ae03bd5c8a8765e80a09469a9bb8875cda4c3917306b2238264736f6c634300081a0033",
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
