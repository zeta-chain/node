// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package teststaking

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

// Validator is an auto generated low-level Go binding around an user-defined struct.
type Validator struct {
	OperatorAddress string
	ConsensusPubKey string
	Jailed          bool
	BondStatus      uint8
}

// TestStakingMetaData contains all meta data concerning the TestStaking contract.
var TestStakingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wzeta\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"depositWZETA\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllValidators\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"operatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"consensusPubKey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"jailed\",\"type\":\"bool\"},{\"internalType\":\"enumBondStatus\",\"name\":\"bondStatus\",\"type\":\"uint8\"}],\"internalType\":\"structValidator[]\",\"name\":\"validators\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"getShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validatorSrc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorDst\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"moveStake\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"withdrawWZETA\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405260666000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055503480156200005257600080fd5b50604051620014d3380380620014d383398181016040528101906200007891906200016b565b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555033600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550506200019d565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620001338262000106565b9050919050565b620001458162000126565b81146200015157600080fd5b50565b60008151905062000165816200013a565b92915050565b60006020828403121562000184576200018362000101565b5b6000620001948482850162000154565b91505092915050565b61132680620001ad6000396000f3fe6080604052600436106100745760003560e01c806390b8436f1161004e57806390b8436f146101015780639a0fb6731461013e578063d11a93d014610167578063f3513a37146101a45761007b565b80630d1b3daf1461007d5780632c5d24ae146100ba57806357c6ea3e146100c45761007b565b3661007b57005b005b34801561008957600080fd5b506100a4600480360381019061009f91906109a4565b6101cf565b6040516100b19190610a19565b60405180910390f35b6100c2610276565b005b3480156100d057600080fd5b506100eb60048036038101906100e69190610a60565b610355565b6040516100f89190610aeb565b60405180910390f35b34801561010d57600080fd5b5061012860048036038101906101239190610a60565b61045a565b6040516101359190610b21565b60405180910390f35b34801561014a57600080fd5b5061016560048036038101906101609190610b3c565b61055f565b005b34801561017357600080fd5b5061018e60048036038101906101899190610b69565b610649565b60405161019b9190610aeb565b60405180910390f35b3480156101b057600080fd5b506101b9610751565b6040516101c69190610e42565b60405180910390f35b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630d1b3daf84846040518363ffffffff1660e01b815260040161022d929190610ebd565b602060405180830381865afa15801561024a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061026e9190610f02565b905092915050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102d057600080fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d0e30db0346040518263ffffffff1660e01b81526004016000604051808303818588803b15801561033a57600080fd5b505af115801561034e573d6000803e3d6000fd5b5050505050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146103b157600080fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166357c6ea3e8585856040518463ffffffff1660e01b815260040161040e93929190610f2f565b6020604051808303816000875af115801561042d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104519190610f99565b90509392505050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146104b657600080fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166390b8436f8585856040518463ffffffff1660e01b815260040161051393929190610f2f565b6020604051808303816000875af1158015610532573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105569190610ff2565b90509392505050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146105b957600080fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632e1a7d4d826040518263ffffffff1660e01b81526004016106149190610a19565b600060405180830381600087803b15801561062e57600080fd5b505af1158015610642573d6000803e3d6000fd5b5050505050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106a557600080fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d11a93d0868686866040518563ffffffff1660e01b8152600401610704949392919061101f565b6020604051808303816000875af1158015610723573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107479190610f99565b9050949350505050565b606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f3513a376040518163ffffffff1660e01b8152600401600060405180830381865afa1580156107be573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052508101906107e791906112a7565b905090565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061082b82610800565b9050919050565b61083b81610820565b811461084657600080fd5b50565b60008135905061085881610832565b92915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6108b182610868565b810181811067ffffffffffffffff821117156108d0576108cf610879565b5b80604052505050565b60006108e36107ec565b90506108ef82826108a8565b919050565b600067ffffffffffffffff82111561090f5761090e610879565b5b61091882610868565b9050602081019050919050565b82818337600083830152505050565b6000610947610942846108f4565b6108d9565b90508281526020810184848401111561096357610962610863565b5b61096e848285610925565b509392505050565b600082601f83011261098b5761098a61085e565b5b813561099b848260208601610934565b91505092915050565b600080604083850312156109bb576109ba6107f6565b5b60006109c985828601610849565b925050602083013567ffffffffffffffff8111156109ea576109e96107fb565b5b6109f685828601610976565b9150509250929050565b6000819050919050565b610a1381610a00565b82525050565b6000602082019050610a2e6000830184610a0a565b92915050565b610a3d81610a00565b8114610a4857600080fd5b50565b600081359050610a5a81610a34565b92915050565b600080600060608486031215610a7957610a786107f6565b5b6000610a8786828701610849565b935050602084013567ffffffffffffffff811115610aa857610aa76107fb565b5b610ab486828701610976565b9250506040610ac586828701610a4b565b9150509250925092565b60008160070b9050919050565b610ae581610acf565b82525050565b6000602082019050610b006000830184610adc565b92915050565b60008115159050919050565b610b1b81610b06565b82525050565b6000602082019050610b366000830184610b12565b92915050565b600060208284031215610b5257610b516107f6565b5b6000610b6084828501610a4b565b91505092915050565b60008060008060808587031215610b8357610b826107f6565b5b6000610b9187828801610849565b945050602085013567ffffffffffffffff811115610bb257610bb16107fb565b5b610bbe87828801610976565b935050604085013567ffffffffffffffff811115610bdf57610bde6107fb565b5b610beb87828801610976565b9250506060610bfc87828801610a4b565b91505092959194509250565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610c6e578082015181840152602081019050610c53565b83811115610c7d576000848401525b50505050565b6000610c8e82610c34565b610c988185610c3f565b9350610ca8818560208601610c50565b610cb181610868565b840191505092915050565b610cc581610b06565b82525050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60048110610d0b57610d0a610ccb565b5b50565b6000819050610d1c82610cfa565b919050565b6000610d2c82610d0e565b9050919050565b610d3c81610d21565b82525050565b60006080830160008301518482036000860152610d5f8282610c83565b91505060208301518482036020860152610d798282610c83565b9150506040830151610d8e6040860182610cbc565b506060830151610da16060860182610d33565b508091505092915050565b6000610db88383610d42565b905092915050565b6000602082019050919050565b6000610dd882610c08565b610de28185610c13565b935083602082028501610df485610c24565b8060005b85811015610e305784840389528151610e118582610dac565b9450610e1c83610dc0565b925060208a01995050600181019050610df8565b50829750879550505050505092915050565b60006020820190508181036000830152610e5c8184610dcd565b905092915050565b610e6d81610820565b82525050565b600082825260208201905092915050565b6000610e8f82610c34565b610e998185610e73565b9350610ea9818560208601610c50565b610eb281610868565b840191505092915050565b6000604082019050610ed26000830185610e64565b8181036020830152610ee48184610e84565b90509392505050565b600081519050610efc81610a34565b92915050565b600060208284031215610f1857610f176107f6565b5b6000610f2684828501610eed565b91505092915050565b6000606082019050610f446000830186610e64565b8181036020830152610f568185610e84565b9050610f656040830184610a0a565b949350505050565b610f7681610acf565b8114610f8157600080fd5b50565b600081519050610f9381610f6d565b92915050565b600060208284031215610faf57610fae6107f6565b5b6000610fbd84828501610f84565b91505092915050565b610fcf81610b06565b8114610fda57600080fd5b50565b600081519050610fec81610fc6565b92915050565b600060208284031215611008576110076107f6565b5b600061101684828501610fdd565b91505092915050565b60006080820190506110346000830187610e64565b81810360208301526110468186610e84565b9050818103604083015261105a8185610e84565b90506110696060830184610a0a565b95945050505050565b600067ffffffffffffffff82111561108d5761108c610879565b5b602082029050602081019050919050565b600080fd5b600080fd5b600080fd5b60006110c06110bb846108f4565b6108d9565b9050828152602081018484840111156110dc576110db610863565b5b6110e7848285610c50565b509392505050565b600082601f8301126111045761110361085e565b5b81516111148482602086016110ad565b91505092915050565b6004811061112a57600080fd5b50565b60008151905061113c8161111d565b92915050565b600060808284031215611158576111576110a3565b5b61116260806108d9565b9050600082015167ffffffffffffffff811115611182576111816110a8565b5b61118e848285016110ef565b600083015250602082015167ffffffffffffffff8111156111b2576111b16110a8565b5b6111be848285016110ef565b60208301525060406111d284828501610fdd565b60408301525060606111e68482850161112d565b60608301525092915050565b600061120561120084611072565b6108d9565b905080838252602082019050602084028301858111156112285761122761109e565b5b835b8181101561126f57805167ffffffffffffffff81111561124d5761124c61085e565b5b80860161125a8982611142565b8552602085019450505060208101905061122a565b5050509392505050565b600082601f83011261128e5761128d61085e565b5b815161129e8482602086016111f2565b91505092915050565b6000602082840312156112bd576112bc6107f6565b5b600082015167ffffffffffffffff8111156112db576112da6107fb565b5b6112e784828501611279565b9150509291505056fea2646970667358221220aba1378d6378c171ad8af1a40b9be85595c4153d6c0d1476a13be035b4e9e36a64736f6c634300080a0033",
}

// TestStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use TestStakingMetaData.ABI instead.
var TestStakingABI = TestStakingMetaData.ABI

// TestStakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestStakingMetaData.Bin instead.
var TestStakingBin = TestStakingMetaData.Bin

// DeployTestStaking deploys a new Ethereum contract, binding an instance of TestStaking to it.
func DeployTestStaking(auth *bind.TransactOpts, backend bind.ContractBackend, _wzeta common.Address) (common.Address, *types.Transaction, *TestStaking, error) {
	parsed, err := TestStakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestStakingBin), backend, _wzeta)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestStaking{TestStakingCaller: TestStakingCaller{contract: contract}, TestStakingTransactor: TestStakingTransactor{contract: contract}, TestStakingFilterer: TestStakingFilterer{contract: contract}}, nil
}

// TestStaking is an auto generated Go binding around an Ethereum contract.
type TestStaking struct {
	TestStakingCaller     // Read-only binding to the contract
	TestStakingTransactor // Write-only binding to the contract
	TestStakingFilterer   // Log filterer for contract events
}

// TestStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestStakingSession struct {
	Contract     *TestStaking      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestStakingCallerSession struct {
	Contract *TestStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// TestStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestStakingTransactorSession struct {
	Contract     *TestStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// TestStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestStakingRaw struct {
	Contract *TestStaking // Generic contract binding to access the raw methods on
}

// TestStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestStakingCallerRaw struct {
	Contract *TestStakingCaller // Generic read-only contract binding to access the raw methods on
}

// TestStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestStakingTransactorRaw struct {
	Contract *TestStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestStaking creates a new instance of TestStaking, bound to a specific deployed contract.
func NewTestStaking(address common.Address, backend bind.ContractBackend) (*TestStaking, error) {
	contract, err := bindTestStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestStaking{TestStakingCaller: TestStakingCaller{contract: contract}, TestStakingTransactor: TestStakingTransactor{contract: contract}, TestStakingFilterer: TestStakingFilterer{contract: contract}}, nil
}

// NewTestStakingCaller creates a new read-only instance of TestStaking, bound to a specific deployed contract.
func NewTestStakingCaller(address common.Address, caller bind.ContractCaller) (*TestStakingCaller, error) {
	contract, err := bindTestStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestStakingCaller{contract: contract}, nil
}

// NewTestStakingTransactor creates a new write-only instance of TestStaking, bound to a specific deployed contract.
func NewTestStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*TestStakingTransactor, error) {
	contract, err := bindTestStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestStakingTransactor{contract: contract}, nil
}

// NewTestStakingFilterer creates a new log filterer instance of TestStaking, bound to a specific deployed contract.
func NewTestStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*TestStakingFilterer, error) {
	contract, err := bindTestStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestStakingFilterer{contract: contract}, nil
}

// bindTestStaking binds a generic wrapper to an already deployed contract.
func bindTestStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestStakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestStaking *TestStakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestStaking.Contract.TestStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestStaking *TestStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestStaking.Contract.TestStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestStaking *TestStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestStaking.Contract.TestStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestStaking *TestStakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestStaking *TestStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestStaking *TestStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestStaking.Contract.contract.Transact(opts, method, params...)
}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((string,string,bool,uint8)[] validators)
func (_TestStaking *TestStakingCaller) GetAllValidators(opts *bind.CallOpts) ([]Validator, error) {
	var out []interface{}
	err := _TestStaking.contract.Call(opts, &out, "getAllValidators")

	if err != nil {
		return *new([]Validator), err
	}

	out0 := *abi.ConvertType(out[0], new([]Validator)).(*[]Validator)

	return out0, err

}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((string,string,bool,uint8)[] validators)
func (_TestStaking *TestStakingSession) GetAllValidators() ([]Validator, error) {
	return _TestStaking.Contract.GetAllValidators(&_TestStaking.CallOpts)
}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((string,string,bool,uint8)[] validators)
func (_TestStaking *TestStakingCallerSession) GetAllValidators() ([]Validator, error) {
	return _TestStaking.Contract.GetAllValidators(&_TestStaking.CallOpts)
}

// GetShares is a free data retrieval call binding the contract method 0x0d1b3daf.
//
// Solidity: function getShares(address staker, string validator) view returns(uint256 shares)
func (_TestStaking *TestStakingCaller) GetShares(opts *bind.CallOpts, staker common.Address, validator string) (*big.Int, error) {
	var out []interface{}
	err := _TestStaking.contract.Call(opts, &out, "getShares", staker, validator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetShares is a free data retrieval call binding the contract method 0x0d1b3daf.
//
// Solidity: function getShares(address staker, string validator) view returns(uint256 shares)
func (_TestStaking *TestStakingSession) GetShares(staker common.Address, validator string) (*big.Int, error) {
	return _TestStaking.Contract.GetShares(&_TestStaking.CallOpts, staker, validator)
}

// GetShares is a free data retrieval call binding the contract method 0x0d1b3daf.
//
// Solidity: function getShares(address staker, string validator) view returns(uint256 shares)
func (_TestStaking *TestStakingCallerSession) GetShares(staker common.Address, validator string) (*big.Int, error) {
	return _TestStaking.Contract.GetShares(&_TestStaking.CallOpts, staker, validator)
}

// DepositWZETA is a paid mutator transaction binding the contract method 0x2c5d24ae.
//
// Solidity: function depositWZETA() payable returns()
func (_TestStaking *TestStakingTransactor) DepositWZETA(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestStaking.contract.Transact(opts, "depositWZETA")
}

// DepositWZETA is a paid mutator transaction binding the contract method 0x2c5d24ae.
//
// Solidity: function depositWZETA() payable returns()
func (_TestStaking *TestStakingSession) DepositWZETA() (*types.Transaction, error) {
	return _TestStaking.Contract.DepositWZETA(&_TestStaking.TransactOpts)
}

// DepositWZETA is a paid mutator transaction binding the contract method 0x2c5d24ae.
//
// Solidity: function depositWZETA() payable returns()
func (_TestStaking *TestStakingTransactorSession) DepositWZETA() (*types.Transaction, error) {
	return _TestStaking.Contract.DepositWZETA(&_TestStaking.TransactOpts)
}

// MoveStake is a paid mutator transaction binding the contract method 0xd11a93d0.
//
// Solidity: function moveStake(address staker, string validatorSrc, string validatorDst, uint256 amount) returns(int64 completionTime)
func (_TestStaking *TestStakingTransactor) MoveStake(opts *bind.TransactOpts, staker common.Address, validatorSrc string, validatorDst string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.contract.Transact(opts, "moveStake", staker, validatorSrc, validatorDst, amount)
}

// MoveStake is a paid mutator transaction binding the contract method 0xd11a93d0.
//
// Solidity: function moveStake(address staker, string validatorSrc, string validatorDst, uint256 amount) returns(int64 completionTime)
func (_TestStaking *TestStakingSession) MoveStake(staker common.Address, validatorSrc string, validatorDst string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.MoveStake(&_TestStaking.TransactOpts, staker, validatorSrc, validatorDst, amount)
}

// MoveStake is a paid mutator transaction binding the contract method 0xd11a93d0.
//
// Solidity: function moveStake(address staker, string validatorSrc, string validatorDst, uint256 amount) returns(int64 completionTime)
func (_TestStaking *TestStakingTransactorSession) MoveStake(staker common.Address, validatorSrc string, validatorDst string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.MoveStake(&_TestStaking.TransactOpts, staker, validatorSrc, validatorDst, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x90b8436f.
//
// Solidity: function stake(address staker, string validator, uint256 amount) returns(bool)
func (_TestStaking *TestStakingTransactor) Stake(opts *bind.TransactOpts, staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.contract.Transact(opts, "stake", staker, validator, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x90b8436f.
//
// Solidity: function stake(address staker, string validator, uint256 amount) returns(bool)
func (_TestStaking *TestStakingSession) Stake(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.Stake(&_TestStaking.TransactOpts, staker, validator, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x90b8436f.
//
// Solidity: function stake(address staker, string validator, uint256 amount) returns(bool)
func (_TestStaking *TestStakingTransactorSession) Stake(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.Stake(&_TestStaking.TransactOpts, staker, validator, amount)
}

// Unstake is a paid mutator transaction binding the contract method 0x57c6ea3e.
//
// Solidity: function unstake(address staker, string validator, uint256 amount) returns(int64 completionTime)
func (_TestStaking *TestStakingTransactor) Unstake(opts *bind.TransactOpts, staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.contract.Transact(opts, "unstake", staker, validator, amount)
}

// Unstake is a paid mutator transaction binding the contract method 0x57c6ea3e.
//
// Solidity: function unstake(address staker, string validator, uint256 amount) returns(int64 completionTime)
func (_TestStaking *TestStakingSession) Unstake(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.Unstake(&_TestStaking.TransactOpts, staker, validator, amount)
}

// Unstake is a paid mutator transaction binding the contract method 0x57c6ea3e.
//
// Solidity: function unstake(address staker, string validator, uint256 amount) returns(int64 completionTime)
func (_TestStaking *TestStakingTransactorSession) Unstake(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.Unstake(&_TestStaking.TransactOpts, staker, validator, amount)
}

// WithdrawWZETA is a paid mutator transaction binding the contract method 0x9a0fb673.
//
// Solidity: function withdrawWZETA(uint256 wad) returns()
func (_TestStaking *TestStakingTransactor) WithdrawWZETA(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error) {
	return _TestStaking.contract.Transact(opts, "withdrawWZETA", wad)
}

// WithdrawWZETA is a paid mutator transaction binding the contract method 0x9a0fb673.
//
// Solidity: function withdrawWZETA(uint256 wad) returns()
func (_TestStaking *TestStakingSession) WithdrawWZETA(wad *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.WithdrawWZETA(&_TestStaking.TransactOpts, wad)
}

// WithdrawWZETA is a paid mutator transaction binding the contract method 0x9a0fb673.
//
// Solidity: function withdrawWZETA(uint256 wad) returns()
func (_TestStaking *TestStakingTransactorSession) WithdrawWZETA(wad *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.WithdrawWZETA(&_TestStaking.TransactOpts, wad)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestStaking *TestStakingTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TestStaking.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestStaking *TestStakingSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TestStaking.Contract.Fallback(&_TestStaking.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TestStaking *TestStakingTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TestStaking.Contract.Fallback(&_TestStaking.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestStaking *TestStakingTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestStaking.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestStaking *TestStakingSession) Receive() (*types.Transaction, error) {
	return _TestStaking.Contract.Receive(&_TestStaking.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_TestStaking *TestStakingTransactorSession) Receive() (*types.Transaction, error) {
	return _TestStaking.Contract.Receive(&_TestStaking.TransactOpts)
}
