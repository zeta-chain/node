package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/zetaobserver/types"
)

func (k Keeper) AddVoteToBallot(ctx sdk.Context, ballot zetaObserverTypes.Ballot, address string, observationType zetaObserverTypes.VoteType) (zetaObserverTypes.Ballot, error) {
	ballot, err := ballot.AddVote(address, observationType)
	if err != nil {
		return ballot, err
	}
	ctx.Logger().Info(fmt.Sprintf("Vote Added | Voter :%s, ballot idetifier %s", address, ballot.BallotIdentifier))
	k.zetaObserverKeeper.SetBallot(ctx, &ballot)
	return ballot, err
}
func (k Keeper) CheckIfBallotIsFinalized(ctx sdk.Context, ballot zetaObserverTypes.Ballot) (zetaObserverTypes.Ballot, bool) {
	ballot, isFinalized := ballot.IsBallotFinalized()
	if !isFinalized {
		return ballot, false
	}
	k.zetaObserverKeeper.SetBallot(ctx, &ballot)
	return ballot, true
}

func (k Keeper) IsAuthorized(ctx sdk.Context, address string, senderChain zetaObserverTypes.ObserverChain, observationType string) (bool, error) {
	observerMapper, found := k.zetaObserverKeeper.GetObserverMapper(ctx, senderChain, observationType)
	if !found {
		return false, errors.Wrap(types.ErrNotAuthorized, fmt.Sprintf("Chain/Observation type not supported Chain : %s , Observation type : %s", senderChain, observationType))
	}
	for _, obs := range observerMapper.ObserverList {
		if obs == address {
			return true, nil
		}
	}
	return false, errors.Wrap(types.ErrNotAuthorized, fmt.Sprintf("address: %s", address))
}

func (k Keeper) CheckCCTXExists(ctx sdk.Context, ballotIdentifier, cctxIdentifier string) (cctx types.CrossChainTx, err error) {
	cctx, isFound := k.GetCctxByIndexAndStatuses(ctx,
		cctxIdentifier,
		[]types.CctxStatus{
			types.CctxStatus_PendingOutbound,
			types.CctxStatus_PendingRevert,
		})
	if !isFound {
		return cctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("Cannot find cctx hash %s", cctxIdentifier))
	}
	if cctx.OutBoundTxParams.OutBoundTXBallotIndex == "" {
		cctx.OutBoundTxParams.OutBoundTXBallotIndex = ballotIdentifier
		k.SetCrossChainTx(ctx, cctx)
	}
	return
}
func (k Keeper) GetBallot(ctx sdk.Context, index string, chain zetaObserverTypes.ObserverChain, observationType zetaObserverTypes.ObservationType) (ballot zetaObserverTypes.Ballot, err error) {
	ballot, found := k.zetaObserverKeeper.GetBallot(ctx, index)
	if !found {
		if !k.zetaObserverKeeper.IsChainSupported(ctx, chain) {
			return ballot, sdkerrors.Wrap(types.ErrUnsupportedChain, fmt.Sprintf("Chain %s, Observation %s", chain.String(), observationType.String()))
		}
		observerMapper, _ := k.zetaObserverKeeper.GetObserverMapper(ctx, chain, observationType.String())
		threshohold, found := k.zetaObserverKeeper.GetParams(ctx).GetVotingThreshold(chain, observationType)
		if !found {
			err = errors.Wrap(zetaObserverTypes.ErrSupportedChains, fmt.Sprintf("Thresholds not set for Chain %s and Observation %s", chain.String(), observationType))
			return
		}

		ballot = zetaObserverTypes.Ballot{
			Index:            "",
			BallotIdentifier: index,
			VoterList:        observerMapper.ObserverList,
			Votes:            zetaObserverTypes.CreateVotes(len(observerMapper.ObserverList)),
			ObservationType:  observationType,
			BallotThreshold:  threshohold.Threshold,
			BallotStatus:     zetaObserverTypes.BallotStatus_BallotInProgress,
		}
	}
	return
}

func (k Keeper) UpdatePrices(ctx sdk.Context, receiveChain string, cctx *types.CrossChainTx) error {
	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, receiveChain)
	if !isFound {
		return sdkerrors.Wrap(types.ErrUnableToGetGasPrice, fmt.Sprintf(" chain %s | Identifiers : %s ", cctx.OutBoundTxParams.ReceiverChain, cctx.LogIdentifierForCCTX()))
	}
	gasLimit := sdk.NewUint(cctx.OutBoundTxParams.OutBoundTxGasLimit)
	rate, isFound := k.GetZetaConversionRate(ctx, receiveChain)
	if !isFound {
		return sdkerrors.Wrap(types.ErrUnableToGetConversionRate, fmt.Sprintf(" chain %s | Identifiers : %s ", cctx.OutBoundTxParams.ReceiverChain, cctx.LogIdentifierForCCTX()))
	}
	medianRate := rate.ZetaConversionRates[rate.MedianIndex]
	uintmedianRate := sdk.NewUintFromString(medianRate)
	// Calculate Gas FEE
	gasFeeInZeta := CalculateFee(medianGasPrice, gasLimit, uintmedianRate)

	cctx.OutBoundTxParams.OutBoundTxGasPrice = medianGasPrice.String()

	// Set ZetaBurnt and ZetaMint
	zetaBurnt := cctx.ZetaBurnt
	if gasFeeInZeta.GT(zetaBurnt) {
		return sdkerrors.Wrap(types.ErrNotEnoughZetaBurnt, fmt.Sprintf("feeInZeta(%s) more than mBurnt (%s) | Identifiers : %s ", gasFeeInZeta, zetaBurnt, cctx.LogIdentifierForCCTX()))
	}
	cctx.ZetaFees = cctx.ZetaFees.Add(gasFeeInZeta)
	cctx.ZetaMint = zetaBurnt.Sub(gasFeeInZeta)

	return nil
}
func (k Keeper) UpdateNonce(ctx sdk.Context, receiveChain string, cctx *types.CrossChainTx) error {
	nonce, found := k.GetChainNonces(ctx, receiveChain)
	if !found {
		return sdkerrors.Wrap(types.ErrCannotFindReceiverNonce, fmt.Sprintf("Chain(%s) | Identifiers : %s ", receiveChain, cctx.LogIdentifierForCCTX()))
	}

	// SET nonce
	cctx.OutBoundTxParams.OutBoundTxTSSNonce = nonce.Nonce
	nonce.Nonce++
	k.SetChainNonces(ctx, nonce)
	return nil
}
func CalculateFee(price, gasLimit, rate sdk.Uint) sdk.Uint {
	gasFee := price.Mul(gasLimit).Mul(rate)
	gasFee = reducePrecision(gasFee)
	return gasFee.Add(types.GetProtocolFee())
}
