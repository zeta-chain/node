package sui

import (
	"strconv"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
)

// CancelTx represents data for a Sui cancelled tx nonce increase event
type CancelTx struct {
	Sender string
	Nonce  uint64
}

// TokenAmount returns hardcoded zero amount
// because nonce increase does not involve any token transfer
func (n CancelTx) TokenAmount() math.Uint {
	return math.NewUint(0)
}

// TxNonce returns the nonce of the cancelled tx
// Nonce in the event is the incremented value, so we need to subtract 1
// see: https://github.com/zeta-chain/protocol-contracts-sui/blob/e5a756e473da884dcbc59b574b387a7a365ac823/sources/gateway.move#L140
func (n CancelTx) TxNonce() uint64 {
	return n.Nonce - 1
}

func parseCancelTx(event models.SuiEventResponse, eventType EventType) (CancelTx, error) {
	if eventType != CancelTxEvent {
		return CancelTx{}, errors.Errorf("invalid event type %q", eventType)
	}

	parsedJSON := event.ParsedJson

	sender, err := extractStr(parsedJSON, "sender")
	if err != nil {
		return CancelTx{}, errors.Wrap(err, "unable to extract sender")
	}

	nonceRaw, err := extractStr(parsedJSON, "nonce")
	if err != nil {
		return CancelTx{}, errors.Wrap(err, "unable to extract nonce")
	}

	nonce, err := strconv.ParseUint(nonceRaw, 10, 64)
	if err != nil {
		return CancelTx{}, errors.Wrap(err, "unable to parse nonce")
	}
	if nonce <= 0 {
		return CancelTx{}, errors.New("nonce must be positive")
	}

	return CancelTx{
		Sender: sender,
		Nonce:  nonce,
	}, nil
}
