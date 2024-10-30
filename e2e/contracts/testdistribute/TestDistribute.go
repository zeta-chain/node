// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testdistribute

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

// TestDistributeMetaData contains all meta data concerning the TestDistribute contract.
var TestDistributeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zrc20_distributor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zrc20_token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Distributed\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"distributeThroughContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60a060405260666000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561005157600080fd5b503373ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505060805161034d6100a06000396000606c015261034d6000f3fe6080604052600436106100225760003560e01c806350b54e841461002b57610029565b3661002957005b005b34801561003757600080fd5b50610052600480360381019061004d9190610201565b610068565b60405161005f919061025c565b60405180910390f35b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146100c257600080fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fb93210884846040518363ffffffff1660e01b815260040161011d929190610295565b6020604051808303816000875af115801561013c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061016091906102ea565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006101988261016d565b9050919050565b6101a88161018d565b81146101b357600080fd5b50565b6000813590506101c58161019f565b92915050565b6000819050919050565b6101de816101cb565b81146101e957600080fd5b50565b6000813590506101fb816101d5565b92915050565b6000806040838503121561021857610217610168565b5b6000610226858286016101b6565b9250506020610237858286016101ec565b9150509250929050565b60008115159050919050565b61025681610241565b82525050565b6000602082019050610271600083018461024d565b92915050565b6102808161018d565b82525050565b61028f816101cb565b82525050565b60006040820190506102aa6000830185610277565b6102b76020830184610286565b9392505050565b6102c781610241565b81146102d257600080fd5b50565b6000815190506102e4816102be565b92915050565b600060208284031215610300576102ff610168565b5b600061030e848285016102d5565b9150509291505056fea26469706673582212205443ec313ecb8c2e08ca8a30687daed4c3b666f9318ae72ccbe9033479c8b8be64736f6c634300080a0033",
}

// TestDistributeABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDistributeMetaData.ABI instead.
var TestDistributeABI = TestDistributeMetaData.ABI

// TestDistributeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestDistributeMetaData.Bin instead.
var TestDistributeBin = TestDistributeMetaData.Bin

// DeployTestDistribute deploys a new Ethereum contract, binding an instance of TestDistribute to it.
func DeployTestDistribute(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestDistribute, error) {
	parsed, err := TestDistributeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDistributeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestDistribute{TestDistributeCaller: TestDistributeCaller{contract: contract}, TestDistributeTransactor: TestDistributeTransactor{contract: contract}, TestDistributeFilterer: TestDistributeFilterer{contract: contract}}, nil
}

// TestDistribute is an auto generated Go binding around an Ethereum contract.
type TestDistribute struct {
	TestDistributeCaller     // Read-only binding to the contract
	TestDistributeTransactor // Write-only binding to the contract
	TestDistributeFilterer   // Log filterer for contract events
}

// TestDistributeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestDistributeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDistributeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestDistributeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDistributeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestDistributeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDistributeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestDistributeSession struct {
	Contract     *TestDistribute   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestDistributeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestDistributeCallerSession struct {
	Contract *TestDistributeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// TestDistributeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestDistributeTransactorSession struct {
	Contract     *TestDistributeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// TestDistributeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestDistributeRaw struct {
	Contract *TestDistribute // Generic contract binding to access the raw methods on
}

// TestDistributeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestDistributeCallerRaw struct {
	Contract *TestDistributeCaller // Generic read-only contract binding to access the raw methods on
}

// TestDistributeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestDistributeTransactorRaw struct {
	Contract *TestDistributeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestDistribute creates a new instance of TestDistribute, bound to a specific deployed contract.
func NewTestDistribute(address common.Address, backend bind.ContractBackend) (*TestDistribute, error) {
	contract, err := bindTestDistribute(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestDistribute{TestDistributeCaller: TestDistributeCaller{contract: contract}, TestDistributeTransactor: TestDistributeTransactor{contract: contract}, TestDistributeFilterer: TestDistributeFilterer{contract: contract}}, nil
}

// NewTestDistributeCaller creates a new read-only instance of TestDistribute, bound to a specific deployed contract.
func NewTestDistributeCaller(address common.Address, caller bind.ContractCaller) (*TestDistributeCaller, error) {
	contract, err := bindTestDistribute(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestDistributeCaller{contract: contract}, nil
}

// NewTestDistributeTransactor creates a new write-only instance of TestDistribute, bound to a specific deployed contract.
func NewTestDistributeTransactor(address common.Address, transactor bind.ContractTransactor) (*TestDistributeTransactor, error) {
	contract, err := bindTestDistribute(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestDistributeTransactor{contract: contract}, nil
}

// NewTestDistributeFilterer creates a new log filterer instance of TestDistribute, bound to a specific deployed contract.
func NewTestDistributeFilterer(address common.Address, filterer bind.ContractFilterer) (*TestDistributeFilterer, error) {
	contract, err := bindTestDistribute(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestDistributeFilterer{contract: contract}, nil
}

// bindTestDistribute binds a generic wrapper to an already deployed contract.
func bindTestDistribute(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestDistributeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDistribute *TestDistributeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDistribute.Contract.TestDistributeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDistribute *TestDistributeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDistribute.Contract.TestDistributeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDistribute *TestDistributeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDistribute.Contract.TestDistributeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDistribute *TestDistributeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDistribute.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDistribute *TestDistributeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDistribute.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDistribute *TestDistributeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDistribute.Contract.contract.Transact(opts, method, params...)
}

// DistributeThroughContract is a paid mutator transaction binding the contract method 0x50b54e84.
//
// Solidity: function distributeThroughContract(address zrc20, uint256 amount) returns(bool)
func (_TestDistribute *TestDistributeTransactor) DistributeThroughContract(opts *bind.TransactOpts, zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestDistribute.contract.Transact(opts, "distributeThroughContract", zrc20, amount)
}

// DistributeThroughContract is a paid mutator transaction binding the contract method 0x50b54e84.
//
// Solidity: function distributeThroughContract(address zrc20, uint256 amount) returns(bool)
func (_TestDistribute *TestDistributeSession) DistributeThroughContract(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestDistribute.Contract.DistributeThroughContract(&_TestDistribute.TransactOpts, zrc20, amount)
}

// DistributeThroughContract is a paid mutator transaction binding the contract method 0x50b54e84.
//
// Solidity: function distributeThroughContract(address zrc20, uint256 amount) returns(bool)
func (_TestDistribute *TestDistributeTransactorSession) DistributeThroughContract(zrc20 common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestDistribute.Contract.DistributeThroughContract(&_TestDistribute.TransactOpts, zrc20, amount)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestDistribute *TestDistributeTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TestDistribute.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestDistribute *TestDistributeSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TestDistribute.Contract.Fallback(&_TestDistribute.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestDistribute *TestDistributeTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TestDistribute.Contract.Fallback(&_TestDistribute.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDistribute *TestDistributeTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDistribute.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDistribute *TestDistributeSession) Receive() (*types.Transaction, error) {
	return _TestDistribute.Contract.Receive(&_TestDistribute.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDistribute *TestDistributeTransactorSession) Receive() (*types.Transaction, error) {
	return _TestDistribute.Contract.Receive(&_TestDistribute.TransactOpts)
}

// TestDistributeDistributedIterator is returned from FilterDistributed and is used to iterate over the raw logs and unpacked data for Distributed events raised by the TestDistribute contract.
type TestDistributeDistributedIterator struct {
	Event *TestDistributeDistributed // Event containing the contract specifics and raw log

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
func (it *TestDistributeDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestDistributeDistributed)
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
		it.Event = new(TestDistributeDistributed)
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
func (it *TestDistributeDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestDistributeDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestDistributeDistributed represents a Distributed event raised by the TestDistribute contract.
type TestDistributeDistributed struct {
	Zrc20Distributor common.Address
	Zrc20Token       common.Address
	Amount           *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDistributed is a free log retrieval operation binding the contract event 0xad4a9acf26d8bba7a8cf1a41160d59be042ee554578e256c98d2ab74cdd43542.
//
// Solidity: event Distributed(address indexed zrc20_distributor, address indexed zrc20_token, uint256 amount)
func (_TestDistribute *TestDistributeFilterer) FilterDistributed(opts *bind.FilterOpts, zrc20_distributor []common.Address, zrc20_token []common.Address) (*TestDistributeDistributedIterator, error) {

	var zrc20_distributorRule []interface{}
	for _, zrc20_distributorItem := range zrc20_distributor {
		zrc20_distributorRule = append(zrc20_distributorRule, zrc20_distributorItem)
	}
	var zrc20_tokenRule []interface{}
	for _, zrc20_tokenItem := range zrc20_token {
		zrc20_tokenRule = append(zrc20_tokenRule, zrc20_tokenItem)
	}

	logs, sub, err := _TestDistribute.contract.FilterLogs(opts, "Distributed", zrc20_distributorRule, zrc20_tokenRule)
	if err != nil {
		return nil, err
	}
	return &TestDistributeDistributedIterator{contract: _TestDistribute.contract, event: "Distributed", logs: logs, sub: sub}, nil
}

// WatchDistributed is a free log subscription operation binding the contract event 0xad4a9acf26d8bba7a8cf1a41160d59be042ee554578e256c98d2ab74cdd43542.
//
// Solidity: event Distributed(address indexed zrc20_distributor, address indexed zrc20_token, uint256 amount)
func (_TestDistribute *TestDistributeFilterer) WatchDistributed(opts *bind.WatchOpts, sink chan<- *TestDistributeDistributed, zrc20_distributor []common.Address, zrc20_token []common.Address) (event.Subscription, error) {

	var zrc20_distributorRule []interface{}
	for _, zrc20_distributorItem := range zrc20_distributor {
		zrc20_distributorRule = append(zrc20_distributorRule, zrc20_distributorItem)
	}
	var zrc20_tokenRule []interface{}
	for _, zrc20_tokenItem := range zrc20_token {
		zrc20_tokenRule = append(zrc20_tokenRule, zrc20_tokenItem)
	}

	logs, sub, err := _TestDistribute.contract.WatchLogs(opts, "Distributed", zrc20_distributorRule, zrc20_tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestDistributeDistributed)
				if err := _TestDistribute.contract.UnpackLog(event, "Distributed", log); err != nil {
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

// ParseDistributed is a log parse operation binding the contract event 0xad4a9acf26d8bba7a8cf1a41160d59be042ee554578e256c98d2ab74cdd43542.
//
// Solidity: event Distributed(address indexed zrc20_distributor, address indexed zrc20_token, uint256 amount)
func (_TestDistribute *TestDistributeFilterer) ParseDistributed(log types.Log) (*TestDistributeDistributed, error) {
	event := new(TestDistributeDistributed)
	if err := _TestDistribute.contract.UnpackLog(event, "Distributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
