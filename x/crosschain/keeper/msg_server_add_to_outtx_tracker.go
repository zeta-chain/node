package keeper

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	cosmoserrors "cosmossdk.io/errors"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// AddToOutTxTracker adds a new record to the outbound transaction tracker.
// only the admin policy account and the observer validators are authorized to broadcast this message without proof.
// If no pending cctx is found, the tracker is removed, if there is an existed tracker with the nonce & chainID.
func (k msgServer) AddToOutTxTracker(goCtx context.Context, msg *types.MsgAddToOutTxTracker) (*types.MsgAddToOutTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, observertypes.ErrSupportedChains
	}

	cctx, err := k.CctxByNonce(ctx, &types.QueryGetCctxByNonceRequest{
		ChainID: msg.ChainId,
		Nonce:   msg.Nonce,
	})
	if err != nil || cctx == nil || cctx.CrossChainTx == nil {
		return nil, cosmoserrors.Wrap(types.ErrCannotFindCctx, "cannot add out tx: no corresponding cctx found")
	}
	if !IsPending(*cctx.CrossChainTx) {
		k.RemoveOutTxTracker(ctx, msg.ChainId, msg.Nonce)
		return &types.MsgAddToOutTxTrackerResponse{}, nil
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

	isProven := false
	if msg.Proof != nil { // verify proof when it is provided
		txBytes, err := k.VerifyProof(ctx, msg.Proof, msg.ChainId, msg.BlockHash, msg.TxIndex)
		if err != nil {
			return nil, types.ErrProofVerificationFail.Wrapf(err.Error())
		}
		err = k.VerifyOutTxBody(ctx, msg, txBytes)
		if err != nil {
			return nil, types.ErrTxBodyVerificationFail.Wrapf(err.Error())
		}
		isProven = true
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
			if isProven {
				hash.Proved = true
				k.SetOutTxTracker(ctx, tracker)
				k.Logger(ctx).Info("Proof'd outbound transaction")
				return &types.MsgAddToOutTxTrackerResponse{}, nil
			}
			break
		}
	}
	if !isDup {
		if isProven {
			hash.Proved = true
			tracker.HashList = append([]*types.TxHashList{&hash}, tracker.HashList...)
			k.Logger(ctx).Info("Proof'd outbound transaction")
		} else if len(tracker.HashList) < 2 {
			tracker.HashList = append(tracker.HashList, &hash)
		}
		k.SetOutTxTracker(ctx, tracker)
	}
	return &types.MsgAddToOutTxTrackerResponse{}, nil
}

func (k Keeper) VerifyOutTxBody(ctx sdk.Context, msg *types.MsgAddToOutTxTracker, txBytes []byte) error {
	// get tss address
	tss, err := k.GetTssAddress(ctx, &types.QueryGetTssAddressRequest{})
	if err != nil {
		return err
	}

	// verify message against transaction body
	if common.IsEVMChain(msg.ChainId) {
		err = VerifyEVMOutTxBody(msg, txBytes, tss.Eth)
	} else if common.IsBitcoinChain(msg.ChainId) {
		err = VerifyBTCOutTxBody(msg, txBytes, tss.Btc)
	} else {
		return fmt.Errorf("cannot verify outTx body for chain %d", msg.ChainId)
	}
	return err
}

// VerifyEVMOutTxBody validates the sender address, nonce, chain id and tx hash.
// Note: 'msg' may contain fabricated information
func VerifyEVMOutTxBody(msg *types.MsgAddToOutTxTracker, txBytes []byte, tssEth string) error {
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

// VerifyBTCOutTxBody validates the SegWit sender address, nonce and chain id and tx hash
// Note: 'msg' may contain fabricated information
func VerifyBTCOutTxBody(msg *types.MsgAddToOutTxTracker, txBytes []byte, tssBtc string) error {
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
