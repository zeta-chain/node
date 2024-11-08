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

// MessageContext is an auto generated low-level Go binding around an user-defined struct.
type MessageContext struct {
	Sender common.Address
}

// RevertContext is an auto generated low-level Go binding around an user-defined struct.
type RevertContext struct {
	Sender        common.Address
	Asset         common.Address
	Amount        *big.Int
	RevertMessage []byte
}

// ZContext is an auto generated low-level Go binding around an user-defined struct.
type ZContext struct {
	Origin  []byte
	Sender  common.Address
	ChainID *big.Int
}

// TestDAppV2MetaData contains all meta data concerning the TestDAppV2 contract.
var TestDAppV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"NO_MESSAGE_CALL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAW\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"amountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"calledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"erc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"erc20Call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"gasCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getAmountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getCalledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getNoMessageIndex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structzContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structMessageContext\",\"name\":\"messageContext\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"}],\"internalType\":\"structRevertContext\",\"name\":\"revertContext\",\"type\":\"tuple\"}],\"name\":\"onRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"senderWithMessage\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"simpleCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052348015600f57600080fd5b50611dda8061001f6000396000f3fe6080604052600436106100e15760003560e01c8063ad23b28b1161007f578063c9028a3611610059578063c9028a36146102c1578063e2842ed7146102ea578063f592cbfb14610327578063f936ae8514610364576100e8565b8063ad23b28b14610230578063c7a339a91461026d578063c85f843414610296576100e8565b80635bcfd616116100bb5780635bcfd6161461017e578063676cc054146101a75780639291fe26146101d7578063a799911f14610214576100e8565b806316ba7197146100ed57806336e980a0146101185780634297a26314610141576100e8565b366100e857005b600080fd5b3480156100f957600080fd5b506101026103a1565b60405161010f9190610fda565b60405180910390f35b34801561012457600080fd5b5061013f600480360381019061013a9190611145565b6103da565b005b34801561014d57600080fd5b50610168600480360381019061016391906111c4565b610404565b604051610175919061120a565b60405180910390f35b34801561018a57600080fd5b506101a560048036038101906101a09190611333565b61041c565b005b6101c160048036038101906101bc91906113f6565b610862565b6040516101ce91906114ab565b60405180910390f35b3480156101e357600080fd5b506101fe60048036038101906101f99190611145565b610974565b60405161020b919061120a565b60405180910390f35b61022e60048036038101906102299190611145565b6109b7565b005b34801561023c57600080fd5b50610257600480360381019061025291906114cd565b6109e0565b6040516102649190610fda565b60405180910390f35b34801561027957600080fd5b50610294600480360381019061028f9190611538565b610a40565b005b3480156102a257600080fd5b506102ab610af4565b6040516102b89190610fda565b60405180910390f35b3480156102cd57600080fd5b506102e860048036038101906102e391906115c6565b610b2d565b005b3480156102f657600080fd5b50610311600480360381019061030c91906111c4565b610c67565b60405161031e919061162a565b60405180910390f35b34801561033357600080fd5b5061034e60048036038101906103499190611145565b610c87565b60405161035b919061162a565b60405180910390f35b34801561037057600080fd5b5061038b600480360381019061038691906116e6565b610cd7565b604051610398919061173e565b60405180910390f35b6040518060400160405280600881526020017f776974686472617700000000000000000000000000000000000000000000000081525081565b6103e381610d20565b156103ed57600080fd5b6103f681610d76565b610401816000610dca565b50565b60036020528060005260406000206000915090505481565b61046982828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610d20565b1561047357600080fd5b6104c082828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610e0c565b156107d2576000808573ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b81526004016040805180830381865afa158015610512573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105369190611783565b915091508573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16146105a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161059f9061180f565b60405180910390fd5b848111156105eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105e2906118a1565b60405180910390fd5b600081866105f991906118f0565b905061063a6040518060400160405280600781526020017f6761736c656674000000000000000000000000000000000000000000000000008152505a610dca565b610642610e99565b8673ffffffffffffffffffffffffffffffffffffffff1663095ea7b333886040518363ffffffff1660e01b815260040161067d929190611924565b6020604051808303816000875af115801561069c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106c09190611979565b503373ffffffffffffffffffffffffffffffffffffffff16637c0dcb5f8960200160208101906106f091906114cd565b604051602001610700919061173e565b604051602081830303815290604052838a6040518060a00160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518563ffffffff1660e01b81526004016107989493929190611a93565b600060405180830381600087803b1580156107b257600080fd5b505af11580156107c6573d6000803e3d6000fd5b5050505050505061085b565b60008083839050146108285782828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610844565b61084386602001602081019061083e91906114cd565b6109e0565b5b905061084f81610d76565b6108598185610dca565b505b5050505050565b606060008084849050146108ba5783838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506108d6565b6108d58560000160208101906108d091906114cd565b6109e0565b5b90506108e181610d76565b6108eb8134610dca565b8460000160208101906108fe91906114cd565b60028260405161090e9190611b22565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604051806020016040528060008152509150509392505050565b6000600360008360405160200161098b9190611b75565b604051602081830303815290604052805190602001208152602001908152602001600020549050919050565b6109c081610d20565b156109ca57600080fd5b6109d381610d76565b6109dd8134610dca565b50565b60606040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525082604051602001610a2a929190611bd4565b6040516020818303038152906040529050919050565b610a4981610d20565b15610a5357600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b8152600401610a9093929190611bfc565b6020604051808303816000875af1158015610aaf573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ad39190611979565b610adc57600080fd5b610ae581610d76565b610aef8183610dca565b505050565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525081565b610b88818060600190610b409190611c42565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610d76565b610be5818060600190610b9b9190611c42565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506000610dca565b806000016020810190610bf891906114cd565b6002828060600190610c0a9190611c42565b604051610c18929190611cca565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60016020528060005260406000206000915054906101000a900460ff1681565b60006001600083604051602001610c9e9190611b75565b60405160208183030381529060405280519060200120815260200190815260200160002060009054906101000a900460ff169050919050565b6002818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000604051602001610d3190611d2f565b6040516020818303038152906040528051906020012082604051602001610d589190611b75565b60405160208183030381529060405280519060200120149050919050565b600180600083604051602001610d8c9190611b75565b60405160208183030381529060405280519060200120815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b806003600084604051602001610de09190611b75565b604051602081830303815290604052805190602001208152602001908152602001600020819055505050565b60006040518060400160405280600881526020017f7769746864726177000000000000000000000000000000000000000000000000815250604051602001610e549190611b75565b6040516020818303038152906040528051906020012082604051602001610e7b9190611b75565b60405160208183030381529060405280519060200120149050919050565b6000621e848090506000614e20905060008183610eb69190611d73565b905060005b81811015610ef95760008190806001815401808255809150506001900390600052602060002001600090919091909150558080600101915050610ebb565b50600080610f079190610f0c565b505050565b5080546000825590600052602060002090810190610f2a9190610f2d565b50565b5b80821115610f46576000816000905550600101610f2e565b5090565b600081519050919050565b600082825260208201905092915050565b60005b83811015610f84578082015181840152602081019050610f69565b60008484015250505050565b6000601f19601f8301169050919050565b6000610fac82610f4a565b610fb68185610f55565b9350610fc6818560208601610f66565b610fcf81610f90565b840191505092915050565b60006020820190508181036000830152610ff48184610fa1565b905092915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61105282610f90565b810181811067ffffffffffffffff821117156110715761107061101a565b5b80604052505050565b6000611084610ffc565b90506110908282611049565b919050565b600067ffffffffffffffff8211156110b0576110af61101a565b5b6110b982610f90565b9050602081019050919050565b82818337600083830152505050565b60006110e86110e384611095565b61107a565b90508281526020810184848401111561110457611103611015565b5b61110f8482856110c6565b509392505050565b600082601f83011261112c5761112b611010565b5b813561113c8482602086016110d5565b91505092915050565b60006020828403121561115b5761115a611006565b5b600082013567ffffffffffffffff8111156111795761117861100b565b5b61118584828501611117565b91505092915050565b6000819050919050565b6111a18161118e565b81146111ac57600080fd5b50565b6000813590506111be81611198565b92915050565b6000602082840312156111da576111d9611006565b5b60006111e8848285016111af565b91505092915050565b6000819050919050565b611204816111f1565b82525050565b600060208201905061121f60008301846111fb565b92915050565b600080fd5b6000606082840312156112405761123f611225565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061127482611249565b9050919050565b61128481611269565b811461128f57600080fd5b50565b6000813590506112a18161127b565b92915050565b6112b0816111f1565b81146112bb57600080fd5b50565b6000813590506112cd816112a7565b92915050565b600080fd5b600080fd5b60008083601f8401126112f3576112f2611010565b5b8235905067ffffffffffffffff8111156113105761130f6112d3565b5b60208301915083600182028301111561132c5761132b6112d8565b5b9250929050565b60008060008060006080868803121561134f5761134e611006565b5b600086013567ffffffffffffffff81111561136d5761136c61100b565b5b6113798882890161122a565b955050602061138a88828901611292565b945050604061139b888289016112be565b935050606086013567ffffffffffffffff8111156113bc576113bb61100b565b5b6113c8888289016112dd565b92509250509295509295909350565b6000602082840312156113ed576113ec611225565b5b81905092915050565b60008060006040848603121561140f5761140e611006565b5b600061141d868287016113d7565b935050602084013567ffffffffffffffff81111561143e5761143d61100b565b5b61144a868287016112dd565b92509250509250925092565b600081519050919050565b600082825260208201905092915050565b600061147d82611456565b6114878185611461565b9350611497818560208601610f66565b6114a081610f90565b840191505092915050565b600060208201905081810360008301526114c58184611472565b905092915050565b6000602082840312156114e3576114e2611006565b5b60006114f184828501611292565b91505092915050565b600061150582611269565b9050919050565b611515816114fa565b811461152057600080fd5b50565b6000813590506115328161150c565b92915050565b60008060006060848603121561155157611550611006565b5b600061155f86828701611523565b9350506020611570868287016112be565b925050604084013567ffffffffffffffff8111156115915761159061100b565b5b61159d86828701611117565b9150509250925092565b6000608082840312156115bd576115bc611225565b5b81905092915050565b6000602082840312156115dc576115db611006565b5b600082013567ffffffffffffffff8111156115fa576115f961100b565b5b611606848285016115a7565b91505092915050565b60008115159050919050565b6116248161160f565b82525050565b600060208201905061163f600083018461161b565b92915050565b600067ffffffffffffffff8211156116605761165f61101a565b5b61166982610f90565b9050602081019050919050565b600061168961168484611645565b61107a565b9050828152602081018484840111156116a5576116a4611015565b5b6116b08482856110c6565b509392505050565b600082601f8301126116cd576116cc611010565b5b81356116dd848260208601611676565b91505092915050565b6000602082840312156116fc576116fb611006565b5b600082013567ffffffffffffffff81111561171a5761171961100b565b5b611726848285016116b8565b91505092915050565b61173881611269565b82525050565b6000602082019050611753600083018461172f565b92915050565b6000815190506117688161127b565b92915050565b60008151905061177d816112a7565b92915050565b6000806040838503121561179a57611799611006565b5b60006117a885828601611759565b92505060206117b98582860161176e565b9150509250929050565b7f7a72633230206973206e6f742067617320746f6b656e00000000000000000000600082015250565b60006117f9601683610f55565b9150611804826117c3565b602082019050919050565b60006020820190508181036000830152611828816117ec565b9050919050565b7f66656520616d6f756e7420697320686967686572207468616e2074686520616d60008201527f6f756e7400000000000000000000000000000000000000000000000000000000602082015250565b600061188b602483610f55565b91506118968261182f565b604082019050919050565b600060208201905081810360008301526118ba8161187e565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006118fb826111f1565b9150611906836111f1565b925082820390508181111561191e5761191d6118c1565b5b92915050565b6000604082019050611939600083018561172f565b61194660208301846111fb565b9392505050565b6119568161160f565b811461196157600080fd5b50565b6000815190506119738161194d565b92915050565b60006020828403121561198f5761198e611006565b5b600061199d84828501611964565b91505092915050565b6119af81611269565b82525050565b6119be8161160f565b82525050565b600082825260208201905092915050565b60006119e082611456565b6119ea81856119c4565b93506119fa818560208601610f66565b611a0381610f90565b840191505092915050565b611a17816111f1565b82525050565b600060a083016000830151611a3560008601826119a6565b506020830151611a4860208601826119b5565b506040830151611a5b60408601826119a6565b5060608301518482036060860152611a7382826119d5565b9150506080830151611a886080860182611a0e565b508091505092915050565b60006080820190508181036000830152611aad8187611472565b9050611abc60208301866111fb565b611ac9604083018561172f565b8181036060830152611adb8184611a1d565b905095945050505050565b600081905092915050565b6000611afc82611456565b611b068185611ae6565b9350611b16818560208601610f66565b80840191505092915050565b6000611b2e8284611af1565b915081905092915050565b600081905092915050565b6000611b4f82610f4a565b611b598185611b39565b9350611b69818560208601610f66565b80840191505092915050565b6000611b818284611b44565b915081905092915050565b60008160601b9050919050565b6000611ba482611b8c565b9050919050565b6000611bb682611b99565b9050919050565b611bce611bc982611269565b611bab565b82525050565b6000611be08285611b44565b9150611bec8284611bbd565b6014820191508190509392505050565b6000606082019050611c11600083018661172f565b611c1e602083018561172f565b611c2b60408301846111fb565b949350505050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112611c5f57611c5e611c33565b5b80840192508235915067ffffffffffffffff821115611c8157611c80611c38565b5b602083019250600182023603831315611c9d57611c9c611c3d565b5b509250929050565b6000611cb18385611ae6565b9350611cbe8385846110c6565b82840190509392505050565b6000611cd7828486611ca5565b91508190509392505050565b7f7265766572740000000000000000000000000000000000000000000000000000600082015250565b6000611d19600683611b39565b9150611d2482611ce3565b600682019050919050565b6000611d3a82611d0c565b9150819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000611d7e826111f1565b9150611d89836111f1565b925082611d9957611d98611d44565b5b82820490509291505056fea2646970667358221220e3fea42a51c464a130be48d80080d5d2c30f8ad7b31ed185435e1071da4ab72c64736f6c634300081a0033",
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

// WITHDRAW is a free data retrieval call binding the contract method 0x16ba7197.
//
// Solidity: function WITHDRAW() view returns(string)
func (_TestDAppV2 *TestDAppV2Caller) WITHDRAW(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "WITHDRAW")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// WITHDRAW is a free data retrieval call binding the contract method 0x16ba7197.
//
// Solidity: function WITHDRAW() view returns(string)
func (_TestDAppV2 *TestDAppV2Session) WITHDRAW() (string, error) {
	return _TestDAppV2.Contract.WITHDRAW(&_TestDAppV2.CallOpts)
}

// WITHDRAW is a free data retrieval call binding the contract method 0x16ba7197.
//
// Solidity: function WITHDRAW() view returns(string)
func (_TestDAppV2 *TestDAppV2CallerSession) WITHDRAW() (string, error) {
	return _TestDAppV2.Contract.WITHDRAW(&_TestDAppV2.CallOpts)
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
// Solidity: function onCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2Transactor) OnCall(opts *bind.TransactOpts, context ZContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "onCall", context, zrc20, amount, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2Session) OnCall(context ZContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCall(&_TestDAppV2.TransactOpts, context, zrc20, amount, message)
}

// OnCall is a paid mutator transaction binding the contract method 0x5bcfd616.
//
// Solidity: function onCall((bytes,address,uint256) context, address zrc20, uint256 amount, bytes message) returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) OnCall(context ZContext, zrc20 common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCall(&_TestDAppV2.TransactOpts, context, zrc20, amount, message)
}

// OnCall0 is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppV2 *TestDAppV2Transactor) OnCall0(opts *bind.TransactOpts, messageContext MessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "onCall0", messageContext, message)
}

// OnCall0 is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppV2 *TestDAppV2Session) OnCall0(messageContext MessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCall0(&_TestDAppV2.TransactOpts, messageContext, message)
}

// OnCall0 is a paid mutator transaction binding the contract method 0x676cc054.
//
// Solidity: function onCall((address) messageContext, bytes message) payable returns(bytes)
func (_TestDAppV2 *TestDAppV2TransactorSession) OnCall0(messageContext MessageContext, message []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnCall0(&_TestDAppV2.TransactOpts, messageContext, message)
}

// OnRevert is a paid mutator transaction binding the contract method 0xc9028a36.
//
// Solidity: function onRevert((address,address,uint256,bytes) revertContext) returns()
func (_TestDAppV2 *TestDAppV2Transactor) OnRevert(opts *bind.TransactOpts, revertContext RevertContext) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "onRevert", revertContext)
}

// OnRevert is a paid mutator transaction binding the contract method 0xc9028a36.
//
// Solidity: function onRevert((address,address,uint256,bytes) revertContext) returns()
func (_TestDAppV2 *TestDAppV2Session) OnRevert(revertContext RevertContext) (*types.Transaction, error) {
	return _TestDAppV2.Contract.OnRevert(&_TestDAppV2.TransactOpts, revertContext)
}

// OnRevert is a paid mutator transaction binding the contract method 0xc9028a36.
//
// Solidity: function onRevert((address,address,uint256,bytes) revertContext) returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) OnRevert(revertContext RevertContext) (*types.Transaction, error) {
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
