// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc20custody

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
)

// ERC20CustodyMetaData contains all meta data concerning the ERC20Custody contract.
var ERC20CustodyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"TSSAddress_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"TSSAddressUpdater_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"zetaFee_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"zetaMaxFee_\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20\",\"name\":\"zeta_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTSSUpdater\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IsPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotWhitelisted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZetaMaxFeeExceeded\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"TSSAddressUpdater_\",\"type\":\"address\"}],\"name\":\"RenouncedTSSUpdater\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"Unwhitelisted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"TSSAddress_\",\"type\":\"address\"}],\"name\":\"UpdatedTSSAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"zetaFee_\",\"type\":\"uint256\"}],\"name\":\"UpdatedZetaFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"Whitelisted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"TSSAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TSSAddressUpdater\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"},{\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceTSSAddressUpdater\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"unwhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"TSSAddress_\",\"type\":\"address\"}],\"name\":\"updateTSSAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"zetaFee_\",\"type\":\"uint256\"}],\"name\":\"updateZetaFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"whitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"whitelisted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zeta\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zetaFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zetaMaxFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162001bf138038062001bf183398181016040528101906200003791906200014f565b84600060016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555083600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550826002819055508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1660601b81525050816080818152505050505050506200027c565b6000815190506200011b816200022e565b92915050565b600081519050620001328162000248565b92915050565b600081519050620001498162000262565b92915050565b600080600080600060a086880312156200016e576200016d62000229565b5b60006200017e888289016200010a565b955050602062000191888289016200010a565b9450506040620001a48882890162000138565b9350506060620001b78882890162000138565b9250506080620001ca8882890162000121565b9150509295509295909350565b6000620001e482620001ff565b9050919050565b6000620001f882620001d7565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600080fd5b6200023981620001d7565b81146200024557600080fd5b50565b6200025381620001eb565b81146200025f57600080fd5b50565b6200026d816200021f565b81146200027957600080fd5b50565b60805160a05160601c611937620002ba60003960008181610e6001528181610e9f015261118301526000818161042b0152610cd101526119376000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80639b19251a11610097578063e5408cfa11610066578063e5408cfa1461024d578063e609055e1461026b578063e8f9cb3a14610287578063ed11692b146102a557610100565b80639b19251a146101c9578063d936547e146101e5578063d9caed1214610215578063de2f6c5e1461023157610100565b80637bdaded3116100d35780637bdaded3146101695780638456cb5914610187578063950837aa146101915780639a590427146101ad57610100565b80633f4ba83a1461010557806353ee30a31461010f57806354b61e811461012d5780635c975abb1461014b575b600080fd5b61010d6102af565b005b6101176103cc565b6040516101249190611653565b60405180910390f35b6101356103f2565b6040516101429190611653565b60405180910390f35b610153610418565b60405161016091906116ce565b60405180910390f35b610171610429565b60405161017e919061174d565b60405180910390f35b61018f61044d565b005b6101ab60048036038101906101a6919061140f565b6105f4565b005b6101c760048036038101906101c29190611563565b61075c565b005b6101e360048036038101906101de9190611563565b610881565b005b6101ff60048036038101906101fa9190611563565b6109a6565b60405161020c91906116ce565b60405180910390f35b61022f600480360381019061022a919061143c565b6109c6565b005b61024b60048036038101906102469190611590565b610c0d565b005b610255610d6a565b604051610262919061174d565b60405180910390f35b610285600480360381019061028091906114bc565b610d70565b005b61028f611181565b60405161029c9190611732565b60405180910390f35b6102ad6111a5565b005b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610335576040517e611fa600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008054906101000a900460ff16610379576040517f6cd6020100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060006101000a81548160ff0219169083151502179055507f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa336040516103c29190611653565b60405180910390a1565b600060019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008054906101000a900460ff1681565b7f000000000000000000000000000000000000000000000000000000000000000081565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146104d3576040517e611fa600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008054906101000a900460ff1615610518576040517f1309a56300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16600060019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614156105a1576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016000806101000a81548160ff0219169083151502179055507f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258336040516105ea9190611653565b60405180910390a1565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461067a576040517e611fa600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156106e1576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600060016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507fcd2958db8285a532edf298cbe1aa28ea3fb5ec82461253f9a8c1477924e29947816040516107519190611653565b60405180910390a150565b600060019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146107e3576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508073ffffffffffffffffffffffffffffffffffffffff167f51085ddf9ebdded84b76e829eb58c4078e4b5bdf97d9a94723f336039da4679160405160405180910390a250565b600060019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610908576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508073ffffffffffffffffffffffffffffffffffffffff167faab7954e9d246b167ef88aeddad35209ca2489d95a8aeb59e288d9b19fae5a5460405160405180910390a250565b60036020528060005260406000206000915054906101000a900460ff1681565b600060019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610a4d576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008054906101000a900460ff1615610a92576040517f1309a56300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610b15576040517f584a793800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84836040518363ffffffff1660e01b8152600401610b509291906116a5565b602060405180830381600087803b158015610b6a57600080fd5b505af1158015610b7e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ba2919061148f565b508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb83604051610c00919061174d565b60405180910390a3505050565b600060019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610c94576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000811415610ccf576040517faf13986d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f0000000000000000000000000000000000000000000000000000000000000000811115610d29576040517fc1be451300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806002819055507f6d2d8e313fbaf76898bb9fa55e4b52525e49c7d7182d0874f97bd9076e81d52381604051610d5f919061174d565b60405180910390a150565b60025481565b60008054906101000a900460ff1615610db5576040517f1309a56300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610e38576040517f584a793800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060025414158015610e985750600073ffffffffffffffffffffffffffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1614155b15610f72577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166323b872dd33600060019054906101000a900473ffffffffffffffffffffffffffffffffffffffff166002546040518463ffffffff1660e01b8152600401610f1e9392919061166e565b602060405180830381600087803b158015610f3857600080fd5b505af1158015610f4c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f70919061148f565b505b60008473ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610fad9190611653565b60206040518083038186803b158015610fc557600080fd5b505afa158015610fd9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ffd91906115bd565b90508473ffffffffffffffffffffffffffffffffffffffff166323b872dd3330876040518463ffffffff1660e01b815260040161103c9392919061166e565b602060405180830381600087803b15801561105657600080fd5b505af115801561106a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061108e919061148f565b508473ffffffffffffffffffffffffffffffffffffffff167f1dafa057cc5c3bccb5ad974129a2bccd3c74002d9dfd7062404ba9523b18d6ae8888848973ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016111039190611653565b60206040518083038186803b15801561111b57600080fd5b505afa15801561112f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061115391906115bd565b61115d9190611779565b87876040516111709594939291906116e9565b60405180910390a250505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461122b576040517e611fa600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16600060019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614156112b4576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060019054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507f39ac266bfc011581be62c138d96e4e8782812013bb66fffb4cd207f4bfee3777336040516113469190611653565b60405180910390a1565b60008135905061135f816118a5565b92915050565b600081519050611374816118bc565b92915050565b60008083601f8401126113905761138f611880565b5b8235905067ffffffffffffffff8111156113ad576113ac61187b565b5b6020830191508360018202830111156113c9576113c8611885565b5b9250929050565b6000813590506113df816118d3565b92915050565b6000813590506113f4816118ea565b92915050565b600081519050611409816118ea565b92915050565b6000602082840312156114255761142461188f565b5b600061143384828501611350565b91505092915050565b6000806000606084860312156114555761145461188f565b5b600061146386828701611350565b9350506020611474868287016113d0565b9250506040611485868287016113e5565b9150509250925092565b6000602082840312156114a5576114a461188f565b5b60006114b384828501611365565b91505092915050565b600080600080600080608087890312156114d9576114d861188f565b5b600087013567ffffffffffffffff8111156114f7576114f661188a565b5b61150389828a0161137a565b9650965050602061151689828a016113d0565b945050604061152789828a016113e5565b935050606087013567ffffffffffffffff8111156115485761154761188a565b5b61155489828a0161137a565b92509250509295509295509295565b6000602082840312156115795761157861188f565b5b6000611587848285016113d0565b91505092915050565b6000602082840312156115a6576115a561188f565b5b60006115b4848285016113e5565b91505092915050565b6000602082840312156115d3576115d261188f565b5b60006115e1848285016113fa565b91505092915050565b6115f3816117ad565b82525050565b611602816117bf565b82525050565b60006116148385611768565b935061162183858461183d565b61162a83611894565b840190509392505050565b61163e81611807565b82525050565b61164d816117fd565b82525050565b600060208201905061166860008301846115ea565b92915050565b600060608201905061168360008301866115ea565b61169060208301856115ea565b61169d6040830184611644565b949350505050565b60006040820190506116ba60008301856115ea565b6116c76020830184611644565b9392505050565b60006020820190506116e360008301846115f9565b92915050565b60006060820190508181036000830152611704818789611608565b90506117136020830186611644565b8181036040830152611726818486611608565b90509695505050505050565b60006020820190506117476000830184611635565b92915050565b60006020820190506117626000830184611644565b92915050565b600082825260208201905092915050565b6000611784826117fd565b915061178f836117fd565b9250828210156117a2576117a161184c565b5b828203905092915050565b60006117b8826117dd565b9050919050565b60008115159050919050565b60006117d6826117ad565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600061181282611819565b9050919050565b60006118248261182b565b9050919050565b6000611836826117dd565b9050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b6118ae816117ad565b81146118b957600080fd5b50565b6118c5816117bf565b81146118d057600080fd5b50565b6118dc816117cb565b81146118e757600080fd5b50565b6118f3816117fd565b81146118fe57600080fd5b5056fea2646970667358221220609c0dba85766835e03db62b6b3126c773f1f033fe41eb56ed4ea5b8a2243ca664736f6c63430008070033",
}

// ERC20CustodyABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20CustodyMetaData.ABI instead.
var ERC20CustodyABI = ERC20CustodyMetaData.ABI

// ERC20CustodyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC20CustodyMetaData.Bin instead.
var ERC20CustodyBin = ERC20CustodyMetaData.Bin

// DeployERC20Custody deploys a new Ethereum contract, binding an instance of ERC20Custody to it.
func DeployERC20Custody(auth *bind.TransactOpts, backend bind.ContractBackend, TSSAddress_ common.Address, TSSAddressUpdater_ common.Address, zetaFee_ *big.Int, zetaMaxFee_ *big.Int, zeta_ common.Address) (common.Address, *types.Transaction, *ERC20Custody, error) {
	parsed, err := ERC20CustodyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC20CustodyBin), backend, TSSAddress_, TSSAddressUpdater_, zetaFee_, zetaMaxFee_, zeta_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20Custody{ERC20CustodyCaller: ERC20CustodyCaller{contract: contract}, ERC20CustodyTransactor: ERC20CustodyTransactor{contract: contract}, ERC20CustodyFilterer: ERC20CustodyFilterer{contract: contract}}, nil
}

// ERC20Custody is an auto generated Go binding around an Ethereum contract.
type ERC20Custody struct {
	ERC20CustodyCaller     // Read-only binding to the contract
	ERC20CustodyTransactor // Write-only binding to the contract
	ERC20CustodyFilterer   // Log filterer for contract events
}

// ERC20CustodyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20CustodyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20CustodyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20CustodyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20CustodyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20CustodyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20CustodySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20CustodySession struct {
	Contract     *ERC20Custody     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20CustodyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20CustodyCallerSession struct {
	Contract *ERC20CustodyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ERC20CustodyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20CustodyTransactorSession struct {
	Contract     *ERC20CustodyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ERC20CustodyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20CustodyRaw struct {
	Contract *ERC20Custody // Generic contract binding to access the raw methods on
}

// ERC20CustodyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20CustodyCallerRaw struct {
	Contract *ERC20CustodyCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20CustodyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20CustodyTransactorRaw struct {
	Contract *ERC20CustodyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Custody creates a new instance of ERC20Custody, bound to a specific deployed contract.
func NewERC20Custody(address common.Address, backend bind.ContractBackend) (*ERC20Custody, error) {
	contract, err := bindERC20Custody(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Custody{ERC20CustodyCaller: ERC20CustodyCaller{contract: contract}, ERC20CustodyTransactor: ERC20CustodyTransactor{contract: contract}, ERC20CustodyFilterer: ERC20CustodyFilterer{contract: contract}}, nil
}

// NewERC20CustodyCaller creates a new read-only instance of ERC20Custody, bound to a specific deployed contract.
func NewERC20CustodyCaller(address common.Address, caller bind.ContractCaller) (*ERC20CustodyCaller, error) {
	contract, err := bindERC20Custody(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyCaller{contract: contract}, nil
}

// NewERC20CustodyTransactor creates a new write-only instance of ERC20Custody, bound to a specific deployed contract.
func NewERC20CustodyTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20CustodyTransactor, error) {
	contract, err := bindERC20Custody(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyTransactor{contract: contract}, nil
}

// NewERC20CustodyFilterer creates a new log filterer instance of ERC20Custody, bound to a specific deployed contract.
func NewERC20CustodyFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20CustodyFilterer, error) {
	contract, err := bindERC20Custody(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyFilterer{contract: contract}, nil
}

// bindERC20Custody binds a generic wrapper to an already deployed contract.
func bindERC20Custody(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20CustodyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Custody *ERC20CustodyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Custody.Contract.ERC20CustodyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Custody *ERC20CustodyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Custody.Contract.ERC20CustodyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Custody *ERC20CustodyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Custody.Contract.ERC20CustodyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Custody *ERC20CustodyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Custody.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Custody *ERC20CustodyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Custody.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Custody *ERC20CustodyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Custody.Contract.contract.Transact(opts, method, params...)
}

// TSSAddress is a free data retrieval call binding the contract method 0x53ee30a3.
//
// Solidity: function TSSAddress() view returns(address)
func (_ERC20Custody *ERC20CustodyCaller) TSSAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20Custody.contract.Call(opts, &out, "TSSAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TSSAddress is a free data retrieval call binding the contract method 0x53ee30a3.
//
// Solidity: function TSSAddress() view returns(address)
func (_ERC20Custody *ERC20CustodySession) TSSAddress() (common.Address, error) {
	return _ERC20Custody.Contract.TSSAddress(&_ERC20Custody.CallOpts)
}

// TSSAddress is a free data retrieval call binding the contract method 0x53ee30a3.
//
// Solidity: function TSSAddress() view returns(address)
func (_ERC20Custody *ERC20CustodyCallerSession) TSSAddress() (common.Address, error) {
	return _ERC20Custody.Contract.TSSAddress(&_ERC20Custody.CallOpts)
}

// TSSAddressUpdater is a free data retrieval call binding the contract method 0x54b61e81.
//
// Solidity: function TSSAddressUpdater() view returns(address)
func (_ERC20Custody *ERC20CustodyCaller) TSSAddressUpdater(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20Custody.contract.Call(opts, &out, "TSSAddressUpdater")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TSSAddressUpdater is a free data retrieval call binding the contract method 0x54b61e81.
//
// Solidity: function TSSAddressUpdater() view returns(address)
func (_ERC20Custody *ERC20CustodySession) TSSAddressUpdater() (common.Address, error) {
	return _ERC20Custody.Contract.TSSAddressUpdater(&_ERC20Custody.CallOpts)
}

// TSSAddressUpdater is a free data retrieval call binding the contract method 0x54b61e81.
//
// Solidity: function TSSAddressUpdater() view returns(address)
func (_ERC20Custody *ERC20CustodyCallerSession) TSSAddressUpdater() (common.Address, error) {
	return _ERC20Custody.Contract.TSSAddressUpdater(&_ERC20Custody.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ERC20Custody *ERC20CustodyCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ERC20Custody.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ERC20Custody *ERC20CustodySession) Paused() (bool, error) {
	return _ERC20Custody.Contract.Paused(&_ERC20Custody.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ERC20Custody *ERC20CustodyCallerSession) Paused() (bool, error) {
	return _ERC20Custody.Contract.Paused(&_ERC20Custody.CallOpts)
}

// Whitelisted is a free data retrieval call binding the contract method 0xd936547e.
//
// Solidity: function whitelisted(address ) view returns(bool)
func (_ERC20Custody *ERC20CustodyCaller) Whitelisted(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ERC20Custody.contract.Call(opts, &out, "whitelisted", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Whitelisted is a free data retrieval call binding the contract method 0xd936547e.
//
// Solidity: function whitelisted(address ) view returns(bool)
func (_ERC20Custody *ERC20CustodySession) Whitelisted(arg0 common.Address) (bool, error) {
	return _ERC20Custody.Contract.Whitelisted(&_ERC20Custody.CallOpts, arg0)
}

// Whitelisted is a free data retrieval call binding the contract method 0xd936547e.
//
// Solidity: function whitelisted(address ) view returns(bool)
func (_ERC20Custody *ERC20CustodyCallerSession) Whitelisted(arg0 common.Address) (bool, error) {
	return _ERC20Custody.Contract.Whitelisted(&_ERC20Custody.CallOpts, arg0)
}

// Zeta is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function zeta() view returns(address)
func (_ERC20Custody *ERC20CustodyCaller) Zeta(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20Custody.contract.Call(opts, &out, "zeta")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Zeta is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function zeta() view returns(address)
func (_ERC20Custody *ERC20CustodySession) Zeta() (common.Address, error) {
	return _ERC20Custody.Contract.Zeta(&_ERC20Custody.CallOpts)
}

// Zeta is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function zeta() view returns(address)
func (_ERC20Custody *ERC20CustodyCallerSession) Zeta() (common.Address, error) {
	return _ERC20Custody.Contract.Zeta(&_ERC20Custody.CallOpts)
}

// ZetaFee is a free data retrieval call binding the contract method 0xe5408cfa.
//
// Solidity: function zetaFee() view returns(uint256)
func (_ERC20Custody *ERC20CustodyCaller) ZetaFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Custody.contract.Call(opts, &out, "zetaFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ZetaFee is a free data retrieval call binding the contract method 0xe5408cfa.
//
// Solidity: function zetaFee() view returns(uint256)
func (_ERC20Custody *ERC20CustodySession) ZetaFee() (*big.Int, error) {
	return _ERC20Custody.Contract.ZetaFee(&_ERC20Custody.CallOpts)
}

// ZetaFee is a free data retrieval call binding the contract method 0xe5408cfa.
//
// Solidity: function zetaFee() view returns(uint256)
func (_ERC20Custody *ERC20CustodyCallerSession) ZetaFee() (*big.Int, error) {
	return _ERC20Custody.Contract.ZetaFee(&_ERC20Custody.CallOpts)
}

// ZetaMaxFee is a free data retrieval call binding the contract method 0x7bdaded3.
//
// Solidity: function zetaMaxFee() view returns(uint256)
func (_ERC20Custody *ERC20CustodyCaller) ZetaMaxFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Custody.contract.Call(opts, &out, "zetaMaxFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ZetaMaxFee is a free data retrieval call binding the contract method 0x7bdaded3.
//
// Solidity: function zetaMaxFee() view returns(uint256)
func (_ERC20Custody *ERC20CustodySession) ZetaMaxFee() (*big.Int, error) {
	return _ERC20Custody.Contract.ZetaMaxFee(&_ERC20Custody.CallOpts)
}

// ZetaMaxFee is a free data retrieval call binding the contract method 0x7bdaded3.
//
// Solidity: function zetaMaxFee() view returns(uint256)
func (_ERC20Custody *ERC20CustodyCallerSession) ZetaMaxFee() (*big.Int, error) {
	return _ERC20Custody.Contract.ZetaMaxFee(&_ERC20Custody.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xe609055e.
//
// Solidity: function deposit(bytes recipient, address asset, uint256 amount, bytes message) returns()
func (_ERC20Custody *ERC20CustodyTransactor) Deposit(opts *bind.TransactOpts, recipient []byte, asset common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ERC20Custody.contract.Transact(opts, "deposit", recipient, asset, amount, message)
}

// Deposit is a paid mutator transaction binding the contract method 0xe609055e.
//
// Solidity: function deposit(bytes recipient, address asset, uint256 amount, bytes message) returns()
func (_ERC20Custody *ERC20CustodySession) Deposit(recipient []byte, asset common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ERC20Custody.Contract.Deposit(&_ERC20Custody.TransactOpts, recipient, asset, amount, message)
}

// Deposit is a paid mutator transaction binding the contract method 0xe609055e.
//
// Solidity: function deposit(bytes recipient, address asset, uint256 amount, bytes message) returns()
func (_ERC20Custody *ERC20CustodyTransactorSession) Deposit(recipient []byte, asset common.Address, amount *big.Int, message []byte) (*types.Transaction, error) {
	return _ERC20Custody.Contract.Deposit(&_ERC20Custody.TransactOpts, recipient, asset, amount, message)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ERC20Custody *ERC20CustodyTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Custody.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ERC20Custody *ERC20CustodySession) Pause() (*types.Transaction, error) {
	return _ERC20Custody.Contract.Pause(&_ERC20Custody.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ERC20Custody *ERC20CustodyTransactorSession) Pause() (*types.Transaction, error) {
	return _ERC20Custody.Contract.Pause(&_ERC20Custody.TransactOpts)
}

// RenounceTSSAddressUpdater is a paid mutator transaction binding the contract method 0xed11692b.
//
// Solidity: function renounceTSSAddressUpdater() returns()
func (_ERC20Custody *ERC20CustodyTransactor) RenounceTSSAddressUpdater(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Custody.contract.Transact(opts, "renounceTSSAddressUpdater")
}

// RenounceTSSAddressUpdater is a paid mutator transaction binding the contract method 0xed11692b.
//
// Solidity: function renounceTSSAddressUpdater() returns()
func (_ERC20Custody *ERC20CustodySession) RenounceTSSAddressUpdater() (*types.Transaction, error) {
	return _ERC20Custody.Contract.RenounceTSSAddressUpdater(&_ERC20Custody.TransactOpts)
}

// RenounceTSSAddressUpdater is a paid mutator transaction binding the contract method 0xed11692b.
//
// Solidity: function renounceTSSAddressUpdater() returns()
func (_ERC20Custody *ERC20CustodyTransactorSession) RenounceTSSAddressUpdater() (*types.Transaction, error) {
	return _ERC20Custody.Contract.RenounceTSSAddressUpdater(&_ERC20Custody.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ERC20Custody *ERC20CustodyTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Custody.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ERC20Custody *ERC20CustodySession) Unpause() (*types.Transaction, error) {
	return _ERC20Custody.Contract.Unpause(&_ERC20Custody.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ERC20Custody *ERC20CustodyTransactorSession) Unpause() (*types.Transaction, error) {
	return _ERC20Custody.Contract.Unpause(&_ERC20Custody.TransactOpts)
}

// Unwhitelist is a paid mutator transaction binding the contract method 0x9a590427.
//
// Solidity: function unwhitelist(address asset) returns()
func (_ERC20Custody *ERC20CustodyTransactor) Unwhitelist(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _ERC20Custody.contract.Transact(opts, "unwhitelist", asset)
}

// Unwhitelist is a paid mutator transaction binding the contract method 0x9a590427.
//
// Solidity: function unwhitelist(address asset) returns()
func (_ERC20Custody *ERC20CustodySession) Unwhitelist(asset common.Address) (*types.Transaction, error) {
	return _ERC20Custody.Contract.Unwhitelist(&_ERC20Custody.TransactOpts, asset)
}

// Unwhitelist is a paid mutator transaction binding the contract method 0x9a590427.
//
// Solidity: function unwhitelist(address asset) returns()
func (_ERC20Custody *ERC20CustodyTransactorSession) Unwhitelist(asset common.Address) (*types.Transaction, error) {
	return _ERC20Custody.Contract.Unwhitelist(&_ERC20Custody.TransactOpts, asset)
}

// UpdateTSSAddress is a paid mutator transaction binding the contract method 0x950837aa.
//
// Solidity: function updateTSSAddress(address TSSAddress_) returns()
func (_ERC20Custody *ERC20CustodyTransactor) UpdateTSSAddress(opts *bind.TransactOpts, TSSAddress_ common.Address) (*types.Transaction, error) {
	return _ERC20Custody.contract.Transact(opts, "updateTSSAddress", TSSAddress_)
}

// UpdateTSSAddress is a paid mutator transaction binding the contract method 0x950837aa.
//
// Solidity: function updateTSSAddress(address TSSAddress_) returns()
func (_ERC20Custody *ERC20CustodySession) UpdateTSSAddress(TSSAddress_ common.Address) (*types.Transaction, error) {
	return _ERC20Custody.Contract.UpdateTSSAddress(&_ERC20Custody.TransactOpts, TSSAddress_)
}

// UpdateTSSAddress is a paid mutator transaction binding the contract method 0x950837aa.
//
// Solidity: function updateTSSAddress(address TSSAddress_) returns()
func (_ERC20Custody *ERC20CustodyTransactorSession) UpdateTSSAddress(TSSAddress_ common.Address) (*types.Transaction, error) {
	return _ERC20Custody.Contract.UpdateTSSAddress(&_ERC20Custody.TransactOpts, TSSAddress_)
}

// UpdateZetaFee is a paid mutator transaction binding the contract method 0xde2f6c5e.
//
// Solidity: function updateZetaFee(uint256 zetaFee_) returns()
func (_ERC20Custody *ERC20CustodyTransactor) UpdateZetaFee(opts *bind.TransactOpts, zetaFee_ *big.Int) (*types.Transaction, error) {
	return _ERC20Custody.contract.Transact(opts, "updateZetaFee", zetaFee_)
}

// UpdateZetaFee is a paid mutator transaction binding the contract method 0xde2f6c5e.
//
// Solidity: function updateZetaFee(uint256 zetaFee_) returns()
func (_ERC20Custody *ERC20CustodySession) UpdateZetaFee(zetaFee_ *big.Int) (*types.Transaction, error) {
	return _ERC20Custody.Contract.UpdateZetaFee(&_ERC20Custody.TransactOpts, zetaFee_)
}

// UpdateZetaFee is a paid mutator transaction binding the contract method 0xde2f6c5e.
//
// Solidity: function updateZetaFee(uint256 zetaFee_) returns()
func (_ERC20Custody *ERC20CustodyTransactorSession) UpdateZetaFee(zetaFee_ *big.Int) (*types.Transaction, error) {
	return _ERC20Custody.Contract.UpdateZetaFee(&_ERC20Custody.TransactOpts, zetaFee_)
}

// Whitelist is a paid mutator transaction binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address asset) returns()
func (_ERC20Custody *ERC20CustodyTransactor) Whitelist(opts *bind.TransactOpts, asset common.Address) (*types.Transaction, error) {
	return _ERC20Custody.contract.Transact(opts, "whitelist", asset)
}

// Whitelist is a paid mutator transaction binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address asset) returns()
func (_ERC20Custody *ERC20CustodySession) Whitelist(asset common.Address) (*types.Transaction, error) {
	return _ERC20Custody.Contract.Whitelist(&_ERC20Custody.TransactOpts, asset)
}

// Whitelist is a paid mutator transaction binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address asset) returns()
func (_ERC20Custody *ERC20CustodyTransactorSession) Whitelist(asset common.Address) (*types.Transaction, error) {
	return _ERC20Custody.Contract.Whitelist(&_ERC20Custody.TransactOpts, asset)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address recipient, address asset, uint256 amount) returns()
func (_ERC20Custody *ERC20CustodyTransactor) Withdraw(opts *bind.TransactOpts, recipient common.Address, asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Custody.contract.Transact(opts, "withdraw", recipient, asset, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address recipient, address asset, uint256 amount) returns()
func (_ERC20Custody *ERC20CustodySession) Withdraw(recipient common.Address, asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Custody.Contract.Withdraw(&_ERC20Custody.TransactOpts, recipient, asset, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address recipient, address asset, uint256 amount) returns()
func (_ERC20Custody *ERC20CustodyTransactorSession) Withdraw(recipient common.Address, asset common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Custody.Contract.Withdraw(&_ERC20Custody.TransactOpts, recipient, asset, amount)
}

// ERC20CustodyDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the ERC20Custody contract.
type ERC20CustodyDepositedIterator struct {
	Event *ERC20CustodyDeposited // Event containing the contract specifics and raw log

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
func (it *ERC20CustodyDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20CustodyDeposited)
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
		it.Event = new(ERC20CustodyDeposited)
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
func (it *ERC20CustodyDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20CustodyDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20CustodyDeposited represents a Deposited event raised by the ERC20Custody contract.
type ERC20CustodyDeposited struct {
	Recipient []byte
	Asset     common.Address
	Amount    *big.Int
	Message   []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x1dafa057cc5c3bccb5ad974129a2bccd3c74002d9dfd7062404ba9523b18d6ae.
//
// Solidity: event Deposited(bytes recipient, address indexed asset, uint256 amount, bytes message)
func (_ERC20Custody *ERC20CustodyFilterer) FilterDeposited(opts *bind.FilterOpts, asset []common.Address) (*ERC20CustodyDepositedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _ERC20Custody.contract.FilterLogs(opts, "Deposited", assetRule)
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyDepositedIterator{contract: _ERC20Custody.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x1dafa057cc5c3bccb5ad974129a2bccd3c74002d9dfd7062404ba9523b18d6ae.
//
// Solidity: event Deposited(bytes recipient, address indexed asset, uint256 amount, bytes message)
func (_ERC20Custody *ERC20CustodyFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *ERC20CustodyDeposited, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _ERC20Custody.contract.WatchLogs(opts, "Deposited", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20CustodyDeposited)
				if err := _ERC20Custody.contract.UnpackLog(event, "Deposited", log); err != nil {
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

// ParseDeposited is a log parse operation binding the contract event 0x1dafa057cc5c3bccb5ad974129a2bccd3c74002d9dfd7062404ba9523b18d6ae.
//
// Solidity: event Deposited(bytes recipient, address indexed asset, uint256 amount, bytes message)
func (_ERC20Custody *ERC20CustodyFilterer) ParseDeposited(log types.Log) (*ERC20CustodyDeposited, error) {
	event := new(ERC20CustodyDeposited)
	if err := _ERC20Custody.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20CustodyPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ERC20Custody contract.
type ERC20CustodyPausedIterator struct {
	Event *ERC20CustodyPaused // Event containing the contract specifics and raw log

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
func (it *ERC20CustodyPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20CustodyPaused)
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
		it.Event = new(ERC20CustodyPaused)
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
func (it *ERC20CustodyPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20CustodyPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20CustodyPaused represents a Paused event raised by the ERC20Custody contract.
type ERC20CustodyPaused struct {
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address sender)
func (_ERC20Custody *ERC20CustodyFilterer) FilterPaused(opts *bind.FilterOpts) (*ERC20CustodyPausedIterator, error) {

	logs, sub, err := _ERC20Custody.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyPausedIterator{contract: _ERC20Custody.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address sender)
func (_ERC20Custody *ERC20CustodyFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ERC20CustodyPaused) (event.Subscription, error) {

	logs, sub, err := _ERC20Custody.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20CustodyPaused)
				if err := _ERC20Custody.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address sender)
func (_ERC20Custody *ERC20CustodyFilterer) ParsePaused(log types.Log) (*ERC20CustodyPaused, error) {
	event := new(ERC20CustodyPaused)
	if err := _ERC20Custody.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20CustodyRenouncedTSSUpdaterIterator is returned from FilterRenouncedTSSUpdater and is used to iterate over the raw logs and unpacked data for RenouncedTSSUpdater events raised by the ERC20Custody contract.
type ERC20CustodyRenouncedTSSUpdaterIterator struct {
	Event *ERC20CustodyRenouncedTSSUpdater // Event containing the contract specifics and raw log

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
func (it *ERC20CustodyRenouncedTSSUpdaterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20CustodyRenouncedTSSUpdater)
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
		it.Event = new(ERC20CustodyRenouncedTSSUpdater)
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
func (it *ERC20CustodyRenouncedTSSUpdaterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20CustodyRenouncedTSSUpdaterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20CustodyRenouncedTSSUpdater represents a RenouncedTSSUpdater event raised by the ERC20Custody contract.
type ERC20CustodyRenouncedTSSUpdater struct {
	TSSAddressUpdater common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRenouncedTSSUpdater is a free log retrieval operation binding the contract event 0x39ac266bfc011581be62c138d96e4e8782812013bb66fffb4cd207f4bfee3777.
//
// Solidity: event RenouncedTSSUpdater(address TSSAddressUpdater_)
func (_ERC20Custody *ERC20CustodyFilterer) FilterRenouncedTSSUpdater(opts *bind.FilterOpts) (*ERC20CustodyRenouncedTSSUpdaterIterator, error) {

	logs, sub, err := _ERC20Custody.contract.FilterLogs(opts, "RenouncedTSSUpdater")
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyRenouncedTSSUpdaterIterator{contract: _ERC20Custody.contract, event: "RenouncedTSSUpdater", logs: logs, sub: sub}, nil
}

// WatchRenouncedTSSUpdater is a free log subscription operation binding the contract event 0x39ac266bfc011581be62c138d96e4e8782812013bb66fffb4cd207f4bfee3777.
//
// Solidity: event RenouncedTSSUpdater(address TSSAddressUpdater_)
func (_ERC20Custody *ERC20CustodyFilterer) WatchRenouncedTSSUpdater(opts *bind.WatchOpts, sink chan<- *ERC20CustodyRenouncedTSSUpdater) (event.Subscription, error) {

	logs, sub, err := _ERC20Custody.contract.WatchLogs(opts, "RenouncedTSSUpdater")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20CustodyRenouncedTSSUpdater)
				if err := _ERC20Custody.contract.UnpackLog(event, "RenouncedTSSUpdater", log); err != nil {
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

// ParseRenouncedTSSUpdater is a log parse operation binding the contract event 0x39ac266bfc011581be62c138d96e4e8782812013bb66fffb4cd207f4bfee3777.
//
// Solidity: event RenouncedTSSUpdater(address TSSAddressUpdater_)
func (_ERC20Custody *ERC20CustodyFilterer) ParseRenouncedTSSUpdater(log types.Log) (*ERC20CustodyRenouncedTSSUpdater, error) {
	event := new(ERC20CustodyRenouncedTSSUpdater)
	if err := _ERC20Custody.contract.UnpackLog(event, "RenouncedTSSUpdater", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20CustodyUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ERC20Custody contract.
type ERC20CustodyUnpausedIterator struct {
	Event *ERC20CustodyUnpaused // Event containing the contract specifics and raw log

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
func (it *ERC20CustodyUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20CustodyUnpaused)
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
		it.Event = new(ERC20CustodyUnpaused)
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
func (it *ERC20CustodyUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20CustodyUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20CustodyUnpaused represents a Unpaused event raised by the ERC20Custody contract.
type ERC20CustodyUnpaused struct {
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address sender)
func (_ERC20Custody *ERC20CustodyFilterer) FilterUnpaused(opts *bind.FilterOpts) (*ERC20CustodyUnpausedIterator, error) {

	logs, sub, err := _ERC20Custody.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyUnpausedIterator{contract: _ERC20Custody.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address sender)
func (_ERC20Custody *ERC20CustodyFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ERC20CustodyUnpaused) (event.Subscription, error) {

	logs, sub, err := _ERC20Custody.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20CustodyUnpaused)
				if err := _ERC20Custody.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address sender)
func (_ERC20Custody *ERC20CustodyFilterer) ParseUnpaused(log types.Log) (*ERC20CustodyUnpaused, error) {
	event := new(ERC20CustodyUnpaused)
	if err := _ERC20Custody.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20CustodyUnwhitelistedIterator is returned from FilterUnwhitelisted and is used to iterate over the raw logs and unpacked data for Unwhitelisted events raised by the ERC20Custody contract.
type ERC20CustodyUnwhitelistedIterator struct {
	Event *ERC20CustodyUnwhitelisted // Event containing the contract specifics and raw log

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
func (it *ERC20CustodyUnwhitelistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20CustodyUnwhitelisted)
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
		it.Event = new(ERC20CustodyUnwhitelisted)
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
func (it *ERC20CustodyUnwhitelistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20CustodyUnwhitelistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20CustodyUnwhitelisted represents a Unwhitelisted event raised by the ERC20Custody contract.
type ERC20CustodyUnwhitelisted struct {
	Asset common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterUnwhitelisted is a free log retrieval operation binding the contract event 0x51085ddf9ebdded84b76e829eb58c4078e4b5bdf97d9a94723f336039da46791.
//
// Solidity: event Unwhitelisted(address indexed asset)
func (_ERC20Custody *ERC20CustodyFilterer) FilterUnwhitelisted(opts *bind.FilterOpts, asset []common.Address) (*ERC20CustodyUnwhitelistedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _ERC20Custody.contract.FilterLogs(opts, "Unwhitelisted", assetRule)
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyUnwhitelistedIterator{contract: _ERC20Custody.contract, event: "Unwhitelisted", logs: logs, sub: sub}, nil
}

// WatchUnwhitelisted is a free log subscription operation binding the contract event 0x51085ddf9ebdded84b76e829eb58c4078e4b5bdf97d9a94723f336039da46791.
//
// Solidity: event Unwhitelisted(address indexed asset)
func (_ERC20Custody *ERC20CustodyFilterer) WatchUnwhitelisted(opts *bind.WatchOpts, sink chan<- *ERC20CustodyUnwhitelisted, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _ERC20Custody.contract.WatchLogs(opts, "Unwhitelisted", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20CustodyUnwhitelisted)
				if err := _ERC20Custody.contract.UnpackLog(event, "Unwhitelisted", log); err != nil {
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

// ParseUnwhitelisted is a log parse operation binding the contract event 0x51085ddf9ebdded84b76e829eb58c4078e4b5bdf97d9a94723f336039da46791.
//
// Solidity: event Unwhitelisted(address indexed asset)
func (_ERC20Custody *ERC20CustodyFilterer) ParseUnwhitelisted(log types.Log) (*ERC20CustodyUnwhitelisted, error) {
	event := new(ERC20CustodyUnwhitelisted)
	if err := _ERC20Custody.contract.UnpackLog(event, "Unwhitelisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20CustodyUpdatedTSSAddressIterator is returned from FilterUpdatedTSSAddress and is used to iterate over the raw logs and unpacked data for UpdatedTSSAddress events raised by the ERC20Custody contract.
type ERC20CustodyUpdatedTSSAddressIterator struct {
	Event *ERC20CustodyUpdatedTSSAddress // Event containing the contract specifics and raw log

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
func (it *ERC20CustodyUpdatedTSSAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20CustodyUpdatedTSSAddress)
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
		it.Event = new(ERC20CustodyUpdatedTSSAddress)
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
func (it *ERC20CustodyUpdatedTSSAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20CustodyUpdatedTSSAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20CustodyUpdatedTSSAddress represents a UpdatedTSSAddress event raised by the ERC20Custody contract.
type ERC20CustodyUpdatedTSSAddress struct {
	TSSAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpdatedTSSAddress is a free log retrieval operation binding the contract event 0xcd2958db8285a532edf298cbe1aa28ea3fb5ec82461253f9a8c1477924e29947.
//
// Solidity: event UpdatedTSSAddress(address TSSAddress_)
func (_ERC20Custody *ERC20CustodyFilterer) FilterUpdatedTSSAddress(opts *bind.FilterOpts) (*ERC20CustodyUpdatedTSSAddressIterator, error) {

	logs, sub, err := _ERC20Custody.contract.FilterLogs(opts, "UpdatedTSSAddress")
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyUpdatedTSSAddressIterator{contract: _ERC20Custody.contract, event: "UpdatedTSSAddress", logs: logs, sub: sub}, nil
}

// WatchUpdatedTSSAddress is a free log subscription operation binding the contract event 0xcd2958db8285a532edf298cbe1aa28ea3fb5ec82461253f9a8c1477924e29947.
//
// Solidity: event UpdatedTSSAddress(address TSSAddress_)
func (_ERC20Custody *ERC20CustodyFilterer) WatchUpdatedTSSAddress(opts *bind.WatchOpts, sink chan<- *ERC20CustodyUpdatedTSSAddress) (event.Subscription, error) {

	logs, sub, err := _ERC20Custody.contract.WatchLogs(opts, "UpdatedTSSAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20CustodyUpdatedTSSAddress)
				if err := _ERC20Custody.contract.UnpackLog(event, "UpdatedTSSAddress", log); err != nil {
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

// ParseUpdatedTSSAddress is a log parse operation binding the contract event 0xcd2958db8285a532edf298cbe1aa28ea3fb5ec82461253f9a8c1477924e29947.
//
// Solidity: event UpdatedTSSAddress(address TSSAddress_)
func (_ERC20Custody *ERC20CustodyFilterer) ParseUpdatedTSSAddress(log types.Log) (*ERC20CustodyUpdatedTSSAddress, error) {
	event := new(ERC20CustodyUpdatedTSSAddress)
	if err := _ERC20Custody.contract.UnpackLog(event, "UpdatedTSSAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20CustodyUpdatedZetaFeeIterator is returned from FilterUpdatedZetaFee and is used to iterate over the raw logs and unpacked data for UpdatedZetaFee events raised by the ERC20Custody contract.
type ERC20CustodyUpdatedZetaFeeIterator struct {
	Event *ERC20CustodyUpdatedZetaFee // Event containing the contract specifics and raw log

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
func (it *ERC20CustodyUpdatedZetaFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20CustodyUpdatedZetaFee)
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
		it.Event = new(ERC20CustodyUpdatedZetaFee)
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
func (it *ERC20CustodyUpdatedZetaFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20CustodyUpdatedZetaFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20CustodyUpdatedZetaFee represents a UpdatedZetaFee event raised by the ERC20Custody contract.
type ERC20CustodyUpdatedZetaFee struct {
	ZetaFee *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUpdatedZetaFee is a free log retrieval operation binding the contract event 0x6d2d8e313fbaf76898bb9fa55e4b52525e49c7d7182d0874f97bd9076e81d523.
//
// Solidity: event UpdatedZetaFee(uint256 zetaFee_)
func (_ERC20Custody *ERC20CustodyFilterer) FilterUpdatedZetaFee(opts *bind.FilterOpts) (*ERC20CustodyUpdatedZetaFeeIterator, error) {

	logs, sub, err := _ERC20Custody.contract.FilterLogs(opts, "UpdatedZetaFee")
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyUpdatedZetaFeeIterator{contract: _ERC20Custody.contract, event: "UpdatedZetaFee", logs: logs, sub: sub}, nil
}

// WatchUpdatedZetaFee is a free log subscription operation binding the contract event 0x6d2d8e313fbaf76898bb9fa55e4b52525e49c7d7182d0874f97bd9076e81d523.
//
// Solidity: event UpdatedZetaFee(uint256 zetaFee_)
func (_ERC20Custody *ERC20CustodyFilterer) WatchUpdatedZetaFee(opts *bind.WatchOpts, sink chan<- *ERC20CustodyUpdatedZetaFee) (event.Subscription, error) {

	logs, sub, err := _ERC20Custody.contract.WatchLogs(opts, "UpdatedZetaFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20CustodyUpdatedZetaFee)
				if err := _ERC20Custody.contract.UnpackLog(event, "UpdatedZetaFee", log); err != nil {
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

// ParseUpdatedZetaFee is a log parse operation binding the contract event 0x6d2d8e313fbaf76898bb9fa55e4b52525e49c7d7182d0874f97bd9076e81d523.
//
// Solidity: event UpdatedZetaFee(uint256 zetaFee_)
func (_ERC20Custody *ERC20CustodyFilterer) ParseUpdatedZetaFee(log types.Log) (*ERC20CustodyUpdatedZetaFee, error) {
	event := new(ERC20CustodyUpdatedZetaFee)
	if err := _ERC20Custody.contract.UnpackLog(event, "UpdatedZetaFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20CustodyWhitelistedIterator is returned from FilterWhitelisted and is used to iterate over the raw logs and unpacked data for Whitelisted events raised by the ERC20Custody contract.
type ERC20CustodyWhitelistedIterator struct {
	Event *ERC20CustodyWhitelisted // Event containing the contract specifics and raw log

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
func (it *ERC20CustodyWhitelistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20CustodyWhitelisted)
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
		it.Event = new(ERC20CustodyWhitelisted)
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
func (it *ERC20CustodyWhitelistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20CustodyWhitelistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20CustodyWhitelisted represents a Whitelisted event raised by the ERC20Custody contract.
type ERC20CustodyWhitelisted struct {
	Asset common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterWhitelisted is a free log retrieval operation binding the contract event 0xaab7954e9d246b167ef88aeddad35209ca2489d95a8aeb59e288d9b19fae5a54.
//
// Solidity: event Whitelisted(address indexed asset)
func (_ERC20Custody *ERC20CustodyFilterer) FilterWhitelisted(opts *bind.FilterOpts, asset []common.Address) (*ERC20CustodyWhitelistedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _ERC20Custody.contract.FilterLogs(opts, "Whitelisted", assetRule)
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyWhitelistedIterator{contract: _ERC20Custody.contract, event: "Whitelisted", logs: logs, sub: sub}, nil
}

// WatchWhitelisted is a free log subscription operation binding the contract event 0xaab7954e9d246b167ef88aeddad35209ca2489d95a8aeb59e288d9b19fae5a54.
//
// Solidity: event Whitelisted(address indexed asset)
func (_ERC20Custody *ERC20CustodyFilterer) WatchWhitelisted(opts *bind.WatchOpts, sink chan<- *ERC20CustodyWhitelisted, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _ERC20Custody.contract.WatchLogs(opts, "Whitelisted", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20CustodyWhitelisted)
				if err := _ERC20Custody.contract.UnpackLog(event, "Whitelisted", log); err != nil {
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

// ParseWhitelisted is a log parse operation binding the contract event 0xaab7954e9d246b167ef88aeddad35209ca2489d95a8aeb59e288d9b19fae5a54.
//
// Solidity: event Whitelisted(address indexed asset)
func (_ERC20Custody *ERC20CustodyFilterer) ParseWhitelisted(log types.Log) (*ERC20CustodyWhitelisted, error) {
	event := new(ERC20CustodyWhitelisted)
	if err := _ERC20Custody.contract.UnpackLog(event, "Whitelisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20CustodyWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the ERC20Custody contract.
type ERC20CustodyWithdrawnIterator struct {
	Event *ERC20CustodyWithdrawn // Event containing the contract specifics and raw log

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
func (it *ERC20CustodyWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20CustodyWithdrawn)
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
		it.Event = new(ERC20CustodyWithdrawn)
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
func (it *ERC20CustodyWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20CustodyWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20CustodyWithdrawn represents a Withdrawn event raised by the ERC20Custody contract.
type ERC20CustodyWithdrawn struct {
	Recipient common.Address
	Asset     common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address indexed recipient, address indexed asset, uint256 amount)
func (_ERC20Custody *ERC20CustodyFilterer) FilterWithdrawn(opts *bind.FilterOpts, recipient []common.Address, asset []common.Address) (*ERC20CustodyWithdrawnIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _ERC20Custody.contract.FilterLogs(opts, "Withdrawn", recipientRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &ERC20CustodyWithdrawnIterator{contract: _ERC20Custody.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address indexed recipient, address indexed asset, uint256 amount)
func (_ERC20Custody *ERC20CustodyFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *ERC20CustodyWithdrawn, recipient []common.Address, asset []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _ERC20Custody.contract.WatchLogs(opts, "Withdrawn", recipientRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20CustodyWithdrawn)
				if err := _ERC20Custody.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address indexed recipient, address indexed asset, uint256 amount)
func (_ERC20Custody *ERC20CustodyFilterer) ParseWithdrawn(log types.Log) (*ERC20CustodyWithdrawn, error) {
	event := new(ERC20CustodyWithdrawn)
	if err := _ERC20Custody.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
