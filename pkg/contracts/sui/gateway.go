package sui

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
)

// CoinType represents the coin type for the inbound
type CoinType string

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

// SUI is the coin type for SUI, native gas token
const SUI CoinType = "0000000000000000000000000000000000000000000000000000000000000002::sui::SUI"

// Event types
const (
	DepositEvent        EventType = "DepositEvent"
	DepositAndCallEvent EventType = "DepositAndCallEvent"
	WithdrawEvent       EventType = "WithdrawEvent"
	NonceIncreaseEvent  EventType = "NonceIncreaseEvent"
)

// Error codes
// https://github.com/zeta-chain/protocol-contracts-sui/blob/e5a756e473da884dcbc59b574b387a7a365ac823/sources/gateway.move#L14-L21
const (
	ErrCodeAlreadyWhitelisted     uint64 = 0
	ErrCodeInvalidReceiverAddress uint64 = 1
	ErrCodeNotWhitelisted         uint64 = 2
	ErrCodeNonceMismatch          uint64 = 3
	ErrCodePayloadTooLong         uint64 = 4
	ErrCodeInactiveWithdrawCap    uint64 = 5
	ErrCodeInactiveWhitelistCap   uint64 = 6
	ErrCodeDepositPaused          uint64 = 7
)

const moduleName = "gateway"

var (
	// ErrParseEvent event parse error
	ErrParseEvent = errors.New("event parse error")

	// retryableOutboundErrCodes are the outbound execution (if failed) error codes that are retryable.
	// The list is used to determine if a withdraw_and_call should fallback if rejected by the network.
	// Note: keep this list in sync with the actual implementation in `gateway.move`
	retryableOutboundErrCodes = []uint64{
		ErrCodeNotWhitelisted,
		ErrCodeNonceMismatch,
		ErrCodeInactiveWithdrawCap,
	}
)

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

// Withdrawal extract Withdraw data.
func (e *Event) Withdrawal() (Withdrawal, error) {
	v, ok := e.content.(Withdrawal)
	if !ok {
		return Withdrawal{}, errors.Errorf("invalid content type %T", e.content)
	}

	return v, nil
}

func (e *Event) IsNonceIncrease() bool {
	return e.EventType == NonceIncreaseEvent
}

// NonceIncrease extract NonceIncrease data.
func (e *Event) NonceIncrease() (NonceIncrease, error) {
	v, ok := e.content.(NonceIncrease)
	if !ok {
		return NonceIncrease{}, errors.Errorf("invalid content type %T", e.content)
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

// Module returns Gateway's module name
func (gw *Gateway) Module() string {
	return moduleName
}

// WithdrawCapType returns struct type of the WithdrawCap
func (gw *Gateway) WithdrawCapType() string {
	return fmt.Sprintf("%s::%s::WithdrawCap", gw.PackageID(), moduleName)
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
	if descriptor.module != moduleName {
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
	case NonceIncreaseEvent:
		content, err = parseNonceIncrease(event, eventType)
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
	if len(res.Events) == 0 {
		return event, nil, errors.New("missing events")
	}

	event, err = gw.ParseEvent(res.Events[0])
	if err != nil {
		return event, nil, err
	}

	// filter outbound events
	switch {
	case event.IsWithdraw():
		withdrawal, err := event.Withdrawal()
		if err != nil {
			return event, nil, errors.Wrap(err, "unable to extract withdraw event")
		}
		return event, withdrawal, nil
	case event.IsNonceIncrease():
		nonceIncrease, err := event.NonceIncrease()
		if err != nil {
			return event, nil, errors.Wrap(err, "unable to extract nonce increase event")
		}
		return event, nonceIncrease, nil
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

// IsRetryableMoveAbort checks if the error message is a retryable 'MoveAbort' error.
func IsRetryableMoveAbort(errorMsg string) (bool, error) {
	if !strings.HasPrefix(errorMsg, "MoveAbort") {
		return false, nil // not MoveAbort error
	}

	code, err := parseExecutionErrorMoveAbortCode(errorMsg)
	if err != nil {
		return false, errors.Wrap(err, "unable to extract move abort code")
	}

	if slices.Contains(retryableOutboundErrCodes, code) {
		return true, nil
	}

	return false, nil
}

// parseExecutionErrorMoveAbortCode parses the error code from Sui 'ExecutionError::MoveAbort' execution error.
// see: https://github.com/MystenLabs/sui-rust-sdk/blob/65eb9f3ad63b98f5b04465963d340e53b301a149/crates/sui-sdk-types/src/execution_status.rs#L173
//
// Example error message:
// "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier("gateway") }, function: 11, instruction: 37, function_name: Some("withdraw_impl") }, 3) in command 0"
func parseExecutionErrorMoveAbortCode(errorMsg string) (uint64, error) {
	// build regex to match pattern: MoveAbort(..., <code>) ...
	re := regexp.MustCompile(`MoveAbort\(.+?,\s*(\d+)\)`)
	matches := re.FindStringSubmatch(errorMsg)
	if len(matches) != 2 {
		return 0, errors.Errorf("unable to extract code from error string: %s", errorMsg)
	}

	codeStr := matches[1]
	code, err := strconv.ParseUint(codeStr, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to convert code to uint64: %s", codeStr)
	}
	return code, nil
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

	// Sui encode bytes in base64
	decodedPayload, err := base64.StdEncoding.DecodeString(string(payload))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode payload from base64")
	}

	return decodedPayload, nil
}

func parsePair(pair string) (string, string, error) {
	parts := strings.Split(pair, ",")
	if len(parts) != 2 {
		return "", "", errors.Errorf("invalid pair %q", pair)
	}

	return parts[0], parts[1], nil
}
