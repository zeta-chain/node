package keeper

import (
	"fmt"
	"sort"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) MigrateTSSFundsForChain(ctx sdk.Context, chainID int64, amount sdkmath.Uint) error {
	currentTss, found := k.GetTSS(ctx)
	if !found {
		return types.ErrCannotFindTSSKeys
	}
	tssList := k.GetAllTSS(ctx)
	if len(tssList) < 2 {
		return errorsmod.Wrap(types.ErrCannotMigrateTss, "only one TSS found")
	}
	// Sort tssList by FinalizedZetaHeight
	sort.SliceStable(tssList, func(i, j int) bool {
		return tssList[i].FinalizedZetaHeight < tssList[j].FinalizedZetaHeight
	})
	newTss := tssList[len(tssList)-1]
	ethAddressOld, err := getTssAddrEVM(currentTss.TssPubkey)
	if err != nil {
		return err
	}
	btcAddressOld, err := getTssAddrBTC(currentTss.TssPubkey)
	if err != nil {
		return err
	}
	ethAddressNew, err := getTssAddrEVM(newTss.TssPubkey)
	if err != nil {
		return err
	}
	btcAddressNew, err := getTssAddrBTC(newTss.TssPubkey)
	if err != nil {
		return err
	}

	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chainID)
	if !isFound {
		return types.ErrUnableToGetGasPrice
	}
	index := fmt.Sprintf("%s-%s-%d-%s", currentTss.TssPubkey, newTss.TssPubkey, chainID, amount.String())

	cctx := types.CrossChainTx{
		Creator:        "",
		Index:          index,
		ZetaFees:       sdkmath.Uint{},
		RelayedMessage: fmt.Sprintf("%s:%s", common.CmdMigrateTssFunds, "Funds Migrator Admin Cmd"),
		CctxStatus: &types.Status{
			Status:              types.CctxStatus_PendingOutbound,
			StatusMessage:       "",
			LastUpdateTimestamp: 0,
		},
		InboundTxParams: &types.InboundTxParams{
			Sender:                          "",
			SenderChainId:                   chainID,
			TxOrigin:                        "",
			CoinType:                        common.CoinType_Cmd,
			Asset:                           "",
			Amount:                          amount,
			InboundTxObservedHash:           tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash()).String(),
			InboundTxObservedExternalHeight: 0,
			InboundTxBallotIndex:            "",
			InboundTxFinalizedZetaHeight:    0,
		},
		OutboundTxParams: []*types.OutboundTxParams{{
			Receiver:                         "",
			ReceiverChainId:                  chainID,
			CoinType:                         common.CoinType_Cmd,
			Amount:                           amount,
			OutboundTxTssNonce:               0,
			OutboundTxGasLimit:               1_000_000,
			OutboundTxGasPrice:               medianGasPrice.MulUint64(2).String(),
			OutboundTxHash:                   "",
			OutboundTxBallotIndex:            "",
			OutboundTxObservedExternalHeight: 0,
			OutboundTxGasUsed:                0,
			OutboundTxEffectiveGasPrice:      sdkmath.Int{},
			OutboundTxEffectiveGasLimit:      0,
			TssPubkey:                        currentTss.TssPubkey,
		}}}

	if common.IsEVMChain(chainID) {
		cctx.InboundTxParams.Sender = ethAddressOld.String()
		cctx.GetCurrentOutTxParam().Receiver = ethAddressNew.String()
	}
	if common.IsBitcoinChain(chainID) {
		cctx.InboundTxParams.Sender = btcAddressOld
		cctx.GetCurrentOutTxParam().Receiver = btcAddressNew
	}
	if cctx.GetCurrentOutTxParam().Receiver == "" {
		return errorsmod.Wrap(types.ErrCannotMigrateTss, fmt.Sprintf("chain %d is not supported", chainID))
	}
	err = k.UpdateNonce(ctx, chainID, &cctx)
	if err != nil {
		return err
	}
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
	EmitEventInboundFinalized(ctx, &cctx)
	return nil
}
