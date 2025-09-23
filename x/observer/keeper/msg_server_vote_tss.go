package keeper

import (
	"context"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/observer/types"
)

const voteTSSid = "Vote TSS"

// VoteTSS votes on creating a TSS key and recording the information about it (public
// key, participant and operator addresses, finalized and keygen heights).
//
// If the vote passes, the information about the TSS key is recorded on chain
// and the status of the keygen is set to "success".
//
// Fails if the keygen does not exist, the keygen has been already
// completed, or the keygen has failed.
//
// Only node accounts are authorized to broadcast this message.
func (k msgServer) VoteTSS(goCtx context.Context, msg *types.MsgVoteTSS) (*types.MsgVoteTSSResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks whether a signer is authorized to sign, by checking if the signer has a node account.
	_, found := k.GetNodeAccount(ctx, msg.Creator)
	if !found {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrorInvalidSigner,
			"%s, signer %s does not have a node account set", voteTSSid, msg.Creator)
	}

	// No need to create a ballot if keygen does not exist.
	keygen, found := k.GetKeygen(ctx)
	if !found {
		return &types.MsgVoteTSSResponse{}, errorsmod.Wrap(types.ErrKeygenNotFound, voteTSSid)
	}

	// GetBallot checks against the supported chains list before querying for Ballot.
	ballotCreated := false
	index := msg.Digest()
	ballot, found := k.GetBallot(ctx, index)
	if !found {
		// If ballot does not exist, create a new ballot.
		var voterList []string

		for _, nodeAccount := range k.GetAllNodeAccount(ctx) {
			voterList = append(voterList, nodeAccount.Operator)
		}

		ballot = types.Ballot{
			BallotIdentifier:     index,
			VoterList:            voterList,
			Votes:                types.CreateVotes(len(voterList)),
			ObservationType:      types.ObservationType_TSSKeyGen,
			BallotThreshold:      sdkmath.LegacyMustNewDecFromStr("1.00"),
			BallotStatus:         types.BallotStatus_BallotInProgress,
			BallotCreationHeight: ctx.BlockHeight(),
		}
		k.AddBallotToList(ctx, ballot)

		EmitEventBallotCreated(ctx, ballot, msg.TssPubkey, "Common-TSS-For-All-Chain")
		ballotCreated = true
	}

	vote := types.VoteType_SuccessObservation
	if msg.Status == chains.ReceiveStatus_failed {
		vote = types.VoteType_FailureObservation
	}

	ballot, err := k.AddVoteToBallot(ctx, ballot, msg.Creator, vote)
	if err != nil {
		return &types.MsgVoteTSSResponse{}, errorsmod.Wrap(err, voteTSSid)
	}

	ballot, isFinalized := k.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		return &types.MsgVoteTSSResponse{
			VoteFinalized: isFinalized,
			BallotCreated: ballotCreated,
			KeygenSuccess: false,
		}, nil
	}

	// The ballot is finalized, we check if this is the correct ballot for updating the TSS
	// The requirements are
	// 1. The keygen is still pending
	// 2. The keygen block number matches the ballot block number ,which makes sure this the correct ballot for the current keygen

	// Return without an error so the vote is added to the ballot
	if keygen.Status != types.KeygenStatus_PendingKeygen {
		// The response is used for testing only.Setting false for keygen success as the keygen has already been finalized and it doesnt matter what the final status is.We are just asserting that the keygen was previously finalized and is not in pending status.
		return &types.MsgVoteTSSResponse{
			VoteFinalized: isFinalized,
			BallotCreated: ballotCreated,
			KeygenSuccess: false,
		}, nil
	}

	// For cases when an observer tries to vote for an older pending ballot, associated with a keygen that was discarded, we would return at this check while still adding the vote to the ballot
	if msg.KeygenZetaHeight != keygen.BlockNumber {
		return &types.MsgVoteTSSResponse{
			VoteFinalized: isFinalized,
			BallotCreated: ballotCreated,
			KeygenSuccess: false,
		}, nil
	}

	// Set TSS only on success, set keygen either way.
	// Keygen block can be updated using a policy transaction if keygen fails.
	keygenSuccess := false
	if ballot.BallotStatus == types.BallotStatus_BallotFinalized_FailureObservation {
		keygen.Status = types.KeygenStatus_KeyGenFailed
		keygen.BlockNumber = math.MaxInt64
	} else {
		tss := types.TSS{
			TssPubkey:           msg.TssPubkey,
			TssParticipantList:  keygen.GetGranteePubkeys(),
			OperatorAddressList: ballot.VoterList,
			FinalizedZetaHeight: ctx.BlockHeight(),
			KeyGenZetaHeight:    msg.KeygenZetaHeight,
		}

		// Set TSS history only, current TSS is updated via admin transaction.
		// In the case this is the first TSS address update both current and history.
		tssList := k.GetAllTSS(ctx)
		if len(tssList) == 0 {
			k.SetTssAndUpdateNonce(ctx, tss)
		}
		k.SetTSSHistory(ctx, tss)
		keygen.Status = types.KeygenStatus_KeyGenSuccess
		keygen.BlockNumber = ctx.BlockHeight()
		keygenSuccess = true
	}

	k.SetKeygen(ctx, keygen)

	return &types.MsgVoteTSSResponse{
		VoteFinalized: isFinalized,
		BallotCreated: ballotCreated,
		KeygenSuccess: keygenSuccess,
	}, nil
}
