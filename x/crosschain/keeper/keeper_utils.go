package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
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

// IsAuthorized checks whether a signer is authorized to sign , by checking their address against the observer mapper which contains the observer list for the chain and type
func (k Keeper) IsAuthorized(ctx sdk.Context, address string, chain *common.Chain) (bool, error) {
	observerMapper, found := k.zetaObserverKeeper.GetObserverMapper(ctx, chain)
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

// IsAuthorized checks whether a signer is authorized to sign , by checking their address against the observer mapper which contains the observer list for the chain and type
func (k Keeper) IsAuthorizedNodeAccount(ctx sdk.Context, address string) bool {
	_, found := k.GetNodeAccount(ctx, address)
	if found {
		return true
	}
	return false
}

func (k Keeper) GetBallot(ctx sdk.Context, index string, chain *common.Chain, observationType zetaObserverTypes.ObservationType) (ballot zetaObserverTypes.Ballot, isNew bool, err error) {
	isNew = false
	ballot, found := k.zetaObserverKeeper.GetBallot(ctx, index)
	if !found {
		observerMapper, _ := k.zetaObserverKeeper.GetObserverMapper(ctx, chain)
		obsParams := k.zetaObserverKeeper.GetParams(ctx).GetParamsForChain(chain)
		if !obsParams.IsSupported {
			err = errors.Wrap(zetaObserverTypes.ErrSupportedChains, fmt.Sprintf("Thresholds not set for Chain %s and Observation %s", chain.String(), observationType))
			return
		}
		ballot = zetaObserverTypes.Ballot{
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
		return sdkerrors.Wrap(types.ErrUnableToGetGasPrice, fmt.Sprintf(" chain %d | Identifiers : %s ", cctx.GetCurrentOutTxParam().ReceiverChainId, cctx.LogIdentifierForCCTX()))
	}
	cctx.GetCurrentOutTxParam().OutboundTxGasPrice = medianGasPrice.String()
	gasLimit := sdk.NewUint(cctx.GetCurrentOutTxParam().OutboundTxGasLimit)
	outTxGasFee := gasLimit.Mul(medianGasPrice)

	// the following logic computes outbound tx gas fee, and convert into ZETA using system uniswapv2 pool wzeta/gasZRC20
	gasZRC20, err := k.fungibleKeeper.QuerySystemContractGasCoinZRC4(ctx, big.NewInt(chain.ChainId))
	if err != nil {
		return sdkerrors.Wrap(err, "UpdatePrices: unable to get system contract gas coin")
	}
	outTxGasFeeInZeta, err := k.fungibleKeeper.QueryUniswapv2RouterGetAmountsIn(ctx, outTxGasFee.BigInt(), gasZRC20)
	if err != nil {
		return sdkerrors.Wrap(err, "UpdatePrices: unable to QueryUniswapv2RouterGetAmountsIn")
	}
	feeInZeta := types.GetProtocolFee().Add(math.NewUintFromBigInt(outTxGasFeeInZeta))

	// swap the outTxGasFeeInZeta portion of zeta to the real gas ZRC20 and burn it
	coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewIntFromBigInt(feeInZeta.BigInt())))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return sdkerrors.Wrap(err, "UpdatePrices: unable to mint coins")
	}
	amounts, err := k.fungibleKeeper.CallUniswapv2RouterSwapExactETHForToken(ctx, types.ModuleAddressEVM, types.ModuleAddressEVM, outTxGasFeeInZeta, gasZRC20)
	if err != nil {
		return sdkerrors.Wrap(err, "UpdatePrices: unable to CallUniswapv2RouterSwapExactETHForToken")
	}
	ctx.Logger().Info("gas fee", "outTxGasFee", outTxGasFee, "outTxGasFeeInZeta", outTxGasFeeInZeta)
	ctx.Logger().Info("CallUniswapv2RouterSwapExactETHForToken", "zetaAmountIn", amounts[0], "zrc20AmountOut", amounts[1])
	err = k.fungibleKeeper.CallZRC20Burn(ctx, types.ModuleAddressEVM, gasZRC20, amounts[1])
	if err != nil {
		return sdkerrors.Wrap(err, "UpdatePrices: unable to CallZRC20Burn")
	}

	cctx.ZetaFees = cctx.ZetaFees.Add(feeInZeta)

	if cctx.ZetaFees.GT(cctx.InboundTxParams.Amount) && cctx.InboundTxParams.CoinType == common.CoinType_Zeta {
		return sdkerrors.Wrap(types.ErrNotEnoughZetaBurnt, fmt.Sprintf("feeInZeta(%s) more than zetaBurnt (%s) | Identifiers : %s ", cctx.ZetaFees, cctx.InboundTxParams.Amount, cctx.LogIdentifierForCCTX()))
	}
	cctx.GetCurrentOutTxParam().Amount = cctx.InboundTxParams.Amount.Sub(cctx.ZetaFees)

	return nil
}

func (k Keeper) UpdateNonce(ctx sdk.Context, receiveChainID int64, cctx *types.CrossChainTx) error {
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(receiveChainID)

	nonce, found := k.GetChainNonces(ctx, chain.ChainName.String())
	if !found {
		return sdkerrors.Wrap(types.ErrCannotFindReceiverNonce, fmt.Sprintf("Chain(%s) | Identifiers : %s ", chain.ChainName.String(), cctx.LogIdentifierForCCTX()))
	}

	// SET nonce
	cctx.GetCurrentOutTxParam().OutboundTxTssNonce = nonce.Nonce
	tss, found := k.GetTSS(ctx)
	if !found {
		return sdkerrors.Wrap(types.ErrCannotFindTSSKeys, fmt.Sprintf("Chain(%s) | Identifiers : %s ", chain.ChainName.String(), cctx.LogIdentifierForCCTX()))
	}

	p, found := k.GetPendingNonces(ctx, tss.TssPubkey, uint64(receiveChainID))
	if !found {
		return sdkerrors.Wrap(types.ErrCannotFindPendingNonces, fmt.Sprintf("chain_id %d, nonce %d", receiveChainID, nonce.Nonce))
	}

	if p.NonceHigh != int64(nonce.Nonce) {
		return sdkerrors.Wrap(types.ErrNonceMismatch, fmt.Sprintf("chain_id %d, high nonce %d, current nonce %d", receiveChainID, p.NonceHigh, nonce.Nonce))
	}

	nonce.Nonce++
	p.NonceHigh++
	k.SetChainNonces(ctx, nonce)
	k.SetPendingNonces(ctx, p)
	return nil
}
func CalculateFee(price, gasLimit, rate sdk.Uint) sdk.Uint {
	gasFee := price.Mul(gasLimit).Mul(rate)
	gasFee = reducePrecision(gasFee)
	return gasFee.Add(types.GetProtocolFee())
}
