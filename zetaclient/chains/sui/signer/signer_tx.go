package signer

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
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
		coinType = string(sui.SUI)
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
	gasBudget := strconv.FormatUint(gasPrice*params.CallOptions.GasLimit, 10)

	// Retrieve withdraw cap ID
	withdrawCapIDStr, err := s.getWithdrawCapIDCached(ctx)
	if err != nil {
		return tx, errors.Wrap(err, "unable to get withdraw cap ID")
	}

	// build tx depending on the type of transaction
	if cctx.IsWithdrawAndCall() {
		return s.buildWithdrawAndCallTx(ctx, params, coinType, gasBudget, withdrawCapIDStr, cctx.RelayedMessage)
	}

	return s.buildWithdrawTx(ctx, params, coinType, gasBudget, withdrawCapIDStr)
}

// buildWithdrawTx builds unsigned withdraw transaction
func (s *Signer) buildWithdrawTx(
	ctx context.Context,
	params *cctypes.OutboundParams,
	coinType, gasBudget, withdrawCapID string,
) (models.TxnMetaData, error) {
	var (
		nonce     = strconv.FormatUint(params.TssNonce, 10)
		recipient = params.Receiver
		amount    = params.Amount.String()
	)

	req := models.MoveCallRequest{
		Signer:          s.TSS().PubKey().AddressSui(),
		PackageObjectId: s.gateway.PackageID(),
		Module:          s.gateway.Module(),
		Function:        funcWithdraw,
		TypeArguments:   []any{coinType},
		Arguments:       []any{s.gateway.ObjectID(), amount, nonce, recipient, gasBudget, withdrawCapID},
		GasBudget:       gasBudget,
	}

	return s.client.MoveCall(ctx, req)
}

// buildWithdrawAndCallTx builds unsigned withdrawAndCall
// a withdrawAndCall is a PTB transaction that contains a withdraw_impl call and a on_call call
func (s *Signer) buildWithdrawAndCallTx(
	ctx context.Context,
	params *cctypes.OutboundParams,
	coinType,
	gasBudget,
	withdrawCapIDStr,
	payload string,
) (models.TxnMetaData, error) {
	// decode and parse the payload to object the on_call arguments
	payloadBytes, err := hex.DecodeString(payload)
	if err != nil {
		return models.TxnMetaData{}, errors.Wrap(err, "unable to decode payload hex bytes")
	}

	var cp sui.CallPayload
	if err := cp.UnpackABI(payloadBytes); err != nil {
		return models.TxnMetaData{}, errors.Wrap(err, "unable to parse withdrawAndCall payload")
	}

	// get latest TSS SUI coin object ref for gas payment
	suiCoinObjRef, err := s.client.GetSuiCoinObjectRef(ctx, s.TSS().PubKey().AddressSui())
	if err != nil {
		return models.TxnMetaData{}, errors.Wrap(err, "unable to get TSS SUI coin object")
	}

	// get all other object references: [gateway, withdrawCap, onCallObjects]
	gatewayObjRef, withdrawCapObjRef, onCallObjectRefs, err := getWithdrawAndCallObjectRefs(
		ctx,
		s.client,
		s.gateway.ObjectID(),
		withdrawCapIDStr,
		cp.ObjectIDs,
	)
	if err != nil {
		return models.TxnMetaData{}, errors.Wrap(err, "unable to get object references")
	}

	// print PTB transaction parameters
	s.Logger().Std.Info().
		Str(logs.FieldMethod, "buildWithdrawAndCallTx").
		Str(logs.FieldCoinType, coinType).
		Str("amount", params.Amount.String()).
		Uint64(logs.FieldNonce, params.TssNonce).
		Str("receiver", params.Receiver).
		Str("gas_budget", gasBudget).
		Any("type_args", cp.TypeArgs).
		Any("object_ids", cp.ObjectIDs).
		Hex("message", cp.Message).
		Msg("calling withdrawAndCallPTB")

	// TODO: check all object IDs are share object here
	// https://github.com/zeta-chain/node/issues/3755

	// build the PTB transaction
	return withdrawAndCallPTB(
		s.TSS().PubKey().AddressSui(),
		s.gateway.PackageID(),
		s.gateway.Module(),
		gatewayObjRef,
		suiCoinObjRef,
		withdrawCapObjRef,
		onCallObjectRefs,
		coinType,
		params.Amount.String(),
		strconv.FormatUint(params.TssNonce, 10),
		gasBudget,
		params.Receiver,
		cp,
	)
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
		Module:          s.gateway.Module(),
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
	if err != nil {
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
		logger.Info().Str(logs.FieldTx, res.Digest).Msg("Executed sui tx block successfully")
		return res.Digest, nil
	}

	// check if the error is a retryable MoveAbort
	// if it is, skip the cancel tx and let the scheduler retry the outbound
	isRetryable, err := sui.IsRetryableExecutionError(res.Effects.Status.Error)
	switch {
	case err != nil:
		return "", errors.Wrapf(err, "unable to check tx execution status error code: %s", res.Effects.Status.Error)
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
	logger.Info().Str(logs.FieldTx, res.Digest).Msg("Executed sui cancel tx block")

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
