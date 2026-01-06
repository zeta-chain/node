package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// UpdateTssAddress updates the TSS address.
func (k msgServer) UpdateTssAddress(
	goCtx context.Context,
	msg *types.MsgUpdateTssAddress,
) (*types.MsgUpdateTssAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if authorized
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
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

	// Each connected chain that needs funds migration should have its own tss migrator
	// Solana, Sui, Ton do not need funds migration; updating TSS address is enough
	if len(k.TSSFundsMigrationChains(ctx)) != len(tssMigrators) {
		return nil, errorsmod.Wrap(
			types.ErrUnableToUpdateTss,
			"cannot update tss address incorrect number of migrations have been created and completed",
		)
	}

	// GetAllTssFundMigrators would return the migrators created for the current migration
	// if any of the migrations is still pending we should not allow the tss address to be updated
	// we can wait for all migrations to complete before updating; this includes btc and evm chains.
	for _, tssMigrator := range tssMigrators {
		migratorTx, found := k.GetCrossChainTx(ctx, tssMigrator.MigrationCctxIndex)
		if !found {
			return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "migration cross chain tx not found")
		}
		if migratorTx.CctxStatus.Status != types.CctxStatus_OutboundMined {
			return nil, errorsmod.Wrapf(
				types.ErrUnableToUpdateTss,
				"cannot update tss address while there are pending migrations , current status of migration cctx : %s ",
				migratorTx.CctxStatus.Status.String(),
			)
		}
	}

	k.GetObserverKeeper().SetTssAndUpdateNonce(ctx, tss)

	// Remove all migrators once the tss address has been updated successfully,
	// A new set of migrators will be created when the next migration is triggered
	k.zetaObserverKeeper.RemoveAllExistingMigrators(ctx)

	return &types.MsgUpdateTssAddressResponse{}, nil
}

// TSSFundsMigrationChains returns the chains that support tss migration.
// Chains that support tss migration are chains that have the following properties:
// 1. External chains
// 2. Gateway observer
// 3. VM is EVM, or Consensus is bitcoin (Solana, Sui, TON are excluded as they do not require funds migration)
func (k *Keeper) TSSFundsMigrationChains(ctx sdk.Context) []chains.Chain {
	supportedChains := k.zetaObserverKeeper.GetSupportedChains(ctx)
	return chains.CombineFilterChains([][]chains.Chain{
		chains.FilterChains(supportedChains, []chains.ChainFilter{
			chains.FilterExternalChains,
			chains.FilterByGateway(chains.CCTXGateway_observers),
			chains.FilterByVM(chains.Vm_evm),
		}...),
		chains.FilterChains(supportedChains, []chains.ChainFilter{
			chains.FilterExternalChains,
			chains.FilterByGateway(chains.CCTXGateway_observers),
			chains.FilterByConsensus(chains.Consensus_bitcoin),
		}...),
	}...)
}
