package signer

import (
	"context"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
)

const funcWithdraw = "withdraw"

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
	case cctx.IsWithdrawAndCall():
		return tx, errors.New("withdrawAndCall is not supported yet")
	case params.CoinType == coin.CoinType_Gas:
		coinType = string(sui.SUI)
	case params.CoinType == coin.CoinType_ERC20:
		// NOTE: 0x prefix is required for coin type other than SUI
		coinType = "0x" + cctx.InboundParams.Asset
	default:
		return tx, errors.Errorf("unsupported coin type %q", params.CoinType.String())
	}

	var (
		nonce     = strconv.FormatUint(params.TssNonce, 10)
		recipient = params.Receiver
		amount    = params.Amount.String()
	)

	// Gas budget is gas limit * gas price
	gasPrice, err := strconv.ParseUint(params.GasPrice, 10, 64)
	if err != nil {
		return tx, errors.Wrap(err, "unable to parse gas price")
	}
	gasBudget := strconv.FormatUint(gasPrice*params.CallOptions.GasLimit, 10)

	withdrawCapID, err := s.getWithdrawCapIDCached(ctx)
	if err != nil {
		return tx, errors.Wrap(err, "unable to get withdraw cap ID")
	}

	// PATCH V29: there is currently no validation of the withdraw address, if the address is not valid then ZetaClient will continuously try to broadcast the tx, blocking outbound
	// This patch will redirect the funds to the TSS to make the withdraw succeeding
	// invalid recipient are considered like a donation to the TSS
	// TODO: add validation of the withdraw address to prevent this issue
	// https://github.com/zeta-chain/node/issues/3798
	// https://github.com/zeta-chain/node/issues/3799
	if sui.ValidAddress(recipient) != nil {
		s.Logger().Std.Warn().Str("recipient", recipient).Msg("Invalid recipient address, redirecting to TSS")
		recipient = s.TSS().PubKey().AddressSui()
	}

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
