// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gatewayzevmcaller

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

// CallOptions is an auto generated low-level Go binding around an user-defined struct.
type CallOptions struct {
	GasLimit        *big.Int
	IsArbitraryCall bool
}

// RevertOptions is an auto generated low-level Go binding around an user-defined struct.
type RevertOptions struct {
	RevertAddress    common.Address
	CallOnRevert     bool
	AbortAddress     common.Address
	RevertMessage    []byte
	OnRevertGasLimit *big.Int
}

// GatewayZEVMCallerMetaData contains all meta data concerning the GatewayZEVMCaller contract.
var GatewayZEVMCallerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"gatewayZEVMAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"wzetaAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isArbitraryCall\",\"type\":\"bool\"}],\"internalType\":\"structCallOptions\",\"name\":\"callOptions\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"revertAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"callOnRevert\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"abortAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"onRevertGasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structRevertOptions\",\"name\":\"revertOptions\",\"type\":\"tuple\"}],\"name\":\"callGatewayZEVM\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositWZETA\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isArbitraryCall\",\"type\":\"bool\"}],\"internalType\":\"structCallOptions\",\"name\":\"callOptions\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"revertAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"callOnRevert\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"abortAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"onRevertGasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structRevertOptions\",\"name\":\"revertOptions\",\"type\":\"tuple\"}],\"name\":\"withdrawAndCallGatewayZEVM\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isArbitraryCall\",\"type\":\"bool\"}],\"internalType\":\"structCallOptions\",\"name\":\"callOptions\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"revertAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"callOnRevert\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"abortAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"onRevertGasLimit\",\"type\":\"uint256\"}],\"internalType\":\"structRevertOptions\",\"name\":\"revertOptions\",\"type\":\"tuple\"}],\"name\":\"withdrawAndCallGatewayZEVM\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620011613803806200116183398181016040528101906200003791906200012a565b816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505062000171565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620000f282620000c5565b9050919050565b6200010481620000e5565b81146200011057600080fd5b50565b6000815190506200012481620000f9565b92915050565b60008060408385031215620001445762000143620000c0565b5b6000620001548582860162000113565b9250506020620001678582860162000113565b9150509250929050565b610fe080620001816000396000f3fe60806040526004361061003f5760003560e01c806325859e62146100445780632c5d24ae1461006d57806362543ae714610077578063f66f4625146100a0575b600080fd5b34801561005057600080fd5b5061006b60048036038101906100669190610795565b6100c9565b005b61007561020d565b005b34801561008357600080fd5b5061009e6004803603810190610099919061089d565b610292565b005b3480156100ac57600080fd5b506100c760048036038101906100c29190610984565b6103d9565b005b8473ffffffffffffffffffffffffffffffffffffffff1663095ea7b360008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1667016345785d8a00006040518363ffffffff1660e01b815260040161012c929190610abf565b6020604051808303816000875af115801561014b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061016f9190610b20565b5060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166306cb89838787878787876040518763ffffffff1660e01b81526004016101d396959493929190610e18565b600060405180830381600087803b1580156101ed57600080fd5b505af1158015610201573d6000803e3d6000fd5b50505050505050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d0e30db0346040518263ffffffff1660e01b81526004016000604051808303818588803b15801561027757600080fd5b505af115801561028b573d6000803e3d6000fd5b5050505050565b8473ffffffffffffffffffffffffffffffffffffffff1663095ea7b360008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1667016345785d8a00006040518363ffffffff1660e01b81526004016102f5929190610abf565b6020604051808303816000875af1158015610314573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103389190610b20565b5060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16637b15118b888888888888886040518863ffffffff1660e01b815260040161039e9796959493929190610e91565b600060405180830381600087803b1580156103b857600080fd5b505af11580156103cc573d6000803e3d6000fd5b5050505050505050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663095ea7b360008054906101000a900473ffffffffffffffffffffffffffffffffffffffff16886040518363ffffffff1660e01b8152600401610456929190610f09565b6020604051808303816000875af1158015610475573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104999190610b20565b5060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632810ae63888888888888886040518863ffffffff1660e01b81526004016104ff9796959493929190610f32565b600060405180830381600087803b15801561051957600080fd5b505af115801561052d573d6000803e3d6000fd5b5050505050505050505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105a182610558565b810181811067ffffffffffffffff821117156105c0576105bf610569565b5b80604052505050565b60006105d361053a565b90506105df8282610598565b919050565b600067ffffffffffffffff8211156105ff576105fe610569565b5b61060882610558565b9050602081019050919050565b82818337600083830152505050565b6000610637610632846105e4565b6105c9565b90508281526020810184848401111561065357610652610553565b5b61065e848285610615565b509392505050565b600082601f83011261067b5761067a61054e565b5b813561068b848260208601610624565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006106bf82610694565b9050919050565b6106cf816106b4565b81146106da57600080fd5b50565b6000813590506106ec816106c6565b92915050565b600080fd5b600080fd5b60008083601f8401126107125761071161054e565b5b8235905067ffffffffffffffff81111561072f5761072e6106f2565b5b60208301915083600182028301111561074b5761074a6106f7565b5b9250929050565b600080fd5b60006040828403121561076d5761076c610752565b5b81905092915050565b600060a0828403121561078c5761078b610752565b5b81905092915050565b60008060008060008060c087890312156107b2576107b1610544565b5b600087013567ffffffffffffffff8111156107d0576107cf610549565b5b6107dc89828a01610666565b96505060206107ed89828a016106dd565b955050604087013567ffffffffffffffff81111561080e5761080d610549565b5b61081a89828a016106fc565b9450945050606061082d89828a01610757565b92505060a087013567ffffffffffffffff81111561084e5761084d610549565b5b61085a89828a01610776565b9150509295509295509295565b6000819050919050565b61087a81610867565b811461088557600080fd5b50565b60008135905061089781610871565b92915050565b600080600080600080600060e0888a0312156108bc576108bb610544565b5b600088013567ffffffffffffffff8111156108da576108d9610549565b5b6108e68a828b01610666565b97505060206108f78a828b01610888565b96505060406109088a828b016106dd565b955050606088013567ffffffffffffffff81111561092957610928610549565b5b6109358a828b016106fc565b945094505060806109488a828b01610757565b92505060c088013567ffffffffffffffff81111561096957610968610549565b5b6109758a828b01610776565b91505092959891949750929550565b600080600080600080600060e0888a0312156109a3576109a2610544565b5b600088013567ffffffffffffffff8111156109c1576109c0610549565b5b6109cd8a828b01610666565b97505060206109de8a828b01610888565b96505060406109ef8a828b01610888565b955050606088013567ffffffffffffffff811115610a1057610a0f610549565b5b610a1c8a828b016106fc565b94509450506080610a2f8a828b01610757565b92505060c088013567ffffffffffffffff811115610a5057610a4f610549565b5b610a5c8a828b01610776565b91505092959891949750929550565b610a74816106b4565b82525050565b6000819050919050565b6000819050919050565b6000610aa9610aa4610a9f84610a7a565b610a84565b610867565b9050919050565b610ab981610a8e565b82525050565b6000604082019050610ad46000830185610a6b565b610ae16020830184610ab0565b9392505050565b60008115159050919050565b610afd81610ae8565b8114610b0857600080fd5b50565b600081519050610b1a81610af4565b92915050565b600060208284031215610b3657610b35610544565b5b6000610b4484828501610b0b565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610b87578082015181840152602081019050610b6c565b83811115610b96576000848401525b50505050565b6000610ba782610b4d565b610bb18185610b58565b9350610bc1818560208601610b69565b610bca81610558565b840191505092915050565b6000610be18385610b58565b9350610bee838584610615565b610bf783610558565b840190509392505050565b6000610c116020840184610888565b905092915050565b610c2281610867565b82525050565b600081359050610c3781610af4565b92915050565b6000610c4c6020840184610c28565b905092915050565b610c5d81610ae8565b82525050565b60408201610c746000830183610c02565b610c816000850182610c19565b50610c8f6020830183610c3d565b610c9c6020850182610c54565b50505050565b6000610cb160208401846106dd565b905092915050565b610cc2816106b4565b82525050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112610cf457610cf3610cd2565b5b83810192508235915060208301925067ffffffffffffffff821115610d1c57610d1b610cc8565b5b600182023603841315610d3257610d31610ccd565b5b509250929050565b600082825260208201905092915050565b6000610d578385610d3a565b9350610d64838584610615565b610d6d83610558565b840190509392505050565b600060a08301610d8b6000840184610ca2565b610d986000860182610cb9565b50610da66020840184610c3d565b610db36020860182610c54565b50610dc16040840184610ca2565b610dce6040860182610cb9565b50610ddc6060840184610cd7565b8583036060870152610def838284610d4b565b92505050610e006080840184610c02565b610e0d6080860182610c19565b508091505092915050565b600060c0820190508181036000830152610e328189610b9c565b9050610e416020830188610a6b565b8181036040830152610e54818688610bd5565b9050610e636060830185610c63565b81810360a0830152610e758184610d78565b9050979650505050505050565b610e8b81610867565b82525050565b600060e0820190508181036000830152610eab818a610b9c565b9050610eba6020830189610e82565b610ec76040830188610a6b565b8181036060830152610eda818688610bd5565b9050610ee96080830185610c63565b81810360c0830152610efb8184610d78565b905098975050505050505050565b6000604082019050610f1e6000830185610a6b565b610f2b6020830184610e82565b9392505050565b600060e0820190508181036000830152610f4c818a610b9c565b9050610f5b6020830189610e82565b610f686040830188610e82565b8181036060830152610f7b818688610bd5565b9050610f8a6080830185610c63565b81810360c0830152610f9c8184610d78565b90509897505050505050505056fea26469706673582212204b70dead8f68839447fc08bd73294d57367a484ba9b048960ffcc202a737ae5664736f6c634300080a0033",
}

// GatewayZEVMCallerABI is the input ABI used to generate the binding from.
// Deprecated: Use GatewayZEVMCallerMetaData.ABI instead.
var GatewayZEVMCallerABI = GatewayZEVMCallerMetaData.ABI

// GatewayZEVMCallerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GatewayZEVMCallerMetaData.Bin instead.
var GatewayZEVMCallerBin = GatewayZEVMCallerMetaData.Bin

// DeployGatewayZEVMCaller deploys a new Ethereum contract, binding an instance of GatewayZEVMCaller to it.
func DeployGatewayZEVMCaller(auth *bind.TransactOpts, backend bind.ContractBackend, gatewayZEVMAddress common.Address, wzetaAddress common.Address) (common.Address, *types.Transaction, *GatewayZEVMCaller, error) {
	parsed, err := GatewayZEVMCallerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GatewayZEVMCallerBin), backend, gatewayZEVMAddress, wzetaAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GatewayZEVMCaller{GatewayZEVMCallerCaller: GatewayZEVMCallerCaller{contract: contract}, GatewayZEVMCallerTransactor: GatewayZEVMCallerTransactor{contract: contract}, GatewayZEVMCallerFilterer: GatewayZEVMCallerFilterer{contract: contract}}, nil
}

// GatewayZEVMCaller is an auto generated Go binding around an Ethereum contract.
type GatewayZEVMCaller struct {
	GatewayZEVMCallerCaller     // Read-only binding to the contract
	GatewayZEVMCallerTransactor // Write-only binding to the contract
	GatewayZEVMCallerFilterer   // Log filterer for contract events
}

// GatewayZEVMCallerCaller is an auto generated read-only Go binding around an Ethereum contract.
type GatewayZEVMCallerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GatewayZEVMCallerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GatewayZEVMCallerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GatewayZEVMCallerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GatewayZEVMCallerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GatewayZEVMCallerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GatewayZEVMCallerSession struct {
	Contract     *GatewayZEVMCaller // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// GatewayZEVMCallerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GatewayZEVMCallerCallerSession struct {
	Contract *GatewayZEVMCallerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// GatewayZEVMCallerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GatewayZEVMCallerTransactorSession struct {
	Contract     *GatewayZEVMCallerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// GatewayZEVMCallerRaw is an auto generated low-level Go binding around an Ethereum contract.
type GatewayZEVMCallerRaw struct {
	Contract *GatewayZEVMCaller // Generic contract binding to access the raw methods on
}

// GatewayZEVMCallerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GatewayZEVMCallerCallerRaw struct {
	Contract *GatewayZEVMCallerCaller // Generic read-only contract binding to access the raw methods on
}

// GatewayZEVMCallerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GatewayZEVMCallerTransactorRaw struct {
	Contract *GatewayZEVMCallerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGatewayZEVMCaller creates a new instance of GatewayZEVMCaller, bound to a specific deployed contract.
func NewGatewayZEVMCaller(address common.Address, backend bind.ContractBackend) (*GatewayZEVMCaller, error) {
	contract, err := bindGatewayZEVMCaller(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GatewayZEVMCaller{GatewayZEVMCallerCaller: GatewayZEVMCallerCaller{contract: contract}, GatewayZEVMCallerTransactor: GatewayZEVMCallerTransactor{contract: contract}, GatewayZEVMCallerFilterer: GatewayZEVMCallerFilterer{contract: contract}}, nil
}

// NewGatewayZEVMCallerCaller creates a new read-only instance of GatewayZEVMCaller, bound to a specific deployed contract.
func NewGatewayZEVMCallerCaller(address common.Address, caller bind.ContractCaller) (*GatewayZEVMCallerCaller, error) {
	contract, err := bindGatewayZEVMCaller(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GatewayZEVMCallerCaller{contract: contract}, nil
}

// NewGatewayZEVMCallerTransactor creates a new write-only instance of GatewayZEVMCaller, bound to a specific deployed contract.
func NewGatewayZEVMCallerTransactor(address common.Address, transactor bind.ContractTransactor) (*GatewayZEVMCallerTransactor, error) {
	contract, err := bindGatewayZEVMCaller(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GatewayZEVMCallerTransactor{contract: contract}, nil
}

// NewGatewayZEVMCallerFilterer creates a new log filterer instance of GatewayZEVMCaller, bound to a specific deployed contract.
func NewGatewayZEVMCallerFilterer(address common.Address, filterer bind.ContractFilterer) (*GatewayZEVMCallerFilterer, error) {
	contract, err := bindGatewayZEVMCaller(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GatewayZEVMCallerFilterer{contract: contract}, nil
}

// bindGatewayZEVMCaller binds a generic wrapper to an already deployed contract.
func bindGatewayZEVMCaller(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GatewayZEVMCallerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GatewayZEVMCaller *GatewayZEVMCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GatewayZEVMCaller.Contract.GatewayZEVMCallerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GatewayZEVMCaller *GatewayZEVMCallerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.GatewayZEVMCallerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GatewayZEVMCaller *GatewayZEVMCallerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.GatewayZEVMCallerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GatewayZEVMCaller *GatewayZEVMCallerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GatewayZEVMCaller.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GatewayZEVMCaller *GatewayZEVMCallerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GatewayZEVMCaller *GatewayZEVMCallerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.contract.Transact(opts, method, params...)
}

// CallGatewayZEVM is a paid mutator transaction binding the contract method 0x25859e62.
//
// Solidity: function callGatewayZEVM(bytes receiver, address zrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerTransactor) CallGatewayZEVM(opts *bind.TransactOpts, receiver []byte, zrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayZEVMCaller.contract.Transact(opts, "callGatewayZEVM", receiver, zrc20, message, callOptions, revertOptions)
}

// CallGatewayZEVM is a paid mutator transaction binding the contract method 0x25859e62.
//
// Solidity: function callGatewayZEVM(bytes receiver, address zrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerSession) CallGatewayZEVM(receiver []byte, zrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.CallGatewayZEVM(&_GatewayZEVMCaller.TransactOpts, receiver, zrc20, message, callOptions, revertOptions)
}

// CallGatewayZEVM is a paid mutator transaction binding the contract method 0x25859e62.
//
// Solidity: function callGatewayZEVM(bytes receiver, address zrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerTransactorSession) CallGatewayZEVM(receiver []byte, zrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.CallGatewayZEVM(&_GatewayZEVMCaller.TransactOpts, receiver, zrc20, message, callOptions, revertOptions)
}

// DepositWZETA is a paid mutator transaction binding the contract method 0x2c5d24ae.
//
// Solidity: function depositWZETA() payable returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerTransactor) DepositWZETA(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GatewayZEVMCaller.contract.Transact(opts, "depositWZETA")
}

// DepositWZETA is a paid mutator transaction binding the contract method 0x2c5d24ae.
//
// Solidity: function depositWZETA() payable returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerSession) DepositWZETA() (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.DepositWZETA(&_GatewayZEVMCaller.TransactOpts)
}

// DepositWZETA is a paid mutator transaction binding the contract method 0x2c5d24ae.
//
// Solidity: function depositWZETA() payable returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerTransactorSession) DepositWZETA() (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.DepositWZETA(&_GatewayZEVMCaller.TransactOpts)
}

// WithdrawAndCallGatewayZEVM is a paid mutator transaction binding the contract method 0x62543ae7.
//
// Solidity: function withdrawAndCallGatewayZEVM(bytes receiver, uint256 amount, address zrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerTransactor) WithdrawAndCallGatewayZEVM(opts *bind.TransactOpts, receiver []byte, amount *big.Int, zrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayZEVMCaller.contract.Transact(opts, "withdrawAndCallGatewayZEVM", receiver, amount, zrc20, message, callOptions, revertOptions)
}

// WithdrawAndCallGatewayZEVM is a paid mutator transaction binding the contract method 0x62543ae7.
//
// Solidity: function withdrawAndCallGatewayZEVM(bytes receiver, uint256 amount, address zrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerSession) WithdrawAndCallGatewayZEVM(receiver []byte, amount *big.Int, zrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.WithdrawAndCallGatewayZEVM(&_GatewayZEVMCaller.TransactOpts, receiver, amount, zrc20, message, callOptions, revertOptions)
}

// WithdrawAndCallGatewayZEVM is a paid mutator transaction binding the contract method 0x62543ae7.
//
// Solidity: function withdrawAndCallGatewayZEVM(bytes receiver, uint256 amount, address zrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerTransactorSession) WithdrawAndCallGatewayZEVM(receiver []byte, amount *big.Int, zrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.WithdrawAndCallGatewayZEVM(&_GatewayZEVMCaller.TransactOpts, receiver, amount, zrc20, message, callOptions, revertOptions)
}

// WithdrawAndCallGatewayZEVM0 is a paid mutator transaction binding the contract method 0xf66f4625.
//
// Solidity: function withdrawAndCallGatewayZEVM(bytes receiver, uint256 amount, uint256 chainId, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerTransactor) WithdrawAndCallGatewayZEVM0(opts *bind.TransactOpts, receiver []byte, amount *big.Int, chainId *big.Int, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayZEVMCaller.contract.Transact(opts, "withdrawAndCallGatewayZEVM0", receiver, amount, chainId, message, callOptions, revertOptions)
}

// WithdrawAndCallGatewayZEVM0 is a paid mutator transaction binding the contract method 0xf66f4625.
//
// Solidity: function withdrawAndCallGatewayZEVM(bytes receiver, uint256 amount, uint256 chainId, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerSession) WithdrawAndCallGatewayZEVM0(receiver []byte, amount *big.Int, chainId *big.Int, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.WithdrawAndCallGatewayZEVM0(&_GatewayZEVMCaller.TransactOpts, receiver, amount, chainId, message, callOptions, revertOptions)
}

// WithdrawAndCallGatewayZEVM0 is a paid mutator transaction binding the contract method 0xf66f4625.
//
// Solidity: function withdrawAndCallGatewayZEVM(bytes receiver, uint256 amount, uint256 chainId, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayZEVMCaller *GatewayZEVMCallerTransactorSession) WithdrawAndCallGatewayZEVM0(receiver []byte, amount *big.Int, chainId *big.Int, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayZEVMCaller.Contract.WithdrawAndCallGatewayZEVM0(&_GatewayZEVMCaller.TransactOpts, receiver, amount, chainId, message, callOptions, revertOptions)
}
