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
	Bin: "0x60a060405260676000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561005157600080fd5b503373ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250506080516106a16100ae6000396000818160fc015281816101fc01526102fc01526106a16000f3fe6080604052600436106100385760003560e01c806347e7ef2414610041578063f3fef3a31461007e578063f7888aec146100bb5761003f565b3661003f57005b005b34801561004d57600080fd5b506100686004803603810190610063919061048f565b6100f8565b60405161007591906104ea565b60405180910390f35b34801561008a57600080fd5b506100a560048036038101906100a0919061048f565b6101f8565b6040516100b291906104ea565b60405180910390f35b3480156100c757600080fd5b506100e260048036038101906100dd9190610505565b6102f8565b6040516100ef9190610554565b60405180910390f35b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461015257600080fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166347e7ef2484846040518363ffffffff1660e01b81526004016101ad92919061057e565b6020604051808303816000875af11580156101cc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101f091906105d3565b905092915050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461025257600080fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f3fef3a384846040518363ffffffff1660e01b81526004016102ad92919061057e565b6020604051808303816000875af11580156102cc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102f091906105d3565b905092915050565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461035257600080fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f7888aec84846040518363ffffffff1660e01b81526004016103ad929190610600565b602060405180830381865afa1580156103ca573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103ee919061063e565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610426826103fb565b9050919050565b6104368161041b565b811461044157600080fd5b50565b6000813590506104538161042d565b92915050565b6000819050919050565b61046c81610459565b811461047757600080fd5b50565b60008135905061048981610463565b92915050565b600080604083850312156104a6576104a56103f6565b5b60006104b485828601610444565b92505060206104c58582860161047a565b9150509250929050565b60008115159050919050565b6104e4816104cf565b82525050565b60006020820190506104ff60008301846104db565b92915050565b6000806040838503121561051c5761051b6103f6565b5b600061052a85828601610444565b925050602061053b85828601610444565b9150509250929050565b61054e81610459565b82525050565b60006020820190506105696000830184610545565b92915050565b6105788161041b565b82525050565b6000604082019050610593600083018561056f565b6105a06020830184610545565b9392505050565b6105b0816104cf565b81146105bb57600080fd5b50565b6000815190506105cd816105a7565b92915050565b6000602082840312156105e9576105e86103f6565b5b60006105f7848285016105be565b91505092915050565b6000604082019050610615600083018561056f565b610622602083018461056f565b9392505050565b60008151905061063881610463565b92915050565b600060208284031215610654576106536103f6565b5b600061066284828501610629565b9150509291505056fea2646970667358221220c0c585590967d576c91dec72feb553b8a6186d045ddc558ecc1de043c068bcb764736f6c634300080a0033",
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
