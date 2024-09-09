package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// UpdateERC20CustodyPauseStatus creates a admin cmd cctx to update the pause status of the ERC20 custody contract
func (k msgServer) UpdateERC20CustodyPauseStatus(
	goCtx context.Context,
	msg *types.MsgUpdateERC20CustodyPauseStatus,
) (*types.MsgUpdateERC20CustodyPauseStatusResponse, error) {
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
	medianGasPrice = medianGasPrice.MulUint64(types.ERC20CustodyPausingGasMultiplierEVM)
	priorityFee = priorityFee.MulUint64(types.ERC20CustodyPausingGasMultiplierEVM)

	// should not happen
	if priorityFee.GT(medianGasPrice) {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidGasAmount,
			"priorityFee %s is greater than median gasPrice %s",
			priorityFee.String(),
			medianGasPrice.String(),
		)
	}

	// create the CCTX that allows to sign the ERC20 custody pause status update
	cctx := types.UpdateERC20CustodyPauseStatusCmdCCTX(
		msg.Creator,
		params.Erc20CustodyContractAddress,
		msg.ChainId,
		msg.Pause,
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
		&types.EventERC20CustodyPausing{
			ChainId:   msg.ChainId,
			Pause:     msg.Pause,
			CctxIndex: cctx.Index,
		},
	)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to emit event")
	}

	return &types.MsgUpdateERC20CustodyPauseStatusResponse{
		CctxIndex: cctx.Index,
	}, nil
}
