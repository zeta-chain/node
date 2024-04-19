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
)

// ZetaInterfacesZetaMessage is an auto generated low-level Go binding around an user-defined struct.
type ZetaInterfacesZetaMessage struct {
	ZetaTxSenderAddress []byte
	SourceChainId       *big.Int
	DestinationAddress  common.Address
	ZetaValue           *big.Int
	Message             []byte
}

// ZetaInterfacesZetaRevert is an auto generated low-level Go binding around an user-defined struct.
type ZetaInterfacesZetaRevert struct {
	ZetaTxSenderAddress common.Address
	SourceChainId       *big.Int
	DestinationAddress  []byte
	DestinationChainId  *big.Int
	RemainingZetaValue  *big.Int
	Message             []byte
}

// DappMetaData contains all meta data concerning the Dapp contract.
var DappMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"destinationAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"destinationChainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"message\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"zetaTxSenderAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"destinationAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"zetaValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"internalType\":\"structZetaInterfaces.ZetaMessage\",\"name\":\"zetaMessage\",\"type\":\"tuple\"}],\"name\":\"onZetaMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"zetaTxSenderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destinationAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingZetaValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"internalType\":\"structZetaInterfaces.ZetaRevert\",\"name\":\"zetaRevert\",\"type\":\"tuple\"}],\"name\":\"onZetaRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sourceChainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zetaTxSenderAddress\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zetaValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051806020016040528060008152506000908161002f91906102fe565b5060006001819055506000600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600381905550600060048190555060405180602001604052806000815250600590816100a891906102fe565b506103d0565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061012f57607f821691505b602082108103610142576101416100e8565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026101aa7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8261016d565b6101b4868361016d565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b60006101fb6101f66101f1846101cc565b6101d6565b6101cc565b9050919050565b6000819050919050565b610215836101e0565b61022961022182610202565b84845461017a565b825550505050565b600090565b61023e610231565b61024981848461020c565b505050565b5b8181101561026d57610262600082610236565b60018101905061024f565b5050565b601f8211156102b25761028381610148565b61028c8461015d565b8101602085101561029b578190505b6102af6102a78561015d565b83018261024e565b50505b505050565b600082821c905092915050565b60006102d5600019846008026102b7565b1980831691505092915050565b60006102ee83836102c4565b9150826002028217905092915050565b610307826100ae565b67ffffffffffffffff8111156103205761031f6100b9565b5b61032a8254610117565b610335828285610271565b600060209050601f8311600181146103685760008415610356578287015190505b61036085826102e2565b8655506103c8565b601f19841661037686610148565b60005b8281101561039e57848901518255600182019150602085019450602081019050610379565b868310156103bb57848901516103b7601f8916826102c4565b8355505b6001600288020188555050505b505050505050565b610c2c806103df6000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063b07506111161005b578063b075061114610101578063ca3254691461011f578063e21f37ce1461013d578063ed6b866b1461015b57610088565b8063050337a21461008d5780631544298e146100ab5780633749c51a146100c95780633ff0693c146100e5575b600080fd5b610095610179565b6040516100a29190610480565b60405180910390f35b6100b361017f565b6040516100c09190610480565b60405180910390f35b6100e360048036038101906100de91906104c9565b610185565b005b6100ff60048036038101906100fa9190610531565b610231565b005b61010961031f565b6040516101169190610480565b60405180910390f35b610127610325565b60405161013491906105bb565b60405180910390f35b61014561034b565b6040516101529190610666565b60405180910390f35b6101636103d9565b6040516101709190610666565b60405180910390f35b60045481565b60015481565b8080600001906101959190610697565b600091826101a4929190610940565b5080602001356001819055508060400160208101906101c39190610a3c565b600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550806060013560048190555080806080019061021e9190610697565b6005918261022d929190610940565b5050565b8060000160208101906102449190610a3c565b6040516020016102549190610ab1565b604051602081830303815290604052600090816102719190610acc565b50806020013560018190555080806040019061028d9190610697565b60405161029b929190610bdd565b604051809103902060001c600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080606001356003819055508060800135600481905550808060a0019061030c9190610697565b6005918261031b929190610940565b5050565b60035481565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6005805461035890610763565b80601f016020809104026020016040519081016040528092919081815260200182805461038490610763565b80156103d15780601f106103a6576101008083540402835291602001916103d1565b820191906000526020600020905b8154815290600101906020018083116103b457829003601f168201915b505050505081565b600080546103e690610763565b80601f016020809104026020016040519081016040528092919081815260200182805461041290610763565b801561045f5780601f106104345761010080835404028352916020019161045f565b820191906000526020600020905b81548152906001019060200180831161044257829003601f168201915b505050505081565b6000819050919050565b61047a81610467565b82525050565b60006020820190506104956000830184610471565b92915050565b600080fd5b600080fd5b600080fd5b600060a082840312156104c0576104bf6104a5565b5b81905092915050565b6000602082840312156104df576104de61049b565b5b600082013567ffffffffffffffff8111156104fd576104fc6104a0565b5b610509848285016104aa565b91505092915050565b600060c08284031215610528576105276104a5565b5b81905092915050565b6000602082840312156105475761054661049b565b5b600082013567ffffffffffffffff811115610565576105646104a0565b5b61057184828501610512565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006105a58261057a565b9050919050565b6105b58161059a565b82525050565b60006020820190506105d060008301846105ac565b92915050565b600081519050919050565b600082825260208201905092915050565b60005b838110156106105780820151818401526020810190506105f5565b60008484015250505050565b6000601f19601f8301169050919050565b6000610638826105d6565b61064281856105e1565b93506106528185602086016105f2565b61065b8161061c565b840191505092915050565b60006020820190508181036000830152610680818461062d565b905092915050565b600080fd5b600080fd5b600080fd5b600080833560016020038436030381126106b4576106b3610688565b5b80840192508235915067ffffffffffffffff8211156106d6576106d561068d565b5b6020830192506001820236038313156106f2576106f1610692565b5b509250929050565b600082905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061077b57607f821691505b60208210810361078e5761078d610734565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026107f67fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826107b9565b61080086836107b9565b95508019841693508086168417925050509392505050565b6000819050919050565b600061083d61083861083384610467565b610818565b610467565b9050919050565b6000819050919050565b61085783610822565b61086b61086382610844565b8484546107c6565b825550505050565b600090565b610880610873565b61088b81848461084e565b505050565b5b818110156108af576108a4600082610878565b600181019050610891565b5050565b601f8211156108f4576108c581610794565b6108ce846107a9565b810160208510156108dd578190505b6108f16108e9856107a9565b830182610890565b50505b505050565b600082821c905092915050565b6000610917600019846008026108f9565b1980831691505092915050565b60006109308383610906565b9150826002028217905092915050565b61094a83836106fa565b67ffffffffffffffff81111561096357610962610705565b5b61096d8254610763565b6109788282856108b3565b6000601f8311600181146109a75760008415610995578287013590505b61099f8582610924565b865550610a07565b601f1984166109b586610794565b60005b828110156109dd578489013582556001820191506020850194506020810190506109b8565b868310156109fa57848901356109f6601f891682610906565b8355505b6001600288020188555050505b50505050505050565b610a198161059a565b8114610a2457600080fd5b50565b600081359050610a3681610a10565b92915050565b600060208284031215610a5257610a5161049b565b5b6000610a6084828501610a27565b91505092915050565b60008160601b9050919050565b6000610a8182610a69565b9050919050565b6000610a9382610a76565b9050919050565b610aab610aa68261059a565b610a88565b82525050565b6000610abd8284610a9a565b60148201915081905092915050565b610ad5826105d6565b67ffffffffffffffff811115610aee57610aed610705565b5b610af88254610763565b610b038282856108b3565b600060209050601f831160018114610b365760008415610b24578287015190505b610b2e8582610924565b865550610b96565b601f198416610b4486610794565b60005b82811015610b6c57848901518255600182019150602085019450602081019050610b47565b86831015610b895784890151610b85601f891682610906565b8355505b6001600288020188555050505b505050505050565b600081905092915050565b82818337600083830152505050565b6000610bc48385610b9e565b9350610bd1838584610ba9565b82840190509392505050565b6000610bea828486610bb8565b9150819050939250505056fea26469706673582212207c5ed9b805f7e22d563799330b1c6e310dd5b1625dc35c9598846e17a8a5686664736f6c63430008190033",
}

// DappABI is the input ABI used to generate the binding from.
// Deprecated: Use DappMetaData.ABI instead.
var DappABI = DappMetaData.ABI

// DappBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DappMetaData.Bin instead.
var DappBin = DappMetaData.Bin

// DeployDapp deploys a new Ethereum contract, binding an instance of Dapp to it.
func DeployDapp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Dapp, error) {
	parsed, err := DappMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DappBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Dapp{DappCaller: DappCaller{contract: contract}, DappTransactor: DappTransactor{contract: contract}, DappFilterer: DappFilterer{contract: contract}}, nil
}

// Dapp is an auto generated Go binding around an Ethereum contract.
type Dapp struct {
	DappCaller     // Read-only binding to the contract
	DappTransactor // Write-only binding to the contract
	DappFilterer   // Log filterer for contract events
}

// DappCaller is an auto generated read-only Go binding around an Ethereum contract.
type DappCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DappTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DappTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DappFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DappFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DappSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DappSession struct {
	Contract     *Dapp             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DappCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DappCallerSession struct {
	Contract *DappCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DappTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DappTransactorSession struct {
	Contract     *DappTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DappRaw is an auto generated low-level Go binding around an Ethereum contract.
type DappRaw struct {
	Contract *Dapp // Generic contract binding to access the raw methods on
}

// DappCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DappCallerRaw struct {
	Contract *DappCaller // Generic read-only contract binding to access the raw methods on
}

// DappTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DappTransactorRaw struct {
	Contract *DappTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDapp creates a new instance of Dapp, bound to a specific deployed contract.
func NewDapp(address common.Address, backend bind.ContractBackend) (*Dapp, error) {
	contract, err := bindDapp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Dapp{DappCaller: DappCaller{contract: contract}, DappTransactor: DappTransactor{contract: contract}, DappFilterer: DappFilterer{contract: contract}}, nil
}

// NewDappCaller creates a new read-only instance of Dapp, bound to a specific deployed contract.
func NewDappCaller(address common.Address, caller bind.ContractCaller) (*DappCaller, error) {
	contract, err := bindDapp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DappCaller{contract: contract}, nil
}

// NewDappTransactor creates a new write-only instance of Dapp, bound to a specific deployed contract.
func NewDappTransactor(address common.Address, transactor bind.ContractTransactor) (*DappTransactor, error) {
	contract, err := bindDapp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DappTransactor{contract: contract}, nil
}

// NewDappFilterer creates a new log filterer instance of Dapp, bound to a specific deployed contract.
func NewDappFilterer(address common.Address, filterer bind.ContractFilterer) (*DappFilterer, error) {
	contract, err := bindDapp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DappFilterer{contract: contract}, nil
}

// bindDapp binds a generic wrapper to an already deployed contract.
func bindDapp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DappABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dapp *DappRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dapp.Contract.DappCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dapp *DappRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dapp.Contract.DappTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dapp *DappRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dapp.Contract.DappTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dapp *DappCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dapp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dapp *DappTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dapp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dapp *DappTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dapp.Contract.contract.Transact(opts, method, params...)
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() view returns(address)
func (_Dapp *DappCaller) DestinationAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Dapp.contract.Call(opts, &out, "destinationAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() view returns(address)
func (_Dapp *DappSession) DestinationAddress() (common.Address, error) {
	return _Dapp.Contract.DestinationAddress(&_Dapp.CallOpts)
}

// DestinationAddress is a free data retrieval call binding the contract method 0xca325469.
//
// Solidity: function destinationAddress() view returns(address)
func (_Dapp *DappCallerSession) DestinationAddress() (common.Address, error) {
	return _Dapp.Contract.DestinationAddress(&_Dapp.CallOpts)
}

// DestinationChainId is a free data retrieval call binding the contract method 0xb0750611.
//
// Solidity: function destinationChainId() view returns(uint256)
func (_Dapp *DappCaller) DestinationChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Dapp.contract.Call(opts, &out, "destinationChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DestinationChainId is a free data retrieval call binding the contract method 0xb0750611.
//
// Solidity: function destinationChainId() view returns(uint256)
func (_Dapp *DappSession) DestinationChainId() (*big.Int, error) {
	return _Dapp.Contract.DestinationChainId(&_Dapp.CallOpts)
}

// DestinationChainId is a free data retrieval call binding the contract method 0xb0750611.
//
// Solidity: function destinationChainId() view returns(uint256)
func (_Dapp *DappCallerSession) DestinationChainId() (*big.Int, error) {
	return _Dapp.Contract.DestinationChainId(&_Dapp.CallOpts)
}

// Message is a free data retrieval call binding the contract method 0xe21f37ce.
//
// Solidity: function message() view returns(bytes)
func (_Dapp *DappCaller) Message(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Dapp.contract.Call(opts, &out, "message")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Message is a free data retrieval call binding the contract method 0xe21f37ce.
//
// Solidity: function message() view returns(bytes)
func (_Dapp *DappSession) Message() ([]byte, error) {
	return _Dapp.Contract.Message(&_Dapp.CallOpts)
}

// Message is a free data retrieval call binding the contract method 0xe21f37ce.
//
// Solidity: function message() view returns(bytes)
func (_Dapp *DappCallerSession) Message() ([]byte, error) {
	return _Dapp.Contract.Message(&_Dapp.CallOpts)
}

// SourceChainId is a free data retrieval call binding the contract method 0x1544298e.
//
// Solidity: function sourceChainId() view returns(uint256)
func (_Dapp *DappCaller) SourceChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Dapp.contract.Call(opts, &out, "sourceChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SourceChainId is a free data retrieval call binding the contract method 0x1544298e.
//
// Solidity: function sourceChainId() view returns(uint256)
func (_Dapp *DappSession) SourceChainId() (*big.Int, error) {
	return _Dapp.Contract.SourceChainId(&_Dapp.CallOpts)
}

// SourceChainId is a free data retrieval call binding the contract method 0x1544298e.
//
// Solidity: function sourceChainId() view returns(uint256)
func (_Dapp *DappCallerSession) SourceChainId() (*big.Int, error) {
	return _Dapp.Contract.SourceChainId(&_Dapp.CallOpts)
}

// ZetaTxSenderAddress is a free data retrieval call binding the contract method 0xed6b866b.
//
// Solidity: function zetaTxSenderAddress() view returns(bytes)
func (_Dapp *DappCaller) ZetaTxSenderAddress(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Dapp.contract.Call(opts, &out, "zetaTxSenderAddress")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ZetaTxSenderAddress is a free data retrieval call binding the contract method 0xed6b866b.
//
// Solidity: function zetaTxSenderAddress() view returns(bytes)
func (_Dapp *DappSession) ZetaTxSenderAddress() ([]byte, error) {
	return _Dapp.Contract.ZetaTxSenderAddress(&_Dapp.CallOpts)
}

// ZetaTxSenderAddress is a free data retrieval call binding the contract method 0xed6b866b.
//
// Solidity: function zetaTxSenderAddress() view returns(bytes)
func (_Dapp *DappCallerSession) ZetaTxSenderAddress() ([]byte, error) {
	return _Dapp.Contract.ZetaTxSenderAddress(&_Dapp.CallOpts)
}

// ZetaValue is a free data retrieval call binding the contract method 0x050337a2.
//
// Solidity: function zetaValue() view returns(uint256)
func (_Dapp *DappCaller) ZetaValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Dapp.contract.Call(opts, &out, "zetaValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ZetaValue is a free data retrieval call binding the contract method 0x050337a2.
//
// Solidity: function zetaValue() view returns(uint256)
func (_Dapp *DappSession) ZetaValue() (*big.Int, error) {
	return _Dapp.Contract.ZetaValue(&_Dapp.CallOpts)
}

// ZetaValue is a free data retrieval call binding the contract method 0x050337a2.
//
// Solidity: function zetaValue() view returns(uint256)
func (_Dapp *DappCallerSession) ZetaValue() (*big.Int, error) {
	return _Dapp.Contract.ZetaValue(&_Dapp.CallOpts)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onZetaMessage((bytes,uint256,address,uint256,bytes) zetaMessage) returns()
func (_Dapp *DappTransactor) OnZetaMessage(opts *bind.TransactOpts, zetaMessage ZetaInterfacesZetaMessage) (*types.Transaction, error) {
	return _Dapp.contract.Transact(opts, "onZetaMessage", zetaMessage)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onZetaMessage((bytes,uint256,address,uint256,bytes) zetaMessage) returns()
func (_Dapp *DappSession) OnZetaMessage(zetaMessage ZetaInterfacesZetaMessage) (*types.Transaction, error) {
	return _Dapp.Contract.OnZetaMessage(&_Dapp.TransactOpts, zetaMessage)
}

// OnZetaMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onZetaMessage((bytes,uint256,address,uint256,bytes) zetaMessage) returns()
func (_Dapp *DappTransactorSession) OnZetaMessage(zetaMessage ZetaInterfacesZetaMessage) (*types.Transaction, error) {
	return _Dapp.Contract.OnZetaMessage(&_Dapp.TransactOpts, zetaMessage)
}

// OnZetaRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onZetaRevert((address,uint256,bytes,uint256,uint256,bytes) zetaRevert) returns()
func (_Dapp *DappTransactor) OnZetaRevert(opts *bind.TransactOpts, zetaRevert ZetaInterfacesZetaRevert) (*types.Transaction, error) {
	return _Dapp.contract.Transact(opts, "onZetaRevert", zetaRevert)
}

// OnZetaRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onZetaRevert((address,uint256,bytes,uint256,uint256,bytes) zetaRevert) returns()
func (_Dapp *DappSession) OnZetaRevert(zetaRevert ZetaInterfacesZetaRevert) (*types.Transaction, error) {
	return _Dapp.Contract.OnZetaRevert(&_Dapp.TransactOpts, zetaRevert)
}

// OnZetaRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onZetaRevert((address,uint256,bytes,uint256,uint256,bytes) zetaRevert) returns()
func (_Dapp *DappTransactorSession) OnZetaRevert(zetaRevert ZetaInterfacesZetaRevert) (*types.Transaction, error) {
	return _Dapp.Contract.OnZetaRevert(&_Dapp.TransactOpts, zetaRevert)
}
