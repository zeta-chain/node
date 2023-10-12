package keeper

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	cosmoserrors "cosmossdk.io/errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getOutTrackerIndex(chainID int64, nonce uint64) string {
	return fmt.Sprintf("%d-%d", chainID, nonce)
}

// SetOutTxTracker set a specific outTxTracker in the store from its index
func (k Keeper) SetOutTxTracker(ctx sdk.Context, outTxTracker types.OutTxTracker) {
	outTxTracker.Index = getOutTrackerIndex(outTxTracker.ChainId, outTxTracker.Nonce)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutTxTrackerKeyPrefix))
	b := k.cdc.MustMarshal(&outTxTracker)
	store.Set(types.OutTxTrackerKey(
		outTxTracker.Index,
	), b)
}

// GetOutTxTracker returns a outTxTracker from its index
func (k Keeper) GetOutTxTracker(
	ctx sdk.Context,
	chainID int64,
	nonce uint64,

) (val types.OutTxTracker, found bool) {
	index := getOutTrackerIndex(chainID, nonce)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutTxTrackerKeyPrefix))

	b := store.Get(types.OutTxTrackerKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOutTxTracker removes a outTxTracker from the store
func (k Keeper) RemoveOutTxTracker(
	ctx sdk.Context,
	chainID int64,
	nonce uint64,

) {
	index := getOutTrackerIndex(chainID, nonce)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutTxTrackerKeyPrefix))
	store.Delete(types.OutTxTrackerKey(
		index,
	))
}

// GetAllOutTxTracker returns all outTxTracker
func (k Keeper) GetAllOutTxTracker(ctx sdk.Context) (list []types.OutTxTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutTxTrackerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OutTxTracker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Queries

func (k Keeper) OutTxTrackerAll(c context.Context, req *types.QueryAllOutTxTrackerRequest) (*types.QueryAllOutTxTrackerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var outTxTrackers []types.OutTxTracker
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	outTxTrackerStore := prefix.NewStore(store, types.KeyPrefix(types.OutTxTrackerKeyPrefix))
	pageRes, err := query.Paginate(outTxTrackerStore, req.Pagination, func(key []byte, value []byte) error {
		var outTxTracker types.OutTxTracker
		if err := k.cdc.Unmarshal(value, &outTxTracker); err != nil {
			return err
		}

		outTxTrackers = append(outTxTrackers, outTxTracker)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllOutTxTrackerResponse{OutTxTracker: outTxTrackers, Pagination: pageRes}, nil
}

func (k Keeper) OutTxTrackerAllByChain(c context.Context, req *types.QueryAllOutTxTrackerByChainRequest) (*types.QueryAllOutTxTrackerByChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var outTxTrackers []types.OutTxTracker
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutTxTrackerKeyPrefix))
	chainStore := prefix.NewStore(store, types.KeyPrefix(fmt.Sprintf("%d-", req.Chain)))

	pageRes, err := query.Paginate(chainStore, req.Pagination, func(key []byte, value []byte) error {
		var outTxTracker types.OutTxTracker
		if err := k.cdc.Unmarshal(value, &outTxTracker); err != nil {
			return err
		}
		if outTxTracker.ChainId == req.Chain {
			outTxTrackers = append(outTxTrackers, outTxTracker)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOutTxTrackerByChainResponse{OutTxTracker: outTxTrackers, Pagination: pageRes}, nil
}

func (k Keeper) OutTxTracker(c context.Context, req *types.QueryGetOutTxTrackerRequest) (*types.QueryGetOutTxTrackerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	val, found := k.GetOutTxTracker(
		ctx,
		req.ChainID,
		req.Nonce,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetOutTxTrackerResponse{OutTxTracker: val}, nil
}

// Messages

// AddToOutTxTracker adds a new record to the outbound transaction tracker.
// only the admin policy account and the observer validators are authorized to broadcast this message.
func (k msgServer) AddToOutTxTracker(goCtx context.Context, msg *types.MsgAddToOutTxTracker) (*types.MsgAddToOutTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, observertypes.ErrSupportedChains
	}

	if msg.Proof == nil { // without proof, only certain accounts can send this message
		adminPolicyAccount := k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_group1)
		isAdmin := msg.Creator == adminPolicyAccount

		isObserver := k.zetaObserverKeeper.IsAuthorized(ctx, msg.Creator, chain)

		// Sender needs to be either the admin policy account or an observer
		if !(isAdmin || isObserver) {
			return nil, cosmoserrors.Wrap(observertypes.ErrNotAuthorized, fmt.Sprintf("Creator %s", msg.Creator))
		}
	}

	proven := false
	if msg.Proof != nil {
		blockHash, err := common.StringToHash(msg.ChainId, msg.BlockHash)
		if err != nil {
			return nil, cosmoserrors.Wrap(err, "block hash conversion failed")
		}
		res, found := k.zetaObserverKeeper.GetBlockHeader(ctx, blockHash)
		if !found {
			return nil, cosmoserrors.Wrap(observertypes.ErrBlockHeaderNotFound, fmt.Sprintf("block header not found %s", msg.BlockHash))
		}

		// verify outTx merkle proof
		txBytes, err := msg.Proof.Verify(res.Header, int(msg.TxIndex))
		if err != nil && !common.IsErrorInvalidProof(err) {
			return nil, err
		}
		if err == nil {
			tss, err := k.GetTssAddress(ctx, &types.QueryGetTssAddressRequest{})
			if err != nil {
				return nil, err
			}
			// verify outTx transaction body
			if common.IsEVMChain(msg.ChainId) {
				err = ValidateEVMOutTxBody(msg, txBytes, tss.Eth)
			} else if common.IsBitcoinChain(msg.ChainId) {
				err = ValidateBTCOutTxBody(msg, txBytes, tss.Btc)
			} else {
				return nil, fmt.Errorf("unsupported chain id %d", msg.ChainId)
			}
			if err != nil {
				return nil, err
			}
		}

		if !proven {
			return nil, fmt.Errorf("proof failed")
		}
	}

	tracker, found := k.GetOutTxTracker(ctx, msg.ChainId, msg.Nonce)
	hash := types.TxHashList{
		TxHash:   msg.TxHash,
		TxSigner: msg.Creator,
	}
	if !found {
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			Index:    "",
			ChainId:  chain.ChainId,
			Nonce:    msg.Nonce,
			HashList: []*types.TxHashList{&hash},
		})
		return &types.MsgAddToOutTxTrackerResponse{}, nil
	}

	var isDup = false
	for _, hash := range tracker.HashList {
		if strings.EqualFold(hash.TxHash, msg.TxHash) {
			isDup = true
			if proven {
				hash.Proved = true
				k.SetOutTxTracker(ctx, tracker)
				k.Logger(ctx).Info("Proof'd outbound transaction")
				return &types.MsgAddToOutTxTrackerResponse{}, nil
			}
			break
		}
	}
	if !isDup {
		if proven {
			hash.Proved = true
			tracker.HashList = append([]*types.TxHashList{&hash}, tracker.HashList...)
			k.Logger(ctx).Info("Proof'd outbound transaction")
		} else {
			tracker.HashList = append(tracker.HashList, &hash)
		}
		k.SetOutTxTracker(ctx, tracker)
	}
	return &types.MsgAddToOutTxTrackerResponse{}, nil
}

// ValidateEVMOutTxBody validates the sender address, nonce and chain ID.
// Note: 'msg' may contain fabricated information
func ValidateEVMOutTxBody(msg *types.MsgAddToOutTxTracker, txBytes []byte, tssEth string) error {
	var txx ethtypes.Transaction
	err := txx.UnmarshalBinary(txBytes)
	if err != nil {
		return err
	}
	signer := ethtypes.NewLondonSigner(txx.ChainId())
	sender, err := ethtypes.Sender(signer, &txx)
	if err != nil {
		return err
	}
	tssAddr := eth.HexToAddress(tssEth)
	if tssAddr == (eth.Address{}) {
		return fmt.Errorf("tss address not found")
	}
	if sender != tssAddr {
		return fmt.Errorf("sender %s is not tss address", sender)
	}
	if txx.ChainId().Cmp(big.NewInt(msg.ChainId)) != 0 {
		return fmt.Errorf("want evm chain id %d, got %d", txx.ChainId(), msg.ChainId)
	}
	if txx.Nonce() != msg.Nonce {
		return fmt.Errorf("want nonce %d, got %d", txx.Nonce(), msg.Nonce)
	}
	if txx.Hash().Hex() != msg.TxHash {
		return fmt.Errorf("want tx hash %s, got %s", txx.Hash().Hex(), msg.TxHash)
	}
	return nil
}

// ValidateBTCOutTxBody validates the SegWit sender address, nonce and chain ID.
// Note: 'msg' may contain fabricated information
func ValidateBTCOutTxBody(msg *types.MsgAddToOutTxTracker, txBytes []byte, tssBtc string) error {
	tx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return err
	}
	for _, vin := range tx.MsgTx().TxIn {
		if len(vin.Witness) != 2 { // outTx is SegWit transaction for now
			return fmt.Errorf("not a SegWit transaction")
		}
		pubKey, err := btcec.ParsePubKey(vin.Witness[1], btcec.S256())
		if err != nil {
			return fmt.Errorf("failed to parse public key")
		}
		addrP2WPKH, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), config.BitconNetParams)
		if err != nil {
			return fmt.Errorf("failed to create P2WPKH address")
		}
		if addrP2WPKH.EncodeAddress() != tssBtc {
			return fmt.Errorf("sender %s is not tss address", addrP2WPKH.EncodeAddress())
		}
	}
	if common.BtcChainID() != msg.ChainId {
		return fmt.Errorf("want btc chain id %d, got %d", common.BtcChainID(), msg.ChainId)
	}
	if len(tx.MsgTx().TxOut) < 1 {
		return fmt.Errorf("outTx should have at least one output")
	}
	if tx.MsgTx().TxOut[0].Value != common.NonceMarkAmount(msg.Nonce) {
		return fmt.Errorf("want nonce mark %d, got %d", tx.MsgTx().TxOut[0].Value, common.NonceMarkAmount(msg.Nonce))
	}
	if tx.MsgTx().TxHash().String() != msg.TxHash {
		return fmt.Errorf("want tx hash %s, got %s", tx.MsgTx().TxHash(), msg.TxHash)
	}
	return nil
}

// RemoveFromOutTxTracker removes a record from the outbound transaction tracker by chain ID and nonce.
// only the admin policy account is authorized to broadcast this message.
func (k msgServer) RemoveFromOutTxTracker(goCtx context.Context, msg *types.MsgRemoveFromOutTxTracker) (*types.MsgRemoveFromOutTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_group1) {
		return &types.MsgRemoveFromOutTxTrackerResponse{}, observertypes.ErrNotAuthorizedPolicy
	}

	k.RemoveOutTxTracker(ctx, msg.ChainId, msg.Nonce)
	return &types.MsgRemoveFromOutTxTrackerResponse{}, nil
}
