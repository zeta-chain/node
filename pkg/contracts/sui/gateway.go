package sui

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/pkg/errors"
)

// SUI is the coin type for SUI, native gas token
const SUI CoinType = "0000000000000000000000000000000000000000000000000000000000000002::sui::SUI"

const (
	eventDeposit        = "DepositEvent"
	eventDepositAndCall = "DepositAndCallEvent"
	moduleName          = "gateway"
)

// CoinType represents the coin type for the inbound
type CoinType string

// Gateway contains the API to read inbounds and sign outbounds to the Sui gateway
type Gateway struct {
	client    sui.ISuiAPI
	packageID string
}

//go:embed bin/gateway.mv
var gatewayBinary []byte

// ErrParseEvent event parse error
var ErrParseEvent = errors.New("event parse error")

// NewGateway creates a new SUI gateway
// Note: packageID is the equivalent for gateway address or program ID on Solana
// It's what will be set in gateway chain params
func NewGateway(client sui.ISuiAPI, packageID string) *Gateway {
	return &Gateway{client, packageID}
}

// GatewayBytecodeBase64 gets the gateway binary encoded as base64 for deployment
func GatewayBytecodeBase64() string {
	return base64.StdEncoding.EncodeToString(gatewayBinary)
}

// QueryDepositInbounds queries the inbounds from deposit events from the Sui gateway
// from and to represents time range in Unix time in milliseconds
func (g *Gateway) QueryDepositInbounds(ctx context.Context, from, to uint64) ([]Inbound, error) {
	return g.queryInbounds(ctx, from, to, eventDeposit)
}

// QueryDepositAndCallInbounds queries the inbounds from depositAndCall events from the Sui gateway
// from and to represents time range in Unix time in milliseconds
func (g *Gateway) QueryDepositAndCallInbounds(ctx context.Context, from, to uint64) ([]Inbound, error) {
	return g.queryInbounds(ctx, from, to, eventDepositAndCall)
}

func (g *Gateway) queryInbounds(ctx context.Context, _, _ uint64, event string) ([]Inbound, error) {
	// make the query
	// TODO: Support pagination
	res, err := g.client.SuiXQueryEvents(ctx, models.SuiXQueryEventsRequest{
		SuiEventFilter: map[string]any{
			// TODO: Fix the error
			// using TimeRange causes the following error when sending the query:
			// {"code":-32602,"message":"Invalid params","data":"expected value at line 1 column 108"}
			// commenting out for new and querying all events
			//"TimeRange": models.TimeRange{
			//	StartTime: from,
			//	EndTime:   to,
			//},
			//"TimeRange": map[string]interface{}{
			//"startTime": from,
			//"endTime":   to,
			//},
			"MoveEventType": eventType(g.packageID, moduleName, event),
		},
		Limit: 5,
	})
	if err != nil {
		return []Inbound{}, err
	}

	// parse the events
	inbounds := make([]Inbound, 0, len(res.Data))
	for _, eventData := range res.Data {
		// TODO: events that fail to be parsed should still be observed as invalid (if at least the tx hash and event seq are present)
		// Example: if the receiver is not a valid ETH address, the observation should fail and be reverted to the sender
		// https://github.com/zeta-chain/node/issues/3502
		inbound, err := parseInbound(eventData, event)
		if err != nil {
			return []Inbound{}, errors.Wrap(ErrParseEvent, err.Error())
		}

		inbounds = append(inbounds, inbound)
	}

	return inbounds, nil
}

func eventType(packageID, module, event string) string {
	return fmt.Sprintf("%s::%s::%s", packageID, module, event)
}
