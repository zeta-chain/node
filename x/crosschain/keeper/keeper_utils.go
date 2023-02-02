package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"math/big"
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

// IsAuthorized checks whether a signer is authorized to sign , by checking thier address against the observer mapper which contains the observer list for the chain and type
func (k Keeper) IsAuthorized(ctx sdk.Context, address string, chain *common.Chain, observationType zetaObserverTypes.ObservationType) (bool, error) {
	observerMapper, found := k.zetaObserverKeeper.GetObserverMapper(ctx, chain, observationType)
	if !found {
		return false, errors.Wrap(types.ErrNotAuthorized, fmt.Sprintf("Mapper Not present | Chain-Observation  %s-%s", chain.String(), observationType))
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
	if cctx.OutboundTxParams.OutboundTxBallotIndex == "" {
		cctx.OutboundTxParams.OutboundTxBallotIndex = ballotIdentifier
		k.SetCrossChainTx(ctx, cctx)
	}
	return
}
func (k Keeper) GetBallot(ctx sdk.Context, index string, chain *common.Chain, observationType zetaObserverTypes.ObservationType) (ballot zetaObserverTypes.Ballot, isNew bool, err error) {
	isNew = false
	ballot, found := k.zetaObserverKeeper.GetBallot(ctx, index)
	if !found {
		observerMapper, _ := k.zetaObserverKeeper.GetObserverMapper(ctx, chain, observationType)
		obsParams := k.zetaObserverKeeper.GetParams(ctx).GetParamsForChain(chain)
		if !obsParams.IsSupported {
			err = errors.Wrap(zetaObserverTypes.ErrSupportedChains, fmt.Sprintf("Thresholds not set for Chain %s and Observation %s", chain.String(), observationType))
			return
		}
		ballot = zetaObserverTypes.Ballot{
			Index:            "",
			BallotIdentifier: index,
			VoterList:        observerMapper.ObserverList,
			Votes:            zetaObserverTypes.CreateVotes(len(observerMapper.ObserverList)),
			ObservationType:  observationType,
			BallotThreshold:  obsParams.BallotThreshold,
			BallotStatus:     zetaObserverTypes.BallotStatus_BallotInProgress,
		}
		isNew = true
	}
	return
}

func (k Keeper) UpdatePrices(ctx sdk.Context, chainID int64, cctx *types.CrossChainTx) error {
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(chainID)
	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chain.ChainId)
	if !isFound {
		return sdkerrors.Wrap(types.ErrUnableToGetGasPrice, fmt.Sprintf(" chain %s | Identifiers : %s ", cctx.OutboundTxParams.ReceiverChain, cctx.LogIdentifierForCCTX()))
	}
	cctx.OutboundTxParams.OutboundTxGasPrice = medianGasPrice.String()
	gasLimit := sdk.NewUint(cctx.OutboundTxParams.OutboundTxGasLimit)

	outTxGasFee := gasLimit.Mul(medianGasPrice)

	zrc20, err := k.fungibleKeeper.QuerySystemContractGasCoinZRC4(ctx, big.NewInt(chain.ChainId))
	if err != nil {
		return sdkerrors.Wrap(err, "UpdatePrices: unable to get system contract gas coin")
	}
	outTxGasFeeInZeta, err := k.fungibleKeeper.QueryUniswapv2RouterGetAmountsIn(ctx, outTxGasFee.BigInt(), zrc20)
	if err != nil {
		return sdkerrors.Wrap(err, "UpdatePrices: unable to QueryUniswapv2RouterGetAmountsIn")
	}
	feeInZeta := types.GetProtocolFee().Add(math.NewUintFromBigInt(outTxGasFeeInZeta))

	// swap the outTxGasFeeInZeta portion of zeta to the real gas ZRC20 and burn it
	coins := sdk.NewCoins(sdk.NewCoin("azeta", sdk.NewIntFromBigInt(feeInZeta.BigInt())))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return sdkerrors.Wrap(err, "UpdatePrices: unable to mint coins")
	}
	amounts, err := k.fungibleKeeper.CallUniswapv2RouterSwapExactETHForToken(ctx, types.ModuleAddressEVM, types.ModuleAddressEVM, outTxGasFeeInZeta, zrc20)
	if err != nil {
		return sdkerrors.Wrap(err, "UpdatePrices: unable to CallUniswapv2RouterSwapExactETHForToken")
	}
	ctx.Logger().Info("gas fee", "outTxGasFee", outTxGasFee, "outTxGasFeeInZeta", outTxGasFeeInZeta)
	ctx.Logger().Info("CallUniswapv2RouterSwapExactETHForToken", "zetaAmountIn", amounts[0], "zrc20AmountOut", amounts[1])
	err = k.fungibleKeeper.CallZRC20Burn(ctx, types.ModuleAddressEVM, zrc20, amounts[1])
	if err != nil {
		return sdkerrors.Wrap(err, "UpdatePrices: unable to CallZRC20Burn")
	}

	cctx.ZetaFees = cctx.ZetaFees.Add(feeInZeta)

	if cctx.ZetaFees.GT(cctx.ZetaBurnt) {
		return sdkerrors.Wrap(types.ErrNotEnoughZetaBurnt, fmt.Sprintf("feeInZeta(%s) more than zetaBurnt (%s) | Identifiers : %s ", cctx.ZetaFees, cctx.ZetaBurnt, cctx.LogIdentifierForCCTX()))
	}
	cctx.ZetaMint = cctx.ZetaBurnt.Sub(cctx.ZetaFees)

	return nil
}

// TODO : USE CHAIN ID
func (k Keeper) UpdateNonce(ctx sdk.Context, receiveChain string, cctx *types.CrossChainTx) error {
	nonce, found := k.GetChainNonces(ctx, receiveChain)
	if !found {
		return sdkerrors.Wrap(types.ErrCannotFindReceiverNonce, fmt.Sprintf("Chain(%s) | Identifiers : %s ", receiveChain, cctx.LogIdentifierForCCTX()))
	}

	// SET nonce
	cctx.OutboundTxParams.OutboundTxTssNonce = nonce.Nonce
	nonce.Nonce++
	k.SetChainNonces(ctx, nonce)
	return nil
}
func CalculateFee(price, gasLimit, rate sdk.Uint) sdk.Uint {
	gasFee := price.Mul(gasLimit).Mul(rate)
	gasFee = reducePrecision(gasFee)
	return gasFee.Add(types.GetProtocolFee())
}
