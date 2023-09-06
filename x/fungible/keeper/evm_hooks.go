package keeper

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

var _ evmtypes.EvmHooks = EVMHooks{}

type EVMHooks struct {
	k Keeper
}

func (k Keeper) EVMHooks() EVMHooks {
	return EVMHooks{k}
}

// PostTxProcessing is a wrapper for calling the EVM PostTxProcessing hook on the module keeper
func (h EVMHooks) PostTxProcessing(ctx sdk.Context, _ core.Message, receipt *ethtypes.Receipt) error {
	return h.k.CheckPausedZRC20(ctx, receipt)
}

// CheckPausedZRC20 checks the events of the receipt
// if an event is emitted from a paused ZRC20 contract it will revert the transaction
func (k Keeper) CheckPausedZRC20(ctx sdk.Context, receipt *ethtypes.Receipt) error {
	// get non-duplicated list of all addresses that emitted logs
	var addresses []ethcommon.Address
	addressExist := make(map[ethcommon.Address]struct{})
	for _, log := range receipt.Logs {
		if _, ok := addressExist[log.Address]; !ok {
			addressExist[log.Address] = struct{}{}
			addresses = append(addresses, log.Address)
		}
	}

	// check if any of the addresses are from a paused ZRC20 contract
	for _, address := range addresses {
		fc, found := k.GetForeignCoins(ctx, address.Hex())
		if found {
			if fc.Paused {
				return cosmoserrors.Wrap(types.ErrPausedZRC20, address.Hex())
			}
		}
	}

	return nil
}
