package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetTSSVoter set a specific tSSVoter in the store from its index
func (k Keeper) SetTSSVoter(ctx sdk.Context, tSSVoter types.TSSVoter) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSVoterKey))
	b := k.cdc.MustMarshal(&tSSVoter)
	store.Set(types.KeyPrefix(tSSVoter.Index), b)
}

// GetTSSVoter returns a tSSVoter from its index
func (k Keeper) GetTSSVoter(ctx sdk.Context, index string) (val types.TSSVoter, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSVoterKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveTSSVoter removes a tSSVoter from the store
func (k Keeper) RemoveTSSVoter(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSVoterKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllTSSVoter returns all tSSVoter
func (k Keeper) GetAllTSSVoter(ctx sdk.Context) (list []types.TSSVoter) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSVoterKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TSSVoter
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

//Queries

func (k Keeper) TSSVoterAll(c context.Context, req *types.QueryAllTSSVoterRequest) (*types.QueryAllTSSVoterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var tSSVoters []*types.TSSVoter
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	tSSVoterStore := prefix.NewStore(store, types.KeyPrefix(types.TSSVoterKey))

	pageRes, err := query.Paginate(tSSVoterStore, req.Pagination, func(key []byte, value []byte) error {
		var tSSVoter types.TSSVoter
		if err := k.cdc.Unmarshal(value, &tSSVoter); err != nil {
			return err
		}

		tSSVoters = append(tSSVoters, &tSSVoter)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTSSVoterResponse{TSSVoter: tSSVoters, Pagination: pageRes}, nil
}

func (k Keeper) TSSVoter(c context.Context, req *types.QueryGetTSSVoterRequest) (*types.QueryGetTSSVoterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTSSVoter(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetTSSVoterResponse{TSSVoter: &val}, nil
}

// MESSAGES

func (k msgServer) CreateTSSVoter(goCtx context.Context, msg *types.MsgCreateTSSVoter) (*types.MsgCreateTSSVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !IsBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	msgDigest := msg.Digest()
	sessionID := ctx.BlockHeight() / 1000 * 1000
	index := crypto.Keccak256Hash([]byte(msgDigest), []byte(fmt.Sprintf("%d", sessionID)))
	// Check if the value already exists
	tssVoter, isFound := k.GetTSSVoter(ctx, index.Hex())

	if isDuplicateSigner(msg.Creator, tssVoter.Signers) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s double signing!!", msg.Creator))
	}

	if isFound {
		tssVoter.Signers = append(tssVoter.Signers, msg.Creator)
	} else {
		tssVoter = types.TSSVoter{
			Creator:         msg.Creator,
			Index:           index.Hex(),
			Chain:           msg.Chain,
			Address:         msg.Address,
			Pubkey:          msg.Pubkey,
			Signers:         []string{msg.Creator},
			FinalizedHeight: 0,
		}
	}

	k.SetTSSVoter(ctx, tssVoter)

	// this needs full consensus on all validators.
	if len(tssVoter.Signers) == len(validators) {
		tss := types.TSS{
			Creator:             "",
			Index:               tssVoter.Chain,
			Chain:               tssVoter.Chain,
			Address:             tssVoter.Address,
			Pubkey:              tssVoter.Pubkey,
			Signer:              tssVoter.Signers,
			FinalizedZetaHeight: uint64(ctx.BlockHeader().Height),
		}
		k.SetTSS(ctx, tss)
	}

	return &types.MsgCreateTSSVoterResponse{}, nil
}
