package signer

import (
	"context"
	"encoding/hex"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
)

const funcWithdraw = "withdraw"

// TODO: use these functions in PTB building
// https://github.com/zeta-chain/node/issues/3741
//const funcWithdrawImpl = "withdraw_impl"
//const funcOnCall = "on_call"

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
	withdrawCapID, err := s.getWithdrawCapIDCached(ctx)
	if err != nil {
		return tx, errors.Wrap(err, "unable to get withdraw cap ID")
	}

	// build tx depending on the type of transaction
	if cctx.IsWithdrawAndCall() {
		return s.buildWithdrawAndCallTx(ctx, params, coinType, gasBudget, withdrawCapID, cctx.RelayedMessage)
	}
	return s.buildWithdrawTx(ctx, params, coinType, gasBudget, withdrawCapID)
}

// buildWithdrawTx builds unsigned withdraw transaction
func (s *Signer) buildWithdrawTx(
	ctx context.Context,
	params *cctypes.OutboundParams,
	coinType,
	gasBudget,
	withdrawCapID string,
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
	withdrawCapID,
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

	// Note: logs not formatted in standard, it's a temporary log
	s.Logger().Std.Info().Msgf(
		"WithdrawAndCall called with type arguments %v, object IDs %v, message %v",
		cp.TypeArgs,
		cp.ObjectIDs,
		cp.Message,
	)

	// keep lint quiet without using _ in params
	_ = ctx
	_ = params
	_ = coinType
	_ = gasBudget
	_ = withdrawCapID

	// TODO: check all object IDs are share object here
	// https://github.com/zeta-chain/node/issues/3755

	// TODO: build PTB here
	// https://github.com/zeta-chain/node/issues/3741

	return models.TxnMetaData{}, errors.New("not implemented")
}

// broadcast attaches signature to tx and broadcasts it to Sui network. Returns tx digest.
func (s *Signer) broadcast(ctx context.Context, tx models.TxnMetaData, sig [65]byte) (string, error) {
	sigBase64, err := sui.SerializeSignatureECDSA(sig, s.TSS().PubKey().AsECDSA())
	if err != nil {
		return "", errors.Wrap(err, "unable to serialize signature")
	}

	req := models.SuiExecuteTransactionBlockRequest{
		TxBytes:   tx.TxBytes,
		Signature: []string{sigBase64},
	}

	res, err := s.client.SuiExecuteTransactionBlock(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "unable to execute tx block")
	}

	return res.Digest, nil
}
