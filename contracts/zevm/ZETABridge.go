// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zevm

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
)

// ZETABridgeMetaData contains all meta data concerning the ZETABridge contract.
var ZETABridgeMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"toChainID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"ZetaSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FUNGIBLE_MODULE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"toChainID\",\"type\":\"uint256\"}],\"name\":\"sendZeta\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// ZETABridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use ZETABridgeMetaData.ABI instead.
var ZETABridgeABI = ZETABridgeMetaData.ABI

// ZETABridge is an auto generated Go binding around an Ethereum contract.
type ZETABridge struct {
	ZETABridgeCaller     // Read-only binding to the contract
	ZETABridgeTransactor // Write-only binding to the contract
	ZETABridgeFilterer   // Log filterer for contract events
}

// ZETABridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZETABridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZETABridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZETABridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZETABridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZETABridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZETABridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZETABridgeSession struct {
	Contract     *ZETABridge       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZETABridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZETABridgeCallerSession struct {
	Contract *ZETABridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ZETABridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZETABridgeTransactorSession struct {
	Contract     *ZETABridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ZETABridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZETABridgeRaw struct {
	Contract *ZETABridge // Generic contract binding to access the raw methods on
}

// ZETABridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZETABridgeCallerRaw struct {
	Contract *ZETABridgeCaller // Generic read-only contract binding to access the raw methods on
}

// ZETABridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZETABridgeTransactorRaw struct {
	Contract *ZETABridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZETABridge creates a new instance of ZETABridge, bound to a specific deployed contract.
func NewZETABridge(address common.Address, backend bind.ContractBackend) (*ZETABridge, error) {
	contract, err := bindZETABridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZETABridge{ZETABridgeCaller: ZETABridgeCaller{contract: contract}, ZETABridgeTransactor: ZETABridgeTransactor{contract: contract}, ZETABridgeFilterer: ZETABridgeFilterer{contract: contract}}, nil
}

// NewZETABridgeCaller creates a new read-only instance of ZETABridge, bound to a specific deployed contract.
func NewZETABridgeCaller(address common.Address, caller bind.ContractCaller) (*ZETABridgeCaller, error) {
	contract, err := bindZETABridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZETABridgeCaller{contract: contract}, nil
}

// NewZETABridgeTransactor creates a new write-only instance of ZETABridge, bound to a specific deployed contract.
func NewZETABridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*ZETABridgeTransactor, error) {
	contract, err := bindZETABridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZETABridgeTransactor{contract: contract}, nil
}

// NewZETABridgeFilterer creates a new log filterer instance of ZETABridge, bound to a specific deployed contract.
func NewZETABridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*ZETABridgeFilterer, error) {
	contract, err := bindZETABridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZETABridgeFilterer{contract: contract}, nil
}

// bindZETABridge binds a generic wrapper to an already deployed contract.
func bindZETABridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ZETABridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZETABridge *ZETABridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZETABridge.Contract.ZETABridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZETABridge *ZETABridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZETABridge.Contract.ZETABridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZETABridge *ZETABridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZETABridge.Contract.ZETABridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZETABridge *ZETABridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZETABridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZETABridge *ZETABridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZETABridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZETABridge *ZETABridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZETABridge.Contract.contract.Transact(opts, method, params...)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZETABridge *ZETABridgeCaller) FUNGIBLEMODULEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZETABridge.contract.Call(opts, &out, "FUNGIBLE_MODULE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZETABridge *ZETABridgeSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _ZETABridge.Contract.FUNGIBLEMODULEADDRESS(&_ZETABridge.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZETABridge *ZETABridgeCallerSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _ZETABridge.Contract.FUNGIBLEMODULEADDRESS(&_ZETABridge.CallOpts)
}

// SendZeta is a paid mutator transaction binding the contract method 0x1b2a853e.
//
// Solidity: function sendZeta(bytes to, uint256 toChainID) payable returns()
func (_ZETABridge *ZETABridgeTransactor) SendZeta(opts *bind.TransactOpts, to []byte, toChainID *big.Int) (*types.Transaction, error) {
	return _ZETABridge.contract.Transact(opts, "sendZeta", to, toChainID)
}

// SendZeta is a paid mutator transaction binding the contract method 0x1b2a853e.
//
// Solidity: function sendZeta(bytes to, uint256 toChainID) payable returns()
func (_ZETABridge *ZETABridgeSession) SendZeta(to []byte, toChainID *big.Int) (*types.Transaction, error) {
	return _ZETABridge.Contract.SendZeta(&_ZETABridge.TransactOpts, to, toChainID)
}

// SendZeta is a paid mutator transaction binding the contract method 0x1b2a853e.
//
// Solidity: function sendZeta(bytes to, uint256 toChainID) payable returns()
func (_ZETABridge *ZETABridgeTransactorSession) SendZeta(to []byte, toChainID *big.Int) (*types.Transaction, error) {
	return _ZETABridge.Contract.SendZeta(&_ZETABridge.TransactOpts, to, toChainID)
}

// ZETABridgeZetaSentIterator is returned from FilterZetaSent and is used to iterate over the raw logs and unpacked data for ZetaSent events raised by the ZETABridge contract.
type ZETABridgeZetaSentIterator struct {
	Event *ZETABridgeZetaSent // Event containing the contract specifics and raw log

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
func (it *ZETABridgeZetaSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZETABridgeZetaSent)
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
		it.Event = new(ZETABridgeZetaSent)
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
func (it *ZETABridgeZetaSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZETABridgeZetaSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZETABridgeZetaSent represents a ZetaSent event raised by the ZETABridge contract.
type ZETABridgeZetaSent struct {
	To        []byte
	ToChainID *big.Int
	Value     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterZetaSent is a free log retrieval operation binding the contract event 0xdd9d4f9ff5d3318a370825c2a4e875bd316005549a1bf69f340ea266d173389d.
//
// Solidity: event ZetaSent(bytes to, uint256 toChainID, uint256 value)
func (_ZETABridge *ZETABridgeFilterer) FilterZetaSent(opts *bind.FilterOpts) (*ZETABridgeZetaSentIterator, error) {

	logs, sub, err := _ZETABridge.contract.FilterLogs(opts, "ZetaSent")
	if err != nil {
		return nil, err
	}
	return &ZETABridgeZetaSentIterator{contract: _ZETABridge.contract, event: "ZetaSent", logs: logs, sub: sub}, nil
}

// WatchZetaSent is a free log subscription operation binding the contract event 0xdd9d4f9ff5d3318a370825c2a4e875bd316005549a1bf69f340ea266d173389d.
//
// Solidity: event ZetaSent(bytes to, uint256 toChainID, uint256 value)
func (_ZETABridge *ZETABridgeFilterer) WatchZetaSent(opts *bind.WatchOpts, sink chan<- *ZETABridgeZetaSent) (event.Subscription, error) {

	logs, sub, err := _ZETABridge.contract.WatchLogs(opts, "ZetaSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZETABridgeZetaSent)
				if err := _ZETABridge.contract.UnpackLog(event, "ZetaSent", log); err != nil {
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

// ParseZetaSent is a log parse operation binding the contract event 0xdd9d4f9ff5d3318a370825c2a4e875bd316005549a1bf69f340ea266d173389d.
//
// Solidity: event ZetaSent(bytes to, uint256 toChainID, uint256 value)
func (_ZETABridge *ZETABridgeFilterer) ParseZetaSent(log types.Log) (*ZETABridgeZetaSent, error) {
	event := new(ZETABridgeZetaSent)
	if err := _ZETABridge.contract.UnpackLog(event, "ZetaSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
