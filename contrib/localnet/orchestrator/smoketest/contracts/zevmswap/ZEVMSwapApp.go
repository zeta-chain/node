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
	_ = abi.ConvertType
)

// Context is an auto generated low-level Go binding around an user-defined struct.
type Context struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}

// ZEVMSwapAppMetaData contains all meta data concerning the ZEVMSwapApp contract.
var ZEVMSwapAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router02_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"systemContract_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LowAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"decodeMemo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"targetZRC20\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"}],\"name\":\"encodeMemo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structContext\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router02\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"systemContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b50604051620012ee380380620012ee833981810160405281019062000037919062000111565b8173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1681525050505062000158565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620000d982620000ac565b9050919050565b620000eb81620000cc565b8114620000f757600080fd5b50565b6000815190506200010b81620000e0565b92915050565b600080604083850312156200012b576200012a620000a7565b5b60006200013b85828601620000fa565b92505060206200014e85828601620000fa565b9150509250929050565b60805160a05161115b62000193600039600081816101b101526101f90152600081816101d50152818161039b0152610420015261115b6000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c8063a06ea8bc1461005c578063bb88b7691461008d578063bd00c9c4146100ab578063de43156e146100c9578063df73044e146100e5575b600080fd5b610076600480360381019061007191906107bc565b610115565b6040516100849291906108da565b60405180910390f35b6100956101af565b6040516100a2919061090a565b60405180910390f35b6100b36101d3565b6040516100c0919061090a565b60405180910390f35b6100e360048036038101906100de91906109ab565b6101f7565b005b6100ff60048036038101906100fa9190610a4f565b610714565b60405161010c9190610aaf565b60405180910390f35b6000606080600080868690509150868660009060149261013793929190610adb565b906101429190610b5a565b60601c90508686601490809261015a93929190610adb565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505092508083945094505050509250929050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461027c576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060608061028b8585610115565b8093508194505050600267ffffffffffffffff8111156102ae576102ad610bb9565b5b6040519080825280602002602001820160405280156102dc5781602001602082028036833780820191505090505b50905086816000815181106102f4576102f3610be8565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050828160018151811061034357610342610be8565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508673ffffffffffffffffffffffffffffffffffffffff1663095ea7b37f0000000000000000000000000000000000000000000000000000000000000000886040518363ffffffff1660e01b81526004016103d8929190610c26565b6020604051808303816000875af11580156103f7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061041b9190610c87565b5060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166338ed17398860008530680100000000000000006040518663ffffffff1660e01b8152600401610489959493929190610db7565b6000604051808303816000875af11580156104a8573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052508101906104d19190610f35565b90506000808573ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b81526004016040805180830381865afa158015610520573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105449190610f93565b915091508173ffffffffffffffffffffffffffffffffffffffff1663095ea7b387836040518363ffffffff1660e01b8152600401610583929190610c26565b6020604051808303816000875af11580156105a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105c69190610c87565b508573ffffffffffffffffffffffffffffffffffffffff1663095ea7b387600a866001815181106105fa576105f9610be8565b5b602002602001015161060c9190611002565b6040518363ffffffff1660e01b8152600401610629929190610c26565b6020604051808303816000875af1158015610648573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061066c9190610c87565b508573ffffffffffffffffffffffffffffffffffffffff1663c7012626868560018151811061069e5761069d610be8565b5b60200260200101516040518363ffffffff1660e01b81526004016106c3929190611044565b6020604051808303816000875af11580156106e2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107069190610c87565b505050505050505050505050565b606083838360405160200161072b939291906110fb565b60405160208183030381529060405290509392505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f84011261077c5761077b610757565b5b8235905067ffffffffffffffff8111156107995761079861075c565b5b6020830191508360018202830111156107b5576107b4610761565b5b9250929050565b600080602083850312156107d3576107d261074d565b5b600083013567ffffffffffffffff8111156107f1576107f0610752565b5b6107fd85828601610766565b92509250509250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061083482610809565b9050919050565b61084481610829565b82525050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610884578082015181840152602081019050610869565b60008484015250505050565b6000601f19601f8301169050919050565b60006108ac8261084a565b6108b68185610855565b93506108c6818560208601610866565b6108cf81610890565b840191505092915050565b60006040820190506108ef600083018561083b565b818103602083015261090181846108a1565b90509392505050565b600060208201905061091f600083018461083b565b92915050565b600080fd5b6000606082840312156109405761093f610925565b5b81905092915050565b61095281610829565b811461095d57600080fd5b50565b60008135905061096f81610949565b92915050565b6000819050919050565b61098881610975565b811461099357600080fd5b50565b6000813590506109a58161097f565b92915050565b6000806000806000608086880312156109c7576109c661074d565b5b600086013567ffffffffffffffff8111156109e5576109e4610752565b5b6109f18882890161092a565b9550506020610a0288828901610960565b9450506040610a1388828901610996565b935050606086013567ffffffffffffffff811115610a3457610a33610752565b5b610a4088828901610766565b92509250509295509295909350565b600080600060408486031215610a6857610a6761074d565b5b6000610a7686828701610960565b935050602084013567ffffffffffffffff811115610a9757610a96610752565b5b610aa386828701610766565b92509250509250925092565b60006020820190508181036000830152610ac981846108a1565b905092915050565b600080fd5b600080fd5b60008085851115610aef57610aee610ad1565b5b83861115610b0057610aff610ad6565b5b6001850283019150848603905094509492505050565b600082905092915050565b60007fffffffffffffffffffffffffffffffffffffffff00000000000000000000000082169050919050565b600082821b905092915050565b6000610b668383610b16565b82610b718135610b21565b92506014821015610bb157610bac7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000083601403600802610b4d565b831692505b505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b610c2081610975565b82525050565b6000604082019050610c3b600083018561083b565b610c486020830184610c17565b9392505050565b60008115159050919050565b610c6481610c4f565b8114610c6f57600080fd5b50565b600081519050610c8181610c5b565b92915050565b600060208284031215610c9d57610c9c61074d565b5b6000610cab84828501610c72565b91505092915050565b6000819050919050565b6000819050919050565b6000610ce3610cde610cd984610cb4565b610cbe565b610975565b9050919050565b610cf381610cc8565b82525050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b610d2e81610829565b82525050565b6000610d408383610d25565b60208301905092915050565b6000602082019050919050565b6000610d6482610cf9565b610d6e8185610d04565b9350610d7983610d15565b8060005b83811015610daa578151610d918882610d34565b9750610d9c83610d4c565b925050600181019050610d7d565b5085935050505092915050565b600060a082019050610dcc6000830188610c17565b610dd96020830187610cea565b8181036040830152610deb8186610d59565b9050610dfa606083018561083b565b610e076080830184610c17565b9695505050505050565b610e1a82610890565b810181811067ffffffffffffffff82111715610e3957610e38610bb9565b5b80604052505050565b6000610e4c610743565b9050610e588282610e11565b919050565b600067ffffffffffffffff821115610e7857610e77610bb9565b5b602082029050602081019050919050565b600081519050610e988161097f565b92915050565b6000610eb1610eac84610e5d565b610e42565b90508083825260208201905060208402830185811115610ed457610ed3610761565b5b835b81811015610efd5780610ee98882610e89565b845260208401935050602081019050610ed6565b5050509392505050565b600082601f830112610f1c57610f1b610757565b5b8151610f2c848260208601610e9e565b91505092915050565b600060208284031215610f4b57610f4a61074d565b5b600082015167ffffffffffffffff811115610f6957610f68610752565b5b610f7584828501610f07565b91505092915050565b600081519050610f8d81610949565b92915050565b60008060408385031215610faa57610fa961074d565b5b6000610fb885828601610f7e565b9250506020610fc985828601610e89565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061100d82610975565b915061101883610975565b925082820261102681610975565b9150828204841483151761103d5761103c610fd3565b5b5092915050565b6000604082019050818103600083015261105e81856108a1565b905061106d6020830184610c17565b9392505050565b60008160601b9050919050565b600061108c82611074565b9050919050565b600061109e82611081565b9050919050565b6110b66110b182610829565b611093565b82525050565b600081905092915050565b82818337600083830152505050565b60006110e283856110bc565b93506110ef8385846110c7565b82840190509392505050565b600061110782866110a5565b6014820191506111188284866110d6565b915081905094935050505056fea2646970667358221220e59a6599851c1ec787995687e8a35035524f5024bc92edd5c91793d9c6e77adf64736f6c63430008170033",
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
	parsed, err := ZEVMSwapAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// EncodeMemo is a free data retrieval call binding the contract method 0xdf73044e.
//
// Solidity: function encodeMemo(address targetZRC20, bytes recipient) pure returns(bytes)
func (_ZEVMSwapApp *ZEVMSwapAppCaller) EncodeMemo(opts *bind.CallOpts, targetZRC20 common.Address, recipient []byte) ([]byte, error) {
	var out []interface{}
	err := _ZEVMSwapApp.contract.Call(opts, &out, "encodeMemo", targetZRC20, recipient)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeMemo is a free data retrieval call binding the contract method 0xdf73044e.
//
// Solidity: function encodeMemo(address targetZRC20, bytes recipient) pure returns(bytes)
func (_ZEVMSwapApp *ZEVMSwapAppSession) EncodeMemo(targetZRC20 common.Address, recipient []byte) ([]byte, error) {
	return _ZEVMSwapApp.Contract.EncodeMemo(&_ZEVMSwapApp.CallOpts, targetZRC20, recipient)
}

// EncodeMemo is a free data retrieval call binding the contract method 0xdf73044e.
//
// Solidity: function encodeMemo(address targetZRC20, bytes recipient) pure returns(bytes)
func (_ZEVMSwapApp *ZEVMSwapAppCallerSession) EncodeMemo(targetZRC20 common.Address, recipient []byte) ([]byte, error) {
	return _ZEVMSwapApp.Contract.EncodeMemo(&_ZEVMSwapApp.CallOpts, targetZRC20, recipient)
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

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) , address zrc20, uint256 amount, bytes message) returns()
func (_ZEVMSwapApp *ZEVMSwapAppTransactor) OnCrossChainCall(opts *bind.TransactOpts, arg0 Context, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ZEVMSwapApp.contract.Transact(opts, "onCrossChainCall", arg0, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) , address zrc20, uint256 amount, bytes message) returns()
func (_ZEVMSwapApp *ZEVMSwapAppSession) OnCrossChainCall(arg0 Context, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ZEVMSwapApp.Contract.OnCrossChainCall(&_ZEVMSwapApp.TransactOpts, arg0, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) , address zrc20, uint256 amount, bytes message) returns()
func (_ZEVMSwapApp *ZEVMSwapAppTransactorSession) OnCrossChainCall(arg0 Context, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ZEVMSwapApp.Contract.OnCrossChainCall(&_ZEVMSwapApp.TransactOpts, arg0, zrc20, amount, message)
}
