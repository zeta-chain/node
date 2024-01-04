package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// Authorized: admin policy group 2.
func (k msgServer) UpdateTssAddress(goCtx context.Context, msg *types.MsgUpdateTssAddress) (*types.MsgUpdateTssAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// TODO : Add a new policy type for updating the TSS address
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observerTypes.Policy_Type_group2) {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Update can only be executed by the correct policy account")
	}
	currentTss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "cannot find current TSS")
	}
	if currentTss.TssPubkey == msg.TssPubkey {
		return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "no new tss address has been generated")
	}
	tss, ok := k.zetaObserverKeeper.CheckIfTssPubkeyHasBeenGenerated(ctx, msg.TssPubkey)
	if !ok {
		return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "tss pubkey has not been generated")
	}

	tssMigrators := k.zetaObserverKeeper.GetAllTssFundMigrators(ctx)
	// Each connected chain should have its own tss migrator
	if len(k.zetaObserverKeeper.GetParams(ctx).GetSupportedChains()) != len(tssMigrators) {
		return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "cannot update tss address not enough migrations have been created and completed")
	}
	// GetAllTssFundMigrators would return the migrators created for the current migration
	// if any of the migrations is still pending we should not allow the tss address to be updated
	// we can wait for all migrations to complete before updating; this includes btc and eth chains.
	for _, tssMigrator := range tssMigrators {
		migratorTx, found := k.GetCrossChainTx(ctx, tssMigrator.MigrationCctxIndex)
		if !found {
			return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "migration cross chain tx not found")
		}
		if migratorTx.CctxStatus.Status != types.CctxStatus_OutboundMined {
			return nil, errorsmod.Wrapf(types.ErrUnableToUpdateTss,
				"cannot update tss address while there are pending migrations , current status of migration cctx : %s ", migratorTx.CctxStatus.Status.String())
		}

	}

	k.GetObserverKeeper().SetTssAndUpdateNonce(ctx, tss)
	// Remove all migrators once the tss address has been updated successfully,
	// A new set of migrators will be created when the next migration is triggered
	k.zetaObserverKeeper.RemoveAllExistingMigrators(ctx)

	return &types.MsgUpdateTssAddressResponse{}, nil
}
