package keeper

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/zeta-chain/node/x/fungible/types"
)

var _ evmtypes.EvmHooks = EVMHooks{}

type EVMHooks struct {
	k Keeper
}

func (k Keeper) EVMHooks() EVMHooks {
	return EVMHooks{k}
}

// PostTxProcessing is a wrapper for calling the EVM PostTxProcessing hook on the module keeper
func (h EVMHooks) PostTxProcessing(
	ctx sdk.Context,
	_ ethcommon.Address,
	_ core.Message,
	receipt *ethtypes.Receipt,
) error {
	return h.k.checkPausedZRC20(ctx, receipt)
}

// checkPausedZRC20 checks the events of the receipt
// if an event is emitted from a paused ZRC20 contract it will revert the transaction
func (k Keeper) checkPausedZRC20(ctx sdk.Context, receipt *ethtypes.Receipt) error {
	if receipt == nil {
		return nil
	}

	// get non-duplicated list of all addresses that emitted logs
	var addresses []ethcommon.Address
	addressExist := make(map[ethcommon.Address]struct{})
	for _, log := range receipt.Logs {
		if log == nil {
			continue
		}
		if _, ok := addressExist[log.Address]; !ok {
			addressExist[log.Address] = struct{}{}
			addresses = append(addresses, log.Address)
		}
	}

	// check if any of the addresses are from a paused ZRC20 contract
	for _, address := range addresses {
		fc, found := k.GetForeignCoins(ctx, address.Hex())
		if found && fc.Paused {
			return cosmoserrors.Wrap(types.ErrPausedZRC20, address.Hex())
		}
	}

	return nil
}
