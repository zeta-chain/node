package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
)

// These functions are exported for testing purposes

func (k Keeper) UpdateZetaAccounting(ctx sdk.Context, cctx types.CrossChainTx) {
	k.updateZetaAccounting(ctx, cctx)
}

func (k Keeper) UpdateInboundHashToCCTX(ctx sdk.Context, cctx types.CrossChainTx) {
	k.updateInboundHashToCCTX(ctx, cctx)
}

func (k Keeper) SetNonceToCCTX(ctx sdk.Context, cctx types.CrossChainTx, tssPubkey string) {
	k.setNonceToCCTX(ctx, cctx, tssPubkey)
}

func (k Keeper) GetNextCctxCounter(ctx sdk.Context) uint64 {
	return k.getNextCctxCounter(ctx)
}

func (k Keeper) UseRemainingGasFee(ctx sdk.Context, cctx *types.CrossChainTx) error {
	return k.useRemainingGasFee(ctx, cctx)
}

func (k Keeper) GetZETAInboundDetails(
	ctx sdk.Context,
	receiverChainID *big.Int,
	callOptions gatewayzevm.CallOptions,
) (InboundDetails, error) {
	return k.getZETAInboundDetails(ctx, receiverChainID, callOptions)
}

func (k Keeper) GetERC20InboundDetails(
	ctx sdk.Context,
	zrc20 ethcommon.Address,
	callEvent bool,
) (InboundDetails, error) {
	return k.getZRC20InboundDetails(ctx, zrc20, callEvent)
}
