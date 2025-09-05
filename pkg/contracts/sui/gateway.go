package sui

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
)

// EventType represents Gateway event type (both inbound & outbound)
type EventType string

// Gateway contains the API to read inbounds and sign outbounds to the Sui gateway
type Gateway struct {
	// packageID is the package ID of the gateway
	packageID string

	// gatewayObjectID is the object ID of the gateway struct
	objectID string

	mu sync.RWMutex
}

// OutboundEventContent is the interface of gateway outbound event content
type OutboundEventContent interface {
	// TokenAmount returns the amount of the outbound
	TokenAmount() math.Uint

	// TxNonce returns the nonce of the outbound
	TxNonce() uint64
}

// Event types
const (
	DepositEvent        EventType = "DepositEvent"
	DepositAndCallEvent EventType = "DepositAndCallEvent"
	WithdrawEvent       EventType = "WithdrawEvent"

	// this event does not exist on gateway, we define it to make the outbound processing consistent
	WithdrawAndCallEvent EventType = "WithdrawAndCallEvent"

	// the gateway.move uses name "NonceIncreaseEvent", but here uses a more descriptive name
	CancelTxEvent EventType = "NonceIncreaseEvent"
)

const GatewayModule = "gateway"

// ActiveMessageContextDynamicFieldName returns the dynamic field name of the active message context
func ActiveMessageContextDynamicFieldName() (json.RawMessage, error) {
	return dynamicFieldNameToJSONArray("active_message_context")
}

// NewGatewayFromPairID creates a new Sui Gateway
// from pair of `$packageID,$gatewayObjectID`
func NewGatewayFromPairID(pair string) (*Gateway, error) {
	packageID, gatewayObjectID, err := parsePair(pair)
	if err != nil {
		return nil, err
	}

	return NewGateway(packageID, gatewayObjectID), nil
}

// NewGateway creates a new Sui Gateway.
func NewGateway(packageID string, gatewayObjectID string) *Gateway {
	return &Gateway{packageID: packageID, objectID: gatewayObjectID}
}

// Event represents generic event wrapper
type Event struct {
	TxHash     string
	EventIndex uint64
	EventType  EventType

	content any
}

func (e *Event) IsDeposit() bool {
	return e.EventType == DepositEvent || e.EventType == DepositAndCallEvent
}

// Deposit extract DepositData.
func (e *Event) Deposit() (Deposit, error) {
	v, ok := e.content.(Deposit)
	if !ok {
		return Deposit{}, errors.Errorf("invalid content type %T", e.content)
	}

	return v, nil
}

func (e *Event) IsWithdraw() bool {
	return e.EventType == WithdrawEvent
}

// Withdrawal extract withdraw data.
func (e *Event) Withdrawal() (Withdrawal, error) {
	v, ok := e.content.(Withdrawal)
	if !ok {
		return Withdrawal{}, errors.Errorf("invalid content type %T", e.content)
	}

	return v, nil
}

func (e *Event) IsCancelTx() bool {
	return e.EventType == CancelTxEvent
}

// CancelTx extract cancel tx data.
func (e *Event) CancelTx() (CancelTx, error) {
	v, ok := e.content.(CancelTx)
	if !ok {
		return CancelTx{}, errors.Errorf("invalid content type %T", e.content)
	}

	return v, nil
}

// PackageID returns object id of Gateway code
func (gw *Gateway) PackageID() string {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	return gw.packageID
}

// ObjectID returns Gateway's struct object id
func (gw *Gateway) ObjectID() string {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	return gw.objectID
}

// WithdrawCapType returns struct type of the WithdrawCap
func (gw *Gateway) WithdrawCapType() string {
	return fmt.Sprintf("%s::%s::WithdrawCap", gw.PackageID(), GatewayModule)
}

// MessageContextType returns struct type of the MessageContext
func (gw *Gateway) MessageContextType() string {
	return fmt.Sprintf("%s::%s::MessageContext", gw.PackageID(), GatewayModule)
}

// UpdateIDs updates packageID and objectID.
func (gw *Gateway) UpdateIDs(pair string) error {
	packageID, gatewayObjectID, err := parsePair(pair)
	if err != nil {
		return err
	}
	gw.mu.Lock()
	defer gw.mu.Unlock()

	gw.packageID = packageID
	gw.objectID = gatewayObjectID

	return nil
}

// ParseEvent parses Event.
func (gw *Gateway) ParseEvent(event models.SuiEventResponse) (Event, error) {
	// basic validation
	switch {
	case event.Id.TxDigest == "":
		return Event{}, errors.Wrap(ErrParseEvent, "empty tx hash")
	case event.Id.EventSeq == "":
		return Event{}, errors.Wrap(ErrParseEvent, "empty event id")
	case event.PackageId != gw.packageID:
		return Event{}, errors.Wrapf(
			ErrParseEvent,
			"package id mismatch (got %s, want %s)",
			event.PackageId,
			gw.packageID,
		)
	}

	// Extract common fields
	txHash := event.Id.TxDigest
	eventID, err := strconv.ParseUint(event.Id.EventSeq, 10, 64)
	if err != nil {
		return Event{}, errors.Wrapf(ErrParseEvent, "failed to parse event id %q", event.Id.EventSeq)
	}

	descriptor, err := parseEventDescriptor(event.Type)
	if err != nil {
		return Event{}, errors.Wrap(ErrParseEvent, err.Error())
	}

	// Note that event.TransactionModule can be different because it represents
	// the module BY WHICH the gateway was called.
	if descriptor.module != GatewayModule {
		return Event{}, errors.Wrapf(ErrParseEvent, "module mismatch %q", descriptor.module)
	}

	var (
		eventType = descriptor.eventType
		content   any
	)

	// Parse specific events
	switch eventType {
	case DepositEvent, DepositAndCallEvent:
		content, err = parseDeposit(event, eventType)
	case WithdrawEvent:
		content, err = parseWithdrawal(event, eventType)
	case CancelTxEvent:
		content, err = parseCancelTx(event, eventType)
	default:
		return Event{}, errors.Wrapf(ErrParseEvent, "unknown event %q", eventType)
	}

	if err != nil {
		return Event{}, errors.Wrapf(ErrParseEvent, "%s: %s", eventType, err.Error())
	}

	return Event{
		TxHash:     txHash,
		EventIndex: eventID,
		EventType:  eventType,
		content:    content,
	}, nil
}

// ParseOutboundEvent parses outbound event from transaction block response.
func (gw *Gateway) ParseOutboundEvent(
	res models.SuiTransactionBlockResponse,
) (event Event, content OutboundEventContent, err error) {
	// a simple withdraw contains one single command, if it contains 5 commands,
	// we try passing the transaction as a withdraw and call with PTB
	if len(res.Transaction.Data.Transaction.Transactions) == ptbWithdrawAndCallCmdCount {
		return gw.parseWithdrawAndCallPTB(res)
	}

	if len(res.Events) == 0 {
		return event, nil, errors.New("missing events")
	}

	event, err = gw.ParseEvent(res.Events[0])
	if err != nil {
		return event, nil, errors.Wrap(err, "unable to parse event")
	}

	// filter outbound events
	switch event.EventType {
	case WithdrawEvent:
		withdrawal, err := event.Withdrawal()
		if err != nil {
			return event, nil, errors.Wrap(err, "unable to extract withdraw event")
		}
		return event, withdrawal, nil
	case CancelTxEvent:
		cancelTx, err := event.CancelTx()
		if err != nil {
			return event, nil, errors.Wrap(err, "unable to extract cancel tx event")
		}
		return event, cancelTx, nil
	default:
		return event, nil, errors.Errorf("unsupported outbound event type %s", event.EventType)
	}
}

// ParseTxWithdrawal a syntax sugar around ParseEvent and Withdrawal.
func (gw *Gateway) ParseTxWithdrawal(tx models.SuiTransactionBlockResponse) (event Event, w Withdrawal, err error) {
	if len(tx.Events) == 0 {
		err = errors.New("missing events")
		return event, w, err
	}

	event, err = gw.ParseEvent(tx.Events[0])
	if err != nil {
		return event, w, err
	}

	if !event.IsWithdraw() {
		err = errors.Errorf("invalid event type %s", event.EventType)
		return event, w, err
	}

	w, err = event.Withdrawal()
	if err != nil {
		return event, w, err
	}

	return event, w, err
}

// ParseDynamicFieldValueStr parses the dynamic field's value from object data as string
func ParseDynamicFieldValueStr(data models.SuiParsedData) (string, error) {
	// dynamic field object contains 3 fields: id, name, value
	// the 'value' is what the dynamic field stores
	rawValue, ok := data.Fields["value"]
	if !ok {
		return "", errors.New("missing value field")
	}

	value, ok := rawValue.(string)
	if !ok {
		return "", errors.Errorf("want string, got %T for dynamic field value", rawValue)
	}

	return value, nil
}

// ParseGatewayNonce parses gateway nonce from event.
func ParseGatewayNonce(data models.SuiParsedData) (uint64, error) {
	fields := data.Fields

	// extract nonce field from the object content
	rawNonce, ok := fields["nonce"]
	if !ok {
		return 0, errors.New("missing nonce field")
	}

	v, ok := rawNonce.(string)
	if !ok {
		return 0, errors.Errorf("want string, got %T for nonce", rawNonce)
	}

	// #nosec G115 always in range
	nonce, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "unable to parse nonce")
	}

	return nonce, nil
}

type eventDescriptor struct {
	packageID string
	module    string
	eventType EventType
}

func parseEventDescriptor(typeString string) (eventDescriptor, error) {
	parts := strings.Split(typeString, "::")
	if len(parts) != 3 {
		return eventDescriptor{}, errors.Errorf("invalid event type %q", typeString)
	}

	return eventDescriptor{
		packageID: parts[0],
		module:    parts[1],
		eventType: EventType(parts[2]),
	}, nil
}

func extractStr(kv map[string]any, key string) (string, error) {
	if _, ok := kv[key]; !ok {
		return "", errors.Errorf("missing %s", key)
	}

	v, ok := kv[key].(string)
	if !ok {
		return "", errors.Errorf("invalid %s", key)
	}

	return v, nil
}

// extractInteger extracts a float64 value from a map and converts it to any integer type
func extractInteger[T constraints.Integer](kv map[string]any, key string) (T, error) {
	rawValue, ok := kv[key]
	if !ok {
		return 0, errors.Errorf("missing %s", key)
	}

	v, ok := rawValue.(float64)
	if !ok {
		return 0, errors.Errorf("want float64, got %T for %s", rawValue, key)
	}

	// #nosec G115 always in range
	return T(v), nil
}

func convertPayload(data []any) ([]byte, error) {
	payload := make([]byte, len(data))

	for idx, something := range data {
		// parsed bytes are represented as float64
		f, ok := something.(float64)
		switch {
		case !ok:
			return nil, errors.Errorf("not a float64 at index %d", idx)
		case f < 0 || f > 255:
			return nil, errors.Errorf("not a byte (%f) at index %d", f, idx)
		default:
			payload[idx] = byte(f)
		}
	}

	// Try decoding the payload from base64
	// If the payload is not base64 encoded, directly return the payload bytes as is
	// Currently the localnet Sui RPC will return bytes in Base64 for the payload data while live network return the actual bytes of the payload directly
	// TODO: fix this discrepancy
	// https://github.com/zeta-chain/node/issues/3919
	base64DecodedPayload, err := base64.StdEncoding.DecodeString(string(payload))
	if err == nil {
		return base64DecodedPayload, nil
	}
	return payload, nil
}

func parsePair(pair string) (string, string, error) {
	parts := strings.Split(pair, ",")
	if len(parts) != 2 {
		return "", "", errors.Errorf("invalid pair %q", pair)
	}

	// each part should be a valid Sui address
	for _, part := range parts {
		if err := ValidateAddress(part); err != nil {
			return "", "", errors.Wrapf(err, "invalid Sui address %q", part)
		}
	}

	return parts[0], parts[1], nil
}

// dynamicFieldNameToJSONArray converts a string dynamic field name to a JSON array of integer values
//
// This conversion is necessary when interacting with Sui Move functions that expect vector<u8> parameters.
// In Sui's JSON-RPC API, byte arrays (vector<u8>) must be passed as JSON arrays of integers representing each byte's numeric value.
//
// For example:
//   - Input string: "active_message_context"
//   - Output JSON: [97,99,116,105,118,101,95,109,101,115,115,97,103,101,95,99,111,110,116,101,120,116]
//
// But Go's json.Marshal([]byte) produces base64 strings, not integer arrays, which is not accepted by Sui JSON-RPC API
// For example:
//   - Input string: "active_message_context"
//   - Output JSON: "YWN0aXZlX21lc3NhZ2VfY29udGV4dA=="
func dynamicFieldNameToJSONArray(name string) (json.RawMessage, error) {
	bytes := []byte(name)

	// convert byte slice to int array
	intArray := make([]int, len(bytes))
	for i, b := range bytes {
		intArray[i] = int(b)
	}

	// marshal the int array to JSON
	jsonBytes, err := json.Marshal(intArray)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(jsonBytes), nil
}
