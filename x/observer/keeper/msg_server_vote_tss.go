package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

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

	// checks whether a signer is authorized to sign , by checking their address against the observer mapper which contains the observer list for the chain and type
	_, found := k.GetNodeAccount(ctx, msg.Creator)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s does not have a node account set", msg.Creator))
	}
	// no need to create a ballot if keygen does not exist
	keygen, found := k.GetKeygen(ctx)
	if !found {
		return &types.MsgVoteTSSResponse{}, types.ErrKeygenNotFound
	}
	// use a separate transaction to update KEYGEN status to pending when trying to change the TSS address
	if keygen.Status == types.KeygenStatus_KeyGenSuccess {
		return &types.MsgVoteTSSResponse{}, types.ErrKeygenCompleted
	}

	// add votes and set Ballot
	// GetBallot checks against the supported chains list before querying for Ballot
	// TODO : https://github.com/zeta-chain/node/issues/896
	ballotCreated := false
	index := msg.Digest()
	ballot, found := k.GetBallot(ctx, index)
	if !found {

		// if ballot does not exist, create a new ballot
		var voterList []string

		for _, nodeAccount := range k.GetAllNodeAccount(ctx) {
			voterList = append(voterList, nodeAccount.Operator)
		}
		ballot = types.Ballot{
			Index:                "",
			BallotIdentifier:     index,
			VoterList:            voterList,
			Votes:                types.CreateVotes(len(voterList)),
			ObservationType:      types.ObservationType_TSSKeyGen,
			BallotThreshold:      sdk.MustNewDecFromStr("1.00"),
			BallotStatus:         types.BallotStatus_BallotInProgress,
			BallotCreationHeight: ctx.BlockHeight(),
		}
		k.AddBallotToList(ctx, ballot)

		EmitEventBallotCreated(ctx, ballot, msg.TssPubkey, "Common-TSS-For-All-Chain")
		ballotCreated = true
	}

	// vote the ballot
	var err error
	vote := types.VoteType_SuccessObservation
	if msg.Status == chains.ReceiveStatus_Failed {
		vote = types.VoteType_FailureObservation
	}
	ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, vote)
	if err != nil {
		return &types.MsgVoteTSSResponse{}, err
	}

	// returns here if the ballot is not finalized
	ballot, isFinalized := k.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		return &types.MsgVoteTSSResponse{
			VoteFinalized: false,
			BallotCreated: ballotCreated,
			KeygenSuccess: false,
		}, nil
	}

	// set TSS only on success, set Keygen either way.
	// keygen block can be updated using a policy transaction if keygen fails
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
		// set TSS history only, current TSS is updated via admin transaction
		// in Case this is the first TSS address update both current and history
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
		VoteFinalized: true,
		BallotCreated: ballotCreated,
		KeygenSuccess: keygenSuccess,
	}, nil
}
