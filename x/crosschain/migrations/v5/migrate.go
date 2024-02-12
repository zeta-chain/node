package v5

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// crosschainKeeper is an interface to prevent cyclic dependency
type crosschainKeeper interface {
	GetStoreKey() storetypes.StoreKey
	GetCodec() codec.Codec
	GetAllCrossChainTx(ctx sdk.Context) []types.CrossChainTx
	AddFinalizedInbound(ctx sdk.Context, inboundTxHash string, senderChainID int64, height uint64)
}

// MigrateStore migrates the x/crosschain module state from the consensus version 4 to 5
// It resets the aborted zeta amount to use the inbound tx amount instead in situations where the outbound cctx is never created.
func MigrateStore(
	ctx sdk.Context,
	crosschainKeeper crosschainKeeper,
) error {
	SetZetaAccounting(ctx, crosschainKeeper.GetStoreKey(), crosschainKeeper.GetCodec())

	return nil
}

func SetZetaAccounting(
	ctx sdk.Context,
	crossChainStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
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
		if val.CctxStatus.Status == types.CctxStatus_Aborted && val.GetCurrentOutTxParam().CoinType == common.CoinType_Zeta {
			abortedValue := keeper.GetAbortedAmount(val)
			abortedAmountZeta = abortedAmountZeta.Add(abortedValue)
		}
	}
	b := cdc.MustMarshal(&types.ZetaAccounting{
		AbortedZetaAmount: abortedAmountZeta,
	})
	store := ctx.KVStore(crossChainStoreKey)
	store.Set([]byte(types.ZetaAccountingKey), b)
}
