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

func (k Keeper) GetZetaInboundDetails(
	ctx sdk.Context,
	receiverChainID *big.Int,
	callOptions gatewayzevm.CallOptions,
) (InboundDetails, error) {
	return k.getZetaInboundDetails(ctx, receiverChainID, callOptions)
}

func (k Keeper) GetErc20InboundDetails(
	ctx sdk.Context,
	zrc20 ethcommon.Address,
	callEvent bool,
) (InboundDetails, error) {
	return k.getErc20InboundDetails(ctx, zrc20, callEvent)
}
