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
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"isZetaChain_\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"gateway_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"NO_MESSAGE_CALL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"amountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"calledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"erc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"erc20Call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"gasCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gateway\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"gatewayCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"}],\"name\":\"gatewayDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"gatewayDepositAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getAmountWithMessage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"getCalledWithMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"getNoMessageIndex\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isZetaChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"origin\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"internalType\":\"structTestDAppV2.zContext\",\"name\":\"_context\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_zrc20\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"internalType\":\"structTestDAppV2.MessageContext\",\"name\":\"messageContext\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"onCall\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"revertMessage\",\"type\":\"bytes\"}],\"internalType\":\"structTestDAppV2.RevertContext\",\"name\":\"revertContext\",\"type\":\"tuple\"}],\"name\":\"onRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"senderWithMessage\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"simpleCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60c060405234801561001057600080fd5b5060405161252f38038061252f83398181016040528101906100329190610114565b8115156080811515815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff16815250505050610154565b600080fd5b60008115159050919050565b6100938161007e565b811461009e57600080fd5b50565b6000815190506100b08161008a565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100e1826100b6565b9050919050565b6100f1816100d6565b81146100fc57600080fd5b50565b60008151905061010e816100e8565b92915050565b6000806040838503121561012b5761012a610079565b5b6000610139858286016100a1565b925050602061014a858286016100ff565b9150509250929050565b60805160a0516123856101aa6000396000818161045b015281816104d40152818161084e0152610fe90152600081816104a90152818161082301528181610abb01528181610f9a0152610fbe01526123856000f3fe60806040526004361061010d5760003560e01c8063ad23b28b11610095578063c91f356711610064578063c91f35671461035b578063deb3b1e414610386578063e2842ed7146103a2578063f592cbfb146103df578063f936ae851461041c57610114565b8063ad23b28b146102a1578063c7a339a9146102de578063c85f843414610307578063c9028a361461033257610114565b80635bcfd616116100dc5780635bcfd616146101d3578063676cc054146101fc5780639291fe261461022c5780639ca016ed14610269578063a799911f1461028557610114565b8063116191b61461011957806336e980a01461014457806341a3cd4a1461016d5780634297a2631461019657610114565b3661011457005b600080fd5b34801561012557600080fd5b5061012e610459565b60405161013b91906113d1565b60405180910390f35b34801561015057600080fd5b5061016b60048036038101906101669190611546565b61047d565b005b34801561017957600080fd5b50610194600480360381019061018f919061161b565b6104a7565b005b3480156101a257600080fd5b506101bd60048036038101906101b891906116b1565b6105ce565b6040516101ca91906116f7565b60405180910390f35b3480156101df57600080fd5b506101fa60048036038101906101f59190611762565b6105e6565b005b61021660048036038101906102119190611825565b6106cc565b6040516102239190611904565b60405180910390f35b34801561023857600080fd5b50610253600480360381019061024e9190611546565b6107de565b60405161026091906116f7565b60405180910390f35b610283600480360381019061027e9190611926565b610821565b005b61029f600480360381019061029a9190611546565b610943565b005b3480156102ad57600080fd5b506102c860048036038101906102c39190611926565b61096c565b6040516102d591906119a8565b60405180910390f35b3480156102ea57600080fd5b5061030560048036038101906103009190611a08565b6109cc565b005b34801561031357600080fd5b5061031c610a80565b60405161032991906119a8565b60405180910390f35b34801561033e57600080fd5b5061035960048036038101906103549190611a96565b610ab9565b005b34801561036757600080fd5b50610370610f98565b60405161037d9190611afa565b60405180910390f35b6103a0600480360381019061039b919061161b565b610fbc565b005b3480156103ae57600080fd5b506103c960048036038101906103c491906116b1565b6110e4565b6040516103d69190611afa565b60405180910390f35b3480156103eb57600080fd5b5061040660048036038101906104019190611546565b611104565b6040516104139190611afa565b60405180910390f35b34801561042857600080fd5b50610443600480360381019061043e9190611bb6565b611154565b60405161045091906113d1565b60405180910390f35b7f000000000000000000000000000000000000000000000000000000000000000081565b6104868161119d565b1561049057600080fd5b610499816111f3565b6104a4816000611247565b50565b7f0000000000000000000000000000000000000000000000000000000000000000156104d257600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16631becceb48484846040518060a001604052803373ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518563ffffffff1660e01b81526004016105979493929190611d19565b600060405180830381600087803b1580156105b157600080fd5b505af11580156105c5573d6000803e3d6000fd5b50505050505050565b60036020528060005260406000206000915090505481565b61063382828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505061119d565b1561063d57600080fd5b60008083839050146106935782828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506106af565b6106ae8660200160208101906106a99190611926565b61096c565b5b90506106ba816111f3565b6106c48185611247565b505050505050565b606060008084849050146107245783838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610740565b61073f85600001602081019061073a9190611926565b61096c565b5b905061074b816111f3565b6107558134611247565b8460000160208101906107689190611926565b6002826040516107789190611d9c565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604051806020016040528060008152509150509392505050565b600060036000836040516020016107f59190611def565b604051602081830303815290604052805190602001208152602001908152602001600020549050919050565b7f00000000000000000000000000000000000000000000000000000000000000001561084c57600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663726ac97c34836040518060a001604052803373ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518463ffffffff1660e01b815260040161090e929190611e06565b6000604051808303818588803b15801561092757600080fd5b505af115801561093b573d6000803e3d6000fd5b505050505050565b61094c8161119d565b1561095657600080fd5b61095f816111f3565b6109698134611247565b50565b60606040518060400160405280601681526020017f63616c6c65642077697468206e6f206d65737361676500000000000000000000815250826040516020016109b6929190611e7e565b6040516020818303038152906040529050919050565b6109d58161119d565b156109df57600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b8152600401610a1c93929190611ea6565b6020604051808303816000875af1158015610a3b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a5f9190611f09565b610a6857600080fd5b610a71816111f3565b610a7b8183611247565b505050565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525081565b7f000000000000000000000000000000000000000000000000000000000000000015610e5e57610ae7611289565b610b42818060600190610afa9190611f45565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506112fc565b15610e5d57600080826020016020810190610b5d9190611926565b73ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b81526004016040805180830381865afa158015610ba6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bca9190611fd2565b91509150826020016020810190610be19190611926565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614610c4e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c459061205e565b60405180910390fd5b8260400135811115610c95576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c8c906120f0565b60405180910390fd5b6000818460400135610ca7919061213f565b9050836020016020810190610cbc9190611926565b73ffffffffffffffffffffffffffffffffffffffff1663095ea7b33386604001356040518363ffffffff1660e01b8152600401610cfa929190612173565b6020604051808303816000875af1158015610d19573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d3d9190611f09565b503373ffffffffffffffffffffffffffffffffffffffff16637c0dcb5f856000016020810190610d6d9190611926565b604051602001610d7d91906113d1565b60405160208183030381529060405283876020016020810190610da09190611926565b6040518060a00160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518563ffffffff1660e01b8152600401610e27949392919061219c565b600060405180830381600087803b158015610e4157600080fd5b505af1158015610e55573d6000803e3d6000fd5b505050505050505b5b610eb9818060600190610e719190611f45565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506111f3565b610f16818060600190610ecc9190611f45565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506000611247565b806000016020810190610f299190611926565b6002828060600190610f3b9190611f45565b604051610f49929190612214565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000015610fe757600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663744b9b8b348585856040518060a001604052803373ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518663ffffffff1660e01b81526004016110ad9493929190611d19565b6000604051808303818588803b1580156110c657600080fd5b505af11580156110da573d6000803e3d6000fd5b5050505050505050565b60016020528060005260406000206000915054906101000a900460ff1681565b6000600160008360405160200161111b9190611def565b60405160208183030381529060405280519060200120815260200190815260200160002060009054906101000a900460ff169050919050565b6002818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60006040516020016111ae90612279565b60405160208183030381529060405280519060200120826040516020016111d59190611def565b60405160208183030381529060405280519060200120149050919050565b6001806000836040516020016112099190611def565b60405160208183030381529060405280519060200120815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b80600360008460405160200161125d9190611def565b604051602081830303815290604052805190602001208152602001908152602001600020819055505050565b60006207a12090506000614e209050600081836112a691906122bd565b905060005b818110156112e957600081908060018154018082558091505060019003906000526020600020016000909190919091505580806001019150506112ab565b506000806112f79190611352565b505050565b600060405160200161130d9061233a565b60405160208183030381529060405280519060200120826040516020016113349190611def565b60405160208183030381529060405280519060200120149050919050565b50805460008255906000526020600020908101906113709190611373565b50565b5b8082111561138c576000816000905550600101611374565b5090565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006113bb82611390565b9050919050565b6113cb816113b0565b82525050565b60006020820190506113e660008301846113c2565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6114538261140a565b810181811067ffffffffffffffff821117156114725761147161141b565b5b80604052505050565b60006114856113ec565b9050611491828261144a565b919050565b600067ffffffffffffffff8211156114b1576114b061141b565b5b6114ba8261140a565b9050602081019050919050565b82818337600083830152505050565b60006114e96114e484611496565b61147b565b90508281526020810184848401111561150557611504611405565b5b6115108482856114c7565b509392505050565b600082601f83011261152d5761152c611400565b5b813561153d8482602086016114d6565b91505092915050565b60006020828403121561155c5761155b6113f6565b5b600082013567ffffffffffffffff81111561157a576115796113fb565b5b61158684828501611518565b91505092915050565b611598816113b0565b81146115a357600080fd5b50565b6000813590506115b58161158f565b92915050565b600080fd5b600080fd5b60008083601f8401126115db576115da611400565b5b8235905067ffffffffffffffff8111156115f8576115f76115bb565b5b602083019150836001820283011115611614576116136115c0565b5b9250929050565b600080600060408486031215611634576116336113f6565b5b6000611642868287016115a6565b935050602084013567ffffffffffffffff811115611663576116626113fb565b5b61166f868287016115c5565b92509250509250925092565b6000819050919050565b61168e8161167b565b811461169957600080fd5b50565b6000813590506116ab81611685565b92915050565b6000602082840312156116c7576116c66113f6565b5b60006116d58482850161169c565b91505092915050565b6000819050919050565b6116f1816116de565b82525050565b600060208201905061170c60008301846116e8565b92915050565b600080fd5b60006060828403121561172d5761172c611712565b5b81905092915050565b61173f816116de565b811461174a57600080fd5b50565b60008135905061175c81611736565b92915050565b60008060008060006080868803121561177e5761177d6113f6565b5b600086013567ffffffffffffffff81111561179c5761179b6113fb565b5b6117a888828901611717565b95505060206117b9888289016115a6565b94505060406117ca8882890161174d565b935050606086013567ffffffffffffffff8111156117eb576117ea6113fb565b5b6117f7888289016115c5565b92509250509295509295909350565b60006020828403121561181c5761181b611712565b5b81905092915050565b60008060006040848603121561183e5761183d6113f6565b5b600061184c86828701611806565b935050602084013567ffffffffffffffff81111561186d5761186c6113fb565b5b611879868287016115c5565b92509250509250925092565b600081519050919050565b600082825260208201905092915050565b60005b838110156118bf5780820151818401526020810190506118a4565b60008484015250505050565b60006118d682611885565b6118e08185611890565b93506118f08185602086016118a1565b6118f98161140a565b840191505092915050565b6000602082019050818103600083015261191e81846118cb565b905092915050565b60006020828403121561193c5761193b6113f6565b5b600061194a848285016115a6565b91505092915050565b600081519050919050565b600082825260208201905092915050565b600061197a82611953565b611984818561195e565b93506119948185602086016118a1565b61199d8161140a565b840191505092915050565b600060208201905081810360008301526119c2818461196f565b905092915050565b60006119d5826113b0565b9050919050565b6119e5816119ca565b81146119f057600080fd5b50565b600081359050611a02816119dc565b92915050565b600080600060608486031215611a2157611a206113f6565b5b6000611a2f868287016119f3565b9350506020611a408682870161174d565b925050604084013567ffffffffffffffff811115611a6157611a606113fb565b5b611a6d86828701611518565b9150509250925092565b600060808284031215611a8d57611a8c611712565b5b81905092915050565b600060208284031215611aac57611aab6113f6565b5b600082013567ffffffffffffffff811115611aca57611ac96113fb565b5b611ad684828501611a77565b91505092915050565b60008115159050919050565b611af481611adf565b82525050565b6000602082019050611b0f6000830184611aeb565b92915050565b600067ffffffffffffffff821115611b3057611b2f61141b565b5b611b398261140a565b9050602081019050919050565b6000611b59611b5484611b15565b61147b565b905082815260208101848484011115611b7557611b74611405565b5b611b808482856114c7565b509392505050565b600082601f830112611b9d57611b9c611400565b5b8135611bad848260208601611b46565b91505092915050565b600060208284031215611bcc57611bcb6113f6565b5b600082013567ffffffffffffffff811115611bea57611be96113fb565b5b611bf684828501611b88565b91505092915050565b6000611c0b8385611890565b9350611c188385846114c7565b611c218361140a565b840190509392505050565b611c35816113b0565b82525050565b611c4481611adf565b82525050565b600082825260208201905092915050565b6000611c6682611885565b611c708185611c4a565b9350611c808185602086016118a1565b611c898161140a565b840191505092915050565b611c9d816116de565b82525050565b600060a083016000830151611cbb6000860182611c2c565b506020830151611cce6020860182611c3b565b506040830151611ce16040860182611c2c565b5060608301518482036060860152611cf98282611c5b565b9150506080830151611d0e6080860182611c94565b508091505092915050565b6000606082019050611d2e60008301876113c2565b8181036020830152611d41818587611bff565b90508181036040830152611d558184611ca3565b905095945050505050565b600081905092915050565b6000611d7682611885565b611d808185611d60565b9350611d908185602086016118a1565b80840191505092915050565b6000611da88284611d6b565b915081905092915050565b600081905092915050565b6000611dc982611953565b611dd38185611db3565b9350611de38185602086016118a1565b80840191505092915050565b6000611dfb8284611dbe565b915081905092915050565b6000604082019050611e1b60008301856113c2565b8181036020830152611e2d8184611ca3565b90509392505050565b60008160601b9050919050565b6000611e4e82611e36565b9050919050565b6000611e6082611e43565b9050919050565b611e78611e73826113b0565b611e55565b82525050565b6000611e8a8285611dbe565b9150611e968284611e67565b6014820191508190509392505050565b6000606082019050611ebb60008301866113c2565b611ec860208301856113c2565b611ed560408301846116e8565b949350505050565b611ee681611adf565b8114611ef157600080fd5b50565b600081519050611f0381611edd565b92915050565b600060208284031215611f1f57611f1e6113f6565b5b6000611f2d84828501611ef4565b91505092915050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112611f6257611f61611f36565b5b80840192508235915067ffffffffffffffff821115611f8457611f83611f3b565b5b602083019250600182023603831315611fa057611f9f611f40565b5b509250929050565b600081519050611fb78161158f565b92915050565b600081519050611fcc81611736565b92915050565b60008060408385031215611fe957611fe86113f6565b5b6000611ff785828601611fa8565b925050602061200885828601611fbd565b9150509250929050565b7f7a72633230206973206e6f742067617320746f6b656e00000000000000000000600082015250565b600061204860168361195e565b915061205382612012565b602082019050919050565b600060208201905081810360008301526120778161203b565b9050919050565b7f66656520616d6f756e7420697320686967686572207468616e2074686520616d60008201527f6f756e7400000000000000000000000000000000000000000000000000000000602082015250565b60006120da60248361195e565b91506120e58261207e565b604082019050919050565b60006020820190508181036000830152612109816120cd565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061214a826116de565b9150612155836116de565b925082820390508181111561216d5761216c612110565b5b92915050565b600060408201905061218860008301856113c2565b61219560208301846116e8565b9392505050565b600060808201905081810360008301526121b681876118cb565b90506121c560208301866116e8565b6121d260408301856113c2565b81810360608301526121e48184611ca3565b905095945050505050565b60006121fb8385611d60565b93506122088385846114c7565b82840190509392505050565b60006122218284866121ef565b91508190509392505050565b7f7265766572740000000000000000000000000000000000000000000000000000600082015250565b6000612263600683611db3565b915061226e8261222d565b600682019050919050565b600061228482612256565b9150819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006122c8826116de565b91506122d3836116de565b9250826122e3576122e261228e565b5b828204905092915050565b7f7769746864726177000000000000000000000000000000000000000000000000600082015250565b6000612324600883611db3565b915061232f826122ee565b600882019050919050565b600061234582612317565b915081905091905056fea2646970667358221220cc1d392f0803f3c30b8da3699bebbbdf5e1711d9c77c2bd74f9a96c12336d8f564736f6c634300081a0033",
}

// TestDAppV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDAppV2MetaData.ABI instead.
var TestDAppV2ABI = TestDAppV2MetaData.ABI

// TestDAppV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestDAppV2MetaData.Bin instead.
var TestDAppV2Bin = TestDAppV2MetaData.Bin

// DeployTestDAppV2 deploys a new Ethereum contract, binding an instance of TestDAppV2 to it.
func DeployTestDAppV2(auth *bind.TransactOpts, backend bind.ContractBackend, isZetaChain_ bool, gateway_ common.Address) (common.Address, *types.Transaction, *TestDAppV2, error) {
	parsed, err := TestDAppV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDAppV2Bin), backend, isZetaChain_, gateway_)
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

// Gateway is a free data retrieval call binding the contract method 0x116191b6.
//
// Solidity: function gateway() view returns(address)
func (_TestDAppV2 *TestDAppV2Caller) Gateway(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestDAppV2.contract.Call(opts, &out, "gateway")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Gateway is a free data retrieval call binding the contract method 0x116191b6.
//
// Solidity: function gateway() view returns(address)
func (_TestDAppV2 *TestDAppV2Session) Gateway() (common.Address, error) {
	return _TestDAppV2.Contract.Gateway(&_TestDAppV2.CallOpts)
}

// Gateway is a free data retrieval call binding the contract method 0x116191b6.
//
// Solidity: function gateway() view returns(address)
func (_TestDAppV2 *TestDAppV2CallerSession) Gateway() (common.Address, error) {
	return _TestDAppV2.Contract.Gateway(&_TestDAppV2.CallOpts)
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

// GatewayCall is a paid mutator transaction binding the contract method 0x41a3cd4a.
//
// Solidity: function gatewayCall(address dst, bytes payload) returns()
func (_TestDAppV2 *TestDAppV2Transactor) GatewayCall(opts *bind.TransactOpts, dst common.Address, payload []byte) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "gatewayCall", dst, payload)
}

// GatewayCall is a paid mutator transaction binding the contract method 0x41a3cd4a.
//
// Solidity: function gatewayCall(address dst, bytes payload) returns()
func (_TestDAppV2 *TestDAppV2Session) GatewayCall(dst common.Address, payload []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.GatewayCall(&_TestDAppV2.TransactOpts, dst, payload)
}

// GatewayCall is a paid mutator transaction binding the contract method 0x41a3cd4a.
//
// Solidity: function gatewayCall(address dst, bytes payload) returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) GatewayCall(dst common.Address, payload []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.GatewayCall(&_TestDAppV2.TransactOpts, dst, payload)
}

// GatewayDeposit is a paid mutator transaction binding the contract method 0x9ca016ed.
//
// Solidity: function gatewayDeposit(address dst) payable returns()
func (_TestDAppV2 *TestDAppV2Transactor) GatewayDeposit(opts *bind.TransactOpts, dst common.Address) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "gatewayDeposit", dst)
}

// GatewayDeposit is a paid mutator transaction binding the contract method 0x9ca016ed.
//
// Solidity: function gatewayDeposit(address dst) payable returns()
func (_TestDAppV2 *TestDAppV2Session) GatewayDeposit(dst common.Address) (*types.Transaction, error) {
	return _TestDAppV2.Contract.GatewayDeposit(&_TestDAppV2.TransactOpts, dst)
}

// GatewayDeposit is a paid mutator transaction binding the contract method 0x9ca016ed.
//
// Solidity: function gatewayDeposit(address dst) payable returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) GatewayDeposit(dst common.Address) (*types.Transaction, error) {
	return _TestDAppV2.Contract.GatewayDeposit(&_TestDAppV2.TransactOpts, dst)
}

// GatewayDepositAndCall is a paid mutator transaction binding the contract method 0xdeb3b1e4.
//
// Solidity: function gatewayDepositAndCall(address dst, bytes payload) payable returns()
func (_TestDAppV2 *TestDAppV2Transactor) GatewayDepositAndCall(opts *bind.TransactOpts, dst common.Address, payload []byte) (*types.Transaction, error) {
	return _TestDAppV2.contract.Transact(opts, "gatewayDepositAndCall", dst, payload)
}

// GatewayDepositAndCall is a paid mutator transaction binding the contract method 0xdeb3b1e4.
//
// Solidity: function gatewayDepositAndCall(address dst, bytes payload) payable returns()
func (_TestDAppV2 *TestDAppV2Session) GatewayDepositAndCall(dst common.Address, payload []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.GatewayDepositAndCall(&_TestDAppV2.TransactOpts, dst, payload)
}

// GatewayDepositAndCall is a paid mutator transaction binding the contract method 0xdeb3b1e4.
//
// Solidity: function gatewayDepositAndCall(address dst, bytes payload) payable returns()
func (_TestDAppV2 *TestDAppV2TransactorSession) GatewayDepositAndCall(dst common.Address, payload []byte) (*types.Transaction, error) {
	return _TestDAppV2.Contract.GatewayDepositAndCall(&_TestDAppV2.TransactOpts, dst, payload)
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
