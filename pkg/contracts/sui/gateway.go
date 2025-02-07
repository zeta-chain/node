package sui

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
)

const DepositEvent = "DepositEvent"
const DepositAndCallEvent = "DepositAndCallEvent"
const GatewayModule = "gateway"

// Gateway contains the API to read inbounds and sign outbounds to the Sui gateway
type Gateway struct {
	client    sui.ISuiAPI
	packageID string
}

// NewGateway creates a new Sui gateway
// Note: packageID is the equivalent for gateway address or program ID on Solana
// It's what will be stared in gateway chain params
func NewGateway(client sui.ISuiAPI, packageID string) Gateway {
	return Gateway{
		client:    client,
		packageID: packageID,
	}
}

// QueryDepositInbounds queries the inbounds from deposit events from the Sui gateway
func (g Gateway) QueryDepositInbounds(ctx context.Context, from, to uint64) ([]Inbound, error) {
	return g.queryInbounds(ctx, from, to, false)
}

// QueryDepositAndCallInbounds queries the inbounds from depositAndCall events from the Sui gateway
func (g Gateway) QueryDepositAndCallInbounds(ctx context.Context, from, to uint64) ([]Inbound, error) {
	return g.queryInbounds(ctx, from, to, true)
}

func (g Gateway) queryInbounds(ctx context.Context, _, _ uint64, depositAndCall bool) ([]Inbound, error) {
	// event to filter
	event := DepositEvent
	if depositAndCall {
		event = DepositAndCallEvent
	}

	// make the query
	// TODO: Support pagination
	res, err := g.client.SuiXQueryEvents(ctx, models.SuiXQueryEventsRequest{
		SuiEventFilter: map[string]interface{}{
			// TODO: Fix the error
			// using TimeRange causes the following error when sending the query:
			// {"code":-32602,"message":"Invalid params","data":"expected value at line 1 column 108"}
			// commenting out for new and querying all events
			//"TimeRange": models.TimeRange{
			//	StartTime: from,
			//	EndTime:   to,
			//},
			"MoveEventType": fmt.Sprintf("%s::%s::%s", g.packageID, GatewayModule, event),
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
		inbound, err := parseInbound(eventData, depositAndCall)
		if err != nil {
			return []Inbound{}, err
		}
		inbounds = append(inbounds, inbound)
	}

	return inbounds, nil
}
