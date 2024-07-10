package v5

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// crosschainKeeper is an interface to prevent cyclic dependency
type crosschainKeeper interface {
	GetStoreKey() storetypes.StoreKey
	GetCodec() codec.Codec
	GetAllCrossChainTx(ctx sdk.Context) []types.CrossChainTx

	SetCrossChainTx(ctx sdk.Context, cctx types.CrossChainTx)
	AddFinalizedInbound(ctx sdk.Context, inboundTxHash string, senderChainID int64, height uint64)

	SetZetaAccounting(ctx sdk.Context, accounting types.ZetaAccounting)
}

// MigrateStore migrates the x/crosschain module state from the consensus version 4 to 5
// It resets the aborted zeta amount to use the inbound tx amount instead in situations where the outbound cctx is never created.
func MigrateStore(ctx sdk.Context, crosschainKeeper crosschainKeeper, observerKeeper types.ObserverKeeper) error {
	err := SetZetaAccounting(ctx, crosschainKeeper, observerKeeper)
	if err != nil {
		return err
	}
	ResetTestnetNonce(ctx, observerKeeper)

	return nil
}

func ResetTestnetNonce(
	ctx sdk.Context,
	observerKeeper types.ObserverKeeper,
) {
	tss, found := observerKeeper.GetTSS(ctx)
	if !found {
		ctx.Logger().Info("ResetTestnetNonce: TSS not found")
		return
	}
	for _, chainNonce := range CurrentTestnetChains() {
		cn, found := observerKeeper.GetChainNonces(ctx, chainNonce.chain.ChainName.String())
		if !found {
			ctx.Logger().Info("ResetTestnetNonce: Chain nonce not found", "chain", chainNonce.chain.ChainName.String())
			continue
		}

		ctx.Logger().Info("ResetTestnetNonce: Resetting chain nonce", "chain", chainNonce.chain.ChainName.String())

		cn.Nonce = chainNonce.nonceHigh
		observerKeeper.SetChainNonces(ctx, cn)

		pn, found := observerKeeper.GetPendingNonces(ctx, tss.TssPubkey, chainNonce.chain.ChainId)
		if !found {
			continue
		}
		// #nosec G115 always in range for testnet chains
		pn.NonceLow = int64(chainNonce.nonceLow)
		// #nosec G115 always in range for testnet chains
		pn.NonceHigh = int64(chainNonce.nonceHigh)
		observerKeeper.SetPendingNonces(ctx, pn)
	}
}

type TestnetNonce struct {
	chain     chains.Chain
	nonceHigh uint64
	nonceLow  uint64
}

func CurrentTestnetChains() []TestnetNonce {
	return []TestnetNonce{
		{chain: chains.Goerli, nonceHigh: 226841, nonceLow: 226841},
		{chain: chains.Mumbai, nonceHigh: 200599, nonceLow: 200599},
		{chain: chains.BscTestnet, nonceHigh: 110454, nonceLow: 110454},
		{chain: chains.BitcoinTestnet, nonceHigh: 4881, nonceLow: 4881},
	}
}

func SetZetaAccounting(
	ctx sdk.Context,
	crosschainKeeper crosschainKeeper,
	observerKeeper types.ObserverKeeper,
) error {
	ccctxList := crosschainKeeper.GetAllCrossChainTx(ctx)
	abortedAmountZeta := sdkmath.ZeroUint()
	for _, cctx := range ccctxList {
		if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
			switch cctx.InboundParams.CoinType {
			case coin.CoinType_ERC20:
				{
					receiverChain, found := observerKeeper.GetSupportedChainFromChainID(
						ctx,
						cctx.GetCurrentOutboundParam().ReceiverChainId,
					)
					if !found {
						ctx.Logger().
							Error(fmt.Sprintf("Error getting chain from chain id: %d , cctx index", cctx.GetCurrentOutboundParam().ReceiverChainId), cctx.Index)
						continue
					}
					// There is a chance that this cctx has already been refunded, so we set the isRefunded flag to true.
					// Even though, there is a slight possibility that the refund tx failed when doing an auto refund; there is no way for us to know. Which is why we can mark this type of cctx as non-refundable
					// Auto refunds are done for ERC20 cctx's when the receiver chain is a zeta chain.
					if receiverChain.IsZetaChain() {
						cctx.CctxStatus.IsAbortRefunded = true
					} else {
						cctx.CctxStatus.IsAbortRefunded = false
					}
				}
			case coin.CoinType_Zeta:
				{
					// add the required amount into the zeta accounting.
					// GetAbortedAmount replaces using Outbound Amount directly, to make sure we refund the amount deposited by the user if the outbound is never created and the cctx is aborted.
					// For these cctx's we allow the refund to be processed later and the Aborted amount would be adjusted when the refund is processed.
					abortedValue := GetAbortedAmount(cctx)
					abortedAmountZeta = abortedAmountZeta.Add(abortedValue)
					cctx.CctxStatus.IsAbortRefunded = false
				}
			case coin.CoinType_Gas:
				{
					// CointType gas can be processed as normal and we can issue the refund using the admin refund tx .
					cctx.CctxStatus.IsAbortRefunded = false
				}
			}
			crosschainKeeper.SetCrossChainTx(ctx, cctx)
		}
	}
	crosschainKeeper.SetZetaAccounting(ctx, types.ZetaAccounting{AbortedZetaAmount: abortedAmountZeta})

	return nil
}
func GetAbortedAmount(cctx types.CrossChainTx) sdkmath.Uint {
	if cctx.OutboundParams != nil && !cctx.GetCurrentOutboundParam().Amount.IsZero() {
		return cctx.GetCurrentOutboundParam().Amount
	}
	if cctx.InboundParams != nil {
		return cctx.InboundParams.Amount
	}

	return sdkmath.ZeroUint()
}
