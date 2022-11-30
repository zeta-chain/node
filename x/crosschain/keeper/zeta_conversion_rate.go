package keeper

//
//func (k Keeper) ZetaConversionRate(ctx context.Context, request *types.QueryZetaConversionRateRequest) (*types.QueryZetaConversionRateResponse, error) {
//	goCtx := sdk.UnwrapSDKContext(ctx)
//	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, receiveChain)
//	if !isFound {
//		return sdkerrors.Wrap(types.ErrUnableToGetGasPrice, fmt.Sprintf(" chain %s | Identifiers : %s ", cctx.OutBoundTxParams.ReceiverChain, cctx.LogIdentifierForCCTX()))
//	}
//	cctx.OutBoundTxParams.OutBoundTxGasPrice = medianGasPrice.String()
//	gasLimit := sdk.NewUint(cctx.OutBoundTxParams.OutBoundTxGasLimit)
//
//	outTxGasFee := gasLimit.Mul(medianGasPrice)
//	recvChain, err := common.ParseChain(receiveChain)
//	if err != nil {
//		return sdkerrors.Wrap(err, "UpdatePrices: unable to parse chain")
//	}
//	chainID := config.Chains[recvChain.String()].ChainID
//	zrc20, err := k.fungibleKeeper.QuerySystemContractGasCoinZRC4(ctx, chainID)
//	if err != nil {
//		return sdkerrors.Wrap(err, "UpdatePrices: unable to get system contract gas coin")
//	}
//	outTxGasFeeInZeta, err := k.fungibleKeeper.QueryUniswapv2RouterGetAmountsIn(ctx, outTxGasFee.BigInt(), zrc20)
//	if err != nil {
//		return sdkerrors.Wrap(err, "UpdatePrices: unable to QueryUniswapv2RouterGetAmountsIn")
//	}
//	feeInZeta := types.GetProtocolFee().Add(sdk.NewUintFromBigInt(outTxGasFeeInZeta))
//
//	// swap the outTxGasFeeInZeta portion of zeta to the real gas ZRC20 and burn it
//	coins := sdk.NewCoins(sdk.NewCoin("azeta", sdk.NewIntFromBigInt(feeInZeta.BigInt())))
//	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
//	if err != nil {
//		return sdkerrors.Wrap(err, "UpdatePrices: unable to mint coins")
//	}
//
//}
