package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k msgServer) CreateTSSVoter(goCtx context.Context, msg *types.MsgCreateTSSVoter) (*types.MsgCreateTSSVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !isBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	index := msg.Digest()
	// Check if the value already exists
	tssVoter, isFound := k.GetTSSVoter(ctx, index)

	if isDuplicateSigner(msg.Creator, tssVoter.Signers) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s double signing!!", msg.Creator))
	}

	if isFound {
		tssVoter.Signers = append(tssVoter.Signers, msg.Creator)
	} else {
		tssVoter = types.TSSVoter{
			Creator:         msg.Creator,
			Index:           index,
			Chain:           msg.Chain,
			Address:         msg.Address,
			Pubkey:          msg.Pubkey,
			Signers:         []string{msg.Creator},
			FinalizedHeight: 0,
		}
	}

	k.SetTSSVoter(ctx, tssVoter)

	// this needs full consensus on all validators.
	if len(tssVoter.Signers) == len(validators) {
		tss, _ := k.GetTSS(ctx, tssVoter.Chain)
		tss = types.TSS{
			Creator:             "",
			Index:               tssVoter.Chain,
			Chain:               tssVoter.Chain,
			Address:             tssVoter.Address,
			Pubkey:              tssVoter.Pubkey,
			Signer:              tssVoter.Signers,
			FinalizedZetaHeight: uint64(ctx.BlockHeader().Height),
		}
		k.SetTSS(ctx, tss)
	}

	return &types.MsgCreateTSSVoterResponse{}, nil
}
