package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	suiptb "github.com/pattonkan/sui-go/sui"
	"github.com/pkg/errors"

	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
)

// Client Sui client.
type Client struct {
	sui.ISuiAPI
}

const (
	// DefaultEventsLimit is the default limit for querying gateway module events.
	DefaultEventsLimit = 50

	// TxStatusSuccess is the success status for a transaction.
	TxStatusSuccess = "success"

	// TxStatusFailure is the failure status for a transaction.
	TxStatusFailure = "failure"

	// filterMoveEventModule is the event filter for querying events for specified move module.
	// @see https://docs.sui.io/guides/developer/sui-101/using-events#filtering-event-queries
	filterMoveEventModule = "MoveEventModule"

	// immutableOwner is the owner type for immutable objects.
	immutableOwner = "Immutable"

	// sharedOwner is the owner type for shared objects.
	sharedOwner = "Shared"
)

// NewFromEndpoint Client constructor based on endpoint string.
func NewFromEndpoint(endpoint string) *Client {
	return New(sui.NewSuiClient(endpoint))
}

// New Client constructor.
func New(client sui.ISuiAPI) *Client {
	return &Client{ISuiAPI: client}
}

// HealthCheck queries latest checkpoint and returns its timestamp.
func (c *Client) HealthCheck(ctx context.Context) (time.Time, error) {
	checkpoint, err := c.GetLatestCheckpoint(ctx)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "unable to get latest checkpoint")
	}

	ts, err := strconv.ParseInt(checkpoint.TimestampMs, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to parse checkpoint timestamp")
	}

	return time.UnixMilli(ts).UTC(), nil
}

// GetLatestCheckpoint returns the latest checkpoint.
// See https://docs.sui.io/concepts/cryptography/system/checkpoint-verification
func (c *Client) GetLatestCheckpoint(ctx context.Context) (models.CheckpointResponse, error) {
	seqNum, err := c.SuiGetLatestCheckpointSequenceNumber(ctx)
	if err != nil {
		return models.CheckpointResponse{}, errors.Wrap(err, "unable to get latest seq num")
	}

	return c.SuiGetCheckpoint(ctx, models.SuiGetCheckpointRequest{
		CheckpointID: fmt.Sprintf("%d", seqNum),
	})
}

// EventQuery represents pagination options
type EventQuery struct {
	PackageID string
	Module    string
	Cursor    string
	Limit     uint64
}

// QueryModuleEvents queries module events. Return events and the next pagination cursor.
// If cursor is empty, then the end of scroll reached.
func (c *Client) QueryModuleEvents(ctx context.Context, q EventQuery) ([]models.SuiEventResponse, string, error) {
	if q.Limit == 0 {
		q.Limit = DefaultEventsLimit
	}

	if err := q.validate(); err != nil {
		return nil, "", errors.Wrap(err, "invalid request")
	}

	req, err := q.asRequest()
	if err != nil {
		return nil, "", errors.Wrap(err, "unable to create request")
	}

	res, err := c.SuiXQueryEvents(ctx, req)
	switch {
	case err != nil:
		return nil, "", errors.Wrap(err, "unable to query events")
	case !res.HasNextPage:
		return res.Data, "", nil
	default:
		return res.Data, EncodeCursor(res.NextCursor), nil
	}
}

func (p *EventQuery) validate() error {
	switch {
	case p.PackageID == "":
		return errors.New("package id is empty")
	case p.Module == "":
		return errors.New("module is empty")
	case p.Limit == 0:
		return errors.New("limit is empty")
	case p.Limit > 1000:
		return errors.New("limit exceeded")
	default:
		return nil
	}
}

func (p *EventQuery) asRequest() (models.SuiXQueryEventsRequest, error) {
	filter := map[string]any{
		filterMoveEventModule: map[string]any{
			"package": p.PackageID,
			"module":  p.Module,
		},
	}

	cursor, err := DecodeCursor(p.Cursor)
	if err != nil {
		return models.SuiXQueryEventsRequest{}, err
	}

	return models.SuiXQueryEventsRequest{
		SuiEventFilter:  filter,
		Cursor:          cursor,
		Limit:           p.Limit,
		DescendingOrder: false,
	}, nil
}

// GetOwnedObjectID returns the first owned object ID by owner address and struct type.
// If no objects found or multiple objects found, returns error.
func (c *Client) GetOwnedObjectID(ctx context.Context, ownerAddress, structType string) (string, error) {
	res, err := c.SuiXGetOwnedObjects(ctx, models.SuiXGetOwnedObjectsRequest{
		Address: ownerAddress,
		Query: models.SuiObjectResponseQuery{
			Filter: map[string]any{
				"StructType": structType,
			},
		},
		Limit: 1,
	})

	switch {
	case err != nil:
		return "", errors.Wrap(err, "unable to get owned objects")
	case len(res.Data) == 0:
		return "", errors.New("no objects found")
	case len(res.Data) > 1:
		return "", errors.New("multiple objects found")
	}

	return res.Data[0].Data.ObjectId, nil
}

// InspectTransactionBlock manual implementation of ISuiAPI.InspectTransactionBlock
// Don't use this function at this moment because Sui SDK currently returns deserialization error.
// TODO: https://github.com/zeta-chain/node/issues/3775
//
// @see sui.(*suiReadTransactionImpl).InspectTransactionBlock
// @see https://docs.sui.io/sui-api-ref#sui_devinspecttransactionblock
func (c *Client) InspectTransactionBlock(
	ctx context.Context,
	req models.SuiDevInspectTransactionBlockRequest,
) (models.SuiTransactionBlockResponse, error) {
	const method = "sui_devInspectTransactionBlock"

	params := []any{
		req.Sender,
		req.TxBytes,
		any(nil), // gas_price
		any(nil), // epoch
		any(nil), // additional_args
	}

	resRaw, err := c.SuiCall(ctx, method, params...)
	if err != nil {
		return models.SuiTransactionBlockResponse{}, errors.Wrap(err, method)
	}

	resString, ok := resRaw.(string)
	if !ok {
		return models.SuiTransactionBlockResponse{}, errors.New("invalid response type")
	}

	return parseRPCResponse[models.SuiTransactionBlockResponse]([]byte(resString))
}

// SuiExecuteTransactionBlock manual implementation of ISuiAPI.SuiExecuteTransactionBlock
// That uses proper parameters signature (original has a bug in marshaling)
//
// @see sui.(*suiWriteTransactionImpl).SuiExecuteTransactionBlock
// @see https://docs.sui.io/sui-api-ref#sui_executetransactionblock
func (c *Client) SuiExecuteTransactionBlock(
	ctx context.Context,
	req models.SuiExecuteTransactionBlockRequest,
) (models.SuiTransactionBlockResponse, error) {
	const method = "sui_executeTransactionBlock"

	responseOptionsNullable := any(nil)
	if req.Options != (models.SuiTransactionBlockOptions{}) {
		responseOptionsNullable = req.Options
	}

	requestTypeNullable := any(nil)
	if req.RequestType != "" {
		requestTypeNullable = req.RequestType
	}

	params := []any{
		req.TxBytes,
		req.Signature,
		responseOptionsNullable,
		requestTypeNullable,
	}

	resRaw, err := c.SuiCall(ctx, method, params...)
	if err != nil {
		return models.SuiTransactionBlockResponse{}, errors.Wrap(err, method)
	}

	resString, ok := resRaw.(string)
	if !ok {
		return models.SuiTransactionBlockResponse{}, errors.New("invalid response type")
	}

	return parseRPCResponse[models.SuiTransactionBlockResponse]([]byte(resString))
}

// GetSuiCoinObjectRefs returns a subset of SUI coin objects for given owner address and minimum balance
func (c *Client) GetSuiCoinObjectRefs(
	ctx context.Context,
	owner string,
	minBalanceMist uint64,
) ([]*suiptb.ObjectRef, error) {
	var (
		totalBalance uint64
		cursor       = any(nil)
		suiCoins     = make([]models.CoinData, 0)
	)

	// select SUI coins to cover the given minimum balance
	// there can be multiple pages of SUI coins, we query 50 coins per page
	for {
		resp, err := c.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
			Owner:    owner,
			CoinType: string(zetasui.SUI),
			Cursor:   cursor,
		})
		if err != nil {
			return nil, errors.Wrap(err, "unable to get TSS coins")
		}

		// sort coins by object ID to make the result deterministic across observres
		sort.SliceStable(resp.Data, func(i, j int) bool {
			return strings.Compare(resp.Data[i].CoinObjectId, resp.Data[j].CoinObjectId) < 0
		})

		// append coins until covering the minimum balance
		for _, coin := range resp.Data {
			balance, err := strconv.ParseUint(coin.Balance, 10, 64)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid balance %s", coin.Balance)
			}

			suiCoins = append(suiCoins, coin)
			totalBalance += balance
			if totalBalance >= minBalanceMist {
				break
			}
		}

		if totalBalance >= minBalanceMist || !resp.HasNextPage {
			break
		}

		cursor = resp.NextCursor
	}

	// this is a rare case that can be resolved by sending funds to owner (TSS)
	if totalBalance < minBalanceMist {
		return nil, fmt.Errorf("SUI balance is too low: %d, min balance: %d", totalBalance, minBalanceMist)
	}

	// convert coin data to object references
	suiCoinRefs := make([]*suiptb.ObjectRef, len(suiCoins))
	for i, coin := range suiCoins {
		suiCoinRef, err := coinToObjectRef(coin)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to convert coin to object reference")
		}

		suiCoinRefs[i] = suiCoinRef
	}

	return suiCoinRefs, nil
}

// GetObjectParsedData queries the parsed data of an object.
func (c *Client) GetObjectParsedData(ctx context.Context, objectID string) (models.SuiParsedData, error) {
	resp, err := c.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: objectID,
		Options:  models.SuiObjectDataOptions{ShowContent: true},
	})

	switch {
	case err != nil:
		return models.SuiParsedData{}, errors.Wrap(err, "unable to get gateway object")
	case resp.Error != nil:
		return models.SuiParsedData{}, fmt.Errorf("gateway object response error: %s", resp.Error.Error)
	case resp.Data == nil:
		return models.SuiParsedData{}, errors.New("gateway object data is nil")
	case resp.Data.Content == nil:
		return models.SuiParsedData{}, errors.New("gateway object content is nil")
	default:
		return *resp.Data.Content, nil
	}
}

// EncodeCursor encodes event ID into cursor.
func EncodeCursor(id models.EventId) string {
	return fmt.Sprintf("%s,%s", id.TxDigest, id.EventSeq)
}

// DecodeCursor decodes cursor into event ID.
func DecodeCursor(cursor string) (*models.EventId, error) {
	if cursor == "" {
		return nil, nil
	}

	parts := strings.Split(cursor, ",")
	if len(parts) != 2 {
		return nil, errors.New("invalid cursor format")
	}

	return &models.EventId{
		TxDigest: parts[0],
		EventSeq: parts[1],
	}, nil
}

// parseRPCResponse RPC response into a given type.
func parseRPCResponse[T any](raw []byte) (T, error) {
	// {
	//   "jsonrpc": "2.0",
	//   "id": 1,
	//   "result": { ...}
	// }
	type response struct {
		Result json.RawMessage `json:"result"`
	}

	var (
		outer response
		tt    T
	)

	if err := json.Unmarshal(raw, &outer); err != nil {
		return tt, errors.Wrap(err, "unable to parse rpc response")
	}

	if err := json.Unmarshal(outer.Result, &tt); err != nil {
		return tt, errors.Wrapf(err, "unable to parse result into %T", tt)
	}

	return tt, nil
}

// CheckObjectIDsShared checks if the provided object ID list represents Sui shared or immmutable objects
func (c *Client) CheckObjectIDsShared(ctx context.Context, objectIDs []string) error {
	if len(objectIDs) == 0 {
		return nil
	}

	res, err := c.SuiMultiGetObjects(ctx, models.SuiMultiGetObjectsRequest{
		ObjectIds: objectIDs,
		Options: models.SuiObjectDataOptions{
			ShowOwner: true,
		},
	})
	if err != nil {
		return errors.Wrap(err, "unable to get objects")
	}

	// should always be the case, we add this check as a extra safety measure to ensure an object is not skipped
	if len(res) != len(objectIDs) {
		return fmt.Errorf("expected %d objects, but got %d", len(objectIDs), len(res))
	}

	return CheckContainOwnedObject(res)
}

// CheckContainOwnedObject checks if the provided object list represents Sui shared or immmutable objects
func CheckContainOwnedObject(res []*models.SuiObjectResponse) error {
	for i, obj := range res {
		if obj.Data == nil {
			return errors.Wrapf(zetasui.ErrObjectOwnership, "object %d is missing data", i)
		}

		switch owner := obj.Data.Owner.(type) {
		case string:
			if owner != immutableOwner {
				return errors.Wrapf(zetasui.ErrObjectOwnership, "object %d has unexpected string owner: %s", i, owner)
			}
			// Immutable is valid, continue
		case map[string]any:
			if _, isShared := owner[sharedOwner]; !isShared {
				return errors.Wrapf(zetasui.ErrObjectOwnership, "object %d is not shared or immutable: owner = %+v", i, owner)
			}
			// Shared is valid, continue
		default:
			return errors.Wrapf(zetasui.ErrObjectOwnership, "object %d has unknown owner type: %+v", i, obj.Data.Owner)
		}
	}

	return nil
}

// coinToObjectRef converts a SUI coin data to an object reference.
func coinToObjectRef(coin models.CoinData) (*suiptb.ObjectRef, error) {
	objectID, err := suiptb.ObjectIdFromHex(coin.CoinObjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid coin object ID: %s", coin.CoinObjectId)
	}

	digest, err := suiptb.NewBase58(coin.Digest)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid coin digest: %s", coin.Digest)
	}

	version, err := strconv.ParseUint(coin.Version, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid coin version: %s", coin.Version)
	}

	return &suiptb.ObjectRef{
		ObjectId: objectID,
		Version:  version,
		Digest:   digest,
	}, nil
}
