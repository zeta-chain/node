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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wzeta\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorSrc\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorDst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MoveStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Unstake\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositWZETA\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllValidators\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"operatorAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"consensusPubKey\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"jailed\",\"type\":\"bool\"},{\"internalType\":\"enumBondStatus\",\"name\":\"bondStatus\",\"type\":\"uint8\"}],\"internalType\":\"structValidator[]\",\"name\":\"validators\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"}],\"name\":\"getShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validatorSrc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"validatorDst\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"moveStake\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stakeAndRevert\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stakeWithStateUpdate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"staker\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"completionTime\",\"type\":\"int64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"withdrawWZETA\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405260666000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060006003553480156200005757600080fd5b506040516200192d3803806200192d83398181016040528101906200007d919062000170565b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555033600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050620001a2565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600062000138826200010b565b9050919050565b6200014a816200012b565b81146200015657600080fd5b50565b6000815190506200016a816200013f565b92915050565b60006020828403121562000189576200018862000106565b5b6000620001998482850162000159565b91505092915050565b61177b80620001b26000396000f3fe6080604052600436106100955760003560e01c806390b8436f1161005957806390b8436f1461018a5780639a0fb673146101c7578063bca8f527146101f0578063d11a93d01461022d578063f3513a371461026a5761009c565b80630d1b3daf1461009e5780632c5d24ae146100db57806357c6ea3e146100e557806361bc221a146101225780636be8916c1461014d5761009c565b3661009c57005b005b3480156100aa57600080fd5b506100c560048036038101906100c09190610d08565b610295565b6040516100d29190610d7d565b60405180910390f35b6100e361033c565b005b3480156100f157600080fd5b5061010c60048036038101906101079190610dc4565b61041b565b6040516101199190610e4f565b60405180910390f35b34801561012e57600080fd5b50610137610520565b6040516101449190610d7d565b60405180910390f35b34801561015957600080fd5b50610174600480360381019061016f9190610dc4565b610526565b6040516101819190610e85565b60405180910390f35b34801561019657600080fd5b506101b160048036038101906101ac9190610dc4565b61065c565b6040516101be9190610e85565b60405180910390f35b3480156101d357600080fd5b506101ee60048036038101906101e99190610ea0565b610761565b005b3480156101fc57600080fd5b5061021760048036038101906102129190610dc4565b61084b565b6040516102249190610e85565b60405180910390f35b34801561023957600080fd5b50610254600480360381019061024f9190610ecd565b6109ad565b6040516102619190610e4f565b60405180910390f35b34801561027657600080fd5b5061027f610ab5565b60405161028c91906111a6565b60405180910390f35b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630d1b3daf84846040518363ffffffff1660e01b81526004016102f3929190611221565b602060405180830381865afa158015610310573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103349190611266565b905092915050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461039657600080fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d0e30db0346040518263ffffffff1660e01b81526004016000604051808303818588803b15801561040057600080fd5b505af1158015610414573d6000803e3d6000fd5b5050505050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461047757600080fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166357c6ea3e8585856040518463ffffffff1660e01b81526004016104d493929190611293565b6020604051808303816000875af11580156104f3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061051791906112fd565b90509392505050565b60035481565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461058257600080fd5b60016003546105919190611359565b60038190555060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166390b8436f8686866040518463ffffffff1660e01b81526004016105f793929190611293565b6020604051808303816000875af1158015610616573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061063a91906113db565b9050600160035461064b9190611359565b600381905550809150509392505050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106b857600080fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166390b8436f8585856040518463ffffffff1660e01b815260040161071593929190611293565b6020604051808303816000875af1158015610734573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061075891906113db565b90509392505050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146107bb57600080fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16632e1a7d4d826040518263ffffffff1660e01b81526004016108169190610d7d565b600060405180830381600087803b15801561083057600080fd5b505af1158015610844573d6000803e3d6000fd5b5050505050565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146108a757600080fd5b60016003546108b69190611359565b60038190555060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166390b8436f8585856040518463ffffffff1660e01b815260040161091993929190611293565b6020604051808303816000875af1158015610938573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061095c91906113db565b50600160035461096c9190611359565b6003819055506040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109a490611454565b60405180910390fd5b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610a0957600080fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d11a93d0868686866040518563ffffffff1660e01b8152600401610a689493929190611474565b6020604051808303816000875af1158015610a87573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610aab91906112fd565b9050949350505050565b606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f3513a376040518163ffffffff1660e01b8152600401600060405180830381865afa158015610b22573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250810190610b4b91906116fc565b905090565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610b8f82610b64565b9050919050565b610b9f81610b84565b8114610baa57600080fd5b50565b600081359050610bbc81610b96565b92915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610c1582610bcc565b810181811067ffffffffffffffff82111715610c3457610c33610bdd565b5b80604052505050565b6000610c47610b50565b9050610c538282610c0c565b919050565b600067ffffffffffffffff821115610c7357610c72610bdd565b5b610c7c82610bcc565b9050602081019050919050565b82818337600083830152505050565b6000610cab610ca684610c58565b610c3d565b905082815260208101848484011115610cc757610cc6610bc7565b5b610cd2848285610c89565b509392505050565b600082601f830112610cef57610cee610bc2565b5b8135610cff848260208601610c98565b91505092915050565b60008060408385031215610d1f57610d1e610b5a565b5b6000610d2d85828601610bad565b925050602083013567ffffffffffffffff811115610d4e57610d4d610b5f565b5b610d5a85828601610cda565b9150509250929050565b6000819050919050565b610d7781610d64565b82525050565b6000602082019050610d926000830184610d6e565b92915050565b610da181610d64565b8114610dac57600080fd5b50565b600081359050610dbe81610d98565b92915050565b600080600060608486031215610ddd57610ddc610b5a565b5b6000610deb86828701610bad565b935050602084013567ffffffffffffffff811115610e0c57610e0b610b5f565b5b610e1886828701610cda565b9250506040610e2986828701610daf565b9150509250925092565b60008160070b9050919050565b610e4981610e33565b82525050565b6000602082019050610e646000830184610e40565b92915050565b60008115159050919050565b610e7f81610e6a565b82525050565b6000602082019050610e9a6000830184610e76565b92915050565b600060208284031215610eb657610eb5610b5a565b5b6000610ec484828501610daf565b91505092915050565b60008060008060808587031215610ee757610ee6610b5a565b5b6000610ef587828801610bad565b945050602085013567ffffffffffffffff811115610f1657610f15610b5f565b5b610f2287828801610cda565b935050604085013567ffffffffffffffff811115610f4357610f42610b5f565b5b610f4f87828801610cda565b9250506060610f6087828801610daf565b91505092959194509250565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610fd2578082015181840152602081019050610fb7565b83811115610fe1576000848401525b50505050565b6000610ff282610f98565b610ffc8185610fa3565b935061100c818560208601610fb4565b61101581610bcc565b840191505092915050565b61102981610e6a565b82525050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6004811061106f5761106e61102f565b5b50565b60008190506110808261105e565b919050565b600061109082611072565b9050919050565b6110a081611085565b82525050565b600060808301600083015184820360008601526110c38282610fe7565b915050602083015184820360208601526110dd8282610fe7565b91505060408301516110f26040860182611020565b5060608301516111056060860182611097565b508091505092915050565b600061111c83836110a6565b905092915050565b6000602082019050919050565b600061113c82610f6c565b6111468185610f77565b93508360208202850161115885610f88565b8060005b8581101561119457848403895281516111758582611110565b945061118083611124565b925060208a0199505060018101905061115c565b50829750879550505050505092915050565b600060208201905081810360008301526111c08184611131565b905092915050565b6111d181610b84565b82525050565b600082825260208201905092915050565b60006111f382610f98565b6111fd81856111d7565b935061120d818560208601610fb4565b61121681610bcc565b840191505092915050565b600060408201905061123660008301856111c8565b818103602083015261124881846111e8565b90509392505050565b60008151905061126081610d98565b92915050565b60006020828403121561127c5761127b610b5a565b5b600061128a84828501611251565b91505092915050565b60006060820190506112a860008301866111c8565b81810360208301526112ba81856111e8565b90506112c96040830184610d6e565b949350505050565b6112da81610e33565b81146112e557600080fd5b50565b6000815190506112f7816112d1565b92915050565b60006020828403121561131357611312610b5a565b5b6000611321848285016112e8565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061136482610d64565b915061136f83610d64565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156113a4576113a361132a565b5b828201905092915050565b6113b881610e6a565b81146113c357600080fd5b50565b6000815190506113d5816113af565b92915050565b6000602082840312156113f1576113f0610b5a565b5b60006113ff848285016113c6565b91505092915050565b7f7465737472657665727400000000000000000000000000000000000000000000600082015250565b600061143e600a836111d7565b915061144982611408565b602082019050919050565b6000602082019050818103600083015261146d81611431565b9050919050565b600060808201905061148960008301876111c8565b818103602083015261149b81866111e8565b905081810360408301526114af81856111e8565b90506114be6060830184610d6e565b95945050505050565b600067ffffffffffffffff8211156114e2576114e1610bdd565b5b602082029050602081019050919050565b600080fd5b600080fd5b600080fd5b600061151561151084610c58565b610c3d565b90508281526020810184848401111561153157611530610bc7565b5b61153c848285610fb4565b509392505050565b600082601f83011261155957611558610bc2565b5b8151611569848260208601611502565b91505092915050565b6004811061157f57600080fd5b50565b60008151905061159181611572565b92915050565b6000608082840312156115ad576115ac6114f8565b5b6115b76080610c3d565b9050600082015167ffffffffffffffff8111156115d7576115d66114fd565b5b6115e384828501611544565b600083015250602082015167ffffffffffffffff811115611607576116066114fd565b5b61161384828501611544565b6020830152506040611627848285016113c6565b604083015250606061163b84828501611582565b60608301525092915050565b600061165a611655846114c7565b610c3d565b9050808382526020820190506020840283018581111561167d5761167c6114f3565b5b835b818110156116c457805167ffffffffffffffff8111156116a2576116a1610bc2565b5b8086016116af8982611597565b8552602085019450505060208101905061167f565b5050509392505050565b600082601f8301126116e3576116e2610bc2565b5b81516116f3848260208601611647565b91505092915050565b60006020828403121561171257611711610b5a565b5b600082015167ffffffffffffffff8111156117305761172f610b5f565b5b61173c848285016116ce565b9150509291505056fea26469706673582212201d48926082b046ce0a5b741a926eefe8bb0873c1d992ea4a6c4210a133dbaae764736f6c634300080a0033",
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

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_TestStaking *TestStakingCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestStaking.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_TestStaking *TestStakingSession) Counter() (*big.Int, error) {
	return _TestStaking.Contract.Counter(&_TestStaking.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_TestStaking *TestStakingCallerSession) Counter() (*big.Int, error) {
	return _TestStaking.Contract.Counter(&_TestStaking.CallOpts)
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

// StakeAndRevert is a paid mutator transaction binding the contract method 0xbca8f527.
//
// Solidity: function stakeAndRevert(address staker, string validator, uint256 amount) returns(bool)
func (_TestStaking *TestStakingTransactor) StakeAndRevert(opts *bind.TransactOpts, staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.contract.Transact(opts, "stakeAndRevert", staker, validator, amount)
}

// StakeAndRevert is a paid mutator transaction binding the contract method 0xbca8f527.
//
// Solidity: function stakeAndRevert(address staker, string validator, uint256 amount) returns(bool)
func (_TestStaking *TestStakingSession) StakeAndRevert(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.StakeAndRevert(&_TestStaking.TransactOpts, staker, validator, amount)
}

// StakeAndRevert is a paid mutator transaction binding the contract method 0xbca8f527.
//
// Solidity: function stakeAndRevert(address staker, string validator, uint256 amount) returns(bool)
func (_TestStaking *TestStakingTransactorSession) StakeAndRevert(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.StakeAndRevert(&_TestStaking.TransactOpts, staker, validator, amount)
}

// StakeWithStateUpdate is a paid mutator transaction binding the contract method 0x6be8916c.
//
// Solidity: function stakeWithStateUpdate(address staker, string validator, uint256 amount) returns(bool)
func (_TestStaking *TestStakingTransactor) StakeWithStateUpdate(opts *bind.TransactOpts, staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.contract.Transact(opts, "stakeWithStateUpdate", staker, validator, amount)
}

// StakeWithStateUpdate is a paid mutator transaction binding the contract method 0x6be8916c.
//
// Solidity: function stakeWithStateUpdate(address staker, string validator, uint256 amount) returns(bool)
func (_TestStaking *TestStakingSession) StakeWithStateUpdate(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.StakeWithStateUpdate(&_TestStaking.TransactOpts, staker, validator, amount)
}

// StakeWithStateUpdate is a paid mutator transaction binding the contract method 0x6be8916c.
//
// Solidity: function stakeWithStateUpdate(address staker, string validator, uint256 amount) returns(bool)
func (_TestStaking *TestStakingTransactorSession) StakeWithStateUpdate(staker common.Address, validator string, amount *big.Int) (*types.Transaction, error) {
	return _TestStaking.Contract.StakeWithStateUpdate(&_TestStaking.TransactOpts, staker, validator, amount)
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

// TestStakingMoveStakeIterator is returned from FilterMoveStake and is used to iterate over the raw logs and unpacked data for MoveStake events raised by the TestStaking contract.
type TestStakingMoveStakeIterator struct {
	Event *TestStakingMoveStake // Event containing the contract specifics and raw log

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
func (it *TestStakingMoveStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestStakingMoveStake)
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
		it.Event = new(TestStakingMoveStake)
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
func (it *TestStakingMoveStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestStakingMoveStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestStakingMoveStake represents a MoveStake event raised by the TestStaking contract.
type TestStakingMoveStake struct {
	Staker       common.Address
	ValidatorSrc common.Address
	ValidatorDst common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMoveStake is a free log retrieval operation binding the contract event 0x4dda2c731d442025256e6e47fbb109592bcd8baf3cf25996ebd09f1da7ec902b.
//
// Solidity: event MoveStake(address indexed staker, address indexed validatorSrc, address indexed validatorDst, uint256 amount)
func (_TestStaking *TestStakingFilterer) FilterMoveStake(opts *bind.FilterOpts, staker []common.Address, validatorSrc []common.Address, validatorDst []common.Address) (*TestStakingMoveStakeIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorSrcRule []interface{}
	for _, validatorSrcItem := range validatorSrc {
		validatorSrcRule = append(validatorSrcRule, validatorSrcItem)
	}
	var validatorDstRule []interface{}
	for _, validatorDstItem := range validatorDst {
		validatorDstRule = append(validatorDstRule, validatorDstItem)
	}

	logs, sub, err := _TestStaking.contract.FilterLogs(opts, "MoveStake", stakerRule, validatorSrcRule, validatorDstRule)
	if err != nil {
		return nil, err
	}
	return &TestStakingMoveStakeIterator{contract: _TestStaking.contract, event: "MoveStake", logs: logs, sub: sub}, nil
}

// WatchMoveStake is a free log subscription operation binding the contract event 0x4dda2c731d442025256e6e47fbb109592bcd8baf3cf25996ebd09f1da7ec902b.
//
// Solidity: event MoveStake(address indexed staker, address indexed validatorSrc, address indexed validatorDst, uint256 amount)
func (_TestStaking *TestStakingFilterer) WatchMoveStake(opts *bind.WatchOpts, sink chan<- *TestStakingMoveStake, staker []common.Address, validatorSrc []common.Address, validatorDst []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorSrcRule []interface{}
	for _, validatorSrcItem := range validatorSrc {
		validatorSrcRule = append(validatorSrcRule, validatorSrcItem)
	}
	var validatorDstRule []interface{}
	for _, validatorDstItem := range validatorDst {
		validatorDstRule = append(validatorDstRule, validatorDstItem)
	}

	logs, sub, err := _TestStaking.contract.WatchLogs(opts, "MoveStake", stakerRule, validatorSrcRule, validatorDstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestStakingMoveStake)
				if err := _TestStaking.contract.UnpackLog(event, "MoveStake", log); err != nil {
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

// ParseMoveStake is a log parse operation binding the contract event 0x4dda2c731d442025256e6e47fbb109592bcd8baf3cf25996ebd09f1da7ec902b.
//
// Solidity: event MoveStake(address indexed staker, address indexed validatorSrc, address indexed validatorDst, uint256 amount)
func (_TestStaking *TestStakingFilterer) ParseMoveStake(log types.Log) (*TestStakingMoveStake, error) {
	event := new(TestStakingMoveStake)
	if err := _TestStaking.contract.UnpackLog(event, "MoveStake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestStakingStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the TestStaking contract.
type TestStakingStakeIterator struct {
	Event *TestStakingStake // Event containing the contract specifics and raw log

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
func (it *TestStakingStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestStakingStake)
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
		it.Event = new(TestStakingStake)
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
func (it *TestStakingStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestStakingStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestStakingStake represents a Stake event raised by the TestStaking contract.
type TestStakingStake struct {
	Staker    common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStake is a free log retrieval operation binding the contract event 0x99039fcf0a98f484616c5196ee8b2ecfa971babf0b519848289ea4db381f85f7.
//
// Solidity: event Stake(address indexed staker, address indexed validator, uint256 amount)
func (_TestStaking *TestStakingFilterer) FilterStake(opts *bind.FilterOpts, staker []common.Address, validator []common.Address) (*TestStakingStakeIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _TestStaking.contract.FilterLogs(opts, "Stake", stakerRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &TestStakingStakeIterator{contract: _TestStaking.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0x99039fcf0a98f484616c5196ee8b2ecfa971babf0b519848289ea4db381f85f7.
//
// Solidity: event Stake(address indexed staker, address indexed validator, uint256 amount)
func (_TestStaking *TestStakingFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *TestStakingStake, staker []common.Address, validator []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _TestStaking.contract.WatchLogs(opts, "Stake", stakerRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestStakingStake)
				if err := _TestStaking.contract.UnpackLog(event, "Stake", log); err != nil {
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

// ParseStake is a log parse operation binding the contract event 0x99039fcf0a98f484616c5196ee8b2ecfa971babf0b519848289ea4db381f85f7.
//
// Solidity: event Stake(address indexed staker, address indexed validator, uint256 amount)
func (_TestStaking *TestStakingFilterer) ParseStake(log types.Log) (*TestStakingStake, error) {
	event := new(TestStakingStake)
	if err := _TestStaking.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestStakingUnstakeIterator is returned from FilterUnstake and is used to iterate over the raw logs and unpacked data for Unstake events raised by the TestStaking contract.
type TestStakingUnstakeIterator struct {
	Event *TestStakingUnstake // Event containing the contract specifics and raw log

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
func (it *TestStakingUnstakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestStakingUnstake)
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
		it.Event = new(TestStakingUnstake)
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
func (it *TestStakingUnstakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestStakingUnstakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestStakingUnstake represents a Unstake event raised by the TestStaking contract.
type TestStakingUnstake struct {
	Staker    common.Address
	Validator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnstake is a free log retrieval operation binding the contract event 0x390b1276974b9463e5d66ab10df69b6f3d7b930eb066a0e66df327edd2cc811c.
//
// Solidity: event Unstake(address indexed staker, address indexed validator, uint256 amount)
func (_TestStaking *TestStakingFilterer) FilterUnstake(opts *bind.FilterOpts, staker []common.Address, validator []common.Address) (*TestStakingUnstakeIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _TestStaking.contract.FilterLogs(opts, "Unstake", stakerRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return &TestStakingUnstakeIterator{contract: _TestStaking.contract, event: "Unstake", logs: logs, sub: sub}, nil
}

// WatchUnstake is a free log subscription operation binding the contract event 0x390b1276974b9463e5d66ab10df69b6f3d7b930eb066a0e66df327edd2cc811c.
//
// Solidity: event Unstake(address indexed staker, address indexed validator, uint256 amount)
func (_TestStaking *TestStakingFilterer) WatchUnstake(opts *bind.WatchOpts, sink chan<- *TestStakingUnstake, staker []common.Address, validator []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _TestStaking.contract.WatchLogs(opts, "Unstake", stakerRule, validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestStakingUnstake)
				if err := _TestStaking.contract.UnpackLog(event, "Unstake", log); err != nil {
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

// ParseUnstake is a log parse operation binding the contract event 0x390b1276974b9463e5d66ab10df69b6f3d7b930eb066a0e66df327edd2cc811c.
//
// Solidity: event Unstake(address indexed staker, address indexed validator, uint256 amount)
func (_TestStaking *TestStakingFilterer) ParseUnstake(log types.Log) (*TestStakingUnstake, error) {
	event := new(TestStakingUnstake)
	if err := _TestStaking.contract.UnpackLog(event, "Unstake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
