// Copyright 2021 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/evmos/ethermint/blob/main/LICENSE

package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	ibcante "github.com/cosmos/ibc-go/v6/modules/core/ante"
	ethante "github.com/evmos/ethermint/app/ante"
	ethermint "github.com/evmos/ethermint/types"
)

func NewLegacyCosmosAnteHandlerEip712(options ethante.HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		NewAuthzLimiterDecorator(options.DisabledAuthzMsgs...),
		NewVestingAccountDecorator(),
		ante.NewSetUpContextDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ethante.NewMinGasPriceDecorator(options.FeeMarketKeeper, options.EvmKeeper),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		// Note: signature verification uses EIP instead of the cosmos signature validator
		ethante.NewLegacyEip712SigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
		ethante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
	)
}

func newEthAnteHandler(options ethante.HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethante.NewEthSetUpContextDecorator(options.EvmKeeper),                         // outermost AnteDecorator. SetUpContext must be called first
		ethante.NewEthMempoolFeeDecorator(options.EvmKeeper),                           // Check eth effective gas price against minimal-gas-prices
		ethante.NewEthMinGasPriceDecorator(options.FeeMarketKeeper, options.EvmKeeper), // Check eth effective gas price against the global MinGasPrice
		ethante.NewEthValidateBasicDecorator(options.EvmKeeper),
		ethante.NewEthSigVerificationDecorator(options.EvmKeeper),
		ethante.NewEthAccountVerificationDecorator(options.AccountKeeper, options.EvmKeeper),
		ethante.NewCanTransferDecorator(options.EvmKeeper),
		ethante.NewEthGasConsumeDecorator(options.EvmKeeper, options.MaxTxGasWanted),
		ethante.NewEthIncrementSenderSequenceDecorator(options.AccountKeeper), // innermost AnteDecorator.
		ethante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
		ethante.NewEthEmitEventDecorator(options.EvmKeeper), // emit eth tx hash and index at the very last ante handler.
	)
}

func newCosmosAnteHandler(options ethante.HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		NewAuthzLimiterDecorator(options.DisabledAuthzMsgs...),
		NewVestingAccountDecorator(),
		ante.NewSetUpContextDecorator(),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ethante.NewMinGasPriceDecorator(options.FeeMarketKeeper, options.EvmKeeper),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ethante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
	)
}

// this applies to special cosmos tx that calls EVM, in which case the EVM overrides the gas limit
func newCosmosAnteHandlerNoGasLimit(options ethante.HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		NewAuthzLimiterDecorator(options.DisabledAuthzMsgs...),
		NewVestingAccountDecorator(),
		NewSetUpContextDecorator(),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		NewMinGasPriceDecorator(options.FeeMarketKeeper, options.EvmKeeper),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ethante.NewGasWantedDecorator(options.EvmKeeper, options.FeeMarketKeeper),
	)
}

var _ GasTx = (*legacytx.StdTx)(nil) // assert StdTx implements GasTx

// GasTx defines a Tx with a GetGas() method which is needed to use SetUpContextDecorator
type GasTx interface {
	sdk.Tx
	GetGas() uint64
}

// SetUpContextDecorator sets the GasMeter in the Context and wraps the next AnteHandler with a defer clause
// to recover from any downstream OutOfGas panics in the AnteHandler chain to return an error with information
// on gas provided and gas used.
// CONTRACT: Must be first decorator in the chain
// CONTRACT: Tx must implement GasTx interface
type SetUpContextDecorator struct{}

func NewSetUpContextDecorator() SetUpContextDecorator {
	return SetUpContextDecorator{}
}

func (sud SetUpContextDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// all transactions must implement GasTx
	gasTx, ok := tx.(GasTx)
	if !ok {
		// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
		// during runTx.
		newCtx = SetGasMeter(simulate, ctx, 0)
		return newCtx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be GasTx")
	}

	newCtx = SetGasMeter(simulate, ctx, gasTx.GetGas())

	// Decorator will catch an OutOfGasPanic caused in the next antehandler
	// AnteHandlers must have their own defer/recover in order for the BaseApp
	// to know how much gas was used! This is because the GasMeter is created in
	// the AnteHandler, but if it panics the context won't be set properly in
	// runTx's recover call.
	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				log := fmt.Sprintf(
					"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
					rType.Descriptor, gasTx.GetGas(), newCtx.GasMeter().GasConsumed())

				err = sdkerrors.Wrap(sdkerrors.ErrOutOfGas, log)
			default:
				panic(r)
			}
		}
	}()

	return next(newCtx, tx, simulate)
}

// SetGasMeter returns a new context with a gas meter set from a given context.
func SetGasMeter(_ bool, ctx sdk.Context, gasLimit uint64) sdk.Context {
	// In various cases such as simulation and during the genesis block, we do not
	// meter any gas utilization.
	//if simulate || ctx.BlockHeight() == 0 {
	//	return ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	//}
	//
	//return ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))

	return ctx.WithGasMeter(ethermint.NewInfiniteGasMeterWithLimit(gasLimit))
}
