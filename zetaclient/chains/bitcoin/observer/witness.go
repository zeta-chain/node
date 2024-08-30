package observer

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

// GetBtcEventWithWitness either returns a valid BTCInboundEvent or nil.
// This method supports data with more than 80 bytes by scanning the witness for possible presence of a tapscript.
// It will first prioritize OP_RETURN over tapscript.
func GetBtcEventWithWitness(
	client interfaces.BTCRPCClient,
	tx btcjson.TxRawResult,
	tssAddress string,
	blockNumber uint64,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
	depositorFee float64,
) (*BTCInboundEvent, error) {
	if len(tx.Vout) < 1 {
		logger.Debug().Msgf("no output %s", tx.Txid)
		return nil, nil
	}
	if len(tx.Vin) == 0 {
		logger.Debug().Msgf("no input found for inbound: %s", tx.Txid)
		return nil, nil
	}

	if err := isValidRecipient(tx.Vout[0].ScriptPubKey.Hex, tssAddress, netParams); err != nil {
		logger.Debug().Msgf("irrelevant recipient %s for tx %s, err: %s", tx.Vout[0].ScriptPubKey.Hex, tx.Txid, err)
		return nil, nil
	}

	isAmountValid, amount := isValidAmount(tx.Vout[0].Value, depositorFee)
	if !isAmountValid {
		logger.Info().
			Msgf("GetBtcEventWithWitness: btc deposit amount %v in txid %s is less than depositor fee %v", tx.Vout[0].Value, tx.Txid, depositorFee)
		return nil, nil
	}

	// Try to extract the memo from the BTC txn. First try to extract from OP_RETURN
	// if not found then try to extract from inscription. Return nil if the above two
	// cannot find the memo.
	var memo []byte
	if candidate := tryExtractOpRet(tx, logger); candidate != nil {
		memo = candidate
		logger.Debug().
			Msgf("GetBtcEventWithWitness: found OP_RETURN memo %s in tx %s", hex.EncodeToString(memo), tx.Txid)
	} else if candidate = tryExtractInscription(tx, logger); candidate != nil {
		memo = candidate
		logger.Debug().Msgf("GetBtcEventWithWitness: found inscription memo %s in tx %s", hex.EncodeToString(memo), tx.Txid)
	} else {
		return nil, errors.Errorf("error getting memo for inbound: %s", tx.Txid)
	}

	// event found, get sender address
	fromAddress, err := GetSenderAddressByVin(client, tx.Vin[0], netParams)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting sender address for inbound: %s", tx.Txid)
	}

	return &BTCInboundEvent{
		FromAddress: fromAddress,
		ToAddress:   tssAddress,
		Value:       amount,
		MemoBytes:   memo,
		BlockNumber: blockNumber,
		TxHash:      tx.Txid,
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
	if length >= 2 && len(lastElement) > 0 && lastElement[0] == 0x50 {
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

// / Try to extract the memo from the OP_RETURN
func tryExtractOpRet(tx btcjson.TxRawResult, logger zerolog.Logger) []byte {
	if len(tx.Vout) < 2 {
		logger.Debug().Msgf("txn %s has fewer than 2 outputs, not target OP_RETURN txn", tx.Txid)
		return nil
	}

	memo, found, err := bitcoin.DecodeOpReturnMemo(tx.Vout[1].ScriptPubKey.Hex, tx.Txid)
	if err != nil {
		logger.Error().Err(err).Msgf("tryExtractOpRet: error decoding OP_RETURN memo: %s", tx.Vout[1].ScriptPubKey.Hex)
		return nil
	}

	if found {
		return memo
	}
	return nil
}

// / Try to extract the memo from inscription
func tryExtractInscription(tx btcjson.TxRawResult, logger zerolog.Logger) []byte {
	for i, input := range tx.Vin {
		script := ParseScriptFromWitness(input.Witness, logger)
		if script == nil {
			continue
		}

		logger.Debug().Msgf("potential witness script, tx %s, input idx %d", tx.Txid, i)

		memo, found, err := bitcoin.DecodeScript(script)
		if err != nil || !found {
			logger.Debug().Msgf("invalid witness script, tx %s, input idx %d", tx.Txid, i)
			continue
		}

		logger.Debug().Msgf("found memo in inscription, tx %s, input idx %d", tx.Txid, i)
		return memo
	}

	return nil
}

func isValidAmount(
	incoming float64,
	minimal float64,
) (bool, float64) {
	if incoming < minimal {
		return false, 0
	}
	return true, incoming - minimal
}

func isValidRecipient(
	script string,
	tssAddress string,
	netParams *chaincfg.Params,
) error {
	receiver, err := bitcoin.DecodeScriptP2WPKH(script, netParams)
	if err != nil {
		return fmt.Errorf("invalid p2wpkh script detected, %s", err)
	}

	// skip irrelevant tx to us
	if receiver != tssAddress {
		return fmt.Errorf("irrelevant recipient, %s", receiver)
	}

	return nil
}
