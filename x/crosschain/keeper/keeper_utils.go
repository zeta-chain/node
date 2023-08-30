package keeper

import (
	"fmt"
	"math/big"

	"github.com/zeta-chain/zetacore/cmd/zetacored/config"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// IsAuthorized checks whether a signer is authorized to sign , by checking their address against the observer mapper which contains the observer list for the chain and type
func (k Keeper) IsAuthorizedNodeAccount(ctx sdk.Context, address string) bool {
	_, found := k.zetaObserverKeeper.GetNodeAccount(ctx, address)
	if found {
		return true
	}
	return false
}

// PayGasInZetaAndUpdateCctx updates parameter cctx with the gas price and gas fee for the outbound tx;
// it also makes a trade to fulfill the outbound tx gas fee in ZETA by swapping ZETA for some gas ZRC20 balances
// The gas ZRC20 balance is subsequently burned to account for the expense of TSS address gas fee payment in the outbound tx.
// zetaBurnt represents the amount of Zeta that has been burnt for the tx, the final amount for the tx is zetaBurnt - gasFee
// **Caller should feed temporary ctx into this function**
func (k Keeper) PayGasInZetaAndUpdateCctx(
	ctx sdk.Context,
	chainID int64,
	cctx *types.CrossChainTx,
	zetaBurnt math.Uint,
	noEthereumTxEvent bool,
) error {
	// TODO: this check currently break the integration tests, enable it when integration tests are fixed
	// https://github.com/zeta-chain/node/issues/1022
	//if cctx.InboundTxParams.CoinType != common.CoinType_Zeta {
	//	return sdkerrors.Wrapf(zetaObserverTypes.ErrInvalidCoinType, "can't pay gas in zeta with %s", cctx.InboundTxParams.CoinType.String())
	//}

	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(chainID)
	if chain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}
	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chain.ChainId)
	if !isFound {
		return sdkerrors.Wrap(types.ErrUnableToGetGasPrice, fmt.Sprintf(" chain %d | Identifiers : %s ",
			chain.ChainId,
			cctx.LogIdentifierForCCTX()),
		)
	}
	medianGasPrice = medianGasPrice.MulUint64(2) // overpays gas price by 2x
	cctx.GetCurrentOutTxParam().OutboundTxGasPrice = medianGasPrice.String()
	gasLimit := sdk.NewUint(cctx.GetCurrentOutTxParam().OutboundTxGasLimit)
	outTxGasFee := gasLimit.Mul(medianGasPrice)

	// the following logic computes outbound tx gas fee, and convert into ZETA using system uniswapv2 pool wzeta/gasZRC20
	gasZRC20, err := k.fungibleKeeper.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chain.ChainId))
	if err != nil {
		return sdkerrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to get system contract gas coin")
	}

	outTxGasFeeInZeta, err := k.fungibleKeeper.QueryUniswapv2RouterGetAmountsIn(ctx, outTxGasFee.BigInt(), gasZRC20)
	if err != nil {
		return sdkerrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to QueryUniswapv2RouterGetAmountsIn")
	}
	feeInZeta := types.GetProtocolFee().Add(math.NewUintFromBigInt(outTxGasFeeInZeta))

	cctx.ZetaFees = cctx.ZetaFees.Add(feeInZeta)

	if cctx.ZetaFees.GT(zetaBurnt) {
		return sdkerrors.Wrap(types.ErrNotEnoughZetaBurnt, fmt.Sprintf("feeInZeta(%s) more than zetaBurnt (%s) | Identifiers : %s ",
			cctx.ZetaFees,
			zetaBurnt,
			cctx.LogIdentifierForCCTX()),
		)
	}

	ctx.Logger().Info("Subtracting amount from inbound tx", "amount", zetaBurnt.String(), "feeInZeta",
		cctx.ZetaFees.String())
	cctx.GetCurrentOutTxParam().Amount = zetaBurnt.Sub(cctx.ZetaFees)

	// ** The following logic converts the outTxGasFeeInZeta into gasZRC20 and burns it **
	// swap the outTxGasFeeInZeta portion of zeta to the real gas ZRC20 and burn it, in a temporary context.
	{
		coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewIntFromBigInt(feeInZeta.BigInt())))
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
		if err != nil {
			return sdkerrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to mint coins")
		}

		amounts, err := k.fungibleKeeper.CallUniswapv2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			outTxGasFeeInZeta,
			gasZRC20,
			noEthereumTxEvent,
		)
		if err != nil {
			return sdkerrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to CallUniswapv2RouterSwapExactETHForToken")
		}
		ctx.Logger().Info("gas fee", "outTxGasFee", outTxGasFee, "outTxGasFeeInZeta", outTxGasFeeInZeta)
		ctx.Logger().Info("CallUniswapv2RouterSwapExactETHForToken", "zetaAmountIn", amounts[0], "zrc20AmountOut", amounts[1])
		err = k.fungibleKeeper.CallZRC20Burn(ctx, types.ModuleAddressEVM, gasZRC20, amounts[1], noEthereumTxEvent)
		if err != nil {
			return sdkerrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to CallZRC20Burn")
		}
	}

	return nil
}

// UpdateNonce sets the CCTX outbound nonce to the next nonce, and updates the nonce of blockchain state.
// It also updates the PendingNonces that is used to track the unfulfilled outbound txs.
func (k Keeper) UpdateNonce(ctx sdk.Context, receiveChainID int64, cctx *types.CrossChainTx) error {
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(receiveChainID)
	if chain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}

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
