package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
  
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// MESSAGES

// The CreateTSSVoter function creates a threshold signature scheme (TSS) voter
// and adds it to the TSS voter store. It first checks if the message creator is
// a bonded validator. It then calculates the sessionID based on the current
// block height and creates an index using the message digest and the sessionID.
// It retrieves a TSS voter based on the index and checks if the creator has
// already signed. If the creator has not signed, the method adds the creator to
// the Signers list in the TSS voter. If the TSS voter is not found, the method
// creates a new TSS voter with the provided information and initializes the
// Signers list with the creator. The method then sets the TSS voter in the
// store using the SetTSSVoter function. If the Signers list in the TSS voter is
// equal to the number of validators, the method creates a new TSS using the TSS
// voter information and sets it in the TSS store using the SetTSS function. The
// function returns an empty MsgCreateTSSVoterResponse and no error.
func (k msgServer) CreateTSSVoter(goCtx context.Context, msg *types.MsgCreateTSSVoter) (*types.MsgCreateTSSVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAuthorizedNodeAccount(ctx, msg.Creator) {
		return nil, errorsmod.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s does not have a node account set", msg.Creator))
	}
	// No need to create a ballot if keygen does not exist
	keygen, found := k.GetKeygen(ctx)
	if !found {
		return &types.MsgCreateTSSVoterResponse{}, types.ErrKeygenNotFound
	}
	// USE a separate transaction to update KEYGEN status to pending when trying to change the TSS address
	if keygen.Status == types.KeygenStatus_KeyGenSuccess {
		return &types.MsgCreateTSSVoterResponse{}, types.ErrKeygenNotFound
	}
	index := msg.Digest()
	// Add votes and Set Ballot
	// GetBallot checks against the supported chains list before querying for Ballot
	ballot, found := k.zetaObserverKeeper.GetBallot(ctx, index)
	if !found {
		var voterList []string

		for _, nodeAccount := range k.GetAllNodeAccount(ctx) {
			voterList = append(voterList, nodeAccount.Operator)
		}
		ballot = zetaObserverTypes.Ballot{
			Index:            "",
			BallotIdentifier: index,
			VoterList:        voterList,
			Votes:            zetaObserverTypes.CreateVotes(len(msg.Creator)),
			ObservationType:  zetaObserverTypes.ObservationType_TSSKeyGen,
			BallotThreshold:  sdk.MustNewDecFromStr("1.00"),
			BallotStatus:     zetaObserverTypes.BallotStatus_BallotInProgress,
		}
	}
	err := error(nil)
	if msg.Status == common.ReceiveStatus_Success {
		ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, zetaObserverTypes.VoteType_SuccessObservation)
		if err != nil {
			return &types.MsgCreateTSSVoterResponse{}, err
		}
	} else if msg.Status == common.ReceiveStatus_Failed {
		ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, zetaObserverTypes.VoteType_FailureObservation)
		if err != nil {
			return &types.MsgCreateTSSVoterResponse{}, err
		}
	}

	ballot, isFinalized := k.CheckIfBallotIsFinalized(ctx, ballot)
	if !isFinalized {
		return &types.MsgCreateTSSVoterResponse{}, nil
	}
	// Set TSS only on success , set Keygen either way.
	// Keygen block cna be updated using a policy transaction if keygen fails
	if ballot.BallotStatus != zetaObserverTypes.BallotStatus_BallotFinalized_FailureObservation {
		k.SetTSS(ctx, types.TSS{
			TssPubkey:           msg.TssPubkey,
			TssParticipantList:  keygen.GetGranteePubkeys(),
			OperatorAddressList: ballot.VoterList,
			FinalizedZetaHeight: ctx.BlockHeight(),
			KeyGenZetaHeight:    msg.KeyGenZetaHeight,
		})
		keygen.Status = types.KeygenStatus_KeyGenSuccess
		// initialize the nonces and pending nonces of all enabled chain
		supportedChains := k.zetaObserverKeeper.GetParams(ctx).GetSupportedChains()
		for _, chain := range supportedChains {
			chainNonce := types.ChainNonces{Index: chain.ChainName.String(), ChainId: chain.ChainId, Nonce: 0, FinalizedHeight: uint64(ctx.BlockHeight())}
			k.SetChainNonces(ctx, chainNonce)

			p := types.PendingNonces{
				NonceLow:  0,
				NonceHigh: 0,
				ChainId:   chain.ChainId,
				Tss:       msg.TssPubkey,
			}
			k.SetPendingNonces(ctx, p)
		}
	} else if ballot.BallotStatus == zetaObserverTypes.BallotStatus_BallotFinalized_FailureObservation {
		keygen.Status = types.KeygenStatus_KeyGenFailed
	}
	k.SetKeygen(ctx, keygen)
	// Remove ballot
	return &types.MsgCreateTSSVoterResponse{}, nil
}
