package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"

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
	ctx.Logger().Debug("vote added",
		"voter", address,
		"ballot_identifier", ballot.BallotIdentifier)
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

// CheckObserverCanVote checks if the address is a valid observer
func (k Keeper) CheckObserverCanVote(ctx sdk.Context, address string) error {
	isActiveObserver := k.IsAddressPartOfObserverSet(ctx, address)
	if !isActiveObserver {
		return sdkerrors.Wrapf(types.ErrNotObserver, "address is not part of the observer set: %s", address)
	}
	valAddress, err := types.GetOperatorAddressFromAccAddress(address)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidAddress, "invalid operator address for observer : %s", address)
	}
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddress)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrNotValidator, "observer is not a validator: %s", address)
	}
	if validator.Jailed {
		return sdkerrors.Wrapf(types.ErrValidatorJailed, "observer is jailed: %s", address)
	}
	if validator.Status != stakingtypes.Bonded {
		return sdkerrors.Wrapf(types.ErrValidatorStatus, "observer is not bonded: %s", address)
	}
	consAddress, err := validator.GetConsAddr()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidAddress, "invalid consensus address for observer: %s", address)
	}
	if k.slashingKeeper.IsTombstoned(ctx, consAddress) {
		return sdkerrors.Wrapf(types.ErrValidatorTombstoned, "observer is tombstoned: %s", address)
	}
	return nil
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
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddress)
	if err != nil {
		return err
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
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddress)
	if err != nil {
		return false, err
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
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddress)
	if err != nil {
		return errors.Wrapf(err, "validator : %s", valAddress)
	}

	delegation, err := k.stakingKeeper.GetDelegation(ctx, selfdelAddr, valAddress)
	if err != nil {
		return errors.Wrapf(types.ErrSelfDelegation, "self delegation : %s , valAddres : %s", selfdelAddr, valAddress)
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
