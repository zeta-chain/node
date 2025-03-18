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
	Bin: "0x60c060405234801561001057600080fd5b5060405161259438038061259483398181016040528101906100329190610114565b8115156080811515815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff16815250505050610154565b600080fd5b60008115159050919050565b6100938161007e565b811461009e57600080fd5b50565b6000815190506100b08161008a565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100e1826100b6565b9050919050565b6100f1816100d6565b81146100fc57600080fd5b50565b60008151905061010e816100e8565b92915050565b6000806040838503121561012b5761012a610079565b5b6000610139858286016100a1565b925050602061014a858286016100ff565b9150509250929050565b60805160a0516123ea6101aa6000396000818161045b015281816104d40152818161084e015261104e0152600081816104a90152818161082301528181610b2001528181610fff015261102301526123ea6000f3fe60806040526004361061010d5760003560e01c8063ad23b28b11610095578063c91f356711610064578063c91f35671461035b578063deb3b1e414610386578063e2842ed7146103a2578063f592cbfb146103df578063f936ae851461041c57610114565b8063ad23b28b146102a1578063c7a339a9146102de578063c85f843414610307578063c9028a361461033257610114565b80635bcfd616116100dc5780635bcfd616146101d3578063676cc054146101fc5780639291fe261461022c5780639ca016ed14610269578063a799911f1461028557610114565b8063116191b61461011957806336e980a01461014457806341a3cd4a1461016d5780634297a2631461019657610114565b3661011457005b600080fd5b34801561012557600080fd5b5061012e610459565b60405161013b9190611436565b60405180910390f35b34801561015057600080fd5b5061016b600480360381019061016691906115ab565b61047d565b005b34801561017957600080fd5b50610194600480360381019061018f9190611680565b6104a7565b005b3480156101a257600080fd5b506101bd60048036038101906101b89190611716565b6105ce565b6040516101ca919061175c565b60405180910390f35b3480156101df57600080fd5b506101fa60048036038101906101f591906117c7565b6105e6565b005b6102166004803603810190610211919061188a565b6106cc565b6040516102239190611969565b60405180910390f35b34801561023857600080fd5b50610253600480360381019061024e91906115ab565b6107de565b604051610260919061175c565b60405180910390f35b610283600480360381019061027e919061198b565b610821565b005b61029f600480360381019061029a91906115ab565b610943565b005b3480156102ad57600080fd5b506102c860048036038101906102c3919061198b565b61096c565b6040516102d59190611a0d565b60405180910390f35b3480156102ea57600080fd5b5061030560048036038101906103009190611a6d565b6109cc565b005b34801561031357600080fd5b5061031c610a80565b6040516103299190611a0d565b60405180910390f35b34801561033e57600080fd5b5061035960048036038101906103549190611afb565b610ab9565b005b34801561036757600080fd5b50610370610ffd565b60405161037d9190611b5f565b60405180910390f35b6103a0600480360381019061039b9190611680565b611021565b005b3480156103ae57600080fd5b506103c960048036038101906103c49190611716565b611149565b6040516103d69190611b5f565b60405180910390f35b3480156103eb57600080fd5b50610406600480360381019061040191906115ab565b611169565b6040516104139190611b5f565b60405180910390f35b34801561042857600080fd5b50610443600480360381019061043e9190611c1b565b6111b9565b6040516104509190611436565b60405180910390f35b7f000000000000000000000000000000000000000000000000000000000000000081565b61048681611202565b1561049057600080fd5b61049981611258565b6104a48160006112ac565b50565b7f0000000000000000000000000000000000000000000000000000000000000000156104d257600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16631becceb48484846040518060a001604052803373ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518563ffffffff1660e01b81526004016105979493929190611d7e565b600060405180830381600087803b1580156105b157600080fd5b505af11580156105c5573d6000803e3d6000fd5b50505050505050565b60036020528060005260406000206000915090505481565b61063382828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050611202565b1561063d57600080fd5b60008083839050146106935782828080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050506106af565b6106ae8660200160208101906106a9919061198b565b61096c565b5b90506106ba81611258565b6106c481856112ac565b505050505050565b606060008084849050146107245783838080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050610740565b61073f85600001602081019061073a919061198b565b61096c565b5b905061074b81611258565b61075581346112ac565b846000016020810190610768919061198b565b6002826040516107789190611e01565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550604051806020016040528060008152509150509392505050565b600060036000836040516020016107f59190611e54565b604051602081830303815290604052805190602001208152602001908152602001600020549050919050565b7f00000000000000000000000000000000000000000000000000000000000000001561084c57600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663726ac97c34836040518060a001604052803373ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518463ffffffff1660e01b815260040161090e929190611e6b565b6000604051808303818588803b15801561092757600080fd5b505af115801561093b573d6000803e3d6000fd5b505050505050565b61094c81611202565b1561095657600080fd5b61095f81611258565b61096981346112ac565b50565b60606040518060400160405280601681526020017f63616c6c65642077697468206e6f206d65737361676500000000000000000000815250826040516020016109b6929190611ee3565b6040516020818303038152906040529050919050565b6109d581611202565b156109df57600080fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b8152600401610a1c93929190611f0b565b6020604051808303816000875af1158015610a3b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a5f9190611f6e565b610a6857600080fd5b610a7181611258565b610a7b81836112ac565b505050565b6040518060400160405280601681526020017f63616c6c65642077697468206e6f206d6573736167650000000000000000000081525081565b610b14818060600190610acc9190611faa565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050611202565b15610b1e57600080fd5b7f000000000000000000000000000000000000000000000000000000000000000015610ec357610b4c6112ee565b610ba7818060600190610b5f9190611faa565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050611361565b15610ec257600080826020016020810190610bc2919061198b565b73ffffffffffffffffffffffffffffffffffffffff1663d9eeebed6040518163ffffffff1660e01b81526004016040805180830381865afa158015610c0b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c2f9190612037565b91509150826020016020810190610c46919061198b565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614610cb3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610caa906120c3565b60405180910390fd5b8260400135811115610cfa576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cf190612155565b60405180910390fd5b6000818460400135610d0c91906121a4565b9050836020016020810190610d21919061198b565b73ffffffffffffffffffffffffffffffffffffffff1663095ea7b33386604001356040518363ffffffff1660e01b8152600401610d5f9291906121d8565b6020604051808303816000875af1158015610d7e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610da29190611f6e565b503373ffffffffffffffffffffffffffffffffffffffff16637c0dcb5f856000016020810190610dd2919061198b565b604051602001610de29190611436565b60405160208183030381529060405283876020016020810190610e05919061198b565b6040518060a00160405280600073ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518563ffffffff1660e01b8152600401610e8c9493929190612201565b600060405180830381600087803b158015610ea657600080fd5b505af1158015610eba573d6000803e3d6000fd5b505050505050505b5b610f1e818060600190610ed69190611faa565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050611258565b610f7b818060600190610f319190611faa565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505060006112ac565b806000016020810190610f8e919061198b565b6002828060600190610fa09190611faa565b604051610fae929190612279565b908152602001604051809103902060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f00000000000000000000000000000000000000000000000000000000000000001561104c57600080fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663744b9b8b348585856040518060a001604052803373ffffffffffffffffffffffffffffffffffffffff168152602001600015158152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160405180602001604052806000815250815260200160008152506040518663ffffffff1660e01b81526004016111129493929190611d7e565b6000604051808303818588803b15801561112b57600080fd5b505af115801561113f573d6000803e3d6000fd5b5050505050505050565b60016020528060005260406000206000915054906101000a900460ff1681565b600060016000836040516020016111809190611e54565b60405160208183030381529060405280519060200120815260200190815260200160002060009054906101000a900460ff169050919050565b6002818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000604051602001611213906122de565b604051602081830303815290604052805190602001208260405160200161123a9190611e54565b60405160208183030381529060405280519060200120149050919050565b60018060008360405160200161126e9190611e54565b60405160208183030381529060405280519060200120815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b8060036000846040516020016112c29190611e54565b604051602081830303815290604052805190602001208152602001908152602001600020819055505050565b60006207a12090506000614e2090506000818361130b9190612322565b905060005b8181101561134e5760008190806001815401808255809150506001900390600052602060002001600090919091909150558080600101915050611310565b5060008061135c91906113b7565b505050565b60006040516020016113729061239f565b60405160208183030381529060405280519060200120826040516020016113999190611e54565b60405160208183030381529060405280519060200120149050919050565b50805460008255906000526020600020908101906113d591906113d8565b50565b5b808211156113f15760008160009055506001016113d9565b5090565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000611420826113f5565b9050919050565b61143081611415565b82525050565b600060208201905061144b6000830184611427565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6114b88261146f565b810181811067ffffffffffffffff821117156114d7576114d6611480565b5b80604052505050565b60006114ea611451565b90506114f682826114af565b919050565b600067ffffffffffffffff82111561151657611515611480565b5b61151f8261146f565b9050602081019050919050565b82818337600083830152505050565b600061154e611549846114fb565b6114e0565b90508281526020810184848401111561156a5761156961146a565b5b61157584828561152c565b509392505050565b600082601f83011261159257611591611465565b5b81356115a284826020860161153b565b91505092915050565b6000602082840312156115c1576115c061145b565b5b600082013567ffffffffffffffff8111156115df576115de611460565b5b6115eb8482850161157d565b91505092915050565b6115fd81611415565b811461160857600080fd5b50565b60008135905061161a816115f4565b92915050565b600080fd5b600080fd5b60008083601f8401126116405761163f611465565b5b8235905067ffffffffffffffff81111561165d5761165c611620565b5b60208301915083600182028301111561167957611678611625565b5b9250929050565b6000806000604084860312156116995761169861145b565b5b60006116a78682870161160b565b935050602084013567ffffffffffffffff8111156116c8576116c7611460565b5b6116d48682870161162a565b92509250509250925092565b6000819050919050565b6116f3816116e0565b81146116fe57600080fd5b50565b600081359050611710816116ea565b92915050565b60006020828403121561172c5761172b61145b565b5b600061173a84828501611701565b91505092915050565b6000819050919050565b61175681611743565b82525050565b6000602082019050611771600083018461174d565b92915050565b600080fd5b60006060828403121561179257611791611777565b5b81905092915050565b6117a481611743565b81146117af57600080fd5b50565b6000813590506117c18161179b565b92915050565b6000806000806000608086880312156117e3576117e261145b565b5b600086013567ffffffffffffffff81111561180157611800611460565b5b61180d8882890161177c565b955050602061181e8882890161160b565b945050604061182f888289016117b2565b935050606086013567ffffffffffffffff8111156118505761184f611460565b5b61185c8882890161162a565b92509250509295509295909350565b60006020828403121561188157611880611777565b5b81905092915050565b6000806000604084860312156118a3576118a261145b565b5b60006118b18682870161186b565b935050602084013567ffffffffffffffff8111156118d2576118d1611460565b5b6118de8682870161162a565b92509250509250925092565b600081519050919050565b600082825260208201905092915050565b60005b83811015611924578082015181840152602081019050611909565b60008484015250505050565b600061193b826118ea565b61194581856118f5565b9350611955818560208601611906565b61195e8161146f565b840191505092915050565b600060208201905081810360008301526119838184611930565b905092915050565b6000602082840312156119a1576119a061145b565b5b60006119af8482850161160b565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60006119df826119b8565b6119e981856119c3565b93506119f9818560208601611906565b611a028161146f565b840191505092915050565b60006020820190508181036000830152611a2781846119d4565b905092915050565b6000611a3a82611415565b9050919050565b611a4a81611a2f565b8114611a5557600080fd5b50565b600081359050611a6781611a41565b92915050565b600080600060608486031215611a8657611a8561145b565b5b6000611a9486828701611a58565b9350506020611aa5868287016117b2565b925050604084013567ffffffffffffffff811115611ac657611ac5611460565b5b611ad28682870161157d565b9150509250925092565b600060808284031215611af257611af1611777565b5b81905092915050565b600060208284031215611b1157611b1061145b565b5b600082013567ffffffffffffffff811115611b2f57611b2e611460565b5b611b3b84828501611adc565b91505092915050565b60008115159050919050565b611b5981611b44565b82525050565b6000602082019050611b746000830184611b50565b92915050565b600067ffffffffffffffff821115611b9557611b94611480565b5b611b9e8261146f565b9050602081019050919050565b6000611bbe611bb984611b7a565b6114e0565b905082815260208101848484011115611bda57611bd961146a565b5b611be584828561152c565b509392505050565b600082601f830112611c0257611c01611465565b5b8135611c12848260208601611bab565b91505092915050565b600060208284031215611c3157611c3061145b565b5b600082013567ffffffffffffffff811115611c4f57611c4e611460565b5b611c5b84828501611bed565b91505092915050565b6000611c7083856118f5565b9350611c7d83858461152c565b611c868361146f565b840190509392505050565b611c9a81611415565b82525050565b611ca981611b44565b82525050565b600082825260208201905092915050565b6000611ccb826118ea565b611cd58185611caf565b9350611ce5818560208601611906565b611cee8161146f565b840191505092915050565b611d0281611743565b82525050565b600060a083016000830151611d206000860182611c91565b506020830151611d336020860182611ca0565b506040830151611d466040860182611c91565b5060608301518482036060860152611d5e8282611cc0565b9150506080830151611d736080860182611cf9565b508091505092915050565b6000606082019050611d936000830187611427565b8181036020830152611da6818587611c64565b90508181036040830152611dba8184611d08565b905095945050505050565b600081905092915050565b6000611ddb826118ea565b611de58185611dc5565b9350611df5818560208601611906565b80840191505092915050565b6000611e0d8284611dd0565b915081905092915050565b600081905092915050565b6000611e2e826119b8565b611e388185611e18565b9350611e48818560208601611906565b80840191505092915050565b6000611e608284611e23565b915081905092915050565b6000604082019050611e806000830185611427565b8181036020830152611e928184611d08565b90509392505050565b60008160601b9050919050565b6000611eb382611e9b565b9050919050565b6000611ec582611ea8565b9050919050565b611edd611ed882611415565b611eba565b82525050565b6000611eef8285611e23565b9150611efb8284611ecc565b6014820191508190509392505050565b6000606082019050611f206000830186611427565b611f2d6020830185611427565b611f3a604083018461174d565b949350505050565b611f4b81611b44565b8114611f5657600080fd5b50565b600081519050611f6881611f42565b92915050565b600060208284031215611f8457611f8361145b565b5b6000611f9284828501611f59565b91505092915050565b600080fd5b600080fd5b600080fd5b60008083356001602003843603038112611fc757611fc6611f9b565b5b80840192508235915067ffffffffffffffff821115611fe957611fe8611fa0565b5b60208301925060018202360383131561200557612004611fa5565b5b509250929050565b60008151905061201c816115f4565b92915050565b6000815190506120318161179b565b92915050565b6000806040838503121561204e5761204d61145b565b5b600061205c8582860161200d565b925050602061206d85828601612022565b9150509250929050565b7f7a72633230206973206e6f742067617320746f6b656e00000000000000000000600082015250565b60006120ad6016836119c3565b91506120b882612077565b602082019050919050565b600060208201905081810360008301526120dc816120a0565b9050919050565b7f66656520616d6f756e7420697320686967686572207468616e2074686520616d60008201527f6f756e7400000000000000000000000000000000000000000000000000000000602082015250565b600061213f6024836119c3565b915061214a826120e3565b604082019050919050565b6000602082019050818103600083015261216e81612132565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006121af82611743565b91506121ba83611743565b92508282039050818111156121d2576121d1612175565b5b92915050565b60006040820190506121ed6000830185611427565b6121fa602083018461174d565b9392505050565b6000608082019050818103600083015261221b8187611930565b905061222a602083018661174d565b6122376040830185611427565b81810360608301526122498184611d08565b905095945050505050565b60006122608385611dc5565b935061226d83858461152c565b82840190509392505050565b6000612286828486612254565b91508190509392505050565b7f7265766572740000000000000000000000000000000000000000000000000000600082015250565b60006122c8600683611e18565b91506122d382612292565b600682019050919050565b60006122e9826122bb565b9150819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061232d82611743565b915061233883611743565b925082612348576123476122f3565b5b828204905092915050565b7f7769746864726177000000000000000000000000000000000000000000000000600082015250565b6000612389600883611e18565b915061239482612353565b600882019050919050565b60006123aa8261237c565b915081905091905056fea2646970667358221220e0e7dfdd94bdbc9796b718d6cbbf5d94e34a34e70ed21a1f60997215fbaed17964736f6c634300081a0033",
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
