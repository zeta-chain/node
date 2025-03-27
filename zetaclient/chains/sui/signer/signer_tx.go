package signer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// txBuilder is a function that returns the tx and the signature
type txBuilder func() (models.TxnMetaData, string, error)

const (
	funcWithdraw      = "withdraw"
	funcExecute       = "execute"
	funcIncreaseNonce = "increase_nonce"
)

// buildWithdraw builds unsigned 'withdraw' transaction using CCTX and Sui RPC
// https://github.com/zeta-chain/protocol-contracts-sui/blob/0245ad3a2eb4001381625070fd76c87c165589b2/sources/gateway.move#L117
func (s *Signer) buildWithdraw(ctx context.Context, cctx *cctypes.CrossChainTx) (tx models.TxnMetaData, err error) {
	var (
		params    = cctx.GetCurrentOutboundParam()
		nonce     = strconv.FormatUint(params.TssNonce, 10)
		recipient = params.Receiver
		amount    = params.Amount.String()
	)

	// perform basic validation on CCTX
	coinType, err := validateCCTX(s.Chain().ChainId, cctx)
	if err != nil {
		return tx, errors.Wrap(err, "CCTX validation failed")
	}

	// get gas budget from CCTX
	gasBudget, err := gasBudgetFromCCTX(cctx)
	if err != nil {
		return tx, errors.Wrap(err, "unable to get gas budget")
	}

	withdrawCapID, err := s.getWithdrawCapIDCached(ctx)
	if err != nil {
		return tx, errors.Wrap(err, "unable to get withdraw cap ID")
	}

	// build withdraw tx
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

// buildExecute builds unsigned 'execute' transaction using CCTX and Sui RPC
//
// TODO: replace it with a real implementation
// this function is a fake implementation for testing the tx cancellation
func (s *Signer) buildExecute(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
) (tx models.TxnMetaData, err error) {
	var (
		params    = cctx.GetCurrentOutboundParam()
		nonce     = strconv.FormatUint(params.TssNonce, 10)
		recipient = params.Receiver
		amount    = params.Amount.String()
	)

	// perform basic validation on CCTX
	coinType, err := validateCCTX(s.Chain().ChainId, cctx)
	if err != nil {
		return tx, errors.Wrap(err, "CCTX validation failed")
	}

	// use a fake coin type to fail the withdraw tx
	// this will force the withdraw to be cancelled
	if cctx.InboundParams.CoinType == coin.CoinType_Gas {
		coinType = coinType[0:64] + "::fake::FAKE"
	} else {
		coinType = coinType[0:66] + "::fake::FAKE"
	}

	// get gas budget from CCTX
	gasBudget, err := gasBudgetFromCCTX(cctx)
	if err != nil {
		return tx, errors.Wrap(err, "unable to get gas budget")
	}

	withdrawCapID, err := s.getWithdrawCapIDCached(ctx)
	if err != nil {
		return tx, errors.Wrap(err, "unable to get withdraw cap ID")
	}

	// build a withdraw tx as a fake 'execute' because we don't have a real one yet
	req := models.MoveCallRequest{
		Signer:          s.TSS().PubKey().AddressSui(),
		PackageObjectId: s.gateway.PackageID(),
		Module:          s.gateway.Module(),
		Function:        funcWithdraw, // TODO: change to funcExecute
		TypeArguments:   []any{coinType},
		Arguments:       []any{s.gateway.ObjectID(), amount, nonce, recipient, gasBudget, withdrawCapID},
		GasBudget:       gasBudget,
	}

	return s.client.MoveCall(ctx, req)
}

// buildIncreaseNonce builds unsigned 'increase_nonce' transaction using CCTX and Sui RPC
func (s *Signer) buildIncreaseNonce(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
) (tx models.TxnMetaData, err error) {
	nonce := strconv.FormatUint(cctx.GetCurrentOutboundParam().TssNonce, 10)

	// perform basic validation on CCTX
	_, err = validateCCTX(s.Chain().ChainId, cctx)
	if err != nil {
		return tx, errors.Wrap(err, "CCTX validation failed")
	}

	// get gas budget from CCTX
	gasBudget, err := gasBudgetFromCCTX(cctx)
	if err != nil {
		return tx, errors.Wrap(err, "unable to get gas budget")
	}

	withdrawCapID, err := s.getWithdrawCapIDCached(ctx)
	if err != nil {
		return tx, errors.Wrap(err, "unable to get withdraw cap ID")
	}

	// build increase nonce tx
	req := models.MoveCallRequest{
		Signer:          s.TSS().PubKey().AddressSui(),
		PackageObjectId: s.gateway.PackageID(),
		Module:          s.gateway.Module(),
		Function:        funcIncreaseNonce,
		TypeArguments:   []any{},
		Arguments:       []any{s.gateway.ObjectID(), nonce, withdrawCapID},
		GasBudget:       gasBudget,
	}

	return s.client.MoveCall(ctx, req)
}

// buildExecuteWithCancelTxBuilder builds both unsigned 'execute' and a cancel tx
func (s *Signer) buildExecuteWithCancelTxBuilder(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
	zetaHeight uint64,
) (models.TxnMetaData, txBuilder, error) {
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	tx, err := s.buildExecute(ctx, cctx)
	if err != nil {
		return tx, nil, errors.Wrap(err, "unable to build execute tx")
	}

	// tx builder for cancel tx
	cancelTxBuilder := func() (models.TxnMetaData, string, error) {
		txCancel, err := s.buildIncreaseNonce(ctx, cctx)
		if err != nil {
			return models.TxnMetaData{}, "", errors.Wrap(err, "unable to build cancel tx")
		}

		sigCancel, err := s.signTx(ctx, txCancel, zetaHeight, nonce)
		if err != nil {
			return models.TxnMetaData{}, "", errors.Wrap(err, "unable to sign cancel tx")
		}

		return txCancel, sigCancel, nil
	}

	return tx, cancelTxBuilder, nil
}

// broadcastWithCancelTx attaches signature to tx and broadcasts it to Sui network. Returns tx digest.
// If the tx execution is rejected, the cancel tx will be used and broadcasted if provided.
func (s *Signer) broadcastWithCancelTx(
	ctx context.Context,
	sig string,
	tx models.TxnMetaData,
	cancelTxBuilder txBuilder,
) (string, error) {
	logger := zerolog.Ctx(ctx).With().Str(logs.FieldMethod, "broadcastWithCancelTx").Logger()

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
	res, err := s.client.SuiExecuteTransactionBlock(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "unable to execute tx block")
	}

	// tx succeeded, return the digest
	if res.Effects.Status.Status == "success" {
		logger.Info().Str(logs.FieldTx, res.Digest).Msg("Executed sui tx block successfully")
		return res.Digest, nil
	}

	// tx failed, return error if no cancel tx provided
	if cancelTxBuilder == nil {
		return "", fmt.Errorf("tx execution status failed: %s", res.Effects.Status.Error)
	}

	// check if the error is a retryable MoveAbort
	// if it is, skip the cancel tx and let the scheduler retry the outbound
	isRetryableAbort, err := sui.IsRetryableMoveAbort(res.Effects.Status.Error)
	switch {
	case err != nil:
		return "", errors.Wrapf(err, "unable to check tx execution status error code: %s", res.Effects.Status.Error)
	case isRetryableAbort:
		return "", fmt.Errorf("tx execution status failed, retry later: %s", res.Effects.Status.Error)
	default:
		// cancel tx if the tx execution failed for all other reasons
		// wait for gateway object version bump to avoid version mismatch
		time.Sleep(2 * time.Second)
		logger.Info().Any("Err", res.Effects.Status.Error).Msg("cancelling tx")
	}

	// build cancel tx
	txCancel, sigCancel, err := cancelTxBuilder()
	if err != nil {
		return "", errors.Wrap(err, "unable to build cancel tx")
	}
	reqCancel := models.SuiExecuteTransactionBlockRequest{
		TxBytes:   txCancel.TxBytes,
		Signature: []string{sigCancel},
	}

	// broadcast cancel tx
	res, err = s.client.SuiExecuteTransactionBlock(ctx, reqCancel)
	if err != nil {
		return "", errors.Wrap(err, "unable to execute cancel tx block")
	}
	logger.Info().Str(logs.FieldTx, res.Digest).Msg("Executed sui cancel tx block")

	return res.Digest, nil
}

// validateCCTX performs basic common-sense validation on CCTX and determines coin type
func validateCCTX(signerChainID int64, cctx *cctypes.CrossChainTx) (coinType string, err error) {
	params := cctx.GetCurrentOutboundParam()

	switch {
	case params.ReceiverChainId != signerChainID:
		return "", errors.Errorf("invalid receiver chain id %d", params.ReceiverChainId)
	case cctx.ProtocolContractVersion != cctypes.ProtocolContractVersion_V2:
		return "", errors.Errorf("invalid protocol version %q", cctx.ProtocolContractVersion)
	case cctx.InboundParams == nil:
		return "", errors.New("inbound params are nil")
	case cctx.InboundParams.CoinType == coin.CoinType_Gas:
		coinType = string(sui.SUI)
	case cctx.InboundParams.CoinType == coin.CoinType_ERC20:
		// NOTE: 0x prefix is required for coin type other than SUI
		coinType = "0x" + cctx.InboundParams.Asset
	default:
		return "", errors.Errorf("unsupported coin type %q", cctx.InboundParams.CoinType.String())
	}

	return coinType, nil
}

// gasBudgetFromCCTX returns gas budget from CCTX
func gasBudgetFromCCTX(cctx *cctypes.CrossChainTx) (string, error) {
	params := cctx.GetCurrentOutboundParam()

	// Gas budget is gas limit * gas price
	gasPrice, err := strconv.ParseUint(params.GasPrice, 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse gas price")
	}

	return strconv.FormatUint(gasPrice*params.CallOptions.GasLimit, 10), nil
}
