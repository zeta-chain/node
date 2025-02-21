package sui

import (
	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// DepositData represents data for a Sui deposit/depositAndCall event
type DepositData struct {
	// Note: CoinType is what is used as Asset field in the ForeignCoin object
	CoinType         CoinType
	Amount           math.Uint
	Sender           string
	Receiver         ethcommon.Address
	Payload          []byte
	IsCrossChainCall bool
}

func (d *DepositData) IsGas() bool {
	return d.CoinType == SUI
}

func parseDeposit(event models.SuiEventResponse, eventType EventType) (DepositData, error) {
	parsedJSON := event.ParsedJson

	coinType, err := extractStr(parsedJSON, "coin_type")
	if err != nil {
		return DepositData{}, err
	}

	amountRaw, err := extractStr(parsedJSON, "amount")
	if err != nil {
		return DepositData{}, err
	}

	amount, err := math.ParseUint(amountRaw)
	if err != nil {
		return DepositData{}, errors.Wrap(err, "unable to parse amount")
	}

	sender, err := extractStr(parsedJSON, "sender")
	if err != nil {
		return DepositData{}, err
	}

	receiverRaw, err := extractStr(parsedJSON, "receiver")
	if err != nil {
		return DepositData{}, err
	}

	receiver := ethcommon.HexToAddress(receiverRaw)
	if receiver == (ethcommon.Address{}) {
		return DepositData{}, errors.Errorf("invalid receiver address %q", receiverRaw)
	}

	var isCrosschainCall bool
	var payload []byte

	if eventType == DepositAndCall {
		isCrosschainCall = true

		payloadRaw, ok := parsedJSON["payload"].([]any)
		if !ok {
			return DepositData{}, errors.New("invalid payload")
		}

		payload, err = convertPayload(payloadRaw)
		if err != nil {
			return DepositData{}, errors.Wrap(err, "unable to convert payload")
		}
	}

	return DepositData{
		CoinType:         CoinType(coinType),
		Amount:           amount,
		Sender:           sender,
		Receiver:         receiver,
		Payload:          payload,
		IsCrossChainCall: isCrosschainCall,
	}, nil
}
