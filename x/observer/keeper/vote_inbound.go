package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/observer/types"
)

// VoteOnInboundBallot casts a vote on an inbound transaction observed on a connected chain. If this
// is the first vote, a new ballot is created. When a threshold of votes is
// reached, the ballot is finalized.
func (k Keeper) VoteOnInboundBallot(
	ctx sdk.Context,
	senderChainID int64,
	receiverChainID int64,
	coinType coin.CoinType,
	voter string,
	ballotIndex string,
	inboundHash string,
) (isFinalized bool, isNew bool, err error) {
	if !k.IsInboundEnabled(ctx) {
		return false, false, types.ErrInboundDisabled
	}

	// makes sure we are getting only supported chains
	// if a chain support has been turned on using gov proposal
	// this function returns nil
	senderChain, found := k.GetSupportedChainFromChainID(ctx, senderChainID)
	if !found {
		return false, false, sdkerrors.Wrapf(types.ErrSupportedChains,
			"ChainID %d, Observation %s",
			senderChainID,
			types.ObservationType_InboundTx.String())
	}

	// checks the voter is authorized to vote
	err = k.CheckObserverCanVote(ctx, voter)
	if err != nil {
		return false, false, err
	}

	// makes sure we are getting only supported chains
	receiverChain, found := k.GetSupportedChainFromChainID(ctx, receiverChainID)
	if !found {
		return false, false, sdkerrors.Wrapf(types.ErrSupportedChains,
			"ChainID %d, Observation %s",
			receiverChainID,
			types.ObservationType_InboundTx.String())
	}

	// check if we want to send ZETA to external chain, but there is no ZETA token.
	if receiverChain.IsExternalChain() {
		coreParams, found := k.GetChainParamsByChainID(ctx, receiverChain.ChainId)
		if !found {
			return false, false, types.ErrChainParamsNotFound
		}
		if coreParams.ZetaTokenContractAddress == "" && coinType == coin.CoinType_Zeta {
			return false, false, types.ErrInvalidZetaCoinTypes
		}
	}

	ballot, isFinalized, isNew, err := k.VoteOnBallot(
		ctx,
		senderChain,
		ballotIndex,
		types.ObservationType_InboundTx,
		voter,
		types.VoteType_SuccessObservation,
	)
	if err != nil {
		return false, false, sdkerrors.Wrap(err, msgVoteOnBallot)
	}

	if isNew {
		EmitEventBallotCreated(ctx, ballot, inboundHash, senderChain.String())
	}

	return isFinalized, isNew, nil
}
