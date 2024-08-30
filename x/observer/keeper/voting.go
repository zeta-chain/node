package keeper

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/observer/types"
)

const (
	msgVoteOnBallot = "error while voting on ballot"
)

func (k Keeper) AddVoteToBallot(
	ctx sdk.Context,
	ballot types.Ballot,
	address string,
	observationType types.VoteType,
) (types.Ballot, error) {
	ballot, err := ballot.AddVote(address, observationType)
	if err != nil {
		return ballot, err
	}
	ctx.Logger().Info(fmt.Sprintf("Vote Added | Voter :%s, ballot identifier %s", address, ballot.BallotIdentifier))
	k.SetBallot(ctx, &ballot)
	return ballot, nil
}

// CheckIfFinalizingVote checks if the ballot is finalized in this block and if it is, it sets the ballot in the store
// This function with only return true if the ballot moves for pending to success or failed status with this vote.
// If the ballot is already finalized in the previous vote , it will return false
func (k Keeper) CheckIfFinalizingVote(ctx sdk.Context, ballot types.Ballot) (types.Ballot, bool) {
	ballot, isFinalized := ballot.IsFinalizingVote()
	if !isFinalized {
		return ballot, false
	}
	k.SetBallot(ctx, &ballot)
	return ballot, true
}

// IsNonTombstonedObserver checks whether a signer is authorized to sign
// This function checks if the signer is present in the observer set
// and also checks if the signer is not tombstoned
func (k Keeper) IsNonTombstonedObserver(ctx sdk.Context, address string) bool {
	isPresentInMapper := k.IsAddressPartOfObserverSet(ctx, address)
	if !isPresentInMapper {
		return false
	}
	isTombstoned, err := k.IsOperatorTombstoned(ctx, address)
	if err != nil || isTombstoned {
		return false
	}
	return true
}

// FindBallot finds the ballot for the given index
// If the ballot is not found, it creates a new ballot and returns it
func (k Keeper) FindBallot(
	ctx sdk.Context,
	index string,
	chain chains.Chain,
	observationType types.ObservationType,
) (ballot types.Ballot, isNew bool, err error) {
	isNew = false
	ballot, found := k.GetBallot(ctx, index)
	if !found {
		observerSet, _ := k.GetObserverSet(ctx)

		cp, found := k.GetChainParamsByChainID(ctx, chain.ChainId)
		if !found || cp == nil || !cp.IsSupported {
			return types.Ballot{}, false, types.ErrSupportedChains
		}

		ballot = types.Ballot{
			Index:                "",
			BallotIdentifier:     index,
			VoterList:            observerSet.ObserverList,
			Votes:                types.CreateVotes(len(observerSet.ObserverList)),
			ObservationType:      observationType,
			BallotThreshold:      cp.BallotThreshold,
			BallotStatus:         types.BallotStatus_BallotInProgress,
			BallotCreationHeight: ctx.BlockHeight(),
		}
		isNew = true
		k.AddBallotToList(ctx, ballot)
	}
	return
}

func (k Keeper) IsValidator(ctx sdk.Context, creator string) error {
	valAddress, err := types.GetOperatorAddressFromAccAddress(creator)
	if err != nil {
		return err
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return types.ErrNotValidator
	}

	if validator.Jailed || !validator.IsBonded() {
		return types.ErrValidatorStatus
	}
	return nil
}

func (k Keeper) IsOperatorTombstoned(ctx sdk.Context, creator string) (bool, error) {
	valAddress, err := types.GetOperatorAddressFromAccAddress(creator)
	if err != nil {
		return false, err
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return false, types.ErrNotValidator
	}

	consAddress, err := validator.GetConsAddr()
	if err != nil {
		return false, err
	}
	return k.slashingKeeper.IsTombstoned(ctx, consAddress), nil
}

func (k Keeper) CheckObserverSelfDelegation(ctx sdk.Context, accAddress string) error {
	selfdelAddr, err := sdk.AccAddressFromBech32(accAddress)
	if err != nil {
		return err
	}
	valAddress, err := types.GetOperatorAddressFromAccAddress(accAddress)
	if err != nil {
		return err
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return types.ErrNotValidator
	}

	delegation, found := k.stakingKeeper.GetDelegation(ctx, selfdelAddr, valAddress)
	if !found {
		return types.ErrSelfDelegation
	}

	minDelegation, err := types.GetMinObserverDelegationDec()
	if err != nil {
		return err
	}
	tokens := validator.TokensFromShares(delegation.Shares)
	if tokens.LT(minDelegation) {
		k.RemoveObserverFromSet(ctx, accAddress)
	}
	return nil
}

// VoteOnBallot finds a ballot or creates a new one if not found,
// and casts a vote on it. Then proceed to check if the vote has been finalized.
// This function holds generic logic for all types of votes.
func (k Keeper) VoteOnBallot(
	ctx sdk.Context,
	chain chains.Chain,
	ballotIndex string,
	observationType types.ObservationType,
	voter string,
	voteType types.VoteType,
) (
	ballot types.Ballot,
	isFinalized bool,
	isNew bool,
	err error) {
	ballot, isNew, err = k.FindBallot(ctx, ballotIndex, chain, observationType)
	if err != nil {
		return ballot, false, false, sdkerrors.Wrap(err, msgVoteOnBallot)
	}

	ballot, err = k.AddVoteToBallot(ctx, ballot, voter, voteType)
	if err != nil {
		return ballot, false, isNew, sdkerrors.Wrap(err, msgVoteOnBallot)
	}

	ballot, isFinalized = k.CheckIfFinalizingVote(ctx, ballot)

	return ballot, isFinalized, isNew, nil
}
