// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zetanoneth

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

// ZetaNonEthMetaData contains all meta data concerning the ZetaNonEth contract.
var ZetaNonEthMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tssAddress_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tssAddressUpdater_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotConnector\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotTss\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotTssOrUpdater\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotTssUpdater\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZetaTransferError\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"burnee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burnt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"mintee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"internalSendHash\",\"type\":\"bytes32\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"connectorAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"mintee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"internalSendHash\",\"type\":\"bytes32\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceTssAddressUpdater\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tssAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tssAddressUpdater\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tssAddress_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"connectorAddress_\",\"type\":\"address\"}],\"name\":\"updateTssAndConnectorAddresses\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200237c3803806200237c8339818101604052810190620000379190620002c8565b6040518060400160405280600481526020017f5a657461000000000000000000000000000000000000000000000000000000008152506040518060400160405280600481526020017f5a455441000000000000000000000000000000000000000000000000000000008152508160039080519060200190620000bb92919062000201565b508060049080519060200190620000d492919062000201565b505050600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614806200013f5750600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16145b1562000177576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050620003c7565b8280546200020f9062000343565b90600052602060002090601f0160209004810192826200023357600085556200027f565b82601f106200024e57805160ff19168380011785556200027f565b828001600101855582156200027f579182015b828111156200027e57825182559160200191906001019062000261565b5b5090506200028e919062000292565b5090565b5b80821115620002ad57600081600090555060010162000293565b5090565b600081519050620002c281620003ad565b92915050565b60008060408385031215620002e257620002e1620003a8565b5b6000620002f285828601620002b1565b92505060206200030585828601620002b1565b9150509250929050565b60006200031c8262000323565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600060028204905060018216806200035c57607f821691505b6020821081141562000373576200037262000379565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600080fd5b620003b8816200030f565b8114620003c457600080fd5b50565b611fa580620003d76000396000f3fe608060405234801561001057600080fd5b50600436106101215760003560e01c806342966c68116100ad57806395d89b411161007157806395d89b41146102f6578063a457c2d714610314578063a9059cbb14610344578063bff9662a14610374578063dd62ed3e1461039257610121565b806342966c68146102665780635b1125911461028257806370a08231146102a0578063779e3b63146102d057806379cc6790146102da57610121565b80631e458bee116100f45780631e458bee146101ae57806323b872dd146101ca578063313ce567146101fa578063328a01d014610218578063395093511461023657610121565b806306fdde0314610126578063095ea7b31461014457806315d57fd41461017457806318160ddd14610190575b600080fd5b61012e6103c2565b60405161013b91906118ea565b60405180910390f35b61015e60048036038101906101599190611621565b610454565b60405161016b91906118cf565b60405180910390f35b61018e6004803603810190610189919061158e565b610477565b005b610198610689565b6040516101a59190611a4c565b60405180910390f35b6101c860048036038101906101c39190611661565b610693565b005b6101e460048036038101906101df91906115ce565b610783565b6040516101f191906118cf565b60405180910390f35b6102026107b2565b60405161020f9190611a67565b60405180910390f35b6102206107bb565b60405161022d91906118b4565b60405180910390f35b610250600480360381019061024b9190611621565b6107e1565b60405161025d91906118cf565b60405180910390f35b610280600480360381019061027b91906116b4565b610818565b005b61028a61082c565b60405161029791906118b4565b60405180910390f35b6102ba60048036038101906102b59190611561565b610852565b6040516102c79190611a4c565b60405180910390f35b6102d861089a565b005b6102f460048036038101906102ef9190611621565b610a1a565b005b6102fe610b08565b60405161030b91906118ea565b60405180910390f35b61032e60048036038101906103299190611621565b610b9a565b60405161033b91906118cf565b60405180910390f35b61035e60048036038101906103599190611621565b610c11565b60405161036b91906118cf565b60405180910390f35b61037c610c34565b60405161038991906118b4565b60405180910390f35b6103ac60048036038101906103a7919061158e565b610c5a565b6040516103b99190611a4c565b60405180910390f35b6060600380546103d190611bba565b80601f01602080910402602001604051908101604052809291908181526020018280546103fd90611bba565b801561044a5780601f1061041f5761010080835404028352916020019161044a565b820191906000526020600020905b81548152906001019060200180831161042d57829003601f168201915b5050505050905090565b60008061045f610ce1565b905061046c818585610ce9565b600191505092915050565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141580156105235750600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b1561056557336040517fcdfcef9700000000000000000000000000000000000000000000000000000000815260040161055c91906118b4565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614806105cc5750600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16145b15610603576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050565b6000600254905090565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461072557336040517f3fe32fba00000000000000000000000000000000000000000000000000000000815260040161071c91906118b4565b60405180910390fd5b61072f8383610eb4565b808373ffffffffffffffffffffffffffffffffffffffff167fc263b302aec62d29105026245f19e16f8e0137066ccd4a8bd941f716bd4096bb846040516107769190611a4c565b60405180910390a3505050565b60008061078e610ce1565b905061079b858285611014565b6107a68585856110a0565b60019150509392505050565b60006012905090565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000806107ec610ce1565b905061080d8185856107fe8589610c5a565b6108089190611a9e565b610ce9565b600191505092915050565b610829610823610ce1565b82611321565b50565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461092c57336040517fe700765e00000000000000000000000000000000000000000000000000000000815260040161092391906118b4565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614156109b5576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610aac57336040517f3fe32fba000000000000000000000000000000000000000000000000000000008152600401610aa391906118b4565b60405180910390fd5b610ab682826114f8565b8173ffffffffffffffffffffffffffffffffffffffff167f919f7e2092ffcc9d09f599be18d8152860b0c054df788a33bc549cdd9d0f15b182604051610afc9190611a4c565b60405180910390a25050565b606060048054610b1790611bba565b80601f0160208091040260200160405190810160405280929190818152602001828054610b4390611bba565b8015610b905780601f10610b6557610100808354040283529160200191610b90565b820191906000526020600020905b815481529060010190602001808311610b7357829003601f168201915b5050505050905090565b600080610ba5610ce1565b90506000610bb38286610c5a565b905083811015610bf8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bef90611a0c565b60405180910390fd5b610c058286868403610ce9565b60019250505092915050565b600080610c1c610ce1565b9050610c298185856110a0565b600191505092915050565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610d59576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d50906119ec565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610dc9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dc09061194c565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92583604051610ea79190611a4c565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610f24576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f1b90611a2c565b60405180910390fd5b610f3060008383611518565b8060026000828254610f429190611a9e565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610f979190611a9e565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610ffc9190611a4c565b60405180910390a36110106000838361151d565b5050565b60006110208484610c5a565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811461109a578181101561108c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110839061196c565b60405180910390fd5b6110998484848403610ce9565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415611110576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611107906119cc565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611180576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111779061190c565b60405180910390fd5b61118b838383611518565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015611211576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112089061198c565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546112a49190611a9e565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516113089190611a4c565b60405180910390a361131b84848461151d565b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611391576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611388906119ac565b60405180910390fd5b61139d82600083611518565b60008060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015611423576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161141a9061192c565b60405180910390fd5b8181036000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816002600082825461147a9190611af4565b92505081905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516114df9190611a4c565b60405180910390a36114f38360008461151d565b505050565b61150a82611504610ce1565b83611014565b6115148282611321565b5050565b505050565b505050565b60008135905061153181611f2a565b92915050565b60008135905061154681611f41565b92915050565b60008135905061155b81611f58565b92915050565b60006020828403121561157757611576611c4a565b5b600061158584828501611522565b91505092915050565b600080604083850312156115a5576115a4611c4a565b5b60006115b385828601611522565b92505060206115c485828601611522565b9150509250929050565b6000806000606084860312156115e7576115e6611c4a565b5b60006115f586828701611522565b935050602061160686828701611522565b92505060406116178682870161154c565b9150509250925092565b6000806040838503121561163857611637611c4a565b5b600061164685828601611522565b92505060206116578582860161154c565b9150509250929050565b60008060006060848603121561167a57611679611c4a565b5b600061168886828701611522565b93505060206116998682870161154c565b92505060406116aa86828701611537565b9150509250925092565b6000602082840312156116ca576116c9611c4a565b5b60006116d88482850161154c565b91505092915050565b6116ea81611b28565b82525050565b6116f981611b3a565b82525050565b600061170a82611a82565b6117148185611a8d565b9350611724818560208601611b87565b61172d81611c4f565b840191505092915050565b6000611745602383611a8d565b915061175082611c60565b604082019050919050565b6000611768602283611a8d565b915061177382611caf565b604082019050919050565b600061178b602283611a8d565b915061179682611cfe565b604082019050919050565b60006117ae601d83611a8d565b91506117b982611d4d565b602082019050919050565b60006117d1602683611a8d565b91506117dc82611d76565b604082019050919050565b60006117f4602183611a8d565b91506117ff82611dc5565b604082019050919050565b6000611817602583611a8d565b915061182282611e14565b604082019050919050565b600061183a602483611a8d565b915061184582611e63565b604082019050919050565b600061185d602583611a8d565b915061186882611eb2565b604082019050919050565b6000611880601f83611a8d565b915061188b82611f01565b602082019050919050565b61189f81611b70565b82525050565b6118ae81611b7a565b82525050565b60006020820190506118c960008301846116e1565b92915050565b60006020820190506118e460008301846116f0565b92915050565b6000602082019050818103600083015261190481846116ff565b905092915050565b6000602082019050818103600083015261192581611738565b9050919050565b600060208201905081810360008301526119458161175b565b9050919050565b600060208201905081810360008301526119658161177e565b9050919050565b60006020820190508181036000830152611985816117a1565b9050919050565b600060208201905081810360008301526119a5816117c4565b9050919050565b600060208201905081810360008301526119c5816117e7565b9050919050565b600060208201905081810360008301526119e58161180a565b9050919050565b60006020820190508181036000830152611a058161182d565b9050919050565b60006020820190508181036000830152611a2581611850565b9050919050565b60006020820190508181036000830152611a4581611873565b9050919050565b6000602082019050611a616000830184611896565b92915050565b6000602082019050611a7c60008301846118a5565b92915050565b600081519050919050565b600082825260208201905092915050565b6000611aa982611b70565b9150611ab483611b70565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611ae957611ae8611bec565b5b828201905092915050565b6000611aff82611b70565b9150611b0a83611b70565b925082821015611b1d57611b1c611bec565b5b828203905092915050565b6000611b3382611b50565b9050919050565b60008115159050919050565b6000819050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600060ff82169050919050565b60005b83811015611ba5578082015181840152602081019050611b8a565b83811115611bb4576000848401525b50505050565b60006002820490506001821680611bd257607f821691505b60208210811415611be657611be5611c1b565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600080fd5b6000601f19601f8301169050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60008201527f6365000000000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360008201527f7300000000000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b611f3381611b28565b8114611f3e57600080fd5b50565b611f4a81611b46565b8114611f5557600080fd5b50565b611f6181611b70565b8114611f6c57600080fd5b5056fea264697066735822122008f6e523d84c166be747205bdbde5d8ba13c478304202ab40df209426973ae7564736f6c63430008070033",
}

// ZetaNonEthABI is the input ABI used to generate the binding from.
// Deprecated: Use ZetaNonEthMetaData.ABI instead.
var ZetaNonEthABI = ZetaNonEthMetaData.ABI

// ZetaNonEthBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ZetaNonEthMetaData.Bin instead.
var ZetaNonEthBin = ZetaNonEthMetaData.Bin

// DeployZetaNonEth deploys a new Ethereum contract, binding an instance of ZetaNonEth to it.
func DeployZetaNonEth(auth *bind.TransactOpts, backend bind.ContractBackend, tssAddress_ common.Address, tssAddressUpdater_ common.Address) (common.Address, *types.Transaction, *ZetaNonEth, error) {
	parsed, err := ZetaNonEthMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ZetaNonEthBin), backend, tssAddress_, tssAddressUpdater_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ZetaNonEth{ZetaNonEthCaller: ZetaNonEthCaller{contract: contract}, ZetaNonEthTransactor: ZetaNonEthTransactor{contract: contract}, ZetaNonEthFilterer: ZetaNonEthFilterer{contract: contract}}, nil
}

// ZetaNonEth is an auto generated Go binding around an Ethereum contract.
type ZetaNonEth struct {
	ZetaNonEthCaller     // Read-only binding to the contract
	ZetaNonEthTransactor // Write-only binding to the contract
	ZetaNonEthFilterer   // Log filterer for contract events
}

// ZetaNonEthCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZetaNonEthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZetaNonEthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZetaNonEthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZetaNonEthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZetaNonEthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZetaNonEthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZetaNonEthSession struct {
	Contract     *ZetaNonEth       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZetaNonEthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZetaNonEthCallerSession struct {
	Contract *ZetaNonEthCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ZetaNonEthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZetaNonEthTransactorSession struct {
	Contract     *ZetaNonEthTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ZetaNonEthRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZetaNonEthRaw struct {
	Contract *ZetaNonEth // Generic contract binding to access the raw methods on
}

// ZetaNonEthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZetaNonEthCallerRaw struct {
	Contract *ZetaNonEthCaller // Generic read-only contract binding to access the raw methods on
}

// ZetaNonEthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZetaNonEthTransactorRaw struct {
	Contract *ZetaNonEthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZetaNonEth creates a new instance of ZetaNonEth, bound to a specific deployed contract.
func NewZetaNonEth(address common.Address, backend bind.ContractBackend) (*ZetaNonEth, error) {
	contract, err := bindZetaNonEth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZetaNonEth{ZetaNonEthCaller: ZetaNonEthCaller{contract: contract}, ZetaNonEthTransactor: ZetaNonEthTransactor{contract: contract}, ZetaNonEthFilterer: ZetaNonEthFilterer{contract: contract}}, nil
}

// NewZetaNonEthCaller creates a new read-only instance of ZetaNonEth, bound to a specific deployed contract.
func NewZetaNonEthCaller(address common.Address, caller bind.ContractCaller) (*ZetaNonEthCaller, error) {
	contract, err := bindZetaNonEth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZetaNonEthCaller{contract: contract}, nil
}

// NewZetaNonEthTransactor creates a new write-only instance of ZetaNonEth, bound to a specific deployed contract.
func NewZetaNonEthTransactor(address common.Address, transactor bind.ContractTransactor) (*ZetaNonEthTransactor, error) {
	contract, err := bindZetaNonEth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZetaNonEthTransactor{contract: contract}, nil
}

// NewZetaNonEthFilterer creates a new log filterer instance of ZetaNonEth, bound to a specific deployed contract.
func NewZetaNonEthFilterer(address common.Address, filterer bind.ContractFilterer) (*ZetaNonEthFilterer, error) {
	contract, err := bindZetaNonEth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZetaNonEthFilterer{contract: contract}, nil
}

// bindZetaNonEth binds a generic wrapper to an already deployed contract.
func bindZetaNonEth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ZetaNonEthABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZetaNonEth *ZetaNonEthRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZetaNonEth.Contract.ZetaNonEthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZetaNonEth *ZetaNonEthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.ZetaNonEthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZetaNonEth *ZetaNonEthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.ZetaNonEthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZetaNonEth *ZetaNonEthCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZetaNonEth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZetaNonEth *ZetaNonEthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZetaNonEth *ZetaNonEthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZetaNonEth *ZetaNonEthCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZetaNonEth.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZetaNonEth *ZetaNonEthSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZetaNonEth.Contract.Allowance(&_ZetaNonEth.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ZetaNonEth *ZetaNonEthCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ZetaNonEth.Contract.Allowance(&_ZetaNonEth.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZetaNonEth *ZetaNonEthCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZetaNonEth.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZetaNonEth *ZetaNonEthSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZetaNonEth.Contract.BalanceOf(&_ZetaNonEth.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ZetaNonEth *ZetaNonEthCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ZetaNonEth.Contract.BalanceOf(&_ZetaNonEth.CallOpts, account)
}

// ConnectorAddress is a free data retrieval call binding the contract method 0xbff9662a.
//
// Solidity: function connectorAddress() view returns(address)
func (_ZetaNonEth *ZetaNonEthCaller) ConnectorAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZetaNonEth.contract.Call(opts, &out, "connectorAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ConnectorAddress is a free data retrieval call binding the contract method 0xbff9662a.
//
// Solidity: function connectorAddress() view returns(address)
func (_ZetaNonEth *ZetaNonEthSession) ConnectorAddress() (common.Address, error) {
	return _ZetaNonEth.Contract.ConnectorAddress(&_ZetaNonEth.CallOpts)
}

// ConnectorAddress is a free data retrieval call binding the contract method 0xbff9662a.
//
// Solidity: function connectorAddress() view returns(address)
func (_ZetaNonEth *ZetaNonEthCallerSession) ConnectorAddress() (common.Address, error) {
	return _ZetaNonEth.Contract.ConnectorAddress(&_ZetaNonEth.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZetaNonEth *ZetaNonEthCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ZetaNonEth.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZetaNonEth *ZetaNonEthSession) Decimals() (uint8, error) {
	return _ZetaNonEth.Contract.Decimals(&_ZetaNonEth.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ZetaNonEth *ZetaNonEthCallerSession) Decimals() (uint8, error) {
	return _ZetaNonEth.Contract.Decimals(&_ZetaNonEth.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZetaNonEth *ZetaNonEthCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZetaNonEth.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZetaNonEth *ZetaNonEthSession) Name() (string, error) {
	return _ZetaNonEth.Contract.Name(&_ZetaNonEth.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ZetaNonEth *ZetaNonEthCallerSession) Name() (string, error) {
	return _ZetaNonEth.Contract.Name(&_ZetaNonEth.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZetaNonEth *ZetaNonEthCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZetaNonEth.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZetaNonEth *ZetaNonEthSession) Symbol() (string, error) {
	return _ZetaNonEth.Contract.Symbol(&_ZetaNonEth.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ZetaNonEth *ZetaNonEthCallerSession) Symbol() (string, error) {
	return _ZetaNonEth.Contract.Symbol(&_ZetaNonEth.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZetaNonEth *ZetaNonEthCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZetaNonEth.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZetaNonEth *ZetaNonEthSession) TotalSupply() (*big.Int, error) {
	return _ZetaNonEth.Contract.TotalSupply(&_ZetaNonEth.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ZetaNonEth *ZetaNonEthCallerSession) TotalSupply() (*big.Int, error) {
	return _ZetaNonEth.Contract.TotalSupply(&_ZetaNonEth.CallOpts)
}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_ZetaNonEth *ZetaNonEthCaller) TssAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZetaNonEth.contract.Call(opts, &out, "tssAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_ZetaNonEth *ZetaNonEthSession) TssAddress() (common.Address, error) {
	return _ZetaNonEth.Contract.TssAddress(&_ZetaNonEth.CallOpts)
}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_ZetaNonEth *ZetaNonEthCallerSession) TssAddress() (common.Address, error) {
	return _ZetaNonEth.Contract.TssAddress(&_ZetaNonEth.CallOpts)
}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_ZetaNonEth *ZetaNonEthCaller) TssAddressUpdater(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZetaNonEth.contract.Call(opts, &out, "tssAddressUpdater")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_ZetaNonEth *ZetaNonEthSession) TssAddressUpdater() (common.Address, error) {
	return _ZetaNonEth.Contract.TssAddressUpdater(&_ZetaNonEth.CallOpts)
}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_ZetaNonEth *ZetaNonEthCallerSession) TssAddressUpdater() (common.Address, error) {
	return _ZetaNonEth.Contract.TssAddressUpdater(&_ZetaNonEth.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ZetaNonEth *ZetaNonEthTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ZetaNonEth *ZetaNonEthSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.Approve(&_ZetaNonEth.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ZetaNonEth *ZetaNonEthTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.Approve(&_ZetaNonEth.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ZetaNonEth *ZetaNonEthTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ZetaNonEth *ZetaNonEthSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.Burn(&_ZetaNonEth.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ZetaNonEth *ZetaNonEthTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.Burn(&_ZetaNonEth.TransactOpts, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ZetaNonEth *ZetaNonEthTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.contract.Transact(opts, "burnFrom", account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ZetaNonEth *ZetaNonEthSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.BurnFrom(&_ZetaNonEth.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ZetaNonEth *ZetaNonEthTransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.BurnFrom(&_ZetaNonEth.TransactOpts, account, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ZetaNonEth *ZetaNonEthTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ZetaNonEth *ZetaNonEthSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.DecreaseAllowance(&_ZetaNonEth.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ZetaNonEth *ZetaNonEthTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.DecreaseAllowance(&_ZetaNonEth.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ZetaNonEth *ZetaNonEthTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ZetaNonEth *ZetaNonEthSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.IncreaseAllowance(&_ZetaNonEth.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ZetaNonEth *ZetaNonEthTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.IncreaseAllowance(&_ZetaNonEth.TransactOpts, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x1e458bee.
//
// Solidity: function mint(address mintee, uint256 value, bytes32 internalSendHash) returns()
func (_ZetaNonEth *ZetaNonEthTransactor) Mint(opts *bind.TransactOpts, mintee common.Address, value *big.Int, internalSendHash [32]byte) (*types.Transaction, error) {
	return _ZetaNonEth.contract.Transact(opts, "mint", mintee, value, internalSendHash)
}

// Mint is a paid mutator transaction binding the contract method 0x1e458bee.
//
// Solidity: function mint(address mintee, uint256 value, bytes32 internalSendHash) returns()
func (_ZetaNonEth *ZetaNonEthSession) Mint(mintee common.Address, value *big.Int, internalSendHash [32]byte) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.Mint(&_ZetaNonEth.TransactOpts, mintee, value, internalSendHash)
}

// Mint is a paid mutator transaction binding the contract method 0x1e458bee.
//
// Solidity: function mint(address mintee, uint256 value, bytes32 internalSendHash) returns()
func (_ZetaNonEth *ZetaNonEthTransactorSession) Mint(mintee common.Address, value *big.Int, internalSendHash [32]byte) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.Mint(&_ZetaNonEth.TransactOpts, mintee, value, internalSendHash)
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_ZetaNonEth *ZetaNonEthTransactor) RenounceTssAddressUpdater(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZetaNonEth.contract.Transact(opts, "renounceTssAddressUpdater")
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_ZetaNonEth *ZetaNonEthSession) RenounceTssAddressUpdater() (*types.Transaction, error) {
	return _ZetaNonEth.Contract.RenounceTssAddressUpdater(&_ZetaNonEth.TransactOpts)
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_ZetaNonEth *ZetaNonEthTransactorSession) RenounceTssAddressUpdater() (*types.Transaction, error) {
	return _ZetaNonEth.Contract.RenounceTssAddressUpdater(&_ZetaNonEth.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ZetaNonEth *ZetaNonEthTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ZetaNonEth *ZetaNonEthSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.Transfer(&_ZetaNonEth.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ZetaNonEth *ZetaNonEthTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.Transfer(&_ZetaNonEth.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ZetaNonEth *ZetaNonEthTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ZetaNonEth *ZetaNonEthSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.TransferFrom(&_ZetaNonEth.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ZetaNonEth *ZetaNonEthTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.TransferFrom(&_ZetaNonEth.TransactOpts, from, to, amount)
}

// UpdateTssAndConnectorAddresses is a paid mutator transaction binding the contract method 0x15d57fd4.
//
// Solidity: function updateTssAndConnectorAddresses(address tssAddress_, address connectorAddress_) returns()
func (_ZetaNonEth *ZetaNonEthTransactor) UpdateTssAndConnectorAddresses(opts *bind.TransactOpts, tssAddress_ common.Address, connectorAddress_ common.Address) (*types.Transaction, error) {
	return _ZetaNonEth.contract.Transact(opts, "updateTssAndConnectorAddresses", tssAddress_, connectorAddress_)
}

// UpdateTssAndConnectorAddresses is a paid mutator transaction binding the contract method 0x15d57fd4.
//
// Solidity: function updateTssAndConnectorAddresses(address tssAddress_, address connectorAddress_) returns()
func (_ZetaNonEth *ZetaNonEthSession) UpdateTssAndConnectorAddresses(tssAddress_ common.Address, connectorAddress_ common.Address) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.UpdateTssAndConnectorAddresses(&_ZetaNonEth.TransactOpts, tssAddress_, connectorAddress_)
}

// UpdateTssAndConnectorAddresses is a paid mutator transaction binding the contract method 0x15d57fd4.
//
// Solidity: function updateTssAndConnectorAddresses(address tssAddress_, address connectorAddress_) returns()
func (_ZetaNonEth *ZetaNonEthTransactorSession) UpdateTssAndConnectorAddresses(tssAddress_ common.Address, connectorAddress_ common.Address) (*types.Transaction, error) {
	return _ZetaNonEth.Contract.UpdateTssAndConnectorAddresses(&_ZetaNonEth.TransactOpts, tssAddress_, connectorAddress_)
}

// ZetaNonEthApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ZetaNonEth contract.
type ZetaNonEthApprovalIterator struct {
	Event *ZetaNonEthApproval // Event containing the contract specifics and raw log

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
func (it *ZetaNonEthApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZetaNonEthApproval)
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
		it.Event = new(ZetaNonEthApproval)
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
func (it *ZetaNonEthApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZetaNonEthApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZetaNonEthApproval represents a Approval event raised by the ZetaNonEth contract.
type ZetaNonEthApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZetaNonEth *ZetaNonEthFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ZetaNonEthApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZetaNonEth.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ZetaNonEthApprovalIterator{contract: _ZetaNonEth.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ZetaNonEth *ZetaNonEthFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ZetaNonEthApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ZetaNonEth.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZetaNonEthApproval)
				if err := _ZetaNonEth.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ZetaNonEth *ZetaNonEthFilterer) ParseApproval(log types.Log) (*ZetaNonEthApproval, error) {
	event := new(ZetaNonEthApproval)
	if err := _ZetaNonEth.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZetaNonEthBurntIterator is returned from FilterBurnt and is used to iterate over the raw logs and unpacked data for Burnt events raised by the ZetaNonEth contract.
type ZetaNonEthBurntIterator struct {
	Event *ZetaNonEthBurnt // Event containing the contract specifics and raw log

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
func (it *ZetaNonEthBurntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZetaNonEthBurnt)
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
		it.Event = new(ZetaNonEthBurnt)
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
func (it *ZetaNonEthBurntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZetaNonEthBurntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZetaNonEthBurnt represents a Burnt event raised by the ZetaNonEth contract.
type ZetaNonEthBurnt struct {
	Burnee common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBurnt is a free log retrieval operation binding the contract event 0x919f7e2092ffcc9d09f599be18d8152860b0c054df788a33bc549cdd9d0f15b1.
//
// Solidity: event Burnt(address indexed burnee, uint256 amount)
func (_ZetaNonEth *ZetaNonEthFilterer) FilterBurnt(opts *bind.FilterOpts, burnee []common.Address) (*ZetaNonEthBurntIterator, error) {

	var burneeRule []interface{}
	for _, burneeItem := range burnee {
		burneeRule = append(burneeRule, burneeItem)
	}

	logs, sub, err := _ZetaNonEth.contract.FilterLogs(opts, "Burnt", burneeRule)
	if err != nil {
		return nil, err
	}
	return &ZetaNonEthBurntIterator{contract: _ZetaNonEth.contract, event: "Burnt", logs: logs, sub: sub}, nil
}

// WatchBurnt is a free log subscription operation binding the contract event 0x919f7e2092ffcc9d09f599be18d8152860b0c054df788a33bc549cdd9d0f15b1.
//
// Solidity: event Burnt(address indexed burnee, uint256 amount)
func (_ZetaNonEth *ZetaNonEthFilterer) WatchBurnt(opts *bind.WatchOpts, sink chan<- *ZetaNonEthBurnt, burnee []common.Address) (event.Subscription, error) {

	var burneeRule []interface{}
	for _, burneeItem := range burnee {
		burneeRule = append(burneeRule, burneeItem)
	}

	logs, sub, err := _ZetaNonEth.contract.WatchLogs(opts, "Burnt", burneeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZetaNonEthBurnt)
				if err := _ZetaNonEth.contract.UnpackLog(event, "Burnt", log); err != nil {
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

// ParseBurnt is a log parse operation binding the contract event 0x919f7e2092ffcc9d09f599be18d8152860b0c054df788a33bc549cdd9d0f15b1.
//
// Solidity: event Burnt(address indexed burnee, uint256 amount)
func (_ZetaNonEth *ZetaNonEthFilterer) ParseBurnt(log types.Log) (*ZetaNonEthBurnt, error) {
	event := new(ZetaNonEthBurnt)
	if err := _ZetaNonEth.contract.UnpackLog(event, "Burnt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZetaNonEthMintedIterator is returned from FilterMinted and is used to iterate over the raw logs and unpacked data for Minted events raised by the ZetaNonEth contract.
type ZetaNonEthMintedIterator struct {
	Event *ZetaNonEthMinted // Event containing the contract specifics and raw log

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
func (it *ZetaNonEthMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZetaNonEthMinted)
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
		it.Event = new(ZetaNonEthMinted)
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
func (it *ZetaNonEthMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZetaNonEthMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZetaNonEthMinted represents a Minted event raised by the ZetaNonEth contract.
type ZetaNonEthMinted struct {
	Mintee           common.Address
	Amount           *big.Int
	InternalSendHash [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterMinted is a free log retrieval operation binding the contract event 0xc263b302aec62d29105026245f19e16f8e0137066ccd4a8bd941f716bd4096bb.
//
// Solidity: event Minted(address indexed mintee, uint256 amount, bytes32 indexed internalSendHash)
func (_ZetaNonEth *ZetaNonEthFilterer) FilterMinted(opts *bind.FilterOpts, mintee []common.Address, internalSendHash [][32]byte) (*ZetaNonEthMintedIterator, error) {

	var minteeRule []interface{}
	for _, minteeItem := range mintee {
		minteeRule = append(minteeRule, minteeItem)
	}

	var internalSendHashRule []interface{}
	for _, internalSendHashItem := range internalSendHash {
		internalSendHashRule = append(internalSendHashRule, internalSendHashItem)
	}

	logs, sub, err := _ZetaNonEth.contract.FilterLogs(opts, "Minted", minteeRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return &ZetaNonEthMintedIterator{contract: _ZetaNonEth.contract, event: "Minted", logs: logs, sub: sub}, nil
}

// WatchMinted is a free log subscription operation binding the contract event 0xc263b302aec62d29105026245f19e16f8e0137066ccd4a8bd941f716bd4096bb.
//
// Solidity: event Minted(address indexed mintee, uint256 amount, bytes32 indexed internalSendHash)
func (_ZetaNonEth *ZetaNonEthFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *ZetaNonEthMinted, mintee []common.Address, internalSendHash [][32]byte) (event.Subscription, error) {

	var minteeRule []interface{}
	for _, minteeItem := range mintee {
		minteeRule = append(minteeRule, minteeItem)
	}

	var internalSendHashRule []interface{}
	for _, internalSendHashItem := range internalSendHash {
		internalSendHashRule = append(internalSendHashRule, internalSendHashItem)
	}

	logs, sub, err := _ZetaNonEth.contract.WatchLogs(opts, "Minted", minteeRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZetaNonEthMinted)
				if err := _ZetaNonEth.contract.UnpackLog(event, "Minted", log); err != nil {
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

// ParseMinted is a log parse operation binding the contract event 0xc263b302aec62d29105026245f19e16f8e0137066ccd4a8bd941f716bd4096bb.
//
// Solidity: event Minted(address indexed mintee, uint256 amount, bytes32 indexed internalSendHash)
func (_ZetaNonEth *ZetaNonEthFilterer) ParseMinted(log types.Log) (*ZetaNonEthMinted, error) {
	event := new(ZetaNonEthMinted)
	if err := _ZetaNonEth.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZetaNonEthTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ZetaNonEth contract.
type ZetaNonEthTransferIterator struct {
	Event *ZetaNonEthTransfer // Event containing the contract specifics and raw log

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
func (it *ZetaNonEthTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZetaNonEthTransfer)
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
		it.Event = new(ZetaNonEthTransfer)
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
func (it *ZetaNonEthTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZetaNonEthTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZetaNonEthTransfer represents a Transfer event raised by the ZetaNonEth contract.
type ZetaNonEthTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZetaNonEth *ZetaNonEthFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ZetaNonEthTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZetaNonEth.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ZetaNonEthTransferIterator{contract: _ZetaNonEth.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ZetaNonEth *ZetaNonEthFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ZetaNonEthTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ZetaNonEth.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZetaNonEthTransfer)
				if err := _ZetaNonEth.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ZetaNonEth *ZetaNonEthFilterer) ParseTransfer(log types.Log) (*ZetaNonEthTransfer, error) {
	event := new(ZetaNonEthTransfer)
	if err := _ZetaNonEth.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
