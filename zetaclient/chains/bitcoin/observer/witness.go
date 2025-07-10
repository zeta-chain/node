package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/logs"
)

const (
	// noMemoFound is a placeholder to indicates no memo is found in Bitcoin inbound
	noMemoFound = "no memo found"
)

// GetBtcEventWithWitness either returns a valid BTCInboundEvent or nil.
//
// This method supports two types of memo:
// 1. OP_RETURN based memo:
//   - the default memo type that can carry up to 80 bytes of data
//
// 2. Tapscript based memo:
//   - allow data with more than 80 bytes by scanning the witness for possible presence of a tapscript.
//
// Note:  OP_RETURN based memo is prioritized over tapscript memo if both are present.
func GetBtcEventWithWitness(
	ctx context.Context,
	rpc RPC,
	tx btcjson.TxRawResult,
	tssAddress string,
	blockNumber uint64,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
	feeCalculator common.DepositorFeeCalculator,
) (*BTCInboundEvent, error) {
	lf := map[string]any{
		logs.FieldMethod: "GetBtcEventWithWitness",
		logs.FieldTx:     tx.Txid,
	}

	if len(tx.Vout) < 1 {
		logger.Debug().Fields(lf).Msg("no output")
		return nil, nil
	}
	if len(tx.Vin) == 0 {
		logger.Debug().Fields(lf).Msg("no input found for inbound")
		return nil, nil
	}

	if err := isValidRecipient(tx.Vout[0].ScriptPubKey.Hex, tssAddress, netParams); err != nil {
		logger.Debug().Err(err).Fields(lf).Msgf("irrelevant recipient: %s", tx.Vout[0].ScriptPubKey.Hex)
		return nil, nil
	}

	// event found, get sender address
	fromAddress, err := rpc.GetTransactionInputSpender(ctx, tx.Vin[0].Txid, tx.Vin[0].Vout)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting sender address for inbound: %s", tx.Txid)
	}

	// skip this tx if one of the two conditions is met
	// 1. sender is empty, we don't know whom to refund if this tx gets reverted in zetacore
	// 2. the tx is an outbound (sender is TSS) and we should not process it as an inbound
	if fromAddress == "" || strings.EqualFold(fromAddress, tssAddress) {
		logger.Info().Fields(lf).Msgf("skipping transaction for sender: %s", fromAddress)
		return nil, nil
	}

	// calculate depositor fee
	depositorFee, err := feeCalculator(ctx, rpc, &tx, netParams)
	if err != nil {
		return nil, errors.Wrapf(err, "error calculating depositor fee for inbound %s", tx.Txid)
	}

	// deduct depositor fee
	// to allow developers to track failed deposit caused by insufficient depositor fee,
	// the error message will be forwarded to zetacore to register a failed CCTX
	status := types.InboundStatus_SUCCESS
	amount, err := DeductDepositorFee(tx.Vout[0].Value, depositorFee)
	if err != nil {
		amount = 0
		status = types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE
		logger.Error().Err(err).Fields(lf).Msgf("unable to deduct depositor fee")
	}

	// Try to extract the memo from the BTC txn. First try to extract from OP_RETURN
	// if not found then try to extract from inscription. If no memo is provided,
	// set the 'noMemoFound' placeholder to indicate the inbound requires a revert.
	var memo []byte
	if candidate := tryExtractOpRet(tx, logger); candidate != nil {
		memo = candidate
		logger.Debug().Fields(lf).Msgf("found OP_RETURN memo: %s", hex.EncodeToString(memo))
	} else if candidate = tryExtractInscription(tx, logger); candidate != nil {
		memo = candidate
		logger.Debug().Fields(lf).Msgf("found inscription memo: %s", hex.EncodeToString(memo))

		// override the sender address with the initiator of the inscription's commit tx
		if fromAddress, err = rpc.GetTransactionInitiator(ctx, tx.Vin[0].Txid); err != nil {
			return nil, errors.Wrap(err, "unable to get inscription initiator")
		}
	} else {
		memo = []byte(noMemoFound)
	}

	return &BTCInboundEvent{
		FromAddress:  fromAddress,
		ToAddress:    tssAddress,
		Value:        amount,
		DepositorFee: depositorFee,
		MemoBytes:    memo,
		BlockNumber:  blockNumber,
		TxHash:       tx.Txid,
		Status:       status,
	}, nil
}

// ParseScriptFromWitness attempts to parse the script from the witness data. Ideally it should be handled by
// bitcoin library, however, it's not found in existing library version. Replace this with actual library implementation
// if libraries are updated.
func ParseScriptFromWitness(witness []string, logger zerolog.Logger) []byte {
	length := len(witness)

	if length == 0 {
		return nil
	}

	lastElement, err := hex.DecodeString(witness[length-1])
	if err != nil {
		logger.Debug().Msgf("invalid witness element")
		return nil
	}

	// From BIP341:
	// If there are at least two witness elements, and the first byte of
	// the last element is 0x50, this last element is called annex a
	// and is removed from the witness stack.
	if length >= 2 && len(lastElement) > 0 && lastElement[0] == txscript.TaprootAnnexTag {
		// account for the extra item removed from the end
		witness = witness[:length-1]
	}

	if len(witness) < 2 {
		logger.Debug().Msgf("not script path spending detected, ignore")
		return nil
	}

	// only the script is the focus here, ignore checking control block or whatever else
	script, err := hex.DecodeString(witness[len(witness)-2])
	if err != nil {
		logger.Debug().Msgf("witness script cannot be decoded from hex, ignore")
		return nil
	}
	return script
}

// Try to extract the memo from the OP_RETURN
func tryExtractOpRet(tx btcjson.TxRawResult, logger zerolog.Logger) []byte {
	if len(tx.Vout) < 2 {
		logger.Debug().Msgf("txn %s has fewer than 2 outputs, not target OP_RETURN txn", tx.Txid)
		return nil
	}

	memo, found, err := common.DecodeOpReturnMemo(tx.Vout[1].ScriptPubKey.Hex)
	if err != nil {
		logger.Error().Err(err).Msgf("tryExtractOpRet: error decoding OP_RETURN memo: %s", tx.Vout[1].ScriptPubKey.Hex)
		return nil
	}

	if found {
		return memo
	}
	return nil
}

// Try to extract the memo from inscription
func tryExtractInscription(tx btcjson.TxRawResult, logger zerolog.Logger) []byte {
	for i, input := range tx.Vin {
		script := ParseScriptFromWitness(input.Witness, logger)
		if script == nil {
			continue
		}

		logger.Debug().Msgf("potential witness script, tx %s, input idx %d", tx.Txid, i)

		memo, found, err := common.DecodeScript(script)
		if err != nil || !found {
			logger.Debug().Msgf("invalid witness script, tx %s, input idx %d", tx.Txid, i)
			continue
		}

		logger.Debug().Msgf("found memo in inscription, tx %s, input idx %d", tx.Txid, i)
		return memo
	}

	return nil
}

// DeductDepositorFee returns the inbound amount after deducting the depositor fee.
// returns an error if the deposited amount is lower than the depositor fee.
func DeductDepositorFee(deposited, depositorFee float64) (float64, error) {
	if deposited < depositorFee {
		return 0, fmt.Errorf("deposited amount %v is less than depositor fee %v", deposited, depositorFee)
	}
	return deposited - depositorFee, nil
}

func isValidRecipient(
	script string,
	tssAddress string,
	netParams *chaincfg.Params,
) error {
	receiver, err := common.DecodeScriptP2WPKH(script, netParams)
	if err != nil {
		return fmt.Errorf("invalid p2wpkh script detected, %s", err)
	}

	// skip irrelevant tx to us
	if receiver != tssAddress {
		return fmt.Errorf("irrelevant recipient, %s", receiver)
	}

	return nil
}
