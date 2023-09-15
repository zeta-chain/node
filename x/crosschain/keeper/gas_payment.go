package keeper

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// PayGasAndUpdateCctx updates the outbound tx with the new amount after paying the gas fee
// **Caller should feed temporary ctx into this function**
func (k Keeper) PayGasAndUpdateCctx(
	ctx sdk.Context,
	chainID int64,
	cctx *types.CrossChainTx,
	inputAmount math.Uint,
	noEthereumTxEvent bool,
) error {
	// Dispatch to the correct function based on the coin type
	switch cctx.InboundTxParams.CoinType {
	case common.CoinType_Zeta:
		return k.PayGasInZetaAndUpdateCctx(ctx, chainID, cctx, inputAmount, noEthereumTxEvent)
	case common.CoinType_Gas:
		return k.PayGasNativeAndUpdateCctx(ctx, chainID, cctx, inputAmount)
	case common.CoinType_ERC20:
		return k.PayGasInERC20AndUpdateCctx(ctx, chainID, cctx, inputAmount, noEthereumTxEvent)
	default:
		// can't pay gas with coin type
		return fmt.Errorf("can't pay gas with coin type %s", cctx.InboundTxParams.CoinType.String())
	}
}

// ChainGasParams returns the params to calculates the fees for gas for a chain
// tha gas address, the gas limit, gas price and protocol flat fee are returned
func (k Keeper) ChainGasParams(
	ctx sdk.Context,
	chainID int64,
) (gasZRC20 ethcommon.Address, gasLimit, gasPrice, protocolFee math.Uint, err error) {
	gasZRC20, err = k.fungibleKeeper.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return gasZRC20, gasLimit, gasPrice, protocolFee, err
	}

	// get the gas limit
	gasLimitQueried, err := k.fungibleKeeper.QueryGasLimit(ctx, gasZRC20)
	if err != nil {
		return gasZRC20, gasLimit, gasPrice, protocolFee, err
	}
	if gasLimitQueried == nil {
		return gasZRC20, gasLimit, gasPrice, protocolFee, errors.New("gas limit is nil")
	}
	gasLimit = math.NewUintFromBigInt(gasLimitQueried)

	// get the protocol flat fee
	protocolFlatFeeQueried, err := k.fungibleKeeper.QueryProtocolFlatFee(ctx, gasZRC20)
	if err != nil {
		return gasZRC20, gasLimit, gasPrice, protocolFee, err
	}
	if protocolFlatFeeQueried == nil {
		return gasZRC20, gasLimit, gasPrice, protocolFee, errors.New("protocol flat fee is nil")
	}
	protocolFee = math.NewUintFromBigInt(protocolFlatFeeQueried)

	// get the gas price
	gasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chainID)
	if !isFound {
		return gasZRC20, gasLimit, gasPrice, protocolFee, types.ErrUnableToGetGasPrice
	}

	return
}

// PayGasNativeAndUpdateCctx updates the outbound tx with the new amount subtracting the gas fee
// **Caller should feed temporary ctx into this function**
func (k Keeper) PayGasNativeAndUpdateCctx(
	ctx sdk.Context,
	chainID int64,
	cctx *types.CrossChainTx,
	inputAmount math.Uint,
) error {
	// preliminary checks
	if cctx.InboundTxParams.CoinType != common.CoinType_Gas {
		return sdkerrors.Wrapf(zetaObserverTypes.ErrInvalidCoinType, "can't pay gas in native gas with %s", cctx.InboundTxParams.CoinType.String())
	}
	if chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(chainID); chain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}

	// get gas params
	_, gasLimit, gasPrice, protocolFlatFee, err := k.ChainGasParams(ctx, chainID)
	if err != nil {
		return err
	}

	// calculate the final gas fee
	outTxGasFee := gasLimit.Mul(gasPrice).Add(protocolFlatFee)

	// subtract the withdraw fee from the input amount
	if outTxGasFee.GT(inputAmount) {
		return sdkerrors.Wrap(types.ErrNotEnoughGas, fmt.Sprintf("outTxGasFee(%s) more than available gas for tx (%s) | Identifiers : %s ",
			outTxGasFee,
			inputAmount,
			cctx.LogIdentifierForCCTX()),
		)
	}
	ctx.Logger().Info("Subtracting amount from inbound tx", "amount", inputAmount.String(), "fee", outTxGasFee.String())
	newAmount := inputAmount.Sub(outTxGasFee)

	// update cctx
	cctx.GetCurrentOutTxParam().Amount = newAmount
	cctx.GetCurrentOutTxParam().OutboundTxGasLimit = gasLimit.Uint64()
	cctx.GetCurrentOutTxParam().OutboundTxGasPrice = gasPrice.String()

	return nil
}

// PayGasInERC20AndUpdateCctx updates parameter cctx amount subtracting the gas fee
// the gas fee in ERC20 is calculated by swapping ERC20 -> Zeta -> Gas
// if the route is not available, the gas payment will fail
func (k Keeper) PayGasInERC20AndUpdateCctx(
	ctx sdk.Context,
	chainID int64,
	cctx *types.CrossChainTx,
	inputAmount math.Uint,
	noEthereumTxEvent bool,
) error {
	// preliminary checks
	if cctx.InboundTxParams.CoinType != common.CoinType_ERC20 {
		return sdkerrors.Wrapf(zetaObserverTypes.ErrInvalidCoinType, "can't pay gas in erc20 with %s", cctx.InboundTxParams.CoinType.String())
	}
	if chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(chainID); chain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}

	// get gas params
	gasZRC20, gasLimit, gasPrice, protocolFlatFee, err := k.ChainGasParams(ctx, chainID)
	if err != nil {
		return err
	}

	// calculate the final gas fee
	outTxGasFee := gasLimit.Mul(gasPrice).Add(protocolFlatFee)

	// get the necessary ERC20 amount for gas
	_, err = k.fungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, outTxGasFee.BigInt(), gasZRC20)
	if err != nil {
		return err
	}

	return nil
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
	// preliminary checks
	if cctx.InboundTxParams.CoinType != common.CoinType_Zeta {
		return sdkerrors.Wrapf(zetaObserverTypes.ErrInvalidCoinType, "can't pay gas in zeta with %s", cctx.InboundTxParams.CoinType.String())
	}
	if chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(chainID); chain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}

	gasZRC20, err := k.fungibleKeeper.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return sdkerrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to get system contract gas coin")
	}

	// get the gas price
	gasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chainID)
	if !isFound {
		return sdkerrors.Wrap(types.ErrUnableToGetGasPrice, fmt.Sprintf(" chain %d | Identifiers : %s ",
			chainID,
			cctx.LogIdentifierForCCTX()),
		)
	}
	gasPrice = gasPrice.MulUint64(2) // overpays gas price by 2x

	// get the gas fee in gas token
	gasLimit := sdk.NewUint(cctx.GetCurrentOutTxParam().OutboundTxGasLimit)
	outTxGasFee := gasLimit.Mul(gasPrice)

	// get the gas fee in Zeta using system uniswapv2 pool wzeta/gasZRC20 and adding the protocol fee
	outTxGasFeeInZeta, err := k.fungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, outTxGasFee.BigInt(), gasZRC20)
	if err != nil {
		return sdkerrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to QueryUniswapv2RouterGetAmountsIn")
	}
	feeInZeta := types.GetProtocolFee().Add(math.NewUintFromBigInt(outTxGasFeeInZeta))

	// reduce the amount of the outbound tx
	if feeInZeta.GT(zetaBurnt) {
		return sdkerrors.Wrap(types.ErrNotEnoughZetaBurnt, fmt.Sprintf("feeInZeta(%s) more than zetaBurnt (%s) | Identifiers : %s ",
			feeInZeta,
			zetaBurnt,
			cctx.LogIdentifierForCCTX()),
		)
	}
	ctx.Logger().Info("Subtracting amount from inbound tx", "amount", zetaBurnt.String(), "feeInZeta", feeInZeta.String())
	newAmount := zetaBurnt.Sub(feeInZeta)

	// ** The following logic converts the outTxGasFeeInZeta into gasZRC20 and burns it **
	// swap the outTxGasFeeInZeta portion of zeta to the real gas ZRC20 and burn it, in a temporary context.
	{
		coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewIntFromBigInt(feeInZeta.BigInt())))
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
		if err != nil {
			return sdkerrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to mint coins")
		}

		amounts, err := k.fungibleKeeper.CallUniswapV2RouterSwapExactETHForToken(
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

	// Update the cctx
	cctx.GetCurrentOutTxParam().OutboundTxGasPrice = gasPrice.String()
	cctx.GetCurrentOutTxParam().Amount = newAmount
	if cctx.ZetaFees.IsNil() {
		cctx.ZetaFees = feeInZeta
	} else {
		cctx.ZetaFees = cctx.ZetaFees.Add(feeInZeta)
	}

	return nil
}
