package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// getEVMCallerAddress returns the caller address.
// Usually the caller is the contract.CallerAddress, which is the address of the contract that called the precompiled contract.
// If contract.CallerAddress != evm.Origin is true, it means the call was made through a contract,
// on which case there is a need to set the caller to the evm.Origin.
func GetEVMCallerAddress(evm *vm.EVM, contract *vm.Contract) (common.Address, error) {
	caller := contract.CallerAddress
	if contract.CallerAddress != evm.Origin {
		caller = evm.Origin
	}

	return caller, nil
}

// getCosmosAddress returns the counterpart cosmos address of the given ethereum address.
// It checks if the address is empty or blocked by the bank keeper.
func GetCosmosAddress(bankKeeper bank.Keeper, addr common.Address) (sdk.AccAddress, error) {
	toAddr := sdk.AccAddress(addr.Bytes())
	if toAddr.Empty() {
		return nil, &ErrInvalidAddr{
			Got:    toAddr.String(),
			Reason: "empty address",
		}
	}

	if bankKeeper.BlockedAddr(toAddr) {
		return nil, &ErrInvalidAddr{
			Got:    toAddr.String(),
			Reason: "destination address blocked by bank keeper",
		}
	}

	return toAddr, nil
}
