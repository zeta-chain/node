package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
)

// These functions are exported for testing purposes

func (k Keeper) ExecuteWithMintedZeta(
	ctx sdk.Context,
	amount *big.Int,
	operation func(sdk.Context) (*evmtypes.MsgEthereumTxResponse, bool, error),
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	return k.executeWithMintedZeta(ctx, amount, operation)
}
