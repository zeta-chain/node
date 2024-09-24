package ton

import (
	"bytes"

	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

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

type Gateway struct {
	accountID ton.AccountID
}

type Donation struct {
	Sender ton.AccountID
	Amount math.Uint
}

type Deposit struct {
	Sender    ton.AccountID
	Amount    math.Uint
	Recipient eth.Address
}

type DepositAndCall struct {
	Deposit
	CallData []byte
}

const (
	sizeOpCode  = 32
	sizeQueryID = 64
)

var (
	ErrParse     = errors.New("unable to parse tx")
	ErrUnknownOp = errors.New("unknown op")
	ErrCast      = errors.New("unable to cast tx content")
)

func NewGateway(accountID ton.AccountID) *Gateway {
	return &Gateway{accountID}
}

func (gw *Gateway) ParseTransaction(tx ton.Transaction) (*Transaction, error) {
	if !tx.IsSuccess() {
		exitCode := tx.Description.TransOrd.ComputePh.TrPhaseComputeVm.Vm.ExitCode
		return nil, errors.Wrapf(ErrParse, "tx %s is not successful (exit code %d)", tx.Hash().Hex(), exitCode)
	}

	if tx.Msgs.InMsg.Exists {
		inbound, err := gw.parseInbound(tx)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse inbound tx %s", tx.Hash().Hex())
		}

		return inbound, nil
	}

	outbound, err := gw.parseOutbound(tx)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse outbound tx %s", tx.Hash().Hex())
	}

	return outbound, nil
}

// ParseAndFilter parses transaction and applies filter to it. Returns (tx, skip?, error)
// If parse fails due to known error, skip is set to true
func (gw *Gateway) ParseAndFilter(tx ton.Transaction, filter func(*Transaction) bool) (*Transaction, bool, error) {
	parsedTX, err := gw.ParseTransaction(tx)
	switch {
	case errors.Is(err, ErrParse):
		return nil, true, nil
	case errors.Is(err, ErrUnknownOp):
		return nil, true, nil
	case err != nil:
		return nil, false, err
	}

	if !filter(parsedTX) {
		return nil, true, nil
	}

	return parsedTX, false, nil
}

// FilterDeposit filters transactions with deposit operations
func FilterDeposit(tx *Transaction) bool {
	return tx.Operation == OpDeposit || tx.Operation == OpDepositAndCall
}

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
		return nil, errors.Wrap(err, "unable to read op code")
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
		content = Donation{Sender: sender, Amount: gramToUint(amount)}
	case OpDeposit:
		content, errContent = parseDeposit(tx, sender, body)
	case OpDepositAndCall:
		content, errContent = parseDepositAndCall(tx, sender, body)
	default:
		return nil, errors.Wrapf(ErrUnknownOp, "op code %d", int64(op))
	}

	if errContent != nil {
		return nil, errors.Wrapf(ErrParse, "unable to parse content for op code %d: %s", int64(op), errContent.Error())
	}

	return &Transaction{
		Transaction: tx,
		Operation:   opCode,

		content: content,
		inbound: true,
	}, nil
}

func parseDeposit(tx ton.Transaction, sender ton.AccountID, body *boc.Cell) (Deposit, error) {
	// skip query id
	if err := body.Skip(sizeQueryID); err != nil {
		return Deposit{}, err
	}

	recipient, err := readEVMAddress(body)
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
	Amount math.Uint
}

func parseDepositLog(tx ton.Transaction) (depositLog, error) {
	messages := tx.Msgs.OutMsgs.Values()
	if len(messages) == 0 {
		return depositLog{}, errors.Wrap(ErrParse, "no out messages")
	}

	// stored as ref
	// cell log = begin_cell()
	//        .store_uint(op::internal::deposit, size::op_code_size)
	//        .store_uint(0, size::query_id_size)
	//        .store_slice(sender)
	//        .store_coins(deposit_amount)
	//        .store_uint(evm_recipient, size::evm_address)
	//        .end_cell();

	body, err := marshalCellRef(messages[0].Value.Body)
	if err != nil {
		return depositLog{}, errors.Wrap(err, "unable to read body cell")
	}

	if err := body.Skip(sizeOpCode + sizeQueryID); err != nil {
		return depositLog{}, errors.Wrap(err, "unable to skip bits")
	}

	// skip msg address (ton sender)
	if err := unmarshalTLB(&tlb.MsgAddress{}, body); err != nil {
		return depositLog{}, errors.Wrap(err, "unable to read sender address")
	}

	var deposited tlb.Grams
	if err := unmarshalTLB(&deposited, body); err != nil {
		return depositLog{}, errors.Wrap(err, "unable to read deposited amount")
	}

	return depositLog{Amount: gramToUint(deposited)}, nil
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

	var sd tlb.SnakeData
	if err = unmarshalTLB(&sd, callDataCell); err != nil {
		return DepositAndCall{}, errors.Wrap(err, "unable to unmarshal call data")
	}

	cd := boc.BitString(sd)

	// TLB operates with bits, so we (might) need to trim some "leftovers" (null chars)
	callData := bytes.Trim(cd.Buffer(), "\x00")

	return DepositAndCall{Deposit: deposit, CallData: callData}, nil
}

func (gw *Gateway) parseOutbound(_ ton.Transaction) (*Transaction, error) {
	return nil, errors.New("not implemented")
}

func parseInternalMessageBody(tx ton.Transaction) (*boc.Cell, error) {
	if !tx.Msgs.InMsg.Exists {
		return nil, errors.Wrap(ErrParse, "tx should have an internal message")
	}

	either := tx.Msgs.InMsg.Value.Value.Body
	if either.IsRight {
		return nil, errors.Wrap(ErrParse, "tx body should not be a Ref")
	}

	body, err := marshalCell(&either.Value)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read body cell")
	}

	return body, nil
}

func marshalCell(v tlb.MarshalerTLB) (*boc.Cell, error) {
	cell := boc.NewCell()

	if err := v.MarshalTLB(cell, &tlb.Encoder{}); err != nil {
		return nil, err
	}

	return cell, nil
}

func marshalCellRef(v tlb.MarshalerTLB) (*boc.Cell, error) {
	c, err := marshalCell(v)
	if err != nil {
		return nil, err
	}

	c, err = c.NextRef()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create ref cell")
	}

	return c, nil
}

func unmarshalTLB(t tlb.UnmarshalerTLB, cell *boc.Cell) error {
	return t.UnmarshalTLB(cell, &tlb.Decoder{})
}

func gramToUint(g tlb.Grams) math.Uint {
	return math.NewUint(uint64(g))
}

func readEVMAddress(cell *boc.Cell) (eth.Address, error) {
	const evmAddrBits = 20 * 8

	s, err := cell.ReadBits(evmAddrBits)
	if err != nil {
		return eth.Address{}, err
	}

	return eth.BytesToAddress(s.Buffer()), nil
}
