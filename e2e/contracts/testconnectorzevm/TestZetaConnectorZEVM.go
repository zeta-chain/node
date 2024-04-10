// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testconnectorzevm

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

// TestZetaConnectorZEVMMetaData contains all meta data concerning the TestZetaConnectorZEVM contract.
var TestZetaConnectorZEVMMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wzeta_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"FailedZetaSent\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyFungibleModule\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyWZETA\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WZETATransferFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"wzeta_\",\"type\":\"address\"}],\"name\":\"SetWZETA\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sourceTxOriginAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"zetaTxSenderAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"destinationAddress\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"zetaValueAndGas\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationGasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"zetaParams\",\"type\":\"bytes\"}],\"name\":\"ZetaSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FUNGIBLE_MODULE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"foo\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destinationAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"destinationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"zetaValueAndGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"zetaParams\",\"type\":\"bytes\"}],\"internalType\":\"structZetaInterfaces.SendInput\",\"name\":\"input\",\"type\":\"tuple\"}],\"name\":\"send\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wzeta_\",\"type\":\"address\"}],\"name\":\"setWzetaAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wzeta\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610b56380380610b56833981810160405281019061003291906100db565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050610108565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100a88261007d565b9050919050565b6100b88161009d565b81146100c357600080fd5b50565b6000815190506100d5816100af565b92915050565b6000602082840312156100f1576100f0610078565b5b60006100ff848285016100c6565b91505092915050565b610a3f806101176000396000f3fe60806040526004361061004d5760003560e01c8062173d46146100de5780633ce4a5bc14610109578063c298557814610134578063eb3bacbd1461015f578063ec02690114610188576100d9565b366100d95760008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146100d7576040517f6e6b6de700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b005b600080fd5b3480156100ea57600080fd5b506100f36101b1565b604051610100919061061b565b60405180910390f35b34801561011557600080fd5b5061011e6101d5565b60405161012b919061061b565b60405180910390f35b34801561014057600080fd5b506101496101ed565b60405161015691906106c6565b60405180910390f35b34801561016b57600080fd5b506101866004803603810190610181919061071e565b61022a565b005b34801561019457600080fd5b506101af60048036038101906101aa919061076f565b61031d565b005b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b73735b14bb79463307aacbed86daf3322b1e6226ab81565b60606040518060400160405280600381526020017f666f6f0000000000000000000000000000000000000000000000000000000000815250905090565b73735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102a3576040517fea02b3f300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f7325870b05f8f3412c318a35fc6a74feca51ea15811ec7a257676ca4db9d417681604051610312919061061b565b60405180910390a150565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd333084608001356040518463ffffffff1660e01b815260040161037e939291906107d1565b6020604051808303816000875af115801561039d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103c19190610840565b6103f7576040517fa8c6fd4a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632e1a7d4d82608001356040518263ffffffff1660e01b8152600401610454919061086d565b600060405180830381600087803b15801561046e57600080fd5b505af1158015610482573d6000803e3d6000fd5b50505050600073735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff1682608001356040516104c4906108b9565b60006040518083038185875af1925050503d8060008114610501576040519150601f19603f3d011682016040523d82523d6000602084013e610506565b606091505b5050905080610541576040517fc7ffc47b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81600001353373ffffffffffffffffffffffffffffffffffffffff167f7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e43285806020019061058f91906108dd565b876080013588604001358980606001906105a991906108dd565b8b8060a001906105b991906108dd565b6040516105ce9998979695949392919061098d565b60405180910390a35050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610605826105da565b9050919050565b610615816105fa565b82525050565b6000602082019050610630600083018461060c565b92915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610670578082015181840152602081019050610655565b60008484015250505050565b6000601f19601f8301169050919050565b600061069882610636565b6106a28185610641565b93506106b2818560208601610652565b6106bb8161067c565b840191505092915050565b600060208201905081810360008301526106e0818461068d565b905092915050565b600080fd5b600080fd5b6106fb816105fa565b811461070657600080fd5b50565b600081359050610718816106f2565b92915050565b600060208284031215610734576107336106e8565b5b600061074284828501610709565b91505092915050565b600080fd5b600060c082840312156107665761076561074b565b5b81905092915050565b600060208284031215610785576107846106e8565b5b600082013567ffffffffffffffff8111156107a3576107a26106ed565b5b6107af84828501610750565b91505092915050565b6000819050919050565b6107cb816107b8565b82525050565b60006060820190506107e6600083018661060c565b6107f3602083018561060c565b61080060408301846107c2565b949350505050565b60008115159050919050565b61081d81610808565b811461082857600080fd5b50565b60008151905061083a81610814565b92915050565b600060208284031215610856576108556106e8565b5b60006108648482850161082b565b91505092915050565b600060208201905061088260008301846107c2565b92915050565b600081905092915050565b50565b60006108a3600083610888565b91506108ae82610893565b600082019050919050565b60006108c482610896565b9150819050919050565b600080fd5b600080fd5b600080fd5b600080833560016020038436030381126108fa576108f96108ce565b5b80840192508235915067ffffffffffffffff82111561091c5761091b6108d3565b5b602083019250600182023603831315610938576109376108d8565b5b509250929050565b600082825260208201905092915050565b82818337600083830152505050565b600061096c8385610940565b9350610979838584610951565b6109828361067c565b840190509392505050565b600060c0820190506109a2600083018c61060c565b81810360208301526109b5818a8c610960565b90506109c460408301896107c2565b6109d160608301886107c2565b81810360808301526109e4818688610960565b905081810360a08301526109f9818486610960565b90509a995050505050505050505056fea26469706673582212206647922040def2c6972690c7b621d67ba2619c6888ae0b5c33ce88c440cebffa64736f6c63430008170033",
}

// TestZetaConnectorZEVMABI is the input ABI used to generate the binding from.
// Deprecated: Use TestZetaConnectorZEVMMetaData.ABI instead.
var TestZetaConnectorZEVMABI = TestZetaConnectorZEVMMetaData.ABI

// TestZetaConnectorZEVMBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestZetaConnectorZEVMMetaData.Bin instead.
var TestZetaConnectorZEVMBin = TestZetaConnectorZEVMMetaData.Bin

// DeployTestZetaConnectorZEVM deploys a new Ethereum contract, binding an instance of TestZetaConnectorZEVM to it.
func DeployTestZetaConnectorZEVM(auth *bind.TransactOpts, backend bind.ContractBackend, wzeta_ common.Address) (common.Address, *types.Transaction, *TestZetaConnectorZEVM, error) {
	parsed, err := TestZetaConnectorZEVMMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestZetaConnectorZEVMBin), backend, wzeta_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestZetaConnectorZEVM{TestZetaConnectorZEVMCaller: TestZetaConnectorZEVMCaller{contract: contract}, TestZetaConnectorZEVMTransactor: TestZetaConnectorZEVMTransactor{contract: contract}, TestZetaConnectorZEVMFilterer: TestZetaConnectorZEVMFilterer{contract: contract}}, nil
}

// TestZetaConnectorZEVM is an auto generated Go binding around an Ethereum contract.
type TestZetaConnectorZEVM struct {
	TestZetaConnectorZEVMCaller     // Read-only binding to the contract
	TestZetaConnectorZEVMTransactor // Write-only binding to the contract
	TestZetaConnectorZEVMFilterer   // Log filterer for contract events
}

// TestZetaConnectorZEVMCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestZetaConnectorZEVMCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestZetaConnectorZEVMTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestZetaConnectorZEVMTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestZetaConnectorZEVMFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestZetaConnectorZEVMFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestZetaConnectorZEVMSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestZetaConnectorZEVMSession struct {
	Contract     *TestZetaConnectorZEVM // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// TestZetaConnectorZEVMCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestZetaConnectorZEVMCallerSession struct {
	Contract *TestZetaConnectorZEVMCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// TestZetaConnectorZEVMTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestZetaConnectorZEVMTransactorSession struct {
	Contract     *TestZetaConnectorZEVMTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// TestZetaConnectorZEVMRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestZetaConnectorZEVMRaw struct {
	Contract *TestZetaConnectorZEVM // Generic contract binding to access the raw methods on
}

// TestZetaConnectorZEVMCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestZetaConnectorZEVMCallerRaw struct {
	Contract *TestZetaConnectorZEVMCaller // Generic read-only contract binding to access the raw methods on
}

// TestZetaConnectorZEVMTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestZetaConnectorZEVMTransactorRaw struct {
	Contract *TestZetaConnectorZEVMTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestZetaConnectorZEVM creates a new instance of TestZetaConnectorZEVM, bound to a specific deployed contract.
func NewTestZetaConnectorZEVM(address common.Address, backend bind.ContractBackend) (*TestZetaConnectorZEVM, error) {
	contract, err := bindTestZetaConnectorZEVM(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestZetaConnectorZEVM{TestZetaConnectorZEVMCaller: TestZetaConnectorZEVMCaller{contract: contract}, TestZetaConnectorZEVMTransactor: TestZetaConnectorZEVMTransactor{contract: contract}, TestZetaConnectorZEVMFilterer: TestZetaConnectorZEVMFilterer{contract: contract}}, nil
}

// NewTestZetaConnectorZEVMCaller creates a new read-only instance of TestZetaConnectorZEVM, bound to a specific deployed contract.
func NewTestZetaConnectorZEVMCaller(address common.Address, caller bind.ContractCaller) (*TestZetaConnectorZEVMCaller, error) {
	contract, err := bindTestZetaConnectorZEVM(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestZetaConnectorZEVMCaller{contract: contract}, nil
}

// NewTestZetaConnectorZEVMTransactor creates a new write-only instance of TestZetaConnectorZEVM, bound to a specific deployed contract.
func NewTestZetaConnectorZEVMTransactor(address common.Address, transactor bind.ContractTransactor) (*TestZetaConnectorZEVMTransactor, error) {
	contract, err := bindTestZetaConnectorZEVM(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestZetaConnectorZEVMTransactor{contract: contract}, nil
}

// NewTestZetaConnectorZEVMFilterer creates a new log filterer instance of TestZetaConnectorZEVM, bound to a specific deployed contract.
func NewTestZetaConnectorZEVMFilterer(address common.Address, filterer bind.ContractFilterer) (*TestZetaConnectorZEVMFilterer, error) {
	contract, err := bindTestZetaConnectorZEVM(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestZetaConnectorZEVMFilterer{contract: contract}, nil
}

// bindTestZetaConnectorZEVM binds a generic wrapper to an already deployed contract.
func bindTestZetaConnectorZEVM(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestZetaConnectorZEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestZetaConnectorZEVM.Contract.TestZetaConnectorZEVMCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.Contract.TestZetaConnectorZEVMTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.Contract.TestZetaConnectorZEVMTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestZetaConnectorZEVM.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.Contract.contract.Transact(opts, method, params...)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMCaller) FUNGIBLEMODULEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestZetaConnectorZEVM.contract.Call(opts, &out, "FUNGIBLE_MODULE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _TestZetaConnectorZEVM.Contract.FUNGIBLEMODULEADDRESS(&_TestZetaConnectorZEVM.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMCallerSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _TestZetaConnectorZEVM.Contract.FUNGIBLEMODULEADDRESS(&_TestZetaConnectorZEVM.CallOpts)
}

// Foo is a free data retrieval call binding the contract method 0xc2985578.
//
// Solidity: function foo() pure returns(string)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMCaller) Foo(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestZetaConnectorZEVM.contract.Call(opts, &out, "foo")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Foo is a free data retrieval call binding the contract method 0xc2985578.
//
// Solidity: function foo() pure returns(string)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMSession) Foo() (string, error) {
	return _TestZetaConnectorZEVM.Contract.Foo(&_TestZetaConnectorZEVM.CallOpts)
}

// Foo is a free data retrieval call binding the contract method 0xc2985578.
//
// Solidity: function foo() pure returns(string)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMCallerSession) Foo() (string, error) {
	return _TestZetaConnectorZEVM.Contract.Foo(&_TestZetaConnectorZEVM.CallOpts)
}

// Wzeta is a free data retrieval call binding the contract method 0x00173d46.
//
// Solidity: function wzeta() view returns(address)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMCaller) Wzeta(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestZetaConnectorZEVM.contract.Call(opts, &out, "wzeta")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Wzeta is a free data retrieval call binding the contract method 0x00173d46.
//
// Solidity: function wzeta() view returns(address)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMSession) Wzeta() (common.Address, error) {
	return _TestZetaConnectorZEVM.Contract.Wzeta(&_TestZetaConnectorZEVM.CallOpts)
}

// Wzeta is a free data retrieval call binding the contract method 0x00173d46.
//
// Solidity: function wzeta() view returns(address)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMCallerSession) Wzeta() (common.Address, error) {
	return _TestZetaConnectorZEVM.Contract.Wzeta(&_TestZetaConnectorZEVM.CallOpts)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMTransactor) Send(opts *bind.TransactOpts, input ZetaInterfacesSendInput) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.contract.Transact(opts, "send", input)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMSession) Send(input ZetaInterfacesSendInput) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.Contract.Send(&_TestZetaConnectorZEVM.TransactOpts, input)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMTransactorSession) Send(input ZetaInterfacesSendInput) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.Contract.Send(&_TestZetaConnectorZEVM.TransactOpts, input)
}

// SetWzetaAddress is a paid mutator transaction binding the contract method 0xeb3bacbd.
//
// Solidity: function setWzetaAddress(address wzeta_) returns()
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMTransactor) SetWzetaAddress(opts *bind.TransactOpts, wzeta_ common.Address) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.contract.Transact(opts, "setWzetaAddress", wzeta_)
}

// SetWzetaAddress is a paid mutator transaction binding the contract method 0xeb3bacbd.
//
// Solidity: function setWzetaAddress(address wzeta_) returns()
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMSession) SetWzetaAddress(wzeta_ common.Address) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.Contract.SetWzetaAddress(&_TestZetaConnectorZEVM.TransactOpts, wzeta_)
}

// SetWzetaAddress is a paid mutator transaction binding the contract method 0xeb3bacbd.
//
// Solidity: function setWzetaAddress(address wzeta_) returns()
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMTransactorSession) SetWzetaAddress(wzeta_ common.Address) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.Contract.SetWzetaAddress(&_TestZetaConnectorZEVM.TransactOpts, wzeta_)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMSession) Receive() (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.Contract.Receive(&_TestZetaConnectorZEVM.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMTransactorSession) Receive() (*types.Transaction, error) {
	return _TestZetaConnectorZEVM.Contract.Receive(&_TestZetaConnectorZEVM.TransactOpts)
}

// TestZetaConnectorZEVMSetWZETAIterator is returned from FilterSetWZETA and is used to iterate over the raw logs and unpacked data for SetWZETA events raised by the TestZetaConnectorZEVM contract.
type TestZetaConnectorZEVMSetWZETAIterator struct {
	Event *TestZetaConnectorZEVMSetWZETA // Event containing the contract specifics and raw log

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
func (it *TestZetaConnectorZEVMSetWZETAIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestZetaConnectorZEVMSetWZETA)
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
		it.Event = new(TestZetaConnectorZEVMSetWZETA)
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
func (it *TestZetaConnectorZEVMSetWZETAIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestZetaConnectorZEVMSetWZETAIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestZetaConnectorZEVMSetWZETA represents a SetWZETA event raised by the TestZetaConnectorZEVM contract.
type TestZetaConnectorZEVMSetWZETA struct {
	Wzeta common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSetWZETA is a free log retrieval operation binding the contract event 0x7325870b05f8f3412c318a35fc6a74feca51ea15811ec7a257676ca4db9d4176.
//
// Solidity: event SetWZETA(address wzeta_)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMFilterer) FilterSetWZETA(opts *bind.FilterOpts) (*TestZetaConnectorZEVMSetWZETAIterator, error) {

	logs, sub, err := _TestZetaConnectorZEVM.contract.FilterLogs(opts, "SetWZETA")
	if err != nil {
		return nil, err
	}
	return &TestZetaConnectorZEVMSetWZETAIterator{contract: _TestZetaConnectorZEVM.contract, event: "SetWZETA", logs: logs, sub: sub}, nil
}

// WatchSetWZETA is a free log subscription operation binding the contract event 0x7325870b05f8f3412c318a35fc6a74feca51ea15811ec7a257676ca4db9d4176.
//
// Solidity: event SetWZETA(address wzeta_)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMFilterer) WatchSetWZETA(opts *bind.WatchOpts, sink chan<- *TestZetaConnectorZEVMSetWZETA) (event.Subscription, error) {

	logs, sub, err := _TestZetaConnectorZEVM.contract.WatchLogs(opts, "SetWZETA")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestZetaConnectorZEVMSetWZETA)
				if err := _TestZetaConnectorZEVM.contract.UnpackLog(event, "SetWZETA", log); err != nil {
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

// ParseSetWZETA is a log parse operation binding the contract event 0x7325870b05f8f3412c318a35fc6a74feca51ea15811ec7a257676ca4db9d4176.
//
// Solidity: event SetWZETA(address wzeta_)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMFilterer) ParseSetWZETA(log types.Log) (*TestZetaConnectorZEVMSetWZETA, error) {
	event := new(TestZetaConnectorZEVMSetWZETA)
	if err := _TestZetaConnectorZEVM.contract.UnpackLog(event, "SetWZETA", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestZetaConnectorZEVMZetaSentIterator is returned from FilterZetaSent and is used to iterate over the raw logs and unpacked data for ZetaSent events raised by the TestZetaConnectorZEVM contract.
type TestZetaConnectorZEVMZetaSentIterator struct {
	Event *TestZetaConnectorZEVMZetaSent // Event containing the contract specifics and raw log

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
func (it *TestZetaConnectorZEVMZetaSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestZetaConnectorZEVMZetaSent)
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
		it.Event = new(TestZetaConnectorZEVMZetaSent)
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
func (it *TestZetaConnectorZEVMZetaSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestZetaConnectorZEVMZetaSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestZetaConnectorZEVMZetaSent represents a ZetaSent event raised by the TestZetaConnectorZEVM contract.
type TestZetaConnectorZEVMZetaSent struct {
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
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMFilterer) FilterZetaSent(opts *bind.FilterOpts, zetaTxSenderAddress []common.Address, destinationChainId []*big.Int) (*TestZetaConnectorZEVMZetaSentIterator, error) {

	var zetaTxSenderAddressRule []interface{}
	for _, zetaTxSenderAddressItem := range zetaTxSenderAddress {
		zetaTxSenderAddressRule = append(zetaTxSenderAddressRule, zetaTxSenderAddressItem)
	}
	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	logs, sub, err := _TestZetaConnectorZEVM.contract.FilterLogs(opts, "ZetaSent", zetaTxSenderAddressRule, destinationChainIdRule)
	if err != nil {
		return nil, err
	}
	return &TestZetaConnectorZEVMZetaSentIterator{contract: _TestZetaConnectorZEVM.contract, event: "ZetaSent", logs: logs, sub: sub}, nil
}

// WatchZetaSent is a free log subscription operation binding the contract event 0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4.
//
// Solidity: event ZetaSent(address sourceTxOriginAddress, address indexed zetaTxSenderAddress, uint256 indexed destinationChainId, bytes destinationAddress, uint256 zetaValueAndGas, uint256 destinationGasLimit, bytes message, bytes zetaParams)
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMFilterer) WatchZetaSent(opts *bind.WatchOpts, sink chan<- *TestZetaConnectorZEVMZetaSent, zetaTxSenderAddress []common.Address, destinationChainId []*big.Int) (event.Subscription, error) {

	var zetaTxSenderAddressRule []interface{}
	for _, zetaTxSenderAddressItem := range zetaTxSenderAddress {
		zetaTxSenderAddressRule = append(zetaTxSenderAddressRule, zetaTxSenderAddressItem)
	}
	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	logs, sub, err := _TestZetaConnectorZEVM.contract.WatchLogs(opts, "ZetaSent", zetaTxSenderAddressRule, destinationChainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestZetaConnectorZEVMZetaSent)
				if err := _TestZetaConnectorZEVM.contract.UnpackLog(event, "ZetaSent", log); err != nil {
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
func (_TestZetaConnectorZEVM *TestZetaConnectorZEVMFilterer) ParseZetaSent(log types.Log) (*TestZetaConnectorZEVMZetaSent, error) {
	event := new(TestZetaConnectorZEVMZetaSent)
	if err := _TestZetaConnectorZEVM.contract.UnpackLog(event, "ZetaSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
