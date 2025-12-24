package ton

import (
	"context"

	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
)

// Sender TON tx sender. Usually an interface to a wallet.
type Sender interface {
	Send(ctx context.Context, messages ...wallet.Sendable) error
}

// Client represents a sender that allows sending an arbitrary external message to the network.
type Client interface {
	SendMessage(ctx context.Context, payload []byte) (uint32, error)
}

type ExternalMsg interface {
	Hash() ([32]byte, error)
	AsBody() (*boc.Cell, error)

	Signature() [65]byte
	SetSignature(sig [65]byte)

	emptySig() bool
}

// See https://docs.ton.org/v3/documentation/smart-contracts/message-management/message-modes-cookbook
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

// SendCall sends `call` operation to the gateway on behalf of the sender.
// The amount should be >= calculate_gas_fee(op::call)
func (gw *Gateway) SendCall(
	ctx context.Context,
	s Sender,
	amount math.Uint,
	zevmRecipient eth.Address,
	callData []byte,
	sendMode uint8,
) error {
	body := boc.NewCell()

	if err := writeCallBody(body, zevmRecipient, callData); err != nil {
		return errors.Wrap(err, "failed to write call body")
	}

	return gw.send(ctx, s, amount, body, sendMode)
}

// SendUpdateTSS sends an admin operation to update the TSS address on the gateway.
func (gw *Gateway) SendUpdateTSS(
	ctx context.Context,
	s Sender,
	amount math.Uint,
	newTSSAddress eth.Address,
	sendMode uint8,
) error {
	updateTSS := UpdateTSS{NewTSSAddress: newTSSAddress}

	body, err := updateTSS.AsBody()
	if err != nil {
		return errors.Wrap(err, "failed to create update_tss body")
	}

	return gw.send(ctx, s, amount, body, sendMode)
}

// SendResetSeqno sends an admin operation to reset the gateway's seqno (nonce).
func (gw *Gateway) SendResetSeqno(
	ctx context.Context,
	s Sender,
	amount math.Uint,
	newSeqno uint32,
	sendMode uint8,
) error {
	resetSeqno := ResetSeqno{NewSeqno: newSeqno}

	body, err := resetSeqno.AsBody()
	if err != nil {
		return errors.Wrap(err, "failed to create reset_seqno body")
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

// SendExternalMessage sends an external message to the Gateway.
func (gw *Gateway) SendExternalMessage(ctx context.Context, s Client, msg ExternalMsg) (uint32, error) {
	return sendExternalMessage(ctx, s, gw.accountID, msg)
}

// inspired by tongo's wallet.Wallet{}.RawSendV2()
func sendExternalMessage(ctx context.Context, via Client, to ton.AccountID, msg ExternalMsg) (uint32, error) {
	if msg.emptySig() {
		return 0, errors.New("empty signature")
	}

	body, err := msg.AsBody()
	if err != nil {
		return 0, errors.Wrap(err, "unable to encode msg as cell")
	}

	extMsg, err := ton.CreateExternalMessage(to, body, nil, tlb.VarUInteger16{})
	if err != nil {
		return 0, errors.Wrap(err, "unable to create external message")
	}

	extMsgCell := boc.NewCell()
	err = tlb.Marshal(extMsgCell, extMsg)
	if err != nil {
		return 0, errors.Wrap(err, "can not marshal wallet external message")
	}

	payload, err := extMsgCell.ToBocCustom(false, false, false, 0)
	if err != nil {
		return 0, errors.Wrap(err, "can not serialize external message cell")
	}

	return via.SendMessage(ctx, payload)
}
