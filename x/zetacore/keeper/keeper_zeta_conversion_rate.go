package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"sort"
)

// SetZetaConversionRate set a specific zetaConversionRate in the store from its index
func (k Keeper) SetZetaConversionRate(ctx sdk.Context, zetaConversionRate types.ZetaConversionRate) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaConversionRateKeyPrefix))
	b := k.cdc.MustMarshal(&zetaConversionRate)
	store.Set(types.ZetaConversionRateKey(
		zetaConversionRate.Index,
	), b)
}

// GetZetaConversionRate returns a zetaConversionRate from its index
func (k Keeper) GetZetaConversionRate(
	ctx sdk.Context,
	index string,

) (val types.ZetaConversionRate, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaConversionRateKeyPrefix))

	b := store.Get(types.ZetaConversionRateKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveZetaConversionRate removes a zetaConversionRate from the store
func (k Keeper) RemoveZetaConversionRate(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaConversionRateKeyPrefix))
	store.Delete(types.ZetaConversionRateKey(
		index,
	))
}

// GetAllZetaConversionRate returns all zetaConversionRate
func (k Keeper) GetAllZetaConversionRate(ctx sdk.Context) (list []types.ZetaConversionRate) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaConversionRateKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ZetaConversionRate
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Queries

func (k Keeper) ZetaConversionRateAll(c context.Context, req *types.QueryAllZetaConversionRateRequest) (*types.QueryAllZetaConversionRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var zetaConversionRates []types.ZetaConversionRate
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	zetaConversionRateStore := prefix.NewStore(store, types.KeyPrefix(types.ZetaConversionRateKeyPrefix))

	pageRes, err := query.Paginate(zetaConversionRateStore, req.Pagination, func(key []byte, value []byte) error {
		var zetaConversionRate types.ZetaConversionRate
		if err := k.cdc.Unmarshal(value, &zetaConversionRate); err != nil {
			return err
		}

		zetaConversionRates = append(zetaConversionRates, zetaConversionRate)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllZetaConversionRateResponse{ZetaConversionRate: zetaConversionRates, Pagination: pageRes}, nil
}

func (k Keeper) ZetaConversionRate(c context.Context, req *types.QueryGetZetaConversionRateRequest) (*types.QueryGetZetaConversionRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetZetaConversionRate(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetZetaConversionRateResponse{ZetaConversionRate: val}, nil
}

// MESSAGES

func (k msgServer) ZetaConversionRateVoter(goCtx context.Context, msg *types.MsgZetaConversionRateVoter) (*types.MsgZetaConversionRateVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !IsBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	chain, err := common.ParseChain(msg.Chain)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrUnsupportedChain, fmt.Sprintf("chain %s not supported", msg.Chain))
	}
	rate, isFound := k.GetZetaConversionRate(ctx, chain.String())
	nativeTokenSymbol := chain.GetNativeTokenSymbol()
	if !isFound {
		rate = types.ZetaConversionRate{
			Index:               chain.String(),
			Chain:               chain.String(),
			Signers:             []string{msg.Creator},
			BlockNums:           []uint64{msg.BlockNumber},
			ZetaConversionRates: []string{msg.ZetaConversionRate},
			NativeTokenSymbol:   nativeTokenSymbol,
			MedianIndex:         0,
		}
	} else {
		signers := rate.Signers
		exist := false
		for i, s := range signers {
			if s == msg.Creator { // update existing entry
				rate.BlockNums[i] = msg.BlockNumber
				rate.ZetaConversionRates[i] = msg.ZetaConversionRate
				exist = true
				break
			}
		}
		if !exist {
			rate.Signers = append(rate.Signers, msg.Creator)
			rate.BlockNums = append(rate.BlockNums, msg.BlockNumber)
			rate.ZetaConversionRates = append(rate.ZetaConversionRates, msg.ZetaConversionRate)
		}
		mi := medianOfArrayFloat(rate.ZetaConversionRates)
		rate.MedianIndex = uint64(mi)
	}
	k.SetZetaConversionRate(ctx, rate)

	return &types.MsgZetaConversionRateVoterResponse{}, nil
}

type indexFloatValue struct {
	Index int
	Value *big.Int
}

func medianOfArrayFloat(values []string) int {
	var array []indexFloatValue
	for i, v := range values {
		f, ok := big.NewInt(0).SetString(v, 0) // should be less than 256bit
		if ok {
			array = append(array, indexFloatValue{Index: i, Value: f})
		} else {
			log.Error().Msgf("parse big.Int error")
		}
	}
	sort.SliceStable(array, func(i, j int) bool {
		return array[i].Value.Cmp(array[j].Value) < 0
	})
	l := len(array)
	return array[l/2].Index
}
