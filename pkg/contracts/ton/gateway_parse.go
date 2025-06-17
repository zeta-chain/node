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
		// #nosec G115 always in range
		opCode = Op(op)
		sender = *sourceID

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
	case OpCall:
		content, errContent = parseCall(tx, sender, body)
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
		ExitCode:    exitCodeFromTx(tx),

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

func parseCall(_ ton.Transaction, sender ton.AccountID, body *boc.Cell) (Call, error) {
	// skip query id
	if err := body.Skip(sizeQueryID); err != nil {
		return Call{}, err
	}

	recipient, err := UnmarshalEVMAddress(body)
	if err != nil {
		return Call{}, errors.Wrap(err, "unable to read recipient")
	}

	callDataCell, err := body.NextRef()
	if err != nil {
		return Call{}, errors.Wrap(err, "unable to read call data cell")
	}

	callData, err := UnmarshalSnakeCell(callDataCell)
	if err != nil {
		return Call{}, errors.Wrap(err, "unable to unmarshal call data")
	}

	return Call{
		Sender:    sender,
		Recipient: recipient,
		CallData:  callData,
	}, nil
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

	// #nosec G115 always in range
	opCode := Op(op)
	content := any(nil)

	switch opCode {
	case OpWithdraw:
		content, err = parseWithdrawal(tx, sig, payload)
		if err != nil {
			return nil, errParse(err, "unable to parse 'withdraw'")
		}
	case OpIncreaseSeqno:
		content, err = parseIncreaseSeqno(tx, sig, payload)
		if err != nil {
			return nil, errParse(err, "unable to parse 'increase seqno'")
		}
	default:
		return nil, errors.Wrapf(ErrUnknownOp, "op code %d", op)
	}

	return &Transaction{
		Transaction: tx,
		Operation:   opCode,
		ExitCode:    exitCodeFromTx(tx),
		content:     content,
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

func parseWithdrawal(tx ton.Transaction, sig [65]byte, payload *boc.Cell) (Withdrawal, error) {
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

	recipientAddr, err := parseAccount(recipient)
	if err != nil {
		return Withdrawal{}, errors.Wrap(err, "unable to parse recipient from payload")
	}

	// ensure a single outgoing message for the withdrawal
	if tx.OutMsgCnt != 1 {
		return Withdrawal{}, errors.Wrap(ErrParse, "invalid out messages count")
	}

	// tlb.Message{}
	outMsg := tx.Msgs.OutMsgs.Values()[0].Value
	if outMsg.Info.SumType != "IntMsgInfo" || outMsg.Info.IntMsgInfo == nil {
		return Withdrawal{}, errors.Wrap(ErrParse, "invalid out message")
	}

	msgRecipientAddr, err := parseAccount(outMsg.Info.IntMsgInfo.Dest)

	switch {
	case err != nil:
		return Withdrawal{}, errors.Wrap(err, "unable to parse recipient from out msg")
	case recipientAddr != msgRecipientAddr:
		// should not happen
		return Withdrawal{}, errors.Wrap(ErrParse, "recipient mismatch")
	case amount != outMsg.Info.IntMsgInfo.Value.Grams:
		// should not happen
		return Withdrawal{}, errors.Wrap(ErrParse, "amount mismatch")
	}

	return Withdrawal{
		Recipient: recipientAddr,
		Amount:    math.NewUint(uint64(amount)),
		Seqno:     seqno,
		Sig:       shiftSignature(sig),
	}, nil
}

func parseIncreaseSeqno(_ ton.Transaction, sig [65]byte, payload *boc.Cell) (IncreaseSeqno, error) {
	reasonCode, err := payload.ReadUint(sizeSeqno)
	if err != nil {
		return IncreaseSeqno{}, errors.Wrap(err, "unable to read reason code")
	}

	seqno, err := payload.ReadUint(sizeSeqno)
	if err != nil {
		return IncreaseSeqno{}, errors.Wrap(err, "unable to read seqno")
	}

	// #nosec G115 always in range
	return IncreaseSeqno{
		Seqno:      uint32(seqno),
		ReasonCode: uint32(reasonCode),
		Sig:        shiftSignature(sig),
	}, nil
}

func parseAccount(raw tlb.MsgAddress) (ton.AccountID, error) {
	if raw.SumType != "AddrStd" {
		return ton.AccountID{}, errors.Wrapf(ErrParse, "invalid address type %s", raw.SumType)
	}

	return ton.AccountID{
		// #nosec G115 always in range
		Workchain: int32(raw.AddrStd.WorkchainId),
		Address:   raw.AddrStd.Address,
	}, nil
}

func errParse(err error, msg string) error {
	return errors.Wrapf(ErrParse, "%s (%s)", msg, err.Error())
}

func exitCodeFromTx(tx ton.Transaction) int32 {
	return tx.Description.TransOrd.ComputePh.TrPhaseComputeVm.Vm.ExitCode
}

// shiftSignature: shifts bytes: ECDSA sig has the following order: (v, r, s) but in EVM we have (r, s, v)
func shiftSignature(sig [65]byte) [65]byte {
	var sigFlipped [65]byte

	copy(sigFlipped[:64], sig[1:])
	sigFlipped[64] = sig[0]

	return sigFlipped
}
