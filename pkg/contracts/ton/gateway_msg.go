package ton

import (
	"errors"

	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// Op operation code
type Op uint32

// github.com/zeta-chain/protocol-contracts-ton/blob/main/contracts/gateway.fc
// Inbound operations
const (
	OpDonate Op = 100 + iota
	OpDeposit
	OpDepositAndCall
)

// Outbound operations
const (
	OpWithdraw Op = 200
)

// Donation represents a donation operation
type Donation struct {
	Sender ton.AccountID
	Amount math.Uint
}

// AsBody casts struct as internal message body.
func (d Donation) AsBody() (*boc.Cell, error) {
	b := boc.NewCell()
	err := ErrCollect(
		b.WriteUint(uint64(OpDonate), sizeOpCode),
		b.WriteUint(0, sizeQueryID),
	)

	return b, err
}

// Deposit represents a deposit operation
type Deposit struct {
	Sender    ton.AccountID
	Amount    math.Uint
	Recipient eth.Address
}

// Memo casts deposit to memo bytes
func (d Deposit) Memo() []byte {
	return d.Recipient.Bytes()
}

// AsBody casts struct as internal message body.
func (d Deposit) AsBody() (*boc.Cell, error) {
	b := boc.NewCell()

	return b, writeDepositBody(b, d.Recipient)
}

// DepositAndCall represents a deposit and call operation
type DepositAndCall struct {
	Deposit
	CallData []byte
}

// Memo casts deposit to call to memo bytes
func (d DepositAndCall) Memo() []byte {
	recipient := d.Recipient.Bytes()
	out := make([]byte, 0, len(recipient)+len(d.CallData))

	out = append(out, recipient...)
	out = append(out, d.CallData...)

	return out
}

// AsBody casts struct to internal message body.
func (d DepositAndCall) AsBody() (*boc.Cell, error) {
	b := boc.NewCell()

	return b, writeDepositAndCallBody(b, d.Recipient, d.CallData)
}

func writeDepositBody(b *boc.Cell, recipient eth.Address) error {
	return ErrCollect(
		b.WriteUint(uint64(OpDeposit), sizeOpCode),
		b.WriteUint(0, sizeQueryID),
		b.WriteBytes(recipient.Bytes()),
	)
}

func writeDepositAndCallBody(b *boc.Cell, recipient eth.Address, callData []byte) error {
	if len(callData) == 0 {
		return errors.New("call data is empty")
	}

	callDataCell, err := MarshalSnakeCell(callData)
	if err != nil {
		return err
	}

	return ErrCollect(
		b.WriteUint(uint64(OpDepositAndCall), sizeOpCode),
		b.WriteUint(0, sizeQueryID),
		b.WriteBytes(recipient.Bytes()),
		b.AddRef(callDataCell),
	)
}

// Withdrawal represents a withdrawal external message
type Withdrawal struct {
	Recipient ton.AccountID
	Amount    math.Uint
	Seqno     uint32
	Sig       [65]byte
}

func (w *Withdrawal) emptySig() bool {
	return w.Sig == [65]byte{}
}

// Hash returns hash of the withdrawal message. (used for signing)
func (w *Withdrawal) Hash() ([32]byte, error) {
	payload, err := w.payload()
	if err != nil {
		return [32]byte{}, err
	}

	return payload.Hash256()
}

// SetSignature sets signature to the withdrawal message.
// Note that signature has the following order: [R, S, V (recovery ID)]
func (w *Withdrawal) SetSignature(sig [65]byte) {
	copy(w.Sig[:], sig[:])
}

func (w *Withdrawal) AsBody() (*boc.Cell, error) {
	payload, err := w.payload()
	if err != nil {
		return nil, err
	}

	var (
		body    = boc.NewCell()
		v, r, s = splitSignature(w.Sig)
	)

	// note that in TVM, the order of signature is different (v, r, s)
	err = ErrCollect(
		body.WriteUint(uint64(v), 8),
		body.WriteBytes(r[:]),
		body.WriteBytes(s[:]),
		body.AddRef(payload),
	)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (w *Withdrawal) payload() (*boc.Cell, error) {
	payload := boc.NewCell()

	err := ErrCollect(
		payload.WriteUint(uint64(OpWithdraw), sizeOpCode),
		tlb.Marshal(payload, w.Recipient.ToMsgAddress()),
		tlb.Marshal(payload, tlb.Coins(w.Amount.Uint64())),
		payload.WriteUint(uint64(w.Seqno), sizeSeqno),
	)

	if err != nil {
		return nil, errors.New("unable to marshal payload as cell")
	}

	return payload, nil
}

// Ton Virtual Machine (TVM) uses different order of signature params (v,r,s)
// let's split them as required
func splitSignature(sig [65]byte) (v byte, r [32]byte, s [32]byte) {
	v = sig[64]
	copy(r[:], sig[:32])
	copy(s[:], sig[32:64])

	return v, r, s
}
