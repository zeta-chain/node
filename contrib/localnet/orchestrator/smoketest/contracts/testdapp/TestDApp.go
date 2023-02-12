// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testdapp

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

// ZetaInterfacesZetaRevert is an auto generated low-level Go binding around an user-defined struct.
type ZetaInterfacesZetaRevert struct {
	ZetaTxSenderAddress common.Address
	SourceChainId       *big.Int
	DestinationAddress  []byte
	DestinationChainId  *big.Int
	RemainingZetaValue  *big.Int
	Message             []byte
}

// TestDAppMetaData contains all meta data concerning the TestDApp contract.
var TestDAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_connector\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidMessageType\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"HelloWorldEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"RevertedHelloWorldEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"HELLO_WORLD_MESSAGE_TYPE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"zetaTxSenderAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"destinationAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"zetaValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"internalType\":\"structZetaInterfaces.ZetaMessage\",\"name\":\"zetaMessage\",\"type\":\"tuple\"}],\"name\":\"onZetaMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"zetaTxSenderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destinationAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingZetaValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"internalType\":\"structZetaInterfaces.ZetaRevert\",\"name\":\"zetaRevert\",\"type\":\"tuple\"}],\"name\":\"onZetaRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"destinationAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"doRevert\",\"type\":\"bool\"}],\"name\":\"sendHelloWorld\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0xnull",
}

// TestDAppABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDAppMetaData.ABI instead.
var TestDAppABI = TestDAppMetaData.ABI

// TestDAppBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestDAppMetaData.Bin instead.
var TestDAppBin = TestDAppMetaData.Bin

// DeployTestDApp deploys a new Ethereum contract, binding an instance of TestDApp to it.
func DeployTestDApp(auth *bind.TransactOpts, backend bind.ContractBackend, _connector common.Address) (common.Address, *types.Transaction, *TestDApp, error) {
	parsed, err := TestDAppMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDAppBin), backend, _connector)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestDApp{TestDAppCaller: TestDAppCaller{contract: contract}, TestDAppTransactor: TestDAppTransactor{contract: contract}, TestDAppFilterer: TestDAppFilterer{contract: contract}}, nil
}

// TestDApp is an auto generated Go binding around an Ethereum contract.
type TestDApp struct {
	TestDAppCaller     // Read-only binding to the contract
	TestDAppTransactor // Write-only binding to the contract
	TestDAppFilterer   // Log filterer for contract events
}

// TestDAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestDAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestDAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestDAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestDAppSession struct {
	Contract     *TestDApp         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestDAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestDAppCallerSession struct {
	Contract *TestDAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TestDAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestDAppTransactorSession struct {
	Contract     *TestDAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TestDAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestDAppRaw struct {
	Contract *TestDApp // Generic contract binding to access the raw methods on
}

// TestDAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestDAppCallerRaw struct {
	Contract *TestDAppCaller // Generic read-only contract binding to access the raw methods on
}

// TestDAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestDAppTransactorRaw struct {
	Contract *TestDAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestDApp creates a new instance of TestDApp, bound to a specific deployed contract.
func NewTestDApp(address common.Address, backend bind.ContractBackend) (*TestDApp, error) {
	contract, err := bindTestDApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestDApp{TestDAppCaller: TestDAppCaller{contract: contract}, TestDAppTransactor: TestDAppTransactor{contract: contract}, TestDAppFilterer: TestDAppFilterer{contract: contract}}, nil
}

// NewTestDAppCaller creates a new read-only instance of TestDApp, bound to a specific deployed contract.
func NewTestDAppCaller(address common.Address, caller bind.ContractCaller) (*TestDAppCaller, error) {
	contract, err := bindTestDApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppCaller{contract: contract}, nil
}

// NewTestDAppTransactor creates a new write-only instance of TestDApp, bound to a specific deployed contract.
func NewTestDAppTransactor(address common.Address, transactor bind.ContractTransactor) (*TestDAppTransactor, error) {
	contract, err := bindTestDApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppTransactor{contract: contract}, nil
}

// NewTestDAppFilterer creates a new log filterer instance of TestDApp, bound to a specific deployed contract.
func NewTestDAppFilterer(address common.Address, filterer bind.ContractFilterer) (*TestDAppFilterer, error) {
	contract, err := bindTestDApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestDAppFilterer{contract: contract}, nil
}

// bindTestDApp binds a generic wrapper to an already deployed contract.
func bindTestDApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestDAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDApp *TestDAppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDApp.Contract.TestDAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDApp *TestDAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDApp.Contract.TestDAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDApp *TestDAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDApp.Contract.TestDAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDApp *TestDAppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDApp *TestDAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDApp *TestDAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDApp.Contract.contract.Transact(opts, method, params...)
}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDApp *TestDAppCaller) HELLOWORLDMESSAGETYPE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TestDApp.contract.Call(opts, &out, "HELLO_WORLD_MESSAGE_TYPE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDApp *TestDAppSession) HELLOWORLDMESSAGETYPE() ([32]byte, error) {
	return _TestDApp.Contract.HELLOWORLDMESSAGETYPE(&_TestDApp.CallOpts)
}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDApp *TestDAppCallerSession) HELLOWORLDMESSAGETYPE() ([32]byte, error) {
	return _TestDApp.Contract.HELLOWORLDMESSAGETYPE(&_TestDApp.CallOpts)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onZetaMessage((bytes,uint256,address,uint256,bytes) zetaMessage) returns()
func (_TestDApp *TestDAppTransactor) OnZetaMessage(opts *bind.TransactOpts, zetaMessage ZetaInterfacesZetaMessage) (*types.Transaction, error) {
	return _TestDApp.contract.Transact(opts, "onZetaMessage", zetaMessage)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onZetaMessage((bytes,uint256,address,uint256,bytes) zetaMessage) returns()
func (_TestDApp *TestDAppSession) OnZetaMessage(zetaMessage ZetaInterfacesZetaMessage) (*types.Transaction, error) {
	return _TestDApp.Contract.OnZetaMessage(&_TestDApp.TransactOpts, zetaMessage)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onZetaMessage((bytes,uint256,address,uint256,bytes) zetaMessage) returns()
func (_TestDApp *TestDAppTransactorSession) OnZetaMessage(zetaMessage ZetaInterfacesZetaMessage) (*types.Transaction, error) {
	return _TestDApp.Contract.OnZetaMessage(&_TestDApp.TransactOpts, zetaMessage)
}

// OnZetaRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onZetaRevert((address,uint256,bytes,uint256,uint256,bytes) zetaRevert) returns()
func (_TestDApp *TestDAppTransactor) OnZetaRevert(opts *bind.TransactOpts, zetaRevert ZetaInterfacesZetaRevert) (*types.Transaction, error) {
	return _TestDApp.contract.Transact(opts, "onZetaRevert", zetaRevert)
}

// OnZetaRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onZetaRevert((address,uint256,bytes,uint256,uint256,bytes) zetaRevert) returns()
func (_TestDApp *TestDAppSession) OnZetaRevert(zetaRevert ZetaInterfacesZetaRevert) (*types.Transaction, error) {
	return _TestDApp.Contract.OnZetaRevert(&_TestDApp.TransactOpts, zetaRevert)
}

// OnZetaRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onZetaRevert((address,uint256,bytes,uint256,uint256,bytes) zetaRevert) returns()
func (_TestDApp *TestDAppTransactorSession) OnZetaRevert(zetaRevert ZetaInterfacesZetaRevert) (*types.Transaction, error) {
	return _TestDApp.Contract.OnZetaRevert(&_TestDApp.TransactOpts, zetaRevert)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDApp *TestDAppTransactor) SendHelloWorld(opts *bind.TransactOpts, destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDApp.contract.Transact(opts, "sendHelloWorld", destinationAddress, destinationChainId, value, doRevert)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDApp *TestDAppSession) SendHelloWorld(destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDApp.Contract.SendHelloWorld(&_TestDApp.TransactOpts, destinationAddress, destinationChainId, value, doRevert)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDApp *TestDAppTransactorSession) SendHelloWorld(destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDApp.Contract.SendHelloWorld(&_TestDApp.TransactOpts, destinationAddress, destinationChainId, value, doRevert)
}

// TestDAppHelloWorldEventIterator is returned from FilterHelloWorldEvent and is used to iterate over the raw logs and unpacked data for HelloWorldEvent events raised by the TestDApp contract.
type TestDAppHelloWorldEventIterator struct {
	Event *TestDAppHelloWorldEvent // Event containing the contract specifics and raw log

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
func (it *TestDAppHelloWorldEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestDAppHelloWorldEvent)
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
		it.Event = new(TestDAppHelloWorldEvent)
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
func (it *TestDAppHelloWorldEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestDAppHelloWorldEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestDAppHelloWorldEvent represents a HelloWorldEvent event raised by the TestDApp contract.
type TestDAppHelloWorldEvent struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterHelloWorldEvent is a free log retrieval operation binding the contract event 0x3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d1.
//
// Solidity: event HelloWorldEvent()
func (_TestDApp *TestDAppFilterer) FilterHelloWorldEvent(opts *bind.FilterOpts) (*TestDAppHelloWorldEventIterator, error) {

	logs, sub, err := _TestDApp.contract.FilterLogs(opts, "HelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return &TestDAppHelloWorldEventIterator{contract: _TestDApp.contract, event: "HelloWorldEvent", logs: logs, sub: sub}, nil
}

// WatchHelloWorldEvent is a free log subscription operation binding the contract event 0x3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d1.
//
// Solidity: event HelloWorldEvent()
func (_TestDApp *TestDAppFilterer) WatchHelloWorldEvent(opts *bind.WatchOpts, sink chan<- *TestDAppHelloWorldEvent) (event.Subscription, error) {

	logs, sub, err := _TestDApp.contract.WatchLogs(opts, "HelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestDAppHelloWorldEvent)
				if err := _TestDApp.contract.UnpackLog(event, "HelloWorldEvent", log); err != nil {
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
func (_TestDApp *TestDAppFilterer) ParseHelloWorldEvent(log types.Log) (*TestDAppHelloWorldEvent, error) {
	event := new(TestDAppHelloWorldEvent)
	if err := _TestDApp.contract.UnpackLog(event, "HelloWorldEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestDAppRevertedHelloWorldEventIterator is returned from FilterRevertedHelloWorldEvent and is used to iterate over the raw logs and unpacked data for RevertedHelloWorldEvent events raised by the TestDApp contract.
type TestDAppRevertedHelloWorldEventIterator struct {
	Event *TestDAppRevertedHelloWorldEvent // Event containing the contract specifics and raw log

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
func (it *TestDAppRevertedHelloWorldEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestDAppRevertedHelloWorldEvent)
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
		it.Event = new(TestDAppRevertedHelloWorldEvent)
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
func (it *TestDAppRevertedHelloWorldEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestDAppRevertedHelloWorldEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestDAppRevertedHelloWorldEvent represents a RevertedHelloWorldEvent event raised by the TestDApp contract.
type TestDAppRevertedHelloWorldEvent struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRevertedHelloWorldEvent is a free log retrieval operation binding the contract event 0x4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a521.
//
// Solidity: event RevertedHelloWorldEvent()
func (_TestDApp *TestDAppFilterer) FilterRevertedHelloWorldEvent(opts *bind.FilterOpts) (*TestDAppRevertedHelloWorldEventIterator, error) {

	logs, sub, err := _TestDApp.contract.FilterLogs(opts, "RevertedHelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return &TestDAppRevertedHelloWorldEventIterator{contract: _TestDApp.contract, event: "RevertedHelloWorldEvent", logs: logs, sub: sub}, nil
}

// WatchRevertedHelloWorldEvent is a free log subscription operation binding the contract event 0x4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a521.
//
// Solidity: event RevertedHelloWorldEvent()
func (_TestDApp *TestDAppFilterer) WatchRevertedHelloWorldEvent(opts *bind.WatchOpts, sink chan<- *TestDAppRevertedHelloWorldEvent) (event.Subscription, error) {

	logs, sub, err := _TestDApp.contract.WatchLogs(opts, "RevertedHelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestDAppRevertedHelloWorldEvent)
				if err := _TestDApp.contract.UnpackLog(event, "RevertedHelloWorldEvent", log); err != nil {
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
func (_TestDApp *TestDAppFilterer) ParseRevertedHelloWorldEvent(log types.Log) (*TestDAppRevertedHelloWorldEvent, error) {
	event := new(TestDAppRevertedHelloWorldEvent)
	if err := _TestDApp.contract.UnpackLog(event, "RevertedHelloWorldEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
