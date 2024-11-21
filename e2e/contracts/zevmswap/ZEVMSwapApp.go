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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router02_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"systemContract_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LowAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"decodeMemo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"targetZRC20\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"}],\"name\":\"encodeMemo\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structContext\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structContext\",\"name\":\"\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router02\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"systemContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b506040516117a63803806117a683398181016040528101906100329190610104565b8173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff16815250505050610144565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100d1826100a6565b9050919050565b6100e1816100c6565b81146100ec57600080fd5b50565b6000815190506100fe816100d8565b92915050565b6000806040838503121561011b5761011a6100a1565b5b6000610129858286016100ef565b925050602061013a858286016100ef565b9150509250929050565b60805160a05161161a61018c6000396000818161067001526106b801526000818161025b015281816102e0015281816106940152818161085a01526108df015261161a6000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80635bcfd61614610067578063a06ea8bc14610083578063bb88b769146100b4578063bd00c9c4146100d2578063de43156e146100f0578063df73044e1461010c575b600080fd5b610081600480360381019061007c9190610d33565b61013c565b005b61009d60048036038101906100989190610dd7565b6105d4565b6040516100ab929190610ec3565b60405180910390f35b6100bc61066e565b6040516100c99190610ef3565b60405180910390f35b6100da610692565b6040516100e79190610ef3565b60405180910390f35b61010a60048036038101906101059190610d33565b6106b6565b005b61012660048036038101906101219190610f0e565b610bd3565b6040516101339190610f6e565b60405180910390f35b600060608061014b85856105d4565b8093508194505050600267ffffffffffffffff81111561016e5761016d610f90565b5b60405190808252806020026020018201604052801561019c5781602001602082028036833780820191505090505b50905086816000815181106101b4576101b3610fbf565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050828160018151811061020357610202610fbf565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508673ffffffffffffffffffffffffffffffffffffffff1663095ea7b37f0000000000000000000000000000000000000000000000000000000000000000886040518363ffffffff1660e01b8152600401610298929190610ffd565b6020604051808303816000875af11580156102b7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102db919061105e565b5060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166338ed17398860008530680100000000000000006040518663ffffffff1660e01b815260040161034995949392919061118e565b6000604051808303816000875af1158015610368573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250810190610391919061130c565b90506000808573ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b81526004016040805180830381865afa1580156103e0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610404919061136a565b915091508173ffffffffffffffffffffffffffffffffffffffff1663095ea7b387836040518363ffffffff1660e01b8152600401610443929190610ffd565b6020604051808303816000875af1158015610462573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610486919061105e565b508573ffffffffffffffffffffffffffffffffffffffff1663095ea7b387600a866001815181106104ba576104b9610fbf565b5b60200260200101516104cc91906113d9565b6040518363ffffffff1660e01b81526004016104e9929190610ffd565b6020604051808303816000875af1158015610508573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061052c919061105e565b508573ffffffffffffffffffffffffffffffffffffffff1663c7012626868560018151811061055e5761055d610fbf565b5b60200260200101516040518363ffffffff1660e01b815260040161058392919061141b565b6020604051808303816000875af11580156105a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105c6919061105e565b505050505050505050505050565b600060608060008086869050915086866000906014926105f693929190611455565b9061060191906114d4565b60601c90508686601490809261061993929190611455565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505092508083945094505050509250929050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461073b576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060608061074a85856105d4565b8093508194505050600267ffffffffffffffff81111561076d5761076c610f90565b5b60405190808252806020026020018201604052801561079b5781602001602082028036833780820191505090505b50905086816000815181106107b3576107b2610fbf565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050828160018151811061080257610801610fbf565b5b602002602001019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff16815250508673ffffffffffffffffffffffffffffffffffffffff1663095ea7b37f0000000000000000000000000000000000000000000000000000000000000000886040518363ffffffff1660e01b8152600401610897929190610ffd565b6020604051808303816000875af11580156108b6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108da919061105e565b5060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166338ed17398860008530680100000000000000006040518663ffffffff1660e01b815260040161094895949392919061118e565b6000604051808303816000875af1158015610967573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250810190610990919061130c565b90506000808573ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b81526004016040805180830381865afa1580156109df573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a03919061136a565b915091508173ffffffffffffffffffffffffffffffffffffffff1663095ea7b387836040518363ffffffff1660e01b8152600401610a42929190610ffd565b6020604051808303816000875af1158015610a61573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a85919061105e565b508573ffffffffffffffffffffffffffffffffffffffff1663095ea7b387600a86600181518110610ab957610ab8610fbf565b5b6020026020010151610acb91906113d9565b6040518363ffffffff1660e01b8152600401610ae8929190610ffd565b6020604051808303816000875af1158015610b07573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b2b919061105e565b508573ffffffffffffffffffffffffffffffffffffffff1663c70126268685600181518110610b5d57610b5c610fbf565b5b60200260200101516040518363ffffffff1660e01b8152600401610b8292919061141b565b6020604051808303816000875af1158015610ba1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bc5919061105e565b505050505050505050505050565b6060838383604051602001610bea939291906115ba565b60405160208183030381529060405290509392505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600060608284031215610c3157610c30610c16565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610c6582610c3a565b9050919050565b610c7581610c5a565b8114610c8057600080fd5b50565b600081359050610c9281610c6c565b92915050565b6000819050919050565b610cab81610c98565b8114610cb657600080fd5b50565b600081359050610cc881610ca2565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610cf357610cf2610cce565b5b8235905067ffffffffffffffff811115610d1057610d0f610cd3565b5b602083019150836001820283011115610d2c57610d2b610cd8565b5b9250929050565b600080600080600060808688031215610d4f57610d4e610c0c565b5b600086013567ffffffffffffffff811115610d6d57610d6c610c11565b5b610d7988828901610c1b565b9550506020610d8a88828901610c83565b9450506040610d9b88828901610cb9565b935050606086013567ffffffffffffffff811115610dbc57610dbb610c11565b5b610dc888828901610cdd565b92509250509295509295909350565b60008060208385031215610dee57610ded610c0c565b5b600083013567ffffffffffffffff811115610e0c57610e0b610c11565b5b610e1885828601610cdd565b92509250509250929050565b610e2d81610c5a565b82525050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610e6d578082015181840152602081019050610e52565b60008484015250505050565b6000601f19601f8301169050919050565b6000610e9582610e33565b610e9f8185610e3e565b9350610eaf818560208601610e4f565b610eb881610e79565b840191505092915050565b6000604082019050610ed86000830185610e24565b8181036020830152610eea8184610e8a565b90509392505050565b6000602082019050610f086000830184610e24565b92915050565b600080600060408486031215610f2757610f26610c0c565b5b6000610f3586828701610c83565b935050602084013567ffffffffffffffff811115610f5657610f55610c11565b5b610f6286828701610cdd565b92509250509250925092565b60006020820190508181036000830152610f888184610e8a565b905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b610ff781610c98565b82525050565b60006040820190506110126000830185610e24565b61101f6020830184610fee565b9392505050565b60008115159050919050565b61103b81611026565b811461104657600080fd5b50565b60008151905061105881611032565b92915050565b60006020828403121561107457611073610c0c565b5b600061108284828501611049565b91505092915050565b6000819050919050565b6000819050919050565b60006110ba6110b56110b08461108b565b611095565b610c98565b9050919050565b6110ca8161109f565b82525050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b61110581610c5a565b82525050565b600061111783836110fc565b60208301905092915050565b6000602082019050919050565b600061113b826110d0565b61114581856110db565b9350611150836110ec565b8060005b83811015611181578151611168888261110b565b975061117383611123565b925050600181019050611154565b5085935050505092915050565b600060a0820190506111a36000830188610fee565b6111b060208301876110c1565b81810360408301526111c28186611130565b90506111d16060830185610e24565b6111de6080830184610fee565b9695505050505050565b6111f182610e79565b810181811067ffffffffffffffff821117156112105761120f610f90565b5b80604052505050565b6000611223610c02565b905061122f82826111e8565b919050565b600067ffffffffffffffff82111561124f5761124e610f90565b5b602082029050602081019050919050565b60008151905061126f81610ca2565b92915050565b600061128861128384611234565b611219565b905080838252602082019050602084028301858111156112ab576112aa610cd8565b5b835b818110156112d457806112c08882611260565b8452602084019350506020810190506112ad565b5050509392505050565b600082601f8301126112f3576112f2610cce565b5b8151611303848260208601611275565b91505092915050565b60006020828403121561132257611321610c0c565b5b600082015167ffffffffffffffff8111156113405761133f610c11565b5b61134c848285016112de565b91505092915050565b60008151905061136481610c6c565b92915050565b6000806040838503121561138157611380610c0c565b5b600061138f85828601611355565b92505060206113a085828601611260565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006113e482610c98565b91506113ef83610c98565b92508282026113fd81610c98565b91508282048414831517611414576114136113aa565b5b5092915050565b600060408201905081810360008301526114358185610e8a565b90506114446020830184610fee565b9392505050565b600080fd5b600080fd5b600080858511156114695761146861144b565b5b8386111561147a57611479611450565b5b6001850283019150848603905094509492505050565b600082905092915050565b60007fffffffffffffffffffffffffffffffffffffffff00000000000000000000000082169050919050565b600082821b905092915050565b60006114e08383611490565b826114eb813561149b565b9250601482101561152b576115267fffffffffffffffffffffffffffffffffffffffff000000000000000000000000836014036008026114c7565b831692505b505092915050565b60008160601b9050919050565b600061154b82611533565b9050919050565b600061155d82611540565b9050919050565b61157561157082610c5a565b611552565b82525050565b600081905092915050565b82818337600083830152505050565b60006115a1838561157b565b93506115ae838584611586565b82840190509392505050565b60006115c68286611564565b6014820191506115d7828486611595565b915081905094935050505056fea2646970667358221220154bbadb87b49f8568829220c413c264d9405e5acb49ef710acf293adbd4f01564736f6c634300081a0033",
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

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) , address zrc20, uint256 amount, bytes message) returns()
func (_ZEVMSwapApp *ZEVMSwapAppTransactor) OnCall(opts *bind.TransactOpts, arg0 Context, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ZEVMSwapApp.contract.Transact(opts, "onCall", arg0, zrc20, amount, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) , address zrc20, uint256 amount, bytes message) returns()
func (_ZEVMSwapApp *ZEVMSwapAppSession) OnCall(arg0 Context, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ZEVMSwapApp.Contract.OnCall(&_ZEVMSwapApp.TransactOpts, arg0, zrc20, amount, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) , address zrc20, uint256 amount, bytes message) returns()
func (_ZEVMSwapApp *ZEVMSwapAppTransactorSession) OnCall(arg0 Context, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ZEVMSwapApp.Contract.OnCall(&_ZEVMSwapApp.TransactOpts, arg0, zrc20, amount, message)
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
