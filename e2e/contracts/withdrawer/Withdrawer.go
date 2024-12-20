// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package withdrawer

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

// WithdrawerMetaData contains all meta data concerning the Withdrawer contract.
var WithdrawerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_withdrawAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a0604052348015600f57600080fd5b506040516106733803806106738339818101604052810190602f91906072565b806080818152505050609a565b600080fd5b6000819050919050565b6052816041565b8114605c57600080fd5b50565b600081519050606c81604b565b92915050565b6000602082840312156085576084603c565b5b6000609184828501605f565b91505092915050565b6080516105b86100bb600039600081816077015261016601526105b86000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063534844a21461003b578063de43156e14610059575b600080fd5b610043610075565b6040516100509190610207565b60405180910390f35b610073600480360381019061006e919061033f565b610099565b005b7f000000000000000000000000000000000000000000000000000000000000000081565b8373ffffffffffffffffffffffffffffffffffffffff1663095ea7b3857fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6040518363ffffffff1660e01b81526004016100f49291906103f2565b6020604051808303816000875af1158015610113573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101379190610453565b508373ffffffffffffffffffffffffffffffffffffffff1663c7012626868060000190610164919061048f565b7f00000000000000000000000000000000000000000000000000000000000000006040518463ffffffff1660e01b81526004016101a393929190610550565b6020604051808303816000875af11580156101c2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101e69190610453565b505050505050565b6000819050919050565b610201816101ee565b82525050565b600060208201905061021c60008301846101f8565b92915050565b600080fd5b600080fd5b600080fd5b6000606082840312156102475761024661022c565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061027b82610250565b9050919050565b61028b81610270565b811461029657600080fd5b50565b6000813590506102a881610282565b92915050565b6102b7816101ee565b81146102c257600080fd5b50565b6000813590506102d4816102ae565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126102ff576102fe6102da565b5b8235905067ffffffffffffffff81111561031c5761031b6102df565b5b602083019150836001820283011115610338576103376102e4565b5b9250929050565b60008060008060006080868803121561035b5761035a610222565b5b600086013567ffffffffffffffff81111561037957610378610227565b5b61038588828901610231565b955050602061039688828901610299565b94505060406103a7888289016102c5565b935050606086013567ffffffffffffffff8111156103c8576103c7610227565b5b6103d4888289016102e9565b92509250509295509295909350565b6103ec81610270565b82525050565b600060408201905061040760008301856103e3565b61041460208301846101f8565b9392505050565b60008115159050919050565b6104308161041b565b811461043b57600080fd5b50565b60008151905061044d81610427565b92915050565b60006020828403121561046957610468610222565b5b60006104778482850161043e565b91505092915050565b600080fd5b600080fd5b600080fd5b600080833560016020038436030381126104ac576104ab610480565b5b80840192508235915067ffffffffffffffff8211156104ce576104cd610485565b5b6020830192506001820236038313156104ea576104e961048a565b5b509250929050565b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b600061052f83856104f2565b935061053c838584610503565b61054583610512565b840190509392505050565b6000604082019050818103600083015261056b818587610523565b905061057a60208301846101f8565b94935050505056fea26469706673582212208d3db17c081600e0a0203ac712a53a290f24ef219417697847e27c563ac3d09a64736f6c634300081a0033",
}

// WithdrawerABI is the input ABI used to generate the binding from.
// Deprecated: Use WithdrawerMetaData.ABI instead.
var WithdrawerABI = WithdrawerMetaData.ABI

// WithdrawerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use WithdrawerMetaData.Bin instead.
var WithdrawerBin = WithdrawerMetaData.Bin

// DeployWithdrawer deploys a new Ethereum contract, binding an instance of Withdrawer to it.
func DeployWithdrawer(auth *bind.TransactOpts, backend bind.ContractBackend, _withdrawAmount *big.Int) (common.Address, *types.Transaction, *Withdrawer, error) {
	parsed, err := WithdrawerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WithdrawerBin), backend, _withdrawAmount)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Withdrawer{WithdrawerCaller: WithdrawerCaller{contract: contract}, WithdrawerTransactor: WithdrawerTransactor{contract: contract}, WithdrawerFilterer: WithdrawerFilterer{contract: contract}}, nil
}

// Withdrawer is an auto generated Go binding around an Ethereum contract.
type Withdrawer struct {
	WithdrawerCaller     // Read-only binding to the contract
	WithdrawerTransactor // Write-only binding to the contract
	WithdrawerFilterer   // Log filterer for contract events
}

// WithdrawerCaller is an auto generated read-only Go binding around an Ethereum contract.
type WithdrawerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WithdrawerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WithdrawerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WithdrawerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WithdrawerSession struct {
	Contract     *Withdrawer       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WithdrawerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WithdrawerCallerSession struct {
	Contract *WithdrawerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// WithdrawerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WithdrawerTransactorSession struct {
	Contract     *WithdrawerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// WithdrawerRaw is an auto generated low-level Go binding around an Ethereum contract.
type WithdrawerRaw struct {
	Contract *Withdrawer // Generic contract binding to access the raw methods on
}

// WithdrawerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WithdrawerCallerRaw struct {
	Contract *WithdrawerCaller // Generic read-only contract binding to access the raw methods on
}

// WithdrawerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WithdrawerTransactorRaw struct {
	Contract *WithdrawerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWithdrawer creates a new instance of Withdrawer, bound to a specific deployed contract.
func NewWithdrawer(address common.Address, backend bind.ContractBackend) (*Withdrawer, error) {
	contract, err := bindWithdrawer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Withdrawer{WithdrawerCaller: WithdrawerCaller{contract: contract}, WithdrawerTransactor: WithdrawerTransactor{contract: contract}, WithdrawerFilterer: WithdrawerFilterer{contract: contract}}, nil
}

// NewWithdrawerCaller creates a new read-only instance of Withdrawer, bound to a specific deployed contract.
func NewWithdrawerCaller(address common.Address, caller bind.ContractCaller) (*WithdrawerCaller, error) {
	contract, err := bindWithdrawer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WithdrawerCaller{contract: contract}, nil
}

// NewWithdrawerTransactor creates a new write-only instance of Withdrawer, bound to a specific deployed contract.
func NewWithdrawerTransactor(address common.Address, transactor bind.ContractTransactor) (*WithdrawerTransactor, error) {
	contract, err := bindWithdrawer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WithdrawerTransactor{contract: contract}, nil
}

// NewWithdrawerFilterer creates a new log filterer instance of Withdrawer, bound to a specific deployed contract.
func NewWithdrawerFilterer(address common.Address, filterer bind.ContractFilterer) (*WithdrawerFilterer, error) {
	contract, err := bindWithdrawer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WithdrawerFilterer{contract: contract}, nil
}

// bindWithdrawer binds a generic wrapper to an already deployed contract.
func bindWithdrawer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WithdrawerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Withdrawer *WithdrawerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Withdrawer.Contract.WithdrawerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Withdrawer *WithdrawerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Withdrawer.Contract.WithdrawerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Withdrawer *WithdrawerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Withdrawer.Contract.WithdrawerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Withdrawer *WithdrawerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Withdrawer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Withdrawer *WithdrawerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Withdrawer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Withdrawer *WithdrawerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Withdrawer.Contract.contract.Transact(opts, method, params...)
}

// WithdrawAmount is a free data retrieval call binding the contract method 0x534844a2.
//
// Solidity: function withdrawAmount() view returns(uint256)
func (_Withdrawer *WithdrawerCaller) WithdrawAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Withdrawer.contract.Call(opts, &out, "withdrawAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawAmount is a free data retrieval call binding the contract method 0x534844a2.
//
// Solidity: function withdrawAmount() view returns(uint256)
func (_Withdrawer *WithdrawerSession) WithdrawAmount() (*big.Int, error) {
	return _Withdrawer.Contract.WithdrawAmount(&_Withdrawer.CallOpts)
}

// WithdrawAmount is a free data retrieval call binding the contract method 0x534844a2.
//
// Solidity: function withdrawAmount() view returns(uint256)
func (_Withdrawer *WithdrawerCallerSession) WithdrawAmount() (*big.Int, error) {
	return _Withdrawer.Contract.WithdrawAmount(&_Withdrawer.CallOpts)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 , bytes ) returns()
func (_Withdrawer *WithdrawerTransactor) OnCrossChainCall(opts *bind.TransactOpts, context Context, zrc20 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _Withdrawer.contract.Transact(opts, "onCrossChainCall", context, zrc20, arg2, arg3)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 , bytes ) returns()
func (_Withdrawer *WithdrawerSession) OnCrossChainCall(context Context, zrc20 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _Withdrawer.Contract.OnCrossChainCall(&_Withdrawer.TransactOpts, context, zrc20, arg2, arg3)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 , bytes ) returns()
func (_Withdrawer *WithdrawerTransactorSession) OnCrossChainCall(context Context, zrc20 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _Withdrawer.Contract.OnCrossChainCall(&_Withdrawer.TransactOpts, context, zrc20, arg2, arg3)
}
