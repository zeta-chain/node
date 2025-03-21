package signer

import (
	"context"
	"strconv"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/logs"
)

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
// Note: this function is a fake implementation for testing the fallback mechanism
// TODO: replace it with a real implementation
func (s *Signer) buildExecute(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
) (models.TxnMetaData, *models.TxnMetaData, error) {
	var (
		err        error
		tx         models.TxnMetaData
		txFallback models.TxnMetaData
	)

	// fake out the execute tx with a withdrawal
	tx, err = s.buildWithdraw(ctx, cctx)
	if err != nil {
		return tx, nil, errors.Wrap(err, "unable to build withdraw tx")
	}

	txFallback, err = s.buildIncreaseNonce(ctx, cctx)
	if err != nil {
		return tx, nil, errors.Wrap(err, "unable to build withdraw tx")
	}

	return tx, &txFallback, nil
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

// broadcastWithFallback attaches signature to tx and broadcasts it to Sui network. Returns tx digest.
// If the tx execution fails, the fallback tx will be used and broadcasted if provided.
func (s *Signer) broadcastWithFallback(
	ctx context.Context,
	tx models.TxnMetaData,
	txFallback *models.TxnMetaData,
	sig string,
	sigFallback *string,
) (string, error) {
	req := models.SuiExecuteTransactionBlockRequest{
		TxBytes:   tx.TxBytes,
		Signature: []string{sig},
	}

	var reqFallback *models.SuiExecuteTransactionBlockRequest
	if txFallback != nil && sigFallback != nil {
		reqFallback = &models.SuiExecuteTransactionBlockRequest{
			TxBytes:   txFallback.TxBytes,
			Signature: []string{*sigFallback},
		}
	}

	// broadcast original tx
	res, err := s.client.SuiExecuteTransactionBlock(ctx, req)
	if err == nil {
		s.Logger().Std.Info().Str(logs.FieldTx, res.Digest).Msg("Broadcasted Sui tx successfully")
		return res.Digest, nil
	}

	// decide whether to broadcast fallback tx
	shouldFallback := reqFallback != nil && strings.Contains(err.Error(), "some specific error")
	if !shouldFallback {
		return "", errors.Wrap(err, "unable to execute tx block")
	}

	// broadcast fallback tx
	res, err = s.client.SuiExecuteTransactionBlock(ctx, *reqFallback)
	if err != nil {
		return "", errors.Wrap(err, "unable to execute fallback tx block")
	}
	s.Logger().Std.Info().Str(logs.FieldTx, res.Digest).Msg("Broadcasted Sui fallback tx successfully")

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
