package sui

import (
	"fmt"
	"slices"
	"strconv"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
)

const (
	// FuncWithdrawImpl is the gateway function name withdraw_impl
	FuncWithdrawImpl = "withdraw_impl"

	// FuncIssueMessageContext is the gateway function name issue_message_context
	FuncIssueMessageContext = "issue_message_context"

	// FuncSetMessageContext is the gateway function name set_message_context
	FuncSetMessageContext = "set_message_context"

	// FuncResetMessageContext is the gateway function name reset_message_context
	FuncResetMessageContext = "reset_message_context"

	// ModuleConnected is the Sui connected module name
	ModuleConnected = "connected"

	// FuncOnCall is the Sui connected module function name on_call
	FuncOnCall = "on_call"

	// typeSeparator is the separator for Sui package and module
	typeSeparator = "::"

	// ptbWithdrawAndCallCmdCount is the number of commands in the PTB withdraw and call
	// the five commands are: [withdraw_impl, transfer_objects, set_message_context, on_call, reset_message_context]
	ptbWithdrawAndCallCmdCount = 5

	// ptbWithdrawImplInputCount is the number of inputs in the withdraw_impl command
	// the inputs are: [gatewayObject, amount, nonce, gasBudget, withdrawCap]
	ptbWithdrawImplInputCount = 5
)

// ptbWithdrawImplArgIndexes is the indexes of the inputs for the withdraw_impl command
// these are the corresponding indexes for arguments: [gatewayObject, amount, nonce, gasBudget, withdrawCap]
// the withdraw_impl command is the first command in the PTB, so the indexes will always be [0, 1, 2, 3, 4]
var ptbWithdrawImplArgIndexes = []int{0, 1, 2, 3, 4}

// MoveCall represents a Sui Move call with package ID, module and function
type MoveCall struct {
	PackageID  string
	Module     string
	Function   string
	ArgIndexes []int
}

// WithdrawAndCallPTB represents data for a Sui withdraw and call event
type WithdrawAndCallPTB struct {
	MoveCall
	Amount math.Uint
	Nonce  uint64
}

// TokenAmount returns the amount of the withdraw and call
func (d WithdrawAndCallPTB) TokenAmount() math.Uint {
	return d.Amount
}

// TxNonce returns the nonce of the withdraw and call
func (d WithdrawAndCallPTB) TxNonce() uint64 {
	return d.Nonce
}

// ExtractInitialSharedVersion extracts the object initial shared version from the object data
//
// Objects referenced for on_call are shared objects, initial shared version is required to build
// the withdraw and call using PTB.
// see: https://docs.sui.io/concepts/transactions/prog-txn-blocks#inputs
func ExtractInitialSharedVersion(objData models.SuiObjectData) (uint64, error) {
	owner, ok := objData.Owner.(map[string]any)
	if !ok {
		return 0, fmt.Errorf("invalid object owner type %T", objData.Owner)
	}

	shared, ok := owner["Shared"]
	if !ok {
		return 0, fmt.Errorf("missing shared object")
	}

	sharedMap, ok := shared.(map[string]any)
	if !ok {
		return 0, fmt.Errorf("invalid shared object type %T", shared)
	}

	return extractInteger[uint64](sharedMap, "initial_shared_version")
}

// parseWithdrawAndCallPTB parses withdraw and call with PTB.
// There is no actual event on gateway for withdraw and call, but we construct our own event to make the logic consistent.
func (gw *Gateway) parseWithdrawAndCallPTB(
	res models.SuiTransactionBlockResponse,
) (event Event, content OutboundEventContent, err error) {
	tx := res.Transaction.Data.Transaction

	// the number of PTB inputs should be >= 5
	if len(tx.Inputs) < ptbWithdrawImplInputCount {
		return event, nil, errors.Wrapf(
			ErrParseEvent,
			"invalid number of inputs(%d) in the PTB",
			len(tx.Inputs),
		)
	}

	// parse withdraw_impl at command 0
	moveCall, err := extractMoveCall(tx.Transactions[0])
	if err != nil {
		return event, nil, errors.Wrap(ErrParseEvent, "unable to parse withdraw_impl command in the PTB")
	}

	if moveCall.PackageID != gw.packageID {
		return event, nil, errors.Wrapf(ErrParseEvent, "invalid package id %s in the PTB", moveCall.PackageID)
	}

	if moveCall.Module != GatewayModule {
		return event, nil, errors.Wrapf(ErrParseEvent, "invalid module name %s in the PTB", moveCall.Module)
	}

	if moveCall.Function != FuncWithdrawImpl {
		return event, nil, errors.Wrapf(ErrParseEvent, "invalid function name %s in the PTB", moveCall.Function)
	}

	// ensure the argument indexes are matching the expected indexes
	if !slices.Equal(moveCall.ArgIndexes, ptbWithdrawImplArgIndexes) {
		return event, nil, errors.Wrapf(ErrParseEvent, "invalid argument indexes %v", moveCall.ArgIndexes)
	}

	// parse withdraw_impl arguments
	// argument1: amount
	amountStr, err := extractStr(tx.Inputs[1], "value")
	if err != nil {
		return Event{}, nil, errors.Wrap(ErrParseEvent, "unable to extract amount")
	}
	amount, err := math.ParseUint(amountStr)
	if err != nil {
		return Event{}, nil, errors.Wrap(ErrParseEvent, "unable to parse amount")
	}

	// argument2: nonce
	nonceStr, err := extractStr(tx.Inputs[2], "value")
	if err != nil {
		return Event{}, nil, errors.Wrap(ErrParseEvent, "unable to extract nonce")
	}
	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		return Event{}, nil, errors.Wrap(ErrParseEvent, "unable to parse nonce")
	}

	content = WithdrawAndCallPTB{
		MoveCall: moveCall,
		Amount:   amount,
		Nonce:    nonce,
	}

	event = Event{
		TxHash:     res.Digest,
		EventIndex: 0,
		EventType:  WithdrawAndCallEvent,
		content:    content,
	}

	return event, content, nil
}

// extractMoveCall extracts the MoveCall information from the PTB transaction command
func extractMoveCall(transaction any) (MoveCall, error) {
	commands, ok := transaction.(map[string]any)
	if !ok {
		return MoveCall{}, errors.Wrap(ErrParseEvent, "invalid command type")
	}

	// parse MoveCall info
	moveCall, ok := commands["MoveCall"].(map[string]any)
	if !ok {
		return MoveCall{}, errors.Wrap(ErrParseEvent, "missing MoveCall")
	}

	packageID, err := extractStr(moveCall, "package")
	if err != nil {
		return MoveCall{}, errors.Wrap(ErrParseEvent, "missing package ID")
	}

	module, err := extractStr(moveCall, "module")
	if err != nil {
		return MoveCall{}, errors.Wrap(ErrParseEvent, "missing module name")
	}

	function, err := extractStr(moveCall, "function")
	if err != nil {
		return MoveCall{}, errors.Wrap(ErrParseEvent, "missing function name")
	}

	// parse MoveCall data
	data, ok := moveCall["arguments"]
	if !ok {
		return MoveCall{}, errors.Wrap(ErrParseEvent, "missing arguments")
	}

	arguments, ok := data.([]any)
	if !ok {
		return MoveCall{}, errors.Wrap(ErrParseEvent, "arguments should be of slice type")
	}

	// extract MoveCall argument indexes
	argIndexes := make([]int, len(arguments))
	for i, arg := range arguments {
		indexes, ok := arg.(map[string]any)
		if !ok {
			return MoveCall{}, errors.Wrap(ErrParseEvent, "invalid argument type")
		}

		index, err := extractInteger[int](indexes, "Input")
		if err != nil {
			return MoveCall{}, errors.Wrap(ErrParseEvent, "missing argument index")
		}
		argIndexes[i] = index
	}

	return MoveCall{
		PackageID:  packageID,
		Module:     module,
		Function:   function,
		ArgIndexes: argIndexes,
	}, nil
}
