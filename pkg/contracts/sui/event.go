package sui

import (
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
)

// CoinTypeSUI is the coin type for SUI, native gas token
const CoinTypeSUI = "0000000000000000000000000000000000000000000000000000000000000002::sui::SUI"

// Inbound represents data for a Sui inbound, it is parsed from a deposit or depositAndCall event
type Inbound struct {
	TxHash           string
	EventIndex       uint64
	CoinType         string
	Amount           string
	Sender           string
	Receiver         string
	Payload          []byte
	IsCrossChainCall bool
}

func (d Inbound) IsGasDeposit() bool {
	return d.CoinType == CoinTypeSUI
}

// parseInbound parses an inbound from a JSON read in the SUI event
// depositAndCall is a flag to indicate if the event is a depositAndCall event otherwise deposit event
func parseInbound(event models.SuiEventResponse, depositAndCall bool) (Inbound, error) {
	eventIndex, err := strconv.ParseUint(event.Id.EventSeq, 10, 64)
	if err != nil {
		return Inbound{}, errors.Wrap(err, "failed to parse event index")
	}

	parsedJSON := event.ParsedJson

	coinType, ok := parsedJSON["coin_type"].(string)
	if !ok {
		return Inbound{}, errors.New("invalid coin type")
	}

	amount, ok := parsedJSON["amount"].(string)
	if !ok {
		return Inbound{}, errors.New("invalid amount")
	}

	sender, ok := parsedJSON["sender"].(string)
	if !ok {
		return Inbound{}, errors.New("invalid sender")
	}

	receiver, ok := parsedJSON["receiver"].(string)
	if !ok {
		return Inbound{}, errors.New("invalid receiver")
	}

	payload := []byte{}
	isCrossChainCall := false
	if depositAndCall {
		parsedPayload, ok := parsedJSON["payload"].(string)
		if !ok {
			return Inbound{}, errors.New("invalid payload")
		}
		payload = []byte(parsedPayload)
		isCrossChainCall = true
	}

	return Inbound{
		TxHash:           event.Id.TxDigest,
		EventIndex:       eventIndex,
		CoinType:         coinType,
		Amount:           amount,
		Sender:           sender,
		Receiver:         receiver,
		IsCrossChainCall: isCrossChainCall,
		Payload:          payload,
	}, nil
}
