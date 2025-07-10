package sample

import (
	"crypto/rand"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"

	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
)

const (
	tonWorkchainID   = 0
	tonShardID       = 123
	tonSampleTxFee   = 006_500_000 // 0.0065 TON
	tonSampleGasUsed = 8500
)

type TONTransactionProps struct {
	Account      ton.AccountID
	GasUsed      uint64
	TotalTONFees uint64
	BlockID      ton.BlockIDExt

	// For simplicity let's have only one input
	// and one output (both optional)
	Input  *tlb.Message
	Output *tlb.Message
}

type intMsgInfo struct {
	IhrDisabled bool
	Bounce      bool
	Bounced     bool
	Src         tlb.MsgAddress
	Dest        tlb.MsgAddress
	Value       tlb.CurrencyCollection
	IhrFee      tlb.Grams
	FwdFee      tlb.Grams
	CreatedLt   uint64
	CreatedAt   uint32
}

func TONDonation(t *testing.T, acc ton.AccountID, d toncontracts.Donation) ton.Transaction {
	return TONTransaction(t, TONDonateProps(t, acc, d))
}

func TONDonateProps(t *testing.T, acc ton.AccountID, d toncontracts.Donation) TONTransactionProps {
	body, err := d.AsBody()
	require.NoError(t, err)

	tonSent := tonSampleTxFee + d.Amount.Uint64()

	return TONTransactionProps{
		Account: acc,
		Input: &tlb.Message{
			Info: internalMessageInfo(&intMsgInfo{
				Bounce: true,
				Src:    d.Sender.ToMsgAddress(),
				Dest:   acc.ToMsgAddress(),
				Value:  tlb.CurrencyCollection{Grams: tlb.Grams(tonSent)},
			}),
			Body: tlb.EitherRef[tlb.Any]{Value: tlb.Any(*body)},
		},
	}
}

func TONDeposit(t *testing.T, acc ton.AccountID, d toncontracts.Deposit) ton.Transaction {
	return TONTransaction(t, TONDepositProps(t, acc, d))
}

func TONDepositProps(t *testing.T, acc ton.AccountID, d toncontracts.Deposit) TONTransactionProps {
	body, err := d.AsBody()
	require.NoError(t, err)

	logBody := depositLogMock(t, d.Amount.Uint64(), tonSampleTxFee)

	return TONTransactionProps{
		Account: acc,
		Input: &tlb.Message{
			Info: internalMessageInfo(&intMsgInfo{
				Bounce: true,
				Src:    d.Sender.ToMsgAddress(),
				Dest:   acc.ToMsgAddress(),
				Value:  tlb.CurrencyCollection{Grams: fakeDepositAmount(d.Amount)},
			}),
			Body: tlb.EitherRef[tlb.Any]{Value: tlb.Any(*body)},
		},
		Output: &tlb.Message{
			Body: tlb.EitherRef[tlb.Any]{IsRight: true, Value: tlb.Any(*logBody)},
		},
	}
}

func TONDepositAndCall(t *testing.T, acc ton.AccountID, d toncontracts.DepositAndCall) ton.Transaction {
	return TONTransaction(t, TONDepositAndCallProps(t, acc, d))
}

func TONDepositAndCallProps(t *testing.T, acc ton.AccountID, d toncontracts.DepositAndCall) TONTransactionProps {
	body, err := d.AsBody()
	require.NoError(t, err)

	logBody := depositLogMock(t, d.Amount.Uint64(), tonSampleTxFee)

	return TONTransactionProps{
		Account: acc,
		Input: &tlb.Message{
			Info: internalMessageInfo(&intMsgInfo{
				Bounce: true,
				Src:    d.Sender.ToMsgAddress(),
				Dest:   acc.ToMsgAddress(),
				Value:  tlb.CurrencyCollection{Grams: fakeDepositAmount(d.Amount)},
			}),
			Body: tlb.EitherRef[tlb.Any]{Value: tlb.Any(*body)},
		},
		Output: &tlb.Message{
			Body: tlb.EitherRef[tlb.Any]{IsRight: true, Value: tlb.Any(*logBody)},
		},
	}
}

func TONCall(t *testing.T, acc ton.AccountID, c toncontracts.Call) ton.Transaction {
	return TONTransaction(t, TONCallProps(t, acc, c))
}

func TONCallProps(t *testing.T, acc ton.AccountID, c toncontracts.Call) TONTransactionProps {
	body, err := c.AsBody()
	require.NoError(t, err)

	return TONTransactionProps{
		Account: acc,
		Input: &tlb.Message{
			Info: internalMessageInfo(&intMsgInfo{
				Bounce: true,
				Src:    c.Sender.ToMsgAddress(),
				Dest:   acc.ToMsgAddress(),
				Value:  tlb.CurrencyCollection{Grams: tlb.Coins(tonSampleTxFee)},
			}),
			Body: tlb.EitherRef[tlb.Any]{Value: tlb.Any(*body)},
		},
	}
}

func TONWithdrawal(t *testing.T, acc ton.AccountID, w toncontracts.Withdrawal) ton.Transaction {
	return TONTransaction(t, TONWithdrawalProps(t, acc, w))
}

func TONWithdrawalProps(t *testing.T, acc ton.AccountID, w toncontracts.Withdrawal) TONTransactionProps {
	body, err := w.AsBody()
	require.NoError(t, err)

	return TONTransactionProps{
		Account: acc,
		Input: &tlb.Message{
			Info: externalMessageInfo(acc),
			Body: tlb.EitherRef[tlb.Any]{Value: tlb.Any(*body)},
		},
		Output: &tlb.Message{
			Info: internalMessageInfo(&intMsgInfo{
				IhrDisabled: true,
				Src:         acc.ToMsgAddress(),
				Dest:        w.Recipient.ToMsgAddress(),
				Value:       tlb.CurrencyCollection{Grams: tlb.Coins(w.Amount.Uint64())},
			}),
		},
	}
}

func TONIncreaseSeqno(t *testing.T, acc ton.AccountID, is toncontracts.IncreaseSeqno) ton.Transaction {
	return TONTransaction(t, TONIncreaseSeqnoProps(t, acc, is))
}

func TONIncreaseSeqnoProps(t *testing.T, acc ton.AccountID, is toncontracts.IncreaseSeqno) TONTransactionProps {
	body, err := is.AsBody()
	require.NoError(t, err)

	return TONTransactionProps{
		Account: acc,
		Input: &tlb.Message{
			Info: externalMessageInfo(acc),
			Body: tlb.EitherRef[tlb.Any]{Value: tlb.Any(*body)},
		},
	}
}

// TONTransaction creates a sample TON transaction.
func TONTransaction(t *testing.T, p TONTransactionProps) ton.Transaction {
	require.False(t, p.Account.IsZero(), "account address is empty")
	require.False(t, p.Input == nil && p.Output == nil, "both input and output are empty")

	now := time.Now().UTC()

	if p.GasUsed == 0 {
		p.GasUsed = tonSampleGasUsed
	}

	if p.TotalTONFees == 0 {
		p.TotalTONFees = tonSampleTxFee
	}

	if p.BlockID.BlockID.Seqno == 0 {
		p.BlockID = tonBlockID(now)
	}

	// Simulate logical time as `2 * now()`
	lt := uint64(2 * now.Unix())

	input := tlb.Maybe[tlb.Ref[tlb.Message]]{}
	if p.Input != nil {
		input.Exists = true
		input.Value.Value = *p.Input
	}

	var outputs tlb.HashmapE[tlb.Uint15, tlb.Ref[tlb.Message]]
	if p.Output != nil {
		outputs = tlb.NewHashmapE(
			[]tlb.Uint15{0},
			[]tlb.Ref[tlb.Message]{{*p.Output}},
		)
	}

	type messages struct {
		InMsg   tlb.Maybe[tlb.Ref[tlb.Message]]
		OutMsgs tlb.HashmapE[tlb.Uint15, tlb.Ref[tlb.Message]]
	}

	tx := ton.Transaction{
		BlockID: p.BlockID,
		Transaction: tlb.Transaction{
			AccountAddr: p.Account.Address,
			Lt:          lt,
			Now:         uint32(now.Unix()),
			OutMsgCnt:   tlb.Uint15(len(outputs.Keys())),
			TotalFees:   tlb.CurrencyCollection{Grams: tlb.Grams(p.TotalTONFees)},
			Msgs:        messages{InMsg: input, OutMsgs: outputs},
		},
	}

	setTXHash(&tx.Transaction, Hash())

	return tx
}

func GenerateTONAccountID() ton.AccountID {
	var addr [32]byte

	//nolint:errcheck // test code
	rand.Read(addr[:])

	return *ton.NewAccountID(0, addr)
}

func internalMessageInfo(info *intMsgInfo) tlb.CommonMsgInfo {
	return tlb.CommonMsgInfo{
		SumType: "IntMsgInfo",
		IntMsgInfo: (*struct {
			IhrDisabled bool
			Bounce      bool
			Bounced     bool
			Src         tlb.MsgAddress
			Dest        tlb.MsgAddress
			Value       tlb.CurrencyCollection
			IhrFee      tlb.Grams
			FwdFee      tlb.Grams
			CreatedLt   uint64
			CreatedAt   uint32
		})(info),
	}
}

func externalMessageInfo(dest ton.AccountID) tlb.CommonMsgInfo {
	ext := struct {
		Src       tlb.MsgAddress
		Dest      tlb.MsgAddress
		ImportFee tlb.VarUInteger16
	}{
		Src:       tlb.MsgAddress{SumType: "AddrNone"},
		Dest:      dest.ToMsgAddress(),
		ImportFee: tlb.VarUInteger16{},
	}

	return tlb.CommonMsgInfo{SumType: "ExtInMsgInfo", ExtInMsgInfo: &ext}
}

func tonBlockID(now time.Time) ton.BlockIDExt {
	// simulate shard seqno as unix timestamp
	seqno := uint32(now.Unix())

	return ton.BlockIDExt{
		BlockID: ton.BlockID{
			Workchain: tonWorkchainID,
			Shard:     tonShardID,
			Seqno:     seqno,
		},
	}
}

func fakeDepositAmount(v math.Uint) tlb.Grams {
	return tlb.Grams(v.Uint64() + tonSampleTxFee)
}

func depositLogMock(t *testing.T, depositAmount, txFee uint64) *boc.Cell {
	//  cell log = begin_cell()
	//    .store_coins(deposit_amount)
	//    .store_coins(tx_fee)
	//    .end_cell();

	b := boc.NewCell()

	require.NoError(t, tlb.Grams(depositAmount).MarshalTLB(b, nil))
	require.NoError(t, tlb.Grams(txFee).MarshalTLB(b, nil))

	return b
}

// well, tlb.Transaction has unexported field `hash` that we need to set OUTSIDE tlb package.
// It's a hack, but it works for testing purposes.
func setTXHash(tx *tlb.Transaction, hash [32]byte) {
	field := reflect.ValueOf(tx).Elem().FieldByName("hash")
	ptr := unsafe.Pointer(field.UnsafeAddr())

	arrPtr := (*[32]byte)(ptr)
	*arrPtr = hash
}
