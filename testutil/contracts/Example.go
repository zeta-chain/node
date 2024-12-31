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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Foo\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"bar\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevertWithMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doRevertWithRequire\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"doSucceed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMessage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structExample.zContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b50600080819055506108e0806100266000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063dd8e556c1161005b578063dd8e556c146100b4578063de43156e146100be578063fd5ad965146100da578063febb0f7e146100e45761007d565b80633297071014610082578063afc874d2146100a0578063d720cb45146100aa575b600080fd5b61008a610102565b6040516100979190610300565b60405180910390f35b6100a8610190565b005b6100b26101c2565b005b6100bc6101fd565b005b6100d860048036038101906100d39190610449565b610240565b005b6100e2610260565b005b6100ec61026a565b6040516100f991906104fc565b60405180910390f35b6001805461010f90610546565b80601f016020809104026020016040519081016040528092919081815260200182805461013b90610546565b80156101885780601f1061015d57610100808354040283529160200191610188565b820191906000526020600020905b81548152906001019060200180831161016b57829003601f168201915b505050505081565b6040517fbfb4ebcf00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101f4906105d4565b60405180910390fd5b600061023e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610235906105d4565b60405180910390fd5b565b826000819055508181600191826102589291906107da565b505050505050565b6001600081905550565b60005481565b600081519050919050565b600082825260208201905092915050565b60005b838110156102aa57808201518184015260208101905061028f565b60008484015250505050565b6000601f19601f8301169050919050565b60006102d282610270565b6102dc818561027b565b93506102ec81856020860161028c565b6102f5816102b6565b840191505092915050565b6000602082019050818103600083015261031a81846102c7565b905092915050565b600080fd5b600080fd5b600080fd5b6000606082840312156103475761034661032c565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061037b82610350565b9050919050565b61038b81610370565b811461039657600080fd5b50565b6000813590506103a881610382565b92915050565b6000819050919050565b6103c1816103ae565b81146103cc57600080fd5b50565b6000813590506103de816103b8565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610409576104086103e4565b5b8235905067ffffffffffffffff811115610426576104256103e9565b5b602083019150836001820283011115610442576104416103ee565b5b9250929050565b60008060008060006080868803121561046557610464610322565b5b600086013567ffffffffffffffff81111561048357610482610327565b5b61048f88828901610331565b95505060206104a088828901610399565b94505060406104b1888289016103cf565b935050606086013567ffffffffffffffff8111156104d2576104d1610327565b5b6104de888289016103f3565b92509250509295509295909350565b6104f6816103ae565b82525050565b600060208201905061051160008301846104ed565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061055e57607f821691505b60208210810361057157610570610517565b5b50919050565b600082825260208201905092915050565b7f666f6f0000000000000000000000000000000000000000000000000000000000600082015250565b60006105be600383610577565b91506105c982610588565b602082019050919050565b600060208201905081810360008301526105ed816105b1565b9050919050565b600082905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026106907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610653565b61069a8683610653565b95508019841693508086168417925050509392505050565b6000819050919050565b60006106d76106d26106cd846103ae565b6106b2565b6103ae565b9050919050565b6000819050919050565b6106f1836106bc565b6107056106fd826106de565b848454610660565b825550505050565b600090565b61071a61070d565b6107258184846106e8565b505050565b5b818110156107495761073e600082610712565b60018101905061072b565b5050565b601f82111561078e5761075f8161062e565b61076884610643565b81016020851015610777578190505b61078b61078385610643565b83018261072a565b50505b505050565b600082821c905092915050565b60006107b160001984600802610793565b1980831691505092915050565b60006107ca83836107a0565b9150826002028217905092915050565b6107e483836105f4565b67ffffffffffffffff8111156107fd576107fc6105ff565b5b6108078254610546565b61081282828561074d565b6000601f831160018114610841576000841561082f578287013590505b61083985826107be565b8655506108a1565b601f19841661084f8661062e565b60005b8281101561087757848901358255600182019150602085019450602081019050610852565b868310156108945784890135610890601f8916826107a0565b8355505b6001600288020188555050505b5050505050505056fea26469706673582212202506bee512af2f3acd98556baafd9afc04b2d8c95f09bbb5e387fce570c9a32764736f6c634300081a0033",
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
