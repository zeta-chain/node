package sui

import (
	"strconv"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
)

// WithdrawData represents data for a Sui withdraw event
type WithdrawData struct {
	CoinType CoinType
	Amount   math.Uint
	Sender   string
	Receiver string
	Nonce    uint64
}

func (d *WithdrawData) IsGas() bool {
	return d.CoinType == SUI
}

func parseWithdrawal(event models.SuiEventResponse, eventType EventType) (WithdrawData, error) {
	if eventType != Withdraw {
		return WithdrawData{}, errors.Errorf("invalid event type %q", eventType)
	}

	parsedJSON := event.ParsedJson

	coinType, err := extractStr(parsedJSON, "coin_type")
	if err != nil {
		return WithdrawData{}, err
	}

	amountRaw, err := extractStr(parsedJSON, "amount")
	if err != nil {
		return WithdrawData{}, err
	}

	amount, err := math.ParseUint(amountRaw)
	if err != nil {
		return WithdrawData{}, errors.Wrap(err, "unable to parse amount")
	}

	sender, err := extractStr(parsedJSON, "sender")
	if err != nil {
		return WithdrawData{}, err
	}

	receiver, err := extractStr(parsedJSON, "receiver")
	if err != nil {
		return WithdrawData{}, err
	}

	nonceRaw, err := extractStr(parsedJSON, "nonce")
	if err != nil {
		return WithdrawData{}, err
	}

	nonce, err := strconv.ParseUint(nonceRaw, 10, 64)
	if err != nil {
		return WithdrawData{}, errors.Wrap(err, "unable to parse nonce")
	}

	return WithdrawData{
		CoinType: CoinType(coinType),
		Amount:   amount,
		Sender:   sender,
		Receiver: receiver,
		Nonce:    nonce,
	}, nil
}
