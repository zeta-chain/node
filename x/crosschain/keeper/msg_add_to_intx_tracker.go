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

func (k msgServer) AddToInTxTracker(goCtx context.Context, msg *types.MsgAddToInTxTracker) (*types.MsgAddToInTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, observertypes.ErrSupportedChains
	}

	adminPolicyAccount := k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observertypes.Policy_Type_group1)
	isAdmin := msg.Creator == adminPolicyAccount

	isObserver, err := k.zetaObserverKeeper.IsAuthorized(ctx, msg.Creator, chain)
	if err != nil {
		ctx.Logger().Error("Error while checking if the account is an observer", err)
	}
	isProven := false
	if !(isAdmin || isObserver) && msg.Proof != nil {
		txx, err := k.VerifyProof(ctx, msg.Proof, msg.BlockHash, msg.TxIndex, msg.ChainId)
		if err != nil {
			return nil, types.ErrCannotVerifyProof.Wrapf(err.Error())
		}
		err = k.VerifyInTxTrackerProof(ctx, txx, msg.ChainId, msg.CoinType)
		if err != nil {
			return nil, types.ErrCannotVerifyProof.Wrapf(err.Error())
		}
		isProven = true
	}

	// Sender needs to be either the admin policy account or an observer
	if !(isAdmin || isObserver || isProven) {
		return nil, errorsmod.Wrap(observertypes.ErrNotAuthorized, fmt.Sprintf("Creator %s", msg.Creator))
	}

	k.Keeper.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  msg.ChainId,
		TxHash:   msg.TxHash,
		CoinType: msg.CoinType,
	})
	return &types.MsgAddToInTxTrackerResponse{}, nil
}

// https://github.com/zeta-chain/node/issues/1254
func (k Keeper) VerifyInTxTrackerProof(ctx sdk.Context, txx ethtypes.Transaction, chainID int64, coinType common.CoinType) error {

	coreParams, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, chainID)
	if !found {
		return types.ErrUnsupportedChain.Wrapf("core params not found for chain %d", chainID)
	}
	tssRes, err := k.GetTssAddress(ctx, &types.QueryGetTssAddressRequest{})
	if err != nil {
		return err
	}
	tssAddr := eth.HexToAddress(tssRes.Eth)
	if tssAddr == (eth.Address{}) {
		return fmt.Errorf("tss address not found")
	}
	if common.IsEVMChain(chainID) {
		switch coinType {
		case common.CoinType_Zeta:
			if txx.To().Hex() != coreParams.ConnectorContractAddress {
				return fmt.Errorf("receiver is not connector contract for coin type %s", coinType)
			}
			return nil
		case common.CoinType_ERC20:
			if txx.To().Hex() != coreParams.Erc20CustodyContractAddress {
				return fmt.Errorf("receiver is not erc20Custory contract for coin type %s", coinType)
			}
			return nil
		case common.CoinType_Gas:
			if txx.To().Hex() != tssAddr.Hex() {
				return fmt.Errorf("receiver is not tssAddress contract for coin type %s", coinType)
			}
			return nil
		}
	}

	return fmt.Errorf("proof failed")
}
