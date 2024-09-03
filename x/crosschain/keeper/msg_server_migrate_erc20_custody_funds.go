package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// MigrateERC20CustodyFunds migrates the funds from the current ERC20Custody contract to the new ERC20Custody contract
func (k msgServer) MigrateERC20CustodyFunds(
	goCtx context.Context,
	msg *types.MsgMigrateERC20CustodyFunds,
) (*types.MsgMigrateERC20CustodyFundsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if authorized
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// get the current TSS nonce allow to set a unique index for the CCTX
	chainNonce, found := k.GetObserverKeeper().GetChainNonces(ctx, msg.ChainId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrInvalidChainID, "cannot find current chain nonce")
	}
	currentNonce := chainNonce.Nonce

	// get the current TSS
	tss, found := k.GetObserverKeeper().GetTSS(ctx)
	if !found {
		return nil, errorsmod.Wrap(types.ErrCannotFindTSSKeys, "cannot find current TSS")
	}

	// get necessary parameters to create the cctx
	params, found := k.zetaObserverKeeper.GetChainParamsByChainID(ctx, msg.ChainId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrInvalidChainID, "chain params not found for chain id (%d)", msg.ChainId)
	}
	medianGasPrice, priorityFee, isFound := k.GetMedianGasValues(ctx, msg.ChainId)
	if !isFound {
		return nil, errorsmod.Wrapf(
			types.ErrUnableToGetGasPrice,
			"median gas price not found for chain id (%d)",
			msg.ChainId,
		)
	}

	// overpays gas price by 2x
	medianGasPrice = medianGasPrice.MulUint64(types.ERC20CustodyMigrationGasMultiplierEVM)
	priorityFee = priorityFee.MulUint64(types.ERC20CustodyMigrationGasMultiplierEVM)

	// should not happen
	if priorityFee.GT(medianGasPrice) {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidGasAmount,
			"priorityFee %s is greater than median gasPrice %s",
			priorityFee.String(),
			medianGasPrice.String(),
		)
	}

	// create the CCTX that allows to sign the fund migration
	cctx := types.MigrateERC20CustodyFundsCmdCCTX(
		msg.Creator,
		msg.Erc20Address,
		params.Erc20CustodyContractAddress,
		msg.NewCustodyAddress,
		msg.ChainId,
		msg.Amount,
		medianGasPrice.String(),
		priorityFee.String(),
		tss.TssPubkey,
		currentNonce,
	)

	// save the cctx
	err = k.SetObserverOutboundInfo(ctx, msg.ChainId, &cctx)
	if err != nil {
		return nil, err
	}
	k.SetCctxAndNonceToCctxAndInboundHashToCctx(ctx, cctx, tss.TssPubkey)

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventERC20CustodyFundsMigration{
			NewCustodyAddress: msg.NewCustodyAddress,
			Erc20Address:      msg.Erc20Address,
			Amount:            msg.Amount.String(),
			CctxIndex:         cctx.Index,
		},
	)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to emit event")
	}

	return &types.MsgMigrateERC20CustodyFundsResponse{
		CctxIndex: cctx.Index,
	}, nil
}
