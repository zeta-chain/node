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
	Bin: "0x608060405234801561001057600080fd5b50610756806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063e3be6f6814610030575b600080fd5b61004a600480360381019061004591906103ff565b61004c565b005b6000808473ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b815260040160408051808303816000875af115801561009b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100bf91906104c8565b915091508173ffffffffffffffffffffffffffffffffffffffff166323b872dd33306001876100ee9190610537565b856100f9919061056b565b6040518463ffffffff1660e01b8152600401610117939291906105cb565b600060405180830381600087803b15801561013157600080fd5b505af1158015610145573d6000803e3d6000fd5b505050508173ffffffffffffffffffffffffffffffffffffffff1663095ea7b3866001866101739190610537565b8461017e919061056b565b6040518363ffffffff1660e01b815260040161019b929190610602565b600060405180830381600087803b1580156101b557600080fd5b505af11580156101c9573d6000803e3d6000fd5b505050508473ffffffffffffffffffffffffffffffffffffffff166323b872dd333086886101f7919061056b565b6040518463ffffffff1660e01b8152600401610215939291906105cb565b600060405180830381600087803b15801561022f57600080fd5b505af1158015610243573d6000803e3d6000fd5b5050505060005b838110156102e0578573ffffffffffffffffffffffffffffffffffffffff1663c70126268989886040518463ffffffff1660e01b815260040161028f93929190610689565b6020604051808303816000875af11580156102ae573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102d291906106f3565b50808060010191505061024a565b5050505050505050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f840112610319576103186102f4565b5b8235905067ffffffffffffffff811115610336576103356102f9565b5b602083019150836001820283011115610352576103516102fe565b5b9250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061038482610359565b9050919050565b600061039682610379565b9050919050565b6103a68161038b565b81146103b157600080fd5b50565b6000813590506103c38161039d565b92915050565b6000819050919050565b6103dc816103c9565b81146103e757600080fd5b50565b6000813590506103f9816103d3565b92915050565b60008060008060006080868803121561041b5761041a6102ea565b5b600086013567ffffffffffffffff811115610439576104386102ef565b5b61044588828901610303565b95509550506020610458888289016103b4565b9350506040610469888289016103ea565b925050606061047a888289016103ea565b9150509295509295909350565b61049081610379565b811461049b57600080fd5b50565b6000815190506104ad81610487565b92915050565b6000815190506104c2816103d3565b92915050565b600080604083850312156104df576104de6102ea565b5b60006104ed8582860161049e565b92505060206104fe858286016104b3565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610542826103c9565b915061054d836103c9565b925082820190508082111561056557610564610508565b5b92915050565b6000610576826103c9565b9150610581836103c9565b925082820261058f816103c9565b915082820484148315176105a6576105a5610508565b5b5092915050565b6105b681610379565b82525050565b6105c5816103c9565b82525050565b60006060820190506105e060008301866105ad565b6105ed60208301856105ad565b6105fa60408301846105bc565b949350505050565b600060408201905061061760008301856105ad565b61062460208301846105bc565b9392505050565b600082825260208201905092915050565b82818337600083830152505050565b6000601f19601f8301169050919050565b6000610668838561062b565b935061067583858461063c565b61067e8361064b565b840190509392505050565b600060408201905081810360008301526106a481858761065c565b90506106b360208301846105bc565b949350505050565b60008115159050919050565b6106d0816106bb565b81146106db57600080fd5b50565b6000815190506106ed816106c7565b92915050565b600060208284031215610709576107086102ea565b5b6000610717848285016106de565b9150509291505056fea26469706673582212204ea5580fd884f66b08516c0cea46d9e71553f966a77e2954832ccf2d9abe611b64736f6c63430008170033",
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
