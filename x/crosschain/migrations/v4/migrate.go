package v4

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// crosschainKeeper is an interface to prevent cyclic dependency
type crosschainKeeper interface {
	GetStoreKey() storetypes.StoreKey
	GetCodec() codec.Codec
	GetAllCrossChainTx(ctx sdk.Context) []types.CrossChainTx
	AddFinalizedInbound(ctx sdk.Context, inboundTxHash string, senderChainID int64, height uint64)
}

// MigrateStore migrates the x/crosschain module state from the consensus version 3 to 4
// It initializes the aborted zeta amount to 0
func MigrateStore(
	ctx sdk.Context,
	observerKeeper types.ObserverKeeper,
	crosschainKeeper crosschainKeeper,
) error {
	SetZetaAccounting(ctx, crosschainKeeper.GetStoreKey(), crosschainKeeper.GetCodec())
	MoveTssToObserverModule(ctx, observerKeeper, crosschainKeeper.GetStoreKey(), crosschainKeeper.GetCodec())
	MoveNonceToObserverModule(ctx, observerKeeper, crosschainKeeper.GetStoreKey(), crosschainKeeper.GetCodec())
	SetBitcoinFinalizedInbound(ctx, crosschainKeeper)

	return nil
}

func SetZetaAccounting(
	ctx sdk.Context,
	crossChainStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.CCTXKey))
	prefixedStore := prefix.NewStore(ctx.KVStore(crossChainStoreKey), p)
	abortedAmountZeta := sdkmath.ZeroUint()
	iterator := sdk.KVStorePrefixIterator(prefixedStore, []byte{})
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var val types.CrossChainTx
		cdc.MustUnmarshal(iterator.Value(), &val)
		if val.CctxStatus.Status == types.CctxStatus_Aborted && val.InboundParams.CoinType == coin.CoinType_Zeta {
			abortedAmountZeta = abortedAmountZeta.Add(val.GetCurrentOutboundParam().Amount)
		}
	}
	b := cdc.MustMarshal(&types.ZetaAccounting{
		AbortedZetaAmount: abortedAmountZeta,
	})
	store := ctx.KVStore(crossChainStoreKey)
	store.Set([]byte(types.ZetaAccountingKey), b)
}

func MoveNonceToObserverModule(
	ctx sdk.Context,
	observerKeeper types.ObserverKeeper,
	crossChainStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) {
	var chainNonces []observertypes.ChainNonces
	var pendingNonces []observertypes.PendingNonces
	var nonceToCcTx []observertypes.NonceToCctx
	store := prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(observertypes.ChainNoncesKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			return
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var val observertypes.ChainNonces
		err := cdc.Unmarshal(iterator.Value(), &val)
		if err == nil {
			chainNonces = append(chainNonces, val)
		}
	}
	store = prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(observertypes.PendingNoncesKeyPrefix))
	iterator = sdk.KVStorePrefixIterator(store, []byte{})
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			return
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var val observertypes.PendingNonces
		err := cdc.Unmarshal(iterator.Value(), &val)
		if err == nil {
			pendingNonces = append(pendingNonces, val)
		}
	}
	store = prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(observertypes.NonceToCctxKeyPrefix))
	iterator = sdk.KVStorePrefixIterator(store, []byte{})
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			return
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var val observertypes.NonceToCctx
		err := cdc.Unmarshal(iterator.Value(), &val)
		if err == nil {
			nonceToCcTx = append(nonceToCcTx, val)
		}
	}
	for _, c := range chainNonces {
		observerKeeper.SetChainNonces(ctx, c)
	}
	for _, p := range pendingNonces {
		observerKeeper.SetPendingNonces(ctx, p)
	}
	for _, n := range nonceToCcTx {
		observerKeeper.SetNonceToCctx(ctx, n)
	}

}

func MoveTssToObserverModule(ctx sdk.Context,
	observerKeeper types.ObserverKeeper,
	crossChainStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) {
	// Using New Types from observer module as the structure is the same
	var tss observertypes.TSS
	var tssHistory []observertypes.TSS

	writeTss := false

	// Fetch data from cross chain store using the legacy keys directly
	store := prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(observertypes.TSSKey))
	b := store.Get([]byte{0})
	if b != nil {
		err := cdc.Unmarshal(b, &tss)
		if err == nil {
			writeTss = true
		}
	}

	store = prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(observertypes.TSSHistoryKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val observertypes.TSS
		err := cdc.Unmarshal(iterator.Value(), &val)
		if err == nil {
			tssHistory = append(tssHistory, val)
		}
	}

	for _, t := range tssHistory {
		observerKeeper.SetTSSHistory(ctx, t)
	}
	if writeTss {
		observerKeeper.SetTSS(ctx, tss)
	}
}

// SetBitcoinFinalizedInbound sets the finalized inbound for bitcoin chains to prevent new ballots from being created with same intxhash
func SetBitcoinFinalizedInbound(ctx sdk.Context, crosschainKeeper crosschainKeeper) {
	for _, cctx := range crosschainKeeper.GetAllCrossChainTx(ctx) {
		if cctx.InboundParams != nil {
			// check if bitcoin inbound
			if chains.IsBitcoinChain(cctx.InboundParams.SenderChainId) {
				// add finalized inbound
				crosschainKeeper.AddFinalizedInbound(
					ctx,
					cctx.InboundParams.ObservedHash,
					cctx.InboundParams.SenderChainId,
					0,
				)
			}
		}
	}
}
