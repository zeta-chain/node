package ton

import (
	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func (gw *Gateway) parseInbound(tx ton.Transaction) (*Transaction, error) {
	body, err := parseInternalMessageBody(tx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse body")
	}

	intMsgInfo := tx.Msgs.InMsg.Value.Value.Info.IntMsgInfo
	if intMsgInfo == nil {
		return nil, errors.Wrap(ErrParse, "no internal message info")
	}

	sourceID, err := ton.AccountIDFromTlb(intMsgInfo.Src)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse source account")
	}

	destinationID, err := ton.AccountIDFromTlb(intMsgInfo.Dest)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse destination account")
	}

	if gw.accountID != *destinationID {
		return nil, errors.Wrap(ErrParse, "destination account is not gateway")
	}

	op, err := body.ReadUint(sizeOpCode)
	if err != nil {
		return nil, errParse(err, "unable to read op code")
	}

	var (
		sender = *sourceID
		opCode = Op(op)

		content    any
		errContent error
	)

	switch opCode {
	case OpDonate:
		amount := intMsgInfo.Value.Grams - tx.TotalFees.Grams
		content = Donation{Sender: sender, Amount: GramsToUint(amount)}
	case OpDeposit:
		content, errContent = parseDeposit(tx, sender, body)
	case OpDepositAndCall:
		content, errContent = parseDepositAndCall(tx, sender, body)
	default:
		// #nosec G115 always in range
		return nil, errors.Wrapf(ErrUnknownOp, "op code %d", int64(op))
	}

	if errContent != nil {
		// #nosec G115 always in range
		return nil, errors.Wrapf(ErrParse, "unable to parse content for op code %d: %s", int64(op), errContent.Error())
	}

	return &Transaction{
		Transaction: tx,
		Operation:   opCode,
		ExitCode:    exitCode(tx),

		content: content,
		inbound: true,
	}, nil
}

func parseDeposit(tx ton.Transaction, sender ton.AccountID, body *boc.Cell) (Deposit, error) {
	// skip query id
	if err := body.Skip(sizeQueryID); err != nil {
		return Deposit{}, err
	}

	recipient, err := UnmarshalEVMAddress(body)
	if err != nil {
		return Deposit{}, errors.Wrap(err, "unable to read recipient")
	}

	dl, err := parseDepositLog(tx)
	if err != nil {
		return Deposit{}, errors.Wrap(err, "unable to parse deposit log")
	}

	return Deposit{
		Sender:    sender,
		Amount:    dl.Amount,
		Recipient: recipient,
	}, nil
}

type depositLog struct {
	Amount     math.Uint
	DepositFee math.Uint
}

func parseDepositLog(tx ton.Transaction) (depositLog, error) {
	messages := tx.Msgs.OutMsgs.Values()
	if len(messages) == 0 {
		return depositLog{}, errors.Wrap(ErrParse, "no out messages")
	}

	// stored as ref
	// cell log = begin_cell()
	//     .store_coins(deposit_amount)
	//     .store_coins(tx_fee)
	//     .end_cell();

	var (
		bodyValue = boc.Cell(messages[0].Value.Body.Value)
		body      = &bodyValue
	)

	var deposited tlb.Grams
	if err := UnmarshalTLB(&deposited, body); err != nil {
		return depositLog{}, errors.Wrap(err, "unable to read deposited amount")
	}

	var depositFee tlb.Grams
	if err := UnmarshalTLB(&depositFee, body); err != nil {
		return depositLog{}, errors.Wrap(err, "unable to read deposit fee")
	}

	return depositLog{
		Amount:     GramsToUint(deposited),
		DepositFee: GramsToUint(depositFee),
	}, nil
}

func parseDepositAndCall(tx ton.Transaction, sender ton.AccountID, body *boc.Cell) (DepositAndCall, error) {
	deposit, err := parseDeposit(tx, sender, body)
	if err != nil {
		return DepositAndCall{}, err
	}

	callDataCell, err := body.NextRef()
	if err != nil {
		return DepositAndCall{}, errors.Wrap(err, "unable to read call data cell")
	}

	callData, err := UnmarshalSnakeCell(callDataCell)
	if err != nil {
		return DepositAndCall{}, errors.Wrap(err, "unable to unmarshal call data")
	}

	return DepositAndCall{Deposit: deposit, CallData: callData}, nil
}

// an outbound is a tx that was initiated by TSS signature with external message
func isOutbound(tx ton.Transaction) bool {
	return tx.Msgs.InMsg.Exists &&
		tx.Msgs.InMsg.Value.Value.Info.SumType == "ExtInMsgInfo"
}

func (gw *Gateway) parseOutbound(tx ton.Transaction) (*Transaction, error) {
	if !isOutbound(tx) {
		return nil, errors.Wrap(ErrParse, "not an outbound transaction")
	}

	extMsgBodyCell := boc.Cell(tx.Msgs.InMsg.Value.Value.Body.Value)

	sig, payload, err := parseExternalMessage(&extMsgBodyCell)
	if err != nil {
		return nil, errParse(err, "unable to parse external message")
	}

	op, err := payload.ReadUint(sizeOpCode)
	if err != nil {
		return nil, errParse(err, "unable to read op code")
	}

	opCode := Op(op)

	if opCode != OpWithdraw {
		return nil, errors.Wrapf(ErrUnknownOp, "op code %d", op)
	}

	withdrawal, err := parseWithdrawal(sig, payload)
	if err != nil {
		return nil, errParse(err, "unable to parse withdrawal")
	}

	return &Transaction{
		Transaction: tx,
		Operation:   opCode,
		ExitCode:    exitCode(tx),
		content:     withdrawal,
	}, nil
}

// external message is essentially a cell with 65 bytes of ECDSA sig + cell_ref to payload
func parseExternalMessage(b *boc.Cell) ([65]byte, *boc.Cell, error) {
	sig, err := b.ReadBytes(65)
	if err != nil {
		return [65]byte{}, nil, err
	}

	var sigArray [65]byte
	copy(sigArray[:], sig)

	payload, err := b.NextRef()

	return sigArray, payload, err
}

func parseWithdrawal(sig [65]byte, payload *boc.Cell) (Withdrawal, error) {
	// Note that ECDSA sig has the following order: (v, r, s) but in EVM we have (r, s, v)
	var sigFlipped [65]byte

	copy(sigFlipped[:64], sig[1:])
	sigFlipped[64] = sig[0]

	var (
		recipient tlb.MsgAddress
		amount    tlb.Coins
		seqno     uint32
	)

	err := ErrCollect(
		tlb.Unmarshal(payload, &recipient),
		tlb.Unmarshal(payload, &amount),
		tlb.Unmarshal(payload, &seqno),
	)
	if err != nil {
		return Withdrawal{}, errors.Wrap(err, "unable to unmarshal payload")
	}

	return Withdrawal{
		Recipient: ton.AccountID{
			Workchain: int32(recipient.AddrStd.WorkchainId),
			Address:   recipient.AddrStd.Address,
		},
		Amount: math.NewUint(uint64(amount)),
		Seqno:  seqno,
		Sig:    sigFlipped,
	}, nil
}

func errParse(err error, msg string) error {
	return errors.Wrapf(ErrParse, msg+" (%s)", err.Error())
}

func exitCode(tx ton.Transaction) int32 {
	return tx.Description.TransOrd.ComputePh.TrPhaseComputeVm.Vm.ExitCode
}
