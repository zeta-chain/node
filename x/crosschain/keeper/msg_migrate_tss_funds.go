package keeper

import (
	"context"
	"fmt"
	"sort"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) MigrateTssFunds(goCtx context.Context, msg *types.MsgMigrateTssFunds) (*types.MsgMigrateTssFundsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observerTypes.Policy_Type_group2) {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Update can only be executed by the correct policy account")
	}
	if k.zetaObserverKeeper.IsInboundEnabled(ctx) {
		return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "cannot migrate funds while inbound is enabled")
	}
	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "cannot find current TSS")
	}
	pendingNonces, found := k.GetPendingNonces(ctx, tss.TssPubkey, msg.ChainId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "cannot find pending nonces for chain")
	}
	if pendingNonces.NonceLow != pendingNonces.NonceHigh {
		return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "cannot migrate funds when there are pending nonces")
	}
	err := k.MigrateTSSFundsForChain(ctx, msg.ChainId, msg.Amount, tss)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, err.Error())
	}
	return &types.MsgMigrateTssFundsResponse{}, nil
}

func (k Keeper) MigrateTSSFundsForChain(ctx sdk.Context, chainID int64, amount sdkmath.Uint, currentTss types.TSS) error {
	tssList := k.GetAllTSS(ctx)
	if len(tssList) < 2 {
		return errorsmod.Wrap(types.ErrCannotMigrateTss, "only one TSS found")
	}
	// Sort tssList by FinalizedZetaHeight
	sort.SliceStable(tssList, func(i, j int) bool {
		return tssList[i].FinalizedZetaHeight < tssList[j].FinalizedZetaHeight
	})
	// Always migrate to the latest TSS if multiple TSS addresses have been generated
	newTss := tssList[len(tssList)-1]
	ethAddressOld, err := getTssAddrEVM(currentTss.TssPubkey)
	if err != nil {
		return err
	}
	btcAddressOld, err := getTssAddrBTC(currentTss.TssPubkey, common.BitcoinNetParamsFromChainID(chainID))
	if err != nil {
		return err
	}
	ethAddressNew, err := getTssAddrEVM(newTss.TssPubkey)
	if err != nil {
		return err
	}
	btcAddressNew, err := getTssAddrBTC(newTss.TssPubkey, common.BitcoinNetParamsFromChainID(chainID))
	if err != nil {
		return err
	}

	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chainID)
	if !isFound {
		return types.ErrUnableToGetGasPrice
	}
	indexString := fmt.Sprintf("%s-%s-%d-%s-%d", currentTss.TssPubkey, newTss.TssPubkey, chainID, amount.String(), ctx.BlockHeight())

	hash := crypto.Keccak256Hash([]byte(indexString))
	index := hash.Hex()

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
