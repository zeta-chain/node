package keeper

import (
	"bytes"
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NonceVoter(goCtx context.Context, msg *types.MsgNonceVoter) (*types.MsgNonceVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !isBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	chain := msg.Chain
	chainNonce, isFound := k.GetChainNonces(ctx, chain)
	if isDuplicateSigner(msg.Creator, chainNonce.Signers) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s double signing!!", msg.Creator))
	}
	if isFound && chainNonce.Nonce == msg.Nonce {
		chainNonce.Signers = append(chainNonce.Signers, msg.Creator)
	} else if !isFound {
		chainNonce = types.ChainNonces{
			Creator: msg.Creator,
			Index:   msg.Chain,
			Chain:   msg.Chain,
			Nonce:   msg.Nonce,
			Signers: []string{msg.Creator},
		}
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("chainNonce vote msg does not match state: %v vs %v", msg, chainNonce))
	}

	if hasSuperMajorityValidators(len(chainNonce.Signers), validators) {
		chainNonce.FinalizedHeight = uint64(ctx.BlockHeader().Height)
	}

	k.SetChainNonces(ctx, chainNonce)
	return &types.MsgNonceVoterResponse{}, nil
}

func isDuplicateSigner(creator string, signers []string) bool {
	for _, v := range signers {
		if creator == v {
			return true
		}
	}
	return false
}

func isBondedValidator(creator string, validators []stakingtypes.Validator) bool {
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {return false}
	for _, v := range validators {
		valAddr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
		if err != nil {continue}
		//TODO: How about Jailed?
		if v.IsBonded() &&  bytes.Compare(creatorAddr.Bytes(), valAddr.Bytes()) == 0 {
			return true
		}
	}
	return false
}

func hasSuperMajorityValidators(numSigners int, validators []stakingtypes.Validator) bool {
	numValidValidators := 0
	for _, v := range validators {
		if v.IsBonded() {
			numValidValidators += 1
		}
	}
	return numSigners > numValidValidators*2/3
}