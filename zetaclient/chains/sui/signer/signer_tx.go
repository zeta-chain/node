package signer

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// txBuilder represents a wrapper that returns Sui tx with TSS signature (base64)
//
// We use a "builder" pattern to delay RPC 'MoveCallRequest'.
// This simplifies downstream usage of tx broadcasting & signing.
// It also help to avoid Gateway's object version mismatch error.
type txBuilder func(ctx context.Context) (models.TxnMetaData, string, error)

const (
	funcWithdraw      = "withdraw"
	funcIncreaseNonce = "increase_nonce"

	// minGasBudgetCancelTx is the minimum gas budget for the cancel tx
	minGasBudgetCancelTx = 2_000_000
)

// TODO: use these functions in PTB building
// https://github.com/zeta-chain/node/issues/3741
//const funcWithdrawImpl = "withdraw_impl"
//const funcOnCall = "on_call"

func (s *Signer) createWithdrawTxBuilder(cctx *cctypes.CrossChainTx, zetaHeight uint64) (txBuilder, error) {
	return func(ctx context.Context) (models.TxnMetaData, string, error) {
		tx, err := s.buildWithdrawal(ctx, cctx)
		if err != nil {
			return models.TxnMetaData{}, "", errors.Wrap(err, "unable to build withdrawal tx")
		}

		nonce := cctx.GetCurrentOutboundParam().TssNonce

		sigBase64, err := s.signTx(ctx, tx, zetaHeight, nonce)
		if err != nil {
			return models.TxnMetaData{}, "", errors.Wrap(err, "unable to sign tx")
		}

		return tx, sigBase64, nil
	}, nil
}

// buildWithdrawal builds unsigned withdrawal transaction using CCTX and Sui RPC
// https://github.com/zeta-chain/protocol-contracts-sui/blob/0245ad3a2eb4001381625070fd76c87c165589b2/sources/gateway.move#L117
func (s *Signer) buildWithdrawal(ctx context.Context, cctx *cctypes.CrossChainTx) (tx models.TxnMetaData, err error) {
	params := cctx.GetCurrentOutboundParam()

	coinType := ""

	// Basic common-sense validation & coin-type determination
	switch {
	case params.ReceiverChainId != s.Chain().ChainId:
		return tx, errors.Errorf("invalid receiver chain id %d", params.ReceiverChainId)
	case cctx.ProtocolContractVersion != cctypes.ProtocolContractVersion_V2:
		return tx, errors.Errorf("invalid protocol version %q", cctx.ProtocolContractVersion)
	case cctx.InboundParams == nil:
		return tx, errors.New("inbound params are nil")
	case params.CoinType == coin.CoinType_Gas:
		coinType = string(zetasui.SUI)
	case params.CoinType == coin.CoinType_ERC20:
		// NOTE: 0x prefix is required for coin type other than SUI
		coinType = "0x" + cctx.InboundParams.Asset
	default:
		return tx, errors.Errorf("unsupported coin type %q", params.CoinType.String())
	}

	// Gas budget is gas limit * gas price
	gasPrice, err := strconv.ParseUint(params.GasPrice, 10, 64)
	if err != nil {
		return tx, errors.Wrap(err, "unable to parse gas price")
	}
	gasBudget := gasPrice * params.CallOptions.GasLimit

	// Retrieve withdraw cap ID
	withdrawCapID, err := s.getWithdrawCapIDCached(ctx)
	if err != nil {
		return tx, errors.Wrap(err, "unable to get withdraw cap ID")
	}

	// Retrieve message context ID
	// TODO: https://github.com/zeta-chain/node/issues/4066
	// bring back this query after re-enabling authenticated call
	msgContextID := ""
	// msgContextID, err := s.getMessageContextIDCached(ctx)
	// if err != nil {
	// 	return tx, errors.Wrap(err, "unable to get message context ID")
	// }

	// build tx depending on the type of transaction
	if cctx.IsWithdrawAndCall() {
		return s.buildWithdrawAndCallTx(
			ctx,
			cctx,
			coinType,
			gasBudget,
			withdrawCapID,
			msgContextID,
			cctx.RelayedMessage,
		)
	}

	return s.buildWithdrawTx(ctx, params, coinType, gasBudget, withdrawCapID)
}

// buildWithdrawTx builds unsigned withdraw transaction
func (s *Signer) buildWithdrawTx(
	ctx context.Context,
	params *cctypes.OutboundParams,
	coinType string,
	gasBudget uint64,
	withdrawCapID string,
) (models.TxnMetaData, error) {
	var (
		nonce        = strconv.FormatUint(params.TssNonce, 10)
		recipient    = params.Receiver
		amount       = params.Amount.String()
		gasBudgetStr = strconv.FormatUint(gasBudget, 10)
	)

	req := models.MoveCallRequest{
		Signer:          s.TSS().PubKey().AddressSui(),
		PackageObjectId: s.gateway.PackageID(),
		Module:          zetasui.GatewayModule,
		Function:        funcWithdraw,
		TypeArguments:   []any{coinType},
		Arguments:       []any{s.gateway.ObjectID(), amount, nonce, recipient, gasBudgetStr, withdrawCapID},
		GasBudget:       gasBudgetStr,
	}

	return s.client.MoveCall(ctx, req)
}

// buildWithdrawAndCallTx builds unsigned withdrawAndCall
// a withdrawAndCall is a PTB transaction that contains a withdraw_impl call and a on_call call
func (s *Signer) buildWithdrawAndCallTx(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
	coinType string,
	gasBudget uint64,
	withdrawCapID string,
	msgContextID string,
	payloadHex string,
) (models.TxnMetaData, error) {
	params := cctx.GetCurrentOutboundParam()

	// decode and parse the payload into object IDs and on_call arguments
	payloadBytes, err := hex.DecodeString(payloadHex)
	if err != nil {
		return models.TxnMetaData{}, errors.Wrap(err, "unable to decode payload hex bytes")
	}

	var cp zetasui.CallPayload
	if err := cp.UnpackABI(payloadBytes); err != nil {
		return models.TxnMetaData{}, errors.Wrap(err, "unable to parse withdrawAndCall payload")
	}

	// get all needed object references
	wacRefs, err := s.getWithdrawAndCallObjectRefs(ctx, withdrawCapID, msgContextID, cp.ObjectIDs, gasBudget)
	if err != nil {
		return models.TxnMetaData{}, errors.Wrap(err, "unable to get object references")
	}

	// all PTB arguments
	args := withdrawAndCallPTBArgs{
		withdrawAndCallObjRefs: wacRefs,
		coinType:               coinType,
		amount:                 params.Amount.Uint64(),
		nonce:                  params.TssNonce,
		gasBudget:              gasBudget,
		sender:                 ethcommon.HexToAddress(cctx.InboundParams.Sender).Hex(),
		target:                 params.Receiver,
		payload:                cp,
	}

	// print PTB transaction parameters
	s.Logger().Std.Info().
		Str(logs.FieldMethod, "buildWithdrawAndCallTx").
		Uint64(logs.FieldNonce, args.nonce).
		Str(logs.FieldCoinType, args.coinType).
		Uint64("tx_amount", args.amount).
		Str("tx_sender", args.sender).
		Str("tx_target", args.target).
		Uint64("tx_gas_budget", args.gasBudget).
		Strs("tx_type_args", args.payload.TypeArgs).
		Strs("tx_object_ids", args.payload.ObjectIDs).
		Hex("tx_payload", args.payload.Message).
		Int("tx_sui_coins", len(args.withdrawAndCallObjRefs.suiCoins)).
		Msg("calling withdrawAndCallPTB")

	// build the PTB transaction
	return s.withdrawAndCallPTB(args)
}

// createCancelTxBuilder creates a cancel tx builder for given CCTX
// The tx cancellation is done by calling the 'increase_nonce' function on the gateway
// The goal or a "builder" instead of regular TxMetaData is to
// delay the 'MoveCall' to the last moment to avoid gateway object version mismatch
func (s *Signer) createCancelTxBuilder(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
	zetaHeight uint64,
) (txBuilder, error) {
	var (
		params = cctx.GetCurrentOutboundParam()
		nonce  = strconv.FormatUint(params.TssNonce, 10)
	)

	// get gas budget from CCTX
	gasBudget, err := getCancelTxGasBudget(params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get gas budget")
	}

	// retrieve withdraw cap ID
	withdrawCapID, err := s.getWithdrawCapIDCached(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get withdraw cap ID")
	}

	req := models.MoveCallRequest{
		Signer:          s.TSS().PubKey().AddressSui(),
		PackageObjectId: s.gateway.PackageID(),
		Module:          zetasui.GatewayModule,
		Function:        funcIncreaseNonce,
		TypeArguments:   []any{},
		Arguments:       []any{s.gateway.ObjectID(), nonce, withdrawCapID},
		GasBudget:       gasBudget,
	}

	return func(ctx context.Context) (models.TxnMetaData, string, error) {
		tx, err := s.client.MoveCall(ctx, req)
		if err != nil {
			return models.TxnMetaData{}, "", errors.Wrap(err, "unable to build cancel tx")
		}

		sigBase64, err := s.signTx(ctx, tx, zetaHeight, params.TssNonce)
		if err != nil {
			return models.TxnMetaData{}, "", errors.Wrap(err, "unable to sign cancel tx")
		}

		return tx, sigBase64, nil
	}, nil
}

// broadcastWithdrawalWithFallback broadcasts withdrawal tx to Sui network.
// If the tx execution is rejected, the cancel tx will be used and broadcasted (if provided).
func (s *Signer) broadcastWithdrawalWithFallback(
	ctx context.Context,
	withdrawTxBuilder, cancelTxBuilder txBuilder,
) (string, error) {
	logger := zerolog.Ctx(ctx).With().Str(logs.FieldMethod, "broadcastWithCancelTx").Logger()

	// should not happen
	if withdrawTxBuilder == nil || cancelTxBuilder == nil {
		return "", errors.New("withdrawal tx builder or cancel tx builder is nil")
	}

	tx, sig, err := withdrawTxBuilder(ctx)

	// we should cancel withdrawAndCall if user provided objects are not shared or immutable
	switch {
	case errors.Is(err, zetasui.ErrObjectOwnership):
		logger.Info().Err(err).Msg("cancelling tx due to wrong object ownership")
		return s.broadcastCancelTx(ctx, cancelTxBuilder)
	case err != nil:
		return "", errors.Wrap(err, "unable to build withdraw tx")
	}

	req := models.SuiExecuteTransactionBlockRequest{
		TxBytes:   tx.TxBytes,
		Signature: []string{sig},
		// we need to wait for the effects to be available and then look into
		// the error code to decide whether to cancel the tx or not
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: true,
		},
		RequestType: "WaitForEffectsCert",
	}

	// broadcast tx
	// Note: this is the place where the gateway object version mismatch error happens
	res, err := s.client.SuiExecuteTransactionBlock(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "unable to execute tx block")
	}

	// tx succeeded, return the digest
	if res.Effects.Status.Status == client.TxStatusSuccess {
		logger.Info().Str(logs.FieldTx, res.Digest).Msg("executed sui tx block successfully")
		return res.Digest, nil
	}

	// check if the error is a retryable MoveAbort
	// if it is, skip the cancel tx and let the scheduler retry the outbound
	// TODO: https://github.com/zeta-chain/node/issues/4066
	// use IsRetryableExecutionError instead after re-enabling authenticated call
	isRetryable, err := zetasui.IsRetryableExecutionErrorLegacy(res.Effects.Status.Error)
	switch {
	case err != nil:
		return "", errors.Wrapf(err, "unable to check tx execution status error: %s", res.Effects.Status.Error)
	case isRetryable:
		return "", fmt.Errorf("tx execution status failed, retry later: %s", res.Effects.Status.Error)
	default:
		// cancel tx if the tx execution failed for all other reasons
		// wait for gateway object version bump to avoid version mismatch
		time.Sleep(2 * time.Second)
		logger.Info().Any("Err", res.Effects.Status.Error).Msg("cancelling tx")
	}

	return s.broadcastCancelTx(ctx, cancelTxBuilder)
}

// broadcastCancelTx broadcasts a cancel tx and returns the tx digest
func (s *Signer) broadcastCancelTx(ctx context.Context, cancelTxBuilder txBuilder) (string, error) {
	logger := zerolog.Ctx(ctx).With().Str(logs.FieldMethod, "broadcastCancelTx").Logger()

	// build cancel tx
	txCancel, sigCancel, err := cancelTxBuilder(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to build cancel tx")
	}

	// create tx request
	reqCancel := models.SuiExecuteTransactionBlockRequest{
		TxBytes:   txCancel.TxBytes,
		Signature: []string{sigCancel},
	}

	// broadcast cancel tx
	res, err := s.client.SuiExecuteTransactionBlock(ctx, reqCancel)
	if err != nil {
		return "", errors.Wrap(err, "unable to execute cancel tx block")
	}
	logger.Info().Str(logs.FieldTx, res.Digest).Msg("executed sui cancel tx block")

	return res.Digest, nil
}

// getCancelTxGasBudget returns gas budget for a cancel tx
func getCancelTxGasBudget(params *cctypes.OutboundParams) (string, error) {
	gasPrice, err := strconv.ParseUint(params.GasPrice, 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse gas price")
	}

	// If it is a cancel tx, we need to use the bigger one
	// because the cancelled tx may be caused by insufficient gas
	gasBudget := max(gasPrice*params.CallOptions.GasLimit, minGasBudgetCancelTx)

	return strconv.FormatUint(gasBudget, 10), nil
}

// getGatewayNonce reads the nonce of the gateway object
func (s *Signer) getGatewayNonce(ctx context.Context) (uint64, error) {
	data, err := s.client.GetObjectParsedData(ctx, s.gateway.ObjectID())
	if err != nil {
		return 0, errors.Wrap(err, "unable to get parsed data of gateway object")
	}

	return zetasui.ParseGatewayNonce(data)
}
