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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router02_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"systemContract_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LowAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"decodeMemo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"targetZRC20\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"minAmountOut\",\"type\":\"uint256\"}],\"name\":\"encodeMemo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router02\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"systemContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162001269380380620012698339818101604052810190620000379190620000c4565b8173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff1660601b815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1660601b8152505050506200015e565b600081519050620000be8162000144565b92915050565b60008060408385031215620000de57620000dd6200013f565b5b6000620000ee85828601620000ad565b92505060206200010185828601620000ad565b9150509250929050565b600062000118826200011f565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600080fd5b6200014f816200010b565b81146200015b57600080fd5b50565b60805160601c60a05160601c6110ca6200019f600039600081816101e4015261022c015260008181610208015281816103cf015261046301526110ca6000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c8063678f75dc1461005c578063a06ea8bc1461008c578063bb88b769146100bd578063bd00c9c4146100db578063c8522691146100f9575b600080fd5b610076600480360381019061007191906108a6565b610115565b6040516100839190610c5d565b60405180910390f35b6100a660048036038101906100a19190610a44565b610147565b6040516100b4929190610c04565b60405180910390f35b6100c56101e2565b6040516100d29190610ba9565b60405180910390f35b6100e3610206565b6040516100f09190610ba9565b60405180910390f35b610113600480360381019061010e919061095a565b61022a565b005b60608484848460405160200161012e9493929190610bc4565b6040516020818303038152906040529050949350505050565b600060608060008585905090506000868660009060149261016a93929190610dba565b906101759190610e9d565b60601c90508686601490809261018d93929190610dba565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505092508083945094505050509250929050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102af576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060606102bd8484610147565b80925081935050506060600267ffffffffffffffff8111156102e2576102e1610fdf565b5b6040519080825280602002602001820160405280156103105781602001602082028036833780820191505090505b509050868160008151811061032857610327610fb0565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050828160018151811061037757610376610fb0565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508673ffffffffffffffffffffffffffffffffffffffff1663095ea7b37f0000000000000000000000000000000000000000000000000000000000000000886040518363ffffffff1660e01b815260040161040c929190610c34565b602060405180830381600087803b15801561042657600080fd5b505af115801561043a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061045e9190610a17565b5060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166338ed17398860008530680100000000000000006040518663ffffffff1660e01b81526004016104cc959493929190610caf565b600060405180830381600087803b1580156104e657600080fd5b505af11580156104fa573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525081019061052391906109ce565b905060008473ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b8152600401604080518083038186803b15801561056c57600080fd5b505afa158015610580573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105a4919061091a565b915050816001815181106105bb576105ba610fb0565b5b60200260200101518111156105fc576040517f4e8ed22600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8473ffffffffffffffffffffffffffffffffffffffff1663095ea7b386836040518363ffffffff1660e01b8152600401610637929190610c34565b602060405180830381600087803b15801561065157600080fd5b505af1158015610665573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106899190610a17565b508473ffffffffffffffffffffffffffffffffffffffff1663c70126268583856001815181106106bc576106bb610fb0565b5b60200260200101516106ce9190610df5565b6040518363ffffffff1660e01b81526004016106eb929190610c7f565b602060405180830381600087803b15801561070557600080fd5b505af1158015610719573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061073d9190610a17565b50505050505050505050565b600061075c61075784610d2e565b610d09565b9050808382526020820190508285602086028201111561077f5761077e611022565b5b60005b858110156107af57816107958882610891565b845260208401935060208301925050600181019050610782565b5050509392505050565b6000813590506107c88161104f565b92915050565b6000815190506107dd8161104f565b92915050565b600082601f8301126107f8576107f7611013565b5b8151610808848260208601610749565b91505092915050565b60008151905061082081611066565b92915050565b60008083601f84011261083c5761083b611013565b5b8235905067ffffffffffffffff8111156108595761085861100e565b5b60208301915083600182028301111561087557610874611022565b5b9250929050565b60008135905061088b8161107d565b92915050565b6000815190506108a08161107d565b92915050565b600080600080606085870312156108c0576108bf61102c565b5b60006108ce878288016107b9565b945050602085013567ffffffffffffffff8111156108ef576108ee611027565b5b6108fb87828801610826565b9350935050604061090e8782880161087c565b91505092959194509250565b600080604083850312156109315761093061102c565b5b600061093f858286016107ce565b925050602061095085828601610891565b9150509250929050565b600080600080606085870312156109745761097361102c565b5b6000610982878288016107b9565b94505060206109938782880161087c565b935050604085013567ffffffffffffffff8111156109b4576109b3611027565b5b6109c087828801610826565b925092505092959194509250565b6000602082840312156109e4576109e361102c565b5b600082015167ffffffffffffffff811115610a0257610a01611027565b5b610a0e848285016107e3565b91505092915050565b600060208284031215610a2d57610a2c61102c565b5b6000610a3b84828501610811565b91505092915050565b60008060208385031215610a5b57610a5a61102c565b5b600083013567ffffffffffffffff811115610a7957610a78611027565b5b610a8585828601610826565b92509250509250929050565b6000610a9d8383610aa9565b60208301905092915050565b610ab281610e29565b82525050565b610ac181610e29565b82525050565b6000610ad282610d6a565b610adc8185610d98565b9350610ae783610d5a565b8060005b83811015610b18578151610aff8882610a91565b9750610b0a83610d8b565b925050600181019050610aeb565b5085935050505092915050565b6000610b318385610da9565b9350610b3e838584610f0e565b610b4783611031565b840190509392505050565b6000610b5d82610d80565b610b678185610da9565b9350610b77818560208601610f1d565b610b8081611031565b840191505092915050565b610b9481610efc565b82525050565b610ba381610e93565b82525050565b6000602082019050610bbe6000830184610ab8565b92915050565b6000606082019050610bd96000830187610ab8565b8181036020830152610bec818587610b25565b9050610bfb6040830184610b9a565b95945050505050565b6000604082019050610c196000830185610ab8565b8181036020830152610c2b8184610b52565b90509392505050565b6000604082019050610c496000830185610ab8565b610c566020830184610b9a565b9392505050565b60006020820190508181036000830152610c778184610b52565b905092915050565b60006040820190508181036000830152610c998185610b52565b9050610ca86020830184610b9a565b9392505050565b600060a082019050610cc46000830188610b9a565b610cd16020830187610b8b565b8181036040830152610ce38186610ac7565b9050610cf26060830185610ab8565b610cff6080830184610b9a565b9695505050505050565b6000610d13610d24565b9050610d1f8282610f50565b919050565b6000604051905090565b600067ffffffffffffffff821115610d4957610d48610fdf565b5b602082029050602081019050919050565b6000819050602082019050919050565b600081519050919050565b600082905092915050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b60008085851115610dce57610dcd61101d565b5b83861115610ddf57610dde611018565b5b6001850283019150848603905094509492505050565b6000610e0082610e93565b9150610e0b83610e93565b925082821015610e1e57610e1d610f81565b5b828203905092915050565b6000610e3482610e73565b9050919050565b60008115159050919050565b60007fffffffffffffffffffffffffffffffffffffffff00000000000000000000000082169050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b6000610ea98383610d75565b82610eb48135610e47565b92506014821015610ef457610eef7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000083601403600802611042565b831692505b505092915050565b6000610f0782610e93565b9050919050565b82818337600083830152505050565b60005b83811015610f3b578082015181840152602081019050610f20565b83811115610f4a576000848401525b50505050565b610f5982611031565b810181811067ffffffffffffffff82111715610f7857610f77610fdf565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b600082821b905092915050565b61105881610e29565b811461106357600080fd5b50565b61106f81610e3b565b811461107a57600080fd5b50565b61108681610e93565b811461109157600080fd5b5056fea2646970667358221220927360e88db3660053aa80817fc73c4c86ba3c8d7c976c49976f970c66ddcc5464736f6c63430008070033",
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

// DecodeMemo is a free data retrieval call binding the contract method 0xa06ea8bc.
//
// Solidity: function decodeMemo(bytes data) pure returns(address, bytes)
func (_ZEVMSwapApp *ZEVMSwapAppCaller) DecodeMemo(opts *bind.CallOpts, data []byte) (common.Address, []byte, error) {
	var out []interface{}
	err := _ZEVMSwapApp.contract.Call(opts, &out, "decodeMemo", data)

	if err != nil {
		return *new(common.Address), *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return out0, out1, err

}

// DecodeMemo is a free data retrieval call binding the contract method 0xa06ea8bc.
//
// Solidity: function decodeMemo(bytes data) pure returns(address, bytes)
func (_ZEVMSwapApp *ZEVMSwapAppSession) DecodeMemo(data []byte) (common.Address, []byte, error) {
	return _ZEVMSwapApp.Contract.DecodeMemo(&_ZEVMSwapApp.CallOpts, data)
}

// DecodeMemo is a free data retrieval call binding the contract method 0xa06ea8bc.
//
// Solidity: function decodeMemo(bytes data) pure returns(address, bytes)
func (_ZEVMSwapApp *ZEVMSwapAppCallerSession) DecodeMemo(data []byte) (common.Address, []byte, error) {
	return _ZEVMSwapApp.Contract.DecodeMemo(&_ZEVMSwapApp.CallOpts, data)
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
