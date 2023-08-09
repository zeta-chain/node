// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contextapp

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

// Context is an auto generated low-level Go binding around an user-defined struct.
type Context struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}

// ContextAppMetaData contains all meta data concerning the ContextApp contract.
var ContextAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"msgSender\",\"type\":\"address\"}],\"name\":\"ContextData\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610420806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063de43156e14610030575b600080fd5b61004a60048036038101906100459190610182565b61004c565b005b7fcde88c82509f7dbeaae2782de64879ac731556f65d4474e9afc4ea01cca4498885806000019061007d91906102bf565b8760200160208101906100909190610155565b8860400135336040516100a7959493929190610271565b60405180910390a15050505050565b6000813590506100c5816103bc565b92915050565b60008083601f8401126100e1576100e0610383565b5b8235905067ffffffffffffffff8111156100fe576100fd61037e565b5b60208301915083600182028301111561011a57610119610397565b5b9250929050565b6000606082840312156101375761013661038d565b5b81905092915050565b60008135905061014f816103d3565b92915050565b60006020828403121561016b5761016a6103a6565b5b6000610179848285016100b6565b91505092915050565b60008060008060006080868803121561019e5761019d6103a6565b5b600086013567ffffffffffffffff8111156101bc576101bb6103a1565b5b6101c888828901610121565b95505060206101d9888289016100b6565b94505060406101ea88828901610140565b935050606086013567ffffffffffffffff81111561020b5761020a6103a1565b5b610217888289016100cb565b92509250509295509295909350565b61022f81610333565b82525050565b60006102418385610322565b935061024e83858461036f565b610257836103ab565b840190509392505050565b61026b81610365565b82525050565b6000608082019050818103600083015261028c818789610235565b905061029b6020830186610226565b6102a86040830185610262565b6102b56060830184610226565b9695505050505050565b600080833560016020038436030381126102dc576102db610392565b5b80840192508235915067ffffffffffffffff8211156102fe576102fd610388565b5b60208301925060018202360383131561031a5761031961039c565b5b509250929050565b600082825260208201905092915050565b600061033e82610345565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b6103c581610333565b81146103d057600080fd5b50565b6103dc81610365565b81146103e757600080fd5b5056fea26469706673582212203bf8ff0cc81cde452b74dd59e10fc546b8344cb0b268971059b45bc40191115d64736f6c63430008070033",
}

// ContextAppABI is the input ABI used to generate the binding from.
// Deprecated: Use ContextAppMetaData.ABI instead.
var ContextAppABI = ContextAppMetaData.ABI

// ContextAppBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContextAppMetaData.Bin instead.
var ContextAppBin = ContextAppMetaData.Bin

// DeployContextApp deploys a new Ethereum contract, binding an instance of ContextApp to it.
func DeployContextApp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ContextApp, error) {
	parsed, err := ContextAppMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContextAppBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContextApp{ContextAppCaller: ContextAppCaller{contract: contract}, ContextAppTransactor: ContextAppTransactor{contract: contract}, ContextAppFilterer: ContextAppFilterer{contract: contract}}, nil
}

// ContextApp is an auto generated Go binding around an Ethereum contract.
type ContextApp struct {
	ContextAppCaller     // Read-only binding to the contract
	ContextAppTransactor // Write-only binding to the contract
	ContextAppFilterer   // Log filterer for contract events
}

// ContextAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContextAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContextAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContextAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContextAppSession struct {
	Contract     *ContextApp       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContextAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContextAppCallerSession struct {
	Contract *ContextAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ContextAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContextAppTransactorSession struct {
	Contract     *ContextAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ContextAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContextAppRaw struct {
	Contract *ContextApp // Generic contract binding to access the raw methods on
}

// ContextAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContextAppCallerRaw struct {
	Contract *ContextAppCaller // Generic read-only contract binding to access the raw methods on
}

// ContextAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContextAppTransactorRaw struct {
	Contract *ContextAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContextApp creates a new instance of ContextApp, bound to a specific deployed contract.
func NewContextApp(address common.Address, backend bind.ContractBackend) (*ContextApp, error) {
	contract, err := bindContextApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContextApp{ContextAppCaller: ContextAppCaller{contract: contract}, ContextAppTransactor: ContextAppTransactor{contract: contract}, ContextAppFilterer: ContextAppFilterer{contract: contract}}, nil
}

// NewContextAppCaller creates a new read-only instance of ContextApp, bound to a specific deployed contract.
func NewContextAppCaller(address common.Address, caller bind.ContractCaller) (*ContextAppCaller, error) {
	contract, err := bindContextApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContextAppCaller{contract: contract}, nil
}

// NewContextAppTransactor creates a new write-only instance of ContextApp, bound to a specific deployed contract.
func NewContextAppTransactor(address common.Address, transactor bind.ContractTransactor) (*ContextAppTransactor, error) {
	contract, err := bindContextApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContextAppTransactor{contract: contract}, nil
}

// NewContextAppFilterer creates a new log filterer instance of ContextApp, bound to a specific deployed contract.
func NewContextAppFilterer(address common.Address, filterer bind.ContractFilterer) (*ContextAppFilterer, error) {
	contract, err := bindContextApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContextAppFilterer{contract: contract}, nil
}

// bindContextApp binds a generic wrapper to an already deployed contract.
func bindContextApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContextAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContextApp *ContextAppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContextApp.Contract.ContextAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContextApp *ContextAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContextApp.Contract.ContextAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContextApp *ContextAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContextApp.Contract.ContextAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContextApp *ContextAppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContextApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContextApp *ContextAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContextApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContextApp *ContextAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContextApp.Contract.contract.Transact(opts, method, params...)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_ContextApp *ContextAppTransactor) OnCrossChainCall(opts *bind.TransactOpts, context Context, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ContextApp.contract.Transact(opts, "onCrossChainCall", context, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_ContextApp *ContextAppSession) OnCrossChainCall(context Context, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ContextApp.Contract.OnCrossChainCall(&_ContextApp.TransactOpts, context, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_ContextApp *ContextAppTransactorSession) OnCrossChainCall(context Context, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ContextApp.Contract.OnCrossChainCall(&_ContextApp.TransactOpts, context, zrc20, amount, message)
}

// ContextAppContextDataIterator is returned from FilterContextData and is used to iterate over the raw logs and unpacked data for ContextData events raised by the ContextApp contract.
type ContextAppContextDataIterator struct {
	Event *ContextAppContextData // Event containing the contract specifics and raw log

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
func (it *ContextAppContextDataIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContextAppContextData)
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
		it.Event = new(ContextAppContextData)
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
func (it *ContextAppContextDataIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContextAppContextDataIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContextAppContextData represents a ContextData event raised by the ContextApp contract.
type ContextAppContextData struct {
	Origin    []byte
	Sender    common.Address
	ChainID   *big.Int
	MsgSender common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterContextData is a free log retrieval operation binding the contract event 0xcde88c82509f7dbeaae2782de64879ac731556f65d4474e9afc4ea01cca44988.
//
// Solidity: event ContextData(bytes origin, address sender, uint256 chainID, address msgSender)
func (_ContextApp *ContextAppFilterer) FilterContextData(opts *bind.FilterOpts) (*ContextAppContextDataIterator, error) {

	logs, sub, err := _ContextApp.contract.FilterLogs(opts, "ContextData")
	if err != nil {
		return nil, err
	}
	return &ContextAppContextDataIterator{contract: _ContextApp.contract, event: "ContextData", logs: logs, sub: sub}, nil
}

// WatchContextData is a free log subscription operation binding the contract event 0xcde88c82509f7dbeaae2782de64879ac731556f65d4474e9afc4ea01cca44988.
//
// Solidity: event ContextData(bytes origin, address sender, uint256 chainID, address msgSender)
func (_ContextApp *ContextAppFilterer) WatchContextData(opts *bind.WatchOpts, sink chan<- *ContextAppContextData) (event.Subscription, error) {

	logs, sub, err := _ContextApp.contract.WatchLogs(opts, "ContextData")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContextAppContextData)
				if err := _ContextApp.contract.UnpackLog(event, "ContextData", log); err != nil {
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

// ParseContextData is a log parse operation binding the contract event 0xcde88c82509f7dbeaae2782de64879ac731556f65d4474e9afc4ea01cca44988.
//
// Solidity: event ContextData(bytes origin, address sender, uint256 chainID, address msgSender)
func (_ContextApp *ContextAppFilterer) ParseContextData(log types.Log) (*ContextAppContextData, error) {
	event := new(ContextAppContextData)
	if err := _ContextApp.contract.UnpackLog(event, "ContextData", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
