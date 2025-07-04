package ton

import (
	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/ton"
)

// Transaction represents a Gateway transaction.
type Transaction struct {
	ton.Transaction
	Operation Op
	ExitCode  int32

	content any
	inbound bool
}

// OutboundAuth contains the outbound seqno, signature and signer.
type OutboundAuth struct {
	Seqno  uint32
	Sig    [65]byte
	Signer eth.Address
}

// IsInbound returns true if the transaction is inbound.
func (tx *Transaction) IsInbound() bool {
	return tx.inbound
}

// IsOutbound returns true if the transaction is outbound.
func (tx *Transaction) IsOutbound() bool {
	return !tx.inbound
}

// GasUsed returns the amount of gas used by the transaction.
func (tx *Transaction) GasUsed() math.Uint {
	return math.NewUint(uint64(tx.TotalFees.Grams))
}

// Donation casts the transaction content to a Donation.
func (tx *Transaction) Donation() (Donation, error) {
	return retrieveContent[Donation](tx)
}

// Deposit casts the transaction content to a Deposit.
func (tx *Transaction) Deposit() (Deposit, error) {
	return retrieveContent[Deposit](tx)
}

// DepositAndCall casts the transaction content to a DepositAndCall.
func (tx *Transaction) DepositAndCall() (DepositAndCall, error) {
	return retrieveContent[DepositAndCall](tx)
}

// Call casts the transaction content to a Call.
func (tx *Transaction) Call() (Call, error) {
	return retrieveContent[Call](tx)
}

// Withdrawal casts the transaction content to a Withdrawal.
func (tx *Transaction) Withdrawal() (Withdrawal, error) {
	return retrieveContent[Withdrawal](tx)
}

func (tx *Transaction) IncreaseSeqno() (IncreaseSeqno, error) {
	return retrieveContent[IncreaseSeqno](tx)
}

// OutboundAuth returns the outbound seqno and signature
func (tx *Transaction) OutboundAuth() (OutboundAuth, error) {
	if !tx.IsOutbound() {
		return OutboundAuth{}, errors.New("not an outbound")
	}

	if tx.Operation == OpWithdraw {
		w, err := tx.Withdrawal()
		if err != nil {
			return OutboundAuth{}, errors.Wrap(err, "unable to get withdrawal")
		}

		signer, err := w.Signer()
		if err != nil {
			return OutboundAuth{}, errors.Wrap(err, "unable to get signer")
		}

		return OutboundAuth{
			Seqno:  w.Seqno,
			Sig:    w.Sig,
			Signer: signer,
		}, nil
	}

	if tx.Operation == OpIncreaseSeqno {
		is, err := tx.IncreaseSeqno()
		if err != nil {
			return OutboundAuth{}, errors.Wrap(err, "unable to get increase seqno")
		}

		signer, err := is.Signer()
		if err != nil {
			return OutboundAuth{}, errors.Wrap(err, "unable to get signer")
		}

		return OutboundAuth{
			Seqno:  is.Seqno,
			Sig:    is.Sig,
			Signer: signer,
		}, nil
	}

	return OutboundAuth{}, errors.Wrapf(ErrUnknownOp, "op %d", tx.Operation)
}

func retrieveContent[T any](tx *Transaction) (T, error) {
	typed, ok := tx.content.(T)
	if !ok {
		var tt T
		return tt, errors.Wrapf(ErrCast, "not a %T (op %d)", tt, int(tx.Operation))
	}

	return typed, nil
}
