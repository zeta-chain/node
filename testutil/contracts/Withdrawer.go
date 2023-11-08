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

// WithdrawerMetaData contains all meta data concerning the Withdrawer contract.
var WithdrawerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"},{\"internalType\":\"contractIZRC20\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"runWithdraws\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061075a806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063e3be6f6814610030575b600080fd5b61004a600480360381019061004591906103ef565b61004c565b005b6000808473ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b815260040160408051808303816000875af115801561009b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100bf91906104b8565b915091508173ffffffffffffffffffffffffffffffffffffffff166323b872dd3330600a856100ee9190610527565b6040518463ffffffff1660e01b815260040161010c93929190610587565b600060405180830381600087803b15801561012657600080fd5b505af115801561013a573d6000803e3d6000fd5b505050508173ffffffffffffffffffffffffffffffffffffffff1663095ea7b386600a846101689190610527565b6040518363ffffffff1660e01b81526004016101859291906105be565b600060405180830381600087803b15801561019f57600080fd5b505af11580156101b3573d6000803e3d6000fd5b505050508473ffffffffffffffffffffffffffffffffffffffff166323b872dd333086886101e19190610527565b6040518463ffffffff1660e01b81526004016101ff93929190610587565b600060405180830381600087803b15801561021957600080fd5b505af115801561022d573d6000803e3d6000fd5b5050505060005b838110156102d0578573ffffffffffffffffffffffffffffffffffffffff1663c70126268989886040518463ffffffff1660e01b815260040161027993929190610645565b6020604051808303816000875af1158015610298573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102bc91906106af565b5080806102c8906106dc565b915050610234565b5050505050505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f840112610309576103086102e4565b5b8235905067ffffffffffffffff811115610326576103256102e9565b5b602083019150836001820283011115610342576103416102ee565b5b9250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061037482610349565b9050919050565b600061038682610369565b9050919050565b6103968161037b565b81146103a157600080fd5b50565b6000813590506103b38161038d565b92915050565b6000819050919050565b6103cc816103b9565b81146103d757600080fd5b50565b6000813590506103e9816103c3565b92915050565b60008060008060006080868803121561040b5761040a6102da565b5b600086013567ffffffffffffffff811115610429576104286102df565b5b610435888289016102f3565b95509550506020610448888289016103a4565b9350506040610459888289016103da565b925050606061046a888289016103da565b9150509295509295909350565b61048081610369565b811461048b57600080fd5b50565b60008151905061049d81610477565b92915050565b6000815190506104b2816103c3565b92915050565b600080604083850312156104cf576104ce6102da565b5b60006104dd8582860161048e565b92505060206104ee858286016104a3565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610532826103b9565b915061053d836103b9565b925082820261054b816103b9565b91508282048414831517610562576105616104f8565b5b5092915050565b61057281610369565b82525050565b610581816103b9565b82525050565b600060608201905061059c6000830186610569565b6105a96020830185610569565b6105b66040830184610578565b949350505050565b60006040820190506105d36000830185610569565b6105e06020830184610578565b9392505050565b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b600061062483856105e7565b93506106318385846105f8565b61063a83610607565b840190509392505050565b60006040820190508181036000830152610660818587610618565b905061066f6020830184610578565b949350505050565b60008115159050919050565b61068c81610677565b811461069757600080fd5b50565b6000815190506106a981610683565b92915050565b6000602082840312156106c5576106c46102da565b5b60006106d38482850161069a565b91505092915050565b60006106e7826103b9565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610719576107186104f8565b5b60018201905091905056fea2646970667358221220facf63a05707c7ee9b4b53c5ddd9a2d42e779dd5790d95ba8e4526d32232b9d664736f6c63430008150033",
}

// WithdrawerABI is the input ABI used to generate the binding from.
// Deprecated: Use WithdrawerMetaData.ABI instead.
var WithdrawerABI = WithdrawerMetaData.ABI

// WithdrawerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use WithdrawerMetaData.Bin instead.
var WithdrawerBin = WithdrawerMetaData.Bin

// DeployWithdrawer deploys a new Ethereum contract, binding an instance of Withdrawer to it.
func DeployWithdrawer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Withdrawer, error) {
	parsed, err := WithdrawerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WithdrawerBin), backend)
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

// RunWithdraws is a paid mutator transaction binding the contract method 0xe3be6f68.
//
// Solidity: function runWithdraws(bytes recipient, address asset, uint256 amount, uint256 count) returns()
func (_Withdrawer *WithdrawerTransactor) RunWithdraws(opts *bind.TransactOpts, recipient []byte, asset common.Address, amount *big.Int, count *big.Int) (*types.Transaction, error) {
	return _Withdrawer.contract.Transact(opts, "runWithdraws", recipient, asset, amount, count)
}

// RunWithdraws is a paid mutator transaction binding the contract method 0xe3be6f68.
//
// Solidity: function runWithdraws(bytes recipient, address asset, uint256 amount, uint256 count) returns()
func (_Withdrawer *WithdrawerSession) RunWithdraws(recipient []byte, asset common.Address, amount *big.Int, count *big.Int) (*types.Transaction, error) {
	return _Withdrawer.Contract.RunWithdraws(&_Withdrawer.TransactOpts, recipient, asset, amount, count)
}

// RunWithdraws is a paid mutator transaction binding the contract method 0xe3be6f68.
//
// Solidity: function runWithdraws(bytes recipient, address asset, uint256 amount, uint256 count) returns()
func (_Withdrawer *WithdrawerTransactorSession) RunWithdraws(recipient []byte, asset common.Address, amount *big.Int, count *big.Int) (*types.Transaction, error) {
	return _Withdrawer.Contract.RunWithdraws(&_Withdrawer.TransactOpts, recipient, asset, amount, count)
}
