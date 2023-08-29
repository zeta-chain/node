package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	math2 "github.com/ethereum/go-ethereum/common/math"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerKeeper "github.com/zeta-chain/zetacore/x/observer/keeper"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// MESSAGES

// Vote on creating a TSS key and recording the information about it (public
// key, participant and operator addresses, finalized and keygen heights).
//
// If the vote passes, the information about the TSS key is recorded on chain
// and the status of the keygen is set to "success".
//
// Fails if the keygen does not exist, the keygen has been already
// completed, or the keygen has failed.
//
// Only node accounts are authorized to broadcast this message.
func (k msgServer) CreateTSSVoter(goCtx context.Context, msg *types.MsgCreateTSSVoter) (*types.MsgCreateTSSVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsAuthorizedNodeAccount(ctx, msg.Creator) {
		return nil, errorsmod.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s does not have a node account set", msg.Creator))
	}
	// No need to create a ballot if keygen does not exist
	keygen, found := k.zetaObserverKeeper.GetKeygen(ctx)
	if !found {
		return &types.MsgCreateTSSVoterResponse{}, observerTypes.ErrKeygenNotFound
	}
	// USE a separate transaction to update KEYGEN status to pending when trying to change the TSS address
	if keygen.Status == observerTypes.KeygenStatus_KeyGenSuccess {
		return &types.MsgCreateTSSVoterResponse{}, observerTypes.ErrKeygenCompleted
	}
	index := msg.Digest()
	// Add votes and Set Ballot
	// GetBallot checks against the supported chains list before querying for Ballot
	// TODO : https://github.com/zeta-chain/node/issues/896
	ballot, found := k.zetaObserverKeeper.GetBallot(ctx, index)
	if !found {
		var voterList []string

		for _, nodeAccount := range k.zetaObserverKeeper.GetAllNodeAccount(ctx) {
			voterList = append(voterList, nodeAccount.Operator)
		}
		ballot = observerTypes.Ballot{
			Index:                "",
			BallotIdentifier:     index,
			VoterList:            voterList,
			Votes:                observerTypes.CreateVotes(len(voterList)),
			ObservationType:      observerTypes.ObservationType_TSSKeyGen,
			BallotThreshold:      sdk.MustNewDecFromStr("1.00"),
			BallotStatus:         observerTypes.BallotStatus_BallotInProgress,
			BallotCreationHeight: ctx.BlockHeight(),
		}
		k.zetaObserverKeeper.AddBallotToList(ctx, ballot)
	}
	err := error(nil)
	if msg.Status == common.ReceiveStatus_Success {
		ballot, err = k.zetaObserverKeeper.AddVoteToBallot(ctx, ballot, msg.Creator, observerTypes.VoteType_SuccessObservation)
		if err != nil {
			return &types.MsgCreateTSSVoterResponse{}, err
		}
	} else if msg.Status == common.ReceiveStatus_Failed {
		ballot, err = k.zetaObserverKeeper.AddVoteToBallot(ctx, ballot, msg.Creator, observerTypes.VoteType_FailureObservation)
		if err != nil {
			return &types.MsgCreateTSSVoterResponse{}, err
		}
	}
	if !found {
		observerKeeper.EmitEventBallotCreated(ctx, ballot, msg.TssPubkey, "Common-TSS-For-All-Chain")
	}

	ballot, isFinalized := k.zetaObserverKeeper.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		// Return nil here to add vote to ballot and commit state
		return &types.MsgCreateTSSVoterResponse{}, nil
	}
	// Set TSS only on success, set Keygen either way.
	// Keygen block can be updated using a policy transaction if keygen fails
	if ballot.BallotStatus != observerTypes.BallotStatus_BallotFinalized_FailureObservation {
		tss := types.TSS{
			TssPubkey:           msg.TssPubkey,
			TssParticipantList:  keygen.GetGranteePubkeys(),
			OperatorAddressList: ballot.VoterList,
			FinalizedZetaHeight: ctx.BlockHeight(),
			KeyGenZetaHeight:    msg.KeyGenZetaHeight,
		}
		// Set TSS history only, current TSS is updated via admin transaction
		k.SetTSSHistory(ctx, tss)
		keygen.Status = observerTypes.KeygenStatus_KeyGenSuccess
		keygen.BlockNumber = ctx.BlockHeight()

	} else if ballot.BallotStatus == observerTypes.BallotStatus_BallotFinalized_FailureObservation {
		keygen.Status = observerTypes.KeygenStatus_KeyGenFailed
		keygen.BlockNumber = math2.MaxInt64
	}
	k.zetaObserverKeeper.SetKeygen(ctx, keygen)
	return &types.MsgCreateTSSVoterResponse{}, nil
}

func (k msgServer) UpdateTssAddress(goCtx context.Context, msg *types.MsgUpdateTssAddress) (*types.MsgUpdateTssAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// TODO : Add a new policy type for updating the TSS address
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observerTypes.Policy_Type_update_keygen_block) {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Update can only be executed by the correct policy account")
	}
	tss, ok := k.CheckIfTssPubkeyHasBeenGenerated(ctx, msg.TssPubkey)
	if !ok {
		return nil, errorsmod.Wrap(types.ErrUnableToUpdateTss, "tss pubkey has not been generated")
	}
	k.SetTSS(ctx, tss)
	// initialize the nonces and pending nonces of all enabled chains
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
	return &types.MsgUpdateTssAddressResponse{}, nil
}
