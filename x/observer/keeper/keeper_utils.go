package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) AddVoteToBallot(ctx sdk.Context, ballot types.Ballot, address string, observationType types.VoteType) (types.Ballot, error) {
	ballot, err := ballot.AddVote(address, observationType)
	if err != nil {
		return ballot, err
	}
	ctx.Logger().Info(fmt.Sprintf("Vote Added | Voter :%s, ballot idetifier %s", address, ballot.BallotIdentifier))
	k.SetBallot(ctx, &ballot)
	return ballot, err
}

// CheckIfFinalizingVote checks if the ballot is finalized in this block and if it is, it sets the ballot in the store
// This function with only return true if the ballot moves for pending to success or failed status with this vote.
// If the ballot is already finalized in the previous vote , it will return false
func (k Keeper) CheckIfFinalizingVote(ctx sdk.Context, ballot types.Ballot) (types.Ballot, bool) {
	ballot, isFinalized := ballot.IsBallotFinalized()
	if !isFinalized {
		return ballot, false
	}
	k.SetBallot(ctx, &ballot)
	return ballot, true
}

// IsAuthorized checks whether a signer is authorized to sign , by checking their address against the observer mapper which contains the observer list for the chain and type
func (k Keeper) IsAuthorized(ctx sdk.Context, address string, chain *common.Chain) (bool, error) {
	observerMapper, found := k.GetObserverMapper(ctx, chain)
	if !found {
		return false, errors.Wrap(types.ErrNotAuthorized, fmt.Sprintf("observer list not present for chain %s", chain.String()))
	}
	for _, obs := range observerMapper.ObserverList {
		if obs == address {
			return true, nil
		}
	}
	return false, errors.Wrap(types.ErrNotAuthorized, fmt.Sprintf("address: %s", address))
}

func (k Keeper) FindBallot(ctx sdk.Context, index string, chain *common.Chain, observationType types.ObservationType) (ballot types.Ballot, isNew bool, err error) {
	isNew = false
	ballot, found := k.GetBallot(ctx, index)
	if !found {
		observerMapper, _ := k.GetObserverMapper(ctx, chain)
		obsParams := k.GetParams(ctx).GetParamsForChain(chain)
		if !obsParams.IsSupported {
			err = errors.Wrap(types.ErrSupportedChains, fmt.Sprintf("Thresholds not set for Chain %s and Observation %s", chain.String(), observationType))
			return
		}
		ballot = types.Ballot{
			Index:            "",
			BallotIdentifier: index,
			VoterList:        observerMapper.ObserverList,
			Votes:            types.CreateVotes(len(observerMapper.ObserverList)),
			ObservationType:  observationType,
			BallotThreshold:  obsParams.BallotThreshold,
			BallotStatus:     types.BallotStatus_BallotInProgress,
		}
		isNew = true
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

	if validator.Jailed == true || validator.IsBonded() == false {
		return types.ErrValidatorStatus
	}
	return nil

}

func (k Keeper) CheckObserverDelegation(ctx sdk.Context, accAddress string, chain *common.Chain) error {
	selfdelAddr, _ := sdk.AccAddressFromBech32(accAddress)
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
	obsParams := k.GetParams(ctx).GetParamsForChain(chain)
	if !obsParams.IsSupported {
		return errors.Wrap(types.ErrSupportedChains, fmt.Sprintf("Chain not suported %s ", chain.String()))
	}

	tokens := validator.TokensFromShares(delegation.Shares)
	if tokens.LT(obsParams.MinObserverDelegation) {
		return types.ErrCheckObserverDelegation
	}
	return nil
}
