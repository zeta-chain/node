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
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"isZetaChain_\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"NO_MESSAGE_CALL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"amountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"calledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"erc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"erc20Call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"gasCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getAmountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getCalledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getNoMessageIndex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isZetaChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structTestDAppV2.zContext\",\"name\":\"_context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structTestDAppV2.MessageContext\",\"name\":\"messageContext\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"}],\"internalType\":\"structTestDAppV2.RevertContext\",\"name\":\"revertContext\",\"type\":\"tuple\"}],\"name\":\"onRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"senderWithMessage\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"simpleCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60a0604052348015600f57600080fd5b506040516117953803806117958339818101604052810190602f91906078565b8015156080811515815250505060a0565b600080fd5b60008115159050919050565b6058816045565b8114606257600080fd5b50565b6000815190506072816051565b92915050565b600060208284031215608b57608a6040565b5b60006097848285016065565b91505092915050565b6080516116d36100c26000396000818161079601526108ff01526116d36000f3fe6080604052600436106100e15760003560e01c8063c7a339a91161007f578063c91f356711610059578063c91f3567146102bf578063e2842ed7146102ea578063f592cbfb14610327578063f936ae8514610364576100e8565b8063c7a339a914610242578063c85f84341461026b578063c9028a3614610296576100e8565b8063676cc054116100bb578063676cc0541461017c5780639291fe26146101ac578063a799911f146101e9578063ad23b28b14610205576100e8565b806336e980a0146100ed5780634297a263146101165780635bcfd61614610153576100e8565b366100e857005b600080fd5b3480156100f957600080fd5b50610114600480360381019061010f9190610cd1565b6103a1565b005b34801561012257600080fd5b5061013d60048036038101906101389190610d50565b6103cb565b60405161014a9190610d96565b60405180910390f35b34801561015f57600080fd5b5061017a60048036038101906101759190610ebf565b6103e3565b005b61019660048036038101906101919190610f82565b6104c9565b6040516101a39190611061565b60405180910390f35b3480156101b857600080fd5b506101d360048036038101906101ce9190610cd1565b6105db565b6040516101e09190610d96565b60405180910390f35b61020360048036038101906101fe9190610cd1565b61061e565b005b34801561021157600080fd5b5061022c60048036038101906102279190611083565b610647565b6040516102399190611105565b60405180910390f35b34801561024e57600080fd5b5061026960048036038101906102649190611165565b6106a7565b005b34801561027757600080fd5b5061028061075b565b60405161028d9190611105565b60405180910390f35b3480156102a257600080fd5b506102bd60048036038101906102b891906111f3565b610794565b005b3480156102cb57600080fd5b506102d46108fd565b6040516102e19190611257565b60405180910390f35b3480156102f657600080fd5b50610311600480360381019061030c9190610d50565b610921565b60405161031e9190611257565b60405180910390f35b34801561033357600080fd5b5061034e60048036038101906103499190610cd1565b610941565b60405161035b9190611257565b60405180910390f35b34801561037057600080fd5b5061038b60048036038101906103869190611313565b610991565b604051610398919061136b565b60405180910390f35b6103aa816109da565b156103b457600080fd5b6103bd81610a30565b6103c8816000610a84565b50565b60036020528060005260406000206000915090505481565b61043082828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506109da565b1561043a57600080fd5b60008083839050146104905782828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506104ac565b6104ab8660200160208101906104a69190611083565b610647565b5b90506104b781610a30565b6104c18185610a84565b505050505050565b606060008084849050146105215783838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505061053d565b61053c8560000160208101906105379190611083565b610647565b5b905061054881610a30565b6105528134610a84565b8460000160208101906105659190611083565b60028260405161057591906113c2565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604051806020016040528060008152509150509392505050565b600060036000836040516020016105f29190611415565b604051602081830303815290604052805190602001208152602001908152602001600020549050919050565b610627816109da565b1561063157600080fd5b61063a81610a30565b6106448134610a84565b50565b60606040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525082604051602001610691929190611474565b6040516020818303038152906040529050919050565b6106b0816109da565b156106ba57600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b81526004016106f79392919061149c565b6020604051808303816000875af1158015610716573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061073a91906114ff565b61074357600080fd5b61074c81610a30565b6107568183610a84565b505050565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525081565b7f0000000000000000000000000000000000000000000000000000000000000000156107c3576107c2610ac6565b5b61081e8180606001906107d6919061153b565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610a30565b61087b818060600190610831919061153b565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506000610a84565b80600001602081019061088e9190611083565b60028280606001906108a0919061153b565b6040516108ae9291906115c3565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b7f000000000000000000000000000000000000000000000000000000000000000081565b60016020528060005260406000206000915054906101000a900460ff1681565b600060016000836040516020016109589190611415565b60405160208183030381529060405280519060200120815260200190815260200160002060009054906101000a900460ff169050919050565b6002818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60006040516020016109eb90611628565b6040516020818303038152906040528051906020012082604051602001610a129190611415565b60405160208183030381529060405280519060200120149050919050565b600180600083604051602001610a469190611415565b60405160208183030381529060405280519060200120815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b806003600084604051602001610a9a9190611415565b604051602081830303815290604052805190602001208152602001908152602001600020819055505050565b60006207a12090506000614e20905060008183610ae3919061166c565b905060005b81811015610b265760008190806001815401808255809150506001900390600052602060002001600090919091909150558080600101915050610ae8565b50600080610b349190610b39565b505050565b5080546000825590600052602060002090810190610b579190610b5a565b50565b5b80821115610b73576000816000905550600101610b5b565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610bde82610b95565b810181811067ffffffffffffffff82111715610bfd57610bfc610ba6565b5b80604052505050565b6000610c10610b77565b9050610c1c8282610bd5565b919050565b600067ffffffffffffffff821115610c3c57610c3b610ba6565b5b610c4582610b95565b9050602081019050919050565b82818337600083830152505050565b6000610c74610c6f84610c21565b610c06565b905082815260208101848484011115610c9057610c8f610b90565b5b610c9b848285610c52565b509392505050565b600082601f830112610cb857610cb7610b8b565b5b8135610cc8848260208601610c61565b91505092915050565b600060208284031215610ce757610ce6610b81565b5b600082013567ffffffffffffffff811115610d0557610d04610b86565b5b610d1184828501610ca3565b91505092915050565b6000819050919050565b610d2d81610d1a565b8114610d3857600080fd5b50565b600081359050610d4a81610d24565b92915050565b600060208284031215610d6657610d65610b81565b5b6000610d7484828501610d3b565b91505092915050565b6000819050919050565b610d9081610d7d565b82525050565b6000602082019050610dab6000830184610d87565b92915050565b600080fd5b600060608284031215610dcc57610dcb610db1565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610e0082610dd5565b9050919050565b610e1081610df5565b8114610e1b57600080fd5b50565b600081359050610e2d81610e07565b92915050565b610e3c81610d7d565b8114610e4757600080fd5b50565b600081359050610e5981610e33565b92915050565b600080fd5b600080fd5b60008083601f840112610e7f57610e7e610b8b565b5b8235905067ffffffffffffffff811115610e9c57610e9b610e5f565b5b602083019150836001820283011115610eb857610eb7610e64565b5b9250929050565b600080600080600060808688031215610edb57610eda610b81565b5b600086013567ffffffffffffffff811115610ef957610ef8610b86565b5b610f0588828901610db6565b9550506020610f1688828901610e1e565b9450506040610f2788828901610e4a565b935050606086013567ffffffffffffffff811115610f4857610f47610b86565b5b610f5488828901610e69565b92509250509295509295909350565b600060208284031215610f7957610f78610db1565b5b81905092915050565b600080600060408486031215610f9b57610f9a610b81565b5b6000610fa986828701610f63565b935050602084013567ffffffffffffffff811115610fca57610fc9610b86565b5b610fd686828701610e69565b92509250509250925092565b600081519050919050565b600082825260208201905092915050565b60005b8381101561101c578082015181840152602081019050611001565b60008484015250505050565b600061103382610fe2565b61103d8185610fed565b935061104d818560208601610ffe565b61105681610b95565b840191505092915050565b6000602082019050818103600083015261107b8184611028565b905092915050565b60006020828403121561109957611098610b81565b5b60006110a784828501610e1e565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60006110d7826110b0565b6110e181856110bb565b93506110f1818560208601610ffe565b6110fa81610b95565b840191505092915050565b6000602082019050818103600083015261111f81846110cc565b905092915050565b600061113282610df5565b9050919050565b61114281611127565b811461114d57600080fd5b50565b60008135905061115f81611139565b92915050565b60008060006060848603121561117e5761117d610b81565b5b600061118c86828701611150565b935050602061119d86828701610e4a565b925050604084013567ffffffffffffffff8111156111be576111bd610b86565b5b6111ca86828701610ca3565b9150509250925092565b6000608082840312156111ea576111e9610db1565b5b81905092915050565b60006020828403121561120957611208610b81565b5b600082013567ffffffffffffffff81111561122757611226610b86565b5b611233848285016111d4565b91505092915050565b60008115159050919050565b6112518161123c565b82525050565b600060208201905061126c6000830184611248565b92915050565b600067ffffffffffffffff82111561128d5761128c610ba6565b5b61129682610b95565b9050602081019050919050565b60006112b66112b184611272565b610c06565b9050828152602081018484840111156112d2576112d1610b90565b5b6112dd848285610c52565b509392505050565b600082601f8301126112fa576112f9610b8b565b5b813561130a8482602086016112a3565b91505092915050565b60006020828403121561132957611328610b81565b5b600082013567ffffffffffffffff81111561134757611346610b86565b5b611353848285016112e5565b91505092915050565b61136581610df5565b82525050565b6000602082019050611380600083018461135c565b92915050565b600081905092915050565b600061139c82610fe2565b6113a68185611386565b93506113b6818560208601610ffe565b80840191505092915050565b60006113ce8284611391565b915081905092915050565b600081905092915050565b60006113ef826110b0565b6113f981856113d9565b9350611409818560208601610ffe565b80840191505092915050565b600061142182846113e4565b915081905092915050565b60008160601b9050919050565b60006114448261142c565b9050919050565b600061145682611439565b9050919050565b61146e61146982610df5565b61144b565b82525050565b600061148082856113e4565b915061148c828461145d565b6014820191508190509392505050565b60006060820190506114b1600083018661135c565b6114be602083018561135c565b6114cb6040830184610d87565b949350505050565b6114dc8161123c565b81146114e757600080fd5b50565b6000815190506114f9816114d3565b92915050565b60006020828403121561151557611514610b81565b5b6000611523848285016114ea565b91505092915050565b600080fd5b600080fd5b600080fd5b600080833560016020038436030381126115585761155761152c565b5b80840192508235915067ffffffffffffffff82111561157a57611579611531565b5b60208301925060018202360383131561159657611595611536565b5b509250929050565b60006115aa8385611386565b93506115b7838584610c52565b82840190509392505050565b60006115d082848661159e565b91508190509392505050565b7f7265766572740000000000000000000000000000000000000000000000000000600082015250565b60006116126006836113d9565b915061161d826115dc565b600682019050919050565b600061163382611605565b9150819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061167782610d7d565b915061168283610d7d565b9250826116925761169161163d565b5b82820490509291505056fea2646970667358221220eabca6f5232e330b7bfacfcaa51280bd616323d8694afc461b86845313fcbed364736f6c634300081a0033",
}

// TestDAppV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDAppV2MetaData.ABI instead.
var TestDAppV2ABI = TestDAppV2MetaData.ABI

// TestDAppV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestDAppV2MetaData.Bin instead.
var TestDAppV2Bin = TestDAppV2MetaData.Bin

// DeployTestDAppV2 deploys a new Ethereum contract, binding an instance of TestDAppV2 to it.
func DeployTestDAppV2(auth *bind.TransactOpts, backend bind.ContractBackend, isZetaChain_ bool) (common.Address, *types.Transaction, *TestDAppV2, error) {
	parsed, err := TestDAppV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDAppV2Bin), backend, isZetaChain_)
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

// GetNoMessageIndex is a free data retrieval call binding the contract method 0xad23b28b.
//
// Solidity: function getNoMessageIndex(address sender) pure returns(string)
func (_TestDAppV2 *TestDAppV2Caller) GetNoMessageIndex(opts *bind.CallOpts, sender common.Address) (string, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "getNoMessageIndex", sender)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetNoMessageIndex is a free data retrieval call binding the contract method 0xad23b28b.
//
// Solidity: function getNoMessageIndex(address sender) pure returns(string)
func (_TestDAppV2 *TestDAppV2Session) GetNoMessageIndex(sender common.Address) (string, error) {
	return _TestDAppV2.Contract.GetNoMessageIndex(&_TestDAppV2.CallOpts, sender)
}

// GetNoMessageIndex is a free data retrieval call binding the contract method 0xad23b28b.
//
// Solidity: function getNoMessageIndex(address sender) pure returns(string)
func (_TestDAppV2 *TestDAppV2CallerSession) GetNoMessageIndex(sender common.Address) (string, error) {
	return _TestDAppV2.Contract.GetNoMessageIndex(&_TestDAppV2.CallOpts, sender)
}

// IsZetaChain is a free data retrieval call binding the contract method 0xc91f3567.
//
// Solidity: function isZetaChain() view returns(bool)
func (_TestDAppV2 *TestDAppV2Caller) IsZetaChain(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "isZetaChain")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsZetaChain is a free data retrieval call binding the contract method 0xc91f3567.
//
// Solidity: function isZetaChain() view returns(bool)
func (_TestDAppV2 *TestDAppV2Session) IsZetaChain() (bool, error) {
	return _TestDAppV2.Contract.IsZetaChain(&_TestDAppV2.CallOpts)
}

// IsZetaChain is a free data retrieval call binding the contract method 0xc91f3567.
//
// Solidity: function isZetaChain() view returns(bool)
func (_TestDAppV2 *TestDAppV2CallerSession) IsZetaChain() (bool, error) {
	return _TestDAppV2.Contract.IsZetaChain(&_TestDAppV2.CallOpts)
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
