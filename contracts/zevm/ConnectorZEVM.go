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
	_ = abi.ConvertType
)

// ZetaInterfacesSendInput is an auto generated low-level Go binding around an user-defined struct.
type ZetaInterfacesSendInput struct {
	DestinationChainId  *big.Int
	DestinationAddress  []byte
	DestinationGasLimit *big.Int
	Message             []byte
	ZetaValueAndGas     *big.Int
	ZetaParams          []byte
}

// ZetaConnectorZEVMMetaData contains all meta data concerning the ZetaConnectorZEVM contract.
var ZetaConnectorZEVMMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wzeta\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sourceTxOriginAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zetaTxSenderAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"destinationAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"zetaValueAndGas\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationGasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"zetaParams\",\"type\":\"bytes\"}],\"name\":\"ZetaSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FUNGIBLE_MODULE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destinationAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"destinationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"zetaValueAndGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"zetaParams\",\"type\":\"bytes\"}],\"internalType\":\"structZetaInterfaces.SendInput\",\"name\":\"input\",\"type\":\"tuple\"}],\"name\":\"send\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wzeta\",\"type\":\"address\"}],\"name\":\"setWzetaAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wzeta\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// ZetaConnectorZEVMABI is the input ABI used to generate the binding from.
// Deprecated: Use ZetaConnectorZEVMMetaData.ABI instead.
var ZetaConnectorZEVMABI = ZetaConnectorZEVMMetaData.ABI

// ZetaConnectorZEVM is an auto generated Go binding around an Ethereum contract.
type ZetaConnectorZEVM struct {
	ZetaConnectorZEVMCaller     // Read-only binding to the contract
	ZetaConnectorZEVMTransactor // Write-only binding to the contract
	ZetaConnectorZEVMFilterer   // Log filterer for contract events
}

// ZetaConnectorZEVMCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZetaConnectorZEVMCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZetaConnectorZEVMTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZetaConnectorZEVMTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZetaConnectorZEVMFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZetaConnectorZEVMFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZetaConnectorZEVMSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZetaConnectorZEVMSession struct {
	Contract     *ZetaConnectorZEVM // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ZetaConnectorZEVMCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZetaConnectorZEVMCallerSession struct {
	Contract *ZetaConnectorZEVMCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ZetaConnectorZEVMTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZetaConnectorZEVMTransactorSession struct {
	Contract     *ZetaConnectorZEVMTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ZetaConnectorZEVMRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZetaConnectorZEVMRaw struct {
	Contract *ZetaConnectorZEVM // Generic contract binding to access the raw methods on
}

// ZetaConnectorZEVMCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZetaConnectorZEVMCallerRaw struct {
	Contract *ZetaConnectorZEVMCaller // Generic read-only contract binding to access the raw methods on
}

// ZetaConnectorZEVMTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZetaConnectorZEVMTransactorRaw struct {
	Contract *ZetaConnectorZEVMTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZetaConnectorZEVM creates a new instance of ZetaConnectorZEVM, bound to a specific deployed contract.
func NewZetaConnectorZEVM(address common.Address, backend bind.ContractBackend) (*ZetaConnectorZEVM, error) {
	contract, err := bindZetaConnectorZEVM(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZetaConnectorZEVM{ZetaConnectorZEVMCaller: ZetaConnectorZEVMCaller{contract: contract}, ZetaConnectorZEVMTransactor: ZetaConnectorZEVMTransactor{contract: contract}, ZetaConnectorZEVMFilterer: ZetaConnectorZEVMFilterer{contract: contract}}, nil
}

// NewZetaConnectorZEVMCaller creates a new read-only instance of ZetaConnectorZEVM, bound to a specific deployed contract.
func NewZetaConnectorZEVMCaller(address common.Address, caller bind.ContractCaller) (*ZetaConnectorZEVMCaller, error) {
	contract, err := bindZetaConnectorZEVM(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZetaConnectorZEVMCaller{contract: contract}, nil
}

// NewZetaConnectorZEVMTransactor creates a new write-only instance of ZetaConnectorZEVM, bound to a specific deployed contract.
func NewZetaConnectorZEVMTransactor(address common.Address, transactor bind.ContractTransactor) (*ZetaConnectorZEVMTransactor, error) {
	contract, err := bindZetaConnectorZEVM(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZetaConnectorZEVMTransactor{contract: contract}, nil
}

// NewZetaConnectorZEVMFilterer creates a new log filterer instance of ZetaConnectorZEVM, bound to a specific deployed contract.
func NewZetaConnectorZEVMFilterer(address common.Address, filterer bind.ContractFilterer) (*ZetaConnectorZEVMFilterer, error) {
	contract, err := bindZetaConnectorZEVM(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZetaConnectorZEVMFilterer{contract: contract}, nil
}

// bindZetaConnectorZEVM binds a generic wrapper to an already deployed contract.
func bindZetaConnectorZEVM(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZetaConnectorZEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZetaConnectorZEVM *ZetaConnectorZEVMRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZetaConnectorZEVM.Contract.ZetaConnectorZEVMCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZetaConnectorZEVM *ZetaConnectorZEVMRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.Contract.ZetaConnectorZEVMTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZetaConnectorZEVM *ZetaConnectorZEVMRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.Contract.ZetaConnectorZEVMTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZetaConnectorZEVM *ZetaConnectorZEVMCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZetaConnectorZEVM.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZetaConnectorZEVM *ZetaConnectorZEVMTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZetaConnectorZEVM *ZetaConnectorZEVMTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.Contract.contract.Transact(opts, method, params...)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZetaConnectorZEVM *ZetaConnectorZEVMCaller) FUNGIBLEMODULEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZetaConnectorZEVM.contract.Call(opts, &out, "FUNGIBLE_MODULE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZetaConnectorZEVM *ZetaConnectorZEVMSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _ZetaConnectorZEVM.Contract.FUNGIBLEMODULEADDRESS(&_ZetaConnectorZEVM.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_ZetaConnectorZEVM *ZetaConnectorZEVMCallerSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _ZetaConnectorZEVM.Contract.FUNGIBLEMODULEADDRESS(&_ZetaConnectorZEVM.CallOpts)
}

// Wzeta is a free data retrieval call binding the contract method 0x00173d46.
//
// Solidity: function wzeta() view returns(address)
func (_ZetaConnectorZEVM *ZetaConnectorZEVMCaller) Wzeta(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZetaConnectorZEVM.contract.Call(opts, &out, "wzeta")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Wzeta is a free data retrieval call binding the contract method 0x00173d46.
//
// Solidity: function wzeta() view returns(address)
func (_ZetaConnectorZEVM *ZetaConnectorZEVMSession) Wzeta() (common.Address, error) {
	return _ZetaConnectorZEVM.Contract.Wzeta(&_ZetaConnectorZEVM.CallOpts)
}

// Wzeta is a free data retrieval call binding the contract method 0x00173d46.
//
// Solidity: function wzeta() view returns(address)
func (_ZetaConnectorZEVM *ZetaConnectorZEVMCallerSession) Wzeta() (common.Address, error) {
	return _ZetaConnectorZEVM.Contract.Wzeta(&_ZetaConnectorZEVM.CallOpts)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_ZetaConnectorZEVM *ZetaConnectorZEVMTransactor) Send(opts *bind.TransactOpts, input ZetaInterfacesSendInput) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.contract.Transact(opts, "send", input)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_ZetaConnectorZEVM *ZetaConnectorZEVMSession) Send(input ZetaInterfacesSendInput) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.Contract.Send(&_ZetaConnectorZEVM.TransactOpts, input)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_ZetaConnectorZEVM *ZetaConnectorZEVMTransactorSession) Send(input ZetaInterfacesSendInput) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.Contract.Send(&_ZetaConnectorZEVM.TransactOpts, input)
}

// SetWzetaAddress is a paid mutator transaction binding the contract method 0xeb3bacbd.
//
// Solidity: function setWzetaAddress(address _wzeta) returns()
func (_ZetaConnectorZEVM *ZetaConnectorZEVMTransactor) SetWzetaAddress(opts *bind.TransactOpts, _wzeta common.Address) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.contract.Transact(opts, "setWzetaAddress", _wzeta)
}

// SetWzetaAddress is a paid mutator transaction binding the contract method 0xeb3bacbd.
//
// Solidity: function setWzetaAddress(address _wzeta) returns()
func (_ZetaConnectorZEVM *ZetaConnectorZEVMSession) SetWzetaAddress(_wzeta common.Address) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.Contract.SetWzetaAddress(&_ZetaConnectorZEVM.TransactOpts, _wzeta)
}

// SetWzetaAddress is a paid mutator transaction binding the contract method 0xeb3bacbd.
//
// Solidity: function setWzetaAddress(address _wzeta) returns()
func (_ZetaConnectorZEVM *ZetaConnectorZEVMTransactorSession) SetWzetaAddress(_wzeta common.Address) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.Contract.SetWzetaAddress(&_ZetaConnectorZEVM.TransactOpts, _wzeta)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ZetaConnectorZEVM *ZetaConnectorZEVMTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZetaConnectorZEVM.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ZetaConnectorZEVM *ZetaConnectorZEVMSession) Receive() (*types.Transaction, error) {
	return _ZetaConnectorZEVM.Contract.Receive(&_ZetaConnectorZEVM.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ZetaConnectorZEVM *ZetaConnectorZEVMTransactorSession) Receive() (*types.Transaction, error) {
	return _ZetaConnectorZEVM.Contract.Receive(&_ZetaConnectorZEVM.TransactOpts)
}

// ZetaConnectorZEVMZetaSentIterator is returned from FilterZetaSent and is used to iterate over the raw logs and unpacked data for ZetaSent events raised by the ZetaConnectorZEVM contract.
type ZetaConnectorZEVMZetaSentIterator struct {
	Event *ZetaConnectorZEVMZetaSent // Event containing the contract specifics and raw log

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
func (it *ZetaConnectorZEVMZetaSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZetaConnectorZEVMZetaSent)
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
		it.Event = new(ZetaConnectorZEVMZetaSent)
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
func (it *ZetaConnectorZEVMZetaSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZetaConnectorZEVMZetaSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZetaConnectorZEVMZetaSent represents a ZetaSent event raised by the ZetaConnectorZEVM contract.
type ZetaConnectorZEVMZetaSent struct {
	SourceTxOriginAddress common.Address
	ZetaTxSenderAddress   common.Address
	DestinationChainId    *big.Int
	DestinationAddress    []byte
	ZetaValueAndGas       *big.Int
	DestinationGasLimit   *big.Int
	Message               []byte
	ZetaParams            []byte
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterZetaSent is a free log retrieval operation binding the contract event 0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4.
//
// Solidity: event ZetaSent(address sourceTxOriginAddress, address indexed zetaTxSenderAddress, uint256 indexed destinationChainId, bytes destinationAddress, uint256 zetaValueAndGas, uint256 destinationGasLimit, bytes message, bytes zetaParams)
func (_ZetaConnectorZEVM *ZetaConnectorZEVMFilterer) FilterZetaSent(opts *bind.FilterOpts, zetaTxSenderAddress []common.Address, destinationChainId []*big.Int) (*ZetaConnectorZEVMZetaSentIterator, error) {

	var zetaTxSenderAddressRule []interface{}
	for _, zetaTxSenderAddressItem := range zetaTxSenderAddress {
		zetaTxSenderAddressRule = append(zetaTxSenderAddressRule, zetaTxSenderAddressItem)
	}
	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	logs, sub, err := _ZetaConnectorZEVM.contract.FilterLogs(opts, "ZetaSent", zetaTxSenderAddressRule, destinationChainIdRule)
	if err != nil {
		return nil, err
	}
	return &ZetaConnectorZEVMZetaSentIterator{contract: _ZetaConnectorZEVM.contract, event: "ZetaSent", logs: logs, sub: sub}, nil
}

// WatchZetaSent is a free log subscription operation binding the contract event 0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4.
//
// Solidity: event ZetaSent(address sourceTxOriginAddress, address indexed zetaTxSenderAddress, uint256 indexed destinationChainId, bytes destinationAddress, uint256 zetaValueAndGas, uint256 destinationGasLimit, bytes message, bytes zetaParams)
func (_ZetaConnectorZEVM *ZetaConnectorZEVMFilterer) WatchZetaSent(opts *bind.WatchOpts, sink chan<- *ZetaConnectorZEVMZetaSent, zetaTxSenderAddress []common.Address, destinationChainId []*big.Int) (event.Subscription, error) {

	var zetaTxSenderAddressRule []interface{}
	for _, zetaTxSenderAddressItem := range zetaTxSenderAddress {
		zetaTxSenderAddressRule = append(zetaTxSenderAddressRule, zetaTxSenderAddressItem)
	}
	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	logs, sub, err := _ZetaConnectorZEVM.contract.WatchLogs(opts, "ZetaSent", zetaTxSenderAddressRule, destinationChainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZetaConnectorZEVMZetaSent)
				if err := _ZetaConnectorZEVM.contract.UnpackLog(event, "ZetaSent", log); err != nil {
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

// ParseZetaSent is a log parse operation binding the contract event 0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4.
//
// Solidity: event ZetaSent(address sourceTxOriginAddress, address indexed zetaTxSenderAddress, uint256 indexed destinationChainId, bytes destinationAddress, uint256 zetaValueAndGas, uint256 destinationGasLimit, bytes message, bytes zetaParams)
func (_ZetaConnectorZEVM *ZetaConnectorZEVMFilterer) ParseZetaSent(log types.Log) (*ZetaConnectorZEVMZetaSent, error) {
	event := new(ZetaConnectorZEVMZetaSent)
	if err := _ZetaConnectorZEVM.contract.UnpackLog(event, "ZetaSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
