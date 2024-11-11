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
	Bin: "0x6080604052348015600f57600080fd5b50611ff28061001f6000396000f3fe6080604052600436106100e15760003560e01c8063ad23b28b1161007f578063c9028a3611610059578063c9028a36146102c1578063e2842ed7146102ea578063f592cbfb14610327578063f936ae8514610364576100e8565b8063ad23b28b14610230578063c7a339a91461026d578063c85f843414610296576100e8565b80635bcfd616116100bb5780635bcfd6161461017e578063676cc054146101a75780639291fe26146101d7578063a799911f14610214576100e8565b806316ba7197146100ed57806336e980a0146101185780634297a26314610141576100e8565b366100e857005b600080fd5b3480156100f957600080fd5b506101026103a1565b60405161010f9190611116565b60405180910390f35b34801561012457600080fd5b5061013f600480360381019061013a9190611281565b6103da565b005b34801561014d57600080fd5b5061016860048036038101906101639190611300565b610404565b6040516101759190611346565b60405180910390f35b34801561018a57600080fd5b506101a560048036038101906101a0919061146f565b61041c565b005b6101c160048036038101906101bc9190611532565b61099e565b6040516101ce91906115e7565b60405180910390f35b3480156101e357600080fd5b506101fe60048036038101906101f99190611281565b610ab0565b60405161020b9190611346565b60405180910390f35b61022e60048036038101906102299190611281565b610af3565b005b34801561023c57600080fd5b5061025760048036038101906102529190611609565b610b1c565b6040516102649190611116565b60405180910390f35b34801561027957600080fd5b50610294600480360381019061028f9190611674565b610b7c565b005b3480156102a257600080fd5b506102ab610c30565b6040516102b89190611116565b60405180910390f35b3480156102cd57600080fd5b506102e860048036038101906102e39190611702565b610c69565b005b3480156102f657600080fd5b50610311600480360381019061030c9190611300565b610da3565b60405161031e9190611766565b60405180910390f35b34801561033357600080fd5b5061034e60048036038101906103499190611281565b610dc3565b60405161035b9190611766565b60405180910390f35b34801561037057600080fd5b5061038b60048036038101906103869190611822565b610e13565b604051610398919061187a565b60405180910390f35b6040518060400160405280600881526020017f776974686472617700000000000000000000000000000000000000000000000081525081565b6103e381610e5c565b156103ed57600080fd5b6103f681610eb2565b610401816000610f06565b50565b60036020528060005260406000206000915090505481565b61046982828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610e5c565b1561047357600080fd5b6104c082828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610f48565b1561090e576000808573ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b81526004016040805180830381865afa158015610512573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061053691906118bf565b915091508573ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16146105a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161059f9061194b565b60405180910390fd5b848111156105eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105e2906119dd565b60405180910390fd5b600081866105f99190611a2c565b905061063a6040518060400160405280600781526020017f6761736c656674000000000000000000000000000000000000000000000000008152505a610f06565b610642610fd5565b8673ffffffffffffffffffffffffffffffffffffffff1663095ea7b333886040518363ffffffff1660e01b815260040161067d929190611a60565b6020604051808303816000875af115801561069c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106c09190611ab5565b503373ffffffffffffffffffffffffffffffffffffffff16637c0dcb5f8960200160208101906106f09190611609565b604051602001610700919061187a565b604051602081830303815290604052838a6040518060a00160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518563ffffffff1660e01b81526004016107989493929190611bcf565b600060405180830381600087803b1580156107b257600080fd5b505af11580156107c6573d6000803e3d6000fd5b5050505060006040518060a001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020016000151581526020013373ffffffffffffffffffffffffffffffffffffffff1681526020016040518060400160405280600e81526020017f726576657274206d6573736167650000000000000000000000000000000000008152508152602001620186a081525090503373ffffffffffffffffffffffffffffffffffffffff16631cb5ea758a602001602081019061088b9190611609565b60405160200161089b919061187a565b6040516020818303038152906040528a8989620186a0876040518763ffffffff1660e01b81526004016108d396959493929190611c94565b600060405180830381600087803b1580156108ed57600080fd5b505af1158015610901573d6000803e3d6000fd5b5050505050505050610997565b60008083839050146109645782828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610980565b61097f86602001602081019061097a9190611609565b610b1c565b5b905061098b81610eb2565b6109958185610f06565b505b5050505050565b606060008084849050146109f65783838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610a12565b610a11856000016020810190610a0c9190611609565b610b1c565b5b9050610a1d81610eb2565b610a278134610f06565b846000016020810190610a3a9190611609565b600282604051610a4a9190611d3a565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604051806020016040528060008152509150509392505050565b60006003600083604051602001610ac79190611d8d565b604051602081830303815290604052805190602001208152602001908152602001600020549050919050565b610afc81610e5c565b15610b0657600080fd5b610b0f81610eb2565b610b198134610f06565b50565b60606040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525082604051602001610b66929190611dec565b6040516020818303038152906040529050919050565b610b8581610e5c565b15610b8f57600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b8152600401610bcc93929190611e14565b6020604051808303816000875af1158015610beb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c0f9190611ab5565b610c1857600080fd5b610c2181610eb2565b610c2b8183610f06565b505050565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525081565b610cc4818060600190610c7c9190611e5a565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610eb2565b610d21818060600190610cd79190611e5a565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506000610f06565b806000016020810190610d349190611609565b6002828060600190610d469190611e5a565b604051610d54929190611ee2565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60016020528060005260406000206000915054906101000a900460ff1681565b60006001600083604051602001610dda9190611d8d565b60405160208183030381529060405280519060200120815260200190815260200160002060009054906101000a900460ff169050919050565b6002818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000604051602001610e6d90611f47565b6040516020818303038152906040528051906020012082604051602001610e949190611d8d565b60405160208183030381529060405280519060200120149050919050565b600180600083604051602001610ec89190611d8d565b60405160208183030381529060405280519060200120815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b806003600084604051602001610f1c9190611d8d565b604051602081830303815290604052805190602001208152602001908152602001600020819055505050565b60006040518060400160405280600881526020017f7769746864726177000000000000000000000000000000000000000000000000815250604051602001610f909190611d8d565b6040516020818303038152906040528051906020012082604051602001610fb79190611d8d565b60405160208183030381529060405280519060200120149050919050565b6000620493e090506000614e20905060008183610ff29190611f8b565b905060005b818110156110355760008190806001815401808255809150506001900390600052602060002001600090919091909150558080600101915050610ff7565b506000806110439190611048565b505050565b50805460008255906000526020600020908101906110669190611069565b50565b5b8082111561108257600081600090555060010161106a565b5090565b600081519050919050565b600082825260208201905092915050565b60005b838110156110c05780820151818401526020810190506110a5565b60008484015250505050565b6000601f19601f8301169050919050565b60006110e882611086565b6110f28185611091565b93506111028185602086016110a2565b61110b816110cc565b840191505092915050565b6000602082019050818103600083015261113081846110dd565b905092915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61118e826110cc565b810181811067ffffffffffffffff821117156111ad576111ac611156565b5b80604052505050565b60006111c0611138565b90506111cc8282611185565b919050565b600067ffffffffffffffff8211156111ec576111eb611156565b5b6111f5826110cc565b9050602081019050919050565b82818337600083830152505050565b600061122461121f846111d1565b6111b6565b9050828152602081018484840111156112405761123f611151565b5b61124b848285611202565b509392505050565b600082601f8301126112685761126761114c565b5b8135611278848260208601611211565b91505092915050565b60006020828403121561129757611296611142565b5b600082013567ffffffffffffffff8111156112b5576112b4611147565b5b6112c184828501611253565b91505092915050565b6000819050919050565b6112dd816112ca565b81146112e857600080fd5b50565b6000813590506112fa816112d4565b92915050565b60006020828403121561131657611315611142565b5b6000611324848285016112eb565b91505092915050565b6000819050919050565b6113408161132d565b82525050565b600060208201905061135b6000830184611337565b92915050565b600080fd5b60006060828403121561137c5761137b611361565b5b81905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006113b082611385565b9050919050565b6113c0816113a5565b81146113cb57600080fd5b50565b6000813590506113dd816113b7565b92915050565b6113ec8161132d565b81146113f757600080fd5b50565b600081359050611409816113e3565b92915050565b600080fd5b600080fd5b60008083601f84011261142f5761142e61114c565b5b8235905067ffffffffffffffff81111561144c5761144b61140f565b5b60208301915083600182028301111561146857611467611414565b5b9250929050565b60008060008060006080868803121561148b5761148a611142565b5b600086013567ffffffffffffffff8111156114a9576114a8611147565b5b6114b588828901611366565b95505060206114c6888289016113ce565b94505060406114d7888289016113fa565b935050606086013567ffffffffffffffff8111156114f8576114f7611147565b5b61150488828901611419565b92509250509295509295909350565b60006020828403121561152957611528611361565b5b81905092915050565b60008060006040848603121561154b5761154a611142565b5b600061155986828701611513565b935050602084013567ffffffffffffffff81111561157a57611579611147565b5b61158686828701611419565b92509250509250925092565b600081519050919050565b600082825260208201905092915050565b60006115b982611592565b6115c3818561159d565b93506115d38185602086016110a2565b6115dc816110cc565b840191505092915050565b6000602082019050818103600083015261160181846115ae565b905092915050565b60006020828403121561161f5761161e611142565b5b600061162d848285016113ce565b91505092915050565b6000611641826113a5565b9050919050565b61165181611636565b811461165c57600080fd5b50565b60008135905061166e81611648565b92915050565b60008060006060848603121561168d5761168c611142565b5b600061169b8682870161165f565b93505060206116ac868287016113fa565b925050604084013567ffffffffffffffff8111156116cd576116cc611147565b5b6116d986828701611253565b9150509250925092565b6000608082840312156116f9576116f8611361565b5b81905092915050565b60006020828403121561171857611717611142565b5b600082013567ffffffffffffffff81111561173657611735611147565b5b611742848285016116e3565b91505092915050565b60008115159050919050565b6117608161174b565b82525050565b600060208201905061177b6000830184611757565b92915050565b600067ffffffffffffffff82111561179c5761179b611156565b5b6117a5826110cc565b9050602081019050919050565b60006117c56117c084611781565b6111b6565b9050828152602081018484840111156117e1576117e0611151565b5b6117ec848285611202565b509392505050565b600082601f8301126118095761180861114c565b5b81356118198482602086016117b2565b91505092915050565b60006020828403121561183857611837611142565b5b600082013567ffffffffffffffff81111561185657611855611147565b5b611862848285016117f4565b91505092915050565b611874816113a5565b82525050565b600060208201905061188f600083018461186b565b92915050565b6000815190506118a4816113b7565b92915050565b6000815190506118b9816113e3565b92915050565b600080604083850312156118d6576118d5611142565b5b60006118e485828601611895565b92505060206118f5858286016118aa565b9150509250929050565b7f7a72633230206973206e6f742067617320746f6b656e00000000000000000000600082015250565b6000611935601683611091565b9150611940826118ff565b602082019050919050565b6000602082019050818103600083015261196481611928565b9050919050565b7f66656520616d6f756e7420697320686967686572207468616e2074686520616d60008201527f6f756e7400000000000000000000000000000000000000000000000000000000602082015250565b60006119c7602483611091565b91506119d28261196b565b604082019050919050565b600060208201905081810360008301526119f6816119ba565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611a378261132d565b9150611a428361132d565b9250828203905081811115611a5a57611a596119fd565b5b92915050565b6000604082019050611a75600083018561186b565b611a826020830184611337565b9392505050565b611a928161174b565b8114611a9d57600080fd5b50565b600081519050611aaf81611a89565b92915050565b600060208284031215611acb57611aca611142565b5b6000611ad984828501611aa0565b91505092915050565b611aeb816113a5565b82525050565b611afa8161174b565b82525050565b600082825260208201905092915050565b6000611b1c82611592565b611b268185611b00565b9350611b368185602086016110a2565b611b3f816110cc565b840191505092915050565b611b538161132d565b82525050565b600060a083016000830151611b716000860182611ae2565b506020830151611b846020860182611af1565b506040830151611b976040860182611ae2565b5060608301518482036060860152611baf8282611b11565b9150506080830151611bc46080860182611b4a565b508091505092915050565b60006080820190508181036000830152611be981876115ae565b9050611bf86020830186611337565b611c05604083018561186b565b8181036060830152611c178184611b59565b905095945050505050565b6000611c2e838561159d565b9350611c3b838584611202565b611c44836110cc565b840190509392505050565b6000819050919050565b6000819050919050565b6000611c7e611c79611c7484611c4f565b611c59565b61132d565b9050919050565b611c8e81611c63565b82525050565b600060a0820190508181036000830152611cae81896115ae565b9050611cbd602083018861186b565b8181036040830152611cd0818688611c22565b9050611cdf6060830185611c85565b8181036080830152611cf18184611b59565b9050979650505050505050565b600081905092915050565b6000611d1482611592565b611d1e8185611cfe565b9350611d2e8185602086016110a2565b80840191505092915050565b6000611d468284611d09565b915081905092915050565b600081905092915050565b6000611d6782611086565b611d718185611d51565b9350611d818185602086016110a2565b80840191505092915050565b6000611d998284611d5c565b915081905092915050565b60008160601b9050919050565b6000611dbc82611da4565b9050919050565b6000611dce82611db1565b9050919050565b611de6611de1826113a5565b611dc3565b82525050565b6000611df88285611d5c565b9150611e048284611dd5565b6014820191508190509392505050565b6000606082019050611e29600083018661186b565b611e36602083018561186b565b611e436040830184611337565b949350505050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112611e7757611e76611e4b565b5b80840192508235915067ffffffffffffffff821115611e9957611e98611e50565b5b602083019250600182023603831315611eb557611eb4611e55565b5b509250929050565b6000611ec98385611cfe565b9350611ed6838584611202565b82840190509392505050565b6000611eef828486611ebd565b91508190509392505050565b7f7265766572740000000000000000000000000000000000000000000000000000600082015250565b6000611f31600683611d51565b9150611f3c82611efb565b600682019050919050565b6000611f5282611f24565b9150819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000611f968261132d565b9150611fa18361132d565b925082611fb157611fb0611f5c565b5b82820490509291505056fea264697066735822122086320805b66b54f60d9089b0be9b6ae44a1cc8b8aeb25c31eed548878bf1215e64736f6c634300081a0033",
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
