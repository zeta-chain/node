package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// TODO https://github.com/zeta-chain/node/issues/1269
func (k Keeper) AddToInTxTracker(goCtx context.Context, msg *types.MsgAddToInTxTracker) (*types.MsgAddToInTxTrackerResponse, error) {
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
		err = k.VerifyInTxBody(ctx, msg, txBytes)
		if err != nil {
			return nil, types.ErrCannotVerifyProof.Wrapf(err.Error())
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

// https://github.com/zeta-chain/node/issues/1254
func (k Keeper) VerifyInTxBody(ctx sdk.Context, msg *types.MsgAddToInTxTracker, txBytes []byte) error {
	// get core params and tss address
	coreParams, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, msg.ChainId)
	if !found {
		return types.ErrUnsupportedChain.Wrapf("core params not found for chain %d", msg.ChainId)
	}
	err := error(nil)
	// verify message against transaction body
	if common.IsEVMChain(msg.ChainId) {
		err = k.VerifyEVMInTxBody(ctx, coreParams, msg, txBytes)
	} else {
		return fmt.Errorf("cannot verify inTx body for chain %d", msg.ChainId)
	}
	return err
}

func (k Keeper) VerifyEVMInTxBody(ctx sdk.Context, coreParams *observertypes.CoreParams, msg *types.MsgAddToInTxTracker, txBytes []byte) error {
	var txx ethtypes.Transaction
	err := txx.UnmarshalBinary(txBytes)
	if err != nil {
		return err
	}
	switch msg.CoinType {
	case common.CoinType_Zeta:
		if txx.To().Hex() != coreParams.ConnectorContractAddress {
			return fmt.Errorf("receiver is not connector contract for coin type %s", msg.CoinType)
		}
		return nil
	case common.CoinType_ERC20:
		if txx.To().Hex() != coreParams.Erc20CustodyContractAddress {
			return fmt.Errorf("receiver is not erc20Custory contract for coin type %s", msg.CoinType)
		}
		return nil
	case common.CoinType_Gas:
		tss, err := k.GetTssAddress(ctx, &types.QueryGetTssAddressRequest{})
		if err != nil {
			return err
		}
		tssAddr := eth.HexToAddress(tss.Eth)
		if tssAddr == (eth.Address{}) {
			return fmt.Errorf("tss address not found")
		}
		if txx.To().Hex() != tssAddr.Hex() {
			return fmt.Errorf("receiver is not tssAddress contract for coin type %s", msg.CoinType)
		}
		return nil
	default:
		return fmt.Errorf("coin type %s not supported", msg.CoinType)
	}
}
