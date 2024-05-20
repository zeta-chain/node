package keeper

import (
	"errors"
	"fmt"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// PayGasAndUpdateCctx updates the outbound tx with the new amount after paying the gas fee
// **Caller should feed temporary ctx into this function**
// chainID is the outbound chain chain id , this can be receiver chain for regular transactions and sender-chain to reverted transactions
func (k Keeper) PayGasAndUpdateCctx(
	ctx sdk.Context,
	chainID int64,
	cctx *types.CrossChainTx,
	inputAmount math.Uint,
	noEthereumTxEvent bool,
) error {
	// Dispatch to the correct function based on the coin type
	switch cctx.InboundParams.CoinType {
	case coin.CoinType_Zeta:
		return k.PayGasInZetaAndUpdateCctx(ctx, chainID, cctx, inputAmount, noEthereumTxEvent)
	case coin.CoinType_Gas:
		return k.PayGasNativeAndUpdateCctx(ctx, chainID, cctx, inputAmount)
	case coin.CoinType_ERC20:
		return k.PayGasInERC20AndUpdateCctx(ctx, chainID, cctx, inputAmount, noEthereumTxEvent)
	default:
		// can't pay gas with coin type
		return fmt.Errorf("can't pay gas with coin type %s", cctx.InboundParams.CoinType.String())
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
	if cctx.InboundParams.CoinType != coin.CoinType_Gas {
		return cosmoserrors.Wrapf(types.ErrInvalidCoinType, "can't pay gas in native gas with %s", cctx.InboundParams.CoinType.String())
	}
	if chain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, chainID); chain == nil {
		return observertypes.ErrSupportedChains
	}

	// get gas params
	_, gasLimit, gasPrice, protocolFlatFee, err := k.ChainGasParams(ctx, chainID)
	if err != nil {
		return cosmoserrors.Wrap(types.ErrCannotFindGasParams, err.Error())
	}

	// calculate the final gas fee
	outTxGasFee := gasLimit.Mul(gasPrice).Add(protocolFlatFee)

	// subtract the withdraw fee from the input amount
	if outTxGasFee.GT(inputAmount) {
		return cosmoserrors.Wrap(types.ErrNotEnoughGas, fmt.Sprintf("outTxGasFee(%s) more than available gas for tx (%s) | Identifiers : %s ",
			outTxGasFee,
			inputAmount,
			cctx.LogIdentifierForCCTX()),
		)
	}
	ctx.Logger().Info("Subtracting amount from inbound tx", "amount", inputAmount.String(), "fee", outTxGasFee.String())
	newAmount := inputAmount.Sub(outTxGasFee)

	// update cctx
	cctx.GetCurrentOutboundParam().Amount = newAmount
	cctx.GetCurrentOutboundParam().GasLimit = gasLimit.Uint64()
	cctx.GetCurrentOutboundParam().GasPrice = gasPrice.String()

	return nil
}

// PayGasInERC20AndUpdateCctx updates parameter cctx amount subtracting the gas fee
// the gas fee in ERC20 is calculated by swapping ERC20 -> Zeta -> Gas
// if the route is not available, the gas payment will fail
// **Caller should feed temporary ctx into this function**
func (k Keeper) PayGasInERC20AndUpdateCctx(
	ctx sdk.Context,
	chainID int64,
	cctx *types.CrossChainTx,
	inputAmount math.Uint,
	noEthereumTxEvent bool,
) error {
	// preliminary checks
	if cctx.InboundParams.CoinType != coin.CoinType_ERC20 {
		return cosmoserrors.Wrapf(types.ErrInvalidCoinType, "can't pay gas in erc20 with %s", cctx.InboundParams.CoinType.String())
	}

	if chain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, chainID); chain == nil {
		return observertypes.ErrSupportedChains
	}
	// get gas params
	gasZRC20, gasLimit, gasPrice, protocolFlatFee, err := k.ChainGasParams(ctx, chainID)
	if err != nil {
		return cosmoserrors.Wrap(types.ErrCannotFindGasParams, err.Error())
	}
	outTxGasFee := gasLimit.Mul(gasPrice).Add(protocolFlatFee)
	// get address of the zrc20
	fc, found := k.fungibleKeeper.GetForeignCoinFromAsset(ctx, cctx.InboundParams.Asset, chainID)
	if !found {
		return cosmoserrors.Wrapf(types.ErrForeignCoinNotFound, "zrc20 from asset %s not found", cctx.InboundParams.Asset)
	}
	zrc20 := ethcommon.HexToAddress(fc.Zrc20ContractAddress)
	if zrc20 == (ethcommon.Address{}) {
		return cosmoserrors.Wrapf(types.ErrForeignCoinNotFound, "zrc20 from asset %s invalid address", cctx.InboundParams.Asset)
	}

	// get the necessary ERC20 amount for gas
	feeInZRC20, err := k.fungibleKeeper.QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(ctx, outTxGasFee.BigInt(), zrc20, gasZRC20)
	if err != nil {
		// NOTE: this is the first method that fails when a liquidity pool is not set for the gas ZRC20, so we return a specific error
		return cosmoserrors.Wrap(types.ErrNoLiquidityPool, err.Error())
	}

	// subtract the withdraw fee from the input amount
	if math.NewUintFromBigInt(feeInZRC20).GT(inputAmount) {
		return cosmoserrors.Wrap(types.ErrNotEnoughGas, fmt.Sprintf("feeInZRC20(%s) more than available gas for tx (%s) | Identifiers : %s ",
			feeInZRC20,
			inputAmount,
			cctx.LogIdentifierForCCTX()),
		)
	}
	newAmount := inputAmount.Sub(math.NewUintFromBigInt(feeInZRC20))

	// mint the amount of ERC20 to be burnt as gas fee
	_, err = k.fungibleKeeper.DepositZRC20(ctx, zrc20, types.ModuleAddressEVM, feeInZRC20)
	if err != nil {
		return cosmoserrors.Wrap(fungibletypes.ErrContractCall, err.Error())
	}
	ctx.Logger().Info("Minted ERC20 for gas fee",
		"zrc20", zrc20.Hex(),
		"amount", feeInZRC20,
	)

	// approve the uniswapv2 router to spend the ERC20
	routerAddress, err := k.fungibleKeeper.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return cosmoserrors.Wrap(fungibletypes.ErrContractCall, err.Error())
	}
	err = k.fungibleKeeper.CallZRC20Approve(
		ctx,
		types.ModuleAddressEVM,
		zrc20,
		routerAddress,
		feeInZRC20,
		noEthereumTxEvent,
	)
	if err != nil {
		return cosmoserrors.Wrap(fungibletypes.ErrContractCall, err.Error())
	}

	// swap the fee in ERC20 into gas passing through Zeta and burn the gas ZRC20
	amounts, err := k.fungibleKeeper.CallUniswapV2RouterSwapExactTokensForTokens(
		ctx,
		types.ModuleAddressEVM,
		types.ModuleAddressEVM,
		feeInZRC20,
		zrc20,
		gasZRC20,
		noEthereumTxEvent,
	)
	if err != nil {
		return cosmoserrors.Wrap(fungibletypes.ErrContractCall, err.Error())
	}
	ctx.Logger().Info("CallUniswapV2RouterSwapExactTokensForTokens",
		"zrc20AmountIn", amounts[0],
		"gasAmountOut", amounts[2],
	)
	gasObtained := amounts[2]

	// FIXME: investigate small mismatches between gasObtained and outTxGasFee
	// https://github.com/zeta-chain/node/issues/1303
	// check if the final gas received after swap matches the gas fee defined
	// if not there might be issues with the pool liquidity and it is safer from an accounting perspective to return an error
	if gasObtained.Cmp(outTxGasFee.BigInt()) == -1 {
		return cosmoserrors.Wrapf(types.ErrInvalidGasAmount, "gas obtained for burn (%s) is lower than gas fee(%s)", gasObtained, outTxGasFee)
	}

	// burn the gas ZRC20
	err = k.fungibleKeeper.CallZRC20Burn(ctx, types.ModuleAddressEVM, gasZRC20, gasObtained, noEthereumTxEvent)
	if err != nil {
		return cosmoserrors.Wrap(fungibletypes.ErrContractCall, err.Error())
	}
	ctx.Logger().Info("Burning gas ZRC20",
		"zrc20", gasZRC20.Hex(),
		"amount", gasObtained,
	)

	// update cctx
	cctx.GetCurrentOutboundParam().Amount = newAmount
	cctx.GetCurrentOutboundParam().GasLimit = gasLimit.Uint64()
	cctx.GetCurrentOutboundParam().GasPrice = gasPrice.String()

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
	if cctx.InboundParams.CoinType != coin.CoinType_Zeta {
		return cosmoserrors.Wrapf(types.ErrInvalidCoinType, "can't pay gas in zeta with %s", cctx.InboundParams.CoinType.String())
	}

	if chain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, chainID); chain == nil {
		return observertypes.ErrSupportedChains
	}

	gasZRC20, err := k.fungibleKeeper.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return cosmoserrors.Wrapf(err, "PayGasInZetaAndUpdateCctx: unable to get system contract gas coin, chaind ID %d", chainID)
	}

	// get the gas price
	gasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chainID)
	if !isFound {
		return cosmoserrors.Wrap(types.ErrUnableToGetGasPrice, fmt.Sprintf(" chain %d | Identifiers : %s ",
			chainID,
			cctx.LogIdentifierForCCTX()),
		)
	}
	gasPrice = gasPrice.MulUint64(2) // overpays gas price by 2x

	// get the gas fee in gas token
	gasLimit := sdk.NewUint(cctx.GetCurrentOutboundParam().GasLimit)
	outTxGasFee := gasLimit.Mul(gasPrice)

	// get the gas fee in Zeta using system uniswapv2 pool wzeta/gasZRC20 and adding the protocol fee
	outTxGasFeeInZeta, err := k.fungibleKeeper.QueryUniswapV2RouterGetZetaAmountsIn(ctx, outTxGasFee.BigInt(), gasZRC20)
	if err != nil {
		return cosmoserrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to QueryUniswapV2RouterGetZetaAmountsIn")
	}
	feeInZeta := types.GetProtocolFee().Add(math.NewUintFromBigInt(outTxGasFeeInZeta))
	// reduce the amount of the outbound tx
	if feeInZeta.GT(zetaBurnt) {
		return cosmoserrors.Wrap(types.ErrNotEnoughZetaBurnt, fmt.Sprintf("feeInZeta(%s) more than zetaBurnt (%s) | Identifiers : %s ",
			feeInZeta,
			zetaBurnt,
			cctx.LogIdentifierForCCTX()),
		)
	}
	ctx.Logger().Info("Subtracting amount from inbound tx",
		"amount", zetaBurnt.String(),
		"feeInZeta", feeInZeta.String(),
	)
	newAmount := zetaBurnt.Sub(feeInZeta)

	// ** The following logic converts the outTxGasFeeInZeta into gasZRC20 and burns it **
	// swap the outTxGasFeeInZeta portion of zeta to the real gas ZRC20 and burn it, in a temporary context.
	{
		coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewIntFromBigInt(feeInZeta.BigInt())))
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
		if err != nil {
			return cosmoserrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to mint coins")
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
			return cosmoserrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to CallUniswapv2RouterSwapExactETHForToken")
		}

		ctx.Logger().Info("gas fee", "outTxGasFee", outTxGasFee, "outTxGasFeeInZeta", outTxGasFeeInZeta)
		ctx.Logger().Info("CallUniswapv2RouterSwapExactETHForToken",
			"zetaAmountIn", amounts[0],
			"zrc20AmountOut", amounts[1],
		)

		// FIXME: investigate small mismatches between amounts[1] and outTxGasFee
		// https://github.com/zeta-chain/node/issues/1303
		err = k.fungibleKeeper.CallZRC20Burn(ctx, types.ModuleAddressEVM, gasZRC20, amounts[1], noEthereumTxEvent)
		if err != nil {
			return cosmoserrors.Wrap(err, "PayGasInZetaAndUpdateCctx: unable to CallZRC20Burn")
		}
	}

	// Update the cctx
	cctx.GetCurrentOutboundParam().GasPrice = gasPrice.String()
	cctx.GetCurrentOutboundParam().Amount = newAmount
	if cctx.ZetaFees.IsNil() {
		cctx.ZetaFees = feeInZeta
	} else {
		cctx.ZetaFees = cctx.ZetaFees.Add(feeInZeta)
	}

	return nil
}
