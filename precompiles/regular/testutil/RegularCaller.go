// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testutils

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

// RegularCallerMetaData contains all meta data concerning the RegularCaller contract.
var RegularCallerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"testBech32ToHexAddr\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testBech32ify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"method\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"testRegularCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610918806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806342875b1d1461004657806348fc7db114610064578063d7c9bd0214610094575b600080fd5b61004e6100b2565b60405161005b9190610675565b60405180910390f35b61007e6004803603810190610079919061053e565b6101ae565b60405161008b91906106e2565b60405180910390f35b61009c61024b565b6040516100a99190610675565b60405180910390f35b6000806040518060600160405280602b81526020016108b8602b91399050600073b9dbc229bf588a613c00bee8e662727ab8121cfe90506000606573ffffffffffffffffffffffffffffffffffffffff1663e4e2a4ec846040518263ffffffff1660e01b81526004016101259190610690565b60206040518083038186803b15801561013d57600080fd5b505afa158015610151573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061017591906104c8565b90508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614935050505090565b600080606573ffffffffffffffffffffffffffffffffffffffff166393e3663d85856040518363ffffffff1660e01b81526004016101ed9291906106b2565b602060405180830381600087803b15801561020757600080fd5b505af115801561021b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061023f919061059a565b90508091505092915050565b6000806040518060400160405280600481526020017f7a657461000000000000000000000000000000000000000000000000000000008152509050600073b9dbc229bf588a613c00bee8e662727ab8121cfe905060006040518060600160405280602b81526020016108b8602b913990506000606573ffffffffffffffffffffffffffffffffffffffff16630615b74e85856040518363ffffffff1660e01b81526004016102fa9291906106b2565b60006040518083038186803b15801561031257600080fd5b505afa158015610326573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525081019061034f91906104f5565b905081604051602001610362919061065e565b6040516020818303038152906040528051906020012081604051602001610389919061065e565b604051602081830303815290604052805190602001201494505050505090565b60006103bc6103b784610722565b6106fd565b9050828152602081018484840111156103d8576103d7610869565b5b6103e38482856107c2565b509392505050565b60006103fe6103f984610722565b6106fd565b90508281526020810184848401111561041a57610419610869565b5b6104258482856107d1565b509392505050565b60008135905061043c81610889565b92915050565b60008151905061045181610889565b92915050565b600082601f83011261046c5761046b610864565b5b813561047c8482602086016103a9565b91505092915050565b600082601f83011261049a57610499610864565b5b81516104aa8482602086016103eb565b91505092915050565b6000815190506104c2816108a0565b92915050565b6000602082840312156104de576104dd610873565b5b60006104ec84828501610442565b91505092915050565b60006020828403121561050b5761050a610873565b5b600082015167ffffffffffffffff8111156105295761052861086e565b5b61053584828501610485565b91505092915050565b6000806040838503121561055557610554610873565b5b600083013567ffffffffffffffff8111156105735761057261086e565b5b61057f85828601610457565b92505060206105908582860161042d565b9150509250929050565b6000602082840312156105b0576105af610873565b5b60006105be848285016104b3565b91505092915050565b6105d08161077a565b82525050565b6105df8161078c565b82525050565b60006105f082610753565b6105fa818561075e565b935061060a8185602086016107d1565b61061381610878565b840191505092915050565b600061062982610753565b610633818561076f565b93506106438185602086016107d1565b80840191505092915050565b610658816107b8565b82525050565b600061066a828461061e565b915081905092915050565b600060208201905061068a60008301846105d6565b92915050565b600060208201905081810360008301526106aa81846105e5565b905092915050565b600060408201905081810360008301526106cc81856105e5565b90506106db60208301846105c7565b9392505050565b60006020820190506106f7600083018461064f565b92915050565b6000610707610718565b90506107138282610804565b919050565b6000604051905090565b600067ffffffffffffffff82111561073d5761073c610835565b5b61074682610878565b9050602081019050919050565b600081519050919050565b600082825260208201905092915050565b600081905092915050565b600061078582610798565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b838110156107ef5780820151818401526020810190506107d4565b838111156107fe576000848401525b50505050565b61080d82610878565b810181811067ffffffffffffffff8211171561082c5761082b610835565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b6108928161077a565b811461089d57600080fd5b50565b6108a9816107b8565b81146108b457600080fd5b5056fe7a65746131683864757932646c747a39787a307171686d357776636e6a303275707938383766796e343375a26469706673582212205d1317177480909fff852a3871dc5d7fbdcb2fec3c47b91fc472fa775517c70264736f6c63430008070033",
}

// RegularCallerABI is the input ABI used to generate the binding from.
// Deprecated: Use RegularCallerMetaData.ABI instead.
var RegularCallerABI = RegularCallerMetaData.ABI

// RegularCallerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RegularCallerMetaData.Bin instead.
var RegularCallerBin = RegularCallerMetaData.Bin

// DeployRegularCaller deploys a new Ethereum contract, binding an instance of RegularCaller to it.
func DeployRegularCaller(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *RegularCaller, error) {
	parsed, err := RegularCallerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RegularCallerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RegularCaller{RegularCallerCaller: RegularCallerCaller{contract: contract}, RegularCallerTransactor: RegularCallerTransactor{contract: contract}, RegularCallerFilterer: RegularCallerFilterer{contract: contract}}, nil
}

// RegularCaller is an auto generated Go binding around an Ethereum contract.
type RegularCaller struct {
	RegularCallerCaller     // Read-only binding to the contract
	RegularCallerTransactor // Write-only binding to the contract
	RegularCallerFilterer   // Log filterer for contract events
}

// RegularCallerCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegularCallerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegularCallerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegularCallerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegularCallerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegularCallerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegularCallerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegularCallerSession struct {
	Contract     *RegularCaller    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegularCallerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegularCallerCallerSession struct {
	Contract *RegularCallerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// RegularCallerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegularCallerTransactorSession struct {
	Contract     *RegularCallerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// RegularCallerRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegularCallerRaw struct {
	Contract *RegularCaller // Generic contract binding to access the raw methods on
}

// RegularCallerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegularCallerCallerRaw struct {
	Contract *RegularCallerCaller // Generic read-only contract binding to access the raw methods on
}

// RegularCallerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegularCallerTransactorRaw struct {
	Contract *RegularCallerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegularCaller creates a new instance of RegularCaller, bound to a specific deployed contract.
func NewRegularCaller(address common.Address, backend bind.ContractBackend) (*RegularCaller, error) {
	contract, err := bindRegularCaller(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RegularCaller{RegularCallerCaller: RegularCallerCaller{contract: contract}, RegularCallerTransactor: RegularCallerTransactor{contract: contract}, RegularCallerFilterer: RegularCallerFilterer{contract: contract}}, nil
}

// NewRegularCallerCaller creates a new read-only instance of RegularCaller, bound to a specific deployed contract.
func NewRegularCallerCaller(address common.Address, caller bind.ContractCaller) (*RegularCallerCaller, error) {
	contract, err := bindRegularCaller(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegularCallerCaller{contract: contract}, nil
}

// NewRegularCallerTransactor creates a new write-only instance of RegularCaller, bound to a specific deployed contract.
func NewRegularCallerTransactor(address common.Address, transactor bind.ContractTransactor) (*RegularCallerTransactor, error) {
	contract, err := bindRegularCaller(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegularCallerTransactor{contract: contract}, nil
}

// NewRegularCallerFilterer creates a new log filterer instance of RegularCaller, bound to a specific deployed contract.
func NewRegularCallerFilterer(address common.Address, filterer bind.ContractFilterer) (*RegularCallerFilterer, error) {
	contract, err := bindRegularCaller(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegularCallerFilterer{contract: contract}, nil
}

// bindRegularCaller binds a generic wrapper to an already deployed contract.
func bindRegularCaller(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RegularCallerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RegularCaller *RegularCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RegularCaller.Contract.RegularCallerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RegularCaller *RegularCallerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RegularCaller.Contract.RegularCallerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RegularCaller *RegularCallerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RegularCaller.Contract.RegularCallerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RegularCaller *RegularCallerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RegularCaller.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RegularCaller *RegularCallerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RegularCaller.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RegularCaller *RegularCallerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RegularCaller.Contract.contract.Transact(opts, method, params...)
}

// TestBech32ToHexAddr is a free data retrieval call binding the contract method 0x42875b1d.
//
// Solidity: function testBech32ToHexAddr() view returns(bool)
func (_RegularCaller *RegularCallerCaller) TestBech32ToHexAddr(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RegularCaller.contract.Call(opts, &out, "testBech32ToHexAddr")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// TestBech32ToHexAddr is a free data retrieval call binding the contract method 0x42875b1d.
//
// Solidity: function testBech32ToHexAddr() view returns(bool)
func (_RegularCaller *RegularCallerSession) TestBech32ToHexAddr() (bool, error) {
	return _RegularCaller.Contract.TestBech32ToHexAddr(&_RegularCaller.CallOpts)
}

// TestBech32ToHexAddr is a free data retrieval call binding the contract method 0x42875b1d.
//
// Solidity: function testBech32ToHexAddr() view returns(bool)
func (_RegularCaller *RegularCallerCallerSession) TestBech32ToHexAddr() (bool, error) {
	return _RegularCaller.Contract.TestBech32ToHexAddr(&_RegularCaller.CallOpts)
}

// TestBech32ify is a free data retrieval call binding the contract method 0xd7c9bd02.
//
// Solidity: function testBech32ify() view returns(bool)
func (_RegularCaller *RegularCallerCaller) TestBech32ify(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RegularCaller.contract.Call(opts, &out, "testBech32ify")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// TestBech32ify is a free data retrieval call binding the contract method 0xd7c9bd02.
//
// Solidity: function testBech32ify() view returns(bool)
func (_RegularCaller *RegularCallerSession) TestBech32ify() (bool, error) {
	return _RegularCaller.Contract.TestBech32ify(&_RegularCaller.CallOpts)
}

// TestBech32ify is a free data retrieval call binding the contract method 0xd7c9bd02.
//
// Solidity: function testBech32ify() view returns(bool)
func (_RegularCaller *RegularCallerCallerSession) TestBech32ify() (bool, error) {
	return _RegularCaller.Contract.TestBech32ify(&_RegularCaller.CallOpts)
}

// TestRegularCall is a paid mutator transaction binding the contract method 0x48fc7db1.
//
// Solidity: function testRegularCall(string method, address addr) returns(uint256)
func (_RegularCaller *RegularCallerTransactor) TestRegularCall(opts *bind.TransactOpts, method string, addr common.Address) (*types.Transaction, error) {
	return _RegularCaller.contract.Transact(opts, "testRegularCall", method, addr)
}

// TestRegularCall is a paid mutator transaction binding the contract method 0x48fc7db1.
//
// Solidity: function testRegularCall(string method, address addr) returns(uint256)
func (_RegularCaller *RegularCallerSession) TestRegularCall(method string, addr common.Address) (*types.Transaction, error) {
	return _RegularCaller.Contract.TestRegularCall(&_RegularCaller.TransactOpts, method, addr)
}

// TestRegularCall is a paid mutator transaction binding the contract method 0x48fc7db1.
//
// Solidity: function testRegularCall(string method, address addr) returns(uint256)
func (_RegularCaller *RegularCallerTransactorSession) TestRegularCall(method string, addr common.Address) (*types.Transaction, error) {
	return _RegularCaller.Contract.TestRegularCall(&_RegularCaller.TransactOpts, method, addr)
}
