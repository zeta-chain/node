// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testbank

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

// TestBankMetaData contains all meta data concerning the TestBank contract.
var TestBankMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60a060405260675f806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550348015604e575f80fd5b503373ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250506080516106826100a85f395f818160f7015281816101f701526102f701526106825ff3fe608060405260043610610037575f3560e01c806347e7ef2414610040578063f3fef3a31461007c578063f7888aec146100b85761003e565b3661003e57005b005b34801561004b575f80fd5b5061006660048036038101906100619190610484565b6100f4565b60405161007391906104dc565b60405180910390f35b348015610087575f80fd5b506100a2600480360381019061009d9190610484565b6101f4565b6040516100af91906104dc565b60405180910390f35b3480156100c3575f80fd5b506100de60048036038101906100d991906104f5565b6102f4565b6040516100eb9190610542565b60405180910390f35b5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461014c575f80fd5b5f805f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166347e7ef2485856040518363ffffffff1660e01b81526004016101a892919061056a565b6020604051808303815f875af11580156101c4573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101e891906105bb565b90508091505092915050565b5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461024c575f80fd5b5f805f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f3fef3a385856040518363ffffffff1660e01b81526004016102a892919061056a565b6020604051808303815f875af11580156102c4573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906102e891906105bb565b90508091505092915050565b5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461034c575f80fd5b5f805f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f7888aec85856040518363ffffffff1660e01b81526004016103a89291906105e6565b602060405180830381865afa1580156103c3573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103e79190610621565b90508091505092915050565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610420826103f7565b9050919050565b61043081610416565b811461043a575f80fd5b50565b5f8135905061044b81610427565b92915050565b5f819050919050565b61046381610451565b811461046d575f80fd5b50565b5f8135905061047e8161045a565b92915050565b5f806040838503121561049a576104996103f3565b5b5f6104a78582860161043d565b92505060206104b885828601610470565b9150509250929050565b5f8115159050919050565b6104d6816104c2565b82525050565b5f6020820190506104ef5f8301846104cd565b92915050565b5f806040838503121561050b5761050a6103f3565b5b5f6105188582860161043d565b92505060206105298582860161043d565b9150509250929050565b61053c81610451565b82525050565b5f6020820190506105555f830184610533565b92915050565b61056481610416565b82525050565b5f60408201905061057d5f83018561055b565b61058a6020830184610533565b9392505050565b61059a816104c2565b81146105a4575f80fd5b50565b5f815190506105b581610591565b92915050565b5f602082840312156105d0576105cf6103f3565b5b5f6105dd848285016105a7565b91505092915050565b5f6040820190506105f95f83018561055b565b610606602083018461055b565b9392505050565b5f8151905061061b8161045a565b92915050565b5f60208284031215610636576106356103f3565b5b5f6106438482850161060d565b9150509291505056fea264697066735822122058606fa72b81a1490986349908a18577e434b840b15c1299c0a4ae455e29088864736f6c634300081a0033",
}

// TestBankABI is the input ABI used to generate the binding from.
// Deprecated: Use TestBankMetaData.ABI instead.
var TestBankABI = TestBankMetaData.ABI

// TestBankBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestBankMetaData.Bin instead.
var TestBankBin = TestBankMetaData.Bin

// DeployTestBank deploys a new Ethereum contract, binding an instance of TestBank to it.
func DeployTestBank(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestBank, error) {
	parsed, err := TestBankMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestBankBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestBank{TestBankCaller: TestBankCaller{contract: contract}, TestBankTransactor: TestBankTransactor{contract: contract}, TestBankFilterer: TestBankFilterer{contract: contract}}, nil
}

// TestBank is an auto generated Go binding around an Ethereum contract.
type TestBank struct {
	TestBankCaller     // Read-only binding to the contract
	TestBankTransactor // Write-only binding to the contract
	TestBankFilterer   // Log filterer for contract events
}

// TestBankCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestBankCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestBankTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestBankTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestBankFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestBankFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestBankSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestBankSession struct {
	Contract     *TestBank         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestBankCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestBankCallerSession struct {
	Contract *TestBankCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TestBankTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestBankTransactorSession struct {
	Contract     *TestBankTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TestBankRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestBankRaw struct {
	Contract *TestBank // Generic contract binding to access the raw methods on
}

// TestBankCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestBankCallerRaw struct {
	Contract *TestBankCaller // Generic read-only contract binding to access the raw methods on
}

// TestBankTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestBankTransactorRaw struct {
	Contract *TestBankTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestBank creates a new instance of TestBank, bound to a specific deployed contract.
func NewTestBank(address common.Address, backend bind.ContractBackend) (*TestBank, error) {
	contract, err := bindTestBank(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestBank{TestBankCaller: TestBankCaller{contract: contract}, TestBankTransactor: TestBankTransactor{contract: contract}, TestBankFilterer: TestBankFilterer{contract: contract}}, nil
}

// NewTestBankCaller creates a new read-only instance of TestBank, bound to a specific deployed contract.
func NewTestBankCaller(address common.Address, caller bind.ContractCaller) (*TestBankCaller, error) {
	contract, err := bindTestBank(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestBankCaller{contract: contract}, nil
}

// NewTestBankTransactor creates a new write-only instance of TestBank, bound to a specific deployed contract.
func NewTestBankTransactor(address common.Address, transactor bind.ContractTransactor) (*TestBankTransactor, error) {
	contract, err := bindTestBank(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestBankTransactor{contract: contract}, nil
}

// NewTestBankFilterer creates a new log filterer instance of TestBank, bound to a specific deployed contract.
func NewTestBankFilterer(address common.Address, filterer bind.ContractFilterer) (*TestBankFilterer, error) {
	contract, err := bindTestBank(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestBankFilterer{contract: contract}, nil
}

// bindTestBank binds a generic wrapper to an already deployed contract.
func bindTestBank(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestBankMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestBank *TestBankRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestBank.Contract.TestBankCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestBank *TestBankRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestBank.Contract.TestBankTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestBank *TestBankRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestBank.Contract.TestBankTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestBank *TestBankCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestBank.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestBank *TestBankTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestBank.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestBank *TestBankTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestBank.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0xf7888aec.
//
// Solidity: function balanceOf(address zrc20, address user) view returns(uint256)
func (_TestBank *TestBankCaller) BalanceOf(opts *bind.CallOpts, zrc20 common.Address, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestBank.contract.Call(opts, &out, "balanceOf", zrc20, user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0xf7888aec.
//
// Solidity: function balanceOf(address zrc20, address user) view returns(uint256)
func (_TestBank *TestBankSession) BalanceOf(zrc20 common.Address, user common.Address) (*big.Int, error) {
	return _TestBank.Contract.BalanceOf(&_TestBank.CallOpts, zrc20, user)
}

// BalanceOf is a free data retrieval call binding the contract method 0xf7888aec.
//
// Solidity: function balanceOf(address zrc20, address user) view returns(uint256)
func (_TestBank *TestBankCallerSession) BalanceOf(zrc20 common.Address, user common.Address) (*big.Int, error) {
	return _TestBank.Contract.BalanceOf(&_TestBank.CallOpts, zrc20, user)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address zrc20, uint256 amount) returns(bool)
func (_TestBank *TestBankTransactor) Deposit(opts *bind.TransactOpts, zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestBank.contract.Transact(opts, "deposit", zrc20, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address zrc20, uint256 amount) returns(bool)
func (_TestBank *TestBankSession) Deposit(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestBank.Contract.Deposit(&_TestBank.TransactOpts, zrc20, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address zrc20, uint256 amount) returns(bool)
func (_TestBank *TestBankTransactorSession) Deposit(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestBank.Contract.Deposit(&_TestBank.TransactOpts, zrc20, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address zrc20, uint256 amount) returns(bool)
func (_TestBank *TestBankTransactor) Withdraw(opts *bind.TransactOpts, zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestBank.contract.Transact(opts, "withdraw", zrc20, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address zrc20, uint256 amount) returns(bool)
func (_TestBank *TestBankSession) Withdraw(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestBank.Contract.Withdraw(&_TestBank.TransactOpts, zrc20, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address zrc20, uint256 amount) returns(bool)
func (_TestBank *TestBankTransactorSession) Withdraw(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestBank.Contract.Withdraw(&_TestBank.TransactOpts, zrc20, amount)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestBank *TestBankTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TestBank.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestBank *TestBankSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TestBank.Contract.Fallback(&_TestBank.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestBank *TestBankTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TestBank.Contract.Fallback(&_TestBank.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestBank *TestBankTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestBank.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestBank *TestBankSession) Receive() (*types.Transaction, error) {
	return _TestBank.Contract.Receive(&_TestBank.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestBank *TestBankTransactorSession) Receive() (*types.Transaction, error) {
	return _TestBank.Contract.Receive(&_TestBank.TransactOpts)
}
