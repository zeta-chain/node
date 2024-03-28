package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	math2 "github.com/ethereum/go-ethereum/common/math"
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
	// No need to create a ballot if keygen does not exist
	keygen, found := k.GetKeygen(ctx)
	if !found {
		return &types.MsgVoteTSSResponse{}, types.ErrKeygenNotFound
	}
	// USE a separate transaction to update KEYGEN status to pending when trying to change the TSS address
	if keygen.Status == types.KeygenStatus_KeyGenSuccess {
		return &types.MsgVoteTSSResponse{}, types.ErrKeygenCompleted
	}
	index := msg.Digest()
	// Add votes and Set Ballot
	// GetBallot checks against the supported chains list before querying for Ballot
	// TODO : https://github.com/zeta-chain/node/issues/896
	ballot, found := k.GetBallot(ctx, index)
	if !found {
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
	}
	var err error
	if msg.Status == chains.ReceiveStatus_Success {
		ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, types.VoteType_SuccessObservation)
		if err != nil {
			return &types.MsgVoteTSSResponse{}, err
		}
	} else if msg.Status == chains.ReceiveStatus_Failed {
		ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, types.VoteType_FailureObservation)
		if err != nil {
			return &types.MsgVoteTSSResponse{}, err
		}
	}
	if !found {
		keeper.EmitEventBallotCreated(ctx, ballot, msg.TssPubkey, "Common-TSS-For-All-Chain")
	}

	ballot, isFinalized := k.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		// Return nil here to add vote to ballot and commit state
		return &types.MsgVoteTSSResponse{}, nil
	}
	// Set TSS only on success, set Keygen either way.
	// Keygen block can be updated using a policy transaction if keygen fails
	if ballot.BallotStatus != types.BallotStatus_BallotFinalized_FailureObservation {
		tss := types.TSS{
			TssPubkey:           msg.TssPubkey,
			TssParticipantList:  keygen.GetGranteePubkeys(),
			OperatorAddressList: ballot.VoterList,
			FinalizedZetaHeight: ctx.BlockHeight(),
			KeyGenZetaHeight:    msg.KeyGenZetaHeight,
		}
		// Set TSS history only, current TSS is updated via admin transaction
		// In Case this is the first TSS address update both current and history
		tssList := k.GetAllTSS(ctx)
		if len(tssList) == 0 {
			k.SetTssAndUpdateNonce(ctx, tss)
		}
		k.SetTSSHistory(ctx, tss)
		keygen.Status = types.KeygenStatus_KeyGenSuccess
		keygen.BlockNumber = ctx.BlockHeight()

	} else if ballot.BallotStatus == types.BallotStatus_BallotFinalized_FailureObservation {
		keygen.Status = types.KeygenStatus_KeyGenFailed
		keygen.BlockNumber = math2.MaxInt64
	}
	k.SetKeygen(ctx, keygen)
	return &types.MsgVoteTSSResponse{}, nil
}
