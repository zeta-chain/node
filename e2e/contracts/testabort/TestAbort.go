// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testabort

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

// AbortContext is an auto generated low-level Go binding around an user-defined struct.
type AbortContext struct {
	Sender        []byte
	Asset         common.Address
	Amount        *big.Int
	Outgoing      bool
	ChainID       *big.Int
	RevertMessage []byte
}

// TestAbortMetaData contains all meta data concerning the TestAbort contract.
var TestAbortMetaData = &bind.MetaData{
	ABI: "[{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"aborted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"abortedWithMessage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"outgoing\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getAbortedWithMessage\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"outgoing\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"}],\"internalType\":\"structAbortContext\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isAborted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"outgoing\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"}],\"internalType\":\"structAbortContext\",\"name\":\"abortContext\",\"type\":\"tuple\"}],\"name\":\"onAbort\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506111c28061001f6000396000f3fe60806040526004361061004e5760003560e01c80632d4cfb7e1461005757806372748f7d1461008057806380b62b70146100c25780639e59f463146100ed578063fe4caa641461012a57610055565b3661005557005b005b34801561006357600080fd5b5061007e600480360381019061007991906106e2565b610155565b005b34801561008c57600080fd5b506100a760048036038101906100a29190610761565b6101bd565b6040516100b996959493929190610893565b60405180910390f35b3480156100ce57600080fd5b506100d7610336565b6040516100e49190610902565b60405180910390f35b3480156100f957600080fd5b50610114600480360381019061010f9190610a52565b610349565b6040516101219190610ba2565b60405180910390f35b34801561013657600080fd5b5061013f610544565b60405161014c9190610902565b60405180910390f35b6101ba818060a001906101689190610bd3565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050826101b590610e3d565b61055b565b50565b60006020528060005260406000206000915090508060000180546101e090610e7f565b80601f016020809104026020016040519081016040528092919081815260200182805461020c90610e7f565b80156102595780601f1061022e57610100808354040283529160200191610259565b820191906000526020600020905b81548152906001019060200180831161023c57829003601f168201915b5050505050908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030160009054906101000a900460ff16908060040154908060050180546102b390610e7f565b80601f01602080910402602001604051908101604052809291908181526020018280546102df90610e7f565b801561032c5780601f106103015761010080835404028352916020019161032c565b820191906000526020600020905b81548152906001019060200180831161030f57829003601f168201915b5050505050905086565b600160009054906101000a900460ff1681565b61035161065c565b600080836040516020016103659190610ef7565b6040516020818303038152906040528051906020012081526020019081526020016000206040518060c00160405290816000820180546103a490610e7f565b80601f01602080910402602001604051908101604052809291908181526020018280546103d090610e7f565b801561041d5780601f106103f25761010080835404028352916020019161041d565b820191906000526020600020905b81548152906001019060200180831161040057829003601f168201915b505050505081526020016001820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600282015481526020016003820160009054906101000a900460ff16151515158152602001600482015481526020016005820180546104bb90610e7f565b80601f01602080910402602001604051908101604052809291908181526020018280546104e790610e7f565b80156105345780601f1061050957610100808354040283529160200191610534565b820191906000526020600020905b81548152906001019060200180831161051757829003601f168201915b5050505050815250509050919050565b6000600160009054906101000a900460ff16905090565b80600080846040516020016105709190610ef7565b60405160208183030381529060405280519060200120815260200190815260200160002060008201518160000190816105a991906110ba565b5060208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030160006101000a81548160ff0219169083151502179055506080820151816004015560a082015181600501908161063a91906110ba565b5090505060018060006101000a81548160ff0219169083151502179055505050565b6040518060c0016040528060608152602001600073ffffffffffffffffffffffffffffffffffffffff1681526020016000815260200160001515815260200160008152602001606081525090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600060c082840312156106d9576106d86106be565b5b81905092915050565b6000602082840312156106f8576106f76106b4565b5b600082013567ffffffffffffffff811115610716576107156106b9565b5b610722848285016106c3565b91505092915050565b6000819050919050565b61073e8161072b565b811461074957600080fd5b50565b60008135905061075b81610735565b92915050565b600060208284031215610777576107766106b4565b5b60006107858482850161074c565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b838110156107c85780820151818401526020810190506107ad565b60008484015250505050565b6000601f19601f8301169050919050565b60006107f08261078e565b6107fa8185610799565b935061080a8185602086016107aa565b610813816107d4565b840191505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006108498261081e565b9050919050565b6108598161083e565b82525050565b6000819050919050565b6108728161085f565b82525050565b60008115159050919050565b61088d81610878565b82525050565b600060c08201905081810360008301526108ad81896107e5565b90506108bc6020830188610850565b6108c96040830187610869565b6108d66060830186610884565b6108e36080830185610869565b81810360a08301526108f581846107e5565b9050979650505050505050565b60006020820190506109176000830184610884565b92915050565b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61095f826107d4565b810181811067ffffffffffffffff8211171561097e5761097d610927565b5b80604052505050565b60006109916106aa565b905061099d8282610956565b919050565b600067ffffffffffffffff8211156109bd576109bc610927565b5b6109c6826107d4565b9050602081019050919050565b82818337600083830152505050565b60006109f56109f0846109a2565b610987565b905082815260208101848484011115610a1157610a10610922565b5b610a1c8482856109d3565b509392505050565b600082601f830112610a3957610a3861091d565b5b8135610a498482602086016109e2565b91505092915050565b600060208284031215610a6857610a676106b4565b5b600082013567ffffffffffffffff811115610a8657610a856106b9565b5b610a9284828501610a24565b91505092915050565b600082825260208201905092915050565b6000610ab78261078e565b610ac18185610a9b565b9350610ad18185602086016107aa565b610ada816107d4565b840191505092915050565b610aee8161083e565b82525050565b610afd8161085f565b82525050565b610b0c81610878565b82525050565b600060c0830160008301518482036000860152610b2f8282610aac565b9150506020830151610b446020860182610ae5565b506040830151610b576040860182610af4565b506060830151610b6a6060860182610b03565b506080830151610b7d6080860182610af4565b5060a083015184820360a0860152610b958282610aac565b9150508091505092915050565b60006020820190508181036000830152610bbc8184610b12565b905092915050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112610bf057610bef610bc4565b5b80840192508235915067ffffffffffffffff821115610c1257610c11610bc9565b5b602083019250600182023603831315610c2e57610c2d610bce565b5b509250929050565b600080fd5b600080fd5b600067ffffffffffffffff821115610c5b57610c5a610927565b5b610c64826107d4565b9050602081019050919050565b6000610c84610c7f84610c40565b610987565b905082815260208101848484011115610ca057610c9f610922565b5b610cab8482856109d3565b509392505050565b600082601f830112610cc857610cc761091d565b5b8135610cd8848260208601610c71565b91505092915050565b610cea8161083e565b8114610cf557600080fd5b50565b600081359050610d0781610ce1565b92915050565b610d168161085f565b8114610d2157600080fd5b50565b600081359050610d3381610d0d565b92915050565b610d4281610878565b8114610d4d57600080fd5b50565b600081359050610d5f81610d39565b92915050565b600060c08284031215610d7b57610d7a610c36565b5b610d8560c0610987565b9050600082013567ffffffffffffffff811115610da557610da4610c3b565b5b610db184828501610cb3565b6000830152506020610dc584828501610cf8565b6020830152506040610dd984828501610d24565b6040830152506060610ded84828501610d50565b6060830152506080610e0184828501610d24565b60808301525060a082013567ffffffffffffffff811115610e2557610e24610c3b565b5b610e3184828501610cb3565b60a08301525092915050565b6000610e493683610d65565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610e9757607f821691505b602082108103610eaa57610ea9610e50565b5b50919050565b600081519050919050565b600081905092915050565b6000610ed182610eb0565b610edb8185610ebb565b9350610eeb8185602086016107aa565b80840191505092915050565b6000610f038284610ec6565b915081905092915050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302610f707fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610f33565b610f7a8683610f33565b95508019841693508086168417925050509392505050565b6000819050919050565b6000610fb7610fb2610fad8461085f565b610f92565b61085f565b9050919050565b6000819050919050565b610fd183610f9c565b610fe5610fdd82610fbe565b848454610f40565b825550505050565b600090565b610ffa610fed565b611005818484610fc8565b505050565b5b818110156110295761101e600082610ff2565b60018101905061100b565b5050565b601f82111561106e5761103f81610f0e565b61104884610f23565b81016020851015611057578190505b61106b61106385610f23565b83018261100a565b50505b505050565b600082821c905092915050565b600061109160001984600802611073565b1980831691505092915050565b60006110aa8383611080565b9150826002028217905092915050565b6110c38261078e565b67ffffffffffffffff8111156110dc576110db610927565b5b6110e68254610e7f565b6110f182828561102d565b600060209050601f8311600181146111245760008415611112578287015190505b61111c858261109e565b865550611184565b601f19841661113286610f0e565b60005b8281101561115a57848901518255600182019150602085019450602081019050611135565b868310156111775784890151611173601f891682611080565b8355505b6001600288020188555050505b50505050505056fea264697066735822122073b672b76653cfb3aa192e6bbbea81b40276e37f9ae4e034167d2012fb55ebaa64736f6c634300081a0033",
}

// TestAbortABI is the input ABI used to generate the binding from.
// Deprecated: Use TestAbortMetaData.ABI instead.
var TestAbortABI = TestAbortMetaData.ABI

// TestAbortBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestAbortMetaData.Bin instead.
var TestAbortBin = TestAbortMetaData.Bin

// DeployTestAbort deploys a new Ethereum contract, binding an instance of TestAbort to it.
func DeployTestAbort(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestAbort, error) {
	parsed, err := TestAbortMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestAbortBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestAbort{TestAbortCaller: TestAbortCaller{contract: contract}, TestAbortTransactor: TestAbortTransactor{contract: contract}, TestAbortFilterer: TestAbortFilterer{contract: contract}}, nil
}

// TestAbort is an auto generated Go binding around an Ethereum contract.
type TestAbort struct {
	TestAbortCaller     // Read-only binding to the contract
	TestAbortTransactor // Write-only binding to the contract
	TestAbortFilterer   // Log filterer for contract events
}

// TestAbortCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestAbortCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestAbortTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestAbortTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestAbortFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestAbortFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestAbortSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestAbortSession struct {
	Contract     *TestAbort        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestAbortCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestAbortCallerSession struct {
	Contract *TestAbortCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TestAbortTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestAbortTransactorSession struct {
	Contract     *TestAbortTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TestAbortRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestAbortRaw struct {
	Contract *TestAbort // Generic contract binding to access the raw methods on
}

// TestAbortCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestAbortCallerRaw struct {
	Contract *TestAbortCaller // Generic read-only contract binding to access the raw methods on
}

// TestAbortTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestAbortTransactorRaw struct {
	Contract *TestAbortTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestAbort creates a new instance of TestAbort, bound to a specific deployed contract.
func NewTestAbort(address common.Address, backend bind.ContractBackend) (*TestAbort, error) {
	contract, err := bindTestAbort(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestAbort{TestAbortCaller: TestAbortCaller{contract: contract}, TestAbortTransactor: TestAbortTransactor{contract: contract}, TestAbortFilterer: TestAbortFilterer{contract: contract}}, nil
}

// NewTestAbortCaller creates a new read-only instance of TestAbort, bound to a specific deployed contract.
func NewTestAbortCaller(address common.Address, caller bind.ContractCaller) (*TestAbortCaller, error) {
	contract, err := bindTestAbort(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestAbortCaller{contract: contract}, nil
}

// NewTestAbortTransactor creates a new write-only instance of TestAbort, bound to a specific deployed contract.
func NewTestAbortTransactor(address common.Address, transactor bind.ContractTransactor) (*TestAbortTransactor, error) {
	contract, err := bindTestAbort(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestAbortTransactor{contract: contract}, nil
}

// NewTestAbortFilterer creates a new log filterer instance of TestAbort, bound to a specific deployed contract.
func NewTestAbortFilterer(address common.Address, filterer bind.ContractFilterer) (*TestAbortFilterer, error) {
	contract, err := bindTestAbort(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestAbortFilterer{contract: contract}, nil
}

// bindTestAbort binds a generic wrapper to an already deployed contract.
func bindTestAbort(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestAbortMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestAbort *TestAbortRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestAbort.Contract.TestAbortCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestAbort *TestAbortRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestAbort.Contract.TestAbortTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestAbort *TestAbortRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestAbort.Contract.TestAbortTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestAbort *TestAbortCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestAbort.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestAbort *TestAbortTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestAbort.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestAbort *TestAbortTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestAbort.Contract.contract.Transact(opts, method, params...)
}

// Aborted is a free data retrieval call binding the contract method 0x80b62b70.
//
// Solidity: function aborted() view returns(bool)
func (_TestAbort *TestAbortCaller) Aborted(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TestAbort.contract.Call(opts, &out, "aborted")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Aborted is a free data retrieval call binding the contract method 0x80b62b70.
//
// Solidity: function aborted() view returns(bool)
func (_TestAbort *TestAbortSession) Aborted() (bool, error) {
	return _TestAbort.Contract.Aborted(&_TestAbort.CallOpts)
}

// Aborted is a free data retrieval call binding the contract method 0x80b62b70.
//
// Solidity: function aborted() view returns(bool)
func (_TestAbort *TestAbortCallerSession) Aborted() (bool, error) {
	return _TestAbort.Contract.Aborted(&_TestAbort.CallOpts)
}

// AbortedWithMessage is a free data retrieval call binding the contract method 0x72748f7d.
//
// Solidity: function abortedWithMessage(bytes32 ) view returns(bytes sender, address asset, uint256 amount, bool outgoing, uint256 chainID, bytes revertMessage)
func (_TestAbort *TestAbortCaller) AbortedWithMessage(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Sender        []byte
	Asset         common.Address
	Amount        *big.Int
	Outgoing      bool
	ChainID       *big.Int
	RevertMessage []byte
}, error) {
	var out []interface{}
	err := _TestAbort.contract.Call(opts, &out, "abortedWithMessage", arg0)

	outstruct := new(struct {
		Sender        []byte
		Asset         common.Address
		Amount        *big.Int
		Outgoing      bool
		ChainID       *big.Int
		RevertMessage []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Sender = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Asset = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Outgoing = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.ChainID = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.RevertMessage = *abi.ConvertType(out[5], new([]byte)).(*[]byte)

	return *outstruct, err

}

// AbortedWithMessage is a free data retrieval call binding the contract method 0x72748f7d.
//
// Solidity: function abortedWithMessage(bytes32 ) view returns(bytes sender, address asset, uint256 amount, bool outgoing, uint256 chainID, bytes revertMessage)
func (_TestAbort *TestAbortSession) AbortedWithMessage(arg0 [32]byte) (struct {
	Sender        []byte
	Asset         common.Address
	Amount        *big.Int
	Outgoing      bool
	ChainID       *big.Int
	RevertMessage []byte
}, error) {
	return _TestAbort.Contract.AbortedWithMessage(&_TestAbort.CallOpts, arg0)
}

// AbortedWithMessage is a free data retrieval call binding the contract method 0x72748f7d.
//
// Solidity: function abortedWithMessage(bytes32 ) view returns(bytes sender, address asset, uint256 amount, bool outgoing, uint256 chainID, bytes revertMessage)
func (_TestAbort *TestAbortCallerSession) AbortedWithMessage(arg0 [32]byte) (struct {
	Sender        []byte
	Asset         common.Address
	Amount        *big.Int
	Outgoing      bool
	ChainID       *big.Int
	RevertMessage []byte
}, error) {
	return _TestAbort.Contract.AbortedWithMessage(&_TestAbort.CallOpts, arg0)
}

// GetAbortedWithMessage is a free data retrieval call binding the contract method 0x9e59f463.
//
// Solidity: function getAbortedWithMessage(string message) view returns((bytes,address,uint256,bool,uint256,bytes))
func (_TestAbort *TestAbortCaller) GetAbortedWithMessage(opts *bind.CallOpts, message string) (AbortContext, error) {
	var out []interface{}
	err := _TestAbort.contract.Call(opts, &out, "getAbortedWithMessage", message)

	if err != nil {
		return *new(AbortContext), err
	}

	out0 := *abi.ConvertType(out[0], new(AbortContext)).(*AbortContext)

	return out0, err

}

// GetAbortedWithMessage is a free data retrieval call binding the contract method 0x9e59f463.
//
// Solidity: function getAbortedWithMessage(string message) view returns((bytes,address,uint256,bool,uint256,bytes))
func (_TestAbort *TestAbortSession) GetAbortedWithMessage(message string) (AbortContext, error) {
	return _TestAbort.Contract.GetAbortedWithMessage(&_TestAbort.CallOpts, message)
}

// GetAbortedWithMessage is a free data retrieval call binding the contract method 0x9e59f463.
//
// Solidity: function getAbortedWithMessage(string message) view returns((bytes,address,uint256,bool,uint256,bytes))
func (_TestAbort *TestAbortCallerSession) GetAbortedWithMessage(message string) (AbortContext, error) {
	return _TestAbort.Contract.GetAbortedWithMessage(&_TestAbort.CallOpts, message)
}

// IsAborted is a free data retrieval call binding the contract method 0xfe4caa64.
//
// Solidity: function isAborted() view returns(bool)
func (_TestAbort *TestAbortCaller) IsAborted(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TestAbort.contract.Call(opts, &out, "isAborted")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAborted is a free data retrieval call binding the contract method 0xfe4caa64.
//
// Solidity: function isAborted() view returns(bool)
func (_TestAbort *TestAbortSession) IsAborted() (bool, error) {
	return _TestAbort.Contract.IsAborted(&_TestAbort.CallOpts)
}

// IsAborted is a free data retrieval call binding the contract method 0xfe4caa64.
//
// Solidity: function isAborted() view returns(bool)
func (_TestAbort *TestAbortCallerSession) IsAborted() (bool, error) {
	return _TestAbort.Contract.IsAborted(&_TestAbort.CallOpts)
}

// OnAbort is a paid mutator transaction binding the contract method 0x2d4cfb7e.
//
// Solidity: function onAbort((bytes,address,uint256,bool,uint256,bytes) abortContext) returns()
func (_TestAbort *TestAbortTransactor) OnAbort(opts *bind.TransactOpts, abortContext AbortContext) (*types.Transaction, error) {
	return _TestAbort.contract.Transact(opts, "onAbort", abortContext)
}

// OnAbort is a paid mutator transaction binding the contract method 0x2d4cfb7e.
//
// Solidity: function onAbort((bytes,address,uint256,bool,uint256,bytes) abortContext) returns()
func (_TestAbort *TestAbortSession) OnAbort(abortContext AbortContext) (*types.Transaction, error) {
	return _TestAbort.Contract.OnAbort(&_TestAbort.TransactOpts, abortContext)
}

// OnAbort is a paid mutator transaction binding the contract method 0x2d4cfb7e.
//
// Solidity: function onAbort((bytes,address,uint256,bool,uint256,bytes) abortContext) returns()
func (_TestAbort *TestAbortTransactorSession) OnAbort(abortContext AbortContext) (*types.Transaction, error) {
	return _TestAbort.Contract.OnAbort(&_TestAbort.TransactOpts, abortContext)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestAbort *TestAbortTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TestAbort.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestAbort *TestAbortSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TestAbort.Contract.Fallback(&_TestAbort.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestAbort *TestAbortTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TestAbort.Contract.Fallback(&_TestAbort.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestAbort *TestAbortTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestAbort.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestAbort *TestAbortSession) Receive() (*types.Transaction, error) {
	return _TestAbort.Contract.Receive(&_TestAbort.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestAbort *TestAbortTransactorSession) Receive() (*types.Transaction, error) {
	return _TestAbort.Contract.Receive(&_TestAbort.TransactOpts)
}
