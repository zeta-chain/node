package sui

import (
	"strconv"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
)

// CoinType represents the coin type for the inbound
type CoinType string

// EventType represents Gateway event type (both inbound & outbound)
type EventType string

// Gateway contains the API to read inbounds and sign outbounds to the Sui gateway
type Gateway struct {
	packageID string
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

// NewGateway creates a new Sui gateway
// Note: packageID is the equivalent for gateway address or program ID on Solana
// It's what will be set in gateway chain params
func NewGateway(packageID string) *Gateway {
	return &Gateway{packageID}
}

// Event represents generic event wrapper
type Event struct {
	TxHash     string
	EventIndex uint64
	EventType  EventType

	content any
	inbound bool
}

// IsInbound checks whether event is Inbound.
func (e *Event) IsInbound() bool { return e.inbound }

// Inbound extract Inbound.
func (e *Event) Inbound() (Inbound, error) {
	if !e.inbound {
		return Inbound{}, errors.New("not an inbound")
	}

	return e.content.(Inbound), nil
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
	case event.TransactionModule != moduleName:
		return Event{}, errors.Wrapf(
			ErrParseEvent,
			"module mismatch (got %s, want %s)",
			event.TransactionModule,
			moduleName,
		)
	}

	// Extract common fields
	txHash := event.Id.TxDigest
	eventID, err := strconv.ParseUint(event.Id.EventSeq, 10, 64)
	if err != nil {
		return Event{}, errors.Wrapf(ErrParseEvent, "failed to parse event id %q", event.Id.EventSeq)
	}

	eventType, err := extractEventType(event.Type)
	if err != nil {
		return Event{}, errors.Wrap(ErrParseEvent, err.Error())
	}

	var (
		inbound bool
		content any
	)

	// Parse specific events
	switch eventType {
	case Deposit, DepositAndCall:
		inbound = true
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

		content: content,
		inbound: inbound,
	}, nil
}

func extractEventType(typeString string) (EventType, error) {
	// e.g. 0x3e9fb7c....d6cc443cf::gateway::DepositEvent
	parts := strings.Split(typeString, "::")
	if len(parts) != 3 {
		return "", errors.Errorf("invalid event type %s", typeString)
	}

	return EventType(parts[2]), nil
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
