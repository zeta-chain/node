package ton

import (
	"errors"

	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	OpCall
)

// Outbound operations
const (
	OpWithdraw      Op = 200
	OpIncreaseSeqno Op = 205
)

// Authority operations (admin ops sent internally)
const (
	OpUpdateTSS  Op = 202
	OpResetSeqno Op = 206
)

// ExitCode represents an error code. Might be TVM or custom.
// TVM: https://docs.ton.org/v3/documentation/tvm/tvm-exit-codes
// Zeta: https://github.com/zeta-chain/protocol-contracts-ton/blob/main/contracts/common/errors.fc
type ExitCode uint32

const (
	ExitCodeInvalidSeqno ExitCode = 109
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

// UpdateTSS represents an admin operation to update the TSS address on the gateway.
// This is an authority operation that must be sent by the gateway admin.
type UpdateTSS struct {
	NewTSSAddress eth.Address
}

// AsBody casts struct as internal message body.
func (u UpdateTSS) AsBody() (*boc.Cell, error) {
	b := boc.NewCell()
	err := ErrCollect(
		b.WriteUint(uint64(OpUpdateTSS), sizeOpCode),
		b.WriteUint(0, sizeQueryID),
		b.WriteBytes(u.NewTSSAddress.Bytes()),
	)

	return b, err
}

// ResetSeqno represents an admin operation to reset the gateway's seqno (nonce).
type ResetSeqno struct {
	NewSeqno uint32
}

// AsBody casts struct as internal message body.
func (r ResetSeqno) AsBody() (*boc.Cell, error) {
	b := boc.NewCell()
	err := ErrCollect(
		b.WriteUint(uint64(OpResetSeqno), sizeOpCode),
		b.WriteUint(0, sizeQueryID),
		b.WriteUint(uint64(r.NewSeqno), sizeSeqno),
	)

	return b, err
}

// Deposit represents a deposit operation
type Deposit struct {
	Sender    ton.AccountID
	Amount    math.Uint
	Recipient eth.Address
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

// AsBody casts struct to internal message body.
func (d DepositAndCall) AsBody() (*boc.Cell, error) {
	b := boc.NewCell()

	return b, writeDepositAndCallBody(b, d.Recipient, d.CallData)
}

// Call represents a call operation
type Call struct {
	Sender    ton.AccountID
	Recipient eth.Address
	CallData  []byte
}

func (c Call) AsBody() (*boc.Cell, error) {
	b := boc.NewCell()

	return b, writeCallBody(b, c.Recipient, c.CallData)
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

func writeCallBody(b *boc.Cell, recipient eth.Address, callData []byte) error {
	if len(callData) == 0 {
		return errors.New("call data is empty")
	}

	callDataCell, err := MarshalSnakeCell(callData)
	if err != nil {
		return err
	}

	return ErrCollect(
		b.WriteUint(uint64(OpCall), sizeOpCode),
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

// SetSignature sets signature to the withdrawal message.
// Note that signature has the following order: [R, S, V (recovery ID)]
func (w *Withdrawal) SetSignature(sig [65]byte) { copy(w.Sig[:], sig[:]) }
func (w *Withdrawal) Signature() [65]byte       { return w.Sig }
func (w *Withdrawal) emptySig() bool            { return w.Sig == [65]byte{} }

// Hash returns hash of the withdrawal message. (used for signing)
func (w *Withdrawal) Hash() ([32]byte, error) {
	payload, err := w.payload()
	if err != nil {
		return [32]byte{}, err
	}

	return payload.Hash256()
}

// Signer returns EVM address of the signer (e.g. TSS)
func (w *Withdrawal) Signer() (eth.Address, error) {
	hash, err := w.Hash()
	if err != nil {
		return eth.Address{}, err
	}

	return deriveSigner(hash, w.Sig)
}

func (w *Withdrawal) AsBody() (*boc.Cell, error) {
	payload, err := w.payload()
	if err != nil {
		return nil, err
	}

	return messageToBody(payload, w.Sig)
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

// IncreaseSeqno represents an external message (an alternative to Withdrawal) that only
// increases seqno (nonce) and might contain reason code. Used as a factual tx for "canceling" CCTX.
type IncreaseSeqno struct {
	Seqno      uint32
	ReasonCode uint32
	Sig        [65]byte
}

func (is *IncreaseSeqno) SetSignature(sig [65]byte) { copy(is.Sig[:], sig[:]) }
func (is *IncreaseSeqno) Signature() [65]byte       { return is.Sig }
func (is *IncreaseSeqno) emptySig() bool            { return is.Sig == [65]byte{} }

func (is *IncreaseSeqno) Hash() ([32]byte, error) {
	payload, err := is.payload()
	if err != nil {
		return [32]byte{}, err
	}

	return payload.Hash256()
}

// Signer returns EVM address of the signer (e.g. TSS)
func (is *IncreaseSeqno) Signer() (eth.Address, error) {
	hash, err := is.Hash()
	if err != nil {
		return eth.Address{}, err
	}

	return deriveSigner(hash, is.Sig)
}

func (is *IncreaseSeqno) AsBody() (*boc.Cell, error) {
	payload, err := is.payload()
	if err != nil {
		return nil, err
	}

	return messageToBody(payload, is.Sig)
}

func (is *IncreaseSeqno) payload() (*boc.Cell, error) {
	payload := boc.NewCell()

	err := ErrCollect(
		payload.WriteUint(uint64(OpIncreaseSeqno), sizeOpCode),
		payload.WriteUint(uint64(is.ReasonCode), sizeOpCode),
		payload.WriteUint(uint64(is.Seqno), sizeSeqno),
	)

	if err != nil {
		return nil, errors.New("unable to marshal payload as cell")
	}

	return payload, nil
}

func deriveSigner(hash [32]byte, sig [65]byte) (eth.Address, error) {
	var sigCopy [65]byte
	copy(sigCopy[:], sig[:])

	// recovery id
	// https://bitcoin.stackexchange.com/questions/38351/ecdsa-v-r-s-what-is-v
	if sigCopy[64] >= 27 {
		sigCopy[64] -= 27
	}

	pub, err := crypto.SigToPub(hash[:], sigCopy[:])
	if err != nil {
		return eth.Address{}, err
	}

	return crypto.PubkeyToAddress(*pub), nil
}

func messageToBody(payload *boc.Cell, sig [65]byte) (*boc.Cell, error) {
	var (
		body    = boc.NewCell()
		v, r, s = splitSignature(sig)
	)

	// note that in TVM, the order of signature is different (v, r, s)
	err := ErrCollect(
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

// Ton Virtual Machine (TVM) uses different order of signature params (v,r,s) instead of (r,s,v);
// Let's split them as required.
func splitSignature(sig [65]byte) (v byte, r [32]byte, s [32]byte) {
	copy(r[:], sig[:32])
	copy(s[:], sig[32:64])
	v = sig[64]

	return v, r, s
}
