package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// TODO https://github.com/zeta-chain/node/issues/1269
func (k msgServer) AddToInTxTracker(goCtx context.Context, msg *types.MsgAddToInTxTracker) (*types.MsgAddToInTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, observertypes.ErrSupportedChains
	}

	adminPolicyAccount := k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_group1)
	isAdmin := msg.Creator == adminPolicyAccount
	isObserver := k.zetaObserverKeeper.IsAuthorized(ctx, msg.Creator, chain)

	isProven := false
	if !(isAdmin || isObserver) && msg.Proof != nil {
		txBytes, err := k.VerifyProof(ctx, msg.Proof, msg.ChainId, msg.BlockHash, msg.TxIndex)
		if err != nil {
			return nil, types.ErrCannotVerifyProof.Wrapf(err.Error())
		}

		if common.IsEVMChain(msg.ChainId) {
			err = k.VerifyEVMInTxBody(ctx, msg, txBytes)
			if err != nil {
				return nil, types.ErrCannotVerifyProof.Wrapf(err.Error())
			}
		} else {
			return nil, types.ErrCannotVerifyProof.Wrapf(fmt.Sprintf("cannot verify inTx body for chain %d", msg.ChainId))
		}
		isProven = true
	}

	// Sender needs to be either the admin policy account or an observer
	if !(isAdmin || isObserver || isProven) {
		return nil, errorsmod.Wrap(observertypes.ErrNotAuthorized, fmt.Sprintf("Creator %s", msg.Creator))
	}

	k.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  msg.ChainId,
		TxHash:   msg.TxHash,
		CoinType: msg.CoinType,
	})
	return &types.MsgAddToInTxTrackerResponse{}, nil
}

