package sui

import (
	"fmt"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// Inbound represents data for a Sui inbound, it is parsed from a deposit or depositAndCall event
type Inbound struct {
	TxHash     string
	EventIndex uint64
	// Note: CoinType is what is used as Asset field in the ForeignCoin object
	CoinType         CoinType
	Amount           uint64
	Sender           string
	Receiver         ethcommon.Address
	Payload          []byte
	IsCrossChainCall bool
}

func (d *Inbound) IsGasDeposit() bool {
	return d.CoinType == SUI
}

// parseInbound parses an inbound from a JSON read in the SUI event
// depositAndCall is a flag to indicate if the event is a depositAndCall event otherwise deposit event
// TODO: add specific error that can be handled when the event data is invalid
// https://github.com/zeta-chain/node/issues/3502
func parseInbound(event models.SuiEventResponse, eventType string) (Inbound, error) {
	eventIndex, err := strconv.ParseUint(event.Id.EventSeq, 10, 64)
	if err != nil {
		return Inbound{}, errors.Wrap(err, "failed to parse event index")
	}

	parsedJSON := event.ParsedJson

	coinType, ok := parsedJSON["coin_type"].(string)
	if !ok {
		return Inbound{}, errors.New("invalid coin type")
	}

	parsedAmount, ok := parsedJSON["amount"].(string)
	if !ok {
		return Inbound{}, errors.New("invalid amount")
	}
	amount, err := strconv.ParseUint(parsedAmount, 10, 64)
	if err != nil {
		return Inbound{}, errors.Wrap(err, "failed to parse amount")
	}

	sender, ok := parsedJSON["sender"].(string)
	if !ok {
		return Inbound{}, errors.New("invalid sender")
	}

	parsedReceiver, ok := parsedJSON["receiver"].(string)
	if !ok {
		return Inbound{}, errors.New("invalid receiver")
	}

	receiver := ethcommon.HexToAddress(parsedReceiver)
	if receiver == (ethcommon.Address{}) {
		return Inbound{}, errors.New("can't parse receiver address")
	}

	isCrosschainCall := false

	var payload []byte
	if eventType == eventDepositAndCall {
		isCrosschainCall = true

		parsedPayload, ok := parsedJSON["payload"].([]interface{})
		if !ok {
			return Inbound{}, errors.New("invalid payload")
		}
		payload, err = convertPayload(parsedPayload)
		if err != nil {
			return Inbound{}, errors.Wrap(err, "failed to convert payload")
		}
	}

	return Inbound{
		TxHash:           event.Id.TxDigest,
		EventIndex:       eventIndex,
		CoinType:         CoinType(coinType),
		Amount:           amount,
		Sender:           sender,
		Receiver:         receiver,
		Payload:          payload,
		IsCrossChainCall: isCrosschainCall,
	}, nil
}

// convertPayload
func convertPayload(data []interface{}) ([]byte, error) {
	payload := make([]byte, 0, len(data))
	for i, d := range data {
		// parsed bytes are represented as float64
		if b, ok := d.(float64); ok {
			if b < 0 || b > 255 {
				return nil, fmt.Errorf("invalid payload data at index %d, not a byte (%f)", i, b)
			}
			payload = append(payload, byte(b))
		} else {
			return nil, fmt.Errorf("invalid payload data at index %d, not a float64", i)
		}
	}
	return payload, nil
}
