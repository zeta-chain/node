package ton

import (
	"context"

	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/wallet"
)

// Sender TON tx sender.
type Sender interface {
	Send(ctx context.Context, messages ...wallet.Sendable) error
}

// see https://docs.ton.org/develop/smart-contracts/messages#message-modes
const (
	SendFlagSeparateFees = uint8(1)
	SendFlagIgnoreErrors = uint8(2)
)

// SendDeposit sends a deposit operation to the gateway on behalf of the sender.
func (gw *Gateway) SendDeposit(
	ctx context.Context,
	s Sender,
	amount math.Uint,
	zevmRecipient eth.Address,
	sendMode uint8,
) error {
	body := boc.NewCell()

	if err := writeDepositBody(body, zevmRecipient); err != nil {
		return errors.Wrap(err, "failed to write deposit body")
	}

	return gw.send(ctx, s, amount, body, sendMode)
}

// SendDepositAndCall sends a deposit operation to the gateway on behalf of the sender
// with a callData to the recipient.
func (gw *Gateway) SendDepositAndCall(
	ctx context.Context,
	s Sender,
	amount math.Uint,
	zevmRecipient eth.Address,
	callData []byte,
	sendMode uint8,
) error {
	body := boc.NewCell()

	if err := writeDepositAndCallBody(body, zevmRecipient, callData); err != nil {
		return errors.Wrap(err, "failed to write depositAndCall body")
	}

	return gw.send(ctx, s, amount, body, sendMode)
}

func (gw *Gateway) send(ctx context.Context, s Sender, amount math.Uint, body *boc.Cell, sendMode uint8) error {
	if body == nil {
		return errors.New("body is nil")
	}

	return s.Send(ctx, wallet.Message{
		Amount:  tlb.Coins(amount.Uint64()),
		Address: gw.accountID,
		Body:    body,
		Mode:    sendMode,
	})
}
