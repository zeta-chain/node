package sui

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
)

const (
	// gatewayPairIDPartsV1 is the number of parts in the gateway pair ID for v1
	gatewayPairIDPartsV1 = 2

	// gatewayPairIDPartsV2 is the number of parts in the gateway pair ID for v2 and later
	gatewayPairIDPartsV2 = 5

	// previousPackageIDIndex is the index of the previous package ID in the pair ID for v2 and later
	previousPackageIDIndex = 3
)

// EventType represents Gateway event type (both inbound & outbound)
type EventType string

// Gateway contains the API to read inbounds and sign outbounds to the Sui gateway
type Gateway struct {
	// packageID is the package ID of the latest gateway
	// For example, upon gateway upgrade v1 -> v2, the packageID will point to the v2 gateway package.
	packageID string

	// gatewayObjectID is the object ID of the gateway struct
	// The gateway object ID will remain unchanged across upgrades.
	objectID string

	// withdrawCapID was introduced in the authenticated call gateway upgrade (v2).
	// It explicitly specifies the withdrawCap object ID to avoid uncertainty across upgrades.
	withdrawCapID string

	// previousPackageID was introduced in the authenticated call gateway upgrade (v2).
	// previousPackageID is an optional field (can be empty) that points to the previous gateway package.
	// To achieve seamless upgrade, the protocol needs to know the previous package ID and continue to
	// support it for a period of time (so users have time to migrate) before fully deprecating it.
	//
	// For example:
	//  - on upgrade v1 -> v2, the previous package ID is v1 (same as originalPackageID)
	//  - on upgrade v2 -> v3, the previous package ID is v2
	//
	// To deprecate previous package, this field must be set to empty string in the chain params' gateway address.
	previousPackageID string

	// originalPackageID was introduced in the authenticated call gateway upgrade (v2).
	// originalPackageID is the original (v1) gateway package ID and will remain unchanged across upgrades.
	// The reason we need this field is because all the events were initially defined in the original gateway
	// package, so the observers MUST pass this packageID to Sui RPC 'QueryModuleEvents' to query events.
	originalPackageID string

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

// NewGatewayFromPairID creates a new Sui Gateway struct from the chain params gateway address,
// which is a pair ID of `$packageID,$gatewayObjectID[,$withdrawCapID,previousPackageID,originalPackageID]`
//
// An example of 5-part gateway address string in the Sui chain params looks like:
/*
"0x5d4b302506645c37ff133b98fff50a5ae14841659738d6d733d59d0d217a9fff,
0xba477ad7b87a31fde3d29c4e4512329d7340ec23e61f130ebb4d0169ba37e189,
0x4a367f98d9299019e3d5bbc6ee1d41b5789172e14f7c63a881377766902438e2,
0x1db3a54b99c2741bf8b8aaa8266d6e7b6daf0c702a5ef5b0d6e9e6cf12527a90,
0x1db3a54b99c2741bf8b8aaa8266d6e7b6daf0c702a5ef5b0d6e9e6cf12527a90"
*/
func NewGatewayFromPairID(pair string) (*Gateway, error) {
	packageID, gatewayObjectID, withdrawCapID, previousPackageID, originalPackageID, err := parsePair(pair)
	if err != nil {
		return nil, err
	}

	return &Gateway{
		packageID:         packageID,
		objectID:          gatewayObjectID,
		withdrawCapID:     withdrawCapID,
		previousPackageID: previousPackageID,
		originalPackageID: originalPackageID,
	}, nil
}

// NewGateway creates a new Sui Gateway.
func NewGateway(packageID string, gatewayObjectID string) *Gateway {
	return &Gateway{packageID: packageID, objectID: gatewayObjectID}
}

// MakePairID makes a pair ID of the form `$packageID,$gatewayObjectID[,$withdrawCapID,previousPackageID,originalPackageID]`
// Note: It is only used for tests at the moment.
func MakePairID(packageID, gatewayObjectID, withdrawCapID, previousPackageID, originalPackageID string) string {
	if withdrawCapID == "" || originalPackageID == "" {
		return fmt.Sprintf("%s,%s", packageID, gatewayObjectID)
	}
	return fmt.Sprintf(
		"%s,%s,%s,%s,%s",
		packageID,
		gatewayObjectID,
		withdrawCapID,
		previousPackageID,
		originalPackageID,
	)
}

// ToPairID return a pair ID of `$packageID,$gatewayObjectID[,$withdrawCapID,previousPackageID,originalPackageID]`
// Note: It is only used for tests at the moment.
func (gw *Gateway) ToPairID() string {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	return MakePairID(gw.packageID, gw.objectID, gw.withdrawCapID, gw.previousPackageID, gw.originalPackageID)
}

// Previous creates a Gateway struct that points to the previous gateway.
// Note: this method is not used for now, but we keep it for future use.
func (gw *Gateway) Previous() *Gateway {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	// previous package ID does not exist, return nil
	if gw.previousPackageID == "" {
		return nil
	}

	return &Gateway{packageID: gw.previousPackageID, objectID: gw.objectID, withdrawCapID: gw.withdrawCapID}
}

// Original creates a Gateway struct that points to the original gateway.
//
// Note:
//   - Gateway events were defined in the original gateway package, so the original package ID should be used for event queries.
//     Event queries on upgraded gateway package ID will return empty events and lead to missed deposits.
//   - This method allows the observer to make a switch and work with the original package after upgrade.
func (gw *Gateway) Original() *Gateway {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	// return self if original package ID is not specified
	if gw.originalPackageID == "" {
		return gw
	}
	return &Gateway{packageID: gw.originalPackageID, objectID: gw.objectID, withdrawCapID: gw.withdrawCapID}
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

// SupportedPackageIDs returns slice of supported package IDs from which emitted events can be observed
//
// There are two cases:
//   - There is only one supported package ID, which is the current package ID.
//   - There are two supported package IDs, which are the current package ID and previous package ID.
//     This happens during gateway upgrades before fully deprecating the previous package.
//
// Note: empty previous package ID means the previous package is deprecated.
func (gw *Gateway) SupportedPackageIDs() []string {
	gw.mu.RLock()
	defer gw.mu.RUnlock()

	if gw.previousPackageID == "" {
		return []string{gw.packageID}
	}
	return []string{gw.packageID, gw.previousPackageID}
}

// ObjectID returns Gateway's struct object id
func (gw *Gateway) ObjectID() string {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	return gw.objectID
}

// WithdrawCapID returns Gateway's withdraw cap object id
func (gw *Gateway) WithdrawCapID() string {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	return gw.withdrawCapID
}

// WithdrawCapType returns struct type of the WithdrawCap
// Note: the withdraw cap was defined in the original package, so original package ID should be used
func (gw *Gateway) WithdrawCapType() string {
	return fmt.Sprintf("%s::%s::WithdrawCap", gw.Original().PackageID(), GatewayModule)
}

// UpdateIDs updates packageID, objectID, withdrawCapID, previousPackageID and originalPackageID.
func (gw *Gateway) UpdateIDs(pair string) error {
	packageID, gatewayObjectID, withdrawCapID, previousPackageID, originalPackageID, err := parsePair(pair)
	if err != nil {
		return err
	}
	gw.mu.Lock()
	defer gw.mu.Unlock()

	gw.packageID = packageID
	gw.objectID = gatewayObjectID
	gw.withdrawCapID = withdrawCapID
	gw.previousPackageID = previousPackageID
	gw.originalPackageID = originalPackageID

	return nil
}

// ParseEvent parses Event.
func (gw *Gateway) ParseEvent(event models.SuiEventResponse) (Event, error) {
	// event may carry different package IDs, depending on which gateway was called
	packageIDs := gw.SupportedPackageIDs()

	// basic validation
	switch {
	case event.Id.TxDigest == "":
		return Event{}, errors.Wrap(ErrParseEvent, "empty tx hash")
	case event.Id.EventSeq == "":
		return Event{}, errors.Wrap(ErrParseEvent, "empty event id")
	case !slices.Contains(packageIDs, event.PackageId):
		return Event{}, errors.Wrapf(
			ErrParseEvent,
			"package id mismatch (got %s, want one of %s)",
			event.PackageId,
			strings.Join(packageIDs, ","),
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

// parsePair parses a pair of IDs from `$packageID,$gatewayObjectID[,$withdrawCapID,previousPackageID,originalPackageID]`
// There are two cases:
//   - `$packageID,$gatewayObjectID`, the first version of the gateway address
//   - `$packageID,$gatewayObjectID,$withdrawCapID,$previousPackageID,$originalPackageID`, gateway address after upgrade
func parsePair(gatewayAddress string) (string, string, string, string, string, error) {
	parts := strings.Split(gatewayAddress, ",")
	if len(parts) != gatewayPairIDPartsV1 && len(parts) != gatewayPairIDPartsV2 {
		return "", "", "", "", "", errors.Errorf("invalid pair %q", gatewayAddress)
	}

	// each part should be a valid Sui address
	for i, part := range parts {
		// empty previous package ID is valid and it means the previous package is deprecated
		if i == previousPackageIDIndex && part == "" {
			continue
		}

		if err := ValidateAddress(part); err != nil {
			return "", "", "", "", "", errors.Wrapf(err, "invalid Sui address %q", part)
		}
	}

	// for first version of the gateway address
	if len(parts) == gatewayPairIDPartsV1 {
		return parts[0], parts[1], "", "", "", nil
	}

	// after upgrade
	return parts[0], parts[1], parts[2], parts[3], parts[4], nil
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
