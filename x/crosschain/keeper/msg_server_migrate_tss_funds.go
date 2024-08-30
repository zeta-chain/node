package keeper

import (
	"context"
	"sort"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// MigrateTssFunds migrates the funds from the current TSS to the new TSS
func (k msgServer) MigrateTssFunds(
	goCtx context.Context,
	msg *types.MsgMigrateTssFunds,
) (*types.MsgMigrateTssFundsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if authorized
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
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

	err = k.initiateMigrateTSSFundsCCTX(ctx, msg.Creator, msg.ChainId, msg.Amount, tss, tssHistory)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrCannotMigrateTssFunds, err.Error())
	}

	return &types.MsgMigrateTssFundsResponse{}, nil
}

// initiateMigrateTSSFundsCCTX sets the CCTX for migrating the funds to initiate the migration outbound
func (k Keeper) initiateMigrateTSSFundsCCTX(
	ctx sdk.Context,
	creator string,
	chainID int64,
	amount sdkmath.Uint,
	currentTss observertypes.TSS,
	tssList []observertypes.TSS,
) error {
	// Always migrate to the latest TSS if multiple TSS addresses have been generated
	newTss := tssList[len(tssList)-1]
	medianGasPrice, priorityFee, isFound := k.GetMedianGasValues(ctx, chainID)
	if !isFound {
		return types.ErrUnableToGetGasPrice
	}

	// initialize the cmd CCTX
	cctx, err := types.MigrateFundCmdCCTX(
		ctx.BlockHeight(),
		creator,
		tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash()).String(),
		chainID,
		amount,
		medianGasPrice,
		priorityFee,
		currentTss.TssPubkey,
		newTss.TssPubkey,
		k.GetAuthorityKeeper().GetAdditionalChainList(ctx),
	)
	if err != nil {
		return err
	}

	// Set the CCTX and the nonce for the outbound migration
	err = k.SetObserverOutboundInfo(ctx, chainID, &cctx)
	if err != nil {
		return errorsmod.Wrap(types.ErrUnableToSetOutboundInfo, err.Error())
	}

	// The migrate funds can be run again to update the migration cctx index if the migration fails
	// This should be used after carefully calculating the amount again
	existingMigrationInfo, found := k.zetaObserverKeeper.GetFundMigrator(ctx, chainID)
	if found {
		olderMigrationCctx, found := k.GetCrossChainTx(ctx, existingMigrationInfo.MigrationCctxIndex)
		if !found {
			return errorsmod.Wrapf(
				types.ErrCannotFindCctx,
				"cannot find existing migration cctx but migration info is present for chainID %d , migrator info : %s",
				chainID,
				existingMigrationInfo.String(),
			)
		}
		if olderMigrationCctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
			return errorsmod.Wrapf(
				types.ErrUnsupportedStatus,
				"cannot migrate funds while there are pending migrations , migrator info :  %s",
				existingMigrationInfo.String(),
			)
		}
	}

	k.SetCctxAndNonceToCctxAndInboundHashToCctx(ctx, cctx, currentTss.TssPubkey)
	k.zetaObserverKeeper.SetFundMigrator(ctx, observertypes.TssFundMigratorInfo{
		ChainId:            chainID,
		MigrationCctxIndex: cctx.Index,
	})
	EmitEventInboundFinalized(ctx, &cctx)

	return nil
}
