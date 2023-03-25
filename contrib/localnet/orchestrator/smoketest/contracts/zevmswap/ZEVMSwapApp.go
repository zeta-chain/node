// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zevmswap

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

// ZEVMSwapAppMetaData contains all meta data concerning the ZEVMSwapApp contract.
var ZEVMSwapAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router02_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"systemContract_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LowAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"targetZRC20\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"minAmountOut\",\"type\":\"uint256\"}],\"name\":\"encodeMemo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router02\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"systemContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b506040516200118c3803806200118c8339818101604052810190620000379190620000c4565b8173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff1660601b815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1660601b8152505050506200015e565b600081519050620000be8162000144565b92915050565b60008060408385031215620000de57620000dd6200013f565b5b6000620000ee85828601620000ad565b92505060206200010185828601620000ad565b9150509250929050565b600062000118826200011f565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600080fd5b6200014f816200010b565b81146200015b57600080fd5b50565b60805160601c60a05160601c610fed6200019f6000396000818161010d01526101550152600081816101310152818161030301526103970152610fed6000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c8063678f75dc14610051578063bb88b76914610081578063bd00c9c41461009f578063c8522691146100bd575b600080fd5b61006b600480360381019061006691906108cf565b6100d9565b6040516100789190610c09565b60405180910390f35b61008961010b565b6040516100969190610b85565b60405180910390f35b6100a761012f565b6040516100b49190610b85565b60405180910390f35b6100d760048036038101906100d29190610983565b610153565b005b6060848484846040516020016100f29493929190610ba0565b6040516020818303038152906040529050949350505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101d8576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006060600084848101906101ed9190610860565b8093508194508295505050506060600267ffffffffffffffff81111561021657610215610efd565b5b6040519080825280602002602001820160405280156102445781602001602082028036833780820191505090505b509050878160008151811061025c5761025b610ece565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505083816001815181106102ab576102aa610ece565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508773ffffffffffffffffffffffffffffffffffffffff1663095ea7b37f0000000000000000000000000000000000000000000000000000000000000000896040518363ffffffff1660e01b8152600401610340929190610be0565b602060405180830381600087803b15801561035a57600080fd5b505af115801561036e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103929190610a40565b5060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166338ed17398960008530680100000000000000006040518663ffffffff1660e01b8152600401610400959493929190610c5b565b600060405180830381600087803b15801561041a57600080fd5b505af115801561042e573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525081019061045791906109f7565b905060008573ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b8152600401604080518083038186803b1580156104a057600080fd5b505afa1580156104b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104d89190610943565b915050816001815181106104ef576104ee610ece565b5b6020026020010151811115610530576040517f4e8ed22600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8573ffffffffffffffffffffffffffffffffffffffff1663095ea7b387836040518363ffffffff1660e01b815260040161056b929190610be0565b602060405180830381600087803b15801561058557600080fd5b505af1158015610599573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105bd9190610a40565b508573ffffffffffffffffffffffffffffffffffffffff1663c70126268683856001815181106105f0576105ef610ece565b5b60200260200101516106029190610d8c565b6040518363ffffffff1660e01b815260040161061f929190610c2b565b602060405180830381600087803b15801561063957600080fd5b505af115801561064d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106719190610a40565b5050505050505050505050565b600061069161068c84610cda565b610cb5565b905080838252602082019050828560208602820111156106b4576106b3610f36565b5b60005b858110156106e457816106ca888261084b565b8452602084019350602083019250506001810190506106b7565b5050509392505050565b60006107016106fc84610d06565b610cb5565b90508281526020810184848401111561071d5761071c610f3b565b5b610728848285610e2c565b509392505050565b60008135905061073f81610f5b565b92915050565b60008151905061075481610f5b565b92915050565b60008135905061076981610f72565b92915050565b600082601f83011261078457610783610f31565b5b815161079484826020860161067e565b91505092915050565b6000815190506107ac81610f89565b92915050565b60008083601f8401126107c8576107c7610f31565b5b8235905067ffffffffffffffff8111156107e5576107e4610f2c565b5b60208301915083600182028301111561080157610800610f36565b5b9250929050565b600082601f83011261081d5761081c610f31565b5b813561082d8482602086016106ee565b91505092915050565b60008135905061084581610fa0565b92915050565b60008151905061085a81610fa0565b92915050565b60008060006060848603121561087957610878610f45565b5b60006108878682870161075a565b935050602084013567ffffffffffffffff8111156108a8576108a7610f40565b5b6108b486828701610808565b92505060406108c586828701610836565b9150509250925092565b600080600080606085870312156108e9576108e8610f45565b5b60006108f787828801610730565b945050602085013567ffffffffffffffff81111561091857610917610f40565b5b610924878288016107b2565b9350935050604061093787828801610836565b91505092959194509250565b6000806040838503121561095a57610959610f45565b5b600061096885828601610745565b92505060206109798582860161084b565b9150509250929050565b6000806000806060858703121561099d5761099c610f45565b5b60006109ab87828801610730565b94505060206109bc87828801610836565b935050604085013567ffffffffffffffff8111156109dd576109dc610f40565b5b6109e9878288016107b2565b925092505092959194509250565b600060208284031215610a0d57610a0c610f45565b5b600082015167ffffffffffffffff811115610a2b57610a2a610f40565b5b610a378482850161076f565b91505092915050565b600060208284031215610a5657610a55610f45565b5b6000610a648482850161079d565b91505092915050565b6000610a798383610a85565b60208301905092915050565b610a8e81610dc0565b82525050565b610a9d81610dc0565b82525050565b6000610aae82610d47565b610ab88185610d6a565b9350610ac383610d37565b8060005b83811015610af4578151610adb8882610a6d565b9750610ae683610d5d565b925050600181019050610ac7565b5085935050505092915050565b6000610b0d8385610d7b565b9350610b1a838584610e2c565b610b2383610f4a565b840190509392505050565b6000610b3982610d52565b610b438185610d7b565b9350610b53818560208601610e3b565b610b5c81610f4a565b840191505092915050565b610b7081610e1a565b82525050565b610b7f81610e10565b82525050565b6000602082019050610b9a6000830184610a94565b92915050565b6000606082019050610bb56000830187610a94565b8181036020830152610bc8818587610b01565b9050610bd76040830184610b76565b95945050505050565b6000604082019050610bf56000830185610a94565b610c026020830184610b76565b9392505050565b60006020820190508181036000830152610c238184610b2e565b905092915050565b60006040820190508181036000830152610c458185610b2e565b9050610c546020830184610b76565b9392505050565b600060a082019050610c706000830188610b76565b610c7d6020830187610b67565b8181036040830152610c8f8186610aa3565b9050610c9e6060830185610a94565b610cab6080830184610b76565b9695505050505050565b6000610cbf610cd0565b9050610ccb8282610e6e565b919050565b6000604051905090565b600067ffffffffffffffff821115610cf557610cf4610efd565b5b602082029050602081019050919050565b600067ffffffffffffffff821115610d2157610d20610efd565b5b610d2a82610f4a565b9050602081019050919050565b6000819050602082019050919050565b600081519050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b6000610d9782610e10565b9150610da283610e10565b925082821015610db557610db4610e9f565b5b828203905092915050565b6000610dcb82610df0565b9050919050565b6000610ddd82610df0565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b6000610e2582610e10565b9050919050565b82818337600083830152505050565b60005b83811015610e59578082015181840152602081019050610e3e565b83811115610e68576000848401525b50505050565b610e7782610f4a565b810181811067ffffffffffffffff82111715610e9657610e95610efd565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b610f6481610dc0565b8114610f6f57600080fd5b50565b610f7b81610dd2565b8114610f8657600080fd5b50565b610f9281610de4565b8114610f9d57600080fd5b50565b610fa981610e10565b8114610fb457600080fd5b5056fea264697066735822122056a11840b156186209e474d6bda7743202c4672f770de2e1e95a2d4f4bb7429264736f6c63430008070033",
}

// ZEVMSwapAppABI is the input ABI used to generate the binding from.
// Deprecated: Use ZEVMSwapAppMetaData.ABI instead.
var ZEVMSwapAppABI = ZEVMSwapAppMetaData.ABI

// ZEVMSwapAppBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ZEVMSwapAppMetaData.Bin instead.
var ZEVMSwapAppBin = ZEVMSwapAppMetaData.Bin

// DeployZEVMSwapApp deploys a new Ethereum contract, binding an instance of ZEVMSwapApp to it.
func DeployZEVMSwapApp(auth *bind.TransactOpts, backend bind.ContractBackend, router02_ common.Address, systemContract_ common.Address) (common.Address, *types.Transaction, *ZEVMSwapApp, error) {
	parsed, err := ZEVMSwapAppMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ZEVMSwapAppBin), backend, router02_, systemContract_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ZEVMSwapApp{ZEVMSwapAppCaller: ZEVMSwapAppCaller{contract: contract}, ZEVMSwapAppTransactor: ZEVMSwapAppTransactor{contract: contract}, ZEVMSwapAppFilterer: ZEVMSwapAppFilterer{contract: contract}}, nil
}

// ZEVMSwapApp is an auto generated Go binding around an Ethereum contract.
type ZEVMSwapApp struct {
	ZEVMSwapAppCaller     // Read-only binding to the contract
	ZEVMSwapAppTransactor // Write-only binding to the contract
	ZEVMSwapAppFilterer   // Log filterer for contract events
}

// ZEVMSwapAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZEVMSwapAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZEVMSwapAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZEVMSwapAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZEVMSwapAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZEVMSwapAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZEVMSwapAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZEVMSwapAppSession struct {
	Contract     *ZEVMSwapApp      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZEVMSwapAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZEVMSwapAppCallerSession struct {
	Contract *ZEVMSwapAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ZEVMSwapAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZEVMSwapAppTransactorSession struct {
	Contract     *ZEVMSwapAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ZEVMSwapAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZEVMSwapAppRaw struct {
	Contract *ZEVMSwapApp // Generic contract binding to access the raw methods on
}

// ZEVMSwapAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZEVMSwapAppCallerRaw struct {
	Contract *ZEVMSwapAppCaller // Generic read-only contract binding to access the raw methods on
}

// ZEVMSwapAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZEVMSwapAppTransactorRaw struct {
	Contract *ZEVMSwapAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZEVMSwapApp creates a new instance of ZEVMSwapApp, bound to a specific deployed contract.
func NewZEVMSwapApp(address common.Address, backend bind.ContractBackend) (*ZEVMSwapApp, error) {
	contract, err := bindZEVMSwapApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZEVMSwapApp{ZEVMSwapAppCaller: ZEVMSwapAppCaller{contract: contract}, ZEVMSwapAppTransactor: ZEVMSwapAppTransactor{contract: contract}, ZEVMSwapAppFilterer: ZEVMSwapAppFilterer{contract: contract}}, nil
}

// NewZEVMSwapAppCaller creates a new read-only instance of ZEVMSwapApp, bound to a specific deployed contract.
func NewZEVMSwapAppCaller(address common.Address, caller bind.ContractCaller) (*ZEVMSwapAppCaller, error) {
	contract, err := bindZEVMSwapApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZEVMSwapAppCaller{contract: contract}, nil
}

// NewZEVMSwapAppTransactor creates a new write-only instance of ZEVMSwapApp, bound to a specific deployed contract.
func NewZEVMSwapAppTransactor(address common.Address, transactor bind.ContractTransactor) (*ZEVMSwapAppTransactor, error) {
	contract, err := bindZEVMSwapApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZEVMSwapAppTransactor{contract: contract}, nil
}

// NewZEVMSwapAppFilterer creates a new log filterer instance of ZEVMSwapApp, bound to a specific deployed contract.
func NewZEVMSwapAppFilterer(address common.Address, filterer bind.ContractFilterer) (*ZEVMSwapAppFilterer, error) {
	contract, err := bindZEVMSwapApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZEVMSwapAppFilterer{contract: contract}, nil
}

// bindZEVMSwapApp binds a generic wrapper to an already deployed contract.
func bindZEVMSwapApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ZEVMSwapAppABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZEVMSwapApp *ZEVMSwapAppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZEVMSwapApp.Contract.ZEVMSwapAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZEVMSwapApp *ZEVMSwapAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZEVMSwapApp.Contract.ZEVMSwapAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZEVMSwapApp *ZEVMSwapAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZEVMSwapApp.Contract.ZEVMSwapAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZEVMSwapApp *ZEVMSwapAppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZEVMSwapApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZEVMSwapApp *ZEVMSwapAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZEVMSwapApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZEVMSwapApp *ZEVMSwapAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZEVMSwapApp.Contract.contract.Transact(opts, method, params...)
}

// EncodeMemo is a free data retrieval call binding the contract method 0x678f75dc.
//
// Solidity: function encodeMemo(address targetZRC20, bytes recipient, uint256 minAmountOut) pure returns(bytes)
func (_ZEVMSwapApp *ZEVMSwapAppCaller) EncodeMemo(opts *bind.CallOpts, targetZRC20 common.Address, recipient []byte, minAmountOut *big.Int) ([]byte, error) {
	var out []interface{}
	err := _ZEVMSwapApp.contract.Call(opts, &out, "encodeMemo", targetZRC20, recipient, minAmountOut)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeMemo is a free data retrieval call binding the contract method 0x678f75dc.
//
// Solidity: function encodeMemo(address targetZRC20, bytes recipient, uint256 minAmountOut) pure returns(bytes)
func (_ZEVMSwapApp *ZEVMSwapAppSession) EncodeMemo(targetZRC20 common.Address, recipient []byte, minAmountOut *big.Int) ([]byte, error) {
	return _ZEVMSwapApp.Contract.EncodeMemo(&_ZEVMSwapApp.CallOpts, targetZRC20, recipient, minAmountOut)
}

// EncodeMemo is a free data retrieval call binding the contract method 0x678f75dc.
//
// Solidity: function encodeMemo(address targetZRC20, bytes recipient, uint256 minAmountOut) pure returns(bytes)
func (_ZEVMSwapApp *ZEVMSwapAppCallerSession) EncodeMemo(targetZRC20 common.Address, recipient []byte, minAmountOut *big.Int) ([]byte, error) {
	return _ZEVMSwapApp.Contract.EncodeMemo(&_ZEVMSwapApp.CallOpts, targetZRC20, recipient, minAmountOut)
}

// Router02 is a free data retrieval call binding the contract method 0xbd00c9c4.
//
// Solidity: function router02() view returns(address)
func (_ZEVMSwapApp *ZEVMSwapAppCaller) Router02(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZEVMSwapApp.contract.Call(opts, &out, "router02")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Router02 is a free data retrieval call binding the contract method 0xbd00c9c4.
//
// Solidity: function router02() view returns(address)
func (_ZEVMSwapApp *ZEVMSwapAppSession) Router02() (common.Address, error) {
	return _ZEVMSwapApp.Contract.Router02(&_ZEVMSwapApp.CallOpts)
}

// Router02 is a free data retrieval call binding the contract method 0xbd00c9c4.
//
// Solidity: function router02() view returns(address)
func (_ZEVMSwapApp *ZEVMSwapAppCallerSession) Router02() (common.Address, error) {
	return _ZEVMSwapApp.Contract.Router02(&_ZEVMSwapApp.CallOpts)
}

// SystemContract is a free data retrieval call binding the contract method 0xbb88b769.
//
// Solidity: function systemContract() view returns(address)
func (_ZEVMSwapApp *ZEVMSwapAppCaller) SystemContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZEVMSwapApp.contract.Call(opts, &out, "systemContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SystemContract is a free data retrieval call binding the contract method 0xbb88b769.
//
// Solidity: function systemContract() view returns(address)
func (_ZEVMSwapApp *ZEVMSwapAppSession) SystemContract() (common.Address, error) {
	return _ZEVMSwapApp.Contract.SystemContract(&_ZEVMSwapApp.CallOpts)
}

// SystemContract is a free data retrieval call binding the contract method 0xbb88b769.
//
// Solidity: function systemContract() view returns(address)
func (_ZEVMSwapApp *ZEVMSwapAppCallerSession) SystemContract() (common.Address, error) {
	return _ZEVMSwapApp.Contract.SystemContract(&_ZEVMSwapApp.CallOpts)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xc8522691.
//
// Solidity: function onCrossChainCall(address zrc20, uint256 amount, bytes message) returns()
func (_ZEVMSwapApp *ZEVMSwapAppTransactor) OnCrossChainCall(opts *bind.TransactOpts, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ZEVMSwapApp.contract.Transact(opts, "onCrossChainCall", zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xc8522691.
//
// Solidity: function onCrossChainCall(address zrc20, uint256 amount, bytes message) returns()
func (_ZEVMSwapApp *ZEVMSwapAppSession) OnCrossChainCall(zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ZEVMSwapApp.Contract.OnCrossChainCall(&_ZEVMSwapApp.TransactOpts, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xc8522691.
//
// Solidity: function onCrossChainCall(address zrc20, uint256 amount, bytes message) returns()
func (_ZEVMSwapApp *ZEVMSwapAppTransactorSession) OnCrossChainCall(zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ZEVMSwapApp.Contract.OnCrossChainCall(&_ZEVMSwapApp.TransactOpts, zrc20, amount, message)
}
