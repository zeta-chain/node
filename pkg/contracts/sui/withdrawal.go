package sui

import (
	"strconv"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
)

// Withdrawal represents data for a Sui withdraw event
type Withdrawal struct {
	CoinType CoinType
	Amount   math.Uint
	Sender   string
	Receiver string
	Nonce    uint64
}

// TokenAmount returns the amount of the withdrawal
func (d Withdrawal) TokenAmount() math.Uint {
	return d.Amount
}

// TxNonce returns the nonce of the withdrawal
func (d Withdrawal) TxNonce() uint64 {
	return d.Nonce
}

func (d Withdrawal) IsGas() bool {
	return d.CoinType == SUI
}

func parseWithdrawal(event models.SuiEventResponse, eventType EventType) (Withdrawal, error) {
	if eventType != WithdrawEvent {
		return Withdrawal{}, errors.Errorf("invalid event type %q", eventType)
	}

	parsedJSON := event.ParsedJson

	coinType, err := extractStr(parsedJSON, "coin_type")
	if err != nil {
		return Withdrawal{}, err
	}

	amountRaw, err := extractStr(parsedJSON, "amount")
	if err != nil {
		return Withdrawal{}, err
	}

	amount, err := math.ParseUint(amountRaw)
	if err != nil {
		return Withdrawal{}, errors.Wrap(err, "unable to parse amount")
	}

	sender, err := extractStr(parsedJSON, "sender")
	if err != nil {
		return Withdrawal{}, err
	}

	receiver, err := extractStr(parsedJSON, "receiver")
	if err != nil {
		return Withdrawal{}, err
	}

	nonceRaw, err := extractStr(parsedJSON, "nonce")
	if err != nil {
		return Withdrawal{}, err
	}

	nonce, err := strconv.ParseUint(nonceRaw, 10, 64)
	if err != nil {
		return Withdrawal{}, errors.Wrap(err, "unable to parse nonce")
	}

	return Withdrawal{
		CoinType: CoinType(coinType),
		Amount:   amount,
		Sender:   sender,
		Receiver: receiver,
		Nonce:    nonce,
	}, nil
}
