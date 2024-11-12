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
	Bin: "0x6080604052348015600f57600080fd5b506120f28061001f6000396000f3fe6080604052600436106100e15760003560e01c8063ad23b28b1161007f578063c9028a3611610059578063c9028a36146102c1578063e2842ed7146102ea578063f592cbfb14610327578063f936ae8514610364576100e8565b8063ad23b28b14610230578063c7a339a91461026d578063c85f843414610296576100e8565b80635bcfd616116100bb5780635bcfd6161461017e578063676cc054146101a75780639291fe26146101d7578063a799911f14610214576100e8565b806316ba7197146100ed57806336e980a0146101185780634297a26314610141576100e8565b366100e857005b600080fd5b3480156100f957600080fd5b506101026103a1565b60405161010f919061114b565b60405180910390f35b34801561012457600080fd5b5061013f600480360381019061013a91906112b6565b6103da565b005b34801561014d57600080fd5b5061016860048036038101906101639190611335565b610404565b604051610175919061137b565b60405180910390f35b34801561018a57600080fd5b506101a560048036038101906101a091906114a4565b61041c565b005b6101c160048036038101906101bc9190611567565b6109d3565b6040516101ce919061161c565b60405180910390f35b3480156101e357600080fd5b506101fe60048036038101906101f991906112b6565b610ae5565b60405161020b919061137b565b60405180910390f35b61022e600480360381019061022991906112b6565b610b28565b005b34801561023c57600080fd5b506102576004803603810190610252919061163e565b610b51565b604051610264919061114b565b60405180910390f35b34801561027957600080fd5b50610294600480360381019061028f91906116a9565b610bb1565b005b3480156102a257600080fd5b506102ab610c65565b6040516102b8919061114b565b60405180910390f35b3480156102cd57600080fd5b506102e860048036038101906102e39190611737565b610c9e565b005b3480156102f657600080fd5b50610311600480360381019061030c9190611335565b610dd8565b60405161031e919061179b565b60405180910390f35b34801561033357600080fd5b5061034e600480360381019061034991906112b6565b610df8565b60405161035b919061179b565b60405180910390f35b34801561037057600080fd5b5061038b60048036038101906103869190611857565b610e48565b60405161039891906118af565b60405180910390f35b6040518060400160405280600881526020017f776974686472617700000000000000000000000000000000000000000000000081525081565b6103e381610e91565b156103ed57600080fd5b6103f681610ee7565b610401816000610f3b565b50565b60036020528060005260406000206000915090505481565b61046982828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610e91565b1561047357600080fd5b6104c082828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610f7d565b15610943576000808573ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b81526004016040805180830381865afa158015610512573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061053691906118f4565b915091508573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16146105a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161059f90611980565b60405180910390fd5b848111156105eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105e290611a12565b60405180910390fd5b600081866105f99190611a61565b905061063a6040518060400160405280600781526020017f6761736c656674000000000000000000000000000000000000000000000000008152505a610f3b565b61064261100a565b8673ffffffffffffffffffffffffffffffffffffffff1663095ea7b333886040518363ffffffff1660e01b815260040161067d929190611a95565b6020604051808303816000875af115801561069c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106c09190611aea565b503373ffffffffffffffffffffffffffffffffffffffff16637c0dcb5f8960200160208101906106f0919061163e565b60405160200161070091906118af565b604051602081830303815290604052838a6040518060a00160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518563ffffffff1660e01b81526004016107989493929190611c04565b600060405180830381600087803b1580156107b257600080fd5b505af11580156107c6573d6000803e3d6000fd5b5050505060006040518060a001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020016000151581526020013373ffffffffffffffffffffffffffffffffffffffff1681526020016040518060400160405280600e81526020017f726576657274206d6573736167650000000000000000000000000000000000008152508152602001620186a081525090503373ffffffffffffffffffffffffffffffffffffffff16631cb5ea758a602001602081019061088b919061163e565b60405160200161089b91906118af565b6040516020818303038152906040528a8989620186a0876040518763ffffffff1660e01b81526004016108d396959493929190611cc9565b600060405180830381600087803b1580156108ed57600080fd5b505af1158015610901573d6000803e3d6000fd5b505050507f39f8c79736fed93bca390bb3d6ff7da07482edb61cd7dafcfba496821d6ab7a360405161093290611dcb565b60405180910390a1505050506109cc565b60008083839050146109995782828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506109b5565b6109b48660200160208101906109af919061163e565b610b51565b5b90506109c081610ee7565b6109ca8185610f3b565b505b5050505050565b60606000808484905014610a2b5783838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610a47565b610a46856000016020810190610a41919061163e565b610b51565b5b9050610a5281610ee7565b610a5c8134610f3b565b846000016020810190610a6f919061163e565b600282604051610a7f9190611e3a565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604051806020016040528060008152509150509392505050565b60006003600083604051602001610afc9190611e8d565b604051602081830303815290604052805190602001208152602001908152602001600020549050919050565b610b3181610e91565b15610b3b57600080fd5b610b4481610ee7565b610b4e8134610f3b565b50565b60606040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525082604051602001610b9b929190611eec565b6040516020818303038152906040529050919050565b610bba81610e91565b15610bc457600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b8152600401610c0193929190611f14565b6020604051808303816000875af1158015610c20573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c449190611aea565b610c4d57600080fd5b610c5681610ee7565b610c608183610f3b565b505050565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525081565b610cf9818060600190610cb19190611f5a565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610ee7565b610d56818060600190610d0c9190611f5a565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506000610f3b565b806000016020810190610d69919061163e565b6002828060600190610d7b9190611f5a565b604051610d89929190611fe2565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60016020528060005260406000206000915054906101000a900460ff1681565b60006001600083604051602001610e0f9190611e8d565b60405160208183030381529060405280519060200120815260200190815260200160002060009054906101000a900460ff169050919050565b6002818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000604051602001610ea290612047565b6040516020818303038152906040528051906020012082604051602001610ec99190611e8d565b60405160208183030381529060405280519060200120149050919050565b600180600083604051602001610efd9190611e8d565b60405160208183030381529060405280519060200120815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b806003600084604051602001610f519190611e8d565b604051602081830303815290604052805190602001208152602001908152602001600020819055505050565b60006040518060400160405280600881526020017f7769746864726177000000000000000000000000000000000000000000000000815250604051602001610fc59190611e8d565b6040516020818303038152906040528051906020012082604051602001610fec9190611e8d565b60405160208183030381529060405280519060200120149050919050565b6000620493e090506000614e20905060008183611027919061208b565b905060005b8181101561106a576000819080600181540180825580915050600190039060005260206000200160009091909190915055808060010191505061102c565b50600080611078919061107d565b505050565b508054600082559060005260206000209081019061109b919061109e565b50565b5b808211156110b757600081600090555060010161109f565b5090565b600081519050919050565b600082825260208201905092915050565b60005b838110156110f55780820151818401526020810190506110da565b60008484015250505050565b6000601f19601f8301169050919050565b600061111d826110bb565b61112781856110c6565b93506111378185602086016110d7565b61114081611101565b840191505092915050565b600060208201905081810360008301526111658184611112565b905092915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6111c382611101565b810181811067ffffffffffffffff821117156111e2576111e161118b565b5b80604052505050565b60006111f561116d565b905061120182826111ba565b919050565b600067ffffffffffffffff8211156112215761122061118b565b5b61122a82611101565b9050602081019050919050565b82818337600083830152505050565b600061125961125484611206565b6111eb565b90508281526020810184848401111561127557611274611186565b5b611280848285611237565b509392505050565b600082601f83011261129d5761129c611181565b5b81356112ad848260208601611246565b91505092915050565b6000602082840312156112cc576112cb611177565b5b600082013567ffffffffffffffff8111156112ea576112e961117c565b5b6112f684828501611288565b91505092915050565b6000819050919050565b611312816112ff565b811461131d57600080fd5b50565b60008135905061132f81611309565b92915050565b60006020828403121561134b5761134a611177565b5b600061135984828501611320565b91505092915050565b6000819050919050565b61137581611362565b82525050565b6000602082019050611390600083018461136c565b92915050565b600080fd5b6000606082840312156113b1576113b0611396565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006113e5826113ba565b9050919050565b6113f5816113da565b811461140057600080fd5b50565b600081359050611412816113ec565b92915050565b61142181611362565b811461142c57600080fd5b50565b60008135905061143e81611418565b92915050565b600080fd5b600080fd5b60008083601f84011261146457611463611181565b5b8235905067ffffffffffffffff81111561148157611480611444565b5b60208301915083600182028301111561149d5761149c611449565b5b9250929050565b6000806000806000608086880312156114c0576114bf611177565b5b600086013567ffffffffffffffff8111156114de576114dd61117c565b5b6114ea8882890161139b565b95505060206114fb88828901611403565b945050604061150c8882890161142f565b935050606086013567ffffffffffffffff81111561152d5761152c61117c565b5b6115398882890161144e565b92509250509295509295909350565b60006020828403121561155e5761155d611396565b5b81905092915050565b6000806000604084860312156115805761157f611177565b5b600061158e86828701611548565b935050602084013567ffffffffffffffff8111156115af576115ae61117c565b5b6115bb8682870161144e565b92509250509250925092565b600081519050919050565b600082825260208201905092915050565b60006115ee826115c7565b6115f881856115d2565b93506116088185602086016110d7565b61161181611101565b840191505092915050565b6000602082019050818103600083015261163681846115e3565b905092915050565b60006020828403121561165457611653611177565b5b600061166284828501611403565b91505092915050565b6000611676826113da565b9050919050565b6116868161166b565b811461169157600080fd5b50565b6000813590506116a38161167d565b92915050565b6000806000606084860312156116c2576116c1611177565b5b60006116d086828701611694565b93505060206116e18682870161142f565b925050604084013567ffffffffffffffff8111156117025761170161117c565b5b61170e86828701611288565b9150509250925092565b60006080828403121561172e5761172d611396565b5b81905092915050565b60006020828403121561174d5761174c611177565b5b600082013567ffffffffffffffff81111561176b5761176a61117c565b5b61177784828501611718565b91505092915050565b60008115159050919050565b61179581611780565b82525050565b60006020820190506117b0600083018461178c565b92915050565b600067ffffffffffffffff8211156117d1576117d061118b565b5b6117da82611101565b9050602081019050919050565b60006117fa6117f5846117b6565b6111eb565b90508281526020810184848401111561181657611815611186565b5b611821848285611237565b509392505050565b600082601f83011261183e5761183d611181565b5b813561184e8482602086016117e7565b91505092915050565b60006020828403121561186d5761186c611177565b5b600082013567ffffffffffffffff81111561188b5761188a61117c565b5b61189784828501611829565b91505092915050565b6118a9816113da565b82525050565b60006020820190506118c460008301846118a0565b92915050565b6000815190506118d9816113ec565b92915050565b6000815190506118ee81611418565b92915050565b6000806040838503121561190b5761190a611177565b5b6000611919858286016118ca565b925050602061192a858286016118df565b9150509250929050565b7f7a72633230206973206e6f742067617320746f6b656e00000000000000000000600082015250565b600061196a6016836110c6565b915061197582611934565b602082019050919050565b600060208201905081810360008301526119998161195d565b9050919050565b7f66656520616d6f756e7420697320686967686572207468616e2074686520616d60008201527f6f756e7400000000000000000000000000000000000000000000000000000000602082015250565b60006119fc6024836110c6565b9150611a07826119a0565b604082019050919050565b60006020820190508181036000830152611a2b816119ef565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611a6c82611362565b9150611a7783611362565b9250828203905081811115611a8f57611a8e611a32565b5b92915050565b6000604082019050611aaa60008301856118a0565b611ab7602083018461136c565b9392505050565b611ac781611780565b8114611ad257600080fd5b50565b600081519050611ae481611abe565b92915050565b600060208284031215611b0057611aff611177565b5b6000611b0e84828501611ad5565b91505092915050565b611b20816113da565b82525050565b611b2f81611780565b82525050565b600082825260208201905092915050565b6000611b51826115c7565b611b5b8185611b35565b9350611b6b8185602086016110d7565b611b7481611101565b840191505092915050565b611b8881611362565b82525050565b600060a083016000830151611ba66000860182611b17565b506020830151611bb96020860182611b26565b506040830151611bcc6040860182611b17565b5060608301518482036060860152611be48282611b46565b9150506080830151611bf96080860182611b7f565b508091505092915050565b60006080820190508181036000830152611c1e81876115e3565b9050611c2d602083018661136c565b611c3a60408301856118a0565b8181036060830152611c4c8184611b8e565b905095945050505050565b6000611c6383856115d2565b9350611c70838584611237565b611c7983611101565b840190509392505050565b6000819050919050565b6000819050919050565b6000611cb3611cae611ca984611c84565b611c8e565b611362565b9050919050565b611cc381611c98565b82525050565b600060a0820190508181036000830152611ce381896115e3565b9050611cf260208301886118a0565b8181036040830152611d05818688611c57565b9050611d146060830185611cba565b8181036080830152611d268184611b8e565b9050979650505050505050565b7f48656c6c6f000000000000000000000000000000000000000000000000000000600082015250565b6000611d696005836110c6565b9150611d7482611d33565b602082019050919050565b7f576f726c64000000000000000000000000000000000000000000000000000000600082015250565b6000611db56005836110c6565b9150611dc082611d7f565b602082019050919050565b60006040820190508181036000830152611de481611d5c565b90508181036020830152611df781611da8565b9050919050565b600081905092915050565b6000611e14826115c7565b611e1e8185611dfe565b9350611e2e8185602086016110d7565b80840191505092915050565b6000611e468284611e09565b915081905092915050565b600081905092915050565b6000611e67826110bb565b611e718185611e51565b9350611e818185602086016110d7565b80840191505092915050565b6000611e998284611e5c565b915081905092915050565b60008160601b9050919050565b6000611ebc82611ea4565b9050919050565b6000611ece82611eb1565b9050919050565b611ee6611ee1826113da565b611ec3565b82525050565b6000611ef88285611e5c565b9150611f048284611ed5565b6014820191508190509392505050565b6000606082019050611f2960008301866118a0565b611f3660208301856118a0565b611f43604083018461136c565b949350505050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112611f7757611f76611f4b565b5b80840192508235915067ffffffffffffffff821115611f9957611f98611f50565b5b602083019250600182023603831315611fb557611fb4611f55565b5b509250929050565b6000611fc98385611dfe565b9350611fd6838584611237565b82840190509392505050565b6000611fef828486611fbd565b91508190509392505050565b7f7265766572740000000000000000000000000000000000000000000000000000600082015250565b6000612031600683611e51565b915061203c82611ffb565b600682019050919050565b600061205282612024565b9150819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061209682611362565b91506120a183611362565b9250826120b1576120b061205c565b5b82820490509291505056fea2646970667358221220f219795264cc24c0fa5aad0aed682601b617e7e547c4844d06da6c607db3bd7e64736f6c634300081a0033",
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
