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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"HelloEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"NO_MESSAGE_CALL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITHDRAW\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"amountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"calledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"erc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"erc20Call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"gasCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getAmountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getCalledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getNoMessageIndex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structzContext\",\"name\":\"context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structMessageContext\",\"name\":\"messageContext\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"}],\"internalType\":\"structRevertContext\",\"name\":\"revertContext\",\"type\":\"tuple\"}],\"name\":\"onRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"senderWithMessage\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"simpleCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052348015600f57600080fd5b50611eda8061001f6000396000f3fe6080604052600436106100e15760003560e01c8063ad23b28b1161007f578063c9028a3611610059578063c9028a36146102c1578063e2842ed7146102ea578063f592cbfb14610327578063f936ae8514610364576100e8565b8063ad23b28b14610230578063c7a339a91461026d578063c85f843414610296576100e8565b80635bcfd616116100bb5780635bcfd6161461017e578063676cc054146101a75780639291fe26146101d7578063a799911f14610214576100e8565b806316ba7197146100ed57806336e980a0146101185780634297a26314610141576100e8565b366100e857005b600080fd5b3480156100f957600080fd5b506101026103a1565b60405161010f919061100f565b60405180910390f35b34801561012457600080fd5b5061013f600480360381019061013a919061117a565b6103da565b005b34801561014d57600080fd5b50610168600480360381019061016391906111f9565b610404565b604051610175919061123f565b60405180910390f35b34801561018a57600080fd5b506101a560048036038101906101a09190611368565b61041c565b005b6101c160048036038101906101bc919061142b565b610897565b6040516101ce91906114e0565b60405180910390f35b3480156101e357600080fd5b506101fe60048036038101906101f9919061117a565b6109a9565b60405161020b919061123f565b60405180910390f35b61022e6004803603810190610229919061117a565b6109ec565b005b34801561023c57600080fd5b5061025760048036038101906102529190611502565b610a15565b604051610264919061100f565b60405180910390f35b34801561027957600080fd5b50610294600480360381019061028f919061156d565b610a75565b005b3480156102a257600080fd5b506102ab610b29565b6040516102b8919061100f565b60405180910390f35b3480156102cd57600080fd5b506102e860048036038101906102e391906115fb565b610b62565b005b3480156102f657600080fd5b50610311600480360381019061030c91906111f9565b610c9c565b60405161031e919061165f565b60405180910390f35b34801561033357600080fd5b5061034e6004803603810190610349919061117a565b610cbc565b60405161035b919061165f565b60405180910390f35b34801561037057600080fd5b5061038b6004803603810190610386919061171b565b610d0c565b6040516103989190611773565b60405180910390f35b6040518060400160405280600881526020017f776974686472617700000000000000000000000000000000000000000000000081525081565b6103e381610d55565b156103ed57600080fd5b6103f681610dab565b610401816000610dff565b50565b60036020528060005260406000206000915090505481565b61046982828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610d55565b1561047357600080fd5b6104c082828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610e41565b15610807576000808573ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b81526004016040805180830381865afa158015610512573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061053691906117b8565b915091508573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16146105a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161059f90611844565b60405180910390fd5b848111156105eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105e2906118d6565b60405180910390fd5b600081866105f99190611925565b905061063a6040518060400160405280600781526020017f6761736c656674000000000000000000000000000000000000000000000000008152505a610dff565b610642610ece565b8673ffffffffffffffffffffffffffffffffffffffff1663095ea7b333886040518363ffffffff1660e01b815260040161067d929190611959565b6020604051808303816000875af115801561069c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106c091906119ae565b503373ffffffffffffffffffffffffffffffffffffffff16637c0dcb5f8960200160208101906106f09190611502565b6040516020016107009190611773565b604051602081830303815290604052838a6040518060a00160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518563ffffffff1660e01b81526004016107989493929190611ac8565b600060405180830381600087803b1580156107b257600080fd5b505af11580156107c6573d6000803e3d6000fd5b505050507f39f8c79736fed93bca390bb3d6ff7da07482edb61cd7dafcfba496821d6ab7a36040516107f790611bb3565b60405180910390a1505050610890565b600080838390501461085d5782828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610879565b6108788660200160208101906108739190611502565b610a15565b5b905061088481610dab565b61088e8185610dff565b505b5050505050565b606060008084849050146108ef5783838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505061090b565b61090a8560000160208101906109059190611502565b610a15565b5b905061091681610dab565b6109208134610dff565b8460000160208101906109339190611502565b6002826040516109439190611c22565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604051806020016040528060008152509150509392505050565b600060036000836040516020016109c09190611c75565b604051602081830303815290604052805190602001208152602001908152602001600020549050919050565b6109f581610d55565b156109ff57600080fd5b610a0881610dab565b610a128134610dff565b50565b60606040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525082604051602001610a5f929190611cd4565b6040516020818303038152906040529050919050565b610a7e81610d55565b15610a8857600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b8152600401610ac593929190611cfc565b6020604051808303816000875af1158015610ae4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b0891906119ae565b610b1157600080fd5b610b1a81610dab565b610b248183610dff565b505050565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525081565b610bbd818060600190610b759190611d42565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610dab565b610c1a818060600190610bd09190611d42565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506000610dff565b806000016020810190610c2d9190611502565b6002828060600190610c3f9190611d42565b604051610c4d929190611dca565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60016020528060005260406000206000915054906101000a900460ff1681565b60006001600083604051602001610cd39190611c75565b60405160208183030381529060405280519060200120815260200190815260200160002060009054906101000a900460ff169050919050565b6002818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000604051602001610d6690611e2f565b6040516020818303038152906040528051906020012082604051602001610d8d9190611c75565b60405160208183030381529060405280519060200120149050919050565b600180600083604051602001610dc19190611c75565b60405160208183030381529060405280519060200120815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b806003600084604051602001610e159190611c75565b604051602081830303815290604052805190602001208152602001908152602001600020819055505050565b60006040518060400160405280600881526020017f7769746864726177000000000000000000000000000000000000000000000000815250604051602001610e899190611c75565b6040516020818303038152906040528051906020012082604051602001610eb09190611c75565b60405160208183030381529060405280519060200120149050919050565b6000620493e090506000614e20905060008183610eeb9190611e73565b905060005b81811015610f2e5760008190806001815401808255809150506001900390600052602060002001600090919091909150558080600101915050610ef0565b50600080610f3c9190610f41565b505050565b5080546000825590600052602060002090810190610f5f9190610f62565b50565b5b80821115610f7b576000816000905550600101610f63565b5090565b600081519050919050565b600082825260208201905092915050565b60005b83811015610fb9578082015181840152602081019050610f9e565b60008484015250505050565b6000601f19601f8301169050919050565b6000610fe182610f7f565b610feb8185610f8a565b9350610ffb818560208601610f9b565b61100481610fc5565b840191505092915050565b600060208201905081810360008301526110298184610fd6565b905092915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61108782610fc5565b810181811067ffffffffffffffff821117156110a6576110a561104f565b5b80604052505050565b60006110b9611031565b90506110c5828261107e565b919050565b600067ffffffffffffffff8211156110e5576110e461104f565b5b6110ee82610fc5565b9050602081019050919050565b82818337600083830152505050565b600061111d611118846110ca565b6110af565b9050828152602081018484840111156111395761113861104a565b5b6111448482856110fb565b509392505050565b600082601f83011261116157611160611045565b5b813561117184826020860161110a565b91505092915050565b6000602082840312156111905761118f61103b565b5b600082013567ffffffffffffffff8111156111ae576111ad611040565b5b6111ba8482850161114c565b91505092915050565b6000819050919050565b6111d6816111c3565b81146111e157600080fd5b50565b6000813590506111f3816111cd565b92915050565b60006020828403121561120f5761120e61103b565b5b600061121d848285016111e4565b91505092915050565b6000819050919050565b61123981611226565b82525050565b60006020820190506112546000830184611230565b92915050565b600080fd5b6000606082840312156112755761127461125a565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006112a98261127e565b9050919050565b6112b98161129e565b81146112c457600080fd5b50565b6000813590506112d6816112b0565b92915050565b6112e581611226565b81146112f057600080fd5b50565b600081359050611302816112dc565b92915050565b600080fd5b600080fd5b60008083601f84011261132857611327611045565b5b8235905067ffffffffffffffff81111561134557611344611308565b5b6020830191508360018202830111156113615761136061130d565b5b9250929050565b6000806000806000608086880312156113845761138361103b565b5b600086013567ffffffffffffffff8111156113a2576113a1611040565b5b6113ae8882890161125f565b95505060206113bf888289016112c7565b94505060406113d0888289016112f3565b935050606086013567ffffffffffffffff8111156113f1576113f0611040565b5b6113fd88828901611312565b92509250509295509295909350565b6000602082840312156114225761142161125a565b5b81905092915050565b6000806000604084860312156114445761144361103b565b5b60006114528682870161140c565b935050602084013567ffffffffffffffff81111561147357611472611040565b5b61147f86828701611312565b92509250509250925092565b600081519050919050565b600082825260208201905092915050565b60006114b28261148b565b6114bc8185611496565b93506114cc818560208601610f9b565b6114d581610fc5565b840191505092915050565b600060208201905081810360008301526114fa81846114a7565b905092915050565b6000602082840312156115185761151761103b565b5b6000611526848285016112c7565b91505092915050565b600061153a8261129e565b9050919050565b61154a8161152f565b811461155557600080fd5b50565b60008135905061156781611541565b92915050565b6000806000606084860312156115865761158561103b565b5b600061159486828701611558565b93505060206115a5868287016112f3565b925050604084013567ffffffffffffffff8111156115c6576115c5611040565b5b6115d28682870161114c565b9150509250925092565b6000608082840312156115f2576115f161125a565b5b81905092915050565b6000602082840312156116115761161061103b565b5b600082013567ffffffffffffffff81111561162f5761162e611040565b5b61163b848285016115dc565b91505092915050565b60008115159050919050565b61165981611644565b82525050565b60006020820190506116746000830184611650565b92915050565b600067ffffffffffffffff8211156116955761169461104f565b5b61169e82610fc5565b9050602081019050919050565b60006116be6116b98461167a565b6110af565b9050828152602081018484840111156116da576116d961104a565b5b6116e58482856110fb565b509392505050565b600082601f83011261170257611701611045565b5b81356117128482602086016116ab565b91505092915050565b6000602082840312156117315761173061103b565b5b600082013567ffffffffffffffff81111561174f5761174e611040565b5b61175b848285016116ed565b91505092915050565b61176d8161129e565b82525050565b60006020820190506117886000830184611764565b92915050565b60008151905061179d816112b0565b92915050565b6000815190506117b2816112dc565b92915050565b600080604083850312156117cf576117ce61103b565b5b60006117dd8582860161178e565b92505060206117ee858286016117a3565b9150509250929050565b7f7a72633230206973206e6f742067617320746f6b656e00000000000000000000600082015250565b600061182e601683610f8a565b9150611839826117f8565b602082019050919050565b6000602082019050818103600083015261185d81611821565b9050919050565b7f66656520616d6f756e7420697320686967686572207468616e2074686520616d60008201527f6f756e7400000000000000000000000000000000000000000000000000000000602082015250565b60006118c0602483610f8a565b91506118cb82611864565b604082019050919050565b600060208201905081810360008301526118ef816118b3565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061193082611226565b915061193b83611226565b9250828203905081811115611953576119526118f6565b5b92915050565b600060408201905061196e6000830185611764565b61197b6020830184611230565b9392505050565b61198b81611644565b811461199657600080fd5b50565b6000815190506119a881611982565b92915050565b6000602082840312156119c4576119c361103b565b5b60006119d284828501611999565b91505092915050565b6119e48161129e565b82525050565b6119f381611644565b82525050565b600082825260208201905092915050565b6000611a158261148b565b611a1f81856119f9565b9350611a2f818560208601610f9b565b611a3881610fc5565b840191505092915050565b611a4c81611226565b82525050565b600060a083016000830151611a6a60008601826119db565b506020830151611a7d60208601826119ea565b506040830151611a9060408601826119db565b5060608301518482036060860152611aa88282611a0a565b9150506080830151611abd6080860182611a43565b508091505092915050565b60006080820190508181036000830152611ae281876114a7565b9050611af16020830186611230565b611afe6040830185611764565b8181036060830152611b108184611a52565b905095945050505050565b7f48656c6c6f000000000000000000000000000000000000000000000000000000600082015250565b6000611b51600583610f8a565b9150611b5c82611b1b565b602082019050919050565b7f576f726c64000000000000000000000000000000000000000000000000000000600082015250565b6000611b9d600583610f8a565b9150611ba882611b67565b602082019050919050565b60006040820190508181036000830152611bcc81611b44565b90508181036020830152611bdf81611b90565b9050919050565b600081905092915050565b6000611bfc8261148b565b611c068185611be6565b9350611c16818560208601610f9b565b80840191505092915050565b6000611c2e8284611bf1565b915081905092915050565b600081905092915050565b6000611c4f82610f7f565b611c598185611c39565b9350611c69818560208601610f9b565b80840191505092915050565b6000611c818284611c44565b915081905092915050565b60008160601b9050919050565b6000611ca482611c8c565b9050919050565b6000611cb682611c99565b9050919050565b611cce611cc98261129e565b611cab565b82525050565b6000611ce08285611c44565b9150611cec8284611cbd565b6014820191508190509392505050565b6000606082019050611d116000830186611764565b611d1e6020830185611764565b611d2b6040830184611230565b949350505050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112611d5f57611d5e611d33565b5b80840192508235915067ffffffffffffffff821115611d8157611d80611d38565b5b602083019250600182023603831315611d9d57611d9c611d3d565b5b509250929050565b6000611db18385611be6565b9350611dbe8385846110fb565b82840190509392505050565b6000611dd7828486611da5565b91508190509392505050565b7f7265766572740000000000000000000000000000000000000000000000000000600082015250565b6000611e19600683611c39565b9150611e2482611de3565b600682019050919050565b6000611e3a82611e0c565b9150819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000611e7e82611226565b9150611e8983611226565b925082611e9957611e98611e44565b5b82820490509291505056fea26469706673582212205b8ad2702dcba66c08992683150d98f82b204b950f48efd2b79f1ac03e90f2cd64736f6c634300081a0033",
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

// TestDAppV2HelloEventIterator is returned from FilterHelloEvent and is used to iterate over the raw logs and unpacked data for HelloEvent events raised by the TestDAppV2 contract.
type TestDAppV2HelloEventIterator struct {
	Event *TestDAppV2HelloEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestDAppV2HelloEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestDAppV2HelloEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestDAppV2HelloEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestDAppV2HelloEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestDAppV2HelloEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestDAppV2HelloEvent represents a HelloEvent event raised by the TestDAppV2 contract.
type TestDAppV2HelloEvent struct {
	Arg0 string
	Arg1 string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterHelloEvent is a free log retrieval operation binding the contract event 0x39f8c79736fed93bca390bb3d6ff7da07482edb61cd7dafcfba496821d6ab7a3.
//
// Solidity: event HelloEvent(string arg0, string arg1)
func (_TestDAppV2 *TestDAppV2Filterer) FilterHelloEvent(opts *bind.FilterOpts) (*TestDAppV2HelloEventIterator, error) {

	logs, sub, err := _TestDAppV2.contract.FilterLogs(opts, "HelloEvent")
	if err != nil {
		return nil, err
	}
	return &TestDAppV2HelloEventIterator{contract: _TestDAppV2.contract, event: "HelloEvent", logs: logs, sub: sub}, nil
}

// WatchHelloEvent is a free log subscription operation binding the contract event 0x39f8c79736fed93bca390bb3d6ff7da07482edb61cd7dafcfba496821d6ab7a3.
//
// Solidity: event HelloEvent(string arg0, string arg1)
func (_TestDAppV2 *TestDAppV2Filterer) WatchHelloEvent(opts *bind.WatchOpts, sink chan<- *TestDAppV2HelloEvent) (event.Subscription, error) {

	logs, sub, err := _TestDAppV2.contract.WatchLogs(opts, "HelloEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestDAppV2HelloEvent)
				if err := _TestDAppV2.contract.UnpackLog(event, "HelloEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseHelloEvent is a log parse operation binding the contract event 0x39f8c79736fed93bca390bb3d6ff7da07482edb61cd7dafcfba496821d6ab7a3.
//
// Solidity: event HelloEvent(string arg0, string arg1)
func (_TestDAppV2 *TestDAppV2Filterer) ParseHelloEvent(log types.Log) (*TestDAppV2HelloEvent, error) {
	event := new(TestDAppV2HelloEvent)
	if err := _TestDAppV2.contract.UnpackLog(event, "HelloEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
