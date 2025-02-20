package sui

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

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

// SUI is the coin type for SUI, native gas token
const SUI CoinType = "0000000000000000000000000000000000000000000000000000000000000002::sui::SUI"

// Event types
const (
	Deposit        EventType = "DepositEvent"
	DepositAndCall EventType = "DepositAndCallEvent"
)

const moduleName = "gateway"

// ErrParseEvent event parse error
var ErrParseEvent = errors.New("event parse error")

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

// Inbound extract Inbound.
func (e *Event) Inbound() (Inbound, error) {
	v, ok := e.content.(Inbound)
	if !ok {
		return Inbound{}, errors.Errorf("invalid content type %T", e.content)
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
	case Deposit, DepositAndCall:
		content, err = parseInbound(event, eventType)
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

	return payload, nil
}

func parsePair(pair string) (string, string, error) {
	parts := strings.Split(pair, ",")
	if len(parts) != 2 {
		return "", "", errors.Errorf("invalid pair %q", pair)
	}

	return parts[0], parts[1], nil
}
