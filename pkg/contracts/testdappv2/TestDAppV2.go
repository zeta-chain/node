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

// TestDAppV2zContext is an auto generated low-level Go binding around an user-defined struct.
type TestDAppV2zContext struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}

// TestDAppV2MetaData contains all meta data concerning the TestDAppV2 contract.
var TestDAppV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"erc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"erc20Call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"gasCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastContext\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMessage\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastZRC20\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structTestDAppV2.zContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"simpleCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611277806100206000396000f3fe60806040526004361061007b5760003560e01c8063b2f79b031161004e578063b2f79b031461011b578063b73f7eb114610146578063c7a339a914610173578063de43156e1461019c5761007b565b8063329707101461008057806336e980a0146100ab578063829a86d9146100d4578063a799911f146100ff575b600080fd5b34801561008c57600080fd5b506100956101c5565b6040516100a29190610aba565b60405180910390f35b3480156100b757600080fd5b506100d260048036038101906100cd9190610875565b610253565b005b3480156100e057600080fd5b506100e9610288565b6040516100f69190610adc565b60405180910390f35b61011960048036038101906101149190610875565b61028e565b005b34801561012757600080fd5b506101306102c2565b60405161013d9190610a2a565b60405180910390f35b34801561015257600080fd5b5061015b6102e8565b60405161016a93929190610a7c565b60405180910390f35b34801561017f57600080fd5b5061019a60048036038101906101959190610806565b6103a8565b005b3480156101a857600080fd5b506101c360048036038101906101be91906108be565b610476565b005b600580546101d290610ea7565b80601f01602080910402602001604051908101604052809291908181526020018280546101fe90610ea7565b801561024b5780601f106102205761010080835404028352916020019161024b565b820191906000526020600020905b81548152906001019060200180831161022e57829003601f168201915b505050505081565b61025c81610540565b1561026657600080fd5b806005908051906020019061027c929190610577565b50600060048190555050565b60045481565b61029781610540565b156102a157600080fd5b80600590805190602001906102b7929190610577565b503460048190555050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008060000180546102f990610ea7565b80601f016020809104026020016040519081016040528092919081815260200182805461032590610ea7565b80156103725780601f1061034757610100808354040283529160200191610372565b820191906000526020600020905b81548152906001019060200180831161035557829003601f168201915b5050505050908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154905083565b6103b181610540565b156103bb57600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b81526004016103f893929190610a45565b602060405180830381600087803b15801561041257600080fd5b505af1158015610426573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061044a91906107d9565b61045357600080fd5b8060059080519060200190610469929190610577565b5081600481905550505050565b6104c382828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610540565b156104cd57600080fd5b84600081816104dc919061118a565b90505083600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550826004819055508181600591906105389291906105fd565b505050505050565b600060405160200161055190610a15565b604051602081830303815290604052805190602001208280519060200120149050919050565b82805461058390610ea7565b90600052602060002090601f0160209004810192826105a557600085556105ec565b82601f106105be57805160ff19168380011785556105ec565b828001600101855582156105ec579182015b828111156105eb5782518255916020019190600101906105d0565b5b5090506105f99190610683565b5090565b82805461060990610ea7565b90600052602060002090601f01602090048101928261062b5760008555610672565b82601f1061064457803560ff1916838001178555610672565b82800160010185558215610672579182015b82811115610671578235825591602001919060010190610656565b5b50905061067f9190610683565b5090565b5b8082111561069c576000816000905550600101610684565b5090565b60006106b36106ae84610b7f565b610b5a565b9050828152602081018484840111156106cf576106ce611005565b5b6106da848285610dee565b509392505050565b6000813590506106f1816111e0565b92915050565b600081519050610706816111f7565b92915050565b60008083601f84011261072257610721610fe7565b5b8235905067ffffffffffffffff81111561073f5761073e610fe2565b5b60208301915083600182028301111561075b5761075a610ffb565b5b9250929050565b6000813590506107718161120e565b92915050565b600082601f83011261078c5761078b610fe7565b5b813561079c8482602086016106a0565b91505092915050565b6000606082840312156107bb576107ba610ff1565b5b81905092915050565b6000813590506107d381611225565b92915050565b6000602082840312156107ef576107ee61100f565b5b60006107fd848285016106f7565b91505092915050565b60008060006060848603121561081f5761081e61100f565b5b600061082d86828701610762565b935050602061083e868287016107c4565b925050604084013567ffffffffffffffff81111561085f5761085e61100a565b5b61086b86828701610777565b9150509250925092565b60006020828403121561088b5761088a61100f565b5b600082013567ffffffffffffffff8111156108a9576108a861100a565b5b6108b584828501610777565b91505092915050565b6000806000806000608086880312156108da576108d961100f565b5b600086013567ffffffffffffffff8111156108f8576108f761100a565b5b610904888289016107a5565b9550506020610915888289016106e2565b9450506040610926888289016107c4565b935050606086013567ffffffffffffffff8111156109475761094661100a565b5b6109538882890161070c565b92509250509295509295909350565b61096b81610c59565b82525050565b600061097c82610bd0565b6109868185610be6565b9350610996818560208601610dfd565b61099f81611014565b840191505092915050565b60006109b582610bdb565b6109bf8185610bf7565b93506109cf818560208601610dfd565b6109d881611014565b840191505092915050565b60006109f0600683610c08565b91506109fb82611064565b600682019050919050565b610a0f81610ca9565b82525050565b6000610a20826109e3565b9150819050919050565b6000602082019050610a3f6000830184610962565b92915050565b6000606082019050610a5a6000830186610962565b610a676020830185610962565b610a746040830184610a06565b949350505050565b60006060820190508181036000830152610a968186610971565b9050610aa56020830185610962565b610ab26040830184610a06565b949350505050565b60006020820190508181036000830152610ad481846109aa565b905092915050565b6000602082019050610af16000830184610a06565b92915050565b60008083356001602003843603038112610b1457610b13610ff6565b5b80840192508235915067ffffffffffffffff821115610b3657610b35610fec565b5b602083019250600182023603831315610b5257610b51611000565b5b509250929050565b6000610b64610b75565b9050610b708282610ef5565b919050565b6000604051905090565b600067ffffffffffffffff821115610b9a57610b99610f73565b5b610ba382611014565b9050602081019050919050565b60008190508160005260206000209050919050565b600082905092915050565b600081519050919050565b600081519050919050565b600082825260208201905092915050565b600082825260208201905092915050565b600081905092915050565b601f821115610c5457610c2581610bb0565b610c2e84610e97565b81016020851015610c3d578190505b610c51610c4985610e97565b830182610cb3565b50505b505050565b6000610c6482610c89565b9050919050565b60008115159050919050565b6000610c8282610c59565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b5b81811015610cd257610cc760008261104c565b600181019050610cb4565b5050565b6000610ce182610ce8565b9050919050565b6000610cf382610cfa565b9050919050565b6000610d0582610c89565b9050919050565b6000610d1782610ca9565b9050919050565b610d288383610bc5565b67ffffffffffffffff811115610d4157610d40610f73565b5b610d4b8254610ea7565b610d56828285610c13565b6000601f831160018114610d855760008415610d73578287013590505b610d7d8582610ed9565b865550610de5565b601f198416610d9386610bb0565b60005b82811015610dbb57848901358255600182019150602085019450602081019050610d96565b86831015610dd85784890135610dd4601f891682610f26565b8355505b6001600288020188555050505b50505050505050565b82818337600083830152505050565b60005b83811015610e1b578082015181840152602081019050610e00565b83811115610e2a576000848401525b50505050565b6000810160008301610e428185610af7565b610e4d81838661117a565b50505050600181016020830180610e6381610fb6565b9050610e6f8184611157565b505050600281016040830180610e8481610fcc565b9050610e908184611198565b5050505050565b60006020601f8301049050919050565b60006002820490506001821680610ebf57607f821691505b60208210811415610ed357610ed2610f44565b5b50919050565b6000610ee58383610f26565b9150826002028217905092915050565b610efe82611014565b810181811067ffffffffffffffff82111715610f1d57610f1c610f73565b5b80604052505050565b6000610f376000198460080261103f565b1980831691505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000819050919050565b6000819050919050565b60008135610fc3816111e0565b80915050919050565b60008135610fd981611225565b80915050919050565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b60008160001b9050919050565b600082821b905092915050565b600082821c905092915050565b61105461123c565b61105f8184846111bb565b505050565b7f7265766572740000000000000000000000000000000000000000000000000000600082015250565b600073ffffffffffffffffffffffffffffffffffffffff6110ad84611025565b9350801983169250808416831791505092915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6110ef84611025565b9350801983169250808416831791505092915050565b6000600883026111357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82611032565b61113f8683611032565b95508019841693508086168417925050509392505050565b61116082610cd6565b61117361116c82610fa2565b835461108d565b8255505050565b611185838383610d1e565b505050565b6111948282610e30565b5050565b6111a182610d0c565b6111b46111ad82610fac565b83546110c3565b8255505050565b6111c483610d0c565b6111d86111d082610fac565b848454611105565b825550505050565b6111e981610c59565b81146111f457600080fd5b50565b61120081610c6b565b811461120b57600080fd5b50565b61121781610c77565b811461122257600080fd5b50565b61122e81610ca9565b811461123957600080fd5b50565b60009056fea2646970667358221220b1ce7c9b2fc96e1650beda0d57f966fa383464e64542b3eb5e2c049c85a3743d64736f6c63430008070033",
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

// LastAmount is a free data retrieval call binding the contract method 0x829a86d9.
//
// Solidity: function lastAmount() view returns(uint256)
func (_TestDAppV2 *TestDAppV2Caller) LastAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "lastAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastAmount is a free data retrieval call binding the contract method 0x829a86d9.
//
// Solidity: function lastAmount() view returns(uint256)
func (_TestDAppV2 *TestDAppV2Session) LastAmount() (*big.Int, error) {
	return _TestDAppV2.Contract.LastAmount(&_TestDAppV2.CallOpts)
}

// LastAmount is a free data retrieval call binding the contract method 0x829a86d9.
//
// Solidity: function lastAmount() view returns(uint256)
func (_TestDAppV2 *TestDAppV2CallerSession) LastAmount() (*big.Int, error) {
	return _TestDAppV2.Contract.LastAmount(&_TestDAppV2.CallOpts)
}

// LastContext is a free data retrieval call binding the contract method 0xb73f7eb1.
//
// Solidity: function lastContext() view returns(bytes origin, address sender, uint256 chainID)
func (_TestDAppV2 *TestDAppV2Caller) LastContext(opts *bind.CallOpts) (struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "lastContext")

	outstruct := new(struct {
		Origin  []byte
		Sender  common.Address
		ChainID *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Origin = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Sender = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.ChainID = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// LastContext is a free data retrieval call binding the contract method 0xb73f7eb1.
//
// Solidity: function lastContext() view returns(bytes origin, address sender, uint256 chainID)
func (_TestDAppV2 *TestDAppV2Session) LastContext() (struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}, error) {
	return _TestDAppV2.Contract.LastContext(&_TestDAppV2.CallOpts)
}

// LastContext is a free data retrieval call binding the contract method 0xb73f7eb1.
//
// Solidity: function lastContext() view returns(bytes origin, address sender, uint256 chainID)
func (_TestDAppV2 *TestDAppV2CallerSession) LastContext() (struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}, error) {
	return _TestDAppV2.Contract.LastContext(&_TestDAppV2.CallOpts)
}

// LastMessage is a free data retrieval call binding the contract method 0x32970710.
//
// Solidity: function lastMessage() view returns(string)
func (_TestDAppV2 *TestDAppV2Caller) LastMessage(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "lastMessage")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// LastMessage is a free data retrieval call binding the contract method 0x32970710.
//
// Solidity: function lastMessage() view returns(string)
func (_TestDAppV2 *TestDAppV2Session) LastMessage() (string, error) {
	return _TestDAppV2.Contract.LastMessage(&_TestDAppV2.CallOpts)
}

// LastMessage is a free data retrieval call binding the contract method 0x32970710.
//
// Solidity: function lastMessage() view returns(string)
func (_TestDAppV2 *TestDAppV2CallerSession) LastMessage() (string, error) {
	return _TestDAppV2.Contract.LastMessage(&_TestDAppV2.CallOpts)
}

// LastZRC20 is a free data retrieval call binding the contract method 0xb2f79b03.
//
// Solidity: function lastZRC20() view returns(address)
func (_TestDAppV2 *TestDAppV2Caller) LastZRC20(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "lastZRC20")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LastZRC20 is a free data retrieval call binding the contract method 0xb2f79b03.
//
// Solidity: function lastZRC20() view returns(address)
func (_TestDAppV2 *TestDAppV2Session) LastZRC20() (common.Address, error) {
	return _TestDAppV2.Contract.LastZRC20(&_TestDAppV2.CallOpts)
}

// LastZRC20 is a free data retrieval call binding the contract method 0xb2f79b03.
//
// Solidity: function lastZRC20() view returns(address)
func (_TestDAppV2 *TestDAppV2CallerSession) LastZRC20() (common.Address, error) {
	return _TestDAppV2.Contract.LastZRC20(&_TestDAppV2.CallOpts)
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

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2Transactor) OnCrossChainCall(opts *bind.TransactOpts, context TestDAppV2zContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "onCrossChainCall", context, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2Session) OnCrossChainCall(context TestDAppV2zContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCrossChainCall(&_TestDAppV2.TransactOpts, context, zrc20, amount, message)
}

// OnCrossChainCall is a paid mutator transaction binding the contract method 0xde43156e.
//
// Solidity: function onCrossChainCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) OnCrossChainCall(context TestDAppV2zContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCrossChainCall(&_TestDAppV2.TransactOpts, context, zrc20, amount, message)
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
