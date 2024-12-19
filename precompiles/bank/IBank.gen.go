// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bank

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

// IBankMetaData contains all meta data concerning the IBank contract.
var IBankMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zrc20_depositor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zrc20_token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"cosmos_token\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cosmos_address\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zrc20_withdrawer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zrc20_token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"cosmos_token\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cosmos_address\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IBankABI is the input ABI used to generate the binding from.
// Deprecated: Use IBankMetaData.ABI instead.
var IBankABI = IBankMetaData.ABI

// IBank is an auto generated Go binding around an Ethereum contract.
type IBank struct {
	IBankCaller     // Read-only binding to the contract
	IBankTransactor // Write-only binding to the contract
	IBankFilterer   // Log filterer for contract events
}

// IBankCaller is an auto generated read-only Go binding around an Ethereum contract.
type IBankCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBankTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IBankTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBankFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IBankFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBankSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IBankSession struct {
	Contract     *IBank            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBankCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IBankCallerSession struct {
	Contract *IBankCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IBankTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IBankTransactorSession struct {
	Contract     *IBankTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBankRaw is an auto generated low-level Go binding around an Ethereum contract.
type IBankRaw struct {
	Contract *IBank // Generic contract binding to access the raw methods on
}

// IBankCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IBankCallerRaw struct {
	Contract *IBankCaller // Generic read-only contract binding to access the raw methods on
}

// IBankTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IBankTransactorRaw struct {
	Contract *IBankTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBank creates a new instance of IBank, bound to a specific deployed contract.
func NewIBank(address common.Address, backend bind.ContractBackend) (*IBank, error) {
	contract, err := bindIBank(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBank{IBankCaller: IBankCaller{contract: contract}, IBankTransactor: IBankTransactor{contract: contract}, IBankFilterer: IBankFilterer{contract: contract}}, nil
}

// NewIBankCaller creates a new read-only instance of IBank, bound to a specific deployed contract.
func NewIBankCaller(address common.Address, caller bind.ContractCaller) (*IBankCaller, error) {
	contract, err := bindIBank(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBankCaller{contract: contract}, nil
}

// NewIBankTransactor creates a new write-only instance of IBank, bound to a specific deployed contract.
func NewIBankTransactor(address common.Address, transactor bind.ContractTransactor) (*IBankTransactor, error) {
	contract, err := bindIBank(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBankTransactor{contract: contract}, nil
}

// NewIBankFilterer creates a new log filterer instance of IBank, bound to a specific deployed contract.
func NewIBankFilterer(address common.Address, filterer bind.ContractFilterer) (*IBankFilterer, error) {
	contract, err := bindIBank(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBankFilterer{contract: contract}, nil
}

// bindIBank binds a generic wrapper to an already deployed contract.
func bindIBank(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBankMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBank *IBankRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBank.Contract.IBankCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBank *IBankRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBank.Contract.IBankTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBank *IBankRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBank.Contract.IBankTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBank *IBankCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBank.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBank *IBankTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBank.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBank *IBankTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBank.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0xf7888aec.
//
// Solidity: function balanceOf(address zrc20, address user) view returns(uint256 balance)
func (_IBank *IBankCaller) BalanceOf(opts *bind.CallOpts, zrc20 common.Address, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IBank.contract.Call(opts, &out, "balanceOf", zrc20, user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0xf7888aec.
//
// Solidity: function balanceOf(address zrc20, address user) view returns(uint256 balance)
func (_IBank *IBankSession) BalanceOf(zrc20 common.Address, user common.Address) (*big.Int, error) {
	return _IBank.Contract.BalanceOf(&_IBank.CallOpts, zrc20, user)
}

// BalanceOf is a free data retrieval call binding the contract method 0xf7888aec.
//
// Solidity: function balanceOf(address zrc20, address user) view returns(uint256 balance)
func (_IBank *IBankCallerSession) BalanceOf(zrc20 common.Address, user common.Address) (*big.Int, error) {
	return _IBank.Contract.BalanceOf(&_IBank.CallOpts, zrc20, user)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address zrc20, uint256 amount) returns(bool success)
func (_IBank *IBankTransactor) Deposit(opts *bind.TransactOpts, zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.contract.Transact(opts, "deposit", zrc20, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address zrc20, uint256 amount) returns(bool success)
func (_IBank *IBankSession) Deposit(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.Deposit(&_IBank.TransactOpts, zrc20, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address zrc20, uint256 amount) returns(bool success)
func (_IBank *IBankTransactorSession) Deposit(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.Deposit(&_IBank.TransactOpts, zrc20, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address zrc20, uint256 amount) returns(bool success)
func (_IBank *IBankTransactor) Withdraw(opts *bind.TransactOpts, zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.contract.Transact(opts, "withdraw", zrc20, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address zrc20, uint256 amount) returns(bool success)
func (_IBank *IBankSession) Withdraw(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.Withdraw(&_IBank.TransactOpts, zrc20, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address zrc20, uint256 amount) returns(bool success)
func (_IBank *IBankTransactorSession) Withdraw(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IBank.Contract.Withdraw(&_IBank.TransactOpts, zrc20, amount)
}

// IBankDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the IBank contract.
type IBankDepositIterator struct {
	Event *IBankDeposit // Event containing the contract specifics and raw log

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
func (it *IBankDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBankDeposit)
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
		it.Event = new(IBankDeposit)
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
func (it *IBankDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBankDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBankDeposit represents a Deposit event raised by the IBank contract.
type IBankDeposit struct {
	Zrc20Depositor common.Address
	Zrc20Token     common.Address
	CosmosToken    common.Hash
	CosmosAddress  string
	Amount         *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xbd7d4de0b30a306221956a420cad57737ae9c1ee63072c96a4f1ab81e6eea264.
//
// Solidity: event Deposit(address indexed zrc20_depositor, address indexed zrc20_token, string indexed cosmos_token, string cosmos_address, uint256 amount)
func (_IBank *IBankFilterer) FilterDeposit(opts *bind.FilterOpts, zrc20_depositor []common.Address, zrc20_token []common.Address, cosmos_token []string) (*IBankDepositIterator, error) {

	var zrc20_depositorRule []interface{}
	for _, zrc20_depositorItem := range zrc20_depositor {
		zrc20_depositorRule = append(zrc20_depositorRule, zrc20_depositorItem)
	}
	var zrc20_tokenRule []interface{}
	for _, zrc20_tokenItem := range zrc20_token {
		zrc20_tokenRule = append(zrc20_tokenRule, zrc20_tokenItem)
	}
	var cosmos_tokenRule []interface{}
	for _, cosmos_tokenItem := range cosmos_token {
		cosmos_tokenRule = append(cosmos_tokenRule, cosmos_tokenItem)
	}

	logs, sub, err := _IBank.contract.FilterLogs(opts, "Deposit", zrc20_depositorRule, zrc20_tokenRule, cosmos_tokenRule)
	if err != nil {
		return nil, err
	}
	return &IBankDepositIterator{contract: _IBank.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xbd7d4de0b30a306221956a420cad57737ae9c1ee63072c96a4f1ab81e6eea264.
//
// Solidity: event Deposit(address indexed zrc20_depositor, address indexed zrc20_token, string indexed cosmos_token, string cosmos_address, uint256 amount)
func (_IBank *IBankFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *IBankDeposit, zrc20_depositor []common.Address, zrc20_token []common.Address, cosmos_token []string) (event.Subscription, error) {

	var zrc20_depositorRule []interface{}
	for _, zrc20_depositorItem := range zrc20_depositor {
		zrc20_depositorRule = append(zrc20_depositorRule, zrc20_depositorItem)
	}
	var zrc20_tokenRule []interface{}
	for _, zrc20_tokenItem := range zrc20_token {
		zrc20_tokenRule = append(zrc20_tokenRule, zrc20_tokenItem)
	}
	var cosmos_tokenRule []interface{}
	for _, cosmos_tokenItem := range cosmos_token {
		cosmos_tokenRule = append(cosmos_tokenRule, cosmos_tokenItem)
	}

	logs, sub, err := _IBank.contract.WatchLogs(opts, "Deposit", zrc20_depositorRule, zrc20_tokenRule, cosmos_tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBankDeposit)
				if err := _IBank.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xbd7d4de0b30a306221956a420cad57737ae9c1ee63072c96a4f1ab81e6eea264.
//
// Solidity: event Deposit(address indexed zrc20_depositor, address indexed zrc20_token, string indexed cosmos_token, string cosmos_address, uint256 amount)
func (_IBank *IBankFilterer) ParseDeposit(log types.Log) (*IBankDeposit, error) {
	event := new(IBankDeposit)
	if err := _IBank.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IBankWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the IBank contract.
type IBankWithdrawIterator struct {
	Event *IBankWithdraw // Event containing the contract specifics and raw log

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
func (it *IBankWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBankWithdraw)
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
		it.Event = new(IBankWithdraw)
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
func (it *IBankWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBankWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBankWithdraw represents a Withdraw event raised by the IBank contract.
type IBankWithdraw struct {
	Zrc20Withdrawer common.Address
	Zrc20Token      common.Address
	CosmosToken     common.Hash
	CosmosAddress   string
	Amount          *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x1ad70707c91d850319aeab00514a0166569359f0b8dc5285bdd6e6b9c464b18e.
//
// Solidity: event Withdraw(address indexed zrc20_withdrawer, address indexed zrc20_token, string indexed cosmos_token, string cosmos_address, uint256 amount)
func (_IBank *IBankFilterer) FilterWithdraw(opts *bind.FilterOpts, zrc20_withdrawer []common.Address, zrc20_token []common.Address, cosmos_token []string) (*IBankWithdrawIterator, error) {

	var zrc20_withdrawerRule []interface{}
	for _, zrc20_withdrawerItem := range zrc20_withdrawer {
		zrc20_withdrawerRule = append(zrc20_withdrawerRule, zrc20_withdrawerItem)
	}
	var zrc20_tokenRule []interface{}
	for _, zrc20_tokenItem := range zrc20_token {
		zrc20_tokenRule = append(zrc20_tokenRule, zrc20_tokenItem)
	}
	var cosmos_tokenRule []interface{}
	for _, cosmos_tokenItem := range cosmos_token {
		cosmos_tokenRule = append(cosmos_tokenRule, cosmos_tokenItem)
	}

	logs, sub, err := _IBank.contract.FilterLogs(opts, "Withdraw", zrc20_withdrawerRule, zrc20_tokenRule, cosmos_tokenRule)
	if err != nil {
		return nil, err
	}
	return &IBankWithdrawIterator{contract: _IBank.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x1ad70707c91d850319aeab00514a0166569359f0b8dc5285bdd6e6b9c464b18e.
//
// Solidity: event Withdraw(address indexed zrc20_withdrawer, address indexed zrc20_token, string indexed cosmos_token, string cosmos_address, uint256 amount)
func (_IBank *IBankFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *IBankWithdraw, zrc20_withdrawer []common.Address, zrc20_token []common.Address, cosmos_token []string) (event.Subscription, error) {

	var zrc20_withdrawerRule []interface{}
	for _, zrc20_withdrawerItem := range zrc20_withdrawer {
		zrc20_withdrawerRule = append(zrc20_withdrawerRule, zrc20_withdrawerItem)
	}
	var zrc20_tokenRule []interface{}
	for _, zrc20_tokenItem := range zrc20_token {
		zrc20_tokenRule = append(zrc20_tokenRule, zrc20_tokenItem)
	}
	var cosmos_tokenRule []interface{}
	for _, cosmos_tokenItem := range cosmos_token {
		cosmos_tokenRule = append(cosmos_tokenRule, cosmos_tokenItem)
	}

	logs, sub, err := _IBank.contract.WatchLogs(opts, "Withdraw", zrc20_withdrawerRule, zrc20_tokenRule, cosmos_tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBankWithdraw)
				if err := _IBank.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x1ad70707c91d850319aeab00514a0166569359f0b8dc5285bdd6e6b9c464b18e.
//
// Solidity: event Withdraw(address indexed zrc20_withdrawer, address indexed zrc20_token, string indexed cosmos_token, string cosmos_address, uint256 amount)
func (_IBank *IBankFilterer) ParseWithdraw(log types.Log) (*IBankWithdraw, error) {
	event := new(IBankWithdraw)
	if err := _IBank.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
