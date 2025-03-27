package sui

import (
	"strconv"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
)

// CanceTxNonceEvent represents data for a Sui cancelled tx nonce increase event
type CanceTxNonceEvent struct {
	Sender string
	Nonce  uint64
}

// TokenAmount returns hardcoded zero amount
// because nonce increase does not involve any token transfer
func (n CanceTxNonceEvent) TokenAmount() math.Uint {
	return math.NewUint(0)
}

// TxNonce returns the nonce of the tx
// Nonce is the incremented nonce value (CCTX.nonce + 1)
// see: https://github.com/zeta-chain/protocol-contracts-sui/blob/e5a756e473da884dcbc59b574b387a7a365ac823/sources/gateway.move#L140
func (n CanceTxNonceEvent) TxNonce() uint64 {
	return n.Nonce - 1
}

func parseNonceIncrease(event models.SuiEventResponse, eventType EventType) (CanceTxNonceEvent, error) {
	if eventType != CancelTxNonceEvent {
		return CanceTxNonceEvent{}, errors.Errorf("invalid event type %q", eventType)
	}

	parsedJSON := event.ParsedJson

	sender, err := extractStr(parsedJSON, "sender")
	if err != nil {
		return CanceTxNonceEvent{}, err
	}

	nonceRaw, err := extractStr(parsedJSON, "nonce")
	if err != nil {
		return CanceTxNonceEvent{}, err
	}

	nonce, err := strconv.ParseUint(nonceRaw, 10, 64)
	if err != nil {
		return CanceTxNonceEvent{}, errors.Wrap(err, "unable to parse nonce")
	}
	if nonce <= 0 {
		return CanceTxNonceEvent{}, errors.New("nonce must be positive")
	}

	return CanceTxNonceEvent{
		Sender: sender,
		Nonce:  nonce,
	}, nil
}
