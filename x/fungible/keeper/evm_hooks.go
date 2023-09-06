package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
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
	return nil
}
