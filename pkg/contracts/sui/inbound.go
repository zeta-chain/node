package sui

import (
	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// Inbound represents data for a Sui inbound,
// it is parsed from a deposit/depositAndCall event
type Inbound struct {
	// Note: CoinType is what is used as Asset field in the ForeignCoin object
	CoinType         CoinType
	Amount           math.Uint
	Sender           string
	Receiver         ethcommon.Address
	Payload          []byte
	IsCrossChainCall bool
}

func (d *Inbound) IsGasDeposit() bool {
	return d.CoinType == SUI
}

func parseInbound(event models.SuiEventResponse, eventType EventType) (Inbound, error) {
	parsedJSON := event.ParsedJson

	coinType, err := extractStr(parsedJSON, "coin_type")
	if err != nil {
		return Inbound{}, err
	}

	amountRaw, err := extractStr(parsedJSON, "amount")
	if err != nil {
		return Inbound{}, err
	}

	amount, err := math.ParseUint(amountRaw)
	if err != nil {
		return Inbound{}, errors.Wrap(err, "unable to parse amount")
	}

	sender, err := extractStr(parsedJSON, "sender")
	if err != nil {
		return Inbound{}, err
	}

	receiverRaw, err := extractStr(parsedJSON, "receiver")
	if err != nil {
		return Inbound{}, err
	}

	receiver := ethcommon.HexToAddress(receiverRaw)
	if receiver == (ethcommon.Address{}) {
		return Inbound{}, errors.Errorf("invalid receiver address %q", receiverRaw)
	}

	var isCrosschainCall bool
	var payload []byte

	if eventType == DepositAndCall {
		isCrosschainCall = true

		payloadRaw, ok := parsedJSON["payload"].([]any)
		if !ok {
			return Inbound{}, errors.New("invalid payload")
		}

		payload, err = convertPayload(payloadRaw)
		if err != nil {
			return Inbound{}, errors.Wrap(err, "unable to convert payload")
		}
	}

	return Inbound{
		CoinType:         CoinType(coinType),
		Amount:           amount,
		Sender:           sender,
		Receiver:         receiver,
		Payload:          payload,
		IsCrossChainCall: isCrosschainCall,
	}, nil
}
