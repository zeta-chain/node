package sui

import (
	"strconv"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
)

// NonceIncrease represents data for a Sui nonce increase event
type NonceIncrease struct {
	Sender string
	Nonce  uint64
}

// TokenAmount returns hardcoded zero amount
// because nonce increase does not involve any token transfer
func (n NonceIncrease) TokenAmount() math.Uint {
	return math.NewUint(0)
}

// GatewayNonce returns the nonce of the nonce increase
func (n NonceIncrease) GatewayNonce() uint64 {
	return n.Nonce
}

func parseNonceIncrease(event models.SuiEventResponse, eventType EventType) (NonceIncrease, error) {
	if eventType != NonceIncreaseEvent {
		return NonceIncrease{}, errors.Errorf("invalid event type %q", eventType)
	}

	parsedJSON := event.ParsedJson

	sender, err := extractStr(parsedJSON, "sender")
	if err != nil {
		return NonceIncrease{}, err
	}

	nonceRaw, err := extractStr(parsedJSON, "nonce")
	if err != nil {
		return NonceIncrease{}, err
	}

	nonce, err := strconv.ParseUint(nonceRaw, 10, 64)
	if err != nil {
		return NonceIncrease{}, errors.Wrap(err, "unable to parse nonce")
	}

	return NonceIncrease{
		Sender: sender,
		Nonce:  nonce,
	}, nil
}
