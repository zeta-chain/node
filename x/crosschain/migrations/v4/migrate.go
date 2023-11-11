package v4

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// MigrateStore migrates the x/crosschain module state from the consensus version 3 to 4
// It initializes the aborted zeta amount to 0
func MigrateStore(
	ctx sdk.Context,
	crossChainStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) error {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	prefixedStore := prefix.NewStore(ctx.KVStore(crossChainStoreKey), p)
	abortedAmount := sdkmath.ZeroUint()
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
		if val.CctxStatus.Status == types.CctxStatus_Aborted {
			abortedAmount = abortedAmount.Add(val.GetCurrentOutTxParam().Amount)
		}
	}
	b := cdc.MustMarshal(&types.AbortedZetaAmount{Amount: abortedAmount})
	store := ctx.KVStore(crossChainStoreKey)
	store.Set([]byte(types.AbortedZetaAmountKey), b)
	return nil
}
