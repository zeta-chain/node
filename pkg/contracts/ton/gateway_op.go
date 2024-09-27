package ton

import (
	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/tonkeeper/tongo/boc"
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
	OpWithdraw Op = 200 + iota
	SetDepositsEnabled
	UpdateTSS
	UpdateCode
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
	err := ErrCollect(
		b.WriteUint(uint64(OpDeposit), sizeOpCode),
		b.WriteUint(0, sizeQueryID),
		b.WriteBytes(d.Recipient.Bytes()),
	)

	return b, err
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
	callDataCell, err := MarshalSnakeCell(d.CallData)
	if err != nil {
		return nil, err
	}

	b := boc.NewCell()
	err = ErrCollect(
		b.WriteUint(uint64(OpDepositAndCall), sizeOpCode),
		b.WriteUint(0, sizeQueryID),
		b.WriteBytes(d.Recipient.Bytes()),
		b.AddRef(callDataCell),
	)

	return b, err
}
