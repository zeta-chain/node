package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// VoteOnInboundBallot casts a vote on an inbound transaction observed on a connected chain. If this
// is the first vote, a new ballot is created. When a threshold of votes is
// reached, the ballot is finalized.
func (k Keeper) VoteOnInboundBallot(
	ctx sdk.Context,
	senderChainID int64,
	receiverChainID int64,
	coinType common.CoinType,
	voter string,
	ballotIndex string,
	inTxHash string,
) (bool, error) {
	if !k.IsInboundEnabled(ctx) {
		return false, types.ErrNotEnoughPermissions
	}

	// makes sure we are getting only supported chains
	// if a chain support has been turned on using gov proposal
	// this function returns nil
	senderChain := k.GetParams(ctx).GetChainFromChainID(senderChainID)
	if senderChain == nil {
		return false, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf(
			"ChainID %d, Observation %s",
			senderChainID,
			observertypes.ObservationType_InBoundTx.String()),
		)
	}

	// checks the voter is authorized to vote on the observation chain
	if ok := k.IsAuthorized(ctx, voter, senderChain); !ok {
		return false, observertypes.ErrNotAuthorizedPolicy
	}

	// makes sure we are getting only supported chains
	receiverChain := k.GetParams(ctx).GetChainFromChainID(receiverChainID)
	if receiverChain == nil {
		return false, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf(
			"ChainID %d, Observation %s",
			receiverChain.ChainId,
			observertypes.ObservationType_InBoundTx.String()),
		)
	}

	// check if we want to send ZETA to external chain, but there is no ZETA token.
	if receiverChain.IsExternalChain() {
		coreParams, found := k.GetCoreParamsByChainID(ctx, receiverChain.ChainId)
		if !found {
			return false, types.ErrNotFoundCoreParams
		}
		if coreParams.ZetaTokenContractAddress == "" && coinType == common.CoinType_Zeta {
			return false, types.ErrUnableToSendCoinType
		}
	}

	// checks against the supported chains list before querying for Ballot
	ballot, isNew, err := k.FindBallot(ctx, ballotIndex, senderChain, observertypes.ObservationType_InBoundTx)
	if err != nil {
		return false, err
	}
	if isNew {
		EmitEventBallotCreated(ctx, ballot, inTxHash, senderChain.String())
	}

	// adds a vote and sets the ballot
	ballot, err = k.AddVoteToBallot(ctx, ballot, voter, observertypes.VoteType_SuccessObservation)
	if err != nil {
		return false, err
	}

	// checks if the ballot is finalized
	_, isFinalized := k.CheckIfFinalizingVote(ctx, ballot)
	return isFinalized, nil
}
