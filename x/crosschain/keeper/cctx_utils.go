package keeper

import (
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateNonce sets the CCTX outbound nonce to the next nonce, and updates the nonce of blockchain state.
// It also updates the PendingNonces that is used to track the unfulfilled outbound txs.
func (k Keeper) UpdateNonce(ctx sdk.Context, receiveChainID int64, cctx *types.CrossChainTx) error {
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(receiveChainID)
	if chain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}

	nonce, found := k.GetChainNonces(ctx, chain.ChainName.String())
	if !found {
		return cosmoserrors.Wrap(types.ErrCannotFindReceiverNonce, fmt.Sprintf("Chain(%s) | Identifiers : %s ", chain.ChainName.String(), cctx.LogIdentifierForCCTX()))
	}

	// SET nonce
	cctx.GetCurrentOutTxParam().OutboundTxTssNonce = nonce.Nonce
	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return cosmoserrors.Wrap(types.ErrCannotFindTSSKeys, fmt.Sprintf("Chain(%s) | Identifiers : %s ", chain.ChainName.String(), cctx.LogIdentifierForCCTX()))
	}

	p, found := k.GetPendingNonces(ctx, tss.TssPubkey, receiveChainID)
	if !found {
		return cosmoserrors.Wrap(types.ErrCannotFindPendingNonces, fmt.Sprintf("chain_id %d, nonce %d", receiveChainID, nonce.Nonce))
	}

	// #nosec G701 always in range
	if p.NonceHigh != int64(nonce.Nonce) {
		return cosmoserrors.Wrap(types.ErrNonceMismatch, fmt.Sprintf("chain_id %d, high nonce %d, current nonce %d", receiveChainID, p.NonceHigh, nonce.Nonce))
	}

	nonce.Nonce++
	p.NonceHigh++
	k.SetChainNonces(ctx, nonce)
	k.SetPendingNonces(ctx, p)
	return nil
}

// RefundAmountOnZetaChain refunds the amount of the cctx on ZetaChain in case of aborted cctx
// NOTE: GetCurrentOutTxParam should contain the last up to date cctx amount
func (k Keeper) RefundAmountOnZetaChain(ctx sdk.Context, cctx types.CrossChainTx, inputAmount math.Uint) error {
	// preliminary checks
	if cctx.InboundTxParams.CoinType != common.CoinType_ERC20 {
		return errors.New("unsupported coin type for refund on ZetaChain")
	}
	if !common.IsEVMChain(cctx.InboundTxParams.SenderChainId) {
		return errors.New("only EVM chains are supported for refund on ZetaChain")
	}
	sender := ethcommon.HexToAddress(cctx.InboundTxParams.Sender)
	if sender == (ethcommon.Address{}) {
		return errors.New("invalid sender address")
	}
	if inputAmount.IsNil() || inputAmount.IsZero() {
		return errors.New("no amount to refund")
	}

	// get address of the zrc20
	fc, found := k.fungibleKeeper.GetForeignCoinFromAsset(ctx, cctx.InboundTxParams.Asset, cctx.InboundTxParams.SenderChainId)
	if !found {
		return fmt.Errorf("asset %s zrc not found", cctx.InboundTxParams.Asset)
	}
	zrc20 := ethcommon.HexToAddress(fc.Zrc20ContractAddress)
	if zrc20 == (ethcommon.Address{}) {
		return fmt.Errorf("asset %s invalid zrc address", cctx.InboundTxParams.Asset)
	}

	// deposit the amount to the sender
	if _, err := k.fungibleKeeper.DepositZRC20(ctx, zrc20, sender, inputAmount.BigInt()); err != nil {
		return errors.New("failed to deposit zrc20 on ZetaChain" + err.Error())
	}

	return nil
}

// GetRevertGasLimit returns the gas limit for the revert transaction in a CCTX
// It returns 0 if there is no error but the gas limit can't be determined from the CCTX data
func (k Keeper) GetRevertGasLimit(ctx sdk.Context, cctx types.CrossChainTx) (uint64, error) {
	if cctx.InboundTxParams == nil {
		return 0, nil
	}

	if cctx.InboundTxParams.CoinType == common.CoinType_Gas {
		// get the gas limit of the gas token
		fc, found := k.fungibleKeeper.GetGasCoinForForeignCoin(ctx, cctx.InboundTxParams.SenderChainId)
		if !found {
			return 0, types.ErrForeignCoinNotFound
		}
		gasLimit, err := k.fungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(fc.Zrc20ContractAddress))
		if err != nil {
			return 0, errors.Wrap(fungibletypes.ErrContractCall, err.Error())
		}
		return gasLimit.Uint64(), nil
	} else if cctx.InboundTxParams.CoinType == common.CoinType_ERC20 {
		// get the gas limit of the associated asset
		fc, found := k.fungibleKeeper.GetForeignCoinFromAsset(ctx, cctx.InboundTxParams.Asset, cctx.InboundTxParams.SenderChainId)
		if !found {
			return 0, types.ErrForeignCoinNotFound
		}
		gasLimit, err := k.fungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(fc.Zrc20ContractAddress))
		if err != nil {
			return 0, errors.Wrap(fungibletypes.ErrContractCall, err.Error())
		}
		return gasLimit.Uint64(), nil
	}

	return 0, nil
}

func IsPending(cctx types.CrossChainTx) bool {
	// pending inbound is not considered a "pending" state because it has not reached consensus yet
	if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound || cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
		return true
	}
	return false
}
