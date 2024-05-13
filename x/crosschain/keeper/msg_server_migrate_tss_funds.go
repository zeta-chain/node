package keeper

import (
	"context"
	"fmt"
	"sort"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/constant"
	zetacrypto "github.com/zeta-chain/zetacore/pkg/crypto"
	"github.com/zeta-chain/zetacore/pkg/gas"

	"github.com/zeta-chain/zetacore/pkg/coin"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// MigrateTssFunds migrates the funds from the current TSS to the new TSS
func (k msgServer) MigrateTssFunds(goCtx context.Context, msg *types.MsgMigrateTssFunds) (*types.MsgMigrateTssFundsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if authorized
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	if k.zetaObserverKeeper.IsInboundEnabled(ctx) {
		return nil, errorsmod.Wrap(types.ErrCannotMigrateTssFunds, "cannot migrate funds while inbound is enabled")
	}

	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return nil, errorsmod.Wrap(types.ErrCannotMigrateTssFunds, "cannot find current TSS")
	}

	tssHistory := k.zetaObserverKeeper.GetAllTSS(ctx)
	if len(tssHistory) == 0 {
		return nil, errorsmod.Wrap(types.ErrCannotMigrateTssFunds, "empty TSS history")
	}

	sort.SliceStable(tssHistory, func(i, j int) bool {
		return tssHistory[i].FinalizedZetaHeight < tssHistory[j].FinalizedZetaHeight
	})

	if tss.TssPubkey == tssHistory[len(tssHistory)-1].TssPubkey {
		return nil, errorsmod.Wrap(types.ErrCannotMigrateTssFunds, "no new tss address has been generated")
	}

	// This check is to deal with an edge case where the current TSS is not part of the TSS history list at all
	if tss.FinalizedZetaHeight >= tssHistory[len(tssHistory)-1].FinalizedZetaHeight {
		return nil, errorsmod.Wrap(types.ErrCannotMigrateTssFunds, "current tss is the latest")
	}

	pendingNonces, found := k.GetObserverKeeper().GetPendingNonces(ctx, tss.TssPubkey, msg.ChainId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrCannotMigrateTssFunds, "cannot find pending nonces for chain")
	}

	if pendingNonces.NonceLow != pendingNonces.NonceHigh {
		return nil, errorsmod.Wrap(types.ErrCannotMigrateTssFunds, "cannot migrate funds when there are pending nonces")
	}

	err = k.MigrateTSSFundsForChain(ctx, msg.ChainId, msg.Amount, tss, tssHistory)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrCannotMigrateTssFunds, err.Error())
	}

	return &types.MsgMigrateTssFundsResponse{}, nil
}

func (k Keeper) MigrateTSSFundsForChain(ctx sdk.Context, chainID int64, amount sdkmath.Uint, currentTss observertypes.TSS, tssList []observertypes.TSS) error {
	// Always migrate to the latest TSS if multiple TSS addresses have been generated
	newTss := tssList[len(tssList)-1]
	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chainID)
	if !isFound {
		return types.ErrUnableToGetGasPrice
	}
	indexString := GetIndexStringForTssMigration(currentTss.TssPubkey, newTss.TssPubkey, chainID, amount, ctx.BlockHeight())

	hash := crypto.Keccak256Hash([]byte(indexString))
	index := hash.Hex()

	// TODO : Use the `NewCCTX` method to create the cctx
	// https://github.com/zeta-chain/node/issues/1909
	cctx := types.CrossChainTx{
		Creator:        "",
		Index:          index,
		ZetaFees:       sdkmath.Uint{},
		RelayedMessage: fmt.Sprintf("%s:%s", constant.CmdMigrateTssFunds, "Funds Migrator Admin Cmd"),
		CctxStatus: &types.Status{
			Status:              types.CctxStatus_PendingOutbound,
			StatusMessage:       "",
			LastUpdateTimestamp: 0,
		},
		InboundTxParams: &types.InboundTxParams{
			Sender:                          "",
			SenderChainId:                   chainID,
			TxOrigin:                        "",
			CoinType:                        coin.CoinType_Cmd,
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
			CoinType:                         coin.CoinType_Cmd,
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
	// Set the sender and receiver addresses for EVM chain
	if chains.IsEVMChain(chainID) {
		ethAddressOld, err := zetacrypto.GetTssAddrEVM(currentTss.TssPubkey)
		if err != nil {
			return err
		}
		ethAddressNew, err := zetacrypto.GetTssAddrEVM(newTss.TssPubkey)
		if err != nil {
			return err
		}
		cctx.InboundTxParams.Sender = ethAddressOld.String()
		cctx.GetCurrentOutTxParam().Receiver = ethAddressNew.String()
		// Tss migration is a send transaction, so the gas limit is set to 21000
		cctx.GetCurrentOutTxParam().OutboundTxGasLimit = gas.EVMSend
		// Multiple current gas price with standard multiplier to add some buffer
		multipliedGasPrice, err := gas.MultiplyGasPrice(medianGasPrice, types.TssMigrationGasMultiplierEVM)
		if err != nil {
			return err
		}
		cctx.GetCurrentOutTxParam().OutboundTxGasPrice = multipliedGasPrice.String()
		evmFee := sdkmath.NewUint(cctx.GetCurrentOutTxParam().OutboundTxGasLimit).Mul(multipliedGasPrice)
		if evmFee.GT(amount) {
			return errorsmod.Wrap(types.ErrInsufficientFundsTssMigration, fmt.Sprintf("insufficient funds to pay for gas fee, amount: %s, gas fee: %s, chainid: %d", amount.String(), evmFee.String(), chainID))
		}
		cctx.GetCurrentOutTxParam().Amount = amount.Sub(evmFee)
	}
	// Set the sender and receiver addresses for Bitcoin chain
	if chains.IsBitcoinChain(chainID) {
		bitcoinNetParams, err := chains.BitcoinNetParamsFromChainID(chainID)
		if err != nil {
			return err
		}
		btcAddressOld, err := zetacrypto.GetTssAddrBTC(currentTss.TssPubkey, bitcoinNetParams)
		if err != nil {
			return err
		}
		btcAddressNew, err := zetacrypto.GetTssAddrBTC(newTss.TssPubkey, bitcoinNetParams)
		if err != nil {
			return err
		}
		cctx.InboundTxParams.Sender = btcAddressOld
		cctx.GetCurrentOutTxParam().Receiver = btcAddressNew
	}

	if cctx.GetCurrentOutTxParam().Receiver == "" {
		return errorsmod.Wrap(types.ErrReceiverIsEmpty, fmt.Sprintf("chain %d is not supported", chainID))
	}

	err := k.UpdateNonce(ctx, chainID, &cctx)
	if err != nil {
		return err
	}
	// The migrate funds can be run again to update the migration cctx index if the migration fails
	// This should be used after carefully calculating the amount again
	existingMigrationInfo, found := k.zetaObserverKeeper.GetFundMigrator(ctx, chainID)
	if found {
		olderMigrationCctx, found := k.GetCrossChainTx(ctx, existingMigrationInfo.MigrationCctxIndex)
		if !found {
			return errorsmod.Wrapf(types.ErrCannotFindCctx, "cannot find existing migration cctx but migration info is present for chainID %d , migrator info : %s", chainID, existingMigrationInfo.String())
		}
		if olderMigrationCctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
			return errorsmod.Wrapf(types.ErrUnsupportedStatus, "cannot migrate funds while there are pending migrations , migrator info :  %s", existingMigrationInfo.String())
		}
	}

	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
	k.zetaObserverKeeper.SetFundMigrator(ctx, observertypes.TssFundMigratorInfo{
		ChainId:            chainID,
		MigrationCctxIndex: index,
	})
	EmitEventInboundFinalized(ctx, &cctx)

	return nil
}

func GetIndexStringForTssMigration(currentTssPubkey, newTssPubkey string, chainID int64, amount sdkmath.Uint, height int64) string {
	return fmt.Sprintf("%s-%s-%d-%s-%d", currentTssPubkey, newTssPubkey, chainID, amount.String(), height)
}
