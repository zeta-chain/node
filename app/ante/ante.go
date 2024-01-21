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
	"runtime/debug"

	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

func ValidateHandlerOptions(options HandlerOptions) error {
	if options.AccountKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.FeeMarketKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "fee market keeper is required for AnteHandler")
	}
	if options.EvmKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "evm keeper is required for AnteHandler")
	}
	if options.ObserverKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "observer keeper is required for AnteHandler")
	}
	return nil
}

// NewAnteHandler returns an ante handler responsible for attempting to route an
// Ethereum or SDK transaction to an internal ante handler for performing
// transaction-level processing (e.g. fee payment, signature verification) before
// being passed onto it's respective handler.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if err := ValidateHandlerOptions(options); err != nil {
		return nil, err
	}

	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		defer Recover(ctx.Logger(), &err)

		txWithExtensions, ok := tx.(authante.HasExtensionOptionsTx)
		if ok {
			opts := txWithExtensions.GetExtensionOptions()
			if len(opts) > 0 {
				switch typeURL := opts[0].GetTypeUrl(); typeURL {
				case "/ethermint.evm.v1.ExtensionOptionsEthereumTx":
					// handle as *evmtypes.MsgEthereumTx
					anteHandler = newEthAnteHandler(options)
				case "/ethermint.types.v1.ExtensionOptionsWeb3Tx":
					// Deprecated: Handle as normal Cosmos SDK tx, except signature is checked for Legacy EIP712 representation
					anteHandler = NewLegacyCosmosAnteHandlerEip712(options)
				case "/ethermint.types.v1.ExtensionOptionDynamicFeeTx":
					// cosmos-sdk tx with dynamic fee extension
					anteHandler = newCosmosAnteHandler(options)
				default:
					return ctx, errorsmod.Wrapf(
						errortypes.ErrUnknownExtensionOptions,
						"rejecting tx with unsupported extension option: %s", typeURL,
					)
				}

				return anteHandler(ctx, tx, sim)
			}
		}

		// handle as totally normal Cosmos SDK tx
		switch tx.(type) {
		case sdk.Tx:
			// default: handle as normal Cosmos SDK tx
			anteHandler = newCosmosAnteHandler(options)

			// the following determines whether the tx is a system tx which will uses different handler
			// System txs are always single Msg txs, optionally wrapped by one level of MsgExec
			if len(tx.GetMsgs()) != 1 { // this is not a system tx
				break
			}
			msg := tx.GetMsgs()[0]

			// if wrapped inside a MsgExec, unwrap it and reveal the innerMsg.
			var innerMsg sdk.Msg
			innerMsg = msg
			if mm, ok := msg.(*authz.MsgExec); ok { // authz tx; look inside it
				msgs, err := mm.GetMessages()
				if err == nil && len(msgs) == 1 {
					innerMsg = msgs[0]
				}
			}

			// is authorized checks if the creator of the tx is in the observer set
			isAuthorized := options.ObserverKeeper.IsAuthorized
			if mm, ok := innerMsg.(*cctxtypes.MsgGasPriceVoter); ok && isAuthorized(ctx, mm.Creator) {
				anteHandler = newCosmosAnteHandlerForSystemTx(options)
			} else if mm, ok := innerMsg.(*cctxtypes.MsgVoteOnObservedInboundTx); ok && isAuthorized(ctx, mm.Creator) {
				anteHandler = newCosmosAnteHandlerForSystemTx(options)
			} else if mm, ok := innerMsg.(*cctxtypes.MsgVoteOnObservedOutboundTx); ok && isAuthorized(ctx, mm.Creator) {
				anteHandler = newCosmosAnteHandlerForSystemTx(options)
			} else if mm, ok := innerMsg.(*cctxtypes.MsgAddToOutTxTracker); ok && isAuthorized(ctx, mm.Creator) {
				anteHandler = newCosmosAnteHandlerForSystemTx(options)
			} else if mm, ok := innerMsg.(*cctxtypes.MsgCreateTSSVoter); ok && isAuthorized(ctx, mm.Creator) {
				anteHandler = newCosmosAnteHandlerForSystemTx(options)
			} else if mm, ok := innerMsg.(*observertypes.MsgAddBlockHeader); ok && isAuthorized(ctx, mm.Creator) {
				anteHandler = newCosmosAnteHandlerForSystemTx(options)
			} else if mm, ok := innerMsg.(*observertypes.MsgAddBlameVote); ok && isAuthorized(ctx, mm.Creator) {
				anteHandler = newCosmosAnteHandlerForSystemTx(options)
			} else if _, ok := innerMsg.(*stakingtypes.MsgCreateValidator); ok && ctx.BlockHeight() == 0 {
				anteHandler = newCosmosAnteHandlerForSystemTx(options)
			}

		default:
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}

		return anteHandler(ctx, tx, sim)
	}, nil
}

func Recover(logger tmlog.Logger, err *error) {
	if r := recover(); r != nil {
		if err != nil {
			// #nosec G703 err is checked non-nil above
			*err = errorsmod.Wrapf(errortypes.ErrPanic, "%v", r)
		}

		if e, ok := r.(error); ok {
			logger.Error(
				"ante handler panicked",
				"error", e,
				"stack trace", string(debug.Stack()),
			)
		} else {
			logger.Error(
				"ante handler panicked",
				"recover", fmt.Sprintf("%v", r),
			)
		}
	}
}
