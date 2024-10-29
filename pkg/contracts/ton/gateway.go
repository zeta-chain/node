// Package ton provider bindings for TON blockchain including Gateway contract wrapper.
package ton

import (
	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

// Gateway represents bindings for Zeta Gateway contract on TON
//
// Gateway.ParseTransaction parses Gateway transaction.
// The parser reads tx body cell and decodes it based on Operation code (op)
//   - inbound transactions: deposit, donate, depositAndCall
//   - outbound transactions: not implemented yet
//   - errors for all other transactions
//
// `Send*` methods work the same way by constructing (& signing) tx body cell that is expected by the contract
//
// @see https://github.com/zeta-chain/protocol-contracts-ton/blob/main/contracts/gateway.fc
type Gateway struct {
	accountID ton.AccountID
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

// NewGateway Gateway constructor
func NewGateway(accountID ton.AccountID) *Gateway {
	return &Gateway{accountID}
}

// AccountID returns gateway address
func (gw *Gateway) AccountID() ton.AccountID {
	return gw.accountID
}

// ParseTransaction parses transaction to Transaction
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

// FilterInbounds filters transactions with deposit operations
func FilterInbounds(tx *Transaction) bool { return tx.IsInbound() }

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
		// #nosec G115 always in range
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

	var (
		bodyValue = boc.Cell(messages[0].Value.Body.Value)
		body      = &bodyValue
	)

	if err := body.Skip(sizeOpCode + sizeQueryID); err != nil {
		return depositLog{}, errors.Wrap(err, "unable to skip bits")
	}

	// skip msg address (ton sender)
	if err := UnmarshalTLB(&tlb.MsgAddress{}, body); err != nil {
		return depositLog{}, errors.Wrap(err, "unable to read sender address")
	}

	var deposited tlb.Grams
	if err := UnmarshalTLB(&deposited, body); err != nil {
		return depositLog{}, errors.Wrap(err, "unable to read deposited amount")
	}

	return depositLog{Amount: GramsToUint(deposited)}, nil
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

func (gw *Gateway) parseOutbound(_ ton.Transaction) (*Transaction, error) {
	return nil, errors.New("not implemented")
}

func parseInternalMessageBody(tx ton.Transaction) (*boc.Cell, error) {
	if !tx.Msgs.InMsg.Exists {
		return nil, errors.Wrap(ErrParse, "tx should have an internal message")
	}

	var (
		inMsg = tx.Msgs.InMsg.Value.Value
		body  = boc.Cell(inMsg.Body.Value)
	)

	return &body, nil
}
