// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testzrc20

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

// TestZRC20MetaData contains all meta data concerning the TestZRC20 contract.
var TestZRC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainid_\",\"type\":\"uint256\"},{\"internalType\":\"enumCoinType\",\"name\":\"coinType_\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CallerIsNotFungibleModule\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasFeeTransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LowAllowance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LowBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroGasCoin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroGasPrice\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"UpdatedGasLimit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"protocolFlatFee\",\"type\":\"uint256\"}],\"name\":\"UpdatedProtocolFlatFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"systemContract\",\"type\":\"address\"}],\"name\":\"UpdatedSystemContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasfee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"protocolFlatFee\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CHAIN_ID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"COIN_TYPE\",\"outputs\":[{\"internalType\":\"enumCoinType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FUNGIBLE_MODULE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GAS_LIMIT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_FLAT_FEE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SYSTEM_CONTRACT_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newField\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newPublicField\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"updateGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newField_\",\"type\":\"uint256\"}],\"name\":\"updateNewField\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"protocolFlatFee\",\"type\":\"uint256\"}],\"name\":\"updateProtocolFlatFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"updateSystemContractAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawGasFee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162002566380380620025668339818101604052810190620000379190620000e1565b816080818152505080600281111562000055576200005462000128565b5b60a08160028111156200006d576200006c62000128565b5b81525050505062000157565b600080fd5b6000819050919050565b62000093816200007e565b81146200009f57600080fd5b50565b600081519050620000b38162000088565b92915050565b60038110620000c757600080fd5b50565b600081519050620000db81620000b9565b92915050565b60008060408385031215620000fb57620000fa62000079565b5b60006200010b85828601620000a2565b92505060206200011e85828601620000ca565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60805160a0516123db6200018b6000396000610aa40152600081816109ee01528181610f59015261107e01526123db6000f3fe608060405234801561001057600080fd5b50600436106101a95760003560e01c806385e1f4d0116100f9578063c701262611610097578063dd62ed3e11610071578063dd62ed3e146104ff578063eddeb1231461052f578063f2441b321461054b578063f687d12a14610569576101a9565b8063c701262614610494578063c835d7cc146104c4578063d9eeebed146104e0576101a9565b8063a457c2d7116100d3578063a457c2d7146103f8578063a7605f4514610428578063a9059cbb14610446578063b92894ba14610476576101a9565b806385e1f4d01461039e57806395d89b41146103bc578063a3413d03146103da576101a9565b8063395093511161016657806347e7ef241161014057806347e7ef24146103045780634d8943bb1461033457806370a0823114610352578063732bb0e414610382576101a9565b806339509351146102865780633ce4a5bc146102b657806342966c68146102d4576101a9565b806306fdde03146101ae578063091d2788146101cc578063095ea7b3146101ea57806318160ddd1461021a57806323b872dd14610238578063313ce56714610268575b600080fd5b6101b6610585565b6040516101c39190611b20565b60405180910390f35b6101d4610617565b6040516101e19190611b5b565b60405180910390f35b61020460048036038101906101ff9190611c14565b61061d565b6040516102119190611c6f565b60405180910390f35b61022261063b565b60405161022f9190611b5b565b60405180910390f35b610252600480360381019061024d9190611c8a565b610645565b60405161025f9190611c6f565b60405180910390f35b61027061073d565b60405161027d9190611cf9565b60405180910390f35b6102a0600480360381019061029b9190611c14565b610754565b6040516102ad9190611c6f565b60405180910390f35b6102be6107fa565b6040516102cb9190611d23565b60405180910390f35b6102ee60048036038101906102e99190611d3e565b610812565b6040516102fb9190611c6f565b60405180910390f35b61031e60048036038101906103199190611c14565b610827565b60405161032b9190611c6f565b60405180910390f35b61033c610993565b6040516103499190611b5b565b60405180910390f35b61036c60048036038101906103679190611d6b565b610999565b6040516103799190611b5b565b60405180910390f35b61039c60048036038101906103979190611d3e565b6109e2565b005b6103a66109ec565b6040516103b39190611b5b565b60405180910390f35b6103c4610a10565b6040516103d19190611b20565b60405180910390f35b6103e2610aa2565b6040516103ef9190611e0f565b60405180910390f35b610412600480360381019061040d9190611c14565b610ac6565b60405161041f9190611c6f565b60405180910390f35b610430610c29565b60405161043d9190611b5b565b60405180910390f35b610460600480360381019061045b9190611c14565b610c2f565b60405161046d9190611c6f565b60405180910390f35b61047e610c4d565b60405161048b9190611b20565b60405180910390f35b6104ae60048036038101906104a99190611f5f565b610cdb565b6040516104bb9190611c6f565b60405180910390f35b6104de60048036038101906104d99190611d6b565b610e22565b005b6104e8610f15565b6040516104f6929190611fbb565b60405180910390f35b61051960048036038101906105149190611fe4565b611162565b6040516105269190611b5b565b60405180910390f35b61054960048036038101906105449190611d3e565b6111e9565b005b6105536112a3565b6040516105609190611d23565b60405180910390f35b610583600480360381019061057e9190611d3e565b6112c7565b005b60606006805461059490612053565b80601f01602080910402602001604051908101604052809291908181526020018280546105c090612053565b801561060d5780601f106105e25761010080835404028352916020019161060d565b820191906000526020600020905b8154815290600101906020018083116105f057829003601f168201915b5050505050905090565b60015481565b600061063161062a611381565b8484611389565b6001905092915050565b6000600554905090565b6000610652848484611540565b6000600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600061069d611381565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905082811015610714576040517f10bad14700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61073185610720611381565b858461072c91906120b3565b611389565b60019150509392505050565b6000600860009054906101000a900460ff16905090565b600081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006107a0611381565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546107e991906120e7565b925050819055506001905092915050565b73735b14bb79463307aacbed86daf3322b1e6226ab81565b600061081e338361179a565b60019050919050565b600073735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141580156108c5575060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b156108fc576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6109068383611951565b8273ffffffffffffffffffffffffffffffffffffffff167f67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab373735b14bb79463307aacbed86daf3322b1e6226ab6040516020016109639190612163565b604051602081830303815290604052846040516109819291906121d3565b60405180910390a26001905092915050565b60025481565b6000600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b8060098190555050565b7f000000000000000000000000000000000000000000000000000000000000000081565b606060078054610a1f90612053565b80601f0160208091040260200160405190810160405280929190818152602001828054610a4b90612053565b8015610a985780601f10610a6d57610100808354040283529160200191610a98565b820191906000526020600020905b815481529060010190602001808311610a7b57829003601f168201915b5050505050905090565b7f000000000000000000000000000000000000000000000000000000000000000081565b600081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000610b12611381565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015610b85576040517f10bad14700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000610bcf611381565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610c1891906120b3565b925050819055506001905092915050565b60095481565b6000610c43610c3c611381565b8484611540565b6001905092915050565b600a8054610c5a90612053565b80601f0160208091040260200160405190810160405280929190818152602001828054610c8690612053565b8015610cd35780601f10610ca857610100808354040283529160200191610cd3565b820191906000526020600020905b815481529060010190602001808311610cb657829003601f168201915b505050505081565b6000806000610ce8610f15565b915091508173ffffffffffffffffffffffffffffffffffffffff166323b872dd3373735b14bb79463307aacbed86daf3322b1e6226ab846040518463ffffffff1660e01b8152600401610d3d93929190612203565b6020604051808303816000875af1158015610d5c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d809190612266565b610db6576040517f0a7cd6d600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610dc0338561179a565b3373ffffffffffffffffffffffffffffffffffffffff167f9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955868684600254604051610e0e9493929190612293565b60405180910390a260019250505092915050565b73735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610e9b576040517f2b2add3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507fd55614e962c5fd6ece71614f6348d702468a997a394dd5e5c1677950226d97ae81604051610f0a9190611d23565b60405180910390a150565b60008060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630be155477f00000000000000000000000000000000000000000000000000000000000000006040518263ffffffff1660e01b8152600401610f949190611b5b565b602060405180830381865afa158015610fb1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fd591906122f4565b9050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361103d576040517f78fff39600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d7fd7afb7f00000000000000000000000000000000000000000000000000000000000000006040518263ffffffff1660e01b81526004016110b99190611b5b565b602060405180830381865afa1580156110d6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110fa9190612336565b905060008103611136576040517fe661aed000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600254600154836111499190612363565b61115391906120e7565b90508281945094505050509091565b6000600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b73735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611262576040517f2b2add3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806002819055507fef13af88e424b5d15f49c77758542c1938b08b8b95b91ed0751f98ba99000d8f816040516112989190611b5b565b60405180910390a150565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b73735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611340576040517f2b2add3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806001819055507fff5788270f43bfc1ca41c503606d2594aa3023a1a7547de403a3e2f146a4a80a816040516113769190611b5b565b60405180910390a150565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036113ef576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611455576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040516115339190611b5b565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036115a6576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361160c576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561168a576040517ffe382aa700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818161169691906120b3565b600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555081600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461172891906120e7565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161178c9190611b5b565b60405180910390a350505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603611800576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561187e576040517ffe382aa700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818161188a91906120b3565b600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555081600560008282546118df91906120b3565b92505081905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516119449190611b5b565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036119b7576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600560008282546119c991906120e7565b9250508190555080600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254611a1f91906120e7565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051611a849190611b5b565b60405180910390a35050565b600081519050919050565b600082825260208201905092915050565b60005b83811015611aca578082015181840152602081019050611aaf565b60008484015250505050565b6000601f19601f8301169050919050565b6000611af282611a90565b611afc8185611a9b565b9350611b0c818560208601611aac565b611b1581611ad6565b840191505092915050565b60006020820190508181036000830152611b3a8184611ae7565b905092915050565b6000819050919050565b611b5581611b42565b82525050565b6000602082019050611b706000830184611b4c565b92915050565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000611bb582611b8a565b9050919050565b611bc581611baa565b8114611bd057600080fd5b50565b600081359050611be281611bbc565b92915050565b611bf181611b42565b8114611bfc57600080fd5b50565b600081359050611c0e81611be8565b92915050565b60008060408385031215611c2b57611c2a611b80565b5b6000611c3985828601611bd3565b9250506020611c4a85828601611bff565b9150509250929050565b60008115159050919050565b611c6981611c54565b82525050565b6000602082019050611c846000830184611c60565b92915050565b600080600060608486031215611ca357611ca2611b80565b5b6000611cb186828701611bd3565b9350506020611cc286828701611bd3565b9250506040611cd386828701611bff565b9150509250925092565b600060ff82169050919050565b611cf381611cdd565b82525050565b6000602082019050611d0e6000830184611cea565b92915050565b611d1d81611baa565b82525050565b6000602082019050611d386000830184611d14565b92915050565b600060208284031215611d5457611d53611b80565b5b6000611d6284828501611bff565b91505092915050565b600060208284031215611d8157611d80611b80565b5b6000611d8f84828501611bd3565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60038110611dd857611dd7611d98565b5b50565b6000819050611de982611dc7565b919050565b6000611df982611ddb565b9050919050565b611e0981611dee565b82525050565b6000602082019050611e246000830184611e00565b92915050565b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b611e6c82611ad6565b810181811067ffffffffffffffff82111715611e8b57611e8a611e34565b5b80604052505050565b6000611e9e611b76565b9050611eaa8282611e63565b919050565b600067ffffffffffffffff821115611eca57611ec9611e34565b5b611ed382611ad6565b9050602081019050919050565b82818337600083830152505050565b6000611f02611efd84611eaf565b611e94565b905082815260208101848484011115611f1e57611f1d611e2f565b5b611f29848285611ee0565b509392505050565b600082601f830112611f4657611f45611e2a565b5b8135611f56848260208601611eef565b91505092915050565b60008060408385031215611f7657611f75611b80565b5b600083013567ffffffffffffffff811115611f9457611f93611b85565b5b611fa085828601611f31565b9250506020611fb185828601611bff565b9150509250929050565b6000604082019050611fd06000830185611d14565b611fdd6020830184611b4c565b9392505050565b60008060408385031215611ffb57611ffa611b80565b5b600061200985828601611bd3565b925050602061201a85828601611bd3565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061206b57607f821691505b60208210810361207e5761207d612024565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006120be82611b42565b91506120c983611b42565b92508282039050818111156120e1576120e0612084565b5b92915050565b60006120f282611b42565b91506120fd83611b42565b925082820190508082111561211557612114612084565b5b92915050565b60008160601b9050919050565b60006121338261211b565b9050919050565b600061214582612128565b9050919050565b61215d61215882611baa565b61213a565b82525050565b600061216f828461214c565b60148201915081905092915050565b600081519050919050565b600082825260208201905092915050565b60006121a58261217e565b6121af8185612189565b93506121bf818560208601611aac565b6121c881611ad6565b840191505092915050565b600060408201905081810360008301526121ed818561219a565b90506121fc6020830184611b4c565b9392505050565b60006060820190506122186000830186611d14565b6122256020830185611d14565b6122326040830184611b4c565b949350505050565b61224381611c54565b811461224e57600080fd5b50565b6000815190506122608161223a565b92915050565b60006020828403121561227c5761227b611b80565b5b600061228a84828501612251565b91505092915050565b600060808201905081810360008301526122ad818761219a565b90506122bc6020830186611b4c565b6122c96040830185611b4c565b6122d66060830184611b4c565b95945050505050565b6000815190506122ee81611bbc565b92915050565b60006020828403121561230a57612309611b80565b5b6000612318848285016122df565b91505092915050565b60008151905061233081611be8565b92915050565b60006020828403121561234c5761234b611b80565b5b600061235a84828501612321565b91505092915050565b600061236e82611b42565b915061237983611b42565b925082820261238781611b42565b9150828204841483151761239e5761239d612084565b5b509291505056fea264697066735822122011fda4fba218fc7f911b68b0418884c426f5ddb324eb1572e2e793fcb60edf6c64736f6c63430008150033",
}

// TestZRC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use TestZRC20MetaData.ABI instead.
var TestZRC20ABI = TestZRC20MetaData.ABI

// TestZRC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestZRC20MetaData.Bin instead.
var TestZRC20Bin = TestZRC20MetaData.Bin

// DeployTestZRC20 deploys a new Ethereum contract, binding an instance of TestZRC20 to it.
func DeployTestZRC20(auth *bind.TransactOpts, backend bind.ContractBackend, chainid_ *big.Int, coinType_ uint8) (common.Address, *types.Transaction, *TestZRC20, error) {
	parsed, err := TestZRC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestZRC20Bin), backend, chainid_, coinType_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestZRC20{TestZRC20Caller: TestZRC20Caller{contract: contract}, TestZRC20Transactor: TestZRC20Transactor{contract: contract}, TestZRC20Filterer: TestZRC20Filterer{contract: contract}}, nil
}

// TestZRC20 is an auto generated Go binding around an Ethereum contract.
type TestZRC20 struct {
	TestZRC20Caller     // Read-only binding to the contract
	TestZRC20Transactor // Write-only binding to the contract
	TestZRC20Filterer   // Log filterer for contract events
}

// TestZRC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type TestZRC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestZRC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TestZRC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestZRC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestZRC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestZRC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestZRC20Session struct {
	Contract     *TestZRC20        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestZRC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestZRC20CallerSession struct {
	Contract *TestZRC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TestZRC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestZRC20TransactorSession struct {
	Contract     *TestZRC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TestZRC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type TestZRC20Raw struct {
	Contract *TestZRC20 // Generic contract binding to access the raw methods on
}

// TestZRC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestZRC20CallerRaw struct {
	Contract *TestZRC20Caller // Generic read-only contract binding to access the raw methods on
}

// TestZRC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestZRC20TransactorRaw struct {
	Contract *TestZRC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTestZRC20 creates a new instance of TestZRC20, bound to a specific deployed contract.
func NewTestZRC20(address common.Address, backend bind.ContractBackend) (*TestZRC20, error) {
	contract, err := bindTestZRC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestZRC20{TestZRC20Caller: TestZRC20Caller{contract: contract}, TestZRC20Transactor: TestZRC20Transactor{contract: contract}, TestZRC20Filterer: TestZRC20Filterer{contract: contract}}, nil
}

// NewTestZRC20Caller creates a new read-only instance of TestZRC20, bound to a specific deployed contract.
func NewTestZRC20Caller(address common.Address, caller bind.ContractCaller) (*TestZRC20Caller, error) {
	contract, err := bindTestZRC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestZRC20Caller{contract: contract}, nil
}

// NewTestZRC20Transactor creates a new write-only instance of TestZRC20, bound to a specific deployed contract.
func NewTestZRC20Transactor(address common.Address, transactor bind.ContractTransactor) (*TestZRC20Transactor, error) {
	contract, err := bindTestZRC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestZRC20Transactor{contract: contract}, nil
}

// NewTestZRC20Filterer creates a new log filterer instance of TestZRC20, bound to a specific deployed contract.
func NewTestZRC20Filterer(address common.Address, filterer bind.ContractFilterer) (*TestZRC20Filterer, error) {
	contract, err := bindTestZRC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestZRC20Filterer{contract: contract}, nil
}

// bindTestZRC20 binds a generic wrapper to an already deployed contract.
func bindTestZRC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestZRC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestZRC20 *TestZRC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestZRC20.Contract.TestZRC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestZRC20 *TestZRC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestZRC20.Contract.TestZRC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestZRC20 *TestZRC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestZRC20.Contract.TestZRC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestZRC20 *TestZRC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestZRC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestZRC20 *TestZRC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestZRC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestZRC20 *TestZRC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestZRC20.Contract.contract.Transact(opts, method, params...)
}

// CHAINID is a free data retrieval call binding the contract method 0x85e1f4d0.
//
// Solidity: function CHAIN_ID() view returns(uint256)
func (_TestZRC20 *TestZRC20Caller) CHAINID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "CHAIN_ID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CHAINID is a free data retrieval call binding the contract method 0x85e1f4d0.
//
// Solidity: function CHAIN_ID() view returns(uint256)
func (_TestZRC20 *TestZRC20Session) CHAINID() (*big.Int, error) {
	return _TestZRC20.Contract.CHAINID(&_TestZRC20.CallOpts)
}

// CHAINID is a free data retrieval call binding the contract method 0x85e1f4d0.
//
// Solidity: function CHAIN_ID() view returns(uint256)
func (_TestZRC20 *TestZRC20CallerSession) CHAINID() (*big.Int, error) {
	return _TestZRC20.Contract.CHAINID(&_TestZRC20.CallOpts)
}

// COINTYPE is a free data retrieval call binding the contract method 0xa3413d03.
//
// Solidity: function COIN_TYPE() view returns(uint8)
func (_TestZRC20 *TestZRC20Caller) COINTYPE(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "COIN_TYPE")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// COINTYPE is a free data retrieval call binding the contract method 0xa3413d03.
//
// Solidity: function COIN_TYPE() view returns(uint8)
func (_TestZRC20 *TestZRC20Session) COINTYPE() (uint8, error) {
	return _TestZRC20.Contract.COINTYPE(&_TestZRC20.CallOpts)
}

// COINTYPE is a free data retrieval call binding the contract method 0xa3413d03.
//
// Solidity: function COIN_TYPE() view returns(uint8)
func (_TestZRC20 *TestZRC20CallerSession) COINTYPE() (uint8, error) {
	return _TestZRC20.Contract.COINTYPE(&_TestZRC20.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_TestZRC20 *TestZRC20Caller) FUNGIBLEMODULEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "FUNGIBLE_MODULE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_TestZRC20 *TestZRC20Session) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _TestZRC20.Contract.FUNGIBLEMODULEADDRESS(&_TestZRC20.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_TestZRC20 *TestZRC20CallerSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _TestZRC20.Contract.FUNGIBLEMODULEADDRESS(&_TestZRC20.CallOpts)
}

// GASLIMIT is a free data retrieval call binding the contract method 0x091d2788.
//
// Solidity: function GAS_LIMIT() view returns(uint256)
func (_TestZRC20 *TestZRC20Caller) GASLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "GAS_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GASLIMIT is a free data retrieval call binding the contract method 0x091d2788.
//
// Solidity: function GAS_LIMIT() view returns(uint256)
func (_TestZRC20 *TestZRC20Session) GASLIMIT() (*big.Int, error) {
	return _TestZRC20.Contract.GASLIMIT(&_TestZRC20.CallOpts)
}

// GASLIMIT is a free data retrieval call binding the contract method 0x091d2788.
//
// Solidity: function GAS_LIMIT() view returns(uint256)
func (_TestZRC20 *TestZRC20CallerSession) GASLIMIT() (*big.Int, error) {
	return _TestZRC20.Contract.GASLIMIT(&_TestZRC20.CallOpts)
}

// PROTOCOLFLATFEE is a free data retrieval call binding the contract method 0x4d8943bb.
//
// Solidity: function PROTOCOL_FLAT_FEE() view returns(uint256)
func (_TestZRC20 *TestZRC20Caller) PROTOCOLFLATFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "PROTOCOL_FLAT_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PROTOCOLFLATFEE is a free data retrieval call binding the contract method 0x4d8943bb.
//
// Solidity: function PROTOCOL_FLAT_FEE() view returns(uint256)
func (_TestZRC20 *TestZRC20Session) PROTOCOLFLATFEE() (*big.Int, error) {
	return _TestZRC20.Contract.PROTOCOLFLATFEE(&_TestZRC20.CallOpts)
}

// PROTOCOLFLATFEE is a free data retrieval call binding the contract method 0x4d8943bb.
//
// Solidity: function PROTOCOL_FLAT_FEE() view returns(uint256)
func (_TestZRC20 *TestZRC20CallerSession) PROTOCOLFLATFEE() (*big.Int, error) {
	return _TestZRC20.Contract.PROTOCOLFLATFEE(&_TestZRC20.CallOpts)
}

// SYSTEMCONTRACTADDRESS is a free data retrieval call binding the contract method 0xf2441b32.
//
// Solidity: function SYSTEM_CONTRACT_ADDRESS() view returns(address)
func (_TestZRC20 *TestZRC20Caller) SYSTEMCONTRACTADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "SYSTEM_CONTRACT_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SYSTEMCONTRACTADDRESS is a free data retrieval call binding the contract method 0xf2441b32.
//
// Solidity: function SYSTEM_CONTRACT_ADDRESS() view returns(address)
func (_TestZRC20 *TestZRC20Session) SYSTEMCONTRACTADDRESS() (common.Address, error) {
	return _TestZRC20.Contract.SYSTEMCONTRACTADDRESS(&_TestZRC20.CallOpts)
}

// SYSTEMCONTRACTADDRESS is a free data retrieval call binding the contract method 0xf2441b32.
//
// Solidity: function SYSTEM_CONTRACT_ADDRESS() view returns(address)
func (_TestZRC20 *TestZRC20CallerSession) SYSTEMCONTRACTADDRESS() (common.Address, error) {
	return _TestZRC20.Contract.SYSTEMCONTRACTADDRESS(&_TestZRC20.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestZRC20 *TestZRC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestZRC20 *TestZRC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TestZRC20.Contract.Allowance(&_TestZRC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestZRC20 *TestZRC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TestZRC20.Contract.Allowance(&_TestZRC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestZRC20 *TestZRC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestZRC20 *TestZRC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _TestZRC20.Contract.BalanceOf(&_TestZRC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestZRC20 *TestZRC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TestZRC20.Contract.BalanceOf(&_TestZRC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestZRC20 *TestZRC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestZRC20 *TestZRC20Session) Decimals() (uint8, error) {
	return _TestZRC20.Contract.Decimals(&_TestZRC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestZRC20 *TestZRC20CallerSession) Decimals() (uint8, error) {
	return _TestZRC20.Contract.Decimals(&_TestZRC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestZRC20 *TestZRC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestZRC20 *TestZRC20Session) Name() (string, error) {
	return _TestZRC20.Contract.Name(&_TestZRC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestZRC20 *TestZRC20CallerSession) Name() (string, error) {
	return _TestZRC20.Contract.Name(&_TestZRC20.CallOpts)
}

// NewField is a free data retrieval call binding the contract method 0xa7605f45.
//
// Solidity: function newField() view returns(uint256)
func (_TestZRC20 *TestZRC20Caller) NewField(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "newField")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NewField is a free data retrieval call binding the contract method 0xa7605f45.
//
// Solidity: function newField() view returns(uint256)
func (_TestZRC20 *TestZRC20Session) NewField() (*big.Int, error) {
	return _TestZRC20.Contract.NewField(&_TestZRC20.CallOpts)
}

// NewField is a free data retrieval call binding the contract method 0xa7605f45.
//
// Solidity: function newField() view returns(uint256)
func (_TestZRC20 *TestZRC20CallerSession) NewField() (*big.Int, error) {
	return _TestZRC20.Contract.NewField(&_TestZRC20.CallOpts)
}

// NewPublicField is a free data retrieval call binding the contract method 0xb92894ba.
//
// Solidity: function newPublicField() view returns(string)
func (_TestZRC20 *TestZRC20Caller) NewPublicField(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "newPublicField")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NewPublicField is a free data retrieval call binding the contract method 0xb92894ba.
//
// Solidity: function newPublicField() view returns(string)
func (_TestZRC20 *TestZRC20Session) NewPublicField() (string, error) {
	return _TestZRC20.Contract.NewPublicField(&_TestZRC20.CallOpts)
}

// NewPublicField is a free data retrieval call binding the contract method 0xb92894ba.
//
// Solidity: function newPublicField() view returns(string)
func (_TestZRC20 *TestZRC20CallerSession) NewPublicField() (string, error) {
	return _TestZRC20.Contract.NewPublicField(&_TestZRC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestZRC20 *TestZRC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestZRC20 *TestZRC20Session) Symbol() (string, error) {
	return _TestZRC20.Contract.Symbol(&_TestZRC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestZRC20 *TestZRC20CallerSession) Symbol() (string, error) {
	return _TestZRC20.Contract.Symbol(&_TestZRC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestZRC20 *TestZRC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestZRC20 *TestZRC20Session) TotalSupply() (*big.Int, error) {
	return _TestZRC20.Contract.TotalSupply(&_TestZRC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestZRC20 *TestZRC20CallerSession) TotalSupply() (*big.Int, error) {
	return _TestZRC20.Contract.TotalSupply(&_TestZRC20.CallOpts)
}

// WithdrawGasFee is a free data retrieval call binding the contract method 0xd9eeebed.
//
// Solidity: function withdrawGasFee() view returns(address, uint256)
func (_TestZRC20 *TestZRC20Caller) WithdrawGasFee(opts *bind.CallOpts) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _TestZRC20.contract.Call(opts, &out, "withdrawGasFee")

	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// WithdrawGasFee is a free data retrieval call binding the contract method 0xd9eeebed.
//
// Solidity: function withdrawGasFee() view returns(address, uint256)
func (_TestZRC20 *TestZRC20Session) WithdrawGasFee() (common.Address, *big.Int, error) {
	return _TestZRC20.Contract.WithdrawGasFee(&_TestZRC20.CallOpts)
}

// WithdrawGasFee is a free data retrieval call binding the contract method 0xd9eeebed.
//
// Solidity: function withdrawGasFee() view returns(address, uint256)
func (_TestZRC20 *TestZRC20CallerSession) WithdrawGasFee() (common.Address, *big.Int, error) {
	return _TestZRC20.Contract.WithdrawGasFee(&_TestZRC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.Approve(&_TestZRC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.Approve(&_TestZRC20.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Transactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Session) Burn(amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.Burn(&_TestZRC20.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20TransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.Burn(&_TestZRC20.TransactOpts, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "decreaseAllowance", spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Session) DecreaseAllowance(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.DecreaseAllowance(&_TestZRC20.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20TransactorSession) DecreaseAllowance(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.DecreaseAllowance(&_TestZRC20.TransactOpts, spender, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address to, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Transactor) Deposit(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "deposit", to, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address to, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Session) Deposit(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.Deposit(&_TestZRC20.TransactOpts, to, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address to, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20TransactorSession) Deposit(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.Deposit(&_TestZRC20.TransactOpts, to, amount)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "increaseAllowance", spender, amount)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Session) IncreaseAllowance(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.IncreaseAllowance(&_TestZRC20.TransactOpts, spender, amount)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20TransactorSession) IncreaseAllowance(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.IncreaseAllowance(&_TestZRC20.TransactOpts, spender, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.Transfer(&_TestZRC20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.Transfer(&_TestZRC20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.TransferFrom(&_TestZRC20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.TransferFrom(&_TestZRC20.TransactOpts, sender, recipient, amount)
}

// UpdateGasLimit is a paid mutator transaction binding the contract method 0xf687d12a.
//
// Solidity: function updateGasLimit(uint256 gasLimit) returns()
func (_TestZRC20 *TestZRC20Transactor) UpdateGasLimit(opts *bind.TransactOpts, gasLimit *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "updateGasLimit", gasLimit)
}

// UpdateGasLimit is a paid mutator transaction binding the contract method 0xf687d12a.
//
// Solidity: function updateGasLimit(uint256 gasLimit) returns()
func (_TestZRC20 *TestZRC20Session) UpdateGasLimit(gasLimit *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.UpdateGasLimit(&_TestZRC20.TransactOpts, gasLimit)
}

// UpdateGasLimit is a paid mutator transaction binding the contract method 0xf687d12a.
//
// Solidity: function updateGasLimit(uint256 gasLimit) returns()
func (_TestZRC20 *TestZRC20TransactorSession) UpdateGasLimit(gasLimit *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.UpdateGasLimit(&_TestZRC20.TransactOpts, gasLimit)
}

// UpdateNewField is a paid mutator transaction binding the contract method 0x732bb0e4.
//
// Solidity: function updateNewField(uint256 newField_) returns()
func (_TestZRC20 *TestZRC20Transactor) UpdateNewField(opts *bind.TransactOpts, newField_ *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "updateNewField", newField_)
}

// UpdateNewField is a paid mutator transaction binding the contract method 0x732bb0e4.
//
// Solidity: function updateNewField(uint256 newField_) returns()
func (_TestZRC20 *TestZRC20Session) UpdateNewField(newField_ *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.UpdateNewField(&_TestZRC20.TransactOpts, newField_)
}

// UpdateNewField is a paid mutator transaction binding the contract method 0x732bb0e4.
//
// Solidity: function updateNewField(uint256 newField_) returns()
func (_TestZRC20 *TestZRC20TransactorSession) UpdateNewField(newField_ *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.UpdateNewField(&_TestZRC20.TransactOpts, newField_)
}

// UpdateProtocolFlatFee is a paid mutator transaction binding the contract method 0xeddeb123.
//
// Solidity: function updateProtocolFlatFee(uint256 protocolFlatFee) returns()
func (_TestZRC20 *TestZRC20Transactor) UpdateProtocolFlatFee(opts *bind.TransactOpts, protocolFlatFee *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "updateProtocolFlatFee", protocolFlatFee)
}

// UpdateProtocolFlatFee is a paid mutator transaction binding the contract method 0xeddeb123.
//
// Solidity: function updateProtocolFlatFee(uint256 protocolFlatFee) returns()
func (_TestZRC20 *TestZRC20Session) UpdateProtocolFlatFee(protocolFlatFee *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.UpdateProtocolFlatFee(&_TestZRC20.TransactOpts, protocolFlatFee)
}

// UpdateProtocolFlatFee is a paid mutator transaction binding the contract method 0xeddeb123.
//
// Solidity: function updateProtocolFlatFee(uint256 protocolFlatFee) returns()
func (_TestZRC20 *TestZRC20TransactorSession) UpdateProtocolFlatFee(protocolFlatFee *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.UpdateProtocolFlatFee(&_TestZRC20.TransactOpts, protocolFlatFee)
}

// UpdateSystemContractAddress is a paid mutator transaction binding the contract method 0xc835d7cc.
//
// Solidity: function updateSystemContractAddress(address addr) returns()
func (_TestZRC20 *TestZRC20Transactor) UpdateSystemContractAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "updateSystemContractAddress", addr)
}

// UpdateSystemContractAddress is a paid mutator transaction binding the contract method 0xc835d7cc.
//
// Solidity: function updateSystemContractAddress(address addr) returns()
func (_TestZRC20 *TestZRC20Session) UpdateSystemContractAddress(addr common.Address) (*types.Transaction, error) {
	return _TestZRC20.Contract.UpdateSystemContractAddress(&_TestZRC20.TransactOpts, addr)
}

// UpdateSystemContractAddress is a paid mutator transaction binding the contract method 0xc835d7cc.
//
// Solidity: function updateSystemContractAddress(address addr) returns()
func (_TestZRC20 *TestZRC20TransactorSession) UpdateSystemContractAddress(addr common.Address) (*types.Transaction, error) {
	return _TestZRC20.Contract.UpdateSystemContractAddress(&_TestZRC20.TransactOpts, addr)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc7012626.
//
// Solidity: function withdraw(bytes to, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Transactor) Withdraw(opts *bind.TransactOpts, to []byte, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.contract.Transact(opts, "withdraw", to, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc7012626.
//
// Solidity: function withdraw(bytes to, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20Session) Withdraw(to []byte, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.Withdraw(&_TestZRC20.TransactOpts, to, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc7012626.
//
// Solidity: function withdraw(bytes to, uint256 amount) returns(bool)
func (_TestZRC20 *TestZRC20TransactorSession) Withdraw(to []byte, amount *big.Int) (*types.Transaction, error) {
	return _TestZRC20.Contract.Withdraw(&_TestZRC20.TransactOpts, to, amount)
}

// TestZRC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TestZRC20 contract.
type TestZRC20ApprovalIterator struct {
	Event *TestZRC20Approval // Event containing the contract specifics and raw log

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
func (it *TestZRC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestZRC20Approval)
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
		it.Event = new(TestZRC20Approval)
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
func (it *TestZRC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestZRC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestZRC20Approval represents a Approval event raised by the TestZRC20 contract.
type TestZRC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestZRC20 *TestZRC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TestZRC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TestZRC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TestZRC20ApprovalIterator{contract: _TestZRC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestZRC20 *TestZRC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TestZRC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TestZRC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestZRC20Approval)
				if err := _TestZRC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestZRC20 *TestZRC20Filterer) ParseApproval(log types.Log) (*TestZRC20Approval, error) {
	event := new(TestZRC20Approval)
	if err := _TestZRC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestZRC20DepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the TestZRC20 contract.
type TestZRC20DepositIterator struct {
	Event *TestZRC20Deposit // Event containing the contract specifics and raw log

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
func (it *TestZRC20DepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestZRC20Deposit)
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
		it.Event = new(TestZRC20Deposit)
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
func (it *TestZRC20DepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestZRC20DepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestZRC20Deposit represents a Deposit event raised by the TestZRC20 contract.
type TestZRC20Deposit struct {
	From  []byte
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab3.
//
// Solidity: event Deposit(bytes from, address indexed to, uint256 value)
func (_TestZRC20 *TestZRC20Filterer) FilterDeposit(opts *bind.FilterOpts, to []common.Address) (*TestZRC20DepositIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestZRC20.contract.FilterLogs(opts, "Deposit", toRule)
	if err != nil {
		return nil, err
	}
	return &TestZRC20DepositIterator{contract: _TestZRC20.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab3.
//
// Solidity: event Deposit(bytes from, address indexed to, uint256 value)
func (_TestZRC20 *TestZRC20Filterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *TestZRC20Deposit, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestZRC20.contract.WatchLogs(opts, "Deposit", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestZRC20Deposit)
				if err := _TestZRC20.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab3.
//
// Solidity: event Deposit(bytes from, address indexed to, uint256 value)
func (_TestZRC20 *TestZRC20Filterer) ParseDeposit(log types.Log) (*TestZRC20Deposit, error) {
	event := new(TestZRC20Deposit)
	if err := _TestZRC20.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestZRC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TestZRC20 contract.
type TestZRC20TransferIterator struct {
	Event *TestZRC20Transfer // Event containing the contract specifics and raw log

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
func (it *TestZRC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestZRC20Transfer)
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
		it.Event = new(TestZRC20Transfer)
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
func (it *TestZRC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestZRC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestZRC20Transfer represents a Transfer event raised by the TestZRC20 contract.
type TestZRC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestZRC20 *TestZRC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestZRC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestZRC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestZRC20TransferIterator{contract: _TestZRC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestZRC20 *TestZRC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TestZRC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestZRC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestZRC20Transfer)
				if err := _TestZRC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestZRC20 *TestZRC20Filterer) ParseTransfer(log types.Log) (*TestZRC20Transfer, error) {
	event := new(TestZRC20Transfer)
	if err := _TestZRC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestZRC20UpdatedGasLimitIterator is returned from FilterUpdatedGasLimit and is used to iterate over the raw logs and unpacked data for UpdatedGasLimit events raised by the TestZRC20 contract.
type TestZRC20UpdatedGasLimitIterator struct {
	Event *TestZRC20UpdatedGasLimit // Event containing the contract specifics and raw log

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
func (it *TestZRC20UpdatedGasLimitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestZRC20UpdatedGasLimit)
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
		it.Event = new(TestZRC20UpdatedGasLimit)
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
func (it *TestZRC20UpdatedGasLimitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestZRC20UpdatedGasLimitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestZRC20UpdatedGasLimit represents a UpdatedGasLimit event raised by the TestZRC20 contract.
type TestZRC20UpdatedGasLimit struct {
	GasLimit *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUpdatedGasLimit is a free log retrieval operation binding the contract event 0xff5788270f43bfc1ca41c503606d2594aa3023a1a7547de403a3e2f146a4a80a.
//
// Solidity: event UpdatedGasLimit(uint256 gasLimit)
func (_TestZRC20 *TestZRC20Filterer) FilterUpdatedGasLimit(opts *bind.FilterOpts) (*TestZRC20UpdatedGasLimitIterator, error) {

	logs, sub, err := _TestZRC20.contract.FilterLogs(opts, "UpdatedGasLimit")
	if err != nil {
		return nil, err
	}
	return &TestZRC20UpdatedGasLimitIterator{contract: _TestZRC20.contract, event: "UpdatedGasLimit", logs: logs, sub: sub}, nil
}

// WatchUpdatedGasLimit is a free log subscription operation binding the contract event 0xff5788270f43bfc1ca41c503606d2594aa3023a1a7547de403a3e2f146a4a80a.
//
// Solidity: event UpdatedGasLimit(uint256 gasLimit)
func (_TestZRC20 *TestZRC20Filterer) WatchUpdatedGasLimit(opts *bind.WatchOpts, sink chan<- *TestZRC20UpdatedGasLimit) (event.Subscription, error) {

	logs, sub, err := _TestZRC20.contract.WatchLogs(opts, "UpdatedGasLimit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestZRC20UpdatedGasLimit)
				if err := _TestZRC20.contract.UnpackLog(event, "UpdatedGasLimit", log); err != nil {
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

// ParseUpdatedGasLimit is a log parse operation binding the contract event 0xff5788270f43bfc1ca41c503606d2594aa3023a1a7547de403a3e2f146a4a80a.
//
// Solidity: event UpdatedGasLimit(uint256 gasLimit)
func (_TestZRC20 *TestZRC20Filterer) ParseUpdatedGasLimit(log types.Log) (*TestZRC20UpdatedGasLimit, error) {
	event := new(TestZRC20UpdatedGasLimit)
	if err := _TestZRC20.contract.UnpackLog(event, "UpdatedGasLimit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestZRC20UpdatedProtocolFlatFeeIterator is returned from FilterUpdatedProtocolFlatFee and is used to iterate over the raw logs and unpacked data for UpdatedProtocolFlatFee events raised by the TestZRC20 contract.
type TestZRC20UpdatedProtocolFlatFeeIterator struct {
	Event *TestZRC20UpdatedProtocolFlatFee // Event containing the contract specifics and raw log

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
func (it *TestZRC20UpdatedProtocolFlatFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestZRC20UpdatedProtocolFlatFee)
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
		it.Event = new(TestZRC20UpdatedProtocolFlatFee)
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
func (it *TestZRC20UpdatedProtocolFlatFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestZRC20UpdatedProtocolFlatFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestZRC20UpdatedProtocolFlatFee represents a UpdatedProtocolFlatFee event raised by the TestZRC20 contract.
type TestZRC20UpdatedProtocolFlatFee struct {
	ProtocolFlatFee *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpdatedProtocolFlatFee is a free log retrieval operation binding the contract event 0xef13af88e424b5d15f49c77758542c1938b08b8b95b91ed0751f98ba99000d8f.
//
// Solidity: event UpdatedProtocolFlatFee(uint256 protocolFlatFee)
func (_TestZRC20 *TestZRC20Filterer) FilterUpdatedProtocolFlatFee(opts *bind.FilterOpts) (*TestZRC20UpdatedProtocolFlatFeeIterator, error) {

	logs, sub, err := _TestZRC20.contract.FilterLogs(opts, "UpdatedProtocolFlatFee")
	if err != nil {
		return nil, err
	}
	return &TestZRC20UpdatedProtocolFlatFeeIterator{contract: _TestZRC20.contract, event: "UpdatedProtocolFlatFee", logs: logs, sub: sub}, nil
}

// WatchUpdatedProtocolFlatFee is a free log subscription operation binding the contract event 0xef13af88e424b5d15f49c77758542c1938b08b8b95b91ed0751f98ba99000d8f.
//
// Solidity: event UpdatedProtocolFlatFee(uint256 protocolFlatFee)
func (_TestZRC20 *TestZRC20Filterer) WatchUpdatedProtocolFlatFee(opts *bind.WatchOpts, sink chan<- *TestZRC20UpdatedProtocolFlatFee) (event.Subscription, error) {

	logs, sub, err := _TestZRC20.contract.WatchLogs(opts, "UpdatedProtocolFlatFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestZRC20UpdatedProtocolFlatFee)
				if err := _TestZRC20.contract.UnpackLog(event, "UpdatedProtocolFlatFee", log); err != nil {
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

// ParseUpdatedProtocolFlatFee is a log parse operation binding the contract event 0xef13af88e424b5d15f49c77758542c1938b08b8b95b91ed0751f98ba99000d8f.
//
// Solidity: event UpdatedProtocolFlatFee(uint256 protocolFlatFee)
func (_TestZRC20 *TestZRC20Filterer) ParseUpdatedProtocolFlatFee(log types.Log) (*TestZRC20UpdatedProtocolFlatFee, error) {
	event := new(TestZRC20UpdatedProtocolFlatFee)
	if err := _TestZRC20.contract.UnpackLog(event, "UpdatedProtocolFlatFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestZRC20UpdatedSystemContractIterator is returned from FilterUpdatedSystemContract and is used to iterate over the raw logs and unpacked data for UpdatedSystemContract events raised by the TestZRC20 contract.
type TestZRC20UpdatedSystemContractIterator struct {
	Event *TestZRC20UpdatedSystemContract // Event containing the contract specifics and raw log

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
func (it *TestZRC20UpdatedSystemContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestZRC20UpdatedSystemContract)
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
		it.Event = new(TestZRC20UpdatedSystemContract)
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
func (it *TestZRC20UpdatedSystemContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestZRC20UpdatedSystemContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestZRC20UpdatedSystemContract represents a UpdatedSystemContract event raised by the TestZRC20 contract.
type TestZRC20UpdatedSystemContract struct {
	SystemContract common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpdatedSystemContract is a free log retrieval operation binding the contract event 0xd55614e962c5fd6ece71614f6348d702468a997a394dd5e5c1677950226d97ae.
//
// Solidity: event UpdatedSystemContract(address systemContract)
func (_TestZRC20 *TestZRC20Filterer) FilterUpdatedSystemContract(opts *bind.FilterOpts) (*TestZRC20UpdatedSystemContractIterator, error) {

	logs, sub, err := _TestZRC20.contract.FilterLogs(opts, "UpdatedSystemContract")
	if err != nil {
		return nil, err
	}
	return &TestZRC20UpdatedSystemContractIterator{contract: _TestZRC20.contract, event: "UpdatedSystemContract", logs: logs, sub: sub}, nil
}

// WatchUpdatedSystemContract is a free log subscription operation binding the contract event 0xd55614e962c5fd6ece71614f6348d702468a997a394dd5e5c1677950226d97ae.
//
// Solidity: event UpdatedSystemContract(address systemContract)
func (_TestZRC20 *TestZRC20Filterer) WatchUpdatedSystemContract(opts *bind.WatchOpts, sink chan<- *TestZRC20UpdatedSystemContract) (event.Subscription, error) {

	logs, sub, err := _TestZRC20.contract.WatchLogs(opts, "UpdatedSystemContract")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestZRC20UpdatedSystemContract)
				if err := _TestZRC20.contract.UnpackLog(event, "UpdatedSystemContract", log); err != nil {
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

// ParseUpdatedSystemContract is a log parse operation binding the contract event 0xd55614e962c5fd6ece71614f6348d702468a997a394dd5e5c1677950226d97ae.
//
// Solidity: event UpdatedSystemContract(address systemContract)
func (_TestZRC20 *TestZRC20Filterer) ParseUpdatedSystemContract(log types.Log) (*TestZRC20UpdatedSystemContract, error) {
	event := new(TestZRC20UpdatedSystemContract)
	if err := _TestZRC20.contract.UnpackLog(event, "UpdatedSystemContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestZRC20WithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the TestZRC20 contract.
type TestZRC20WithdrawalIterator struct {
	Event *TestZRC20Withdrawal // Event containing the contract specifics and raw log

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
func (it *TestZRC20WithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestZRC20Withdrawal)
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
		it.Event = new(TestZRC20Withdrawal)
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
func (it *TestZRC20WithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestZRC20WithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestZRC20Withdrawal represents a Withdrawal event raised by the TestZRC20 contract.
type TestZRC20Withdrawal struct {
	From            common.Address
	To              []byte
	Value           *big.Int
	Gasfee          *big.Int
	ProtocolFlatFee *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955.
//
// Solidity: event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee, uint256 protocolFlatFee)
func (_TestZRC20 *TestZRC20Filterer) FilterWithdrawal(opts *bind.FilterOpts, from []common.Address) (*TestZRC20WithdrawalIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _TestZRC20.contract.FilterLogs(opts, "Withdrawal", fromRule)
	if err != nil {
		return nil, err
	}
	return &TestZRC20WithdrawalIterator{contract: _TestZRC20.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955.
//
// Solidity: event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee, uint256 protocolFlatFee)
func (_TestZRC20 *TestZRC20Filterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *TestZRC20Withdrawal, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _TestZRC20.contract.WatchLogs(opts, "Withdrawal", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestZRC20Withdrawal)
				if err := _TestZRC20.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955.
//
// Solidity: event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee, uint256 protocolFlatFee)
func (_TestZRC20 *TestZRC20Filterer) ParseWithdrawal(log types.Log) (*TestZRC20Withdrawal, error) {
	event := new(TestZRC20Withdrawal)
	if err := _TestZRC20.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
