package keeper

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetGasPrice set a specific gasPrice in the store from its index
func (k Keeper) SetGasPrice(ctx sdk.Context, gasPrice types.GasPrice) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasPriceKey))
	b := k.cdc.MustMarshal(&gasPrice)
	gasPrice.Index = strconv.FormatInt(gasPrice.ChainId, 10)
	store.Set(types.KeyPrefix(gasPrice.Index), b)
}

// GetGasPrice returns a gasPrice from its index
func (k Keeper) GetGasPrice(ctx sdk.Context, chainID int64) (val types.GasPrice, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasPriceKey))
	b := store.Get(types.KeyPrefix(strconv.FormatInt(chainID, 10)))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetMedianGasPriceInUint(ctx sdk.Context, chainID int64) (sdk.Uint, bool) {
	gasPrice, isFound := k.GetGasPrice(ctx, chainID)
	if !isFound {
		return math.ZeroUint(), isFound
	}
	mi := gasPrice.MedianIndex
	return sdk.NewUint(gasPrice.Prices[mi]), true
}

// RemoveGasPrice removes a gasPrice from the store
func (k Keeper) RemoveGasPrice(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasPriceKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllGasPrice returns all gasPrice
func (k Keeper) GetAllGasPrice(ctx sdk.Context) (list []types.GasPrice) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasPriceKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.GasPrice
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Queries

func (k Keeper) GasPriceAll(c context.Context, req *types.QueryAllGasPriceRequest) (*types.QueryAllGasPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var gasPrices []*types.GasPrice
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	gasPriceStore := prefix.NewStore(store, types.KeyPrefix(types.GasPriceKey))

	pageRes, err := query.Paginate(gasPriceStore, req.Pagination, func(key []byte, value []byte) error {
		var gasPrice types.GasPrice
		if err := k.cdc.Unmarshal(value, &gasPrice); err != nil {
			return err
		}

		gasPrices = append(gasPrices, &gasPrice)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllGasPriceResponse{GasPrice: gasPrices, Pagination: pageRes}, nil
}

func (k Keeper) GasPrice(c context.Context, req *types.QueryGetGasPriceRequest) (*types.QueryGetGasPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	chainID, err := strconv.Atoi(req.Index)
	if err != nil {
		return nil, err
	}
	val, found := k.GetGasPrice(ctx, int64(chainID))
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetGasPriceResponse{GasPrice: &val}, nil
}

// MESSAGES

// Submit information about the connected chain's gas price at a specific block
// height. Gas price submitted by each validator is recorded separately and a
// median index is updated.
//
// Only observer validators are authorized to broadcast this message.
func (k msgServer) GasPriceVoter(goCtx context.Context, msg *types.MsgGasPriceVoter) (*types.MsgGasPriceVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, zetaObserverTypes.ErrSupportedChains
	}
	ok, err := k.IsAuthorized(ctx, msg.Creator, chain)
	if !ok {
		return nil, err
	}
	if chain == nil {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("ChainID : %d ", msg.ChainId))
	}

	gasPrice, isFound := k.GetGasPrice(ctx, chain.ChainId)
	if !isFound {
		gasPrice = types.GasPrice{
			Creator:     msg.Creator,
			Index:       strconv.FormatInt(chain.ChainId, 10), // TODO : Not needed index set at keeper
			ChainId:     chain.ChainId,
			Prices:      []uint64{msg.Price},
			BlockNums:   []uint64{msg.BlockNumber},
			Signers:     []string{msg.Creator},
			MedianIndex: 0,
		}
	} else {
		signers := gasPrice.Signers
		exist := false
		for i, s := range signers {
			if s == msg.Creator { // update existing entry
				gasPrice.BlockNums[i] = msg.BlockNumber
				gasPrice.Prices[i] = msg.Price
				exist = true
				break
			}
		}
		if !exist {
			gasPrice.Signers = append(gasPrice.Signers, msg.Creator)
			gasPrice.BlockNums = append(gasPrice.BlockNums, msg.BlockNumber)
			gasPrice.Prices = append(gasPrice.Prices, msg.Price)
		}
		// recompute the median gas price
		mi := medianOfArray(gasPrice.Prices)
		gasPrice.MedianIndex = uint64(mi)
	}
	k.SetGasPrice(ctx, gasPrice)
	chainIDBigINT := big.NewInt(chain.ChainId)
	gasUsed, err := k.fungibleKeeper.SetGasPrice(ctx, chainIDBigINT, big.NewInt(int64(gasPrice.Prices[gasPrice.MedianIndex])))
	if err != nil {
		return nil, err
	}

	// reset the gas count
	k.ResetGasMeterAndConsumeGas(ctx, gasUsed)

	return &types.MsgGasPriceVoterResponse{}, nil
}

type indexValue struct {
	Index int
	Value uint64
}

func medianOfArray(values []uint64) int {
	array := make([]indexValue, len(values))
	for i, v := range values {
		array[i] = indexValue{Index: i, Value: v}
	}
	sort.SliceStable(array, func(i, j int) bool {
		return array[i].Value < array[j].Value
	})
	l := len(array)
	return array[l/2].Index
}

// ResetGasMeterAndConsumeGas reset first the gas meter consumed value to zero and set it back to the new value
// 'gasUsed'
func (k *Keeper) ResetGasMeterAndConsumeGas(ctx sdk.Context, gasUsed uint64) {
	// reset the gas count
	ctx.GasMeter().RefundGas(ctx.GasMeter().GasConsumed(), "reset the gas count")
	ctx.GasMeter().ConsumeGas(gasUsed, "apply evm transaction")
}
