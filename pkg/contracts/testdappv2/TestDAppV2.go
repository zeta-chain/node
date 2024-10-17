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
	Sender        common.Address
	Asset         common.Address
	Amount        *big.Int
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
	ABI: "[{\"inputs\":[],\"name\":\"NO_MESSAGE_CALL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"amountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"calledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"erc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"erc20Call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"expectedOnCallSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"gasCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getAmountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getCalledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structTestDAppV2.zContext\",\"name\":\"_context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structTestDAppV2.MessageContext\",\"name\":\"messageContext\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structTestDAppV2.zContext\",\"name\":\"_context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCrossChainCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"}],\"internalType\":\"structTestDAppV2.RevertContext\",\"name\":\"revertContext\",\"type\":\"tuple\"}],\"name\":\"onRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"senderWithMessage\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_expectedOnCallSender\",\"type\":\"address\"}],\"name\":\"setExpectedOnCallSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"simpleCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506117698061001f6000396000f3fe6080604052600436106100ec5760003560e01c8063c234fecf1161008a578063de43156e11610059578063de43156e146102e1578063e2842ed71461030a578063f592cbfb14610347578063f936ae8514610384576100f3565b8063c234fecf1461023b578063c7a339a914610264578063c85f84341461028d578063c9028a36146102b8576100f3565b80635bcfd616116100c65780635bcfd61614610189578063676cc054146101b25780639291fe26146101e2578063a799911f1461021f576100f3565b806336e980a0146100f85780634297a2631461012157806359f4a7771461015e576100f3565b366100f357005b600080fd5b34801561010457600080fd5b5061011f600480360381019061011a9190610dcb565b6103c1565b005b34801561012d57600080fd5b5061014860048036038101906101439190610e4a565b6103eb565b6040516101559190610e90565b60405180910390f35b34801561016a57600080fd5b50610173610403565b6040516101809190610eec565b60405180910390f35b34801561019557600080fd5b506101b060048036038101906101ab9190610fe3565b610427565b005b6101cc60048036038101906101c791906110a6565b610528565b6040516101d99190611185565b60405180910390f35b3480156101ee57600080fd5b5061020960048036038101906102049190610dcb565b6106f5565b6040516102169190610e90565b60405180910390f35b61023960048036038101906102349190610dcb565b610738565b005b34801561024757600080fd5b50610262600480360381019061025d91906111a7565b610761565b005b34801561027057600080fd5b5061028b60048036038101906102869190611212565b6107a4565b005b34801561029957600080fd5b506102a2610858565b6040516102af91906112d6565b60405180910390f35b3480156102c457600080fd5b506102df60048036038101906102da9190611317565b610891565b005b3480156102ed57600080fd5b5061030860048036038101906103039190610fe3565b6109cb565b005b34801561031657600080fd5b50610331600480360381019061032c9190610e4a565b610acc565b60405161033e919061137b565b60405180910390f35b34801561035357600080fd5b5061036e60048036038101906103699190610dcb565b610aec565b60405161037b919061137b565b60405180910390f35b34801561039057600080fd5b506103ab60048036038101906103a69190611437565b610b3c565b6040516103b89190610eec565b60405180910390f35b6103ca81610b85565b156103d457600080fd5b6103dd81610bdb565b6103e8816000610c2f565b50565b60036020528060005260406000206000915090505481565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b61047482828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610b85565b1561047e57600080fd5b60008083839050146104d45782828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505061050b565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d657373616765000000000000000000008152505b905061051681610bdb565b6105208185610c2f565b505050505050565b606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1684600001602081019061057491906111a7565b73ffffffffffffffffffffffffffffffffffffffff16146105ca576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105c1906114cc565b60405180910390fd5b60008084849050146106205783838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610657565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d657373616765000000000000000000008152505b905061066281610bdb565b61066c8134610c2f565b84600001602081019061067f91906111a7565b60028260405161068f9190611528565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604051806020016040528060008152509150509392505050565b6000600360008360405160200161070c919061157b565b604051602081830303815290604052805190602001208152602001908152602001600020549050919050565b61074181610b85565b1561074b57600080fd5b61075481610bdb565b61075e8134610c2f565b50565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6107ad81610b85565b156107b757600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b81526004016107f493929190611592565b6020604051808303816000875af1158015610813573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061083791906115f5565b61084057600080fd5b61084981610bdb565b6108538183610c2f565b505050565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525081565b6108ec8180606001906108a49190611631565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610bdb565b6109498180606001906108ff9190611631565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506000610c2f565b80600001602081019061095c91906111a7565b600282806060019061096e9190611631565b60405161097c9291906116b9565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b610a1882828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610b85565b15610a2257600080fd5b6000808383905014610a785782828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610aaf565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d657373616765000000000000000000008152505b9050610aba81610bdb565b610ac48185610c2f565b505050505050565b60016020528060005260406000206000915054906101000a900460ff1681565b60006001600083604051602001610b03919061157b565b60405160208183030381529060405280519060200120815260200190815260200160002060009054906101000a900460ff169050919050565b6002818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000604051602001610b969061171e565b6040516020818303038152906040528051906020012082604051602001610bbd919061157b565b60405160208183030381529060405280519060200120149050919050565b600180600083604051602001610bf1919061157b565b60405160208183030381529060405280519060200120815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b806003600084604051602001610c45919061157b565b604051602081830303815290604052805190602001208152602001908152602001600020819055505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610cd882610c8f565b810181811067ffffffffffffffff82111715610cf757610cf6610ca0565b5b80604052505050565b6000610d0a610c71565b9050610d168282610ccf565b919050565b600067ffffffffffffffff821115610d3657610d35610ca0565b5b610d3f82610c8f565b9050602081019050919050565b82818337600083830152505050565b6000610d6e610d6984610d1b565b610d00565b905082815260208101848484011115610d8a57610d89610c8a565b5b610d95848285610d4c565b509392505050565b600082601f830112610db257610db1610c85565b5b8135610dc2848260208601610d5b565b91505092915050565b600060208284031215610de157610de0610c7b565b5b600082013567ffffffffffffffff811115610dff57610dfe610c80565b5b610e0b84828501610d9d565b91505092915050565b6000819050919050565b610e2781610e14565b8114610e3257600080fd5b50565b600081359050610e4481610e1e565b92915050565b600060208284031215610e6057610e5f610c7b565b5b6000610e6e84828501610e35565b91505092915050565b6000819050919050565b610e8a81610e77565b82525050565b6000602082019050610ea56000830184610e81565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610ed682610eab565b9050919050565b610ee681610ecb565b82525050565b6000602082019050610f016000830184610edd565b92915050565b600080fd5b600060608284031215610f2257610f21610f07565b5b81905092915050565b610f3481610ecb565b8114610f3f57600080fd5b50565b600081359050610f5181610f2b565b92915050565b610f6081610e77565b8114610f6b57600080fd5b50565b600081359050610f7d81610f57565b92915050565b600080fd5b600080fd5b60008083601f840112610fa357610fa2610c85565b5b8235905067ffffffffffffffff811115610fc057610fbf610f83565b5b602083019150836001820283011115610fdc57610fdb610f88565b5b9250929050565b600080600080600060808688031215610fff57610ffe610c7b565b5b600086013567ffffffffffffffff81111561101d5761101c610c80565b5b61102988828901610f0c565b955050602061103a88828901610f42565b945050604061104b88828901610f6e565b935050606086013567ffffffffffffffff81111561106c5761106b610c80565b5b61107888828901610f8d565b92509250509295509295909350565b60006020828403121561109d5761109c610f07565b5b81905092915050565b6000806000604084860312156110bf576110be610c7b565b5b60006110cd86828701611087565b935050602084013567ffffffffffffffff8111156110ee576110ed610c80565b5b6110fa86828701610f8d565b92509250509250925092565b600081519050919050565b600082825260208201905092915050565b60005b83811015611140578082015181840152602081019050611125565b60008484015250505050565b600061115782611106565b6111618185611111565b9350611171818560208601611122565b61117a81610c8f565b840191505092915050565b6000602082019050818103600083015261119f818461114c565b905092915050565b6000602082840312156111bd576111bc610c7b565b5b60006111cb84828501610f42565b91505092915050565b60006111df82610ecb565b9050919050565b6111ef816111d4565b81146111fa57600080fd5b50565b60008135905061120c816111e6565b92915050565b60008060006060848603121561122b5761122a610c7b565b5b6000611239868287016111fd565b935050602061124a86828701610f6e565b925050604084013567ffffffffffffffff81111561126b5761126a610c80565b5b61127786828701610d9d565b9150509250925092565b600081519050919050565b600082825260208201905092915050565b60006112a882611281565b6112b2818561128c565b93506112c2818560208601611122565b6112cb81610c8f565b840191505092915050565b600060208201905081810360008301526112f0818461129d565b905092915050565b60006080828403121561130e5761130d610f07565b5b81905092915050565b60006020828403121561132d5761132c610c7b565b5b600082013567ffffffffffffffff81111561134b5761134a610c80565b5b611357848285016112f8565b91505092915050565b60008115159050919050565b61137581611360565b82525050565b6000602082019050611390600083018461136c565b92915050565b600067ffffffffffffffff8211156113b1576113b0610ca0565b5b6113ba82610c8f565b9050602081019050919050565b60006113da6113d584611396565b610d00565b9050828152602081018484840111156113f6576113f5610c8a565b5b611401848285610d4c565b509392505050565b600082601f83011261141e5761141d610c85565b5b813561142e8482602086016113c7565b91505092915050565b60006020828403121561144d5761144c610c7b565b5b600082013567ffffffffffffffff81111561146b5761146a610c80565b5b61147784828501611409565b91505092915050565b7f756e61757468656e746963617465642073656e64657200000000000000000000600082015250565b60006114b660168361128c565b91506114c182611480565b602082019050919050565b600060208201905081810360008301526114e5816114a9565b9050919050565b600081905092915050565b600061150282611106565b61150c81856114ec565b935061151c818560208601611122565b80840191505092915050565b600061153482846114f7565b915081905092915050565b600081905092915050565b600061155582611281565b61155f818561153f565b935061156f818560208601611122565b80840191505092915050565b6000611587828461154a565b915081905092915050565b60006060820190506115a76000830186610edd565b6115b46020830185610edd565b6115c16040830184610e81565b949350505050565b6115d281611360565b81146115dd57600080fd5b50565b6000815190506115ef816115c9565b92915050565b60006020828403121561160b5761160a610c7b565b5b6000611619848285016115e0565b91505092915050565b600080fd5b600080fd5b600080fd5b6000808335600160200384360303811261164e5761164d611622565b5b80840192508235915067ffffffffffffffff8211156116705761166f611627565b5b60208301925060018202360383131561168c5761168b61162c565b5b509250929050565b60006116a083856114ec565b93506116ad838584610d4c565b82840190509392505050565b60006116c6828486611694565b91508190509392505050565b7f7265766572740000000000000000000000000000000000000000000000000000600082015250565b600061170860068361153f565b9150611713826116d2565b600682019050919050565b6000611729826116fb565b915081905091905056fea2646970667358221220e61be08abc7bb184a794396b2202787fdb9ea3ead91ff0c14520235faf2cb7b164736f6c634300081a0033",
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

// NOMESSAGECALL is a free data retrieval call binding the contract method 0xc85f8434.
//
// Solidity: function NO_MESSAGE_CALL() view returns(string)
func (_TestDAppV2 *TestDAppV2Caller) NOMESSAGECALL(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "NO_MESSAGE_CALL")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NOMESSAGECALL is a free data retrieval call binding the contract method 0xc85f8434.
//
// Solidity: function NO_MESSAGE_CALL() view returns(string)
func (_TestDAppV2 *TestDAppV2Session) NOMESSAGECALL() (string, error) {
	return _TestDAppV2.Contract.NOMESSAGECALL(&_TestDAppV2.CallOpts)
}

// NOMESSAGECALL is a free data retrieval call binding the contract method 0xc85f8434.
//
// Solidity: function NO_MESSAGE_CALL() view returns(string)
func (_TestDAppV2 *TestDAppV2CallerSession) NOMESSAGECALL() (string, error) {
	return _TestDAppV2.Contract.NOMESSAGECALL(&_TestDAppV2.CallOpts)
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

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) _context, address _zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2Transactor) OnCall(opts *bind.TransactOpts, _context TestDAppV2zContext, _zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "onCall", _context, _zrc20, amount, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) _context, address _zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2Session) OnCall(_context TestDAppV2zContext, _zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCall(&_TestDAppV2.TransactOpts, _context, _zrc20, amount, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) _context, address _zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) OnCall(_context TestDAppV2zContext, _zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCall(&_TestDAppV2.TransactOpts, _context, _zrc20, amount, message)
}

// OnCall0 is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppV2 *TestDAppV2Transactor) OnCall0(opts *bind.TransactOpts, messageContext TestDAppV2MessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "onCall0", messageContext, message)
}

// OnCall0 is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppV2 *TestDAppV2Session) OnCall0(messageContext TestDAppV2MessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCall0(&_TestDAppV2.TransactOpts, messageContext, message)
}

// OnCall0 is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppV2 *TestDAppV2TransactorSession) OnCall0(messageContext TestDAppV2MessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCall0(&_TestDAppV2.TransactOpts, messageContext, message)
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

// OnRevert is a paid mutator transaction binding the contract method 0xc9028a36.
//
// Solidity: function onRevert((address,address,uint256,bytes) revertContext) returns()
func (_TestDAppV2 *TestDAppV2Transactor) OnRevert(opts *bind.TransactOpts, revertContext TestDAppV2RevertContext) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "onRevert", revertContext)
}

// OnRevert is a paid mutator transaction binding the contract method 0xc9028a36.
//
// Solidity: function onRevert((address,address,uint256,bytes) revertContext) returns()
func (_TestDAppV2 *TestDAppV2Session) OnRevert(revertContext TestDAppV2RevertContext) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnRevert(&_TestDAppV2.TransactOpts, revertContext)
}

// OnRevert is a paid mutator transaction binding the contract method 0xc9028a36.
//
// Solidity: function onRevert((address,address,uint256,bytes) revertContext) returns()
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
