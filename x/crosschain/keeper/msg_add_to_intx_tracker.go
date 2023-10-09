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
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) AddToInTxTracker(goCtx context.Context, msg *types.MsgAddToInTxTracker) (*types.MsgAddToInTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, observerTypes.ErrSupportedChains
	}

	adminPolicyAccount := k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observerTypes.Policy_Type_group1)
	isAdmin := msg.Creator == adminPolicyAccount

	isObserver, err := k.zetaObserverKeeper.IsAuthorized(ctx, msg.Creator, chain)
	if err != nil {
		ctx.Logger().Error("Error while checking if the account is an observer", err)
	}
	isProven := false
	if !(isAdmin || isObserver) && msg.Proof != nil {
		isProven, err = k.VerifyInTxTrackerProof(ctx, msg.Proof, msg.BlockHash, msg.TxIndex, msg.ChainId, msg.CoinType)
		if err != nil {
			return nil, types.ErrCannotVerifyProof.Wrapf(err.Error())
		}
	}

	// Sender needs to be either the admin policy account or an observer
	if !(isAdmin || isObserver || isProven) {
		return nil, errorsmod.Wrap(observerTypes.ErrNotAuthorized, fmt.Sprintf("Creator %s", msg.Creator))
	}

	k.Keeper.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  msg.ChainId,
		TxHash:   msg.TxHash,
		CoinType: msg.CoinType,
	})
	return &types.MsgAddToInTxTrackerResponse{}, nil
}

// https://github.com/zeta-chain/node/issues/1254
func (k Keeper) VerifyInTxTrackerProof(ctx sdk.Context, proof *common.Proof, hash string, txIndex int64, chainID int64, coinType common.CoinType) (bool, error) {
	if common.IsBitcoinChain(chainID) {
		return false, fmt.Errorf("cannot verify proof for bitcoin chain %d", chainID)
	}

	senderChain := common.GetChainFromChainID(chainID)
	if senderChain == nil {
		return false, types.ErrUnsupportedChain
	}

	if !senderChain.IsProvable() {
		return false, fmt.Errorf("chain %d does not support block header verification", chainID)
	}

	blockHash := eth.HexToHash(hash)
	res, found := k.zetaObserverKeeper.GetBlockHeader(ctx, blockHash.Bytes())
	if !found {
		return false, fmt.Errorf("block header not found %s", blockHash)
	}

	// verify and process the proof
	val, err := proof.Verify(res.Header, int(txIndex))
	if err != nil && !common.IsErrorInvalidProof(err) {
		return false, err
	}
	var txx ethtypes.Transaction
	err = txx.UnmarshalBinary(val)
	if err != nil {
		return false, err
	}

	coreParams, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, senderChain.ChainId)
	if !found {
		return false, types.ErrUnsupportedChain.Wrapf("core params not found for chain %d", senderChain.ChainId)
	}
	tssRes, err := k.GetTssAddress(ctx, &types.QueryGetTssAddressRequest{})
	if err != nil {
		return false, err
	}
	tssAddr := eth.HexToAddress(tssRes.Eth)
	if tssAddr == (eth.Address{}) {
		return false, fmt.Errorf("tss address not found")
	}
	if common.IsEVMChain(chainID) {
		switch coinType {
		case common.CoinType_Zeta:
			if txx.To().Hex() != coreParams.ConnectorContractAddress {
				return false, fmt.Errorf("receiver is not connector contract for coin type %s", coinType)
			}
			return true, nil
		case common.CoinType_ERC20:
			if txx.To().Hex() != coreParams.Erc20CustodyContractAddress {
				return false, fmt.Errorf("receiver is not erc20Custory contract for coin type %s", coinType)
			}
			return true, nil
		case common.CoinType_Gas:
			if txx.To().Hex() != tssAddr.Hex() {
				return false, fmt.Errorf("receiver is not tssAddress contract for coin type %s", coinType)
			}
			return true, nil
		}
	}

	return false, fmt.Errorf("proof failed")
}
