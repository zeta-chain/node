package sui

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// Deposit represents data for a Sui deposit/depositAndCall event
type Deposit struct {
	// Note: CoinType is what is used as Asset field in the ForeignCoin object
	CoinType              CoinType
	Amount                math.Uint
	Sender                string
	Receiver              ethcommon.Address
	Payload               []byte
	IsCrossChainCall      bool
	InvalidDepositMessage string
}

func (d *Deposit) IsGas() bool {
	return d.CoinType == SUI
}

func (d *Deposit) IsInvalid() bool {
	return d.InvalidDepositMessage != ""
}

func parseDeposit(event models.SuiEventResponse, eventType EventType) (Deposit, error) {
	parsedJSON := event.ParsedJson

	coinType, err := extractStr(parsedJSON, "coin_type")
	if err != nil {
		return Deposit{}, err
	}

	amountRaw, err := extractStr(parsedJSON, "amount")
	if err != nil {
		return Deposit{}, err
	}

	amount, err := math.ParseUint(amountRaw)
	if err != nil {
		return Deposit{}, errors.Wrap(err, "unable to parse amount")
	}

	sender, err := extractStr(parsedJSON, "sender")
	if err != nil {
		return Deposit{}, err
	}

	receiverRaw, err := extractStr(parsedJSON, "receiver")
	if err != nil {
		return Deposit{}, err
	}

	receiver := ethcommon.HexToAddress(receiverRaw)
	var invalidMessage string
	if receiver == (ethcommon.Address{}) {
		// receiver is data set by the user, if the format is invalid we don't return an error but set the deposit as invalid
		// so it is observed and reverted to the sender
		invalidMessage = fmt.Sprintf("invalid receiver address %q", receiverRaw)
	}

	var isCrosschainCall bool
	var payload []byte

	if eventType == DepositAndCallEvent {
		isCrosschainCall = true

		payloadRaw, ok := parsedJSON["payload"].([]any)
		if !ok {
			return Deposit{}, errors.New("invalid payload")
		}

		payload, err = convertPayload(payloadRaw)
		if err != nil {
			return Deposit{}, errors.Wrap(err, "unable to convert payload")
		}
	}

	return Deposit{
		CoinType:              CoinType(coinType),
		Amount:                amount,
		Sender:                sender,
		Receiver:              receiver,
		Payload:               payload,
		IsCrossChainCall:      isCrosschainCall,
		InvalidDepositMessage: invalidMessage,
	}, nil
}
