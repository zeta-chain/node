// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc20

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

// USDTMetaData contains all meta data concerning the USDT contract.
var USDTMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000bb138038062000bb183398101604081905262000034916200014c565b600362000042848262000260565b50600462000051838262000260565b506005805460ff191660ff92909216919091179055505033600090815260208190526040902066038d7ea4c6800090556200032c565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620000af57600080fd5b81516001600160401b0380821115620000cc57620000cc62000087565b604051601f8301601f19908116603f01168101908282118183101715620000f757620000f762000087565b816040528381526020925086838588010111156200011457600080fd5b600091505b8382101562000138578582018301518183018401529082019062000119565b600093810190920192909252949350505050565b6000806000606084860312156200016257600080fd5b83516001600160401b03808211156200017a57600080fd5b62000188878388016200009d565b945060208601519150808211156200019f57600080fd5b50620001ae868287016200009d565b925050604084015160ff81168114620001c657600080fd5b809150509250925092565b600181811c90821680620001e657607f821691505b6020821081036200020757634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200025b57600081815260208120601f850160051c81016020861015620002365750805b601f850160051c820191505b81811015620002575782815560010162000242565b5050505b505050565b81516001600160401b038111156200027c576200027c62000087565b62000294816200028d8454620001d1565b846200020d565b602080601f831160018114620002cc5760008415620002b35750858301515b600019600386901b1c1916600185901b17855562000257565b600085815260208120601f198616915b82811015620002fd57888601518255948401946001909101908401620002dc565b50858210156200031c5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b610875806200033c6000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c806342966c681161007157806342966c681461012957806370a082311461013e57806395d89b4114610167578063a0712d681461016f578063a9059cbb14610182578063dd62ed3e1461019557600080fd5b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100ef57806323b872dd14610101578063313ce56714610114575b600080fd5b6100b66101ce565b6040516100c3919061068b565b60405180910390f35b6100df6100da3660046106f5565b610260565b60405190151581526020016100c3565b6002545b6040519081526020016100c3565b6100df61010f36600461071f565b610277565b60055460405160ff90911681526020016100c3565b61013c61013736600461075b565b61032d565b005b6100f361014c366004610774565b6001600160a01b031660009081526020819052604090205490565b6100b6610354565b61013c61017d36600461075b565b610363565b6100df6101903660046106f5565b610382565b6100f36101a3366004610796565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6060600380546101dd906107c9565b80601f0160208091040260200160405190810160405280929190818152602001828054610209906107c9565b80156102565780601f1061022b57610100808354040283529160200191610256565b820191906000526020600020905b81548152906001019060200180831161023957829003601f168201915b5050505050905090565b600061026d33848461038f565b5060015b92915050565b60006102848484846104b3565b6001600160a01b03841660009081526001602090815260408083203384529091529020548281101561030e5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b60648201526084015b60405180910390fd5b610322853361031d8685610819565b61038f565b506001949350505050565b336000908152602081905260408120805483929061034c908490610819565b909155505050565b6060600480546101dd906107c9565b336000908152602081905260408120805483929061034c90849061082c565b600061026d3384846104b3565b6001600160a01b0383166103f15760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b6064820152608401610305565b6001600160a01b0382166104525760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b6064820152608401610305565b6001600160a01b0383811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b6001600160a01b0383166105175760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f206164604482015264647265737360d81b6064820152608401610305565b6001600160a01b0382166105795760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201526265737360e81b6064820152608401610305565b6001600160a01b038316600090815260208190526040902054818110156105f15760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e7420657863656564732062604482015265616c616e636560d01b6064820152608401610305565b6105fb8282610819565b6001600160a01b03808616600090815260208190526040808220939093559085168152908120805484929061063190849061082c565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161067d91815260200190565b60405180910390a350505050565b600060208083528351808285015260005b818110156106b85785810183015185820160400152820161069c565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b03811681146106f057600080fd5b919050565b6000806040838503121561070857600080fd5b610711836106d9565b946020939093013593505050565b60008060006060848603121561073457600080fd5b61073d846106d9565b925061074b602085016106d9565b9150604084013590509250925092565b60006020828403121561076d57600080fd5b5035919050565b60006020828403121561078657600080fd5b61078f826106d9565b9392505050565b600080604083850312156107a957600080fd5b6107b2836106d9565b91506107c0602084016106d9565b90509250929050565b600181811c908216806107dd57607f821691505b6020821081036107fd57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b8181038181111561027157610271610803565b808201808211156102715761027161080356fea26469706673582212202be4005b8bbb29132ec10f55a7ffcfb97074f1c7504254b1a3048764592668bb64736f6c63430008110033",
}

// USDTABI is the input ABI used to generate the binding from.
// Deprecated: Use USDTMetaData.ABI instead.
var USDTABI = USDTMetaData.ABI

// USDTBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use USDTMetaData.Bin instead.
var USDTBin = USDTMetaData.Bin

// DeployUSDT deploys a new Ethereum contract, binding an instance of USDT to it.
func DeployUSDT(auth *bind.TransactOpts, backend bind.ContractBackend, name_ string, symbol_ string, decimals_ uint8) (common.Address, *types.Transaction, *USDT, error) {
	parsed, err := USDTMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(USDTBin), backend, name_, symbol_, decimals_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &USDT{USDTCaller: USDTCaller{contract: contract}, USDTTransactor: USDTTransactor{contract: contract}, USDTFilterer: USDTFilterer{contract: contract}}, nil
}

// USDT is an auto generated Go binding around an Ethereum contract.
type USDT struct {
	USDTCaller     // Read-only binding to the contract
	USDTTransactor // Write-only binding to the contract
	USDTFilterer   // Log filterer for contract events
}

// USDTCaller is an auto generated read-only Go binding around an Ethereum contract.
type USDTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// USDTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type USDTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// USDTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type USDTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// USDTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type USDTSession struct {
	Contract     *USDT             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// USDTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type USDTCallerSession struct {
	Contract *USDTCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// USDTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type USDTTransactorSession struct {
	Contract     *USDTTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// USDTRaw is an auto generated low-level Go binding around an Ethereum contract.
type USDTRaw struct {
	Contract *USDT // Generic contract binding to access the raw methods on
}

// USDTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type USDTCallerRaw struct {
	Contract *USDTCaller // Generic read-only contract binding to access the raw methods on
}

// USDTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type USDTTransactorRaw struct {
	Contract *USDTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUSDT creates a new instance of USDT, bound to a specific deployed contract.
func NewUSDT(address common.Address, backend bind.ContractBackend) (*USDT, error) {
	contract, err := bindUSDT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &USDT{USDTCaller: USDTCaller{contract: contract}, USDTTransactor: USDTTransactor{contract: contract}, USDTFilterer: USDTFilterer{contract: contract}}, nil
}

// NewUSDTCaller creates a new read-only instance of USDT, bound to a specific deployed contract.
func NewUSDTCaller(address common.Address, caller bind.ContractCaller) (*USDTCaller, error) {
	contract, err := bindUSDT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &USDTCaller{contract: contract}, nil
}

// NewUSDTTransactor creates a new write-only instance of USDT, bound to a specific deployed contract.
func NewUSDTTransactor(address common.Address, transactor bind.ContractTransactor) (*USDTTransactor, error) {
	contract, err := bindUSDT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &USDTTransactor{contract: contract}, nil
}

// NewUSDTFilterer creates a new log filterer instance of USDT, bound to a specific deployed contract.
func NewUSDTFilterer(address common.Address, filterer bind.ContractFilterer) (*USDTFilterer, error) {
	contract, err := bindUSDT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &USDTFilterer{contract: contract}, nil
}

// bindUSDT binds a generic wrapper to an already deployed contract.
func bindUSDT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(USDTABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_USDT *USDTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDT.Contract.USDTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_USDT *USDTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDT.Contract.USDTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_USDT *USDTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDT.Contract.USDTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_USDT *USDTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_USDT *USDTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_USDT *USDTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDT.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_USDT *USDTCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _USDT.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_USDT *USDTSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _USDT.Contract.Allowance(&_USDT.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_USDT *USDTCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _USDT.Contract.Allowance(&_USDT.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_USDT *USDTCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _USDT.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_USDT *USDTSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _USDT.Contract.BalanceOf(&_USDT.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_USDT *USDTCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _USDT.Contract.BalanceOf(&_USDT.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_USDT *USDTCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _USDT.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_USDT *USDTSession) Decimals() (uint8, error) {
	return _USDT.Contract.Decimals(&_USDT.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_USDT *USDTCallerSession) Decimals() (uint8, error) {
	return _USDT.Contract.Decimals(&_USDT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_USDT *USDTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _USDT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_USDT *USDTSession) Name() (string, error) {
	return _USDT.Contract.Name(&_USDT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_USDT *USDTCallerSession) Name() (string, error) {
	return _USDT.Contract.Name(&_USDT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_USDT *USDTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _USDT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_USDT *USDTSession) Symbol() (string, error) {
	return _USDT.Contract.Symbol(&_USDT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_USDT *USDTCallerSession) Symbol() (string, error) {
	return _USDT.Contract.Symbol(&_USDT.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_USDT *USDTCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _USDT.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_USDT *USDTSession) TotalSupply() (*big.Int, error) {
	return _USDT.Contract.TotalSupply(&_USDT.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_USDT *USDTCallerSession) TotalSupply() (*big.Int, error) {
	return _USDT.Contract.TotalSupply(&_USDT.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_USDT *USDTTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _USDT.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_USDT *USDTSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _USDT.Contract.Approve(&_USDT.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_USDT *USDTTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _USDT.Contract.Approve(&_USDT.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_USDT *USDTTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _USDT.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_USDT *USDTSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _USDT.Contract.Burn(&_USDT.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_USDT *USDTTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _USDT.Contract.Burn(&_USDT.TransactOpts, amount)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 amount) returns()
func (_USDT *USDTTransactor) Mint(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _USDT.contract.Transact(opts, "mint", amount)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 amount) returns()
func (_USDT *USDTSession) Mint(amount *big.Int) (*types.Transaction, error) {
	return _USDT.Contract.Mint(&_USDT.TransactOpts, amount)
}

// Mint is a paid mutator transaction binding the contract method 0xa0712d68.
//
// Solidity: function mint(uint256 amount) returns()
func (_USDT *USDTTransactorSession) Mint(amount *big.Int) (*types.Transaction, error) {
	return _USDT.Contract.Mint(&_USDT.TransactOpts, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_USDT *USDTTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _USDT.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_USDT *USDTSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _USDT.Contract.Transfer(&_USDT.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_USDT *USDTTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _USDT.Contract.Transfer(&_USDT.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_USDT *USDTTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _USDT.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_USDT *USDTSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _USDT.Contract.TransferFrom(&_USDT.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_USDT *USDTTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _USDT.Contract.TransferFrom(&_USDT.TransactOpts, sender, recipient, amount)
}

// USDTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the USDT contract.
type USDTApprovalIterator struct {
	Event *USDTApproval // Event containing the contract specifics and raw log

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
func (it *USDTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDTApproval)
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
		it.Event = new(USDTApproval)
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
func (it *USDTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *USDTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// USDTApproval represents a Approval event raised by the USDT contract.
type USDTApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_USDT *USDTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*USDTApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _USDT.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &USDTApprovalIterator{contract: _USDT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_USDT *USDTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *USDTApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _USDT.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(USDTApproval)
				if err := _USDT.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_USDT *USDTFilterer) ParseApproval(log types.Log) (*USDTApproval, error) {
	event := new(USDTApproval)
	if err := _USDT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// USDTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the USDT contract.
type USDTTransferIterator struct {
	Event *USDTTransfer // Event containing the contract specifics and raw log

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
func (it *USDTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDTTransfer)
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
		it.Event = new(USDTTransfer)
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
func (it *USDTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *USDTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// USDTTransfer represents a Transfer event raised by the USDT contract.
type USDTTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_USDT *USDTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*USDTTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _USDT.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &USDTTransferIterator{contract: _USDT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_USDT *USDTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *USDTTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _USDT.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(USDTTransfer)
				if err := _USDT.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_USDT *USDTFilterer) ParseTransfer(log types.Log) (*USDTTransfer, error) {
	event := new(USDTTransfer)
	if err := _USDT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
