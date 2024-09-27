// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testdappv2

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

// TestDAppV2MessageContext is an auto generated low-level Go binding around an user-defined struct.
type TestDAppV2MessageContext struct {
	Sender common.Address
}

// TestDAppV2RevertContext is an auto generated low-level Go binding around an user-defined struct.
type TestDAppV2RevertContext struct {
	Asset         common.Address
	Amount        uint64
	RevertMessage []byte
}

// TestDAppV2zContext is an auto generated low-level Go binding around an user-defined struct.
type TestDAppV2zContext struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}

// TestDAppV2MetaData contains all meta data concerning the TestDAppV2 contract.
var TestDAppV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"amountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"calledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"erc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"erc20Call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"expectedOnCallSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"gasCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getAmountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getCalledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structTestDAppV2.MessageContext\",\"name\":\"messageContext\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structTestDAppV2.zContext\",\"name\":\"_context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"}],\"internalType\":\"structTestDAppV2.RevertContext\",\"name\":\"revertContext\",\"type\":\"tuple\"}],\"name\":\"onRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"senderWithMessage\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_expectedOnCallSender\",\"type\":\"address\"}],\"name\":\"setExpectedOnCallSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"simpleCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061147c806100206000396000f3fe6080604052600436106100c65760003560e01c8063a799911f1161007f578063de43156e11610059578063de43156e14610267578063e2842ed714610290578063f592cbfb146102cd578063f936ae851461030a576100cd565b8063a799911f146101f9578063c234fecf14610215578063c7a339a91461023e576100cd565b806336e980a0146100d25780634297a263146100fb57806359f4a77714610138578063660b9de014610163578063676cc0541461018c5780639291fe26146101bc576100cd565b366100cd57005b600080fd5b3480156100de57600080fd5b506100f960048036038101906100f49190610b78565b610347565b005b34801561010757600080fd5b50610122600480360381019061011d9190610bf7565b610371565b60405161012f9190610c3d565b60405180910390f35b34801561014457600080fd5b5061014d610389565b60405161015a9190610c99565b60405180910390f35b34801561016f57600080fd5b5061018a60048036038101906101859190610cd8565b6103ad565b005b6101a660048036038101906101a19190610da0565b610468565b6040516101b39190610e88565b60405180910390f35b3480156101c857600080fd5b506101e360048036038101906101de9190610b78565b61061d565b6040516101f09190610c3d565b60405180910390f35b610213600480360381019061020e9190610b78565b610660565b005b34801561022157600080fd5b5061023c60048036038101906102379190610ed6565b610689565b005b34801561024a57600080fd5b5061026560048036038101906102609190610f6d565b6106cc565b005b34801561027357600080fd5b5061028e60048036038101906102899190610ffb565b610780565b005b34801561029c57600080fd5b506102b760048036038101906102b29190610bf7565b610879565b6040516102c491906110ba565b60405180910390f35b3480156102d957600080fd5b506102f460048036038101906102ef9190610b78565b610899565b60405161030191906110ba565b60405180910390f35b34801561031657600080fd5b50610331600480360381019061032c9190611176565b6108e9565b60405161033e9190610c99565b60405180910390f35b61035081610932565b1561035a57600080fd5b61036381610988565b61036e8160006109dc565b50565b60036020528060005260406000206000915090505481565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6104088180604001906103c091906111ce565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610988565b61046581806040019061041b91906111ce565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505060006109dc565b50565b606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168460000160208101906104b49190610ed6565b73ffffffffffffffffffffffffffffffffffffffff161461050a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105019061128e565b60405180910390fd5b61055783838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610988565b6105a583838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050346109dc565b8360000160208101906105b89190610ed6565b600284846040516105ca9291906112de565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509392505050565b60006003600083604051602001610634919061133e565b604051602081830303815290604052805190602001208152602001908152602001600020549050919050565b61066981610932565b1561067357600080fd5b61067c81610988565b61068681346109dc565b50565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6106d581610932565b156106df57600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b815260040161071c93929190611355565b6020604051808303816000875af115801561073b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061075f91906113b8565b61076857600080fd5b61077181610988565b61077b81836109dc565b505050565b6107cd82828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610932565b156107d757600080fd5b61082482828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610988565b61087282828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050846109dc565b5050505050565b60016020528060005260406000206000915054906101000a900460ff1681565b600060016000836040516020016108b0919061133e565b60405160208183030381529060405280519060200120815260200190815260200160002060009054906101000a900460ff169050919050565b6002818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600060405160200161094390611431565b604051602081830303815290604052805190602001208260405160200161096a919061133e565b60405160208183030381529060405280519060200120149050919050565b60018060008360405160200161099e919061133e565b60405160208183030381529060405280519060200120815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b8060036000846040516020016109f2919061133e565b604051602081830303815290604052805190602001208152602001908152602001600020819055505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610a8582610a3c565b810181811067ffffffffffffffff82111715610aa457610aa3610a4d565b5b80604052505050565b6000610ab7610a1e565b9050610ac38282610a7c565b919050565b600067ffffffffffffffff821115610ae357610ae2610a4d565b5b610aec82610a3c565b9050602081019050919050565b82818337600083830152505050565b6000610b1b610b1684610ac8565b610aad565b905082815260208101848484011115610b3757610b36610a37565b5b610b42848285610af9565b509392505050565b600082601f830112610b5f57610b5e610a32565b5b8135610b6f848260208601610b08565b91505092915050565b600060208284031215610b8e57610b8d610a28565b5b600082013567ffffffffffffffff811115610bac57610bab610a2d565b5b610bb884828501610b4a565b91505092915050565b6000819050919050565b610bd481610bc1565b8114610bdf57600080fd5b50565b600081359050610bf181610bcb565b92915050565b600060208284031215610c0d57610c0c610a28565b5b6000610c1b84828501610be2565b91505092915050565b6000819050919050565b610c3781610c24565b82525050565b6000602082019050610c526000830184610c2e565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610c8382610c58565b9050919050565b610c9381610c78565b82525050565b6000602082019050610cae6000830184610c8a565b92915050565b600080fd5b600060608284031215610ccf57610cce610cb4565b5b81905092915050565b600060208284031215610cee57610ced610a28565b5b600082013567ffffffffffffffff811115610d0c57610d0b610a2d565b5b610d1884828501610cb9565b91505092915050565b600060208284031215610d3757610d36610cb4565b5b81905092915050565b600080fd5b600080fd5b60008083601f840112610d6057610d5f610a32565b5b8235905067ffffffffffffffff811115610d7d57610d7c610d40565b5b602083019150836001820283011115610d9957610d98610d45565b5b9250929050565b600080600060408486031215610db957610db8610a28565b5b6000610dc786828701610d21565b935050602084013567ffffffffffffffff811115610de857610de7610a2d565b5b610df486828701610d4a565b92509250509250925092565b600081519050919050565b600082825260208201905092915050565b60005b83811015610e3a578082015181840152602081019050610e1f565b83811115610e49576000848401525b50505050565b6000610e5a82610e00565b610e648185610e0b565b9350610e74818560208601610e1c565b610e7d81610a3c565b840191505092915050565b60006020820190508181036000830152610ea28184610e4f565b905092915050565b610eb381610c78565b8114610ebe57600080fd5b50565b600081359050610ed081610eaa565b92915050565b600060208284031215610eec57610eeb610a28565b5b6000610efa84828501610ec1565b91505092915050565b6000610f0e82610c78565b9050919050565b610f1e81610f03565b8114610f2957600080fd5b50565b600081359050610f3b81610f15565b92915050565b610f4a81610c24565b8114610f5557600080fd5b50565b600081359050610f6781610f41565b92915050565b600080600060608486031215610f8657610f85610a28565b5b6000610f9486828701610f2c565b9350506020610fa586828701610f58565b925050604084013567ffffffffffffffff811115610fc657610fc5610a2d565b5b610fd286828701610b4a565b9150509250925092565b600060608284031215610ff257610ff1610cb4565b5b81905092915050565b60008060008060006080868803121561101757611016610a28565b5b600086013567ffffffffffffffff81111561103557611034610a2d565b5b61104188828901610fdc565b955050602061105288828901610ec1565b945050604061106388828901610f58565b935050606086013567ffffffffffffffff81111561108457611083610a2d565b5b61109088828901610d4a565b92509250509295509295909350565b60008115159050919050565b6110b48161109f565b82525050565b60006020820190506110cf60008301846110ab565b92915050565b600067ffffffffffffffff8211156110f0576110ef610a4d565b5b6110f982610a3c565b9050602081019050919050565b6000611119611114846110d5565b610aad565b90508281526020810184848401111561113557611134610a37565b5b611140848285610af9565b509392505050565b600082601f83011261115d5761115c610a32565b5b813561116d848260208601611106565b91505092915050565b60006020828403121561118c5761118b610a28565b5b600082013567ffffffffffffffff8111156111aa576111a9610a2d565b5b6111b684828501611148565b91505092915050565b600080fd5b600080fd5b600080fd5b600080833560016020038436030381126111eb576111ea6111bf565b5b80840192508235915067ffffffffffffffff82111561120d5761120c6111c4565b5b602083019250600182023603831315611229576112286111c9565b5b509250929050565b600082825260208201905092915050565b7f756e61757468656e746963617465642073656e64657200000000000000000000600082015250565b6000611278601683611231565b915061128382611242565b602082019050919050565b600060208201905081810360008301526112a78161126b565b9050919050565b600081905092915050565b60006112c583856112ae565b93506112d2838584610af9565b82840190509392505050565b60006112eb8284866112b9565b91508190509392505050565b600081519050919050565b600081905092915050565b6000611318826112f7565b6113228185611302565b9350611332818560208601610e1c565b80840191505092915050565b600061134a828461130d565b915081905092915050565b600060608201905061136a6000830186610c8a565b6113776020830185610c8a565b6113846040830184610c2e565b949350505050565b6113958161109f565b81146113a057600080fd5b50565b6000815190506113b28161138c565b92915050565b6000602082840312156113ce576113cd610a28565b5b60006113dc848285016113a3565b91505092915050565b7f7265766572740000000000000000000000000000000000000000000000000000600082015250565b600061141b600683611302565b9150611426826113e5565b600682019050919050565b600061143c8261140e565b915081905091905056fea264697066735822122046fb444f8c754142359339f5f1d728822414688169d0d22f23bd25c688c73fdf64736f6c634300080a0033",
}

// TestDAppV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDAppV2MetaData.ABI instead.
var TestDAppV2ABI = TestDAppV2MetaData.ABI

// TestDAppV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestDAppV2MetaData.Bin instead.
var TestDAppV2Bin = TestDAppV2MetaData.Bin

// DeployTestDAppV2 deploys a new Ethereum contract, binding an instance of TestDAppV2 to it.
func DeployTestDAppV2(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestDAppV2, error) {
	parsed, err := TestDAppV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDAppV2Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestDAppV2{TestDAppV2Caller: TestDAppV2Caller{contract: contract}, TestDAppV2Transactor: TestDAppV2Transactor{contract: contract}, TestDAppV2Filterer: TestDAppV2Filterer{contract: contract}}, nil
}

// TestDAppV2 is an auto generated Go binding around an Ethereum contract.
type TestDAppV2 struct {
	TestDAppV2Caller     // Read-only binding to the contract
	TestDAppV2Transactor // Write-only binding to the contract
	TestDAppV2Filterer   // Log filterer for contract events
}

// TestDAppV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type TestDAppV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TestDAppV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestDAppV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestDAppV2Session struct {
	Contract     *TestDAppV2       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestDAppV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestDAppV2CallerSession struct {
	Contract *TestDAppV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// TestDAppV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestDAppV2TransactorSession struct {
	Contract     *TestDAppV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// TestDAppV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type TestDAppV2Raw struct {
	Contract *TestDAppV2 // Generic contract binding to access the raw methods on
}

// TestDAppV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestDAppV2CallerRaw struct {
	Contract *TestDAppV2Caller // Generic read-only contract binding to access the raw methods on
}

// TestDAppV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestDAppV2TransactorRaw struct {
	Contract *TestDAppV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTestDAppV2 creates a new instance of TestDAppV2, bound to a specific deployed contract.
func NewTestDAppV2(address common.Address, backend bind.ContractBackend) (*TestDAppV2, error) {
	contract, err := bindTestDAppV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestDAppV2{TestDAppV2Caller: TestDAppV2Caller{contract: contract}, TestDAppV2Transactor: TestDAppV2Transactor{contract: contract}, TestDAppV2Filterer: TestDAppV2Filterer{contract: contract}}, nil
}

// NewTestDAppV2Caller creates a new read-only instance of TestDAppV2, bound to a specific deployed contract.
func NewTestDAppV2Caller(address common.Address, caller bind.ContractCaller) (*TestDAppV2Caller, error) {
	contract, err := bindTestDAppV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppV2Caller{contract: contract}, nil
}

// NewTestDAppV2Transactor creates a new write-only instance of TestDAppV2, bound to a specific deployed contract.
func NewTestDAppV2Transactor(address common.Address, transactor bind.ContractTransactor) (*TestDAppV2Transactor, error) {
	contract, err := bindTestDAppV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppV2Transactor{contract: contract}, nil
}

// NewTestDAppV2Filterer creates a new log filterer instance of TestDAppV2, bound to a specific deployed contract.
func NewTestDAppV2Filterer(address common.Address, filterer bind.ContractFilterer) (*TestDAppV2Filterer, error) {
	contract, err := bindTestDAppV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestDAppV2Filterer{contract: contract}, nil
}

// bindTestDAppV2 binds a generic wrapper to an already deployed contract.
func bindTestDAppV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestDAppV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDAppV2 *TestDAppV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDAppV2.Contract.TestDAppV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDAppV2 *TestDAppV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDAppV2.Contract.TestDAppV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDAppV2 *TestDAppV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDAppV2.Contract.TestDAppV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDAppV2 *TestDAppV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDAppV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDAppV2 *TestDAppV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDAppV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDAppV2 *TestDAppV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDAppV2.Contract.contract.Transact(opts, method, params...)
}

// AmountWithMessage is a free data retrieval call binding the contract method 0x4297a263.
//
// Solidity: function amountWithMessage(bytes32 ) view returns(uint256)
func (_TestDAppV2 *TestDAppV2Caller) AmountWithMessage(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "amountWithMessage", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AmountWithMessage is a free data retrieval call binding the contract method 0x4297a263.
//
// Solidity: function amountWithMessage(bytes32 ) view returns(uint256)
func (_TestDAppV2 *TestDAppV2Session) AmountWithMessage(arg0 [32]byte) (*big.Int, error) {
	return _TestDAppV2.Contract.AmountWithMessage(&_TestDAppV2.CallOpts, arg0)
}

// AmountWithMessage is a free data retrieval call binding the contract method 0x4297a263.
//
// Solidity: function amountWithMessage(bytes32 ) view returns(uint256)
func (_TestDAppV2 *TestDAppV2CallerSession) AmountWithMessage(arg0 [32]byte) (*big.Int, error) {
	return _TestDAppV2.Contract.AmountWithMessage(&_TestDAppV2.CallOpts, arg0)
}

// CalledWithMessage is a free data retrieval call binding the contract method 0xe2842ed7.
//
// Solidity: function calledWithMessage(bytes32 ) view returns(bool)
func (_TestDAppV2 *TestDAppV2Caller) CalledWithMessage(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "calledWithMessage", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CalledWithMessage is a free data retrieval call binding the contract method 0xe2842ed7.
//
// Solidity: function calledWithMessage(bytes32 ) view returns(bool)
func (_TestDAppV2 *TestDAppV2Session) CalledWithMessage(arg0 [32]byte) (bool, error) {
	return _TestDAppV2.Contract.CalledWithMessage(&_TestDAppV2.CallOpts, arg0)
}

// CalledWithMessage is a free data retrieval call binding the contract method 0xe2842ed7.
//
// Solidity: function calledWithMessage(bytes32 ) view returns(bool)
func (_TestDAppV2 *TestDAppV2CallerSession) CalledWithMessage(arg0 [32]byte) (bool, error) {
	return _TestDAppV2.Contract.CalledWithMessage(&_TestDAppV2.CallOpts, arg0)
}

// ExpectedOnCallSender is a free data retrieval call binding the contract method 0x59f4a777.
//
// Solidity: function expectedOnCallSender() view returns(address)
func (_TestDAppV2 *TestDAppV2Caller) ExpectedOnCallSender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "expectedOnCallSender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ExpectedOnCallSender is a free data retrieval call binding the contract method 0x59f4a777.
//
// Solidity: function expectedOnCallSender() view returns(address)
func (_TestDAppV2 *TestDAppV2Session) ExpectedOnCallSender() (common.Address, error) {
	return _TestDAppV2.Contract.ExpectedOnCallSender(&_TestDAppV2.CallOpts)
}

// ExpectedOnCallSender is a free data retrieval call binding the contract method 0x59f4a777.
//
// Solidity: function expectedOnCallSender() view returns(address)
func (_TestDAppV2 *TestDAppV2CallerSession) ExpectedOnCallSender() (common.Address, error) {
	return _TestDAppV2.Contract.ExpectedOnCallSender(&_TestDAppV2.CallOpts)
}

// GetAmountWithMessage is a free data retrieval call binding the contract method 0x9291fe26.
//
// Solidity: function getAmountWithMessage(string message) view returns(uint256)
func (_TestDAppV2 *TestDAppV2Caller) GetAmountWithMessage(opts *bind.CallOpts, message string) (*big.Int, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "getAmountWithMessage", message)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountWithMessage is a free data retrieval call binding the contract method 0x9291fe26.
//
// Solidity: function getAmountWithMessage(string message) view returns(uint256)
func (_TestDAppV2 *TestDAppV2Session) GetAmountWithMessage(message string) (*big.Int, error) {
	return _TestDAppV2.Contract.GetAmountWithMessage(&_TestDAppV2.CallOpts, message)
}

// GetAmountWithMessage is a free data retrieval call binding the contract method 0x9291fe26.
//
// Solidity: function getAmountWithMessage(string message) view returns(uint256)
func (_TestDAppV2 *TestDAppV2CallerSession) GetAmountWithMessage(message string) (*big.Int, error) {
	return _TestDAppV2.Contract.GetAmountWithMessage(&_TestDAppV2.CallOpts, message)
}

// GetCalledWithMessage is a free data retrieval call binding the contract method 0xf592cbfb.
//
// Solidity: function getCalledWithMessage(string message) view returns(bool)
func (_TestDAppV2 *TestDAppV2Caller) GetCalledWithMessage(opts *bind.CallOpts, message string) (bool, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "getCalledWithMessage", message)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetCalledWithMessage is a free data retrieval call binding the contract method 0xf592cbfb.
//
// Solidity: function getCalledWithMessage(string message) view returns(bool)
func (_TestDAppV2 *TestDAppV2Session) GetCalledWithMessage(message string) (bool, error) {
	return _TestDAppV2.Contract.GetCalledWithMessage(&_TestDAppV2.CallOpts, message)
}

// GetCalledWithMessage is a free data retrieval call binding the contract method 0xf592cbfb.
//
// Solidity: function getCalledWithMessage(string message) view returns(bool)
func (_TestDAppV2 *TestDAppV2CallerSession) GetCalledWithMessage(message string) (bool, error) {
	return _TestDAppV2.Contract.GetCalledWithMessage(&_TestDAppV2.CallOpts, message)
}

// SenderWithMessage is a free data retrieval call binding the contract method 0xf936ae85.
//
// Solidity: function senderWithMessage(bytes ) view returns(address)
func (_TestDAppV2 *TestDAppV2Caller) SenderWithMessage(opts *bind.CallOpts, arg0 []byte) (common.Address, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "senderWithMessage", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SenderWithMessage is a free data retrieval call binding the contract method 0xf936ae85.
//
// Solidity: function senderWithMessage(bytes ) view returns(address)
func (_TestDAppV2 *TestDAppV2Session) SenderWithMessage(arg0 []byte) (common.Address, error) {
	return _TestDAppV2.Contract.SenderWithMessage(&_TestDAppV2.CallOpts, arg0)
}

// SenderWithMessage is a free data retrieval call binding the contract method 0xf936ae85.
//
// Solidity: function senderWithMessage(bytes ) view returns(address)
func (_TestDAppV2 *TestDAppV2CallerSession) SenderWithMessage(arg0 []byte) (common.Address, error) {
	return _TestDAppV2.Contract.SenderWithMessage(&_TestDAppV2.CallOpts, arg0)
}

// Erc20Call is a paid mutator transaction binding the contract method 0xc7a339a9.
//
// Solidity: function erc20Call(address erc20, uint256 amount, string message) returns()
func (_TestDAppV2 *TestDAppV2Transactor) Erc20Call(opts *bind.TransactOpts, erc20 common.Address, amount *big.Int, message string) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "erc20Call", erc20, amount, message)
}

// Erc20Call is a paid mutator transaction binding the contract method 0xc7a339a9.
//
// Solidity: function erc20Call(address erc20, uint256 amount, string message) returns()
func (_TestDAppV2 *TestDAppV2Session) Erc20Call(erc20 common.Address, amount *big.Int, message string) (*types.Transaction, error) {
	return _TestDAppV2.Contract.Erc20Call(&_TestDAppV2.TransactOpts, erc20, amount, message)
}

// Erc20Call is a paid mutator transaction binding the contract method 0xc7a339a9.
//
// Solidity: function erc20Call(address erc20, uint256 amount, string message) returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) Erc20Call(erc20 common.Address, amount *big.Int, message string) (*types.Transaction, error) {
	return _TestDAppV2.Contract.Erc20Call(&_TestDAppV2.TransactOpts, erc20, amount, message)
}

// GasCall is a paid mutator transaction binding the contract method 0xa799911f.
//
// Solidity: function gasCall(string message) payable returns()
func (_TestDAppV2 *TestDAppV2Transactor) GasCall(opts *bind.TransactOpts, message string) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "gasCall", message)
}

// GasCall is a paid mutator transaction binding the contract method 0xa799911f.
//
// Solidity: function gasCall(string message) payable returns()
func (_TestDAppV2 *TestDAppV2Session) GasCall(message string) (*types.Transaction, error) {
	return _TestDAppV2.Contract.GasCall(&_TestDAppV2.TransactOpts, message)
}

// GasCall is a paid mutator transaction binding the contract method 0xa799911f.
//
// Solidity: function gasCall(string message) payable returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) GasCall(message string) (*types.Transaction, error) {
	return _TestDAppV2.Contract.GasCall(&_TestDAppV2.TransactOpts, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppV2 *TestDAppV2Transactor) OnCall(opts *bind.TransactOpts, messageContext TestDAppV2MessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "onCall", messageContext, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppV2 *TestDAppV2Session) OnCall(messageContext TestDAppV2MessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCall(&_TestDAppV2.TransactOpts, messageContext, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppV2 *TestDAppV2TransactorSession) OnCall(messageContext TestDAppV2MessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCall(&_TestDAppV2.TransactOpts, messageContext, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) _context, address _zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2Transactor) OnCrossChainCall(opts *bind.TransactOpts, _context TestDAppV2zContext, _zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "onCrossChainCall", _context, _zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) _context, address _zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2Session) OnCrossChainCall(_context TestDAppV2zContext, _zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCrossChainCall(&_TestDAppV2.TransactOpts, _context, _zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) _context, address _zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) OnCrossChainCall(_context TestDAppV2zContext, _zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCrossChainCall(&_TestDAppV2.TransactOpts, _context, _zrc20, amount, message)
}

// OnRevert is a paid mutator transaction binding the contract method 0x660b9de0.
//
// Solidity: function onRevert((address,uint64,bytes) revertContext) returns()
func (_TestDAppV2 *TestDAppV2Transactor) OnRevert(opts *bind.TransactOpts, revertContext TestDAppV2RevertContext) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "onRevert", revertContext)
}

// OnRevert is a paid mutator transaction binding the contract method 0x660b9de0.
//
// Solidity: function onRevert((address,uint64,bytes) revertContext) returns()
func (_TestDAppV2 *TestDAppV2Session) OnRevert(revertContext TestDAppV2RevertContext) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnRevert(&_TestDAppV2.TransactOpts, revertContext)
}

// OnRevert is a paid mutator transaction binding the contract method 0x660b9de0.
//
// Solidity: function onRevert((address,uint64,bytes) revertContext) returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) OnRevert(revertContext TestDAppV2RevertContext) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnRevert(&_TestDAppV2.TransactOpts, revertContext)
}

// SetExpectedOnCallSender is a paid mutator transaction binding the contract method 0xc234fecf.
//
// Solidity: function setExpectedOnCallSender(address _expectedOnCallSender) returns()
func (_TestDAppV2 *TestDAppV2Transactor) SetExpectedOnCallSender(opts *bind.TransactOpts, _expectedOnCallSender common.Address) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "setExpectedOnCallSender", _expectedOnCallSender)
}

// SetExpectedOnCallSender is a paid mutator transaction binding the contract method 0xc234fecf.
//
// Solidity: function setExpectedOnCallSender(address _expectedOnCallSender) returns()
func (_TestDAppV2 *TestDAppV2Session) SetExpectedOnCallSender(_expectedOnCallSender common.Address) (*types.Transaction, error) {
	return _TestDAppV2.Contract.SetExpectedOnCallSender(&_TestDAppV2.TransactOpts, _expectedOnCallSender)
}

// SetExpectedOnCallSender is a paid mutator transaction binding the contract method 0xc234fecf.
//
// Solidity: function setExpectedOnCallSender(address _expectedOnCallSender) returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) SetExpectedOnCallSender(_expectedOnCallSender common.Address) (*types.Transaction, error) {
	return _TestDAppV2.Contract.SetExpectedOnCallSender(&_TestDAppV2.TransactOpts, _expectedOnCallSender)
}

// SimpleCall is a paid mutator transaction binding the contract method 0x36e980a0.
//
// Solidity: function simpleCall(string message) returns()
func (_TestDAppV2 *TestDAppV2Transactor) SimpleCall(opts *bind.TransactOpts, message string) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "simpleCall", message)
}

// SimpleCall is a paid mutator transaction binding the contract method 0x36e980a0.
//
// Solidity: function simpleCall(string message) returns()
func (_TestDAppV2 *TestDAppV2Session) SimpleCall(message string) (*types.Transaction, error) {
	return _TestDAppV2.Contract.SimpleCall(&_TestDAppV2.TransactOpts, message)
}

// SimpleCall is a paid mutator transaction binding the contract method 0x36e980a0.
//
// Solidity: function simpleCall(string message) returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) SimpleCall(message string) (*types.Transaction, error) {
	return _TestDAppV2.Contract.SimpleCall(&_TestDAppV2.TransactOpts, message)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDAppV2 *TestDAppV2Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDAppV2.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDAppV2 *TestDAppV2Session) Receive() (*types.Transaction, error) {
	return _TestDAppV2.Contract.Receive(&_TestDAppV2.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) Receive() (*types.Transaction, error) {
	return _TestDAppV2.Contract.Receive(&_TestDAppV2.TransactOpts)
}
