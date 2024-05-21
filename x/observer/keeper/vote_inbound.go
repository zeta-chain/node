package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/observer/types"
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
) (bool, bool, error) {
	if !k.IsInboundEnabled(ctx) {
		return false, false, types.ErrInboundDisabled
	}

	// makes sure we are getting only supported chains
	// if a chain support has been turned on using gov proposal
	// this function returns nil
	senderChain := k.GetSupportedChainFromChainID(ctx, senderChainID)
	if senderChain == nil {
		return false, false, sdkerrors.Wrap(types.ErrSupportedChains, fmt.Sprintf(
			"ChainID %d, Observation %s",
			senderChainID,
			types.ObservationType_InBoundTx.String()),
		)
	}

	// checks the voter is authorized to vote on the observation chain
	if ok := k.IsNonTombstonedObserver(ctx, voter); !ok {
		return false, false, types.ErrNotObserver
	}

	// makes sure we are getting only supported chains
	receiverChain := k.GetSupportedChainFromChainID(ctx, receiverChainID)
	if receiverChain == nil {
		return false, false, sdkerrors.Wrap(types.ErrSupportedChains, fmt.Sprintf(
			"ChainID %d, Observation %s",
			receiverChainID,
			types.ObservationType_InBoundTx.String()),
		)
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

	// checks against the supported chains list before querying for Ballot
	ballot, isNew, err := k.FindBallot(ctx, ballotIndex, senderChain, types.ObservationType_InBoundTx)
	if err != nil {
		return false, false, err
	}
	if isNew {
		EmitEventBallotCreated(ctx, ballot, inboundHash, senderChain.String())
	}

	// adds a vote and sets the ballot
	ballot, err = k.AddVoteToBallot(ctx, ballot, voter, types.VoteType_SuccessObservation)
	if err != nil {
		return false, isNew, err
	}

	// checks if the ballot is finalized
	_, isFinalized := k.CheckIfFinalizingVote(ctx, ballot)
	return isFinalized, isNew, nil
}
